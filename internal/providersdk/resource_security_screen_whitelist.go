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

type screenWhiteListOptions struct {
	name    string
	address []string
}

func resourceSecurityScreenWhiteList() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityScreenWhiteListCreate,
		ReadWithoutTimeout:   resourceSecurityScreenWhiteListRead,
		UpdateWithoutTimeout: resourceSecurityScreenWhiteListUpdate,
		DeleteWithoutTimeout: resourceSecurityScreenWhiteListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityScreenWhiteListImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
			},
			"address": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateCIDRNetworkFunc(),
				},
			},
		},
	}
}

func resourceSecurityScreenWhiteListCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityScreenWhiteList(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security screen white-list not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityScreenWhiteListExists, err := checkSecurityScreenWhiteListExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityScreenWhiteListExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security screen white-list %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityScreenWhiteList(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_security_screen_whitelist")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityScreenWhiteListExists, err = checkSecurityScreenWhiteListExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityScreenWhiteListExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security screen white-list %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityScreenWhiteListReadWJunSess(d, junSess)...)
}

func resourceSecurityScreenWhiteListRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityScreenWhiteListReadWJunSess(d, junSess)
}

func resourceSecurityScreenWhiteListReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	whiteListOptions, err := readSecurityScreenWhiteList(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if whiteListOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityScreenWhiteListData(d, whiteListOptions)
	}

	return nil
}

func resourceSecurityScreenWhiteListUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityScreenWhiteList(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityScreenWhiteList(d, junSess); err != nil {
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
	if err := delSecurityScreenWhiteList(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityScreenWhiteList(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_security_screen_whitelist")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityScreenWhiteListReadWJunSess(d, junSess)...)
}

func resourceSecurityScreenWhiteListDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityScreenWhiteList(d.Get("name").(string), junSess); err != nil {
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
	if err := delSecurityScreenWhiteList(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_security_screen_whitelist")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityScreenWhiteListImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	securityScreenWhiteListExists, err := checkSecurityScreenWhiteListExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !securityScreenWhiteListExists {
		return nil, fmt.Errorf("don't find screen white-list with id '%v' (id must be <name>)", d.Id())
	}
	whiteListOptions, err := readSecurityScreenWhiteList(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityScreenWhiteListData(d, whiteListOptions)

	result[0] = d

	return result, nil
}

func checkSecurityScreenWhiteListExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "security screen white-list " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityScreenWhiteList(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security screen white-list " + d.Get("name").(string) + " "

	for _, v := range sortSetOfString(d.Get("address").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"address "+v)
	}

	return junSess.ConfigSet(configSet)
}

func readSecurityScreenWhiteList(name string, junSess *junos.Session,
) (confRead screenWhiteListOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security screen white-list " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if balt.CutPrefixInString(&itemTrim, "address ") {
				confRead.address = append(confRead.address, itemTrim)
			}
		}
	}

	return confRead, nil
}

func delSecurityScreenWhiteList(name string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security screen white-list "+name)

	return junSess.ConfigSet(configSet)
}

func fillSecurityScreenWhiteListData(d *schema.ResourceData, whiteListOptions screenWhiteListOptions) {
	if tfErr := d.Set("name", whiteListOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address", whiteListOptions.address); tfErr != nil {
		panic(tfErr)
	}
}
