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
      match_direction    = "upload"
      match_file_types   = ["zip", "bzip"]
      then_action        = "close-server"
      then_notification_endpoint {
        type = "message"
      }
    }
    rule {
      name               = "rule1"
      match_applications = ["smtp", "http"]
      match_direction    = "upload"
      match_file_types   = ["jmp", "bzip"]
      then_notification_endpoint {
        notify_mail_sender = false
      }
    }
  }
  content_filtering_rule_set {
    name = "cf rule-set #A"
    rule {
      name = "cf rule #44 in rule-set #A"

      match_applications = ["any"]
      match_direction    = "download"
      match_file_types   = ["zip"]
      then_action        = "close-client"
      then_notification_endpoint {
        custom_message     = "a cust@m me$$age"
        notify_mail_sender = true
        type               = "protocol-only"
      }
    }
    rule {
      name = "cf rule #22 in rule-set #A"

      match_applications    = ["any"]
      match_direction       = "any"
      match_file_types      = ["applesingle", "flash", "paquet"]
      then_notification_log = true
    }
    rule {
      name               = "rule1"
      match_applications = ["smtp", "http"]
      match_direction    = "upload"
      match_file_types   = ["jmp", "bzip"]
      then_notification_endpoint {
        notify_mail_sender = false
      }
    }
  }
}
