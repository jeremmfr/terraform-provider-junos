package junos

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type natSourcePoolOptions struct {
	portNoTranslation                  bool
	poolUtilizationAlarmClearThreshold int
	poolUtilizationAlarmRaiseThreshold int
	portOverloadingFactor              int
	addressPooling                     string
	description                        string
	name                               string
	portRange                          string
	routingInstance                    string
	address                            []string
}

func resourceSecurityNatSourcePool() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityNatSourcePoolCreate,
		ReadWithoutTimeout:   resourceSecurityNatSourcePoolRead,
		UpdateWithoutTimeout: resourceSecurityNatSourcePoolUpdate,
		DeleteWithoutTimeout: resourceSecurityNatSourcePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityNatSourcePoolImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
			},
			"address": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateCIDRFunc(),
				},
			},
			"address_pooling": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"no-paired", "paired"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pool_utilization_alarm_raise_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(50, 100),
			},
			"pool_utilization_alarm_clear_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"pool_utilization_alarm_raise_threshold"},
				ValidateFunc: validation.IntBetween(40, 100),
			},
			"port_no_translation": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"port_overloading_factor", "port_range"},
			},
			"port_overloading_factor": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validation.IntBetween(2, 32),
				ConflictsWith: []string{"port_no_translation", "port_range"},
			},
			"port_range": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"port_overloading_factor", "port_no_translation"},
				ValidateDiagFunc: validateSourcePoolPortRange(),
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
		},
	}
}

func resourceSecurityNatSourcePoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSecurityNatSourcePool(d, clt, nil); err != nil {
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
		return diag.FromErr(fmt.Errorf("security nat source pool not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityNatSourcePoolExists, err := checkSecurityNatSourcePoolExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatSourcePoolExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security nat source pool %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityNatSourcePool(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_security_nat_source_pool", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityNatSourcePoolExists, err = checkSecurityNatSourcePoolExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatSourcePoolExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security nat source pool %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityNatSourcePoolReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityNatSourcePoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSecurityNatSourcePoolReadWJunSess(d, clt, junSess)
}

func resourceSecurityNatSourcePoolReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	natSourcePoolOptions, err := readSecurityNatSourcePool(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if natSourcePoolOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityNatSourcePoolData(d, natSourcePoolOptions)
	}

	return nil
}

func resourceSecurityNatSourcePoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSecurityNatSourcePool(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityNatSourcePool(d, clt, nil); err != nil {
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
	if err := delSecurityNatSourcePool(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityNatSourcePool(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_security_nat_source_pool", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityNatSourcePoolReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityNatSourcePoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delSecurityNatSourcePool(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delSecurityNatSourcePool(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_security_nat_source_pool", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityNatSourcePoolImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)

	securityNatSourcePoolExists, err := checkSecurityNatSourcePoolExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !securityNatSourcePoolExists {
		return nil, fmt.Errorf("don't find nat source pool with id '%v' (id must be <name>)", d.Id())
	}
	natSourcePoolOptions, err := readSecurityNatSourcePool(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityNatSourcePoolData(d, natSourcePoolOptions)

	result[0] = d

	return result, nil
}

func checkSecurityNatSourcePoolExists(name string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+
		"security nat source pool "+name+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityNatSourcePool(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := "set security nat source pool " + d.Get("name").(string)
	for _, v := range d.Get("address").([]interface{}) {
		configSet = append(configSet, setPrefix+" address "+v.(string))
	}
	if d.Get("address_pooling").(string) != "" {
		configSet = append(configSet, setPrefix+" address-pooling "+d.Get("address_pooling").(string))
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+" description \""+v+"\"")
	}
	if d.Get("pool_utilization_alarm_clear_threshold").(int) != 0 {
		configSet = append(configSet, setPrefix+" pool-utilization-alarm clear-threshold "+
			strconv.Itoa(d.Get("pool_utilization_alarm_clear_threshold").(int)))
	}
	if d.Get("pool_utilization_alarm_raise_threshold").(int) != 0 {
		configSet = append(configSet, setPrefix+" pool-utilization-alarm raise-threshold "+
			strconv.Itoa(d.Get("pool_utilization_alarm_raise_threshold").(int)))
	}
	if d.Get("port_no_translation").(bool) {
		configSet = append(configSet, setPrefix+" port no-translation ")
	}
	if d.Get("port_overloading_factor").(int) != 0 {
		configSet = append(configSet, setPrefix+" port port-overloading-factor "+
			strconv.Itoa(d.Get("port_overloading_factor").(int)))
	}
	if d.Get("port_range").(string) != "" {
		rangePort := strings.Split(d.Get("port_range").(string), "-")
		if len(rangePort) > 1 {
			configSet = append(configSet, setPrefix+" port range "+rangePort[0]+" to "+rangePort[1])
		} else {
			configSet = append(configSet, setPrefix+" port range "+rangePort[0])
		}
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, setPrefix+" routing-instance "+d.Get("routing_instance").(string))
	}

	return clt.configSet(configSet, junSess)
}

func readSecurityNatSourcePool(name string, clt *Client, junSess *junosSession) (natSourcePoolOptions, error) {
	var confRead natSourcePoolOptions

	showConfig, err := clt.command(cmdShowConfig+
		"security nat source pool "+name+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		var portRange string
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "address "):
				confRead.address = append(confRead.address, strings.TrimPrefix(itemTrim, "address "))
			case strings.HasPrefix(itemTrim, "address-pooling "):
				confRead.addressPooling = strings.TrimPrefix(itemTrim, "address-pooling ")
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "pool-utilization-alarm clear-threshold "):
				confRead.poolUtilizationAlarmClearThreshold, err = strconv.Atoi(
					strings.TrimPrefix(itemTrim, "pool-utilization-alarm clear-threshold "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "pool-utilization-alarm raise-threshold "):
				confRead.poolUtilizationAlarmRaiseThreshold, err = strconv.Atoi(
					strings.TrimPrefix(itemTrim, "pool-utilization-alarm raise-threshold "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "port no-translation":
				confRead.portNoTranslation = true
			case strings.HasPrefix(itemTrim, "port port-overloading-factor"):
				confRead.portOverloadingFactor, err = strconv.Atoi(strings.TrimPrefix(itemTrim,
					"port port-overloading-factor "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "port range to"):
				portRange += "-" + strings.TrimPrefix(itemTrim, "port range to ")
			case strings.HasPrefix(itemTrim, "port range "):
				portRange = strings.TrimPrefix(itemTrim, "port range ")
			case strings.HasPrefix(itemTrim, "routing-instance"):
				confRead.routingInstance = strings.TrimPrefix(itemTrim, "routing-instance ")
			}
		}
		confRead.portRange = portRange
	}

	return confRead, nil
}

func delSecurityNatSourcePool(natSourcePool string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat source pool "+natSourcePool)

	return clt.configSet(configSet, junSess)
}

func fillSecurityNatSourcePoolData(d *schema.ResourceData, natSourcePoolOptions natSourcePoolOptions) {
	if tfErr := d.Set("name", natSourcePoolOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address", natSourcePoolOptions.address); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_pooling", natSourcePoolOptions.addressPooling); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", natSourcePoolOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("pool_utilization_alarm_clear_threshold",
		natSourcePoolOptions.poolUtilizationAlarmClearThreshold); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("pool_utilization_alarm_raise_threshold",
		natSourcePoolOptions.poolUtilizationAlarmRaiseThreshold); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("port_no_translation", natSourcePoolOptions.portNoTranslation); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("port_overloading_factor", natSourcePoolOptions.portOverloadingFactor); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("port_range", natSourcePoolOptions.portRange); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", natSourcePoolOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
}

func validateSourcePoolPortRange() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		v := i.(string)
		if ok := regexp.MustCompile(`^\d+(-\d+)?$`).MatchString(v); !ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf(`expected value of port_range to match regular expression \d+(-\d+)?, got %v`, v),
				AttributePath: path,
			})

			return diags
		}
		vSplit := strings.Split(v, "-")
		low, err := strconv.Atoi(vSplit[0])
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       err.Error(),
				AttributePath: path,
			})

			return diags
		}
		high := low
		if len(vSplit) > 1 {
			var err error
			high, err = strconv.Atoi(vSplit[1])
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       err.Error(),
					AttributePath: path,
				})

				return diags
			}
		}
		if low > high {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("low(%d) in %s bigger than high(%d)", low, v, high),
				AttributePath: path,
			})
		}
		if low < 1024 {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("low(%d) in %s is too small (min 1024)", low, v),
				AttributePath: path,
			})
		}
		if high > 65535 {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("high(%d) in %s is too big (max 65535)", high, v),
				AttributePath: path,
			})
		}

		return diags
	}
}
