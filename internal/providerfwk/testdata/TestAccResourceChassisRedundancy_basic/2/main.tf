resource "junos_chassis_redundancy" "testacc_cred" {
  failover_disk_read_threshold      = 2000
  failover_disk_write_threshold     = 3000
  failover_not_on_disk_underperform = true
  failover_on_disk_failure          = true
  failover_on_loss_of_keepalives    = true
  keepalive_time                    = 300
  routing_engine {
    slot = 0
    role = "master"
  }
}
