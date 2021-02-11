package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type filterOptions struct {
	interfaceSpecific bool
	name              string
	family            string
	term              []map[string]interface{}
}

func resourceFirewallFilter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallFilterCreate,
		ReadContext:   resourceFirewallFilterRead,
		UpdateContext: resourceFirewallFilterUpdate,
		DeleteContext: resourceFirewallFilterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceFirewallFilterImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"family": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					inetWord, inet6Word, "any", "ccc", "mpls", "vpls", "ethernet-switching"}, false),
			},
			"interface_specific": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"term": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
						},
						"filter": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
						},
						"from": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"address_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"destination_address": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"destination_address_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"destination_port": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"destination_port_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"destination_prefix_list": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"destination_prefix_list_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"icmp_code": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"icmp_code_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"icmp_type": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"icmp_type_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"is_fragment": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"next_header": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"next_header_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"port": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"port_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"protocol": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"protocol_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"prefix_list": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"prefix_list_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"source_address": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"source_address_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"source_port": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"source_port_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"source_prefix_list": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"source_prefix_list_except": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"tcp_established": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"tcp_flags": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tcp_initial": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"then": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"accept", "reject", "discard", "next term"}, false),
									},
									"count": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"log": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"policer": {
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
									},
									"port_mirror": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"routing_instance": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"sample": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"service_accounting": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"syslog": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceFirewallFilterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	firewallFilterExists, err := checkFirewallFilterExists(d.Get("name").(string), d.Get("family").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if firewallFilterExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("firewall filter %v already exists", d.Get("name").(string)))
	}

	if err := setFirewallFilter(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_firewall_filter", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	firewallFilterExists, err = checkFirewallFilterExists(d.Get("name").(string), d.Get("family").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if firewallFilterExists {
		d.SetId(d.Get("name").(string) + idSeparator + d.Get("family").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("firewall filter %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceFirewallFilterReadWJnprSess(d, m, jnprSess)...)
}
func resourceFirewallFilterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceFirewallFilterReadWJnprSess(d, m, jnprSess)
}
func resourceFirewallFilterReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	filterOptions, err := readFirewallFilter(d.Get("name").(string), d.Get("family").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if filterOptions.name == "" {
		d.SetId("")
	} else {
		fillFirewallFilterData(d, filterOptions)
	}

	return nil
}
func resourceFirewallFilterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delFirewallFilter(d.Get("name").(string), d.Get("family").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setFirewallFilter(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_firewall_filter", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceFirewallFilterReadWJnprSess(d, m, jnprSess)...)
}
func resourceFirewallFilterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delFirewallFilter(d.Get("name").(string), d.Get("family").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_firewall_filter", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceFirewallFilterImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idList := strings.Split(d.Id(), idSeparator)
	if len(idList) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	firewallFilterExists, err := checkFirewallFilterExists(idList[0], idList[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !firewallFilterExists {
		return nil, fmt.Errorf("don't find firewall filter with id '%v' (id must be <name>"+idSeparator+"<family>)", d.Id())
	}
	filterOptions, err := readFirewallFilter(idList[0], idList[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillFirewallFilterData(d, filterOptions)

	result[0] = d

	return result, nil
}

func checkFirewallFilterExists(name, family string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	filterConfig, err := sess.command("show configuration "+
		"firewall family "+family+" filter "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if filterConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setFirewallFilter(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	var err error
	setPrefix := "set firewall family " + d.Get("family").(string) + " filter " + d.Get("name").(string)

	if d.Get("interface_specific").(bool) {
		configSet = append(configSet, setPrefix+" interface-specific")
	}
	for _, term := range d.Get("term").([]interface{}) {
		termMap := term.(map[string]interface{})
		setPrefixTerm := setPrefix + " term " + termMap["name"].(string)
		if termMap["filter"].(string) != "" {
			configSet = append(configSet, setPrefixTerm+" filter "+termMap["filter"].(string))
		}

		for _, from := range termMap["from"].([]interface{}) {
			configSet, err = setFirewallFilterOptsFrom(setPrefixTerm+" from ", configSet, from.(map[string]interface{}))
			if err != nil {
				return err
			}
		}
		for _, then := range termMap["then"].([]interface{}) {
			configSet = setFirewallFilterOptsThen(setPrefixTerm+" then ", configSet, then.(map[string]interface{}))
		}
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readFirewallFilter(filter, family string, m interface{}, jnprSess *NetconfObject) (filterOptions, error) {
	sess := m.(*Session)
	var confRead filterOptions

	filterConfig, err := sess.command("show configuration "+
		"firewall family "+family+" filter "+filter+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if filterConfig != emptyWord {
		confRead.name = filter
		confRead.family = family
		for _, item := range strings.Split(filterConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case itemTrim == "interface-specific":
				confRead.interfaceSpecific = true
			case strings.HasPrefix(itemTrim, "term "):
				termSplit := strings.Split(strings.TrimPrefix(itemTrim, "term "), " ")
				termOptions := map[string]interface{}{
					"name":   termSplit[0],
					"filter": "",
					"from":   make([]map[string]interface{}, 0),
					"then":   make([]map[string]interface{}, 0),
				}
				itemTrimTerm := strings.TrimPrefix(itemTrim, "term "+termSplit[0]+" ")
				if len(confRead.term) > 0 {
					termOptions, confRead.term = copyAndRemoveItemMapList("name", false, termOptions, confRead.term)
				}
				switch {
				case strings.HasPrefix(itemTrimTerm, "filter "):
					termOptions["filter"] = strings.TrimPrefix(itemTrimTerm, "filter ")
				case strings.HasPrefix(itemTrimTerm, "from "):
					termOptions["from"] = readFirewallFilterOptsFrom(strings.TrimPrefix(itemTrimTerm, "from "),
						termOptions["from"].([]map[string]interface{}))
				case strings.HasPrefix(itemTrimTerm, "then "):
					termOptions["then"] = readFirewallFilterOptsThen(strings.TrimPrefix(itemTrimTerm, "then "),
						termOptions["then"].([]map[string]interface{}))
				}
				confRead.term = append(confRead.term, termOptions)
			}
		}
	}

	return confRead, nil
}

func delFirewallFilter(filter, family string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete firewall family "+family+" filter "+filter)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func fillFirewallFilterData(d *schema.ResourceData, filterOptions filterOptions) {
	if tfErr := d.Set("name", filterOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family", filterOptions.family); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interface_specific", filterOptions.interfaceSpecific); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("term", filterOptions.term); tfErr != nil {
		panic(tfErr)
	}
}

func setFirewallFilterOptsFrom(setPrefixTermFrom string,
	configSet []string, fromMap map[string]interface{}) ([]string, error) {
	for _, address := range fromMap["address"].([]interface{}) {
		err := validateCIDRNetwork(address.(string))
		if err != nil {
			return nil, err
		}
		configSet = append(configSet, setPrefixTermFrom+"address "+address.(string))
	}
	for _, address := range fromMap["address_except"].([]interface{}) {
		err := validateCIDRNetwork(address.(string))
		if err != nil {
			return nil, err
		}
		configSet = append(configSet, setPrefixTermFrom+"address "+address.(string)+" except")
	}
	for _, address := range fromMap["destination_address"].([]interface{}) {
		err := validateCIDRNetwork(address.(string))
		if err != nil {
			return nil, err
		}
		configSet = append(configSet, setPrefixTermFrom+"destination-address "+address.(string))
	}
	for _, address := range fromMap["destination_address_except"].([]interface{}) {
		err := validateCIDRNetwork(address.(string))
		if err != nil {
			return nil, err
		}
		configSet = append(configSet, setPrefixTermFrom+"destination-address "+address.(string)+" except")
	}
	if len(fromMap["destination_port"].([]interface{})) > 0 &&
		len(fromMap["destination_port_except"].([]interface{})) > 0 {
		return configSet, fmt.Errorf("conflict between destination_port and destination_port_except")
	}
	for _, port := range fromMap["destination_port"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"destination-port "+port.(string))
	}
	for _, port := range fromMap["destination_port_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"destination-port-except "+port.(string))
	}
	for _, prefixList := range fromMap["destination_prefix_list"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"destination-prefix-list "+prefixList.(string))
	}
	for _, prefixList := range fromMap["destination_prefix_list_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"destination-prefix-list "+prefixList.(string)+" except")
	}
	if len(fromMap["icmp_code"].([]interface{})) > 0 && len(fromMap["icmp_code_except"].([]interface{})) > 0 {
		return nil, fmt.Errorf("conflict between icmp_code and icmp_code_except")
	}
	for _, icmp := range fromMap["icmp_code"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"icmp-code "+icmp.(string))
	}
	for _, icmp := range fromMap["icmp_code_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"icmp-code-except "+icmp.(string))
	}
	if len(fromMap["icmp_type"].([]interface{})) > 0 && len(fromMap["icmp_type_except"].([]interface{})) > 0 {
		return nil, fmt.Errorf("conflict between icmp_type and icmp_type_except")
	}
	for _, icmp := range fromMap["icmp_type"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"icmp-type "+icmp.(string))
	}
	for _, icmp := range fromMap["icmp_type_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"icmp-type-except "+icmp.(string))
	}
	if fromMap["is_fragment"].(bool) {
		configSet = append(configSet, setPrefixTermFrom+"is-fragment")
	}
	if len(fromMap["next_header"].([]interface{})) > 0 && len(fromMap["next_header_except"].([]interface{})) > 0 {
		return nil, fmt.Errorf("conflict between next_header and next_header_except")
	}
	for _, header := range fromMap["next_header"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"next-header "+header.(string))
	}
	for _, header := range fromMap["next_header_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"next-header-except "+header.(string))
	}
	if len(fromMap["port"].([]interface{})) > 0 && len(fromMap["port_except"].([]interface{})) > 0 {
		return configSet, fmt.Errorf("conflict between port and port_except")
	}
	for _, port := range fromMap["port"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"port "+port.(string))
	}
	for _, port := range fromMap["port_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"port-except "+port.(string))
	}
	for _, prefixList := range fromMap["prefix_list"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"prefix-list "+prefixList.(string))
	}
	for _, prefixList := range fromMap["prefix_list_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"prefix-list "+prefixList.(string)+" except")
	}
	if len(fromMap["protocol"].([]interface{})) > 0 && len(fromMap["protocol_except"].([]interface{})) > 0 {
		return nil, fmt.Errorf("conflict between protocol and protocol_except")
	}
	for _, protocol := range fromMap["protocol"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"protocol "+protocol.(string))
	}
	for _, protocol := range fromMap["protocol_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"protocol-except "+protocol.(string))
	}
	for _, address := range fromMap["source_address"].([]interface{}) {
		err := validateCIDRNetwork(address.(string))
		if err != nil {
			return nil, err
		}
		configSet = append(configSet, setPrefixTermFrom+"source-address "+address.(string))
	}
	for _, address := range fromMap["source_address_except"].([]interface{}) {
		err := validateCIDRNetwork(address.(string))
		if err != nil {
			return nil, err
		}
		configSet = append(configSet, setPrefixTermFrom+"source-address "+address.(string)+" except")
	}
	if len(fromMap["source_port"].([]interface{})) > 0 && len(fromMap["source_port_except"].([]interface{})) > 0 {
		return configSet, fmt.Errorf("conflict between source_port and source_port_except")
	}
	for _, port := range fromMap["source_port"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"source-port "+port.(string))
	}
	for _, port := range fromMap["source_port_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"source-port-except "+port.(string))
	}
	for _, prefixList := range fromMap["source_prefix_list"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"source-prefix-list "+prefixList.(string))
	}
	for _, prefixList := range fromMap["source_prefix_list_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"source-prefix-list "+prefixList.(string)+" except")
	}
	if (fromMap["tcp_established"].(bool) || fromMap["tcp_initial"].(bool)) && fromMap["tcp_flags"].(string) != "" {
		return configSet, fmt.Errorf("conflict between tcp_flags and tcp_initial|tcp_established")
	}
	if fromMap["tcp_established"].(bool) {
		configSet = append(configSet, setPrefixTermFrom+"tcp-established")
	}
	if fromMap["tcp_flags"].(string) != "" {
		configSet = append(configSet, setPrefixTermFrom+"tcp-flags \""+fromMap["tcp_flags"].(string)+"\"")
	}
	if fromMap["tcp_initial"].(bool) {
		configSet = append(configSet, setPrefixTermFrom+"tcp-initial")
	}

	return configSet, nil
}
func setFirewallFilterOptsThen(setPrefixTermThen string, configSet []string, thenMap map[string]interface{}) []string {
	if thenMap["action"].(string) != "" {
		configSet = append(configSet, setPrefixTermThen+thenMap["action"].(string))
	}
	if thenMap["count"].(string) != "" {
		configSet = append(configSet, setPrefixTermThen+"count "+thenMap["count"].(string))
	}
	if thenMap["log"].(bool) {
		configSet = append(configSet, setPrefixTermThen+"log")
	}
	if thenMap["policer"].(string) != "" {
		configSet = append(configSet, setPrefixTermThen+"policer "+thenMap["policer"].(string))
	}
	if thenMap["port_mirror"].(bool) {
		configSet = append(configSet, setPrefixTermThen+"port-mirror")
	}
	if thenMap["routing_instance"].(string) != "" {
		configSet = append(configSet, setPrefixTermThen+"routing-instance "+thenMap["routing_instance"].(string))
	}
	if thenMap["sample"].(bool) {
		configSet = append(configSet, setPrefixTermThen+"sample")
	}
	if thenMap["service_accounting"].(bool) {
		configSet = append(configSet, setPrefixTermThen+"service-accounting")
	}
	if thenMap["syslog"].(bool) {
		configSet = append(configSet, setPrefixTermThen+"syslog")
	}

	return configSet
}
func readFirewallFilterOptsFrom(item string,
	confReadElement []map[string]interface{}) []map[string]interface{} {
	fromMap := genMapFirewallFilterOptsFrom()
	if len(confReadElement) > 0 {
		for k, v := range confReadElement[0] {
			fromMap[k] = v
		}
	}
	switch {
	case strings.HasPrefix(item, "address "):
		if strings.HasSuffix(item, " except") {
			fromMap["address_except"] = append(fromMap["address_except"].([]string),
				strings.TrimSuffix(strings.TrimPrefix(item, "address "), " except"))
		} else {
			fromMap["address"] = append(fromMap["address"].([]string),
				strings.TrimPrefix(item, "address "))
		}
	case strings.HasPrefix(item, "destination-address "):
		if strings.HasSuffix(item, " except") {
			fromMap["destination_address_except"] = append(fromMap["destination_address_except"].([]string),
				strings.TrimSuffix(strings.TrimPrefix(item, "destination-address "), " except"))
		} else {
			fromMap["destination_address"] = append(fromMap["destination_address"].([]string),
				strings.TrimPrefix(item, "destination-address "))
		}
	case strings.HasPrefix(item, "destination-port "):
		fromMap["destination_port"] = append(fromMap["destination_port"].([]string),
			strings.TrimPrefix(item, "destination-port "))
	case strings.HasPrefix(item, "destination-port-except "):
		fromMap["destination_port_except"] = append(fromMap["destination_port_except"].([]string),
			strings.TrimPrefix(item, "destination-port-except "))
	case strings.HasPrefix(item, "destination-prefix-list "):
		if strings.HasSuffix(item, " except") {
			fromMap["destination_prefix_list_except"] = append(fromMap["destination_prefix_list_except"].([]string),
				strings.TrimSuffix(strings.TrimPrefix(item, "destination-prefix-list "), " except"))
		} else {
			fromMap["destination_prefix_list"] = append(fromMap["destination_prefix_list"].([]string),
				strings.TrimPrefix(item, "destination-prefix-list "))
		}
	case strings.HasPrefix(item, "icmp-code "):
		fromMap["icmp_code"] = append(fromMap["icmp_code"].([]string), strings.TrimPrefix(item, "icmp-code "))
	case strings.HasPrefix(item, "icmp-code-except "):
		fromMap["icmp_code_except"] = append(fromMap["icmp_code_except"].([]string),
			strings.TrimPrefix(item, "icmp-code-except "))
	case strings.HasPrefix(item, "icmp-type "):
		fromMap["icmp_type"] = append(fromMap["icmp_type"].([]string), strings.TrimPrefix(item, "icmp-type "))
	case strings.HasPrefix(item, "icmp-type-except "):
		fromMap["icmp_type_except"] = append(fromMap["icmp_type_except"].([]string),
			strings.TrimPrefix(item, "icmp-type-except "))
	case item == "is-fragment":
		fromMap["is_fragment"] = true
	case strings.HasPrefix(item, "next-header "):
		fromMap["next_header"] = append(fromMap["next_header"].([]string),
			strings.TrimPrefix(item, "next-header "))
	case strings.HasPrefix(item, "next-header-except "):
		fromMap["next_header_except"] = append(fromMap["next_header_except"].([]string),
			strings.TrimPrefix(item, "next-header-except "))
	case strings.HasPrefix(item, "port "):
		fromMap["port"] = append(fromMap["port"].([]string),
			strings.TrimPrefix(item, "port "))
	case strings.HasPrefix(item, "port-except "):
		fromMap["port_except"] = append(fromMap["port_except"].([]string),
			strings.TrimPrefix(item, "port-except "))
	case strings.HasPrefix(item, "protocol "):
		fromMap["protocol"] = append(fromMap["protocol"].([]string), strings.TrimPrefix(item, "protocol "))
	case strings.HasPrefix(item, "protocol-except "):
		fromMap["protocol_except"] = append(fromMap["protocol_except"].([]string),
			strings.TrimPrefix(item, "protocol-except "))
	case strings.HasPrefix(item, "prefix-list "):
		if strings.HasSuffix(item, " except") {
			fromMap["prefix_list_except"] = append(fromMap["prefix_list_except"].([]string),
				strings.TrimSuffix(strings.TrimPrefix(item, "prefix-list "), " except"))
		} else {
			fromMap["prefix_list"] = append(fromMap["prefix_list"].([]string),
				strings.TrimPrefix(item, "prefix-list "))
		}
	case strings.HasPrefix(item, "source-address "):
		if strings.HasSuffix(item, " except") {
			fromMap["source_address_except"] = append(fromMap["source_address_except"].([]string),
				strings.TrimSuffix(strings.TrimPrefix(item, "source-address "), " except"))
		} else {
			fromMap["source_address"] = append(fromMap["source_address"].([]string),
				strings.TrimPrefix(item, "source-address "))
		}
	case strings.HasPrefix(item, "source-port "):
		fromMap["source_port"] = append(fromMap["source_port"].([]string),
			strings.TrimPrefix(item, "source-port "))
	case strings.HasPrefix(item, "source-port-except "):
		fromMap["source_port_except"] = append(fromMap["source_port_except"].([]string),
			strings.TrimPrefix(item, "source-port-except "))
	case strings.HasPrefix(item, "source-prefix-list "):
		if strings.HasSuffix(item, " except") {
			fromMap["source_prefix_list_except"] = append(fromMap["source_prefix_list_except"].([]string),
				strings.TrimSuffix(strings.TrimPrefix(item, "source-prefix-list "), " except"))
		} else {
			fromMap["source_prefix_list"] = append(fromMap["source_prefix_list"].([]string),
				strings.TrimPrefix(item, "source-prefix-list "))
		}
	case item == "tcp-established":
		fromMap["tcp_established"] = true
	case strings.HasPrefix(item, "tcp-flags "):
		fromMap["tcp_flags"] = strings.Trim(strings.TrimPrefix(item, "tcp-flags "), "\"")
	case item == "tcp-initial":
		fromMap["tcp_initial"] = true
	}

	// override (maxItem = 1)
	return []map[string]interface{}{fromMap}
}
func readFirewallFilterOptsThen(item string,
	confReadElement []map[string]interface{}) []map[string]interface{} {
	thenMap := genMapFirewallFilterOptsThen()
	if len(confReadElement) > 0 {
		for k, v := range confReadElement[0] {
			thenMap[k] = v
		}
	}
	switch {
	case item == "accept",
		item == "reject",
		item == discardW,
		item == "next term":
		thenMap["action"] = item
	case strings.HasPrefix(item, "count "):
		thenMap["count"] = strings.TrimPrefix(item, "count ")
	case item == "log":
		thenMap["log"] = true
	case strings.HasPrefix(item, "policer "):
		thenMap["policer"] = strings.TrimPrefix(item, "policer ")
	case item == "port-mirror":
		thenMap["port_mirror"] = true
	case strings.HasPrefix(item, "routing-instance "):
		thenMap["routing_instance"] = strings.TrimPrefix(item, "routing-instance ")
	case item == "sample":
		thenMap["sample"] = true
	case item == "service-accounting":
		thenMap["service_accounting"] = true
	case item == "syslog":
		thenMap["syslog"] = true
	}
	// override (maxItem = 1)
	return []map[string]interface{}{thenMap}
}

