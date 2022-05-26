package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
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
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSecurityDynamicAddressFeedServer(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security dynamic-address feed-server "+
			"not compatible with Junos device %s", jnprSess.SystemInformation.HardwareModel))
	}
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityDynamicAddressFeedServerExists, err := checkSecurityDynamicAddressFeedServersExists(
		d.Get("name").(string),
		sess, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityDynamicAddressFeedServerExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security dynamic-address feed-server %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityDynamicAddressFeedServer(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_dynamic_address_feed_server", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityDynamicAddressFeedServerExists, err = checkSecurityDynamicAddressFeedServersExists(
		d.Get("name").(string),
		sess, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityDynamicAddressFeedServerExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security dynamic-address feed-server %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityDynamicAddressFeedServerReadWJnprSess(d, sess, jnprSess)...)
}

func resourceSecurityDynamicAddressFeedServerRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityDynamicAddressFeedServerReadWJnprSess(d, sess, jnprSess)
}

func resourceSecurityDynamicAddressFeedServerReadWJnprSess(
	d *schema.ResourceData, sess *Session, jnprSess *NetconfObject,
) diag.Diagnostics {
	mutex.Lock()
	dynamicAddressFeedServerOptions, err := readSecurityDynamicAddressFeedServer(d.Get("name").(string), sess, jnprSess)
	mutex.Unlock()
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
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delSecurityDynamicAddressFeedServer(d.Get("name").(string), sess, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityDynamicAddressFeedServer(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityDynamicAddressFeedServer(d.Get("name").(string), sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityDynamicAddressFeedServer(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_dynamic_address_feed_server", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityDynamicAddressFeedServerReadWJnprSess(d, sess, jnprSess)...)
}

func resourceSecurityDynamicAddressFeedServerDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delSecurityDynamicAddressFeedServer(d.Get("name").(string), sess, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityDynamicAddressFeedServer(d.Get("name").(string), sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_dynamic_address_feed_server", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityDynamicAddressFeedServerImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	securityDynamicAddressFeedServerExists, err := checkSecurityDynamicAddressFeedServersExists(d.Id(), sess, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityDynamicAddressFeedServerExists {
		return nil, fmt.Errorf("security dynamic-address feed-server with id '%v' (id must be <name>)", d.Id())
	}
	dynamicAddressFeedServerOptions, err := readSecurityDynamicAddressFeedServer(d.Id(), sess, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityDynamicAddressFeedServerData(d, dynamicAddressFeedServerOptions)

	result[0] = d

	return result, nil
}

func checkSecurityDynamicAddressFeedServersExists(name string, sess *Session, jnprSess *NetconfObject) (bool, error) {
	showConfig, err := sess.command(cmdShowConfig+
		"security dynamic-address feed-server "+name+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityDynamicAddressFeedServer(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject) error {
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
		if bchk.StringInSlice(feedName["name"].(string), feedNameList) {
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

	return sess.configSet(configSet, jnprSess)
}

func readSecurityDynamicAddressFeedServer(name string, sess *Session, jnprSess *NetconfObject,
) (dynamicAddressFeedServerOptions, error) {
	var confRead dynamicAddressFeedServerOptions
	// default -1
	confRead.holdInterval = -1

	showConfig, err := sess.command(cmdShowConfig+
		"security dynamic-address feed-server "+name+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "hostname "):
				confRead.hostname = strings.Trim(strings.TrimPrefix(itemTrim, "hostname "), "\"")
			case strings.HasPrefix(itemTrim, "url "):
				confRead.url = strings.Trim(strings.TrimPrefix(itemTrim, "url "), "\"")
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "feed-name "):
				itemTrimFeedSplit := strings.Split(strings.TrimPrefix(itemTrim, "feed-name "), " ")
				feedName := map[string]interface{}{
					"name":            itemTrimFeedSplit[0],
					"path":            "",
					"description":     "",
					"hold_interval":   -1,
					"update_interval": 0,
				}
				confRead.feedName = copyAndRemoveItemMapList("name", feedName, confRead.feedName)
				itemTrimFeedName := strings.TrimPrefix(itemTrim, "feed-name "+itemTrimFeedSplit[0]+" ")
				switch {
				case strings.HasPrefix(itemTrimFeedName, "path "):
					feedName["path"] = strings.Trim(strings.TrimPrefix(itemTrimFeedName, "path "), "\"")
				case strings.HasPrefix(itemTrimFeedName, "description "):
					feedName["description"] = strings.Trim(strings.TrimPrefix(itemTrimFeedName, "description "), "\"")
				case strings.HasPrefix(itemTrimFeedName, "hold-interval "):
					var err error
					feedName["hold_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrimFeedName, "hold-interval "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrimFeedName, "update-interval "):
					var err error
					feedName["update_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrimFeedName, "update-interval "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
				confRead.feedName = append(confRead.feedName, feedName)
			case strings.HasPrefix(itemTrim, "hold-interval "):
				var err error
				confRead.holdInterval, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "hold-interval "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "tls-profile "):
				confRead.tlsProfile = strings.Trim(strings.TrimPrefix(itemTrim, "tls-profile "), "\"")
			case strings.HasPrefix(itemTrim, "update-interval "):
				var err error
				confRead.updateInterval, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "update-interval "))
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

func delSecurityDynamicAddressFeedServer(name string, sess *Session, jnprSess *NetconfObject) error {
	configSet := []string{"delete security dynamic-address feed-server " + name}

	return sess.configSet(configSet, jnprSess)
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
