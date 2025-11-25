resource "junos_group_raw" "testacc_foobar" {
  name   = "testacc foo bar"
  format = "set"
  config = <<EOT
set system time-zone Europe/paris
set system default-address-selection
EOT
}

resource "junos_apply_group_except" "testacc_foobar" {
  name   = junos_group_raw.testacc_foobar.name
  prefix = "system services "
}
