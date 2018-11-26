# Command

## resource 
Create, Update or Delete Alibaba Cloud Assistant Command.

### paramters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|name|*String*|Yes|||Command name. Supporting all the character encoding sets.|
|type|*String*|Yes||"RunBatScript"<br>"RunPowerShellScript"<br>"RunShellScript"|Command type. RunBatScript: Creates a Bat script for Windows instances. RunPowerShellScript: Creates a PowerShell script for Windows instances. RunShellScript: Creates a Shell script for Linux instances. Optional values:<br>RunBatScript: Creates a Bat script for Windows instances.<br>RunPowerShellScript: Creates a PowerShell script for Windows instances.<br>RunShellScript: Creates a Shell script for Linux instances.|
|description|*String*|No|||Command description. Supporting all the character encoding sets.|
|command_content|*String*|Yes|||The Base64-encoded content of the command. You must pass in this parameter at the same time when you pass in the Type request parameter. The parameter value must be Base64-encoded for transmission and the script content size before the Base64 encoding cannot exceed 16 KB.|
|working_dir|*String*|No|||	The directory where your created command runs on the ECS instances. Default value:<br>For Linux instances, commands are performed in the /root directory.<br>For Windows instances, commands are performed in the directory where the cloud assistant client process is located, such as C:\ProgramData\aliyun\assist\$(version).|
|time_out|*String*|No|||	The invocation timeout value of the command. The unit is seconds. When the command fails to run for some reason, the invocation may time out, after which the cloud assistant client forces the command process to stop by canceling the command PID. The parameter value must be greater than or equal to 60. If the value is smaller than 60, the timeout value is 60 seconds by default. Default value: 3600.|

### example
```
variable access_key {}
variable secret_key {}

provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

resource "alicloud_command" "command" {
  name = "TerraformCommand"
  type = "RunShellScript"
  description  = "Shell script created by Terraform"
  command_content = "ps -ef"
}

```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'
```

## data
List Commands.

### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|id|*String*|No|||Command ID.|
|name|*String*|No|||Command name, fuzzy search temporarily not supported.|
|type|*String*|No|||	Command type. Optional values:<br>RunBatScript: The command process is a Bat script for Windows instances.<br>RunPowerShellScript: The command process is a PowerShell script for Windows instances.<br>RunShellScript: The command process is a Shell script for Linux instances.|
|page_number|*String*|No|1||Current page number. Start value: 1. Default value: 1.|
|page_size|*String*|No|50||The number of rows per page for multi-page display. Maximum value: 50.|
|output_file|*String*|No|||Save JSON data in a local file.|

### example
```
variable access_key {}
variable secret_key {}

provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

data "alicloud_commands" "commands" {
  output_file  = "./output.json"
}

```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'
```

### output
```
[
	{
		"command_content": "$client = new-object System.Net.WebClient\n$client.DownloadFile('http://logtail-release.oss-cn-hangzhou.aliyuncs.com/win/logtail_installer.zip', 'c:\\logtail_installer.zip')\n\n$Source = 'c:\\logtail_installer.zip'\n$Destination = 'C:\\'\n$ShowDestinationFolder = $true\n \nif ((Test-Path $Destination) -eq $false)\n{\n  $null = mkdir $Destination\n}\n \n$shell = New-Object -ComObject Shell.Application\n$sourceFolder = $shell.NameSpace($Source)\n$destinationFolder = $shell.NameSpace($Destination)\n$DestinationFolder.CopyHere($sourceFolder.Items())\n \nif ($ShowDestinationFolder)\n{\n  explorer.exe $Destination\n}\n\ncd c:\\logtail_installer\n.\\logtail_installer.exe install cn_hangzhou_vpc",
		"description": "Download Ali LogTail for windows",
		"id": "c-16b8dbe95bf04a60b61ac8bddf0f1dc3",
		"name": "Download_LogTail_Windows_cn-hangzhou_VPC",
		"time_out": 3600,
		"type": "RunPowerShellScript",
		"working_dir": ""
	},
	{
		"command_content": "wget http://logtail-release.vpc100-oss-cn-hangzhou.aliyuncs.com/linux64/logtail.sh -O logtail.sh; chmod 755 logtail.sh; sh logtail.sh install cn-hangzhou_vpc",
		"description": "Download Ali LogTail for linux",
		"id": "c-a93ddcd31f8a4151b90f5bccbf6d7377",
		"name": "Download_LogTail_Linux_cn-hangzhou_VPC",
		"time_out": 3600,
		"type": "RunShellScript",
		"working_dir": ""
	},
	{
		"command_content": "ls",
		"description": "Shell script created by Terraform",
		"id": "c-ec3a2252d7e84abaa5969dc986488932",
		"name": "TerraformCommand",
		"time_out": 0,
		"type": "RunShellScript",
		"working_dir": ""
	}
]
```

