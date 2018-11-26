# Configure the Alicloud Provider
provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

resource "alicloud_auto_snapshot_policy" "test" {
    time_points = "[\"0\"]"
    repeat_weekdays = "[\"1\"]"
    retention_days = 20
    name = "OKiOK"
}

resource "alicloud_disk" "ecs_disk" {
  availability_zone = "cn-hangzhou-g"
  name              = "New-disk"
  description       = "Hello ecs disk."
  category          = "cloud_efficiency"
  size              = "20"

  tags {
    Name = "TerraformTest"
  }
}

resource "alicloud_auto_snapshot_policy_application" "backup" {
  auto_snapshot_policy_id = "${alicloud_auto_snapshot_policy.test.id}"
  disk_id = "${alicloud_disk.ecs_disk.id}"
}


