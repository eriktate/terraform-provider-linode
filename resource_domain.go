package main

import (
	"fmt"
	"strconv"

	"github.com/eriktate/lingo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		Update: resourceDomainUpdate,
		Delete: resourceDomainDelete,

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"soa_email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, m interface{}) error {
	req := lingo.Domain{
		Domain: d.Get("domain").(string),
		Type:   lingo.DomainType(d.Get("type").(string)),
		SOA:    d.Get("soa_email").(string),
	}

	domain, err := linode.CreateDomain(req)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", domain.ID))
	return nil
}

func resourceDomainRead(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.ParseUint(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	domain, err := linode.ViewDomain(uint(id))
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("domain", domain.Domain)
	d.Set("type", domain.Type)
	d.Set("soa_email", domain.SOA)

	return nil
}

func resourceDomainUpdate(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.ParseUint(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	req := lingo.Domain{
		ID:     uint(id),
		Domain: d.Get("domain").(string),
		Type:   lingo.DomainType(d.Get("type").(string)),
		SOA:    d.Get("soa_email").(string),
	}

	if _, err := linode.UpdateDomain(req); err != nil {
		return err
	}

	return nil
}

func resourceDomainDelete(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.ParseUint(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	if err := linode.DeleteDomain(uint(id)); err != nil {
		return err
	}

	return nil
}
