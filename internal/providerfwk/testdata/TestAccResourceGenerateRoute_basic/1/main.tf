resource "junos_routing_instance" "testacc_generateRoute" {
  name = "testacc_generateRoute"
}
resource "junos_policyoptions_policy_statement" "testacc_generateRoute" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_generateRoute"
  then {
    action = "accept"
  }
}

resource "junos_generate_route" "testacc_generateRoute" {
  destination                  = "192.0.2.0/24"
  routing_instance             = junos_routing_instance.testacc_generateRoute.name
  preference                   = 100
  metric                       = 100
  active                       = true
  full                         = true
  discard                      = true
  community                    = ["no-advertise"]
  policy                       = [junos_policyoptions_policy_statement.testacc_generateRoute.name]
  as_path_aggregator_as_number = "65000"
  as_path_aggregator_address   = "192.0.2.1"
  as_path_atomic_aggregate     = true
  as_path_origin               = "igp"
  as_path_path                 = "65000 65000"
}
resource "junos_generate_route" "testacc_generateRoute6" {
  destination                  = "2001:db8:85a3::/48"
  routing_instance             = junos_routing_instance.testacc_generateRoute.name
  preference                   = 100
  metric                       = 100
  active                       = true
  full                         = true
  discard                      = true
  community                    = ["no-advertise"]
  policy                       = [junos_policyoptions_policy_statement.testacc_generateRoute.name]
  as_path_aggregator_as_number = "65000"
  as_path_aggregator_address   = "192.0.2.1"
  as_path_atomic_aggregate     = true
  as_path_origin               = "igp"
  as_path_path                 = "65000 65000"
}
