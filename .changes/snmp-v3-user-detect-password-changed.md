<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_snmp_v3_usm_user**: now provider store the corresponding `authentication-key` and `privacy-key` of `authentication-password` and `privacy-password` in private state of Terraform after create/update resource to be able to detect a change of the password outside of Terraform.
