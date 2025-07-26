package config

import (
	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	DatabasePath string `mapstructure:"database"`
	Verbose      bool   `mapstructure:"verbose"`
	Format       string `mapstructure:"format"`
}

// NewConfig creates a new config instance with defaults
func NewConfig() *Config {
	return &Config{
		DatabasePath: ":memory:",
		Verbose:      false,
		Format:       "text",
	}
}

// App is the global configuration instance
var App = NewConfig()

// Init initializes the global configuration from viper
func (c *Config) Init() error {
	c.DatabasePath = viper.GetString("database")
	c.Verbose = viper.GetBool("verbose")
	c.Format = viper.GetString("format")
	
	// Apply defaults if empty
	if c.DatabasePath == "" {
		c.DatabasePath = ":memory:"
	}
	if c.Format == "" {
		c.Format = "text"
	}
	
	return nil
}

// GetDatabasePath returns the configured database path
func (c *Config) GetDatabasePath() string {
	return c.DatabasePath
}

// GetFormat returns the configured output format
func (c *Config) GetFormat() string {
	return c.Format
}

// IsVerbose returns whether verbose mode is enabled
func (c *Config) IsVerbose() bool {
	return c.Verbose
}