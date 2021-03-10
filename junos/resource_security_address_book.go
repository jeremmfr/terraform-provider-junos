package junos

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type addressBookOptions struct {
	name        string
	description string
	address     []map[string]interface{}
	addressSet  []map[string]interface{}
	attachZone  string
}

func resourceSecurityAddressBook() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceSecurityAddressBookRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true, // Not sure if this is right?
				Default:  "global",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dns_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ExactlyOneOf: []string{"network", "dns_name", "range_address", "wildcard_address"},
						},
						"network": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
							ExactlyOneOf: []string{"network", "dns_name", "range_address", "wildcard_address"},
						},
						"range_address": {
							Type:         schema.TypeSet,
							Optional:     true,
							MaxItems:     1,
							ExactlyOneOf: []string{"network", "dns_name", "range_address", "wildcard_address"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"from": {
										Type:         schema.TypeString,
										Optional:     false,
										ValidateFunc: validation.IsIPAddress,
									},
									"to": {
										Type:         schema.TypeString,
										Optional:     false,
										ValidateFunc: validation.IsIPAddress,
									},
								},
							},
						},
						"wildcard_address": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateWildcardFunc(),
							ExactlyOneOf:     []string{"network", "dns_name", "range_address", "wildcard_address"},
						},
					},
				},
			},
			"address_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
						},
						"address": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"attach_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
