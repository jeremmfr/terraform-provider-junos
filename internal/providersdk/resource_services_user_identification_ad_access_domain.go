package providersdk

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	jdecode "github.com/jeremmfr/junosdecode"
)

type svcUserIdentAdAccessDomainOptions struct {
	name                      string
	userName                  string
	userPassword              string
	domainController          []map[string]interface{}
	ipUserMappingDiscoveryWmi []map[string]interface{}
	userGroupMappingLdap      []map[string]interface{}
}

func resourceServicesUserIdentAdAccessDomain() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceServicesUserIdentAdAccessDomainCreate,
		ReadWithoutTimeout:   resourceServicesUserIdentAdAccessDomainRead,
		UpdateWithoutTimeout: resourceServicesUserIdentAdAccessDomainUpdate,
		DeleteWithoutTimeout: resourceServicesUserIdentAdAccessDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServicesUserIdentAdAccessDomainImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefAndDots),
			},
			"user_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefAndDots),
			},
			"user_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"domain_controller": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefAndDots),
						},
						"address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
					},
				},
			},
			"ip_user_mapping_discovery_wmi": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_log_scanning_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(5, 60),
						},
						"initial_event_log_timespan": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 168),
						},
					},
				},
			},
			"user_group_mapping_ldap": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base": {
							Type:     schema.TypeString,
							Required: true,
						},
						"address": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"auth_algo_simple": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ssl": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"user_name": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefAndDots),
						},
						"user_password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func resourceServicesUserIdentAdAccessDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setServicesUserIdentAdAccessDomain(d, junSess); err != nil {
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
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	svcUserIdentAdAccessDomainExists, err := checkServicesUserIdentAdAccessDomainExists(
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if svcUserIdentAdAccessDomainExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf(
				"services user-identification active-directory-access domain %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesUserIdentAdAccessDomain(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_services_user_identification_ad_access_domain")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	svcUserIdentAdAccessDomainExists, err = checkServicesUserIdentAdAccessDomainExists(
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if svcUserIdentAdAccessDomainExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"services user-identification active-directory-access domain %v "+
				"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesUserIdentAdAccessDomainReadWJunSess(d, junSess)...)
}

func resourceServicesUserIdentAdAccessDomainRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceServicesUserIdentAdAccessDomainReadWJunSess(d, junSess)
}

