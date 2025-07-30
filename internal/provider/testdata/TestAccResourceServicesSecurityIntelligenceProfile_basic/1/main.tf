resource "junos_services_security_intelligence_profile" "testacc_svcSecIntelProfile" {
  name     = "testacc_svcSecIntelProfile@1"
  category = "CC"
  rule {
    name = "test#2"
    match {
      threat_level = [10]
      feed_name    = ["CC_IP"]
    }
    then_action = "block close http redirect-url http://www.test.com/url1.html"
    then_log    = true
  }
}
