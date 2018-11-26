package alicloud

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudCommand() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCommandCreate,
		Read:   resourceAlicloudCommandRead,
		Update: resourceAlicloudCommandUpdate,
		Delete: resourceAlicloudCommandDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"RunBatScript", "RunPowerShellScript", "RunShellScript",
				}),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"command_content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"working_dir": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_out": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			// Computed values
			"alicloud_command": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"working_dir": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_out": {
							Type:     schema.TypeInt,
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

func resourceAlicloudCommandCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	createCommandRequest := ecs.CreateCreateCommandRequest()

	createCommandRequest.Name = d.Get("name").(string)
	createCommandRequest.Type = d.Get("type").(string)
	if description, ok := d.GetOk("description"); ok {
		createCommandRequest.Description = description.(string)
	}
	commandContent := d.Get("command_content").(string)
	createCommandRequest.CommandContent = base64.StdEncoding.EncodeToString([]byte(commandContent))
	if workingDir, ok := d.GetOk("working_dir"); ok {
		createCommandRequest.WorkingDir = workingDir.(string)
	}
	if timeOut, ok := d.GetOk("time_out"); ok {
		createCommandRequest.Timeout = requests.NewInteger(timeOut.(int))
	}

	createCommandResponse, err := client.aliecsconn.CreateCommand(createCommandRequest)
	if err != nil {
		return fmt.Errorf("Creating command got an error: %#v", err)
	}

	d.SetId(createCommandResponse.CommandId)

	return resourceAlicloudCommandUpdate(d, meta)
}

func resourceAlicloudCommandRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	describeCommandsRequest := ecs.CreateDescribeCommandsRequest()
	describeCommandsRequest.CommandId = d.Id()

	describeCommandsResponse, err := client.aliecsconn.DescribeCommands(describeCommandsRequest)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", describeCommandsResponse.Commands.Command[0].Name)
	d.Set("type", describeCommandsResponse.Commands.Command[0].Type)
	if _, ok := d.GetOk("description"); ok {
		d.Set("description", describeCommandsResponse.Commands.Command[0].Description)
	}

	commandContent, err := base64.StdEncoding.DecodeString(describeCommandsResponse.Commands.Command[0].CommandContent)
	if err != nil {
		return fmt.Errorf("Reading command got an error: %#v", err)
	} else {
		d.Set("command_content", string(commandContent))
	}

	if _, ok := d.GetOk("working_dir"); ok {
		d.Set("working_dir", describeCommandsResponse.Commands.Command[0].WorkingDir)
	}
	if _, ok := d.GetOk("time_out"); ok {
		d.Set("time_out", describeCommandsResponse.Commands.Command[0].Timeout)
	}

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":              describeCommandsResponse.Commands.Command[0].CommandId,
		"name":            describeCommandsResponse.Commands.Command[0].Name,
		"status":          "Available",
		"creation_time":   time.Now().Format("2006-01-02 15:04:05"),
		"type":            describeCommandsResponse.Commands.Command[0].Type,
		"description":     describeCommandsResponse.Commands.Command[0].Description,
		"command_content": string(commandContent),
		"working_dir":     describeCommandsResponse.Commands.Command[0].WorkingDir,
		"time_out":        describeCommandsResponse.Commands.Command[0].Timeout,
		"resource_type":   "alicloud_command",
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_command", s); err != nil {
		return fmt.Errorf("Setting alicloud_command got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudCommandUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	d.Partial(true)
	update := false

	modifyCommandRequest := ecs.CreateModifyCommandRequest()
	modifyCommandRequest.CommandId = d.Id()

	if d.HasChange("name") {
		update = true
		modifyCommandRequest.Name = d.Get("name").(string)
		d.SetPartial("name")
	}
	if d.HasChange("type") && !d.IsNewResource() {
		return fmt.Errorf("Updating command got an error: %#v", "Modifying the parameter of Type is not supported.")
	}
	if d.HasChange("description") {
		update = true
		modifyCommandRequest.Description = d.Get("description").(string)
		d.SetPartial("description")
	}
	if d.HasChange("command_content") {
		update = true
		commandContent := d.Get("command_content").(string)
		modifyCommandRequest.CommandContent = base64.StdEncoding.EncodeToString([]byte(commandContent))
		d.SetPartial("command_content")
	}
	if d.HasChange("working_dir") {
		update = true
		modifyCommandRequest.WorkingDir = d.Get("working_dir").(string)
		d.SetPartial("working_dir")
	}
	if d.HasChange("time_out") {
		update = true
		modifyCommandRequest.Timeout = requests.NewInteger(d.Get("time_out").(int))
		d.SetPartial("time_out")
	}

	if !d.IsNewResource() && update {
		if _, err := client.aliecsconn.ModifyCommand(modifyCommandRequest); err != nil {
			return fmt.Errorf("Updating command got an error: %#v", err)
		}
	}

	d.Partial(false)

	return resourceAlicloudCommandRead(d, meta)
}

func resourceAlicloudCommandDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	deleteCommandRequest := ecs.CreateDeleteCommandRequest()
	deleteCommandRequest.CommandId = d.Id()

	describeCommandsRequest := ecs.CreateDescribeCommandsRequest()
	describeCommandsRequest.CommandId = d.Id()

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.aliecsconn.DeleteCommand(deleteCommandRequest)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting command got an error: %#v", err))
		}

		resp, err := client.aliecsconn.DescribeCommands(describeCommandsRequest)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe command got an error: %#v", err))
		}
		if resp.Commands.Command == nil || len(resp.Commands.Command) == 0 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting command got an error: %#v", err))
	})
}
