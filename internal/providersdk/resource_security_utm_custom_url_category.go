package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
		if err := setUtmCustomURLCategory(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if !junos.CheckCompatibilitySecurity(junSess) {
		return diag.FromErr(fmt.Errorf("security utm custom-objects custom-url-category "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	utmCustomURLCategoryExists, err := checkUtmCustomURLCategorysExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmCustomURLCategoryExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security utm custom-objects custom-url-category %v already exists", d.Get("name").(string)))...)
	}

	if err := setUtmCustomURLCategory(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_security_utm_custom_url_category", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	utmCustomURLCategoryExists, err = checkUtmCustomURLCategorysExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmCustomURLCategoryExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security utm custom-objects custom-url-category %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityUtmCustomURLCategoryReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityUtmCustomURLCategoryRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceSecurityUtmCustomURLCategoryReadWJunSess(d, clt, junSess)
}

func resourceSecurityUtmCustomURLCategoryReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	utmCustomURLCategoryOptions, err := readUtmCustomURLCategory(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
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
		if err := delUtmCustomURLCategory(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setUtmCustomURLCategory(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delUtmCustomURLCategory(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setUtmCustomURLCategory(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_security_utm_custom_url_category", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityUtmCustomURLCategoryReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityUtmCustomURLCategoryDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delUtmCustomURLCategory(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delUtmCustomURLCategory(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_security_utm_custom_url_category", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

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
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)
	utmCustomURLCategoryExists, err := checkUtmCustomURLCategorysExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !utmCustomURLCategoryExists {
		return nil, fmt.Errorf(
			"missing security utm custom-objects custom-url-category with id '%v' (id must be <name>)", d.Id())
	}
	utmCustomURLCategoryOptions, err := readUtmCustomURLCategory(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillUtmCustomURLCategoryData(d, utmCustomURLCategoryOptions)

	result[0] = d

	return result, nil
}

func checkUtmCustomURLCategorysExists(urlCategory string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"security utm custom-objects custom-url-category "+urlCategory+junos.PipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setUtmCustomURLCategory(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security utm custom-objects custom-url-category " + d.Get("name").(string) + " "
	for _, v := range d.Get("value").([]interface{}) {
		configSet = append(configSet, setPrefix+"value "+v.(string))
	}

	return clt.ConfigSet(configSet, junSess)
}

func readUtmCustomURLCategory(urlCategory string, clt *junos.Client, junSess *junos.Session,
) (confRead utmCustomURLCategoryOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"security utm custom-objects custom-url-category "+urlCategory+junos.PipeDisplaySetRelative, junSess)
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

func delUtmCustomURLCategory(urlCategory string, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security utm custom-objects custom-url-category "+urlCategory)

	return clt.ConfigSet(configSet, junSess)
}

func fillUtmCustomURLCategoryData(d *schema.ResourceData, utmCustomURLCategoryOptions utmCustomURLCategoryOptions) {
	if tfErr := d.Set("name", utmCustomURLCategoryOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("value", utmCustomURLCategoryOptions.value); tfErr != nil {
		panic(tfErr)
	}
}
