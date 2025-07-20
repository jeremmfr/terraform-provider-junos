resource "junos_firewall_policer" "testacc_fwPolic" {
  name                      = "testacc_fwPolic"
  logical_bandwidth_policer = true
  logical_interface_policer = true
  shared_bandwidth_policer  = true
  if_exceeding {
    bandwidth_limit  = "32k"
    burst_size_limit = "52k"
  }
  then {
    forwarding_class = "best-effort"
    loss_priority    = "high"
  }
}
resource "junos_firewall_policer" "testacc_fwPolic2" {
  name                       = "testacc_fwPolic2"
  physical_interface_policer = true
  if_exceeding_pps {
    packet_burst = "35k"
    pps_limit    = "51k"
  }
  then {
    forwarding_class = "best-effort"
  }
}
