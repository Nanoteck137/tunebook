package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoteck137/tunebook/config"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nrednav/cuid2"
	"golang.org/x/oauth2"
)

var createAuthId, _ = cuid2.Init(cuid2.WithLength(32))

var authErr = NewServiceErrCreator("auth")

var (
	ErrAuthServiceProviderNotFound = authErr.New("provider not found")

	ErrAuthServiceRequestAlreadyExists = authErr.New("request already exists")
	ErrAuthServiceRequestNotFound      = authErr.New("request not found")
	ErrAuthServiceRequestExpired       = authErr.New("request expired")
	ErrAuthServiceRequestNotReady      = authErr.New("request not ready")
	ErrAuthServiceRequestInvalid       = authErr.New("request invalid")
)

const (
	authProviderRequestExpireDuration   = 5 * time.Minute
	authProviderRequestDeletionDuration = authProviderRequestExpireDuration + 10*time.Minute

	authQuickRequestExpireDuration   = 5 * time.Minute
	authQuickRequestDeletionDuration = authQuickRequestExpireDuration + 10*time.Minute
)

type AuthProviderRequestStatus string

const (
	AuthProviderRequestStatusPending   AuthProviderRequestStatus = "pending"
	AuthProviderRequestStatusCompleted AuthProviderRequestStatus = "completed"
	AuthProviderRequestStatusExpired   AuthProviderRequestStatus = "expired"
	AuthProviderRequestStatusFailed    AuthProviderRequestStatus = "failed"
)

type AuthQuickRequestStatus string

const (
	AuthQuickRequestStatusPending   AuthQuickRequestStatus = "pending"
	AuthQuickRequestStatusCompleted AuthQuickRequestStatus = "completed"
	AuthQuickRequestStatusExpired   AuthQuickRequestStatus = "expired"
	AuthQuickRequestStatusFailed    AuthQuickRequestStatus = "failed"
)

type authProviderRequest struct {
	id         string
	providerId string
	status     AuthProviderRequestStatus
	challenge  string
	oauth2Url  string
	oauth2Code string
	expires    time.Time
	delete     time.Time
}

type authProvider struct {
	initialized  bool
	id           string
	displayName  string
	config       config.ConfigOidcProvider
	provider     *oidc.Provider
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func (p *authProvider) init(ctx context.Context) error {
	if p.initialized {
		return nil
	}

	var err error

	p.provider, err = oidc.NewProvider(ctx, p.config.IssuerUrl)
	if err != nil {
		return err
	}

	p.oauth2Config = &oauth2.Config{
		ClientID:     p.config.ClientId,
		ClientSecret: p.config.ClientSecret,
		RedirectURL:  p.config.RedirectUrl,
		Endpoint:     p.provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	p.verifier = p.provider.Verifier(&oidc.Config{ClientID: p.config.ClientId})

	p.initialized = true

	return nil
}

type providerClaim struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
	Sub         string `json:"sub"`
}

func (p *authProvider) claim(ctx context.Context, code string) (providerClaim, error) {
	oauth2Token, err := p.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return providerClaim{}, fmt.Errorf("exchange: %w", err)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return providerClaim{}, errors.New("oauth2 token is missing id_token")
	}

	idToken, err := p.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return providerClaim{}, fmt.Errorf("verify: %w", err)
	}

	var claims providerClaim
	err = idToken.Claims(&claims)
	if err != nil {
		return providerClaim{}, fmt.Errorf("claims: %w", err)
	}

	return claims, nil
}

type authQuickConnectRequest struct {
	status    AuthQuickRequestStatus
	code      string
	challenge string
	userId    string
	expires   time.Time
	delete    time.Time
}

type AuthService struct {
	logger               *slog.Logger
	mu                   sync.Mutex
	db                   *database.Database
	imageService         *ImageService
	jwtSecret            string
	providers            map[string]*authProvider
	providerRequests     map[string]*authProviderRequest
	quickConnectRequests map[string]*authQuickConnectRequest
}

func NewAuthService(
	logger *slog.Logger,
	db *database.Database,
	config *config.Config,
	imageService *ImageService,
) *AuthService {
	providers := make(map[string]*authProvider, len(config.OidcProviders))

	for _, providerConfig := range config.OidcProviders {
		id := providerConfig.Id

		res := &authProvider{
			id:          id,
			displayName: providerConfig.Name,
			config:      providerConfig,
		}

		providers[id] = res
	}

	return &AuthService{
		logger:               logger,
		db:                   db,
		imageService:         imageService,
		jwtSecret:            config.JwtSecret,
		providers:            providers,
		providerRequests:     make(map[string]*authProviderRequest),
		quickConnectRequests: make(map[string]*authQuickConnectRequest),
	}
}

type ProviderRequestResult struct {
	RequestId string
	AuthUrl   string
	Challenge string
	Expires   time.Time
}

