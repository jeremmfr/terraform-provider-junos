package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type proxyProfileOptions struct {
	protocolHTTPPort int
	name             string
	protocolHTTPHost string
}

func resourceServicesProxyProfile() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceServicesProxyProfileCreate,
		ReadWithoutTimeout:   resourceServicesProxyProfileRead,
		UpdateWithoutTimeout: resourceServicesProxyProfileUpdate,
		DeleteWithoutTimeout: resourceServicesProxyProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServicesProxyProfileImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"protocol_http_host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol_http_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
		},
	}
}

func resourceServicesProxyProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := setServicesProxyProfile(d, clt, nil); err != nil {
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
	proxyProfileExists, err := checkServicesProxyProfileExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if proxyProfileExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("services proxy profile %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesProxyProfile(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_services_proxy_profile", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	proxyProfileExists, err = checkServicesProxyProfileExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if proxyProfileExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("services proxy profile %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesProxyProfileReadWJunSess(d, clt, junSess)...)
}

func resourceServicesProxyProfileRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceServicesProxyProfileReadWJunSess(d, clt, junSess)
}

func resourceServicesProxyProfileReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	proxyProfileOptions, err := readServicesProxyProfile(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if proxyProfileOptions.name == "" {
		d.SetId("")
	} else {
		fillServicesProxyProfileData(d, proxyProfileOptions)
	}

	return nil
}

func resourceServicesProxyProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delServicesProxyProfile(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setServicesProxyProfile(d, clt, nil); err != nil {
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
	if err := delServicesProxyProfile(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesProxyProfile(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_services_proxy_profile", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesProxyProfileReadWJunSess(d, clt, junSess)...)
}

func resourceServicesProxyProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delServicesProxyProfile(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delServicesProxyProfile(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_services_proxy_profile", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesProxyProfileImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)
	proxyProfileExists, err := checkServicesProxyProfileExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !proxyProfileExists {
		return nil, fmt.Errorf("don't find services proxy profile with id '%v' (id must be <name>)", d.Id())
	}
	proxyProfileOptions, err := readServicesProxyProfile(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillServicesProxyProfileData(d, proxyProfileOptions)

	result[0] = d

	return result, nil
}

func checkServicesProxyProfileExists(profile string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"services proxy profile \""+profile+"\""+junos.PipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setServicesProxyProfile(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set services proxy profile \"" + d.Get("name").(string) + "\" "
	configSet = append(configSet, setPrefix+"protocol http host "+d.Get("protocol_http_host").(string))
	if d.Get("protocol_http_port").(int) != 0 {
		configSet = append(configSet, setPrefix+"protocol http port "+strconv.Itoa(d.Get("protocol_http_port").(int)))
	}

	return clt.ConfigSet(configSet, junSess)
}

func readServicesProxyProfile(profile string, clt *junos.Client, junSess *junos.Session,
) (confRead proxyProfileOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"services proxy profile \""+profile+"\""+junos.PipeDisplaySetRelative, junSess)
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
			case balt.CutPrefixInString(&itemTrim, "protocol http host "):
				confRead.protocolHTTPHost = itemTrim
			case balt.CutPrefixInString(&itemTrim, "protocol http port "):
				confRead.protocolHTTPPort, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delServicesProxyProfile(profile string, clt *junos.Client, junSess *junos.Session) error {
	configSet := []string{"delete services proxy profile \"" + profile + "\""}

	return clt.ConfigSet(configSet, junSess)
}

func fillServicesProxyProfileData(d *schema.ResourceData, proxyProfileOptions proxyProfileOptions,
) {
	if tfErr := d.Set("name", proxyProfileOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("protocol_http_host", proxyProfileOptions.protocolHTTPHost); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("protocol_http_port", proxyProfileOptions.protocolHTTPPort); tfErr != nil {
		panic(tfErr)
	}
}
