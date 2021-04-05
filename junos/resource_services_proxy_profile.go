package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type proxyProfileOptions struct {
	protocolHTTPPort int
	name             string
	protocolHTTPHost string
}

func resourceServicesProxyProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServicesProxyProfileCreate,
		ReadContext:   resourceServicesProxyProfileRead,
		UpdateContext: resourceServicesProxyProfileUpdate,
		DeleteContext: resourceServicesProxyProfileDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServicesProxyProfileImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"protocol_http_host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol_http_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
		},
	}
}

func resourceServicesProxyProfileCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setServicesProxyProfile(d, m, nil); err != nil {
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
	sess.configLock(jnprSess)
	proxyProfileExists, err := checkServicesProxyProfileExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if proxyProfileExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("services proxy profile %v already exists", d.Get("name").(string)))
	}

	if err := setServicesProxyProfile(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_services_proxy_profile", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	proxyProfileExists, err = checkServicesProxyProfileExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if proxyProfileExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("services proxy profile %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesProxyProfileReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesProxyProfileRead(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceServicesProxyProfileReadWJnprSess(d, m, jnprSess)
}

func resourceServicesProxyProfileReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	proxyProfileOptions, err := readServicesProxyProfile(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if proxyProfileOptions.name == "" {
		d.SetId("")
	} else {
		fillServicesProxyProfileData(d, proxyProfileOptions)
	}

	return nil
}

func resourceServicesProxyProfileUpdate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delServicesProxyProfile(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setServicesProxyProfile(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_services_proxy_profile", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesProxyProfileReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesProxyProfileDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delServicesProxyProfile(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_services_proxy_profile", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesProxyProfileImport(
	d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	proxyProfileExists, err := checkServicesProxyProfileExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !proxyProfileExists {
		return nil, fmt.Errorf("don't find services proxy profile with id '%v' (id must be <name>)", d.Id())
	}
	proxyProfileOptions, err := readServicesProxyProfile(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillServicesProxyProfileData(d, proxyProfileOptions)

	result[0] = d

	return result, nil
}

func checkServicesProxyProfileExists(profile string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	profileConfig, err := sess.command("show configuration services proxy profile \""+
		profile+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if profileConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setServicesProxyProfile(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set services proxy profile \"" + d.Get("name").(string) + "\" "
	configSet = append(configSet, setPrefix+"protocol http host "+d.Get("protocol_http_host").(string))
	if d.Get("protocol_http_port").(int) != 0 {
		configSet = append(configSet, setPrefix+"protocol http port "+strconv.Itoa(d.Get("protocol_http_port").(int)))
	}

	return sess.configSet(configSet, jnprSess)
}

func readServicesProxyProfile(profile string, m interface{}, jnprSess *NetconfObject) (
	proxyProfileOptions, error) {
	sess := m.(*Session)
	var confRead proxyProfileOptions

	profileConfig, err := sess.command("show configuration"+
		" services proxy profile \""+profile+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if profileConfig != emptyWord {
		confRead.name = profile
		for _, item := range strings.Split(profileConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "protocol http host "):
				confRead.protocolHTTPHost = strings.TrimPrefix(itemTrim, "protocol http host ")
			case strings.HasPrefix(itemTrim, "protocol http port "):
				var err error
				confRead.protocolHTTPPort, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "protocol http port "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delServicesProxyProfile(profile string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete services proxy profile \"" + profile + "\""}

	return sess.configSet(configSet, jnprSess)
}

func fillServicesProxyProfileData(
	d *schema.ResourceData, proxyProfileOptions proxyProfileOptions) {
	if tfErr := d.Set("name", proxyProfileOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("protocol_http_host", proxyProfileOptions.protocolHTTPHost); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("protocol_http_port", proxyProfileOptions.protocolHTTPPort); tfErr != nil {
		panic(tfErr)
	}
}
