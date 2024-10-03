resource "junos_applications_ordered" "testacc" {
  application {
    name             = "testacc_app4"
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
