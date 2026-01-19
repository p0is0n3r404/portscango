package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config application configuration
type Config struct {
	Default  DefaultConfig            `yaml:"default"`
	Profiles map[string]ProfileConfig `yaml:"profiles"`
	Webhook  WebhookConfig            `yaml:"webhook"`
	Database DatabaseConfig           `yaml:"database"`
}

// DefaultConfig default settings
type DefaultConfig struct {
	Threads          int           `yaml:"threads"`
	Timeout          time.Duration `yaml:"timeout"`
	ServiceDetection bool          `yaml:"service_detection"`
	BannerGrab       bool          `yaml:"banner_grab"`
	OutputFormat     string        `yaml:"output_format"`
	ColorOutput      bool          `yaml:"color_output"`
}

// ProfileConfig scan profile settings
type ProfileConfig struct {
	Threads     int           `yaml:"threads"`
	Timeout     time.Duration `yaml:"timeout"`
	Delay       time.Duration `yaml:"delay"`
	RandomOrder bool          `yaml:"random_order"`
	Ports       []int         `yaml:"ports"`
}

// WebhookConfig webhook settings
type WebhookConfig struct {
	Discord string `yaml:"discord"`
	Slack   string `yaml:"slack"`
	Custom  string `yaml:"custom"`
}

// DatabaseConfig database settings
type DatabaseConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

// DefaultConfiguration returns default config
func DefaultConfiguration() *Config {
	return &Config{
		Default: DefaultConfig{
			Threads:          100,
			Timeout:          1 * time.Second,
			ServiceDetection: true,
			BannerGrab:       true,
			OutputFormat:     "table",
			ColorOutput:      true,
		},
		Profiles: map[string]ProfileConfig{
			"stealth": {
				Threads:     10,
				Timeout:     2 * time.Second,
				Delay:       500 * time.Millisecond,
				RandomOrder: true,
			},
			"quick": {
				Threads: 200,
				Timeout: 500 * time.Millisecond,
			},
			"aggressive": {
				Threads: 500,
				Timeout: 200 * time.Millisecond,
			},
		},
		Database: DatabaseConfig{
			Enabled: false,
			Path:    "portscango.db",
		},
	}
}

// LoadConfig loads config from file
func LoadConfig(path string) (*Config, error) {
	// If path not specified, try default locations
	if path == "" {
		path = findConfigFile()
	}

	if path == "" {
		return DefaultConfiguration(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfiguration(), nil
	}

	cfg := DefaultConfiguration()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// SaveConfig saves config to file
func SaveConfig(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// findConfigFile looks for config in default locations
func findConfigFile() string {
	// Check current directory
	if _, err := os.Stat(".portscango.yaml"); err == nil {
		return ".portscango.yaml"
	}

	// Check home directory
	home, err := os.UserHomeDir()
	if err == nil {
		homePath := filepath.Join(home, ".portscango.yaml")
		if _, err := os.Stat(homePath); err == nil {
			return homePath
		}
	}

	return ""
}

// GenerateDefaultConfig creates a default config file
func GenerateDefaultConfig(path string) error {
	cfg := DefaultConfiguration()
	return SaveConfig(path, cfg)
}
