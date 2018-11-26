package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudCmsContactGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsContactGroupsRead,

		Schema: map[string]*schema.Schema{
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
			"request_for_existence": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_contact_groups": {
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

func dataSourceAlicloudCmsContactGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := cms.CreateListContactGroupRequest()

	if pageNumber, ok := d.GetOk("page_number"); ok {
		request.PageNumber = requests.NewInteger(pageNumber.(int))
	}
	if pageSize, ok := d.GetOk("page_size"); ok {
		request.PageSize = requests.NewInteger(pageSize.(int))
	}

	listContactGroupResponse, err := client.cmsconn.ListContactGroup(request)
	log.Printf("[DEBUG] alicloud_cms_contact_group - contact groups found: %#v", listContactGroupResponse)
	if err != nil {
		return fmt.Errorf("List contact groups got an error: %#v", err)
	}
	if listContactGroupResponse != nil && !listContactGroupResponse.Success {
		return fmt.Errorf("List contact groups got an error: %#v", listContactGroupResponse.Message)
	}

	if requestForExistence, ok := d.GetOk("request_for_existence"); ok {
		requestForExistence := requestForExistence.(string)
		check := false
		for _, contacntGroup := range listContactGroupResponse.ContactGroups.ContactGroup {
			if contacntGroup == requestForExistence {
				check = true
			}
		}
		if !check {
			return fmt.Errorf("List contact groups got an error: %#v", "The contact group ["+requestForExistence+"] does not exist.")
		}
	}
	var ids []string
	var s []map[string]interface{}
	for _, contacntGroup := range listContactGroupResponse.ContactGroups.ContactGroup {
		mapping := map[string]interface{}{
			"id":            contacntGroup,
			"name":          contacntGroup,
			"status":        "Available",
			"creation_time": "",
			"resource_type": "alicloud_contact_group",
		}

		ids = append(ids, contacntGroup)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alicloud_contact_groups", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
