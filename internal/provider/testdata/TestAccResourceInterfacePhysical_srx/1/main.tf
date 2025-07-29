resource "junos_interface_physical" "testacc_interface" {
  depends_on = [
    junos_chassis_cluster.testacc_interface
  ]
  name        = var.interface
  description = "testacc_interface"
  gigether_opts {
    redundant_parent = "reth0"
  }
}
resource "junos_interface_physical" "testacc_interface2" {
  name = var.interface2
}
resource "junos_chassis_cluster" "testacc_interface" {
  fab0 {
    member_interfaces = [junos_interface_physical.testacc_interface2.name]
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
  }
  reth_count = 1
}
resource "junos_interface_physical" "testacc_interface_reth" {
  depends_on = [
    junos_interface_physical.testacc_interface
  ]
  name        = "reth0"
  description = "testacc_interface_reth"
  parent_ether_opts {
    redundancy_group = 1
    minimum_links    = 1
  }
}
