package alicloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudSecurityGroupRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSecurityGroupRulesRead,

		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nic_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSecurityRuleNicType,
			},
			"direction": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSecurityRuleType,
			},
			"ip_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSecurityRuleIpProtocol,
			},
			"policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSecurityRulePolicy,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alicloud_security_group_rules": {
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
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_cidr_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_group_owner_account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_cidr_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_group_owner_account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSecurityGroupRulesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	attr, err := conn.DescribeSecurityGroupAttribute(
		&ecs.DescribeSecurityGroupAttributeArgs{
			SecurityGroupId: d.Get("security_group_id").(string),
			RegionId:        getRegion(d, meta),
			NicType:         ecs.NicType(d.Get("nic_type").(string)),
			Direction:       ecs.Direction(d.Get("direction").(string)),
		},
	)
	if err != nil {
		return fmt.Errorf("DescribeSecurityGroupAttribute: %#v", err)
	}

	var rules []map[string]interface{}

	log.Printf("alicloud_security_group_rules: total permission rules: %v", len(attr.Permissions.Permission))
	for _, item := range attr.Permissions.Permission {
		if v, ok := d.GetOk("ip_protocol"); ok && strings.ToLower(string(item.IpProtocol)) != v.(string) {
			continue
		}

		if v, ok := d.GetOk("policy"); ok && strings.ToLower(string(item.Policy)) != v.(string) {
			continue
		}

		mapping := map[string]interface{}{
			"id":                         attr.SecurityGroupId,
			"name":                       attr.SecurityGroupName,
			"status":                     "Available",
			"creation_time":              "",
			"resource_type":              "alicloud_security_group_rule",
			"group_desc":                 attr.Description,
			"ip_protocol":                strings.ToLower(string(item.IpProtocol)),
			"port_range":                 item.PortRange,
			"source_cidr_ip":             item.SourceCidrIp,
			"source_group_id":            item.SourceGroupId,
			"source_group_owner_account": item.SourceGroupOwnerAccount,
			"dest_cidr_ip":               item.DestCidrIp,
			"dest_group_id":              item.DestGroupId,
			"dest_group_owner_account":   item.DestGroupOwnerAccount,
			"policy":                     strings.ToLower(string(item.Policy)),
			"nic_type":                   item.NicType,
			"priority":                   item.Priority,
			"direction":                  item.Direction,
			"description":                item.Description,
		}

		log.Printf("alicloud_security_group_rules: adding permission rule: %v", mapping)
		rules = append(rules, mapping)
	}

	d.SetId(attr.SecurityGroupId)

	if err := d.Set("alicloud_security_group_rules", rules); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), rules)
	}
	return nil
}
