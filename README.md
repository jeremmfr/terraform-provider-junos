<!-- markdownlint-disable-file MD013 MD033 MD041 -->
<div>
<a href="https://www.juniper.net"><img src=".github/junos-os.png" alt="Junos logo" title="Junos" align="right" height="50" /></a>
<a href="https://www.terraform.io"><img src=".github/terraform.png" alt="Terraform logo" title="Terraform" align="right" height="50" /></a>
</div>

# terraform-provider-junos

<!-- markdownlint-disable -->
[![Release](https://img.shields.io/github/v/release/jeremmfr/terraform-provider-junos)](https://github.com/jeremmfr/terraform-provider-junos/releases)
[![Installs](https://img.shields.io/badge/dynamic/json?logo=terraform&label=installs&query=$.data.attributes.total&url=https%3A%2F%2Fregistry.terraform.io%2Fv2%2Fproviders%2F713%2Fdownloads%2Fsummary)](https://registry.terraform.io/providers/jeremmfr/junos)
[![Registry](https://img.shields.io/badge/registry-doc%40latest-lightgrey?logo=terraform)](https://registry.terraform.io/providers/jeremmfr/junos/latest/docs)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/jeremmfr/terraform-provider-junos/blob/main/LICENSE)  
[![Go Status](https://github.com/jeremmfr/terraform-provider-junos/actions/workflows/go.yml/badge.svg)](https://github.com/jeremmfr/terraform-provider-junos/actions/workflows/go.yml)
[![Linters Status](https://github.com/jeremmfr/terraform-provider-junos/actions/workflows/linters.yml/badge.svg)](https://github.com/jeremmfr/terraform-provider-junos/actions/workflows/linters.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jeremmfr/terraform-provider-junos)](https://goreportcard.com/report/github.com/jeremmfr/terraform-provider-junos)  
[![Buy Me A Coffee](https://img.shields.io/badge/buy%20me%20a%20coffee-donate-yellow.svg)](https://www.buymeacoffee.com/jeremmfr)
<!-- markdownlint-restore -->

---

This is an **unofficial** Terraform provider for Junos devices with netconf protocol

See [Terraform registry](https://registry.terraform.io/providers/jeremmfr/junos)
for provider and resources documentation.

## Requirements

- [Terraform](https://www.terraform.io/downloads)

### In addition to develop

- [Go](https://golang.org/doc/install) `v1.21` or `v1.22`

## Automatic install

Add source information inside the Terraform configuration block for automatic provider installation:

```hcl
terraform {
  required_providers {
    junos = {
      source = "jeremmfr/junos"
    }
  }
}
```

## Manual install

- Download latest version in [releases](https://github.com/jeremmfr/terraform-provider-junos/releases)

- Extract provider binary in
[local mirror directory](https://www.terraform.io/cli/config/config-file#implied-local-mirror-directories)
with a fake registry (`registry.local`):

```bash
for archive in $(ls terraform-provider-junos*.zip) ; do
  OS_ARCH=$(echo $archive | cut -d'_' -f3-4 | cut -d'.' -f1)
  VERSION=$(echo $archive | cut -d'_' -f2)
  tfPath="${HOME}/.terraform.d/plugins/registry.local/jeremmfr/junos/${VERSION}/${OS_ARCH}/"
  mkdir -p ${tfPath}
  unzip ${archive} -d ${tfPath}
done
```

- Add inside the terraform configuration block:

```hcl
terraform {
  required_providers {
    junos = {
      source = "registry.local/jeremmfr/junos"
    }
  }
}
```

## Missing Junos parameters

Some Junos parameters are not included in provider for various reasons
(time, utility, understanding, ...) but you can create a issue
to request the potential addition of missing features.

## Contributing

To contribute, please read the [contribution guideline](.github/CONTRIBUTING.md)

## Compile a binary from source to use with Terraform

### Build to override automatic install version (Terraform 0.14 and later)

Since Terraform 0.14,
[development overrides for provider developers](https://www.terraform.io/cli/config/config-file#development-overrides-for-provider-developers)
allow to use the provider built from source.  
Use a Terraform [cli configuration file](https://www.terraform.io/cli/config/config-file)
(`~/.terraformrc` by default) with at least the following options:

```hcl
provider_installation {
  dev_overrides {
    "jeremmfr/junos" = "[replace with the GOPATH]/bin"
  }
  direct {}
}
```

and build then install in $GOPATH/bin:

```bash
git clone https://github.com/jeremmfr/terraform-provider-junos.git
cd terraform-provider-junos
go install
```

---

### Build to use with a fake registry (Terraform 0.13)

```bash
git clone https://github.com/jeremmfr/terraform-provider-junos.git
cd terraform-provider-junos && git fetch --tags
latestTag=$(git describe --tags `git rev-list --tags --max-count=1`)
git checkout ${latestTag}
tfPath="${HOME}/.terraform.d/plugins/registry.local/jeremmfr/junos/${latestTag:1}/$(go env GOOS)_$(go env GOARCH)/"
mkdir -p ${tfPath}
go build -o ${tfPath}/terraform-provider-junos_${latestTag}
unset latestTag tfPath
```

and add inside the terraform configuration block:

```hcl
terraform {
  required_providers {
    junos = {
      source = "registry.local/jeremmfr/junos"
    }
  }
}
```
