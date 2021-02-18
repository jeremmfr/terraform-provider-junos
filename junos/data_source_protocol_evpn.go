package junos

import (
	"context"
	//"fmt"
	//"regexp"
	//"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProtocolEvpn() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProtocolEvpnRead,
		Schema: map[string]*schema.Schema{
                        "routing_instance": {
                                Type:     schema.TypeString,
                                Optional: true,
				Default:	defaultWord,
                        },
			"enabled":	{
				Type:	schema.TypeBool,
				Computed: true,
			},
			"encapsulation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"multicast_mode": {
				Type:     schema.TypeString,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"extended_vni_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
                        "route_distinguisher": {
                                Type:             schema.TypeString,
                                Required:         true,
                        },
                        "vrf_import": {
                                Type:     schema.TypeList,
                                Optional: true,
                                MinItems: 1,
                                Elem:     &schema.Schema{Type: schema.TypeString},
                        },
                        "vrf_export": {
                                Type:     schema.TypeList,
                                Optional: true,
                                MinItems: 1,
                                Elem:     &schema.Schema{Type: schema.TypeString},
                        },
                        "vrf_target": {
                                Type:             schema.TypeString,
                                Optional:         true,
                        },
                        "vtep_source_interface": {
                                Type:     schema.TypeString,
                                Optional: true,
                        },
		},
	}
}

func dataSourceProtocolEvpnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	var protocolEvpnEnabled bool
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	mutex.Lock()
	protocolEvpnOpts, err := readProtocolEvpn(d.Get("routing_instance").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("protocol_evpn/" + d.Get("routing_instance").(string))
	if tfErr := d.Set("enabled", protocolEvpnEnabled); tfErr != nil {
		panic(tfErr)
	}
	fillProtocolEvpnData(d, protocolEvpnOpts)

	return nil
}
