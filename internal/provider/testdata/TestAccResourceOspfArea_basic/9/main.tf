resource "junos_ospf_area" "testacc_ospfarea" {
  area_id = "1"
  version = "v3"
  interface {
    name = "all"
  }
  area_range {
    range = "fe80::/64"
  }
  no_context_identifier_advertisement = true
  inter_area_prefix_export = [
    junos_policyoptions_policy_statement.testacc_ospfarea2.name,
  ]
  inter_area_prefix_import = [
    junos_policyoptions_policy_statement.testacc_ospfarea.name,
  ]
  nssa {}
}
resource "junos_policyoptions_policy_statement" "testacc_ospfarea" {
  name = "testacc_ospfarea"
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_ospfarea2" {
  name = "testacc_ospfarea2"
  then {
    action = "reject"
  }
}
resource "junos_ospf_area" "testacc_ospfarea2" {
  area_id = "2"
  version = "v3"
  interface {
    name    = "${var.interface}.0"
    passive = true
  }
  stub {}
}
