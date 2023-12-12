package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type utmCustomURLPatternOptions struct {
	name  string
	value []string
}

func resourceSecurityUtmCustomURLPattern() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityUtmCustomURLPatternCreate,
		ReadWithoutTimeout:   resourceSecurityUtmCustomURLPatternRead,
		UpdateWithoutTimeout: resourceSecurityUtmCustomURLPatternUpdate,
		DeleteWithoutTimeout: resourceSecurityUtmCustomURLPatternDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityUtmCustomURLPatternImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
			},
			"value": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSecurityUtmCustomURLPatternCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setUtmCustomURLPattern(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security utm custom-objects url-pattern "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	utmCustomURLPatternExists, err := checkUtmCustomURLPatternsExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmCustomURLPatternExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security utm custom-objects url-pattern %v already exists", d.Get("name").(string)))...)
	}
	if err := setUtmCustomURLPattern(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_security_utm_custom_url_pattern")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	utmCustomURLPatternExists, err = checkUtmCustomURLPatternsExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmCustomURLPatternExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security utm custom-objects url-pattern %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityUtmCustomURLPatternReadWJunSess(d, junSess)...)
}

func resourceSecurityUtmCustomURLPatternRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityUtmCustomURLPatternReadWJunSess(d, junSess)
}

func resourceSecurityUtmCustomURLPatternReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	utmCustomURLPatternOptions, err := readUtmCustomURLPattern(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if utmCustomURLPatternOptions.name == "" {
		d.SetId("")
	} else {
		fillUtmCustomURLPatternData(d, utmCustomURLPatternOptions)
	}

	return nil
}

func resourceSecurityUtmCustomURLPatternUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delUtmCustomURLPattern(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setUtmCustomURLPattern(d, junSess); err != nil {
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
	if err := delUtmCustomURLPattern(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setUtmCustomURLPattern(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_security_utm_custom_url_pattern")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityUtmCustomURLPatternReadWJunSess(d, junSess)...)
}

func resourceSecurityUtmCustomURLPatternDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delUtmCustomURLPattern(d.Get("name").(string), junSess); err != nil {
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
	if err := delUtmCustomURLPattern(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_security_utm_custom_url_pattern")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityUtmCustomURLPatternImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	utmCustomURLPatternExists, err := checkUtmCustomURLPatternsExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !utmCustomURLPatternExists {
		return nil, fmt.Errorf("don't find security utm custom-objects url-pattern with id '%v' (id must be <name>)", d.Id())
	}
	utmCustomURLPatternOptions, err := readUtmCustomURLPattern(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillUtmCustomURLPatternData(d, utmCustomURLPatternOptions)

	result[0] = d

	return result, nil
}

func checkUtmCustomURLPatternsExists(urlPattern string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm custom-objects url-pattern " + urlPattern + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setUtmCustomURLPattern(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security utm custom-objects url-pattern " + d.Get("name").(string) + " "
	for _, v := range d.Get("value").([]interface{}) {
		configSet = append(configSet, setPrefix+"value "+v.(string))
	}

	return junSess.ConfigSet(configSet)
}

func readUtmCustomURLPattern(urlPattern string, junSess *junos.Session,
) (confRead utmCustomURLPatternOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm custom-objects url-pattern " + urlPattern + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = urlPattern
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if balt.CutPrefixInString(&itemTrim, "value ") {
				confRead.value = append(confRead.value, strings.Trim(itemTrim, "\""))
			}
		}
	}

	return confRead, nil
}

func delUtmCustomURLPattern(urlPattern string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security utm custom-objects url-pattern "+urlPattern)

	return junSess.ConfigSet(configSet)
}

func fillUtmCustomURLPatternData(d *schema.ResourceData, utmCustomURLPatternOptions utmCustomURLPatternOptions) {
	if tfErr := d.Set("name", utmCustomURLPatternOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("value", utmCustomURLPatternOptions.value); tfErr != nil {
		panic(tfErr)
	}
}
