resource "junos_applications" "testacc" {
  application {
    name                 = "testacc_apps"
    protocol             = "tcp"
    destination_port     = "22"
    application_protocol = "ssh"
    description          = "ssh protocol"
    inactivity_timeout   = 900
    source_port          = "1024-65535"
  }
  application {
    name                     = "testacc_apps2"
    protocol                 = "tcp"
    ether_type               = "0x0800"
    rpc_program_number       = "0-0"
    inactivity_timeout_never = true
    uuid                     = "AAAAA0AA-B9B0-CCcc-DDDD-EEEffFFFAAAA"
  }
  application {
    name = "testacc_apps3"
    term {
      name               = "term_B"
      protocol           = "tcp"
      destination_port   = 22
      inactivity_timeout = 600
      source_port        = "1024-65535"
    }
    term {
      name     = "term_ALG"
      protocol = "tcp"
      alg      = "ssh"
    }
  }
  application {
    name = "testacc_apps4"
    term {
      name                     = "term_B"
      protocol                 = "tcp"
      rpc_program_number       = "1-1"
      inactivity_timeout_never = true
      uuid                     = "BBBAA0AA-B9B0-CCcc-DDDD-EEEffFFFAAAA"
    }
  }
  application {
    name = "testacc_apps5"
    term {
      name      = "term_I"
      protocol  = "icmp"
      icmp_code = "1"
      icmp_type = "echo-reply"
    }
    term {
      name       = "term_I6"
      protocol   = "icmp6"
      icmp6_code = "1"
      icmp6_type = "echo-reply"
    }
  }
  application {
    name      = "testacc_apps6"
    protocol  = "icmp"
    icmp_code = "1"
    icmp_type = "echo-reply"
  }
  application {
    name       = "testacc_apps7"
    protocol   = "icmp6"
    icmp6_code = "1"
    icmp6_type = "echo-reply"
  }
  application {
    name                                   = "testacc_apps8"
    application_protocol                   = "dns"
    do_not_translate_a_query_to_aaaa_query = true
  }
  application {
    name                                   = "testacc_apps9"
    application_protocol                   = "dns"
    do_not_translate_aaaa_query_to_a_query = true
  }

  application_set {
    name         = "testacc_apps_set1"
    applications = ["junos-ssh", "junos-telnet"]
    application_set = [
      "testacc_apps_set2"
    ]
  }
  application_set {
    name         = "testacc_apps_set2"
    applications = ["junos-ftp"]
    description  = "testacc appsets2"
  }
}
