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

func resourceAliyunDiskAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDiskAttachmentCreate,
		Read:   resourceAliyunDiskAttachmentRead,
		Delete: resourceAliyunDiskAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delete_with_instance": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"device_name": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "Attribute device_name is deprecated on disk attachment resource. Suggest to remove it from your template.",
			},
			// Computed values
			"alicloud_disk_attachment": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunDiskAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	err := diskAttachment(d, meta)
	if err != nil {
		return err
	}

	d.SetId(d.Get("disk_id").(string) + ":" + d.Get("instance_id").(string))

	return resourceAliyunDiskAttachmentRead(d, meta)
}

func resourceAliyunDiskAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	diskId, instanceId, err := getDiskIDAndInstanceID(d, meta)
	if err != nil {
		return err
	}

	conn := meta.(*AliyunClient).ecsconn
	disks, _, err := conn.DescribeDisks(&ecs.DescribeDisksArgs{
		RegionId:   getRegion(d, meta),
		InstanceId: instanceId,
		DiskIds:    []string{diskId},
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
		return fmt.Errorf("No Disks Found.")
	}

	disk := disks[0]
	d.Set("instance_id", disk.InstanceId)
	d.Set("disk_id", disk.DiskId)
	d.Set("device_name", disk.Device)

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":            d.Id(),
		"name":          d.Id(),
		"status":        "Available",
		"creation_time": time.Now().Format("2006-01-02 15:04:05"),
		"type":          "alicloud_disk_attachment",
		"instance_id":   disk.InstanceId,
		"disk_id":       disk.DiskId,
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_disk_attachment", s); err != nil {
		return fmt.Errorf("Setting alicloud_disk_attachment got an error: %#v.", err)
	}

	return nil
}

func resourceAliyunDiskAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn
	diskID, instanceID, err := getDiskIDAndInstanceID(d, meta)
	if err != nil {
		return err
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.DetachDisk(instanceID, diskID)
		if err != nil {
			if IsExceptedError(err, DiskIncorrectStatus) || IsExceptedError(err, InstanceLockedForSecurity) ||
				IsExceptedError(err, DiskInvalidOperation) {
				return resource.RetryableError(fmt.Errorf("Detach Disk timeout and got an error: %#v", err))
			}
		}

		disks, _, descErr := conn.DescribeDisks(&ecs.DescribeDisksArgs{
			RegionId: getRegion(d, meta),
			DiskIds:  []string{diskID},
		})

		if descErr != nil {
			log.Printf("[ERROR] Disk %s is not detached.", diskID)
			return resource.NonRetryableError(err)
		}

		for _, disk := range disks {
			if disk.Status != ecs.DiskStatusAvailable {
				return resource.RetryableError(fmt.Errorf("Detach Disk timeout and got an error: %#v", err))
			}
		}
		return nil
	})
}

func getDiskIDAndInstanceID(d *schema.ResourceData, meta interface{}) (string, string, error) {
	parts := strings.Split(d.Id(), ":")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}

func diskAttachment(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	diskID := d.Get("disk_id").(string)
	instanceID := d.Get("instance_id").(string)
	deleteWithInstance := d.Get("delete_with_instance").(bool)

	args := &ecs.AttachDiskArgs{
		InstanceId:         instanceID,
		DiskId:             diskID,
		DeleteWithInstance: deleteWithInstance,
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.AttachDisk(args)
		log.Printf("error : %s", err)

		if err != nil {
			if IsExceptedError(err, DiskIncorrectStatus) || IsExceptedError(err, InstanceIncorrectStatus) ||
				IsExceptedError(err, DiskOperationConflict) || IsExceptedError(err, DiskInternalError) ||
				IsExceptedError(err, DiskInvalidOperation) {
				return resource.RetryableError(fmt.Errorf("Attach Disk timeout and got an error: %#v", err))
			}
			return resource.NonRetryableError(err)
		}

		disks, _, descErr := conn.DescribeDisks(&ecs.DescribeDisksArgs{
			RegionId:   getRegion(d, meta),
			InstanceId: instanceID,
			DiskIds:    []string{diskID},
		})

		if descErr != nil {
			log.Printf("[ERROR] Disk %s is not attached.", diskID)
			return resource.NonRetryableError(err)
		}

		if disks == nil || len(disks) <= 0 {
			return resource.RetryableError(fmt.Errorf("Attach Disk timeout and got an error: %#v", err))
		}

		return nil

	})
}
