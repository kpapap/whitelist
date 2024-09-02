package whitelist

import (
	"context"
	"os"
	"os/exec"
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

	// Get the collector's environment variables
	env := os.Environ()

	// Execute script
	r.settings.Logger.Info("Running script...")
	script := "script.sh"
	cmd := exec.Command("/bin/sh", script)
	cmd.Env = env
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	screrr := cmd.Run()
	if screrr != nil {
		r.settings.Logger.Sugar().Errorf("Error running script: %s", screrr.Error())
	}
	r.settings.Logger.Info("Script executed")
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