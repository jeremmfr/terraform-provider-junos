package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	addressBookSet                   []map[string]interface{}
}

func resourceSecurityZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityZoneCreate,
		ReadContext:   resourceSecurityZoneRead,
		UpdateContext: resourceSecurityZoneUpdate,
		DeleteContext: resourceSecurityZoneDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityZoneImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"address_book": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
						},
						"network": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
						},
					},
				},
			},
			"address_book_set": {
				Type:     schema.TypeList,
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
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"inbound_services": {
				Type:     schema.TypeList,
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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security zone not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	securityZoneExists, err := checkSecurityZonesExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if securityZoneExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security zone %v already exists", d.Get("name").(string)))
	}

	if err := setSecurityZone(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_zone", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityZoneExists, err = checkSecurityZonesExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityZoneExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security zone %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityZoneReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityZoneReadWJnprSess(d, m, jnprSess)
}
func resourceSecurityZoneReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	zoneOptions, err := readSecurityZone(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	err = delSecurityZoneOpts(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setSecurityZone(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security_zone", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityZoneReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSecurityZone(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_zone", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceSecurityZoneImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	securityZoneExists, err := checkSecurityZonesExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityZoneExists {
		return nil, fmt.Errorf("don't find zone with id '%v' (id must be <name>)", d.Id())
	}
	zoneOptions, err := readSecurityZone(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityZoneData(d, zoneOptions)

	result[0] = d

	return result, nil
}

func checkSecurityZonesExists(zone string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	zoneConfig, err := sess.command("show configuration security zones security-zone "+zone+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if zoneConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setSecurityZone(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security zones security-zone " + d.Get("name").(string)
	configSet = append(configSet, setPrefix)
	for _, v := range d.Get("address_book").([]interface{}) {
		addressBook := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+" address-book address "+
			addressBook["name"].(string)+" "+addressBook["network"].(string))
	}
	for _, v := range d.Get("address_book_set").([]interface{}) {
		addressBookSet := v.(map[string]interface{})
		for _, addressBookSetAddress := range addressBookSet["address"].([]interface{}) {
			configSet = append(configSet, setPrefix+" address-book address-set "+addressBookSet["name"].(string)+
				" address "+addressBookSetAddress.(string))
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
	for _, v := range d.Get("inbound_protocols").([]interface{}) {
		configSet = append(configSet, setPrefix+" host-inbound-traffic protocols "+v.(string))
	}
	for _, v := range d.Get("inbound_services").([]interface{}) {
		configSet = append(configSet, setPrefix+" host-inbound-traffic system-services "+v.(string))
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

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readSecurityZone(zone string, m interface{}, jnprSess *NetconfObject) (zoneOptions, error) {
	sess := m.(*Session)
	var confRead zoneOptions

	zoneConfig, err := sess.command("show configuration"+
		" security zones security-zone "+zone+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if zoneConfig != emptyWord {
		confRead.name = zone
		for _, item := range strings.Split(zoneConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "address-book address "):
				address := strings.TrimPrefix(itemTrim, "address-book address ")
				addressWords := strings.Split(address, " ")
				// addressWords[0] = name of address
				// addressWords[1] = network
				confRead.addressBook = append(confRead.addressBook, map[string]interface{}{
					"name":    addressWords[0],
					"network": addressWords[1],
				})
			case strings.HasPrefix(itemTrim, "address-book address-set "):
				address := strings.TrimPrefix(itemTrim, "address-book address-set ")
				addressWords := strings.Split(address, " ")
				// addressWords[0] = name of address-set
				// addressWords[1] = "address"
				// addressWords[2] = name of address
				m := map[string]interface{}{
					"name":    addressWords[0],
					"address": make([]string, 0),
				}
				// search if name of address-set already create
				m, confRead.addressBookSet = copyAndRemoveItemMapList("name", false, m, confRead.addressBookSet)
				// append new address find
				m["address"] = append(m["address"].([]string), addressWords[2])
				confRead.addressBookSet = append(confRead.addressBookSet, m)
			case strings.HasPrefix(itemTrim, "advance-policy-based-routing-profile "):
				confRead.advancePolicyBasedRoutingProfile = strings.Trim(strings.TrimPrefix(itemTrim,
					"advance-policy-based-routing-profile "), "\"")
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim,
					"description "), "\"")
			case itemTrim == "application-tracking":
				confRead.appTrack = true
			case strings.HasPrefix(itemTrim, "host-inbound-traffic protocols "):
				confRead.inboundProtocols = append(confRead.inboundProtocols, strings.TrimPrefix(itemTrim,
					"host-inbound-traffic protocols "))
			case strings.HasPrefix(itemTrim, "host-inbound-traffic system-services "):
				confRead.inboundServices = append(confRead.inboundServices, strings.TrimPrefix(itemTrim,
					"host-inbound-traffic system-services "))
			case itemTrim == "enable-reverse-reroute":
				confRead.reverseReroute = true
			case strings.HasPrefix(itemTrim, "screen "):
				confRead.screen = strings.Trim(strings.TrimPrefix(itemTrim, "screen "), "\"")
			case itemTrim == "source-identity-log":
				confRead.sourceIdentityLog = true
			case itemTrim == "tcp-rst":
				confRead.tcpRst = true
			}
		}
	}

	return confRead, nil
}
func delSecurityZoneOpts(zone string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	listLinesToDelete := []string{
		"address-book",
		"advance-policy-based-routing-profile",
		"description",
		"application-tracking",
		"host-inbound-traffic",
		"enable-reverse-reroute",
		"screen",
		"source-identity-log",
		"tcp-rst",
	}
	configSet := make([]string, 0, 1)
	delPrefix := "delete security zones security-zone " + zone + " "
	for _, line := range listLinesToDelete {
		configSet = append(configSet, delPrefix+line)
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func delSecurityZone(zone string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security zones security-zone "+zone)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func fillSecurityZoneData(d *schema.ResourceData, zoneOptions zoneOptions) {
	if tfErr := d.Set("name", zoneOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_book", zoneOptions.addressBook); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_book_set", zoneOptions.addressBookSet); tfErr != nil {
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
