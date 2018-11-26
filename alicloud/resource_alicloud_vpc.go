package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpcCreate,
		Read:   resourceAliyunVpcRead,
		Update: resourceAliyunVpcUpdate,
		Delete: resourceAliyunVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cidr_block": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 128 {
						errors = append(errors, fmt.Errorf("%s cannot be longer than 128 characters", k))
					}

					if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
						errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
					}

					return
				},
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(2, 256),
			},
			"router_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"router_table_id": &schema.Schema{
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute router_table_id has been deprecated and replaced with route_table_id.",
			},
			"route_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			// Computed values
			"alicloud_vpc": {
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
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vrouter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
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
					},
				},
			},
		},
	}
}

func resourceAliyunVpcCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	var vpc *vpc.CreateVpcResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args, err := buildAliyunVpcArgs(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Building CreateVpcRequest got an error: %#v", err))
		}
		resp, err := client.vpcconn.CreateVpc(args)
		if err != nil {
			if IsExceptedError(err, VpcQuotaExceeded) {
				return resource.NonRetryableError(fmt.Errorf("The number of VPC has quota has reached the quota limit in your account, and please use existing VPCs or remove some of them."))
			}
			if IsExceptedError(err, TaskConflict) || IsExceptedError(err, UnknownError) {
				return resource.RetryableError(fmt.Errorf("Create vpc timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(err)
		}
		vpc = resp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Create vpc got an error :%#v", err)
	}

	d.SetId(vpc.VpcId)

	err = client.WaitForVpc(d.Id(), Available, 60)
	if err != nil {
		return fmt.Errorf("Timeout when WaitForVpcAvailable")
	}

	return resourceAliyunVpcUpdate(d, meta)
}

func resourceAliyunVpcRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	resp, err := client.DescribeVpc(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("cidr_block", resp.CidrBlock)
	d.Set("name", resp.VpcName)
	d.Set("description", resp.Description)
	d.Set("router_id", resp.VRouterId)
	request := vpc.CreateDescribeVRoutersRequest()
	request.RegionId = string(getRegion(d, meta))
	request.VRouterId = resp.VRouterId
	response, err := client.vpcconn.DescribeVRouters(request)
	if err != nil {
		return fmt.Errorf("DescribeVRouters got an error: %#v.", err)
	}
	if len(response.VRouters.VRouter) > 0 && len(response.VRouters.VRouter[0].RouteTableIds.RouteTableId) > 0 {
		d.Set("router_table_id", response.VRouters.VRouter[0].RouteTableIds.RouteTableId[0])
		d.Set("route_table_id", response.VRouters.VRouter[0].RouteTableIds.RouteTableId[0])
	} else {
		d.Set("router_table_id", "")
		d.Set("route_table_id", "")
	}

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":             d.Id(),
		"name":           resp.VpcName,
		"region_id":      string(getRegion(d, meta)),
		"status":         resp.Status,
		"vpc_name":       resp.VpcName,
		"vswitch_ids":    resp.VSwitchIds.VSwitchId,
		"cidr_block":     resp.CidrBlock,
		"vrouter_id":     resp.VRouterId,
		"route_table_id": d.Get("route_table_id").(string),
		"description":    resp.Description,
		"is_default":     resp.IsDefault,
		"creation_time":  time.Now().Format("2006-01-02 15:04:05"),
		"resource_type":  "alicloud_vpc",
	}

	s = append(s, mapping)

	if err := d.Set("alicloud_vpc", s); err != nil {
		return fmt.Errorf("Setting alicloud_vpc got an error: %#v.", err)
	}

	return nil
}

func resourceAliyunVpcUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	attributeUpdate := false
	request := vpc.CreateModifyVpcAttributeRequest()
	request.VpcId = d.Id()

	if d.HasChange("name") {
		d.SetPartial("name")
		request.VpcName = d.Get("name").(string)

		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		request.Description = d.Get("description").(string)

		attributeUpdate = true
	}

	if attributeUpdate {
		if _, err := meta.(*AliyunClient).vpcconn.ModifyVpcAttribute(request); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAliyunVpcRead(d, meta)
}

func resourceAliyunVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := vpc.CreateDeleteVpcRequest()
	request.VpcId = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.vpcconn.DeleteVpc(request)

		if err != nil {
			if IsExceptedError(err, InvalidVpcIDNotFound) || IsExceptedError(err, ForbiddenVpcNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete VPC timeout and got an error: %#v.", err))
		}

		if _, err := client.DescribeVpc(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func buildAliyunVpcArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateVpcRequest, error) {
	request := vpc.CreateCreateVpcRequest()
	request.RegionId = string(getRegion(d, meta))
	request.CidrBlock = d.Get("cidr_block").(string)

	if v := d.Get("name").(string); v != "" {
		request.VpcName = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	return request, nil
}
