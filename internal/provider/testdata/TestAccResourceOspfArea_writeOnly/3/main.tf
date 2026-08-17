resource "junos_ospf_area" "testacc_ospfarea_wo" {
  area_id = "0.0.0.0"
  interface {
    name                                      = "all"
    authentication_simple_password_wo         = "simplPa2"
    authentication_simple_password_wo_version = 2
  }
  interface {
    name = junos_interface_logical.testacc_ospfarea_wo.name
    authentication_md5 {
      key_id         = 1
      key_wo         = "md5keyone2"
      key_wo_version = 2
    }
    authentication_md5 {
      key_id         = 2
      key_wo         = "md5keytwo2"
      key_wo_version = 2
    }
  }
}
resource "junos_interface_logical" "testacc_ospfarea_wo" {
  name        = "${var.interface}.0"
  description = "testacc_ospfarea_wo"
}