# Invoke Command

## resource 
Invoke Alibaba Cloud Assistant Command.

### paramters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|command_id|*String*|Yes|||Command ID. You can call the DescribeCommands API to check all the available CommandId.|
|instance_ids|*List*|Yes|||List of instances for command invocation. The parameter value is a formatted JSON array in the format of [InstanceId1, instanceId2, …]. You can specify a maximum of 100 instance IDs separated by commas (,).|
|timed|*Int*|No|||Whether the command is periodically performed or not. Optional values:<br>True: Periodical invocation.<br>False: Non-periodical invocation.|
|frequency|*String*|No|||The invocation period of a periodical task. When the Timed parameter value is True, the Frequency parameter is required.<br>The parameter value observes the Cron expression. For more information, see ![Cron expressions](https://www.alibabacloud.com/help/faq-detail/64769.htm?spm=a2c63.p38356.a3.7.9bf81014ghiYRx).|

### example
```
variable access_key {}
variable secret_key {}

provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

resource "alicloud_command_invoke" "command_invoke" {
  command_id   = "c-c9c16c81747e45eca021f6324fa29be8"
  instance_ids = ["i-bp12a63qihelfqn4g6v5"]
  timed        = true
  frequency    = "0 15 10 ? * *"
}


```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'
```

## data invoke
List invokes.

### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|invoke_id|*String*|No|||Invocation ID of a command process.|
|command_id|*String*|No|||Command ID. You can query all the available CommandId by calling the DescribeCommands|
|instance_id|*String*|No|||Instance ID. When you pass in the parameter, the system queries the invocation-record status of all commands on the instance.|
|command_name|*String*|No|||Name of the command.|
|command_type|*String*|No|||Type of the command. Optional values:<br>RunBatScript: A Bat script for Windows instances.<br>RunPowerShellScript: A PowerShell script for Windows instances.<br>RunShellScript: A Shell script for Linux instances|
|invoke_status|*String*|No|||Specifies the overall invocation status of a command. The overall invocation status is depends on  the invocation status of command processes on all target instances. Optional values:<br>Running: The command process is running.<br>Periodical invocation: Before you manually stop a command on periodical invocation, the command process is always in the Running (Running) status.<br>One-time invocation: The overall invocation status is Running (Running), as long as the command process is in the Running (Running) status on any target instance.<br>Finished: The invocation of the command process is completed.<br>Periodical invocation: The command process can never be in the Finished (Finished) status.<br>One-time invocation: When invocation of command processes on all the target instances managed by the specified command is completed, the overall invocation status is Finished (Finished).<br>Or when you manually stop the invocation of command processes on some target instances (Stopped) and the invocation on other target instances is completed, the overall invocation status is Finished (Finished).<br>Failed: The command process invocation failed, the command process timed out, or encountered other exceptions.<br>Periodical invocation: The command process can never be in the Failed (Failed) status.<br>One-time invocation: When invocation of command processes on all the target instances managed by the specified command fails, the overall invocation status is Failed (Failed).<br>PartialFailed: Part of the invocation of command processes failed.<br>Periodical invocation: The command process can never be in the PartialFailed (PartialFailed) status.<br>One-time invocation: The overall invocation status is PartialFailed (PartialFailed) when any command process is in the Failed (Failed) status on any target instance managed by the specified command.<br>Stopped: The command process is manually stopped.|
|timed|*String*|No|||Whether the command is periodically performed or not. Optional values:<br>True: Periodical invocation. False: Not periodical invocation.Default value: False.|
|page_number|*String*|No|1||	Current page number. Start value: 1.<br>Default value: 1.|
|page_size|*String*|No|50||The number of rows per page for multi-page display. Maximum value: 50.|
|output_file|*String*|No|||Save JSON data in a local file.|

### example
```
variable access_key {}
variable secret_key {}

provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

data "alicloud_command_invokes" "command_invokes" {
  command_id   = "c-ec3a2252d7e84abaa5969dc986488932"
  output_file  = "./output.json"
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'
```

### output
```
[
	{
		"command_id": "c-ec3a2252d7e84abaa5969dc986488932",
		"command_name": "TerraformCommand",
		"command_type": "RunShellScript",
		"frequency": "",
		"invocation_results": null,
		"invoke_id": "t-320542646eae4e418dc50d5050c1b908",
		"invoke_instances": [
			{
				"instance_id": "i-bp12a63qihelfqn4g6v5",
				"instance_invoke_status": "Finished"
			}
		],
		"invoke_status": "Finished",
		"page_number": 0,
		"page_size": 0,
		"timed": false,
		"total_count": 0
	},
	{
		"command_id": "c-ec3a2252d7e84abaa5969dc986488932",
		"command_name": "TerraformCommand",
		"command_type": "RunShellScript",
		"frequency": "0 15 10 ? * *",
		"invocation_results": null,
		"invoke_id": "t-9339318ec6984529a2b02df03f96da65",
		"invoke_instances": [
			{
				"instance_id": "i-bp12a63qihelfqn4g6v5",
				"instance_invoke_status": "Running"
			}
		],
		"invoke_status": "Running",
		"page_number": 0,
		"page_size": 0,
		"timed": true,
		"total_count": 0
	},
	{
		"command_id": "c-ec3a2252d7e84abaa5969dc986488932",
		"command_name": "TerraformCommand",
		"command_type": "RunShellScript",
		"frequency": "",
		"invocation_results": null,
		"invoke_id": "t-96e361db064443b99569b65d621e3232",
		"invoke_instances": [
			{
				"instance_id": "i-bp12a63qihelfqn4g6v5",
				"instance_invoke_status": "Finished"
			}
		],
		"invoke_status": "Finished",
		"page_number": 0,
		"page_size": 0,
		"timed": false,
		"total_count": 0
	},
	{
		"command_id": "c-ec3a2252d7e84abaa5969dc986488932",
		"command_name": "TerraformCommand",
		"command_type": "RunShellScript",
		"frequency": "",
		"invocation_results": null,
		"invoke_id": "t-cc67f10c8d6c4331a688d3954ad235b9",
		"invoke_instances": [
			{
				"instance_id": "i-bp12a63qihelfqn4g6v5",
				"instance_invoke_status": "Finished"
			}
		],
		"invoke_status": "Finished",
		"page_number": 0,
		"page_size": 0,
		"timed": false,
		"total_count": 0
	},
	{
		"command_id": "c-ec3a2252d7e84abaa5969dc986488932",
		"command_name": "TerraformCommand",
		"command_type": "RunShellScript",
		"frequency": "",
		"invocation_results": null,
		"invoke_id": "t-ee9e0235931441b395256d37902673ab",
		"invoke_instances": [
			{
				"instance_id": "i-bp12a63qihelfqn4g6v5",
				"instance_invoke_status": "Finished"
			}
		],
		"invoke_status": "Finished",
		"page_number": 0,
		"page_size": 0,
		"timed": false,
		"total_count": 0
	}
]
```
## data invoke results
List invoke results.

### parameters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|invoke_id|*String*|Yes|||Invocation ID of a command process. You can use the DescribeInvocations API to check all the InvokeId.|
|command_id|*String*|No|||Command ID. You can use the DescribeCommands API to check all the available CommandId.|
|instance_id|*String*|No|||	Instance ID.|
|invoke_record_status|*String*|No|||The status of the command process you want to query. Optional values:<br>Running: The command process is running.<br>Failed: The command process invocation failed, the command process timed out, or encountered exceptions.<br>Finished: The invocation of the command process is completed.<br>Stopped: The command process is manually stopped.|
|page_number|*String*|No|1||	Current page number. Start value: 1.<br>Default value: 1.|
|page_size|*String*|No|50||The number of rows per page for multi-page display. Maximum value: 50.|
|output_file|*String*|No|||Save JSON data in a local file.|

### example
```
variable access_key {}
variable secret_key {}

provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

data "alicloud_command_invoke_results" "command_invoke_results" {
  invoke_id   = "t-c48896db70514c90b3b306b3aeefd510"
  output_file  = "./output.json"
}

```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'
```

### output
```
[
	{
		"command_id": "",
		"command_name": "",
		"command_type": "",
		"frequency": "",
		"invocation_results": [
			{
				"command_id": "c-c9c16c81747e45eca021f6324fa29be8",
				"exit_code": 0,
				"finished_time": "Thu Jul 19 10:01:28 CST 2018",
				"instance_id": "i-bp12a63qihelfqn4g6v5",
				"invoke_id": "",
				"invoke_record_status": "Finished",
				"output": "UID        PID  PPID  C STIME TTY          TIME CMD\nroot         1     0  0 Jul18 ?        00:00:01 /usr/lib/systemd/systemd --switched-root --system --deserialize 21\nroot         2     0  0 Jul18 ?        00:00:00 [kthreadd]\nroot         3     2  0 Jul18 ?        00:00:00 [ksoftirqd/0]\nroot         5     2  0 Jul18 ?        00:00:00 [kworker/0:0H]\nroot         6     2  0 Jul18 ?        00:00:00 [kworker/u2:0]\nroot         7     2  0 Jul18 ?        00:00:00 [migration/0]\nroot         8     2  0 Jul18 ?        00:00:00 [rcu_bh]\nroot         9     2  0 Jul18 ?        00:00:02 [rcu_sched]\nroot        10     2  0 Jul18 ?        00:00:00 [watchdog/0]\nroot        12     2  0 Jul18 ?        00:00:00 [kdevtmpfs]\nroot        13     2  0 Jul18 ?        00:00:00 [netns]\nroot        14     2  0 Jul18 ?        00:00:00 [khungtaskd]\nroot        15     2  0 Jul18 ?        00:00:00 [writeback]\nroot        16     2  0 Jul18 ?        00:00:00 [kintegrityd]\nroot        17     2  0 Jul18 ?        00:00:00 [bioset]\nroot        18     2  0 Jul18 ?        00:00:00 [kblockd]\nroot        19     2  0 Jul18 ?        00:00:00 [md]\nroot        25     2  0 Jul18 ?        00:00:00 [kswapd0]\nroot        26     2  0 Jul18 ?        00:00:00 [ksmd]\nroot        27     2  0 Jul18 ?        00:00:00 [khugepaged]\nroot        28     2  0 Jul18 ?        00:00:00 [crypto]\nroot        36     2  0 Jul18 ?        00:00:00 [kthrotld]\nroot        37     2  0 Jul18 ?        00:00:00 [kworker/u2:1]\nroot        38     2  0 Jul18 ?        00:00:00 [kmpath_rdacd]\nroot        39     2  0 Jul18 ?        00:00:00 [kpsmoused]\nroot        40     2  0 Jul18 ?        00:00:00 [ipv6_addrconf]\nroot        59     2  0 Jul18 ?        00:00:00 [deferwq]\nroot        91     2  0 Jul18 ?        00:00:00 [kauditd]\nroot        93     2  0 Jul18 ?        00:00:01 [kworker/0:2]\nroot       225     2  0 Jul18 ?        00:00:00 [ata_sff]\nroot       233     2  0 Jul18 ?        00:00:00 [scsi_eh_0]\nroot       234     2  0 Jul18 ?        00:00:00 [scsi_tmf_0]\nroot       235     2  0 Jul18 ?        00:00:00 [scsi_eh_1]\nroot       236     2  0 Jul18 ?        00:00:00 [scsi_tmf_1]\nroot       239     2  0 Jul18 ?        00:00:00 [ttm_swap]\nroot       251     2  0 Jul18 ?        00:00:00 [kworker/0:1H]\nroot       256     2  0 Jul18 ?        00:00:00 [jbd2/vda1-8]\nroot       257     2  0 Jul18 ?        00:00:00 [ext4-rsv-conver]\nroot       324     1  0 Jul18 ?        00:00:00 /usr/lib/systemd/systemd-journald\nroot       352     1  0 Jul18 ?        00:00:00 /usr/lib/systemd/systemd-udevd\nroot       368     1  0 Jul18 ?        00:00:00 /sbin/auditd\ndbus       449     1  0 Jul18 ?        00:00:00 /bin/dbus-daemon --system --address=systemd: --nofork --nopidfile --systemd-activation\nroot       451     2  0 Jul18 ?        00:00:00 [edac-poller]\nroot       460     1  0 Jul18 ?        00:00:00 /usr/lib/systemd/systemd-logind\npolkitd    461     1  0 Jul18 ?        00:00:00 /usr/lib/polkit-1/polkitd --no-debug\nroot       462     1  0 Jul18 ?        00:00:02 /usr/sbin/rsyslogd -n\nroot       465     1  0 Jul18 ?        00:00:00 /usr/sbin/atd -f\nroot       467     1  0 Jul18 ?        00:00:00 /usr/sbin/crond -n\nroot       488     1  0 Jul18 ttyS0    00:00:00 /sbin/agetty --keep-baud 115200 38400 9600 ttyS0 vt220\nroot       489     1  0 Jul18 tty1     00:00:00 /sbin/agetty --noclear tty1 linux\nroot       708     1  0 Jul18 ?        00:00:00 /sbin/dhclient -1 -q -lf /var/lib/dhclient/dhclient--eth0.lease -pf /var/run/dhclient-eth0.pid -H izbp12a63qihelfqn4g6v5z eth0\nroot       771     1  0 Jul18 ?        00:00:06 /usr/bin/python -Es /usr/sbin/tuned -l -P\nntp        819     1  0 Jul18 ?        00:00:00 /usr/sbin/ntpd -u ntp:ntp -g\nroot       891     1  0 Jul18 ?        00:00:17 /usr/local/aegis/aegis_update/AliYunDunUpdate\nroot      1054     1  0 Jul18 ?        00:01:33 /usr/local/aegis/aegis_client/aegis_10_47/AliYunDun\nroot      1178     1  0 Jul18 ?        00:00:00 /usr/sbin/sshd -D\nroot      1212     1  0 Jul18 ?        00:00:04 /usr/sbin/aliyun-service\nroot      2337     2  0 09:50 ?        00:00:00 [kworker/0:1]\nroot      2345     2  0 10:00 ?        00:00:00 [kworker/0:0]\nroot      2361  1212  0 10:01 ?        00:00:00 ps -ef\n"
			}
		],
		"invoke_id": "",
		"invoke_instances": null,
		"invoke_status": "",
		"page_number": 1,
		"page_size": 20,
		"timed": false,
		"total_count": 1
	}
]
```