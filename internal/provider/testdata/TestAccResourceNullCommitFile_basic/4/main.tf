resource "junos_interface_physical" "testacc_nullcommitfile" {
  name         = var.interface
  description  = "testacc_null"
  vlan_tagging = true
}

resource "local_file" "hostname" {
  content  = "set interfaces ${var.interface} description testacc_nullfile"
  filename = var.file
}

resource "junos_null_commit_file" "testacc_nullcommitfile" {
  filename                = local_file.hostname.filename
  clear_file_after_commit = true
}
