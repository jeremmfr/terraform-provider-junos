resource "junos_interface_logical" "testacc_interface_logical" {
  name                      = "irb.100"
  proxy_macip_advertisement = true
}
