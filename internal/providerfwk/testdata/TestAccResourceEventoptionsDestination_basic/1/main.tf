resource "junos_eventoptions_destination" "testacc_evtopts_dest" {
  name = "testacc_evtopts_dest#1"
  archive_site {
    url = "https://example.com"
  }
}
