---
page_title: "Junos: junos_load_config"
---

# junos_load_config

Load an arbitrary configuration and commit it.

This action provides a way to load and commit a configuration in different `format`s to
the device without creating a persistent resource in the Terraform state.

<!-- markdownlint-disable -->
-> **Note**
  Actions are a Terraform 1.14+ feature that allow you to perform operations without managing state.
  Unlike the `junos_null_load_config` resource, this action does not create any state entries.

-> **Note**
  For details on how the configuration is committed for each value of `action`, see the NETCONF doc of the [&lt;load-configuration&gt;](https://www.juniper.net/documentation/us/en/software/junos/netconf/junos-xml-protocol/topics/ref/tag/junos-xml-protocol-load-configuration.html#load-configuration__d2375e254) command.
<!-- markdownlint-restore -->

## Example Usage

```hcl
action "junos_load_config" "applications" {
  config {
    action = "replace"
    config = <<EOT
replace: applications {
    application custom-ssh {
        protocol tcp;
        destination-port 22;
    }
}
EOT
  }
}

action "junos_load_config" "set_host-name" {
  config {
    action = "set"
    config = "set system host-name vSRX-1"
  }
}
```

## Argument Reference

The following arguments are supported:

- **config** (Required, String)  
  The configuration to load and apply.
- **action** (Optional, String)  
  Specify how to load the configuration data.  
  Must be `merge`, `override`, `replace`, `set` or `update`.  
  Defaults to `merge`.
- **format** (Optional, String)  
  The format used for the configuration data.  
  Must be `text`, `json` or `xml`.  
  Defaults to `text`.  

## Progress Events

This action sends progress updates during execution:

- Starting session to device
- Locking candidate configuration
- Loading configuration
- Committing configuration
- Configuration loaded and committed
