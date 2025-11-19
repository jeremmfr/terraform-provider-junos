resource "junos_group_raw" "testacc_foo" {
  name   = "testacc foo"
  config = <<EOT
system {
    services {
        dns {
            forwarders {
                192.0.2.3;
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
set system ntp peer 192.0.2.1
EOT
}
