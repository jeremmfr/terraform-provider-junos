resource "terraform_data" "trigger" {
  triggers_replace = "1"
  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.junos_load_config.load-application]
    }
  }
}

action "junos_load_config" "load-application" {
  config {
    config = <<EOT
applications {
    application testacc-load-config {
        protocol tcp;
        source-port 8080;
        destination-port 22;
    }
}
EOT
  }
}
