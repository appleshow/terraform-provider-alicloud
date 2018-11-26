package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudCommandInvoke() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCommandInvokeCreate,
		Read:   resourceAlicloudCommandInvokeRead,
		Update: resourceAlicloudCommandInvokeUpdate,
		Delete: resourceAlicloudCommandInvokeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"command_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_ids": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
				MaxItems: 100,
			},
			"timed": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"frequency": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_command_invoke": {
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

func resourceAlicloudCommandInvokeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	invokeCommandRequest := ecs.CreateInvokeCommandRequest()
	invokeCommandRequest.CommandId = d.Get("command_id").(string)

	var instanceIdsStr []string
	instanceIds := d.Get("instance_ids").([]interface{})
	for _, instanceId := range instanceIds {
		instanceIdsStr = append(instanceIdsStr, instanceId.(string))
	}

	describeInstancesRequest := ecs.CreateDescribeInstancesRequest()
	if bytes, err := json.Marshal(instanceIdsStr); err != nil {
		return fmt.Errorf("Marshaling instanceIds to json string got an error: %#v.", err)
	} else {
		describeInstancesRequest.InstanceIds = string(bytes[:])
		describeInstancesRequest.PageNumber = requests.NewInteger(1)
		describeInstancesRequest.PageSize = requests.NewInteger(100)
	}
	describeInstancesResponse, err := client.aliecsconn.DescribeInstances(describeInstancesRequest)
	if err != nil {
		return fmt.Errorf("DescribeInstances got an error: %#v.", err)
	}
	for _, instance := range describeInstancesResponse.Instances.Instance {
		if "Running" != instance.Status {
			return fmt.Errorf("Instance[%s,%s] is not running.", instance.InstanceName, instance.InstanceId)
		}
	}

	invokeCommandRequest.InstanceId = &instanceIdsStr

	if timed, ok := d.GetOk("timed"); ok {
		invokeCommandRequest.Timed = requests.NewBoolean(timed.(bool))

		if timed.(bool) {
			if frequency, ok := d.GetOk("frequency"); ok {
				invokeCommandRequest.Frequency = frequency.(string)
			} else {
				return fmt.Errorf("Creating command invoke got an error: %#v", "Parameter frequency can not be null when timed is true")
			}
		}
	}

	invokeCommandResponse, err := client.aliecsconn.InvokeCommand(invokeCommandRequest)
	if err != nil {
		return fmt.Errorf("Creating command invoke got an error: %#v", err)
	}

	d.SetId(invokeCommandResponse.InvokeId)

	return resourceAlicloudCommandInvokeUpdate(d, meta)
}

func resourceAlicloudCommandInvokeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	describeInvocationsRequest := ecs.CreateDescribeInvocationsRequest()
	describeInvocationsRequest.CommandId = d.Get("command_id").(string)
	describeInvocationsRequest.InvokeId = d.Id()
	describeInvocationsRequest.PageNumber = requests.NewInteger(1)
	describeInvocationsRequest.PageSize = requests.NewInteger(50)

	describeInvocationsResponse, err := client.aliecsconn.DescribeInvocations(describeInvocationsRequest)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	} else {
		if describeInvocationsResponse.Invocations.Invocation == nil || len(describeInvocationsResponse.Invocations.Invocation) == 0 {
			d.SetId("")
		}
	}

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
			"creation_time":      time.Now().Format("2006-01-02 15:04:05"),
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
		s = append(s, mapping)
	}

	if err := d.Set("alicloud_command_invoke", s); err != nil {
		return fmt.Errorf("Setting alicloud_command_invoke got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudCommandInvokeUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	if d.HasChange("command_id") && !d.IsNewResource() {
		return fmt.Errorf("Updating command invoke got an error: %#v", "Modifying the parameter of command_id is not supported.")
	}
	if d.HasChange("instance_ids") && !d.IsNewResource() {
		return fmt.Errorf("Updating command invoke got an error: %#v", "Modifying the parameter of instance_ids is not supported.")
	}
	if d.HasChange("timed") && !d.IsNewResource() {
		return fmt.Errorf("Updating command invoke got an error: %#v", "Modifying the parameter of timed is not supported.")
	}
	if d.HasChange("frequency") && !d.IsNewResource() {
		return fmt.Errorf("Updating command invoke got an error: %#v", "Modifying the parameter of frequency is not supported.")
	}

	d.Partial(false)

	return resourceAlicloudCommandInvokeRead(d, meta)
}

func resourceAlicloudCommandInvokeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	stopInvocationRequest := ecs.CreateStopInvocationRequest()
	stopInvocationRequest.InvokeId = d.Id()

	var instanceIdsStr []string
	instanceIds := d.Get("instance_ids").([]interface{})
	for _, instanceId := range instanceIds {
		instanceIdsStr = append(instanceIdsStr, instanceId.(string))
	}
	stopInvocationRequest.InstanceId = &instanceIdsStr

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.aliecsconn.StopInvocation(stopInvocationRequest)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting command invoke got an error: %#v", err))
		}

		return nil
	})

}
