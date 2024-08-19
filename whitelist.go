package whitelist

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)


type whitelistReceiver struct {
	cfg      			*Config
	nextConsumer	consumer.Logs
	settings 			receiver.Settings
	shutdownWG  	sync.WaitGroup
}

// newWhitelistReceiver just creates the OpenTelemetry receiver services. It is the caller's
// responsibility to invoke the respective Start*Reception methods as well
// as the various Stop*Reception methods to end it.
func newWhitelistReceiver(cfg *Config, nextConsumer consumer.Logs, settings receiver.Settings) (*whitelistReceiver, error) {
	r := &whitelistReceiver{
		cfg:        	cfg,
		nextConsumer:	nextConsumer,
		settings:			settings,
	}
	return r, nil
}

// Start the receiver
func (r *whitelistReceiver) Start(ctx context.Context, host component.Host) error {
	// Create an http ticket for http checks
	log.Println("Creating HTTP ticker")
	httprepeatTimeStr := "1m"
	if httprepeatTimeStr == "" {
		log.Fatal("HTTP ticker is not set")
	}
	httprepeatTime, err := time.ParseDuration(httprepeatTimeStr)
	if err != nil {
		log.Fatalf("Error parsing http ticker environment variable: %s", err.Error())
	}
	httpticker := time.NewTicker(httprepeatTime)
	defer httpticker.Stop()
	log.Println("HTTP Ticker created")

	// Check connection
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-httpticker.C:
			log.Println("Checking http connection...")
			conn, err := net.DialTimeout("tcp", "www.google.com:80", 3*time.Second)
			if err != nil {
				fmt.Println("port closed")
				return err
			}
			defer conn.Close()
			fmt.Println("port open")
		}
	}
}

// Shutdown the receiver
func (r *whitelistReceiver) Shutdown(ctx context.Context) error {
	var err error
	r.shutdownWG.Wait()
	// Log a message indicating that the receiver is shutting down.
	log.Println("Shutting down receiver")
	// Return err to indicate that the receiver shut down successfully.
	return err
}