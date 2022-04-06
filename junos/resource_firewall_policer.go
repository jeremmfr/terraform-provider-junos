package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type policerOptions struct {
	filterSpecific bool
	name           string
	ifExceeding    []map[string]interface{}
	then           []map[string]interface{}
}

func resourceFirewallPolicer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallPolicerCreate,
		ReadContext:   resourceFirewallPolicerRead,
		UpdateContext: resourceFirewallPolicerUpdate,
		DeleteContext: resourceFirewallPolicerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceFirewallPolicerImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"filter_specific": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"if_exceeding": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"burst_size_limit": {
							Type:     schema.TypeString,
							Required: true,
						},
						"bandwidth_percent": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(1, 100),
							ConflictsWith: []string{"if_exceeding.0.bandwidth_limit"},
						},
						"bandwidth_limit": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"if_exceeding.0.bandwidth_percent"},
						},
					},
				},
			},
			"then": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"discard": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"then.0.out_of_profile", "then.0.forwarding_class"},
						},
						"forwarding_class": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"then.0.discard"},
						},
						"loss_priority": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"then.0.discard"},
						},
						"out_of_profile": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"then.0.discard"},
						},
					},
				},
			},
		},
	}
}

func resourceFirewallPolicerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setFirewallPolicer(d, m, nil); err != nil {
			return diag.FromErr(err)
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
	firewallPolicerExists, err := checkFirewallPolicerExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if firewallPolicerExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("firewall policer %v already exists", d.Get("name").(string)))...)
	}

	if err := setFirewallPolicer(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_firewall_policer", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	firewallPolicerExists, err = checkFirewallPolicerExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if firewallPolicerExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("firewall policer %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceFirewallPolicerReadWJnprSess(d, m, jnprSess)...)
}

func resourceFirewallPolicerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceFirewallPolicerReadWJnprSess(d, m, jnprSess)
}

func resourceFirewallPolicerReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
) diag.Diagnostics {
	mutex.Lock()
	policerOptions, err := readFirewallPolicer(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if policerOptions.name == "" {
		d.SetId("")
	} else {
		fillFirewallPolicerData(d, policerOptions)
	}

	return nil
}

func resourceFirewallPolicerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delFirewallPolicer(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setFirewallPolicer(d, m, nil); err != nil {
			return diag.FromErr(err)
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
	if err := delFirewallPolicer(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setFirewallPolicer(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_firewall_policer", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceFirewallPolicerReadWJnprSess(d, m, jnprSess)...)
}

func resourceFirewallPolicerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delFirewallPolicer(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
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
	if err := delFirewallPolicer(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_firewall_policer", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceFirewallPolicerImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	firewallPolicerExists, err := checkFirewallPolicerExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !firewallPolicerExists {
		return nil, fmt.Errorf("don't find firewall policer with id '%v' (id must be <name>)", d.Id())
	}
	policerOptions, err := readFirewallPolicer(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillFirewallPolicerData(d, policerOptions)

	result[0] = d

	return result, nil
}

func checkFirewallPolicerExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command(cmdShowConfig+"firewall policer "+name+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setFirewallPolicer(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set firewall policer " + d.Get("name").(string)
	if d.Get("filter_specific").(bool) {
		configSet = append(configSet, setPrefix+" filter-specific")
	}
	for _, ifExceeding := range d.Get("if_exceeding").([]interface{}) {
		ifExceedingMap := ifExceeding.(map[string]interface{})
		configSet = append(configSet, setPrefix+
			" if-exceeding burst-size-limit "+ifExceedingMap["burst_size_limit"].(string))
		if ifExceedingMap["bandwidth_percent"].(int) != 0 {
			configSet = append(configSet, setPrefix+
				" if-exceeding bandwidth-percent "+strconv.Itoa(ifExceedingMap["bandwidth_percent"].(int)))
		}
		if ifExceedingMap["bandwidth_limit"].(string) != "" {
			configSet = append(configSet, setPrefix+
				" if-exceeding bandwidth-limit "+ifExceedingMap["bandwidth_limit"].(string))
		}
	}
	for _, then := range d.Get("then").([]interface{}) {
		if then != nil {
			thenMap := then.(map[string]interface{})
			if thenMap["discard"].(bool) {
				configSet = append(configSet, setPrefix+
					" then discard")
			}
			if thenMap["forwarding_class"].(string) != "" {
				configSet = append(configSet, setPrefix+
					" then forwarding-class "+thenMap["forwarding_class"].(string))
			}
			if thenMap["loss_priority"].(string) != "" {
				configSet = append(configSet, setPrefix+
					" then loss-priority "+thenMap["loss_priority"].(string))
			}
			if thenMap["out_of_profile"].(bool) {
				configSet = append(configSet, setPrefix+
					" then out-of-profile")
			}
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func readFirewallPolicer(name string, m interface{}, jnprSess *NetconfObject) (policerOptions, error) {
	sess := m.(*Session)
	var confRead policerOptions

	showConfig, err := sess.command(cmdShowConfig+"firewall policer "+name+pipeDisplaySetRelative, jnprSess)
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
			case itemTrim == "filter-specific":
				confRead.filterSpecific = true
			case strings.HasPrefix(itemTrim, "if-exceeding "):
				if len(confRead.ifExceeding) == 0 {
					confRead.ifExceeding = append(confRead.ifExceeding, map[string]interface{}{
						"burst_size_limit":  "",
						"bandwidth_percent": 0,
						"bandwidth_limit":   "",
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "if-exceeding burst-size-limit "):
					confRead.ifExceeding[0]["burst_size_limit"] = strings.TrimPrefix(itemTrim, "if-exceeding burst-size-limit ")
				case strings.HasPrefix(itemTrim, "if-exceeding bandwidth-percent "):
					confRead.ifExceeding[0]["bandwidth_percent"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrim, "if-exceeding bandwidth-percent "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "if-exceeding bandwidth-limit "):
					confRead.ifExceeding[0]["bandwidth_limit"] = strings.TrimPrefix(itemTrim, "if-exceeding bandwidth-limit ")
				}
			case strings.HasPrefix(itemTrim, "then "):
				if len(confRead.then) == 0 {
					confRead.then = append(confRead.then, map[string]interface{}{
						"discard":          false,
						"forwarding_class": "",
						"loss_priority":    "",
						"out_of_profile":   false,
					})
				}
				switch {
				case itemTrim == "then discard":
					confRead.then[0]["discard"] = true
				case strings.HasPrefix(itemTrim, "then forwarding-class "):
					confRead.then[0]["forwarding_class"] = strings.TrimPrefix(itemTrim, "then forwarding-class ")
				case strings.HasPrefix(itemTrim, "then loss-priority "):
					confRead.then[0]["loss_priority"] = strings.TrimPrefix(itemTrim, "then loss-priority ")
				case itemTrim == "then out-of-profile":
					confRead.then[0]["out_of_profile"] = true
				}
			}
		}
	}

	return confRead, nil
}

func delFirewallPolicer(policer string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete firewall policer "+policer)

	return sess.configSet(configSet, jnprSess)
}

func fillFirewallPolicerData(d *schema.ResourceData, policerOptions policerOptions) {
	if tfErr := d.Set("name", policerOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("filter_specific", policerOptions.filterSpecific); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("if_exceeding", policerOptions.ifExceeding); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("then", policerOptions.then); tfErr != nil {
		panic(tfErr)
	}
}
