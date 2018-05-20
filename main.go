package main

import (
	"os"

	"github.com/eriktate/lingo"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

var linode lingo.Lingo

func main() {
	apiKey := os.Getenv("LINODE_API_KEY")
	linode = lingo.NewLingo(apiKey)

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})
}
