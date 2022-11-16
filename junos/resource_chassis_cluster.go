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

func checkCompatibilityChassisCluster(junSess *junosSession) bool {
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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setChassisCluster(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("cluster")

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if !checkCompatibilityChassisCluster(junSess) {
		return diag.FromErr(fmt.Errorf("chassis cluster "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := setChassisCluster(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_chassis_cluster", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("cluster")

	return append(diagWarns, resourceChassisClusterReadWJunSess(d, clt, junSess)...)
}

func resourceChassisClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceChassisClusterReadWJunSess(d, clt, junSess)
}

func resourceChassisClusterReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	clusterOptions, err := readChassisCluster(clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillChassisCluster(d, clusterOptions)

	return nil
}

func resourceChassisClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delChassisCluster(clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setChassisCluster(d, clt, nil); err != nil {
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
	if err := delChassisCluster(clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setChassisCluster(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_chassis_cluster", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceChassisClusterReadWJunSess(d, clt, junSess)...)
}

func resourceChassisClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delChassisCluster(clt, nil); err != nil {
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
	if err := delChassisCluster(clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_chassis_cluster", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceChassisClusterImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	clusterOptions, err := readChassisCluster(clt, junSess)
	if err != nil {
		return nil, err
	}
	fillChassisCluster(d, clusterOptions)
	d.SetId("cluster")
	result[0] = d

	return result, nil
}

func setChassisCluster(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
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
			if bchk.InSlice(interfaceMonitor["name"].(string), interfaceMonitorNameList) {
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
			return fmt.Errorf("preempt need to be true with preempt_(delay|limit|period) arguments")
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

	return clt.configSet(configSet, junSess)
}

func delChassisCluster(clt *Client, junSess *junosSession) error {
	listLinesToDelete := make([]string, 0)
	listLinesToDelete = append(listLinesToDelete, "chassis cluster")
	listLinesToDelete = append(listLinesToDelete, "interfaces fab0")
	listLinesToDelete = append(listLinesToDelete, "interfaces fab1")
	configSet := make([]string, len(listLinesToDelete))
	for k, line := range listLinesToDelete {
		configSet[k] = deleteLS + line
	}

	return clt.configSet(configSet, junSess)
}

func readChassisCluster(clt *Client, junSess *junosSession) (chassisClusterOptions, error) {
	var confRead chassisClusterOptions

	showConfig, err := clt.command(cmdShowConfig+"chassis cluster"+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
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
			case strings.HasPrefix(itemTrim, "redundancy-group "):
				itemRGTrim := strings.TrimPrefix(itemTrim, "redundancy-group ")
				itemRGTrimSplit := strings.Split(itemRGTrim, " ")
				number, err := strconv.Atoi(itemRGTrimSplit[0])
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemRGTrimSplit[0], err)
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
				itemRGNbTrim := strings.TrimPrefix(itemRGTrim, strconv.Itoa(number)+" ")
				switch {
				case strings.HasPrefix(itemRGNbTrim, "node 0 priority "):
					var err error
					confRead.redundancyGroup[number]["node0_priority"], err = strconv.Atoi(strings.TrimPrefix(
						itemRGNbTrim, "node 0 priority "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemRGNbTrim, err)
					}
				case strings.HasPrefix(itemRGNbTrim, "node 1 priority "):
					var err error
					confRead.redundancyGroup[number]["node1_priority"], err = strconv.Atoi(strings.TrimPrefix(
						itemRGNbTrim, "node 1 priority "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemRGNbTrim, err)
					}
				case strings.HasPrefix(itemRGNbTrim, "gratuitous-arp-count "):
					var err error
					confRead.redundancyGroup[number]["gratuitous_arp_count"], err = strconv.Atoi(strings.TrimPrefix(
						itemRGNbTrim, "gratuitous-arp-count "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemRGNbTrim, err)
					}
				case strings.HasPrefix(itemRGNbTrim, "hold-down-interval "):
					var err error
					confRead.redundancyGroup[number]["hold_down_interval"], err = strconv.Atoi(strings.TrimPrefix(
						itemRGNbTrim, "hold-down-interval "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemRGNbTrim, err)
					}
				case strings.HasPrefix(itemRGNbTrim, "interface-monitor "):
					name := strings.Split(strings.TrimPrefix(itemRGNbTrim, "interface-monitor "), " ")[0]
					weight, err := strconv.Atoi(strings.TrimPrefix(itemRGNbTrim, "interface-monitor "+name+" weight "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemRGNbTrim, err)
					}
					confRead.redundancyGroup[number]["interface_monitor"] = append(
						confRead.redundancyGroup[number]["interface_monitor"].([]map[string]interface{}),
						map[string]interface{}{
							"name":   name,
							"weight": weight,
						})
				case strings.HasPrefix(itemRGNbTrim, "preempt"):
					confRead.redundancyGroup[number]["preempt"] = true
					switch {
					case strings.HasPrefix(itemRGNbTrim, "preempt delay "):
						var err error
						confRead.redundancyGroup[number]["preempt_delay"], err = strconv.Atoi(strings.TrimPrefix(
							itemRGNbTrim, "preempt delay "))
						if err != nil {
							return confRead, fmt.Errorf(failedConvAtoiError, itemRGNbTrim, err)
						}
					case strings.HasPrefix(itemRGNbTrim, "preempt limit "):
						var err error
						confRead.redundancyGroup[number]["preempt_limit"], err = strconv.Atoi(strings.TrimPrefix(
							itemRGNbTrim, "preempt limit "))
						if err != nil {
							return confRead, fmt.Errorf(failedConvAtoiError, itemRGNbTrim, err)
						}
					case strings.HasPrefix(itemRGNbTrim, "preempt period "):
						var err error
						confRead.redundancyGroup[number]["preempt_period"], err = strconv.Atoi(strings.TrimPrefix(
							itemRGNbTrim, "preempt period "))
						if err != nil {
							return confRead, fmt.Errorf(failedConvAtoiError, itemRGNbTrim, err)
						}
					}
				}
			case strings.HasPrefix(itemTrim, "reth-count "):
				var err error
				confRead.rethCount, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "reth-count "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "configuration-synchronize no-secondary-bootup-auto":
				confRead.configSyncNoSecBootAuto = true
			case strings.HasPrefix(itemTrim, "control-ports fpc "):
				itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "control-ports fpc "), " ")
				if len(itemTrimSplit) < 3 {
					return confRead, fmt.Errorf("can't read values for control-ports fpc in '%s'", itemTrim)
				}
				controlPort := make(map[string]interface{})
				var err error
				controlPort["fpc"], err = strconv.Atoi(itemTrimSplit[0])
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrimSplit[0], err)
				}
				controlPort["port"], err = strconv.Atoi(itemTrimSplit[2])
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrimSplit[2], err)
				}
				confRead.controlPorts = append(confRead.controlPorts, controlPort)
			case itemTrim == "control-link-recovery":
				confRead.controlLinkRecovery = true
			case strings.HasPrefix(itemTrim, "heartbeat-interval "):
				var err error
				confRead.heartbeatInterval, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "heartbeat-interval "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "heartbeat-threshold "):
				var err error
				confRead.heartbeatThreshold, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "heartbeat-threshold "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}
	showConfigFab0, err := clt.command(cmdShowConfig+"interfaces fab0"+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfigFab0 != emptyW {
		if len(confRead.fab0) == 0 {
			confRead.fab0 = append(confRead.fab0, map[string]interface{}{
				"member_interfaces": make([]string, 0),
				"description":       "",
			})
		}
		for _, item := range strings.Split(showConfigFab0, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "description "):
				confRead.fab0[0]["description"] = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "fabric-options member-interfaces "):
				confRead.fab0[0]["member_interfaces"] = append(confRead.fab0[0]["member_interfaces"].([]string),
					strings.TrimPrefix(itemTrim, "fabric-options member-interfaces "))
			}
		}
	}
	showConfigFab1, err := clt.command(cmdShowConfig+"interfaces fab1"+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfigFab1 != emptyW {
		if len(confRead.fab1) == 0 {
			confRead.fab1 = append(confRead.fab1, map[string]interface{}{
				"member_interfaces": make([]string, 0),
				"description":       "",
			})
		}
		for _, item := range strings.Split(showConfigFab1, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "description "):
				confRead.fab1[0]["description"] = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "fabric-options member-interfaces "):
				confRead.fab1[0]["member_interfaces"] = append(confRead.fab1[0]["member_interfaces"].([]string),
					strings.TrimPrefix(itemTrim, "fabric-options member-interfaces "))
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
