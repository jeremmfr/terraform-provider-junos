resource "junos_system_login_user" "testacc_wo" {
  name  = "testacc_wo"
  class = "unauthorized"
  authentication {
    encrypted_password_wo         = "test2"
    encrypted_password_wo_version = 2
  }
}
