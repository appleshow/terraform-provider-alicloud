package alicloud

import (
	"fmt"
	"log"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRegionsRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"current": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			//Computed value
			"alicloud_regions": &schema.Schema{
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
						"local_name": {
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

func dataSourceAlicloudRegionsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn
	currentRegion := getRegion(d, meta)

	resp, err := conn.DescribeRegions()
	if err != nil {
		return err
	}
	if resp == nil || len(resp) == 0 {
		return fmt.Errorf("no matching regions found")
	}
	name, nameOk := d.GetOk("name")
	id, idOk := d.GetOk("id")
	current := d.Get("current").(bool)
	var filterRegions []ecs.RegionType
	for _, region := range resp {
		if current {
			if nameOk && common.Region(name.(string)) != currentRegion {
				return fmt.Errorf("name doesn't match current region: %#v, please input again.", currentRegion)
			}
			if region.RegionId == currentRegion {
				filterRegions = append(filterRegions, region)
				break
			}
			continue
		}
		if nameOk {
			if common.Region(name.(string)) == region.RegionId {
				if idOk && common.Region(id.(string)) != region.RegionId {
					continue
				}
				filterRegions = append(filterRegions, region)
				break
			}
			continue
		}
		if idOk && common.Region(id.(string)) != region.RegionId {
			continue
		}
		filterRegions = append(filterRegions, region)
	}
	if len(filterRegions) < 1 {
		return fmt.Errorf("Your query region returned no results. Please change your search criteria and try again.")
	}

	return regionsDescriptionAttributes(d, filterRegions)
}

func regionsDescriptionAttributes(d *schema.ResourceData, regions []ecs.RegionType) error {
	var ids []string
	var s []map[string]interface{}
	for _, region := range regions {
		mapping := map[string]interface{}{
			"id":            region.RegionId,
			"name":          region.LocalName,
			"status":        "Available",
			"creation_time": "",
			"region_id":     region.RegionId,
			"local_name":    region.LocalName,
			"resource_type": "alicloud_region",
		}

		log.Printf("[DEBUG] alicloud_regions - adding region mapping: %v", mapping)
		ids = append(ids, string(region.RegionId))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alicloud_regions", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
