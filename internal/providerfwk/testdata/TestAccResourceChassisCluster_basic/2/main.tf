resource "junos_interface_physical" "testacc_cluster_int2" {
  name        = var.interface2
  description = "testacc_cluster_int2"
  gigether_opts {
    redundant_parent = "reth0"
  }
}
resource "junos_interface_physical" "testacc_cluster" {
  name = var.interface
}
resource "junos_chassis_cluster" "testacc_cluster" {
  fab0 {
    member_interfaces = [junos_interface_physical.testacc_cluster.name]
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
    interface_monitor {
      name   = junos_interface_physical.testacc_cluster_int2.name
      weight = 255
    }
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
    preempt        = true
    preempt_delay  = 2
    preempt_limit  = 3
    preempt_period = 4
  }
  reth_count = 3
}
