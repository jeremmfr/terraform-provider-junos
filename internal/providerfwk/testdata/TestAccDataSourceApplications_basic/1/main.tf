resource "junos_application" "testacc_custom_ssh" {
  name             = "testacc_custom_ssh"
  protocol         = "tcp"
  destination_port = 22
}
resource "junos_application" "testacc_custom_ssh_term" {
  name = "testacc_custom_ssh_term"
  term {
    name             = "1"
    protocol         = "tcp"
    destination_port = 22
  }
}
resource "junos_application" "testacc_custom_multi_term" {
  name = "testacc_custom_multi_term"
  term {
    name             = "1"
    protocol         = "tcp"
    destination_port = 1001
  }
  term {
    name             = "2"
    protocol         = "tcp"
    destination_port = 1002
  }
}
