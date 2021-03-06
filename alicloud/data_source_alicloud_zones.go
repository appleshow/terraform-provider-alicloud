package alicloud

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudZonesRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"available_instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceType,
			},
			"available_resource_creation": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ResourceTypeInstance),
					string(ResourceTypeRds),
					string(ResourceTypeVSwitch),
					string(ResourceTypeDisk),
					string(IoOptimized),
				}),
			},
			"available_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateDiskCategory,
			},

			"multi": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validateInstanceChargeType,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(Vpc), string(Classic)}),
			},
			"spot_strategy": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      ecs.NoSpot,
				ValidateFunc: validateInstanceSpotStrategy,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values.
			"alicloud_zones": {
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
						"local_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"available_resource_creation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"available_disk_categories": {
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
			"zone_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAlicloudZonesRead(d *schema.ResourceData, meta interface{}) error {
	resType, _ := d.Get("available_resource_creation").(string)
	multi := d.Get("multi").(bool)
	client := meta.(*AliyunClient)
	var zoneIds []string
	rdsZones := make(map[string]string)
	if strings.ToLower(Trim(resType)) == strings.ToLower(string(ResourceTypeRds)) {
		request := rds.CreateDescribeRegionsRequest()
		if regions, err := client.rdsconn.DescribeRegions(request); err != nil {
			return fmt.Errorf("[ERROR] DescribeRegions got an error: %#v", err)
		} else if len(regions.Regions.RDSRegion) <= 0 {
			return fmt.Errorf("[ERROR] There is no available region for RDS.")
		} else {
			for _, r := range regions.Regions.RDSRegion {
				if multi && strings.Contains(r.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(getRegion(d, meta)) {
					zoneIds = append(zoneIds, r.ZoneId)
					continue
				}
				rdsZones[r.ZoneId] = r.RegionId
			}
		}
	}
	if len(zoneIds) > 0 {
		sort.Strings(zoneIds)
		return multiZonesDescriptionAttributes(d, zoneIds)
	} else if multi {
		return fmt.Errorf("There is no multi zones in the current region %s. Please change region and try again.", getRegion(d, meta))
	}

	_, validZones, err := client.DescribeAvailableResources(d, meta, ZoneResource)
	if err != nil {
		return err
	}

	zones, err := client.ecsconn.DescribeZones(getRegion(d, meta))
	if err != nil {
		return fmt.Errorf("DescribeZones got an error: %#v", err)
	}

	/*
		if zones == nil || len(zones) < 1 {
			return fmt.Errorf("There are no availability zones in the region: %#v.", getRegion(d, meta))
		}
	*/

	mapZones := make(map[string]ecs.ZoneType)
	insType, _ := d.Get("available_instance_type").(string)
	diskType, _ := d.Get("available_disk_category").(string)

	for _, zone := range zones {
		for _, v := range validZones {
			if zone.ZoneId != v.ZoneId {
				continue
			}
			if len(zone.AvailableInstanceTypes.InstanceTypes) <= 0 ||
				(insType != "" && !constraints(zone.AvailableInstanceTypes.InstanceTypes, insType)) {
				continue
			}
			if len(zone.AvailableDiskCategories.DiskCategories) <= 0 ||
				(diskType != "" && !constraints(zone.AvailableDiskCategories.DiskCategories, diskType)) {
				continue
			}
			if len(rdsZones) > 0 {
				if _, ok := rdsZones[zone.ZoneId]; !ok {
					continue
				}
			}
			zoneIds = append(zoneIds, zone.ZoneId)
			mapZones[zone.ZoneId] = zone
		}
	}

	if len(zoneIds) <= 0 {
		return fmt.Errorf("Your query zones returned no results. Please change your search criteria and try again.")
	}

	// Sort zones before reading
	sort.Strings(zoneIds)

	var s []map[string]interface{}
	id, idOk := d.GetOk("id")
	for _, zoneId := range zoneIds {
		if idOk && id.(string) != zoneId {
			continue
		}
		mapping := map[string]interface{}{
			"id":                          zoneId,
			"name":                        mapZones[zoneId].LocalName,
			"status":                      "Available",
			"creation_time":               "",
			"local_name":                  mapZones[zoneId].LocalName,
			"available_instance_types":    mapZones[zoneId].AvailableInstanceTypes.InstanceTypes,
			"available_resource_creation": mapZones[zoneId].AvailableResourceCreation.ResourceTypes,
			"available_disk_categories":   mapZones[zoneId].AvailableDiskCategories.DiskCategories,
			"resource_type":               "alicloud_zone",
		}
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(zoneIds))
	if err := d.Set("alicloud_zones", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	if err := d.Set("zone_ids", zoneIds); err != nil {
		return err
	}

	return nil
}

// check array constraints str
func constraints(arr interface{}, v string) bool {
	arrs := reflect.ValueOf(arr)
	len := arrs.Len()
	for i := 0; i < len; i++ {
		if arrs.Index(i).String() == v {
			return true
		}
	}
	return false
}

func multiZonesDescriptionAttributes(d *schema.ResourceData, zones []string) error {
	var s []map[string]interface{}
	id, idOk := d.GetOk("id")
	for _, t := range zones {
		if idOk && id.(string) != t {
			continue
		}
		mapping := map[string]interface{}{
			"id": t,
		}
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(zones))
	if err := d.Set("alicloud_zones", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
