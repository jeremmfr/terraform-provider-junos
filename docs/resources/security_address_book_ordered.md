---
page_title: "Junos: junos_security_address_book_ordered"
---

# junos_security_address_book_ordered

It has the same functionality as the `junos_security_address_book` resource
but with `network_address`, `dns_name`, `range_address`, `wildcard_address` and `address_set`
arguments as Block List instead of Block Set.

This provides a workaround for the performance issue on Terraform plan with many Block Sets
(details in GitHub issue [#775](https://github.com/hashicorp/terraform-plugin-framework/issues/775))
but Block List involves:

- a change in the order of the blocks triggers a resource change.
- Terraform plan output can be complex when the number of blocks on the resource changes.

See the [junos_security_address_book](security_address_book) resource
for more details on arguments or attributes.
