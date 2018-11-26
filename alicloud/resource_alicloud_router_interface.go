package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRouterInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRouterInterfaceCreate,
		Read:   resourceAlicloudRouterInterfaceRead,
		Update: resourceAlicloudRouterInterfaceUpdate,
		Delete: resourceAlicloudRouterInterfaceDelete,

		Schema: map[string]*schema.Schema{
			"opposite_region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"router_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(VRouter), string(VBR)}),
				ForceNew: true,
			},
			"router_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(InitiatingSide), string(AcceptingSide)}),
				ForceNew: true,
			},
			"specification": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(GetAllRouterInterfaceSpec()),
			},
			"access_point_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRouterInterfaceDescription,
			},
			"health_check_source_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_target_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Computed values
			"alicloud_router_interface": {
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
						"opposite_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"router_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_point_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
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
					},
				},
			},
		},
	}
}

func resourceAlicloudRouterInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	args, err := buildAlicloudRouterInterfaceCreateArgs(d, meta)
	if err != nil {
		return err
	}

	response, err := client.vpcconn.CreateRouterInterface(args)
	if err != nil {
		return fmt.Errorf("CreateRouterInterface got an error: %#v", err)
	}

	d.SetId(response.RouterInterfaceId)

	if err := client.WaitForRouterInterface(d.Id(), Idle, 300); err != nil {
		return fmt.Errorf("WaitForRouterInterface %s got error: %#v", Idle, err)
	}

	return resourceAlicloudRouterInterfaceUpdate(d, meta)
}

func resourceAlicloudRouterInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	d.Partial(true)

	args, attributeUpdate, err := buildAlicloudRouterInterfaceModifyAttrArgs(d, meta)
	if err != nil {
		return err
	}

	if attributeUpdate {
		if _, err := conn.ModifyRouterInterfaceAttribute(args); err != nil {
			return fmt.Errorf("ModifyRouterInterfaceAttribute got an error: %#v", err)
		}
	}

	if d.HasChange("specification") && !d.IsNewResource() {
		d.SetPartial("specification")
		request := vpc.CreateModifyRouterInterfaceSpecRequest()
		request.RegionId = string(getRegion(d, meta))
		request.RouterInterfaceId = d.Id()
		request.Spec = d.Get("specification").(string)
		if _, err := conn.ModifyRouterInterfaceSpec(request); err != nil {
			return fmt.Errorf("ModifyRouterInterfaceSpec got an error: %#v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRouterInterfaceRead(d, meta)
}

func resourceAlicloudRouterInterfaceRead(d *schema.ResourceData, meta interface{}) error {

	ri, err := meta.(*AliyunClient).DescribeRouterInterface(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	d.Set("role", ri.Role)
	d.Set("specification", ri.Spec)
	d.Set("name", ri.Name)
	d.Set("router_id", ri.RouterId)
	d.Set("router_type", ri.RouterType)
	d.Set("description", ri.Description)
	d.Set("access_point_id", ri.AccessPointId)
	d.Set("health_check_source_ip", ri.HealthCheckSourceIp)
	d.Set("health_check_target_ip", ri.HealthCheckTargetIp)
	d.Set("status", ri.Status)

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":                     d.Id(),
		"name":                   ri.Name,
		"status":                 ri.Status,
		"creation_time":          time.Now().Format("2006-01-02 15:04:05"),
		"type":                   "alicloud_router_interface",
		"opposite_region":        d.Get("opposite_region").(string),
		"router_type":            ri.RouterType,
		"router_id":              ri.RouterId,
		"role":                   ri.Role,
		"specification":          ri.Spec,
		"access_point_id":        ri.AccessPointId,
		"description":            ri.Description,
		"health_check_source_ip": ri.HealthCheckSourceIp,
		"health_check_target_ip": ri.HealthCheckTargetIp,
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_router_interface", s); err != nil {
		return fmt.Errorf("Setting alicloud_router_interface got an error: %#v.", err)
	}

	return nil

}

func resourceAlicloudRouterInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	if ri, err := meta.(*AliyunClient).DescribeRouterInterface(d.Id()); err == nil && "Active" == ri.Status {
		deactivateRouterInterfaceRequest := vpc.CreateDeactivateRouterInterfaceRequest()
		deactivateRouterInterfaceRequest.RouterInterfaceId = d.Id()

		if _, err := conn.DeactivateRouterInterface(deactivateRouterInterfaceRequest); err != nil {
			return fmt.Errorf("Error deactivate router interface %s: %#v", d.Id(), err)
		}
	}

	args := vpc.CreateDeleteRouterInterfaceRequest()
	args.RegionId = string(getRegion(d, meta))
	args.RouterInterfaceId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DeleteRouterInterface(args); err != nil {
			if IsExceptedError(err, RouterInterfaceIncorrectStatus) || IsExceptedError(err, DependencyViolationRouterInterfaceReferedByRouteEntry) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(fmt.Errorf("Delete router interface timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting interface %s: %#v", d.Id(), err))
		}
		return nil
	})
}

func buildAlicloudRouterInterfaceCreateArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateRouterInterfaceRequest, error) {
	client := meta.(*AliyunClient)

	oppositeRegion := common.Region(d.Get("opposite_region").(string))
	if err := client.JudgeRegionValidation("opposite_region", oppositeRegion); err != nil {
		return nil, err
	}

	request := vpc.CreateCreateRouterInterfaceRequest()
	request.RegionId = string(getRegion(d, meta))
	request.RouterType = d.Get("router_type").(string)
	request.RouterId = d.Get("router_id").(string)
	request.Role = d.Get("role").(string)
	request.Spec = d.Get("specification").(string)
	request.OppositeRegionId = string(oppositeRegion)

	if request.RouterType == string(VBR) {
		if request.Role != string(InitiatingSide) {
			return nil, fmt.Errorf("'role': valid value is only 'InitiatingSide' when 'router_type' is 'VBR'.")
		}

		v, ok := d.GetOk("access_point_id")
		if !ok {
			return nil, fmt.Errorf("'access_point_id': required field is not set when 'router_type' is 'VBR'.")
		}
		request.AccessPointId = v.(string)
	}

	if request.Role == string(AcceptingSide) {
		if request.Spec == "" {
			request.Spec = string(Negative)
		} else if request.Spec != string(Negative) {
			return nil, fmt.Errorf("'specification': valid value is only '%s' when 'role' is 'AcceptingSide'.", Negative)
		}
	} else if oppositeRegion == getRegion(d, meta) {
		if request.RouterType == string(VRouter) {
			if request.Spec != string(Large2) {
				return nil, fmt.Errorf("'specification': valid value is only '%s' when 'role' is 'InitiatingSide' and 'region' is equal to 'opposite_region' and 'router_type' is 'VRouter'.", Large2)
			}
		} else {
			if request.Spec != string(Middle1) && request.Spec != string(Middle2) && request.Spec != string(Middle5) && request.Spec != string(Large1) {
				return nil, fmt.Errorf("'specification': valid values are '%s', '%s', '%s' and '%s' when 'role' is 'InitiatingSide' and 'region' is equal to 'opposite_region' and 'router_type' is 'VBR'.", Large1, Middle1, Middle2, Middle5)
			}
		}
	} else if request.Spec == string(Large2) {
		return nil, fmt.Errorf("The 'specification' can not be '%s' when 'role' is 'InitiatingSide' and 'region' is not equal to 'opposite_region'.", Large2)
	}

	return request, nil
}

func buildAlicloudRouterInterfaceModifyAttrArgs(d *schema.ResourceData, meta interface{}) (*vpc.ModifyRouterInterfaceAttributeRequest, bool, error) {

	sourceIp, sourceOk := d.GetOk("health_check_source_ip")
	targetIp, targetOk := d.GetOk("health_check_target_ip")
	if sourceOk && !targetOk || !sourceOk && targetOk {
		return nil, false, fmt.Errorf("The 'health_check_source_ip' and 'health_check_target_ip' should be specified or not at one time.")
	}

	args := vpc.CreateModifyRouterInterfaceAttributeRequest()
	args.RegionId = string(getRegion(d, meta))
	args.RouterInterfaceId = d.Id()

	attributeUpdate := false

	if d.HasChange("health_check_source_ip") {
		d.SetPartial("health_check_source_ip")
		args.HealthCheckSourceIp = sourceIp.(string)
		args.HealthCheckTargetIp = targetIp.(string)
		attributeUpdate = true
	}

	if d.HasChange("health_check_target_ip") {
		d.SetPartial("health_check_target_ip")
		args.HealthCheckTargetIp = targetIp.(string)
		args.HealthCheckSourceIp = sourceIp.(string)
		attributeUpdate = true
	}

	if d.HasChange("name") {
		d.SetPartial("name")
		args.Name = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		args.Description = d.Get("description").(string)
		attributeUpdate = true
	}

	return args, attributeUpdate, nil
}
