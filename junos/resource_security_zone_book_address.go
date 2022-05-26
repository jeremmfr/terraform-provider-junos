package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type zoneBookAddressOptions struct {
	dnsIPv4Only bool
	dnsIPv6Only bool
	cidr        string
	description string
	dnsName     string
	name        string
	rangeFrom   string
	rangeTo     string
	wildcard    string
	zone        string
}

func resourceSecurityZoneBookAddress() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityZoneBookAddressCreate,
		ReadWithoutTimeout:   resourceSecurityZoneBookAddressRead,
		UpdateWithoutTimeout: resourceSecurityZoneBookAddressUpdate,
		DeleteWithoutTimeout: resourceSecurityZoneBookAddressDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityZoneBookAddressImport,
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
			"cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 128),
				ExactlyOneOf: []string{"cidr", "dns_name", "range_from", "wildcard"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_ipv4_only": {
				Type:          schema.TypeBool,
				Optional:      true,
				RequiredWith:  []string{"dns_name"},
				ConflictsWith: []string{"dns_ipv6_only"},
			},
			"dns_ipv6_only": {
				Type:          schema.TypeBool,
				Optional:      true,
				RequiredWith:  []string{"dns_name"},
				ConflictsWith: []string{"dns_ipv4_only"},
			},
			"dns_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"cidr", "dns_name", "range_from", "wildcard"},
			},
			"range_from": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				RequiredWith: []string{"range_to"},
				ExactlyOneOf: []string{"cidr", "dns_name", "range_from", "wildcard"},
			},
			"range_to": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				RequiredWith: []string{"range_from"},
			},
			"wildcard": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateWildcardFunc(),
				ExactlyOneOf:     []string{"cidr", "dns_name", "range_from", "wildcard"},
			},
		},
	}
}

func resourceSecurityZoneBookAddressCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSecurityZoneBookAddress(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("zone").(string) + idSeparator + d.Get("name").(string))

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if !checkCompatibilitySecurity(junSess) {
		return diag.FromErr(fmt.Errorf("security zone address-book address not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	zonesExists, err := checkSecurityZonesExists(d.Get("zone").(string), sess, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !zonesExists {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security zone %v doesn't exist", d.Get("zone").(string)))...)
	}
	securityZoneBookAddressExists, err := checkSecurityZoneBookAddresssExists(
		d.Get("zone").(string),
		d.Get("name").(string),
		sess, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityZoneBookAddressExists {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security zone address-book address %v already exists in zone %s",
			d.Get("name").(string), d.Get("zone").(string)))...)
	}

	if err := setSecurityZoneBookAddress(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_zone_book_address", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityZoneBookAddressExists, err = checkSecurityZoneBookAddresssExists(
		d.Get("zone").(string),
		d.Get("name").(string),
		sess, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityZoneBookAddressExists {
		d.SetId(d.Get("zone").(string) + idSeparator + d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security zone address-book address %v not exists in zone %s after commit "+
				"=> check your config", d.Get("name").(string), d.Get("zone").(string)))...)
	}

	return append(diagWarns, resourceSecurityZoneBookAddressReadWJunSess(d, sess, junSess)...)
}

func resourceSecurityZoneBookAddressRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)

	return resourceSecurityZoneBookAddressReadWJunSess(d, sess, junSess)
}

func resourceSecurityZoneBookAddressReadWJunSess(d *schema.ResourceData, sess *Session, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	zoneBookAddressOptions, err := readSecurityZoneBookAddress(
		d.Get("zone").(string),
		d.Get("name").(string),
		sess, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if zoneBookAddressOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityZoneBookAddressData(d, zoneBookAddressOptions)
	}

	return nil
}

func resourceSecurityZoneBookAddressUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delSecurityZoneBookAddress(d.Get("zone").(string), d.Get("name").(string), sess, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityZoneBookAddress(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityZoneBookAddress(d.Get("zone").(string), d.Get("name").(string), sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityZoneBookAddress(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_zone_book_address", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityZoneBookAddressReadWJunSess(d, sess, junSess)...)
}

func resourceSecurityZoneBookAddressDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delSecurityZoneBookAddress(d.Get("zone").(string), d.Get("name").(string), sess, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityZoneBookAddress(d.Get("zone").(string), d.Get("name").(string), sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_zone_book_address", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityZoneBookAddressImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	idList := strings.Split(d.Id(), idSeparator)
	if len(idList) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	securityZoneBookAddressExists, err := checkSecurityZoneBookAddresssExists(idList[0], idList[1], sess, junSess)
	if err != nil {
		return nil, err
	}
	if !securityZoneBookAddressExists {
		return nil, fmt.Errorf(
			"don't find zone address-book address with id '%v' (id must be <zone>"+idSeparator+"<name>)", d.Id())
	}
	zoneBookAddressOptions, err := readSecurityZoneBookAddress(idList[0], idList[1], sess, junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityZoneBookAddressData(d, zoneBookAddressOptions)

	result[0] = d

	return result, nil
}

func checkSecurityZoneBookAddresssExists(zone, address string, sess *Session, junSess *junosSession) (bool, error) {
	showConfig, err := sess.command(cmdShowConfig+
		"security zones security-zone "+zone+" address-book address "+address+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityZoneBookAddress(d *schema.ResourceData, sess *Session, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := "set security zones security-zone " +
		d.Get("zone").(string) + " address-book address " + d.Get("name").(string) + " "

	if v := d.Get("cidr").(string); v != "" {
		configSet = append(configSet, setPrefix+v)
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := d.Get("dns_name").(string); v != "" {
		configSet = append(configSet, setPrefix+"dns-name "+v)
		if d.Get("dns_ipv4_only").(bool) {
			configSet = append(configSet, setPrefix+"dns-name "+v+" ipv4-only")
		}
		if d.Get("dns_ipv6_only").(bool) {
			configSet = append(configSet, setPrefix+"dns-name "+v+" ipv6-only")
		}
	}
	if v := d.Get("range_from").(string); v != "" {
		configSet = append(configSet, setPrefix+"range-address "+v+" to "+d.Get("range_to").(string))
	}
	if v := d.Get("wildcard").(string); v != "" {
		configSet = append(configSet, setPrefix+"wildcard-address "+v)
	}

	return sess.configSet(configSet, junSess)
}

func readSecurityZoneBookAddress(zone, address string, sess *Session, junSess *junosSession,
) (zoneBookAddressOptions, error) {
	var confRead zoneBookAddressOptions

	showConfig, err := sess.command(cmdShowConfig+
		"security zones security-zone "+zone+" address-book address "+address+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = address
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
			case strings.HasPrefix(itemTrim, "dns-name "):
				dnsValue := strings.TrimPrefix(itemTrim, "dns-name ")
				switch {
				case strings.HasSuffix(itemTrim, " ipv4-only"):
					confRead.dnsIPv4Only = true
					dnsValue = strings.TrimSuffix(strings.TrimPrefix(itemTrim, "dns-name "), " ipv4-only")
				case strings.HasSuffix(itemTrim, " ipv6-only"):
					confRead.dnsIPv6Only = true
					dnsValue = strings.TrimSuffix(strings.TrimPrefix(itemTrim, "dns-name "), " ipv6-only")
				}
				confRead.dnsName = dnsValue
			case strings.HasPrefix(itemTrim, "range-address "):
				itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "range-address "), " ")
				confRead.rangeFrom = itemTrimSplit[0]
				confRead.rangeTo = itemTrimSplit[2]
			case strings.HasPrefix(itemTrim, "wildcard-address "):
				confRead.wildcard = strings.TrimPrefix(itemTrim, "wildcard-address ")
			case strings.Contains(itemTrim, "/"):
				confRead.cidr = itemTrim
			}
		}
	}

	return confRead, nil
}

func delSecurityZoneBookAddress(zone, address string, sess *Session, junSess *junosSession) error {
	configSet := []string{"delete security zones security-zone " + zone + " address-book address " + address}

	return sess.configSet(configSet, junSess)
}

func fillSecurityZoneBookAddressData(d *schema.ResourceData, zoneBookAddressOptions zoneBookAddressOptions) {
	if tfErr := d.Set("name", zoneBookAddressOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("zone", zoneBookAddressOptions.zone); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("cidr", zoneBookAddressOptions.cidr); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", zoneBookAddressOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dns_ipv4_only", zoneBookAddressOptions.dnsIPv4Only); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dns_ipv6_only", zoneBookAddressOptions.dnsIPv6Only); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dns_name", zoneBookAddressOptions.dnsName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("range_from", zoneBookAddressOptions.rangeFrom); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("range_to", zoneBookAddressOptions.rangeTo); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("wildcard", zoneBookAddressOptions.wildcard); tfErr != nil {
		panic(tfErr)
	}
}
