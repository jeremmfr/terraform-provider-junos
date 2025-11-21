resource "junos_null_load_config" "load-application" {
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

data "junos_applications" "testacc" {
  match_name = "^testacc.*$"
}
