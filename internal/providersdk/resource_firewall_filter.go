package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type filterOptions struct {
	interfaceSpecific bool
	name              string
	family            string
	term              []map[string]interface{}
}

func resourceFirewallFilter() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceFirewallFilterCreate,
		ReadWithoutTimeout:   resourceFirewallFilterRead,
		UpdateWithoutTimeout: resourceFirewallFilterUpdate,
		DeleteWithoutTimeout: resourceFirewallFilterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceFirewallFilterImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"family": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{junos.InetW, junos.Inet6W, "any", "ccc", "mpls", "vpls", "ethernet-switching"}, false),
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
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"filter": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"from": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validateCIDRNetworkFunc(),
										},
									},
									"address_except": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validateCIDRNetworkFunc(),
										},
									},
									"destination_address": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validateCIDRNetworkFunc(),
										},
									},
									"destination_address_except": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validateCIDRNetworkFunc(),
										},
									},
									"destination_port": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"destination_port_except": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"destination_prefix_list": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
										},
									},
									"destination_prefix_list_except": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
										},
									},
									"icmp_code": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"icmp_code_except": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"icmp_type": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"icmp_type_except": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"is_fragment": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"next_header": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"next_header_except": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"port": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"port_except": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"prefix_list": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
										},
									},
									"prefix_list_except": {
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
									"protocol_except": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"source_address": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validateCIDRNetworkFunc(),
										},
									},
									"source_address_except": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validateCIDRNetworkFunc(),
										},
									},
									"source_port": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"source_port_except": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"source_prefix_list": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
										},
									},
									"source_prefix_list_except": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
										},
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
									"packet_mode": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"policer": {
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setFirewallFilter(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string) + junos.IDSeparator + d.Get("family").(string))

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
	firewallFilterExists, err := checkFirewallFilterExists(
		d.Get("name").(string),
		d.Get("family").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if firewallFilterExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("firewall filter %v already exists", d.Get("name").(string)))...)
	}

	if err := setFirewallFilter(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_firewall_filter")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	firewallFilterExists, err = checkFirewallFilterExists(d.Get("name").(string), d.Get("family").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if firewallFilterExists {
		d.SetId(d.Get("name").(string) + junos.IDSeparator + d.Get("family").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("firewall filter %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceFirewallFilterReadWJunSess(d, junSess)...)
}

func resourceFirewallFilterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceFirewallFilterReadWJunSess(d, junSess)
}

func resourceFirewallFilterReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	filterOptions, err := readFirewallFilter(d.Get("name").(string), d.Get("family").(string), junSess)
	junos.MutexUnlock()
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
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delFirewallFilter(d.Get("name").(string), d.Get("family").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setFirewallFilter(d, junSess); err != nil {
			return diag.FromErr(err)
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
	if err := delFirewallFilter(d.Get("name").(string), d.Get("family").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setFirewallFilter(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_firewall_filter")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceFirewallFilterReadWJunSess(d, junSess)...)
}

func resourceFirewallFilterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delFirewallFilter(d.Get("name").(string), d.Get("family").(string), junSess); err != nil {
			return diag.FromErr(err)
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
	if err := delFirewallFilter(d.Get("name").(string), d.Get("family").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_firewall_filter")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceFirewallFilterImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	idList := strings.Split(d.Id(), junos.IDSeparator)
	if len(idList) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", junos.IDSeparator)
	}
	firewallFilterExists, err := checkFirewallFilterExists(idList[0], idList[1], junSess)
	if err != nil {
		return nil, err
	}
	if !firewallFilterExists {
		return nil,
			fmt.Errorf(
				"don't find firewall filter with id '%v' (id must be <name>"+junos.IDSeparator+"<family>)",
				d.Id(),
			)
	}
	filterOptions, err := readFirewallFilter(idList[0], idList[1], junSess)
	if err != nil {
		return nil, err
	}
	fillFirewallFilterData(d, filterOptions)

	result[0] = d

	return result, nil
}

func checkFirewallFilterExists(name, family string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"firewall family " + family + " filter " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setFirewallFilter(d *schema.ResourceData, junSess *junos.Session) (err error) {
	configSet := make([]string, 0)
	setPrefix := "set firewall family " + d.Get("family").(string) + " filter " + d.Get("name").(string)

	if d.Get("interface_specific").(bool) {
		configSet = append(configSet, setPrefix+" interface-specific")
	}
	termNameList := make([]string, 0)
	for _, v := range d.Get("term").([]interface{}) {
		term := v.(map[string]interface{})
		if bchk.InSlice(term["name"].(string), termNameList) {
			return fmt.Errorf("multiple blocks term with the same name %s", term["name"].(string))
		}
		termNameList = append(termNameList, term["name"].(string))
		setPrefixTerm := setPrefix + " term " + term["name"].(string)
		if term["filter"].(string) != "" {
			configSet = append(configSet, setPrefixTerm+" filter "+term["filter"].(string))
		}

		for _, from := range term["from"].([]interface{}) {
			configSet, err = setFirewallFilterOptsFrom(setPrefixTerm+" from ", configSet, from.(map[string]interface{}))
			if err != nil {
				return err
			}
		}
		for _, then := range term["then"].([]interface{}) {
			configSet = setFirewallFilterOptsThen(setPrefixTerm+" then ", configSet, then.(map[string]interface{}))
		}
	}

	return junSess.ConfigSet(configSet)
}

func readFirewallFilter(filter, family string, junSess *junos.Session,
) (confRead filterOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"firewall family " + family + " filter " + filter + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = filter
		confRead.family = family
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == "interface-specific":
				confRead.interfaceSpecific = true
			case balt.CutPrefixInString(&itemTrim, "term "):
				itemTrimFields := strings.Split(itemTrim, " ")
				termOptions := map[string]interface{}{
					"name":   itemTrimFields[0],
					"filter": "",
					"from":   make([]map[string]interface{}, 0),
					"then":   make([]map[string]interface{}, 0),
				}
				confRead.term = copyAndRemoveItemMapList("name", termOptions, confRead.term)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "filter "):
					termOptions["filter"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, "from "):
					if len(termOptions["from"].([]map[string]interface{})) == 0 {
						termOptions["from"] = append(termOptions["from"].([]map[string]interface{}),
							genMapFirewallFilterOptsFrom())
					}
					readFirewallFilterOptsFrom(itemTrim, termOptions["from"].([]map[string]interface{})[0])
				case balt.CutPrefixInString(&itemTrim, "then "):
					if len(termOptions["then"].([]map[string]interface{})) == 0 {
						termOptions["then"] = append(termOptions["then"].([]map[string]interface{}),
							genMapFirewallFilterOptsThen())
					}
					readFirewallFilterOptsThen(itemTrim, termOptions["then"].([]map[string]interface{})[0])
				}
				confRead.term = append(confRead.term, termOptions)
			}
		}
	}

	return confRead, nil
}

