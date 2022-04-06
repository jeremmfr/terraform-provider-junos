package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type rstpInterfaceOptions struct {
	accessTrunk            bool
	bpduTimeoutActionAlarm bool
	bpduTimeoutActionBlock bool
	edge                   bool
	noRootPort             bool
	cost                   int
	priority               int
	mode                   string
	name                   string
	routingInstance        string
}

func resourceRstpInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRstpInterfaceCreate,
		ReadContext:   resourceRstpInterfaceRead,
		UpdateContext: resourceRstpInterfaceUpdate,
		DeleteContext: resourceRstpInterfaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRstpInterfaceImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if strings.Count(value, ".") > 0 {
						errors = append(errors, fmt.Errorf(
							"%q in %q cannot have a dot", value, k))
					}

					return
				},
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          defaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"access_trunk": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bpdu_timeout_action_alarm": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bpdu_timeout_action_block": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cost": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 200000000),
			},
			"edge": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"point-to-point", "shared"}, false),
			},
			"no_root_port": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(0, 240),
			},
		},
	}
}

func resourceRstpInterfaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setRstpInterface(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string) + idSeparator + d.Get("routing_instance").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if d.Get("routing_instance").(string) != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), m, jnprSess)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	rstpInterfaceExists, err := checkRstpInterfaceExists(
		d.Get("name").(string), d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if rstpInterfaceExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))
		if d.Get("routing_instance").(string) == defaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf("protocols rstp interface %v already exists",
				d.Get("name").(string)))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"protocols rstp interface %v already exists in routing-instance %v",
			d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}

	if err := setRstpInterface(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_rstp_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	rstpInterfaceExists, err = checkRstpInterfaceExists(
		d.Get("name").(string), d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if rstpInterfaceExists {
		d.SetId(d.Get("name").(string) + idSeparator + d.Get("routing_instance").(string))
	} else {
		if d.Get("routing_instance").(string) == defaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf("protocols rstp interface %v not exists after commit "+
				"=> check your config", d.Get("name").(string)))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"protocols rstp interface %v not exists in routing-instance %v after commit "+
				"=> check your config", d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceRstpInterfaceReadWJnprSess(d, m, jnprSess)...)
}

func resourceRstpInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceRstpInterfaceReadWJnprSess(d, m, jnprSess)
}

func resourceRstpInterfaceReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
) diag.Diagnostics {
	mutex.Lock()
	rstpInterfaceOptions, err := readRstpInterface(
		d.Get("name").(string), d.Get("routing_instance").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if rstpInterfaceOptions.name == "" {
		d.SetId("")
	} else {
		fillRstpInterfaceData(d, rstpInterfaceOptions)
	}

	return nil
}

func resourceRstpInterfaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delRstpInterface(d.Get("name").(string), d.Get("routing_instance").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setRstpInterface(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delRstpInterface(d.Get("name").(string), d.Get("routing_instance").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setRstpInterface(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_rstp_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRstpInterfaceReadWJnprSess(d, m, jnprSess)...)
}

func resourceRstpInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delRstpInterface(d.Get("name").(string), d.Get("routing_instance").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delRstpInterface(d.Get("name").(string), d.Get("routing_instance").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_rstp_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceRstpInterfaceImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	rstpInterfaceExists, err := checkRstpInterfaceExists(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !rstpInterfaceExists {
		return nil, fmt.Errorf("don't find protocols rstp interface with id '%v' "+
			"(id must be <name>%s<routing_instance>)", d.Id(), idSeparator)
	}
	rstpInterfaceOptions, err := readRstpInterface(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillRstpInterfaceData(d, rstpInterfaceOptions)

	result[0] = d

	return result, nil
}

func checkRstpInterfaceExists(name, routingInstance string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	var showConfig string
	var err error
	if routingInstance == defaultW {
		showConfig, err = sess.command(cmdShowConfig+
			"protocols rstp interface "+name+pipeDisplaySet, jnprSess)
	} else {
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols rstp interface "+name+pipeDisplaySet, jnprSess)
	}
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setRstpInterface(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := setLS
	if rI := d.Get("routing_instance").(string); rI != defaultW {
		setPrefix = setRoutingInstances + rI + " "
	}
	setPrefix += "protocols rstp interface " + d.Get("name").(string) + " "

	configSet = append(configSet, setPrefix)
	if d.Get("access_trunk").(bool) {
		configSet = append(configSet, setPrefix+"access-trunk")
	}
	if d.Get("bpdu_timeout_action_alarm").(bool) {
		configSet = append(configSet, setPrefix+"bpdu-timeout-action alarm")
	}
	if d.Get("bpdu_timeout_action_block").(bool) {
		configSet = append(configSet, setPrefix+"bpdu-timeout-action block")
	}
	if v := d.Get("cost").(int); v != 0 {
		configSet = append(configSet, setPrefix+"cost "+strconv.Itoa(v))
	}
	if d.Get("edge").(bool) {
		configSet = append(configSet, setPrefix+"edge")
	}
	if v := d.Get("mode").(string); v != "" {
		configSet = append(configSet, setPrefix+"mode "+v)
	}
	if d.Get("no_root_port").(bool) {
		configSet = append(configSet, setPrefix+"no-root-port")
	}
	if v := d.Get("priority").(int); v != -1 {
		configSet = append(configSet, setPrefix+"priority "+strconv.Itoa(v))
	}

	return sess.configSet(configSet, jnprSess)
}

func readRstpInterface(name, routingInstance string, m interface{}, jnprSess *NetconfObject,
) (rstpInterfaceOptions, error) {
	sess := m.(*Session)
	var confRead rstpInterfaceOptions
	confRead.priority = -1 // default -1
	var showConfig string
	var err error
	if routingInstance == defaultW {
		showConfig, err = sess.command(cmdShowConfig+
			"protocols rstp interface "+name+pipeDisplaySetRelative, jnprSess)
	} else {
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols rstp interface "+name+pipeDisplaySetRelative, jnprSess)
	}
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		confRead.routingInstance = routingInstance
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case itemTrim == "access-trunk":
				confRead.accessTrunk = true
			case itemTrim == "bpdu-timeout-action alarm":
				confRead.bpduTimeoutActionAlarm = true
			case itemTrim == "bpdu-timeout-action block":
				confRead.bpduTimeoutActionBlock = true
			case strings.HasPrefix(itemTrim, "cost "):
				var err error
				confRead.cost, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "cost "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "edge":
				confRead.edge = true
			case strings.HasPrefix(itemTrim, "mode "):
				confRead.mode = strings.TrimPrefix(itemTrim, "mode ")
			case itemTrim == "no-root-port":
				confRead.noRootPort = true
			case strings.HasPrefix(itemTrim, "priority "):
				var err error
				confRead.priority, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "priority "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delRstpInterface(name, routingInstance string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)

	if routingInstance == defaultW {
		configSet = append(configSet, "delete protocols rstp interface "+name)
	} else {
		configSet = append(configSet, delRoutingInstances+routingInstance+" protocols rstp interface "+name)
	}

	return sess.configSet(configSet, jnprSess)
}

func fillRstpInterfaceData(d *schema.ResourceData, rstpInterfaceOptions rstpInterfaceOptions) {
	if tfErr := d.Set("name", rstpInterfaceOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", rstpInterfaceOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("access_trunk", rstpInterfaceOptions.accessTrunk); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bpdu_timeout_action_alarm", rstpInterfaceOptions.bpduTimeoutActionAlarm); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bpdu_timeout_action_block", rstpInterfaceOptions.bpduTimeoutActionBlock); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("cost", rstpInterfaceOptions.cost); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("edge", rstpInterfaceOptions.edge); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("mode", rstpInterfaceOptions.mode); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_root_port", rstpInterfaceOptions.noRootPort); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("priority", rstpInterfaceOptions.priority); tfErr != nil {
		panic(tfErr)
	}
}
