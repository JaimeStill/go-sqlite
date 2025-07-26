package config

import (
	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	Verbose bool   `mapstructure:"verbose"`
	Format  string `mapstructure:"format"`
}

// App is the global configuration instance
var App *Config

// Init initializes the global configuration from viper
func (c *Config) Init() error {
	App = &Config{
		Verbose: viper.GetBool("verbose"),
		Format:  viper.GetString("format"),
	}
	
	// Apply defaults
	if App.Format == "" {
		App.Format = "text"
	}
	
	return nil
}

// GetFormat returns the configured output format
func (c *Config) GetFormat() string {
	return c.Format
}

// IsVerbose returns whether verbose mode is enabled
func (c *Config) IsVerbose() bool {
	return c.Verbose
}