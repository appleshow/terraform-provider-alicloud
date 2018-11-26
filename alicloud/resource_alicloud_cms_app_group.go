package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudCmsAppGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsAppGroupCreate,
		Read:   resourceAlicloudCmsAppGroupRead,
		Update: resourceAlicloudCmsAppGroupUpdate,
		Delete: resourceAlicloudCmsAppGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"bind_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_groups": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"options": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_cms_app_group": {
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
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
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

func resourceAlicloudCmsAppGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	createMyGroupsRequest := cms.CreateCreateMyGroupsRequest()

	createMyGroupsRequest.GroupName = d.Get("group_name").(string)
	if groupType, ok := d.GetOk("type"); ok {
		createMyGroupsRequest.Type = groupType.(string)
	}
	if serviceId, ok := d.GetOk("service_id"); ok {
		createMyGroupsRequest.ServiceId = requests.NewInteger(serviceId.(int))
	}
	if bindUrl, ok := d.GetOk("bind_url"); ok {
		createMyGroupsRequest.BindUrl = bindUrl.(string)
	}
	if contactGroups, ok := d.GetOk("contact_groups"); ok {
		createMyGroupsRequest.ContactGroups = contactGroups.(string)
	}
	if options, ok := d.GetOk("options"); ok {
		createMyGroupsRequest.Options = options.(string)
	}

	createMyGroupsResponse, err := client.cmsconn.CreateMyGroups(createMyGroupsRequest)
	if err != nil {
		return fmt.Errorf("Creating application group got an error: %#v", err)
	}

	d.SetId(strconv.Itoa(createMyGroupsResponse.GroupId))

	return resourceAlicloudCmsAppGroupUpdate(d, meta)
}

func resourceAlicloudCmsAppGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	getMyGroupsRequest := cms.CreateGetMyGroupsRequest()
	groupId, err := strconv.Atoi(d.Id())
	if err == nil {
		getMyGroupsRequest.GroupId = requests.NewInteger(groupId)
	} else {
		return fmt.Errorf("Reading application group got an error: %#v", err)
	}

	getMyGroupsResponse, err := client.cmsconn.GetMyGroups(getMyGroupsRequest)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("group_name", getMyGroupsResponse.Group.GroupName)
	if _, ok := d.GetOk("type"); ok {
		d.Set("type", getMyGroupsResponse.Group.Type)
	}
	if _, ok := d.GetOk("service_id"); ok {
		d.Set("service_id", getMyGroupsResponse.Group.ServiceId)
	}
	if _, ok := d.GetOk("bind_url"); ok {
		d.Set("bind_url", getMyGroupsResponse.Group.BindUrl)
	}
	if _, ok := d.GetOk("contact_groups"); ok {
		var contactGroupStr = ""
		for _, contactGroup := range getMyGroupsResponse.Group.ContactGroups.ContactGroup {
			if contactGroupStr == "" {
				contactGroupStr = contactGroup.Name
			} else {
				contactGroupStr = contactGroupStr + "," + contactGroup.Name
			}
		}

		d.Set("contact_groups", contactGroupStr)
	}

	var s []map[string]interface{}
	var contactGroups []map[string]interface{}
	for _, contactGroup := range getMyGroupsResponse.Group.ContactGroups.ContactGroup {
		mapping := map[string]interface{}{
			"name": contactGroup.Name,
		}
		contactGroups = append(contactGroups, mapping)
	}
	mapping := map[string]interface{}{
		"id":             strconv.Itoa(getMyGroupsResponse.Group.GroupId),
		"name":           getMyGroupsResponse.Group.GroupName,
		"status":         "Available",
		"creation_time":  time.Now().Format("2006-01-02 15:04:05"),
		"group_id":       strconv.Itoa(getMyGroupsResponse.Group.GroupId),
		"group_name":     getMyGroupsResponse.Group.GroupName,
		"contact_groups": contactGroups,
		"resource_type":  "alicloud_cms_app_group",
	}
	s = append(s, mapping)
	if err := d.Set("alicloud_cms_app_group", s); err != nil {
		return fmt.Errorf("Setting alicloud_cms_app_group got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudCmsAppGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	d.Partial(true)
	update := false

	updateMyGroupsRequest := cms.CreateUpdateMyGroupsRequest()
	updateMyGroupsRequest.GroupId = d.Id()

	if d.HasChange("group_name") {
		update = true
		updateMyGroupsRequest.GroupName = d.Get("group_name").(string)
		d.SetPartial("group_name")
	}
	if d.HasChange("service_id") {
		update = true
		updateMyGroupsRequest.ServiceId = requests.NewInteger(d.Get("service_id").(int))
		d.SetPartial("service_id")
	}
	if d.HasChange("bind_url") {
		update = true
		updateMyGroupsRequest.BindUrls = d.Get("bind_url").(string)
		d.SetPartial("bind_url")
	}
	if d.HasChange("contact_groups") {
		update = true
		updateMyGroupsRequest.ContactGroups = d.Get("contact_groups").(string)
		d.SetPartial("contact_groups")
	}

	if !d.IsNewResource() && update {
		if _, err := client.cmsconn.UpdateMyGroups(updateMyGroupsRequest); err != nil {
			return fmt.Errorf("Updating application group got an error: %#v", err)
		}
	}

	d.Partial(false)

	return resourceAlicloudCmsAppGroupRead(d, meta)
}

func resourceAlicloudCmsAppGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	deleteMyGroupsRequest := cms.CreateDeleteMyGroupsRequest()

	groupId, err := strconv.Atoi(d.Id())
	if err == nil {
		deleteMyGroupsRequest.GroupId = requests.NewInteger(groupId)
	} else {
		return fmt.Errorf("Deleting application group got an error: %#v", err)
	}

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.cmsconn.DeleteMyGroups(deleteMyGroupsRequest)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting application group got an error: %#v", err))
		}

		getMyGroupsRequest := cms.CreateGetMyGroupsRequest()
		groupId, err := strconv.Atoi(d.Id())
		if err == nil {
			getMyGroupsRequest.GroupId = requests.NewInteger(groupId)
		} else {
			return resource.NonRetryableError(fmt.Errorf("Deleting application group got an error: %#v", err))
		}

		_, errGet := client.cmsconn.GetMyGroups(getMyGroupsRequest)
		if errGet != nil {
			if NotFoundError(errGet) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe application group got an error: %#v", errGet))
		}

		return resource.RetryableError(fmt.Errorf("Deleting application group got an error: %#v", errGet))
	})
}
