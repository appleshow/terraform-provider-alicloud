package alicloud

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudDummyParameters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDummyParametersRead,

		Schema: map[string]*schema.Schema{
			"map_parameters": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"list_parameters": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"list_map_parameters": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func dataSourceAlicloudDummyParametersRead(d *schema.ResourceData, meta interface{}) error {
	id := time.Now().Format("2006-01-02 15:04:05")
	d.SetId(id)

	return nil
}
