package alicloud

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudCommandInvokeResults() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCommandInvokeResultsRead,

		Schema: map[string]*schema.Schema{
			"invoke_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"command_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"invoke_record_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Running", "Finished", "Stopped", "Failed", "PartialFailed",
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
				Default:  20,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_command_invoke_results": {
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
						"invoke_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"invoke_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"frequency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"page_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"page_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"invoke_instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_invoke_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"invocation_results": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"invoke_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"command_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"finished_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"invoke_record_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"output": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"exit_code": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
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

func dataSourceAlicloudCommandInvokeResultsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	describeInvocationResultsRequest := ecs.CreateDescribeInvocationResultsRequest()
	if invokeId, ok := d.GetOk("invoke_id"); ok {
		describeInvocationResultsRequest.InvokeId = invokeId.(string)
	}
	if commandId, ok := d.GetOk("command_id"); ok {
		describeInvocationResultsRequest.CommandId = commandId.(string)
	}
	if instanceId, ok := d.GetOk("instance_id"); ok {
		describeInvocationResultsRequest.InstanceId = instanceId.(string)
	}
	if invokeRecordStatus, ok := d.GetOk("invoke_record_status"); ok {
		describeInvocationResultsRequest.InvokeRecordStatus = invokeRecordStatus.(string)
	}
	if pageNumber, ok := d.GetOk("page_number"); ok {
		describeInvocationResultsRequest.PageNumber = requests.NewInteger(pageNumber.(int))
	}
	if pageSize, ok := d.GetOk("page_size"); ok {
		describeInvocationResultsRequest.PageSize = requests.NewInteger(pageSize.(int))
	}

	describeInvocationResultsResponse, err := client.aliecsconn.DescribeInvocationResults(describeInvocationResultsRequest)

	if err != nil {
		return fmt.Errorf("List commands invoke results got an error: %#v", err)
	} else {
		log.Printf("[DEBUG] alicloud - command invoke results found: %#v", describeInvocationResultsResponse)

		var ids []string
		var s []map[string]interface{}
		invocation := describeInvocationResultsResponse.Invocation

		var invokeInstances []map[string]interface{}
		var invocationResults []map[string]interface{}

		for _, invokeInstance := range invocation.InvokeInstances.InvokeInstance {
			mapping := map[string]interface{}{
				"instance_id":            invokeInstance.InstanceId,
				"instance_invoke_status": invokeInstance.InstanceInvokeStatus,
			}
			invokeInstances = append(invokeInstances, mapping)
		}
		for _, invocationResult := range invocation.InvocationResults.InvocationResult {
			output, err := base64.StdEncoding.DecodeString(invocationResult.Output)
			if err != nil {
				return fmt.Errorf("Reading command invoke output got an error: %#v", err)
			}

			mapping := map[string]interface{}{
				"invoke_id":            invocationResult.InvokeId,
				"command_id":           invocationResult.CommandId,
				"instance_id":          invocationResult.InstanceId,
				"finished_time":        invocationResult.FinishedTime,
				"invoke_record_status": invocationResult.InvokeRecordStatus,
				"output":               string(output),
				"exit_code":            invocationResult.ExitCode,
			}
			invocationResults = append(invocationResults, mapping)
		}

		mapping := map[string]interface{}{
			"id":                 invocation.InvokeId,
			"name":               invocation.CommandName,
			"status":             "Available",
			"creation_time":      "",
			"invoke_id":          invocation.InvokeId,
			"command_id":         invocation.CommandId,
			"command_name":       invocation.CommandName,
			"command_type":       invocation.CommandType,
			"invoke_status":      invocation.InvokeStatus,
			"timed":              invocation.Timed,
			"frequency":          invocation.Frequency,
			"total_count":        invocation.TotalCount,
			"page_number":        invocation.PageNumber,
			"page_size":          invocation.PageSize,
			"invoke_instances":   invokeInstances,
			"invocation_results": invocationResults,
			"resource_type":      "alicloud_command_invoke_result",
		}
		ids = append(ids, invocation.InvokeId)
		s = append(s, mapping)

		d.SetId(dataResourceIdHash(ids))
		if err := d.Set("alicloud_command_invoke_results", s); err != nil {
			return err
		}

		// create a json file in current directory and write data source to it.
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), s)
		}
	}

	return nil
}
