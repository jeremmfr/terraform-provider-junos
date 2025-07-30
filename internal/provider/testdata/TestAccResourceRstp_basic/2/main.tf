resource "junos_rstp" "testacc_rstp" {
  backup_bridge_priority = "32k"
  bridge_priority        = "16k"
  system_id {
    id = "00:22:33:44:55:aa"
  }
}
