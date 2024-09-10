resource "tls_private_key" "rsa4096" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "junos_system_root_authentication" "root_auth" {
  encrypted_password = "$6$XXXX"
  ssh_public_keys    = [chomp(tls_private_key.rsa4096.public_key_openssh)]
}
