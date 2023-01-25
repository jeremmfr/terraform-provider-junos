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

type applicationSetOptions struct {
	name         string
	applications []string
}

func resourceApplicationSet() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceApplicationSetCreate,
		ReadWithoutTimeout:   resourceApplicationSetRead,
		UpdateWithoutTimeout: resourceApplicationSetUpdate,
		DeleteWithoutTimeout: resourceApplicationSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationSetImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"applications": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
		},
	}
}

func resourceApplicationSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := setApplicationSet(d, clt, nil); err != nil {
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
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	appSetExists, err := checkApplicationSetExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if appSetExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("application-set %v already exists", d.Get("name").(string)))...)
	}
	if err := setApplicationSet(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_application_set", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	appSetExists, err = checkApplicationSetExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if appSetExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("application-set %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceApplicationSetReadWJunSess(d, clt, junSess)...)
}

func resourceApplicationSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceApplicationSetReadWJunSess(d, clt, junSess)
}

func resourceApplicationSetReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	applicationSetOptions, err := readApplicationSet(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if applicationSetOptions.name == "" {
		d.SetId("")
	} else {
		fillApplicationSetData(d, applicationSetOptions)
	}

	return nil
}

func resourceApplicationSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delApplicationSet(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setApplicationSet(d, clt, nil); err != nil {
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
	if err := delApplicationSet(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setApplicationSet(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_application_set", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceApplicationSetReadWJunSess(d, clt, junSess)...)
}

func resourceApplicationSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delApplicationSet(d, clt, nil); err != nil {
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
	if err := delApplicationSet(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_application_set", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceApplicationSetImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)
	appSetExists, err := checkApplicationSetExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !appSetExists {
		return nil, fmt.Errorf("don't find application-set with id '%v' (id must be <name>)", d.Id())
	}
	applicationSetOptions, err := readApplicationSet(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillApplicationSetData(d, applicationSetOptions)
	result[0] = d

	return result, nil
}

func checkApplicationSetExists(applicationSet string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"applications application-set "+applicationSet+junos.PipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setApplicationSet(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set applications application-set " + d.Get("name").(string)
	for _, v := range d.Get("applications").([]interface{}) {
		configSet = append(configSet, setPrefix+" application "+v.(string))
	}

	return clt.ConfigSet(configSet, junSess)
}

func readApplicationSet(applicationSet string, clt *junos.Client, junSess *junos.Session,
) (confRead applicationSetOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"applications application-set "+applicationSet+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = applicationSet
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if balt.CutPrefixInString(&itemTrim, "application ") {
				confRead.applications = append(confRead.applications, itemTrim)
			}
		}
	}

	return confRead, nil
}

func delApplicationSet(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete applications application-set "+d.Get("name").(string))

	return clt.ConfigSet(configSet, junSess)
}

func fillApplicationSetData(d *schema.ResourceData, applicationSetOptions applicationSetOptions) {
	if tfErr := d.Set("name", applicationSetOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("applications", applicationSetOptions.applications); tfErr != nil {
		panic(tfErr)
	}
}