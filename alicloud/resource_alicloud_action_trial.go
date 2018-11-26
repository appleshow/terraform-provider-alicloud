package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudActionTrial() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudActionTrialCreate,
		Read:   resourceAlicloudActionTrialRead,
		Update: resourceAlicloudActionTrialUpdate,
		Delete: resourceAlicloudActionTrialDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"oss_bucket_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"role_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"oss_key_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_logging": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			// Computed values
			"alicloud_action_trial": {
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
						"statistics": {
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
						"oss_bucket_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_bucket_location": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_key_prefix": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudActionTrialCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := CreateCreateTrailRequest()

	request.Name = d.Get("name").(string)
	request.OssBucketName = d.Get("oss_bucket_name").(string)
	request.RoleName = d.Get("role_name").(string)
	if ossKeyPrefix, ok := d.GetOk("oss_key_prefix"); ok {
		request.OssKeyPrefix = ossKeyPrefix.(string)
	}

	_, err := CreateTrail(client.cmsconn, request)
	if err != nil {
		return fmt.Errorf("Creating action trial got an error: %#v", err)
	}

	d.SetId(request.Name)

	if d.Get("is_logging").(bool) {
		request := CreateStartLoggingRequest()
		request.Name = d.Id()

		_, err := StartLogging(client.cmsconn, request)
		if err != nil {
			return fmt.Errorf("Starting action trial got an error: %#v", err)
		}
	}

	return resourceAlicloudActionTrialUpdate(d, meta)
}

func resourceAlicloudActionTrialRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := CreateDescribeTrailsRequest()
	request.NameList = d.Id()

	response, err := DescribeTrails(client.cmsconn, request)
	if err != nil {
		return fmt.Errorf("Describing action trials got an error: %#v", err)
	}

	if response == nil || response.TrailList == nil || len(response.TrailList) == 0 {
		d.SetId("")
		return nil
	}

	actionTrial := response.TrailList[0]
	d.Set("name", actionTrial.Name)
	d.Set("oss_bucket_name", actionTrial.OssBucketName)
	d.Set("role_name", actionTrial.RoleName)
	if _, ok := d.GetOk("oss_key_prefix"); ok {
		d.Set("oss_key_prefix", actionTrial.OssKeyPrefix)
	}

	var s []map[string]interface{}
	var statistics string

	if d.Get("is_logging").(bool) {
		statistics = "Available"
	} else {
		statistics = "Disable"
	}
	mapping := map[string]interface{}{
		"id":                  actionTrial.Name,
		"name":                actionTrial.Name,
		"statistics":          statistics,
		"creation_time":       time.Now().Format("2006-01-02 15:04:05"),
		"resource_type":       "alicloud_action_trial",
		"oss_bucket_name":     actionTrial.OssBucketName,
		"oss_bucket_location": actionTrial.OssBucketLocation,
		"role_name":           actionTrial.RoleName,
		"oss_key_prefix":      actionTrial.OssKeyPrefix,
	}

	s = append(s, mapping)
	if err := d.Set("alicloud_action_trial", s); err != nil {
		return fmt.Errorf("Setting alicloud_action_trial got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudActionTrialUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	d.Partial(true)
	update := false

	request := CreateUpdateTrailRequest()
	request.Name = d.Id()

	if d.HasChange("name") && !d.IsNewResource() {
		return fmt.Errorf("Updating action trial got an error: %#v", "Cannot modify parameter name")
	}

	if d.HasChange("oss_bucket_name") {
		update = true
		request.OssBucketName = d.Get("oss_bucket_name").(string)
		d.SetPartial("oss_bucket_name")
	}

	if d.HasChange("role_name") {
		update = true
		request.RoleName = d.Get("role_name").(string)
		d.SetPartial("role_name")
	}

	if d.HasChange("oss_key_prefix") {
		update = true
		request.OssKeyPrefix = d.Get("oss_key_prefix").(string)
		d.SetPartial("oss_key_prefix")
	}

	if !d.IsNewResource() && update {
		if _, err := UpdateTrail(client.cmsconn, request); err != nil {
			return fmt.Errorf("Updating action trial got an error: %#v", err)
		}
		if d.HasChange("is_logging") {
			if d.Get("is_logging").(bool) {
				request := CreateStartLoggingRequest()
				request.Name = d.Id()

				_, err := StartLogging(client.cmsconn, request)
				if err != nil {
					return fmt.Errorf("Starting action trial got an error: %#v", err)
				}
			} else {
				request := CreateStopLoggingRequest()
				request.Name = d.Id()

				_, err := StopLogging(client.cmsconn, request)
				if err != nil {
					return fmt.Errorf("Stopping action trial got an error: %#v", err)
				}
			}
		}
	}

	d.Partial(false)

	return resourceAlicloudActionTrialRead(d, meta)
}

func resourceAlicloudActionTrialDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := CreateDeleteTrailRequest()

	request.Name = d.Id()

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := DeleteTrail(client.cmsconn, request)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting action trial got an error: %#v", err))
		}

		return nil
	})
}
