provider "junos" {
  alias                    = "fake"
  fake_create_with_setfile = var.file
}
