resource "junos_interface_physical" "testacc_actioncommitfile" {
  provider = junos.fake

  name         = var.interface
  description  = "testacc_fakeupdate"
  vlan_tagging = true

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
