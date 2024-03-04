resource "junos_routing_instance" "testacc_radiusServer" {
  name = "testacc_radiusServer"
}
resource "junos_system_radius_server" "testacc_radiusServer" {
  address                  = "192.0.2.1"
  secret                   = "password"
  preauthentication_secret = "password"
  source_address           = "192.0.2.2"
  port                     = 1645
  accounting_port          = 1646
  dynamic_request_port     = 3799
  preauthentication_port   = 1812
  timeout                  = 10
  accounting_timeout       = 5
  retry                    = 3
  accounting_retry         = 2
  max_outstanding_requests = 1000
  routing_instance         = junos_routing_instance.testacc_radiusServer.name
}
