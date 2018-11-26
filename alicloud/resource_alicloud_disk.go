package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDiskCreate,
		Read:   resourceAliyunDiskRead,
		Update: resourceAliyunDiskUpdate,
		Delete: resourceAliyunDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDiskName,
			},

			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDiskDescription,
			},

			"category": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateDiskCategory,
				Default:      "cloud_efficiency",
			},

			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"encrypted": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"auto_snapshot_policy_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"enable_auto_snapshot": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tagsSchema(),

			"creation_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			// Computed values
			"alicloud_disk": {
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
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"portable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"operation_locks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_auto_snapshot": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"attached_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detached_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						}, /*
							"image_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"source_snapshot_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"product_code": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"delete_auto_snapshot": {
								Type:     schema.TypeBool,
								Computed: true,
							},*/
					},
				},
			},
		},
	}
}

func resourceAliyunDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	conn := client.ecsconn

	availabilityZone, err := client.DescribeZone(d.Get("availability_zone").(string))
	if err != nil {
		return err
	}

	args := &ecs.CreateDiskArgs{
		RegionId: getRegion(d, meta),
		ZoneId:   availabilityZone.ZoneId,
	}

	if v, ok := d.GetOk("category"); ok && v.(string) != "" {
		category := ecs.DiskCategory(v.(string))
		if err := client.DiskAvailable(availabilityZone, category); err != nil {
			return err
		}
		args.DiskCategory = category
	}

	if v, ok := d.GetOk("size"); ok {
		size := v.(int)
		if args.DiskCategory == ecs.DiskCategoryCloud && (size < 5 || size > 2000) {
			return fmt.Errorf("the size of cloud disk must between 5 to 2000")
		}

		if (args.DiskCategory == ecs.DiskCategoryCloudEfficiency ||
			args.DiskCategory == ecs.DiskCategoryCloudSSD) && (size < 20 || size > 32768) {
			return fmt.Errorf("the size of %s disk must between 20 to 32768", args.DiskCategory)
		}
		args.Size = size

		d.Set("size", args.Size)
	}

	if v, ok := d.GetOk("snapshot_id"); ok && v.(string) != "" {
		args.SnapshotId = v.(string)
	}

	if args.Size <= 0 && args.SnapshotId == "" {
		return fmt.Errorf("One of size or snapshot_id is required when specifying an ECS disk.")
	}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		args.DiskName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		args.Description = v.(string)
	}

	if v, ok := d.GetOk("encrypted"); ok {
		args.Encrypted = v.(bool)
	}

	diskID, err := conn.CreateDisk(args)
	if err != nil {
		return fmt.Errorf("CreateDisk got a error: %#v", err)
	}

	d.SetId(diskID)

	return resourceAliyunDiskUpdate(d, meta)
}

func resourceAliyunDiskRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	disks, _, err := conn.DescribeDisks(&ecs.DescribeDisksArgs{
		RegionId: getRegion(d, meta),
		DiskIds:  []string{d.Id()},
	})

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error DescribeDiskAttribute: %#v", err)
	}

	log.Printf("[DEBUG] DescribeDiskAttribute for instance: %#v", disks)

	if disks == nil || len(disks) <= 0 {
		return fmt.Errorf("No disks found.")
	}

	disk := disks[0]
	d.Set("availability_zone", disk.ZoneId)
	d.Set("category", disk.Category)
	d.Set("size", disk.Size)
	d.Set("status", disk.Status)
	d.Set("name", disk.DiskName)
	d.Set("description", disk.Description)
	d.Set("snapshot_id", disk.SourceSnapshotId)
	d.Set("encrypted", disk.Encrypted)
	d.Set("auto_snapshot_policy_id", disk.AutoSnapshotPolicyId)
	d.Set("enable_auto_snapshot", disk.EnableAutoSnapshot)
	d.Set("creation_time", disk.CreationTime)

	tags, _, err := conn.DescribeTags(&ecs.DescribeTagsArgs{
		RegionId:     getRegion(d, meta),
		ResourceType: ecs.TagResourceDisk,
		ResourceId:   d.Id(),
	})

	if err != nil {
		log.Printf("[DEBUG] DescribeTags for disk got error: %#v", err)
	}

	d.Set("tags", tagsToMap(tags))

	var s []map[string]interface{}
	var operationLocks []string
	for _, lockReasonType := range disk.OperationLocks.LockReason {
		operationLocks = append(operationLocks, string(lockReasonType.LockReason))
	}
	mapping := map[string]interface{}{
		"id":                   disk.DiskId,
		"name":                 disk.DiskName,
		"status":               disk.Status,
		"creation_time":        disk.CreationTime.String(),
		"resource_type":        "alicloud_disk",
		"region_id":            disk.RegionId,
		"zone_id":              disk.ZoneId,
		"description":          disk.Description,
		"type":                 disk.Type,
		"encrypted":            disk.Encrypted,
		"category":             disk.Category,
		"size":                 disk.Size,
		"portable":             disk.Portable,
		"operation_locks":      operationLocks,
		"instance_id":          disk.InstanceId,
		"device":               disk.Device,
		"delete_with_instance": disk.DeleteWithInstance,
		"enable_auto_snapshot": disk.EnableAutoSnapshot,
		"attached_time":        disk.AttachedTime.String(),
		"detached_time":        disk.DetachedTime.String(),
		"disk_charge_type":     disk.DiskChargeType,
		/*
			"image_id":             disk.ImageId,
			"source_snapshot_id":   disk.SourceSnapshotId,
			"product_code":         disk.ProductCode,
			"delete_auto_snapshot": disk.DeleteAutoSnapshot,
		*/
	}

	s = append(s, mapping)
	if err := d.Set("alicloud_disk", s); err != nil {
		return fmt.Errorf("Setting alicloud_disk got an error: %#v.", err)
	}

	return nil
}

func resourceAliyunDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.ecsconn

	d.Partial(true)

	if err := setTags(client, ecs.TagResourceDisk, d); err != nil {
		log.Printf("[DEBUG] Set tags for instance got error: %#v", err)
		return fmt.Errorf("Set tags for instance got error: %#v", err)
	} else {
		d.SetPartial("tags")
	}
	attributeUpdate := false
	args := &ecs.ModifyDiskAttributeArgs{
		DiskId: d.Id(),
	}

	if d.HasChange("name") {
		d.SetPartial("name")
		val := d.Get("name").(string)
		args.DiskName = val

		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		val := d.Get("description").(string)
		args.Description = val

		attributeUpdate = true
	}
	if attributeUpdate {
		if err := conn.ModifyDiskAttribute(args); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAliyunDiskRead(d, meta)
}

func resourceAliyunDiskDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.DeleteDisk(d.Id())
		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code == DiskIncorrectStatus || e.ErrorResponse.Code == DiskCreatingSnapshot {
				return resource.RetryableError(fmt.Errorf("Disk in use - trying again while it is deleted."))
			}
		}

		disks, _, descErr := conn.DescribeDisks(&ecs.DescribeDisksArgs{
			RegionId: getRegion(d, meta),
			DiskIds:  []string{d.Id()},
		})

		if descErr != nil {
			log.Printf("[ERROR] Delete disk is failed.")
			return resource.NonRetryableError(descErr)
		}
		if disks == nil || len(disks) < 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Disk in use - trying again while it is deleted."))
	})
}
