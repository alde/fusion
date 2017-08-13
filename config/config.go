package config

import (
	"io/ioutil"
	"os"

	"github.com/kelseyhightower/envconfig"
	yaml "gopkg.in/yaml.v2"
)

// DatabaseConfig object
type DatabaseConfig struct {
	Host     string `envconfig:"database_host"`
	Port     int16  `envconfig:"database_port"`
	Name     string `envconfig:"database_name"`
	User     string `envconnfig:"database_user"`
	Password string `envconnfig:"database_password"`
}

// The Config struct holds the Fusion Configuration
type Config struct {
	// Logging
	LogLevel  string `yaml:"loglevel"`
	LogFormat string `yaml:"logformat"`

	// Server
	Address string `yaml:"address"`
	Port    int16  `yaml:"port"`

	// Database
	Database DatabaseConfig `yaml:"database"`
}

// DefaultConfig returns a Config struct with the default settings
func DefaultConfig() *Config {
	return &Config{
		LogLevel:  "debug",
		LogFormat: "text",

		Database: DatabaseConfig{
			Host:     "127.0.0.1",
			Port:     3306,
			Name:     "fusion",
			User:     "fusion",
			Password: "password",
		},
	}
}

// GetConfigFilePath returns the location of the config file in order of priority:
// 1 ) Specified by --config command-line flag
// 2 ) Global file in /etc/fusion/fusion.yml
func GetConfigFilePath(override string) string {
	if len(override) > 0 {
		return override
	}
	globalPath := "/etc/fusion/fusion.yml"
	if _, err := os.Open(globalPath); err == nil {
		return globalPath
	}

	return ""
}

// ReadConfigFile reads the config file and overrides any values net in both it
// and the DefaultConfig
func ReadConfigFile(conf *Config, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	configFile, _ := ioutil.ReadAll(file)
	return yaml.Unmarshal(configFile, conf)
}

// ReadEnvironment takes environment variables and overrides any values from
// DefaultConfig and the Config file.
func ReadEnvironment(conf *Config) {
	envconfig.Process("", conf)
}
