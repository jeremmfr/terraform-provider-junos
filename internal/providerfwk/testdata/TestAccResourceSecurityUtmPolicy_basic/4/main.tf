resource "junos_security_utm_policy" "testacc_Policy" {
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
