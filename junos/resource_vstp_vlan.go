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
				Default:          defaultW,
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
	sess := m.(*Session)
	routingInstance := d.Get("routing_instance").(string)
	vlanID := d.Get("vlan_id").(string)
	if sess.junosFakeCreateSetFile != "" {
		if err := setVstpVlan(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(vlanID + idSeparator + routingInstance)

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if routingInstance != defaultW {
		instanceExists, err := checkRoutingInstanceExists(routingInstance, m, jnprSess)
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
	vstpVlanExists, err := checkVstpVlanExists(vlanID, routingInstance, m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if vstpVlanExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))
		if routingInstance != defaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp vlan %v already exists in routing-instance %v", vlanID, routingInstance))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp vlan %v already exists", vlanID))...)
	}
	if err := setVstpVlan(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_vstp_vlan", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	vstpVlanExists, err = checkVstpVlanExists(vlanID, routingInstance, m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if vstpVlanExists {
		d.SetId(vlanID + idSeparator + routingInstance)
	} else {
		if routingInstance != defaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"protocols vstp vlan %v not exists in routing-instance %v after commit "+
					"=> check your config", vlanID, routingInstance))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf("protocols vstp vlan %v not exists after commit "+
			"=> check your config", vlanID))...)
	}

	return append(diagWarns, resourceVstpVlanReadWJnprSess(d, m, jnprSess)...)
}

func resourceVstpVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceVstpVlanReadWJnprSess(d, m, jnprSess)
}

func resourceVstpVlanReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
) diag.Diagnostics {
	mutex.Lock()
	vstpVlanOptions, err := readVstpVlan(d.Get("vlan_id").(string), d.Get("routing_instance").(string), m, jnprSess)
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
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delVstpVlan(d.Get("vlan_id").(string), d.Get("routing_instance").(string), false, m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setVstpVlan(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delVstpVlan(d.Get("vlan_id").(string), d.Get("routing_instance").(string), false, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setVstpVlan(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_vstp_vlan", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceVstpVlanReadWJnprSess(d, m, jnprSess)...)
}

func resourceVstpVlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delVstpVlan(d.Get("vlan_id").(string), d.Get("routing_instance").(string), true, m, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delVstpVlan(d.Get("vlan_id").(string), d.Get("routing_instance").(string), true, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_vstp_vlan", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceVstpVlanImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	vstpVlanExists, err := checkVstpVlanExists(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !vstpVlanExists {
		return nil, fmt.Errorf("don't find protocols vstp vlan with id '%v' "+
			"(id must be <vlan_id>%s<routing_instance>", d.Id(), idSeparator)
	}
	vstpVlanOptions, err := readVstpVlan(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillVstpVlanData(d, vstpVlanOptions)

	result[0] = d

	return result, nil
}

func checkVstpVlanExists(vlanID, routingInstance string, m interface{}, jnprSess *NetconfObject,
) (bool, error) {
	sess := m.(*Session)
	var showConfig string
	var err error
	if routingInstance == defaultW {
		showConfig, err = sess.command(cmdShowConfig+
			"protocols vstp vlan "+vlanID+pipeDisplaySet, jnprSess)
	} else {
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols vstp vlan "+vlanID+pipeDisplaySet, jnprSess)
	}
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setVstpVlan(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)

	setPrefix := setLS
	if rI := d.Get("routing_instance").(string); rI != defaultW {
		setPrefix = setRoutingInstances + rI + " "
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

	return sess.configSet(configSet, jnprSess)
}

func readVstpVlan(vlanID, routingInstance string, m interface{}, jnprSess *NetconfObject,
) (vstpVlanOptions, error) {
	sess := m.(*Session)
	var confRead vstpVlanOptions
	var showConfig string
	var err error
	if routingInstance == defaultW {
		showConfig, err = sess.command(cmdShowConfig+
			"protocols vstp vlan "+vlanID+pipeDisplaySetRelative, jnprSess)
	} else {
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols vstp vlan "+vlanID+pipeDisplaySetRelative, jnprSess)
	}
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.vlanID = vlanID
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

func delVstpVlan(vlanID, routingInstance string, deleteAll bool, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	delPrefix := deleteLS
	if routingInstance != defaultW {
		delPrefix = delRoutingInstances + routingInstance + " "
	}
	delPrefix += "protocols vstp vlan " + vlanID + " "

	if deleteAll {
		return sess.configSet([]string{delPrefix}, jnprSess)
	}
	listLinesToDelete := []string{
		"backup-bridge-priority",
		"bridge-priority",
		"forward-delay",
		"hello-time",
		"max-age",
		"system-identifier",
	}
	for _, line := range listLinesToDelete {
		configSet = append(configSet, delPrefix+line)
	}

	return sess.configSet(configSet, jnprSess)
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
