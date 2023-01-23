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

type ntpServerOptions struct {
	prefer          bool
	key             int
	version         int
	address         string
	routingInstance string
}

func resourceSystemNtpServer() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSystemNtpServerCreate,
		ReadWithoutTimeout:   resourceSystemNtpServerRead,
		UpdateWithoutTimeout: resourceSystemNtpServerUpdate,
		DeleteWithoutTimeout: resourceSystemNtpServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSystemNtpServerImport,
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"key": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65534),
			},
			"prefer": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 4),
			},
		},
	}
}

func resourceSystemNtpServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := setSystemNtpServer(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("address").(string))

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
	ntpServerExists, err := checkSystemNtpServerExists(d.Get("address").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ntpServerExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("system ntp server %v already exists", d.Get("address").(string)))...)
	}

	if err := setSystemNtpServer(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_system_ntp_server", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ntpServerExists, err = checkSystemNtpServerExists(d.Get("address").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ntpServerExists {
		d.SetId(d.Get("address").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system ntp server %v not exists after commit "+
			"=> check your config", d.Get("address").(string)))...)
	}

	return append(diagWarns, resourceSystemNtpServerReadWJunSess(d, clt, junSess)...)
}

func resourceSystemNtpServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceSystemNtpServerReadWJunSess(d, clt, junSess)
}

func resourceSystemNtpServerReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	ntpServerOptions, err := readSystemNtpServer(d.Get("address").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ntpServerOptions.address == "" {
		d.SetId("")
	} else {
		fillSystemNtpServerData(d, ntpServerOptions)
	}

	return nil
}

func resourceSystemNtpServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delSystemNtpServer(d.Get("address").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystemNtpServer(d, clt, nil); err != nil {
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
	if err := delSystemNtpServer(d.Get("address").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemNtpServer(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_system_ntp_server", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	d.Partial(false)

	return append(diagWarns, resourceSystemNtpServerReadWJunSess(d, clt, junSess)...)
}

func resourceSystemNtpServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delSystemNtpServer(d.Get("address").(string), clt, nil); err != nil {
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
	if err := delSystemNtpServer(d.Get("address").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_system_ntp_server", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSystemNtpServerImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)

	ntpServerExists, err := checkSystemNtpServerExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !ntpServerExists {
		return nil, fmt.Errorf("don't find system ntp server with id '%v' (id must be <address>)", d.Id())
	}
	ntpServerOptions, err := readSystemNtpServer(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSystemNtpServerData(d, ntpServerOptions)

	result[0] = d

	return result, nil
}

func checkSystemNtpServerExists(address string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+"system ntp server "+address+junos.PipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSystemNtpServer(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	setPrefix := "set system ntp server " + d.Get("address").(string)
	configSet := []string{setPrefix}

	if d.Get("key").(int) != 0 {
		configSet = append(configSet, setPrefix+" key "+strconv.Itoa(d.Get("key").(int)))
	}
	if d.Get("prefer").(bool) {
		configSet = append(configSet, setPrefix+" prefer")
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, setPrefix+" routing-instance "+d.Get("routing_instance").(string))
	}
	if d.Get("version").(int) != 0 {
		configSet = append(configSet, setPrefix+" version "+strconv.Itoa(d.Get("version").(int)))
	}

	return clt.ConfigSet(configSet, junSess)
}

func readSystemNtpServer(address string, clt *junos.Client, junSess *junos.Session,
) (confRead ntpServerOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+"system ntp server "+address+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.address = address
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "key "):
				confRead.key, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "prefer":
				confRead.prefer = true
			case balt.CutPrefixInString(&itemTrim, "routing-instance "):
				confRead.routingInstance = itemTrim
			case balt.CutPrefixInString(&itemTrim, "version "):
				confRead.version, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delSystemNtpServer(address string, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system ntp server "+address)

	return clt.ConfigSet(configSet, junSess)
}

func fillSystemNtpServerData(d *schema.ResourceData, ntpServerOptions ntpServerOptions) {
	if tfErr := d.Set("address", ntpServerOptions.address); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("key", ntpServerOptions.key); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("prefer", ntpServerOptions.prefer); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", ntpServerOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("version", ntpServerOptions.version); tfErr != nil {
		panic(tfErr)
	}
}
