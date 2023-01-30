package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type zoneOptions struct {
	appTrack                         bool
	reverseReroute                   bool
	sourceIdentityLog                bool
	tcpRst                           bool
	advancePolicyBasedRoutingProfile string
	description                      string
	name                             string
	screen                           string
	inboundProtocols                 []string
	inboundServices                  []string
	addressBook                      []map[string]interface{}
	addressBookDNS                   []map[string]interface{}
	addressBookRange                 []map[string]interface{}
	addressBookSet                   []map[string]interface{}
	addressBookWildcard              []map[string]interface{}
	interFace                        []map[string]interface{} // to data_source
}

func resourceSecurityZone() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityZoneCreate,
		ReadWithoutTimeout:   resourceSecurityZoneRead,
		UpdateWithoutTimeout: resourceSecurityZoneUpdate,
		DeleteWithoutTimeout: resourceSecurityZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityZoneImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"address_book": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"network": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			"address_book_configure_singly": {
				Type:     schema.TypeBool,
				Optional: true,
				ConflictsWith: []string{
					"address_book",
					"address_book_dns",
					"address_book_range",
					"address_book_set",
					"address_book_wildcard",
				},
			},
			"address_book_dns": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"fqdn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"ipv4_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"ipv6_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"address_book_range": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"from": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"to": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			"address_book_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"address": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
							},
						},
						"address_set": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
							},
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			"address_book_wildcard": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"network": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateWildcardFunc(),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			"advance_policy_based_routing_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"application_tracking": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inbound_protocols": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"inbound_services": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"reverse_reroute": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"screen": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_identity_log": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tcp_rst": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceSecurityZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityZone(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilitySecurity() {
		return diag.FromErr(fmt.Errorf("security zone not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityZoneExists, err := checkSecurityZonesExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityZoneExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("security zone %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityZone(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_zone")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityZoneExists, err = checkSecurityZonesExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityZoneExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security zone %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityZoneReadWJunSess(d, junSess)...)
}

func resourceSecurityZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityZoneReadWJunSess(d, junSess)
}

func resourceSecurityZoneReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	zoneOptions, err := readSecurityZone(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if zoneOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityZoneData(d, zoneOptions)
	}

	return nil
}

func resourceSecurityZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	var diagWarns diag.Diagnostics
	addressBookConfiguredSingly := d.Get("address_book_configure_singly").(bool)
	if d.HasChange("address_book_configure_singly") {
		if o, _ := d.GetChange("address_book_configure_singly"); o.(bool) {
			addressBookConfiguredSingly = o.(bool)
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary: "Disable address_book_configure_singly on resource already created doesn't " +
					"delete addresses and address-sets already configured.",
				Detail:        "So refresh resource after apply to detect address-book entries that need to be deleted",
				AttributePath: cty.Path{cty.GetAttrStep{Name: "address_book_configure_singly"}},
			})
		} else {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary: "Enable address_book_configure_singly on resource already created doesn't " +
					"delete addresses and address-sets already configured.",
				Detail:        "So import address-book entries in dedicated resource(s) to be able to manage them",
				AttributePath: cty.Path{cty.GetAttrStep{Name: "address_book_configure_singly"}},
			})
		}
	}
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityZoneOpts(
			d.Get("name").(string),
			addressBookConfiguredSingly,
			junSess,
		); err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		if err := setSecurityZone(d, junSess); err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		d.Partial(false)

		return diagWarns
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	if err := delSecurityZoneOpts(
		d.Get("name").(string),
		addressBookConfiguredSingly,
		junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityZone(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_zone")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityZoneReadWJunSess(d, junSess)...)
}

func resourceSecurityZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityZone(d.Get("name").(string), junSess); err != nil {
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
	if err := delSecurityZone(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_zone")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityZoneImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	securityZoneExists, err := checkSecurityZonesExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !securityZoneExists {
		return nil, fmt.Errorf("don't find zone with id '%v' (id must be <name>)", d.Id())
	}
	zoneOptions, err := readSecurityZone(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityZoneData(d, zoneOptions)

	result[0] = d

	return result, nil
}

func checkSecurityZonesExists(zone string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "security zones security-zone " + zone + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityZone(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security zones security-zone " + d.Get("name").(string)
	configSet = append(configSet, setPrefix)
	if !d.Get("address_book_configure_singly").(bool) {
		addressNameList := make([]string, 0)
		for _, v := range d.Get("address_book").(*schema.Set).List() {
			addressBook := v.(map[string]interface{})
			if bchk.InSlice(addressBook["name"].(string), addressNameList) {
				return fmt.Errorf("multiple addresses with the same name %s", addressBook["name"].(string))
			}
			addressNameList = append(addressNameList, addressBook["name"].(string))
			configSet = append(configSet, setPrefix+" address-book address "+
				addressBook["name"].(string)+" "+addressBook["network"].(string))
			if v2 := addressBook["description"].(string); v2 != "" {
				configSet = append(configSet, setPrefix+" address-book address "+
					addressBook["name"].(string)+" description \""+v2+"\"")
			}
		}
		for _, v := range d.Get("address_book_dns").(*schema.Set).List() {
			addressBook := v.(map[string]interface{})
			if bchk.InSlice(addressBook["name"].(string), addressNameList) {
				return fmt.Errorf("multiple addresses with the same name %s", addressBook["name"].(string))
			}
			addressNameList = append(addressNameList, addressBook["name"].(string))
			setLine := setPrefix + " address-book address " + addressBook["name"].(string) +
				" dns-name " + addressBook["fqdn"].(string)
			configSet = append(configSet, setLine)
			if addressBook["ipv4_only"].(bool) {
				configSet = append(configSet, setLine+" ipv4-only")
			}
			if addressBook["ipv6_only"].(bool) {
				configSet = append(configSet, setLine+" ipv6-only")
			}
			if v2 := addressBook["description"].(string); v2 != "" {
				configSet = append(configSet, setPrefix+" address-book address "+
					addressBook["name"].(string)+" description \""+v2+"\"")
			}
		}
		for _, v := range d.Get("address_book_range").(*schema.Set).List() {
			addressBook := v.(map[string]interface{})
			if bchk.InSlice(addressBook["name"].(string), addressNameList) {
				return fmt.Errorf("multiple addresses with the same name %s", addressBook["name"].(string))
			}
			addressNameList = append(addressNameList, addressBook["name"].(string))
			configSet = append(configSet, setPrefix+" address-book address "+
				addressBook["name"].(string)+" range-address "+addressBook["from"].(string)+
				" to "+addressBook["to"].(string))
			if v2 := addressBook["description"].(string); v2 != "" {
				configSet = append(configSet, setPrefix+" address-book address "+
					addressBook["name"].(string)+" description \""+v2+"\"")
			}
		}
		for _, v := range d.Get("address_book_wildcard").(*schema.Set).List() {
			addressBook := v.(map[string]interface{})
			if bchk.InSlice(addressBook["name"].(string), addressNameList) {
				return fmt.Errorf("multiple addresses with the same name %s", addressBook["name"].(string))
			}
			addressNameList = append(addressNameList, addressBook["name"].(string))
			configSet = append(configSet, setPrefix+" address-book address "+
				addressBook["name"].(string)+" wildcard-address "+addressBook["network"].(string))
			if v2 := addressBook["description"].(string); v2 != "" {
				configSet = append(configSet, setPrefix+" address-book address "+
					addressBook["name"].(string)+" description \""+v2+"\"")
			}
		}
		for _, v := range d.Get("address_book_set").(*schema.Set).List() {
			addressBookSet := v.(map[string]interface{})
			if bchk.InSlice(addressBookSet["name"].(string), addressNameList) {
				return fmt.Errorf("multiple addresses or address-sets with the same name %s", addressBookSet["name"].(string))
			}
			addressNameList = append(addressNameList, addressBookSet["name"].(string))
			if len(addressBookSet["address"].(*schema.Set).List()) == 0 &&
				len(addressBookSet["address_set"].(*schema.Set).List()) == 0 {
				return fmt.Errorf("at least one of address or address_set is required "+
					"in address_book_set %s", addressBookSet["name"].(string))
			}
			for _, addressBookSetAddress := range sortSetOfString(addressBookSet["address"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefix+" address-book address-set "+addressBookSet["name"].(string)+
					" address "+addressBookSetAddress)
			}
			for _, addressBookSetAddressSet := range sortSetOfString(addressBookSet["address_set"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefix+" address-book address-set "+addressBookSet["name"].(string)+
					" address-set "+addressBookSetAddressSet)
			}
			if v2 := addressBookSet["description"].(string); v2 != "" {
				configSet = append(configSet, setPrefix+" address-book address-set "+
					addressBookSet["name"].(string)+" description \""+v2+"\"")
			}
		}
	}
	if d.Get("advance_policy_based_routing_profile").(string) != "" {
		configSet = append(configSet, setPrefix+" advance-policy-based-routing-profile \""+
			d.Get("advance_policy_based_routing_profile").(string)+"\"")
	}
	if d.Get("application_tracking").(bool) {
		configSet = append(configSet, setPrefix+" application-tracking")
	}
	if d.Get("description").(string) != "" {
		configSet = append(configSet, setPrefix+" description \""+d.Get("description").(string)+"\"")
	}
	for _, v := range sortSetOfString(d.Get("inbound_protocols").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+" host-inbound-traffic protocols "+v)
	}
	for _, v := range sortSetOfString(d.Get("inbound_services").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+" host-inbound-traffic system-services "+v)
	}
	if d.Get("reverse_reroute").(bool) {
		configSet = append(configSet, setPrefix+" enable-reverse-reroute")
	}
	if d.Get("screen").(string) != "" {
		configSet = append(configSet, setPrefix+" screen \""+d.Get("screen").(string)+"\"")
	}
	if d.Get("source_identity_log").(bool) {
		configSet = append(configSet, setPrefix+" source-identity-log")
	}
	if d.Get("tcp_rst").(bool) {
		configSet = append(configSet, setPrefix+" tcp-rst")
	}

	return junSess.ConfigSet(configSet)
}

func readSecurityZone(zone string, junSess *junos.Session,
) (confRead zoneOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security zones security-zone " + zone + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	descAddressBookMap := make(map[string]string)
	if showConfig != junos.EmptyW {
		confRead.name = zone
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "address-book address "):
				itemTrimFields := strings.Split(itemTrim, " ")
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "description "):
					descAddressBookMap[itemTrimFields[0]] = strings.Trim(itemTrim, "\"")
				case balt.CutPrefixInString(&itemTrim, "dns-name "):
					switch {
					case balt.CutSuffixInString(&itemTrim, " ipv4-only"):
						confRead.addressBookDNS = append(confRead.addressBookDNS, map[string]interface{}{
							"name":        itemTrimFields[0],
							"description": descAddressBookMap[itemTrimFields[0]],
							"fqdn":        itemTrim,
							"ipv4_only":   true,
							"ipv6_only":   false,
						})
					case balt.CutSuffixInString(&itemTrim, " ipv6-only"):
						confRead.addressBookDNS = append(confRead.addressBookDNS, map[string]interface{}{
							"name":        itemTrimFields[0],
							"description": descAddressBookMap[itemTrimFields[0]],
							"fqdn":        itemTrim,
							"ipv4_only":   false,
							"ipv6_only":   true,
						})
					default:
						confRead.addressBookDNS = append(confRead.addressBookDNS, map[string]interface{}{
							"name":        itemTrimFields[0],
							"description": descAddressBookMap[itemTrimFields[0]],
							"fqdn":        itemTrim,
							"ipv4_only":   false,
							"ipv6_only":   false,
						})
					}
				case balt.CutPrefixInString(&itemTrim, "range-address "):
					rangeAddressFields := strings.Split(itemTrim, " ")
					if len(rangeAddressFields) < 3 { // <from> to <to>
						return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "range-address", itemTrim)
					}
					confRead.addressBookRange = append(confRead.addressBookRange, map[string]interface{}{
						"name":        itemTrimFields[0],
						"description": descAddressBookMap[itemTrimFields[0]],
						"from":        rangeAddressFields[0],
						"to":          rangeAddressFields[2],
					})
				case balt.CutPrefixInString(&itemTrim, "wildcard-address "):
					confRead.addressBookWildcard = append(confRead.addressBookWildcard, map[string]interface{}{
						"name":        itemTrimFields[0],
						"description": descAddressBookMap[itemTrimFields[0]],
						"network":     itemTrim,
					})
				default:
					confRead.addressBook = append(confRead.addressBook, map[string]interface{}{
						"name":        itemTrimFields[0],
						"description": descAddressBookMap[itemTrimFields[0]],
						"network":     itemTrim,
					})
				}
			case balt.CutPrefixInString(&itemTrim, "address-book address-set "):
				itemTrimFields := strings.Split(itemTrim, " ")
				adSet := map[string]interface{}{
					"name":        itemTrimFields[0],
					"address":     make([]string, 0),
					"address_set": make([]string, 0),
					"description": "",
				}
				confRead.addressBookSet = copyAndRemoveItemMapList("name", adSet, confRead.addressBookSet)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "description "):
					adSet["description"] = strings.Trim(itemTrim, "\"")
				case balt.CutPrefixInString(&itemTrim, "address "):
					adSet["address"] = append(adSet["address"].([]string), itemTrim)
				case balt.CutPrefixInString(&itemTrim, "address-set "):
					adSet["address_set"] = append(adSet["address_set"].([]string), itemTrim)
				}
				confRead.addressBookSet = append(confRead.addressBookSet, adSet)
			case balt.CutPrefixInString(&itemTrim, "advance-policy-based-routing-profile "):
				confRead.advancePolicyBasedRoutingProfile = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			case itemTrim == "application-tracking":
				confRead.appTrack = true
			case balt.CutPrefixInString(&itemTrim, "host-inbound-traffic protocols "):
				confRead.inboundProtocols = append(confRead.inboundProtocols, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "host-inbound-traffic system-services "):
				confRead.inboundServices = append(confRead.inboundServices, itemTrim)
			case itemTrim == "enable-reverse-reroute":
				confRead.reverseReroute = true
			case balt.CutPrefixInString(&itemTrim, "screen "):
				confRead.screen = strings.Trim(itemTrim, "\"")
			case itemTrim == "source-identity-log":
				confRead.sourceIdentityLog = true
			case itemTrim == "tcp-rst":
				confRead.tcpRst = true
			case balt.CutPrefixInString(&itemTrim, "interfaces "):
				itemTrimFields := strings.Split(itemTrim, " ")
				interFace := map[string]interface{}{
					"name":              itemTrimFields[0],
					"inbound_protocols": make([]string, 0),
					"inbound_services":  make([]string, 0),
				}
				confRead.interFace = copyAndRemoveItemMapList("name", interFace, confRead.interFace)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "host-inbound-traffic protocols "):
					interFace["inbound_protocols"] = append(interFace["inbound_protocols"].([]string), itemTrim)
				case balt.CutPrefixInString(&itemTrim, "host-inbound-traffic system-services "):
					interFace["inbound_services"] = append(interFace["inbound_services"].([]string), itemTrim)
				}
				confRead.interFace = append(confRead.interFace, interFace)
			}
		}
	}
	copySecurityAddressBookAddressDescriptions(descAddressBookMap, confRead.addressBook)
	copySecurityAddressBookAddressDescriptions(descAddressBookMap, confRead.addressBookDNS)
	copySecurityAddressBookAddressDescriptions(descAddressBookMap, confRead.addressBookRange)
	copySecurityAddressBookAddressDescriptions(descAddressBookMap, confRead.addressBookWildcard)

	return confRead, nil
}

func delSecurityZoneOpts(zone string, addressBookSingly bool, junSess *junos.Session) error {
	listLinesToDelete := []string{
		"advance-policy-based-routing-profile",
		"description",
		"application-tracking",
		"host-inbound-traffic",
		"enable-reverse-reroute",
		"screen",
		"source-identity-log",
		"tcp-rst",
	}
	if !addressBookSingly {
		listLinesToDelete = append(listLinesToDelete, "address-book")
	}
	delPrefix := "delete security zones security-zone " + zone + " "
	configSet := make([]string, len(listLinesToDelete))
	for k, line := range listLinesToDelete {
		configSet[k] = delPrefix + line
	}

	return junSess.ConfigSet(configSet)
}

func delSecurityZone(zone string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security zones security-zone "+zone)

	return junSess.ConfigSet(configSet)
}

func fillSecurityZoneData(d *schema.ResourceData, zoneOptions zoneOptions) {
	if tfErr := d.Set("name", zoneOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if !d.Get("address_book_configure_singly").(bool) {
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
