package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEipCreate,
		Read:   resourceAliyunEipRead,
		Update: resourceAliyunEipUpdate,
		Delete: resourceAliyunEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bandwidth": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"internet_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "PayByTraffic",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInternetChargeType,
			},

			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Computed values
			"alicloud_eip": {
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
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunEipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	request := vpc.CreateAllocateEipAddressRequest()
	request.RegionId = string(getRegion(d, meta))
	request.Bandwidth = strconv.Itoa(d.Get("bandwidth").(int))
	request.InternetChargeType = d.Get("internet_charge_type").(string)

	eip, err := client.vpcconn.AllocateEipAddress(request)
	if err != nil {
		if IsExceptedError(err, COMMODITYINVALID_COMPONENT) && request.InternetChargeType == string(PayByBandwidth) {
			return fmt.Errorf("Your account is international and it can only create '%s' elastic IP. Please change it and try again.", PayByTraffic)
		}
		return err
	}

	err = client.WaitForEip(eip.AllocationId, Available, 60)
	if err != nil {
		return fmt.Errorf("Error Waitting for EIP available: %#v", err)
	}

	d.SetId(eip.AllocationId)

	return resourceAliyunEipUpdate(d, meta)
}

func resourceAliyunEipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	eip, err := client.DescribeEipAddress(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe Eip Attribute: %#v", err)
	}

	// Output parameter 'instance' would be deprecated in the next version.
	if eip.InstanceId != "" {
		d.Set("instance", eip.InstanceId)
	} else {
		d.Set("instance", "")
	}

	bandwidth, _ := strconv.Atoi(eip.Bandwidth)
	d.Set("bandwidth", bandwidth)
	d.Set("internet_charge_type", eip.InternetChargeType)
	d.Set("ip_address", eip.IpAddress)
	d.Set("status", eip.Status)

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":                   d.Id(),
		"name":                 d.Id(),
		"status":               eip.Status,
		"creation_time":        time.Now().Format("2006-01-02 15:04:05"),
		"type":                 "alicloud_eip",
		"bandwidth":            bandwidth,
		"internet_charge_type": eip.InternetChargeType,
		"ip_address":           eip.IpAddress,
		"instance":             d.Get("instance").(string),
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_eip", s); err != nil {
		return fmt.Errorf("Setting alicloud_eip got an error: %#v.", err)
	}

	return nil
}

func resourceAliyunEipUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	if d.HasChange("bandwidth") && !d.IsNewResource() {
		request := vpc.CreateModifyEipAddressAttributeRequest()
		request.AllocationId = d.Id()
		request.Bandwidth = strconv.Itoa(d.Get("bandwidth").(int))
		if _, err := meta.(*AliyunClient).vpcconn.ModifyEipAddressAttribute(request); err != nil {
			return err
		}

		d.SetPartial("bandwidth")
	}

	d.Partial(false)

	return resourceAliyunEipRead(d, meta)
}

func resourceAliyunEipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	request := vpc.CreateReleaseEipAddressRequest()
	request.AllocationId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.vpcconn.ReleaseEipAddress(request); err != nil {
			if IsExceptedError(err, EipIncorrectStatus) {
				return resource.RetryableError(fmt.Errorf("Delete EIP timeout and got an error:%#v.", err))
			}
			return resource.NonRetryableError(err)

		}

		eip, descErr := client.DescribeEipAddress(d.Id())

		if descErr != nil {
			if NotFoundError(descErr) {
				return nil
			}
			return resource.NonRetryableError(descErr)
		} else if eip.AllocationId == d.Id() {
			return resource.RetryableError(fmt.Errorf("Delete EIP timeout and it still exists."))
		}
		return nil
	})
}
