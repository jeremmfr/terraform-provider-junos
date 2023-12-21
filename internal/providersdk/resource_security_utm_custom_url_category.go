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

type utmCustomURLCategoryOptions struct {
	name  string
	value []string
}

func resourceSecurityUtmCustomURLCategory() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityUtmCustomURLCategoryCreate,
		ReadWithoutTimeout:   resourceSecurityUtmCustomURLCategoryRead,
		UpdateWithoutTimeout: resourceSecurityUtmCustomURLCategoryUpdate,
		DeleteWithoutTimeout: resourceSecurityUtmCustomURLCategoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityUtmCustomURLCategoryImport,
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

func resourceSecurityUtmCustomURLCategoryCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setUtmCustomURLCategory(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security utm custom-objects custom-url-category "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	utmCustomURLCategoryExists, err := checkUtmCustomURLCategorysExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmCustomURLCategoryExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security utm custom-objects custom-url-category %v already exists", d.Get("name").(string)))...)
	}

	if err := setUtmCustomURLCategory(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_security_utm_custom_url_category")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	utmCustomURLCategoryExists, err = checkUtmCustomURLCategorysExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmCustomURLCategoryExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security utm custom-objects custom-url-category %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityUtmCustomURLCategoryReadWJunSess(d, junSess)...)
}

func resourceSecurityUtmCustomURLCategoryRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityUtmCustomURLCategoryReadWJunSess(d, junSess)
}

func resourceSecurityUtmCustomURLCategoryReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	utmCustomURLCategoryOptions, err := readUtmCustomURLCategory(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if utmCustomURLCategoryOptions.name == "" {
		d.SetId("")
	} else {
		fillUtmCustomURLCategoryData(d, utmCustomURLCategoryOptions)
	}

	return nil
}

func resourceSecurityUtmCustomURLCategoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delUtmCustomURLCategory(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setUtmCustomURLCategory(d, junSess); err != nil {
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
	if err := delUtmCustomURLCategory(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setUtmCustomURLCategory(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_security_utm_custom_url_category")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityUtmCustomURLCategoryReadWJunSess(d, junSess)...)
}

func resourceSecurityUtmCustomURLCategoryDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delUtmCustomURLCategory(d.Get("name").(string), junSess); err != nil {
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
	if err := delUtmCustomURLCategory(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_security_utm_custom_url_category")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityUtmCustomURLCategoryImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	utmCustomURLCategoryExists, err := checkUtmCustomURLCategorysExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !utmCustomURLCategoryExists {
		return nil, fmt.Errorf(
			"missing security utm custom-objects custom-url-category with id '%v' (id must be <name>)", d.Id())
	}
	utmCustomURLCategoryOptions, err := readUtmCustomURLCategory(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillUtmCustomURLCategoryData(d, utmCustomURLCategoryOptions)

	result[0] = d

	return result, nil
}

func checkUtmCustomURLCategorysExists(urlCategory string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm custom-objects custom-url-category " + urlCategory + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setUtmCustomURLCategory(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security utm custom-objects custom-url-category " + d.Get("name").(string) + " "
	for _, v := range d.Get("value").([]interface{}) {
		configSet = append(configSet, setPrefix+"value "+v.(string))
	}

	return junSess.ConfigSet(configSet)
}

func readUtmCustomURLCategory(urlCategory string, junSess *junos.Session,
) (confRead utmCustomURLCategoryOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm custom-objects custom-url-category " + urlCategory + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = urlCategory
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if balt.CutPrefixInString(&itemTrim, "value ") {
				confRead.value = append(confRead.value, itemTrim)
			}
		}
	}

	return confRead, nil
}

func delUtmCustomURLCategory(urlCategory string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security utm custom-objects custom-url-category "+urlCategory)

	return junSess.ConfigSet(configSet)
}

func fillUtmCustomURLCategoryData(d *schema.ResourceData, utmCustomURLCategoryOptions utmCustomURLCategoryOptions) {
	if tfErr := d.Set("name", utmCustomURLCategoryOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("value", utmCustomURLCategoryOptions.value); tfErr != nil {
		panic(tfErr)
	}
}
