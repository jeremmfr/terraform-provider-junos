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
)

type dynamicAddressFeedServerOptions struct {
	validateCertAttrSubOrSan bool
	holdInterval             int
	updateInterval           int
	description              string
	hostname                 string
	name                     string
	url                      string
	tlsProfile               string
	feedName                 []map[string]interface{}
}

func resourceSecurityDynamicAddressFeedServer() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityDynamicAddressFeedServerCreate,
		ReadWithoutTimeout:   resourceSecurityDynamicAddressFeedServerRead,
		UpdateWithoutTimeout: resourceSecurityDynamicAddressFeedServerUpdate,
		DeleteWithoutTimeout: resourceSecurityDynamicAddressFeedServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityDynamicAddressFeedServerImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 16, formatDefault),
			},
			"hostname": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"hostname", "url"},
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"hostname", "url"},
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"feed_name": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hold_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 4294967295),
						},
						"update_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(30, 4294967295),
						},
					},
				},
			},
			"hold_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(0, 4294967295),
			},
			"tls_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(30, 4294967295),
			},
			"validate_certificate_attributes_subject_or_san": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"tls_profile"},
			},
		},
	}
}

func resourceSecurityDynamicAddressFeedServerCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityDynamicAddressFeedServer(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security dynamic-address feed-server "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityDynamicAddressFeedServerExists, err := checkSecurityDynamicAddressFeedServersExists(
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityDynamicAddressFeedServerExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security dynamic-address feed-server %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityDynamicAddressFeedServer(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_security_dynamic_address_feed_server")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityDynamicAddressFeedServerExists, err = checkSecurityDynamicAddressFeedServersExists(
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityDynamicAddressFeedServerExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security dynamic-address feed-server %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityDynamicAddressFeedServerReadWJunSess(d, junSess)...)
}

func resourceSecurityDynamicAddressFeedServerRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityDynamicAddressFeedServerReadWJunSess(d, junSess)
}

func resourceSecurityDynamicAddressFeedServerReadWJunSess(
	d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	dynamicAddressFeedServerOptions, err := readSecurityDynamicAddressFeedServer(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if dynamicAddressFeedServerOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityDynamicAddressFeedServerData(d, dynamicAddressFeedServerOptions)
	}

	return nil
}

func resourceSecurityDynamicAddressFeedServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityDynamicAddressFeedServer(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityDynamicAddressFeedServer(d, junSess); err != nil {
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
	if err := delSecurityDynamicAddressFeedServer(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityDynamicAddressFeedServer(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_security_dynamic_address_feed_server")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityDynamicAddressFeedServerReadWJunSess(d, junSess)...)
}

func resourceSecurityDynamicAddressFeedServerDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityDynamicAddressFeedServer(d.Get("name").(string), junSess); err != nil {
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
	if err := delSecurityDynamicAddressFeedServer(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_security_dynamic_address_feed_server")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityDynamicAddressFeedServerImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	securityDynamicAddressFeedServerExists, err := checkSecurityDynamicAddressFeedServersExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !securityDynamicAddressFeedServerExists {
		return nil, fmt.Errorf("security dynamic-address feed-server with id '%v' (id must be <name>)", d.Id())
	}
	dynamicAddressFeedServerOptions, err := readSecurityDynamicAddressFeedServer(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityDynamicAddressFeedServerData(d, dynamicAddressFeedServerOptions)

	result[0] = d

	return result, nil
}

func checkSecurityDynamicAddressFeedServersExists(name string, junSess *junos.Session,
) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security dynamic-address feed-server " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityDynamicAddressFeedServer(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security dynamic-address feed-server " + d.Get("name").(string) + " "

	if v := d.Get("hostname").(string); v != "" {
		configSet = append(configSet, setPrefix+"hostname \""+v+"\"")
	}
	if v := d.Get("url").(string); v != "" {
		configSet = append(configSet, setPrefix+"url \""+v+"\"")
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	feedNameList := make([]string, 0)
	for _, fn := range d.Get("feed_name").([]interface{}) {
		feedName := fn.(map[string]interface{})
		if slices.Contains(feedNameList, feedName["name"].(string)) {
			return fmt.Errorf("multiple blocks feed_name with the same name %s", feedName["name"].(string))
		}
		feedNameList = append(feedNameList, feedName["name"].(string))
		setPrefixFeedName := setPrefix + "feed-name " + feedName["name"].(string) + " "
		configSet = append(configSet, setPrefixFeedName)
		configSet = append(configSet, setPrefixFeedName+"path \""+feedName["path"].(string)+"\"")
		if v := feedName["description"].(string); v != "" {
			configSet = append(configSet, setPrefixFeedName+"description \""+v+"\"")
		}
		if v := feedName["hold_interval"].(int); v != -1 {
			configSet = append(configSet, setPrefixFeedName+"hold-interval "+strconv.Itoa(v))
		}
		if v := feedName["update_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixFeedName+"update-interval "+strconv.Itoa(v))
		}
	}
	if v := d.Get("hold_interval").(int); v != -1 {
		configSet = append(configSet, setPrefix+"hold-interval "+strconv.Itoa(v))
	}
	if v := d.Get("tls_profile").(string); v != "" {
		configSet = append(configSet, setPrefix+"tls-profile \""+v+"\"")
	}
	if v := d.Get("update_interval").(int); v != 0 {
		configSet = append(configSet, setPrefix+"update-interval "+strconv.Itoa(v))
	}
	if d.Get("validate_certificate_attributes_subject_or_san").(bool) {
		configSet = append(configSet, setPrefix+"validate-certificate-attributes subject-or-subject-alternative-names")
	}

	return junSess.ConfigSet(configSet)
}

func readSecurityDynamicAddressFeedServer(name string, junSess *junos.Session,
) (confRead dynamicAddressFeedServerOptions, err error) {
	// default -1
	confRead.holdInterval = -1
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security dynamic-address feed-server " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "hostname "):
				confRead.hostname = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "url "):
				confRead.url = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "feed-name "):
				itemTrimFields := strings.Split(itemTrim, " ")
				feedName := map[string]interface{}{
					"name":            itemTrimFields[0],
					"path":            "",
					"description":     "",
					"hold_interval":   -1,
					"update_interval": 0,
				}
				confRead.feedName = copyAndRemoveItemMapList("name", feedName, confRead.feedName)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "path "):
					feedName["path"] = strings.Trim(itemTrim, "\"")
				case balt.CutPrefixInString(&itemTrim, "description "):
					feedName["description"] = strings.Trim(itemTrim, "\"")
				case balt.CutPrefixInString(&itemTrim, "hold-interval "):
					feedName["hold_interval"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "update-interval "):
					feedName["update_interval"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
				confRead.feedName = append(confRead.feedName, feedName)
			case balt.CutPrefixInString(&itemTrim, "hold-interval "):
				confRead.holdInterval, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "tls-profile "):
				confRead.tlsProfile = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "update-interval "):
				confRead.updateInterval, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "validate-certificate-attributes subject-or-subject-alternative-names":
				confRead.validateCertAttrSubOrSan = true
			}
		}
	}

	return confRead, nil
}

func delSecurityDynamicAddressFeedServer(name string, junSess *junos.Session) error {
	configSet := []string{"delete security dynamic-address feed-server " + name}

	return junSess.ConfigSet(configSet)
}

func fillSecurityDynamicAddressFeedServerData(
	d *schema.ResourceData, dynamicAddressFeedServerOptions dynamicAddressFeedServerOptions,
) {
	if tfErr := d.Set("name", dynamicAddressFeedServerOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("hostname", dynamicAddressFeedServerOptions.hostname); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("url", dynamicAddressFeedServerOptions.url); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", dynamicAddressFeedServerOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("feed_name", dynamicAddressFeedServerOptions.feedName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("hold_interval", dynamicAddressFeedServerOptions.holdInterval); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("tls_profile", dynamicAddressFeedServerOptions.tlsProfile); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("update_interval", dynamicAddressFeedServerOptions.updateInterval); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"validate_certificate_attributes_subject_or_san",
		dynamicAddressFeedServerOptions.validateCertAttrSubOrSan,
	); tfErr != nil {
		panic(tfErr)
	}
}
