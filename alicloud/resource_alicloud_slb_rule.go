package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunSlbRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbRuleCreate,
		Read:   resourceAliyunSlbRuleRead,
		Update: resourceAliyunSlbRuleUpdate,
		Delete: resourceAliyunSlbRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"frontend_port": &schema.Schema{
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerInRange(1, 65535),
				Required:     true,
				ForceNew:     true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "tf-slb-rule",
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// Computed values
			"alicloud_slb_rule": {
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
						"load_balancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frontend_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunSlbRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	slb_id := d.Get("load_balancer_id").(string)
	port := d.Get("frontend_port").(int)
	name := strings.Trim(d.Get("name").(string), " ")
	group_id := strings.Trim(d.Get("server_group_id").(string), " ")

	var domain, url, rule string
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}
	if v, ok := d.GetOk("url"); ok {
		url = v.(string)
	}

	if domain == "" && url == "" {
		return fmt.Errorf("At least one 'domain' or 'url' must be set.")
	} else if domain == "" {
		rule = fmt.Sprintf("[{'RuleName':'%s','Url':'%s','VServerGroupId':'%s'}]", name, url, group_id)
	} else if url == "" {
		rule = fmt.Sprintf("[{'RuleName':'%s','Domain':'%s','VServerGroupId':'%s'}]", name, domain, group_id)
	} else {
		rule = fmt.Sprintf("[{'RuleName':'%s','Domain':'%s','Url':'%s','VServerGroupId':'%s'}]", name, domain, url, group_id)
	}

	if err := client.slbconn.CreateRules(&slb.CreateRulesArgs{
		RegionId:       getRegion(d, meta),
		LoadBalancerId: slb_id,
		ListenerPort:   port,
		RuleList:       rule,
	}); err != nil {
		if IsExceptedError(err, RuleDomainExist) {
			if ruleId, err := client.DescribeLoadBalancerRuleId(slb_id, port, domain, url); err != nil {
				return err
			} else {
				return fmt.Errorf("The rule with same domain and url already exists. "+
					"Please import it using ID '%s' to import it or specify a different 'domain' or 'url' and then try again.", ruleId)
			}
		}
		return fmt.Errorf("CreateRule got an error: %#v", err)
	}

	ruleId, err := client.DescribeLoadBalancerRuleId(slb_id, port, domain, url)
	if err != nil {
		return err
	}

	if ruleId == "" {
		return fmt.Errorf("There is not found any rules in the load balancer %s and listener port %d.", slb_id, port)
	}

	d.SetId(ruleId)

	return resourceAliyunSlbRuleUpdate(d, meta)
}

func resourceAliyunSlbRuleRead(d *schema.ResourceData, meta interface{}) error {

	rule, err := meta.(*AliyunClient).slbconn.DescribeRuleAttribute(&slb.DescribeRuleAttributeArgs{
		RegionId: getRegion(d, meta),
		RuleId:   d.Id(),
	})

	if err != nil {
		if IsExceptedError(err, InvalidRuleIdNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("DescribeRuleAttribute got an error: %#v", err)
	}

	d.Set("name", rule.RuleName)
	d.Set("load_balancer_id", rule.LoadBalancerId)
	d.Set("frontend_port", rule.ListenerPort)
	d.Set("domain", rule.Domain)
	d.Set("url", rule.Url)
	d.Set("server_group_id", rule.VServerGroupId)

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":               d.Id(),
		"name":             rule.RuleName,
		"status":           "Available",
		"creation_time":    time.Now().Format("2006-01-02 15:04:05"),
		"type":             "alicloud_slb_rule",
		"load_balancer_id": rule.LoadBalancerId,
		"frontend_port":    rule.ListenerPort,
		"domain":           rule.Domain,
		"url":              rule.Url,
		"server_group_id":  rule.VServerGroupId,
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_slb_rule", s); err != nil {
		return fmt.Errorf("Setting alicloud_slb_rule got an error: %#v.", err)
	}

	return nil
}

func resourceAliyunSlbRuleUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	if d.HasChange("server_group_id") && !d.IsNewResource() {
		if err := meta.(*AliyunClient).slbconn.SetRule(&slb.SetRuleArgs{
			RegionId:       getRegion(d, meta),
			RuleId:         d.Id(),
			VServerGroupId: d.Get("server_group_id").(string),
		}); err != nil {
			return fmt.Errorf("Modify rule %s server group got an error: %#v", d.Id(), err)
		}
		d.SetPartial("server_group_id")
	}

	d.Partial(false)

	return resourceAliyunSlbRuleRead(d, meta)
}

func resourceAliyunSlbRuleDelete(d *schema.ResourceData, meta interface{}) error {
	slbconn := meta.(*AliyunClient).slbconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if err := slbconn.DeleteRules(&slb.DeleteRulesArgs{
			RegionId: getRegion(d, meta),
			RuleIds:  fmt.Sprintf("['%s']", d.Id()),
		}); err != nil {
			if IsExceptedError(err, InvalidRuleIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		rule, err := meta.(*AliyunClient).slbconn.DescribeRuleAttribute(&slb.DescribeRuleAttributeArgs{
			RegionId: getRegion(d, meta),
			RuleId:   d.Id(),
		})

		if err != nil {
			if IsExceptedError(err, InvalidRuleIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("While deleting rule, DescribeRuleAttribute got an error: %#v", err))
		}

		if rule != nil {
			return resource.RetryableError(fmt.Errorf("DeleteRule got an error: %#v", err))
		}
		return nil
	})
}
