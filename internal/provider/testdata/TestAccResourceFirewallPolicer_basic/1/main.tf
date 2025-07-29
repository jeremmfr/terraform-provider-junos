resource "junos_firewall_policer" "testacc_fwPolic" {
  name = "testacc_fwPolic"
  if_exceeding {
    bandwidth_percent = 80
    burst_size_limit  = "50k"
  }
  then {
    discard = true
  }
}
resource "junos_firewall_policer" "testacc_fwPolic2" {
  name                      = "testacc_fwPolic2"
  filter_specific           = true
  logical_interface_policer = true
  if_exceeding_pps {
    packet_burst = "33k"
    pps_limit    = "51k"
  }
  then {
    forwarding_class = "best-effort"
  }
}