func resourceServicesUserIdentAdAccessDomainReadWJunSess(
	d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	svcUserIdentAdAccessDomainOptions, err := readServicesUserIdentAdAccessDomain(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if svcUserIdentAdAccessDomainOptions.name == "" {
		d.SetId("")
	} else {
		fillServicesUserIdentAdAccessDomainData(d, svcUserIdentAdAccessDomainOptions)
	}

	return nil
}

func resourceServicesUserIdentAdAccessDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delServicesUserIdentAdAccessDomain(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setServicesUserIdentAdAccessDomain(d, junSess); err != nil {
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
	if err := delServicesUserIdentAdAccessDomain(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesUserIdentAdAccessDomain(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_services_user_identification_ad_access_domain")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesUserIdentAdAccessDomainReadWJunSess(d, junSess)...)
}

func resourceServicesUserIdentAdAccessDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delServicesUserIdentAdAccessDomain(d.Get("name").(string), junSess); err != nil {
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
	if err := delServicesUserIdentAdAccessDomain(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_services_user_identification_ad_access_domain")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesUserIdentAdAccessDomainImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	svcUserIdentAdAccessDomainExists, err := checkServicesUserIdentAdAccessDomainExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !svcUserIdentAdAccessDomainExists {
		return nil, fmt.Errorf("don't find services user-identification "+
			"active-directory-access domain with id '%v' (id must be <name>)", d.Id())
	}
	svcUserIdentAdAccessDomainOptions, err := readServicesUserIdentAdAccessDomain(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillServicesUserIdentAdAccessDomainData(d, svcUserIdentAdAccessDomainOptions)

	result[0] = d

	return result, nil
}

func checkServicesUserIdentAdAccessDomainExists(domain string, junSess *junos.Session,
) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services user-identification active-directory-access domain " + domain + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setServicesUserIdentAdAccessDomain(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set services user-identification active-directory-access domain " + d.Get("name").(string) + " "
	configSet = append(configSet, setPrefix+"user "+d.Get("user_name").(string))
	configSet = append(configSet, setPrefix+"user password \""+d.Get("user_password").(string)+"\"")
	domainControllerNameList := make([]string, 0)
	for _, v := range d.Get("domain_controller").([]interface{}) {
		domainController := v.(map[string]interface{})
		if slices.Contains(domainControllerNameList, domainController["name"].(string)) {
			return fmt.Errorf("multiple blocks domain_controller with the same name %s", domainController["name"].(string))
		}
		domainControllerNameList = append(domainControllerNameList, domainController["name"].(string))
		configSet = append(configSet, setPrefix+"domain-controller "+domainController["name"].(string)+
			" address "+domainController["address"].(string))
	}
	for _, v := range d.Get("ip_user_mapping_discovery_wmi").([]interface{}) {
		configSet = append(configSet, setPrefix+"ip-user-mapping discovery-method wmi")
		if v != nil {
			wmi := v.(map[string]interface{})
			if v2 := wmi["event_log_scanning_interval"].(int); v2 != 0 {
				configSet = append(configSet,
					setPrefix+"ip-user-mapping discovery-method wmi event-log-scanning-interval "+strconv.Itoa(v2))
			}
			if v2 := wmi["initial_event_log_timespan"].(int); v2 != 0 {
				configSet = append(configSet,
					setPrefix+"ip-user-mapping discovery-method wmi initial-event-log-timespan "+strconv.Itoa(v2))
			}
		}
	}
	for _, v := range d.Get("user_group_mapping_ldap").([]interface{}) {
		ldap := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+"user-group-mapping ldap base \""+ldap["base"].(string)+"\"")
		for _, v2 := range ldap["address"].([]interface{}) {
			configSet = append(configSet, setPrefix+"user-group-mapping ldap address "+v2.(string))
		}
		if ldap["auth_algo_simple"].(bool) {
			configSet = append(configSet, setPrefix+"user-group-mapping ldap authentication-algorithm simple")
		}
		if ldap["ssl"].(bool) {
			configSet = append(configSet, setPrefix+"user-group-mapping ldap ssl")
		}
		if v2 := ldap["user_name"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"user-group-mapping ldap user "+v2)
		}
		if v2 := ldap["user_password"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"user-group-mapping ldap user password \""+v2+"\"")
		}
	}

	return junSess.ConfigSet(configSet)
}

func readServicesUserIdentAdAccessDomain(domain string, junSess *junos.Session,
) (confRead svcUserIdentAdAccessDomainOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services user-identification active-directory-access domain " + domain + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = domain
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "user password "):
				confRead.userPassword, err = jdecode.Decode(strings.Trim(itemTrim, "\""))
				if err != nil {
					return confRead, fmt.Errorf("decoding user password: %w", err)
				}
			case balt.CutPrefixInString(&itemTrim, "user "):
				confRead.userName = itemTrim
			case balt.CutPrefixInString(&itemTrim, "domain-controller "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 3 { // <name> address <address>
					return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "domain-controller", itemTrim)
				}
				confRead.domainController = append(confRead.domainController, map[string]interface{}{
					"name":    itemTrimFields[0],
					"address": itemTrimFields[2],
				})
			case balt.CutPrefixInString(&itemTrim, "ip-user-mapping discovery-method wmi"):
				if len(confRead.ipUserMappingDiscoveryWmi) == 0 {
					confRead.ipUserMappingDiscoveryWmi = append(confRead.ipUserMappingDiscoveryWmi, map[string]interface{}{
						"event_log_scanning_interval": 0,
						"initial_event_log_timespan":  0,
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " event-log-scanning-interval "):
					confRead.ipUserMappingDiscoveryWmi[0]["event_log_scanning_interval"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, " initial-event-log-timespan "):
					confRead.ipUserMappingDiscoveryWmi[0]["initial_event_log_timespan"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			case balt.CutPrefixInString(&itemTrim, "user-group-mapping ldap "):
				if len(confRead.userGroupMappingLdap) == 0 {
					confRead.userGroupMappingLdap = append(confRead.userGroupMappingLdap, map[string]interface{}{
						"base":             "",
						"address":          make([]string, 0),
						"auth_algo_simple": false,
						"ssl":              false,
						"user_name":        "",
						"user_password":    "",
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "base "):
					confRead.userGroupMappingLdap[0]["base"] = strings.Trim(itemTrim, "\"")
				case balt.CutPrefixInString(&itemTrim, "address "):
					confRead.userGroupMappingLdap[0]["address"] = append(
						confRead.userGroupMappingLdap[0]["address"].([]string),
						itemTrim,
					)
				case itemTrim == "authentication-algorithm simple":
					confRead.userGroupMappingLdap[0]["auth_algo_simple"] = true
				case itemTrim == "ssl":
					confRead.userGroupMappingLdap[0]["ssl"] = true
				case balt.CutPrefixInString(&itemTrim, "user password "):
					confRead.userGroupMappingLdap[0]["user_password"], err = jdecode.Decode(strings.Trim(itemTrim, "\""))
					if err != nil {
						return confRead, fmt.Errorf("decoding user password: %w", err)
					}
				case balt.CutPrefixInString(&itemTrim, "user "):
					confRead.userGroupMappingLdap[0]["user_name"] = itemTrim
				}
			}
		}
	}

	return confRead, nil
}

func delServicesUserIdentAdAccessDomain(domain string, junSess *junos.Session) error {
	configSet := []string{
		"delete services user-identification active-directory-access domain " + domain,
	}

	return junSess.ConfigSet(configSet)
}

func fillServicesUserIdentAdAccessDomainData(
	d *schema.ResourceData, svcUserIdentAdAccessDomainOptions svcUserIdentAdAccessDomainOptions,
) {
	if tfErr := d.Set("name", svcUserIdentAdAccessDomainOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("user_name", svcUserIdentAdAccessDomainOptions.userName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("user_password", svcUserIdentAdAccessDomainOptions.userPassword); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("domain_controller", svcUserIdentAdAccessDomainOptions.domainController); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ip_user_mapping_discovery_wmi",
		svcUserIdentAdAccessDomainOptions.ipUserMappingDiscoveryWmi); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("user_group_mapping_ldap", svcUserIdentAdAccessDomainOptions.userGroupMappingLdap); tfErr != nil {
		panic(tfErr)
	}
}
