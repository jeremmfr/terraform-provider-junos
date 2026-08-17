resource "junos_system" "testacc_system_wo" {
  host_name = "testacc-terraform-wo"

  accounting {
    events              = ["login"]
    destination_radius  = true
    destination_tacplus = true
    destination_radius_server {
      address                             = "192.0.2.53"
      secret_wo                           = "aSecret"
      secret_wo_version                   = 1
      preauthentication_secret_wo         = "aPreauthenticationSecret"
      preauthentication_secret_wo_version = 1
    }
    destination_tacplus_server {
      address           = "192.0.2.55"
      secret_wo         = "aTacplusSecret"
      secret_wo_version = 1
    }
  }
  archival_configuration {
    transfer_on_commit = true
    archive_site {
      url                 = "scp://juniper-configs@192.0.2.30:/dir"
      password_wo         = "anArchivePassword"
      password_wo_version = 1
    }
  }
  license {
    autoupdate                     = true
    autoupdate_url                 = "https://ae1.juniper.net/junos/key_retrieval"
    autoupdate_password_wo         = "aLicensePassword"
    autoupdate_password_wo_version = 1
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
    web_management_http {
      interface = ["fxp0.0"]
      port      = 80
    }
    web_management_https {
      interface                    = ["fxp0.0"]
      system_generated_certificate = true
      port                         = 443
    }
  }
  time_zone = "Europe/Paris"
}
