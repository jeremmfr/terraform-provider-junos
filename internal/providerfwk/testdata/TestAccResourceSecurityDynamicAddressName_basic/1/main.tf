resource "junos_security_dynamic_address_feed_server" "testacc_dyn_add_name" {
  lifecycle {
    create_before_destroy = true
  }

  name     = "tfacc_dynadd"
  hostname = "example.com"
  feed_name {
    name = "feedtfacc"
    path = "/srx/"
  }
}
resource "junos_security_dynamic_address_name" "testacc_dyn_add_name" {
  name              = "tfacc_dynadd"
  description       = "desc tfacc dynamic-address"
  profile_feed_name = junos_security_dynamic_address_feed_server.testacc_dyn_add_name.feed_name.0.name
}
resource "junos_security_dynamic_address_name" "testacc_dyn_add_name2" {
  name        = "tfacc_dynadd2"
  description = "desc tfacc dynamic-address2"
  profile_category {
    name = "IPFilter"
    property {
      name   = "others #1"
      string = ["test#2", "test#1"]
    }
    property {
      name   = "country"
      string = ["AU", "CN"]
    }
  }
}
