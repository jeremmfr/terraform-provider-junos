package junos

import (
	"context"
	"fmt"
	"strings"

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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSecurityScreenWhiteList(d, clt, nil); err != nil {
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
		return diag.FromErr(fmt.Errorf("security screen white-list not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityScreenWhiteListExists, err := checkSecurityScreenWhiteListExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityScreenWhiteListExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security screen white-list %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityScreenWhiteList(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_security_screen_whitelist", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityScreenWhiteListExists, err = checkSecurityScreenWhiteListExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityScreenWhiteListExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security screen white-list %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityScreenWhiteListReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityScreenWhiteListRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSecurityScreenWhiteListReadWJunSess(d, clt, junSess)
}

func resourceSecurityScreenWhiteListReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	whiteListOptions, err := readSecurityScreenWhiteList(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
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
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSecurityScreenWhiteList(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityScreenWhiteList(d, clt, nil); err != nil {
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
	if err := delSecurityScreenWhiteList(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityScreenWhiteList(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_security_screen_whitelist", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityScreenWhiteListReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityScreenWhiteListDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delSecurityScreenWhiteList(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delSecurityScreenWhiteList(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_security_screen_whitelist", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityScreenWhiteListImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	securityScreenWhiteListExists, err := checkSecurityScreenWhiteListExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !securityScreenWhiteListExists {
		return nil, fmt.Errorf("don't find screen white-list with id '%v' (id must be <name>)", d.Id())
	}
	whiteListOptions, err := readSecurityScreenWhiteList(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityScreenWhiteListData(d, whiteListOptions)

	result[0] = d

	return result, nil
}

func checkSecurityScreenWhiteListExists(name string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+"security screen white-list "+name+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityScreenWhiteList(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := "set security screen white-list " + d.Get("name").(string) + " "

	for _, v := range sortSetOfString(d.Get("address").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"address "+v)
	}

	return clt.configSet(configSet, junSess)
}

func readSecurityScreenWhiteList(name string, clt *Client, junSess *junosSession,
) (confRead screenWhiteListOptions, err error) {
	showConfig, err := clt.command(cmdShowConfig+
		"security screen white-list "+name+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			if balt.CutPrefixInString(&itemTrim, "address ") {
				confRead.address = append(confRead.address, itemTrim)
			}
		}
	}

	return confRead, nil
}

func delSecurityScreenWhiteList(name string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security screen white-list "+name)

	return clt.configSet(configSet, junSess)
}

func fillSecurityScreenWhiteListData(d *schema.ResourceData, whiteListOptions screenWhiteListOptions) {
	if tfErr := d.Set("name", whiteListOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address", whiteListOptions.address); tfErr != nil {
		panic(tfErr)
	}
}
