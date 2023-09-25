resource "junos_routing_instance" "testacc_aggregateRoute" {
  name = "testacc_aggregateRoute"
}
resource "junos_policyoptions_policy_statement" "testacc_aggregateRoute" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_aggregateRoute"
  then {
    action = "accept"
  }
}

resource "junos_aggregate_route" "testacc_aggregateRoute" {
  destination                  = "192.0.2.0/24"
  routing_instance             = junos_routing_instance.testacc_aggregateRoute.name
  preference                   = 100
  metric                       = 100
  active                       = true
  full                         = true
  discard                      = true
  community                    = ["no-advertise"]
  policy                       = [junos_policyoptions_policy_statement.testacc_aggregateRoute.name]
  as_path_aggregator_as_number = "65000"
  as_path_aggregator_address   = "192.0.2.1"
  as_path_atomic_aggregate     = true
  as_path_origin               = "igp"
  as_path_path                 = "65000 65000"
}
resource "junos_aggregate_route" "testacc_aggregateRoute6" {
  destination                  = "2001:db8:85a3::/48"
  routing_instance             = junos_routing_instance.testacc_aggregateRoute.name
  preference                   = 100
  metric                       = 100
  active                       = true
  full                         = true
  discard                      = true
  community                    = ["no-advertise"]
  policy                       = [junos_policyoptions_policy_statement.testacc_aggregateRoute.name]
  as_path_aggregator_as_number = "65000"
  as_path_aggregator_address   = "192.0.2.1"
  as_path_atomic_aggregate     = true
  as_path_origin               = "igp"
  as_path_path                 = "65000 65000"
}
