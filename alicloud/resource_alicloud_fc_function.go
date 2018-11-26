package alicloud

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudFcFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFcFunctionCreate,
		Read:   resourceAlicloudFcFunctionRead,
		Update: resourceAlicloudFcFunctionUpdate,
		Delete: resourceAlicloudFcFunctionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"function_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"runtime": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"handler": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"memory_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  128,
			},
			"code": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_variables": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			// Computed values
			"alicloud_fc_function": {
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
						"runtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"environment_variables": {
							Type:     schema.TypeMap,
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

func resourceAlicloudFcFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	serviceName := d.Get("service_name").(string)
	functionName := d.Get("function_name").(string)
	createFunctionInput := fc.NewCreateFunctionInput(serviceName)

	createFunctionInput.WithFunctionName(functionName)
	if description, ok := d.GetOk("description"); ok {
		createFunctionInput.WithDescription(description.(string))
	}
	createFunctionInput.WithRuntime(d.Get("runtime").(string))
	createFunctionInput.WithHandler(d.Get("handler").(string))
	if timeout, ok := d.GetOk("timeout"); ok {
		if timeout, err := strconv.ParseInt(strconv.Itoa(timeout.(int)), 10, 32); err != nil {
			return fmt.Errorf("Creating function of function compute got an error: %#v", err)
		} else {
			createFunctionInput.WithTimeout(int32(timeout))
		}
	}
	if memorySize, ok := d.GetOk("memory_size"); ok {
		if memorySize, err := strconv.ParseInt(strconv.Itoa(memorySize.(int)), 10, 32); err != nil {
			return fmt.Errorf("Creating function of function compute got an error: %#v", err)
		} else {
			createFunctionInput.WithMemorySize(int32(memorySize))
		}
	}
	code := d.Get("code").(string)
	if _, err := os.Stat(code); err != nil {
		if os.IsNotExist(err) {
			// Does not exist
			return fmt.Errorf("Creating function of function compute got an error: %#v", "File or directory does not exist")
		} else {
			// Other error
			return fmt.Errorf("Creating function of function compute got an error: %#v", err)
		}
	}
	if strings.HasSuffix(code, ".zip") || strings.HasSuffix(code, ".jar") {
		if zipFile, err := ioutil.ReadFile(code); err != nil {
			return fmt.Errorf("Creating function of function compute got an error: %#v", err)
		} else {
			createFunctionInput.WithCode(fc.NewCode().WithZipFile(zipFile))
		}
	} else {
		createFunctionInput.WithCode(fc.NewCode().WithDir(code))
	}
	if environmentVariables, ok := d.GetOk("environment_variables"); ok {
		environmentVariablesStr := make(map[string]string)
		for k, v := range environmentVariables.(map[string]interface{}) {
			environmentVariablesStr[k] = v.(string)
		}
		createFunctionInput.WithEnvironmentVariables(environmentVariablesStr)
	}

	_, err := client.fcconn.CreateFunction(createFunctionInput)
	if err != nil {
		return fmt.Errorf("Creating function of function compute got an error: %#v", err)
	}

	d.SetId(serviceName + COMMA_SEPARATED + functionName)

	return resourceAlicloudFcFunctionUpdate(d, meta)
}

func resourceAlicloudFcFunctionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parameters := strings.Split(d.Id(), COMMA_SEPARATED)
	getFunctionOutput, err := client.fcconn.GetFunction(fc.NewGetFunctionInput(parameters[0], parameters[1]))

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("service_name", parameters[0])
	d.Set("function_name", getFunctionOutput.FunctionName)
	d.Set("description", getFunctionOutput.Description)
	d.Set("runtime", getFunctionOutput.Runtime)
	d.Set("handler", getFunctionOutput.Handler)
	d.Set("timeout", getFunctionOutput.Timeout)
	d.Set("memory_size", getFunctionOutput.MemorySize)
	d.Set("environment_variables", getFunctionOutput.EnvironmentVariables)
	d.Set("code", "")

	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":                    *getFunctionOutput.FunctionID,
		"name":                  *getFunctionOutput.FunctionName,
		"status":                "Available",
		"creation_time":         *getFunctionOutput.CreatedTime,
		"description":           *getFunctionOutput.Description,
		"runtime":               *getFunctionOutput.Runtime,
		"handler":               *getFunctionOutput.Handler,
		"timeout":               strconv.Itoa(int(*getFunctionOutput.Timeout)),
		"memory_size":           strconv.Itoa(int(*getFunctionOutput.MemorySize)),
		"code_size":             strconv.FormatInt(*getFunctionOutput.CodeSize, 10),
		"last_modified_time":    *getFunctionOutput.LastModifiedTime,
		"environment_variables": getFunctionOutput.EnvironmentVariables,
		"resource_type":         "alicloud_fc_function",
	}
	s = append(s, mapping)

	if err := d.Set("alicloud_fc_function", s); err != nil {
		return fmt.Errorf("Setting alicloud_fc_function got an error: %#v.", err)
	}

	return nil
}

func resourceAlicloudFcFunctionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parameters := strings.Split(d.Id(), COMMA_SEPARATED)

	d.Partial(true)
	update := false

	updateFunctionInput := fc.NewUpdateFunctionInput(parameters[0], parameters[1])

	if d.HasChange("description") {
		update = true
		code := d.Get("code").(string)
		if _, err := os.Stat(code); err != nil {
			if os.IsNotExist(err) {
				// Does not exist
				return fmt.Errorf("Creating function of function compute got an error: %#v", "File or directory does not exist")
			} else {
				// Other error
				return fmt.Errorf("Creating function of function compute got an error: %#v", err)
			}
		}
		if strings.HasSuffix(code, ".zip") {
			if zipFile, err := ioutil.ReadFile(code); err != nil {
				return fmt.Errorf("Creating function of function compute got an error: %#v", err)
			} else {
				updateFunctionInput.WithCode(fc.NewCode().WithZipFile(zipFile))
			}
		} else {
			updateFunctionInput.WithCode(fc.NewCode().WithDir(code))
		}
		d.SetPartial("code")
	}

	if d.HasChange("description") {
		update = true
		updateFunctionInput.WithDescription(d.Get("description").(string))
		d.SetPartial("description")
	}
	if d.HasChange("runtime") {
		update = true
		updateFunctionInput.WithRuntime(d.Get("runtime").(string))
		d.SetPartial("runtime")
	}
	if d.HasChange("handler") {
		update = true
		updateFunctionInput.WithHandler(d.Get("handler").(string))
		d.SetPartial("handler")
	}
	if d.HasChange("timeout") {
		update = true
		if timeout, err := strconv.ParseInt(strconv.Itoa(d.Get("timeout").(int)), 10, 32); err != nil {
			return fmt.Errorf("UPdating function of function compute got an error: %#v", err)
		} else {
			updateFunctionInput.WithTimeout(int32(timeout))
		}
		d.SetPartial("timeout")
	}
	if d.HasChange("memory_size") {
		update = true
		if memorySize, err := strconv.ParseInt(strconv.Itoa(d.Get("memory_size").(int)), 10, 32); err != nil {
			return fmt.Errorf("Updating function of function compute got an error: %#v", err)
		} else {
			updateFunctionInput.WithMemorySize(int32(memorySize))
		}
		d.SetPartial("memory_size")
	}
	if d.HasChange("environment_variables") {
		update = true
		environmentVariables := d.Get("environment_variables")
		environmentVariablesStr := make(map[string]string)
		for k, v := range environmentVariables.(map[string]interface{}) {
			environmentVariablesStr[k] = v.(string)
		}
		updateFunctionInput.WithEnvironmentVariables(environmentVariablesStr)
		d.SetPartial("environment_variables")
	}

	if !d.IsNewResource() && update {
		if _, err := client.fcconn.UpdateFunction(updateFunctionInput); err != nil {
			return fmt.Errorf("Updating function of function compute got an error: %#v", err)
		}
	}

	d.Partial(false)

	return resourceAlicloudFcFunctionRead(d, meta)
}

func resourceAlicloudFcFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		parameters := strings.Split(d.Id(), COMMA_SEPARATED)
		_, err := client.fcconn.DeleteFunction(fc.NewDeleteFunctionInput(parameters[0], parameters[1]))

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting function of function compute got an error: %#v", err))
		}

		resp, err := client.fcconn.GetFunction(fc.NewGetFunctionInput(parameters[0], parameters[1]))
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe function of function compute got an error: %#v", err))
		}
		if resp.FunctionName == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting function of function compute got an error: %#v", err))
	})
}
