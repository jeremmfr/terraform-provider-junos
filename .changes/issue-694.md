<!-- markdownlint-disable-file MD013 MD041 -->
FEATURES:

* add **junos_applications** resource (Fix [#694](https://github.com/jeremmfr/terraform-provider-junos/issues/694))

ENHANCEMENTS:

* **resource/junos_application**: add `do_not_translate_a_query_to_aaaa_query`, `do_not_translate_aaaa_query_to_a_query`, `icmp_code`, `icmp_type`, `icmp6_code` and `icmp6_type` arguments
* **data-source/junos_applications**:
  * add `do_not_translate_a_query_to_aaaa_query`, `do_not_translate_aaaa_query_to_a_query`, `icmp_code`, `icmp_type`, `icmp6_code` and `icmp6_type` attributes in `applications` attribute
  * add `do_not_translate_a_query_to_aaaa_query` and `do_not_translate_aaaa_query_to_a_query` arguments inside `match_options` block argument
