data "junos_applications" "testacc" {
  match_name = "^testacc.*$"
}

data "junos_applications" "test_acc" {
  match_name = "^test_acc.*$"
}
