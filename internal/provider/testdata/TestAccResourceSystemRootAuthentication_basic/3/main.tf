resource "junos_system_root_authentication" "root_auth" {
  encrypted_password = "$6$XXX"
  no_public_keys     = true
}
