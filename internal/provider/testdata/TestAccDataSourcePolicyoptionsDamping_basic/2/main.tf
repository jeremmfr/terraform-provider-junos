resource "junos_policyoptions_damping" "testacc_dataDamping" {
  name         = "testacc_dataDamping"
  half_life    = 15
  max_suppress = 60
  reuse        = 750
  suppress     = 3000
}

data "junos_policyoptions_damping" "testacc_dataDamping" {
  name = junos_policyoptions_damping.testacc_dataDamping.name
}
