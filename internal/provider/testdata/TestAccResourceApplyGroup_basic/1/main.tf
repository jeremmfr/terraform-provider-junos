resource "junos_group_raw" "testacc_foobar" {
  name   = "testacc foo bar"
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

resource "junos_apply_group" "testacc_foobar" {
  name   = junos_group_raw.testacc_foobar.name
  prefix = "system "
}

resource "junos_group_raw" "testacc_barfoo" {
  name   = "testacc bar foo"
  format = "set"
  config = <<EOT
set system time-zone Europe/paris
set system default-address-selection
EOT
}

resource "junos_apply_group" "barfoo" {
  name = junos_group_raw.testacc_barfoo.name
}
