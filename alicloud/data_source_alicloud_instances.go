package alicloud

import (
	"log"
	"regexp"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudInstancesRead,

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
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceStatus,
				ForceNew:     true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tags": tagsSchema(),

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Computed values
			"alicloud_instances": {
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
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_max_bandwidth_out": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"spot_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_device_mappings": {
							Type:     schema.TypeList,
							Computed: true,
							//Set:      imageDiskDeviceMappingHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tags": tagsSchema(),
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
func dataSourceAlicloudInstancesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	args := &ecs.DescribeInstancesArgs{
		Status: ecs.InstanceStatus(d.Get("status").(string)),
	}

	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		args.InstanceIds = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("region_id"); ok && v != "" {
		args.RegionId = common.Region(v.(string))
	} else {
		args.RegionId = getRegion(d, meta)
	}
	if v, ok := d.GetOk("zone_id"); ok && v != "" {
		args.ZoneId = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok && v != "" {
		args.VpcId = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok && v != "" {
		args.VSwitchId = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok && v != "" {
		args.ZoneId = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		mapping := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			mapping[key] = value.(string)
		}
		args.Tag = mapping
	}

	var allInstances []ecs.InstanceAttributesType

	for {
		instances, paginationResult, err := conn.DescribeInstances(args)
		if err != nil {
			return err
		}

		allInstances = append(allInstances, instances...)

		pagination := paginationResult.NextPage()
		if pagination == nil {
			break
		}

		args.Pagination = *pagination
	}

	var filteredInstancesTemp []ecs.InstanceAttributesType

	nameRegex, ok := d.GetOk("name_regex")
	imageId, okImg := d.GetOk("image_id")
	if (ok && nameRegex.(string) != "") || (okImg && imageId.(string) != "") {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, inst := range allInstances {
			if r != nil && !r.MatchString(inst.InstanceName) {
				continue
			}
			if imageId.(string) != "" && inst.ImageId != imageId.(string) {
				continue
			}
			filteredInstancesTemp = append(filteredInstancesTemp, inst)
		}
	} else {
		filteredInstancesTemp = allInstances
	}

	/*
		if len(filteredInstancesTemp) < 1 {
			return fmt.Errorf("Your query returned no results[alicloud_instances]. Please change your search criteria and try again.")
		}
	*/

	log.Printf("[DEBUG] alicloud_instances - Instances found: %#v", filteredInstancesTemp)

	return instancessDescriptionAttributes(d, filteredInstancesTemp, meta)
}

// populate the numerous fields that the instance description returns.
func instancessDescriptionAttributes(d *schema.ResourceData, instances []ecs.InstanceAttributesType, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, inst := range instances {
		mapping := map[string]interface{}{
			"id":                         inst.InstanceId,
			"region_id":                  inst.RegionId,
			"availability_zone":          inst.ZoneId,
			"status":                     inst.Status,
			"name":                       inst.InstanceName,
			"instance_type":              inst.InstanceType,
			"instance_name":              inst.InstanceName,
			"host_name":                  inst.HostName,
			"vpc_id":                     inst.VpcAttributes.VpcId,
			"vswitch_id":                 inst.VpcAttributes.VSwitchId,
			"image_id":                   inst.ImageId,
			"description":                inst.Description,
			"security_groups":            inst.SecurityGroupIds.SecurityGroupId,
			"eip":                        inst.EipAddress.IpAddress,
			"key_name":                   inst.KeyPairName,
			"spot_strategy":              inst.SpotStrategy,
			"creation_time":              inst.CreationTime.String(),
			"instance_charge_type":       inst.InstanceChargeType,
			"internet_charge_type":       inst.InternetChargeType,
			"internet_max_bandwidth_out": inst.InternetMaxBandwidthOut,
			// Complex types get their own functions
			"disk_device_mappings": instanceDisksMappings(d, inst.InstanceId, meta),
			"tags":                 tagsToMap(inst.Tags.Tag),
			"resource_type":        "alicloud_instance",
		}
		if len(inst.InnerIpAddress.IpAddress) > 0 {
			mapping["private_ip"] = inst.InnerIpAddress.IpAddress[0]
		} else {
			mapping["private_ip"] = inst.VpcAttributes.PrivateIpAddress.IpAddress[0]
		}
		if len(inst.PublicIpAddress.IpAddress) > 0 {
			mapping["public_ip"] = inst.PublicIpAddress.IpAddress[0]
		} else {
			mapping["public_ip"] = inst.VpcAttributes.NatIpAddress
		}

		log.Printf("[DEBUG] alicloud_instance - adding instance mapping: %v", mapping)
		ids = append(ids, inst.InstanceId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alicloud_instances", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

//Returns a mapping of instance disks
func instanceDisksMappings(d *schema.ResourceData, instanceId string, meta interface{}) []map[string]interface{} {

	disks, _, err := meta.(*AliyunClient).ecsconn.DescribeDisks(&ecs.DescribeDisksArgs{
		RegionId:   getRegion(d, meta),
		InstanceId: instanceId,
	})

	if err != nil {
		log.Printf("[ERROR] DescribeDisks for instance got error: %#v", err)
		return nil
	}

	var s []map[string]interface{}

	for _, v := range disks {
		mapping := map[string]interface{}{
			"device":   v.Device,
			"size":     v.Size,
			"category": v.Category,
			"type":     v.Type,
		}

		log.Printf("[DEBUG] alicloud_instances - adding disk device mapping: %v", mapping)
		s = append(s, mapping)
	}

	return s
}
