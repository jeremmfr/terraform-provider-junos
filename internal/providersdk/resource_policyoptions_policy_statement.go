package providersdk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type policyStatementOptions struct {
	name string
	from []map[string]interface{}
	then []map[string]interface{}
	to   []map[string]interface{}
	term []map[string]interface{}
}

func resourcePolicyoptionsPolicyStatement() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourcePolicyoptionsPolicyStatementCreate,
		ReadWithoutTimeout:   resourcePolicyoptionsPolicyStatementRead,
		UpdateWithoutTimeout: resourcePolicyoptionsPolicyStatementUpdate,
		DeleteWithoutTimeout: resourcePolicyoptionsPolicyStatementDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyoptionsPolicyStatementImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"add_it_to_forwarding_table_export": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"from": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: schemaPolicyoptionsPolicyStatementFrom(),
				},
			},
			"then": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: schemaPolicyoptionsPolicyStatementThen(),
				},
			},
			"to": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: schemaPolicyoptionsPolicyStatementTo(),
				},
			},
			"term": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"from": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: schemaPolicyoptionsPolicyStatementFrom(),
							},
						},
						"then": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: schemaPolicyoptionsPolicyStatementThen(),
							},
						},
						"to": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: schemaPolicyoptionsPolicyStatementTo(),
							},
						},
					},
				},
			},
		},
	}
}

