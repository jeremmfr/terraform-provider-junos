resource "junos_services_flowmonitoring_v9_template" "testacc_flow_v9_template_ipv4" {
  name                         = "testacc_template@1"
  type                         = "ipv4-template"
  ip_template_export_extension = ["app-id", "flow-dir"]
  flow_active_timeout          = 60
  flow_inactive_timeout        = 30
  flow_key_flow_direction      = true
  flow_key_vlan_id             = true
  nexthop_learning_enable      = true
  option_refresh_rate {}
  option_template_id = 2000
  template_id        = 2001
  template_refresh_rate {
    packets = 200
    seconds = 120
  }
}
resource "junos_services_flowmonitoring_v9_template" "testacc_flow_v9_template_ipv6" {
  name                         = "testacc_template@3"
  type                         = "ipv6-template"
  ip_template_export_extension = ["app-id", "flow-dir"]
  flow_active_timeout          = 60
  flow_inactive_timeout        = 30
}
resource "junos_services_flowmonitoring_v9_template" "testacc_flow_v9_template_mpls" {
  name                         = "testacc_template@2"
  type                         = "mpls-template"
  flow_active_timeout          = 60
  flow_inactive_timeout        = 30
  flow_key_flow_direction      = true
  flow_key_vlan_id             = true
  mpls_template_label_position = [8, 4]
  nexthop_learning_enable      = true
  option_refresh_rate {
    packets = 100
    seconds = 60
  }
  option_template_id = 2002
  source_id          = 10
  template_id        = 2003
  template_refresh_rate {}
  tunnel_observation_ipv4 = true
  tunnel_observation_ipv6 = true
}
resource "junos_services_flowmonitoring_v9_template" "testacc_flow_v9_template_bridge" {
  name                      = "testacc_template@4"
  type                      = "bridge-template"
  flow_key_output_interface = true
}
resource "junos_services_flowmonitoring_v9_template" "testacc_flow_v9_template_peerBill" {
  name                     = "testacc_template@5"
  type                     = "peer-as-billing-template"
  nexthop_learning_disable = true
  template_refresh_rate {
    seconds = 150
  }
}