func delFirewallFilter(filter, family string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete firewall family "+family+" filter "+filter)

	return junSess.ConfigSet(configSet)
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

func setFirewallFilterOptsFrom(setPrefixTermFrom string, configSet []string, fromMap map[string]interface{},
) ([]string, error) {
	for _, address := range sortSetOfString(fromMap["address"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"address "+address)
	}
	for _, address := range sortSetOfString(fromMap["address_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"address "+address+" except")
	}
	for _, address := range sortSetOfString(fromMap["destination_address"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"destination-address "+address)
	}
	for _, address := range sortSetOfString(fromMap["destination_address_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"destination-address "+address+" except")
	}
	if len(fromMap["destination_port"].(*schema.Set).List()) > 0 &&
		len(fromMap["destination_port_except"].(*schema.Set).List()) > 0 {
		return configSet, fmt.Errorf("conflict between destination_port and destination_port_except")
	}
	for _, port := range sortSetOfString(fromMap["destination_port"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"destination-port "+port)
	}
	for _, port := range sortSetOfString(fromMap["destination_port_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"destination-port-except "+port)
	}
	for _, prefixList := range sortSetOfString(fromMap["destination_prefix_list"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"destination-prefix-list "+prefixList)
	}
	for _, prefixList := range sortSetOfString(fromMap["destination_prefix_list_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"destination-prefix-list "+prefixList+" except")
	}
	if len(fromMap["icmp_code"].(*schema.Set).List()) > 0 &&
		len(fromMap["icmp_code_except"].(*schema.Set).List()) > 0 {
		return nil, fmt.Errorf("conflict between icmp_code and icmp_code_except")
	}
	for _, icmp := range sortSetOfString(fromMap["icmp_code"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"icmp-code "+icmp)
	}
	for _, icmp := range sortSetOfString(fromMap["icmp_code_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"icmp-code-except "+icmp)
	}
	if len(fromMap["icmp_type"].(*schema.Set).List()) > 0 &&
		len(fromMap["icmp_type_except"].(*schema.Set).List()) > 0 {
		return nil, fmt.Errorf("conflict between icmp_type and icmp_type_except")
	}
	for _, icmp := range sortSetOfString(fromMap["icmp_type"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"icmp-type "+icmp)
	}
	for _, icmp := range sortSetOfString(fromMap["icmp_type_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"icmp-type-except "+icmp)
	}
	if fromMap["is_fragment"].(bool) {
		configSet = append(configSet, setPrefixTermFrom+"is-fragment")
	}
	if len(fromMap["next_header"].(*schema.Set).List()) > 0 &&
		len(fromMap["next_header_except"].(*schema.Set).List()) > 0 {
		return nil, fmt.Errorf("conflict between next_header and next_header_except")
	}
	for _, header := range sortSetOfString(fromMap["next_header"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"next-header "+header)
	}
	for _, header := range sortSetOfString(fromMap["next_header_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"next-header-except "+header)
	}
	if len(fromMap["port"].(*schema.Set).List()) > 0 &&
		len(fromMap["port_except"].(*schema.Set).List()) > 0 {
		return configSet, fmt.Errorf("conflict between port and port_except")
	}
	for _, port := range sortSetOfString(fromMap["port"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"port "+port)
	}
	for _, port := range sortSetOfString(fromMap["port_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"port-except "+port)
	}
	for _, prefixList := range sortSetOfString(fromMap["prefix_list"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"prefix-list "+prefixList)
	}
	for _, prefixList := range sortSetOfString(fromMap["prefix_list_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"prefix-list "+prefixList+" except")
	}
	if len(fromMap["protocol"].(*schema.Set).List()) > 0 &&
		len(fromMap["protocol_except"].(*schema.Set).List()) > 0 {
		return nil, fmt.Errorf("conflict between protocol and protocol_except")
	}
	for _, protocol := range sortSetOfString(fromMap["protocol"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"protocol "+protocol)
	}
	for _, protocol := range sortSetOfString(fromMap["protocol_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"protocol-except "+protocol)
	}
	for _, address := range sortSetOfString(fromMap["source_address"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"source-address "+address)
	}
	for _, address := range sortSetOfString(fromMap["source_address_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"source-address "+address+" except")
	}
	if len(fromMap["source_port"].(*schema.Set).List()) > 0 &&
		len(fromMap["source_port_except"].(*schema.Set).List()) > 0 {
		return configSet, fmt.Errorf("conflict between source_port and source_port_except")
	}
	for _, port := range sortSetOfString(fromMap["source_port"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"source-port "+port)
	}
	for _, port := range sortSetOfString(fromMap["source_port_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"source-port-except "+port)
	}
	for _, prefixList := range sortSetOfString(fromMap["source_prefix_list"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"source-prefix-list "+prefixList)
	}
	for _, prefixList := range sortSetOfString(fromMap["source_prefix_list_except"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefixTermFrom+"source-prefix-list "+prefixList+" except")
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
	if thenMap["packet_mode"].(bool) {
		configSet = append(configSet, setPrefixTermThen+"packet-mode")
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

func readFirewallFilterOptsFrom(itemTrim string, fromMap map[string]interface{}) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "address "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			fromMap["address_except"] = append(fromMap["address_except"].([]string), itemTrim)
		} else {
			fromMap["address"] = append(fromMap["address"].([]string), itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "destination-address "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			fromMap["destination_address_except"] = append(fromMap["destination_address_except"].([]string), itemTrim)
		} else {
			fromMap["destination_address"] = append(fromMap["destination_address"].([]string), itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "destination-port "):
		fromMap["destination_port"] = append(fromMap["destination_port"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination-port-except "):
		fromMap["destination_port_except"] = append(fromMap["destination_port_except"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination-prefix-list "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			fromMap["destination_prefix_list_except"] = append(
				fromMap["destination_prefix_list_except"].([]string),
				itemTrim,
			)
		} else {
			fromMap["destination_prefix_list"] = append(fromMap["destination_prefix_list"].([]string), itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "icmp-code "):
		fromMap["icmp_code"] = append(fromMap["icmp_code"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "icmp-code-except "):
		fromMap["icmp_code_except"] = append(fromMap["icmp_code_except"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "icmp-type "):
		fromMap["icmp_type"] = append(fromMap["icmp_type"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "icmp-type-except "):
		fromMap["icmp_type_except"] = append(fromMap["icmp_type_except"].([]string), itemTrim)
	case itemTrim == "is-fragment":
		fromMap["is_fragment"] = true
	case balt.CutPrefixInString(&itemTrim, "next-header "):
		fromMap["next_header"] = append(fromMap["next_header"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "next-header-except "):
		fromMap["next_header_except"] = append(fromMap["next_header_except"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "port "):
		fromMap["port"] = append(fromMap["port"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "port-except "):
		fromMap["port_except"] = append(fromMap["port_except"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "protocol "):
		fromMap["protocol"] = append(fromMap["protocol"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "protocol-except "):
		fromMap["protocol_except"] = append(fromMap["protocol_except"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "prefix-list "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			fromMap["prefix_list_except"] = append(fromMap["prefix_list_except"].([]string), itemTrim)
		} else {
			fromMap["prefix_list"] = append(fromMap["prefix_list"].([]string), itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			fromMap["source_address_except"] = append(
				fromMap["source_address_except"].([]string),
				itemTrim,
			)
		} else {
			fromMap["source_address"] = append(fromMap["source_address"].([]string), itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "source-port "):
		fromMap["source_port"] = append(fromMap["source_port"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-port-except "):
		fromMap["source_port_except"] = append(fromMap["source_port_except"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-prefix-list "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			fromMap["source_prefix_list_except"] = append(
				fromMap["source_prefix_list_except"].([]string),
				itemTrim,
			)
		} else {
			fromMap["source_prefix_list"] = append(fromMap["source_prefix_list"].([]string), itemTrim)
		}
	case itemTrim == "tcp-established":
		fromMap["tcp_established"] = true
	case balt.CutPrefixInString(&itemTrim, "tcp-flags "):
		fromMap["tcp_flags"] = strings.Trim(itemTrim, "\"")
	case itemTrim == "tcp-initial":
		fromMap["tcp_initial"] = true
	}
}

func readFirewallFilterOptsThen(itemTrim string, thenMap map[string]interface{}) {
	switch {
	case itemTrim == "accept",
		itemTrim == "reject",
		itemTrim == junos.DiscardW,
		itemTrim == "next term":
		thenMap["action"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "count "):
		thenMap["count"] = itemTrim
	case itemTrim == "log":
		thenMap["log"] = true
	case itemTrim == "packet-mode":
		thenMap["packet_mode"] = true
	case balt.CutPrefixInString(&itemTrim, "policer "):
		thenMap["policer"] = itemTrim
	case itemTrim == "port-mirror":
		thenMap["port_mirror"] = true
	case balt.CutPrefixInString(&itemTrim, "routing-instance "):
		thenMap["routing_instance"] = itemTrim
	case itemTrim == "sample":
		thenMap["sample"] = true
	case itemTrim == "service-accounting":
		thenMap["service_accounting"] = true
	case itemTrim == "syslog":
		thenMap["syslog"] = true
	}
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
		"packet_mode":        false,
		"policer":            "",
		"port_mirror":        false,
		"routing_instance":   "",
		"sample":             false,
		"service_accounting": false,
		"syslog":             false,
	}
}
