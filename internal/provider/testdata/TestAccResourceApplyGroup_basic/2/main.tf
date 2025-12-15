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
  name = junos_group_raw.testacc_foobar.name
}

resource "junos_group_raw" "testacc_barfoo" {
  name   = "testacc bar foo"
  format = "set"
  config = <<EOT
set system time-zone Europe/paris
set system default-address-selection
EOT
}

resource "junos_apply_group" "testacc_barfoo" {
  name   = junos_group_raw.testacc_barfoo.name
  prefix = "system "
}

resource "junos_snmp_view" "testacc_barfoo" {
  name        = "testacc snmpview"
  oid_include = [".1"]
}

resource "junos_apply_group" "testacc_barfoo2" {
  name   = junos_group_raw.testacc_barfoo.name
  prefix = "snmp view \"${junos_snmp_view.testacc_barfoo.name}\" "
}

