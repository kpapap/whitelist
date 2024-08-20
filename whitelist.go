package whitelist

import (
	"context"
	"net"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver"
)


type whitelistReceiver struct {
	ctx						context.Context
	cfg      			*Config
	nextConsumer	consumer.Logs
	settings 			receiver.Settings
	shutdownWG  	sync.WaitGroup
}

// newWhitelistReceiver just creates the OpenTelemetry receiver services. It is the caller's
// responsibility to invoke the respective Start*Reception methods as well
// as the various Stop*Reception methods to end it.
func newWhitelistReceiver(ctx context.Context, cfg *Config, nextConsumer consumer.Logs, settings receiver.Settings) (*whitelistReceiver, error) {
	r := &whitelistReceiver{
		ctx:        	ctx,
		cfg:        	cfg,
		nextConsumer:	nextConsumer,
		settings:			settings,
	}
	// Proceed to the next receiver.
	logs := plog.NewLogs()
	err := r.nextConsumer.ConsumeLogs(ctx, logs)
	if err != nil {
		// Handle the error
		r.settings.Logger.Sugar().Errorf("error consuming logs: %v", err.Error())
	}
	return r, nil
}

// Start the receiver
func (r *whitelistReceiver) Start(ctx context.Context, host component.Host) error {
	// Create an http ticket for http checks
	r.settings.Logger.Info("Creating HTTP ticker")
	httprepeatTimeStr := "1m"
	if httprepeatTimeStr == "" {
		r.settings.Logger.Error("HTTP ticker is not set")
	}
	httprepeatTime, err := time.ParseDuration(httprepeatTimeStr)
	if err != nil {
		r.settings.Logger.Sugar().Errorf("Error parsing http ticker environment variable: %s", err.Error())
	}
	httpticker := time.NewTicker(httprepeatTime)
	defer httpticker.Stop()
	r.settings.Logger.Info("HTTP Ticker created")

	// Check connection
			r.settings.Logger.Info("Checking http connection...")
			conn, err := net.DialTimeout("tcp", "www.google.com:80", 3*time.Second)
			if err != nil {
				r.settings.Logger.Info("port closed")
				return err
			}
			defer conn.Close()
			r.settings.Logger.Info("port open")
			return nil
		}


// Shutdown the receiver
func (r *whitelistReceiver) Shutdown(ctx context.Context) error {
	var err error
	r.shutdownWG.Wait()
	// Log a message indicating that the receiver is shutting down.
	r.settings.Logger.Info("Shutting down receiver")
	// Return err to indicate that the receiver shut down successfully.
	return err
}