package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudRamRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamRolesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamPolicyName,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyType,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"alicloud_ram_roles": {
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
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assume_role_policy_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_date": {
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

func dataSourceAlicloudRamRolesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	allRoles := []interface{}{}

	allRolesMap := make(map[string]interface{})
	policyFilterRolesMap := make(map[string]interface{})

	dataMap := []map[string]interface{}{}

	policyName, policyNameOk := d.GetOk("policy_name")
	policyType, policyTypeOk := d.GetOk("policy_type")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if policyTypeOk && !policyNameOk {
		return fmt.Errorf("You must set 'policy_name' at one time when you set 'policy_type'.")
	}

	// all roles
	resp, err := conn.ListRoles()
	if err != nil {
		return fmt.Errorf("ListRoles got an error: %#v", err)
	}
	for _, v := range resp.Roles.Role {
		if nameRegexOk {
			r := regexp.MustCompile(nameRegex.(string))
			if !r.MatchString(v.RoleName) {
				continue
			}
		}
		allRolesMap[v.RoleName] = v
	}

	// roles which attach with this policy
	if policyNameOk {
		pType := ram.System
		if policyTypeOk {
			pType = ram.Type(policyType.(string))
		}
		resp, err := conn.ListEntitiesForPolicy(ram.PolicyRequest{PolicyName: policyName.(string), PolicyType: pType})
		if err != nil {
			return fmt.Errorf("ListEntitiesForPolicy got an error: %#v", err)
		}

		for _, v := range resp.Roles.Role {
			policyFilterRolesMap[v.RoleName] = v
		}
		dataMap = append(dataMap, policyFilterRolesMap)
	}

	// GetIntersection of each map
	allRoles = GetIntersection(dataMap, allRolesMap)

	if len(allRoles) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_ram_roles - Roles found: %#v", allRoles)

	return ramRolesDescriptionAttributes(d, meta, allRoles)
}

func ramRolesDescriptionAttributes(d *schema.ResourceData, meta interface{}, roles []interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, v := range roles {
		role := v.(ram.Role)
		conn := meta.(*AliyunClient).ramconn
		resp, _ := conn.GetRole(ram.RoleQueryRequest{RoleName: role.RoleName})
		mapping := map[string]interface{}{
			"id":                          role.RoleId,
			"name":                        role.RoleName,
			"arn":                         role.Arn,
			"description":                 role.Description,
			"create_date":                 role.CreateDate,
			"update_date":                 role.UpdateDate,
			"assume_role_policy_document": resp.Role.AssumeRolePolicyDocument,
			"document":                    resp.Role.AssumeRolePolicyDocument,
			"resource_type":               "alicloud_ram_role",
		}
		log.Printf("[DEBUG] alicloud_ram_roles - adding role: %v", mapping)
		ids = append(ids, role.RoleId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alicloud_ram_roles", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
