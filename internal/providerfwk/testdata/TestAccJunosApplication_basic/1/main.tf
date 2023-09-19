resource "junos_application" "testacc_app" {
  name             = "testacc_app"
  protocol         = "tcp"
  destination_port = 22
}
resource "junos_application" "testacc_app3" {
  name = "testacc_app3"
  term {
    name             = "term_B"
    protocol         = "tcp"
    destination_port = 22
  }
}
