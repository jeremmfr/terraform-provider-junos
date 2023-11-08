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

data "junos_applications" "testacc_default_any" {
  match_options {
    protocol = "0"
  }
}
data "junos_applications" "testacc_default_ssh-name" {
  match_name = "^j.*-ssh$"
}
data "junos_applications" "testacc_default_ssh" {
  match_name = "^junos-"
  match_options {
    protocol         = "tcp"
    destination_port = 22
  }
}
data "junos_applications" "testacc_all_ssh" {
  match_options {
    protocol         = "tcp"
    destination_port = 22
  }
}
data "junos_applications" "testacc_multi_term" {
  match_options {
    protocol         = "tcp"
    destination_port = 1001
  }
  match_options {
    protocol         = "tcp"
    destination_port = 1002
  }
}
