package main

import (
	"fmt"
	"strconv"

	"github.com/eriktate/lingo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiskCreate,
		Read:   resourceDiskRead,
		Update: resourceDiskUpdate,
		Delete: resourceDiskDelete,

		Schema: map[string]*schema.Schema{
			"linode_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"filesystem": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateFileSystem,
			},
			"read_only": &schema.Schema{
				Type:     schema.TypeBool,
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
			"root_pass": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"stackscript_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceDiskCreate(d *schema.ResourceData, m interface{}) error {
	req := lingo.CreateDiskRequest{
		LinodeID:       GetUint(d, "linode_id"),
		Label:          GetString(d, "label"),
		Size:           GetUint(d, "size"),
		FileSystem:     lingo.FileSystem(GetString(d, "filesystem")),
		ReadOnly:       GetBool(d, "read_only"),
		Image:          GetString(d, "image"),
		AuthorizedKeys: GetSliceString(d, "authorized_keys"),
		RootPass:       GetString(d, "root_pass"),
		StackscriptID:  GetUint(d, "stackscript_id"),
	}

	disk, err := linode.CreateDisk(req)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", disk.ID))
	return nil
}

func resourceDiskRead(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.ParseUint(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	linodeID := GetUint(d, "linode_id")

	disk, err := linode.ViewDisk(linodeID, uint(id))
	if err != nil {
		return err
	}

	d.Set("label", disk.Label)
	d.Set("size", disk.Size)
	d.Set("filesystem", disk.FileSystem)

	return nil
}

func resourceDiskUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	id, err := strconv.ParseUint(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	linodeID := GetUint(d, "linode_id")

	if d.HasChange("size") {
		if _, err := linode.ResizeDisk(linodeID, uint(id), GetUint(d, "size")); err != nil {
			return err
		}

		d.SetPartial("size")
	}

	if d.HasChange("label") || d.HasChange("filesystem") {
		req := lingo.UpdateDiskRequest{
			ID:         uint(id),
			LinodeID:   linodeID,
			Label:      GetString(d, "label"),
			FileSystem: lingo.FileSystem(GetString(d, "filesystem")),
		}

		if _, err := linode.UpdateDisk(req); err != nil {
			return err
		}

		d.SetPartial("label")
		d.SetPartial("filesystem")
	}

	d.Partial(false)
	return nil
}

func resourceDiskDelete(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.ParseUint(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	linodeID := GetUint(d, "linode_id")

	if err := linode.DeleteDisk(linodeID, uint(id)); err != nil {
		return err
	}

	return nil
}

func validateFileSystem(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be a string", k)}
	}

	if !lingo.ValidateFileSystem(v) {
		return nil, []error{fmt.Errorf("'%s' value of '%s' is not a valid FileSystem", v, k)}
	}

	return nil, nil
}
