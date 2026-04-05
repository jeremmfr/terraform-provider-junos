resource "junos_policyoptions_as_path_group" "testacc_dataASPathGroup" {
  name = "testacc_dataASPathGroup"
  as_path {
    name = "testacc_dataASPathGroup"
    path = "5|12|18"
  }
}

data "junos_policyoptions_as_path_group" "testacc_dataASPathGroup" {
  name = junos_policyoptions_as_path_group.testacc_dataASPathGroup.name
}
