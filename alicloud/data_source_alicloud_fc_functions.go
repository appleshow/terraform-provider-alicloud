package alicloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudFcFunctions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFcFunctionsRead,

		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"next_token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_fc_functions": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"environment_variables": {
							Type:     schema.TypeMap,
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

func dataSourceAlicloudFcFunctionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	listFunctionsInput := fc.NewListFunctionsInput(d.Get("service_name").(string))
	if prefix, ok := d.GetOk("prefix"); ok {
		listFunctionsInput.WithPrefix(prefix.(string))
	}
	if startKey, ok := d.GetOk("start_key"); ok {
		listFunctionsInput.WithStartKey(startKey.(string))
	}
	if nextToken, ok := d.GetOk("next_token"); ok {
		listFunctionsInput.WithNextToken(nextToken.(string))
	}
	if limit, ok := d.GetOk("limit"); ok {
		if limit, err := strconv.ParseInt(strconv.Itoa(limit.(int)), 10, 32); err != nil {
			return fmt.Errorf("List functions of function compute got an error: %#v", err)
		} else {
			listFunctionsInput.WithLimit(int32(limit))
		}
	}

	listFunctionsOutput, err := client.fcconn.ListFunctions(listFunctionsInput)
	if err != nil {
		return fmt.Errorf("List functions of function compute got an error: %#v", err)
	} else {
		log.Printf("[DEBUG] alicloud_fc - Functions found: %#v", listFunctionsOutput.Functions)

		var ids []string
		var s []map[string]interface{}
		for _, function := range listFunctionsOutput.Functions {
			mapping := map[string]interface{}{
				"id":                    *function.FunctionID,
				"name":                  *function.FunctionName,
				"status":                "Available",
				"creation_time":         *function.CreatedTime,
				"description":           *function.Description,
				"runtime":               *function.Runtime,
				"handler":               *function.Handler,
				"timeout":               strconv.Itoa(int(*function.Timeout)),
				"memory_size":           strconv.Itoa(int(*function.MemorySize)),
				"code_size":             strconv.FormatInt(*function.CodeSize, 10),
				"last_modified_time":    *function.LastModifiedTime,
				"environment_variables": function.EnvironmentVariables,
				"resource_type":         "alicloud_fc_function",
			}
			ids = append(ids, *function.FunctionID)
			s = append(s, mapping)
		}

		d.SetId(dataResourceIdHash(ids))
		if err := d.Set("alicloud_fc_functions", s); err != nil {
			return err
		}

		// create a json file in current directory and write data source to it.
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), s)
		}
	}

	return nil
}
