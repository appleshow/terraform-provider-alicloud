# VPC

# VSwitch

# Router Interface

## resource

## data

### paramters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|owner_id|*Int*|No||||
|resource_owner_account|*String*|No||||
|resource_owner_id|*Int*|No||||
|page_number|*Int*|No|1|||
|page_size|*Int*|No|50||Max:50|
|output_file|*String*|No||||

### example
```
variable "access_key" {}
variable "secret_key" {}

provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "cn-hangzhou"
}

data "alicloud_router_interfaces" "interfaces" {
  output_file = "./out_put.json"
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```

### output
```
[
	{
		"access_point_id": "",
		"business_status": "Normal",
		"charge_type": "AfterPay",
		"connected_time": "2018-07-02T10:58:17Z",
		"creation_time": "2018-07-02T10:57:49Z",
		"description": "Guanbang Test",
		"end_time": "2999-09-08T16:00:00Z",
		"health_check_source_ip": "",
		"health_check_target_ip": "",
		"name": "route-interfacea-guanbang",
		"opposite_access_point_id": "",
		"opposite_interface_business_status": "Normal",
		"opposite_interface_id": "ri-bp119eijc1y2lz0um3yoc",
		"opposite_interface_owner_id": "1563557888557255",
		"opposite_interface_spec": "Negative",
		"opposite_interface_status": "Active",
		"opposite_region_id": "cn-hangzhou",
		"opposite_router_id": "vrt-bp1irxxvencg7sak3ohwq",
		"opposite_router_type": "VRouter",
		"opposite_vpc_instance_id": "vpc-bp1nuze31437oy7ym0all",
		"role": "InitiatingSide",
		"router_id": "vrt-bp1c9txaqdt0pa5k8y4wf",
		"router_interface_id": "ri-bp1b0lhawm991qys5l8zd",
		"router_type": "VRouter",
		"spec": "Large.2",
		"status": "Active",
		"vpc_instance_id": "vpc-bp1fy4x6oluft9icr1aa6"
	},
	{
		"access_point_id": "",
		"business_status": "Normal",
		"charge_type": "AfterPay",
		"connected_time": "2018-07-02T10:58:11Z",
		"creation_time": "2018-07-02T10:57:48Z",
		"description": "Guanbang Test",
		"end_time": "2999-09-08T16:00:00Z",
		"health_check_source_ip": "",
		"health_check_target_ip": "",
		"name": "route-interfaceb-guanbang",
		"opposite_access_point_id": "",
		"opposite_interface_business_status": "Normal",
		"opposite_interface_id": "ri-bp1b0lhawm991qys5l8zd",
		"opposite_interface_owner_id": "1563557888557255",
		"opposite_interface_spec": "Large.2",
		"opposite_interface_status": "Active",
		"opposite_region_id": "cn-hangzhou",
		"opposite_router_id": "vrt-bp1c9txaqdt0pa5k8y4wf",
		"opposite_router_type": "VRouter",
		"opposite_vpc_instance_id": "vpc-bp1fy4x6oluft9icr1aa6",
		"role": "AcceptingSide",
		"router_id": "vrt-bp1irxxvencg7sak3ohwq",
		"router_interface_id": "ri-bp119eijc1y2lz0um3yoc",
		"router_type": "VRouter",
		"spec": "Negative",
		"status": "Active",
		"vpc_instance_id": "vpc-bp1nuze31437oy7ym0all"
	}
]
```

# Router Interface Connect

## resource 
Connect two router interfaces.

### paramters
|Name|Type|Required|Default|Option|Description|
|:---|:---|:---|:---|:---|:---|
|router_interface_from_id|*String*|Yes|||Initiator's router interface id|
|router_interface_from_owner_id|*String*|Yes|||Initiator's owner id|
|router_interface_from_region_id|*String*|Yes|||Initiator's region|
|router_interface_to_id|*String*|Yes|||Receiver's router interface id|
|router_interface_to_owner_id|*String*|Yes|||Receiver's owner id|
|router_interface_to_region_id|*String*|Yes|||Receiver's region|

### example
```
variable "access_key" {}
variable "secret_key" {}

provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "cn-hangzhou"
}

data "alicloud_zones" "zones" {
}

resource "alicloud_vpc" "vpc_a" {
  name        = "VPC-Guanbang"
  cidr_block  = "192.168.0.0/16"
  description = "Guanbang Test"
}

resource "alicloud_vswitch" "vsw_public_a" {
  vpc_id            = "${alicloud_vpc.vpc_a.id}"
  cidr_block        = "192.168.0.0/24"
  availability_zone = "${data.alicloud_zones.zones.zones.0.id}"
  name              = "vsw-pu-guanbang"
  description       = "Guanbang Test"
}

resource "alicloud_vswitch" "vsw_private_a" {
  vpc_id            = "${alicloud_vpc.vpc_a.id}"
  cidr_block        = "192.168.2.0/24"
  availability_zone = "${data.alicloud_zones.zones.zones.0.id}"
  name              = "vsw-pr-guanbang"
  description       = "Guanbang Test"
}

resource "alicloud_vpc" "vpc_b" {
  name        = "VPC-Guanbang"
  cidr_block  = "192.168.0.0/16"
  description = "Guanbang Test"
}

resource "alicloud_vswitch" "vsw_public_b" {
  vpc_id            = "${alicloud_vpc.vpc_b.id}"
  cidr_block        = "192.168.0.0/24"
  availability_zone = "${data.alicloud_zones.zones.zones.0.id}"
  name              = "vsw-pu-guanbang"
  description       = "Guanbang Test"
}

resource "alicloud_vswitch" "vsw_private_b" {
  vpc_id            = "${alicloud_vpc.vpc_b.id}"
  cidr_block        = "192.168.2.0/24"
  availability_zone = "${data.alicloud_zones.zones.zones.0.id}"
  name              = "vsw-pr-guanbang"
  description       = "Guanbang Test"
}

resource "alicloud_router_interface" "interface_a" {
  opposite_region = "cn-hangzhou"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.vpc_a.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "route-interfacea-guanbang"
  description = "Guanbang Test"
}

resource "alicloud_router_interface" "interface_b" {
  opposite_region = "cn-hangzhou"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.vpc_b.router_id}"
  role = "AcceptingSide"
  specification = "Negative"
  name = "route-interfaceb-guanbang"
  description = "Guanbang Test"
}


resource "alicloud_router_interface_connect" "connect" {
  router_interface_from_id = "${alicloud_router_interface.interface_a.id}"
  router_interface_from_owner_id = "xxxx"
  router_interface_from_region_id = "cn-hangzhou"
  router_interface_to_id = "${alicloud_router_interface.interface_b.id}"
  router_interface_to_owner_id = "xxxxx"
  router_interface_to_region_id = "cn-hangzhou"
}
```
```
terraform apply -var 'access_key=xxx' -var 'secret_key=xxx'  -var 'user_id=xxx' 
```
