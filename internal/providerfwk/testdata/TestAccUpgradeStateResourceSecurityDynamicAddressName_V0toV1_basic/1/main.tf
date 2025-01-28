resource "junos_security_dynamic_address_name" "testacc_dyn_add_name2" {
  name        = "tfacc_dynadd2"
  description = "desc tfacc dynamic-address2"
  profile_category {
    name = "IPFilter"
    property {
      name   = "others#1"
      string = ["test#2", "test#1"]
    }
    property {
      name   = "country"
      string = ["AU", "CN"]
    }
  }
}
