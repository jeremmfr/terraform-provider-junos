resource "junos_services_security_intelligence_profile" "testacc_svcSecIntelProfile" {
  name     = "testacc_svcSecIntelProfile@1"
  category = "CC"
  default_rule_then {
    action = "permit"
    no_log = true
  }
  rule {
    name = "test#3"
    match {
      threat_level = [5, 4]
      feed_name    = ["CC URL"]
    }
    then_action = "permit"
    then_log    = true
  }
  rule {
    name = "test"
    match {
      threat_level = [1]
    }
    then_action = "recommended"
  }
  rule {
    name = "test#2"
    match {
      threat_level = [10]
      feed_name    = ["CC_IP"]
    }
    then_action = "block close http redirect-url http://www.test.com/url1.html"
    then_log    = true
  }
  rule {
    name = "test2"
    match {
      threat_level = [10]
    }
    then_action = "block drop"
  }
}
resource "junos_services_security_intelligence_profile" "testacc_svcSecIntelProfile2" {
  name     = "testacc_svcSecIntelProfile@2"
  category = "DNS"
  default_rule_then {
    action = "sinkhole"
    no_log = true
  }
  rule {
    name = "test 1"
    match {
      threat_level = [1]
    }
    then_action = "sinkhole"
  }
}
