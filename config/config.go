package config

import "github.com/nanoteck137/validate"

type ConfigOidcProvider struct {
	Id           string `mapstructure:"name"`
	Name         string `mapstructure:"name"`
	ClientId     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	IssuerUrl    string `mapstructure:"issuer_url"`
	RedirectUrl  string `mapstructure:"redirect_url"`
}

func (c ConfigOidcProvider) Validate() error {
	return validate.ValidateStruct(&c,
		validate.Field(&c.Name, validate.Required),
		validate.Field(&c.ClientId, validate.Required),
		validate.Field(&c.ClientSecret, validate.Required),
		validate.Field(&c.IssuerUrl, validate.Required),
		validate.Field(&c.RedirectUrl, validate.Required),
	)
}

type MediaQualityConfig struct {
	High   int `mapstructure:"high"`
	Medium int `mapstructure:"medium"`
	Low    int `mapstructure:"low"`
}

type MediaConfig struct {
	Opus   MediaQualityConfig `mapstructure:"opus"`
	Vorbis MediaQualityConfig `mapstructure:"vorbis"`
	Mp3    MediaQualityConfig `mapstructure:"mp3"`
	Aac    MediaQualityConfig `mapstructure:"aac"`

	AudioNormalization bool `mapstructure:"audio_normalization"`
}

type Config struct {
	RunMigrations bool `mapstructure:"run_migrations"`

	ListenAddr string `mapstructure:"listen_addr"`

	DataDir    string `mapstructure:"data_dir"`
	LibraryDir string `mapstructure:"library_dir"`
	WebDir     string `mapstructure:"web"`

	JwtSecret string `mapstructure:"jwt_secret"`

	MeilisearchAddress string `mapstructure:"meilisearch_address"`
	MeilisearchApiKey  string `mapstructure:"meilisearch_api_key"`

	NtfyBaseUrl string `mapstructure:"ntfy_base_url"`
	NtfyTopic   string `mapstructure:"ntfy_topic"`

	OidcProviders []ConfigOidcProvider `mapstructure:"oidc_providers"`

	Media MediaConfig `mapstructure:"media"`
}

func (c Config) Validate() error {
	return validate.ValidateStruct(&c,
		validate.Field(&c.ListenAddr, validate.Required),
		validate.Field(&c.DataDir, validate.Required),
		validate.Field(&c.LibraryDir, validate.Required),
		validate.Field(&c.JwtSecret, validate.Required),

		validate.Field(&c.MeilisearchAddress, validate.Required),
		validate.Field(&c.MeilisearchApiKey, validate.Required),

		validate.Field(&c.OidcProviders,
			validate.Required, validate.Length(1, 0)),
	)
}

func configDefaults() map[string]any {
	return map[string]any{
		"run_migrations": true,
		"listen_addr":    ":3000",
		"media": map[string]any{
			"audio_normalization": true,
			"opus": map[string]any{
				"high":   128,
				"medium": 96,
				"low":    64,
			},
			"vorbis": map[string]any{
				"high":   192,
				"medium": 128,
				"low":    96,
			},
			"mp3": map[string]any{
				"high":   320,
				"medium": 192,
				"low":    128,
			},
			"aac": map[string]any{
				"high":   256,
				"medium": 192,
				"low":    96,
			},
		},
	}
}
