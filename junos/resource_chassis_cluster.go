package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type chassisClusterOptions struct {
	configSyncNoSecBootAuto bool
	controlLinkRecovery     bool
	heartbeatInterval       int
	heartbeatThreshold      int
	rethCount               int
	redundancyGroup         []map[string]interface{}
	fab0                    []map[string]interface{}
	fab1                    []map[string]interface{}
}

func resourceChassisCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChassisClusterCreate,
		ReadContext:   resourceChassisClusterRead,
		UpdateContext: resourceChassisClusterUpdate,
		DeleteContext: resourceChassisClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceChassisClusterImport,
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
							Elem:     &schema.Schema{Type: schema.TypeString},
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
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func checkCompatibilityChassisCluster(jnprSess *NetconfObject) bool {
	if strings.HasPrefix(strings.ToLower(jnprSess.SystemInformation.HardwareModel), "srx") {
		return true
	}
	if strings.HasPrefix(strings.ToLower(jnprSess.SystemInformation.HardwareModel), "vsrx") {
		return true
	}
	if strings.HasPrefix(strings.ToLower(jnprSess.SystemInformation.HardwareModel), "j") {
		return true
	}

	return false
}

func resourceChassisClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setChassisCluster(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("cluster")

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilityChassisCluster(jnprSess) {
		return diag.FromErr(fmt.Errorf("chassis cluster "+
			"not compatible with Junos device %s", jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := setChassisCluster(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_chassis_cluster", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("cluster")

	return append(diagWarns, resourceChassisClusterReadWJnprSess(d, m, jnprSess)...)
}

func resourceChassisClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceChassisClusterReadWJnprSess(d, m, jnprSess)
}

func resourceChassisClusterReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	clusterOptions, err := readChassisCluster(m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillChassisCluster(d, clusterOptions)

	return nil
}

func resourceChassisClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delChassisCluster(m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setChassisCluster(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_chassis_cluster", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceChassisClusterReadWJnprSess(d, m, jnprSess)...)
}

func resourceChassisClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delChassisCluster(m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_chassis_cluster", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceChassisClusterImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	clusterOptions, err := readChassisCluster(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillChassisCluster(d, clusterOptions)
	d.SetId("cluster")
	result[0] = d

	return result, nil
}

func setChassisCluster(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
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
		for _, v2 := range redundancyGroup["interface_monitor"].([]interface{}) {
			interfaceMonitor := v2.(map[string]interface{})
			configSet = append(configSet, setChassisluster+"redundancy-group "+strconv.Itoa(i)+
				" interface-monitor "+interfaceMonitor["name"].(string)+
				" weight "+strconv.Itoa(interfaceMonitor["weight"].(int)))
		}
		if redundancyGroup["preempt"].(bool) {
			configSet = append(configSet, setChassisluster+"redundancy-group "+strconv.Itoa(i)+
				" preempt")
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
	if v := d.Get("heartbeat_interval").(int); v != 0 {
		configSet = append(configSet, setChassisluster+"heartbeat-interval "+
			strconv.Itoa(v))
	}
	if v := d.Get("heartbeat_threshold").(int); v != 0 {
		configSet = append(configSet, setChassisluster+"heartbeat-threshold "+
			strconv.Itoa(v))
	}

	return sess.configSet(configSet, jnprSess)
}

func delChassisCluster(m interface{}, jnprSess *NetconfObject) error {
	listLinesToDelete := make([]string, 0)
	listLinesToDelete = append(listLinesToDelete, "chassis cluster")
	listLinesToDelete = append(listLinesToDelete, "interfaces fab0")
	listLinesToDelete = append(listLinesToDelete, "interfaces fab1")
	sess := m.(*Session)
	configSet := make([]string, 0)
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			"delete "+line)
	}

	return sess.configSet(configSet, jnprSess)
}

func readChassisCluster(m interface{}, jnprSess *NetconfObject) (chassisClusterOptions, error) {
	sess := m.(*Session)
	var confRead chassisClusterOptions

	chassisClusterConfig, err := sess.command("show configuration chassis cluster"+
		" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if chassisClusterConfig != emptyWord {
		for _, item := range strings.Split(chassisClusterConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "redundancy-group "):
				itemRGTrim := strings.TrimPrefix(itemTrim, "redundancy-group ")
				itemRGTrimSplit := strings.Split(itemRGTrim, " ")
				number, err := strconv.Atoi(itemRGTrimSplit[0])
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemRGTrimSplit[0], err)
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
						})
					}
				}
				itemRGNbTrim := strings.TrimPrefix(itemRGTrim, strconv.Itoa(number)+" ")
				switch {
				case strings.HasPrefix(itemRGNbTrim, "node 0 priority "):
					var err error
					confRead.redundancyGroup[number]["node0_priority"], err =
						strconv.Atoi(strings.TrimPrefix(itemRGNbTrim, "node 0 priority "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemRGNbTrim, err)
					}
				case strings.HasPrefix(itemRGNbTrim, "node 1 priority "):
					var err error
					confRead.redundancyGroup[number]["node1_priority"], err =
						strconv.Atoi(strings.TrimPrefix(itemRGNbTrim, "node 1 priority "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemRGNbTrim, err)
					}
				case strings.HasPrefix(itemRGNbTrim, "gratuitous-arp-count "):
					var err error
					confRead.redundancyGroup[number]["gratuitous_arp_count"], err =
						strconv.Atoi(strings.TrimPrefix(itemRGNbTrim, "gratuitous-arp-count "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemRGNbTrim, err)
					}
				case strings.HasPrefix(itemRGNbTrim, "hold-down-interval "):
					var err error
					confRead.redundancyGroup[number]["hold_down_interval"], err =
						strconv.Atoi(strings.TrimPrefix(itemRGNbTrim, "hold-down-interval "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemRGNbTrim, err)
					}
				case strings.HasPrefix(itemRGNbTrim, "interface-monitor "):
					name := strings.Split(strings.TrimPrefix(itemRGNbTrim, "interface-monitor "), " ")[0]
					weight, err := strconv.Atoi(strings.TrimPrefix(itemRGNbTrim, "interface-monitor "+name+" weight "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemRGNbTrim, err)
					}
					confRead.redundancyGroup[number]["interface_monitor"] = append(
						confRead.redundancyGroup[number]["interface_monitor"].([]map[string]interface{}),
						map[string]interface{}{
							"name":   name,
							"weight": weight,
						})
				case itemRGNbTrim == preemptWord:
					confRead.redundancyGroup[number][preemptWord] = true
				}
			case strings.HasPrefix(itemTrim, "reth-count "):
				var err error
				confRead.rethCount, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "reth-count "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case itemTrim == "configuration-synchronize no-secondary-bootup-auto":
				confRead.configSyncNoSecBootAuto = true
			case itemTrim == "control-link-recovery":
				confRead.controlLinkRecovery = true
			case strings.HasPrefix(itemTrim, "heartbeat-interval "):
				var err error
				confRead.heartbeatInterval, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "heartbeat-interval "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "heartbeat-threshold "):
				var err error
				confRead.heartbeatThreshold, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "heartbeat-threshold "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			}
		}
	}
	fab0Config, err := sess.command("show configuration interfaces fab0"+
		" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if fab0Config != emptyWord {
		if len(confRead.fab0) == 0 {
			confRead.fab0 = append(confRead.fab0, map[string]interface{}{
				"member_interfaces": make([]string, 0),
				"description":       "",
			})
		}
		for _, item := range strings.Split(fab0Config, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "description "):
				confRead.fab0[0]["description"] = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "fabric-options member-interfaces "):
				confRead.fab0[0]["member_interfaces"] = append(confRead.fab0[0]["member_interfaces"].([]string),
					strings.TrimPrefix(itemTrim, "fabric-options member-interfaces "))
			}
		}
	}
	fab1Config, err := sess.command("show configuration interfaces fab1"+
		" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if fab1Config != emptyWord {
		if len(confRead.fab1) == 0 {
			confRead.fab1 = append(confRead.fab1, map[string]interface{}{
				"member_interfaces": make([]string, 0),
				"description":       "",
			})
		}
		for _, item := range strings.Split(fab1Config, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
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
