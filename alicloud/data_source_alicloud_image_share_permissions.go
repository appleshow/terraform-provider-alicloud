package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudImageSharePermissions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudImageSharePermissionRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"page_number": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"page_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  50,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_image_share_permissions": {
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
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"page_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"page_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"share_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"accounts": {
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

func dataSourceAlicloudImageSharePermissionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	describeImageSharePermissionRequest := ecs.CreateDescribeImageSharePermissionRequest()
	describeImageSharePermissionRequest.ImageId = d.Get("id").(string)
	describeImageSharePermissionRequest.PageNumber = requests.NewInteger(d.Get("page_number").(int))
	describeImageSharePermissionRequest.PageSize = requests.NewInteger(d.Get("page_size").(int))

	describeImageSharePermissionResponse, err := client.aliecsconn.DescribeImageSharePermission(describeImageSharePermissionRequest)

	if err != nil {
		return fmt.Errorf("List Image Share Permission got an error: %#v", err)
	}

	var ids []string
	var s []map[string]interface{}
	var shareGroups []interface{}
	var accounts []interface{}

	for _, group := range describeImageSharePermissionResponse.ShareGroups.ShareGroup {
		shareGroups = append(shareGroups, group.Group)
	}
	for _, accout := range describeImageSharePermissionResponse.Accounts.Account {
		accounts = append(accounts, accout.AliyunId)
	}

	mapping := map[string]interface{}{
		"id":            describeImageSharePermissionResponse.ImageId,
		"name":          describeImageSharePermissionResponse.ImageId,
		"status":        "Available",
		"creation_time": "",
		"region_id":     describeImageSharePermissionResponse.RegionId,
		"total_count":   describeImageSharePermissionResponse.TotalCount,
		"page_number":   describeImageSharePermissionResponse.PageNumber,
		"page_size":     describeImageSharePermissionResponse.PageSize,
		"share_groups":  shareGroups,
		"accounts":      accounts,
		"resource_type": "alicloud_image_share_permission",
	}

	ids = append(ids, describeImageSharePermissionResponse.ImageId)
	s = append(s, mapping)

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alicloud_image_share_permissions", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
