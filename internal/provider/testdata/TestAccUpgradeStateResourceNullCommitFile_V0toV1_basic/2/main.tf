resource "local_file" "timezone" {
  content  = "set system time-zone Europe/Paris"
  filename = var.file
}

resource "junos_null_commit_file" "testacc_nullcommitfile" {
  filename = local_file.timezone.filename
  triggers = {
    md5 = filemd5(var.file)
  }
}
