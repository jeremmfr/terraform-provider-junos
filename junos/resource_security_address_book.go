package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
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
		CreateWithoutTimeout: resourceSecurityAddressBookCreate,
		ReadWithoutTimeout:   resourceSecurityAddressBookRead,
		UpdateWithoutTimeout: resourceSecurityAddressBookUpdate,
		DeleteWithoutTimeout: resourceSecurityAddressBookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityAddressBookImport,
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
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"network_address": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"value": {
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
			"wildcard_address": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"value": {
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
			"dns_name": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			"range_address": {
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
			"address_set": {
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
		},
	}
}

func resourceSecurityAddressBookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSecurityAddressBook(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if !checkCompatibilitySecurity(junSess) {
		return diag.FromErr(fmt.Errorf("security policy not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	addressBookExists, err := checkSecurityAddressBookExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if addressBookExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security address book %v already exists", d.Get("name").(string)))...)
	}
	if err := setSecurityAddressBook(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_security_address_book", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	addressBookExists, err = checkSecurityAddressBookExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if addressBookExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security address book  %v does not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityAddressBookReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityAddressBookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSecurityAddressBookReadWJunSess(d, clt, junSess)
}

func resourceSecurityAddressBookReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	addressOptions, err := readSecurityAddressBook(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if addressOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityAddressBookData(d, addressOptions)
	}

	return nil
}

func resourceSecurityAddressBookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSecurityAddressBook(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityAddressBook(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityAddressBook(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityAddressBook(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_security_address_book", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityAddressBookReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityAddressBookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delSecurityAddressBook(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityAddressBook(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_security_address_book", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityAddressBookImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	securityAddressBookExists, err := checkSecurityAddressBookExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !securityAddressBookExists {
		return nil, fmt.Errorf("don't find address book with id '%v' (id must be <name>)", d.Id())
	}
	addressOptions, err := readSecurityAddressBook(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityAddressBookData(d, addressOptions)

	result[0] = d

	return result, nil
}

func checkSecurityAddressBookExists(addrBook string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+
		"security address-book "+addrBook+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityAddressBook(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
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
	addressNameList := make([]string, 0)
	for _, v := range d.Get("network_address").(*schema.Set).List() {
		address := v.(map[string]interface{})
		if bchk.InSlice(address["name"].(string), addressNameList) {
			return fmt.Errorf("multiple addresses with the same name %s", address["name"].(string))
		}
		addressNameList = append(addressNameList, address["name"].(string))
		setPrefixAddr := setPrefix + " address " + address["name"].(string) + " "
		configSet = append(configSet, setPrefixAddr+address["value"].(string))
		if address["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+address["description"].(string)+"\"")
		}
	}
	for _, v := range d.Get("wildcard_address").(*schema.Set).List() {
		address := v.(map[string]interface{})
		if bchk.InSlice(address["name"].(string), addressNameList) {
			return fmt.Errorf("multiple addresses with the same name %s", address["name"].(string))
		}
		addressNameList = append(addressNameList, address["name"].(string))
		setPrefixAddr := setPrefix + " address " + address["name"].(string)
		configSet = append(configSet, setPrefixAddr+" wildcard-address "+address["value"].(string))
		if address["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+address["description"].(string)+"\"")
		}
	}
	for _, v := range d.Get("dns_name").(*schema.Set).List() {
		address := v.(map[string]interface{})
		if bchk.InSlice(address["name"].(string), addressNameList) {
			return fmt.Errorf("multiple addresses with the same name %s", address["name"].(string))
		}
		addressNameList = append(addressNameList, address["name"].(string))
		setPrefixAddr := setPrefix + " address " + address["name"].(string)
		configSet = append(configSet, setPrefixAddr+" dns-name "+address["value"].(string))
		if address["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+address["description"].(string)+"\"")
		}
	}
	for _, v := range d.Get("range_address").(*schema.Set).List() {
		address := v.(map[string]interface{})
		if bchk.InSlice(address["name"].(string), addressNameList) {
			return fmt.Errorf("multiple addresses with the same name %s", address["name"].(string))
		}
		addressNameList = append(addressNameList, address["name"].(string))
		setPrefixAddr := setPrefix + " address " + address["name"].(string)
		configSet = append(configSet, setPrefixAddr+" range-address "+address["from"].(string)+" to "+address["to"].(string))
		if address["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+address["description"].(string)+"\"")
		}
	}
	for _, v := range d.Get("address_set").(*schema.Set).List() {
		addressSet := v.(map[string]interface{})
		if bchk.InSlice(addressSet["name"].(string), addressNameList) {
			return fmt.Errorf("multiple addresses or address-sets with the same name %s", addressSet["name"].(string))
		}
		addressNameList = append(addressNameList, addressSet["name"].(string))
		setPrefixAddrSet := setPrefix + " address-set " + addressSet["name"].(string)
		if len(addressSet["address"].(*schema.Set).List()) == 0 &&
			len(addressSet["address_set"].(*schema.Set).List()) == 0 {
			return fmt.Errorf("at least one of address or address_set is required "+
				"in address_set %s", addressSet["name"].(string))
		}
		for _, addr := range sortSetOfString(addressSet["address"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixAddrSet+" address "+addr)
		}
		for _, addrSet := range sortSetOfString(addressSet["address_set"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixAddrSet+" address-set "+addrSet)
		}
		if addressSet["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddrSet+"description \""+addressSet["description"].(string)+"\"")
		}
	}

	return clt.configSet(configSet, junSess)
}