func schemaPolicyoptionsPolicyStatementFrom() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"aggregate_contributor": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"bgp_as_path": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64, formatDefault),
			},
		},
		"bgp_as_path_calc_length": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"count": {
						Type:         schema.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 1024),
					},
					"match": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "orhigher", "orlower"}, false),
					},
				},
			},
		},
		"bgp_as_path_group": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64, formatDefault),
			},
		},
		"bgp_as_path_unique_count": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"count": {
						Type:         schema.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 1024),
					},
					"match": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "orhigher", "orlower"}, false),
					},
				},
			},
		},
		"bgp_community": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64, formatDefault),
			},
		},
		"bgp_community_count": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"count": {
						Type:         schema.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 1024),
					},
					"match": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "orhigher", "orlower"}, false),
					},
				},
			},
		},
		"bgp_origin": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"egp", "igp", "incomplete"}, false),
		},
		"bgp_srte_discriminator": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      -1,
			ValidateFunc: validation.IntBetween(0, 4294967295),
		},
		"color": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      -1,
			ValidateFunc: validation.IntBetween(0, 4294967295),
		},
		"evpn_esi": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^([\d\w]{2}:){9}[\d\w]{2}$`), "bad format or length"),
			},
		},
		"evpn_mac_route": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"mac-ipv4", "mac-ipv6", "mac-only"}, false),
		},
		"evpn_tag": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 4294967295),
			},
		},
		"family": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"evpn", "inet", "inet-mdt", "inet-mvpn", "inet-vpn",
				"inet6", "inet6-mvpn", "inet6-vpn", "iso",
			}, false),
		},
		"local_preference": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"routing_instance": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64, formatDefault),
		},
		"interface": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"metric": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"neighbor": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"next_hop": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"next_hop_type_merged": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"next_hop_weight": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"match": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"equal", "greater-than", "greater-than-equal", "less-than", "less-than-equal",
						}, false),
					},
					"weight": {
						Type:         schema.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 65535),
					},
				},
			},
		},
		"ospf_area": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"policy": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
		},
		"preference": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"prefix_list": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
		},
		"protocol": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"route_filter": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"route": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.IsCIDRNetwork(0, 128),
					},
					"option": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"address-mask", "exact", "longer", "orlonger", "prefix-length-range", "through", "upto",
						}, false),
					},
					"option_value": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"route_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"external", "internal"}, false),
		},
		"srte_color": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      -1,
			ValidateFunc: validation.IntBetween(0, 4294967295),
		},
		"state": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"active", "inactive"}, false),
		},
		"tunnel_type": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"gre", "ipip", "udp"}, false),
			},
		},
		"validation_database": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"invalid", "unknown", "valid"}, false),
		},
	}
}

func schemaPolicyoptionsPolicyStatementThen() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"accept", "reject"}, false),
		},
		"as_path_expand": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"as_path_prepend": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"community": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"action": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"add", "delete", "set"}, false),
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"default_action": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"accept", "reject"}, false),
		},
		"load_balance": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"per-packet", "consistent-hash"}, false),
		},
		"local_preference": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"action": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"add", "subtract", "none"}, false),
					},
					"value": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
		"next": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"policy", "term"}, false),
		},
		"next_hop": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"metric": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"action": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"add", "subtract", "none"}, false),
					},
					"value": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
		"origin": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"preference": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"action": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"add", "subtract", "none"}, false),
					},
					"value": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
	}
}

func schemaPolicyoptionsPolicyStatementTo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bgp_as_path": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64, formatDefault),
			},
		},
		"bgp_as_path_group": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64, formatDefault),
			},
		},
		"bgp_community": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64, formatDefault),
			},
		},
		"bgp_origin": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"egp", "igp", "incomplete"}, false),
		},
		"family": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"evpn", "inet", "inet-mdt", "inet-mvpn", "inet-vpn",
				"inet6", "inet6-mvpn", "inet6-vpn", "iso",
			}, false),
		},
		"local_preference": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"routing_instance": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64, formatDefault),
		},
		"interface": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"metric": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"neighbor": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"next_hop": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"ospf_area": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"policy": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
		},
		"preference": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"protocol": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func resourcePolicyoptionsPolicyStatementCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setPolicyStatement(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if d.Get("add_it_to_forwarding_table_export").(bool) {
			if err := setPolicyStatementFwTableExport(d.Get("name").(string), junSess); err != nil {
				return diag.FromErr(err)
			}
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	policyStatementExists, err := checkPolicyStatementExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyStatementExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("policy-options policy-statement %v already exists", d.Get("name").(string)))...)
	}

	if err := setPolicyStatement(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if d.Get("add_it_to_forwarding_table_export").(bool) {
		if err := setPolicyStatementFwTableExport(d.Get("name").(string), junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	warns, err := junSess.CommitConf("create resource junos_policyoptions_policy_statement")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	policyStatementExists, err = checkPolicyStatementExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyStatementExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("policy-options policy-statement %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourcePolicyoptionsPolicyStatementReadWJunSess(d, junSess)...)
}

func resourcePolicyoptionsPolicyStatementRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourcePolicyoptionsPolicyStatementReadWJunSess(d, junSess)
}

func resourcePolicyoptionsPolicyStatementReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	policyStatementOptions, err := readPolicyStatement(d.Get("name").(string), junSess)
	if err != nil {
		junos.MutexUnlock()

		return diag.FromErr(err)
	}
	if d.Get("add_it_to_forwarding_table_export").(bool) {
		export, err := readPolicyStatementFwTableExport(d.Get("name").(string), junSess)
		if err != nil {
			junos.MutexUnlock()

			return diag.FromErr(err)
		}
		if !export {
			if tfErr := d.Set("add_it_to_forwarding_table_export", false); tfErr != nil {
				panic(tfErr)
			}
		}
	}
	junos.MutexUnlock()

	if policyStatementOptions.name == "" {
		d.SetId("")
	} else {
		fillPolicyStatementData(d, policyStatementOptions)
	}

	return nil
}

func resourcePolicyoptionsPolicyStatementUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delPolicyStatement(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if d.HasChange("add_it_to_forwarding_table_export") {
			if o, _ := d.GetChange("add_it_to_forwarding_table_export"); o.(bool) {
				if err := delPolicyStatementFwTableExport(d.Get("name").(string), junSess); err != nil {
					return diag.FromErr(err)
				}
			}
		}
		if err := setPolicyStatement(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if d.Get("add_it_to_forwarding_table_export").(bool) {
			if err := setPolicyStatementFwTableExport(d.Get("name").(string), junSess); err != nil {
				return diag.FromErr(err)
			}
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delPolicyStatement(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if d.HasChange("add_it_to_forwarding_table_export") {
		if o, _ := d.GetChange("add_it_to_forwarding_table_export"); o.(bool) {
			if err := delPolicyStatementFwTableExport(d.Get("name").(string), junSess); err != nil {
				appendDiagWarns(&diagWarns, junSess.ConfigClear())

				return append(diagWarns, diag.FromErr(err)...)
			}
		}
	}
	if err := setPolicyStatement(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if d.Get("add_it_to_forwarding_table_export").(bool) {
		if err := setPolicyStatementFwTableExport(d.Get("name").(string), junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	warns, err := junSess.CommitConf("update resource junos_policyoptions_policy_statement")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourcePolicyoptionsPolicyStatementReadWJunSess(d, junSess)...)
}

func resourcePolicyoptionsPolicyStatementDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delPolicyStatement(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if d.Get("add_it_to_forwarding_table_export").(bool) {
			if err := delPolicyStatementFwTableExport(d.Get("name").(string), junSess); err != nil {
				return diag.FromErr(err)
			}
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delPolicyStatement(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if d.Get("add_it_to_forwarding_table_export").(bool) {
		if err := delPolicyStatementFwTableExport(d.Get("name").(string), junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	warns, err := junSess.CommitConf("delete resource junos_policyoptions_policy_statement")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourcePolicyoptionsPolicyStatementImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	policyStatementExists, err := checkPolicyStatementExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !policyStatementExists {
		return nil, fmt.Errorf("don't find policy-options policy-statement with id '%v' (id must be <name>)", d.Id())
	}
	policyStatementOptions, err := readPolicyStatement(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillPolicyStatementData(d, policyStatementOptions)

	result[0] = d

	return result, nil
}

func checkPolicyStatementExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"policy-options policy-statement " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setPolicyStatement(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set policy-options policy-statement " + d.Get("name").(string)
	for _, from := range d.Get("from").([]interface{}) {
		if from != nil {
			configSetFrom, err := setPolicyStatementOptsFrom(setPrefix, from.(map[string]interface{}))
			if err != nil {
				return err
			}
			configSet = append(configSet, configSetFrom...)
		}
	}
	for _, then := range d.Get("then").([]interface{}) {
		if then != nil {
			configSetThen, err := setPolicyStatementOptsThen(setPrefix, then.(map[string]interface{}))
			if err != nil {
				return err
			}
			configSet = append(configSet, configSetThen...)
		}
	}
	for _, to := range d.Get("to").([]interface{}) {
		if to != nil {
			configSetTo := setPolicyStatementOptsTo(setPrefix, to.(map[string]interface{}))
			configSet = append(configSet, configSetTo...)
		}
	}
	termNameList := make([]string, 0)
	for _, term := range d.Get("term").([]interface{}) {
		termMap := term.(map[string]interface{})
		if bchk.InSlice(termMap["name"].(string), termNameList) {
			return fmt.Errorf("multiple blocks term with the same name %s", termMap["name"].(string))
		}
		termNameList = append(termNameList, termMap["name"].(string))
		setPrefixTerm := setPrefix + " term " + termMap["name"].(string)
		for _, from := range termMap["from"].([]interface{}) {
			if from != nil {
				configSetFrom, err := setPolicyStatementOptsFrom(setPrefixTerm, from.(map[string]interface{}))
				if err != nil {
					return err
				}
				configSet = append(configSet, configSetFrom...)
			}
		}
		for _, then := range termMap["then"].([]interface{}) {
			if then != nil {
				configSetThen, err := setPolicyStatementOptsThen(setPrefixTerm, then.(map[string]interface{}))
				if err != nil {
					return err
				}
				configSet = append(configSet, configSetThen...)
			}
		}
		for _, to := range termMap["to"].([]interface{}) {
			if to != nil {
				configSetTo := setPolicyStatementOptsTo(setPrefixTerm, to.(map[string]interface{}))
				configSet = append(configSet, configSetTo...)
			}
		}
	}

	return junSess.ConfigSet(configSet)
}

func setPolicyStatementFwTableExport(policyName string, junSess *junos.Session) error {
	configSet := []string{"set routing-options forwarding-table export " + policyName}

	return junSess.ConfigSet(configSet)
}

func readPolicyStatement(name string, junSess *junos.Session,
) (confRead policyStatementOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"policy-options policy-statement " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "term "):
				itemTrimFields := strings.Split(itemTrim, " ")
				termOptions := map[string]interface{}{
					"name": itemTrimFields[0],
					"from": make([]map[string]interface{}, 0),
					"then": make([]map[string]interface{}, 0),
					"to":   make([]map[string]interface{}, 0),
				}
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				confRead.term = copyAndRemoveItemMapList("name", termOptions, confRead.term)
				switch {
				case balt.CutPrefixInString(&itemTrim, "from "):
					if len(termOptions["from"].([]map[string]interface{})) == 0 {
						termOptions["from"] = append(termOptions["from"].([]map[string]interface{}),
							genMapPolicyStatementOptsFrom())
					}
					if err := readPolicyStatementOptsFrom(itemTrim,
						termOptions["from"].([]map[string]interface{})[0]); err != nil {
						return confRead, err
					}
				case balt.CutPrefixInString(&itemTrim, "then "):
					if len(termOptions["then"].([]map[string]interface{})) == 0 {
						termOptions["then"] = append(termOptions["then"].([]map[string]interface{}),
							genMapPolicyStatementOptsThen())
					}
					if err := readPolicyStatementOptsThen(itemTrim,
						termOptions["then"].([]map[string]interface{})[0]); err != nil {
						return confRead, err
					}
				case balt.CutPrefixInString(&itemTrim, "to "):
					if len(termOptions["to"].([]map[string]interface{})) == 0 {
						termOptions["to"] = append(termOptions["to"].([]map[string]interface{}),
							genMapPolicyStatementOptsTo())
					}
					if err := readPolicyStatementOptsTo(itemTrim,
						termOptions["to"].([]map[string]interface{})[0]); err != nil {
						return confRead, err
					}
				}
				confRead.term = append(confRead.term, termOptions)
			case balt.CutPrefixInString(&itemTrim, "from "):
				if len(confRead.from) == 0 {
					confRead.from = append(confRead.from, genMapPolicyStatementOptsFrom())
				}
				if err := readPolicyStatementOptsFrom(itemTrim, confRead.from[0]); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "then "):
				if len(confRead.then) == 0 {
					confRead.then = append(confRead.then, genMapPolicyStatementOptsThen())
				}
				if err := readPolicyStatementOptsThen(itemTrim, confRead.then[0]); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "to "):
				if len(confRead.to) == 0 {
					confRead.to = append(confRead.to, genMapPolicyStatementOptsTo())
				}
				if err := readPolicyStatementOptsTo(itemTrim, confRead.to[0]); err != nil {
					return confRead, err
				}
			}
		}
	}

	return confRead, nil
}

func readPolicyStatementFwTableExport(policyName string, junSess *junos.Session,
) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"routing-options forwarding-table export" + junos.PipeDisplaySetRelative)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}
	for _, item := range strings.Split(showConfig, "\n") {
		if strings.Contains(item, junos.XMLStartTagConfigOut) {
			continue
		}
		if strings.Contains(item, junos.XMLEndTagConfigOut) {
			break
		}
		itemTrim := strings.TrimPrefix(item, junos.SetLS)
		if itemTrim == policyName || itemTrim == policyName+" " {
			return true, nil
		}
	}

	return false, nil
}

func delPolicyStatement(policyName string, junSess *junos.Session) error {
	configSet := []string{"delete policy-options policy-statement " + policyName}

	return junSess.ConfigSet(configSet)
}

func delPolicyStatementFwTableExport(policyName string, junSess *junos.Session) error {
	configSet := []string{"delete routing-options forwarding-table export " + policyName}

	return junSess.ConfigSet(configSet)
}

func fillPolicyStatementData(d *schema.ResourceData, policyStatementOptions policyStatementOptions) {
	if tfErr := d.Set("name", policyStatementOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("from", policyStatementOptions.from); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("then", policyStatementOptions.then); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("to", policyStatementOptions.to); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("term", policyStatementOptions.term); tfErr != nil {
		panic(tfErr)
	}
}

func setPolicyStatementOptsFrom(setPrefix string, opts map[string]interface{}) ([]string, error) {
	configSet := make([]string, 0)
	setPrefixFrom := setPrefix + " from "

	if opts["aggregate_contributor"].(bool) {
		configSet = append(configSet, setPrefixFrom+"aggregate-contributor")
	}
	for _, v := range sortSetOfString(opts["bgp_as_path"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"as-path "+v)
	}
	bgpASPathCalcLengthList := make([]int, 0)
	for _, block := range opts["bgp_as_path_calc_length"].(*schema.Set).List() {
		bgpASPathCalcLength := block.(map[string]interface{})
		count := bgpASPathCalcLength["count"].(int)
		if bchk.InSlice(count, bgpASPathCalcLengthList) {
			return configSet, fmt.Errorf("multiple blocks bgp_as_path_calc_length with the same count %d", count)
		}
		bgpASPathCalcLengthList = append(bgpASPathCalcLengthList, count)
		configSet = append(configSet,
			setPrefixFrom+"as-path-calc-length "+strconv.Itoa(count)+" "+bgpASPathCalcLength["match"].(string))
	}
	for _, v := range sortSetOfString(opts["bgp_as_path_group"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"as-path-group "+v)
	}
	bgpASPathUniqueCountList := make([]int, 0)
	for _, block := range opts["bgp_as_path_unique_count"].(*schema.Set).List() {
		bgpASPathUniqueCount := block.(map[string]interface{})
		count := bgpASPathUniqueCount["count"].(int)
		if bchk.InSlice(count, bgpASPathUniqueCountList) {
			return configSet, fmt.Errorf("multiple blocks bgp_as_path_unique_count with the same count %d", count)
		}
		bgpASPathUniqueCountList = append(bgpASPathUniqueCountList, count)
		configSet = append(configSet,
			setPrefixFrom+"as-path-unique-count "+strconv.Itoa(count)+" "+bgpASPathUniqueCount["match"].(string))
	}
	for _, v := range sortSetOfString(opts["bgp_community"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"community "+v)
	}
	bgpCommunityCountList := make([]int, 0)
	for _, block := range opts["bgp_community_count"].(*schema.Set).List() {
		bgpCommunityCount := block.(map[string]interface{})
		count := bgpCommunityCount["count"].(int)
		if bchk.InSlice(count, bgpCommunityCountList) {
			return configSet, fmt.Errorf("multiple blocks bgp_community_count with the same count %d", count)
		}
		bgpCommunityCountList = append(bgpCommunityCountList, count)
		configSet = append(configSet,
			setPrefixFrom+"community-count "+strconv.Itoa(count)+" "+bgpCommunityCount["match"].(string))
	}
	if opts["bgp_origin"].(string) != "" {
		configSet = append(configSet, setPrefixFrom+"origin "+opts["bgp_origin"].(string))
	}
	if v := opts["bgp_srte_discriminator"].(int); v != -1 {
		configSet = append(configSet, setPrefixFrom+"bgp-srte-discriminator "+strconv.Itoa(v))
	}
	if v := opts["color"].(int); v != -1 {
		configSet = append(configSet, setPrefixFrom+"color "+strconv.Itoa(v))
	}
	for _, v := range sortSetOfString(opts["evpn_esi"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"evpn-esi "+v)
	}
	if v := opts["evpn_mac_route"].(string); v != "" {
		configSet = append(configSet, setPrefixFrom+"evpn-mac-route "+v)
	}
	for _, v := range sortSetOfNumberToString(opts["evpn_tag"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"evpn-tag "+v)
	}
	if opts["family"].(string) != "" {
		configSet = append(configSet, setPrefixFrom+"family "+opts["family"].(string))
	}
	if opts["local_preference"].(int) != 0 {
		configSet = append(configSet, setPrefixFrom+"local-preference "+strconv.Itoa(opts["local_preference"].(int)))
	}
	if opts["routing_instance"].(string) != "" {
		configSet = append(configSet, setPrefixFrom+"instance "+opts["routing_instance"].(string))
	}
	for _, v := range sortSetOfString(opts["interface"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"interface "+v)
	}
	if opts["metric"].(int) != 0 {
		configSet = append(configSet, setPrefixFrom+"metric "+strconv.Itoa(opts["metric"].(int)))
	}
	for _, v := range sortSetOfString(opts["neighbor"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"neighbor "+v)
	}
	for _, v := range sortSetOfString(opts["next_hop"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"next-hop "+v)
	}
	if opts["next_hop_type_merged"].(bool) {
		configSet = append(configSet, setPrefixFrom+"next-hop-type merged")
	}
	for _, block := range opts["next_hop_weight"].(*schema.Set).List() {
		nextHopWeight := block.(map[string]interface{})
		configSet = append(configSet,
			setPrefixFrom+"nexthop-weight "+nextHopWeight["match"].(string)+" "+strconv.Itoa(nextHopWeight["weight"].(int)))
	}
	if opts["ospf_area"].(string) != "" {
		configSet = append(configSet, setPrefixFrom+"area "+opts["ospf_area"].(string))
	}
	for _, v := range opts["policy"].([]interface{}) {
		configSet = append(configSet, setPrefixFrom+"policy "+v.(string))
	}
	if opts["preference"].(int) != 0 {
		configSet = append(configSet, setPrefixFrom+"preference "+strconv.Itoa(opts["preference"].(int)))
	}
	for _, v := range sortSetOfString(opts["prefix_list"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"prefix-list "+v)
	}
	for _, v := range sortSetOfString(opts["protocol"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"protocol "+v)
	}
	for _, v := range opts["route_filter"].([]interface{}) {
		routeFilter := v.(map[string]interface{})
		setRoutFilter := setPrefixFrom + "route-filter " +
			routeFilter["route"].(string) + " " + routeFilter["option"].(string)
		if routeFilter["option_value"].(string) != "" {
			setRoutFilter += " " + routeFilter["option_value"].(string)
		}
		configSet = append(configSet, setRoutFilter)
	}
	if v := opts["route_type"].(string); v != "" {
		configSet = append(configSet, setPrefixFrom+"route-type "+v)
	}
	if v := opts["srte_color"].(int); v != -1 {
		configSet = append(configSet, setPrefixFrom+"srte-color "+strconv.Itoa(v))
	}
	if v := opts["state"].(string); v != "" {
		configSet = append(configSet, setPrefixFrom+"state "+v)
	}
	for _, v := range sortSetOfString(opts["tunnel_type"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"tunnel-type "+v)
	}
	if v := opts["validation_database"].(string); v != "" {
		configSet = append(configSet, setPrefixFrom+"validation-database "+v)
	}

	return configSet, nil
}

func setPolicyStatementOptsThen(setPrefix string, opts map[string]interface{}) ([]string, error) {
	configSet := make([]string, 0)
	setPrefixThen := setPrefix + " then "

	if opts["action"].(string) != "" {
		configSet = append(configSet, setPrefixThen+opts["action"].(string))
	}
	if opts["as_path_expand"].(string) != "" {
		if strings.Contains(opts["as_path_expand"].(string), "last-as") {
			configSet = append(configSet, setPrefixThen+"as-path-expand "+opts["as_path_expand"].(string))
		} else {
			configSet = append(configSet, setPrefixThen+"as-path-expand \""+opts["as_path_expand"].(string)+"\"")
		}
	}
	if opts["as_path_prepend"].(string) != "" {
		configSet = append(configSet, setPrefixThen+"as-path-prepend \""+opts["as_path_prepend"].(string)+"\"")
	}
	communityList := make([]string, 0)
	for _, v := range opts["community"].([]interface{}) {
		community := v.(map[string]interface{})
		setCommunityActVal := "community " + community["action"].(string) + " " + community["value"].(string)
		if bchk.InSlice(setCommunityActVal, communityList) {
			return configSet, fmt.Errorf("multiple blocks community with the same action %s and value %s",
				community["action"].(string), community["value"].(string))
		}
		communityList = append(communityList, setCommunityActVal)
		configSet = append(configSet, setPrefixThen+setCommunityActVal)
	}
	if opts["default_action"].(string) != "" {
		configSet = append(configSet, setPrefixThen+"default-action "+opts["default_action"].(string))
	}
	if opts["load_balance"].(string) != "" {
		configSet = append(configSet, setPrefixThen+"load-balance "+opts["load_balance"].(string))
	}
	for _, v := range opts["local_preference"].([]interface{}) {
		localPreference := v.(map[string]interface{})
		if localPreference["action"] == "none" {
			configSet = append(configSet, setPrefixThen+
				"local-preference "+strconv.Itoa(localPreference["value"].(int)))
		} else {
			configSet = append(configSet, setPrefixThen+
				"local-preference "+localPreference["action"].(string)+
				" "+strconv.Itoa(localPreference["value"].(int)))
		}
	}
	if opts["next"].(string) != "" {
		configSet = append(configSet, setPrefixThen+"next "+opts["next"].(string))
	}
	if opts["next_hop"].(string) != "" {
		configSet = append(configSet, setPrefixThen+"next-hop "+opts["next_hop"].(string))
	}
	for _, v := range opts["metric"].([]interface{}) {
		metric := v.(map[string]interface{})
		if metric["action"] == "none" {
			configSet = append(configSet, setPrefixThen+
				"metric "+strconv.Itoa(metric["value"].(int)))
		} else {
			configSet = append(configSet, setPrefixThen+
				"metric "+metric["action"].(string)+
				" "+strconv.Itoa(metric["value"].(int)))
		}
	}
	if opts["origin"].(string) != "" {
		configSet = append(configSet, setPrefixThen+"origin "+opts["origin"].(string))
	}
	for _, v := range opts["preference"].([]interface{}) {
		preference := v.(map[string]interface{})
		if preference["action"] == "none" {
			configSet = append(configSet, setPrefixThen+
				"preference "+strconv.Itoa(preference["value"].(int)))
		} else {
			configSet = append(configSet, setPrefixThen+
				"preference "+preference["action"].(string)+
				" "+strconv.Itoa(preference["value"].(int)))
		}
	}

	return configSet, nil
}

func setPolicyStatementOptsTo(setPrefix string, opts map[string]interface{}) []string {
	configSet := make([]string, 0)
	setPrefixTo := setPrefix + " to "

	for _, v := range sortSetOfString(opts["bgp_as_path"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTo+"as-path "+v)
	}
	for _, v := range sortSetOfString(opts["bgp_as_path_group"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTo+"as-path-group "+v)
	}
	for _, v := range sortSetOfString(opts["bgp_community"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTo+"community "+v)
	}
	if opts["bgp_origin"].(string) != "" {
		configSet = append(configSet, setPrefixTo+"origin "+opts["bgp_origin"].(string))
	}
	if opts["family"].(string) != "" {
		configSet = append(configSet, setPrefixTo+"family "+opts["family"].(string))
	}
	if opts["local_preference"].(int) != 0 {
		configSet = append(configSet, setPrefixTo+"local-preference "+strconv.Itoa(opts["local_preference"].(int)))
	}
	if opts["routing_instance"].(string) != "" {
		configSet = append(configSet, setPrefixTo+"instance "+opts["routing_instance"].(string))
	}
	for _, v := range sortSetOfString(opts["interface"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTo+"interface "+v)
	}
	if opts["metric"].(int) != 0 {
		configSet = append(configSet, setPrefixTo+"metric "+strconv.Itoa(opts["metric"].(int)))
	}
	for _, v := range sortSetOfString(opts["neighbor"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTo+"neighbor "+v)
	}
	for _, v := range sortSetOfString(opts["next_hop"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTo+"next-hop "+v)
	}
	if opts["ospf_area"].(string) != "" {
		configSet = append(configSet, setPrefixTo+"area "+opts["ospf_area"].(string))
	}
	for _, v := range opts["policy"].([]interface{}) {
		configSet = append(configSet, setPrefixTo+"policy "+v.(string))
	}
	if opts["preference"].(int) != 0 {
		configSet = append(configSet, setPrefixTo+"preference "+strconv.Itoa(opts["preference"].(int)))
	}
	for _, v := range sortSetOfString(opts["protocol"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTo+"protocol "+v)
	}

	return configSet
}

func readPolicyStatementOptsFrom(itemTrim string, fromMap map[string]interface{}) (err error) {
	switch {
	case itemTrim == "aggregate-contributor":
		fromMap["aggregate_contributor"] = true
	case balt.CutPrefixInString(&itemTrim, "as-path "):
		fromMap["bgp_as_path"] = append(fromMap["bgp_as_path"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "as-path-calc-length "):
		itemTrimFields := strings.Split(itemTrim, " ")
		count, err := strconv.Atoi(itemTrimFields[0])
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
		fromMap["bgp_as_path_calc_length"] = append(
			fromMap["bgp_as_path_calc_length"].([]map[string]interface{}),
			map[string]interface{}{
				"count": count,
				"match": strings.TrimPrefix(itemTrim, itemTrimFields[0]+" "),
			},
		)
	case balt.CutPrefixInString(&itemTrim, "as-path-group "):
		fromMap["bgp_as_path_group"] = append(fromMap["bgp_as_path_group"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "as-path-unique-count "):
		itemTrimFields := strings.Split(itemTrim, " ")
		count, err := strconv.Atoi(itemTrimFields[0])
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
		fromMap["bgp_as_path_unique_count"] = append(
			fromMap["bgp_as_path_unique_count"].([]map[string]interface{}),
			map[string]interface{}{
				"count": count,
				"match": strings.TrimPrefix(itemTrim, itemTrimFields[0]+" "),
			},
		)
	case balt.CutPrefixInString(&itemTrim, "community "):
		fromMap["bgp_community"] = append(fromMap["bgp_community"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "community-count "):
		itemTrimFields := strings.Split(itemTrim, " ")
		count, err := strconv.Atoi(itemTrimFields[0])
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
		fromMap["bgp_community_count"] = append(
			fromMap["bgp_community_count"].([]map[string]interface{}),
			map[string]interface{}{
				"count": count,
				"match": strings.TrimPrefix(itemTrim, itemTrimFields[0]+" "),
			},
		)
	case balt.CutPrefixInString(&itemTrim, "origin "):
		fromMap["bgp_origin"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "bgp-srte-discriminator "):
		fromMap["bgp_srte_discriminator"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "color "):
		fromMap["color"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "evpn-esi "):
		fromMap["evpn_esi"] = append(fromMap["evpn_esi"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "evpn-mac-route "):
		fromMap["evpn_mac_route"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "evpn-tag "):
		tag, err := strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
		fromMap["evpn_tag"] = append(fromMap["evpn_tag"].([]int), tag)
	case balt.CutPrefixInString(&itemTrim, "family "):
		fromMap["family"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "local-preference "):
		fromMap["local_preference"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "instance "):
		fromMap["routing_instance"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "interface "):
		fromMap["interface"] = append(fromMap["interface"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "metric "):
		fromMap["metric"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "neighbor "):
		fromMap["neighbor"] = append(fromMap["neighbor"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "next-hop "):
		fromMap["next_hop"] = append(fromMap["next_hop"].([]string), itemTrim)
	case itemTrim == "next-hop-type merged":
		fromMap["next_hop_type_merged"] = true
	case balt.CutPrefixInString(&itemTrim, "nexthop-weight "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 2 { // <match> <weight>
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "nexthop-weight", itemTrim)
		}
		weight, err := strconv.Atoi(itemTrimFields[1])
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
		fromMap["next_hop_weight"] = append(
			fromMap["next_hop_weight"].([]map[string]interface{}),
			map[string]interface{}{
				"match":  itemTrimFields[0],
				"weight": weight,
			},
		)
	case balt.CutPrefixInString(&itemTrim, "area "):
		fromMap["ospf_area"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "policy "):
		fromMap["policy"] = append(fromMap["policy"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "preference "):
		fromMap["preference"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "prefix-list "):
		fromMap["prefix_list"] = append(fromMap["prefix_list"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "protocol "):
		fromMap["protocol"] = append(fromMap["protocol"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "route-filter "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 2 { // <route> <option> <option_value>?
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "route-filter", itemTrim)
		}
		routeFilterMap := map[string]interface{}{
			"route":        itemTrimFields[0],
			"option":       itemTrimFields[1],
			"option_value": "",
		}
		if len(itemTrimFields) > 2 {
			routeFilterMap["option_value"] = itemTrimFields[2]
		}
		fromMap["route_filter"] = append(fromMap["route_filter"].([]map[string]interface{}), routeFilterMap)
	case balt.CutPrefixInString(&itemTrim, "route-type "):
		fromMap["route_type"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "srte-color "):
		fromMap["srte_color"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "state "):
		fromMap["state"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "tunnel-type "):
		fromMap["tunnel_type"] = append(fromMap["tunnel_type"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "validation-database "):
		fromMap["validation_database"] = itemTrim
	}

	return nil
}

func readPolicyStatementOptsThen(itemTrim string, thenMap map[string]interface{}) (err error) {
	switch {
	case itemTrim == "accept", itemTrim == "reject":
		thenMap["action"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "as-path-expand "):
		thenMap["as_path_expand"] = strings.Trim(itemTrim, "\"")
	case balt.CutPrefixInString(&itemTrim, "as-path-prepend "):
		thenMap["as_path_prepend"] = strings.Trim(itemTrim, "\"")
	case balt.CutPrefixInString(&itemTrim, "community "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 2 { // <action> <value>
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "community", itemTrim)
		}
		communityMap := map[string]interface{}{
			"action": itemTrimFields[0],
			"value":  itemTrimFields[1],
		}
		thenMap["community"] = append(thenMap["community"].([]map[string]interface{}), communityMap)
	case balt.CutPrefixInString(&itemTrim, "default-action "):
		thenMap["default_action"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "load-balance "):
		thenMap["load_balance"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "local-preference "):
		localPreferenceMap := map[string]interface{}{
			"action": "",
			"value":  0,
		}
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) == 1 { // <value>
			localPreferenceMap["action"] = "none"
			localPreferenceMap["value"], err = strconv.Atoi(itemTrimFields[0])
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		} else { // <action> <value>
			localPreferenceMap["action"] = itemTrimFields[0]
			localPreferenceMap["value"], err = strconv.Atoi(itemTrimFields[1])
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}

		thenMap["local_preference"] = append(thenMap["local_preference"].([]map[string]interface{}), localPreferenceMap)
	case balt.CutPrefixInString(&itemTrim, "next "):
		thenMap["next"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "next-hop "):
		thenMap["next_hop"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "metric "):
		metricMap := map[string]interface{}{
			"action": "",
			"value":  0,
		}
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) == 1 { // <value>
			metricMap["action"] = "none"
			metricMap["value"], err = strconv.Atoi(itemTrimFields[0])
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		} else { // <action> <value>
			metricMap["action"] = itemTrimFields[0]
			metricMap["value"], err = strconv.Atoi(itemTrimFields[1])
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
		thenMap["metric"] = append(thenMap["metric"].([]map[string]interface{}), metricMap)
	case balt.CutPrefixInString(&itemTrim, "origin "):
		thenMap["origin"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "preference "):
		preferenceMap := map[string]interface{}{
			"action": "",
			"value":  0,
		}
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) == 1 { // <value>
			preferenceMap["action"] = "none"
			preferenceMap["value"], err = strconv.Atoi(itemTrimFields[0])
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		} else { // <action> <value>
			preferenceMap["action"] = itemTrimFields[0]
			preferenceMap["value"], err = strconv.Atoi(itemTrimFields[1])
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
		thenMap["preference"] = append(thenMap["preference"].([]map[string]interface{}), preferenceMap)
	}

	return nil
}

func readPolicyStatementOptsTo(itemTrim string, toMap map[string]interface{}) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "as-path "):
		toMap["bgp_as_path"] = append(toMap["bgp_as_path"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "as-path-group "):
		toMap["bgp_as_path_group"] = append(toMap["bgp_as_path_group"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "community "):
		toMap["bgp_community"] = append(toMap["bgp_community"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "origin "):
		toMap["bgp_origin"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "family "):
		toMap["family"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "local-preference "):
		toMap["local_preference"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "instance "):
		toMap["routing_instance"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "interface "):
		toMap["interface"] = append(toMap["interface"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "metric "):
		toMap["metric"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "neighbor "):
		toMap["neighbor"] = append(toMap["neighbor"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "next-hop "):
		toMap["next_hop"] = append(toMap["next_hop"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "area "):
		toMap["ospf_area"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "policy "):
		toMap["policy"] = append(toMap["policy"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "preference "):
		toMap["preference"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "protocol "):
		toMap["protocol"] = append(toMap["protocol"].([]string), itemTrim)
	}

	return nil
}

func genMapPolicyStatementOptsFrom() map[string]interface{} {
	return map[string]interface{}{
		"aggregate_contributor":    false,
		"bgp_as_path":              make([]string, 0),
		"bgp_as_path_calc_length":  make([]map[string]interface{}, 0),
		"bgp_as_path_group":        make([]string, 0),
		"bgp_as_path_unique_count": make([]map[string]interface{}, 0),
		"bgp_community":            make([]string, 0),
		"bgp_community_count":      make([]map[string]interface{}, 0),
		"bgp_origin":               "",
		"bgp_srte_discriminator":   -1,
		"color":                    -1,
		"evpn_esi":                 make([]string, 0),
		"evpn_mac_route":           "",
		"evpn_tag":                 make([]int, 0),
		"family":                   "",
		"local_preference":         0,
		"routing_instance":         "",
		"interface":                make([]string, 0),
		"metric":                   0,
		"neighbor":                 make([]string, 0),
		"next_hop":                 make([]string, 0),
		"next_hop_type_merged":     false,
		"next_hop_weight":          make([]map[string]interface{}, 0),
		"ospf_area":                "",
		"policy":                   make([]string, 0),
		"preference":               0,
		"prefix_list":              make([]string, 0),
		"protocol":                 make([]string, 0),
		"route_filter":             make([]map[string]interface{}, 0),
		"route_type":               "",
		"srte_color":               -1,
		"state":                    "",
		"tunnel_type":              make([]string, 0),
		"validation_database":      "",
	}
}

func genMapPolicyStatementOptsThen() map[string]interface{} {
	return map[string]interface{}{
		"action":           "",
		"as_path_expand":   "",
		"as_path_prepend":  "",
		"community":        make([]map[string]interface{}, 0),
		"default_action":   "",
		"load_balance":     "",
		"local_preference": make([]map[string]interface{}, 0),
		"next":             "",
		"next_hop":         "",
		"metric":           make([]map[string]interface{}, 0),
		"origin":           "",
		"preference":       make([]map[string]interface{}, 0),
	}
}

func genMapPolicyStatementOptsTo() map[string]interface{} {
	return map[string]interface{}{
		"bgp_as_path":       make([]string, 0),
		"bgp_as_path_group": make([]string, 0),
		"bgp_community":     make([]string, 0),
		"bgp_origin":        "",
		"family":            "",
		"local_preference":  0,
		"routing_instance":  "",
		"interface":         make([]string, 0),
		"metric":            0,
		"neighbor":          make([]string, 0),
		"next_hop":          make([]string, 0),
		"ospf_area":         "",
		"policy":            make([]string, 0),
		"preference":        0,
		"protocol":          make([]string, 0),
	}
}
