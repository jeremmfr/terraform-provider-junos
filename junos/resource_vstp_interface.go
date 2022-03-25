package junos

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type vstpInterfaceOptions struct {
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
	vlan                   string
	vlanGroup              string
}

func resourceVstpInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVstpInterfaceCreate,
		ReadContext:   resourceVstpInterfaceRead,
		UpdateContext: resourceVstpInterfaceUpdate,
		DeleteContext: resourceVstpInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVstpInterfaceImport,
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
			"vlan": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^all|[0-9]{1,4}$`), "must be 'all' or a VLAN id"),
				ConflictsWith: []string{"vlan_group"},
			},
			"vlan_group": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				ConflictsWith:    []string{"vlan"},
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

type vstpInterfaceVIdType int

const (
	vstpInterfaceVIdTypeNone vstpInterfaceVIdType = iota
	vstpInterfaceVIdTypeVlan
	vstpInterfaceVIdTypeVlanGroup
)

func (vType vstpInterfaceVIdType) prefix() string {
	switch vType {
	case vstpInterfaceVIdTypeNone:
		return ""
	case vstpInterfaceVIdTypeVlan:
		return "v_"
	case vstpInterfaceVIdTypeVlanGroup:
		return "vg_"
	}

	return ""
}

func resourceVstpInterfaceNewID(d *schema.ResourceData) string {
	name := d.Get("name").(string)
	routingInstance := d.Get("routing_instance").(string)
	vlan := d.Get("vlan").(string)
	vlanGroup := d.Get("vlan_group").(string)
	switch {
	case vlan != "":
		return name + idSeparator + vstpInterfaceVIdTypeVlan.prefix() + vlan + idSeparator + routingInstance
	case vlanGroup != "":
		return name + idSeparator + vstpInterfaceVIdTypeVlanGroup.prefix() + vlanGroup + idSeparator + routingInstance
	default:
		return name + idSeparator + vstpInterfaceVIdTypeNone.prefix() + idSeparator + routingInstance
	}
}

func resourceVstpInterfaceReadID(resourceID string) (vType vstpInterfaceVIdType, name, vName, routingInstnace string) {
	ressIDSplit := strings.Split(resourceID, idSeparator)
	switch len(ressIDSplit) {
	case 1:
		return vstpInterfaceVIdTypeNone,
			ressIDSplit[0],
			"",
			""
	case 2:
		return vstpInterfaceVIdTypeNone,
			ressIDSplit[0],
			"",
			ressIDSplit[1]
	default:
		switch {
		case strings.HasPrefix(ressIDSplit[1], vstpInterfaceVIdTypeVlan.prefix()):
			return vstpInterfaceVIdTypeVlan,
				ressIDSplit[0],
				strings.TrimPrefix(ressIDSplit[1], vstpInterfaceVIdTypeVlan.prefix()),
				ressIDSplit[2]
		case strings.HasPrefix(ressIDSplit[1], vstpInterfaceVIdTypeVlanGroup.prefix()):
			return vstpInterfaceVIdTypeVlanGroup,
				ressIDSplit[0],
				strings.TrimPrefix(ressIDSplit[1], vstpInterfaceVIdTypeVlanGroup.prefix()),
				ressIDSplit[2]
		default:
			return vstpInterfaceVIdTypeNone,
				ressIDSplit[0],
				ressIDSplit[1],
				ressIDSplit[2]
		}
	}
}

func resourceVstpInterfaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setVstpInterface(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resourceVstpInterfaceNewID(d))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	vstpInterfaceVIdType, name, _, routingInstance := resourceVstpInterfaceReadID(resourceVstpInterfaceNewID(d))
	vlan := d.Get("vlan").(string)
	vlanGroup := d.Get("vlan_group").(string)
	if routingInstance != defaultW {
		instanceExists, err := checkRoutingInstanceExists(routingInstance, m, jnprSess)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", routingInstance))...)
		}
	}
	if vlan != "" {
		vstpVlanExists, err := checkVstpVlanExists(vlan, routingInstance, m, jnprSess)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !vstpVlanExists {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))
			if routingInstance == defaultW {
				return append(diagWarns,
					diag.FromErr(fmt.Errorf("protocol vstp vlan %v doesn't exist", vlan))...)
			}

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("protocol vstp vlan %v in routing-instance %v doesn't exist", vlan, routingInstance))...)
		}
	}
	if vlanGroup != "" {
		vstpVlanGroupExists, err := checkVstpVlanGroupExists(vlanGroup, routingInstance, m, jnprSess)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !vstpVlanGroupExists {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))
			if routingInstance == defaultW {
				return append(diagWarns,
					diag.FromErr(fmt.Errorf("protocol vstp vlan-group group %v doesn't exist", vlanGroup))...)
			}

			return append(diagWarns,
				diag.FromErr(fmt.Errorf(
					"protocol vstp vlan-group group %v in routing-instance %v doesn't exist", vlanGroup, routingInstance))...)
		}
	}
	vstpInterfaceExists, err := checkVstpInterfaceExists(name, routingInstance, vlan, vlanGroup, m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if vstpInterfaceExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))
		switch vstpInterfaceVIdType {
		case vstpInterfaceVIdTypeNone:
			if routingInstance == defaultW {
				return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp interface %v already exists",
					name))...)
			}

			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp interface %v already exists in routing-instance %v", name, routingInstance))...)
		case vstpInterfaceVIdTypeVlan:
			if routingInstance == defaultW {
				return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp interface %v in vlan %v already exists",
					name, vlan))...)
			}

			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp interface %v already exists in vlan %s in routing-instance %v",
				name, vlan, routingInstance))...)
		case vstpInterfaceVIdTypeVlanGroup:
			if routingInstance == defaultW {
				return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp interface %v in vlan-group %v already exists",
					name, vlanGroup))...)
			}

			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp interface %v already exists in vlan-group %v in routing-instance %v",
				name, vlanGroup, routingInstance))...)
		}
	}
	if err := setVstpInterface(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_vstp_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	vstpInterfaceExists, err = checkVstpInterfaceExists(name, routingInstance, vlan, vlanGroup, m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if vstpInterfaceExists {
		d.SetId(resourceVstpInterfaceNewID(d))
	} else {
		switch vstpInterfaceVIdType {
		case vstpInterfaceVIdTypeNone:
			if routingInstance == defaultW {
				return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp interface %v not exists after commit "+
					"=> check your config", name))...)
			}

			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp interface %v not exists in routing-instance %v after commit "+
					"=> check your config", name, routingInstance))...)
		case vstpInterfaceVIdTypeVlan:
			if routingInstance == defaultW {
				return append(diagWarns, diag.FromErr(fmt.Errorf(
					"protocols vstp interface %v in vlan %v not exists after commit "+
						"=> check your config", name, vlan))...)
			}

			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp interface %v not exists in vlan %v in routing-instance %v after commit "+
					"=> check your config", name, vlan, routingInstance))...)
		case vstpInterfaceVIdTypeVlanGroup:
			if routingInstance == defaultW {
				return append(diagWarns, diag.FromErr(fmt.Errorf(
					"protocols vstp interface %v in vlan-group %v not exists after commit "+
						"=> check your config", name, vlanGroup))...)
			}

			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp interface %v not exists in vlan-group %v in routing-instance %v after commit "+
					"=> check your config",
				name, vlanGroup, routingInstance))...)
		}
	}

	return append(diagWarns, resourceVstpInterfaceReadWJnprSess(d, m, jnprSess)...)
}

func resourceVstpInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceVstpInterfaceReadWJnprSess(d, m, jnprSess)
}

func resourceVstpInterfaceReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
) diag.Diagnostics {
	mutex.Lock()
	vstpInterfaceOptions, err := readVstpInterface(
		d.Get("name").(string), d.Get("routing_instance").(string),
		d.Get("vlan").(string), d.Get("vlan_group").(string),
		m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if vstpInterfaceOptions.name == "" {
		d.SetId("")
	} else {
		fillVstpInterfaceData(d, vstpInterfaceOptions)
	}

	return nil
}

func resourceVstpInterfaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delVstpInterface(
			d.Get("name").(string), d.Get("routing_instance").(string),
			d.Get("vlan").(string), d.Get("vlan_group").(string),
			m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setVstpInterface(d, m, nil); err != nil {
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
	if err := delVstpInterface(
		d.Get("name").(string), d.Get("routing_instance").(string),
		d.Get("vlan").(string), d.Get("vlan_group").(string),
		m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setVstpInterface(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_vstp_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceVstpInterfaceReadWJnprSess(d, m, jnprSess)...)
}

func resourceVstpInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delVstpInterface(
			d.Get("name").(string), d.Get("routing_instance").(string),
			d.Get("vlan").(string), d.Get("vlan_group").(string),
			m, nil); err != nil {
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
	if err := delVstpInterface(
		d.Get("name").(string), d.Get("routing_instance").(string),
		d.Get("vlan").(string), d.Get("vlan_group").(string),
		m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_vstp_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceVstpInterfaceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	if len(strings.Split(d.Id(), idSeparator)) < 3 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	vType, name, vName, routingInstance := resourceVstpInterfaceReadID(d.Id())
	var vstpInterfaceExists bool
	switch vType {
	case vstpInterfaceVIdTypeNone:
		vstpInterfaceExists, err = checkVstpInterfaceExists(name, routingInstance, "", "", m, jnprSess)
	case vstpInterfaceVIdTypeVlan:
		vstpInterfaceExists, err = checkVstpInterfaceExists(name, routingInstance, vName, "", m, jnprSess)
	case vstpInterfaceVIdTypeVlanGroup:
		vstpInterfaceExists, err = checkVstpInterfaceExists(name, routingInstance, "", vName, m, jnprSess)
	}
	if err != nil {
		return nil, err
	}
	if !vstpInterfaceExists {
		return nil, fmt.Errorf("don't find protocols vstp interface with id '%v' "+
			"(id must be <name>%s%s<routing_instance>, "+
			"<name>%s%s<vlan>%s<routing_instance> or <name>%s%s<vlan_group>%s<routing_instance>)", d.Id(),
			idSeparator, idSeparator,
			idSeparator, vstpInterfaceVIdTypeVlan.prefix(), idSeparator,
			idSeparator, vstpInterfaceVIdTypeVlanGroup.prefix(), idSeparator)
	}
	var vstpInterfaceOptions vstpInterfaceOptions
	switch vType {
	case vstpInterfaceVIdTypeNone:
		vstpInterfaceOptions, err = readVstpInterface(name, routingInstance, "", "", m, jnprSess)
	case vstpInterfaceVIdTypeVlan:
		vstpInterfaceOptions, err = readVstpInterface(name, routingInstance, vName, "", m, jnprSess)
	case vstpInterfaceVIdTypeVlanGroup:
		vstpInterfaceOptions, err = readVstpInterface(name, routingInstance, "", vName, m, jnprSess)
	}
	if err != nil {
		return nil, err
	}
	fillVstpInterfaceData(d, vstpInterfaceOptions)
	d.SetId(resourceVstpInterfaceNewID(d))
	result[0] = d

	return result, nil
}

func checkVstpInterfaceExists(name, routingInstance, vlan, vlanGroup string, m interface{}, jnprSess *NetconfObject,
) (bool, error) {
	sess := m.(*Session)
	var showConfig string
	var err error
	if vlan != "" && vlanGroup != "" {
		return false, fmt.Errorf("internal error: checkVstpInterfaceExists called with vlan and vlanGroup")
	}
	if routingInstance == defaultW {
		switch {
		case vlan != "":
			showConfig, err = sess.command(cmdShowConfig+
				"protocols vstp vlan "+vlan+" interface "+name+pipeDisplaySet, jnprSess)
		case vlanGroup != "":
			showConfig, err = sess.command(cmdShowConfig+
				"protocols vstp vlan-group group "+vlanGroup+" interface "+name+pipeDisplaySet, jnprSess)
		default:
			showConfig, err = sess.command(cmdShowConfig+
				"protocols vstp interface "+name+pipeDisplaySet, jnprSess)
		}
	} else {
		switch {
		case vlan != "":
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
				"protocols vstp vlan "+vlan+" interface "+name+pipeDisplaySet, jnprSess)
		case vlanGroup != "":
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
				"protocols vstp vlan-group group "+vlanGroup+" interface "+name+pipeDisplaySet, jnprSess)
		default:
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
				"protocols vstp interface "+name+pipeDisplaySet, jnprSess)
		}
	}
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setVstpInterface(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	name := d.Get("name").(string)
	vlan := d.Get("vlan").(string)
	vlanGroup := d.Get("vlan_group").(string)
	setPrefix := setLS
	if rI := d.Get("routing_instance").(string); rI != defaultW {
		setPrefix = setRoutingInstances + rI + " "
	}
	switch {
	case vlan != "":
		setPrefix += "protocols vstp vlan " + vlan + " interface " + name + " "
	case vlanGroup != "":
		setPrefix += "protocols vstp vlan-group group " + vlanGroup + " interface " + name + " "
	default:
		setPrefix += "protocols vstp interface " + name + " "
	}

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

func readVstpInterface(name, routingInstance, vlan, vlanGroup string, m interface{}, jnprSess *NetconfObject,
) (vstpInterfaceOptions, error) {
	sess := m.(*Session)
	var confRead vstpInterfaceOptions
	if vlan != "" && vlanGroup != "" {
		return confRead, fmt.Errorf("internal error: readVstpInterface called with vlan and vlanGroup")
	}
	confRead.priority = -1 // default -1
	var showConfig string
	var err error
	if routingInstance == defaultW {
		switch {
		case vlan != "":
			showConfig, err = sess.command(cmdShowConfig+
				"protocols vstp vlan "+vlan+" interface "+name+pipeDisplaySetRelative, jnprSess)
		case vlanGroup != "":
			showConfig, err = sess.command(cmdShowConfig+
				"protocols vstp vlan-group group "+vlanGroup+" interface "+name+pipeDisplaySetRelative, jnprSess)
		default:
			showConfig, err = sess.command(cmdShowConfig+
				"protocols vstp interface "+name+pipeDisplaySetRelative, jnprSess)
		}
	} else {
		switch {
		case vlan != "":
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
				"protocols vstp vlan "+vlan+" interface "+name+pipeDisplaySetRelative, jnprSess)
		case vlanGroup != "":
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
				"protocols vstp vlan-group group "+vlanGroup+" interface "+name+pipeDisplaySetRelative, jnprSess)
		default:
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
				"protocols vstp interface "+name+pipeDisplaySetRelative, jnprSess)
		}
	}
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		confRead.routingInstance = routingInstance
		confRead.vlan = vlan
		confRead.vlanGroup = vlanGroup
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

func delVstpInterface(name, routingInstance, vlan, vlanGroup string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	if vlan != "" && vlanGroup != "" {
		return fmt.Errorf("internal error: delVstpInterface called with vlan and vlanGroup")
	}
	if routingInstance == defaultW {
		switch {
		case vlan != "":
			configSet = append(configSet, "delete protocols vstp vlan "+vlan+" interface "+name)
		case vlanGroup != "":
			configSet = append(configSet, "delete protocols vstp vlan-group group "+vlanGroup+" interface "+name)
		default:
			configSet = append(configSet, "delete protocols vstp interface "+name)
		}
	} else {
		switch {
		case vlan != "":
			configSet = append(configSet, delRoutingInstances+routingInstance+
				" protocols vstp vlan "+vlan+" interface "+name)
		case vlanGroup != "":
			configSet = append(configSet, delRoutingInstances+routingInstance+
				" protocols vstp vlan-group group "+vlanGroup+" interface "+name)
		default:
			configSet = append(configSet, delRoutingInstances+routingInstance+
				" protocols vstp interface "+name)
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func fillVstpInterfaceData(d *schema.ResourceData, vstpInterfaceOptions vstpInterfaceOptions) {
	if tfErr := d.Set("name", vstpInterfaceOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", vstpInterfaceOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vlan", vstpInterfaceOptions.vlan); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vlan_group", vstpInterfaceOptions.vlanGroup); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("access_trunk", vstpInterfaceOptions.accessTrunk); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bpdu_timeout_action_alarm", vstpInterfaceOptions.bpduTimeoutActionAlarm); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bpdu_timeout_action_block", vstpInterfaceOptions.bpduTimeoutActionBlock); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("cost", vstpInterfaceOptions.cost); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("edge", vstpInterfaceOptions.edge); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("mode", vstpInterfaceOptions.mode); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_root_port", vstpInterfaceOptions.noRootPort); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("priority", vstpInterfaceOptions.priority); tfErr != nil {
		panic(tfErr)
	}
}
