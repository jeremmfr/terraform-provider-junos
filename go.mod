module github.com/jeremmfr/terraform-provider-junos

go 1.16

require (
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.11.0
	github.com/jeremmfr/go-netconf v0.4.3
	github.com/jeremmfr/go-utils v0.4.1
	github.com/jeremmfr/junosdecode v1.1.0
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd
)

replace github.com/hashicorp/terraform-plugin-sdk/v2 => github.com/jeremmfr/terraform-plugin-sdk/v2 v2.11.1-0.20220316083425-f711567b3c5d
