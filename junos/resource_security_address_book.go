package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type addressBookOptions struct {
	name            string
	description     string
	attachZone      []string
	networkAddress  []map[string]interface{}
	wildcardAddress []map[string]interface{}
	dnsName         []map[string]interface{}
	rangeAddress    []map[string]interface{}
	addressSet      []map[string]interface{}
}

func resourceSecurityAddressBook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityAddressBookCreate,
		ReadContext:   resourceSecurityAddressBookRead,
		UpdateContext: resourceSecurityAddressBookUpdate,
		DeleteContext: resourceSecurityAddressBookDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityAddressBookImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "global",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attach_zone": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"network_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
						},
					},
				},
			},
			"wildcard_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateWildcardFunc(),
						},
					},
				},
			},
			"dns_name": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"range_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
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
					},
				},
			},
			"address_set": {
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
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceSecurityAddressBookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security policy not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	addressBookExists, err := checkAddressBookExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if addressBookExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security address book %v already exists", d.Get("name").(string)))
	}
	if err := setAddressBook(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_address_book", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	addressBookExists, err = checkAddressBookExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if addressBookExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security address book  %v does not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityAddressBookReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityAddressBookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityAddressBookReadWJnprSess(d, m, jnprSess)
}
func resourceSecurityAddressBookReadWJnprSess(d *schema.ResourceData, m interface{},
	jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	addressOptions, err := readSecurityAddressBook(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if addressOptions.name == "" {
		d.SetId("")
	} else {
		fillAddressBookData(d, addressOptions)
	}

	return nil
}

func resourceSecurityAddressBookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)

	if err := delSecurityAddressBook(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}

	if err := setAddressBook(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security_address_book", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityAddressBookReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityAddressBookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSecurityAddressBook(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_address_book", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityAddressBookImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	securityAddressBookExists, err := checkAddressBookExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityAddressBookExists {
		return nil, fmt.Errorf("don't find address book with id '%v' (id must be <name>)", d.Id())
	}
	addressOptions, err := readSecurityAddressBook(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillAddressBookData(d, addressOptions)

	result[0] = d

	return result, nil
}

func checkAddressBookExists(addrBook string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)

	addrBookConfig, err := sess.command("show configuration security address-book "+addrBook+
		" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if addrBookConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setAddressBook(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := "set security address-book " + d.Get("name").(string)

	if d.Get("description").(string) != "" {
		configSet = append(configSet, setPrefix+" description \""+d.Get("description").(string)+"\"")
	}
	for _, v := range d.Get("attach_zone").([]interface{}) {
		if d.Get("name").(string) == "global" {
			return fmt.Errorf("cannot attach global address book to a zone")
		}
		attachZone := v.(string)
		configSet = append(configSet, setPrefix+" attach zone "+attachZone)
	}
	for _, v := range d.Get("network_address").([]interface{}) {
		address := v.(map[string]interface{})
		setPrefixAddr := setPrefix + " address " + address["name"].(string) + " "
		if address["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+address["description"].(string)+"\"")
		}
		configSet = append(configSet, setPrefixAddr+address["value"].(string))
	}
	for _, v := range d.Get("wildcard_address").([]interface{}) {
		address := v.(map[string]interface{})
		setPrefixAddr := setPrefix + " address " + address["name"].(string)
		if address["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+address["description"].(string)+"\"")
		}
		configSet = append(configSet, setPrefixAddr+" wildcard-address "+address["value"].(string))
	}
	for _, v := range d.Get("dns_name").([]interface{}) {
		address := v.(map[string]interface{})
		setPrefixAddr := setPrefix + " address " + address["name"].(string)
		if address["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+address["description"].(string)+"\"")
		}
		configSet = append(configSet, setPrefixAddr+" dns-name "+address["value"].(string))
	}
	for _, v := range d.Get("range_address").([]interface{}) {
		address := v.(map[string]interface{})
		setPrefixAddr := setPrefix + " address " + address["name"].(string)
		if address["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+address["description"].(string)+"\"")
		}
		configSet = append(configSet, setPrefixAddr+" range-address "+address["from"].(string)+" to "+address["to"].(string))
	}
	for _, v := range d.Get("address_set").([]interface{}) {
		addressSet := v.(map[string]interface{})
		setPrefixAddrSet := setPrefix + " address-set " + addressSet["name"].(string)
		if addressSet["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddrSet+"description \""+addressSet["description"].(string)+"\"")
		}
		for _, addr := range addressSet["address"].([]interface{}) {
			configSet = append(configSet, setPrefixAddrSet+" address "+addr.(string))
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func readSecurityAddressBook(addrBook string, m interface{}, jnprSess *NetconfObject) (addressBookOptions, error) {
	sess := m.(*Session)
	var confRead addressBookOptions

	securityAddressBookConfig, err := sess.command("show configuration security address-book "+addrBook+
		" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	descMap := make(map[string]string)
	if securityAddressBookConfig != emptyWord {
		confRead.name = addrBook
		for _, item := range strings.Split(securityAddressBookConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "address "):
				addressSplit := strings.Split(itemTrim, " ")
				itemTrimAddress := strings.TrimPrefix(itemTrim, "address "+addressSplit[1]+" ")
				switch {
				case strings.HasPrefix(itemTrimAddress, "description "):
					descMap[addressSplit[1]] = strings.Trim(strings.TrimPrefix(itemTrimAddress, "description "), "\"")
				case strings.HasPrefix(itemTrimAddress, "wildcard-address "):
					confRead.wildcardAddress = append(confRead.wildcardAddress, map[string]interface{}{
						"name":        addressSplit[1],
						"description": descMap[addressSplit[1]],
						"value":       strings.TrimPrefix(itemTrimAddress, "wildcard-address "),
					})
				case strings.HasPrefix(itemTrimAddress, "range-address "):
					rangeAddr := strings.TrimPrefix(itemTrimAddress, "range-address ")
					addresses := strings.Split(rangeAddr, " ")
					confRead.rangeAddress = append(confRead.rangeAddress, map[string]interface{}{
						"name":        addressSplit[1],
						"description": descMap[addressSplit[1]],
						"from":        addresses[0],
						"to":          addresses[2],
					})
				case strings.HasPrefix(itemTrimAddress, "dns-name "):
					confRead.dnsName = append(confRead.dnsName, map[string]interface{}{
						"name":        addressSplit[1],
						"description": descMap[addressSplit[1]],
						"value":       strings.TrimPrefix(itemTrimAddress, "dns-name "),
					})
				default:
					confRead.networkAddress = append(confRead.networkAddress, map[string]interface{}{
						"name":        addressSplit[1],
						"description": descMap[addressSplit[1]],
						"value":       itemTrimAddress,
					})
				}
			case strings.HasPrefix(itemTrim, "address-set "):
				addressSetSplit := strings.Split(strings.TrimPrefix(itemTrim, "address-set "), " ")
				m := map[string]interface{}{
					"name":        addressSetSplit[0],
					"address":     make([]string, 0),
					"description": "",
				}
				m, confRead.addressSet = copyAndRemoveItemMapList("name", false, m, confRead.addressSet)
				if addressSetSplit[1] == "description" {
					m["description"] = strings.Trim(strings.TrimPrefix(
						itemTrim, "address-set "+addressSetSplit[0]+" description "), "\"")
				} else {
					m["address"] = append(m["address"].([]string), addressSetSplit[2])
				}
				confRead.addressSet = append(confRead.addressSet, m)
			case strings.HasPrefix(itemTrim, "attach zone"):
				confRead.attachZone = append(confRead.attachZone, strings.TrimPrefix(itemTrim, "attach zone "))
			}
		}
	}
	confRead.networkAddress = copyAddressDescriptions(descMap, confRead.networkAddress)
	confRead.dnsName = copyAddressDescriptions(descMap, confRead.dnsName)
	confRead.rangeAddress = copyAddressDescriptions(descMap, confRead.rangeAddress)
	confRead.wildcardAddress = copyAddressDescriptions(descMap, confRead.wildcardAddress)

	return confRead, nil
}

func delSecurityAddressBook(addrBook string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security address-book "+addrBook)

	return sess.configSet(configSet, jnprSess)
}

func fillAddressBookData(d *schema.ResourceData, addressOptions addressBookOptions) {
	if tfErr := d.Set("name", addressOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", addressOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("attach_zone", addressOptions.attachZone); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("network_address", addressOptions.networkAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("wildcard_address", addressOptions.wildcardAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dns_name", addressOptions.dnsName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("range_address", addressOptions.rangeAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_set", addressOptions.addressSet); tfErr != nil {
		panic(tfErr)
	}
}

func copyAddressDescriptions(descMap map[string]string,
	addrList []map[string]interface{}) (newList []map[string]interface{}) {
	for _, ele := range addrList {
		ele["description"] = descMap[ele["name"].(string)]
		newList = append(newList, ele)
	}

	return newList
}
