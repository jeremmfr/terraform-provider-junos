resource "junos_security" "testacc_secIdpPolicy" {
  idp_sensor_configuration {
    packet_log {
      source_address = "192.0.2.4"
      host_address   = "192.0.2.5"
      host_port      = 514
    }
  }
  alg {
    dns_disable    = true
    ftp_disable    = true
    h323_disable   = true
    mgcp_disable   = true
    msrpc_disable  = true
    pptp_disable   = true
    rsh_disable    = true
    rtsp_disable   = true
    sccp_disable   = true
    sip_disable    = true
    sql_disable    = true
    sunrpc_disable = true
    talk_disable   = true
    tftp_disable   = true
  }
}
resource "junos_security_idp_policy" "testacc_idp_pol" {
  name = "testacc_idp/#1"
  ips_rule {
    name        = "rules_#B"
    description = "rules _ test #B"
    match {
      application = "junos:telnet"
    }
    then {
      action                      = "class-of-service"
      action_cos_forwarding_class = "best-effort"
      action_dscp_code_point      = 3
    }
  }
  ips_rule {
    name        = "rules_#1"
    description = "rules _ test #1"
    match {
      application = "junos:ssh"
    }
    then {
      action                 = "mark-diffserv"
      action_dscp_code_point = 4
    }
  }
  exempt_rule {
    name        = "rules_#A"
    description = "rules _ test #A"
    match {
      destination_address_except = ["192.0.2.1/32"]
      source_address_except      = ["192.0.2.254/32"]
    }
  }
}
