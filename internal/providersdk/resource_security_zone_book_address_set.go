package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
		CreateWithoutTimeout: resourceSecurityZoneBookAddressSetCreate,
		ReadWithoutTimeout:   resourceSecurityZoneBookAddressSetRead,
		UpdateWithoutTimeout: resourceSecurityZoneBookAddressSetUpdate,
		DeleteWithoutTimeout: resourceSecurityZoneBookAddressSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityZoneBookAddressSetImport,
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
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
				},
			},
			"address_set": {
				Type:         schema.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"address", "address_set"},
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSecurityZoneBookAddressSetCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityZoneBookAddressSet(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("zone").(string) + junos.IDSeparator + d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilitySecurity() {
		return diag.FromErr(fmt.Errorf("security zone address-book address-set not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	zonesExists, err := checkSecurityZonesExists(d.Get("zone").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !zonesExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security zone %v doesn't exist", d.Get("zone").(string)))...)
	}
	securityZoneBookAddressSetExists, err := checkSecurityZoneBookAddressSetsExists(
		d.Get("zone").(string),
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityZoneBookAddressSetExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security zone address-book address-set %v already exists in zone %s",
			d.Get("name").(string), d.Get("zone").(string)))...)
	}

	if err := setSecurityZoneBookAddressSet(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_zone_book_address_set")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityZoneBookAddressSetExists, err = checkSecurityZoneBookAddressSetsExists(
		d.Get("zone").(string),
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityZoneBookAddressSetExists {
		d.SetId(d.Get("zone").(string) + junos.IDSeparator + d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security zone address-book address-set %v not exists in zone %s after commit "+
				"=> check your config", d.Get("name").(string), d.Get("zone").(string)))...)
	}

	return append(diagWarns, resourceSecurityZoneBookAddressSetReadWJunSess(d, junSess)...)
}

func resourceSecurityZoneBookAddressSetRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityZoneBookAddressSetReadWJunSess(d, junSess)
}

func resourceSecurityZoneBookAddressSetReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	zoneBookAddressSetOptions, err := readSecurityZoneBookAddressSet(
		d.Get("zone").(string),
		d.Get("name").(string),
		junSess,
	)
	junos.MutexUnlock()
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

func resourceSecurityZoneBookAddressSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityZoneBookAddressSet(d.Get("zone").(string), d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityZoneBookAddressSet(d, junSess); err != nil {
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
	if err := delSecurityZoneBookAddressSet(d.Get("zone").(string), d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityZoneBookAddressSet(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_zone_book_address_set")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityZoneBookAddressSetReadWJunSess(d, junSess)...)
}

func resourceSecurityZoneBookAddressSetDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityZoneBookAddressSet(d.Get("zone").(string), d.Get("name").(string), junSess); err != nil {
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
	if err := delSecurityZoneBookAddressSet(d.Get("zone").(string), d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_zone_book_address_set")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityZoneBookAddressSetImport(ctx context.Context, d *schema.ResourceData, m interface{},
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
	securityZoneBookAddressSetExists, err := checkSecurityZoneBookAddressSetsExists(idList[0], idList[1], junSess)
	if err != nil {
		return nil, err
	}
	if !securityZoneBookAddressSetExists {
		return nil, fmt.Errorf(
			"don't find zone address-book address-set with id '%v' (id must be <zone>"+junos.IDSeparator+"<name>)", d.Id())
	}
	zoneBookAddressSetOptions, err := readSecurityZoneBookAddressSet(idList[0], idList[1], junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityZoneBookAddressSetData(d, zoneBookAddressSetOptions)

	result[0] = d

	return result, nil
}

func checkSecurityZoneBookAddressSetsExists(zone, addressSet string, junSess *junos.Session,
) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security zones security-zone " + zone + " address-book address-set " + addressSet + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityZoneBookAddressSet(d *schema.ResourceData, junSess *junos.Session) error {
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

	return junSess.ConfigSet(configSet)
}

func readSecurityZoneBookAddressSet(zone, addressSet string, junSess *junos.Session,
) (confRead zoneBookAddressSetOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security zones security-zone " + zone + " address-book address-set " + addressSet + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = addressSet
		confRead.zone = zone
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "address "):
				confRead.address = append(confRead.address, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "address-set "):
				confRead.addressSet = append(confRead.addressSet, itemTrim)
			}
		}
	}

	return confRead, nil
}

func delSecurityZoneBookAddressSet(zone, addressSet string, junSess *junos.Session) error {
	configSet := []string{"delete security zones security-zone " + zone + " address-book address-set " + addressSet}

	return junSess.ConfigSet(configSet)
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
