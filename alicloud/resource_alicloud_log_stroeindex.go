package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudLogStoreIndex() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogStoreIndexCreate,
		Read:   resourceAlicloudLogStoreIndexRead,
		Update: resourceAlicloudLogStoreIndexUpdate,
		Delete: resourceAlicloudLogStoreIndexDelete,
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
			"keys": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"case_sensitive": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"doc_value": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"alias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"chn": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"line": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"case_sensitive": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"include_keys": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"exclude_keys": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"chn": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
				MaxItems: 1,
			},
		},
	}
}

/**
*
 */
func resourceAlicloudLogStoreIndexCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	projectName := d.Get("project_name").(string)
	storeName := d.Get("store_name").(string)

	keys, okKeys := d.GetOk("keys")
	lines, okLine := d.GetOk("line")
	if !okKeys && !okLine {
		index := sls.Index{
			Line: &sls.IndexLine{
				Token:         []string{" ", "\n", "\t", "\r", ",", ";", "[", "]", "{", "}", "(", ")", "&", "^", "*", "#", "@", "~", "=", "<", ">", "/", "\\", "?", ":", "'", "\""},
				CaseSensitive: false,
			},
		}
		err := client.slsconn.CreateIndex(projectName, storeName, index)
		if err != nil {
			return fmt.Errorf("Creating index to logstore of log service got an error: %#v", err)
		}
	} else {
		indexKeyMap := map[string]sls.IndexKey{}
		indexLine := &sls.IndexLine{}
		if okKeys {
			index := 0
			for _, key := range keys.([]map[string]interface{}) {
				indexKey := sls.IndexKey{}

				if v, ok := key["token"]; ok {
					var tokens []string
					for _, token := range v.([]interface{}) {
						tokens = append(tokens, token.(string))
					}
					indexKey.Token = tokens
				}
				if v, ok := key["case_sensitive"]; ok {
					indexKey.CaseSensitive = v.(bool)
				}
				if v, ok := key["type"]; ok {
					indexKey.Type = v.(string)
				}
				if v, ok := key["doc_value"]; ok {
					indexKey.DocValue = v.(bool)
				}
				if v, ok := key["alias"]; ok {
					indexKey.Alias = v.(string)
				}
				if v, ok := key["chn"]; ok {
					indexKey.Chn = v.(bool)
				}

				indexKeyMap["col_"+string(index)] = indexKey
				index++
			}
		}

		if okLine {
			for _, line := range lines.([]map[string]interface{}) {
				if v, ok := line["token"]; ok {
					var tokens []string
					for _, token := range v.([]interface{}) {
						tokens = append(tokens, token.(string))
					}
					indexLine.Token = tokens
				}
				if v, ok := line["case_sensitive"]; ok {
					indexLine.CaseSensitive = v.(bool)
				}
				if v, ok := line["include_keys"]; ok {
					var includeKeys []string
					for _, includeKey := range v.([]interface{}) {
						includeKeys = append(includeKeys, includeKey.(string))
					}
					indexLine.IncludeKeys = includeKeys
				}
				if v, ok := line["exclude_keys"]; ok {
					var excludeKeys []string
					for _, excludeKey := range v.([]interface{}) {
						excludeKeys = append(excludeKeys, excludeKey.(string))
					}
					indexLine.ExcludeKeys = excludeKeys
				}
				if v, ok := line["chn"]; ok {
					indexLine.Chn = v.(bool)
				}

				break
			}
		}

		index := sls.Index{
			Keys: indexKeyMap,
			Line: indexLine,
		}

		err := client.slsconn.CreateIndex(projectName, storeName, index)
		if err != nil {
			return fmt.Errorf("Creating index to logstore of log service got an error: %#v", err)
		}
	}

	d.SetId(projectName + COMMA_SEPARATED + storeName)

	return resourceAlicloudLogStoreIndexUpdate(d, meta)
}

/**
*
 */
func resourceAlicloudLogStoreIndexRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parameters := strings.Split(d.Id(), COMMA_SEPARATED)
	index, err := client.slsconn.GetIndex(parameters[0], parameters[1])

	if err != nil {
		if NotFoundError(err) || index == nil {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("project_name", parameters[0])
	d.Set("store_name", parameters[1])

	_, okKeys := d.GetOk("keys")
	_, okLine := d.GetOk("line")
	if okKeys {
		var indexKeys []map[string]interface{}
		for _, v := range index.Keys {
			var indexKeyMap map[string]interface{}
			indexKeyMap["token"] = v.Token
			indexKeyMap["case_sensitive"] = v.CaseSensitive
			indexKeyMap["type"] = v.Type
			indexKeyMap["doc_value"] = v.DocValue
			indexKeyMap["alias"] = v.Alias
			indexKeyMap["chn"] = v.Chn
			indexKeys = append(indexKeys, indexKeyMap)
		}
		d.Set("keys", indexKeys)
	}
	if okLine {
		var indexLine []map[string]interface{}
		var indexLineMap map[string]interface{}

		indexLineMap["token"] = index.Line.Token
		indexLineMap["case_sensitive"] = index.Line.CaseSensitive
		indexLineMap["include_keys"] = index.Line.IncludeKeys
		indexLineMap["exclude_keys"] = index.Line.ExcludeKeys
		indexLineMap["chn"] = index.Line.Chn

		indexLine = append(indexLine, indexLineMap)

		d.Set("line", indexLine)
	}

	return nil
}

/**
*
 */
func resourceAlicloudLogStoreIndexUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	d.Partial(true)
	update := false

	if d.HasChange("project_name") && !d.IsNewResource() {
		return fmt.Errorf("Updating index to logstore of log service got an error: %#v", "Cannot modify parameter project_name")
	}
	if d.HasChange("store_name") && !d.IsNewResource() {
		return fmt.Errorf("Updating index to logstore of log service got an error: %#v", "Cannot modify parameter store_name")
	}

	if d.HasChange("keys") {
		update = true
		d.SetPartial("keys")
	}
	if d.HasChange("line") {
		update = true
		d.SetPartial("line")
	}

	if !d.IsNewResource() && update {
		projectName := d.Get("project_name").(string)
		storeName := d.Get("store_name").(string)

		keys, okKeys := d.GetOk("keys")
		lines, okLine := d.GetOk("line")

		indexKeyMap := map[string]sls.IndexKey{}
		indexLine := &sls.IndexLine{}
		if okKeys {
			index := 0
			for _, key := range keys.([]map[string]interface{}) {
				indexKey := sls.IndexKey{}

				if v, ok := key["token"]; ok {
					var tokens []string
					for _, token := range v.([]interface{}) {
						tokens = append(tokens, token.(string))
					}
					indexKey.Token = tokens
				}
				if v, ok := key["case_sensitive"]; ok {
					indexKey.CaseSensitive = v.(bool)
				}
				if v, ok := key["type"]; ok {
					indexKey.Type = v.(string)
				}
				if v, ok := key["doc_value"]; ok {
					indexKey.DocValue = v.(bool)
				}
				if v, ok := key["alias"]; ok {
					indexKey.Alias = v.(string)
				}
				if v, ok := key["chn"]; ok {
					indexKey.Chn = v.(bool)
				}

				indexKeyMap["col_"+string(index)] = indexKey
				index++
			}
		}

		if okLine {
			for _, line := range lines.([]map[string]interface{}) {
				if v, ok := line["token"]; ok {
					var tokens []string
					for _, token := range v.([]interface{}) {
						tokens = append(tokens, token.(string))
					}
					indexLine.Token = tokens
				}
				if v, ok := line["case_sensitive"]; ok {
					indexLine.CaseSensitive = v.(bool)
				}
				if v, ok := line["include_keys"]; ok {
					var includeKeys []string
					for _, includeKey := range v.([]interface{}) {
						includeKeys = append(includeKeys, includeKey.(string))
					}
					indexLine.IncludeKeys = includeKeys
				}
				if v, ok := line["exclude_keys"]; ok {
					var excludeKeys []string
					for _, excludeKey := range v.([]interface{}) {
						excludeKeys = append(excludeKeys, excludeKey.(string))
					}
					indexLine.ExcludeKeys = excludeKeys
				}
				if v, ok := line["chn"]; ok {
					indexLine.Chn = v.(bool)
				}

				break
			}

			index := sls.Index{
				Keys: indexKeyMap,
				Line: indexLine,
			}

			err := client.slsconn.UpdateIndex(projectName, storeName, index)
			if err != nil {
				return fmt.Errorf("Updating index to logstore of log service got an error: %#v", err)
			}
		}
	}

	d.Partial(false)

	return resourceAlicloudLogStoreIndexRead(d, meta)
}

func resourceAlicloudLogStoreIndexDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		parameters := strings.Split(d.Id(), COMMA_SEPARATED)
		err := client.slsconn.DeleteIndex(parameters[0], parameters[1])

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Delete index of log service got an error: %#v", err))
		}

		resp, err := client.slsconn.GetIndex(parameters[0], parameters[1])
		if err != nil {
			if NotFoundError(err) || resp == nil {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe index of log service got an error: %#v", err))
		}
		if resp == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Delete index of log service got an error: %#v", err))
	})
}
