package main

import (
	"context"
	"fmt"
	"log"

	"terraform-provider-kubiya/internal/provider"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"

	sentrygo "github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

const (
	version = "dev"
	address = "hashicorp.com/edu/Kubiya"
)

func main() {
	// Initialize Sentry
	if err := kubiyasentry.Initialize(); err != nil {
		// Log the error but continue - Sentry should not prevent the provider from running
		log.Printf("Warning: Failed to initialize Sentry: %v", err)
	}
	defer kubiyasentry.Flush()

	// Set up panic recovery with Sentry
	defer func() {
		if r := recover(); r != nil {
			// Capture panic to Sentry
			kubiyasentry.CaptureMessage("Provider panic", sentrygo.LevelFatal, map[string]string{
				"panic": fmt.Sprintf("%v", r),
			})
			kubiyasentry.Flush()
			panic(r) // Re-panic after capturing
		}
	}()

	ctx := context.Background()
	kubiya := provider.New(version)

	opts := providerserver.ServeOpts{
		Address: address,
	}

	err := providerserver.Serve(ctx, kubiya, opts)
	if err != nil {
		// Capture fatal error to Sentry
		kubiyasentry.CaptureError(err, ctx, map[string]string{
			"error_type": "provider_serve_error",
		})
		kubiyasentry.Flush()
		log.Fatal(err.Error())
	}
}
