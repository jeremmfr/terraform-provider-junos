resource "junos_system_login_user" "testacc3" {
  name  = "test.acc3"
  class = "unauthorized"
  authentication {
    plain_text_password = "test1234"
  }
}

import {
  to = junos_system_login_user.testacc3_copy
  id = "test.acc3"
}

resource "junos_system_login_user" "testacc3_copy" {
  name  = "test.acc3"
  class = "unauthorized"
  authentication {
    plain_text_password = "test4567"
  }
}
