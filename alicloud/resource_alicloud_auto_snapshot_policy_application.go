package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudAutoSnapshotPolicyApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAutoSnapshotPolicyApplicationCreate,
		Read:   resourceAlicloudAutoSnapshotPolicyApplicationRead,
		Delete: resourceAlicloudAutoSnapshotPolicyApplicationDelete,

		Schema: map[string]*schema.Schema{
			"auto_snapshot_policy_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"disk_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			//Computed value
			"alicloud_auto_snapshot_policy_application": &schema.Schema{
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
						"auto_snapshot_policy_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudAutoSnapshotPolicyApplicationCreate(d *schema.ResourceData, meta interface{}) error {

	err := autoSnapshotPolicyApplication(d, meta)
	if err != nil {
		return err
	}

	d.SetId(d.Get("auto_snapshot_policy_id").(string) + ":" + d.Get("disk_id").(string))

	return resourceAlicloudAutoSnapshotPolicyApplicationRead(d, meta)
}

func resourceAlicloudAutoSnapshotPolicyApplicationRead(d *schema.ResourceData, meta interface{}) error {
	_, diskId, err := getAutoSnapshotPolicyAndDiskID(d, meta)
	if err != nil {
		return err
	}

	conn := meta.(*AliyunClient).ecsconn
	disks, _, err := conn.DescribeDisks(&ecs.DescribeDisksArgs{
		RegionId: getRegion(d, meta),
		DiskIds:  []string{diskId},
	})

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error DescribeDiskAttribute: %#v", err)
	}

	if disks == nil || len(disks) <= 0 {
		return fmt.Errorf("No Disks Found.")
	}

	disk := disks[0]
	d.Set("auto_snapshot_policy_id", disk.AutoSnapshotPolicyId)
	d.Set("disk_id", disk.DiskId)

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":                      d.Id(),
		"name":                    d.Id(),
		"status":                  "Available",
		"creation_time":           time.Now().Format("2006-01-02 15:04:05"),
		"resource_type":           "alicloud_auto_snapshot_policy_application",
		"auto_snapshot_policy_id": disk.AutoSnapshotPolicyId,
		"disk_id":                 disk.DiskId,
	}
	log.Printf("[DEBUG] alicloud_auto_snapshot_policy_application - adding alicloud_auto_snapshot_policy_application: %v", mapping)
	s = append(s, mapping)

	if err := d.Set("alicloud_auto_snapshot_policy_application", s); err != nil {
		return fmt.Errorf("Setting alicloud_auto_snapshot_policy_application got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudAutoSnapshotPolicyApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn
	autoSnapshotPolicyId, diskId, err := getAutoSnapshotPolicyAndDiskID(d, meta)
	if err != nil {
		return err
	}
	diskIds := fmt.Sprintf("[\"%s\"]", diskId)

	args := &ecs.AutoSnapshotPolicyApplicationArgs{
		RegionId:             getRegion(d, meta),
		AutoSnapshotPolicyId: autoSnapshotPolicyId,
		DiskIds:              diskIds,
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.CancelAutoSnapshotPolicy(args)
		if err != nil {
			if IsExceptedError(err, OperationConflict) || IsExceptedError(err, InternalError) ||
				IsExceptedError(err, InvalidOperation) {
				return resource.RetryableError(fmt.Errorf("Cancel AutoSnapshotPolicy timeout and got an error: %#v", err))
			}
			return resource.NonRetryableError(err)
		}

		disks, _, err := conn.DescribeDisks(&ecs.DescribeDisksArgs{
			RegionId: getRegion(d, meta),
			DiskIds:  []string{diskId},
		})

		if err != nil {
			log.Printf("[ERROR] Disk %s DescribeDisksArgs failed.", diskId)
			return resource.NonRetryableError(err)
		}

		if disks == nil || len(disks) <= 0 {
			return resource.RetryableError(fmt.Errorf("Cancel AutoSnapshotPolicy timeout and got an error: %#v", err))
		}

		// disk := disks[0]
		// log.Printf("[DEBUG] Disk AutoSnapshotPolicyId is [%s].", disk.AutoSnapshotPolicyId)
		// log.Printf("[ERROR] Disk EnableAutoSnapshot is [%s].", disk.EnableAutoSnapshot)
		// if disk.AutoSnapshotPolicyId != "" || disk.EnableAutoSnapshot != false {
		// 	log.Printf("[ERROR] Cancel AutoSnapshotPolicy failed.", diskId)
		// 	return resource.NonRetryableError(errors.New("AutoSnapshotPolicyId or EnableAutoSnapshot not correct"))
		// }

		return nil
	})
}

func getAutoSnapshotPolicyAndDiskID(d *schema.ResourceData, meta interface{}) (string, string, error) {
	parts := strings.Split(d.Id(), ":")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}

func autoSnapshotPolicyApplication(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	autoSnapshotPolicyId := d.Get("auto_snapshot_policy_id").(string)
	diskId := d.Get("disk_id").(string)
	diskIds := fmt.Sprintf("[\"%s\"]", diskId)

	args := &ecs.AutoSnapshotPolicyApplicationArgs{
		RegionId:             getRegion(d, meta),
		AutoSnapshotPolicyId: autoSnapshotPolicyId,
		DiskIds:              diskIds,
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.ApplyAutoSnapshotPolicy(args)
		log.Printf("error : %s", err)

		if err != nil {
			if IsExceptedError(err, OperationConflict) || IsExceptedError(err, InternalError) ||
				IsExceptedError(err, InvalidOperation) {
				return resource.RetryableError(fmt.Errorf("Apply AutoSnapshotPolicy timeout and got an error: %#v", err))
			}
			return resource.NonRetryableError(err)
		}

		disks, _, err := conn.DescribeDisks(&ecs.DescribeDisksArgs{
			RegionId: getRegion(d, meta),
			DiskIds:  []string{diskId},
		})

		if err != nil {
			log.Printf("[ERROR] Disk %s DescribeDisksArgs failed.", diskId)
			return resource.NonRetryableError(err)
		}

		if disks == nil || len(disks) <= 0 {
			return resource.RetryableError(fmt.Errorf("Apply AutoSnapshotPolicy timeout and got an error: %#v", err))
		}

		// disk := disks[0]
		// if disk.AutoSnapshotPolicyId != autoSnapshotPolicyId || disk.EnableAutoSnapshot != true {
		// 	log.Printf("[ERROR] Apply AutoSnapshotPolicy failed.", diskId)
		// 	return resource.NonRetryableError(errors.New("AutoSnapshotPolicyId or EnableAutoSnapshot not correct"))
		// }

		return nil

	})
}