func (a *AuthService) CreateProviderRequest(providerId string) (ProviderRequestResult, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	provider, exists := a.providers[providerId]
	if !exists {
		return ProviderRequestResult{}, ErrAuthServiceProviderNotFound
	}

	err := provider.init(context.TODO())
	if err != nil {
		return ProviderRequestResult{}, authErr.Newf("initialize AuthProvider(%s): %w", provider.id, err)
	}

	challenge, err := generateAuthChallenge()
	if err != nil {
		return ProviderRequestResult{}, authErr.Newf("generate auth challenge: %w", err)
	}

	id := createAuthId()

	t := time.Now()
	request := &authProviderRequest{
		id:         id,
		providerId: provider.id,
		status:     AuthProviderRequestStatusPending,
		challenge:  challenge,
		expires:    t.Add(authProviderRequestExpireDuration),
		delete:     t.Add(authProviderRequestDeletionDuration),
	}

	request.oauth2Url = provider.oauth2Config.AuthCodeURL(request.id)

	_, exists = a.providerRequests[id]
	if exists {
		return ProviderRequestResult{}, ErrAuthServiceRequestAlreadyExists
	}

	a.providerRequests[id] = request

	return ProviderRequestResult{
		RequestId: request.id,
		AuthUrl:   request.oauth2Url,
		Challenge: request.challenge,
		Expires:   request.expires,
	}, nil
}

type QuickConnectRequestResult struct {
	Code      string
	Challenge string
	Expires   time.Time
}

