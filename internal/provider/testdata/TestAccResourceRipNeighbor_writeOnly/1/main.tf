resource "junos_rip_group" "testacc_ripneigh_wo" {
  name = "test_rip_group_wo"
}
resource "junos_rip_neighbor" "testacc_ripneigh_wo" {
  name                          = "ae0.0"
  group                         = junos_rip_group.testacc_ripneigh_wo.name
  authentication_type           = "md5"
  authentication_key_wo         = "ripKeyOne"
  authentication_key_wo_version = 1
}
resource "junos_rip_neighbor" "testacc_ripneigh_wo_md5" {
  name  = "ae1.0"
  group = junos_rip_group.testacc_ripneigh_wo.name
  authentication_selective_md5 {
    key_id         = 1
    key_wo         = "md5keyone"
    key_wo_version = 1
  }
  authentication_selective_md5 {
    key_id         = 2
    key_wo         = "md5keytwo"
    key_wo_version = 1
  }
}
