resource "junos_policyoptions_community" "testacc_dataCommunity" {
  name    = "testacc_dataCommunity"
  members = ["65000:100"]
}

data "junos_policyoptions_community" "testacc_dataCommunity" {
  name = junos_policyoptions_community.testacc_dataCommunity.name
}
