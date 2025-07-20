resource "junos_firewall_policer" "testacc_fwPolic" {
  name = "testacc_fwPolic"
  if_exceeding {
    bandwidth_limit  = "32k"
    burst_size_limit = "50k"
  }
  then {
    forwarding_class = "best-effort"
    loss_priority    = "high"
    out_of_profile   = true
  }
}
