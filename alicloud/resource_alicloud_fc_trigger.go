package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

type TriggerConfigTimer struct {
	Payload        string `json:"payload"`
	CronExpression string `json:"cronExpression"`
	Enable         bool   `json:"enable"`
}

type TriggerConfigOss struct {
	Events []string `json:"events"`
	Filter struct {
		Key struct {
			Prefix string `json:"prefix"`
			Suffix string `json:"suffix"`
		} `json:"key"`
	} `json:"filter"`
}

type TriggerConfigLog struct {
	SourceConfig struct {
		Logstore string `json:"logstore"`
	} `json:"sourceConfig"`
	JobConfig struct {
		MaxRetryTime    int `json:"maxRetryTime"`
		TriggerInterval int `json:"triggerInterval"`
	} `json:"jobConfig"`
	FunctionParameter *map[string]interface{} `json:"functionParameter"`
	LogConfig         struct {
		Project  string `json:"project"`
		Logstore string `json:"logstore"`
	} `json:"logConfig"`
	Enable bool `json:"enable"`
}

type TriggerConfigHttp struct {
	AuthType string   `json:"authType"`
	Methods  []string `json:"methods"`
}

func resourceAlicloudFcTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFcTriggerCreate,
		Read:   resourceAlicloudFcTriggerRead,
		Update: resourceAlicloudFcTriggerUpdate,
		Delete: resourceAlicloudFcTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"function_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"trigger_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"source_arn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"trigger_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(fc.TRIGGER_TYPE_OSS), string(fc.TRIGGER_TYPE_LOG), string(fc.TRIGGER_TYPE_TIMER),
					string(fc.TRIGGER_TYPE_HTTP),
				}),
			},
			"invocation_role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Comm config
			"config_enable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			// TIMER trigger config
			"config_payload": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_cron_expression": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// OSS trigger config
			"config_events": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_filter_key_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_filter_key_suffix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// LOG trigger config
			"config_source_logstore": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_job_interval": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"config_job_max_retry_time": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"config_function_parameter": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"config_log_project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_log_logstore": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// HTTP trigger config
			"config_auth_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"anonymous", "function"}),
			},
			"config_methods": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_fc_trigger": {
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

func resourceAlicloudFcTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	serviceName := d.Get("service_name").(string)
	functionName := d.Get("function_name").(string)
	triggerName := d.Get("trigger_name").(string)
	createTriggerInput := fc.NewCreateTriggerInput(serviceName, functionName)

	createTriggerInput.WithTriggerName(triggerName)
	createTriggerInput.WithTriggerType(d.Get("trigger_type").(string))

	triggerConfig := make(map[string]interface{})
	switch d.Get("trigger_type").(string) {
	case string(fc.TRIGGER_TYPE_TIMER):
		if enable, ok := d.GetOk("config_enable"); ok {
			triggerConfig["enable"] = enable
		} else {
			d.Set("enable", true)
		}
		if payload, ok := d.GetOk("config_payload"); ok {
			triggerConfig["payload"] = payload
		} else {
			return fmt.Errorf("Creating trigger of function compute got an error: %#v", "Can not find the parameter config_payload")
		}
		if cronExpression, ok := d.GetOk("config_cron_expression"); ok {
			triggerConfig["cronExpression"] = cronExpression
		} else {
			return fmt.Errorf("Creating trigger of function compute got an error: %#v", "Can not find the parameter config_cron_expression")
		}
	case string(fc.TRIGGER_TYPE_OSS):
		createTriggerInput.WithSourceARN(d.Get("source_arn").(string))
		createTriggerInput.WithInvocationRole(d.Get("invocation_role").(string))
		if events, ok := d.GetOk("config_events"); ok {
			triggerConfig["events"] = strings.Split(events.(string), COMMA_SEPARATED)
		} else {
			return fmt.Errorf("Creating trigger of function compute got an error: %#v", "Can not find the parameter config_events")
		}

		filter := make(map[string]interface{})
		key := make(map[string]string)
		if prefix, ok := d.GetOk("config_filter_key_prefix"); ok {
			key["prefix"] = prefix.(string)
		} else {
			return fmt.Errorf("Creating trigger of function compute got an error: %#v", "Can not find the parameter config_filter_key_prefix")
		}
		if suffix, ok := d.GetOk("config_filter_key_suffix"); ok {
			key["suffix"] = suffix.(string)
		} else {
			return fmt.Errorf("Creating trigger of function compute got an error: %#v", "Can not find the parameter config_filter_key_suffix")
		}
		filter["key"] = key
		triggerConfig["filter"] = filter
	case string(fc.TRIGGER_TYPE_LOG):
		createTriggerInput.WithSourceARN(d.Get("source_arn").(string))
		createTriggerInput.WithInvocationRole(d.Get("invocation_role").(string))
		sourceConfig := make(map[string]string)
		jobConfig := make(map[string]int32)
		logConfig := make(map[string]string)
		if enable, ok := d.GetOk("config_enable"); ok {
			triggerConfig["enable"] = enable
		} else {
			d.Set("enable", true)
		}
		if sourceLogstore, ok := d.GetOk("config_source_logstore"); ok {
			sourceConfig["logstore"] = sourceLogstore.(string)
		} else {
			return fmt.Errorf("Creating trigger of function compute got an error: %#v", "Can not find the parameter config_source_logstore")
		}
		if jobInterval, ok := d.GetOk("config_job_interval"); ok {
			if jobInterval32, err := strconv.ParseInt(strconv.Itoa(jobInterval.(int)), 10, 32); err != nil {
				return fmt.Errorf("Creating trigger of function compute got an error: %#v", err)
			} else {
				jobConfig["triggerInterval"] = int32(jobInterval32)
			}
		} else {
			return fmt.Errorf("Creating trigger of function compute got an error: %#v", "Can not find the parameter config_job_interval")
		}
		if maxRetryTime, ok := d.GetOk("config_job_max_retry_time"); ok {
			if maxRetryTime32, err := strconv.ParseInt(strconv.Itoa(maxRetryTime.(int)), 10, 32); err != nil {
				return fmt.Errorf("Creating trigger of function compute got an error: %#v", err)
			} else {
				jobConfig["maxRetryTime"] = int32(maxRetryTime32)
			}
		} else {
			return fmt.Errorf("Creating trigger of function compute got an error: %#v", "Can not find the parameter config_job_max_retry_time")
		}
		if functionParameter, ok := d.GetOk("config_function_parameter"); ok {
			triggerConfig["functionParameter"] = functionParameter
		}
		if project, ok := d.GetOk("config_log_project"); ok {
			logConfig["project"] = project.(string)
		} else {
			return fmt.Errorf("Creating trigger of function compute got an error: %#v", "Can not find the parameter config_log_project")
		}
		if logstore, ok := d.GetOk("config_log_logstore"); ok {
			logConfig["logstore"] = logstore.(string)
		} else {
			return fmt.Errorf("Creating trigger of function compute got an error: %#v", "Can not find the parameter config_log_logstore")
		}
		triggerConfig["sourceConfig"] = sourceConfig
		triggerConfig["jobConfig"] = jobConfig
		triggerConfig["logConfig"] = logConfig
	case string(fc.TRIGGER_TYPE_HTTP):
		if authType, ok := d.GetOk("config_auth_type"); ok {
			triggerConfig["authType"] = authType
		} else {
			return fmt.Errorf("Creating trigger of function compute got an error: %#v", "Can not find the parameter config_auth_type")
		}
		if methods, ok := d.GetOk("config_methods"); ok {
			triggerConfig["methods"] = strings.Split(methods.(string), COMMA_SEPARATED)
		} else {
			return fmt.Errorf("Creating trigger of function compute got an error: %#v", "Can not find the parameter config_methods")
		}
	default:
	}
	createTriggerInput.WithTriggerConfig(triggerConfig)

	_, err := client.fcconn.CreateTrigger(createTriggerInput)
	if err != nil {
		return fmt.Errorf("Creating trigger of function compute got an error: %#v", err)
	}

	d.SetId(serviceName + COMMA_SEPARATED + functionName + COMMA_SEPARATED + triggerName)

	return resourceAlicloudFcTriggerUpdate(d, meta)
}

func resourceAlicloudFcTriggerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parameters := strings.Split(d.Id(), COMMA_SEPARATED)
	getTriggerOutput, err := client.fcconn.GetTrigger(fc.NewGetTriggerInput(parameters[0], parameters[1], parameters[2]))

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("service_name", parameters[0])
	d.Set("function_name", parameters[1])
	d.Set("trigger_name", parameters[2])
	d.Set("trigger_type", getTriggerOutput.TriggerType)

	switch d.Get("trigger_type").(string) {
	case string(fc.TRIGGER_TYPE_TIMER):
		var triggerconfigTimer TriggerConfigTimer

		if err := json.Unmarshal(getTriggerOutput.RawTriggerConfig, &triggerconfigTimer); err != nil {
			return fmt.Errorf("Reading trigger of function compute got an error: %#v", err)
		} else {
			d.Set("config_payload", triggerconfigTimer.Payload)
			d.Set("config_cron_expression", triggerconfigTimer.CronExpression)
			d.Set("enable", triggerconfigTimer.Enable)
		}
	case string(fc.TRIGGER_TYPE_OSS):
		d.Set("source_arn", getTriggerOutput.SourceARN)
		d.Set("invocation_role", getTriggerOutput.InvocationRole)

		var triggerConfigOss TriggerConfigOss

		if err := json.Unmarshal(getTriggerOutput.RawTriggerConfig, &triggerConfigOss); err != nil {
			return fmt.Errorf("Reading trigger of function compute got an error: %#v", err)
		} else {
			d.Set("config_events", strings.Join(triggerConfigOss.Events, COMMA_SEPARATED))
			d.Set("config_filter_key_prefix", triggerConfigOss.Filter.Key.Prefix)
			d.Set("config_filter_key_suffix", triggerConfigOss.Filter.Key.Suffix)
		}
	case string(fc.TRIGGER_TYPE_LOG):
		d.Set("source_arn", getTriggerOutput.SourceARN)
		d.Set("invocation_role", getTriggerOutput.InvocationRole)

		var triggerConfigLog TriggerConfigLog

		if err := json.Unmarshal(getTriggerOutput.RawTriggerConfig, &triggerConfigLog); err != nil {
			return fmt.Errorf("Reading trigger of function compute got an error: %#v", err)
		} else {
			d.Set("config_source_logstore", triggerConfigLog.SourceConfig.Logstore)
			d.Set("config_job_interval", triggerConfigLog.JobConfig.TriggerInterval)
			d.Set("config_job_max_retry_time", triggerConfigLog.JobConfig.MaxRetryTime)
			d.Set("functionParameter", triggerConfigLog.FunctionParameter)
			d.Set("config_log_project", triggerConfigLog.LogConfig.Project)
			d.Set("config_log_logstore", triggerConfigLog.LogConfig.Logstore)
		}
	case string(fc.TRIGGER_TYPE_HTTP):
		var triggerConfigHttp TriggerConfigHttp

		if err := json.Unmarshal(getTriggerOutput.RawTriggerConfig, &triggerConfigHttp); err != nil {
			return fmt.Errorf("Reading trigger of function compute got an error: %#v", err)
		} else {
			d.Set("config_auth_type", triggerConfigHttp.AuthType)
			d.Set("config_methods", strings.Join(triggerConfigHttp.Methods, COMMA_SEPARATED))
		}
	default:
	}

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":                 *getTriggerOutput.TriggerName,
		"name":               *getTriggerOutput.TriggerName,
		"status":             "Available",
		"creation_time":      *getTriggerOutput.CreatedTime,
		"trigger_type":       *getTriggerOutput.TriggerType,
		"raw_rrigger_config": string(getTriggerOutput.RawTriggerConfig[:]),
		"last_modified_time": *getTriggerOutput.LastModifiedTime,
		"resource_type":      "alicloud_fc_trigger",
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_fc_trigger", s); err != nil {
		return fmt.Errorf("Setting alicloud_fc_trigger got an error: %#v.", err)
	}
	return nil
}

func resourceAlicloudFcTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parameters := strings.Split(d.Id(), COMMA_SEPARATED)

	d.Partial(true)
	update := false
	configChange := false

	updateTriggerInput := fc.NewUpdateTriggerInput(parameters[0], parameters[1], parameters[2])

	if d.HasChange("invocation_role") {
		update = true
		updateTriggerInput.WithInvocationRole(d.Get("invocation_role").(string))
		d.SetPartial("invocation_role")
	}

	triggerConfig := make(map[string]interface{})
	triggerConfig["enable"] = d.Get("config_enable")
	if d.HasChange("config_enable") {
		update = true
		configChange = true
		d.SetPartial("config_enable")
	}
	if d.HasChange("config_function_parameter") {
		update = true
		configChange = true
		triggerConfig["functionParameter"] = d.Get("config_function_parameter")
		d.SetPartial("config_function_parameter")
	}
	switch d.Get("trigger_type").(string) {
	case string(fc.TRIGGER_TYPE_TIMER):
		if configChange || d.HasChange("config_payload") || d.HasChange("config_cron_expression") {
			update = true
			configChange = true
			if d.HasChange("config_payload") {
				d.SetPartial("config_payload")
			}
			if d.HasChange("config_cron_expression") {
				d.SetPartial("config_cron_expression")
			}
			if payload, ok := d.GetOk("config_payload"); ok {
				triggerConfig["payload"] = payload
			} else {
				return fmt.Errorf("Updating trigger of function compute got an error: %#v", "Can not find the parameter config_payload")
			}
			if cronExpression, ok := d.GetOk("config_cron_expression"); ok {
				triggerConfig["cronExpression"] = cronExpression
			} else {
				return fmt.Errorf("Updating trigger of function compute got an error: %#v", "Can not find the parameter config_cron_expression")
			}
		}
	case string(fc.TRIGGER_TYPE_OSS):
		if d.HasChange("invocation_role") {
			update = true
			configChange = true
			updateTriggerInput.WithInvocationRole(d.Get("invocation_role").(string))
			d.SetPartial("invocation_role")
		}
		if configChange || d.HasChange("config_events") || d.HasChange("config_filter_key_prefix") || d.HasChange("config_filter_key_suffix") {
			update = true
			if d.HasChange("config_events") {
				d.SetPartial("config_events")
			}
			if d.HasChange("config_filter_key_prefix") {
				d.SetPartial("config_filter_key_prefix")
			}
			if d.HasChange("config_filter_key_suffix") {
				d.SetPartial("config_filter_key_suffix")
			}
			if events, ok := d.GetOk("config_events"); ok {
				triggerConfig["events"] = strings.Split(events.(string), COMMA_SEPARATED)
			} else {
				return fmt.Errorf("Updating trigger of function compute got an error: %#v", "Can not find the parameter config_events")
			}

			filter := make(map[string]interface{})
			key := make(map[string]string)
			if prefix, ok := d.GetOk("config_filter_key_prefix"); ok {
				key["prefix"] = prefix.(string)
			} else {
				return fmt.Errorf("Updating trigger of function compute got an error: %#v", "Can not find the parameter config_filter_key_prefix")
			}
			if suffix, ok := d.GetOk("config_filter_key_suffix"); ok {
				key["suffix"] = suffix.(string)
			} else {
				return fmt.Errorf("Updating trigger of function compute got an error: %#v", "Can not find the parameter config_filter_key_suffix")
			}
			filter["key"] = key
			triggerConfig["filter"] = filter
		}
	case string(fc.TRIGGER_TYPE_LOG):
		if d.HasChange("invocation_role") {
			update = true
			updateTriggerInput.WithInvocationRole(d.Get("invocation_role").(string))
			d.SetPartial("invocation_role")
		}
		if configChange || d.HasChange("config_source_logstore") || d.HasChange("config_job_interval") || d.HasChange("config_job_max_retry_time") || d.HasChange("config_log_project") || d.HasChange("config_log_logstore") {
			update = true
			configChange = true
			if d.HasChange("config_source_logstore") {
				d.SetPartial("config_source_logstore")
			}
			if d.HasChange("config_job_interval") {
				d.SetPartial("config_job_interval")
			}
			if d.HasChange("config_job_max_retry_time") {
				d.SetPartial("config_job_max_retry_time")
			}
			if d.HasChange("config_log_project") {
				d.SetPartial("config_log_project")
			}
			if d.HasChange("config_log_logstore") {
				d.SetPartial("config_log_logstore")
			}
			sourceConfig := make(map[string]string)
			jobConfig := make(map[string]int32)
			logConfig := make(map[string]string)
			if sourceLogstore, ok := d.GetOk("config_source_logstore"); ok {
				sourceConfig["logstore"] = sourceLogstore.(string)
			} else {
				return fmt.Errorf("Updating trigger of function compute got an error: %#v", "Can not find the parameter config_source_logstore")
			}
			if jobInterval, ok := d.GetOk("config_job_interval"); ok {
				if jobInterval32, err := strconv.ParseInt(strconv.Itoa(jobInterval.(int)), 10, 32); err != nil {
					return fmt.Errorf("Updating trigger of function compute got an error: %#v", err)
				} else {
					jobConfig["triggerInterval"] = int32(jobInterval32)
				}
			} else {
				return fmt.Errorf("Updating trigger of function compute got an error: %#v", "Can not find the parameter config_job_interval")
			}
			if maxRetryTime, ok := d.GetOk("config_job_max_retry_time"); ok {
				if maxRetryTime32, err := strconv.ParseInt(strconv.Itoa(maxRetryTime.(int)), 10, 32); err != nil {
					return fmt.Errorf("Updating trigger of function compute got an error: %#v", err)
				} else {
					jobConfig["maxRetryTime"] = int32(maxRetryTime32)
				}
			} else {
				return fmt.Errorf("Updating trigger of function compute got an error: %#v", "Can not find the parameter config_job_max_retry_time")
			}
			if project, ok := d.GetOk("config_log_project"); ok {
				logConfig["project"] = project.(string)
			} else {
				return fmt.Errorf("Updating trigger of function compute got an error: %#v", "Can not find the parameter config_log_project")
			}
			if logstore, ok := d.GetOk("config_log_logstore"); ok {
				logConfig["logstore"] = logstore.(string)
			} else {
				return fmt.Errorf("Updating trigger of function compute got an error: %#v", "Can not find the parameter config_log_logstore")
			}
			triggerConfig["sourceConfig"] = sourceConfig
			triggerConfig["jobConfig"] = jobConfig
			triggerConfig["logConfig"] = logConfig
		}
	case string(fc.TRIGGER_TYPE_HTTP):
		if configChange || d.HasChange("config_auth_type") || d.HasChange("config_methods") {
			update = true
			configChange = true
			if d.HasChange("config_auth_type") {
				d.SetPartial("config_auth_type")
			}
			if d.HasChange("config_methods") {
				d.SetPartial("config_methods")
			}
			if authType, ok := d.GetOk("config_auth_type"); ok {
				triggerConfig["authType"] = authType
			} else {
				return fmt.Errorf("Updating trigger of function compute got an error: %#v", "Can not find the parameter config_auth_type")
			}
			if methods, ok := d.GetOk("config_methods"); ok {
				triggerConfig["methods"] = strings.Split(methods.(string), COMMA_SEPARATED)
			} else {
				return fmt.Errorf("Updating trigger of function compute got an error: %#v", "Can not find the parameter config_methods")
			}
		}
	default:
	}

	if configChange {
		updateTriggerInput.WithTriggerConfig(triggerConfig)
	}
	if !d.IsNewResource() && update {
		if _, err := client.fcconn.UpdateTrigger(updateTriggerInput); err != nil {
			return fmt.Errorf("Updating trigger of function compute got an error: %#v", err)
		}
	}

	d.Partial(false)

	return resourceAlicloudFcTriggerRead(d, meta)
}

func resourceAlicloudFcTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		parameters := strings.Split(d.Id(), COMMA_SEPARATED)
		_, err := client.fcconn.DeleteTrigger(fc.NewDeleteTriggerInput(parameters[0], parameters[1], parameters[2]))

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting trigger of function compute got an error: %#v", err))
		}

		resp, err := client.fcconn.GetTrigger(fc.NewGetTriggerInput(parameters[0], parameters[1], parameters[2]))
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("trigger trigger of function compute got an error: %#v", err))
		}
		if resp.TriggerName == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting trigger of function compute got an error: %#v", err))
	})
}
