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
	jdecode "github.com/jeremmfr/junosdecode"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type radiusServerOptions struct {
	accountingPort          int
	accountingRetry         int
	accountingTimeout       int
	dynamicRequestPort      int
	maxOutstandingRequests  int
	port                    int
	preauthenticationPort   int
	retry                   int
	timeout                 int
	address                 string
	preauthenticationSecret string
	routingInstance         string
	secret                  string
	sourceAddress           string
}

func resourceSystemRadiusServer() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSystemRadiusServerCreate,
		ReadWithoutTimeout:   resourceSystemRadiusServerRead,
		UpdateWithoutTimeout: resourceSystemRadiusServerUpdate,
		DeleteWithoutTimeout: resourceSystemRadiusServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSystemRadiusServerImport,
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"accounting_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"accounting_retry": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"accounting_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(0, 1000),
			},
			"dynamic_request_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"max_outstanding_requests": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(0, 2000),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"preauthentication_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"preauthentication_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"retry": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"source_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
			},
		},
	}
}

func resourceSystemRadiusServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSystemRadiusServer(d, junSess); err != nil {
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
	radiusServerExists, err := checkSystemRadiusServerExists(d.Get("address").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if radiusServerExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("system radius-server %v already exists", d.Get("address").(string)))...)
	}

	if err := setSystemRadiusServer(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_system_radius_server")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	radiusServerExists, err = checkSystemRadiusServerExists(d.Get("address").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if radiusServerExists {
		d.SetId(d.Get("address").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system radius-server %v not exists after commit "+
			"=> check your config", d.Get("address").(string)))...)
	}

	return append(diagWarns, resourceSystemRadiusServerReadWJunSess(d, junSess)...)
}

func resourceSystemRadiusServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSystemRadiusServerReadWJunSess(d, junSess)
}

func resourceSystemRadiusServerReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	radiusServerOptions, err := readSystemRadiusServer(d.Get("address").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if radiusServerOptions.address == "" {
		d.SetId("")
	} else {
		fillSystemRadiusServerData(d, radiusServerOptions)
	}

	return nil
}

func resourceSystemRadiusServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSystemRadiusServer(d.Get("address").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystemRadiusServer(d, junSess); err != nil {
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
	if err := delSystemRadiusServer(d.Get("address").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemRadiusServer(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_system_radius_server")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemRadiusServerReadWJunSess(d, junSess)...)
}

func resourceSystemRadiusServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSystemRadiusServer(d.Get("address").(string), junSess); err != nil {
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
	if err := delSystemRadiusServer(d.Get("address").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_system_radius_server")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSystemRadiusServerImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	radiusServerExists, err := checkSystemRadiusServerExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !radiusServerExists {
		return nil, fmt.Errorf("don't find system radius-server with id '%v' (id must be <address>)", d.Id())
	}
	radiusServerOptions, err := readSystemRadiusServer(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSystemRadiusServerData(d, radiusServerOptions)

	result[0] = d

	return result, nil
}

func checkSystemRadiusServerExists(address string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "system radius-server " + address + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSystemRadiusServer(d *schema.ResourceData, junSess *junos.Session) error {
	setPrefix := "set system radius-server " + d.Get("address").(string)
	configSet := []string{setPrefix + " secret \"" + d.Get("secret").(string) + "\""}

	if d.Get("accounting_port").(int) != 0 {
		configSet = append(configSet, setPrefix+" accounting-port "+
			strconv.Itoa(d.Get("accounting_port").(int)))
	}
	if d.Get("accounting_retry").(int) != -1 {
		configSet = append(configSet, setPrefix+" accounting-retry "+
			strconv.Itoa(d.Get("accounting_retry").(int)))
	}
	if d.Get("accounting_timeout").(int) != -1 {
		configSet = append(configSet, setPrefix+" accounting-timeout "+
			strconv.Itoa(d.Get("accounting_timeout").(int)))
	}
	if d.Get("dynamic_request_port").(int) != 0 {
		configSet = append(configSet, setPrefix+" dynamic-request-port "+
			strconv.Itoa(d.Get("dynamic_request_port").(int)))
	}
	if d.Get("max_outstanding_requests").(int) != -1 {
		configSet = append(configSet, setPrefix+" max-outstanding-requests "+
			strconv.Itoa(d.Get("max_outstanding_requests").(int)))
	}
	if d.Get("port").(int) != 0 {
		configSet = append(configSet, setPrefix+" port "+
			strconv.Itoa(d.Get("port").(int)))
	}
	if d.Get("preauthentication_port").(int) != 0 {
		configSet = append(configSet, setPrefix+" preauthentication-port "+
			strconv.Itoa(d.Get("preauthentication_port").(int)))
	}
	if d.Get("preauthentication_secret").(string) != "" {
		configSet = append(configSet, setPrefix+" preauthentication-secret \""+
			d.Get("preauthentication_secret").(string)+"\"")
	}
	if d.Get("retry").(int) != 0 {
		configSet = append(configSet, setPrefix+" retry "+
			strconv.Itoa(d.Get("retry").(int)))
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, setPrefix+" routing-instance "+
			d.Get("routing_instance").(string))
	}
	if d.Get("source_address").(string) != "" {
		configSet = append(configSet, setPrefix+" source-address "+
			d.Get("source_address").(string))
	}
	if d.Get("timeout").(int) != 0 {
		configSet = append(configSet, setPrefix+" timeout "+
			strconv.Itoa(d.Get("timeout").(int)))
	}

	return junSess.ConfigSet(configSet)
}

func readSystemRadiusServer(address string, junSess *junos.Session,
) (confRead radiusServerOptions, err error) {
	// default -1
	confRead.accountingRetry = -1
	confRead.accountingTimeout = -1
	confRead.maxOutstandingRequests = -1
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system radius-server " + address + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "accounting-port "):
				confRead.accountingPort, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "accounting-retry "):
				confRead.accountingRetry, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "accounting-timeout "):
				confRead.accountingTimeout, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "dynamic-request-port "):
				confRead.dynamicRequestPort, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "max-outstanding-requests "):
				confRead.maxOutstandingRequests, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "port "):
				confRead.port, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "preauthentication-port "):
				confRead.preauthenticationPort, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "preauthentication-secret "):
				confRead.preauthenticationSecret, err = jdecode.Decode(strings.Trim(itemTrim, "\""))
				if err != nil {
					return confRead, fmt.Errorf("decoding preauthentication-secret: %w", err)
				}
			case balt.CutPrefixInString(&itemTrim, "retry "):
				confRead.retry, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "routing-instance "):
				confRead.routingInstance = itemTrim
			case balt.CutPrefixInString(&itemTrim, "secret "):
				confRead.secret, err = jdecode.Decode(strings.Trim(itemTrim, "\""))
				if err != nil {
					return confRead, fmt.Errorf("decoding secret: %w", err)
				}
			case balt.CutPrefixInString(&itemTrim, "source-address "):
				confRead.sourceAddress = itemTrim
			case balt.CutPrefixInString(&itemTrim, "timeout "):
				confRead.timeout, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delSystemRadiusServer(address string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system radius-server "+address)

	return junSess.ConfigSet(configSet)
}

func fillSystemRadiusServerData(d *schema.ResourceData, radiusServerOptions radiusServerOptions) {
	if tfErr := d.Set("address", radiusServerOptions.address); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("accounting_port", radiusServerOptions.accountingPort); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("accounting_retry", radiusServerOptions.accountingRetry); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("accounting_timeout", radiusServerOptions.accountingTimeout); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dynamic_request_port", radiusServerOptions.dynamicRequestPort); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("max_outstanding_requests", radiusServerOptions.maxOutstandingRequests); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("port", radiusServerOptions.port); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("preauthentication_port", radiusServerOptions.preauthenticationPort); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("preauthentication_secret", radiusServerOptions.preauthenticationSecret); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("retry", radiusServerOptions.retry); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", radiusServerOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("secret", radiusServerOptions.secret); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("source_address", radiusServerOptions.sourceAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("timeout", radiusServerOptions.timeout); tfErr != nil {
		panic(tfErr)
	}
}
