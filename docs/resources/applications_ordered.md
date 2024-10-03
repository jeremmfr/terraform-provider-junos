---
page_title: "Junos: junos_applications_ordered"
---

# junos_applications_ordered

It has the same functionality as the `junos_applications` resource
but with `applications` and `application_set` arguments as Block List instead of Block Set.

This provides a workaround for the performance issue on Terraform plan with many Block Sets
(details in GitHub issue [#775](https://github.com/hashicorp/terraform-plugin-framework/issues/775))
but Block List involves:

- a change in the order of the blocks triggers a resource change.
- Terraform plan output can be complex when the number of blocks on the resource changes.

See the [junos_applications](applications) resource
for more details on arguments or attributes.
