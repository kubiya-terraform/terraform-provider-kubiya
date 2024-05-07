package main

import (
	"context"
	"log"

	"terraform-provider-kubiya/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

const (
	version = "dev"
	address = "hashicorp.com/edu/Kubiya"
)

func main() {
	ctx := context.Background()
	kubiya := provider.New(version)

	opts := providerserver.ServeOpts{
		Address: address,
	}

	err := providerserver.Serve(ctx, kubiya, opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
