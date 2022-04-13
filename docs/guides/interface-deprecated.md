---
page_title: "Junos: junos_interface resource deprecated"
---

# junos_interface deprecated

For more consistency, functionalities of `junos_interface` resource have been splitted in two new
resource : `junos_interface_physical` and `junos_interface_logical`.  
The `junos_interface` resource is **deprecated** since v1.11.0.

## Rewrite resource for physical interface

For physical interface (without dot in name) :

- rename the type of resource `junos_interface` to `junos_interface_physical`
- rename `complete_destroy` argument to `no_disable_on_destroy`

For example :

```hcl
# deprecated junos_interface
resource "junos_interface" "interface_physical_demo" {
  name         = "ge-0/0/0"
  description  = "interfacePhysicalDemo"
  trunk        = true
  vlan_members = ["100"]
}

# new junos_interface_physical
resource "junos_interface_physical" "interface_physical_demo" {
  name         = "ge-0/0/0"
  description  = "interfacePhysicalDemo"
  trunk        = true
  vlan_members = ["100"]
}
```

## Rewrite resource for logical interface

For logical interface (with dot in name) :

- rename type of resource `junos_interface`to `junos_interface_logical`
- rename `vlan_tagging_id` argument to `vlan_id`
- rename `complete_destroy` argument to `st0_also_on_destroy`
- move `inet_*` arguments in new `family_inet` block without prefix `inet_`
- move `inet6_*` arguments in new `family_inet6` block without prefix `inet6_`
- rename `address` in old `inet_address` argument to `cidr_ip`
- rename `address` in old `inet6_address` argument to `cidr_ip`

For example :

```hcl
# deprecated junos_interface 
resource "junos_interface" "interface_logical_demo_100" {
  name        = "ge-0/0/2.100"
  description = "interfaceLogicalDemo100"
  inet_address {
    address = "192.0.2.1/25"
  }
  inet_filter_input = "filter_demo"
}
resource "junos_interface" "st0_100" {
  name        = "st0.100"
  description = "st0_100"
  inet        = true
}

# new junos_interface_logical
resource "junos_interface_logical" "interface_logical_demo_100" {
  name        = "ge-0/0/2.100"
  description = "interfaceLogicalDemo100"
  family_inet {
    address {
      cidr_ip = "192.0.2.1/25"
    }
    filter_input = "filter_demo"
  }
}
resource "junos_interface_logical" "st0_100" {
  name        = "st0.100"
  description = "st0_100"
  family_inet {}
}
```

## Upgrade without destroy and create new

For upgrade to the new resource without destroy deprecated resource and recreate resource, you need
to import the new resource and delete the old resource in Terraform state.

After rewrite resource with new type, import each resources :

```shell
terraform import junos_interface_physical.interface_physical_demo ge-0/0/0
terraform import junos_interface_logical.interface_logical_demo_100 ge-0/0/2.100
terraform import junos_interface_logical.st0_100 st0.100
```

then now delete deprecated resource :

```shell
terraform state rm junos_interface.interface_physical_demo
terraform state rm junos_interface.interface_logical_demo_100
terraform state rm junos_interface.st0_100
```
