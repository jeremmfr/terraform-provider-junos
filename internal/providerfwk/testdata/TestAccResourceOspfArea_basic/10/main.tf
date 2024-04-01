resource "junos_ospf_area" "testacc_ospfarea" {
  area_id = "1"
  version = "v3"
  interface {
    name = "all"
  }
  area_range {
    range = "fe80:f::/64"
    exact = true
  }
  area_range {
    range           = "fe80:e::/64"
    exact           = true
    override_metric = 106
  }
  area_range {
    range    = "fe80::/64"
    restrict = true
  }
  context_identifier = ["127.0.0.2", "127.0.0.1"]
  inter_area_prefix_export = [
    junos_policyoptions_policy_statement.testacc_ospfarea.name,
    junos_policyoptions_policy_statement.testacc_ospfarea2.name,
  ]
  inter_area_prefix_import = [
    junos_policyoptions_policy_statement.testacc_ospfarea2.name,
    junos_policyoptions_policy_statement.testacc_ospfarea.name,
  ]
  nssa {
    area_range {
      range = "fe80::/64"
      exact = true
    }
    area_range {
      range           = "fe80:b::/64"
      override_metric = 107
    }
    area_range {
      range    = "fe80:a::/64"
      restrict = true
    }
    default_lsa {
      default_metric = 109
      metric_type    = 2
      type_7         = true
    }
    summaries = true
  }
}
resource "junos_policyoptions_policy_statement" "testacc_ospfarea" {
  name = "testacc_ospfarea"
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_ospfarea2" {
  name = "testacc_ospfarea2"
  then {
    action = "reject"
  }
}
resource "junos_ospf_area" "testacc_ospfarea2" {
  area_id = "2"
  version = "v3"
  interface {
    name    = "${var.interface}.0"
    passive = true
  }
  stub {
    default_metric = 150
    no_summaries   = true
  }
}
resource "junos_ospf_area" "testacc_ospfarea3" {
  area_id = "3"
  version = "v3"
  realm   = "ipv4-unicast"
  interface {
    name    = "${var.interface}.0"
    passive = true
  }
  nssa {
    no_summaries = true
    default_lsa {}
  }
}
