package main

import "github.com/hashicorp/terraform/helper/schema"

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"linode_domain":   resourceDomain(),
			"linode_instance": resourceInstance(),
			"linode_disk":     resourceDisk(),
		},
	}
}
