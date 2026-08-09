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

# read the configuration to check that the write-only password has been sent to the device,
# the resource is unchanged in this step so the data source is read during the plan
data "junos_config_raw" "testacc_evtopts_dest_wo" {
  format = "set"
}
