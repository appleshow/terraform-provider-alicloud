package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudLogConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudLogConfigsRead,

		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_log_configs": {
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
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_path": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_pattern": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_sample": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"keys": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"topic_format": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_storage": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"time_key": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_format": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_begin_regex": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"regex": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_keys": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_regex": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceAlicloudLogConfigsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	projectName := d.Get("project_name").(string)
	offset := d.Get("offset").(int)
	size := d.Get("size").(int)
	var listLogConfigOutput []string

	if groupName, ok := d.GetOk("group_name"); ok {
		if listLogConfigOutputTmp, err := client.slsconn.GetAppliedConfigs(projectName, groupName.(string)); err != nil {
			return fmt.Errorf("List machine configs of log service got an error: %#v", err)
		} else {
			listLogConfigOutput = listLogConfigOutputTmp
		}
	} else {
		if listLogConfigOutputTmp, _, err := client.slsconn.ListConfig(projectName, offset, size); err != nil {
			return fmt.Errorf("List machine configs of log service got an error: %#v", err)
		} else {
			listLogConfigOutput = listLogConfigOutputTmp
		}
	}

	log.Printf("[DEBUG] alicloud_log_service - configs found: %#v", listLogConfigOutput)

	var ids []string
	var s []map[string]interface{}
	for _, name := range listLogConfigOutput {
		logConfig, err := client.slsconn.GetConfig(projectName, name)

		if err != nil {
			return fmt.Errorf("Get config of log service got an error: %#v", err)
		}
		mapping := map[string]interface{}{
			"id":            name,
			"name":          name,
			"status":        "Available",
			"creation_time": "",
			"project_name":  projectName,
			"config_name":   name,
			"resource_type": "alicloud_log_config",
		}

		if inputDetail, ok := sls.ConvertToInputDetail(logConfig.InputDetail); ok {
			mapping["log_path"] = inputDetail.LogPath
			mapping["file_pattern"] = inputDetail.FilePattern
			if _, ok := d.GetOk("keys"); ok {
				mapping["keys"] = inputDetail.Keys
			}
			mapping["topic_format"] = inputDetail.TopicFormat
			mapping["local_storage"] = inputDetail.LocalStorage
			mapping["time_key"] = inputDetail.TimeKey
			mapping["time_format"] = inputDetail.TimeFormat
			mapping["log_begin_regex"] = inputDetail.LogBeginRegex
			mapping["regex"] = inputDetail.Regex
			if _, ok := d.GetOk("filter_keys"); ok {
				mapping["filter_keys"] = inputDetail.FilterKeys
			}
			if _, ok := d.GetOk("filter_regex"); ok {
				mapping["filter_regex"] = inputDetail.FilterRegex
			}
		}

		ids = append(ids, name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alicloud_log_configs", s); err != nil {
		return fmt.Errorf("Setting alicloud_log_configs got an error: %#v.", err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
