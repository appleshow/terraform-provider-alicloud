package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudAutoSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAutoSnapshotPolicyCreate,
		Read:   resourceAlicloudAutoSnapshotPolicyRead,
		Update: resourceAlicloudAutoSnapshotPolicyUpdate,
		Delete: resourceAlicloudAutoSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"time_points": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"repeat_weekdays": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"retention_days": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCommonResourceName,
			},

			"disk_nums": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			// TODO: OperationLocks
			// TODO: validate function

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"creation_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			//Computed value
			"alicloud_auto_snapshot_policy": &schema.Schema{
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

func resourceAlicloudAutoSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	time_points := d.Get("time_points").(string)
	repeat_weekdays := d.Get("repeat_weekdays").(string)
	retention_days := d.Get("retention_days").(int)

	args := &ecs.CreateAutoSnapshotPolocyArgs{
		RegionId:       getRegion(d, meta),
		TimePoints:     time_points,
		RepeatWeekdays: repeat_weekdays,
		RetentionDays:  requests.NewInteger(retention_days),
	}
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		args.AutoSnapshotPolicyName = v.(string)
	}

	autoSnapshotPolicyId, err := conn.CreateAutoSnapshotPolicy(args)
	if err != nil {
		log.Printf("[DEBUG] CreateAutoSnapshotPolicy got error: %#v", err)
		return fmt.Errorf("CreateAutoSnapshotPolicy got error: %#v", err)
	}

	d.SetId(autoSnapshotPolicyId)

	return resourceAlicloudAutoSnapshotPolicyUpdate(d, meta)
}

func resourceAlicloudAutoSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	auto_snapshot_policies, _, err := conn.DescribeAutoSnapshotPolicyEx(&ecs.DescribeAutoSnapshotPolocyExArgs{
		RegionId:             getRegion(d, meta),
		AutoSnapshotPolicyId: d.Id(),
	})

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe AutoSnapshotPolicy Attribute: %#v", err)
	}

	log.Printf("[DEBUG] DescribeAutoSnapshotPolicyEx for instance: %#v", auto_snapshot_policies)

	if auto_snapshot_policies == nil || len(auto_snapshot_policies) <= 0 {
		return fmt.Errorf("No auto_snapshot_policies found.")
	}

	auto_snapshot_policy := auto_snapshot_policies[0]
	d.Set("time_points", auto_snapshot_policy.TimePoints)
	d.Set("name", auto_snapshot_policy.AutoSnapshotPolicyName)
	d.Set("repeat_weekdays", auto_snapshot_policy.RepeatWeekdays)
	d.Set("retention_days", auto_snapshot_policy.RetentionDays)
	d.Set("disk_nums", auto_snapshot_policy.DiskNums)
	d.Set("status", auto_snapshot_policy.Status)
	d.Set("creation_time", auto_snapshot_policy.CreationTime)

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":              auto_snapshot_policy.AutoSnapshotPolicyId,
		"time_points":     auto_snapshot_policy.TimePoints,
		"name":            auto_snapshot_policy.AutoSnapshotPolicyName,
		"repeat_weekdays": auto_snapshot_policy.RepeatWeekdays,
		"retention_days":  auto_snapshot_policy.RetentionDays,
		"disk_nums":       auto_snapshot_policy.DiskNums,
		"status":          auto_snapshot_policy.Status,
		"creation_time":   auto_snapshot_policy.CreationTime,
		"resource_type":   "alicloud_auto_snapshot_policy",
	}
	log.Printf("[DEBUG] auto_snapshot_policies - adding auto_snapshot_policy: %v", mapping)
	s = append(s, mapping)

	if err := d.Set("alicloud_auto_snapshot_policy", s); err != nil {
		return fmt.Errorf("Setting alicloud_auto_snapshot_policy got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudAutoSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn
	d.Partial(true)

	attributeUpdate := false
	args := &ecs.ModifyAutoSnapshotPolicyExArgs{
		RegionId:             getRegion(d, meta),
		AutoSnapshotPolicyId: d.Id(),
	}

	if d.HasChange("name") {
		d.SetPartial("name")
		val := d.Get("name").(string)
		args.AutoSnapshotPolicyName = val

		attributeUpdate = true
	}

	if d.HasChange("time_points") {
		d.SetPartial("time_points")
		val := d.Get("time_points").(string)
		args.TimePoints = val

		attributeUpdate = true
	}

	if d.HasChange("repeat_weekdays") {
		d.SetPartial("repeat_weekdays")
		val := d.Get("repeat_weekdays").(string)
		args.RepeatWeekdays = val

		attributeUpdate = true
	}

	if d.HasChange("retention_days") {
		d.SetPartial("retention_days")
		val := d.Get("retention_days").(int)
		args.RetentionDays = requests.NewInteger(val)

		attributeUpdate = true
	}

	if attributeUpdate {
		if err := conn.ModifyAutoSnapshotPolicyEx(args); err != nil {
			return err
		}
	}

	d.Partial(false)
	return resourceAlicloudAutoSnapshotPolicyRead(d, meta)
}

func resourceAlicloudAutoSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.DeleteAutoSnapshotPolicy(&ecs.DeleteAutoSnapshotPolicyArgs{
			RegionId:             getRegion(d, meta),
			AutoSnapshotPolicyId: d.Id(),
		})
		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code != ParameterInvalid {
				return resource.RetryableError(fmt.Errorf("Delete AutoSnapshotPolicy got an error:%#v.", err))
			}
		}

		auto_snapshot_policies, _, descErr := conn.DescribeAutoSnapshotPolicyEx(&ecs.DescribeAutoSnapshotPolocyExArgs{
			RegionId:             getRegion(d, meta),
			AutoSnapshotPolicyId: d.Id(),
		})

		if descErr != nil {
			log.Printf("[ERROR] Delete AutoSnapshotPolicy is failed.")
			return resource.NonRetryableError(descErr)
		}
		if auto_snapshot_policies == nil || len(auto_snapshot_policies) < 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("AutoSnapshotPolicy delete - trying again."))
	})
}
