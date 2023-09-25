resource "junos_application" "testacc_app" {
  name                 = "testacc_app"
  protocol             = "tcp"
  destination_port     = "22"
  application_protocol = "ssh"
  description          = "ssh protocol"
  inactivity_timeout   = 900
  source_port          = "1024-65535"
}
resource "junos_application" "testacc_app2" {
  name                     = "testacc_app2"
  protocol                 = "tcp"
  ether_type               = "0x0800"
  rpc_program_number       = "0-0"
  inactivity_timeout_never = true
  uuid                     = "AAAAA0AA-B9B0-CCcc-DDDD-EEEffFFFAAAA"
}
resource "junos_application" "testacc_app3" {
  name = "testacc_app3"
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
resource "junos_application" "testacc_app4" {
  name = "testacc_app4"
  term {
    name                     = "term_B"
    protocol                 = "tcp"
    rpc_program_number       = "1-1"
    inactivity_timeout_never = true
    uuid                     = "BBBAA0AA-B9B0-CCcc-DDDD-EEEffFFFAAAA"
  }
}
resource "junos_application" "testacc_app5" {
  name = "testacc_app5"
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
