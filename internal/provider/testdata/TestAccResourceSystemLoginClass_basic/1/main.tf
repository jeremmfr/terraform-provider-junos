resource "junos_system_login_class" "testacc" {
  name                      = "testacc"
  access_start              = "08:00:00"
  access_end                = "18:00:00"
  allow_commands            = ".*"
  allow_configuration       = ".*"
  allow_hidden_commands     = true
  allowed_days              = ["sunday", "monday"]
  cli_prompt                = "prompt cli"
  configuration_breadcrumbs = true
  confirm_commands          = ["confirm commands"]
  deny_commands             = "request"
  deny_configuration        = "system"
  idle_timeout              = 120
  login_alarms              = true
  login_tip                 = true
  permissions               = ["view", "floppy"]
  security_role             = "security-administrator"
}
