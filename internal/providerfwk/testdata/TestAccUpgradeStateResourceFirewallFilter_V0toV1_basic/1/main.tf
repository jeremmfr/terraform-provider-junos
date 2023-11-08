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
    }
    then {
      action             = "next term"
      syslog             = true
      log                = true
      port_mirror        = true
      service_accounting = true
    }
  }
  term {
    name = "testacc_fwFilter_term2"
    from {
      source_address            = ["192.0.2.0/25"]
      source_address_except     = ["192.0.2.128/25"]
      port_except               = ["23"]
      source_prefix_list        = [junos_policyoptions_prefix_list.testacc_fwFilter.name]
      source_prefix_list_except = [junos_policyoptions_prefix_list.testacc_fwFilter2.name]
      tcp_established           = true
      protocol_except           = ["icmp"]
    }
    then {
      policer = junos_firewall_policer.testacc_fwfilter.name
      action  = "accept"
    }
  }
  term {
    name = "testacc_fwFilter_term3"
    from {
      destination_address            = ["192.0.2.0/25"]
      destination_address_except     = ["192.0.2.128/25"]
      destination_port               = ["22-23"]
      source_port_except             = ["23"]
      destination_prefix_list        = [junos_policyoptions_prefix_list.testacc_fwFilter.name]
      destination_prefix_list_except = [junos_policyoptions_prefix_list.testacc_fwFilter2.name]
      tcp_initial                    = true
    }
    then {
      action = "discard"
    }
  }
  term {
    name = "testacc_fwFilter_term4"
    from {
      source_port             = ["22-23"]
      destination_port_except = ["23"]
    }
    then {
      action = "reject"
    }
  }
  term {
    name = "testacc_fwFilter_term5"
    from {
      icmp_code_except = ["network-unreachable"]
      icmp_type_except = ["router-advertisement"]
    }
    then {
      action = "reject"
    }
  }
}
resource "junos_firewall_filter" "testacc_fwFilter6" {
  name   = "testacc_fwFilter6"
  family = "inet6"
  term {
    name = "testacc_fwFilter6_term1"
    from {
      next_header = ["icmp6"]
    }
    then {
      action = "discard"
    }
  }
}
resource "junos_policyoptions_prefix_list" "testacc_fwFilter" {
  name   = "testacc_fwFilter"
  prefix = ["192.0.2.0/25"]
}
resource "junos_policyoptions_prefix_list" "testacc_fwFilter2" {
  name   = "testacc_fwFilter2"
  prefix = ["192.0.2.128/25"]
}
resource "junos_firewall_policer" "testacc_fwfilter" {
  name = "testacc_fwfilter"
  if_exceeding {
    bandwidth_percent = 80
    burst_size_limit  = "50k"
  }
  then {
    discard = true
  }
}
