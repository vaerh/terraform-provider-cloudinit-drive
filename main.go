package main

import (
	"context"
	"flag"
	"log"
	"terraform-provider-cloudinit-drive/cid"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(
		context.Background(),
		cid.NewCloudInitDrive,
		providerserver.ServeOpts{
			Address: "github.com/vaerh/cloudinit-drive",
			Debug:   debug,
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}
