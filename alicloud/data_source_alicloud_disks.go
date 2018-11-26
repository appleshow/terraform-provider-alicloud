package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudDisks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDisksRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "all",
				ValidateFunc: validateAllowedStringValue([]string{
					"all", "system", "data",
				}),
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "all",
				ValidateFunc: validateAllowedStringValue([]string{
					"all", "cloud", "ephemeral", "ephemeral_ssd", "cloud_efficiency", "cloud_ssd",
				}),
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "All",
				ValidateFunc: validateAllowedStringValue([]string{
					"All", "In_use", "Available", "Attaching", "Detaching", "Creating", "ReIniting",
				}),
			},
			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"portable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"delete_with_instance": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"delete_auto_snapshot": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"enable_auto_snapshot": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"disk_charge_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"PrePaid", "PostPaid",
				}),
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Computed values
			"alicloud_disks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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
						"name": {
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
						"portable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
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
						"aelete_auto_snapshot": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_auto_snapshot": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
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
func dataSourceAlicloudDisksRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	args := &ecs.DescribeDisksArgs{
		Status: ecs.DiskStatus(d.Get("status").(string)),
	}

	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		var ids []string
		for _, id := range v.([]interface{}) {
			ids = append(ids, id.(string))
		}
		args.DiskIds = ids
	}
	if v, ok := d.GetOk("region_id"); ok && v != "" {
		args.RegionId = common.Region(v.(string))
	} else {
		args.RegionId = getRegion(d, meta)
	}
	if v, ok := d.GetOk("zone_id"); ok && v != "" {
		args.ZoneId = v.(string)
	}
	if v, ok := d.GetOk("instance_id"); ok && v != "" {
		args.InstanceId = v.(string)
	}
	if v, ok := d.GetOk("disk_type"); ok && v != "" {
		args.DiskType = ecs.DiskType(v.(string))
	}
	if v, ok := d.GetOk("category"); ok && v != "" {
		args.Category = ecs.DiskCategory(v.(string))
	}
	if v, ok := d.GetOk("snapshot_id"); ok && v != "" {
		args.SnapshotId = v.(string)
	}
	if v, ok := d.GetOk("portable"); ok && v != "" {
		vb := v.(bool)
		args.Portable = &vb
	}
	if v, ok := d.GetOk("delete_with_instance"); ok && v != "" {
		vb := v.(bool)
		args.DeleteWithInstance = &vb
	}
	if v, ok := d.GetOk("delete_auto_snapshot"); ok && v != "" {
		vb := v.(bool)
		args.DeleteAutoSnapshot = &vb
	}
	if v, ok := d.GetOk("enable_auto_snapshot"); ok && v != "" {
		vb := v.(bool)
		args.EnableAutoSnapshot = &vb
	}
	if v, ok := d.GetOk("disk_charge_type"); ok && v != "" {
		args.DiskChargeType = ecs.DiskChargeType(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		mapping := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			mapping[key] = value.(string)
		}
		args.Tag = mapping
	}

	var allDisks []ecs.DiskItemType

	for {
		disks, paginationResult, err := conn.DescribeDisks(args)
		if err != nil {
			return fmt.Errorf("List disks got an error: %#v", err)
		}

		allDisks = append(allDisks, disks...)

		pagination := paginationResult.NextPage()
		if pagination == nil {
			break
		}

		args.Pagination = *pagination
	}

	log.Printf("[DEBUG] alicloud_disks - Disks found: %#v", allDisks)

	var ids []string
	var s []map[string]interface{}
	var r *regexp.Regexp

	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	for _, disk := range allDisks {
		if r != nil && !r.MatchString(disk.DiskName) {
			continue
		}

		var operationLocks []string
		for _, lockReasonType := range disk.OperationLocks.LockReason {
			operationLocks = append(operationLocks, string(lockReasonType.LockReason))
		}
		mapping := map[string]interface{}{
			"id":                   disk.DiskId,
			"region_id":            disk.RegionId,
			"zone_id":              disk.ZoneId,
			"name":                 disk.DiskName,
			"description":          disk.Description,
			"type":                 disk.Type,
			"encrypted":            disk.Encrypted,
			"category":             disk.Category,
			"size":                 disk.Size,
			"image_id":             disk.ImageId,
			"source_snapshot_id":   disk.SourceSnapshotId,
			"product_code":         disk.ProductCode,
			"portable":             disk.Portable,
			"status":               disk.Status,
			"operation_locks":      operationLocks,
			"instance_id":          disk.InstanceId,
			"device":               disk.Device,
			"delete_with_instance": disk.DeleteWithInstance,
			"aelete_auto_snapshot": disk.DeleteAutoSnapshot,
			"enable_auto_snapshot": disk.EnableAutoSnapshot,
			"creation_time":        disk.CreationTime.String(),
			"attached_time":        disk.AttachedTime.String(),
			"detached_time":        disk.DetachedTime.String(),
			"disk_charge_type":     disk.DiskChargeType,
			"resource_type":        "alicloud_disk",
		}

		ids = append(ids, disk.DiskId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alicloud_disks", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
