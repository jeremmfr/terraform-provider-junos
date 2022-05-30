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

type vstpVlanGroupOptions struct {
	forwardDelay         int
	helloTime            int
	maxAge               int
	backupBridgePriority string
	bridgePriority       string
	name                 string
	routingInstance      string
	systemIdentifier     string
	vlan                 []string
}

func resourceVstpVlanGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVstpVlanGroupCreate,
		ReadWithoutTimeout:   resourceVstpVlanGroupRead,
		UpdateWithoutTimeout: resourceVstpVlanGroupUpdate,
		DeleteWithoutTimeout: resourceVstpVlanGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVstpVlanGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          defaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"vlan": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceVstpVlanGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	routingInstance := d.Get("routing_instance").(string)
	name := d.Get("name").(string)
	if clt.fakeCreateSetFile != "" {
		if err := setVstpVlanGroup(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(name + idSeparator + routingInstance)

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
	if routingInstance != defaultW {
		instanceExists, err := checkRoutingInstanceExists(routingInstance, clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	vstpVlanGroupExists, err := checkVstpVlanGroupExists(name, routingInstance, clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if vstpVlanGroupExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))
		if routingInstance != defaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp vlan-group group %v already exists in routing-instance %v",
				name, routingInstance))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp vlan-group group %v already exists",
			name))...)
	}
	if err := setVstpVlanGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_vstp_vlan_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	vstpVlanGroupExists, err = checkVstpVlanGroupExists(name, routingInstance, clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if vstpVlanGroupExists {
		d.SetId(name + idSeparator + routingInstance)
	} else {
		if routingInstance != defaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp vlan-group group %v not exists in routing-instance %v after commit "+
					"=> check your config", name, routingInstance))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp vlan-group group %v not exists after commit "+
			"=> check your config", name))...)
	}

	return append(diagWarns, resourceVstpVlanGroupReadWJunSess(d, clt, junSess)...)
}

func resourceVstpVlanGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceVstpVlanGroupReadWJunSess(d, clt, junSess)
}

func resourceVstpVlanGroupReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	vstpVlanGroupOptions, err := readVstpVlanGroup(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if vstpVlanGroupOptions.name == "" {
		d.SetId("")
	} else {
		fillVstpVlanGroupData(d, vstpVlanGroupOptions)
	}

	return nil
}

func resourceVstpVlanGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delVstpVlanGroup(
			d.Get("name").(string),
			d.Get("routing_instance").(string),
			false,
			clt, nil,
		); err != nil {
			return diag.FromErr(err)
		}
		if err := setVstpVlanGroup(d, clt, nil); err != nil {
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
	if err := delVstpVlanGroup(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		false,
		clt, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setVstpVlanGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_vstp_vlan_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceVstpVlanGroupReadWJunSess(d, clt, junSess)...)
}

func resourceVstpVlanGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delVstpVlanGroup(
			d.Get("name").(string),
			d.Get("routing_instance").(string),
			true,
			clt, nil,
		); err != nil {
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
	if err := delVstpVlanGroup(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		true,
		clt, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_vstp_vlan_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceVstpVlanGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	vstpVlanGroupExists, err := checkVstpVlanGroupExists(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	if !vstpVlanGroupExists {
		return nil, fmt.Errorf("don't find protocols vstp vlan-group group with id '%v' "+
			"(id must be <name>%s<routing_instance>", d.Id(), idSeparator)
	}
	vstpVlanGroupOptions, err := readVstpVlanGroup(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	fillVstpVlanGroupData(d, vstpVlanGroupOptions)

	result[0] = d

	return result, nil
}

func checkVstpVlanGroupExists(name, routingInstance string, clt *Client, junSess *junosSession,
) (bool, error) {
	var showConfig string
	var err error
	if routingInstance == defaultW {
		showConfig, err = clt.command(cmdShowConfig+
			"protocols vstp vlan-group group "+name+pipeDisplaySet, junSess)
	} else {
		showConfig, err = clt.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols vstp vlan-group group "+name+pipeDisplaySet, junSess)
	}
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setVstpVlanGroup(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := setLS
	if rI := d.Get("routing_instance").(string); rI != defaultW {
		setPrefix = setRoutingInstances + rI + " "
	}
	setPrefix += "protocols vstp vlan-group group " + d.Get("name").(string) + " "

	for _, vlan := range sortSetOfString(d.Get("vlan").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"vlan "+vlan)
	}
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

	return clt.configSet(configSet, junSess)
}

func readVstpVlanGroup(name, routingInstance string, clt *Client, junSess *junosSession,
) (vstpVlanGroupOptions, error) {
	var confRead vstpVlanGroupOptions
	var showConfig string
	var err error
	if routingInstance == defaultW {
		showConfig, err = clt.command(cmdShowConfig+
			"protocols vstp vlan-group group "+name+pipeDisplaySetRelative, junSess)
	} else {
		showConfig, err = clt.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols vstp vlan-group group "+name+pipeDisplaySetRelative, junSess)
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
			case strings.HasPrefix(itemTrim, "vlan "):
				confRead.vlan = append(confRead.vlan, strings.TrimPrefix(itemTrim, "vlan "))
			case strings.HasPrefix(itemTrim, "backup-bridge-priority "):
				confRead.backupBridgePriority = strings.TrimPrefix(itemTrim, "backup-bridge-priority ")
			case strings.HasPrefix(itemTrim, "bridge-priority "):
				confRead.bridgePriority = strings.TrimPrefix(itemTrim, "bridge-priority ")
			case strings.HasPrefix(itemTrim, "forward-delay "):
				var err error
				confRead.forwardDelay, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "forward-delay "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "hello-time "):
				var err error
				confRead.helloTime, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "hello-time "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "max-age "):
				var err error
				confRead.maxAge, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-age "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "system-identifier "):
				confRead.systemIdentifier = strings.TrimPrefix(itemTrim, "system-identifier ")
			}
		}
	}

	return confRead, nil
}

func delVstpVlanGroup(name, routingInstance string, deleteAll bool, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	delPrefix := deleteLS
	if routingInstance != defaultW {
		delPrefix = delRoutingInstances + routingInstance + " "
	}
	delPrefix += "protocols vstp vlan-group group " + name + " "

	if deleteAll {
		return clt.configSet([]string{delPrefix}, junSess)
	}
	listLinesToDelete := []string{
		"backup-bridge-priority",
		"bridge-priority",
		"forward-delay",
		"hello-time",
		"max-age",
		"system-identifier",
		"vlan",
	}
	for _, line := range listLinesToDelete {
		configSet = append(configSet, delPrefix+line)
	}

	return clt.configSet(configSet, junSess)
}

func fillVstpVlanGroupData(d *schema.ResourceData, vstpVlanGroupOptions vstpVlanGroupOptions) {
	if tfErr := d.Set("name", vstpVlanGroupOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", vstpVlanGroupOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vlan", vstpVlanGroupOptions.vlan); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("backup_bridge_priority", vstpVlanGroupOptions.backupBridgePriority); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bridge_priority", vstpVlanGroupOptions.bridgePriority); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("forward_delay", vstpVlanGroupOptions.forwardDelay); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("hello_time", vstpVlanGroupOptions.helloTime); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("max_age", vstpVlanGroupOptions.maxAge); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("system_identifier", vstpVlanGroupOptions.systemIdentifier); tfErr != nil {
		panic(tfErr)
	}
}
