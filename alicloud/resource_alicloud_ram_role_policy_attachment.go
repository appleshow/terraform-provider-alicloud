package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRamRolePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamRolePolicyAttachmentCreate,
		Read:   resourceAlicloudRamRolePolicyAttachmentRead,
		//Update: resourceAlicloudRamRolePolicyAttachmentUpdate,
		Delete: resourceAlicloudRamRolePolicyAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"role_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
			},
			"policy_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamPolicyName,
			},
			"policy_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyType,
			},
			// Computed values
			"alicloud_ram_role_policy_attachment": {
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
						"role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudRamRolePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	args := ram.AttachPolicyToRoleRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: ram.Type(d.Get("policy_type").(string)),
		},
		RoleName: d.Get("role_name").(string),
	}

	if _, err := conn.AttachPolicyToRole(args); err != nil {
		return fmt.Errorf("AttachPolicyToRole got an error: %#v", err)
	}
	d.SetId("role" + args.PolicyName + string(args.PolicyType) + args.RoleName)

	return resourceAlicloudRamRolePolicyAttachmentRead(d, meta)
}

func resourceAlicloudRamRolePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.RoleQueryRequest{
		RoleName: d.Get("role_name").(string),
	}

	response, err := conn.ListPoliciesForRole(args)
	if err != nil {
		return fmt.Errorf("Get list policies for role got an error: %v", err)
	}

	if len(response.Policies.Policy) > 0 {
		for _, v := range response.Policies.Policy {
			if v.PolicyName == d.Get("policy_name").(string) && v.PolicyType == d.Get("policy_type").(string) {
				d.Set("role_name", args.RoleName)
				d.Set("policy_name", v.PolicyName)
				d.Set("policy_type", v.PolicyType)

				var s []map[string]interface{}
				mapping := map[string]interface{}{
					"id":            d.Id(),
					"name":          d.Id(),
					"status":        "Available",
					"creation_time": "",
					"type":          "alicloud_ram_role_policy_attachment",
					"role_name":     args.RoleName,
					"policy_name":   v.PolicyName,
					"policy_type":   v.PolicyType,
				}
				s = append(s, mapping)

				if err := d.Set("alicloud_ram_role_policy_attachment", s); err != nil {
					return fmt.Errorf("Setting alicloud_ram_role_policy_attachment got an error: %#v.", err)
				}

				return nil
			}
		}
	}

	d.SetId("")
	return nil
}

func resourceAlicloudRamRolePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.AttachPolicyToRoleRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: ram.Type(d.Get("policy_type").(string)),
		},
		RoleName: d.Get("role_name").(string),
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DetachPolicyFromRole(args); err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting role policy attachment: %#v", err))
		}

		response, err := conn.ListPoliciesForRole(ram.RoleQueryRequest{RoleName: args.RoleName})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if len(response.Policies.Policy) < 1 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("Error deleting role policy attachment - trying again while it is deleted."))
	})
}