func genMapFirewallFilterOptsFrom() map[string]interface{} {
	return map[string]interface{}{
		"address":                        make([]string, 0),
		"address_except":                 make([]string, 0),
		"destination_address":            make([]string, 0),
		"destination_address_except":     make([]string, 0),
		"destination_port":               make([]string, 0),
		"destination_port_except":        make([]string, 0),
		"destination_prefix_list":        make([]string, 0),
		"destination_prefix_list_except": make([]string, 0),
		"icmp_code":                      make([]string, 0),
		"icmp_code_except":               make([]string, 0),
		"icmp_type":                      make([]string, 0),
		"icmp_type_except":               make([]string, 0),
		"is_fragment":                    false,
		"next_header":                    make([]string, 0),
		"next_header_except":             make([]string, 0),
		"port":                           make([]string, 0),
		"port_except":                    make([]string, 0),
		"prefix_list":                    make([]string, 0),
		"prefix_list_except":             make([]string, 0),
		"protocol":                       make([]string, 0),
		"protocol_except":                make([]string, 0),
		"source_address":                 make([]string, 0),
		"source_address_except":          make([]string, 0),
		"source_port":                    make([]string, 0),
		"source_port_except":             make([]string, 0),
		"source_prefix_list":             make([]string, 0),
		"source_prefix_list_except":      make([]string, 0),
		"tcp_established":                false,
		"tcp_flags":                      "",
		"tcp_initial":                    false,
	}
}
func genMapFirewallFilterOptsThen() map[string]interface{} {
	return map[string]interface{}{
		"action":             "",
		"count":              "",
		"log":                false,
		"policer":            "",
		"port_mirror":        false,
		"routing_instance":   "",
		"sample":             false,
		"service_accounting": false,
		"syslog":             false,
	}
}
