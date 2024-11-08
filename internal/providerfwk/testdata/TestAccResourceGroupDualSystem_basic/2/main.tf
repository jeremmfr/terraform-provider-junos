resource "junos_group_dual_system" "testacc_node0" {
  name         = "node0"
  apply_groups = false
  interface_fxp0 {
    description = "test_"
    family_inet_address {
      cidr_ip = "192.0.2.193/26"
    }
    family_inet_address {
      cidr_ip     = "192.0.2.194/26"
      master_only = true
      primary     = true
      preferred   = true
    }
    family_inet6_address {
      cidr_ip     = "fe80::2/64"
      master_only = true
      primary     = true
      preferred   = true
    }
  }
  routing_options {
    static_route {
      destination = "192.0.2.0/26"
      next_hop    = ["192.0.2.254"]
    }
  }
  security {
    log_source_address = "192.0.2.128"
  }
  system {
    host_name             = "test_node"
    backup_router_address = "192.0.2.254"
    backup_router_destination = [
      "192.0.2.64/26",
      "192.0.2.0/26",
    ]
    inet6_backup_router_address = "fe80::1"
    inet6_backup_router_destination = [
      "fe80:b::/48",
      "fe80:a::/48",
    ]
  }
}
