resource "junos_interface_physical" "testacc_nullcommitfile" {
  provider = junos.fake

  name         = var.interface
  description  = "testacc_fakeupdate"
  vlan_tagging = true
}

resource "junos_null_commit_file" "setfile2" {
  provider = junos.fake

  depends_on = [
    junos_interface_physical.testacc_nullcommitfile
  ]

  filename                = var.file
  clear_file_after_commit = true
}
