package whitelist

import (
	"fmt"
	"time"
)

// Config represents the receiver config settings within the collector's config.yaml
type Config struct {
	// The name of the configMap.
	ConfigMapName string `mapstructure:"config_map_name"`
	// The interval of the collector.
	Interval string `mapstructure:"config_interval"`
}

// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}

	interval, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		return fmt.Errorf("failed to parse interval: %w", err)
	}

	if interval.Minutes() < 1 {
		return fmt.Errorf("interval has to be set to at least 1 minute (1m)")
	}

	return nil
}