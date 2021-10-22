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
	holdInterval   int
	updateInterval int
	description    string
	hostname       string
	name           string
	feedName       []map[string]interface{}
}

func resourceSecurityDynamicAddressFeedServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityDynamicAddressFeedServerCreate,
		ReadContext:   resourceSecurityDynamicAddressFeedServerRead,
		UpdateContext: resourceSecurityDynamicAddressFeedServerUpdate,
		DeleteContext: resourceSecurityDynamicAddressFeedServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityDynamicAddressFeedServerImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 16, formatDefault),
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
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
			"update_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(30, 4294967295),
			},
		},
	}
}

func resourceSecurityDynamicAddressFeedServerCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSecurityDynamicAddressFeedServer(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security dynamic-address feed-server "+
			"not compatible with Junos device %s", jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	securityDynamicAddressFeedServerExists, err :=
		checkSecurityDynamicAddressFeedServersExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityDynamicAddressFeedServerExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security dynamic-address feed-server %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityDynamicAddressFeedServer(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_dynamic_address_feed_server", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityDynamicAddressFeedServerExists, err =
		checkSecurityDynamicAddressFeedServersExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityDynamicAddressFeedServerExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security dynamic-address feed-server %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityDynamicAddressFeedServerReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityDynamicAddressFeedServerRead(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityDynamicAddressFeedServerReadWJnprSess(d, m, jnprSess)
}

func resourceSecurityDynamicAddressFeedServerReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	dynamicAddressFeedServerOptions, err := readSecurityDynamicAddressFeedServer(d.Get("name").(string), m, jnprSess)
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

func resourceSecurityDynamicAddressFeedServerUpdate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSecurityDynamicAddressFeedServer(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityDynamicAddressFeedServer(d, m, jnprSess); err != nil {
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

	return append(diagWarns, resourceSecurityDynamicAddressFeedServerReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityDynamicAddressFeedServerDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSecurityDynamicAddressFeedServer(d.Get("name").(string), m, jnprSess); err != nil {
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

func resourceSecurityDynamicAddressFeedServerImport(
	d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	securityDynamicAddressFeedServerExists, err := checkSecurityDynamicAddressFeedServersExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityDynamicAddressFeedServerExists {
		return nil, fmt.Errorf("security dynamic-address feed-server with id '%v' (id must be <name>)", d.Id())
	}
	dynamicAddressFeedServerOptions, err := readSecurityDynamicAddressFeedServer(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityDynamicAddressFeedServerData(d, dynamicAddressFeedServerOptions)

	result[0] = d

	return result, nil
}

func checkSecurityDynamicAddressFeedServersExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration"+
		" security dynamic-address feed-server "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setSecurityDynamicAddressFeedServer(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security dynamic-address feed-server " + d.Get("name").(string) + " "

	configSet = append(configSet, setPrefix+"hostname \""+d.Get("hostname").(string)+"\"")
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
	if v := d.Get("update_interval").(int); v != 0 {
		configSet = append(configSet, setPrefix+"update-interval "+strconv.Itoa(v))
	}

	return sess.configSet(configSet, jnprSess)
}

func readSecurityDynamicAddressFeedServer(
	name string, m interface{}, jnprSess *NetconfObject) (dynamicAddressFeedServerOptions, error) {
	sess := m.(*Session)
	var confRead dynamicAddressFeedServerOptions
	// default -1
	confRead.holdInterval = -1

	showConfig, err := sess.command("show configuration"+
		" security dynamic-address feed-server "+name+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "hostname "):
				confRead.hostname = strings.Trim(strings.TrimPrefix(itemTrim, "hostname "), "\"")
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
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrimFeedName, "update-interval "):
					var err error
					feedName["update_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrimFeedName, "update-interval "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				}
				confRead.feedName = append(confRead.feedName, feedName)
			case strings.HasPrefix(itemTrim, "hold-interval "):
				var err error
				confRead.holdInterval, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "hold-interval "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "update-interval "):
				var err error
				confRead.updateInterval, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "update-interval "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delSecurityDynamicAddressFeedServer(name string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete security dynamic-address feed-server " + name}

	return sess.configSet(configSet, jnprSess)
}

func fillSecurityDynamicAddressFeedServerData(
	d *schema.ResourceData, dynamicAddressFeedServerOptions dynamicAddressFeedServerOptions) {
	if tfErr := d.Set("name", dynamicAddressFeedServerOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("hostname", dynamicAddressFeedServerOptions.hostname); tfErr != nil {
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
	if tfErr := d.Set("update_interval", dynamicAddressFeedServerOptions.updateInterval); tfErr != nil {
		panic(tfErr)
	}
}
