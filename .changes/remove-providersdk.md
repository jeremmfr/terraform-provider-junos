<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_null_commit_file**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  `triggers` argument now accept any attribute type (and so still a Map)  
  the permissions of file targeted by `filename` argument are now preserved when using `clear_file_after_commit` argument
