package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type zoneBookAddressSetOptions struct {
	description string
	name        string
	zone        string
	address     []string
	addressSet  []string
}

func resourceSecurityZoneBookAddressSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityZoneBookAddressSetCreate,
		ReadContext:   resourceSecurityZoneBookAddressSetRead,
		UpdateContext: resourceSecurityZoneBookAddressSetUpdate,
		DeleteContext: resourceSecurityZoneBookAddressSetDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityZoneBookAddressSetImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
			},
			"zone": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"address": {
				Type:         schema.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"address", "address_set"},
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"address_set": {
				Type:         schema.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"address", "address_set"},
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSecurityZoneBookAddressSetCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSecurityZoneBookAddressSet(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("zone").(string) + idSeparator + d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security zone address-book address-set not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	zonesExists, err := checkSecurityZonesExists(d.Get("zone").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !zonesExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security zone %v doesn't exist", d.Get("zone").(string)))...)
	}
	securityZoneBookAddressSetExists, err := checkSecurityZoneBookAddressSetsExists(
		d.Get("zone").(string), d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityZoneBookAddressSetExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security zone address-book address-set %v already exists in zone %s",
			d.Get("name").(string), d.Get("zone").(string)))...)
	}

	if err := setSecurityZoneBookAddressSet(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_zone_book_address_set", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityZoneBookAddressSetExists, err = checkSecurityZoneBookAddressSetsExists(
		d.Get("zone").(string), d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityZoneBookAddressSetExists {
		d.SetId(d.Get("zone").(string) + idSeparator + d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security zone address-book address-set %v not exists in zone %s after commit "+
				"=> check your config", d.Get("name").(string), d.Get("zone").(string)))...)
	}

	return append(diagWarns, resourceSecurityZoneBookAddressSetReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityZoneBookAddressSetRead(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityZoneBookAddressSetReadWJnprSess(d, m, jnprSess)
}

func resourceSecurityZoneBookAddressSetReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	zoneBookAddressSetOptions, err := readSecurityZoneBookAddressSet(
		d.Get("zone").(string), d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if zoneBookAddressSetOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityZoneBookAddressSetData(d, zoneBookAddressSetOptions)
	}

	return nil
}

func resourceSecurityZoneBookAddressSetUpdate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delSecurityZoneBookAddressSet(d.Get("zone").(string), d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityZoneBookAddressSet(d, m, nil); err != nil {
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
	if err := delSecurityZoneBookAddressSet(d.Get("zone").(string), d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityZoneBookAddressSet(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_zone_book_address_set", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityZoneBookAddressSetReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityZoneBookAddressSetDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delSecurityZoneBookAddressSet(d.Get("zone").(string), d.Get("name").(string), m, nil); err != nil {
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
	if err := delSecurityZoneBookAddressSet(d.Get("zone").(string), d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_zone_book_address_set", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityZoneBookAddressSetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
	securityZoneBookAddressSetExists, err := checkSecurityZoneBookAddressSetsExists(idList[0], idList[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityZoneBookAddressSetExists {
		return nil, fmt.Errorf(
			"don't find zone address-book address-set with id '%v' (id must be <zone>"+idSeparator+"<name>)", d.Id())
	}
	zoneBookAddressSetOptions, err := readSecurityZoneBookAddressSet(idList[0], idList[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityZoneBookAddressSetData(d, zoneBookAddressSetOptions)

	result[0] = d

	return result, nil
}

func checkSecurityZoneBookAddressSetsExists(
	zone, addressSet string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command(cmdShowConfig+
		"security zones security-zone "+zone+" address-book address-set "+addressSet+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityZoneBookAddressSet(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security zones security-zone " +
		d.Get("zone").(string) + " address-book address-set " + d.Get("name").(string) + " "
	if len(d.Get("address").(*schema.Set).List()) == 0 &&
		len(d.Get("address_set").(*schema.Set).List()) == 0 {
		return fmt.Errorf("at least one element of address or address_set is required")
	}
	for _, v := range sortSetOfString(d.Get("address").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"address "+v)
	}
	for _, v := range sortSetOfString(d.Get("address_set").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"address-set "+v)
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}

	return sess.configSet(configSet, jnprSess)
}

func readSecurityZoneBookAddressSet(
	zone, addressSet string, m interface{}, jnprSess *NetconfObject) (zoneBookAddressSetOptions, error) {
	sess := m.(*Session)
	var confRead zoneBookAddressSetOptions

	showConfig, err := sess.command(cmdShowConfig+
		"security zones security-zone "+zone+" address-book address-set "+addressSet+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = addressSet
		confRead.zone = zone
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "address "):
				confRead.address = append(confRead.address, strings.TrimPrefix(itemTrim, "address "))
			case strings.HasPrefix(itemTrim, "address-set "):
				confRead.addressSet = append(confRead.addressSet, strings.TrimPrefix(itemTrim, "address-set "))
			}
		}
	}

	return confRead, nil
}

func delSecurityZoneBookAddressSet(zone, addressSet string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete security zones security-zone " + zone + " address-book address-set " + addressSet}

	return sess.configSet(configSet, jnprSess)
}

func fillSecurityZoneBookAddressSetData(d *schema.ResourceData, zoneBookAddressSetOptions zoneBookAddressSetOptions) {
	if tfErr := d.Set("name", zoneBookAddressSetOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("zone", zoneBookAddressSetOptions.zone); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address", zoneBookAddressSetOptions.address); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_set", zoneBookAddressSetOptions.addressSet); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", zoneBookAddressSetOptions.description); tfErr != nil {
		panic(tfErr)
	}
}
