resource "junos_services_flowmonitoring_vipfix_template" "testacc_v0toV1_flowtemplate" {
  name                         = "testacc_v0toV1_flowtemplate"
  type                         = "ipv4-template"
  ip_template_export_extension = ["app-id", "flow-dir"]
  flow_active_timeout          = 60
  flow_inactive_timeout        = 30
  flow_key_flow_direction      = true
  flow_key_vlan_id             = true
  nexthop_learning_enable      = true
  observation_domain_id        = 10
  option_template_id           = 2000
  template_id                  = 2001
  option_refresh_rate {
    packets = 100
    seconds = 60
  }
}
