resource "junos_interface_logical" "testacc_v0to1_ikegateway" {
  name = "${var.interface}.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_policy" "testacc_v0to1_ikegateway" {
  name                = "testacc_v0to1_ikegateway"
  proposal_set        = "basic"
  mode                = "aggressive"
  pre_shared_key_text = "thePassWord"
}
resource "junos_security_ike_gateway" "testacc_v0to1_ikegateway" {
  name               = "testacc_v0to1_ikegateway"
  policy             = junos_security_ike_policy.testacc_v0to1_ikegateway.name
  external_interface = junos_interface_logical.testacc_v0to1_ikegateway.name
  dynamic_remote {
    distinguished_name {
      container = "dc=example,dc=com"
    }
    connections_limit = 10
  }
  aaa {
    client_username = "user"
    client_password = "password"
  }
}
