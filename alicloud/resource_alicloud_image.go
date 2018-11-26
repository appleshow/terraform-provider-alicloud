package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	yunecs "github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudImageCreate,
		Read:   resourceAlicloudImageRead,
		Update: resourceAlicloudImageUpdate,
		Delete: resourceAlicloudImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_device_mapping": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"image_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"platform": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"architecture": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MinItems: 0,
				MaxItems: 5,
			},
			// Computed values
			"alicloud_image": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	createImageRequest := ecs.CreateCreateImageRequest()
	snapshotId, snapshotIdOk := d.GetOk("snapshot_id")
	instanceId, instanceIdOk := d.GetOk("instance_id")
	diskDeviceMapping, diskDeviceMappingOk := d.GetOk("disk_device_mapping")

	if (!snapshotIdOk && !instanceIdOk && diskDeviceMappingOk) || (snapshotIdOk && instanceIdOk) || (snapshotIdOk && diskDeviceMappingOk) || (instanceIdOk && diskDeviceMappingOk) {
		return fmt.Errorf("Creating Image got an error: %#v.", "snapshot_id, instance_id and disk_device_mapping are entered at least one and only one can be entered.")
	}
	if snapshotIdOk {
		createImageRequest.SnapshotId = snapshotId.(string)
	}
	if instanceIdOk {
		createImageRequest.InstanceId = instanceId.(string)
	}
	if diskDeviceMappingOk {
		var diskDeviceMappingList []ecs.CreateImageDiskDeviceMapping
		diskDeviceMappingPars := diskDeviceMapping.([]interface{})
		for _, diskDeviceMappingPar := range diskDeviceMappingPars {
			diskDeviceMappingTmp := diskDeviceMappingPar.(schema.ResourceData)

			var createImageDiskDeviceMapping ecs.CreateImageDiskDeviceMapping
			createImageDiskDeviceMapping.SnapshotId = diskDeviceMappingTmp.Get("snapshot_id").(string)
			createImageDiskDeviceMapping.DiskType = diskDeviceMappingTmp.Get("disk_type").(string)
			createImageDiskDeviceMapping.Size = diskDeviceMappingTmp.Get("size").(string)

			diskDeviceMappingList = append(diskDeviceMappingList, createImageDiskDeviceMapping)
		}
		createImageRequest.DiskDeviceMapping = &diskDeviceMappingList
	}
	if imageName, ok := d.GetOk("image_name"); ok {
		createImageRequest.ImageName = imageName.(string)
	}
	if description, ok := d.GetOk("description"); ok {
		createImageRequest.Description = description.(string)
	}
	if platform, ok := d.GetOk("platform"); ok {
		createImageRequest.Platform = platform.(string)
	}
	if architecture, ok := d.GetOk("architecture"); ok {
		createImageRequest.Architecture = architecture.(string)
	}
	if imageVersion, ok := d.GetOk("image_version"); ok {
		createImageRequest.ImageVersion = imageVersion.(string)
	}
	if tag, ok := d.GetOk("tag"); ok {
		tagPars := tag.([]interface{})
		for index, tagPar := range tagPars {
			tagParTemp := tagPar.(schema.ResourceData)
			if index == 0 {
				createImageRequest.Tag1Key = tagParTemp.Get("key").(string)
				createImageRequest.Tag1Value = tagParTemp.Get("value").(string)
			}
			if index == 1 {
				createImageRequest.Tag2Key = tagParTemp.Get("key").(string)
				createImageRequest.Tag2Value = tagParTemp.Get("value").(string)
			}
			if index == 2 {
				createImageRequest.Tag3Key = tagParTemp.Get("key").(string)
				createImageRequest.Tag3Value = tagParTemp.Get("value").(string)
			}
			if index == 3 {
				createImageRequest.Tag4Key = tagParTemp.Get("key").(string)
				createImageRequest.Tag4Value = tagParTemp.Get("value").(string)
			}
			if index == 4 {
				createImageRequest.Tag5Key = tagParTemp.Get("key").(string)
				createImageRequest.Tag5Value = tagParTemp.Get("value").(string)
			}
		}
	}

	createImageResponse, err := client.aliecsconn.CreateImage(createImageRequest)
	if err != nil {
		return fmt.Errorf("Creating Image got an error: %#v.", err)
	}

	d.SetId(createImageResponse.ImageId)

	return resourceAlicloudImageUpdate(d, meta)
}

func resourceAlicloudImageRead(d *schema.ResourceData, meta interface{}) error {
	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":            d.Id(),
		"name":          d.Get("image_name").(string),
		"status":        "Available",
		"creation_time": time.Now().Format("2006-01-02 15:04:05"),
		"resource_type": "alicloud_image",
		"description":   d.Get("description").(string),
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_image", s); err != nil {
		return fmt.Errorf("Setting alicloud_image got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	d.Partial(true)
	update := false

	if d.HasChange("snapshot_id") && !d.IsNewResource() {
		return fmt.Errorf("Updating image got an error: %#v", "Cannot modify parameter snapshot_id")
	}
	if d.HasChange("instance_id") && !d.IsNewResource() {
		return fmt.Errorf("Updating image got an error: %#v", "Cannot modify parameter instance_id")
	}
	if d.HasChange("disk_device_mapping") && !d.IsNewResource() {
		return fmt.Errorf("Updating image got an error: %#v", "Cannot modify parameter disk_device_mapping")
	}
	if d.HasChange("tag") && !d.IsNewResource() {
		return fmt.Errorf("Updating image got an error: %#v", "Cannot modify parameter tag")
	}
	if d.HasChange("platform") && !d.IsNewResource() {
		return fmt.Errorf("Updating image got an error: %#v", "Cannot modify parameter snapshot_id")
	}
	if d.HasChange("architecture") && !d.IsNewResource() {
		return fmt.Errorf("Updating image got an error: %#v", "Cannot modify parameter snapshot_id")
	}
	if d.HasChange("image_version") && !d.IsNewResource() {
		return fmt.Errorf("Updating image got an error: %#v", "Cannot modify parameter snapshot_id")
	}

	if d.HasChange("image_name") {
		update = true
		d.SetPartial("image_name")
	}
	if d.HasChange("description") {
		update = true
		d.SetPartial("description")
	}

	if !d.IsNewResource() && update {
		modifyImageAttributeRequest := ecs.CreateModifyImageAttributeRequest()

		modifyImageAttributeRequest.ImageId = d.Id()
		modifyImageAttributeRequest.ImageName = d.Get("image_name").(string)
		modifyImageAttributeRequest.Description = d.Get("description").(string)

		_, err := client.aliecsconn.ModifyImageAttribute(modifyImageAttributeRequest)
		if err != nil {
			return fmt.Errorf("Updating image attribute got an error: %#v.", err)
		}
	}

	d.Partial(false)

	return resourceAlicloudImageRead(d, meta)
}

func resourceAlicloudImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	deleteImageRequest := ecs.CreateDeleteImageRequest()
	deleteImageRequest.ImageId = d.Id()

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.aliecsconn.DeleteImage(deleteImageRequest)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting image got an error: %#v", err))
		}

		describeImagesArgs := &yunecs.DescribeImagesArgs{
			RegionId: getRegion(d, meta),
		}
		describeImagesArgs.ImageId = d.Id()

		conn := meta.(*AliyunClient).ecsconn
		images, _, err := conn.DescribeImages(describeImagesArgs)
		if err != nil {
			return resource.RetryableError(fmt.Errorf("Deleting Image Share Permission got an error: %#v", err))
		}
		if images == nil || len(images) == 0 {
			return nil
		}

		return nil
	})
}
