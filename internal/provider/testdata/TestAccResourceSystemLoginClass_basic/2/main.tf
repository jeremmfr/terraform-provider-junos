resource "junos_system_login_class" "testacc" {
  name                        = "testacc"
  access_start                = "08:00:00"
  access_end                  = "18:00:00"
  allow_commands_regexps      = [".*"]
  allow_configuration_regexps = [".*"]
  no_hidden_commands_except   = [".*"]
  deny_commands_regexps       = ["request"]
  deny_configuration_regexps  = ["system"]
  idle_timeout                = 120
  login_alarms                = true
  login_tip                   = true
  permissions                 = ["view"]
  security_role               = "security-administrator"
}
