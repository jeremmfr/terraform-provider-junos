resource "local_file" "timezone" {
  content  = "set system time-zone Europe/Paris"
  filename = var.file
}
