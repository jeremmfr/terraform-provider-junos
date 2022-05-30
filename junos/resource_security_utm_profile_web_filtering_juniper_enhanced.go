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

type utmProfileWebFilteringEnhancedOptions struct {
	noSafeSearch            bool
	timeout                 int
	name                    string
	defaultAction           string
	customBlockMessage      string
	quarantineCustomMessage string
	blockMessage            []map[string]interface{}
	category                []map[string]interface{}
	fallbackSettings        []map[string]interface{}
	quarantineMessage       []map[string]interface{}
	siteReputationAction    []map[string]interface{}
}

func resourceSecurityUtmProfileWebFilteringEnhanced() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityUtmProfileWebFilteringEnhancedCreate,
		ReadWithoutTimeout:   resourceSecurityUtmProfileWebFilteringEnhancedRead,
		UpdateWithoutTimeout: resourceSecurityUtmProfileWebFilteringEnhancedUpdate,
		DeleteWithoutTimeout: resourceSecurityUtmProfileWebFilteringEnhancedDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityUtmProfileWebFilteringEnhancedImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"block_message": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type_custom_redirect_url": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"category": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 128, formatDefault),
						},
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"block", "log-and-permit", permitW, "quarantine"}, false),
						},
						"reputation_action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"site_reputation": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"fairly-safe", "harmful", "moderately-safe", "suspicious", "very-safe",
										}, false),
									},
									"action": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"block", "log-and-permit", permitW, "quarantine"}, false),
									},
								},
							},
						},
					},
				},
			},
			"custom_block_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"block", "log-and-permit", permitW, "quarantine"}, false),
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
			"no_safe_search": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"quarantine_custom_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quarantine_message": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type_custom_redirect_url": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"site_reputation_action": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_reputation": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"fairly-safe", "harmful", "moderately-safe", "suspicious", "very-safe",
							}, false),
						},
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"block", "log-and-permit", permitW, "quarantine"}, false),
						},
					},
				},
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 1800),
			},
		},
	}
}

func resourceSecurityUtmProfileWebFilteringEnhancedCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setUtmProfileWebFEnhanced(d, clt, nil); err != nil {
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
		return diag.FromErr(fmt.Errorf("security utm feature-profile web-filtering juniper-enhanced "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	utmProfileWebFEnhancedExists, err := checkUtmProfileWebFEnhancedExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmProfileWebFEnhancedExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("security utm feature-profile web-filtering juniper-enhanced "+
			"%v already exists", d.Get("name").(string)))...)
	}

	if err := setUtmProfileWebFEnhanced(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_security_utm_profile_web_filtering_juniper_enhanced", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	utmProfileWebFEnhancedExists, err = checkUtmProfileWebFEnhancedExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmProfileWebFEnhancedExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security utm feature-profile web-filtering juniper-enhanced %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityUtmProfileWebFilteringEnhancedReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityUtmProfileWebFilteringEnhancedRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSecurityUtmProfileWebFilteringEnhancedReadWJunSess(d, clt, junSess)
}

