terraform-provider-junos
========================
![GitHub release (latest by date)](https://img.shields.io/github/v/release/jeremmfr/terraform-provider-junos)
[![Go Status](https://github.com/jeremmfr/terraform-provider-junos/workflows/Go%20Tests/badge.svg)](https://github.com/jeremmfr/terraform-provider-junos/actions)
[![Lint Status](https://github.com/jeremmfr/terraform-provider-junos/workflows/GolangCI-Lint/badge.svg)](https://github.com/jeremmfr/terraform-provider-junos/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/jeremmfr/terraform-provider-junos)](https://goreportcard.com/report/github.com/jeremmfr/terraform-provider-junos)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/jeremmfr/terraform-provider-junos/blob/master/LICENSE)
<br/><br/>
This is an **unofficial** terraform provider for Junos devices with netconf protocol

Requirements
------------
-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

Building The Provider
---------------------
```
$ git clone https://github.com/jeremmfr/terraform-provider-junos.git
$ cd terraform-provider-junos && git fetch --tags
$ latestTag=$(git describe --tags `git rev-list --tags --max-count=1`)
$ git checkout ${latestTag}
$ tfPath=$(which terraform | rev | cut -d'/' -f2- | rev)
$ go build -o ${tfPath}/terraform-provider-junos_${latestTag}
$ unset latestTag tfPath
```

See [website/docs/index](website/docs/index.html.markdown) for config provider and start add resource

See [website/docs/r/](https://github.com/jeremmfr/terraform-provider-junos/tree/master/website/docs/r) directory for resources documentation.

Some Junos parameters are not included in provider for various reasons (time, utility, understanding, ...)
