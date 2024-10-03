---
page_title: "Junos: junos_security_policy_unordered"
---

# junos_security_policy_unordered

It has the same functionality as the `junos_security_policy` resource
but with `policy` argument as Block Set instead of Block List.

This provides a workaround for too complex plan output when the number of blocks on the resource changes
and if the `policy` order it's not important
(by considering the action of the first policy that the traffic matches is applied to the packet).

Block Set involves a performance issue on Terraform plan with many Block Sets
(details in GitHub issue [#775](https://github.com/hashicorp/terraform-plugin-framework/issues/775)).

See the [junos_security_policy](security_policy) resource
for more details on arguments or attributes.
