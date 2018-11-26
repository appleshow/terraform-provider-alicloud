package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudLogMachineGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogMachineGroupCreate,
		Read:   resourceAlicloudLogMachineGroupRead,
		Update: resourceAlicloudLogMachineGroupUpdate,
		Delete: resourceAlicloudLogMachineGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"machine_id_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ip", "userdefined"}),
			},
			"machine_id_list": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"attribute_external_mame": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"attribute_topic_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			// Computed values
			"alicloud_machine_group": {
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

/**
*
 */
func resourceAlicloudLogMachineGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	projectName := d.Get("project_name").(string)
	groupName := d.Get("group_name").(string)

	attribute := sls.MachinGroupAttribute{
		ExternalName: d.Get("attribute_external_mame").(string),
		TopicName:    d.Get("attribute_topic_name").(string),
	}
	var machineList []string
	for _, v := range d.Get("machine_id_list").([]interface{}) {
		machineList = append(machineList, v.(string))
	}
	var machineGroup = &sls.MachineGroup{
		Name:          groupName,
		Type:          d.Get("type").(string),
		MachineIDType: d.Get("machine_id_type").(string),
		MachineIDList: machineList,
		Attribute:     attribute,
	}

	err := client.slsconn.CreateMachineGroup(projectName, machineGroup)
	if err != nil {
		return fmt.Errorf("Creating machine group of log service got an error: %#v", err)
	}

	d.SetId(projectName + COMMA_SEPARATED + groupName)

	return resourceAlicloudLogMachineGroupUpdate(d, meta)
}

/**
*
 */
func resourceAlicloudLogMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parameters := strings.Split(d.Id(), COMMA_SEPARATED)
	machineGroup, err := client.slsconn.GetMachineGroup(parameters[0], parameters[1])

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("project_name", parameters[0])
	d.Set("group_name", parameters[1])
	d.Set("machine_id_type", machineGroup.MachineIDType)
	d.Set("machine_id_list", machineGroup.MachineIDList)
	d.Set("type", machineGroup.Type)
	d.Set("attribute_external_mame", machineGroup.Attribute.ExternalName)
	d.Set("attribute_topic_name", machineGroup.Attribute.TopicName)

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":            parameters[1],
		"name":          parameters[1],
		"status":        "Available",
		"creation_time": time.Now().Format("2006-01-02 15:04:05"),
		"project_name":  parameters[0],
		"group_name":    parameters[1],
		"resource_type": "alicloud_machine_groups",
	}
	mapping["machine_id_type"] = machineGroup.MachineIDType
	mapping["machine_id_list"] = machineGroup.MachineIDList
	mapping["type"] = machineGroup.Type
	mapping["attribute_external_mame"] = machineGroup.Attribute.ExternalName
	mapping["attribute_topic_name"] = machineGroup.Attribute.TopicName

	s = append(s, mapping)

	if err := d.Set("alicloud_machine_group", s); err != nil {
		return fmt.Errorf("Setting alicloud_machine_group got an error: %#v.", err)
	}

	return nil
}

/**
*
 */
func resourceAlicloudLogMachineGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	d.Partial(true)
	update := false

	if d.HasChange("project_name") && !d.IsNewResource() {
		return fmt.Errorf("Updating machine group of log service got an error: %#v", "Cannot modify parameter project_name")
	}
	if d.HasChange("group_name") && !d.IsNewResource() {
		return fmt.Errorf("Updating machine group of log service got an error: %#v", "Cannot modify parameter group_name")
	}
	if d.HasChange("machine_id_type") {
		update = true
		d.SetPartial("machine_id_type")
	}
	if d.HasChange("machine_id_list") {
		update = true
		d.SetPartial("machine_id_list")
	}
	if d.HasChange("type") {
		update = true
		d.SetPartial("type")
	}
	if d.HasChange("attribute_external_mame") {
		update = true
		d.SetPartial("attribute_external_mame")
	}
	if d.HasChange("attribute_topic_name") {
		update = true
		d.SetPartial("attribute_topic_name")
	}

	if !d.IsNewResource() && update {
		projectName := d.Get("project_name").(string)
		groupName := d.Get("group_name").(string)

		attribute := sls.MachinGroupAttribute{
			ExternalName: d.Get("attribute_external_mame").(string),
			TopicName:    d.Get("attribute_topic_name").(string),
		}
		var machineList []string
		for _, v := range d.Get("machine_id_list").([]interface{}) {
			machineList = append(machineList, v.(string))
		}
		var machineGroup = &sls.MachineGroup{
			Name:          groupName,
			Type:          d.Get("type").(string),
			MachineIDType: d.Get("machine_id_type").(string),
			MachineIDList: machineList,
			Attribute:     attribute,
		}

		if err := client.slsconn.UpdateMachineGroup(projectName, machineGroup); err != nil {
			return fmt.Errorf("Updating machine group of log service got an error: %#v", err)
		}
	}

	d.Partial(false)

	return resourceAlicloudLogMachineGroupRead(d, meta)
}

func resourceAlicloudLogMachineGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		parameters := strings.Split(d.Id(), COMMA_SEPARATED)
		err := client.slsconn.DeleteMachineGroup(parameters[0], parameters[1])

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting machine group of log service got an error: %#v", err))
		}

		resp, err := client.slsconn.GetMachineGroup(parameters[0], parameters[1])
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe machie group of log service got an error: %#v", err))
		}
		if resp == nil || resp.Name == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting machie group of log service got an error: %#v", err))
	})
}
