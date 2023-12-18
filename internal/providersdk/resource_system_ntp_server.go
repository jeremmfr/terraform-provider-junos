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
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSystemNtpServer(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("address").(string))

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
	ntpServerExists, err := checkSystemNtpServerExists(d.Get("address").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ntpServerExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("system ntp server %v already exists", d.Get("address").(string)))...)
	}

	if err := setSystemNtpServer(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_system_ntp_server")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	ntpServerExists, err = checkSystemNtpServerExists(d.Get("address").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ntpServerExists {
		d.SetId(d.Get("address").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system ntp server %v not exists after commit "+
			"=> check your config", d.Get("address").(string)))...)
	}

	return append(diagWarns, resourceSystemNtpServerReadWJunSess(d, junSess)...)
}

func resourceSystemNtpServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSystemNtpServerReadWJunSess(d, junSess)
}

func resourceSystemNtpServerReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	ntpServerOptions, err := readSystemNtpServer(d.Get("address").(string), junSess)
	junos.MutexUnlock()
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
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSystemNtpServer(d.Get("address").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystemNtpServer(d, junSess); err != nil {
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
	if err := delSystemNtpServer(d.Get("address").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemNtpServer(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_system_ntp_server")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	d.Partial(false)

	return append(diagWarns, resourceSystemNtpServerReadWJunSess(d, junSess)...)
}

func resourceSystemNtpServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSystemNtpServer(d.Get("address").(string), junSess); err != nil {
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
	if err := delSystemNtpServer(d.Get("address").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_system_ntp_server")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

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
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	ntpServerExists, err := checkSystemNtpServerExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !ntpServerExists {
		return nil, fmt.Errorf("don't find system ntp server with id '%v' (id must be <address>)", d.Id())
	}
	ntpServerOptions, err := readSystemNtpServer(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSystemNtpServerData(d, ntpServerOptions)

	result[0] = d

	return result, nil
}

func checkSystemNtpServerExists(address string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "system ntp server " + address + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSystemNtpServer(d *schema.ResourceData, junSess *junos.Session) error {
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

	return junSess.ConfigSet(configSet)
}

func readSystemNtpServer(address string, junSess *junos.Session,
) (confRead ntpServerOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "system ntp server " + address + junos.PipeDisplaySetRelative)
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
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "prefer":
				confRead.prefer = true
			case balt.CutPrefixInString(&itemTrim, "routing-instance "):
				confRead.routingInstance = itemTrim
			case balt.CutPrefixInString(&itemTrim, "version "):
				confRead.version, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delSystemNtpServer(address string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system ntp server "+address)

	return junSess.ConfigSet(configSet)
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
