package alicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudLogProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudLogProjectsRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_log_projects": {
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
						"project_name": {
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

func dataSourceAlicloudLogProjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	listProjectOutput, err := client.slsconn.ListProject()
	if err != nil {
		return fmt.Errorf("List projects of log service got an error: %#v", err)
	} else {
		log.Printf("[DEBUG] alicloud_log_service - projects found: %#v", listProjectOutput)

		var ids []string
		var s []map[string]interface{}
		for _, name := range listProjectOutput {
			mapping := map[string]interface{}{
				"id":            name,
				"name":          name,
				"status":        "Available",
				"creation_time": "",
				"project_name":  name,
				"resource_type": "alicloud_log_project",
			}
			ids = append(ids, name)
			s = append(s, mapping)
		}

		d.SetId(dataResourceIdHash(ids))
		if err := d.Set("alicloud_log_projects", s); err != nil {
			return fmt.Errorf("Setting alicloud_log_projects got an error: %#v.", err)
		}

		// create a json file in current directory and write data source to it.
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), s)
		}
	}

	return nil
}
