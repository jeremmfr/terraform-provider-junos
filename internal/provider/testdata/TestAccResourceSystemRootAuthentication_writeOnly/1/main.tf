resource "junos_system_root_authentication" "root_auth_wo" {
  encrypted_password_wo         = "$6$XXXX"
  encrypted_password_wo_version = 1
}
