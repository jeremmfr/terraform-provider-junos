resource "junos_policyoptions_source_address_filter_list" "testacc_dataSAFList" {
  name = "testacc_dataSAFList"
  address {
    address = "192.0.2.0/25"
    option  = "orlonger"
  }
}

data "junos_policyoptions_source_address_filter_list" "testacc_dataSAFList" {
  name = junos_policyoptions_source_address_filter_list.testacc_dataSAFList.name
}
