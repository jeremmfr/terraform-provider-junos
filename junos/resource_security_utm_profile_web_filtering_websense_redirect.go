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

type utmProfileWebFilteringWebsenseOptions struct {
	sockets            int
	timeout            int
	name               string
	account            string
	customBlockMessage string
	fallbackSettings   []map[string]interface{}
	server             []map[string]interface{}
}

func resourceSecurityUtmProfileWebFilteringWebsense() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityUtmProfileWebFilteringWebsenseCreate,
		ReadWithoutTimeout:   resourceSecurityUtmProfileWebFilteringWebsenseRead,
		UpdateWithoutTimeout: resourceSecurityUtmProfileWebFilteringWebsenseUpdate,
		DeleteWithoutTimeout: resourceSecurityUtmProfileWebFilteringWebsenseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityUtmProfileWebFilteringWebsenseImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"account": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"custom_block_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fallback_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"block", "log-and-permit"}, false),
						},
						"server_connectivity": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"block", "log-and-permit"}, false),
						},
						"timeout": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"block", "log-and-permit"}, false),
						},
						"too_many_requests": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"block", "log-and-permit"}, false),
						},
					},
				},
			},
			"server": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1024, 65535),
						},
					},
				},
			},
			"sockets": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 32),
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 1800),
			},
		},
	}
}

func resourceSecurityUtmProfileWebFilteringWebsenseCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setUtmProfileWebFWebsense(d, clt, nil); err != nil {
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
		return diag.FromErr(fmt.Errorf("security utm feature-profile web-filtering websense-redirect "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	utmProfileWebFWebsenseExists, err := checkUtmProfileWebFWebsenseExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmProfileWebFWebsenseExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("security utm feature-profile web-filtering websense-redirect "+
			"%v already exists", d.Get("name").(string)))...)
	}

	if err := setUtmProfileWebFWebsense(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_security_utm_profile_web_filtering_websense_redirect", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	utmProfileWebFWebsenseExists, err = checkUtmProfileWebFWebsenseExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmProfileWebFWebsenseExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security utm feature-profile web-filtering websense-redirect %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityUtmProfileWebFilteringWebsenseReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityUtmProfileWebFilteringWebsenseRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSecurityUtmProfileWebFilteringWebsenseReadWJunSess(d, clt, junSess)
}

