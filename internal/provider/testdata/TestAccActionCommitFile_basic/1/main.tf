resource "junos_interface_physical" "testacc_actioncommitfile" {
  provider = junos.fake

  name         = var.interface
  description  = "testacc_fakecreate"
  vlan_tagging = true
}

// a empty resource between junos resource(s) and action(s) due to
// a bug with lifecycle.action_trigger.events (see https://github.com/hashicorp/terraform/issues/37930)
resource "terraform_data" "setfile" {
  depends_on = [
    junos_interface_physical.testacc_actioncommitfile,
  ]
  triggers_replace = junos_interface_physical.testacc_actioncommitfile.description
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
