data "junos_system_information" "srx" {}
locals {
  content_filtering_rule_set_available = tonumber(replace(data.junos_system_information.srx.os_version, "/\\..*$/", "")) >= 22 ? 1 : 0
}

resource "junos_security_utm_policy" "testacc_Policy" {
  count = local.content_filtering_rule_set_available

  name = "testacc Policy"
  content_filtering_rule_set {
    name = "cf rule-set #P"
    rule {
      name = "cf rule #33 in rule-set #P"

      match_applications = ["http"]
      match_direction    = "any"
      match_file_types   = ["zip"]
      then_action        = "close-server"
    }
  }
}
