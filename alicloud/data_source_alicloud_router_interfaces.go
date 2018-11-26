package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudRouterInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRouterInterfacesRead,

		Schema: map[string]*schema.Schema{
			"owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"resource_owner_account": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  50,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_router_interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"router_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec": {
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
						"router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"router_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connected_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_interface_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_interface_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_interface_business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_router_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_interface_owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_point_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_access_point_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_source_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_target_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_vpc_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_instance_id": {
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

func dataSourceAlicloudRouterInterfacesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	describeRouterInterfacesRequest := vpc.CreateDescribeRouterInterfacesRequest()
	if ownerId, ok := d.GetOk("owner_id"); ok {
		describeRouterInterfacesRequest.OwnerId = requests.NewInteger(ownerId.(int))
	}
	if resourceOwnerAccount, ok := d.GetOk("resource_owner_account"); ok {
		describeRouterInterfacesRequest.ResourceOwnerAccount = resourceOwnerAccount.(string)
	}
	if pageNumber, ok := d.GetOk("page_number"); ok {
		describeRouterInterfacesRequest.PageNumber = requests.NewInteger(pageNumber.(int))
	}
	if pageSize, ok := d.GetOk("page_size"); ok {
		describeRouterInterfacesRequest.PageSize = requests.NewInteger(pageSize.(int))
	}

	describeRouterInterfacesResponse, err := client.vpcconn.DescribeRouterInterfaces(describeRouterInterfacesRequest)
	if err != nil {
		return fmt.Errorf("Describe router interfaces got an error: %#v", err)
	} else {
		log.Printf("[DEBUG] Describe router interfaces - interfaces found: %#v", describeRouterInterfacesResponse)

		var ids []string
		var s []map[string]interface{}
		for _, routerInterface := range describeRouterInterfacesResponse.RouterInterfaceSet.RouterInterfaceType {
			mapping := map[string]interface{}{
				"id": routerInterface.RouterInterfaceId,
				"router_interface_id":                routerInterface.RouterInterfaceId,
				"opposite_region_id":                 routerInterface.OppositeRegionId,
				"role":                               routerInterface.Role,
				"spec":                               routerInterface.Spec,
				"name":                               routerInterface.Name,
				"description":                        routerInterface.Description,
				"router_id":                          routerInterface.RouterId,
				"router_type":                        routerInterface.RouterType,
				"creation_time":                      routerInterface.CreationTime,
				"end_time":                           routerInterface.EndTime,
				"charge_type":                        routerInterface.ChargeType,
				"status":                             routerInterface.Status,
				"business_status":                    routerInterface.BusinessStatus,
				"connected_time":                     routerInterface.ConnectedTime,
				"opposite_interface_id":              routerInterface.OppositeInterfaceId,
				"opposite_interface_spec":            routerInterface.OppositeInterfaceSpec,
				"opposite_interface_status":          routerInterface.OppositeInterfaceStatus,
				"opposite_interface_business_status": routerInterface.OppositeInterfaceBusinessStatus,
				"opposite_router_id":                 routerInterface.OppositeRouterId,
				"opposite_router_type":               routerInterface.OppositeRouterType,
				"opposite_interface_owner_id":        routerInterface.OppositeInterfaceOwnerId,
				"access_point_id":                    routerInterface.AccessPointId,
				"opposite_access_point_id":           routerInterface.OppositeAccessPointId,
				"health_check_source_ip":             routerInterface.HealthCheckSourceIp,
				"health_check_target_ip":             routerInterface.HealthCheckTargetIp,
				"opposite_vpc_instance_id":           routerInterface.OppositeVpcInstanceId,
				"vpc_instance_id":                    routerInterface.VpcInstanceId,
				"resource_type":                      "alicloud_router_interface",
			}
			ids = append(ids, routerInterface.RouterInterfaceId)
			s = append(s, mapping)
		}

		d.SetId(dataResourceIdHash(ids))
		if err := d.Set("alicloud_router_interfaces", s); err != nil {
			return err
		}

		// create a json file in current directory and write data source to it.
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), s)
		}
	}

	return nil
}
