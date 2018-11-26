package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudLogConfigToMachineGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogConfigToMachineGroupCreate,
		Read:   resourceAlicloudLogConfigToMachineGroupRead,
		Update: resourceAlicloudLogConfigToMachineGroupUpdate,
		Delete: resourceAlicloudLogConfigToMachineGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"config_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// Computed values
			"alicloud_log_configtomachinegroup": {
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
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

/**
*
 */
func resourceAlicloudLogConfigToMachineGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	projectName := d.Get("project_name").(string)
	configName := d.Get("config_name").(string)
	groupName := d.Get("group_name").(string)

	err := client.slsconn.ApplyConfigToMachineGroup(projectName, configName, groupName)
	if err != nil {
		return fmt.Errorf("Applying config to machine group of log service got an error: %#v", err)
	}

	d.SetId(projectName + COMMA_SEPARATED + configName + COMMA_SEPARATED + groupName)

	return resourceAlicloudLogConfigToMachineGroupUpdate(d, meta)
}

/**
*
 */
func resourceAlicloudLogConfigToMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parameters := strings.Split(d.Id(), COMMA_SEPARATED)
	machineGroup, err := client.slsconn.GetAppliedMachineGroups(parameters[0], parameters[1])

	if err != nil {
		if NotFoundError(err) || machineGroup == nil || len(machineGroup) == 0 {
			d.SetId("")
			return nil
		}
		return err
	}

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":            d.Id(),
		"name":          d.Id(),
		"status":        "Available",
		"creation_time": time.Now().Format("2006-01-02 15:04:05"),
		"resource_type": "alicloud_log_configtomachinegroup",
		"project_name":  d.Get("project_name").(string),
		"config_name":   d.Get("config_name").(string),
		"group_name":    d.Get("group_name").(string),
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_log_configtomachinegroup", s); err != nil {
		return fmt.Errorf("Setting alicloud_log_configtomachinegroup got an error: %#v.", err)
	}

	return nil
}

/**
*
 */
func resourceAlicloudLogConfigToMachineGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)
	update := false

	if d.HasChange("project_name") && !d.IsNewResource() {
		return fmt.Errorf("Applying config to machine group of log service got an error: %#v", "Cannot modify parameter project_name")
	}
	if d.HasChange("config_name") && !d.IsNewResource() {
		return fmt.Errorf("Applying config to machine group of log service got an error: %#v", "Cannot modify parameter config_name")
	}
	if d.HasChange("group_name") && !d.IsNewResource() {
		return fmt.Errorf("Applying config to machine group of log service got an error: %#v", "Cannot modify parameter group_name")
	}

	if !d.IsNewResource() && update {
		// Nothing to do
	}

	d.Partial(false)

	return resourceAlicloudLogConfigToMachineGroupRead(d, meta)
}

func resourceAlicloudLogConfigToMachineGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		parameters := strings.Split(d.Id(), COMMA_SEPARATED)
		err := client.slsconn.RemoveConfigFromMachineGroup(parameters[0], parameters[1], parameters[2])

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Removing config to machine group of log service got an error: %#v", err))
		}

		resp, err := client.slsconn.GetAppliedMachineGroups(parameters[0], parameters[1])
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe config to machine group of log service got an error: %#v", err))
		}
		if resp == nil || len(resp) == 0 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Removing config to machine of log service got an error: %#v", err))
	})
}
