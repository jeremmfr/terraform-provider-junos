resource "junos_virtual_chassis" "testacc_virtual_chassis" {
  member {
    id                  = 2
    mastership_priority = 10
  }

  auto_sw_update              = true
  auto_sw_update_package_name = "http://package.name/sw"
  graceful_restart_disable    = true
  identifier                  = "5345.6ac8.6ac8"
  mac_persistence_timer       = "disable"
  vcp_no_hold_time            = true
  traceoptions {
    flag = [
      "hello detail",
      "heartbeat",
    ]
    file {
      name           = "trace_#VC"
      files          = 100
      no_stamp       = true
      replace        = true
      size           = 102400
      world_readable = true
    }
  }
}
