resource "terraform_data" "trigger" {
  triggers_replace = "2"
  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.junos_load_config.load-application]
    }
  }
}

action "junos_load_config" "load-application" {
  config {
    action = "replace"
    config = <<EOT
applications {
    replace: application testacc-load-config {
        protocol tcp;
        destination-port 22;
    }
}
EOT
  }
}

data "junos_applications" "testacc" {
  match_name = "^testacc.*$"
}