func (a *AuthService) CreateQuickConnectRequest() (QuickConnectRequestResult, error) {
	code, err := generateCode()
	if err != nil {
		return QuickConnectRequestResult{}, fmt.Errorf("failed to generate code: %w", err)
	}

	challenge, err := generateAuthChallenge()
	if err != nil {
		return QuickConnectRequestResult{}, fmt.Errorf("failed to generate challenge: %w", err)
	}

	t := time.Now()
	request := &authQuickConnectRequest{
		status:    AuthQuickRequestStatusPending,
		code:      code,
		challenge: challenge,
		expires:   t.Add(authQuickRequestExpireDuration),
		delete:    t.Add(authQuickRequestDeletionDuration),
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	_, exists := a.quickConnectRequests[code]
	if exists {
		return QuickConnectRequestResult{}, err
	}

	a.quickConnectRequests[code] = request

	return QuickConnectRequestResult{
		Code:      request.code,
		Challenge: request.challenge,
		Expires:   request.expires,
	}, nil
}

func (a *AuthService) CompleteQuickConnectRequest(requestCode, userId string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	request, exists := a.quickConnectRequests[requestCode]
	if !exists {
		return ErrAuthServiceRequestNotFound
	}

	if time.Now().After(request.expires) {
		request.status = AuthQuickRequestStatusExpired
		return ErrAuthServiceRequestExpired
	}

	if request.status == AuthQuickRequestStatusPending {
		request.status = AuthQuickRequestStatusCompleted
		request.userId = userId
	}

	return nil
}

func (a *AuthService) CompleteProviderRequest(requestId, code string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	request, exists := a.providerRequests[requestId]
	if !exists {
		return ErrAuthServiceRequestNotFound
	}

	if time.Now().After(request.expires) {
		request.status = AuthProviderRequestStatusExpired
		return ErrAuthServiceRequestExpired
	}

	if request.status == AuthProviderRequestStatusPending {
		request.status = AuthProviderRequestStatusCompleted
		request.oauth2Code = code
	}

	return nil
}

func (a *AuthService) CheckProviderRequestStatus(requestId, challenge string) (AuthProviderRequestStatus, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	request, exists := a.providerRequests[requestId]
	if !exists {
		return AuthProviderRequestStatusFailed, ErrAuthServiceRequestNotFound
	}

	if request.challenge != challenge {
		// TODO(patrik): Give this it's own error
		return AuthProviderRequestStatusFailed, ErrAuthServiceRequestNotFound
	}

	now := time.Now()
	if now.After(request.expires) {
		request.status = AuthProviderRequestStatusExpired
	}

	return request.status, nil
}

func (a *AuthService) CheckQuickConnectRequestStatus(requestCode, challenge string) (AuthQuickRequestStatus, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	request, exists := a.quickConnectRequests[requestCode]
	if !exists {
		return AuthQuickRequestStatusFailed, ErrAuthServiceRequestNotFound
	}

	if request.challenge != challenge {
		return AuthQuickRequestStatusFailed, ErrAuthServiceRequestNotFound
	}

	now := time.Now()
	if now.After(request.expires) {
		request.status = AuthQuickRequestStatusExpired
	}

	return request.status, nil
}

func (a *AuthService) CreateAuthTokenForProvider(requestId, challenge string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	request, exists := a.providerRequests[requestId]
	if !exists {
		return "", ErrAuthServiceRequestNotFound
	}

	if request.challenge != challenge {
		return "", ErrAuthServiceRequestNotFound
	}

	if request.status != AuthProviderRequestStatusCompleted {
		return "", ErrAuthServiceRequestNotReady
	}

	// TODO(patrik): Check provider?
	provider := a.providers[request.providerId]

	if request.oauth2Code == "" {
		request.status = AuthProviderRequestStatusFailed
		return "", ErrAuthServiceRequestInvalid
	}

	userId, err := a.getUserFromCode(context.TODO(), provider, request.oauth2Code)
	if err != nil {
		request.status = AuthProviderRequestStatusFailed
		return "", err
	}

	if userId == "" {
		request.status = AuthProviderRequestStatusFailed
		return "", ErrAuthServiceRequestInvalid
	}

	request.status = AuthProviderRequestStatusExpired

	token, err := a.SignUserToken(userId)
	if err != nil {
		request.status = AuthProviderRequestStatusFailed
		return "", err
	}

	return token, nil
}

func (a *AuthService) CreateAuthTokenForQuickConnect(requestCode, challenge string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	request, exists := a.quickConnectRequests[requestCode]
	if !exists {
		return "", ErrAuthServiceRequestNotFound
	}

	if request.challenge != challenge {
		return "", ErrAuthServiceRequestNotFound
	}

	if request.status != AuthQuickRequestStatusCompleted {
		return "", ErrAuthServiceRequestInvalid
	}

	if request.userId == "" {
		return "", ErrAuthServiceRequestInvalid
	}

	request.status = AuthQuickRequestStatusExpired

	token, err := a.SignUserToken(request.userId)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) getUserFromCode(ctx context.Context, provider *authProvider, code string) (string, error) {
	oidcClaims, err := provider.claim(ctx, code)
	if err != nil {
		return "", authErr.Newf("provider claim: %w", err)
	}

	getOrCreateUser := func() (string, error) {
		user, err := a.db.GetUserByEmail(ctx, oidcClaims.Email)
		if err == nil {
			return user.Id, nil
		}

		if errors.Is(err, database.ErrItemNotFound) {
			displayName := oidcClaims.DisplayName
			if displayName == "" {
				displayName = oidcClaims.Name
			}

			userCount, err := a.db.CountUsers(ctx)
			if err != nil {
				return "", err
			}

			role := types.RoleUser
			if userCount == 0 {
				role = types.RoleSuperUser
			}

			userId, err := a.db.CreateUser(ctx, database.CreateUserParams{
				Email:       oidcClaims.Email,
				DisplayName: displayName,
				Role:        role,
			})
			if err != nil {
				return "", authErr.Newf("create user: %w", err)
			}

			if oidcClaims.Picture != "" {
				picture, err := a.imageService.DownloadPictureForUser(
					ctx,
					DownloadPictureForUserParams{
						UserId: userId,
						Url:    oidcClaims.Picture,
					},
				)
				if err != nil {
					return "", err
				}

				err = a.db.UpdateUser(ctx, userId, database.UserChanges{
					Picture: database.Change[sql.NullString]{
						Value: sql.NullString{
							String: picture,
							Valid:  picture != "",
						},
						Changed: true,
					},
				})
				if err != nil {
					return "", fmt.Errorf("failed to update user picture: %w", err)
				}
			}

			return userId, nil
		} else {
			return "", authErr.Newf("get user by email: %w", err)
		}
	}

	identity, err := a.db.GetUserIdentity(ctx, provider.id, oidcClaims.Sub)
	if err == nil {
		return identity.UserId, nil
	}

	if errors.Is(err, database.ErrItemNotFound) {
		userId, err := getOrCreateUser()
		if err != nil {
			return "", err
		}

		err = a.db.CreateUserIdentity(ctx, database.CreateUserIdentityParams{
			Provider:   provider.id,
			ProviderId: oidcClaims.Sub,
			UserId:     userId,
		})
		if err != nil {
			return "", authErr.Newf("create user identity: %w", err)
		}

		return userId, nil
	} else {
		return "", authErr.Newf("get user identity: %w", err)
	}
}

func (a *AuthService) SignUserToken(userId string) (string, error) {
	user, err := a.db.GetUserById(context.Background(), userId)
	if err != nil {
		return "", authErr.Newf("signing token: get user by id: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Id,
		"iat":    time.Now().Unix(),
	})

	tokenString, err := token.SignedString(([]byte)(a.jwtSecret))
	if err != nil {
		return "", authErr.Newf("signing token: jwt sign: %w", err)
	}

	return tokenString, nil
}

func (a *AuthService) Cleanup() {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()

	for k, request := range a.providerRequests {
		if now.After(request.delete) {
			delete(a.providerRequests, k)
		}
	}

	for k, request := range a.quickConnectRequests {
		if now.After(request.delete) {
			delete(a.quickConnectRequests, k)
		}
	}
}

const (
	letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
)

func randomString(charset string, length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}

func generateCode() (string, error) {
	part1, err := randomString(letters, 4)
	if err != nil {
		return "", err
	}

	part2, err := randomString(digits, 4)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", part1, part2), nil
}

func generateAuthChallenge() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