func resourceSecurityUtmProfileWebFilteringWebsenseReadWJunSess(
	d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	utmProfileWebFWebsenseOptions, err := readUtmProfileWebFWebsense(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if utmProfileWebFWebsenseOptions.name == "" {
		d.SetId("")
	} else {
		fillUtmProfileWebFWebsenseData(d, utmProfileWebFWebsenseOptions)
	}

	return nil
}

func resourceSecurityUtmProfileWebFilteringWebsenseUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delUtmProfileWebFWebsense(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setUtmProfileWebFWebsense(d, clt, nil); err != nil {
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
	if err := delUtmProfileWebFWebsense(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setUtmProfileWebFWebsense(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_security_utm_profile_web_filtering_websense_redirect", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityUtmProfileWebFilteringWebsenseReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityUtmProfileWebFilteringWebsenseDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delUtmProfileWebFWebsense(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delUtmProfileWebFWebsense(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_security_utm_profile_web_filtering_websense_redirect", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityUtmProfileWebFilteringWebsenseImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	utmProfileWebFWebsenseExists, err := checkUtmProfileWebFWebsenseExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !utmProfileWebFWebsenseExists {
		return nil, fmt.Errorf("don't find security utm feature-profile web-filtering websense-redirect with id "+
			"'%v' (id must be <name>)", d.Id())
	}
	utmProfileWebFWebsenseOptions, err := readUtmProfileWebFWebsense(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillUtmProfileWebFWebsenseData(d, utmProfileWebFWebsenseOptions)

	result[0] = d

	return result, nil
}

func checkUtmProfileWebFWebsenseExists(profile string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+
		"security utm feature-profile web-filtering websense-redirect profile \""+profile+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setUtmProfileWebFWebsense(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := "set security utm feature-profile web-filtering websense-redirect " +
		"profile \"" + d.Get("name").(string) + "\" "
	if d.Get("account").(string) != "" {
		configSet = append(configSet, setPrefix+"account \""+d.Get("account").(string)+"\"")
	}
	if d.Get("custom_block_message").(string) != "" {
		configSet = append(configSet, setPrefix+"custom-block-message \""+d.Get("custom_block_message").(string)+"\"")
	}
	for _, v := range d.Get("fallback_settings").([]interface{}) {
		if v != nil {
			fSettings := v.(map[string]interface{})
			if fSettings["default"].(string) != "" {
				configSet = append(configSet, setPrefix+"fallback-settings default "+
					fSettings["default"].(string))
			}
			if fSettings["server_connectivity"].(string) != "" {
				configSet = append(configSet, setPrefix+"fallback-settings server-connectivity "+
					fSettings["server_connectivity"].(string))
			}
			if fSettings["timeout"].(string) != "" {
				configSet = append(configSet, setPrefix+"fallback-settings timeout "+
					fSettings["timeout"].(string))
			}
			if fSettings["too_many_requests"].(string) != "" {
				configSet = append(configSet, setPrefix+"fallback-settings too-many-requests "+
					fSettings["too_many_requests"].(string))
			}
		} else {
			configSet = append(configSet, setPrefix+"fallback-settings")
		}
	}
	for _, v := range d.Get("server").([]interface{}) {
		configSet = append(configSet, setPrefix+"server")
		if v != nil {
			server := v.(map[string]interface{})
			if server["host"].(string) != "" {
				configSet = append(configSet, setPrefix+"server host "+server["host"].(string))
			}
			if server["port"].(int) != 0 {
				configSet = append(configSet, setPrefix+"server port "+strconv.Itoa(server["port"].(int)))
			}
		} else {
			configSet = append(configSet, setPrefix+"server")
		}
	}
	if d.Get("sockets").(int) != 0 {
		configSet = append(configSet, setPrefix+"sockets "+strconv.Itoa(d.Get("sockets").(int)))
	}
	if d.Get("timeout").(int) != 0 {
		configSet = append(configSet, setPrefix+"timeout "+strconv.Itoa(d.Get("timeout").(int)))
	}

	return clt.configSet(configSet, junSess)
}

func readUtmProfileWebFWebsense(profile string, clt *Client, junSess *junosSession,
) (utmProfileWebFilteringWebsenseOptions, error) {
	var confRead utmProfileWebFilteringWebsenseOptions

	showConfig, err := clt.command(cmdShowConfig+
		"security utm feature-profile web-filtering websense-redirect"+
		" profile \""+profile+"\""+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = profile
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "account "):
				confRead.customBlockMessage = strings.Trim(strings.TrimPrefix(itemTrim, "account "), "\"")
			case strings.HasPrefix(itemTrim, "custom-block-message "):
				confRead.customBlockMessage = strings.Trim(strings.TrimPrefix(itemTrim, "custom-block-message "), "\"")
			case strings.HasPrefix(itemTrim, "fallback-settings"):
				if len(confRead.fallbackSettings) == 0 {
					confRead.fallbackSettings = append(confRead.fallbackSettings, map[string]interface{}{
						"default":             "",
						"server_connectivity": "",
						"timeout":             "",
						"too_many_requests":   "",
					})
				}
				itemTrimFallback := strings.TrimPrefix(itemTrim, "fallback-settings ")
				switch {
				case strings.HasPrefix(itemTrimFallback, "default "):
					confRead.fallbackSettings[0]["default"] = strings.TrimPrefix(itemTrimFallback, "default ")
				case strings.HasPrefix(itemTrimFallback, "server-connectivity "):
					confRead.fallbackSettings[0]["server_connectivity"] = strings.TrimPrefix(itemTrimFallback, "server-connectivity ")
				case strings.HasPrefix(itemTrimFallback, "timeout "):
					confRead.fallbackSettings[0]["timeout"] = strings.TrimPrefix(itemTrimFallback, "timeout ")
				case strings.HasPrefix(itemTrimFallback, "too-many-requests "):
					confRead.fallbackSettings[0]["too_many_requests"] = strings.TrimPrefix(itemTrimFallback, "too-many-requests ")
				}
			case strings.HasPrefix(itemTrim, "server"):
				if len(confRead.server) == 0 {
					confRead.server = append(confRead.server, map[string]interface{}{
						"host": "",
						"port": 0,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "server host "):
					confRead.server[0]["host"] = strings.TrimPrefix(itemTrim, "server host ")
				case strings.HasPrefix(itemTrim, "server port "):
					var err error
					confRead.server[0]["port"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "server port "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			case strings.HasPrefix(itemTrim, "sockets "):
				var err error
				confRead.sockets, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "sockets "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "timeout "):
				var err error
				confRead.timeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "timeout "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delUtmProfileWebFWebsense(profile string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security utm feature-profile web-filtering websense-redirect "+
		"profile \""+profile+"\"")

	return clt.configSet(configSet, junSess)
}

func fillUtmProfileWebFWebsenseData(
	d *schema.ResourceData, utmProfileWebFWebsenseOptions utmProfileWebFilteringWebsenseOptions,
) {
	if tfErr := d.Set("name", utmProfileWebFWebsenseOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("account", utmProfileWebFWebsenseOptions.account); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("custom_block_message", utmProfileWebFWebsenseOptions.customBlockMessage); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("fallback_settings", utmProfileWebFWebsenseOptions.fallbackSettings); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("server", utmProfileWebFWebsenseOptions.server); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("sockets", utmProfileWebFWebsenseOptions.sockets); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("timeout", utmProfileWebFWebsenseOptions.timeout); tfErr != nil {
		panic(tfErr)
	}
}
