# Alarm

## resource 
Create, Update or Delete Alibaba Cloud Monitoring Alarms.

### paramters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|name|*String*|Yes|||Alarm rule name|
|project|*String*|Yes||"acs_ecs_dashboard"<br>"acs_rds_dashboard"<br>"acs_slb_dashboard"<br>"acs_memcache"<br>"acs_vpc_eip"<br>"acs_kvstore"<br>"acs_messageservice_new"<br>"acs_kvstore"<br>"acs_ads"<br>"acs_kvstore"<br>"acs_express_connect"<br>"acs_fc"<br>"acs_nat_gateway"<br>"acs_sls_dashboard"<br>"acs_containerservice_dashboard"<br>"acs_vpn"<br>"acs_bandwidth_package"<br>"acs_cen"|Product name. For [more information](https://www.alibabacloud.com/help/doc-detail/28619.htm?spm=a2c63.p38356.b99.98.64aa7bb6jxeySR), see the projects for various products, such as ```acs_ecs_dashboard``` and ```acs_rds_dashboard```.|
|metric|*String*|Yes|||Names of the monitoring metrics corresponding to a product. For [more information](https://www.alibabacloud.com/help/doc-detail/28619.htm?spm=a2c63.p38356.b99.98.64aa7bb6jxeySR), see the metric definitions for various products|
|dimensions|*Map*|Yes|||List of instances associated with the alarm rule.For example, ```{"instanceId":"instanceId1,instanceId2,instanceId3"}```. Usring ```{}``` for all instances.|
|period|*Int*|No|300||Index query cycle, which must be consistent with that defined for metrics; default value: ```300```, in seconds.|
|statistics|*String*|No|"Average"|"Average"<br>"Minimum"<br>"Maximum"|Statistical method, for example, ```Average```, which must be consistent with that defined for metrics|
|operator|*String*|No|"=="|">"<br>">="<br>"<"<br>"<="<br>""<br>"=="<br>""<br>"!="|Alarm comparison operator, which must be ```<=```, ```<```, ```>```, ```>=```, ```==```, or ```!=```.|
|threshold|*String*|Yes|||Alarm threshold value, which must be a numeric value currently.|
|triggered_count|*Int*|No|3||Number of consecutive times it has been detected that the values exceed the threshold; default value: three times|
|contact_groups|*Array*|Yes|||The contact group of the alarm rule, which must have been created on the console as a string corresponding to the JSON array, for example, ```["Contact Group 1","Contact Group 2"]```|
|start_time|*Int*|No|0||Start time of the alarm effective period; default value: ```0```, which indicates the time 00:00.|
|end_time|*Int*|No|24||End time of the alarm effective period; default value: ```24```, which indicates the time 24:00.|
|silence_time|*Int*|No|86400||Notification silence period in the alarm state, in seconds; default value: ```86,400```; minimum value: 1 hour|
|notify_type|*Int*|No|0|0<br>1|Notification type. The value ```0``` indicates TradeManager+email, and the value ```1``` indicates that TradeManager+email+SMS|
|enabled|*Bool*|No|true|true<br>false|Enabled or Disabled the alarm.|
|webhook|*String*|No|||If this blank filled with internet accessiable URL, the cloud mointor will send the alarm information to that address via POST request. only HTTP protocol supported currently. For example ```"{'method':'post','url':'https://getman.cn/echo'}"```|

### example
```
variable access_key {}
variable secret_key {}

provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

resource "alicloud_cms_alarm" "CPUUtilization" {
  name = "Default CPUUtilization alarm for All ECS instances"
  project = "acs_ecs_dashboard"
  metric = "CPUUtilization"
  statistics = "Average"
  dimensions = {}
  period = 300
  operator = ">="
  threshold = 85
  triggered_count = 3
  contact_groups = ["Group1"]
  end_time = 24
  start_time = 0
  notify_type = 1
}

resource "alicloud_cms_alarm" "MemoryUsedutilization" {
  name = "Default Memory Usedutilization alarm for All ECS instances"
  project = "acs_ecs_dashboard"
  metric = "memory_usedutilization"
  statistics = "Average"
  dimensions = {}
  period = 300
  operator = ">="
  threshold = 85
  triggered_count = 2
  contact_groups = ["Group1"]
  end_time = 24
  start_time = 0
  notify_type = 1
}

resource "alicloud_cms_alarm" "DiskusageUtilization" {
  name = "Default Diskusage Utilization alarm for All ECS instances"
  project = "acs_ecs_dashboard"
  metric = "diskusage_utilization"
  statistics = "Average"
  dimensions = {}
  period = 300
  operator = ">="
  threshold = 85
  triggered_count = 2
  contact_groups = ["Group1"]
  end_time = 24
  start_time = 0
  notify_type = 1
  webhook = "{'method':'post','url':'https://getman.cn/echo'}"
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

## data
List alarms.

### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|id|*String*|No||||
|name|*String*|No||||
|project|*String*|No||||
|dimension|*String*|No||||
|state|*String*|No||||
|page_number|*Int*|No|1|||
|page_size|*Int*|No|100|||
|output_file|*String*|No|||Save JSON data in a local file.|

### example
```
variable access_key {}
variable secret_key {}

provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

data "alicloud_cms_alarms" "listAlarms" { 
  output_file  = "./output.json"
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

### output
```
[
	{
		"comparison_operator": "\u003e=",
		"contact_groups": "[\"Group1\"]",
		"dimensions": "[\"{}\"]",
		"enable": true,
		"end_time": 24,
		"evaluation_count": 3,
		"id": "14D42849ADAC3ED02E5D3FF764D3B7F025491127",
		"metrie_name": "CPUUtilization",
		"name": "Default CPUUtilization alarm for All ECS instances",
		"name_space": "acs_ecs_dashboard",
		"notify_type": 1,
		"period": 300,
		"silence_time": 86400,
		"start_time": 0,
		"state": "INSUFFICIENT_DATA",
		"statistics": "Average",
		"threshold": "85",
		"webhook": "null"
	},
	{
		"comparison_operator": "\u003e=",
		"contact_groups": "[\"Group1\"]",
		"dimensions": "[\"{}\"]",
		"enable": true,
		"end_time": 24,
		"evaluation_count": 2,
		"id": "2C6F4E320A0E2DFF4488437C34FF839525491127",
		"metrie_name": "memory_usedutilization",
		"name": "Default Memory Usedutilization alarm for All ECS instances",
		"name_space": "acs_ecs_dashboard",
		"notify_type": 1,
		"period": 300,
		"silence_time": 86400,
		"start_time": 0,
		"state": "INSUFFICIENT_DATA",
		"statistics": "Average",
		"threshold": "85",
		"webhook": "null"
	},
	{
		"comparison_operator": "\u003e=",
		"contact_groups": "[\"Group1\"]",
		"dimensions": "[\"{}\"]",
		"enable": true,
		"end_time": 24,
		"evaluation_count": 2,
		"id": "3BF0D3F23FC972DC23F5D9429324E98325491127",
		"metrie_name": "diskusage_utilization",
		"name": "Default Diskusage Utilization alarm for All ECS instances",
		"name_space": "acs_ecs_dashboard",
		"notify_type": 1,
		"period": 300,
		"silence_time": 86400,
		"start_time": 0,
		"state": "INSUFFICIENT_DATA",
		"statistics": "Average",
		"threshold": "85",
		"webhook": "{\"method\":\"post\",\"url\":\"https://getman.cn/echo\"}"
	}
]
```
