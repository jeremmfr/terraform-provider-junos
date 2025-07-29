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
  name   = "testacc_fwFilter#6"
  family = "inet6"
  term {
    name = "testacc_fwFilter#6 term1"
    from {
      interface     = ["fe-*"]
      next_header   = ["icmp6"]
      loss_priority = ["low"]
    }
    then {
      action = "discard"
    }
  }
}
resource "junos_firewall_filter" "testacc_fwFilter62" {
  name   = "testacc_fwFilte #62"
  family = "inet6"
  term {
    name   = "testacc_fwFilter#62 term1"
    filter = junos_firewall_filter.testacc_fwFilter6.name
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
resource "junos_firewall_policer" "testacc_fwfilter" {
  name = "testacc_fwfilter#1"
  if_exceeding {
    bandwidth_percent = 80
    burst_size_limit  = "50k"
  }
  then {
    discard = true
  }
}
resource "junos_firewall_filter" "testacc_fwFilter_vpls" {
  name   = "testacc_fwFilter vpls"
  family = "vpls"
  term {
    name = "testacc_fwFilter vpls term1"
    from {
      forwarding_class_except = [
        "network-control",
      ]
    }
    then {
      loss_priority = "high"
    }
  }
}
resource "junos_firewall_filter" "testacc_fwFilter_any" {
  name   = "testacc_fwFilter any"
  family = "any"
  term {
    name = "testacc_fwFilter any term1"
    from {
      packet_length_except = ["1-500"]
    }
    then {
      forwarding_class = "best-effort"
    }
  }
}
