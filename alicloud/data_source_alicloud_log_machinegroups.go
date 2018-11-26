package alicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudLogMachineGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudLogMachineGroupsRead,

		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"config_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_machine_groups": {
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
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"machine_id_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"machine_id_list": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"attribute_external_mame": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"attribute_topic_name": &schema.Schema{
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

func dataSourceAlicloudLogMachineGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	projectName := d.Get("project_name").(string)
	offset := d.Get("offset").(int)
	size := d.Get("size").(int)
	var listMachineGroup []string

	if configName, ok := d.GetOk("config_name"); ok {
		if listMachineGroupTmp, err := client.slsconn.GetAppliedMachineGroups(projectName, configName.(string)); err != nil {
			return fmt.Errorf("List machine groups of log service got an error: %#v", err)
		} else {
			listMachineGroup = listMachineGroupTmp
		}
	} else {
		if listMachineGroupTmp, _, err := client.slsconn.ListMachineGroup(projectName, offset, size); err != nil {
			return fmt.Errorf("List machine groups of log service got an error: %#v", err)
		} else {
			listMachineGroup = listMachineGroupTmp
		}
	}

	log.Printf("[DEBUG] alicloud_log_service - machine groups found: %#v", listMachineGroup)

	var ids []string
	var s []map[string]interface{}
	for _, name := range listMachineGroup {
		machineGroup, err := client.slsconn.GetMachineGroup(projectName, name)
		if err != nil {
			return fmt.Errorf("Get Machine Group got an error: %#v.", err)
		}

		mapping := map[string]interface{}{
			"id":            name,
			"name":          name,
			"status":        "Available",
			"creation_time": "",
			"project_name":  projectName,
			"group_name":    name,
			"resource_type": "alicloud_machine_group",
		}
		mapping["machine_id_type"] = machineGroup.MachineIDType
		mapping["machine_id_list"] = machineGroup.MachineIDList
		mapping["type"] = machineGroup.Type
		mapping["attribute_external_mame"] = machineGroup.Attribute.ExternalName
		mapping["attribute_topic_name"] = machineGroup.Attribute.TopicName

		ids = append(ids, name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alicloud_machine_groups", s); err != nil {
		return fmt.Errorf("Setting alicloud_machine_groups got an error: %#v.", err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
