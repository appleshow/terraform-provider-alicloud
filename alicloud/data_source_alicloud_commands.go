package alicloud

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudCommands() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCommandsRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"RunBatScript", "RunPowerShellScript", "RunShellScript",
				}),
			},
			"page_number": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"page_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  50,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_commands": {
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

func dataSourceAlicloudCommandsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	describeCommandsRequest := ecs.CreateDescribeCommandsRequest()
	if id, ok := d.GetOk("id"); ok {
		describeCommandsRequest.CommandId = id.(string)
	}
	if name, ok := d.GetOk("name"); ok {
		describeCommandsRequest.Name = name.(string)
	}
	if commandType, ok := d.GetOk("type"); ok {
		describeCommandsRequest.Type = commandType.(string)
	}
	if pageNumber, ok := d.GetOk("page_number"); ok {
		describeCommandsRequest.PageNumber = requests.NewInteger(pageNumber.(int))
	}
	if pageSize, ok := d.GetOk("page_size"); ok {
		describeCommandsRequest.PageSize = requests.NewInteger(pageSize.(int))
	}

	describeCommandsResponse, err := client.aliecsconn.DescribeCommands(describeCommandsRequest)

	if err != nil {
		return fmt.Errorf("List commands got an error: %#v", err)
	} else {
		log.Printf("[DEBUG] alicloud - commands found: %#v", describeCommandsResponse)

		var ids []string
		var s []map[string]interface{}
		for _, command := range describeCommandsResponse.Commands.Command {
			commandContent, err := base64.StdEncoding.DecodeString(command.CommandContent)
			if err != nil {
				return fmt.Errorf("Reading command got an error: %#v", err)
			}

			mapping := map[string]interface{}{
				"id":              command.CommandId,
				"name":            command.Name,
				"status":          "Available",
				"creation_time":   "",
				"type":            command.Type,
				"description":     command.Description,
				"command_content": string(commandContent),
				"working_dir":     command.WorkingDir,
				"time_out":        command.Timeout,
				"resource_type":   "alicloud_command",
			}
			ids = append(ids, command.CommandId)
			s = append(s, mapping)
		}

		d.SetId(dataResourceIdHash(ids))
		if err := d.Set("alicloud_commands", s); err != nil {
			return err
		}

		// create a json file in current directory and write data source to it.
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), s)
		}
	}

	return nil
}
