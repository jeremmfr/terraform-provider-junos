---
layout: "junos"
page_title: "Junos: junos_security"
sidebar_current: "docs-junos-resource-security"
description: |-
  Configure static configuration in security block (when Junos device supports it)
---

# junos_security

-> **Note:** This resource should only create **once**. It's used to configure static (not object) options in `security` block. Destroy this resource as no effect on Junos configuration.

Configure static configuration in `security` block

## Example Usage

```hcl
# Configure security
resource junos_security "security" {
  ike_traceoptions {
    file {
      name  = "ike.log"
      files = 5
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `ike_traceoptions` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'ike traceoptions' configuration.
  * `file` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'file' configuration. See the [`file` argument] (#file-argument) block.
  * `flag` - (Optional)(`ListOfString`) Tracing parameters for IKE.
  * `rate_limit` - (Optional)(`Int`) Limit the incoming rate of trace messages (0..4294967295)
* `utm` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'utm' configuration.
  * `feature_profile_web_filtering_type` - (Optional)(`String`) Configuring feature-profile web-filtering type. Need to be 'juniper-enhanced', 'juniper-local', 'web-filtering-none' or 'websense-redirect'.
* `alg` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'alg' configuration. See the [`alg` argument] (#alg-argument) block.

#### file argument
* `name` - (Optional)(`String`) Name of file in which to write trace information.
* `files` - (Optional)(`Int`) Maximum number of trace files (2..1000).
* `match` - (Optional)(`String`) Regular expression for lines to be logged.
* `no_world_readable` - (Optional)(`Bool`) Don't allow any user to read the log file.
* `size` - (Optional)(`Int`) Maximum trace file size (10240..1073741824)
* `world_readable` - (Optional)(`Bool`) Allow any user to read the log file

#### alg argument
* `dns_disable` - (Optional)(`Bool`) Disable dns alg.
* `ftp_disable` - (Optional)(`Bool`) Disable ftp alg.
* `msrpc_disable` - (Optional)(`Bool`) Disable msrpc alg.
* `pptp_disable` - (Optional)(`Bool`) Disable pptp alg.
* `sunrpc_disable` - (Optional)(`Bool`) Disable sunrpc alg.
* `talk_disable` - (Optional)(`Bool`) Disable talk alg.
* `tftp_disable` - (Optional)(`Bool`) Disable tftp alg.

## Import

Junos security can be imported using any id, e.g.

```
$ terraform import junos_security.security random
```
