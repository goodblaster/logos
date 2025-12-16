package main

import (
	"context"
	"os"

	"github.com/goodblaster/logos"
)

// This demo shows context-based logging for request-scoped loggers.
func main() {
	// Create a logger with request-specific fields
	requestLogger := logos.NewLogger(logos.LevelInfo, logos.ConsoleFormatter(), os.Stdout).
		With("request_id", "req-12345").
		With("user_id", "user-789").
		With("ip", "192.168.1.100")

	// Store logger in context
	ctx := logos.WithLogger(context.Background(), requestLogger)

	// Pass context through your application
	handleRequest(ctx)
	processOrder(ctx)
	sendNotification(ctx)

	// If no logger in context, FromContext returns the default logger
	emptyCtx := context.Background()
	logos.FromContext(emptyCtx).Info("No logger in context, using default")
}

func handleRequest(ctx context.Context) {
	log := logos.FromContext(ctx)
	log.Info("Handling request")
	log.With("endpoint", "/api/orders").Info("Request validated")
}

func processOrder(ctx context.Context) {
	log := logos.FromContext(ctx)
	log.With("order_id", "order-555").Info("Processing order")
	log.With("amount", 99.99).Info("Payment processed")
}

func sendNotification(ctx context.Context) {
	log := logos.FromContext(ctx)
	log.With("notification_type", "email").Info("Sending notification")
}
