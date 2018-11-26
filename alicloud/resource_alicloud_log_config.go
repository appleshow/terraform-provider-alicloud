package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogConfigCreate,
		Read:   resourceAlicloudLogConfigRead,
		Update: resourceAlicloudLogConfigUpdate,
		Delete: resourceAlicloudLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"store_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"config_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"log_path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"file_pattern": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"log_sample": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"keys": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"topic_format": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"local_storage": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"time_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_format": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_begin_regex": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  ".*",
			},
			"regex": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "(.*)",
			},
			"filter_keys": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter_regex": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"alicloud_log_config": {
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

func resourceAlicloudLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	projectName := d.Get("project_name").(string)
	storeName := d.Get("store_name").(string)
	configName := d.Get("config_name").(string)
	logPath := d.Get("log_path").(string)
	filePattern := d.Get("file_pattern").(string)

	// InputDetail parameter
	inputDetail := sls.InputDetail{}
	inputDetail.LogType = "common_reg_log"
	inputDetail.LogPath = logPath
	inputDetail.FilePattern = filePattern
	if topicFormat, ok := d.GetOk("topic_format"); ok {
		inputDetail.TopicFormat = topicFormat.(string)
	}
	inputDetail.LocalStorage = d.Get("local_storage").(bool)
	if timeKey, ok := d.GetOk("time_key"); ok {
		inputDetail.TimeKey = timeKey.(string)
	}
	if timeFormat, ok := d.GetOk("time_format"); ok {
		inputDetail.TimeFormat = timeFormat.(string)
	}
	if logBeginRegex, ok := d.GetOk("log_begin_regex"); ok {
		inputDetail.LogBeginRegex = logBeginRegex.(string)
	}
	if regex, ok := d.GetOk("regex"); ok {
		inputDetail.Regex = regex.(string)
	}
	if keys, ok := d.GetOk("keys"); ok {
		var keysStr []string
		for _, v := range keys.([]interface{}) {
			keysStr = append(keysStr, v.(string))
		}
		inputDetail.Keys = keysStr
	} else {
		inputDetail.Keys = []string{"content"}
	}
	if filterKeys, ok := d.GetOk("filter_keys"); ok {
		var filterKeysStr []string
		for _, v := range filterKeys.([]interface{}) {
			filterKeysStr = append(filterKeysStr, v.(string))
		}
		inputDetail.FilterKeys = filterKeysStr
	} else {
		inputDetail.FilterKeys = make([]string, 1)
	}
	if filterRegex, ok := d.GetOk("filter_regex"); ok {
		var filterRegexStr []string
		for _, v := range filterRegex.([]interface{}) {
			filterRegexStr = append(filterRegexStr, v.(string))
		}
		inputDetail.FilterRegex = filterRegexStr
	} else {
		inputDetail.FilterRegex = make([]string, 1)
	}
	if topicFormat, ok := d.GetOk("topic_format"); ok {
		inputDetail.TopicFormat = topicFormat.(string)
	}

	// OutPutDetail parameter
	outputDetail := sls.OutputDetail{
		ProjectName:  projectName,
		LogStoreName: storeName,
	}

	//LogConfig parameter
	logConfig := &sls.LogConfig{}
	logConfig.Name = configName
	logConfig.InputType = "file"
	logConfig.OutputType = "LogService"
	logConfig.InputDetail = inputDetail
	logConfig.OutputDetail = outputDetail
	if logSample, ok := d.GetOk("log_sample"); ok {
		logConfig.LogSample = logSample.(string)
	}

	err := client.slsconn.CreateConfig(projectName, logConfig)
	if err != nil {
		return fmt.Errorf("Creating config of log service got an error: %#v", err)
	}

	d.SetId(projectName + COMMA_SEPARATED + storeName + COMMA_SEPARATED + configName)

	return resourceAlicloudLogConfigUpdate(d, meta)
}

func resourceAlicloudLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parameters := strings.Split(d.Id(), COMMA_SEPARATED)

	logConfig, err := client.slsconn.GetConfig(parameters[0], parameters[2])

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("project_name", parameters[0])
	d.Set("store_name", parameters[1])
	d.Set("config_name", parameters[2])
	d.Set("log_sample", logConfig.LogSample)
	if inputDetail, ok := sls.ConvertToInputDetail(logConfig.InputDetail); ok {
		d.Set("log_path", inputDetail.LogPath)
		d.Set("file_pattern", inputDetail.FilePattern)
		if _, ok := d.GetOk("keys"); ok {
			d.Set("keys", inputDetail.Keys)
		}
		d.Set("topic_format", inputDetail.TopicFormat)
		d.Set("local_storage", inputDetail.LocalStorage)
		d.Set("time_key", inputDetail.TimeKey)
		d.Set("time_format", inputDetail.TimeFormat)
		d.Set("log_begin_regex", inputDetail.LogBeginRegex)
		d.Set("regex", inputDetail.Regex)
		if _, ok := d.GetOk("filter_keys"); ok {
			d.Set("filter_keys", inputDetail.FilterKeys)
		}
		if _, ok := d.GetOk("filter_regex"); ok {
			d.Set("filter_regex", inputDetail.FilterRegex)
		}
	} else {
		return fmt.Errorf("Reading config of log service got an error: %#v", err)
	}

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":            parameters[2],
		"name":          parameters[2],
		"status":        "Available",
		"creation_time": time.Now().Format("2006-01-02 15:04:05"),
		"project_name":  parameters[0],
		"config_name":   parameters[2],
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

	s = append(s, mapping)

	if err := d.Set("alicloud_log_config", s); err != nil {
		return fmt.Errorf("Setting alicloud_log_config got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	d.Partial(true)
	update := false

	if d.HasChange("project_name") && !d.IsNewResource() {
		return fmt.Errorf("Updating stroe of log service got an error: %#v", "Cannot modify parameter project_name")
	}
	if d.HasChange("store_name") && !d.IsNewResource() {
		return fmt.Errorf("Updating stroe of log service got an error: %#v", "Cannot modify parameter store_name")
	}
	if d.HasChange("config_name") && !d.IsNewResource() {
		return fmt.Errorf("Updating stroe of log service got an error: %#v", "Cannot modify parameter config_name")
	}

	if d.HasChange("log_sample") {
		update = true
		d.SetPartial("log_sample")
	}
	if d.HasChange("log_path") {
		update = true
		d.SetPartial("log_path")
	}
	if d.HasChange("file_pattern") {
		update = true
		d.SetPartial("file_pattern")
	}
	if d.HasChange("keys") {
		update = true
		d.SetPartial("keys")
	}
	if d.HasChange("topic_format") {
		update = true
		d.SetPartial("topic_format")
	}
	if d.HasChange("local_storage") {
		update = true
		d.SetPartial("local_storage")
	}
	if d.HasChange("time_key") {
		update = true
		d.SetPartial("time_key")
	}
	if d.HasChange("time_format") {
		update = true
		d.SetPartial("time_format")
	}
	if d.HasChange("log_begin_regex") {
		update = true
		d.SetPartial("log_begin_regex")
	}
	if d.HasChange("regex") {
		update = true
		d.SetPartial("regex")
	}
	if d.HasChange("filter_keys") {
		update = true
		d.SetPartial("filter_keys")
	}
	if d.HasChange("filter_regex") {
		update = true
		d.SetPartial("filter_regex")
	}

	if !d.IsNewResource() && update {
		projectName := d.Get("project_name").(string)
		storeName := d.Get("store_name").(string)
		configName := d.Get("config_name").(string)
		logPath := d.Get("log_path").(string)
		filePattern := d.Get("file_pattern").(string)

		// InputDetail parameter
		inputDetail := sls.InputDetail{}
		inputDetail.LogType = "common_reg_log"
		inputDetail.LogPath = logPath
		inputDetail.FilePattern = filePattern

		if topicFormat, ok := d.GetOk("topic_format"); ok {
			inputDetail.TopicFormat = topicFormat.(string)
		}
		inputDetail.LocalStorage = d.Get("local_storage").(bool)
		if timeKey, ok := d.GetOk("time_key"); ok {
			inputDetail.TimeKey = timeKey.(string)
		}
		if timeFormat, ok := d.GetOk("time_format"); ok {
			inputDetail.TimeFormat = timeFormat.(string)
		}
		if logBeginRegex, ok := d.GetOk("log_begin_regex"); ok {
			inputDetail.LogBeginRegex = logBeginRegex.(string)
		}
		if regex, ok := d.GetOk("regex"); ok {
			inputDetail.Regex = regex.(string)
		}
		if keys, ok := d.GetOk("keys"); ok {
			var keysStr []string
			for _, v := range keys.([]interface{}) {
				keysStr = append(keysStr, v.(string))
			}
			inputDetail.Keys = keysStr
		}
		if filterKeys, ok := d.GetOk("filter_keys"); ok {
			var filterKeysStr []string
			for _, v := range filterKeys.([]interface{}) {
				filterKeysStr = append(filterKeysStr, v.(string))
			}
			inputDetail.FilterKeys = filterKeysStr
		}
		if filterRegex, ok := d.GetOk("filter_regex"); ok {
			var filterRegexStr []string
			for _, v := range filterRegex.([]interface{}) {
				filterRegexStr = append(filterRegexStr, v.(string))
			}
			inputDetail.FilterRegex = filterRegexStr
		}
		if topicFormat, ok := d.GetOk("topic_format"); ok {
			inputDetail.TopicFormat = topicFormat.(string)
		}

		// OutPutDetail parameter
		outputDetail := sls.OutputDetail{
			ProjectName:  projectName,
			LogStoreName: storeName,
		}

		//LogConfig parameter
		logConfig := &sls.LogConfig{}
		logConfig.Name = configName
		logConfig.InputType = "file"
		logConfig.OutputType = "LogService"
		logConfig.InputDetail = inputDetail
		logConfig.OutputDetail = outputDetail
		if logSample, ok := d.GetOk("log_sample"); ok {
			logConfig.LogSample = logSample.(string)
		}

		if err := client.slsconn.UpdateConfig(projectName, logConfig); err != nil {
			return fmt.Errorf("Updating config of log service got an error: %#v", err)
		}
	}

	d.Partial(false)

	return resourceAlicloudLogConfigRead(d, meta)
}

func resourceAlicloudLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		parameters := strings.Split(d.Id(), COMMA_SEPARATED)
		err := client.slsconn.DeleteConfig(parameters[0], parameters[2])

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting config of log service got an error: %#v", err))
		}

		resp, err := client.slsconn.GetLogStore(parameters[0], parameters[2])
		if err != nil {
			if NotFoundError(err) || resp == nil {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe config of log service got an error: %#v", err))
		}
		if resp == nil || resp.Name == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting config of log service got an error: %#v", err))
	})
}
