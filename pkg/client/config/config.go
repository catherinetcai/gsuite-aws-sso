package config

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/catherinetcai/gsuite-aws-sso/pkg/file"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/logging"
	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v2"
)

const (
	DefaultServerURL = "http://localhost:3030/credentials"
)

// Config wraps all client configs
type Config struct {
	Server string `yaml:"server" json:"server"`
	GCP    GCP    `yaml:"gcp" json:"gcp"`
	AWS    AWS    `yaml:"aws" json:"aws"`
}

// GCP wraps all of the Google Cloud configs
type GCP struct {
	// The filepath of the credentials seeded by GCloud auth
	CredentialFilePath string `yaml:"credential_file_path" json:"credential_file_path"`
}

// AWS wraps all of the AWS configs
type AWS struct {
	// Path where the AWS credential path goes - typically ~/.aws/credentials
	CredentialOutputPath string `yaml:"credential_output_path" json:"credential_output_path"`
}

// DefaultGCloudCredentialPath ...
func DefaultGCloudCredentialPath() string {
	path, err := file.WithUserHomeDir(".config", "gcloud", "application_default_credentials.json")
	if err != nil {
		logging.Logger().Error("error getting home path", zap.Error(err))
		// Just defaults to returning a relative position
		return ""
	}

	return path
}

// DefaultAWSOutputPath ...
func DefaultAWSOutputPath() string {
	path, err := file.WithUserHomeDir(".aws", "credentials")
	if err != nil {
		logging.Logger().Error("error getting home path", zap.Error(err))
		// Just defaults to returning a relative position
		return ""
	}

	return path
}

func configPath() string {
	path, err := file.WithUserHomeDir(".gsuite_aws_sso/config")
	if err != nil {
		logging.Logger().Error("error getting home path", zap.Error(err))
		// Just defaults to returning a relative position
		return ""
	}

	return path
}

// Default returns a default Config instance
func Default() *Config {
	return &Config{
		Server: DefaultServerURL,
		GCP: GCP{
			CredentialFilePath: DefaultGCloudCredentialPath(),
		},
		AWS: AWS{
			CredentialOutputPath: DefaultAWSOutputPath(),
		},
	}
}

// Get Configs...
func Get() (*Config, error) {
	cfgFile, err := os.Open(configPath())
	if err != nil {
		return nil, err
	}

	cfgRaw, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = yaml.Unmarshal(cfgRaw, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// SetConfigs takes in a config file and attempts to write it
func SetConfigs(cfg *Config) error {
	cfgRaw, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	dirPath := path.Dir(configPath())
	if err = os.MkdirAll(dirPath, 0744); err != nil {
		return err
	}

	cfgFile, err := os.OpenFile(configPath(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	// Ensure file is empty
	if err := cfgFile.Truncate(0); err != nil {
		return err
	}
	cfgFile.Seek(0, 0)

	_, err = cfgFile.Write(cfgRaw)
	return err
}
