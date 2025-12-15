resource "junos_interface_physical" "testacc_actioncommitfile" {
  name         = var.interface
  description  = "testacc_null"
  vlan_tagging = true
}

resource "local_file" "hostname" {
  depends_on = [
    junos_interface_physical.testacc_actioncommitfile,
  ]
  content  = "set interfaces ${var.interface} description testacc_action"
  filename = var.file

  lifecycle {
    action_trigger {
      events  = [after_create, after_update]
      actions = [action.junos_commit_file.setfile]
    }
  }
}

action "junos_commit_file" "setfile" {
  provider = junos.fake

  config {
    filename                = var.file
    clear_file_after_commit = true
  }
}
