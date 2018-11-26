package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudImageSharePermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudImageSharePermissionCreate,
		Read:   resourceAlicloudImageSharePermissionRead,
		Update: resourceAlicloudImageSharePermissionUpdate,
		Delete: resourceAlicloudImageSharePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"accounts": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
				MaxItems: 10,
			},
			// Computed values
			"alicloud_image_share_permission": {
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
						"image_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"accounts": &schema.Schema{
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudImageSharePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	modifyImageSharePermissionRequest := ecs.CreateModifyImageSharePermissionRequest()
	modifyImageSharePermissionRequest.ImageId = d.Get("image_id").(string)
	accounts := d.Get("accounts").([]interface{})

	modifyImageSharePermissionRequest.AddAccount1 = accounts[0].(string)
	if len(accounts) >= 2 {
		modifyImageSharePermissionRequest.AddAccount2 = accounts[1].(string)
	}
	if len(accounts) >= 3 {
		modifyImageSharePermissionRequest.AddAccount3 = accounts[2].(string)
	}
	if len(accounts) >= 4 {
		modifyImageSharePermissionRequest.AddAccount4 = accounts[3].(string)
	}
	if len(accounts) >= 5 {
		modifyImageSharePermissionRequest.AddAccount5 = accounts[4].(string)
	}
	if len(accounts) >= 6 {
		modifyImageSharePermissionRequest.AddAccount6 = accounts[5].(string)
	}
	if len(accounts) >= 7 {
		modifyImageSharePermissionRequest.AddAccount7 = accounts[6].(string)
	}
	if len(accounts) >= 8 {
		modifyImageSharePermissionRequest.AddAccount8 = accounts[7].(string)
	}
	if len(accounts) >= 9 {
		modifyImageSharePermissionRequest.AddAccount9 = accounts[7].(string)
	}
	if len(accounts) >= 10 {
		modifyImageSharePermissionRequest.AddAccount10 = accounts[9].(string)
	}

	_, err := client.aliecsconn.ModifyImageSharePermission(modifyImageSharePermissionRequest)
	if err != nil {
		return fmt.Errorf("Creating Image Share Permission got an error: %#v.", err)
	}

	d.SetId(modifyImageSharePermissionRequest.ImageId)

	return resourceAlicloudImageSharePermissionUpdate(d, meta)
}

