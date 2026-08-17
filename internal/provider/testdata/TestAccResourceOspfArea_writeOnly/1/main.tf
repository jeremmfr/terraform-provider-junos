resource "junos_ospf_area" "testacc_ospfarea_wo" {
  area_id = "0.0.0.0"
  interface {
    name                                      = "all"
    authentication_simple_password_wo         = "simplPa1"
    authentication_simple_password_wo_version = 1
  }
  interface {
    name = junos_interface_logical.testacc_ospfarea_wo.name
    authentication_md5 {
      key_id         = 1
      key_wo         = "md5keyone"
      key_wo_version = 1
    }
    authentication_md5 {
      key_id         = 2
      key_wo         = "md5keytwo"
      key_wo_version = 1
    }
  }
}
resource "junos_interface_logical" "testacc_ospfarea_wo" {
  name        = "${var.interface}.0"
  description = "testacc_ospfarea_wo"
}
