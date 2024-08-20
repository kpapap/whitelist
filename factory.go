package whitelist

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (defaultInterval = "1m")
var (typeStr = component.MustNewType("whitelist"))

// NewFactory creates a new receiver factory.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithLogs(createLogsReceiver, component.StabilityLevelUndefined),
	)
}

// createDefaultConfig returns the default configuration for the  receiver.
// This function is used when creating a new factory to provide a default configuration
// for the receiver.
func createDefaultConfig() component.Config {
	return &Config{
		Interval: defaultInterval,
	}
}

// createLogsReceiver creates a new instance of the logs receiver.
// createLogsReceiver creates a log receiver based on provided config.
func createLogsReceiver(ctx context.Context, settings receiver.Settings, cfg component.Config, consumer consumer.Logs) (receiver.Logs, error) {
	// Create the new receiver
	rCfg := cfg.(*Config)
	return newWhitelistReceiver(rCfg, consumer, settings)
}