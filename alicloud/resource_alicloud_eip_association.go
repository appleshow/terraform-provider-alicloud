package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEipAssociationCreate,
		Read:   resourceAliyunEipAssociationRead,
		Delete: resourceAliyunEipAssociationDelete,

		Schema: map[string]*schema.Schema{
			"allocation_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			// Computed values
			"alicloud_eip_association": {
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
						"allocation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	args := vpc.CreateAssociateEipAddressRequest()
	args.AllocationId = Trim(d.Get("allocation_id").(string))
	args.InstanceId = Trim(d.Get("instance_id").(string))
	args.InstanceType = EcsInstance

	if strings.HasPrefix(args.InstanceId, "lb-") {
		args.InstanceType = SlbInstance
	}
	if strings.HasPrefix(args.InstanceId, "ngw-") {
		args.InstanceType = Nat
	}

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ar := args
		if _, err := client.vpcconn.AssociateEipAddress(ar); err != nil {
			if IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(fmt.Errorf("AssociateEip got an error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("AssociateEip got an error: %#v", err))
		}
		return nil
	}); err != nil {
		return err
	}

	if err := client.WaitForEip(args.AllocationId, InUse, 60); err != nil {
		return fmt.Errorf("Error Waitting for EIP allocated: %#v", err)
	}
	// There is at least 30 seconds delay for ecs instance
	if args.InstanceType == EcsInstance {
		time.Sleep(30 * time.Second)
	}

	d.SetId(args.AllocationId + ":" + args.InstanceId)

	return resourceAliyunEipAssociationRead(d, meta)
}

func resourceAliyunEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	allocationId, instanceId, err := getAllocationIdAndInstanceId(d, meta)
	if err != nil {
		return err
	}

	eip, err := client.DescribeEipAddress(allocationId)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe Eip Attribute: %#v", err)
	}

	if eip.InstanceId != instanceId {
		d.SetId("")
		return nil
	}

	d.Set("instance_id", eip.InstanceId)
	d.Set("allocation_id", allocationId)

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":            d.Id(),
		"name":          d.Id(),
		"status":        "Available",
		"creation_time": time.Now().Format("2006-01-02 15:04:05"),
		"type":          "alicloud_eip_association",
		"instance_id":   eip.InstanceId,
		"allocation_id": allocationId,
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_eip_association", s); err != nil {
		return fmt.Errorf("Setting alicloud_eip_association got an error: %#v.", err)
	}

	return nil
}

func resourceAliyunEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	allocationId, instanceId, err := getAllocationIdAndInstanceId(d, meta)
	if err != nil {
		return err
	}

	request := vpc.CreateUnassociateEipAddressRequest()
	request.AllocationId = allocationId
	request.InstanceId = instanceId
	request.InstanceType = EcsInstance

	if strings.HasPrefix(instanceId, "lb-") {
		request.InstanceType = SlbInstance
	}
	if strings.HasPrefix(instanceId, "ngw-") {
		request.InstanceType = Nat
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.vpcconn.UnassociateEipAddress(request); err != nil {
			if IsExceptedError(err, InstanceIncorrectStatus) ||
				IsExceptedError(err, HaVipIncorrectStatus) ||
				IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(fmt.Errorf("Unassociate EIP timeout and got an error:%#v.", err))
			}
		}

		eip, descErr := client.DescribeEipAddress(allocationId)
		if descErr != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(descErr)
		}

		if eip.InstanceId == instanceId {
			return resource.RetryableError(fmt.Errorf("Unassociate EIP timeout and got an error:%#v.", err))
		}

		return nil
	})
}

func getAllocationIdAndInstanceId(d *schema.ResourceData, meta interface{}) (string, string, error) {
	parts := strings.Split(d.Id(), ":")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}
