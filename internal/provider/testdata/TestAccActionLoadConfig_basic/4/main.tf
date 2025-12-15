resource "terraform_data" "trigger" {
  triggers_replace = "4"
  lifecycle {
    action_trigger {
      events = [before_create]
      actions = [
        action.junos_load_config.load-application,
        action.junos_load_config.load-application2,
      ]
    }
  }
}

action "junos_load_config" "load-application" {
  config {
    action = "set"
    format = "text"
    config = <<EOT
delete applications application testacc-load-config
EOT
  }
}

action "junos_load_config" "load-application2" {
  config {
    action = "merge"
    format = "xml"
    config = <<EOT
<applications>
    <application>
        <name>test_acc-load-config</name>
        <protocol>tcp</protocol>
        <destination-port>22</destination-port>
    </application>
</applications>
EOT
  }
}

data "junos_applications" "testacc" {
  match_name = "^testacc.*$"
}
