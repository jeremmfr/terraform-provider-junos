terraform {
  required_providers {
    junos = {
      source  = "registry.terraform.io/jeremmfr/junos"
      version = "1.33.0"
    }
  }
}

provider "junos" {}
