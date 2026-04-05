resource "junos_policyoptions_as_path" "testacc_dataASPath" {
  name = "testacc_dataASPath"
  path = "5|12|18"
}

data "junos_policyoptions_as_path" "testacc_dataASPath" {
  name = junos_policyoptions_as_path.testacc_dataASPath.name
}
