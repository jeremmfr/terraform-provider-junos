resource "junos_vlan" "testacc_vlansw" {
  name         = "testacc_vlansw"
  description  = "testacc_vlansw"
  vlan_id_list = ["1001-1002"]
  private_vlan = "community"
}
