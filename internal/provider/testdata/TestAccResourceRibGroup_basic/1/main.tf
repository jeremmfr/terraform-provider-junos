resource "junos_routing_instance" "testacc_ribGroup1" {
  name = "testacc_ribGroup1"
}
resource "junos_policyoptions_policy_statement" "testacc_ribGroup" {
  name = "testacc ribGroup"
  then {
    action = "accept"
  }
}
resource "junos_rib_group" "testacc_ribGroup" {
  name = "testacc ribGroup test"
  import_policy = [
    junos_policyoptions_policy_statement.testacc_ribGroup.name,
  ]
  import_rib = [
    "${junos_routing_instance.testacc_ribGroup1.name}.inet.0",
  ]
  export_rib = "${junos_routing_instance.testacc_ribGroup1.name}.inet.0"
}
