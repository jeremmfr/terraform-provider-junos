resource "junos_layer2_control" "l2c" {
  bpdu_block {
    interface {
      name = var.interface2
    }
    interface {
      name    = var.interface
      disable = true
      drop    = true
    }
  }
  mac_rewrite_interface {
    name           = var.interface2
    enable_all_ifl = true
    protocol       = ["cdp", "stp"]
  }
  mac_rewrite_interface {
    name = var.interface
  }
}
