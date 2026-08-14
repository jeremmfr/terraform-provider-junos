<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_chassis_fpc**: add `configure_sampling_instance_singly` argument to configure the `sampling-instance` option in an other resource
* **resource/junos_forwardingoptions_sampling_instance**: add `chassis_fpc_slot_numbers` argument to attach the sampling instance to chassis FPC slots with the same commit, required on some Junos systems where the sampling instance and the FPC binding cannot be committed separately (Fix [#912](https://github.com/jeremmfr/terraform-provider-junos/issues/912))
