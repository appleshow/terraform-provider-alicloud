package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudDBInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDBInstancesRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(MySQL),
					string(SQLServer),
					string(PPAS),
					string(PostgreSQL),
				}),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				// please follow the link below to see more details on available statusesplease follow the link below to see more details on available statuses
				// https://help.aliyun.com/document_detail/26315.html
			},
			"db_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Primary",
					"Readonly",
					"Guard",
					"Temp",
				}),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Standard",
					"Safe",
				}),
			},
			"tags": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateJsonString,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"alicloud_db_instances": {
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
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_upgrade_time_start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_upgrade_recovery_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_upgrade_recovery_max_iops": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_disk_used": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"advanced_features": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_cloud_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_max_quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"db_instance_cpu": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"increment_source_db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_upgrade_recovery_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ip_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_upgrade_account_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pay_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"support_create_super_account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"can_temp_upgrade": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"lock_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guard_db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ins_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"db_instance_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guard_db_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_upgrade_time_end": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_upgrade_recovery_memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"account_max_quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"temp_upgrade_recovery_max_connections": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_class_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"read_delay_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replicate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_upgrade_recovery_cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"availability_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDBInstancesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).rdsconn

	args := rds.CreateDescribeDBInstancesRequest()

	if instanceId, ok := d.GetOk("id"); ok {
		args.DBInstanceId = instanceId.(string)
	}
	args.RegionId = getRegionId(d, meta)
	args.Engine = d.Get("engine").(string)
	args.DBInstanceStatus = d.Get("status").(string)
	args.DBInstanceType = d.Get("db_type").(string)
	args.VpcId = d.Get("vpc_id").(string)
	args.VSwitchId = d.Get("vswitch_id").(string)
	args.ConnectionMode = d.Get("connection_mode").(string)
	args.Tags = d.Get("tags").(string)
	args.PageSize = requests.NewInteger(PageSizeLarge)

	var dbi []rds.DBInstance

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	for {
		resp, err := conn.DescribeDBInstances(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.Items.DBInstance) < 1 {
			break
		}

		for _, item := range resp.Items.DBInstance {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.DBInstanceDescription) {
					continue
				}
			}
			dbi = append(dbi, item)
		}

		if len(resp.Items.DBInstance) < PageSizeLarge {
			break
		}

		args.PageNumber = args.PageNumber + requests.NewInteger(1)
	}

	return rdsInstancesDescription(d, dbi, meta)
}

func rdsInstancesDescription(d *schema.ResourceData, dbi []rds.DBInstance, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}

	client := meta.(*AliyunClient)
	for _, item := range dbi {
		instance, err := client.DescribeDBInstanceById(item.DBInstanceId)
		if err != nil {
			return fmt.Errorf("Error Describe DB InstanceAttribute: %#v", err)
		}
		mapping := map[string]interface{}{
			"id":                                    instance.DBInstanceId,
			"name":                                  instance.DBInstanceDescription,
			"status":                                instance.DBInstanceStatus,
			"creation_time":                         instance.CreationTime,
			"resource_type":                         "alicloud_db_instance",
			"engine":                                instance.Engine,
			"temp_upgrade_time_start":               instance.TempUpgradeTimeStart,
			"temp_upgrade_recovery_time":            instance.TempUpgradeRecoveryTime,
			"temp_upgrade_recovery_max_iops":        instance.TempUpgradeRecoveryMaxIOPS,
			"db_instance_disk_used":                 instance.DBInstanceDiskUsed,
			"advanced_features":                     instance.AdvancedFeatures,
			"db_instance_class":                     instance.DBInstanceClass,
			"db_instance_net_type":                  instance.DBInstanceNetType,
			"vpc_cloud_instance_id":                 instance.VpcCloudInstanceId,
			"db_max_quantity":                       instance.DBMaxQuantity,
			"db_instance_cpu":                       instance.DBInstanceCPU,
			"max_connections":                       instance.MaxConnections,
			"increment_source_db_instance_id":       instance.IncrementSourceDBInstanceId,
			"instance_network_type":                 instance.InstanceNetworkType,
			"db_instance_type":                      instance.DBInstanceType,
			"temp_upgrade_recovery_class":           instance.TempUpgradeRecoveryClass,
			"db_instance_memory":                    instance.DBInstanceMemory,
			"vpc_id":                                instance.VpcId,
			"db_instance_storage_type":              instance.DBInstanceStorageType,
			"security_ip_list":                      instance.SecurityIPList,
			"support_upgrade_account_type":          instance.SupportUpgradeAccountType,
			"max_iops":                              instance.MaxIOPS,
			"tags":                                  instance.Tags,
			"engine_version":                        instance.EngineVersion,
			"maintain_time":                         instance.MaintainTime,
			"pay_type":                              instance.PayType,
			"db_instance_storage":                   instance.DBInstanceStorage,
			"support_create_super_account":          instance.SupportCreateSuperAccount,
			"temp_db_instance_id":                   instance.TempDBInstanceId,
			"zone_id":                               instance.ZoneId,
			"connection_mode":                       instance.ConnectionMode,
			"lock_mode":                             instance.LockMode,
			"can_temp_upgrade":                      instance.CanTempUpgrade,
			"lock_reason":                           instance.LockReason,
			"category":                              instance.Category,
			"guard_db_instance_id":                  instance.GuardDBInstanceId,
			"ins_id":                                instance.InsId,
			"db_instance_description":               instance.DBInstanceDescription,
			"account_type":                          instance.AccountType,
			"guard_db_instance_name":                instance.GuardDBInstanceName,
			"region_id":                             instance.RegionId,
			"resource_group_id":                     instance.ResourceGroupId,
			"temp_upgrade_time_end":                 instance.TempUpgradeTimeEnd,
			"expire_time":                           instance.ExpireTime,
			"temp_upgrade_recovery_memory":          instance.TempUpgradeRecoveryMemory,
			"account_max_quantity":                  instance.AccountMaxQuantity,
			"temp_upgrade_recovery_max_connections": instance.TempUpgradeRecoveryMaxConnections,
			"port":                      instance.Port,
			"vswitch_id":                instance.VSwitchId,
			"master_instance_id":        instance.MasterInstanceId,
			"db_instance_class_type":    instance.DBInstanceClassType,
			"read_delay_time":           instance.ReadDelayTime,
			"replicate_id":              instance.ReplicateId,
			"connection_string":         instance.ConnectionString,
			"temp_upgrade_recovery_cpu": instance.TempUpgradeRecoveryCpu,
			"availability_value":        instance.AvailabilityValue,
		}

		ids = append(ids, item.DBInstanceId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("alicloud_db_instances", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
