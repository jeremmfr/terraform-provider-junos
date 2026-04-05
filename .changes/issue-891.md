<!-- markdownlint-disable-file MD013 MD041 -->
FEATURES:

* provider: add `single_session` argument to use a single shared SSH/NETCONF session for all provider operations instead of opening a new session for each resource action. (Fix [#891](https://github.com/jeremmfr/terraform-provider-junos/issues/891))  
  **Be careful with this option, because Terraform shuts down the provider without prior notice at the end of a run, the session will not be properly closed.**  
  It can also be enabled from the `JUNOS_SINGLE_SESSION` environment variable.
