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

type vstpVlanOptions struct {
	forwardDelay         int
	helloTime            int
	maxAge               int
	backupBridgePriority string
	bridgePriority       string
	vlanID               string
	routingInstance      string
	systemIdentifier     string
}

func resourceVstpVlan() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVstpVlanCreate,
		ReadWithoutTimeout:   resourceVstpVlanRead,
		UpdateWithoutTimeout: resourceVstpVlanUpdate,
		DeleteWithoutTimeout: resourceVstpVlanDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVstpVlanImport,
		},
		Schema: map[string]*schema.Schema{
			"vlan_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^all|[0-9]{1,4}$`), "must be 'all' or a VLAN id"),
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          junos.DefaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"backup_bridge_priority": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^\d\d?k$`), "must be a number with increments of 4k - 4k,8k,..60k"),
			},
			"bridge_priority": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^(0|\d\d?k)$`), "must be a number with increments of 4k - 0,4k,8k,..60k"),
			},
			"forward_delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(4, 30),
			},
			"hello_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"max_age": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(6, 40),
			},
			"system_identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsMACAddress,
			},
		},
	}
}

func resourceVstpVlanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	routingInstance := d.Get("routing_instance").(string)
	vlanID := d.Get("vlan_id").(string)
	if clt.FakeCreateSetFile() {
		if err := setVstpVlan(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(vlanID + junos.IDSeparator + routingInstance)

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
	if routingInstance != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(routingInstance, clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	vstpVlanExists, err := checkVstpVlanExists(vlanID, routingInstance, clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if vstpVlanExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))
		if routingInstance != junos.DefaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp vlan %v already exists in routing-instance %v", vlanID, routingInstance))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp vlan %v already exists", vlanID))...)
	}
	if err := setVstpVlan(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_vstp_vlan", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	vstpVlanExists, err = checkVstpVlanExists(vlanID, routingInstance, clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if vstpVlanExists {
		d.SetId(vlanID + junos.IDSeparator + routingInstance)
	} else {
		if routingInstance != junos.DefaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp vlan %v not exists in routing-instance %v after commit "+
					"=> check your config", vlanID, routingInstance))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp vlan %v not exists after commit "+
			"=> check your config", vlanID))...)
	}

	return append(diagWarns, resourceVstpVlanReadWJunSess(d, clt, junSess)...)
}

func resourceVstpVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceVstpVlanReadWJunSess(d, clt, junSess)
}

func resourceVstpVlanReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	vstpVlanOptions, err := readVstpVlan(d.Get("vlan_id").(string), d.Get("routing_instance").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if vstpVlanOptions.vlanID == "" {
		d.SetId("")
	} else {
		fillVstpVlanData(d, vstpVlanOptions)
	}

	return nil
}

func resourceVstpVlanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delVstpVlan(d.Get("vlan_id").(string), d.Get("routing_instance").(string), false, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setVstpVlan(d, clt, nil); err != nil {
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
	if err := delVstpVlan(
		d.Get("vlan_id").(string),
		d.Get("routing_instance").(string),
		false,
		clt, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setVstpVlan(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_vstp_vlan", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceVstpVlanReadWJunSess(d, clt, junSess)...)
}

func resourceVstpVlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delVstpVlan(d.Get("vlan_id").(string), d.Get("routing_instance").(string), true, clt, nil); err != nil {
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
	if err := delVstpVlan(
		d.Get("vlan_id").(string),
		d.Get("routing_instance").(string),
		true,
		clt, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_vstp_vlan", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceVstpVlanImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), junos.IDSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", junos.IDSeparator)
	}
	vstpVlanExists, err := checkVstpVlanExists(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	if !vstpVlanExists {
		return nil, fmt.Errorf("don't find protocols vstp vlan with id '%v' "+
			"(id must be <vlan_id>"+junos.IDSeparator+"<routing_instance>", d.Id())
	}
	vstpVlanOptions, err := readVstpVlan(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	fillVstpVlanData(d, vstpVlanOptions)

	result[0] = d

	return result, nil
}

func checkVstpVlanExists(vlanID, routingInstance string, clt *junos.Client, junSess *junos.Session,
) (_ bool, err error) {
	var showConfig string
	if routingInstance == junos.DefaultW {
		showConfig, err = clt.Command(junos.CmdShowConfig+
			"protocols vstp vlan "+vlanID+junos.PipeDisplaySet, junSess)
	} else {
		showConfig, err = clt.Command(junos.CmdShowConfig+junos.RoutingInstancesWS+routingInstance+" "+
			"protocols vstp vlan "+vlanID+junos.PipeDisplaySet, junSess)
	}
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setVstpVlan(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)

	setPrefix := junos.SetLS
	if rI := d.Get("routing_instance").(string); rI != junos.DefaultW {
		setPrefix = junos.SetRoutingInstances + rI + " "
	}
	setPrefix += "protocols vstp vlan " + d.Get("vlan_id").(string) + " "

	configSet = append(configSet, setPrefix)
	if v := d.Get("backup_bridge_priority").(string); v != "" {
		configSet = append(configSet, setPrefix+"backup-bridge-priority "+v)
	}
	if v := d.Get("bridge_priority").(string); v != "" {
		configSet = append(configSet, setPrefix+"bridge-priority "+v)
	}
	if v := d.Get("forward_delay").(int); v != 0 {
		configSet = append(configSet, setPrefix+"forward-delay "+strconv.Itoa(v))
	}
	if v := d.Get("hello_time").(int); v != 0 {
		configSet = append(configSet, setPrefix+"hello-time "+strconv.Itoa(v))
	}
	if v := d.Get("max_age").(int); v != 0 {
		configSet = append(configSet, setPrefix+"max-age "+strconv.Itoa(v))
	}
	if v := d.Get("system_identifier").(string); v != "" {
		configSet = append(configSet, setPrefix+"system-identifier "+v)
	}

	return clt.ConfigSet(configSet, junSess)
}

func readVstpVlan(vlanID, routingInstance string, clt *junos.Client, junSess *junos.Session,
) (confRead vstpVlanOptions, err error) {
	var showConfig string
	if routingInstance == junos.DefaultW {
		showConfig, err = clt.Command(junos.CmdShowConfig+
			"protocols vstp vlan "+vlanID+junos.PipeDisplaySetRelative, junSess)
	} else {
		showConfig, err = clt.Command(junos.CmdShowConfig+junos.RoutingInstancesWS+routingInstance+" "+
			"protocols vstp vlan "+vlanID+junos.PipeDisplaySetRelative, junSess)
	}
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.vlanID = vlanID
		confRead.routingInstance = routingInstance
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "backup-bridge-priority "):
				confRead.backupBridgePriority = itemTrim
			case balt.CutPrefixInString(&itemTrim, "bridge-priority "):
				confRead.bridgePriority = itemTrim
			case balt.CutPrefixInString(&itemTrim, "forward-delay "):
				confRead.forwardDelay, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "hello-time "):
				confRead.helloTime, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "max-age "):
				confRead.maxAge, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "system-identifier "):
				confRead.systemIdentifier = itemTrim
			}
		}
	}

	return confRead, nil
}

func delVstpVlan(vlanID, routingInstance string, deleteAll bool, clt *junos.Client, junSess *junos.Session) error {
	delPrefix := junos.DeleteLS
	if routingInstance != junos.DefaultW {
		delPrefix = junos.DelRoutingInstances + routingInstance + " "
	}
	delPrefix += "protocols vstp vlan " + vlanID + " "

	if deleteAll {
		return clt.ConfigSet([]string{delPrefix}, junSess)
	}
	listLinesToDelete := []string{
		"backup-bridge-priority",
		"bridge-priority",
		"forward-delay",
		"hello-time",
		"max-age",
		"system-identifier",
	}
	configSet := make([]string, len(listLinesToDelete))
	for k, line := range listLinesToDelete {
		configSet[k] = delPrefix + line
	}

	return clt.ConfigSet(configSet, junSess)
}

func fillVstpVlanData(d *schema.ResourceData, vstpVlanOptions vstpVlanOptions) {
	if tfErr := d.Set("vlan_id", vstpVlanOptions.vlanID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", vstpVlanOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("backup_bridge_priority", vstpVlanOptions.backupBridgePriority); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bridge_priority", vstpVlanOptions.bridgePriority); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("forward_delay", vstpVlanOptions.forwardDelay); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("hello_time", vstpVlanOptions.helloTime); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("max_age", vstpVlanOptions.maxAge); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("system_identifier", vstpVlanOptions.systemIdentifier); tfErr != nil {
		panic(tfErr)
	}
}
