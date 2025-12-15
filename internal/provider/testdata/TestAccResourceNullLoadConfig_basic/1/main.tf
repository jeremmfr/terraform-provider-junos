resource "junos_null_load_config" "load-application" {
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
