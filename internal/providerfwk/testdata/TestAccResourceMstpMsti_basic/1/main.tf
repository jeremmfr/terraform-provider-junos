resource "junos_mstp_msti" "testacc" {
  msti_id = 17
  vlan    = ["32", "33"]
}
