resource "junos_policyoptions_route_filter_list" "testacc_dataRFList" {
  name = "testacc_dataRFList"
  address {
    address = "192.0.2.0/25"
    option  = "orlonger"
  }
}
