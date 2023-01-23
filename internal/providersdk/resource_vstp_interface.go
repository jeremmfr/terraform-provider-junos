package providersdk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
		CreateWithoutTimeout: resourceVstpInterfaceCreate,
		ReadWithoutTimeout:   resourceVstpInterfaceRead,
		UpdateWithoutTimeout: resourceVstpInterfaceUpdate,
		DeleteWithoutTimeout: resourceVstpInterfaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVstpInterfaceImport,
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
				Default:          junos.DefaultW,
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
		return name + junos.IDSeparator +
			vstpInterfaceVIdTypeVlan.prefix() + vlan + junos.IDSeparator +
			routingInstance
	case vlanGroup != "":
		return name + junos.IDSeparator +
			vstpInterfaceVIdTypeVlanGroup.prefix() + vlanGroup + junos.IDSeparator +
			routingInstance
	default:
		return name + junos.IDSeparator +
			vstpInterfaceVIdTypeNone.prefix() + junos.IDSeparator +
			routingInstance
	}
}

func resourceVstpInterfaceReadID(resourceID string) (vType vstpInterfaceVIdType, name, vName, routingInstnace string) {
	ressIDSplit := strings.Split(resourceID, junos.IDSeparator)
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
		case balt.CutPrefixInString(&ressIDSplit[1], vstpInterfaceVIdTypeVlan.prefix()):
			return vstpInterfaceVIdTypeVlan,
				ressIDSplit[0],
				ressIDSplit[1],
				ressIDSplit[2]
		case balt.CutPrefixInString(&ressIDSplit[1], vstpInterfaceVIdTypeVlanGroup.prefix()):
			return vstpInterfaceVIdTypeVlanGroup,
				ressIDSplit[0],
				ressIDSplit[1],
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := setVstpInterface(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resourceVstpInterfaceNewID(d))

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
	vstpInterfaceVIdType, name, _, routingInstance := resourceVstpInterfaceReadID(resourceVstpInterfaceNewID(d))
	vlan := d.Get("vlan").(string)
	vlanGroup := d.Get("vlan_group").(string)
	if routingInstance != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(routingInstance, clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", routingInstance))...)
		}
	}
	if vlan != "" {
		vstpVlanExists, err := checkVstpVlanExists(vlan, routingInstance, clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !vstpVlanExists {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))
			if routingInstance == junos.DefaultW {
				return append(diagWarns,
					diag.FromErr(fmt.Errorf("protocol vstp vlan %v doesn't exist", vlan))...)
			}

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("protocol vstp vlan %v in routing-instance %v doesn't exist", vlan, routingInstance))...)
		}
	}
	if vlanGroup != "" {
		vstpVlanGroupExists, err := checkVstpVlanGroupExists(vlanGroup, routingInstance, clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !vstpVlanGroupExists {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))
			if routingInstance == junos.DefaultW {
				return append(diagWarns,
					diag.FromErr(fmt.Errorf("protocol vstp vlan-group group %v doesn't exist", vlanGroup))...)
			}

			return append(diagWarns,
				diag.FromErr(fmt.Errorf(
					"protocol vstp vlan-group group %v in routing-instance %v doesn't exist", vlanGroup, routingInstance))...)
		}
	}
	vstpInterfaceExists, err := checkVstpInterfaceExists(name, routingInstance, vlan, vlanGroup, clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if vstpInterfaceExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))
		switch vstpInterfaceVIdType {
		case vstpInterfaceVIdTypeNone:
			if routingInstance == junos.DefaultW {
				return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp interface %v already exists",
					name))...)
			}

			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp interface %v already exists in routing-instance %v", name, routingInstance))...)
		case vstpInterfaceVIdTypeVlan:
			if routingInstance == junos.DefaultW {
				return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp interface %v in vlan %v already exists",
					name, vlan))...)
			}

			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp interface %v already exists in vlan %s in routing-instance %v",
				name, vlan, routingInstance))...)
		case vstpInterfaceVIdTypeVlanGroup:
			if routingInstance == junos.DefaultW {
				return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp interface %v in vlan-group %v already exists",
					name, vlanGroup))...)
			}

			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp interface %v already exists in vlan-group %v in routing-instance %v",
				name, vlanGroup, routingInstance))...)
		}
	}
	if err := setVstpInterface(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_vstp_interface", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	vstpInterfaceExists, err = checkVstpInterfaceExists(name, routingInstance, vlan, vlanGroup, clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if vstpInterfaceExists {
		d.SetId(resourceVstpInterfaceNewID(d))
	} else {
		switch vstpInterfaceVIdType {
		case vstpInterfaceVIdTypeNone:
			if routingInstance == junos.DefaultW {
				return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp interface %v not exists after commit "+
					"=> check your config", name))...)
			}

			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp interface %v not exists in routing-instance %v after commit "+
					"=> check your config", name, routingInstance))...)
		case vstpInterfaceVIdTypeVlan:
			if routingInstance == junos.DefaultW {
				return append(diagWarns, diag.FromErr(fmt.Errorf(
					"protocols vstp interface %v in vlan %v not exists after commit "+
						"=> check your config", name, vlan))...)
			}

			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp interface %v not exists in vlan %v in routing-instance %v after commit "+
					"=> check your config", name, vlan, routingInstance))...)
		case vstpInterfaceVIdTypeVlanGroup:
			if routingInstance == junos.DefaultW {
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

	return append(diagWarns, resourceVstpInterfaceReadWJunSess(d, clt, junSess)...)
}

func resourceVstpInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceVstpInterfaceReadWJunSess(d, clt, junSess)
}

func resourceVstpInterfaceReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	vstpInterfaceOptions, err := readVstpInterface(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		d.Get("vlan").(string),
		d.Get("vlan_group").(string),
		clt, junSess)
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
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delVstpInterface(
			d.Get("name").(string),
			d.Get("routing_instance").(string),
			d.Get("vlan").(string),
			d.Get("vlan_group").(string),
			clt, nil,
		); err != nil {
			return diag.FromErr(err)
		}
		if err := setVstpInterface(d, clt, nil); err != nil {
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
	if err := delVstpInterface(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		d.Get("vlan").(string),
		d.Get("vlan_group").(string),
		clt, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setVstpInterface(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_vstp_interface", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceVstpInterfaceReadWJunSess(d, clt, junSess)...)
}

func resourceVstpInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delVstpInterface(
			d.Get("name").(string),
			d.Get("routing_instance").(string),
			d.Get("vlan").(string),
			d.Get("vlan_group").(string),
			clt, nil,
		); err != nil {
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
	if err := delVstpInterface(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		d.Get("vlan").(string),
		d.Get("vlan_group").(string),
		clt, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_vstp_interface", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceVstpInterfaceImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)
	if len(strings.Split(d.Id(), junos.IDSeparator)) < 3 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", junos.IDSeparator)
	}
	vType, name, vName, routingInstance := resourceVstpInterfaceReadID(d.Id())
	var vstpInterfaceExists bool
	switch vType {
	case vstpInterfaceVIdTypeNone:
		vstpInterfaceExists, err = checkVstpInterfaceExists(name, routingInstance, "", "", clt, junSess)
	case vstpInterfaceVIdTypeVlan:
		vstpInterfaceExists, err = checkVstpInterfaceExists(name, routingInstance, vName, "", clt, junSess)
	case vstpInterfaceVIdTypeVlanGroup:
		vstpInterfaceExists, err = checkVstpInterfaceExists(name, routingInstance, "", vName, clt, junSess)
	}
	if err != nil {
		return nil, err
	}
	if !vstpInterfaceExists {
		return nil, fmt.Errorf("don't find protocols vstp interface with id '%v' "+
			"(id must be <name>"+junos.IDSeparator+junos.IDSeparator+"<routing_instance>, "+
			"<name>"+junos.IDSeparator+"%s<vlan>"+junos.IDSeparator+"<routing_instance> or "+
			"<name>"+junos.IDSeparator+"%s<vlan_group>"+junos.IDSeparator+"<routing_instance>)",
			d.Id(),
			vstpInterfaceVIdTypeVlan.prefix(),
			vstpInterfaceVIdTypeVlanGroup.prefix(),
		)
	}
	var vstpInterfaceOptions vstpInterfaceOptions
	switch vType {
	case vstpInterfaceVIdTypeNone:
		vstpInterfaceOptions, err = readVstpInterface(name, routingInstance, "", "", clt, junSess)
	case vstpInterfaceVIdTypeVlan:
		vstpInterfaceOptions, err = readVstpInterface(name, routingInstance, vName, "", clt, junSess)
	case vstpInterfaceVIdTypeVlanGroup:
		vstpInterfaceOptions, err = readVstpInterface(name, routingInstance, "", vName, clt, junSess)
	}
	if err != nil {
		return nil, err
	}
	fillVstpInterfaceData(d, vstpInterfaceOptions)
	d.SetId(resourceVstpInterfaceNewID(d))
	result[0] = d

	return result, nil
}

func checkVstpInterfaceExists(name, routingInstance, vlan, vlanGroup string, clt *junos.Client, junSess *junos.Session,
) (_ bool, err error) {
	var showConfig string
	if vlan != "" && vlanGroup != "" {
		return false, fmt.Errorf("internal error: checkVstpInterfaceExists called with vlan and vlanGroup")
	}
	if routingInstance == junos.DefaultW {
		switch {
		case vlan != "":
			showConfig, err = clt.Command(junos.CmdShowConfig+
				"protocols vstp vlan "+vlan+" interface "+name+junos.PipeDisplaySet, junSess)
		case vlanGroup != "":
			showConfig, err = clt.Command(junos.CmdShowConfig+
				"protocols vstp vlan-group group "+vlanGroup+" interface "+name+junos.PipeDisplaySet, junSess)
		default:
			showConfig, err = clt.Command(junos.CmdShowConfig+
				"protocols vstp interface "+name+junos.PipeDisplaySet, junSess)
		}
	} else {
		switch {
		case vlan != "":
			showConfig, err = clt.Command(junos.CmdShowConfig+junos.RoutingInstancesWS+routingInstance+" "+
				"protocols vstp vlan "+vlan+" interface "+name+junos.PipeDisplaySet, junSess)
		case vlanGroup != "":
			showConfig, err = clt.Command(junos.CmdShowConfig+junos.RoutingInstancesWS+routingInstance+" "+
				"protocols vstp vlan-group group "+vlanGroup+" interface "+name+junos.PipeDisplaySet, junSess)
		default:
			showConfig, err = clt.Command(junos.CmdShowConfig+junos.RoutingInstancesWS+routingInstance+" "+
				"protocols vstp interface "+name+junos.PipeDisplaySet, junSess)
		}
	}
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setVstpInterface(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0)

	name := d.Get("name").(string)
	vlan := d.Get("vlan").(string)
	vlanGroup := d.Get("vlan_group").(string)
	setPrefix := junos.SetLS
	if rI := d.Get("routing_instance").(string); rI != junos.DefaultW {
		setPrefix = junos.SetRoutingInstances + rI + " "
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

	return clt.ConfigSet(configSet, junSess)
}

func readVstpInterface(name, routingInstance, vlan, vlanGroup string, clt *junos.Client, junSess *junos.Session,
) (confRead vstpInterfaceOptions, err error) {
	// default -1
	confRead.priority = -1
	if vlan != "" && vlanGroup != "" {
		return confRead, fmt.Errorf("internal error: readVstpInterface called with vlan and vlanGroup")
	}
	var showConfig string
	if routingInstance == junos.DefaultW {
		switch {
		case vlan != "":
			showConfig, err = clt.Command(junos.CmdShowConfig+
				"protocols vstp vlan "+vlan+" interface "+name+junos.PipeDisplaySetRelative, junSess)
		case vlanGroup != "":
			showConfig, err = clt.Command(junos.CmdShowConfig+
				"protocols vstp vlan-group group "+vlanGroup+" interface "+name+junos.PipeDisplaySetRelative, junSess)
		default:
			showConfig, err = clt.Command(junos.CmdShowConfig+
				"protocols vstp interface "+name+junos.PipeDisplaySetRelative, junSess)
		}
	} else {
		switch {
		case vlan != "":
			showConfig, err = clt.Command(junos.CmdShowConfig+junos.RoutingInstancesWS+routingInstance+" "+
				"protocols vstp vlan "+vlan+" interface "+name+junos.PipeDisplaySetRelative, junSess)
		case vlanGroup != "":
			showConfig, err = clt.Command(junos.CmdShowConfig+junos.RoutingInstancesWS+routingInstance+" "+
				"protocols vstp vlan-group group "+vlanGroup+" interface "+name+junos.PipeDisplaySetRelative, junSess)
		default:
			showConfig, err = clt.Command(junos.CmdShowConfig+junos.RoutingInstancesWS+routingInstance+" "+
				"protocols vstp interface "+name+junos.PipeDisplaySetRelative, junSess)
		}
	}
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		confRead.routingInstance = routingInstance
		confRead.vlan = vlan
		confRead.vlanGroup = vlanGroup
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == "access-trunk":
				confRead.accessTrunk = true
			case itemTrim == "bpdu-timeout-action alarm":
				confRead.bpduTimeoutActionAlarm = true
			case itemTrim == "bpdu-timeout-action block":
				confRead.bpduTimeoutActionBlock = true
			case balt.CutPrefixInString(&itemTrim, "cost "):
				confRead.cost, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "edge":
				confRead.edge = true
			case balt.CutPrefixInString(&itemTrim, "mode "):
				confRead.mode = itemTrim
			case itemTrim == "no-root-port":
				confRead.noRootPort = true
			case balt.CutPrefixInString(&itemTrim, "priority "):
				confRead.priority, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delVstpInterface(name, routingInstance, vlan, vlanGroup string, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	if vlan != "" && vlanGroup != "" {
		return fmt.Errorf("internal error: delVstpInterface called with vlan and vlanGroup")
	}
	if routingInstance == junos.DefaultW {
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
			configSet = append(configSet, junos.DelRoutingInstances+routingInstance+
				" protocols vstp vlan "+vlan+" interface "+name)
		case vlanGroup != "":
			configSet = append(configSet, junos.DelRoutingInstances+routingInstance+
				" protocols vstp vlan-group group "+vlanGroup+" interface "+name)
		default:
			configSet = append(configSet, junos.DelRoutingInstances+routingInstance+
				" protocols vstp interface "+name)
		}
	}

	return clt.ConfigSet(configSet, junSess)
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
