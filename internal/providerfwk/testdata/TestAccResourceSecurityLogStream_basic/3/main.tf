resource "junos_security_log_stream" "testacc_logstream" {
  name = "testacc_logstream"
  file {
    name             = "#File@test.txt"
    allow_duplicates = true
    size             = 3
    rotation         = 3
  }
  filter_threat_attack = true
}
resource "junos_security_log_stream" "testacc_logstream2" {
  name = "testacc_logstream2"
  transport {}
}
