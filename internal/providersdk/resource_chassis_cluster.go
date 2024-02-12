package providersdk

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type chassisClusterOptions struct {
	configSyncNoSecBootAuto bool
	controlLinkRecovery     bool
	heartbeatInterval       int
	heartbeatThreshold      int
	rethCount               int
	controlPorts            []map[string]interface{}
	redundancyGroup         []map[string]interface{}
	fab0                    []map[string]interface{}
	fab1                    []map[string]interface{}
}

func resourceChassisCluster() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceChassisClusterCreate,
		ReadWithoutTimeout:   resourceChassisClusterRead,
		UpdateWithoutTimeout: resourceChassisClusterUpdate,
		DeleteWithoutTimeout: resourceChassisClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceChassisClusterImport,
		},
		Schema: map[string]*schema.Schema{
			"fab0": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member_interfaces": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
									value := v.(string)
									if strings.Count(value, ".") > 0 {
										errors = append(errors, fmt.Errorf(
											"%q in %q cannot have a dot", value, k))
									}

									return
								},
							},
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"fab1": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member_interfaces": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
									value := v.(string)
									if strings.Count(value, ".") > 0 {
										errors = append(errors, fmt.Errorf(
											"%q in %q cannot have a dot", value, k))
									}

									return
								},
							},
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"redundancy_group": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 128,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node0_priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 254),
						},
						"node1_priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 254),
						},
						"gratuitous_arp_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 16),
						},
						"hold_down_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 1800),
						},
						"interface_monitor": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"weight": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 255),
									},
								},
							},
						},
						"preempt": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"preempt_delay": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 21600),
						},
						"preempt_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 50),
						},
						"preempt_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 1400),
						},
					},
				},
			},
			"reth_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 128),
			},
			"config_sync_no_secondary_bootup_auto": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"control_link_recovery": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"control_ports": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fpc": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 23),
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1),
						},
					},
				},
			},
			"heartbeat_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1000, 2000),
			},
			"heartbeat_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(3, 8),
			},
		},
	}
}

func checkCompatibilityChassisCluster(junSess *junos.Session) bool {
	if strings.HasPrefix(strings.ToLower(junSess.SystemInformation.HardwareModel), "srx") {
		return true
	}
	if strings.HasPrefix(strings.ToLower(junSess.SystemInformation.HardwareModel), "vsrx") {
		return true
	}
	if strings.HasPrefix(strings.ToLower(junSess.SystemInformation.HardwareModel), "j") {
		return true
	}

	return false
}

func resourceChassisClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setChassisCluster(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("cluster")

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if !checkCompatibilityChassisCluster(junSess) {
		return diag.FromErr(fmt.Errorf("chassis cluster "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := setChassisCluster(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_chassis_cluster")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("cluster")

	return append(diagWarns, resourceChassisClusterReadWJunSess(d, junSess)...)
}

func resourceChassisClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceChassisClusterReadWJunSess(d, junSess)
}

func resourceChassisClusterReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	clusterOptions, err := readChassisCluster(junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillChassisCluster(d, clusterOptions)

	return nil
}

func resourceChassisClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delChassisCluster(junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setChassisCluster(d, junSess); err != nil {
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
	if err := delChassisCluster(junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setChassisCluster(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_chassis_cluster")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceChassisClusterReadWJunSess(d, junSess)...)
}

func resourceChassisClusterDelete(ctx context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delChassisCluster(junSess); err != nil {
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
	if err := delChassisCluster(junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_chassis_cluster")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceChassisClusterImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	clusterOptions, err := readChassisCluster(junSess)
	if err != nil {
		return nil, err
	}
	fillChassisCluster(d, clusterOptions)
	d.SetId("cluster")
	result[0] = d

	return result, nil
}

func setChassisCluster(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setIntPrefix := "set interfaces "
	for _, v := range d.Get("fab0").([]interface{}) {
		configSet = append(configSet, "delete interfaces fab0 disable")
		fab0 := v.(map[string]interface{})
		for _, v2 := range fab0["member_interfaces"].([]interface{}) {
			configSet = append(configSet, setIntPrefix+
				"fab0 fabric-options member-interfaces "+v2.(string))
		}
		if fab0["description"].(string) != "" {
			configSet = append(configSet, setIntPrefix+"fab0 description \""+
				fab0["description"].(string)+"\"")
		}
	}
	for _, v := range d.Get("fab1").([]interface{}) {
		configSet = append(configSet, "delete interfaces fab1 disable")
		fab0 := v.(map[string]interface{})
		for _, v2 := range fab0["member_interfaces"].([]interface{}) {
			configSet = append(configSet, setIntPrefix+
				"fab1 fabric-options member-interfaces "+v2.(string))
		}
		if fab0["description"].(string) != "" {
			configSet = append(configSet, setIntPrefix+"fab1 description \""+
				fab0["description"].(string)+"\"")
		}
	}

	setChassisluster := "set chassis cluster "
	for i, v := range d.Get("redundancy_group").([]interface{}) {
		redundancyGroup := v.(map[string]interface{})
		configSet = append(configSet, setChassisluster+"redundancy-group "+strconv.Itoa(i)+
			" node 0 priority "+strconv.Itoa(redundancyGroup["node0_priority"].(int)))
		configSet = append(configSet, setChassisluster+"redundancy-group "+strconv.Itoa(i)+
			" node 1 priority "+strconv.Itoa(redundancyGroup["node1_priority"].(int)))
		if redundancyGroup["gratuitous_arp_count"].(int) != 0 {
			configSet = append(configSet, setChassisluster+"redundancy-group "+strconv.Itoa(i)+
				" gratuitous-arp-count "+strconv.Itoa(redundancyGroup["gratuitous_arp_count"].(int)))
		}
		if redundancyGroup["hold_down_interval"].(int) != -1 {
			configSet = append(configSet, setChassisluster+"redundancy-group "+strconv.Itoa(i)+
				" hold-down-interval "+strconv.Itoa(redundancyGroup["hold_down_interval"].(int)))
		}
		interfaceMonitorNameList := make([]string, 0)
		for _, v2 := range redundancyGroup["interface_monitor"].([]interface{}) {
			interfaceMonitor := v2.(map[string]interface{})
			if slices.Contains(interfaceMonitorNameList, interfaceMonitor["name"].(string)) {
				return fmt.Errorf("multiple blocks interface_monitor with the same name %s", interfaceMonitor["name"].(string))
			}
			interfaceMonitorNameList = append(interfaceMonitorNameList, interfaceMonitor["name"].(string))
			configSet = append(configSet, setChassisluster+"redundancy-group "+strconv.Itoa(i)+
				" interface-monitor "+interfaceMonitor["name"].(string)+
				" weight "+strconv.Itoa(interfaceMonitor["weight"].(int)))
		}
		if redundancyGroup["preempt"].(bool) {
			configSet = append(configSet, setChassisluster+"redundancy-group "+strconv.Itoa(i)+
				" preempt")
			if v2 := redundancyGroup["preempt_delay"].(int); v2 != 0 {
				configSet = append(configSet, setChassisluster+"redundancy-group "+strconv.Itoa(i)+
					" preempt delay "+strconv.Itoa(v2))
			}
			if v2 := redundancyGroup["preempt_limit"].(int); v2 != 0 {
				configSet = append(configSet, setChassisluster+"redundancy-group "+strconv.Itoa(i)+
					" preempt limit "+strconv.Itoa(v2))
			}
			if v2 := redundancyGroup["preempt_period"].(int); v2 != 0 {
				configSet = append(configSet, setChassisluster+"redundancy-group "+strconv.Itoa(i)+
					" preempt period "+strconv.Itoa(v2))
			}
		} else if redundancyGroup["preempt_delay"].(int) != 0 ||
			redundancyGroup["preempt_limit"].(int) != 0 ||
			redundancyGroup["preempt_period"].(int) != 0 {
			return errors.New("preempt need to be true with preempt_(delay|limit|period) arguments")
		}
	}
	configSet = append(configSet, setChassisluster+"reth-count "+
		strconv.Itoa(d.Get("reth_count").(int)))
	if d.Get("config_sync_no_secondary_bootup_auto").(bool) {
		configSet = append(configSet, setChassisluster+"configuration-synchronize no-secondary-bootup-auto")
	}
	if d.Get("control_link_recovery").(bool) {
		configSet = append(configSet, setChassisluster+"control-link-recovery")
	}
	for _, cp := range d.Get("control_ports").(*schema.Set).List() {
		controlPort := cp.(map[string]interface{})
		configSet = append(configSet, setChassisluster+"control-ports fpc "+
			strconv.Itoa(controlPort["fpc"].(int))+" port "+strconv.Itoa(controlPort["port"].(int)))
	}
	if v := d.Get("heartbeat_interval").(int); v != 0 {
		configSet = append(configSet, setChassisluster+"heartbeat-interval "+
			strconv.Itoa(v))
	}
	if v := d.Get("heartbeat_threshold").(int); v != 0 {
		configSet = append(configSet, setChassisluster+"heartbeat-threshold "+
			strconv.Itoa(v))
	}

	return junSess.ConfigSet(configSet)
}

func delChassisCluster(junSess *junos.Session) error {
	listLinesToDelete := make([]string, 0)
	listLinesToDelete = append(listLinesToDelete, "chassis cluster")
	listLinesToDelete = append(listLinesToDelete, "interfaces fab0")
	listLinesToDelete = append(listLinesToDelete, "interfaces fab1")
	configSet := make([]string, len(listLinesToDelete))
	for k, line := range listLinesToDelete {
		configSet[k] = junos.DeleteLS + line
	}

	return junSess.ConfigSet(configSet)
}

func readChassisCluster(junSess *junos.Session,
) (confRead chassisClusterOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "chassis cluster" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
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
			case balt.CutPrefixInString(&itemTrim, "redundancy-group "):
				itemTrimFields := strings.Split(itemTrim, " ")
				number, err := strconv.Atoi(itemTrimFields[0])
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrimFields[0], err)
				}
				if len(confRead.redundancyGroup) < number+1 {
					for i := len(confRead.redundancyGroup); i < number+1; i++ {
						confRead.redundancyGroup = append(confRead.redundancyGroup, map[string]interface{}{
							"node0_priority":       0,
							"node1_priority":       0,
							"gratuitous_arp_count": 0,
							"hold_down_interval":   -1,
							"interface_monitor":    make([]map[string]interface{}, 0),
							"preempt":              false,
							"preempt_delay":        0,
							"preempt_limit":        0,
							"preempt_period":       0,
						})
					}
				}
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "node 0 priority "):
					confRead.redundancyGroup[number]["node0_priority"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "node 1 priority "):
					confRead.redundancyGroup[number]["node1_priority"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "gratuitous-arp-count "):
					confRead.redundancyGroup[number]["gratuitous_arp_count"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "hold-down-interval "):
					confRead.redundancyGroup[number]["hold_down_interval"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "interface-monitor "):
					ifaceMonitorFields := strings.Split(itemTrim, " ")
					if len(ifaceMonitorFields) < 3 { // <name> weight <weight>
						return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "interface-monitor", itemTrim)
					}
					weight, err := strconv.Atoi(ifaceMonitorFields[2])
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
					confRead.redundancyGroup[number]["interface_monitor"] = append(
						confRead.redundancyGroup[number]["interface_monitor"].([]map[string]interface{}),
						map[string]interface{}{
							"name":   ifaceMonitorFields[0],
							"weight": weight,
						})
				case balt.CutPrefixInString(&itemTrim, "preempt"):
					confRead.redundancyGroup[number]["preempt"] = true
					switch {
					case balt.CutPrefixInString(&itemTrim, " delay "):
						confRead.redundancyGroup[number]["preempt_delay"], err = strconv.Atoi(itemTrim)
						if err != nil {
							return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
						}
					case balt.CutPrefixInString(&itemTrim, " limit "):
						confRead.redundancyGroup[number]["preempt_limit"], err = strconv.Atoi(itemTrim)
						if err != nil {
							return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
						}
					case balt.CutPrefixInString(&itemTrim, " period "):
						confRead.redundancyGroup[number]["preempt_period"], err = strconv.Atoi(itemTrim)
						if err != nil {
							return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
						}
					}
				}
			case balt.CutPrefixInString(&itemTrim, "reth-count "):
				confRead.rethCount, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "configuration-synchronize no-secondary-bootup-auto":
				confRead.configSyncNoSecBootAuto = true
			case balt.CutPrefixInString(&itemTrim, "control-ports fpc "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 3 { // <fpc> port <port>
					return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "control-ports fpc", itemTrim)
				}
				controlPort := make(map[string]interface{})
				controlPort["fpc"], err = strconv.Atoi(itemTrimFields[0])
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrimFields[0], err)
				}
				controlPort["port"], err = strconv.Atoi(itemTrimFields[2])
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrimFields[2], err)
				}
				confRead.controlPorts = append(confRead.controlPorts, controlPort)
			case itemTrim == "control-link-recovery":
				confRead.controlLinkRecovery = true
			case balt.CutPrefixInString(&itemTrim, "heartbeat-interval "):
				confRead.heartbeatInterval, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "heartbeat-threshold "):
				confRead.heartbeatThreshold, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}
	showConfigFab0, err := junSess.Command(junos.CmdShowConfig + "interfaces fab0" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfigFab0 != junos.EmptyW {
		if len(confRead.fab0) == 0 {
			confRead.fab0 = append(confRead.fab0, map[string]interface{}{
				"member_interfaces": make([]string, 0),
				"description":       "",
			})
		}
		for _, item := range strings.Split(showConfigFab0, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.fab0[0]["description"] = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "fabric-options member-interfaces "):
				confRead.fab0[0]["member_interfaces"] = append(confRead.fab0[0]["member_interfaces"].([]string), itemTrim)
			}
		}
	}
	showConfigFab1, err := junSess.Command(junos.CmdShowConfig + "interfaces fab1" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfigFab1 != junos.EmptyW {
		if len(confRead.fab1) == 0 {
			confRead.fab1 = append(confRead.fab1, map[string]interface{}{
				"member_interfaces": make([]string, 0),
				"description":       "",
			})
		}
		for _, item := range strings.Split(showConfigFab1, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.fab1[0]["description"] = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "fabric-options member-interfaces "):
				confRead.fab1[0]["member_interfaces"] = append(confRead.fab1[0]["member_interfaces"].([]string), itemTrim)
			}
		}
	}

	return confRead, nil
}

func fillChassisCluster(d *schema.ResourceData, chassisClusterOptions chassisClusterOptions) {
	if tfErr := d.Set("fab0", chassisClusterOptions.fab0); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("fab1", chassisClusterOptions.fab1); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("redundancy_group", chassisClusterOptions.redundancyGroup); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("reth_count", chassisClusterOptions.rethCount); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("config_sync_no_secondary_bootup_auto",
		chassisClusterOptions.configSyncNoSecBootAuto); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("control_ports", chassisClusterOptions.controlPorts); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("control_link_recovery", chassisClusterOptions.controlLinkRecovery); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("heartbeat_interval", chassisClusterOptions.heartbeatInterval); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("heartbeat_threshold", chassisClusterOptions.heartbeatThreshold); tfErr != nil {
		panic(tfErr)
	}
}
