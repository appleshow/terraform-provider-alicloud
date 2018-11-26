# Function Compute

## Service

### resource
Create, Update or Delete the service of Function Compute.

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|service_name|*String*|Yes|||The name of the service.<br>- Only letters, numbers, underscores (_), and hyphens (-) are allowed. <br>- The name cannot start with a number or hyphen.<br>- The name has to be between 1 to 128 characters in length. |
|description|*String*|No|||Service description.|
|role|*String*|No|||TThe role grants Function Compute the permission to access user’s cloud resources, such as pushing logs to user’s log store. The temporary STS token generated from this role can be retrieved from function context and used to access cloud resources. Example : ```"acs:ram::1234567890:role/fc-test"```|
|log_config|*Map*|No|||Log configuration. Function Compute pushes function execution logs to the configured log store. The parameter format is as follows: <br>```{```<br> ``` project = "xxx"```<br>```  logstore = "xxx"```<br>```}``` |

#### example
```
variable access_key {}
variable secret_key {}
variable user_id    {}

provider "alicloud" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  user_id     = "${var.user_id}"
  region      = "cn-hangzhou"
  api_version = "2016-08-15"
}

resource "alicloud_fc_service" "fcService" {
  service_name = "TerraformFCService"
  description  = "Service of function compute created by terraform"
  log_config   = {
                   project = "log-project-fc1"
                   logstore = "instance-health-check1"
                 }
  role        = "acs:ram::1234567890:role/fc-logservice"
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

### data
List services of Function Compute.

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|prefix|*String*|No|||Limits the resource names that begin with the specified prefix.|
|start_key|*String*|No|||Specifies the resource name where to start the query.|
|next_token|*String*|No|||Is used to query more resource records. The token is included in the List API response.|
|limit|*Int*|No|20||Limits the number of returned resource records. Defaults to ```20``` and cannot exceed ```100```.|
|output_file|*String*|No|||Save JSON data in a local file.|

#### example
```

variable access_key {}
variable secret_key {}
variable user_id    {}

provider "alicloud" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  user_id     = "${var.user_id}"
  region      = "cn-hangzhou"
  api_version = "2016-08-15"
}

data "alicloud_fc_services" "fcServices" {
  limit       = 100
  output_file = "./output.json" 
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

#### output
```
[
	{
		"serviceName": "TerraformFCService",
		"description": "Service of function compute created by terraform",
		"role": "acs:ram::1234567890:role/fc-logservice",
		"logConfig": {
			"project": "log-project-fc1",
			"logstore": "instance-health-check1"
		},
		"vpcConfig": {
			"vpcId": "",
			"vSwitchIds": [],
			"securityGroupId": ""
		},
		"internetAccess": true,
		"serviceId": "7bd48381-0df1-4f0a-81af-c8f1ce2e9e0d",
		"createdTime": "2018-06-11T07:58:33Z",
		"lastModifiedTime": "2018-06-11T08:01:01Z"
	}
]
```

## Function

### resource
Create, Update or Delete the Function of Function Compute.

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|service_name|*String*|Yes||||
|function_name|*String*|Yes|||Function name.<br>- Only letters, numbers, underscores (_), and hyphens (-) are allowed.<br>- It cannot start with a number or hyphen.<br>- The name must be 1 to 128 characters in length.|
|description|*String*|No||||
|runtime|*String*|Yes||"nodejs6"<br>"nodejs8"<br>"python2.7"<br>"python3"<br>"java8"|The function runtime environment. Supporting nodejs6, nodejs8, python2.7, python3, java8|
|handler|*String*|Yes|||The function execution entry point. For example: ```index.handler```.|
|timeout|*Int*|No|300||The maximum time duration a function can execute, in seconds. After which Function Compute terminates the execution. Defaults to ```3``` seconds, and should be between ```1``` to ```300``` seconds.|
|memory_size|*Int*|No|128||The amount of memory that’s used to execute function, in MB. Function Compute uses this value to allocate CPU resources proportionally. Defaults to ```128MB```. It should be multiple of ```64``` MB and between ```128MB``` and ```3072MB```.
|code|*String*|Yes|||The code that contains the function implementation.|
|environment_variables|*Map*|No|||The script runtime environment variable.|

#### example
```
variable access_key {}
variable secret_key {}
variable user_id    {}

provider "alicloud" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  user_id     = "${var.user_id}"
  region      = "cn-hangzhou"
  api_version = "2016-08-15"
}

resource "alicloud_fc_function" "fcfunction" {
  service_name  = "TerraformFCService"
  function_name = "TerraformFCFunction"  
  description   = "Function of function compute created by terraform"
  runtime       = "python3"
  handler       = "HealthCheck.instance_health"
  code          = "./HealthCheck.zip"  
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

### data
List Functions of Function Compute.

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|service_name|*String*|Yes||||
|prefix|*String*|No|||Limits the resource names that begin with the specified prefix.|
|start_key|*String*|No|||Specifies the resource name where to start the query.|
|next_token|*String*|No|||Is used to query more resource records. The token is included in the List API response.|
|limit|*Int*|No|20||Limits the number of returned resource records. Defaults to ```20``` and cannot exceed ```100```.|
|output_file|*String*|No|||Save JSON data in a local file.|

#### example
```
variable access_key {}
variable secret_key {}
variable user_id    {}

provider "alicloud" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  user_id     = "${var.user_id}"
  region      = "cn-hangzhou"
  api_version = "2016-08-15"
}

