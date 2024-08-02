resource "junos_system_login_user" "testacc" {
  name       = "testacc"
  class      = "unauthorized"
  cli_prompt = "test cli"
  full_name  = "test name"
  authentication {
    encrypted_password = "test"
    no_public_keys     = true
  }
}
resource "junos_system_login_user" "testacc2" {
  name  = "test.acc2"
  class = "unauthorized"
}
resource "junos_system_login_user" "testacc3" {
  name  = "test.acc3"
  class = "unauthorized"
  authentication {
    plain_text_password = "test1234"
  }
}
