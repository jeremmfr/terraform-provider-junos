resource "junos_forwardingoptions_sampling_instance" "testacc_v0toV1_sampInstance" {
  depends_on = [
    junos_interface_logical.testacc_v0toV1_sampInstance,
    junos_system_ntp_server.testacc_v0toV1_sampInstance,
  ]

  name = "testacc_v0toV1_sampInstance"
  input {
    rate = 1
  }
  family_inet_output {
    flow_server {
      hostname = "192.0.2.1"
      port     = 3000
    }
    interface {
      name           = "si-0/1/0"
      source_address = "192.0.2.2"
    }
  }
}
resource "junos_system_ntp_server" "testacc_v0toV1_sampInstance" {
  address = "192.0.2.3"
}
resource "junos_interface_logical" "testacc_v0toV1_sampInstance" {
  name = "si-0/1/0.0"
  family_inet {}
}
