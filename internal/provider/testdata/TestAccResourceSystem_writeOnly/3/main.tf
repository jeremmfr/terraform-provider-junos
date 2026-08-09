resource "junos_system" "testacc_system_wo" {
  host_name = "testacc-terraform-wo"

  accounting {
    events              = ["login"]
    destination_radius  = true
    destination_tacplus = true
    destination_radius_server {
      address                             = "192.0.2.53"
      secret_wo                           = "aSecret2"
      secret_wo_version                   = 2
      preauthentication_secret_wo         = "aPreauthenticationSecret2"
      preauthentication_secret_wo_version = 2
    }
    destination_tacplus_server {
      address           = "192.0.2.55"
      secret_wo         = "aTacplusSecret2"
      secret_wo_version = 2
    }
  }
  archival_configuration {
    transfer_on_commit = true
    archive_site {
      url                 = "scp://juniper-configs@192.0.2.30:/dir"
      password_wo         = "anArchivePassword2"
      password_wo_version = 2
    }
  }
  license {
    autoupdate                     = true
    autoupdate_url                 = "https://ae1.juniper.net/junos/key_retrieval"
    autoupdate_password_wo         = "aLicensePassword2"
    autoupdate_password_wo_version = 2
  }
  name_server_opts {
    address = "192.0.2.10"
  }
  name_server_opts {
    address = "192.0.2.11"
  }
  services {
    ssh {
      root_login = "allow"
    }
  }
  time_zone = "Europe/Paris"
}
