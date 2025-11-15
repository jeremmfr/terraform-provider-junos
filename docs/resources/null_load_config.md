---
page_title: "Junos: junos_null_load_config"
---

# junos_null_load_config

Load an arbitrary configuration and commit it.

<!-- markdownlint-disable -->
-> **Note**
  Does not provide a real resource, just loads a configuration in different `format`s to the
  candidate configuration on the device, and commits it with different types of `action`.  
  For details on how the configuration is committed for each value of `action`, see the NETCONF doc of the [&lt;load-configuration&gt;](https://www.juniper.net/documentation/us/en/software/junos/netconf/junos-xml-protocol/topics/ref/tag/junos-xml-protocol-load-configuration.html#load-configuration__d2375e254) command.
<!-- markdownlint-restore -->

-> **Note**
  Read resource doesn't update the state of resource
  (there is no comparison with the device configuration).
  So if the configuration is **not correct**, the resource cannot detect it.  
  Destroying this resource has no effect on the Junos configuration.

## Example Usage

```hcl
resource "junos_null_load_config" "applications" {
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

resource "junos_null_load_config" "set_host-name" {
  action = "set"
  config = "set system host-name vSRX-1"
}
```

## Argument Reference

The following arguments are supported:

- **config** (Required, String, Forces new resource)  
  The configuration to load and apply.
- **action** (Optional, String, Forces new resource)  
  Specify how to load the configuration data.  
  Need to be `merge`, `override`, `replace`, `set` or `update`.  
  Defaults to `merge`.
- **format** (Optional, String, Forces new resource)  
  The format used for the configuration data.  
  Need to be `text`, `json` or `xml`.  
  Defaults to `text`.
- **triggers** (Optional, Any, Forces new resource)  
  Any value that, when changed, will force the resource to be replaced.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `null_load_config`.