data "alicloud_fc_functions" "fcfunctions" {
  service_name = "TerraformFCService"
  output_file = "./output.json" 
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

#### output
```
[
	{
		"functionId": "039c2867-e5f7-4fa0-b65e-fba2f040ae41",
		"functionName": "TerraformFCFunction",
		"description": "Function of function compute created by terraform",
		"runtime": "python3",
		"handler": "HealthCheck.instance_health",
		"timeout": 300,
		"memorySize": 128,
		"codeSize": 1083654,
		"codeChecksum": "8312311085311512683",
		"environmentVariables": {},
		"createdTime": "2018-06-04T06:38:07Z",
		"lastModifiedTime": "2018-06-04T06:38:07Z"
	}
]
```

## Trigger
Create, Update or Delete the griiger of Function Compute.

### resource

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|service_name|*String*|Yes||||
|function_name|*String*|Yes||||
|trigger_name|*String*|Yes|||Trigger name.<br>- Only letters, numbers, underscores (_), and hyphens (-) are allowed.<br>- The name cannot start with a number or hyphen.<br>- The name can be ```1``` to ```128```characters in length.|
|source_arn|*String*|No|||The Aliyun Resource Name（ARN）of event source. This is optional for some triggers. For example:```"acs:oss:cn-shanghai:12345:mybucket"```|
|trigger_type|*String*|Yes||"oss"<br>"log"<br>"timer"<br>"http"|Trigger type, e.g. oss, timer, logs. This determines how the trigger config is interpreted.For example : ```"oss"```.|
|invocation_role|*String*|No|||The role grants event source the permission to invoke function on behalf of user. This is optional for some triggers. For example:```"acs:ram::1234567890:role/fc-test"```,|
|config_enable|*Bool*|No|true|true<br>false|Enable or disable the trigger.|
|config_payload|*String*|Yes for timer type||"awesome-fc"||
|config_cron_expression|*String*|Yes for timer type|||The frequency of script execution. For example: ```0 2 * * * *```|
|config_events|*String*|Yes for oss type|||OSS event type. For example: ```oss:ObjectCreated:*,oss:ObjectDeleted:*```. Multiple types are separated by ```,```.|
|config_filter_key_prefix|*String*|Yes for oss type||||
|config_filter_key_suffix|*String*|Yes for oss type||||
|config_source_logstore|*String*|Yes for log type|||The LogStore name of the target log service.|
|config_job_interval|*Int*|Yes for log type||||
|config_job_max_retry_time|*Int*|Yes for log type||||
|config_log_project|*String*|Yes for log type|||The name of the project that saved the log.|
|config_log_logstore|*String*|Yes for log type|||The name of the logstore that saved the log.|
|config_auth_type|*String*|Yes for http type||"anonymous"<br>"function"|```"anonymous"``` does not require authorization. ```"function"``` requires authorization.|
|config_methods|*String*|Yes for http type||"GET"<br>"POST"<br>"PUT"<br>"DELETE"<br>"HEAD"|Request method. Multiple methods are separated by ```,```.|

#### example
 - timer
 ```
 variable access_key {}
variable secret_key {}
variable user_id    {}

provider "alicloud" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  user_id     = "${var.user_id}"
  region      = "cn-hangzhou"
  api_version = "2016-08-15"
}

resource "alicloud_fc_trigger" "fctriggertimer" {
  service_name           = "TerraformFCService"
  function_name          = "TerraformFCFunction"
  trigger_name           = "TerraformTriggerTimmer"
  trigger_type           = "timer"
  config_payload         = "awesome-fc"
  config_cron_expression = "0 10 * * * *"
}
```
- log
```
variable access_key {}
variable secret_key {}
variable user_id    {}

provider "alicloud" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  user_id     = "${var.user_id}"
  region      = "cn-hangzhou"
  api_version = "2016-08-15"
}

resource "alicloud_fc_trigger" "fctriggertimer" {
  service_name    = "TerraformFCService"
  function_name   = "TerraformFCFunction"
  trigger_name    = "TerraformTriggerLog"
  trigger_type    = "log"
  source_arn      = "acs:log:cn-hangzhou:1234567890:project/log-project-fc1"
  invocation_role = "acs:ram::1234567890:role/fc-oss"
  config_source_logstore    = "fc-log-trigger"
  config_job_interval       =  60
  config_job_max_retry_time = 10
  config_log_project        = "log-project-fc1"
  config_log_logstore       = "instance-health-check1"
  config_enable             = true
}
```
- oss
```
variable access_key {}
variable secret_key {}
variable user_id    {}

provider "alicloud" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  user_id     = "${var.user_id}"
  region      = "cn-hangzhou"
  api_version = "2016-08-15"
}

resource "alicloud_fc_trigger" "fctriggertimer" {
  service_name    = "TerraformFCService"
  function_name   = "TerraformFCFunction"
  trigger_name    = "TerraformTriggerOSS"
  trigger_type    = "oss"
  source_arn      = "acs:oss:cn-hangzhou:1234567890:test147369"
  invocation_role = "acs:ram::1234567890:role/fc-oss"
  config_events   = "oss:ObjectCreated:*"
  config_filter_key_prefix = "prefix"
  config_filter_key_suffix = "suffix"  
}
```
- http
```
variable access_key {}
variable secret_key {}
variable user_id    {}

provider "alicloud" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  user_id     = "${var.user_id}"
  region      = "cn-hangzhou"
  api_version = "2016-08-15"
}

resource "alicloud_fc_trigger" "fctriggertimer" {
  service_name    = "TerraformFCService"
  function_name   = "TerraformFCFunction"
  trigger_name    = "TerraformTriggerHttp"
  trigger_type    = "http"
  config_auth_type = "anonymous"
  config_methods   = "GET"
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

### data

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|service_name|*String*|Yes||||
|function_name|*String*|Yes||||
|prefix|*String*|No|||Limits the resource names that begin with the specified prefix.|
|start_key|*String*|No|||Specifies the resource name where to start the query.|
|next_token|*String*|No|||Is used to query more resource records. The token is included in the List API response.|
|limit|*Int*|No|20||Limits the number of returned resource records. Defaults to ```20``` and cannot exceed ```100```.|
|output_file|*String*|No|||Save JSON data in a local file.|

#### example
```
variable access_key {}
variable secret_key {}
variable user_id    {}

provider "alicloud" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  user_id     = "${var.user_id}"
  region      = "cn-hangzhou"
  api_version = "2016-08-15"
}

