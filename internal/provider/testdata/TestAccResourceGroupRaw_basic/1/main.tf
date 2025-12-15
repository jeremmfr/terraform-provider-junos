resource "junos_group_raw" "testacc_foo" {
  name   = "testacc foo"
  config = <<EOT
system {
    services {
        dns {
            forwarders {
                192.0.2.3;
                192.0.2.33;
            }
        }
    }
}
EOT
}

resource "junos_group_raw" "testacc_bar" {
  name   = "testacc bar"
  format = "set"
  config = <<EOT
set system time-zone Europe/paris
set system default-address-selection
EOT
}

resource "junos_group_raw" "testacc" {
  name   = "testacc"
  format = "set"
  config = <<EOT
set system host-name junDevice
EOT
}
