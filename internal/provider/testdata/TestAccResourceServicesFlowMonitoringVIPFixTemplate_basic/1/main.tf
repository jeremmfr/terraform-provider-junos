resource "junos_services_flowmonitoring_vipfix_template" "testacc_flow_vipfix_template_ipv4" {
  name = "testacc_template@1"
  type = "ipv4-template"
}
resource "junos_services_flowmonitoring_vipfix_template" "testacc_flow_vipfix_template_ipv6" {
  name = "testacc_template@3"
  type = "ipv6-template"
}
resource "junos_services_flowmonitoring_vipfix_template" "testacc_flow_vipfix_template_mpls" {
  name = "testacc_template@2"
  type = "mpls-template"
}
