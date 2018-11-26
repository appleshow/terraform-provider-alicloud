package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudFcService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFcServiceCreate,
		Read:   resourceAlicloudFcServiceRead,
		Update: resourceAlicloudFcServiceUpdate,
		Delete: resourceAlicloudFcServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_config": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			/* Argument vpcConfig is not supported"
			"vpc_config": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
			},
			*/
			/* Argument internetAccess is not supported
			"internet_access": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
			},
			*/

			// Computed values
			"alicloud_fc_service": {
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

func resourceAlicloudFcServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	serviceName := d.Get("service_name").(string)
	createServiceInput := fc.NewCreateServiceInput()

	createServiceInput.WithServiceName(serviceName)
	if description, ok := d.GetOk("description"); ok {
		createServiceInput.WithDescription(description.(string))
	}
	if role, ok := d.GetOk("role"); ok {
		createServiceInput.WithRole(role.(string))
	}
	if logConfig, ok := d.GetOk("log_config"); ok {
		fcLogConfig := fc.NewLogConfig()
		if project, ok := logConfig.(map[string]interface{})["project"]; ok {
			fcLogConfig.WithProject(project.(string))
		} else {
			return fmt.Errorf("Creating service of function compute got an error: %#v", "Parameter log_config missing item project")
		}
		if logstore, ok := logConfig.(map[string]interface{})["logstore"]; ok {
			fcLogConfig.WithLogstore(logstore.(string))
		} else {
			return fmt.Errorf("Creating service of function compute got an error: %#v", "Parameter log_config missing item logstore")
		}
		createServiceInput.WithLogConfig(fcLogConfig)
	}
	/* Argument vpcConfig is not supported"
	if vpcConfig, ok := d.GetOk("vpc_config"); ok {
		fcVpcConfig := fc.NewVPCConfig()
		if vpcId, ok := vpcConfig.(map[string]interface{})["vpc_id"]; ok {
			fcVpcConfig.WithVPCID(vpcId.(string))
		} else {
			return fmt.Errorf("Creating service of function compute got an error: %#v", "Parameter vpc_config missing item vpc_id")
		}
		if vSwitchIDs, ok := vpcConfig.(map[string]interface{})["vswitch_ids"]; ok {
			fcVpcConfig.WithVSwitchIDs(strings.Split(vSwitchIDs.(string), COMMA_SEPARATED))
		} else {
			return fmt.Errorf("Creating service of function compute got an error: %#v", "Parameter vpc_config missing item vswitch_ids")
		}
		if securityGroupID, ok := vpcConfig.(map[string]interface{})["security_group_id"]; ok {
			fcVpcConfig.WithSecurityGroupID(securityGroupID.(string))
		} else {
			return fmt.Errorf("Creating service of function compute got an error: %#v", "Parameter vpc_config missing item security_group_id")
		}
		createServiceInput.WithVPCConfig(fcVpcConfig)
	}
	*/
	/* Argument internetAccess is not supported
	if internetAccess, ok := d.GetOk("internet_access"); ok {
		createServiceInput.WithInternetAccess(internetAccess.(bool))
	} else {
		createServiceInput.WithInternetAccess(true)
	}
	*/

	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.fcconn.CreateService(createServiceInput)

		if err != nil {
			return resource.RetryableError(fmt.Errorf("Creating service of function compute got an error: %#v", err))
		} else {
			return nil
		}
	})

	if err != nil {
		return fmt.Errorf("Creating service of function compute got an error: %#v", err)
	}

	d.SetId(serviceName)

	return resourceAlicloudFcServiceUpdate(d, meta)
}

func resourceAlicloudFcServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	getServiceOutput, err := client.fcconn.GetService(fc.NewGetServiceInput(d.Id()))

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("service_name", getServiceOutput.ServiceName)
	d.Set("description", getServiceOutput.Description)
	d.Set("role", getServiceOutput.Role)

	logConfig := make(map[string]interface{})
	logConfig["project"] = getServiceOutput.LogConfig.Project
	logConfig["logstore"] = getServiceOutput.LogConfig.Logstore
	d.Set("log_config", logConfig)

	/* Argument vpcConfig is not supported"
	vpcConfig := make(map[string]interface{})
	vpcConfig["vpc_id"] = getServiceOutput.VPCConfig.VPCID
	vpcConfig["vswitch_ids"] = ""
	if len(getServiceOutput.VPCConfig.VSwitchIDs) > 0 {
		for _, v := range getServiceOutput.VPCConfig.VSwitchIDs {
			if vpcConfig["vswitch_ids"] != "" {
				vpcConfig["vswitch_ids"] = vpcConfig["vswitch_ids"].(string) + COMMA_SEPARATED + v
			} else {
				vpcConfig["vswitch_ids"] = v
			}
		}
	}
	vpcConfig["security_group_id"] = getServiceOutput.VPCConfig.SecurityGroupID
	d.Set("vpc_config", vpcConfig)
	*/

	/* Argument internetAccess is not supported
	d.Set("internet_access", getServiceOutput.InternetAccess)
	*/

	var s []map[string]interface{}

	var logConfigs []map[string]interface{}
	var vpcConfigs []map[string]interface{}

	mappingLogConfig := map[string]interface{}{
		"project":   *getServiceOutput.LogConfig.Project,
		"log_store": *getServiceOutput.LogConfig.Logstore,
	}
	logConfigs = append(logConfigs, mappingLogConfig)

	mapping := map[string]interface{}{
		"id":                 *getServiceOutput.ServiceID,
		"name":               *getServiceOutput.ServiceName,
		"status":             "Available",
		"creation_time":      *getServiceOutput.CreatedTime,
		"description":        *getServiceOutput.Description,
		"role":               *getServiceOutput.Role,
		"log_config":         logConfigs,
		"vpc_config":         vpcConfigs,
		"internet_access":    *getServiceOutput.InternetAccess,
		"last_modified_time": *getServiceOutput.LastModifiedTime,
		"resource_type":      "alicloud_fc_service",
	}

	s = append(s, mapping)

	if err := d.Set("alicloud_fc_service", s); err != nil {
		return fmt.Errorf("Setting alicloud_fc_service got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudFcServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	d.Partial(true)
	update := false

	updateServiceInput := fc.NewUpdateServiceInput(d.Id())

	if d.HasChange("description") {
		update = true
		updateServiceInput.WithDescription(d.Get("description").(string))
		d.SetPartial("description")
	}
	if d.HasChange("role") {
		update = true
		updateServiceInput.WithRole(d.Get("role").(string))
		d.SetPartial("role")
	}
	if d.HasChange("log_config") {
		update = true
		logConfig := d.Get("log_config")
		fcLogConfig := fc.NewLogConfig()

		if project, ok := logConfig.(map[string]interface{})["project"]; ok {
			fcLogConfig.WithProject(project.(string))
		} else {
			return fmt.Errorf("Creating service of function compute got an error: %#v", "Parameter log_config missing item project")
		}
		if logstore, ok := logConfig.(map[string]interface{})["logstore"]; ok {
			fcLogConfig.WithLogstore(logstore.(string))
		} else {
			return fmt.Errorf("Creating service of function compute got an error: %#v", "Parameter log_config missing item logstore")
		}
		updateServiceInput.WithLogConfig(fcLogConfig)

		d.SetPartial("log_config")
	}

	/* Argument vpcConfig is not supported"
	if d.HasChange("vpc_config") {
		update = true
		vpcConfig := d.Get("vpc_config")
		fcVpcConfig := fc.NewVPCConfig()

		if vpcId, ok := vpcConfig.(map[string]interface{})["vpc_id"]; ok {
			fcVpcConfig.WithVPCID(vpcId.(string))
		} else {
			return fmt.Errorf("Creating service of function compute got an error: %#v", "Parameter vpc_config missing item vpc_id")
		}
		if vSwitchIDs, ok := vpcConfig.(map[string]interface{})["vswitch_ids"]; ok {
			fcVpcConfig.WithVSwitchIDs(strings.Split(vSwitchIDs.(string), COMMA_SEPARATED))
		} else {
			return fmt.Errorf("Creating service of function compute got an error: %#v", "Parameter vpc_config missing item vswitch_ids")
		}
		if securityGroupID, ok := vpcConfig.(map[string]interface{})["security_group_id"]; ok {
			fcVpcConfig.WithSecurityGroupID(securityGroupID.(string))
		} else {
			return fmt.Errorf("Creating service of function compute got an error: %#v", "Parameter vpc_config missing item security_group_id")
		}
		updateServiceInput.WithVPCConfig(fcVpcConfig)
		d.SetPartial("vpc_config")
	}
	*/
	/* Argument internetAccess is not supported
	if d.HasChange("internet_access") {
		update = true
		updateServiceInput.WithInternetAccess(d.Get("internet_access").(bool))
		d.SetPartial("internet_access")
	}
	*/

	if !d.IsNewResource() && update {
		if _, err := client.fcconn.UpdateService(updateServiceInput); err != nil {
			return fmt.Errorf("Updating service of function compute got an error: %#v", err)
		}
	}

	d.Partial(false)

	return resourceAlicloudFcServiceRead(d, meta)
}

func resourceAlicloudFcServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.fcconn.DeleteService(fc.NewDeleteServiceInput(d.Id()))

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting service of function compute got an error: %#v", err))
		}

		resp, err := client.fcconn.GetService(fc.NewGetServiceInput(d.Id()))
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe service of function compute got an error: %#v", err))
		}
		if resp.ServiceName == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting service of function compute got an error: %#v", err))
	})
}
