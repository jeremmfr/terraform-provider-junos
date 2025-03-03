resource "junos_chassis_redundancy" "l2c" {
  graceful_switchover = true
}
resource "junos_layer2_control" "l2c" {
  depends_on = [
    junos_chassis_redundancy.l2c
  ]
  lifecycle {
    create_before_destroy = true
  }
  bpdu_block {
    disable_timeout = 300
    interface {
      name    = var.interface
      disable = true
      drop    = true
    }
  }
  mac_rewrite_interface {
    name           = var.interface
    enable_all_ifl = true
    protocol       = ["cdp", "stp"]
  }
  nonstop_bridging = true
}
