package alicloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudFcServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFcServicesRead,

		Schema: map[string]*schema.Schema{
			"prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"next_token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_fc_services": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"log_store": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"vpc_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vswitch_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"internet_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"last_modified_time": {
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

func dataSourceAlicloudFcServicesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	listServicesInput := fc.NewListServicesInput()
	if prefix, ok := d.GetOk("prefix"); ok {
		listServicesInput.WithPrefix(prefix.(string))
	}
	if startKey, ok := d.GetOk("start_key"); ok {
		listServicesInput.WithStartKey(startKey.(string))
	}
	if nextToken, ok := d.GetOk("next_token"); ok {
		listServicesInput.WithNextToken(nextToken.(string))
	}
	if limit, ok := d.GetOk("limit"); ok {
		if limit, err := strconv.ParseInt(strconv.Itoa(limit.(int)), 10, 32); err != nil {
			return fmt.Errorf("List services of function compute got an error: %#v", err)
		} else {
			listServicesInput.WithLimit(int32(limit))
		}
	}

	listServicesOutput, err := client.fcconn.ListServices(listServicesInput)
	if err != nil {
		return fmt.Errorf("List services of function compute got an error: %#v", err)
	} else {
		log.Printf("[DEBUG] alicloud_fc - Services found: %#v", listServicesOutput.Services)

		var ids []string
		var s []map[string]interface{}
		for _, service := range listServicesOutput.Services {
			var logConfig []map[string]interface{}
			var vpcConfig []map[string]interface{}

			mappingLogConfig := map[string]interface{}{
				"project":   *service.LogConfig.Project,
				"log_store": *service.LogConfig.Logstore,
			}
			logConfig = append(logConfig, mappingLogConfig)

			mapping := map[string]interface{}{
				"id":                 *service.ServiceID,
				"name":               *service.ServiceName,
				"status":             "Available",
				"creation_time":      *service.CreatedTime,
				"description":        *service.Description,
				"role":               *service.Role,
				"log_config":         logConfig,
				"vpc_config":         vpcConfig,
				"internet_access":    *service.InternetAccess,
				"last_modified_time": *service.LastModifiedTime,
				"resource_type":      "alicloud_fc_service",
			}
			ids = append(ids, *service.ServiceID)
			s = append(s, mapping)
		}

		d.SetId(dataResourceIdHash(ids))
		if err := d.Set("alicloud_fc_services", s); err != nil {
			return err
		}

		// create a json file in current directory and write data source to it.
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), s)
		}
	}

	return nil
}
