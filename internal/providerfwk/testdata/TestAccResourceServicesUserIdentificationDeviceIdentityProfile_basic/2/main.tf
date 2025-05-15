resource "junos_services_user_identification_device_identity_profile" "testacc_devidProf" {
  name   = "testacc_devidProf.1"
  domain = "clearpass"
  attribute {
    name  = "device-identity"
    value = ["device1", "barcode scan"]
  }
  attribute {
    name  = "device-category"
    value = ["category@1"]
  }
}
