resource "junos_system_root_authentication" "root_auth_wo" {
  plain_text_password_wo         = "aPassWord1!"
  plain_text_password_wo_version = 1
}
