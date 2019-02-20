package config

import (
	"strings"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	goconfig "github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
)

// NOTE: I may want to rethink this entire approach.
// I'm just not very happy with the way that Viper does things.
// This is an experiment with go-config.

var (
	gocfg    goconfig.Config
	instance *Config
	once     sync.Once
)

// Config ...
type Config struct {
	GSuite GSuite `json:"gsuite"`
	OAuth  OAuth  `json:"oauth"`
	Server Server `json:"server"`
}

// AWS encapsulates all AWS configs
type AWS struct{}

// GSuite encapsulates all GSuite service info
type GSuite struct {
	// Base64 encoded representation of the service account.
	// Either this or the path need to be set
	ServiceAccountBase64EncodedFile string `json:"service_account_file"`

	// The file path where the service account file can be found.
	// If this is prefixed with s3:// it will go look in S3 for the
	// file. Must be in the format of s3://<bucket>/<path>.
	// Otherwise, this must be the absolute path
	ServiceAccountPath  string `json:"service_account_path"`
	ClientID            string `json:"client_id"`
	ServiceAccountEmail string `json:"service_account_email"`
}

// OAuth encapsulates all OAuth configs
type OAuth struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	// Scopes will come in as a comma delimited string
	Scopes                  []string      `json:"scopes"`
	TokenURL                string        `json:"token_url"`
	AuthURL                 string        `json:"auth_url"`
	RedirectURL             string        `json:"redirect_url"`
	StateParameterGenerator func() string `json:"-"`
}

// Server encapsulates all server configs
type Server struct {
	Port int `json:"port"`
}

// Initialize configs
func Initialize() {
	// With the env source - env variables with _ will be broken up as follows
	// DATABASE_SERVER_HOST = config.Get("database", "server", "host")
	envSrc := env.NewSource()
	gocfg = goconfig.NewConfig(goconfig.WithSource(envSrc))
}

// Get instance of the config
func Get() *Config {
	once.Do(func() {
		instance = &Config{
			GSuite: GSuite{
				ClientID:                        gocfg.Get("gsuite", "client", "id").String(""),
				ServiceAccountBase64EncodedFile: gocfg.Get("gsuite", "service", "account", "base64", "file").String(""),
				ServiceAccountPath:              gocfg.Get("gsuite", "service", "account", "path").String(""),
				ServiceAccountEmail:             gocfg.Get("gsuite", "service", "account", "email").String(""),
			},
			OAuth: OAuth{
				ClientID:     gocfg.Get("oauth", "client", "id").String(""),
				ClientSecret: gocfg.Get("oauth", "client", "secret").String(""),
				Scopes:       strings.Split(gocfg.Get("oauth", "scopes").String(""), ","),
				TokenURL:     gocfg.Get("oauth", "token", "url").String(""),
				AuthURL:      gocfg.Get("oauth", "auth", "url").String(""),
				RedirectURL:  gocfg.Get("oauth", "redirect", "url").String(""),
			},
			Server: Server{
				Port: gocfg.Get("server", "port").Int(3030),
			},
		}
	})

	return instance
}
