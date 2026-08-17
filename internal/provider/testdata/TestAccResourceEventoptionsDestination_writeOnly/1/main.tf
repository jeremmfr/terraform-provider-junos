resource "junos_eventoptions_destination" "testacc_evtopts_dest_wo" {
  name = "testacc_evtopts_dest_wo"
  archive_site {
    url = "https://example.com"
  }
  archive_site {
    url                 = "https://example.fr"
    password_wo         = "thePassword"
    password_wo_version = 1
  }
}
