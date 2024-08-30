resource "junos_applications" "testacc" {
  application {
    name             = "testacc_app"
    protocol         = "tcp"
    destination_port = 22
  }
  application {
    name = "testacc_app3"
    term {
      name             = "term_B"
      protocol         = "tcp"
      destination_port = 22
    }
  }
  application_set {
    name         = "testacc_app_set"
    applications = ["junos-ssh"]
  }
}
