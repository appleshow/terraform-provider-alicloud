package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudLogStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogStoreCreate,
		Read:   resourceAlicloudLogStoreRead,
		Update: resourceAlicloudLogStoreUpdate,
		Delete: resourceAlicloudLogStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"store_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"shard_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			// Computed values
			"alicloud_log_store": {
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
						"project_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"store_name": {
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

func resourceAlicloudLogStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	projectName := d.Get("project_name").(string)
	storeName := d.Get("store_name").(string)
	ttl := d.Get("ttl").(int)
	shardCount := d.Get("shard_count").(int)

	err := client.slsconn.CreateLogStore(projectName, storeName, ttl, shardCount)
	if err != nil {
		return fmt.Errorf("Creating store of log service got an error: %#v", err)
	}

	d.SetId(projectName + COMMA_SEPARATED + storeName)

	return resourceAlicloudLogStoreUpdate(d, meta)
}

func resourceAlicloudLogStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parameters := strings.Split(d.Id(), COMMA_SEPARATED)

	logStore, err := client.slsconn.GetLogStore(parameters[0], parameters[1])

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("project_name", parameters[0])
	d.Set("store_name", parameters[1])
	d.Set("ttl", logStore.TTL)
	d.Set("shard_count", logStore.ShardCount)

	var s []map[string]interface{}

	mapping := map[string]interface{}{
		"id":            parameters[1],
		"name":          parameters[1],
		"status":        "Available",
		"creation_time": time.Now().Format("2006-01-02 15:04:05"),
		"project_name":  parameters[0],
		"store_name":    parameters[1],
		"resource_type": "alicloud_log_store",
	}

	s = append(s, mapping)

	if err := d.Set("alicloud_log_store", s); err != nil {
		return fmt.Errorf("Setting alicloud_log_store got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudLogStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parameters := strings.Split(d.Id(), COMMA_SEPARATED)

	d.Partial(true)
	update := false

	if d.HasChange("project_name") && !d.IsNewResource() {
		return fmt.Errorf("Updating stroe of log service got an error: %#v", "Cannot modify parameter project_name")
	}
	if d.HasChange("store_name") && !d.IsNewResource() {
		return fmt.Errorf("Updating stroe of log service got an error: %#v", "Cannot modify parameter store_name")
	}
	if d.HasChange("ttl") {
		update = true
		d.SetPartial("ttl")
	}
	if d.HasChange("shard_count") {
		update = true
		d.SetPartial("shard_count")
	}

	if !d.IsNewResource() && update {
		ttl := d.Get("ttl").(int)
		shardCount := d.Get("shard_count").(int)

		if err := client.slsconn.UpdateLogStore(parameters[0], parameters[1], ttl, shardCount); err != nil {
			return fmt.Errorf("Updating store of log service got an error: %#v", err)
		}
	}

	d.Partial(false)

	return resourceAlicloudLogStoreRead(d, meta)
}

func resourceAlicloudLogStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		parameters := strings.Split(d.Id(), COMMA_SEPARATED)
		err := client.slsconn.DeleteLogStore(parameters[0], parameters[1])

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting store of log service got an error: %#v", err))
		}

		resp, err := client.slsconn.GetLogStore(parameters[0], parameters[1])
		if err != nil {
			if NotFoundError(err) || resp == nil {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe store of log service got an error: %#v", err))
		}
		if resp == nil || resp.Name == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting store of log service got an error: %#v", err))
	})
}
