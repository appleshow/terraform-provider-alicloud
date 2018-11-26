package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudCommandInvokes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCommandInvokesRead,

		Schema: map[string]*schema.Schema{
			"invoke_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"command_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"command_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"command_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"RunBatScript", "RunPowerShellScript", "RunShellScript",
				}),
			},
			"invoke_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Running", "Finished", "Stopped", "Failed", "PartialFailed",
				}),
			},
			"timed": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
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
			"alicloud_command_invokes": {
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

func dataSourceAlicloudCommandInvokesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	describeInvocationsRequest := ecs.CreateDescribeInvocationsRequest()
	if invokeId, ok := d.GetOk("invoke_id"); ok {
		describeInvocationsRequest.InvokeId = invokeId.(string)
	}
	if commandId, ok := d.GetOk("command_id"); ok {
		describeInvocationsRequest.CommandId = commandId.(string)
	}
	if instanceId, ok := d.GetOk("instance_id"); ok {
		describeInvocationsRequest.InstanceId = instanceId.(string)
	}
	if commandName, ok := d.GetOk("command_name"); ok {
		describeInvocationsRequest.CommandName = commandName.(string)
	}
	if commandType, ok := d.GetOk("command_type"); ok {
		describeInvocationsRequest.CommandType = commandType.(string)
	}
	if invokeStatus, ok := d.GetOk("invoke_status"); ok {
		describeInvocationsRequest.InvokeStatus = invokeStatus.(string)
	}
	if timed, ok := d.GetOk("timed"); ok {
		describeInvocationsRequest.Timed = requests.NewBoolean(timed.(bool))
	}
	if pageNumber, ok := d.GetOk("page_number"); ok {
		describeInvocationsRequest.PageNumber = requests.NewInteger(pageNumber.(int))
	}
	if pageSize, ok := d.GetOk("page_size"); ok {
		describeInvocationsRequest.PageSize = requests.NewInteger(pageSize.(int))
	}

	describeInvocationsResponse, err := client.aliecsconn.DescribeInvocations(describeInvocationsRequest)

	if err != nil {
		return fmt.Errorf("List commands invokes got an error: %#v", err)
	} else {
		log.Printf("[DEBUG] alicloud - command invokes found: %#v", describeInvocationsResponse)

		var ids []string
		var s []map[string]interface{}
		for _, invocation := range describeInvocationsResponse.Invocations.Invocation {
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
				mapping := map[string]interface{}{
					"invoke_id":            invocationResult.InvokeId,
					"command_id":           invocationResult.CommandId,
					"instance_id":          invocationResult.InstanceId,
					"finished_time":        invocationResult.FinishedTime,
					"invoke_record_status": invocationResult.InvokeRecordStatus,
					"output":               invocationResult.Output,
					"exit_code":            invocationResult.ExitCode,
				}
				invocationResults = append(invocationResults, mapping)
			}

			mapping := map[string]interface{}{
				"id":                 invocation.InvokeId,
				"name":               invocation.CommandName,
				"status":             invocation.InvokeStatus,
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
				"resource_type":      "alicloud_command_invoke",
			}
			ids = append(ids, invocation.InvokeId)
			s = append(s, mapping)
		}

		d.SetId(dataResourceIdHash(ids))
		if err := d.Set("alicloud_command_invokes", s); err != nil {
			return err
		}

		// create a json file in current directory and write data source to it.
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), s)
		}
	}

	return nil
}
