package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type utmProfileWebFilteringLocalOptions struct {
	timeout            int
	name               string
	defaultAction      string
	customBlockMessage string
	fallbackSettings   []map[string]interface{}
}

func resourceSecurityUtmProfileWebFilteringLocal() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityUtmProfileWebFilteringLocalCreate,
		ReadWithoutTimeout:   resourceSecurityUtmProfileWebFilteringLocalRead,
		UpdateWithoutTimeout: resourceSecurityUtmProfileWebFilteringLocalUpdate,
		DeleteWithoutTimeout: resourceSecurityUtmProfileWebFilteringLocalDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityUtmProfileWebFilteringLocalImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"custom_block_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"block", "log-and-permit", "permit"}, false),
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
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 1800),
			},
		},
	}
}

func resourceSecurityUtmProfileWebFilteringLocalCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setUtmProfileWebFLocal(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security utm feature-profile web-filtering juniper-local "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	utmProfileWebFLocalExists, err := checkUtmProfileWebFLocalExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmProfileWebFLocalExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("security utm feature-profile web-filtering juniper-local "+
			"%v already exists", d.Get("name").(string)))...)
	}

	if err := setUtmProfileWebFLocal(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_security_utm_profile_web_filtering_juniper_local")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	utmProfileWebFLocalExists, err = checkUtmProfileWebFLocalExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmProfileWebFLocalExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security utm feature-profile web-filtering juniper-local %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityUtmProfileWebFilteringLocalReadWJunSess(d, junSess)...)
}

func resourceSecurityUtmProfileWebFilteringLocalRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityUtmProfileWebFilteringLocalReadWJunSess(d, junSess)
}

func resourceSecurityUtmProfileWebFilteringLocalReadWJunSess(
	d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	utmProfileWebFLocalOptions, err := readUtmProfileWebFLocal(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if utmProfileWebFLocalOptions.name == "" {
		d.SetId("")
	} else {
		fillUtmProfileWebFLocalData(d, utmProfileWebFLocalOptions)
	}

	return nil
}

func resourceSecurityUtmProfileWebFilteringLocalUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delUtmProfileWebFLocal(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setUtmProfileWebFLocal(d, junSess); err != nil {
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
	if err := delUtmProfileWebFLocal(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setUtmProfileWebFLocal(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_security_utm_profile_web_filtering_juniper_local")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityUtmProfileWebFilteringLocalReadWJunSess(d, junSess)...)
}

func resourceSecurityUtmProfileWebFilteringLocalDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delUtmProfileWebFLocal(d.Get("name").(string), junSess); err != nil {
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
	if err := delUtmProfileWebFLocal(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_security_utm_profile_web_filtering_juniper_local")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityUtmProfileWebFilteringLocalImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	utmProfileWebFLocalExists, err := checkUtmProfileWebFLocalExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !utmProfileWebFLocalExists {
		return nil, fmt.Errorf("don't find security utm feature-profile web-filtering juniper-local with id "+
			"'%v' (id must be <name>)", d.Id())
	}
	utmProfileWebFLocalOptions, err := readUtmProfileWebFLocal(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillUtmProfileWebFLocalData(d, utmProfileWebFLocalOptions)

	result[0] = d

	return result, nil
}

func checkUtmProfileWebFLocalExists(profile string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm feature-profile web-filtering juniper-local profile \"" + profile + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setUtmProfileWebFLocal(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security utm feature-profile web-filtering juniper-local " +
		"profile \"" + d.Get("name").(string) + "\" "
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
	if d.Get("timeout").(int) != 0 {
		configSet = append(configSet, setPrefix+"timeout "+strconv.Itoa(d.Get("timeout").(int)))
	}

	return junSess.ConfigSet(configSet)
}

func readUtmProfileWebFLocal(profile string, junSess *junos.Session,
) (confRead utmProfileWebFilteringLocalOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm feature-profile web-filtering juniper-local profile \"" + profile + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = profile
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "custom-block-message "):
				confRead.customBlockMessage = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "default "):
				confRead.defaultAction = itemTrim
			case balt.CutPrefixInString(&itemTrim, "fallback-settings"):
				if len(confRead.fallbackSettings) == 0 {
					confRead.fallbackSettings = append(confRead.fallbackSettings, map[string]interface{}{
						"default":             "",
						"server_connectivity": "",
						"timeout":             "",
						"too_many_requests":   "",
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " default "):
					confRead.fallbackSettings[0]["default"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, " server-connectivity "):
					confRead.fallbackSettings[0]["server_connectivity"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, " timeout "):
					confRead.fallbackSettings[0]["timeout"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, " too-many-requests "):
					confRead.fallbackSettings[0]["too_many_requests"] = itemTrim
				}
			case balt.CutPrefixInString(&itemTrim, "timeout "):
				confRead.timeout, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delUtmProfileWebFLocal(profile string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security utm feature-profile web-filtering juniper-local "+
		"profile \""+profile+"\"")

	return junSess.ConfigSet(configSet)
}

func fillUtmProfileWebFLocalData(
	d *schema.ResourceData, utmProfileWebFLocalOptions utmProfileWebFilteringLocalOptions,
) {
	if tfErr := d.Set("name", utmProfileWebFLocalOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("custom_block_message", utmProfileWebFLocalOptions.customBlockMessage); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("default_action", utmProfileWebFLocalOptions.defaultAction); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("fallback_settings", utmProfileWebFLocalOptions.fallbackSettings); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("timeout", utmProfileWebFLocalOptions.timeout); tfErr != nil {
		panic(tfErr)
	}
}
