module github.com/jeremmfr/terraform-provider-junos

go 1.16

require (
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.8.0
	github.com/jeremmfr/go-netconf v0.4.0
	github.com/jeremmfr/go-utils v0.3.0
	github.com/jeremmfr/junosdecode v1.1.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
)

replace github.com/hashicorp/terraform-plugin-sdk/v2 v2.8.0 => github.com/jeremmfr/terraform-plugin-sdk/v2 v2.8.1-0.20211007115003-2ac7d96a040a
