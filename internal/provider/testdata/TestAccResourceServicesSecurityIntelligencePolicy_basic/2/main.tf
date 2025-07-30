resource "junos_services_security_intelligence_profile" "testacc_svcSecIntelPolicy_CC" {
  name     = "testacc svcSecIntelPolicy_CC"
  category = "CC"
  rule {
    name = "rule_1"
    match {
      threat_level = [1]
    }
    then_action = "permit"
  }
}
resource "junos_services_security_intelligence_profile" "testacc_svcSecIntelPolicy_IPFilter" {
  name     = "testacc_svcSecIntelPolicy_IPFilter"
  category = "IPFilter"
  rule {
    name = "rule_1"
    match {
      threat_level = [1]
    }
    then_action = "permit"
  }
}
resource "junos_services_security_intelligence_policy" "testacc_svcSecIntelPolicy" {
  name = "testacc_svcSecIntelPolicy#1"
  category {
    name         = "CC"
    profile_name = junos_services_security_intelligence_profile.testacc_svcSecIntelPolicy_CC.name
  }
  category {
    name         = "IPFilter"
    profile_name = junos_services_security_intelligence_profile.testacc_svcSecIntelPolicy_IPFilter.name
  }
}
