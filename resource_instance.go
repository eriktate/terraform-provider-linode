package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/eriktate/lingo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,

		Schema: map[string]*schema.Schema{
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"stackscript_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"booted": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"root_pass": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"image": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"authorized_keys": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"backup_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"backups_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"swap_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, m interface{}) error {
	req := lingo.CreateLinodeRequest{
		Label:          GetString(d, "label"),
		Region:         GetString(d, "region"),
		Type:           GetString(d, "type"),
		StackScriptID:  GetUint(d, "stackscript_id"),
		Booted:         GetBool(d, "booted"),
		RootPass:       GetString(d, "root_pass"),
		Image:          GetString(d, "image"),
		AuthorizedKeys: GetSliceString(d, "authorized_keys"),
		BackupID:       GetUint(d, "backup_id"),
		BackupsEnabled: GetBool(d, "backups_enabled"),
		SwapSize:       GetUint(d, "swap_size"),
	}

	linode, err := linode.CreateLinode(req)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", linode.ID))
	return nil
}

func resourceInstanceRead(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.ParseUint(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	linode, err := linode.ViewLinode(uint(id))
	if err != nil {
		d.SetId("")
		return nil
	}

	disks, err := linode.ListDisks(uint(id))
	if err != nil {
		return nil
	}

	var swapSize uint
	for _, disk := range disks {
		if disk.FileSystem == lingo.FileSystemSwap {
			swapSize = disk.Size
			break
		}
	}

	d.Set("label", linode.Label)
	d.Set("image", linode.Image)
	d.Set("region", linode.Region)
	d.Set("region", linode.Region)
	d.Set("type", linode.Type)
	d.Set("swap_size", swapSize)

	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	id, err := strconv.ParseUint(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	if d.HasChange("swap_size") {
		disks, err := linode.ListDisks(uint(id))
		if err != nil {
			return err
		}

		var swap lingo.Disk
		for _, disk := range disks {
			log.Printf("Searching disk: %+v", disk)
			if disk.FileSystem == lingo.FileSystemSwap {
				swap = disk
				break
			}
		}

		if _, err := linode.ResizeDisk(uint(id), swap.ID, GetUint(d, "swap_size")); err != nil {
			return err
		}

		d.SetPartial("swap_size")
	}

	if d.HasChange("label") || d.HasChange("alerts") {
		// TODO: Add alerts later.
		req := lingo.UpdateLinodeRequest{
			ID:    uint(id),
			Label: GetString(d, "label"),
		}

		if _, err := linode.UpdateLinode(req); err != nil {
			return err
		}

		d.SetPartial("label")
		d.SetPartial("alerts")
	}

	d.Partial(false)

	return nil
}

func resourceInstanceDelete(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.ParseUint(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	if err := linode.DeleteLinode(uint(id)); err != nil {
		return err
	}

	return nil
}
