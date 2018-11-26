package alicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudLogStores() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudLogStoresRead,

		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_log_stores": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"store_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudLogStoresRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	projectName := d.Get("project_name").(string)

	listLogStoreOutput, err := client.slsconn.ListLogStore(projectName)
	if err != nil {
		return fmt.Errorf("List stores of log service got an error: %#v", err)
	} else {
		log.Printf("[DEBUG] alicloud_log_service - stores found: %#v", listLogStoreOutput)

		var ids []string
		var s []map[string]interface{}
		for _, name := range listLogStoreOutput {
			mapping := map[string]interface{}{
				"id":            name,
				"name":          name,
				"status":        "Available",
				"creation_time": "",
				"project_name":  projectName,
				"store_name":    name,
				"resource_type": "alicloud_log_store",
			}
			ids = append(ids, name)
			s = append(s, mapping)
		}

		d.SetId(dataResourceIdHash(ids))
		if err := d.Set("alicloud_log_stores", s); err != nil {
			return fmt.Errorf("Setting alicloud_log_stores got an error: %#v.", err)
		}

		// create a json file in current directory and write data source to it.
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), s)
		}
	}

	return nil
}