func resourceAlicloudImageSharePermissionRead(d *schema.ResourceData, meta interface{}) error {
	var accounts_list []string
	accounts := d.Get("accounts").([]interface{})

	for _, account := range accounts {
		accounts_list = append(accounts_list, account.(string))
	}

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":            d.Id(),
		"name":          d.Id(),
		"status":        "Available",
		"creation_time": time.Now().Format("2006-01-02 15:04:05"),
		"resource_type": "alicloud_image_share_permission",
		"image_id":      d.Get("image_id").(string),
		"accounts":      accounts_list,
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_image_share_permission", s); err != nil {
		return fmt.Errorf("Setting alicloud_image_share_permission got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudImageSharePermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	d.Partial(true)
	update := false

	if d.HasChange("image_id") && !d.IsNewResource() {
		return fmt.Errorf("Updating Image Share Permission got an error: %#v", "Cannot modify parameter image_id")
	}
	if d.HasChange("accounts") {
		update = true
		d.SetPartial("accounts")
	}

	if !d.IsNewResource() && update {
		accountOld, accountNew := d.GetChange("accounts")
		modifyImageSharePermissionRequest := ecs.CreateModifyImageSharePermissionRequest()
		modifyImageSharePermissionRequest.ImageId = d.Get("image_id").(string)

		var index = 0
		for _, n := range accountNew.([]interface{}) {
			var find = false
			for _, o := range accountOld.([]interface{}) {
				if n.(string) == o.(string) {
					find = true
				}
			}
			if !find {
				index = index + 1
				if index == 1 {
					modifyImageSharePermissionRequest.AddAccount1 = n.(string)
				}
				if index == 2 {
					modifyImageSharePermissionRequest.AddAccount2 = n.(string)
				}
				if index == 3 {
					modifyImageSharePermissionRequest.AddAccount3 = n.(string)
				}
				if index == 4 {
					modifyImageSharePermissionRequest.AddAccount4 = n.(string)
				}
				if index == 5 {
					modifyImageSharePermissionRequest.AddAccount5 = n.(string)
				}
				if index == 6 {
					modifyImageSharePermissionRequest.AddAccount6 = n.(string)
				}
				if index == 7 {
					modifyImageSharePermissionRequest.AddAccount7 = n.(string)
				}
				if index == 8 {
					modifyImageSharePermissionRequest.AddAccount8 = n.(string)
				}
				if index == 9 {
					modifyImageSharePermissionRequest.AddAccount9 = n.(string)
				}
				if index == 10 {
					modifyImageSharePermissionRequest.AddAccount10 = n.(string)
				}
			}
		}

		index = 0
		for _, o := range accountOld.([]interface{}) {
			var find = false
			for _, n := range accountNew.([]interface{}) {
				if n.(string) == o.(string) {
					find = true
				}
			}
			if !find {
				index = index + 1
				if index == 1 {
					modifyImageSharePermissionRequest.RemoveAccount1 = o.(string)
				}
				if index == 2 {
					modifyImageSharePermissionRequest.RemoveAccount2 = o.(string)
				}
				if index == 3 {
					modifyImageSharePermissionRequest.RemoveAccount3 = o.(string)
				}
				if index == 4 {
					modifyImageSharePermissionRequest.RemoveAccount4 = o.(string)
				}
				if index == 5 {
					modifyImageSharePermissionRequest.RemoveAccount5 = o.(string)
				}
				if index == 6 {
					modifyImageSharePermissionRequest.RemoveAccount6 = o.(string)
				}
				if index == 7 {
					modifyImageSharePermissionRequest.RemoveAccount7 = o.(string)
				}
				if index == 8 {
					modifyImageSharePermissionRequest.RemoveAccount8 = o.(string)
				}
				if index == 9 {
					modifyImageSharePermissionRequest.RemoveAccount9 = o.(string)
				}
				if index == 10 {
					modifyImageSharePermissionRequest.RemoveAccount10 = o.(string)
				}
			}
		}

		_, err := client.aliecsconn.ModifyImageSharePermission(modifyImageSharePermissionRequest)
		if err != nil {
			return fmt.Errorf("Updating Image Share Permission got an error: %#v.", err)
		}
	}

	d.Partial(false)

	return resourceAlicloudImageSharePermissionRead(d, meta)
}

func resourceAlicloudImageSharePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	modifyImageSharePermissionRequest := ecs.CreateModifyImageSharePermissionRequest()
	modifyImageSharePermissionRequest.ImageId = d.Get("image_id").(string)
	accounts := d.Get("accounts").([]interface{})

	modifyImageSharePermissionRequest.RemoveAccount1 = accounts[0].(string)
	if len(accounts) >= 2 {
		modifyImageSharePermissionRequest.RemoveAccount2 = accounts[1].(string)
	}
	if len(accounts) >= 3 {
		modifyImageSharePermissionRequest.RemoveAccount3 = accounts[2].(string)
	}
	if len(accounts) >= 4 {
		modifyImageSharePermissionRequest.RemoveAccount4 = accounts[3].(string)
	}
	if len(accounts) >= 5 {
		modifyImageSharePermissionRequest.RemoveAccount5 = accounts[4].(string)
	}
	if len(accounts) >= 6 {
		modifyImageSharePermissionRequest.RemoveAccount6 = accounts[5].(string)
	}
	if len(accounts) >= 7 {
		modifyImageSharePermissionRequest.RemoveAccount7 = accounts[6].(string)
	}
	if len(accounts) >= 8 {
		modifyImageSharePermissionRequest.RemoveAccount8 = accounts[7].(string)
	}
	if len(accounts) >= 9 {
		modifyImageSharePermissionRequest.RemoveAccount9 = accounts[7].(string)
	}
	if len(accounts) >= 10 {
		modifyImageSharePermissionRequest.RemoveAccount10 = accounts[9].(string)
	}

	describeImageSharePermissionRequest := ecs.CreateDescribeImageSharePermissionRequest()
	describeImageSharePermissionRequest.ImageId = modifyImageSharePermissionRequest.ImageId
	describeImageSharePermissionRequest.PageNumber = requests.NewInteger(1)
	describeImageSharePermissionRequest.PageSize = requests.NewInteger(50)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.aliecsconn.ModifyImageSharePermission(modifyImageSharePermissionRequest)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting Image Share Permission got an error: %#v", err))
		}

		describeImageSharePermissionResponse, err := client.aliecsconn.DescribeImageSharePermission(describeImageSharePermissionRequest)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe Image Share Permission got an error: %#v", err))
		}
		if describeImageSharePermissionResponse.TotalCount == 0 {
			return nil
		} else {
			var allDeleted = true

			for _, account := range describeImageSharePermissionResponse.Accounts.Account {
				for _, accountToDelete := range accounts {
					if account.AliyunId == accountToDelete.(string) {
						allDeleted = false
					}
				}
			}
			if allDeleted {
				return nil
			} else {
				return resource.RetryableError(fmt.Errorf("Deleting Image Share Permission got an error: %#v", err))
			}
		}
	})
}
