resource "junos_null_load_config" "load-application" {
  action = "set"
  format = "text"
  config = <<EOT
delete applications application testacc-load-config
EOT
}

resource "junos_null_load_config" "load-application2" {
  depends_on = [
    junos_null_load_config.load-application,
  ]
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


data "junos_applications" "testacc" {
  match_name = "^testacc.*$"
}
