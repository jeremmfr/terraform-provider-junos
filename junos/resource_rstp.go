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
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type rstpOptions struct {
	bpduBlockOnEdge             bool
	bpduDestMACAddProvBridgeGrp bool
	disable                     bool
	forceVersionStp             bool
	vplsFlushOnTopologyChange   bool
	extendedSystemID            int
	forwardDelay                int
	helloTime                   int
	maxAge                      int
	priorityHoldTime            int
	backupBridgePriority        string
	bridgePriority              string
	routingInstance             string
	systemIdentifier            string
	systemID                    []map[string]interface{}
}

func resourceRstp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRstpCreate,
		ReadContext:   resourceRstpRead,
		UpdateContext: resourceRstpUpdate,
		DeleteContext: resourceRstpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRstpImport,
		},
		Schema: map[string]*schema.Schema{
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
			"bpdu_block_on_edge": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bpdu_destination_mac_address_provider_bridge_group": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bridge_priority": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^(0|\d\d?k)$`), "must be a number with increments of 4k - 0,4k,8k,..60k"),
			},
			"disable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"extended_system_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 4095),
				Default:      -1,
			},
			"force_version_stp": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"system_identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsMACAddress,
			},
			"vpls_flush_on_topology_change": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceRstpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setRstp(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("routing_instance").(string))

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
	if err := setRstp(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_rstp", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId(d.Get("routing_instance").(string))

	return append(diagWarns, resourceRstpReadWJnprSess(d, m, jnprSess)...)
}

func resourceRstpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceRstpReadWJnprSess(d, m, jnprSess)
}

func resourceRstpReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	if d.Get("routing_instance").(string) != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), m, jnprSess)
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
	rstpOptions, err := readRstp(d.Get("routing_instance").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillRstpData(d, rstpOptions)

	return nil
}

func resourceRstpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delRstp(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setRstp(d, m, nil); err != nil {
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
	if err := delRstp(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setRstp(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_rstp", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRstpReadWJnprSess(d, m, jnprSess)...)
}

func resourceRstpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delRstp(d, m, nil); err != nil {
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
	if err := delRstp(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_rstp", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceRstpImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	if d.Id() != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Id(), m, jnprSess)
		if err != nil {
			return nil, err
		}
		if !instanceExists {
			return nil, fmt.Errorf("routing instance %v doesn't exist", d.Id())
		}
	}
	result := make([]*schema.ResourceData, 1)
	rstpOptions, err := readRstp(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillRstpData(d, rstpOptions)
	result[0] = d

	return result, nil
}

func setRstp(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := setLS
	if d.Get("routing_instance").(string) != defaultW {
		setPrefix = setRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	setPrefix += "protocols rstp "

	if v := d.Get("backup_bridge_priority").(string); v != "" {
		configSet = append(configSet, setPrefix+"backup-bridge-priority "+v)
	}
	if d.Get("bpdu_block_on_edge").(bool) {
		configSet = append(configSet, setPrefix+"bpdu-block-on-edge")
	}
	if d.Get("bpdu_destination_mac_address_provider_bridge_group").(bool) {
		configSet = append(configSet, setPrefix+"bpdu-destination-mac-address provider-bridge-group")
	}
	if v := d.Get("bridge_priority").(string); v != "" {
		configSet = append(configSet, setPrefix+"bridge-priority "+v)
	}
	if d.Get("disable").(bool) {
		configSet = append(configSet, setPrefix+"disable")
	}
	if v := d.Get("extended_system_id").(int); v != -1 {
		configSet = append(configSet, setPrefix+"extended-system-id "+strconv.Itoa(v))
	}
	if d.Get("force_version_stp").(bool) {
		configSet = append(configSet, setPrefix+"force-version stp")
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
	if v := d.Get("system_identifier").(string); v != "" {
		configSet = append(configSet, setPrefix+"system-identifier "+v)
	}
	if d.Get("vpls_flush_on_topology_change").(bool) {
		configSet = append(configSet, setPrefix+"vpls-flush-on-topology-change")
	}

	return sess.configSet(configSet, jnprSess)
}

func readRstp(routingInstance string, m interface{}, jnprSess *NetconfObject) (rstpOptions, error) {
	sess := m.(*Session)
	var confRead rstpOptions
	confRead.extendedSystemID = -1

	var showConfig string
	if routingInstance == defaultW {
		var err error
		showConfig, err = sess.command(cmdShowConfig+
			"protocols rstp | display set relative", jnprSess)
		if err != nil {
			return confRead, err
		}
	} else {
		var err error
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols rstp | display set relative", jnprSess)
		if err != nil {
			return confRead, err
		}
	}

	confRead.routingInstance = routingInstance
	if showConfig != emptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "backup-bridge-priority "):
				confRead.backupBridgePriority = strings.TrimPrefix(itemTrim, "backup-bridge-priority ")
			case itemTrim == "bpdu-block-on-edge":
				confRead.bpduBlockOnEdge = true
			case itemTrim == "bpdu-destination-mac-address provider-bridge-group":
				confRead.bpduDestMACAddProvBridgeGrp = true
			case strings.HasPrefix(itemTrim, "bridge-priority "):
				confRead.bridgePriority = strings.TrimPrefix(itemTrim, "bridge-priority ")
			case itemTrim == disableW:
				confRead.disable = true
			case strings.HasPrefix(itemTrim, "extended-system-id "):
				var err error
				confRead.extendedSystemID, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "extended-system-id "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case itemTrim == "force-version stp":
				confRead.forceVersionStp = true
			case strings.HasPrefix(itemTrim, "forward-delay "):
				var err error
				confRead.forwardDelay, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "forward-delay "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "hello-time "):
				var err error
				confRead.helloTime, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "hello-time "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "max-age "):
				var err error
				confRead.maxAge, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-age "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "priority-hold-time "):
				var err error
				confRead.priorityHoldTime, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "priority-hold-time "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
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
			case strings.HasPrefix(itemTrim, "system-identifier "):
				confRead.systemIdentifier = strings.TrimPrefix(itemTrim, "system-identifier ")
			case itemTrim == "vpls-flush-on-topology-change":
				confRead.vplsFlushOnTopologyChange = true
			}
		}
	}

	return confRead, nil
}

func delRstp(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	delPrefix := deleteLS
	if d.Get("routing_instance").(string) != defaultW {
		delPrefix = delRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	delPrefix += "protocols rstp "

	listLinesToDelete := []string{
		"backup-bridge-priority",
		"bpdu-block-on-edge",
		"bpdu-destination-mac-address",
		"bridge-priority",
		"disable",
		"extended-system-id",
		"force-version",
		"forward-delay",
		"hello-time",
		"max-age",
		"priority-hold-time",
		"system-identifier",
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

func fillRstpData(d *schema.ResourceData, rstpOptions rstpOptions) {
	if tfErr := d.Set("routing_instance", rstpOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("backup_bridge_priority", rstpOptions.backupBridgePriority); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bpdu_block_on_edge", rstpOptions.bpduBlockOnEdge); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"bpdu_destination_mac_address_provider_bridge_group",
		rstpOptions.bpduDestMACAddProvBridgeGrp,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bridge_priority", rstpOptions.bridgePriority); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("disable", rstpOptions.disable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("extended_system_id", rstpOptions.extendedSystemID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("force_version_stp", rstpOptions.forceVersionStp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("forward_delay", rstpOptions.forwardDelay); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("hello_time", rstpOptions.helloTime); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("max_age", rstpOptions.maxAge); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("priority_hold_time", rstpOptions.priorityHoldTime); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("system_id", rstpOptions.systemID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("system_identifier", rstpOptions.systemIdentifier); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vpls_flush_on_topology_change", rstpOptions.vplsFlushOnTopologyChange); tfErr != nil {
		panic(tfErr)
	}
}
