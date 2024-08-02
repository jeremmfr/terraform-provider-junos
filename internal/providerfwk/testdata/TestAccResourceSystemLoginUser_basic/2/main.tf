resource "tls_private_key" "rsa4096" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "junos_system_login_user" "testacc" {
  name  = "testacc"
  class = "unauthorized"
  authentication {
    encrypted_password = "test"
    ssh_public_keys    = [chomp(tls_private_key.rsa4096.public_key_openssh)]
  }
}

resource "junos_system_login_user" "testacc2" {
  name  = "test.acc2"
  class = "unauthorized"
  uid   = 5000
}

resource "junos_system_login_user" "testacc3" {
  name  = "test.acc3"
  class = "unauthorized"
  authentication {
    plain_text_password = "test1234"
  }
}
