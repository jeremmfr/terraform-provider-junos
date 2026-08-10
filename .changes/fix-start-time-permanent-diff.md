<!-- markdownlint-disable-file MD013 MD041 -->
BUG FIXES:

* **resource/junos_eventoptions_generate_event**: fix a permanent diff when the `start_time` argument has a leading zero in the month or the day, the device removes it so the value read didn't match the configuration
