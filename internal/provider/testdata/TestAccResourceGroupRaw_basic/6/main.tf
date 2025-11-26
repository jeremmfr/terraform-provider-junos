resource "junos_group_raw" "testacc" {
  name   = "testacc"
  format = "set"

  // intentionally invalid
  // valid config is "set system host-name junDevice"
  config = <<EOT
set system hostname junDevice
EOT
}
