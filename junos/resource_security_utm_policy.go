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

type utmPolicyOptions struct {
	name                     string
	antiSpamSMTPProfile      string
	webFilteringProfile      string
	antiVirus                []map[string]interface{}
	contentFiltering         []map[string]interface{}
	trafficSessionsPerClient []map[string]interface{}
}

func resourceSecurityUtmPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityUtmPolicyCreate,
		ReadWithoutTimeout:   resourceSecurityUtmPolicyRead,
		UpdateWithoutTimeout: resourceSecurityUtmPolicyUpdate,
		DeleteWithoutTimeout: resourceSecurityUtmPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityUtmPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"anti_spam_smtp_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"anti_virus": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ftp_download_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ftp_upload_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"http_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"imap_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"pop3_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"smtp_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"content_filtering": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ftp_download_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ftp_upload_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"http_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"imap_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"pop3_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"smtp_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"traffic_sessions_per_client": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 2000),
							Default:      -1,
						},
						"over_limit": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"block", "log-and-permit"}, false),
						},
					},
				},
			},
			"web_filtering_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSecurityUtmPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setUtmPolicy(d, clt, nil); err != nil {
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
		return diag.FromErr(fmt.Errorf("security utm utm-policy "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	utmPolicyExists, err := checkUtmPolicysExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmPolicyExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security utm utm-policy %v already exists", d.Get("name").(string)))...)
	}

	if err := setUtmPolicy(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_security_utm_policy", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	utmPolicyExists, err = checkUtmPolicysExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmPolicyExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security utm utm-policy %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityUtmPolicyReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityUtmPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSecurityUtmPolicyReadWJunSess(d, clt, junSess)
}

func resourceSecurityUtmPolicyReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	utmPolicyOptions, err := readUtmPolicy(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if utmPolicyOptions.name == "" {
		d.SetId("")
	} else {
		fillUtmPolicyData(d, utmPolicyOptions)
	}

	return nil
}

func resourceSecurityUtmPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delUtmPolicy(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setUtmPolicy(d, clt, nil); err != nil {
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
	if err := delUtmPolicy(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setUtmPolicy(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_security_utm_policy", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityUtmPolicyReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityUtmPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delUtmPolicy(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delUtmPolicy(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_security_utm_policy", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityUtmPolicyImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	utmPolicyExists, err := checkUtmPolicysExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !utmPolicyExists {
		return nil, fmt.Errorf("don't find security utm utm-policy with id '%v' (id must be <name>)", d.Id())
	}
	utmPolicyOptions, err := readUtmPolicy(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillUtmPolicyData(d, utmPolicyOptions)

	result[0] = d

	return result, nil
}

func checkUtmPolicysExists(policy string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+"security utm utm-policy \""+policy+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setUtmPolicy(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := "set security utm utm-policy \"" + d.Get("name").(string) + "\" "
	if d.Get("anti_spam_smtp_profile").(string) != "" {
		configSet = append(configSet, setPrefix+"anti-spam smtp-profile \""+d.Get("anti_spam_smtp_profile").(string)+"\"")
	}
	for _, v := range d.Get("anti_virus").([]interface{}) {
		if v != nil {
			antiVirus := v.(map[string]interface{})
			setPrefixAntiVirus := setPrefix + "anti-virus "
			if antiVirus["ftp_download_profile"].(string) != "" {
				configSet = append(configSet, setPrefixAntiVirus+"ftp download-profile \""+
					antiVirus["ftp_download_profile"].(string)+"\"")
			}
			if antiVirus["ftp_upload_profile"].(string) != "" {
				configSet = append(configSet, setPrefixAntiVirus+"ftp upload-profile \""+
					antiVirus["ftp_upload_profile"].(string)+"\"")
			}
			if antiVirus["http_profile"].(string) != "" {
				configSet = append(configSet, setPrefixAntiVirus+"http-profile \""+
					antiVirus["http_profile"].(string)+"\"")
			}
			if antiVirus["imap_profile"].(string) != "" {
				configSet = append(configSet, setPrefixAntiVirus+"imap-profile \""+
					antiVirus["imap_profile"].(string)+"\"")
			}
			if antiVirus["pop3_profile"].(string) != "" {
				configSet = append(configSet, setPrefixAntiVirus+"pop3-profile \""+
					antiVirus["pop3_profile"].(string)+"\"")
			}
			if antiVirus["smtp_profile"].(string) != "" {
				configSet = append(configSet, setPrefixAntiVirus+"smtp-profile \""+
					antiVirus["smtp_profile"].(string)+"\"")
			}
		} else {
			return fmt.Errorf("anti_virus block is empty")
		}
	}
	for _, v := range d.Get("content_filtering").([]interface{}) {
		if v != nil {
			contentFiltering := v.(map[string]interface{})
			setPrefixContentFiltering := setPrefix + "content-filtering "
			if contentFiltering["ftp_download_profile"].(string) != "" {
				configSet = append(configSet, setPrefixContentFiltering+"ftp download-profile \""+
					contentFiltering["ftp_download_profile"].(string)+"\"")
			}
			if contentFiltering["ftp_upload_profile"].(string) != "" {
				configSet = append(configSet, setPrefixContentFiltering+"ftp upload-profile \""+
					contentFiltering["ftp_upload_profile"].(string)+"\"")
			}
			if contentFiltering["http_profile"].(string) != "" {
				configSet = append(configSet, setPrefixContentFiltering+"http-profile \""+
					contentFiltering["http_profile"].(string)+"\"")
			}
			if contentFiltering["imap_profile"].(string) != "" {
				configSet = append(configSet, setPrefixContentFiltering+"imap-profile \""+
					contentFiltering["imap_profile"].(string)+"\"")
			}
			if contentFiltering["pop3_profile"].(string) != "" {
				configSet = append(configSet, setPrefixContentFiltering+"pop3-profile \""+
					contentFiltering["pop3_profile"].(string)+"\"")
			}
			if contentFiltering["smtp_profile"].(string) != "" {
				configSet = append(configSet, setPrefixContentFiltering+"smtp-profile \""+
					contentFiltering["smtp_profile"].(string)+"\"")
			}
		} else {
			return fmt.Errorf("content_filtering block is empty")
		}
	}
	for _, v := range d.Get("traffic_sessions_per_client").([]interface{}) {
		trafficSessPerClient := v.(map[string]interface{})
		if trafficSessPerClient["limit"].(int) != -1 {
			configSet = append(configSet, setPrefix+"traffic-options sessions-per-client limit "+
				strconv.Itoa(trafficSessPerClient["limit"].(int)))
		}
		if trafficSessPerClient["over_limit"].(string) != "" {
			configSet = append(configSet, setPrefix+"traffic-options sessions-per-client over-limit "+
				trafficSessPerClient["over_limit"].(string))
		}
		if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1],
			setPrefix+"traffic-options sessions-per-client") {
			return fmt.Errorf("traffic_sessions_per_client block is empty")
		}
	}
	if d.Get("web_filtering_profile").(string) != "" {
		configSet = append(configSet, setPrefix+"web-filtering http-profile \""+
			d.Get("web_filtering_profile").(string)+"\"")
	}

	return clt.configSet(configSet, junSess)
}

func readUtmPolicy(policy string, clt *Client, junSess *junosSession,
) (utmPolicyOptions, error) {
	var confRead utmPolicyOptions

	showConfig, err := clt.command(cmdShowConfig+
		"security utm utm-policy \""+policy+"\""+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = policy
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "anti-spam smtp-profile "):
				confRead.antiSpamSMTPProfile = strings.Trim(strings.TrimPrefix(itemTrim, "anti-spam smtp-profile "), "\"")
			case strings.HasPrefix(itemTrim, "anti-virus "):
				if len(confRead.antiVirus) == 0 {
					confRead.antiVirus = append(confRead.antiVirus, genMapUtmPolicyProfile())
				}
				readUtmPolicyProfile(strings.TrimPrefix(itemTrim, "anti-virus "), confRead.antiVirus[0])
			case strings.HasPrefix(itemTrim, "content-filtering "):
				if len(confRead.contentFiltering) == 0 {
					confRead.contentFiltering = append(confRead.contentFiltering, genMapUtmPolicyProfile())
				}
				readUtmPolicyProfile(strings.TrimPrefix(itemTrim, "content-filtering "), confRead.contentFiltering[0])
			case strings.HasPrefix(itemTrim, "traffic-options sessions-per-client "):
				if len(confRead.trafficSessionsPerClient) == 0 {
					confRead.trafficSessionsPerClient = append(confRead.trafficSessionsPerClient, map[string]interface{}{
						"limit":      -1,
						"over_limit": "",
					})
				}
				if strings.HasPrefix(itemTrim, "traffic-options sessions-per-client limit ") {
					var err error
					confRead.trafficSessionsPerClient[0]["limit"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrim, "traffic-options sessions-per-client limit "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
				if strings.HasPrefix(itemTrim, "traffic-options sessions-per-client over-limit ") {
					confRead.trafficSessionsPerClient[0]["over_limit"] = strings.TrimPrefix(
						itemTrim, "traffic-options sessions-per-client over-limit ")
				}
			case strings.HasPrefix(itemTrim, "web-filtering http-profile "):
				confRead.webFilteringProfile = strings.Trim(strings.TrimPrefix(itemTrim, "web-filtering http-profile "), "\"")
			}
		}
	}

	return confRead, nil
}

func genMapUtmPolicyProfile() map[string]interface{} {
	return map[string]interface{}{
		"ftp_download_profile": "",
		"ftp_upload_profile":   "",
		"http_profile":         "",
		"imap_profile":         "",
		"pop3_profile":         "",
		"smtp_profile":         "",
	}
}

func readUtmPolicyProfile(itemTrimPolicyProfile string, profileMap map[string]interface{}) {
	switch {
	case strings.HasPrefix(itemTrimPolicyProfile, "ftp download-profile "):
		profileMap["ftp_download_profile"] = strings.Trim(
			strings.TrimPrefix(itemTrimPolicyProfile, "ftp download-profile "), "\"")
	case strings.HasPrefix(itemTrimPolicyProfile, "ftp upload-profile "):
		profileMap["ftp_upload_profile"] = strings.Trim(
			strings.TrimPrefix(itemTrimPolicyProfile, "ftp upload-profile "), "\"")
	case strings.HasPrefix(itemTrimPolicyProfile, "http-profile "):
		profileMap["http_profile"] = strings.Trim(
			strings.TrimPrefix(itemTrimPolicyProfile, "http-profile "), "\"")
	case strings.HasPrefix(itemTrimPolicyProfile, "imap-profile "):
		profileMap["imap_profile"] = strings.Trim(
			strings.TrimPrefix(itemTrimPolicyProfile, "imap-profile "), "\"")
	case strings.HasPrefix(itemTrimPolicyProfile, "pop3-profile "):
		profileMap["pop3_profile"] = strings.Trim(
			strings.TrimPrefix(itemTrimPolicyProfile, "pop3-profile "), "\"")
	case strings.HasPrefix(itemTrimPolicyProfile, "smtp-profile "):
		profileMap["smtp_profile"] = strings.Trim(
			strings.TrimPrefix(itemTrimPolicyProfile, "smtp-profile "), "\"")
	}
}

func delUtmPolicy(policy string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security utm utm-policy \""+policy+"\"")

	return clt.configSet(configSet, junSess)
}

func fillUtmPolicyData(d *schema.ResourceData, utmPolicyOptions utmPolicyOptions) {
	if tfErr := d.Set("name", utmPolicyOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("anti_spam_smtp_profile", utmPolicyOptions.antiSpamSMTPProfile); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("anti_virus", utmPolicyOptions.antiVirus); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("content_filtering", utmPolicyOptions.contentFiltering); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("traffic_sessions_per_client", utmPolicyOptions.trafficSessionsPerClient); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("web_filtering_profile", utmPolicyOptions.webFilteringProfile); tfErr != nil {
		panic(tfErr)
	}
}
