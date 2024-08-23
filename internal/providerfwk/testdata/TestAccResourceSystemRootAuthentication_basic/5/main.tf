resource "junos_system_root_authentication" "root_auth" {
  plain_text_password = "testPassword1234"
}


resource "junos_system_root_authentication" "root_auth_copy" {
  depends_on = [junos_system_root_authentication.root_auth]

  plain_text_password = "testPassword5678"
}
