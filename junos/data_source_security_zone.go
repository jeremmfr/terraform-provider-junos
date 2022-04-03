package junos

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecurityZone() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceSecurityZoneRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"address_book": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"address_book_dns": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4_only": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ipv6_only": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"address_book_range": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"from": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"to": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"address_book_set": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"address_set": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"address_book_wildcard": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"advance_policy_based_routing_profile": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_tracking": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inbound_protocols": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"inbound_services": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"interface": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inbound_protocols": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"inbound_services": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"reverse_reroute": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"screen": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_identity_log": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tcp_rst": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSecurityZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	mutex.Lock()
	zoneOptions, err := readSecurityZone(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if zoneOptions.name == "" {
		return diag.FromErr(fmt.Errorf("security zone %v doesn't exist", d.Get("name").(string)))
	}
	d.SetId(zoneOptions.name)
	fillSecurityZoneDataSource(d, zoneOptions)

	return nil
}

func fillSecurityZoneDataSource(d *schema.ResourceData, zoneOptions zoneOptions) {
	if tfErr := d.Set("name", zoneOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_book", zoneOptions.addressBook); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_book_dns", zoneOptions.addressBookDNS); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_book_range", zoneOptions.addressBookRange); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_book_set", zoneOptions.addressBookSet); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_book_wildcard", zoneOptions.addressBookWildcard); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("advance_policy_based_routing_profile", zoneOptions.advancePolicyBasedRoutingProfile); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("application_tracking", zoneOptions.appTrack); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", zoneOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inbound_protocols", zoneOptions.inboundProtocols); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inbound_services", zoneOptions.inboundServices); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interface", zoneOptions.interFace); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("reverse_reroute", zoneOptions.reverseReroute); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("screen", zoneOptions.screen); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("source_identity_log", zoneOptions.sourceIdentityLog); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("tcp_rst", zoneOptions.tcpRst); tfErr != nil {
		panic(tfErr)
	}
}
