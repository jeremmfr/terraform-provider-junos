## 1.0.4

BUG FIXES:

* fix ipsec_vpn bind_interface_auto -> search st0 unit not in terse simply ([17](https://github.com/jeremmfr/terraform-provider-junos/pull/17))
* remove commit-check before commit which gives the same error if there is ([16](https://github.com/jeremmfr/terraform-provider-junos/pull/16))
* fix check interface disable and NC ([15](https://github.com/jeremmfr/terraform-provider-junos/pull/15))

## 1.0.3

BUG FIXES:

* fix terraform crash with an empty blocks-mode (no one required) ([14](https://github.com/jeremmfr/terraform-provider-junos/pull/14))

## 1.0.2

ENHANCEMENTS:

* move cmd/debug environnement variables to provider config ([13](https://github.com/jeremmfr/terraform-provider-junos/pull/13))

## 1.0.1

BUG FIXES:

* fix readInterface with empty/disappeared interface

## 1.0.0

First release
