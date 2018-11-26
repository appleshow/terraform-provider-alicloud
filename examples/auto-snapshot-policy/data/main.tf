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
    name = "OK"
}
variable "access_key" {}
variable "secret_key" {}
variable "region" {
  default = "cn-hangzhou"
}
