package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRouterInterfaceConnect() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRouterInterfaceConnectCreate,
		Read:   resourceAlicloudRouterInterfaceConnectRead,
		Update: resourceAlicloudRouterInterfaceConnectUpdate,
		Delete: resourceAlicloudRouterInterfaceConnectDelete,

		Schema: map[string]*schema.Schema{
			"router_interface_from_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"router_interface_from_owner_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"router_interface_from_region_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"router_interface_to_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"router_interface_to_owner_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"router_interface_to_region_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudRouterInterfaceConnectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	routerInterfaceFromId := d.Get("router_interface_from_id").(string)
	routerInterfaceToId := d.Get("router_interface_to_id").(string)
	routerInterfaceFromOwnerId := d.Get("router_interface_from_owner_id").(string)
	routerInterfaceToOwnerId := d.Get("router_interface_to_owner_id").(string)
	routerInterfaceFromRegionId := d.Get("router_interface_from_region_id").(string)
	routerInterfaceToRegionId := d.Get("router_interface_to_region_id").(string)

	from, err := meta.(*AliyunClient).DescribeRouterInterface(routerInterfaceFromId)
	if err != nil {
		return fmt.Errorf("DescribeRouterInterface[ID = %s] got an error: %#v", routerInterfaceFromId, err)
	}
	to, err := meta.(*AliyunClient).DescribeRouterInterface(routerInterfaceToId)
	if err != nil {
		return fmt.Errorf("DescribeRouterInterface[ID = %s] got an error: %#v", routerInterfaceToId, err)
	}

	if "InitiatingSide" != from.Role {
		return fmt.Errorf("The role of the router interface[ID = %s] showed be %s.", routerInterfaceFromId, "InitiatingSide")
	}
	if "AcceptingSide" != to.Role {
		return fmt.Errorf("The role of the router interface[ID = %s] showed be %s.", routerInterfaceToId, "AcceptingSide")
	}
	modifyRequestFrom := vpc.CreateModifyRouterInterfaceAttributeRequest()
	modifyRequestFrom.RegionId = routerInterfaceFromRegionId
	modifyRequestFrom.RouterInterfaceId = routerInterfaceFromId

	modifyRequestFrom.OppositeInterfaceId = routerInterfaceToId
	modifyRequestFrom.OppositeRouterId = to.RouterId
	modifyRequestFrom.OppositeRouterType = to.RouterType
	modifyRequestFrom.OppositeInterfaceOwnerId = requests.Integer(routerInterfaceToOwnerId)

	if _, err := client.vpcconn.ModifyRouterInterfaceAttribute(modifyRequestFrom); err != nil {
		return fmt.Errorf("ModifyRouterInterfaceAttribute[ID = %s] got an error: %#v", routerInterfaceFromId, err)
	}

	modifyRequestTo := vpc.CreateModifyRouterInterfaceAttributeRequest()
	modifyRequestTo.RegionId = routerInterfaceToRegionId
	modifyRequestTo.RouterInterfaceId = routerInterfaceToId

	modifyRequestTo.OppositeInterfaceId = routerInterfaceFromId
	modifyRequestTo.OppositeRouterId = from.RouterId
	modifyRequestTo.OppositeRouterType = from.RouterType
	modifyRequestTo.OppositeInterfaceOwnerId = requests.Integer(routerInterfaceFromOwnerId)

	if _, err := client.vpcconn.ModifyRouterInterfaceAttribute(modifyRequestTo); err != nil {
		return fmt.Errorf("ModifyRouterInterfaceAttribute[ID = %s] got an error: %#v", modifyRequestTo, err)
	}

	connectRouterInterfaceRequest := vpc.CreateConnectRouterInterfaceRequest()
	connectRouterInterfaceRequest.RouterInterfaceId = routerInterfaceFromId

	if _, err := client.vpcconn.ConnectRouterInterface(connectRouterInterfaceRequest); err != nil {
		return fmt.Errorf("ConnectRouterInterface got an error: %#v", err)
	}

	d.SetId(routerInterfaceFromId + COMMA_SEPARATED + routerInterfaceToId)
	return resourceAlicloudRouterInterfaceConnectUpdate(d, meta)
}

func resourceAlicloudRouterInterfaceConnectUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("router_interface_from_id") && !d.IsNewResource() {
		return fmt.Errorf("Updating router interface connect got an error: %#v", "Cannot modify parameter router_interface_from_id")
	}
	if d.HasChange("router_interface_to_id") && !d.IsNewResource() {
		return fmt.Errorf("Updating router interface connect got an error: %#v", "Cannot modify parameter router_interface_to_id")
	}
	if d.HasChange("router_interface_from_owner_id") && !d.IsNewResource() {
		return fmt.Errorf("Updating router interface connect got an error: %#v", "Cannot modify parameter router_interface_from_owner_id")
	}
	if d.HasChange("router_interface_to_owner_id") && !d.IsNewResource() {
		return fmt.Errorf("Updating router interface connect got an error: %#v", "Cannot modify parameter router_interface_to_owner_id")
	}
	if d.HasChange("router_interface_from_region_id") && !d.IsNewResource() {
		return fmt.Errorf("Updating router interface connect got an error: %#v", "Cannot modify parameter router_interface_from_region_id")
	}
	if d.HasChange("router_interface_to_region_id") && !d.IsNewResource() {
		return fmt.Errorf("Updating router interface connect got an error: %#v", "Cannot modify parameter router_interface_to_region_id")
	}

	d.Partial(false)

	return resourceAlicloudRouterInterfaceConnectRead(d, meta)
}

func resourceAlicloudRouterInterfaceConnectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parameters := strings.Split(d.Id(), COMMA_SEPARATED)

	from, err := client.DescribeRouterInterface(parameters[0])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	to, err := client.DescribeRouterInterface(parameters[1])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	d.Set("router_interface_from_id", from.RouterInterfaceId)
	d.Set("router_interface_to_id", to.RouterInterfaceId)
	d.Set("router_interface_from_owner_id", to.OppositeInterfaceOwnerId)
	d.Set("router_interface_to_owner_id", from.OppositeInterfaceOwnerId)
	d.Set("router_interface_from_region_id", to.OppositeRegionId)
	d.Set("router_interface_to_region_id", from.OppositeRegionId)

	return nil
}

func resourceAlicloudRouterInterfaceConnectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parameters := strings.Split(d.Id(), COMMA_SEPARATED)

	modifyRequestFrom := vpc.CreateModifyRouterInterfaceAttributeRequest()
	modifyRequestFrom.RegionId = ""
	modifyRequestFrom.RouterInterfaceId = parameters[0]

	modifyRequestFrom.OppositeInterfaceId = ""
	modifyRequestFrom.OppositeRouterId = ""
	modifyRequestFrom.OppositeRouterType = ""
	modifyRequestFrom.OppositeInterfaceOwnerId = requests.Integer("")

	from, err := client.DescribeRouterInterface(parameters[0])
	if err == nil && "Active" == from.Status {

		if _, err := client.vpcconn.ModifyRouterInterfaceAttribute(modifyRequestFrom); err != nil {
			return fmt.Errorf("ModifyRouterInterfaceAttribute[ID = %s] got an error: %#v", parameters[0], err)
		}

		deactivateRouterInterfaceRequest := vpc.CreateDeactivateRouterInterfaceRequest()
		deactivateRouterInterfaceRequest.RouterInterfaceId = parameters[0]

		if _, err := client.vpcconn.DeactivateRouterInterface(deactivateRouterInterfaceRequest); err != nil {
			return fmt.Errorf("Error deactivate router interface %s: %#v", parameters[0], err)
		}
	}

	modifyRequestTo := vpc.CreateModifyRouterInterfaceAttributeRequest()
	modifyRequestTo.RegionId = ""
	modifyRequestTo.RouterInterfaceId = parameters[1]

	modifyRequestTo.OppositeInterfaceId = ""
	modifyRequestTo.OppositeRouterId = ""
	modifyRequestTo.OppositeRouterType = ""
	modifyRequestTo.OppositeInterfaceOwnerId = requests.Integer("")

	if _, err := client.vpcconn.ModifyRouterInterfaceAttribute(modifyRequestTo); err != nil {
		return fmt.Errorf("ModifyRouterInterfaceAttribute[ID = %s] got an error: %#v", parameters[1], err)
	}

	to, err := client.DescribeRouterInterface(parameters[1])
	if err == nil && "Active" == to.Status {
		deactivateRouterInterfaceRequest := vpc.CreateDeactivateRouterInterfaceRequest()
		deactivateRouterInterfaceRequest.RouterInterfaceId = parameters[1]

		if _, err := client.vpcconn.DeactivateRouterInterface(deactivateRouterInterfaceRequest); err != nil {
			return fmt.Errorf("Error deactivate router interface %s: %#v", parameters[1], err)
		}
	}

	return nil
}
