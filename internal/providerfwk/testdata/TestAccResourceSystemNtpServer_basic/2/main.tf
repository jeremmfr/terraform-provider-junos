resource "junos_routing_instance" "testacc_ntpServer" {
  name = "testacc_ntpServer"
}
resource "junos_system_ntp_server" "testacc_ntpServer" {
  address          = "192.0.2.1"
  prefer           = true
  version          = 4
  routing_instance = junos_routing_instance.testacc_ntpServer.name
}


resource "junos_system_ntp_server" "testacc_ntpServer2" {
  address = "192.0.2.2"
  nts {
    remote_identity_distinguished_name_container = " test acc"
  }
}

resource "junos_system_ntp_server" "testacc_ntpServer3" {
  address = "192.0.2.3"
  nts {
    remote_identity_distinguished_name_wildcard = " test acc"
  }
}

resource "junos_system_ntp_server" "testacc_ntpServer4" {
  address = "192.0.2.4"
  nts {
    remote_identity_hostname = " test acc"
  }
}
