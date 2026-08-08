resource "junos_interface_logical" "testacc_ikegw_wo" {
  name = "${var.interface}.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_proposal" "testacc_ikegw_wo" {
  name                     = "testacc_ikegw_wo"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group2"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikegw_wo" {
  name                = "testacc_ikegw_wo"
  proposals           = [junos_security_ike_proposal.testacc_ikegw_wo.name]
  mode                = "aggressive"
  pre_shared_key_text = "mysecret"
}
resource "junos_security_ike_gateway" "testacc_ikegw_wo" {
  name = "testacc_ikegw_wo"
  dynamic_remote {
    distinguished_name {
      container = "dc=example,dc=com"
    }
    connections_limit = 10
  }
  aaa {
    client_username            = "user"
    client_password_wo         = "thePassWord"
    client_password_wo_version = 1
  }
  policy             = junos_security_ike_policy.testacc_ikegw_wo.name
  external_interface = junos_interface_logical.testacc_ikegw_wo.name
}
# read the configuration to check that the write-only password has been sent to the device,
# the resource is unchanged in this step so the data source is read during the plan
data "junos_config_raw" "testacc_ikegw_wo" {
  format = "set"
}
