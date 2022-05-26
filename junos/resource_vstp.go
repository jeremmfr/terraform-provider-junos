package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type vstpOptions struct {
	bpduBlockOnEdge           bool
	disable                   bool
	forceVersionStp           bool
	vplsFlushOnTopologyChange bool
	priorityHoldTime          int
	routingInstance           string
	systemID                  []map[string]interface{}
}

func resourceVstp() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVstpCreate,
		ReadWithoutTimeout:   resourceVstpRead,
		UpdateWithoutTimeout: resourceVstpUpdate,
		DeleteWithoutTimeout: resourceVstpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVstpImport,
		},
		Schema: map[string]*schema.Schema{
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          defaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"bpdu_block_on_edge": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"force_version_stp": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"priority_hold_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 255),
			},
			"system_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsMACAddress,
						},
						"ip_address": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "",
							ValidateFunc: validation.IsCIDR,
						},
					},
				},
			},
			"vpls_flush_on_topology_change": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceVstpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setVstp(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("routing_instance").(string))

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
	if d.Get("routing_instance").(string) != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), sess, jnprSess)
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
	if err := setVstp(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_vstp", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId(d.Get("routing_instance").(string))

	return append(diagWarns, resourceVstpReadWJnprSess(d, sess, jnprSess)...)
}

func resourceVstpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceVstpReadWJnprSess(d, sess, jnprSess)
}

func resourceVstpReadWJnprSess(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	if d.Get("routing_instance").(string) != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), sess, jnprSess)
		if err != nil {
			mutex.Unlock()

			return diag.FromErr(err)
		}
		if !instanceExists {
			mutex.Unlock()
			d.SetId("")

			return nil
		}
	}
	vstpOptions, err := readVstp(d.Get("routing_instance").(string), sess, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillVstpData(d, vstpOptions)

	return nil
}

func resourceVstpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delVstp(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setVstp(d, sess, nil); err != nil {
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
	if err := delVstp(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setVstp(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_vstp", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceVstpReadWJnprSess(d, sess, jnprSess)...)
}

func resourceVstpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delVstp(d, sess, nil); err != nil {
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
	if err := delVstp(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_vstp", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceVstpImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	if d.Id() != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Id(), sess, jnprSess)
		if err != nil {
			return nil, err
		}
		if !instanceExists {
			return nil, fmt.Errorf("routing instance %v doesn't exist", d.Id())
		}
	}
	result := make([]*schema.ResourceData, 1)
	vstpOptions, err := readVstp(d.Id(), sess, jnprSess)
	if err != nil {
		return nil, err
	}
	fillVstpData(d, vstpOptions)
	result[0] = d

	return result, nil
}

func setVstp(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject) error {
	configSet := make([]string, 0)
	setPrefix := setLS
	if d.Get("routing_instance").(string) != defaultW {
		setPrefix = setRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	setPrefix += "protocols vstp "

	if d.Get("bpdu_block_on_edge").(bool) {
		configSet = append(configSet, setPrefix+"bpdu-block-on-edge")
	}
	if d.Get("disable").(bool) {
		configSet = append(configSet, setPrefix+"disable")
	}
	if d.Get("force_version_stp").(bool) {
		configSet = append(configSet, setPrefix+"force-version stp")
	}
	if v := d.Get("priority_hold_time").(int); v != 0 {
		configSet = append(configSet, setPrefix+"priority-hold-time "+strconv.Itoa(v))
	}
	systemIDList := make([]string, 0)
	for _, mSysID := range d.Get("system_id").(*schema.Set).List() {
		systemID := mSysID.(map[string]interface{})
		if bchk.StringInSlice(systemID["id"].(string), systemIDList) {
			return fmt.Errorf("multiple blocks system_id with the same id '%s'", systemID["id"].(string))
		}
		systemIDList = append(systemIDList, systemID["id"].(string))
		configSet = append(configSet, setPrefix+"system-id "+systemID["id"].(string))
		if ipAdd := systemID["ip_address"].(string); ipAdd != "" {
			configSet = append(configSet, setPrefix+"system-id "+systemID["id"].(string)+" ip-address "+ipAdd)
		}
	}
	if d.Get("vpls_flush_on_topology_change").(bool) {
		configSet = append(configSet, setPrefix+"vpls-flush-on-topology-change")
	}

	return sess.configSet(configSet, jnprSess)
}

func readVstp(routingInstance string, sess *Session, jnprSess *NetconfObject) (vstpOptions, error) {
	var confRead vstpOptions

	var showConfig string
	if routingInstance == defaultW {
		var err error
		showConfig, err = sess.command(cmdShowConfig+
			"protocols vstp"+pipeDisplaySetRelative, jnprSess)
		if err != nil {
			return confRead, err
		}
	} else {
		var err error
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols vstp"+pipeDisplaySetRelative, jnprSess)
		if err != nil {
			return confRead, err
		}
	}
	confRead.routingInstance = routingInstance
	if showConfig != emptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case itemTrim == "bpdu-block-on-edge":
				confRead.bpduBlockOnEdge = true
			case itemTrim == disableW:
				confRead.disable = true
			case itemTrim == "force-version stp":
				confRead.forceVersionStp = true
			case strings.HasPrefix(itemTrim, "priority-hold-time "):
				var err error
				confRead.priorityHoldTime, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "priority-hold-time "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "system-id "):
				itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "system-id "), " ")
				switch len(itemTrimSplit) {
				case 1:
					confRead.systemID = append(confRead.systemID, map[string]interface{}{
						"id":         itemTrimSplit[0],
						"ip_address": "",
					})
				case 3:
					confRead.systemID = append(confRead.systemID, map[string]interface{}{
						"id":         itemTrimSplit[0],
						"ip_address": itemTrimSplit[2],
					})
				default:
					return confRead, fmt.Errorf("can't read value for system_id in '%s'", itemTrim)
				}
			case itemTrim == "vpls-flush-on-topology-change":
				confRead.vplsFlushOnTopologyChange = true
			}
		}
	}

	return confRead, nil
}

func delVstp(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject) error {
	configSet := make([]string, 0, 1)
	delPrefix := deleteLS
	if d.Get("routing_instance").(string) != defaultW {
		delPrefix = delRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	delPrefix += "protocols vstp "

	listLinesToDelete := []string{
		"bpdu-block-on-edge",
		"disable",
		"force-version",
		"priority-hold-time",
		"vpls-flush-on-topology-change",
	}

	for _, line := range listLinesToDelete {
		configSet = append(configSet, delPrefix+line)
	}
	if d.HasChange("system_id") {
		oSysID, _ := d.GetChange("system_id")
		for _, mSysID := range oSysID.(*schema.Set).List() {
			systemID := mSysID.(map[string]interface{})
			configSet = append(configSet, delPrefix+"system-id "+systemID["id"].(string))
		}
	} else {
		for _, mSysID := range d.Get("system_id").(*schema.Set).List() {
			systemID := mSysID.(map[string]interface{})
			configSet = append(configSet, delPrefix+"system-id "+systemID["id"].(string))
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func fillVstpData(d *schema.ResourceData, vstpOptions vstpOptions) {
	if tfErr := d.Set("routing_instance", vstpOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bpdu_block_on_edge", vstpOptions.bpduBlockOnEdge); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("disable", vstpOptions.disable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("force_version_stp", vstpOptions.forceVersionStp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("priority_hold_time", vstpOptions.priorityHoldTime); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("system_id", vstpOptions.systemID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vpls_flush_on_topology_change", vstpOptions.vplsFlushOnTopologyChange); tfErr != nil {
		panic(tfErr)
	}
}
