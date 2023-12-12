package providersdk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
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
		CreateWithoutTimeout: resourceRstpCreate,
		ReadWithoutTimeout:   resourceRstpRead,
		UpdateWithoutTimeout: resourceRstpUpdate,
		DeleteWithoutTimeout: resourceRstpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRstpImport,
		},
		Schema: map[string]*schema.Schema{
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setRstp(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("routing_instance").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if d.Get("routing_instance").(string) != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	if err := setRstp(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_rstp")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId(d.Get("routing_instance").(string))

	return append(diagWarns, resourceRstpReadWJunSess(d, junSess)...)
}

func resourceRstpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceRstpReadWJunSess(d, junSess)
}

func resourceRstpReadWJunSess(d *schema.ResourceData, junSess *junos.Session) diag.Diagnostics {
	junos.MutexLock()
	if d.Get("routing_instance").(string) != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), junSess)
		if err != nil {
			junos.MutexUnlock()

			return diag.FromErr(err)
		}
		if !instanceExists {
			junos.MutexUnlock()
			d.SetId("")

			return nil
		}
	}
	rstpOptions, err := readRstp(d.Get("routing_instance").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillRstpData(d, rstpOptions)

	return nil
}

func resourceRstpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delRstp(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setRstp(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delRstp(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setRstp(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_rstp")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRstpReadWJunSess(d, junSess)...)
}

func resourceRstpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delRstp(d, junSess); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delRstp(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_rstp")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceRstpImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	if d.Id() != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Id(), junSess)
		if err != nil {
			return nil, err
		}
		if !instanceExists {
			return nil, fmt.Errorf("routing instance %v doesn't exist", d.Id())
		}
	}
	result := make([]*schema.ResourceData, 1)
	rstpOptions, err := readRstp(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillRstpData(d, rstpOptions)
	result[0] = d

	return result, nil
}

func setRstp(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	if d.Get("routing_instance").(string) != junos.DefaultW {
		setPrefix = junos.SetRoutingInstances + d.Get("routing_instance").(string) + " "
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
		if bchk.InSlice(systemID["id"].(string), systemIDList) {
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

	return junSess.ConfigSet(configSet)
}

func readRstp(routingInstance string, junSess *junos.Session,
) (confRead rstpOptions, err error) {
	// default -1
	confRead.extendedSystemID = -1
	var showConfig string
	if routingInstance == junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"protocols rstp" + junos.PipeDisplaySetRelative)
		if err != nil {
			return confRead, err
		}
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			"protocols rstp" + junos.PipeDisplaySetRelative)
		if err != nil {
			return confRead, err
		}
	}

	confRead.routingInstance = routingInstance
	if showConfig != junos.EmptyW {
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
			case itemTrim == "bpdu-block-on-edge":
				confRead.bpduBlockOnEdge = true
			case itemTrim == "bpdu-destination-mac-address provider-bridge-group":
				confRead.bpduDestMACAddProvBridgeGrp = true
			case balt.CutPrefixInString(&itemTrim, "bridge-priority "):
				confRead.bridgePriority = itemTrim
			case itemTrim == junos.DisableW:
				confRead.disable = true
			case balt.CutPrefixInString(&itemTrim, "extended-system-id "):
				confRead.extendedSystemID, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "force-version stp":
				confRead.forceVersionStp = true
			case balt.CutPrefixInString(&itemTrim, "forward-delay "):
				confRead.forwardDelay, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "hello-time "):
				confRead.helloTime, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "max-age "):
				confRead.maxAge, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "priority-hold-time "):
				confRead.priorityHoldTime, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "system-id "):
				itemTrimFields := strings.Split(itemTrim, " ")
				switch len(itemTrimFields) { // <id> (ip-address <ip_address>)?
				case 1:
					confRead.systemID = append(confRead.systemID, map[string]interface{}{
						"id":         itemTrimFields[0],
						"ip_address": "",
					})
				case 3:
					confRead.systemID = append(confRead.systemID, map[string]interface{}{
						"id":         itemTrimFields[0],
						"ip_address": itemTrimFields[2],
					})
				default:
					return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "system-id", itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "system-identifier "):
				confRead.systemIdentifier = itemTrim
			case itemTrim == "vpls-flush-on-topology-change":
				confRead.vplsFlushOnTopologyChange = true
			}
		}
	}

	return confRead, nil
}

func delRstp(d *schema.ResourceData, junSess *junos.Session) error {
	delPrefix := junos.DeleteLS
	if d.Get("routing_instance").(string) != junos.DefaultW {
		delPrefix = junos.DelRoutingInstances + d.Get("routing_instance").(string) + " "
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
	configSet := make([]string,
		len(listLinesToDelete), len(listLinesToDelete)+len(d.Get("system_id").(*schema.Set).List()))
	for k, line := range listLinesToDelete {
		configSet[k] = delPrefix + line
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

	return junSess.ConfigSet(configSet)
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
