package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		CreateContext: resourcePolicyoptionsPolicyStatementCreate,
		ReadContext:   resourcePolicyoptionsPolicyStatementRead,
		UpdateContext: resourcePolicyoptionsPolicyStatementUpdate,
		DeleteContext: resourcePolicyoptionsPolicyStatementDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePolicyoptionsPolicyStatementImport,
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
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"bgp_as_path_group": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"bgp_community": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
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
			Type:     schema.TypeString,
			Optional: true,
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
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"preference": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"prefix_list": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
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
						ValidateFunc: validation.StringInSlice([]string{"add", deleteW, setW}, false),
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
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"bgp_as_path_group": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"bgp_community": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
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
			Type:     schema.TypeString,
			Optional: true,
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
			Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourcePolicyoptionsPolicyStatementCreate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setPolicyStatement(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		if d.Get("add_it_to_forwarding_table_export").(bool) {
			if err := setPolicyStatementFwTableExport(d.Get("name").(string), m, nil); err != nil {
				return diag.FromErr(err)
			}
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	policyStatementExists, err := checkPolicyStatementExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyStatementExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("policy-options policy-statement %v already exists", d.Get("name").(string)))...)
	}

	if err := setPolicyStatement(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if d.Get("add_it_to_forwarding_table_export").(bool) {
		if err := setPolicyStatementFwTableExport(d.Get("name").(string), m, jnprSess); err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	warns, err := sess.commitConf("create resource junos_policyoptions_policy_statement", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	policyStatementExists, err = checkPolicyStatementExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyStatementExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("policy-options policy-statement %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourcePolicyoptionsPolicyStatementReadWJnprSess(d, m, jnprSess)...)
}

func resourcePolicyoptionsPolicyStatementRead(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourcePolicyoptionsPolicyStatementReadWJnprSess(d, m, jnprSess)
}

func resourcePolicyoptionsPolicyStatementReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	policyStatementOptions, err := readPolicyStatement(d.Get("name").(string), m, jnprSess)
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	if d.Get("add_it_to_forwarding_table_export").(bool) {
		export, err := readPolicyStatementFwTableExport(d.Get("name").(string), m, jnprSess)
		if err != nil {
			mutex.Unlock()

			return diag.FromErr(err)
		}
		if !export {
			if tfErr := d.Set("add_it_to_forwarding_table_export", false); tfErr != nil {
				panic(tfErr)
			}
		}
	}
	mutex.Unlock()

	if policyStatementOptions.name == "" {
		d.SetId("")
	} else {
		fillPolicyStatementData(d, policyStatementOptions)
	}

	return nil
}

func resourcePolicyoptionsPolicyStatementUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delPolicyStatement(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if d.HasChange("add_it_to_forwarding_table_export") {
			if o, _ := d.GetChange("add_it_to_forwarding_table_export"); o.(bool) {
				if err := delPolicyStatementFwTableExport(d.Get("name").(string), m, nil); err != nil {
					return diag.FromErr(err)
				}
			}
		}
		if err := setPolicyStatement(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		if d.Get("add_it_to_forwarding_table_export").(bool) {
			if err := setPolicyStatementFwTableExport(d.Get("name").(string), m, nil); err != nil {
				return diag.FromErr(err)
			}
		}
		d.Partial(false)

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delPolicyStatement(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if d.HasChange("add_it_to_forwarding_table_export") {
		if o, _ := d.GetChange("add_it_to_forwarding_table_export"); o.(bool) {
			if err := delPolicyStatementFwTableExport(d.Get("name").(string), m, jnprSess); err != nil {
				appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
		}
	}
	if err := setPolicyStatement(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if d.Get("add_it_to_forwarding_table_export").(bool) {
		if err := setPolicyStatementFwTableExport(d.Get("name").(string), m, jnprSess); err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	warns, err := sess.commitConf("update resource junos_policyoptions_policy_statement", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourcePolicyoptionsPolicyStatementReadWJnprSess(d, m, jnprSess)...)
}

func resourcePolicyoptionsPolicyStatementDelete(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delPolicyStatement(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if d.Get("add_it_to_forwarding_table_export").(bool) {
			if err := delPolicyStatementFwTableExport(d.Get("name").(string), m, nil); err != nil {
				return diag.FromErr(err)
			}
		}

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delPolicyStatement(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if d.Get("add_it_to_forwarding_table_export").(bool) {
		if err := delPolicyStatementFwTableExport(d.Get("name").(string), m, jnprSess); err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	warns, err := sess.commitConf("delete resource junos_policyoptions_policy_statement", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourcePolicyoptionsPolicyStatementImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	policyStatementExists, err := checkPolicyStatementExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !policyStatementExists {
		return nil, fmt.Errorf("don't find policy-options policy-statement with id '%v' (id must be <name>)", d.Id())
	}
	policyStatementOptions, err := readPolicyStatement(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillPolicyStatementData(d, policyStatementOptions)

	result[0] = d

	return result, nil
}

func checkPolicyStatementExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command(cmdShowConfig+
		"policy-options policy-statement "+name+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setPolicyStatement(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set policy-options policy-statement " + d.Get("name").(string)
	for _, from := range d.Get("from").([]interface{}) {
		if from != nil {
			configSetFrom := setPolicyStatementOptsFrom(setPrefix, from.(map[string]interface{}))
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
		if bchk.StringInSlice(termMap["name"].(string), termNameList) {
			return fmt.Errorf("multiple blocks term with the same name %s", termMap["name"].(string))
		}
		termNameList = append(termNameList, termMap["name"].(string))
		setPrefixTerm := setPrefix + " term " + termMap["name"].(string)
		for _, from := range termMap["from"].([]interface{}) {
			if from != nil {
				configSetFrom := setPolicyStatementOptsFrom(setPrefixTerm, from.(map[string]interface{}))
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

	return sess.configSet(configSet, jnprSess)
}

func setPolicyStatementFwTableExport(policyName string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"set routing-options forwarding-table export " + policyName}

	return sess.configSet(configSet, jnprSess)
}

func readPolicyStatement(name string, m interface{}, jnprSess *NetconfObject) (policyStatementOptions, error) {
	sess := m.(*Session)
	var confRead policyStatementOptions

	showConfig, err := sess.command(cmdShowConfig+
		"policy-options policy-statement "+name+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "term "):
				itemTermList := strings.Split(strings.TrimPrefix(itemTrim, "term "), " ")
				termOptions := map[string]interface{}{
					"name": itemTermList[0],
					"from": make([]map[string]interface{}, 0),
					"then": make([]map[string]interface{}, 0),
					"to":   make([]map[string]interface{}, 0),
				}
				itemTrimTerm := strings.TrimPrefix(itemTrim, "term "+itemTermList[0]+" ")
				confRead.term = copyAndRemoveItemMapList("name", termOptions, confRead.term)
				switch {
				case strings.HasPrefix(itemTrimTerm, "from "):
					if len(termOptions["from"].([]map[string]interface{})) == 0 {
						termOptions["from"] = append(termOptions["from"].([]map[string]interface{}),
							genMapPolicyStatementOptsFrom())
					}
					if err := readPolicyStatementOptsFrom(strings.TrimPrefix(itemTrimTerm, "from "),
						termOptions["from"].([]map[string]interface{})[0]); err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrimTerm, "then "):
					if len(termOptions["then"].([]map[string]interface{})) == 0 {
						termOptions["then"] = append(termOptions["then"].([]map[string]interface{}),
							genMapPolicyStatementOptsThen())
					}
					if err := readPolicyStatementOptsThen(strings.TrimPrefix(itemTrimTerm, "then "),
						termOptions["then"].([]map[string]interface{})[0]); err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrimTerm, "to "):
					if len(termOptions["to"].([]map[string]interface{})) == 0 {
						termOptions["to"] = append(termOptions["to"].([]map[string]interface{}),
							genMapPolicyStatementOptsTo())
					}
					if err := readPolicyStatementOptsTo(strings.TrimPrefix(itemTrimTerm, "to "),
						termOptions["to"].([]map[string]interface{})[0]); err != nil {
						return confRead, err
					}
				}
				confRead.term = append(confRead.term, termOptions)
			case strings.HasPrefix(itemTrim, "from "):
				if len(confRead.from) == 0 {
					confRead.from = append(confRead.from, genMapPolicyStatementOptsFrom())
				}
				if err := readPolicyStatementOptsFrom(strings.TrimPrefix(itemTrim, "from "), confRead.from[0]); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "then "):
				if len(confRead.then) == 0 {
					confRead.then = append(confRead.then, genMapPolicyStatementOptsThen())
				}
				if err := readPolicyStatementOptsThen(strings.TrimPrefix(itemTrim, "then "), confRead.then[0]); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "to "):
				if len(confRead.to) == 0 {
					confRead.to = append(confRead.to, genMapPolicyStatementOptsTo())
				}
				if err := readPolicyStatementOptsTo(strings.TrimPrefix(itemTrim, "to "), confRead.to[0]); err != nil {
					return confRead, err
				}
			}
		}
	}

	return confRead, nil
}

func readPolicyStatementFwTableExport(policyName string,
	m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command(cmdShowConfig+
		"routing-options forwarding-table export"+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}
	for _, item := range strings.Split(showConfig, "\n") {
		if strings.Contains(item, xmlStartTagConfigOut) {
			continue
		}
		if strings.Contains(item, xmlEndTagConfigOut) {
			break
		}
		itemTrim := strings.TrimPrefix(item, setLS)
		if itemTrim == policyName || itemTrim == policyName+" " {
			return true, nil
		}
	}

	return false, nil
}

func delPolicyStatement(policyName string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete policy-options policy-statement " + policyName}

	return sess.configSet(configSet, jnprSess)
}

func delPolicyStatementFwTableExport(policyName string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete routing-options forwarding-table export " + policyName}

	return sess.configSet(configSet, jnprSess)
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

func setPolicyStatementOptsFrom(setPrefix string, opts map[string]interface{}) []string {
	configSet := make([]string, 0)
	setPrefixFrom := setPrefix + " from "

	if opts["aggregate_contributor"].(bool) {
		configSet = append(configSet, setPrefixFrom+"aggregate-contributor")
	}
	for _, v := range sortSetOfString(opts["bgp_as_path"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"as-path "+v)
	}
	for _, v := range sortSetOfString(opts["bgp_as_path_group"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"as-path-group "+v)
	}
	for _, v := range sortSetOfString(opts["bgp_community"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixFrom+"community "+v)
	}
	if opts["bgp_origin"].(string) != "" {
		configSet = append(configSet, setPrefixFrom+"origin "+opts["bgp_origin"].(string))
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

	return configSet
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
		if bchk.StringInSlice(setCommunityActVal, communityList) {
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

func readPolicyStatementOptsFrom(item string, fromMap map[string]interface{}) error {
	switch {
	case item == "aggregate-contributor":
		fromMap["aggregate_contributor"] = true
	case strings.HasPrefix(item, "as-path "):
		fromMap["bgp_as_path"] = append(fromMap["bgp_as_path"].([]string), strings.TrimPrefix(item, "as-path "))
	case strings.HasPrefix(item, "as-path-group "):
		fromMap["bgp_as_path_group"] = append(fromMap["bgp_as_path_group"].([]string),
			strings.TrimPrefix(item, "as-path-group "))
	case strings.HasPrefix(item, "community "):
		fromMap["bgp_community"] = append(fromMap["bgp_community"].([]string), strings.TrimPrefix(item, "community "))
	case strings.HasPrefix(item, "origin "):
		fromMap["bgp_origin"] = strings.TrimPrefix(item, "origin ")
	case strings.HasPrefix(item, "family "):
		fromMap["family"] = strings.TrimPrefix(item, "family ")
	case strings.HasPrefix(item, "local-preference "):
		var err error
		fromMap["local_preference"], err = strconv.Atoi(strings.TrimPrefix(item, "local-preference "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, item, err)
		}
	case strings.HasPrefix(item, "instance "):
		fromMap["routing_instance"] = strings.TrimPrefix(item, "instance ")
	case strings.HasPrefix(item, "interface "):
		fromMap["interface"] = append(fromMap["interface"].([]string), strings.TrimPrefix(item, "interface "))
	case strings.HasPrefix(item, "metric "):
		var err error
		fromMap["metric"], err = strconv.Atoi(strings.TrimPrefix(item, "metric "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, item, err)
		}
	case strings.HasPrefix(item, "neighbor "):
		fromMap["neighbor"] = append(fromMap["neighbor"].([]string), strings.TrimPrefix(item, "neighbor "))
	case strings.HasPrefix(item, "next-hop "):
		fromMap["next_hop"] = append(fromMap["next_hop"].([]string), strings.TrimPrefix(item, "next-hop "))
	case strings.HasPrefix(item, "area "):
		fromMap["ospf_area"] = strings.TrimPrefix(item, "area ")
	case strings.HasPrefix(item, "policy "):
		fromMap["policy"] = append(fromMap["policy"].([]string), strings.TrimPrefix(item, "policy "))
	case strings.HasPrefix(item, "preference "):
		var err error
		fromMap["preference"], err = strconv.Atoi(strings.TrimPrefix(item, "preference "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, item, err)
		}
	case strings.HasPrefix(item, "prefix-list "):
		fromMap["prefix_list"] = append(fromMap["prefix_list"].([]string), strings.TrimPrefix(item, "prefix-list "))
	case strings.HasPrefix(item, "protocol "):
		fromMap["protocol"] = append(fromMap["protocol"].([]string), strings.TrimPrefix(item, "protocol "))
	case strings.HasPrefix(item, "route-filter "):
		routeFilterMap := map[string]interface{}{
			"route":        "",
			"option":       "",
			"option_value": "",
		}
		itemSplit := strings.Split(item, " ")
		routeFilterMap["route"] = itemSplit[1]
		routeFilterMap["option"] = itemSplit[2]
		if len(itemSplit) > 3 {
			routeFilterMap["option_value"] = itemSplit[3]
		}
		fromMap["route_filter"] = append(fromMap["route_filter"].([]map[string]interface{}), routeFilterMap)
	}

	return nil
}

func readPolicyStatementOptsThen(item string, thenMap map[string]interface{}) error {
	switch {
	case strings.HasPrefix(item, "accept"),
		strings.HasPrefix(item, "reject"):
		thenMap["action"] = item
	case strings.HasPrefix(item, "as-path-expand "):
		thenMap["as_path_expand"] = strings.Trim(strings.TrimPrefix(item, "as-path-expand "), "\"")
	case strings.HasPrefix(item, "as-path-prepend "):
		thenMap["as_path_prepend"] = strings.Trim(strings.TrimPrefix(item, "as-path-prepend "), "\"")
	case strings.HasPrefix(item, "community "):
		communityMap := map[string]interface{}{
			"action": "",
			"value":  "",
		}
		itemSplit := strings.Split(item, " ")
		communityMap["action"] = itemSplit[1]
		communityMap["value"] = itemSplit[2]
		thenMap["community"] = append(thenMap["community"].([]map[string]interface{}), communityMap)
	case strings.HasPrefix(item, "default-action "):
		thenMap["default_action"] = strings.TrimPrefix(item, "default-action ")
	case strings.HasPrefix(item, "load-balance "):
		thenMap["load_balance"] = strings.TrimPrefix(item, "load-balance ")
	case strings.HasPrefix(item, "local-preference "):
		localPreferenceMap := map[string]interface{}{
			"action": "",
			"value":  0,
		}
		itemSplit := strings.Split(item, " ")
		var err error
		if len(itemSplit) == 2 {
			localPreferenceMap["action"] = "none"
			localPreferenceMap["value"], err = strconv.Atoi(itemSplit[1])
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, item, err)
			}
		} else {
			localPreferenceMap["action"] = itemSplit[1]
			localPreferenceMap["value"], err = strconv.Atoi(itemSplit[2])
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, item, err)
			}
		}

		thenMap["local_preference"] = append(thenMap["local_preference"].([]map[string]interface{}), localPreferenceMap)
	case strings.HasPrefix(item, "next "):
		thenMap["next"] = strings.TrimPrefix(item, "next ")
	case strings.HasPrefix(item, "next-hop "):
		thenMap["next_hop"] = strings.TrimPrefix(item, "next-hop ")
	case strings.HasPrefix(item, "metric "):
		metricMap := map[string]interface{}{
			"action": "",
			"value":  0,
		}
		itemSplit := strings.Split(item, " ")
		var err error
		if len(itemSplit) == 2 {
			metricMap["action"] = "none"
			metricMap["value"], err = strconv.Atoi(itemSplit[1])
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, item, err)
			}
		} else {
			metricMap["action"] = itemSplit[1]
			metricMap["value"], err = strconv.Atoi(itemSplit[2])
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, item, err)
			}
		}
		thenMap["metric"] = append(thenMap["metric"].([]map[string]interface{}), metricMap)
	case strings.HasPrefix(item, "origin "):
		thenMap["origin"] = strings.TrimPrefix(item, "origin ")
	case strings.HasPrefix(item, "preference "):
		preferenceMap := map[string]interface{}{
			"action": "",
			"value":  0,
		}
		itemSplit := strings.Split(item, " ")
		var err error
		if len(itemSplit) == 2 {
			preferenceMap["action"] = "none"
			preferenceMap["value"], err = strconv.Atoi(itemSplit[1])
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, item, err)
			}
		} else {
			preferenceMap["action"] = itemSplit[1]
			preferenceMap["value"], err = strconv.Atoi(itemSplit[2])
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, item, err)
			}
		}
		thenMap["preference"] = append(thenMap["preference"].([]map[string]interface{}), preferenceMap)
	}

	return nil
}

func readPolicyStatementOptsTo(item string, toMap map[string]interface{}) error {
	switch {
	case strings.HasPrefix(item, "as-path "):
		toMap["bgp_as_path"] = append(toMap["bgp_as_path"].([]string), strings.TrimPrefix(item, "as-path "))
	case strings.HasPrefix(item, "as-path-group "):
		toMap["bgp_as_path_group"] = append(toMap["bgp_as_path_group"].([]string), strings.TrimPrefix(item, "as-path-group "))
	case strings.HasPrefix(item, "community "):
		toMap["bgp_community"] = append(toMap["bgp_community"].([]string), strings.TrimPrefix(item, "community "))
	case strings.HasPrefix(item, "origin "):
		toMap["bgp_origin"] = strings.TrimPrefix(item, "origin ")
	case strings.HasPrefix(item, "family "):
		toMap["family"] = strings.TrimPrefix(item, "family ")
	case strings.HasPrefix(item, "local-preference "):
		var err error
		toMap["local_preference"], err = strconv.Atoi(strings.TrimPrefix(item, "local-preference "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, item, err)
		}
	case strings.HasPrefix(item, "instance "):
		toMap["routing_instance"] = strings.TrimPrefix(item, "instance ")
	case strings.HasPrefix(item, "interface "):
		toMap["interface"] = append(toMap["interface"].([]string), strings.TrimPrefix(item, "interface "))
	case strings.HasPrefix(item, "metric "):
		var err error
		toMap["metric"], err = strconv.Atoi(strings.TrimPrefix(item, "metric "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, item, err)
		}
	case strings.HasPrefix(item, "neighbor "):
		toMap["neighbor"] = append(toMap["neighbor"].([]string), strings.TrimPrefix(item, "neighbor "))
	case strings.HasPrefix(item, "next-hop "):
		toMap["next_hop"] = append(toMap["next_hop"].([]string), strings.TrimPrefix(item, "next-hop "))
	case strings.HasPrefix(item, "area "):
		toMap["ospf_area"] = strings.TrimPrefix(item, "area ")
	case strings.HasPrefix(item, "policy "):
		toMap["policy"] = append(toMap["policy"].([]string), strings.TrimPrefix(item, "policy "))
	case strings.HasPrefix(item, "preference "):
		var err error
		toMap["preference"], err = strconv.Atoi(strings.TrimPrefix(item, "preference "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, item, err)
		}
	case strings.HasPrefix(item, "protocol "):
		toMap["protocol"] = append(toMap["protocol"].([]string), strings.TrimPrefix(item, "protocol "))
	}

	return nil
}

func genMapPolicyStatementOptsFrom() map[string]interface{} {
	return map[string]interface{}{
		"aggregate_contributor": false,
		"bgp_as_path":           make([]string, 0),
		"bgp_as_path_group":     make([]string, 0),
		"bgp_community":         make([]string, 0),
		"bgp_origin":            "",
		"family":                "",
		"local_preference":      0,
		"routing_instance":      "",
		"interface":             make([]string, 0),
		"metric":                0,
		"neighbor":              make([]string, 0),
		"next_hop":              make([]string, 0),
		"ospf_area":             "",
		"policy":                make([]string, 0),
		"preference":            0,
		"prefix_list":           make([]string, 0),
		"protocol":              make([]string, 0),
		"route_filter":          make([]map[string]interface{}, 0),
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
