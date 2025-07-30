resource "junos_virtual_chassis" "testacc_virtual_chassis" {
  member {
    id            = 0
    serial_number = var.serial_number
    role          = "routing-engine"
  }
  member {
    id            = 2
    serial_number = "ABC123"
    role          = "line-card"
  }
  member {
    id                 = 1
    serial_number      = "9876"
    role               = "line-card"
    no_management_vlan = true
    location           = "In House"
  }

  alias {
    serial_number = "666"
    alias_name    = "Evil"
  }
  alias {
    serial_number = "112"
    alias_name    = "Emergency"
  }

  auto_sw_update        = true
  preprovisioned        = true
  identifier            = "9622.6ac8.5345"
  mac_persistence_timer = "30"
  traceoptions {
    file {
      name              = "virtualChassis"
      no_world_readable = true
    }
  }
}
