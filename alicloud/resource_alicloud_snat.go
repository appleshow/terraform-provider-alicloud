package alicloud

import (
	"fmt"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunSnatEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSnatEntryCreate,
		Read:   resourceAliyunSnatEntryRead,
		Update: resourceAliyunSnatEntryUpdate,
		Delete: resourceAliyunSnatEntryDelete,

		Schema: map[string]*schema.Schema{
			"snat_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snat_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// Computed values
			"alicloud_snat_entry": {
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
						"snat_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snat_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunSnatEntryCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	request := vpc.CreateCreateSnatEntryRequest()
	request.RegionId = string(getRegion(d, meta))
	request.SnatTableId = d.Get("snat_table_id").(string)
	request.SourceVSwitchId = d.Get("source_vswitch_id").(string)
	request.SnatIp = d.Get("snat_ip").(string)

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ar := request
		resp, err := conn.CreateSnatEntry(ar)
		if err != nil {
			if IsExceptedError(err, EIP_NOT_IN_GATEWAY) {
				return resource.RetryableError(fmt.Errorf("CreateSnatEntry timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("CreateSnatEntry got error: %#v.", err))
		}
		d.SetId(resp.SnatEntryId)
		return nil
	}); err != nil {
		return err
	}

	return resourceAliyunSnatEntryRead(d, meta)
}

func resourceAliyunSnatEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	snatEntry, err := client.DescribeSnatEntry(d.Get("snat_table_id").(string), d.Id())

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return err
	}

	d.Set("snat_table_id", snatEntry.SnatTableId)
	d.Set("source_vswitch_id", snatEntry.SourceVSwitchId)
	d.Set("snat_ip", snatEntry.SnatIp)

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":                d.Id(),
		"name":              d.Id(),
		"status":            "Available",
		"creation_time":     time.Now().Format("2006-01-02 15:04:05"),
		"type":              "alicloud_snat_entry",
		"snat_table_id":     snatEntry.SnatTableId,
		"source_vswitch_id": snatEntry.SourceVSwitchId,
		"snat_ip":           snatEntry.SnatIp,
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_snat_entry", s); err != nil {
		return fmt.Errorf("Setting alicloud_snat_entry got an error: %#v.", err)
	}

	return nil
}

func resourceAliyunSnatEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	snatEntry, err := client.DescribeSnatEntry(d.Get("snat_table_id").(string), d.Id())
	if err != nil {
		return err
	}

	d.Partial(true)
	attributeUpdate := false
	request := vpc.CreateModifySnatEntryRequest()
	request.RegionId = string(getRegion(d, meta))
	request.SnatTableId = snatEntry.SnatTableId
	request.SnatEntryId = snatEntry.SnatEntryId

	if d.HasChange("snat_ip") {
		d.SetPartial("snat_ip")
		var snat_ip string
		if v, ok := d.GetOk("snat_ip"); ok {
			snat_ip = v.(string)
		} else {
			return fmt.Errorf("cann't change snap_ip to empty string")
		}
		request.SnatIp = snat_ip

		attributeUpdate = true
	}

	if attributeUpdate {
		if _, err := client.vpcconn.ModifySnatEntry(request); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAliyunSnatEntryRead(d, meta)
}

func resourceAliyunSnatEntryDelete(d *schema.ResourceData, meta interface{}) error {

	request := vpc.CreateDeleteSnatEntryRequest()
	request.RegionId = string(getRegion(d, meta))
	request.SnatTableId = d.Get("snat_table_id").(string)
	request.SnatEntryId = d.Id()

	if _, err := meta.(*AliyunClient).vpcconn.DeleteSnatEntry(request); err != nil {
		if IsExceptedError(err, InvalidSnatTableIdNotFound) {
			return nil
		}
		return err
	}

	return nil
}
