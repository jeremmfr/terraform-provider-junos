resource "junos_eventoptions_generate_event" "testacc_evtopts_genevent" {
  name          = "testacc_evtopts_genevent#1"
  time_interval = 3600
  start_time    = "2024-2-18.01:02:03"
  no_drift      = true
}
