<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_system_login_user**, **resource/junos_system_root_authentication**: mark `encrypted_password` attribute as sensitive (to obscure the value in Terraform output) even if the data is encrypted