data "alicloud_fc_triggers" "fctriggers" {
  service_name  = "TerraformFCService"
  function_name = "TerraformFCFunction"
  output_file   = "./output.json" 
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

#### output
```
[
	{
		"triggerName": "TerraformTriggerOSS",
		"sourceArn": "acs:oss:cn-hangzhou:1234567890:test147369",
		"triggerType": "oss",
		"invocationRole": "acs:ram::1234567890:role/fc-oss",
		"triggerConfig": {
			"events": [
				"oss:ObjectCreated:*"
			],
			"filter": {
				"key": {
					"prefix": "prefix",
					"suffix": "suffix"
				}
			}
		},
		"createdTime": "2018-06-24T11:27:19Z",
		"lastModifiedTime": "2018-06-24T11:27:19Z"
	},
	{
		"triggerName": "TerraformTriggerTimmer",
		"sourceArn": null,
		"triggerType": "timer",
		"invocationRole": null,
		"triggerConfig": {
			"payload": "awesome-fc",
			"cronExpression": "0 10 * * * *",
			"enable": true
		},
		"createdTime": "2018-06-24T11:26:54Z",
		"lastModifiedTime": "2018-06-24T11:26:54Z"
	}
]
```

## Invoke

### data

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|service_name|*String*|Yes||||
|function_name|*String*|Yes||||
|payload|*String*|No||||
|invocation_type|*String*|No||"Async"<br>"Sync"||
|log_type|*String*|No||"Tail"<br>"None"||
|output_file|*String*|No|||Save JSON data in a local file.|
|environment_variables|*Map*|No||||

#### example
````

variable access_key {}
variable secret_key {}
variable user_id    {}

provider "alicloud" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  user_id     = "${var.user_id}"
  region      = "cn-hangzhou"
  api_version = "2016-08-15"
}

data "alicloud_fc_invokes" "fcinvokes" {
  service_name  = "TerraformFCService"
  function_name = "TerraformFCFunction"
  output_file   = "./output.json" 
}
````
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

#### output
````
{
	"Header": {
		"Access-Control-Allow-Origin": [
			""
		],
		"Access-Control-Expose-Headers": [
			"Date,x-fc-request-id,x-fc-error-type,x-fc-code-checksum,x-fc-max-memory-usage,x-fc-log-result,x-fc-invocation-code-version"
		],
		"Content-Length": [
			"50"
		],
		"Content-Type": [
			"application/octet-stream"
		],
		"Date": [
			"Mon, 11 Jun 2018 08:03:32 GMT"
		],
		"X-Fc-Code-Checksum": [
			"8312311085311512683"
		],
		"X-Fc-Max-Memory-Usage": [
			"60.54"
		],
		"X-Fc-Request-Id": [
			"27888f62-1097-3855-d9f6-0716fcf5fced"
		]
	},
	"Payload": "WwogICAgIltGYWlsZWRdIENhbiBub3QgZmluZCBhbnkgRUNTIGluc3RhbmNlcyEiCl0="
}
````
