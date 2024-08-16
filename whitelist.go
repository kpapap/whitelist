package whitelist

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)


type whitelistReceiver struct {
	config *Config
	logger *zap.Logger
	nextConsumer consumer.Logs
}

func (c *whitelistReceiver) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (c *whitelistReceiver) ConsumeLogs(ctx context.Context, ld consumer.Logs) error {
	return nil
}

func (r *whitelistReceiver) Start(ctx context.Context, host component.Host) (context.Context, error) {
	// Start the receiver
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
	log.Println("Checking http connection...")
	conn, err := net.DialTimeout("tcp", "www.google.com:80", 3*time.Second)
		if err != nil {
			fmt.Println("port closed")
		}
	defer conn.Close()
	fmt.Println("port open")
	return ctx, nil
}

// Shutdown shuts down the receiver.
func (r *whitelistReceiver) Shutdown(ctx context.Context) error {
	// Shutdown the receiver
	// Log a message indicating that the receiver is shutting down.
	log.Println("Shutting down receiver")
	// Return nil to indicate that the receiver shut down successfully.
	return nil
}

