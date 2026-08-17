resource "junos_system_login_user" "testacc_wo" {
  name  = "testacc_wo"
  class = "unauthorized"
  authentication {
    plain_text_password_wo         = "test1234"
    plain_text_password_wo_version = 1
  }
}
