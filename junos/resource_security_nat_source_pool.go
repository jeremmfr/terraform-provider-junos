package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type natSourcePoolOptions struct {
	portNoTranslation     bool
	portOverloadingFactor int
	name                  string
	portRange             string
	routingInstance       string
	address               []string
}

func resourceSecurityNatSourcePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityNatSourcePoolCreate,
		ReadContext:   resourceSecurityNatSourcePoolRead,
		UpdateContext: resourceSecurityNatSourcePoolUpdate,
		DeleteContext: resourceSecurityNatSourcePoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityNatSourcePoolImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32),
			},
			"address": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
		},
	}
}

func resourceSecurityNatSourcePoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security nat source pool not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	securityNatSourcePoolExists, err := checkSecurityNatSourcePoolExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if securityNatSourcePoolExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security nat source pool %v already exists", d.Get("name").(string)))
	}

	if err := setSecurityNatSourcePool(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_nat_source_pool", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityNatSourcePoolExists, err = checkSecurityNatSourcePoolExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatSourcePoolExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security nat source pool %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityNatSourcePoolReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityNatSourcePoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityNatSourcePoolReadWJnprSess(d, m, jnprSess)
}
func resourceSecurityNatSourcePoolReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	natSourcePoolOptions, err := readSecurityNatSourcePool(d.Get("name").(string), m, jnprSess)
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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSecurityNatSourcePool(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setSecurityNatSourcePool(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security_nat_source_pool", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityNatSourcePoolReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityNatSourcePoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSecurityNatSourcePool(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_nat_source_pool", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceSecurityNatSourcePoolImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	securityNatSourcePoolExists, err := checkSecurityNatSourcePoolExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityNatSourcePoolExists {
		return nil, fmt.Errorf("don't find nat source pool with id '%v' (id must be <name>)", d.Id())
	}
	natSourcePoolOptions, err := readSecurityNatSourcePool(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityNatSourcePoolData(d, natSourcePoolOptions)

	result[0] = d

	return result, nil
}

func checkSecurityNatSourcePoolExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	natSourcePoolConfig, err := sess.command("show configuration"+
		" security nat source pool "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if natSourcePoolConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setSecurityNatSourcePool(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security nat source pool " + d.Get("name").(string)
	for _, v := range d.Get("address").([]interface{}) {
		err := validateIPwithMask(v.(string))
		if err != nil {
			return err
		}
		configSet = append(configSet, setPrefix+" address "+v.(string))
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
		configSet = append(configSet, setPrefix+" port range "+rangePort[0]+" to "+rangePort[1])
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, setPrefix+" routing-instance "+d.Get("routing_instance").(string))
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readSecurityNatSourcePool(natSourcePool string,
	m interface{}, jnprSess *NetconfObject) (natSourcePoolOptions, error) {
	sess := m.(*Session)
	var confRead natSourcePoolOptions

	natSourcePoolConfig, err := sess.command("show configuration"+
		" security nat source pool "+natSourcePool+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if natSourcePoolConfig != emptyWord {
		confRead.name = natSourcePool
		var portRange string
		for _, item := range strings.Split(natSourcePoolConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "address "):
				confRead.address = append(confRead.address, strings.TrimPrefix(itemTrim, "address "))
			case itemTrim == "port no-translation":
				confRead.portNoTranslation = true
			case strings.HasPrefix(itemTrim, "port port-overloading-factor"):
				confRead.portOverloadingFactor, err = strconv.Atoi(strings.TrimPrefix(itemTrim,
					"port port-overloading-factor "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
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

func delSecurityNatSourcePool(natSourcePool string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat source pool "+natSourcePool)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func fillSecurityNatSourcePoolData(d *schema.ResourceData, natSourcePoolOptions natSourcePoolOptions) {
	if tfErr := d.Set("name", natSourcePoolOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address", natSourcePoolOptions.address); tfErr != nil {
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
		vSplit := strings.Split(v, "-")
		if len(vSplit) < 2 {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("missing range separtor - in %s", v),
				AttributePath: path,
			})
		}
		low, err := strconv.Atoi(vSplit[0])
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       err.Error(),
				AttributePath: path,
			})
		}
		high, err := strconv.Atoi(vSplit[1])
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       err.Error(),
				AttributePath: path,
			})
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
		if high > 63487 {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("high(%d) in %s is too big (max 63487)", high, v),
				AttributePath: path,
			})
		}

		return diags
	}
}
