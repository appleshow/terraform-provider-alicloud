package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudEssAttachment_basic(t *testing.T) {
	var sg ess.ScalingGroupItemType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_attachment.attach",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAttachmentExists(
						"alicloud_ess_attachment.attach", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_ess_attachment.attach",
						"instance_ids.#", "2"),
				),
			},
		},
	})
}

func testAccCheckEssAttachmentExists(n string, d *ess.ScalingGroupItemType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS attachment ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		group, err := client.DescribeScalingGroupById(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error Describe scaling group: %#v", err)
		}

		if group == nil {
			return fmt.Errorf("Scaling group not found")
		}

		instances, _, err := client.essconn.DescribeScalingInstances(&ess.DescribeScalingInstancesArgs{
			RegionId:       client.Region,
			ScalingGroupId: rs.Primary.ID,
			CreationType:   "Attached",
		})

		if err != nil {
			return fmt.Errorf("Error Describe scaling instances: %#v", err)
		}

		if len(instances) < 1 {
			return fmt.Errorf("Scaling instances not found")
		}

		*d = *group
		return nil
	}
}

func testAccCheckEssAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_configuration" {
			continue
		}

		group, err := client.DescribeScalingGroupById(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidScalingGroupIdNotFound) {
				return nil
			}
			return fmt.Errorf("Error Describe scaling group: %#v", err)
		}

		if group != nil {
			return fmt.Errorf("Scaling group still existed.")
		}

		instances, _, err := client.essconn.DescribeScalingInstances(&ess.DescribeScalingInstancesArgs{
			RegionId:       client.Region,
			ScalingGroupId: rs.Primary.ID,
			CreationType:   "Attached",
		})

		if err != nil {
			if IsExceptedError(err, InvalidScalingGroupIdNotFound) {
				return nil
			}
			return fmt.Errorf("Error Describe scaling instances: %#v", err)
		}

		if len(instances) > 0 {
			return fmt.Errorf("There are still ECS instances in the scaling group.")
		}
	}

	return nil
}

const testAccEssAttachmentConfig = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
	vpc_id = "${alicloud_vpc.vpc.id}"
	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	description = "foo"
	vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group_rule" "ssh-in" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 0
	max_size = 2
	scaling_group_name = "scaling-group-for-test-case"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_id = "${alicloud_vswitch.vswitch.id}"
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "ecs.n4.small"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = true
	active = true
  	enable = true
}

resource "alicloud_instance" "instance" {
	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "ecs.n4.small"
	count = "2"
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = "10"
	instance_charge_type = "PostPaid"
	system_disk_category = "cloud_efficiency"
	vswitch_id = "${alicloud_vswitch.vswitch.id}"
}

resource "alicloud_ess_attachment" "attach" {
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	instance_ids = ["${alicloud_instance.instance.*.id}"]
	force = true
}
`
