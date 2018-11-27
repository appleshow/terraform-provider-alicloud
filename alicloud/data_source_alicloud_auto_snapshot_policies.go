package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudAutoSnapshotPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAutoSnapshotPoliciesRead,

		Schema: map[string]*schema.Schema{
			// TODO: filter
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},
			//Computed value
			"alicloud_auto_snapshot_policies": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_points": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repeat_weekdays": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_nums": {
							Type:     schema.TypeInt,
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

func dataSourceAlicloudAutoSnapshotPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	var arg = &ecs.DescribeAutoSnapshotPolocyExArgs{
		RegionId: getRegion(d, meta),
	}
	if autoSnapshotPolicyId, ok := d.GetOk("id"); ok {
		arg.AutoSnapshotPolicyId = autoSnapshotPolicyId.(string)
	}

	auto_snapshot_policies, _, err := conn.DescribeAutoSnapshotPolicyEx(arg)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe AutoSnapshotPolicy Attribute: %#v", err)
	}

	log.Printf("[DEBUG] DescribeAutoSnapshotPolicyEx for instance: %#v", auto_snapshot_policies)

	/*
		if len(auto_snapshot_policies) < 1 {
			return fmt.Errorf("Your query returned no results[auto_snapshot_policies]. Please change your search criteria and try again.")
		}
	*/

	log.Printf("[DEBUG] auto_snapshot_policies found: %#v", auto_snapshot_policies)

	var ids []string
	var s []map[string]interface{}
	var nameRegex *regexp.Regexp
	var exist bool

	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	exist = false
	for _, auto_snapshot_policy := range auto_snapshot_policies {
		if nameRegex != nil {
			if !nameRegex.MatchString(auto_snapshot_policy.AutoSnapshotPolicyName) {
				continue
			}
		}
		exist = true
		mapping := map[string]interface{}{
			"id":              auto_snapshot_policy.AutoSnapshotPolicyId,
			"time_points":     auto_snapshot_policy.TimePoints,
			"name":            auto_snapshot_policy.AutoSnapshotPolicyName,
			"repeat_weekdays": auto_snapshot_policy.RepeatWeekdays,
			"retention_days":  auto_snapshot_policy.RetentionDays,
			"disk_nums":       auto_snapshot_policy.DiskNums,
			"status":          auto_snapshot_policy.Status,
			"creation_time":   auto_snapshot_policy.CreationTime,
			"resource_type":   "alcloud_auto_snapshot_policy",
		}
		log.Printf("[DEBUG] auto_snapshot_policies - adding auto_snapshot_policy: %v", mapping)
		ids = append(ids, auto_snapshot_policy.AutoSnapshotPolicyId)
		s = append(s, mapping)
	}

	if !exist {
		return fmt.Errorf("Your query returned no results[auto_snapshot_policies]. Please change your search criteria and try again.")
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alicloud_auto_snapshot_policies", s); err != nil {
		return err
	}
	return nil
}
