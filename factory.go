package whitelist

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

var (typeStr = component.MustNewType("whitelist"))
const (defaultInterval = "1m")

// createLogsReceiver creates a new instance of the nbcmr receiver.
func createLogsReceiver(ctx context.Context, set receiver.Settings, cfg component.Config, nextConsumer consumer.Logs) (receiver.Logs, error) {
	// Create the new receiver
	logs := &whitelistReceiver{
			config:       cfg.(*Config),
			nextConsumer: nextConsumer,
			logger:       set.Logger,
	}
	return logs, nil
}

// NewFactory creates a new receiver factory for the nbcmr receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithLogs(createLogsReceiver, component.StabilityLevelAlpha))
}

// createDefaultConfig returns the default configuration for the nbcmr receiver.
// This function is used when creating a new factory to provide a default configuration
// for the receiver.
func createDefaultConfig() component.Config {
	return &Config{
		Interval: defaultInterval,
	}
}