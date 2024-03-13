package providersdk

import (
	"context"
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
		CreateWithoutTimeout: resourceChassisRedundancyCreate,
		ReadWithoutTimeout:   resourceChassisRedundancyRead,
		UpdateWithoutTimeout: resourceChassisRedundancyUpdate,
		DeleteWithoutTimeout: resourceChassisRedundancyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceChassisRedundancyImport,
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setChassisRedundancy(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("redundancy")

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
	if err := setChassisRedundancy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_chassis_redundancy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("redundancy")

	return append(diagWarns, resourceChassisRedundancyReadWJunSess(d, junSess)...)
}

func resourceChassisRedundancyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceChassisRedundancyReadWJunSess(d, junSess)
}

func resourceChassisRedundancyReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	redundancyOptions, err := readChassisRedundancy(junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillChassisRedundancy(d, redundancyOptions)

	return nil
}

func resourceChassisRedundancyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delChassisRedundancy(junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setChassisRedundancy(d, junSess); err != nil {
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
	if err := delChassisRedundancy(junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setChassisRedundancy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_chassis_redundancy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceChassisRedundancyReadWJunSess(d, junSess)...)
}

func resourceChassisRedundancyDelete(ctx context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delChassisRedundancy(junSess); err != nil {
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
	if err := delChassisRedundancy(junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_chassis_redundancy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceChassisRedundancyImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	redundancyOptions, err := readChassisRedundancy(junSess)
	if err != nil {
		return nil, err
	}
	fillChassisRedundancy(d, redundancyOptions)
	d.SetId("redundancy")
	result[0] = d

	return result, nil
}

func setChassisRedundancy(d *schema.ResourceData, junSess *junos.Session) error {
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
		if slices.Contains(routingEngineList, routingEngine["slot"].(int)) {
			return fmt.Errorf("multiple blocks routing_engine with the same slot '%d'", routingEngine["slot"].(int))
		}
		routingEngineList = append(routingEngineList, routingEngine["slot"].(int))
		configSet = append(configSet, setPrefix+
			"routing-engine "+strconv.Itoa(routingEngine["slot"].(int))+
			" "+routingEngine["role"].(string))
	}

	return junSess.ConfigSet(configSet)
}

func delChassisRedundancy(junSess *junos.Session) error {
	configSet := []string{"delete chassis redundancy"}

	return junSess.ConfigSet(configSet)
}

func readChassisRedundancy(junSess *junos.Session,
) (confRead chassisRedundancyOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "chassis redundancy" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "failover disk-read-threshold "):
				confRead.failoverDiskReadThreshold, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "failover disk-write-threshold "):
				confRead.failoverDiskWriteThreshold, err = strconv.Atoi(itemTrim)
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
			case balt.CutPrefixInString(&itemTrim, "keepalive-time "):
				confRead.keepaliveTime, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "routing-engine "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 2 { // <slot> <role>
					return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "routing-engine", itemTrim)
				}
				slot, err := strconv.Atoi(itemTrimFields[0])
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrimFields[0], err)
				}
				confRead.routingEngine = append(confRead.routingEngine, map[string]interface{}{
					"slot": slot,
					"role": itemTrimFields[1],
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
