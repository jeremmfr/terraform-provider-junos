resource "junos_interface_logical" "testacc_interface_logical" {
  name = "ip-0/0/0.0"
  tunnel {
    destination         = "192.0.2.10"
    source              = "192.0.2.11"
    allow_fragmentation = true
    path_mtu_discovery  = true
  }
}
