resource "junos_eventoptions_policy" "testacc_evtopts_policy" {
  name   = "testacc_evtopts_policy#1"
  events = ["aaa_infra_fail", "acct_fork_err"]
  then {
    ignore = true
  }
}
