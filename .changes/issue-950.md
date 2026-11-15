<!-- markdownlint-disable-file MD013 MD041 -->
FEATURES:

* add `junos_policyoptions_damping` resource (Fix [#950](https://github.com/jeremmfr/terraform-provider-junos/issues/950))
* add `junos_policyoptions_damping` data source

ENHANCEMENTS:

* **resource/junos_policyoptions_policy_statement**: add `damping` argument in `then` block and in `then` block in `term` block
* **data-source/junos_policyoptions_policy_statement**: add `damping` attribute in `then` block and in `then` block in `term` block like resource
