package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type filterOptions struct {
	interfaceSpecific bool
	name              string
	family            string
	term              []map[string]interface{}
}

func resourceFirewallFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallFilterCreate,
		Read:   resourceFirewallFilterRead,
		Update: resourceFirewallFilterUpdate,
		Delete: resourceFirewallFilterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceFirewallFilterImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"family": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !stringInSlice(value, []string{inetWord, inet6Word, "any", "ccc", "mpls", "vpls", "ethernet-switching"}) {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not valid family", value, k))
					}
					return
				},
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateNameObjectJunos(),
						},
						"filter": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateNameObjectJunos(),
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
									"tcp_flags": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tcp_initial": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"tcp_established": {
										Type:     schema.TypeBool,
										Optional: true,
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
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											value := v.(string)
											if !stringInSlice(value, []string{"accept", "reject", "discard", "next term"}) {
												errors = append(errors, fmt.Errorf(
													"%q for %q is not valid acceptance", value, k))
											}
											return
										},
									},
									"count": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"routing_instance": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"policer": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateNameObjectJunos(),
									},
									"log": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"syslog": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"port_mirror": {
										Type:     schema.TypeBool,
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceFirewallFilterCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	firewallFilterExists, err := checkFirewallFilterExists(d.Get("name").(string), d.Get("family").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if firewallFilterExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("firewall filter %v already exists", d.Get("name").(string))
	}

	err = setFirewallFilter(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	firewallFilterExists, err = checkFirewallFilterExists(d.Get("name").(string), d.Get("family").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if firewallFilterExists {
		d.SetId(d.Get("name").(string) + idSeparator + d.Get("family").(string))
	} else {
		return fmt.Errorf("firewall filter %v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourceFirewallFilterRead(d, m)
}
func resourceFirewallFilterRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	filterOptions, err := readFirewallFilter(d.Get("name").(string), d.Get("family").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if filterOptions.name == "" {
		d.SetId("")
	} else {
		fillFirewallFilterData(d, filterOptions)
	}
	return nil
}
func resourceFirewallFilterUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delFirewallFilter(d.Get("name").(string), d.Get("family").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setFirewallFilter(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourceFirewallFilterRead(d, m)
}
func resourceFirewallFilterDelete(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delFirewallFilter(d.Get("name").(string), d.Get("family").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
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
		configSet = append(configSet, setPrefix+" interface-specific\n")
	}
	for _, term := range d.Get("term").([]interface{}) {
		termMap := term.(map[string]interface{})
		setPrefixTerm := setPrefix + " term " + termMap["name"].(string)
		if termMap["filter"].(string) != "" {
			configSet = append(configSet, setPrefixTerm+" filter "+termMap["filter"].(string)+"\n")
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

	err = sess.configSet(configSet, jnprSess)
	if err != nil {
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
			case strings.HasPrefix(itemTrim, "interface-specific"):
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
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}

func delFirewallFilter(filter, family string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete firewall family "+family+" filter "+filter+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func fillFirewallFilterData(d *schema.ResourceData, filterOptions filterOptions) {
	tfErr := d.Set("name", filterOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("family", filterOptions.family)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("interface_specific", filterOptions.interfaceSpecific)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("term", filterOptions.term)
	if tfErr != nil {
		panic(tfErr)
	}
}

func setFirewallFilterOptsFrom(setPrefixTermFrom string,
	configSet []string, fromMap map[string]interface{}) ([]string, error) {
	for _, address := range fromMap["address"].([]interface{}) {
		err := validateNetwork(address.(string))
		if err != nil {
			return nil, err
		}
		configSet = append(configSet, setPrefixTermFrom+"address "+address.(string)+"\n")
	}
	for _, address := range fromMap["address_except"].([]interface{}) {
		err := validateNetwork(address.(string))
		if err != nil {
			return nil, err
		}
		configSet = append(configSet, setPrefixTermFrom+"address "+address.(string)+" except\n")
	}
	if len(fromMap["port"].([]interface{})) > 0 && len(fromMap["port_except"].([]interface{})) > 0 {
		return configSet, fmt.Errorf("conflict between port and port_except")
	}
	for _, port := range fromMap["port"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"port "+port.(string)+"\n")
	}
	for _, port := range fromMap["port_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"port-except "+port.(string)+"\n")
	}
	for _, prefixList := range fromMap["prefix_list"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"prefix-list "+prefixList.(string)+"\n")
	}
	for _, prefixList := range fromMap["prefix_list_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"prefix-list "+prefixList.(string)+" except\n")
	}
	for _, address := range fromMap["destination_address"].([]interface{}) {
		err := validateNetwork(address.(string))
		if err != nil {
			return nil, err
		}
		configSet = append(configSet, setPrefixTermFrom+"destination-address "+address.(string)+"\n")
	}
	for _, address := range fromMap["destination_address_except"].([]interface{}) {
		err := validateNetwork(address.(string))
		if err != nil {
			return nil, err
		}
		configSet = append(configSet, setPrefixTermFrom+"destination-address "+address.(string)+" except\n")
	}
	if len(fromMap["destination_port"].([]interface{})) > 0 &&
		len(fromMap["destination_port_except"].([]interface{})) > 0 {
		return configSet, fmt.Errorf("conflict between destination_port and destination_port_except")
	}
	for _, port := range fromMap["destination_port"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"destination-port "+port.(string)+"\n")
	}
	for _, port := range fromMap["destination_port_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"destination-port-except "+port.(string)+"\n")
	}
	for _, prefixList := range fromMap["destination_prefix_list"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"destination-prefix-list "+prefixList.(string)+"\n")
	}
	for _, prefixList := range fromMap["destination_prefix_list_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"destination-prefix-list "+prefixList.(string)+" except\n")
	}
	for _, address := range fromMap["source_address"].([]interface{}) {
		err := validateNetwork(address.(string))
		if err != nil {
			return nil, err
		}
		configSet = append(configSet, setPrefixTermFrom+"source-address "+address.(string)+"\n")
	}
	for _, address := range fromMap["source_address_except"].([]interface{}) {
		err := validateNetwork(address.(string))
		if err != nil {
			return nil, err
		}
		configSet = append(configSet, setPrefixTermFrom+"source-address "+address.(string)+" except\n")
	}
	if len(fromMap["source_port"].([]interface{})) > 0 && len(fromMap["source_port_except"].([]interface{})) > 0 {
		return configSet, fmt.Errorf("conflict between source_port and source_port_except")
	}
	for _, port := range fromMap["source_port"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"source-port "+port.(string)+"\n")
	}
	for _, port := range fromMap["source_port_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"source-port-except "+port.(string)+"\n")
	}
	for _, prefixList := range fromMap["source_prefix_list"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"source-prefix-list "+prefixList.(string)+"\n")
	}
	for _, prefixList := range fromMap["source_prefix_list_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"source-prefix-list "+prefixList.(string)+" except\n")
	}
	if len(fromMap["protocol"].([]interface{})) > 0 && len(fromMap["protocol_except"].([]interface{})) > 0 {
		return nil, fmt.Errorf("conflict between protocol and protocol_except")
	}
	for _, protocol := range fromMap["protocol"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"protocol "+protocol.(string)+"\n")
	}
	for _, protocol := range fromMap["protocol_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"protocol-except "+protocol.(string)+"\n")
	}
	if fromMap["tcp_flags"].(string) != "" && (fromMap["tcp_initial"].(bool) || fromMap["tcp_established"].(bool)) {
		return configSet, fmt.Errorf("conflict between tcp_flags and tcp_initial|tcp_established")
	}
	if fromMap["tcp_flags"].(string) != "" {
		configSet = append(configSet, setPrefixTermFrom+"tcp-flags \""+fromMap["tcp_flags"].(string)+"\"\n")
	}
	if fromMap["tcp_initial"].(bool) {
		configSet = append(configSet, setPrefixTermFrom+"tcp-initial\n")
	}
	if fromMap["tcp_established"].(bool) {
		configSet = append(configSet, setPrefixTermFrom+"tcp-established\n")
	}
	if len(fromMap["icmp_type"].([]interface{})) > 0 && len(fromMap["icmp_type_except"].([]interface{})) > 0 {
		return nil, fmt.Errorf("conflict between icmp_type and icmp_type_except")
	}
	for _, icmp := range fromMap["icmp_type"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"icmp-type "+icmp.(string)+"\n")
	}
	for _, icmp := range fromMap["icmp_type_except"].([]interface{}) {
		configSet = append(configSet, setPrefixTermFrom+"icmp-type-except "+icmp.(string)+"\n")
	}
	return configSet, nil
}
func setFirewallFilterOptsThen(setPrefixTermThen string, configSet []string, thenMap map[string]interface{}) []string {
	if thenMap["action"].(string) != "" {
		configSet = append(configSet, setPrefixTermThen+thenMap["action"].(string)+"\n")
	}
	if thenMap["count"].(string) != "" {
		configSet = append(configSet, setPrefixTermThen+"count "+thenMap["count"].(string)+"\n")
	}
	if thenMap["routing_instance"].(string) != "" {
		configSet = append(configSet, setPrefixTermThen+"routing-instance "+thenMap["routing_instance"].(string)+"\n")
	}
	if thenMap["policer"].(string) != "" {
		configSet = append(configSet, setPrefixTermThen+"policer "+thenMap["policer"].(string)+"\n")
	}
	if thenMap["log"].(bool) {
		configSet = append(configSet, setPrefixTermThen+"log\n")
	}
	if thenMap["syslog"].(bool) {
		configSet = append(configSet, setPrefixTermThen+"syslog\n")
	}
	if thenMap["port_mirror"].(bool) {
		configSet = append(configSet, setPrefixTermThen+"port-mirror\n")
	}
	if thenMap["sample"].(bool) {
		configSet = append(configSet, setPrefixTermThen+"sample\n")
	}
	if thenMap["service_accounting"].(bool) {
		configSet = append(configSet, setPrefixTermThen+"service-accounting\n")
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
	case strings.HasPrefix(item, "port "):
		fromMap["port"] = append(fromMap["port"].([]string),
			strings.TrimPrefix(item, "port "))
	case strings.HasPrefix(item, "port-except "):
		fromMap["port_except"] = append(fromMap["port_except"].([]string),
			strings.TrimPrefix(item, "port-except "))
	case strings.HasPrefix(item, "prefix-list "):
		if strings.HasSuffix(item, " except") {
			fromMap["prefix_list_except"] = append(fromMap["prefix_list_except"].([]string),
				strings.TrimSuffix(strings.TrimPrefix(item, "prefix-list "), " except"))
		} else {
			fromMap["prefix_list"] = append(fromMap["prefix_list"].([]string),
				strings.TrimPrefix(item, "prefix-list "))
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
	case strings.HasPrefix(item, "protocol "):
		fromMap["protocol"] = append(fromMap["protocol"].([]string), strings.TrimPrefix(item, "protocol "))
	case strings.HasPrefix(item, "protocol-except "):
		fromMap["protocol_except"] = append(fromMap["protocol_except"].([]string),
			strings.TrimPrefix(item, "protocol-except "))
	case strings.HasPrefix(item, "tcp-flags "):
		fromMap["tcp_flags"] = strings.Trim(strings.TrimPrefix(item, "tcp-flags "), "\"")
	case strings.HasSuffix(item, "tcp-initial"):
		fromMap["tcp_initial"] = true
	case strings.HasSuffix(item, "tcp-established"):
		fromMap["tcp_established"] = true
	case strings.HasPrefix(item, "icmp-type "):
		fromMap["icmp_type"] = append(fromMap["icmp_type"].([]string), strings.TrimPrefix(item, "icmp-type "))
	case strings.HasPrefix(item, "icmp-type-except "):
		fromMap["icmp_type_except"] = append(fromMap["icmp_type_except"].([]string),
			strings.TrimPrefix(item, "icmp-type-except "))
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
	case strings.HasSuffix(item, "accept"),
		strings.HasSuffix(item, "reject"),
		strings.HasSuffix(item, "discard"),
		strings.HasSuffix(item, "next term"):
		thenMap["action"] = item
	case strings.HasPrefix(item, "count "):
		thenMap["count"] = strings.TrimPrefix(item, "count ")
	case strings.HasPrefix(item, "routing-instance "):
		thenMap["routing_instance"] = strings.TrimPrefix(item, "routing-instance ")
	case strings.HasPrefix(item, "policer "):
		thenMap["policer"] = strings.TrimPrefix(item, "policer ")
	case strings.HasSuffix(item, "syslog"):
		thenMap["syslog"] = true
	case strings.HasSuffix(item, "log"):
		thenMap["log"] = true
	case strings.HasSuffix(item, "port-mirror"):
		thenMap["port_mirror"] = true
	case strings.HasSuffix(item, "sample"):
		thenMap["sample"] = true
	case strings.HasSuffix(item, "service-accounting"):
		thenMap["service_accounting"] = true
	}
	// override (maxItem = 1)
	return []map[string]interface{}{thenMap}
}

func genMapFirewallFilterOptsFrom() map[string]interface{} {
	return map[string]interface{}{
		"address":                        make([]string, 0),
		"address_except":                 make([]string, 0),
		"port":                           make([]string, 0),
		"port_except":                    make([]string, 0),
		"prefix_list":                    make([]string, 0),
		"prefix_list_except":             make([]string, 0),
		"destination_address":            make([]string, 0),
		"destination_address_except":     make([]string, 0),
		"destination_port":               make([]string, 0),
		"destination_port_except":        make([]string, 0),
		"destination_prefix_list":        make([]string, 0),
		"destination_prefix_list_except": make([]string, 0),
		"source_address":                 make([]string, 0),
		"source_address_except":          make([]string, 0),
		"source_port":                    make([]string, 0),
		"source_port_except":             make([]string, 0),
		"source_prefix_list":             make([]string, 0),
		"source_prefix_list_except":      make([]string, 0),
		"protocol":                       make([]string, 0),
		"protocol_except":                make([]string, 0),
		"tcp_flags":                      "",
		"tcp_initial":                    false,
		"tcp_established":                false,
		"icmp_type":                      make([]string, 0),
		"icmp_type_except":               make([]string, 0),
	}
}
func genMapFirewallFilterOptsThen() map[string]interface{} {
	return map[string]interface{}{
		"action":             "",
		"count":              "",
		"routing_instance":   "",
		"policer":            "",
		"log":                false,
		"syslog":             false,
		"port_mirror":        false,
		"sample":             false,
		"service_accounting": false,
	}
}
