resource "junos_rstp_interface" "all" {
  name                      = "all"
  access_trunk              = true
  bpdu_timeout_action_alarm = true
  bpdu_timeout_action_block = true
  cost                      = 16
  edge                      = true
  mode                      = "shared"
  priority                  = 240
}