func readSecurityAddressBook(addrBook string, clt *Client, junSess *junosSession) (addressBookOptions, error) {
	var confRead addressBookOptions

	showConfig, err := clt.command(cmdShowConfig+
		"security address-book "+addrBook+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	descMap := make(map[string]string)
	if showConfig != emptyW {
		confRead.name = addrBook
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
				addressSplit := strings.Split(itemTrim, " ")
				itemTrimAddress := strings.TrimPrefix(itemTrim, "address "+addressSplit[1]+" ")
				switch {
				case strings.HasPrefix(itemTrimAddress, "description "):
					descMap[addressSplit[1]] = strings.Trim(strings.TrimPrefix(itemTrimAddress, "description "), "\"")
				case strings.HasPrefix(itemTrimAddress, "wildcard-address "):
					confRead.wildcardAddress = append(confRead.wildcardAddress, map[string]interface{}{
						"name":        addressSplit[1],
						"value":       strings.TrimPrefix(itemTrimAddress, "wildcard-address "),
						"description": descMap[addressSplit[1]],
					})
				case strings.HasPrefix(itemTrimAddress, "range-address "):
					rangeAddr := strings.TrimPrefix(itemTrimAddress, "range-address ")
					addresses := strings.Split(rangeAddr, " ")
					confRead.rangeAddress = append(confRead.rangeAddress, map[string]interface{}{
						"name":        addressSplit[1],
						"from":        addresses[0],
						"to":          addresses[2],
						"description": descMap[addressSplit[1]],
					})
				case strings.HasPrefix(itemTrimAddress, "dns-name "):
					confRead.dnsName = append(confRead.dnsName, map[string]interface{}{
						"name":        addressSplit[1],
						"value":       strings.TrimPrefix(itemTrimAddress, "dns-name "),
						"description": descMap[addressSplit[1]],
					})
				default:
					confRead.networkAddress = append(confRead.networkAddress, map[string]interface{}{
						"name":        addressSplit[1],
						"value":       itemTrimAddress,
						"description": descMap[addressSplit[1]],
					})
				}
			case strings.HasPrefix(itemTrim, "address-set "):
				addressSetSplit := strings.Split(strings.TrimPrefix(itemTrim, "address-set "), " ")
				adSet := map[string]interface{}{
					"name":        addressSetSplit[0],
					"address":     make([]string, 0),
					"address_set": make([]string, 0),
					"description": "",
				}
				confRead.addressSet = copyAndRemoveItemMapList("name", adSet, confRead.addressSet)
				switch {
				case strings.HasPrefix(itemTrim, "address-set "+addressSetSplit[0]+" description "):
					adSet["description"] = strings.Trim(strings.TrimPrefix(
						itemTrim, "address-set "+addressSetSplit[0]+" description "), "\"")
				case strings.HasPrefix(itemTrim, "address-set "+addressSetSplit[0]+" address "):
					adSet["address"] = append(adSet["address"].([]string),
						strings.TrimPrefix(itemTrim, "address-set "+addressSetSplit[0]+" address "))
				case strings.HasPrefix(itemTrim, "address-set "+addressSetSplit[0]+" address-set "):
					adSet["address_set"] = append(adSet["address_set"].([]string),
						strings.TrimPrefix(itemTrim, "address-set "+addressSetSplit[0]+" address-set "))
				}
				confRead.addressSet = append(confRead.addressSet, adSet)
			case strings.HasPrefix(itemTrim, "attach zone"):
				confRead.attachZone = append(confRead.attachZone, strings.TrimPrefix(itemTrim, "attach zone "))
			}
		}
	}
	confRead.networkAddress = copySecurityAddressBookAddressDescriptions(descMap, confRead.networkAddress)
	confRead.dnsName = copySecurityAddressBookAddressDescriptions(descMap, confRead.dnsName)
	confRead.rangeAddress = copySecurityAddressBookAddressDescriptions(descMap, confRead.rangeAddress)
	confRead.wildcardAddress = copySecurityAddressBookAddressDescriptions(descMap, confRead.wildcardAddress)

	return confRead, nil
}

func delSecurityAddressBook(addrBook string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security address-book "+addrBook)

	return clt.configSet(configSet, junSess)
}

func fillSecurityAddressBookData(d *schema.ResourceData, addressOptions addressBookOptions) {
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

func copySecurityAddressBookAddressDescriptions(descMap map[string]string, addrList []map[string]interface{},
) (newList []map[string]interface{}) {
	for _, ele := range addrList {
		ele["description"] = descMap[ele["name"].(string)]
		newList = append(newList, ele)
	}

	return newList
}
