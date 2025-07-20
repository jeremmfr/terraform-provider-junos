resource "junos_routing_instance" "testacc_routingInst" {
  name            = "testacc_routingInst"
  as              = "65000"
  description     = "testacc routingInst"
  instance_export = [junos_policyoptions_policy_statement.testacc_routingInst2.name]
  instance_import = [junos_policyoptions_policy_statement.testacc_routingInst2.name]
  router_id       = "192.0.2.65"
}
resource "junos_policyoptions_community" "testacc_routingInst2" {
  name    = "testacc_routingInst2"
  members = ["target:65000:100"]
}
resource "junos_policyoptions_policy_statement" "testacc_routingInst2" {
  name = "testacc_routingInst2"
  from {
    bgp_community = [junos_policyoptions_community.testacc_routingInst2.name]
  }
  then {
    action = "accept"
  }
}
