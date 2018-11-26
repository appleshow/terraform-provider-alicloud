package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSecurityGroupCreate,
		Read:   resourceAliyunSecurityGroupRead,
		Update: resourceAliyunSecurityGroupUpdate,
		Delete: resourceAliyunSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSecurityGroupName,
			},

			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSecurityGroupDescription,
			},

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"inner_access": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Computed values
			"alicloud_security_group": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inner_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	args, err := buildAliyunSecurityGroupArgs(d, meta)
	if err != nil {
		return err
	}

	securityGroupID, err := conn.CreateSecurityGroup(args)
	if err != nil {
		return err
	}

	d.SetId(securityGroupID)
	return resourceAliyunSecurityGroupUpdate(d, meta)
}

func resourceAliyunSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	args := &ecs.DescribeSecurityGroupAttributeArgs{
		SecurityGroupId: d.Id(),
		RegionId:        getRegion(d, meta),
	}
	//err := resource.Retry(3*time.Minute, func() *resource.RetryError {
	var sg *ecs.DescribeSecurityGroupAttributeResponse
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		group, e := conn.DescribeSecurityGroupAttribute(args)
		if e != nil {
			if IsExceptedError(e, InvalidSecurityGroupIdNotFound) {
				sg = nil
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error DescribeSecurityGroupAttribute: %#v", e))
		}
		if group != nil {
			sg = group
			return nil
		}
		return resource.RetryableError(fmt.Errorf("Create security group timeout and got an error: %#v", e))
	})

	if err != nil {
		return err
	}
	if sg == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", sg.SecurityGroupName)
	d.Set("description", sg.Description)
	d.Set("vpc_id", sg.VpcId)
	d.Set("inner_access", sg.InnerAccessPolicy == ecs.GroupInnerAccept)
	d.Set("status", "Available")
	d.Set("creation_time", "")

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":            d.Id(),
		"name":          sg.SecurityGroupName,
		"status":        "Available",
		"creation_time": time.Now().Format("2006-01-02 15:04:05"),
		"type":          "alicloud_security_group",
		"description":   sg.Description,
		"vpc_id":        sg.VpcId,
		"inner_access":  sg.InnerAccessPolicy == ecs.GroupInnerAccept,
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_security_group", s); err != nil {
		return fmt.Errorf("Setting alicloud_security_group got an error: %#v.", err)
	}

	return nil
}

func resourceAliyunSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*AliyunClient).ecsconn

	d.Partial(true)
	attributeUpdate := false
	args := &ecs.ModifySecurityGroupAttributeArgs{
		SecurityGroupId: d.Id(),
		RegionId:        getRegion(d, meta),
	}

	if d.HasChange("name") && !d.IsNewResource() {
		d.SetPartial("name")
		args.SecurityGroupName = d.Get("name").(string)

		attributeUpdate = true
	}

	if d.HasChange("description") && !d.IsNewResource() {
		d.SetPartial("description")
		args.Description = d.Get("description").(string)

		attributeUpdate = true
	}
	if attributeUpdate {
		if err := conn.ModifySecurityGroupAttribute(args); err != nil {
			return err
		}
	}

	if d.HasChange("inner_access") {
		policy := ecs.GroupInnerAccept
		if !d.Get("inner_access").(bool) {
			policy = ecs.GroupInnerDrop
		}
		if err := conn.ModifySecurityGroupPolicy(&ecs.ModifySecurityGroupPolicyArgs{
			RegionId:          getRegion(d, meta),
			SecurityGroupId:   d.Id(),
			InnerAccessPolicy: policy,
		}); err != nil {
			return fmt.Errorf("ModifySecurityGroupPolicy got an error: %#v.", err)
		}

	}

	d.Partial(false)

	return resourceAliyunSecurityGroupRead(d, meta)
}

func resourceAliyunSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*AliyunClient).ecsconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.DeleteSecurityGroup(getRegion(d, meta), d.Id())

		if err != nil {
			if IsExceptedError(err, SgDependencyViolation) {
				return resource.RetryableError(fmt.Errorf("Delete security group timeout and got an error: %#v", err))
			}
		}

		sg, err := conn.DescribeSecurityGroupAttribute(&ecs.DescribeSecurityGroupAttributeArgs{
			RegionId:        getRegion(d, meta),
			SecurityGroupId: d.Id(),
		})

		if err != nil {
			if IsExceptedError(err, InvalidSecurityGroupIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		} else if sg == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Delete security group timeout and got an error: %#v", err))
	})

}

func buildAliyunSecurityGroupArgs(d *schema.ResourceData, meta interface{}) (*ecs.CreateSecurityGroupArgs, error) {

	args := &ecs.CreateSecurityGroupArgs{
		RegionId: getRegion(d, meta),
	}

	if v := d.Get("name").(string); v != "" {
		args.SecurityGroupName = v
	}

	if v := d.Get("description").(string); v != "" {
		args.Description = v
	}

	if v := d.Get("vpc_id").(string); v != "" {
		args.VpcId = v
	}

	return args, nil
}
