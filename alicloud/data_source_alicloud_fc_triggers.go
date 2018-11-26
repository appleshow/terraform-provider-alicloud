package alicloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudFcTriggers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFcTriggersRead,

		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"function_name": &schema.Schema{
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
			"alicloud_fc_triggers": {
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
						"source_arm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"invocation_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"raw_rrigger_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified_time": {
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

func dataSourceAlicloudFcTriggersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	listTriggersInput := fc.NewListTriggersInput(d.Get("service_name").(string), d.Get("function_name").(string))
	if prefix, ok := d.GetOk("prefix"); ok {
		listTriggersInput.WithPrefix(prefix.(string))
	}
	if startKey, ok := d.GetOk("start_key"); ok {
		listTriggersInput.WithStartKey(startKey.(string))
	}
	if nextToken, ok := d.GetOk("next_token"); ok {
		listTriggersInput.WithNextToken(nextToken.(string))
	}
	if limit, ok := d.GetOk("limit"); ok {
		if limit, err := strconv.ParseInt(strconv.Itoa(limit.(int)), 10, 32); err != nil {
			return fmt.Errorf("List triggers of function compute got an error: %#v", err)
		} else {
			listTriggersInput.WithLimit(int32(limit))
		}
	}

	listTriggersOutput, err := client.fcconn.ListTriggers(listTriggersInput)
	if err != nil {
		return fmt.Errorf("List triggers of function compute got an error: %#v", err)
	} else {
		log.Printf("[DEBUG] alicloud_fc - Triggers found: %#v", listTriggersOutput.Triggers)

		var ids []string
		var s []map[string]interface{}
		for _, trigger := range listTriggersOutput.Triggers {
			mapping := map[string]interface{}{
				"id":            *trigger.TriggerName,
				"name":          *trigger.TriggerName,
				"status":        "Available",
				"creation_time": *trigger.CreatedTime,
				//"source_arm":         *trigger.SourceARN,
				"trigger_type": *trigger.TriggerType,
				//"invocation_role":    *trigger.InvocationRole,
				"raw_rrigger_config": string(trigger.RawTriggerConfig[:]),
				"last_modified_time": *trigger.LastModifiedTime,
				"resource_type":      "alicloud_fc_trigger",
			}
			ids = append(ids, *trigger.TriggerName)
			s = append(s, mapping)
		}

		d.SetId(dataResourceIdHash(ids))
		if err := d.Set("alicloud_fc_triggers", s); err != nil {
			return err
		}

		// create a json file in current directory and write data source to it.
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), s)
		}
	}

	return nil
}
