data "junos_applications" "testacc" {
  match_name = "^testacc.*$"
}

data "junos_applications" "test_acc" {
  match_name = "^test_acc.*$"
}

resource "terraform_data" "trigger" {
  triggers_replace = "5"
  lifecycle {
    action_trigger {
      events = [before_create]
      actions = [
        action.junos_load_config.load-application2,
      ]
    }
  }
}

action "junos_load_config" "load-application2" {
  config {
    action = "set"
    format = "text"
    config = <<EOT
delete applications application test_acc-load-config
EOT
  }
}
