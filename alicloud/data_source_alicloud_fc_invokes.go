package alicloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudFcInvokes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFcInvokesRead,

		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"function_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"payload": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"invocation_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"Async", "Sync"}),
			},
			"log_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"Tail", "None"}),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_variables": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			// Computed values
			"alicloud_fc_invokes": {
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
						"error_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log": {
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

func dataSourceAlicloudFcInvokesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	result := make(map[string]interface{})
	result["id"] = d.Get("function_name").(string)
	result["name"] = d.Get("function_name").(string)

	var ids []string
	var s []map[string]interface{}

	invokeFunctionInput := fc.NewInvokeFunctionInput(d.Get("service_name").(string), d.Get("function_name").(string))
	if payload, ok := d.GetOk("payload"); ok {
		invokeFunctionInput.WithPayload([]byte(payload.(string)))
	}
	if invocationType, ok := d.GetOk("invocation_type"); ok {
		invokeFunctionInput.WithInvocationType(invocationType.(string))
	}
	if logType, ok := d.GetOk("log_type"); ok {
		invokeFunctionInput.WithLogType(logType.(string))
	}
	if environmentVariables, ok := d.GetOk("environment_variables"); ok {
		updateFunctionInput := fc.NewUpdateFunctionInput(d.Get("service_name").(string), d.Get("function_name").(string))
		environmentVariablesStr := make(map[string]string)
		for k, v := range environmentVariables.(map[string]interface{}) {
			environmentVariablesStr[k] = v.(string)
		}
		updateFunctionInput.WithEnvironmentVariables(environmentVariablesStr)
		if _, err := client.fcconn.UpdateFunction(updateFunctionInput); err != nil {
			result["status"] = "Failed"
			result["error_message"] = fmt.Errorf("Updating parameters got an error: %#v", err).Error()
			result["log"] = ""
			result["resource_type"] = "alicloud_fc_invoke"

			ids = append(ids, d.Get("function_name").(string))
			s = append(s, result)
			d.SetId(dataResourceIdHash(ids))
			if err := d.Set("alicloud_fc_invokes", s); err != nil {
				return err
			}

			// create a json file in current directory and write data source to it.
			if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
				writeToFile(output.(string), s)
			}

			return fmt.Errorf("Updating parameters got an error: %#v", err)
		}
	}

	invokeFunctionOutput, err := client.fcconn.InvokeFunction(invokeFunctionInput)
	if err != nil {
		result["status"] = "Failed"
		result["error_message"] = fmt.Errorf("Invoke function of function compute got an error: %#v", err).Error()
		result["log"] = ""

		ids = append(ids, d.Get("function_name").(string))
		s = append(s, result)
		d.SetId(dataResourceIdHash(ids))
		if err := d.Set("alicloud_fc_invokes", s); err != nil {
			return err
		}

		// create a json file in current directory and write data source to it.
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), s)
		}

		return fmt.Errorf("Invoke function of function compute got an error: %#v", err)
	} else {
		result["status"] = "Success"
		result["error_message"] = ""

		logResult, err := invokeFunctionOutput.GetLogResult()
		if err == nil {
			result["log"] = strings.Split(logResult, "\n")
		}
		log.Printf("[DEBUG] alicloud_fc - Invoke function found: %#v", result)

		ids = append(ids, d.Get("function_name").(string))
		s = append(s, result)
		d.SetId(dataResourceIdHash(ids))
		if err := d.Set("alicloud_fc_invokes", s); err != nil {
			return err
		}

		// create a json file in current directory and write data source to it.
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), s)
		}
	}

	return nil
}
