package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudCmsAlarms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsAlarmsRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"dimension": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
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
				Default:  100,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_cms_alarms": {
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

func dataSourceAlicloudCmsAlarmsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := cms.CreateListAlarmRequest()

	if id, ok := d.GetOk("id"); ok {
		request.Id = id.(string)
	}
	if name, ok := d.GetOk("name"); ok {
		request.Name = name.(string)
	}
	if nameSpace, ok := d.GetOk("project"); ok {
		request.Namespace = nameSpace.(string)
	}
	if dimension, ok := d.GetOk("dimension"); ok {
		request.Dimension = dimension.(string)
	}
	if isEnable, ok := d.GetOk("is_enable"); ok {
		request.IsEnable = requests.NewBoolean(isEnable.(bool))
	}
	if state, ok := d.GetOk("state"); ok {
		request.State = state.(string)
	}
	if pageNumber, ok := d.GetOk("page_number"); ok {
		request.PageNumber = requests.NewInteger(pageNumber.(int))
	}
	if pageSize, ok := d.GetOk("page_size"); ok {
		request.PageSize = requests.NewInteger(pageSize.(int))
	}

	listAlarmResponse, err := client.cmsconn.ListAlarm(request)
	log.Printf("[DEBUG] alicloud_cms_alarm - alarms found: %#v", listAlarmResponse)
	if err != nil {
		return fmt.Errorf("List alarms got an error: %#v", err)
	}
	if listAlarmResponse != nil && !listAlarmResponse.Success {
		return fmt.Errorf("List alarms got an error: %#v", listAlarmResponse.Message)
	}

	var ids []string
	var s []map[string]interface{}
	for _, alarm := range listAlarmResponse.AlarmList.Alarm {
		mapping := map[string]interface{}{
			"id":                  alarm.Id,
			"name":                alarm.Name,
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
			"state":               alarm.State,
			"contact_groups":      alarm.ContactGroups,
			"webhook":             alarm.Webhook,
			"creation_time":       "",
			"resource_type":       "alicloud_cms_alarm",
		}

		ids = append(ids, alarm.Id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alicloud_cms_alarms", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
