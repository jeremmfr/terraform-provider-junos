resource "junos_firewall_filter" "testacc_fwFilter" {
  name               = "testacc_fwFilter"
  family             = "inet"
  interface_specific = true
  term {
    name = "testacc_fwFilter_term1"
    from {
      address            = ["192.0.2.0/25"]
      address_except     = ["192.0.2.128/25"]
      port               = ["22-23"]
      prefix_list        = [junos_policyoptions_prefix_list.testacc_fwFilter.name]
      prefix_list_except = [junos_policyoptions_prefix_list.testacc_fwFilter2.name]
      protocol           = ["tcp"]
      tcp_flags          = "!0x3"
      is_fragment        = true
    }
    then {
      action             = "next term"
      syslog             = true
      log                = true
      packet_mode        = true
      port_mirror        = true
      service_accounting = true
    }
  }
  term {
    name = "testacc_fwFilter_term2"
    from {
      icmp_code = ["network-unreachable"]
      icmp_type = ["router-advertisement"]
    }
    then {
      action = "accept"
    }
  }
}
resource "junos_firewall_filter" "testacc_fwFilter_vpls" {
  name   = "testacc_fwFilter vpls"
  family = "vpls"
  term {
    name = "testacc_fwFilter vpls term1"
    from {
      destination_mac_address = [
        "aa:bb:cc:dd:ee:ff/48",
      ]
      destination_mac_address_except = [
        "aa:bb:cc:dd:ee:f0/48",
      ]
      forwarding_class = [
        "best-effort",
      ]
      source_mac_address_except = [
        "aa:bb:cc:dd:ee:01/48",
      ]
      source_mac_address = [
        "aa:bb:cc:dd:ee:02/48",
      ]
    }
  }
}
resource "junos_firewall_filter" "testacc_fwFilter_any" {
  name   = "testacc_fwFilter any"
  family = "any"
  term {
    name = "testacc_fwFilter any term1"
    from {
      packet_length        = ["1-500"]
      loss_priority_except = ["medium-high"]
    }
  }
}

resource "junos_policyoptions_prefix_list" "testacc_fwFilter" {
  name   = "testacc_fwFilter#1"
  prefix = ["192.0.2.0/25"]
}
resource "junos_policyoptions_prefix_list" "testacc_fwFilter2" {
  name   = "testacc_fwFilter#2"
  prefix = ["192.0.2.128/25"]
}
