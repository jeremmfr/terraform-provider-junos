resource "junos_multichassis" "testacc_multichassis" {
  clean_on_destroy                               = true
  mc_lag_consistency_check_comparison_delay_time = 180
}
