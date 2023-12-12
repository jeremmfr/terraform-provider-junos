package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
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
				Default:          junos.DefaultW,
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setVstp(d, junSess); err != nil {
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
	if err := setVstp(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_vstp")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId(d.Get("routing_instance").(string))

	return append(diagWarns, resourceVstpReadWJunSess(d, junSess)...)
}

func resourceVstpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceVstpReadWJunSess(d, junSess)
}

func resourceVstpReadWJunSess(d *schema.ResourceData, junSess *junos.Session) diag.Diagnostics {
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
	vstpOptions, err := readVstp(d.Get("routing_instance").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillVstpData(d, vstpOptions)

	return nil
}

func resourceVstpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delVstp(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setVstp(d, junSess); err != nil {
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
	if err := delVstp(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setVstp(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_vstp")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceVstpReadWJunSess(d, junSess)...)
}

func resourceVstpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delVstp(d, junSess); err != nil {
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
	if err := delVstp(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_vstp")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceVstpImport(ctx context.Context, d *schema.ResourceData, m interface{},
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
	vstpOptions, err := readVstp(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillVstpData(d, vstpOptions)
	result[0] = d

	return result, nil
}

func setVstp(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	if d.Get("routing_instance").(string) != junos.DefaultW {
		setPrefix = junos.SetRoutingInstances + d.Get("routing_instance").(string) + " "
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
		if bchk.InSlice(systemID["id"].(string), systemIDList) {
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

	return junSess.ConfigSet(configSet)
}

func readVstp(routingInstance string, junSess *junos.Session,
) (confRead vstpOptions, err error) {
	var showConfig string
	if routingInstance == junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"protocols vstp" + junos.PipeDisplaySetRelative)
		if err != nil {
			return confRead, err
		}
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			"protocols vstp" + junos.PipeDisplaySetRelative)
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
			case itemTrim == "bpdu-block-on-edge":
				confRead.bpduBlockOnEdge = true
			case itemTrim == junos.DisableW:
				confRead.disable = true
			case itemTrim == "force-version stp":
				confRead.forceVersionStp = true
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
			case itemTrim == "vpls-flush-on-topology-change":
				confRead.vplsFlushOnTopologyChange = true
			}
		}
	}

	return confRead, nil
}

func delVstp(d *schema.ResourceData, junSess *junos.Session) error {
	delPrefix := junos.DeleteLS
	if d.Get("routing_instance").(string) != junos.DefaultW {
		delPrefix = junos.DelRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	delPrefix += "protocols vstp "

	listLinesToDelete := []string{
		"bpdu-block-on-edge",
		"disable",
		"force-version",
		"priority-hold-time",
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
