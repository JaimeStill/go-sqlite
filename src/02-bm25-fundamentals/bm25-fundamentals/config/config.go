package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// App represents the global configuration instance with defaults
var App = NewConfig()

// Config represents the application configuration schema
type Config struct {
	// Global settings
	Database string `mapstructure:"database"`
	Verbose  bool   `mapstructure:"verbose"`
	Format   string `mapstructure:"format"`

	// Corpus configuration
	Corpus CorpusConfig `mapstructure:"corpus"`

	// Search configuration
	Search SearchConfig `mapstructure:"search"`

	// Display configuration
	Display DisplayConfig `mapstructure:"display"`

	// Visualization configuration
	Visualization VisualizationConfig `mapstructure:"visualization"`

	// Analysis configuration
	Analysis AnalysisConfig `mapstructure:"analysis"`
}

// CorpusConfig holds corpus generation settings
type CorpusConfig struct {
	Size      int `mapstructure:"size"`
	BatchSize int `mapstructure:"batch_size"`
}

// SearchConfig holds search-related settings
type SearchConfig struct {
	MaxResults    int `mapstructure:"max_results"`
	TermFreqLimit int `mapstructure:"term_freq_limit"`
}

// DisplayConfig holds display formatting settings
type DisplayConfig struct {
	ScorePrecision int `mapstructure:"score_precision"`
}

// VisualizationConfig holds visualization settings
type VisualizationConfig struct {
	HistogramWidth  int  `mapstructure:"histogram_width"`
	HistogramHeight int  `mapstructure:"histogram_height"`
	ShowLegend      bool `mapstructure:"show_legend"`
}

// AnalysisConfig holds analysis settings
type AnalysisConfig struct {
	MinScoreBuckets int   `mapstructure:"min_score_buckets"`
	Percentiles     []int `mapstructure:"percentiles"`
}

// NewConfig creates a new config with defaults
func NewConfig() *Config {
	return &Config{
		Database: ":memory:",
		Verbose:  false,
		Format:   "text",
		Corpus: CorpusConfig{
			Size:      100,
			BatchSize: 1000,
		},
		Search: SearchConfig{
			MaxResults:    20,
			TermFreqLimit: 10,
		},
		Display: DisplayConfig{
			ScorePrecision: 4,
		},
		Visualization: VisualizationConfig{
			HistogramWidth:  50,
			HistogramHeight: 10,
			ShowLegend:      true,
		},
		Analysis: AnalysisConfig{
			MinScoreBuckets: 10,
			Percentiles:     []int{25, 50, 75, 90, 95, 99},
		},
	}
}

// SetDefaults applies default values to viper
func (c *Config) SetDefaults() {
	viper.SetDefault("database", c.Database)
	viper.SetDefault("verbose", c.Verbose)
	viper.SetDefault("format", c.Format)

	viper.SetDefault("corpus.size", c.Corpus.Size)
	viper.SetDefault("corpus.batch_size", c.Corpus.BatchSize)

	viper.SetDefault("search.max_results", c.Search.MaxResults)
	viper.SetDefault("search.term_freq_limit", c.Search.TermFreqLimit)

	viper.SetDefault("display.score_precision", c.Display.ScorePrecision)

	viper.SetDefault("visualization.histogram_width", c.Visualization.HistogramWidth)
	viper.SetDefault("visualization.histogram_height", c.Visualization.HistogramHeight)
	viper.SetDefault("visualization.show_legend", c.Visualization.ShowLegend)

	viper.SetDefault("analysis.min_score_buckets", c.Analysis.MinScoreBuckets)
	viper.SetDefault("analysis.percentiles", c.Analysis.Percentiles)
}

// Validate checks the configuration for errors
func (c *Config) Validate() error {
	// Validate format
	switch c.Format {
	case "text", "json", "csv":
		// Valid formats
	default:
		return fmt.Errorf("invalid format: %s (must be text, json, or csv)", c.Format)
	}

	// Validate corpus settings
	if c.Corpus.Size < 1 {
		return fmt.Errorf("corpus size must be at least 1")
	}
	if c.Corpus.BatchSize < 1 {
		return fmt.Errorf("corpus batch size must be at least 1")
	}

	// Validate search settings
	if c.Search.MaxResults < 1 {
		return fmt.Errorf("max results must be at least 1")
	}

	// Validate visualization settings
	if c.Visualization.HistogramWidth < 10 {
		return fmt.Errorf("histogram width must be at least 10")
	}
	if c.Visualization.HistogramHeight < 3 {
		return fmt.Errorf("histogram height must be at least 3")
	}

	return nil
}

// Load reads configuration from viper
func (c *Config) Load() error {
	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return c.Validate()
}

// RefreshFromFlags updates config from current flag values
func (c *Config) RefreshFromFlags() error {
	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return nil
}

// GetDatabasePath returns the database path, expanding ~ to home directory
func (c *Config) GetDatabasePath() string {
	dbPath := c.Database
	if dbPath == ":memory:" {
		return dbPath
	}

	// Expand ~ to home directory
	if len(dbPath) > 0 && dbPath[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			return dbPath
		}
		dbPath = filepath.Join(home, dbPath[1:])
	}

	return dbPath
}

// Init initializes the global configuration
func (c *Config) Init() {
	// Apply defaults to viper
	c.SetDefaults()

	if cfgFile := viper.GetString("config"); cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bm25-fundamentals" (without extension)
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".bm25-fundamentals")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

	// Load configuration into struct
	if err := c.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}
}
