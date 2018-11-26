# Log Service

## project

### resource

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|project_name|*String*|Yes|||The name of the project.<br>- The project name can only contain lowercase letters, numbers, hyphen (-) and underscores (_)<br>- It must begin and end with lowercase letters or numbers.<br>- It should contain 3-63 characters.|
|description|*String*|No||||

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

resource "alicloud_log_project" "logproject" {
  project_name = "terraform-log-project"
  description  = "Project of log service created by terraformAAAA"
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

### data

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
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

data "alicloud_log_projects" "logProjects" {
  output_file = "./output.json" 
}
```

#### output
```
[
	{
		"project_name": "terraform-log-project"
	},
	{
		"project_name": "log-project-fc1"
	}
]
```

## stor

### resource

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|project_name|*String*|Yes||||
|store_name|*String*|Yes|||The name of the store.<br>- The project name can only contain lowercase letters, numbers, hyphen (-) and underscores (_)<br>- It must begin and end with lowercase letters or numbers.<br>- It should contain 3-63 characters.|
|ttl|*Int*|No|30||The data retention time (in days).|
|shard_count|*Int*|No|1||The number of shards in this Logstore.|

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

resource "alicloud_log_store" "logstore" {
  project_name = "terraform-log-project"
  store_name   = "terraform-log-store-a"
  ttl          = 7
  shard_count  = 1
}

```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

### data

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|project_name|*String*|Yes||||
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

data "alicloud_log_stores" "logStores" {
  project_name = "terraform-log-project" 
  output_file  = "./output.json" 
}

```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

#### output
```
[
	{
		"store_name": "terraform-log-store-a"
	}
]
```

## config

### resource

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|project_name|*String*|Yes||||
|store_name|*String*|Yes||||
|config_name|*String*|Yes|||The name of configuration.  The name can be 3 to 63 characters in length and contain lowercase letters, numbers, hyphens (-), and underscores (_). It must begin and end with a lowercase letter or number.|
|log_path|*String*|Yes|||The parent directory where the log resides. For example: ```/var/logs/```.|
|file_pattern|*String*|Yes|||The pattern of a log file. For example: ```access*.log```.|
|log_sample|*String*|No|||The log sample of the Logtail configuration. The log size cannot exceed 1,000 bytes.|
|keys|*Array*|No|||The key generated after logs are extracted.|
|topic_format|*String*|No|||The topic generation mode. The four supported modes are as follows:<br>- Use a part of the log file path as the topic. For example, /var/log/(.*).log.<br>- none indicates the topic is empty.<br>- default indicates to use the log file path as the topic.<br>- group_topic indicates to use the topic attribute of the machine group that applies this configuration as the topic.|
|local_storage|*Bool*|No|true|true<br>false|Whether or not to activate the local cache. Logs of 1 GB can be cached locally when the link to Log Service is disconnected.|
|time_key|*String*|No||||
|time_format|*String*|No|||The format of log time. For example: ```%Y/%m/%d %H:%M:%S```.|
|log_begin_regex|*String*|No|".*"||The characteristics (regular expression) of the first log line, which is used to match with logs composed of multiple lines.|
|regex|*String*|No|"(.*)"||The regular expression used for extracting logs.|
|filter_keys|*Array*|No|||The key used for filtering logs. The log meets the requirements only when the key value matches the regular expression specified in the corresponding ```filterRegex``` column.|
|filter_regex|*Array*|No|||The regular expression corresponding to each filterKey. The length of filterRegex must be the same as that of filterKey.|

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

resource "alicloud_log_config" "logstore" {
  project_name = "terraform-log-project"
  store_name   = "terraform-log-store-a"
  config_name  = "terraform-log-config-a"
  log_path     = "/var/log"
  file_pattern = "messages"
}

````
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

### data

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|project_name|*String*|Yes||||
|group_name|*String*|No||||
|offset|*Int*|No|0|||
|size|*Int*|No|100|||
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

data "alicloud_log_configs" "logConfigs" {
  project_name = "terraform-log-project"
  output_file  = "./output.json" 
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

#### output
```
[
	{
		"config_name": "terraform-log-config-a"
	}
]
```

## machine group

### resource

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|project_name|*String*|Yes||||
|group_name|*String*|Yes|||The machine group name, which is unique in the same project.<br>- The name can only contain lowercase letters, numbers, hyphens (-) and underscores (_).<br>- It should begin and end with lowercase letters or numbers.<br>- The name should be 3-128 characters long.|
|machine_id_type|*String*|Yes||"ip"<br>"userdefined"|The machine identification type, including ```ip``` and ```userdefined``` identity.|
|type|*String*|No|""||The machine group type, which is empty by default.|
|attribute_external_mame|*String*|No|""||The external identification that the machine group depends, which is empty by default.|
|attribute_topic_name|*String*|No|""||The topic of a machine group, which is empty by default.|

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

resource "alicloud_log_machinegroup" "logMachineGroup" {
  project_name = "terraform-log-project"
  group_name   = "terraform-machine-group-a"
  machine_id_type = "ip"
  machine_id_list = ["172.22.22.22"]
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

### data
#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|project_name|*String*|Yes||||
|config_name|*String*|No||||
|offset|*Int*|No|0|||
|size|*Int*|No|100|||
|output_file|*String*|No|||Save JSON data in a local file.|

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

data "alicloud_log_machinegroups" "logMachineGroups" {
  project_name = "terraform-log-project"
  config_name  = "terraform-log-config-a" 
  output_file  = "./output.json" 
}
````
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

##### output
```

```

## config to machine group

### resource

#### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|project_name|*Yes*||||
|config_name|*Yes*||||
|group_name|*Yes*||||

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

resource "alicloud_log_configtomachinegroup" "configToMachineGroup" {
  project_name = "terraform-log-project"
  config_name  = "terraform-log-config-a"
  group_name   = "terraform-machine-group-a"
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```
