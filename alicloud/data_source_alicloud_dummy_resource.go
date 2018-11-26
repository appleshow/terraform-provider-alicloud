package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudDummyResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDummyResourceRead,

		Schema: map[string]*schema.Schema{
			"parameters": &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Required: true,
			},
			"alicloud_dummy_resource": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Required: false,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudDummyResourceRead(d *schema.ResourceData, meta interface{}) error {
	var s []map[string]interface{}

	id := time.Now().Format("2006-01-02 15:04:05")
	parameters := d.Get("parameters").(map[string]interface{})

	parameters["id"] = id
	parameters["name"] = id
	parameters["status"] = "Available"
	parameters["creation_time"] = id
	parameters["resource_type"] = "alicloud_dummy_resource"

	d.SetId(id)
	s = append(s, parameters)
	if err := d.Set("alicloud_dummy_resource", s); err != nil {
		return fmt.Errorf("Setting alicloud_dummy_resource got an error: %#v.", err)
	}

	return nil
}
