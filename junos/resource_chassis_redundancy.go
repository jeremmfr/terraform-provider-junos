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

type chassisRedundancyOptions struct {
	failoverNotOnDiskUnderperform bool
	failoverOnDiskFailure         bool
	failoverOnLossOfKeepalives    bool
	gracefulSwitchover            bool
	failoverDiskReadThreshold     int
	failoverDiskWriteThreshold    int
	keepaliveTime                 int
	routingEngine                 []map[string]interface{}
}

func resourceChassisRedundancy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChassisRedundancyCreate,
		ReadContext:   resourceChassisRedundancyRead,
		UpdateContext: resourceChassisRedundancyUpdate,
		DeleteContext: resourceChassisRedundancyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceChassisRedundancyImport,
		},
		Schema: map[string]*schema.Schema{
			"failover_disk_read_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1000, 10000),
			},
			"failover_disk_write_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1000, 10000),
			},
			"failover_not_on_disk_underperform": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"failover_on_disk_failure": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"failover_on_loss_of_keepalives": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"graceful_switchover": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"keepalive_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 10000),
			},
			"routing_engine": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slot": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1),
						},
						"role": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"backup", "disabled", "master"}, false),
						},
					},
				},
			},
		},
	}
}

func resourceChassisRedundancyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setChassisRedundancy(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("redundancy")

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := setChassisRedundancy(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_chassis_redundancy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("redundancy")

	return append(diagWarns, resourceChassisRedundancyReadWJnprSess(d, m, jnprSess)...)
}

func resourceChassisRedundancyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceChassisRedundancyReadWJnprSess(d, m, jnprSess)
}

func resourceChassisRedundancyReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	redundancyOptions, err := readChassisRedundancy(m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillChassisRedundancy(d, redundancyOptions)

	return nil
}

func resourceChassisRedundancyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delChassisRedundancy(m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setChassisRedundancy(d, m, nil); err != nil {
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
	if err := delChassisRedundancy(m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setChassisRedundancy(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_chassis_redundancy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceChassisRedundancyReadWJnprSess(d, m, jnprSess)...)
}

func resourceChassisRedundancyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delChassisRedundancy(m, nil); err != nil {
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
	if err := delChassisRedundancy(m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_chassis_redundancy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceChassisRedundancyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	redundancyOptions, err := readChassisRedundancy(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillChassisRedundancy(d, redundancyOptions)
	d.SetId("redundancy")
	result[0] = d

	return result, nil
}

func setChassisRedundancy(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set chassis redundancy "

	if v := d.Get("failover_disk_read_threshold").(int); v != 0 {
		configSet = append(configSet, setPrefix+"failover disk-read-threshold "+strconv.Itoa(v))
	}
	if v := d.Get("failover_disk_write_threshold").(int); v != 0 {
		configSet = append(configSet, setPrefix+"failover disk-write-threshold "+strconv.Itoa(v))
	}
	if d.Get("failover_not_on_disk_underperform").(bool) {
		configSet = append(configSet, setPrefix+"failover not-on-disk-underperform")
	}
	if d.Get("failover_on_disk_failure").(bool) {
		configSet = append(configSet, setPrefix+"failover on-disk-failure")
	}
	if d.Get("failover_on_loss_of_keepalives").(bool) {
		configSet = append(configSet, setPrefix+"failover on-loss-of-keepalives")
	}
	if d.Get("graceful_switchover").(bool) {
		configSet = append(configSet, setPrefix+"graceful-switchover")
	}
	if v := d.Get("keepalive_time").(int); v != 0 {
		configSet = append(configSet, setPrefix+"keepalive-time "+strconv.Itoa(v))
	}
	routingEngineList := make([]int, 0)
	for _, mRE := range d.Get("routing_engine").(*schema.Set).List() {
		routingEngine := mRE.(map[string]interface{})
		if bchk.IntInSlice(routingEngine["slot"].(int), routingEngineList) {
			return fmt.Errorf("multiple blocks routing_engine with the same slot '%d'", routingEngine["slot"].(int))
		}
		routingEngineList = append(routingEngineList, routingEngine["slot"].(int))
		configSet = append(configSet, setPrefix+
			"routing-engine "+strconv.Itoa(routingEngine["slot"].(int))+
			" "+routingEngine["role"].(string))
	}

	return sess.configSet(configSet, jnprSess)
}

func delChassisRedundancy(m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete chassis redundancy"}

	return sess.configSet(configSet, jnprSess)
}

func readChassisRedundancy(m interface{}, jnprSess *NetconfObject) (chassisRedundancyOptions, error) {
	sess := m.(*Session)
	var confRead chassisRedundancyOptions

	showConfig, err := sess.command(cmdShowConfig+"chassis redundancy"+pipeDisplaySetRelative, jnprSess)
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
			case strings.HasPrefix(itemTrim, "failover disk-read-threshold "):
				var err error
				confRead.failoverDiskReadThreshold, err = strconv.Atoi(strings.TrimPrefix(
					itemTrim, "failover disk-read-threshold "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "failover disk-write-threshold "):
				var err error
				confRead.failoverDiskWriteThreshold, err = strconv.Atoi(strings.TrimPrefix(
					itemTrim, "failover disk-write-threshold "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "failover not-on-disk-underperform":
				confRead.failoverNotOnDiskUnderperform = true
			case itemTrim == "failover on-disk-failure":
				confRead.failoverOnDiskFailure = true
			case itemTrim == "failover on-loss-of-keepalives":
				confRead.failoverOnLossOfKeepalives = true
			case itemTrim == "graceful-switchover":
				confRead.gracefulSwitchover = true
			case strings.HasPrefix(itemTrim, "keepalive-time "):
				var err error
				confRead.keepaliveTime, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "keepalive-time "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "routing-engine "):
				itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "routing-engine "), " ")
				if len(itemTrimSplit) < 2 {
					return confRead, fmt.Errorf("can't find slot and role in %s", itemTrim)
				}
				slot, err := strconv.Atoi(itemTrimSplit[0])
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrimSplit[0], err)
				}
				confRead.routingEngine = append(confRead.routingEngine, map[string]interface{}{
					"slot": slot,
					"role": itemTrimSplit[1],
				})
			}
		}
	}

	return confRead, nil
}

func fillChassisRedundancy(d *schema.ResourceData, chassisRedundancyOptions chassisRedundancyOptions) {
	if tfErr := d.Set("failover_disk_read_threshold",
		chassisRedundancyOptions.failoverDiskReadThreshold); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("failover_disk_write_threshold",
		chassisRedundancyOptions.failoverDiskWriteThreshold); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("failover_not_on_disk_underperform",
		chassisRedundancyOptions.failoverNotOnDiskUnderperform); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("failover_on_disk_failure",
		chassisRedundancyOptions.failoverOnDiskFailure); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("failover_on_loss_of_keepalives",
		chassisRedundancyOptions.failoverOnLossOfKeepalives); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("graceful_switchover",
		chassisRedundancyOptions.gracefulSwitchover); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("keepalive_time",
		chassisRedundancyOptions.keepaliveTime); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_engine",
		chassisRedundancyOptions.routingEngine); tfErr != nil {
		panic(tfErr)
	}
}
