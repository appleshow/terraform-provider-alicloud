package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudCmsAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsAlarmCreate,
		Read:   resourceAlicloudCmsAlarmRead,
		Update: resourceAlicloudCmsAlarmUpdate,
		Delete: resourceAlicloudCmsAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metric": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dimensions": &schema.Schema{
				Type:             schema.TypeMap,
				Required:         true,
				Elem:             schema.TypeString,
				DiffSuppressFunc: cmsDimensionsDiffSuppressFunc,
			},
			"period": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"statistics": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  Average,
				ValidateFunc: validateAllowedStringValue([]string{
					string(Average), string(Minimum), string(Maximum),
				}),
			},
			"operator": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  Equal,
				ValidateFunc: validateAllowedStringValue([]string{
					MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, Equal, NotEqual,
				}),
			},
			"threshold": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"triggered_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"contact_groups": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"start_time": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateIntegerInRange(0, 24),
			},
			"end_time": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      24,
				ValidateFunc: validateIntegerInRange(0, 24),
			},
			"silence_time": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      86400,
				ValidateFunc: validateIntegerInRange(300, 86400),
			},
			"notify_type": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"webhook": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_cms_alarm": {
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
						"project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metrie_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dimensions": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"statistics": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"threshold": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"evaluation_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"silence_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"notify_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_groups": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"webhook": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
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

func resourceAlicloudCmsAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := cms.CreateCreateAlarmRequest()

	request.Name = d.Get("name").(string)
	request.Namespace = d.Get("project").(string)
	request.MetricName = d.Get("metric").(string)
	request.Period = requests.NewInteger(d.Get("period").(int))
	request.Statistics = d.Get("statistics").(string)
	request.ComparisonOperator = d.Get("operator").(string)
	request.Threshold = d.Get("threshold").(string)
	request.EvaluationCount = requests.NewInteger(d.Get("triggered_count").(int))
	request.ContactGroups = convertListToJsonString(d.Get("contact_groups").([]interface{}))
	request.StartTime = requests.NewInteger(d.Get("start_time").(int))
	request.EndTime = requests.NewInteger(d.Get("end_time").(int))
	request.SilenceTime = requests.NewInteger(d.Get("silence_time").(int))
	if webhook, ok := d.GetOk("webhook"); ok && "" != webhook {
		request.Webhook = webhook.(string)
	}
	if v, ok := d.GetOk("notify_type"); ok {
		request.NotifyType = requests.NewInteger(v.(int))
	}

	var dimList []map[string]string
	if dimensions, ok := d.GetOk("dimensions"); ok {
		for k, v := range dimensions.(map[string]interface{}) {
			values := strings.Split(v.(string), COMMA_SEPARATED)
			if len(values) > 0 {
				for _, vv := range values {
					dimList = append(dimList, map[string]string{k: Trim(vv)})
				}
			} else {
				dimList = append(dimList, map[string]string{k: Trim(v.(string))})
			}

		}
	}
	if len(dimList) > 0 {
		if bytes, err := json.Marshal(dimList); err != nil {
			return fmt.Errorf("Marshaling dimensions to json string got an error: %#v.", err)
		} else {
			request.Dimensions = string(bytes[:])
		}
	} else {
		request.Dimensions = string("[{}]")
	}
	response, err := client.cmsconn.CreateAlarm(request)
	if err != nil {
		return fmt.Errorf("Creating alarm got an error: %#v", err)
	}

	d.SetId(response.Data)

	return resourceAlicloudCmsAlarmUpdate(d, meta)
}

func resourceAlicloudCmsAlarmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	alarm, err := client.DescribeAlarm(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", alarm.Name)
	d.Set("project", alarm.Namespace)
	d.Set("metric", alarm.MetricName)
	d.Set("period", alarm.Period)
	d.Set("statistics", alarm.Statistics)
	d.Set("operator", alarm.ComparisonOperator)
	d.Set("threshold", alarm.Threshold)
	d.Set("triggered_count", alarm.EvaluationCount)
	d.Set("start_time", alarm.StartTime)
	d.Set("end_time", alarm.EndTime)
	d.Set("silence_time", alarm.SilenceTime)
	d.Set("notify_type", alarm.NotifyType)
	d.Set("enabled", alarm.Enable)
	if webhook, ok := d.GetOk("webhook"); ok && "" != webhook {
		d.Set("webhook", alarm.Webhook)
	}

	var groups []string
	if err := json.Unmarshal([]byte(alarm.ContactGroups), &groups); err != nil {
		return fmt.Errorf("Unmarshaling contact groups got an error: %#v.", err)
	} else {
		d.Set("contact_groups", groups)
	}

	var dims []string
	if err := json.Unmarshal([]byte(alarm.Dimensions), &dims); err != nil {
		return fmt.Errorf("Unmarshaling Dimensions got an error: %#v.", err)
	}
	d.Set("dimensions", dims)

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":                  alarm.Id,
		"name":                alarm.Name,
		"state":               alarm.State,
		"creation_time":       time.Now().Format("2006-01-02 15:04:05"),
		"project":             alarm.Namespace,
		"metrie_name":         alarm.MetricName,
		"dimensions":          alarm.Dimensions,
		"period":              alarm.Period,
		"statistics":          alarm.Statistics,
		"comparison_operator": alarm.ComparisonOperator,
		"threshold":           alarm.Threshold,
		"evaluation_count":    alarm.EvaluationCount,
		"start_time":          alarm.StartTime,
		"end_time":            alarm.EndTime,
		"silence_time":        alarm.SilenceTime,
		"notify_type":         alarm.NotifyType,
		"enable":              alarm.Enable,
		"contact_groups":      alarm.ContactGroups,
		"webhook":             alarm.Webhook,
		"resource_type":       "alicloud_cms_alarm",
	}

	s = append(s, mapping)
	if err := d.Set("alicloud_cms_alarm", s); err != nil {
		return fmt.Errorf("Setting alicloud_cms_alarm got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudCmsAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	d.Partial(true)
	update := false

	request := cms.CreateUpdateAlarmRequest()
	request.Id = d.Id()

	if d.HasChange("project") && !d.IsNewResource() {
		return fmt.Errorf("Updating alarm got an error: %#v", "Cannot modify parameter project")
	}
	if d.HasChange("dimensions") && !d.IsNewResource() {
		return fmt.Errorf("Updating alarm got an error: %#v", "Cannot modify parameter dimensions")
	}
	if d.HasChange("Name") {
		update = true
		request.Name = d.Get("name").(string)
		d.SetPartial("name")
	}
	if d.HasChange("period") {
		update = true
		request.Period = requests.NewInteger(d.Get("period").(int))
		d.SetPartial("period")
	}
	if d.HasChange("statistics") {
		update = true
		request.Statistics = d.Get("statistics").(string)
		d.SetPartial("statistics")
	}
	if d.HasChange("operator") {
		update = true
		request.ComparisonOperator = d.Get("operator").(string)
		d.SetPartial("operator")
	}
	if d.HasChange("threshold") {
		update = true
		request.Threshold = d.Get("threshold").(string)
		d.SetPartial("threshold")
	}
	if d.HasChange("triggered_count") {
		update = true
		request.EvaluationCount = requests.NewInteger(d.Get("triggered_count").(int))
		d.SetPartial("triggered_count")
	}
	if d.HasChange("contact_groups") {
		update = true
		request.ContactGroups = convertListToJsonString(d.Get("contact_groups").([]interface{}))
		d.SetPartial("contact_groups")
	}
	if d.HasChange("start_time") {
		update = true
		request.StartTime = requests.NewInteger(d.Get("start_time").(int))
		d.SetPartial("start_time")
	}
	if d.HasChange("end_time") {
		update = true
		request.EndTime = requests.NewInteger(d.Get("end_time").(int))
		d.SetPartial("end_time")
	}
	if d.HasChange("silence_time") {
		update = true
		request.SilenceTime = requests.NewInteger(d.Get("silence_time").(int))
		d.SetPartial("silence_time")
	}
	if d.HasChange("notify_type") {
		update = true
		request.NotifyType = requests.NewInteger(d.Get("notify_type").(int))
		d.SetPartial("notify_type")
	}

	if d.HasChange("webhook") {
		update = true
		request.Webhook = d.Get("webhook").(string)
		d.SetPartial("webhook")
	}

	if !d.IsNewResource() && update {
		if _, err := client.cmsconn.UpdateAlarm(request); err != nil {
			return fmt.Errorf("Updating alarm got an error: %#v", err)
		}
	}

	if d.HasChange("enabled") {
		if d.Get("enabled").(bool) {
			request := cms.CreateEnableAlarmRequest()
			request.Id = d.Id()

			if _, err := client.cmsconn.EnableAlarm(request); err != nil {
				return fmt.Errorf("Enabling alarm got an error: %#v", err)
			}
		} else {
			request := cms.CreateDisableAlarmRequest()
			request.Id = d.Id()

			if _, err := client.cmsconn.DisableAlarm(request); err != nil {
				return fmt.Errorf("Disableing alarm got an error: %#v", err)
			}
		}
		if err := client.WaitForCmsAlarm(d.Id(), d.Get("enabled").(bool), 102); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAlicloudCmsAlarmRead(d, meta)
}

func resourceAlicloudCmsAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := cms.CreateDeleteAlarmRequest()

	request.Id = d.Id()

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.cmsconn.DeleteAlarm(request)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting alarm rule got an error: %#v", err))
		}

		resp, err := client.DescribeAlarm(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe alarm rule got an error: %#v", err))
		}
		if resp.Id == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting alarm rule got an error: %#v", err))
	})
}
