resource "junos_security_dynamic_address_name" "testacc_dyn_add_name" {
  name        = "tfacc_dynadd"
  description = "desc tfacc dynamic-address"
  profile_category {
    name = "GeoIP"
    feed = "cat_feed"
    property {
      name   = "country"
      string = ["AU", "CN"]
    }
  }
}
