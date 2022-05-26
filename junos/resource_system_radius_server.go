package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	jdecode "github.com/jeremmfr/junosdecode"
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
			"accouting_timeout": { // old version (typo) of accounting_timeout
				Type:       schema.TypeInt,
				Computed:   true,
				Deprecated: "use accounting_timeout instead",
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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSystemRadiusServer(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("address").(string))

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
	radiusServerExists, err := checkSystemRadiusServerExists(d.Get("address").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if radiusServerExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("system radius-server %v already exists", d.Get("address").(string)))...)
	}

	if err := setSystemRadiusServer(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_system_radius_server", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	radiusServerExists, err = checkSystemRadiusServerExists(d.Get("address").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if radiusServerExists {
		d.SetId(d.Get("address").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system radius-server %v not exists after commit "+
			"=> check your config", d.Get("address").(string)))...)
	}

	return append(diagWarns, resourceSystemRadiusServerReadWJunSess(d, clt, junSess)...)
}

func resourceSystemRadiusServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSystemRadiusServerReadWJunSess(d, clt, junSess)
}

func resourceSystemRadiusServerReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	radiusServerOptions, err := readSystemRadiusServer(d.Get("address").(string), clt, junSess)
	mutex.Unlock()
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
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSystemRadiusServer(d.Get("address").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystemRadiusServer(d, clt, nil); err != nil {
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
	if err := delSystemRadiusServer(d.Get("address").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemRadiusServer(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_system_radius_server", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemRadiusServerReadWJunSess(d, clt, junSess)...)
}

func resourceSystemRadiusServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delSystemRadiusServer(d.Get("address").(string), clt, nil); err != nil {
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
	if err := delSystemRadiusServer(d.Get("address").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_system_radius_server", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSystemRadiusServerImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)

	radiusServerExists, err := checkSystemRadiusServerExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !radiusServerExists {
		return nil, fmt.Errorf("don't find system radius-server with id '%v' (id must be <address>)", d.Id())
	}
	radiusServerOptions, err := readSystemRadiusServer(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSystemRadiusServerData(d, radiusServerOptions)

	result[0] = d

	return result, nil
}

func checkSystemRadiusServerExists(address string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+"system radius-server "+address+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSystemRadiusServer(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
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

	return clt.configSet(configSet, junSess)
}

func readSystemRadiusServer(address string, clt *Client, junSess *junosSession) (radiusServerOptions, error) {
	var confRead radiusServerOptions
	confRead.accountingRetry = -1
	confRead.accountingTimeout = -1
	confRead.maxOutstandingRequests = -1

	showConfig, err := clt.command(cmdShowConfig+"system radius-server "+address+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.address = address
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "accounting-port "):
				var err error
				confRead.accountingPort, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "accounting-port "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "accounting-retry "):
				var err error
				confRead.accountingRetry, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "accounting-retry "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "accounting-timeout "):
				var err error
				confRead.accountingTimeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "accounting-timeout "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "dynamic-request-port "):
				var err error
				confRead.dynamicRequestPort, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "dynamic-request-port "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "max-outstanding-requests "):
				var err error
				confRead.maxOutstandingRequests, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-outstanding-requests "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "port "):
				var err error
				confRead.port, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "port "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "preauthentication-port "):
				var err error
				confRead.preauthenticationPort, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "preauthentication-port "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "preauthentication-secret "):
				var err error
				confRead.preauthenticationSecret, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim,
					"preauthentication-secret "), "\""))
				if err != nil {
					return confRead, fmt.Errorf("failed to decode preauthentication-secret: %w", err)
				}
			case strings.HasPrefix(itemTrim, "retry "):
				var err error
				confRead.retry, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "retry "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "routing-instance "):
				confRead.routingInstance = strings.TrimPrefix(itemTrim, "routing-instance ")
			case strings.HasPrefix(itemTrim, "secret "):
				var err error
				confRead.secret, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim,
					"secret "), "\""))
				if err != nil {
					return confRead, fmt.Errorf("failed to decode secret: %w", err)
				}
			case strings.HasPrefix(itemTrim, "source-address "):
				confRead.sourceAddress = strings.TrimPrefix(itemTrim, "source-address ")
			case strings.HasPrefix(itemTrim, "timeout "):
				var err error
				confRead.timeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "timeout "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delSystemRadiusServer(address string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system radius-server "+address)

	return clt.configSet(configSet, junSess)
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
	if tfErr := d.Set("accouting_timeout", radiusServerOptions.accountingTimeout); tfErr != nil {
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
