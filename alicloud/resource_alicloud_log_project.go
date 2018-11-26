package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudLogProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogProjectCreate,
		Read:   resourceAlicloudLogProjectRead,
		Update: resourceAlicloudLogProjectUpdate,
		Delete: resourceAlicloudLogProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"alicloud_log_project": {
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
						"project_name": {
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

func resourceAlicloudLogProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	projectName := d.Get("project_name").(string)
	description := d.Get("description").(string)

	_, err := client.slsconn.CreateProject(projectName, description)
	if err != nil {
		return fmt.Errorf("Creating project of log service got an error: %#v", err)
	}

	d.SetId(projectName)

	return resourceAlicloudLogProjectUpdate(d, meta)
}

func resourceAlicloudLogProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	logProject, err := client.slsconn.GetProject(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("project_name", logProject.Name)
	d.Set("description", logProject.Description)

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":            logProject.Name,
		"name":          logProject.Name,
		"status":        "Available",
		"creation_time": time.Now().Format("2006-01-02 15:04:05"),
		"project_name":  logProject.Name,
		"resource_type": "alicloud_log_project",
	}

	s = append(s, mapping)
	if err := d.Set("alicloud_log_project", s); err != nil {
		return fmt.Errorf("Setting alicloud_log_project got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudLogProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	projectName := d.Id()
	description := ""

	d.Partial(true)
	update := false

	if d.HasChange("project_name") && !d.IsNewResource() {
		return fmt.Errorf("Updating project of log service got an error: %#v", "Cannot modify parameter project_name")
	}
	if d.HasChange("description") {
		update = true
		description = d.Get("description").(string)
		d.SetPartial("description")
	}

	if !d.IsNewResource() && update {
		if _, err := client.slsconn.UpdateProject(projectName, description); err != nil {
			return fmt.Errorf("Updating project of log service got an error: %#v", err)
		}
	}

	d.Partial(false)

	return resourceAlicloudLogProjectRead(d, meta)
}

func resourceAlicloudLogProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		err := client.slsconn.DeleteProject(d.Id())

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting project of log service got an error: %#v", err))
		}

		resp, err := client.slsconn.GetProject(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe project of log service got an error: %#v", err))
		}
		if resp == nil || resp.Name == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting project of log service got an error: %#v", err))
	})
}
