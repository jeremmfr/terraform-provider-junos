resource "junos_security_screen" "testacc_securityScreen" {
  name               = "testacc 1"
  alarm_without_drop = true
  description        = "desc testacc 1"
  icmp {
    flood {}
    fragment         = true
    icmpv6_malformed = true
    large            = true
    ping_death       = true
    sweep {}
  }
  ip {
    bad_option = true
    block_frag = true
    ipv6_extension_header {
      ah_header  = true
      esp_header = true
      hip_header = true
      destination_header {}
      fragment_header = true
      hop_by_hop_header {}
      mobility_header = true
      no_next_header  = true
      routing_header  = true
      shim6_header    = true
      user_defined_header_type = [
        "10 to 20",
        "2 to 5",
        "1",
      ]
    }
    ipv6_extension_header_limit = 32
    ipv6_malformed_header       = true
    loose_source_route_option   = true
    record_route_option         = true
    security_option             = true
    source_route_option         = true
    spoofing                    = true
    stream_option               = true
    strict_source_route_option  = true
    tear_drop                   = true
    timestamp_option            = true
    tunnel {
      bad_inner_header = true
      gre {
        gre_4in4 = true
        gre_4in6 = true
        gre_6in4 = true
        gre_6in6 = true
      }
      ip_in_udp_teredo = true
      ipip {
        ipip_4in4      = true
        ipip_4in6      = true
        ipip_6in4      = true
        ipip_6in6      = true
        ipip_6over4    = true
        ipip_6to4relay = true
        dslite         = true
        isatap         = true
      }
    }
    unknown_protocol = true
  }
  limit_session {
    destination_ip_based = 2000
    source_ip_based      = 3000
  }
  tcp {
    fin_no_ack = true
    land       = true
    no_flag    = true
    port_scan {}
    syn_ack_ack_proxy {}
    syn_fin = true
    syn_flood {
      alarm_threshold       = 10011
      attack_threshold      = 10012
      destination_threshold = 10013
      source_threshold      = 10014
      timeout               = 10
      whitelist {
        name                = "test3"
        source_address      = ["192.0.2.0/26"]
        destination_address = ["192.0.2.64/26"]
      }
    }
    syn_frag = true
    sweep {}
    winnuke = true
  }
  udp {
    flood {}
    port_scan {}
    sweep {}
  }
}
resource "junos_security_screen_whitelist" "testacc1" {
  name = "testacc1"
  address = [
    "192.0.2.128/26",
    "192.0.2.64/26",
  ]
}
