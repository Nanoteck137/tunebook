package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/nanoteck137/tunebook"
	"github.com/nanoteck137/validate"
	"github.com/pelletier/go-toml/v2"
)

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
}

func (c Config) Validate() error {
	return validate.ValidateStruct(&c,
		validate.Field(&c.ListenAddr, validate.Required),
		validate.Field(&c.DataDir, validate.Required),
		validate.Field(&c.LibraryDir, validate.Required),
		validate.Field(&c.JwtSecret, validate.Required),

		validate.Field(&c.MeilisearchAddress, validate.Required),
		validate.Field(&c.MeilisearchApiKey, validate.Required),

		validate.Field(&c.OidcProviders, validate.Required, validate.Length(1, 0)),
	)
}

func configDefaults() map[string]any {
	return map[string]any{
		"run_migrations": true,
		"listen_addr":    ":3000",
	}
}

func readFileToMap(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	err = toml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func mergeMaps(base, overlay map[string]any) {
	for k, v := range overlay {
		ov, ok := v.(map[string]any)
		if !ok {
			base[k] = v
			continue
		}

		bv, ok := base[k]
		if !ok {
			base[k] = v
			continue
		}

		bvm, ok := bv.(map[string]any)
		if !ok {
			base[k] = v
			continue
		}

		mergeMaps(bvm, ov)
	}
}

func readConfigFromEnv() map[string]any {
	prefix := strings.ToUpper(tunebook.AppName + "_")
	m := make(map[string]any)

	for _, e := range os.Environ() {
		k, v, ok := strings.Cut(e, "=")
		if !ok {
			continue
		}

		if !strings.HasPrefix(k, prefix) {
			continue
		}

		key := strings.ToLower(strings.TrimPrefix(k, prefix))
		m[key] = v
	}

	return m
}

func extractConfigOverride(m map[string]any) string {
	if raw, ok := m["config_override"]; ok {
		s, _ := raw.(string)
		return s
	}
	return ""
}

func Load(cfgFile string) (*Config, error) {
	configMap := configDefaults()

	if cfgFile != "" {
		m, err := readFileToMap(cfgFile)
		if err != nil {
			return nil, fmt.Errorf("reading config: %w", err)
		}

		mergeMaps(configMap, m)
	} else {
		m, err := readFileToMap("config.toml")
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("reading config: %w", err)
			}
		} else {
			mergeMaps(configMap, m)
		}
	}

	key := strings.ToUpper(tunebook.AppName + "_CONFIG_OVERRIDE")
	configOverride := os.Getenv(key)
	if configOverride == "" {
		configOverride = extractConfigOverride(configMap)
	}

	if configOverride != "" {
		m, err := readFileToMap(configOverride)
		if err != nil {
			return nil, fmt.Errorf("reading override config: %w", err)
		}

		mergeMaps(configMap, m)
	}

	envMap := readConfigFromEnv()
	mergeMaps(configMap, envMap)

	var config Config

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           &config,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return nil, fmt.Errorf("new decoder: %w", err)
	}

	err = decoder.Decode(configMap)
	if err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	oldTag := validate.ErrorTag
	validate.ErrorTag = "mapstructure"
	defer func() {
		validate.ErrorTag = oldTag
	}()

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}
