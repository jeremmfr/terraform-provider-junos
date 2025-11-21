data "junos_applications" "testacc" {
  match_name = "^testacc.*$"
}

data "junos_applications" "test_acc" {
  match_name = "^test_acc.*$"
}

resource "junos_null_load_config" "load-application2" {
  action = "set"
  format = "text"
  config = <<EOT
delete applications application test_acc-load-config
EOT
}