func resourceSecurityUtmProfileWebFilteringEnhancedReadWJunSess(
	d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	utmProfileWebFEnhancedOptions, err := readUtmProfileWebFEnhanced(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if utmProfileWebFEnhancedOptions.name == "" {
		d.SetId("")
	} else {
		fillUtmProfileWebFEnhancedData(d, utmProfileWebFEnhancedOptions)
	}

	return nil
}

func resourceSecurityUtmProfileWebFilteringEnhancedUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delUtmProfileWebFEnhanced(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setUtmProfileWebFEnhanced(d, clt, nil); err != nil {
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
	if err := delUtmProfileWebFEnhanced(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setUtmProfileWebFEnhanced(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_security_utm_profile_web_filtering_juniper_enhanced", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityUtmProfileWebFilteringEnhancedReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityUtmProfileWebFilteringEnhancedDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delUtmProfileWebFEnhanced(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delUtmProfileWebFEnhanced(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_security_utm_profile_web_filtering_juniper_enhanced", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityUtmProfileWebFilteringEnhancedImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	utmProfileWebFEnhancedExists, err := checkUtmProfileWebFEnhancedExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !utmProfileWebFEnhancedExists {
		return nil, fmt.Errorf("don't find security utm feature-profile web-filtering juniper-enhanced with id "+
			"'%v' (id must be <name>)", d.Id())
	}
	utmProfileWebFEnhancedOptions, err := readUtmProfileWebFEnhanced(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillUtmProfileWebFEnhancedData(d, utmProfileWebFEnhancedOptions)

	result[0] = d

	return result, nil
}

func checkUtmProfileWebFEnhancedExists(profile string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+
		"security utm feature-profile web-filtering juniper-enhanced profile \""+profile+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setUtmProfileWebFEnhanced(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := "set security utm feature-profile web-filtering juniper-enhanced " +
		"profile \"" + d.Get("name").(string) + "\" "
	for _, v := range d.Get("block_message").([]interface{}) {
		if v != nil {
			message := v.(map[string]interface{})
			if message["url"].(string) != "" {
				configSet = append(configSet, setPrefix+"block-message url \""+message["url"].(string)+"\"")
			}
			if message["type_custom_redirect_url"].(bool) {
				configSet = append(configSet, setPrefix+"block-message type custom-redirect-url")
			}
		} else {
			configSet = append(configSet, setPrefix+"block-message")
		}
	}
	categoryNameList := make([]string, 0)
	for _, v := range d.Get("category").([]interface{}) {
		category := v.(map[string]interface{})
		if bchk.StringInSlice(category["name"].(string), categoryNameList) {
			return fmt.Errorf("multiple blocks category with the same name %s", category["name"].(string))
		}
		categoryNameList = append(categoryNameList, category["name"].(string))
		setPrefixCategory := setPrefix + "category \"" + category["name"].(string) + "\" "
		configSet = append(configSet, setPrefixCategory+"action "+category["action"].(string))
		reputationActionSiteList := make([]string, 0)
		for _, r := range category["reputation_action"].([]interface{}) {
			reputation := r.(map[string]interface{})
			if bchk.StringInSlice(reputation["site_reputation"].(string), reputationActionSiteList) {
				return fmt.Errorf("multiple blocks reputation_action with the same site_reputation %s",
					reputation["site_reputation"].(string))
			}
			reputationActionSiteList = append(reputationActionSiteList, reputation["site_reputation"].(string))
			configSet = append(configSet, setPrefixCategory+"reputation-action "+
				reputation["site_reputation"].(string)+" "+reputation["action"].(string))
		}
	}
	if d.Get("custom_block_message").(string) != "" {
		configSet = append(configSet, setPrefix+"custom-block-message \""+d.Get("custom_block_message").(string)+"\"")
	}
	if d.Get("default_action").(string) != "" {
		configSet = append(configSet, setPrefix+"default "+d.Get("default_action").(string))
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
	if d.Get("no_safe_search").(bool) {
		configSet = append(configSet, setPrefix+"no-safe-search")
	}
	if d.Get("quarantine_custom_message").(string) != "" {
		configSet = append(configSet,
			setPrefix+"quarantine-custom-message \""+d.Get("quarantine_custom_message").(string)+"\"")
	}
	for _, v := range d.Get("quarantine_message").([]interface{}) {
		if v != nil {
			message := v.(map[string]interface{})
			if message["url"].(string) != "" {
				configSet = append(configSet, setPrefix+"quarantine-message url \""+message["url"].(string)+"\"")
			}
			if message["type_custom_redirect_url"].(bool) {
				configSet = append(configSet, setPrefix+"quarantine-message type custom-redirect-url")
			}
		} else {
			configSet = append(configSet, setPrefix+"quarantine-message")
		}
	}
	siteReputationNameList := make([]string, 0)
	for _, v := range d.Get("site_reputation_action").([]interface{}) {
		siteReputation := v.(map[string]interface{})
		if bchk.StringInSlice(siteReputation["site_reputation"].(string), siteReputationNameList) {
			return fmt.Errorf("multiple blocks site_reputation_action with the same site_reputation %s",
				siteReputation["site_reputation"].(string))
		}
		siteReputationNameList = append(siteReputationNameList, siteReputation["site_reputation"].(string))
		configSet = append(configSet, setPrefix+"site-reputation-action "+
			siteReputation["site_reputation"].(string)+" "+siteReputation["action"].(string))
	}
	if d.Get("timeout").(int) != 0 {
		configSet = append(configSet, setPrefix+"timeout "+strconv.Itoa(d.Get("timeout").(int)))
	}

	return clt.configSet(configSet, junSess)
}

func readUtmProfileWebFEnhanced(profile string, clt *Client, junSess *junosSession,
) (utmProfileWebFilteringEnhancedOptions, error) {
	var confRead utmProfileWebFilteringEnhancedOptions

	showConfig, err := clt.command(cmdShowConfig+
		"security utm feature-profile web-filtering juniper-enhanced"+
		" profile \""+profile+"\""+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = profile
		categoryList := make([]map[string]interface{}, 0)
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "block-message"):
				if len(confRead.blockMessage) == 0 {
					confRead.blockMessage = append(confRead.blockMessage, map[string]interface{}{
						"url":                      "",
						"type_custom_redirect_url": false,
					})
				}
				switch {
				case itemTrim == "block-message type custom-redirect-url":
					confRead.blockMessage[0]["type_custom_redirect_url"] = true
				case strings.HasPrefix(itemTrim, "block-message url "):
					confRead.blockMessage[0]["url"] = strings.Trim(strings.TrimPrefix(itemTrim, "block-message url "), "\"")
				}
			case strings.HasPrefix(itemTrim, "category "):
				catergoryLineCut := strings.Split(itemTrim, " ")
				c := map[string]interface{}{
					"name":              catergoryLineCut[1],
					"action":            "",
					"reputation_action": make([]map[string]interface{}, 0),
				}
				categoryList = copyAndRemoveItemMapList("name", c, categoryList)
				itemTrimCategory := strings.TrimPrefix(itemTrim, "category "+catergoryLineCut[1]+" ")
				switch {
				case strings.HasPrefix(itemTrimCategory, "action "):
					c["action"] = strings.TrimPrefix(itemTrimCategory, "action ")
				case strings.HasPrefix(itemTrimCategory, "reputation-action "):
					cutReputationAction := strings.Split(strings.TrimPrefix(itemTrimCategory, "reputation-action "), " ")
					c["reputation_action"] = append(c["reputation_action"].([]map[string]interface{}), map[string]interface{}{
						"site_reputation": cutReputationAction[0],
						"action":          cutReputationAction[1],
					})
				}
				categoryList = append(categoryList, c)
			case strings.HasPrefix(itemTrim, "custom-block-message "):
				confRead.customBlockMessage = strings.Trim(strings.TrimPrefix(itemTrim, "custom-block-message "), "\"")
			case strings.HasPrefix(itemTrim, "default "):
				confRead.defaultAction = strings.TrimPrefix(itemTrim, "default ")
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
			case itemTrim == "no-safe-search":
				confRead.noSafeSearch = true
			case strings.HasPrefix(itemTrim, "quarantine-custom-message "):
				confRead.quarantineCustomMessage = strings.Trim(strings.TrimPrefix(itemTrim, "quarantine-custom-message "), "\"")
			case strings.HasPrefix(itemTrim, "quarantine-message"):
				if len(confRead.quarantineMessage) == 0 {
					confRead.quarantineMessage = append(confRead.quarantineMessage, map[string]interface{}{
						"url":                      "",
						"type_custom_redirect_url": false,
					})
				}
				switch {
				case itemTrim == "quarantine-message type custom-redirect-url":
					confRead.quarantineMessage[0]["type_custom_redirect_url"] = true
				case strings.HasPrefix(itemTrim, "quarantine-message url "):
					confRead.quarantineMessage[0]["url"] = strings.Trim(strings.TrimPrefix(itemTrim, "quarantine-message url "), "\"")
				}
			case strings.HasPrefix(itemTrim, "site-reputation-action "):
				itemTrimSiteReput := strings.Split(strings.TrimPrefix(itemTrim, "site-reputation-action "), " ")
				confRead.siteReputationAction = append(confRead.siteReputationAction, map[string]interface{}{
					"site_reputation": itemTrimSiteReput[0],
					"action":          itemTrimSiteReput[1],
				})
			case strings.HasPrefix(itemTrim, "timeout "):
				var err error
				confRead.timeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "timeout "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
		confRead.category = categoryList
	}

	return confRead, nil
}

func delUtmProfileWebFEnhanced(profile string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security utm feature-profile web-filtering juniper-enhanced "+
		"profile \""+profile+"\"")

	return clt.configSet(configSet, junSess)
}

func fillUtmProfileWebFEnhancedData(
	d *schema.ResourceData, utmProfileWebFEnhancedOptions utmProfileWebFilteringEnhancedOptions,
) {
	if tfErr := d.Set("name", utmProfileWebFEnhancedOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("block_message", utmProfileWebFEnhancedOptions.blockMessage); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("category", utmProfileWebFEnhancedOptions.category); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("custom_block_message", utmProfileWebFEnhancedOptions.customBlockMessage); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("default_action", utmProfileWebFEnhancedOptions.defaultAction); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("fallback_settings", utmProfileWebFEnhancedOptions.fallbackSettings); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_safe_search", utmProfileWebFEnhancedOptions.noSafeSearch); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("quarantine_custom_message", utmProfileWebFEnhancedOptions.quarantineCustomMessage); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("quarantine_message", utmProfileWebFEnhancedOptions.quarantineMessage); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("site_reputation_action", utmProfileWebFEnhancedOptions.siteReputationAction); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("timeout", utmProfileWebFEnhancedOptions.timeout); tfErr != nil {
		panic(tfErr)
	}
}
