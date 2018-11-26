package alicloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudCmsAppGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsAppGroupsRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"bind_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"select_contact_groups": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"keyword": &schema.Schema{
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
			"alicloud_cms_app_groups": {
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

func dataSourceAlicloudCmsAppGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	listMyGroupsRequest := cms.CreateListMyGroupsRequest()

	if instanceId, ok := d.GetOk("instance_id"); ok {
		listMyGroupsRequest.InstanceId = instanceId.(string)
	}
	if groupType, ok := d.GetOk("type"); ok {
		listMyGroupsRequest.Type = groupType.(string)
	}
	if groupName, ok := d.GetOk("group_name"); ok {
		listMyGroupsRequest.GroupName = groupName.(string)
	}
	if bindUrl, ok := d.GetOk("bind_url"); ok {
		listMyGroupsRequest.BindUrls = bindUrl.(string)
	}
	if selectContactGroups, ok := d.GetOk("select_contact_groups"); ok {
		listMyGroupsRequest.SelectContactGroups = requests.NewBoolean(selectContactGroups.(bool))
	}
	if keyword, ok := d.GetOk("keyword"); ok {
		listMyGroupsRequest.Keyword = keyword.(string)
	}
	if pageNumber, ok := d.GetOk("page_number"); ok {
		listMyGroupsRequest.PageNumber = requests.NewInteger(pageNumber.(int))
	}
	if pageSize, ok := d.GetOk("page_size"); ok {
		listMyGroupsRequest.PageSize = requests.NewInteger(pageSize.(int))
	}

	listMyGroupsResponse, err := client.cmsconn.ListMyGroups(listMyGroupsRequest)
	log.Printf("[DEBUG] alicloud_cms - applicatoin groups found: %#v", listMyGroupsResponse)
	if err != nil {
		return fmt.Errorf("List application groups got an error: %#v", err)
	}
	if listMyGroupsResponse != nil && !listMyGroupsResponse.Success {
		return fmt.Errorf("List application groups got an error: %#v", listMyGroupsResponse.ErrorMessage)
	}

	var ids []string
	var s []map[string]interface{}
	for _, group := range listMyGroupsResponse.Resources.Resource {
		if id, ok := d.GetOk("id"); ok {
			if id.(int) != group.GroupId {
				continue
			}
		}

		var contactGroups []map[string]interface{}
		for _, contactGroup := range group.ContactGroups.ContactGroup {
			mapping := map[string]interface{}{
				"name": contactGroup.Name,
			}
			contactGroups = append(contactGroups, mapping)
		}

		mapping := map[string]interface{}{
			"id":             strconv.Itoa(group.GroupId),
			"name":           group.GroupName,
			"status":         "Available",
			"creation_time":  "",
			"group_id":       strconv.Itoa(group.GroupId),
			"group_name":     group.GroupName,
			"contact_groups": contactGroups,
			"resource_type":  "alicloud_cms_app_group",
		}

		ids = append(ids, strconv.Itoa(group.GroupId))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alicloud_cms_app_groups", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
