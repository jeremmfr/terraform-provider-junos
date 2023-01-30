package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type ripGroupOptions struct {
	demandCircuit        bool
	ng                   bool
	maxRetransTime       int
	metricOut            int
	preference           int
	routeTimeout         int
	updateInterval       int
	name                 string
	routingInstance      string
	exportPolicy         []string
	importPolicy         []string
	bfdLivenessDetection []map[string]interface{}
}

func resourceRipGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceRipGroupCreate,
		ReadWithoutTimeout:   resourceRipGroupRead,
		UpdateWithoutTimeout: resourceRipGroupUpdate,
		DeleteWithoutTimeout: resourceRipGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRipGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 48),
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          junos.DefaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"ng": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"bfd_liveness_detection": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"ng"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication_algorithm": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"authentication_key_chain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"authentication_loose_check": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"detection_time_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 4294967295),
						},
						"minimum_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 255000),
						},
						"minimum_receive_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 255000),
						},
						"multiplier": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 255),
						},
						"no_adaptation": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"transmit_interval_minimum_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 255000),
						},
						"transmit_interval_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 4294967295),
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"demand_circuit": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"ng"},
			},
			"export": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"import": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"max_retrans_time": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"ng"},
				ValidateFunc:  validation.IntBetween(5, 180),
			},
			"metric_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 15),
			},
			"preference": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(0, 4294967295),
			},
			"route_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(30, 360),
			},
			"update_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(10, 60),
			},
		},
	}
}

func resourceRipGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	routingInstance := d.Get("routing_instance").(string)
	name := d.Get("name").(string)
	ripNg := d.Get("ng").(bool)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setRipGroup(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if ripNg {
			d.SetId(name + junos.IDSeparator + "ng" + junos.IDSeparator + routingInstance)
		} else {
			d.SetId(name + junos.IDSeparator + routingInstance)
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
	if routingInstance != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(routingInstance, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", routingInstance))...)
		}
	}
	ripGroupExists, err := checkRipGroupExists(name, ripNg, routingInstance, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ripGroupExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())
		protocolsRipGroup := "rip group"
		if ripNg {
			protocolsRipGroup = "ripng group"
		}
		if routingInstance != junos.DefaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				protocolsRipGroup+" %v already exists in routing-instance %v", name, routingInstance))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(protocolsRipGroup+" %v already exists", name))...)
	}
	if err := setRipGroup(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_rip_group")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	ripGroupExists, err = checkRipGroupExists(name, ripNg, routingInstance, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ripGroupExists {
		if ripNg {
			d.SetId(name + junos.IDSeparator + "ng" + junos.IDSeparator + routingInstance)
		} else {
			d.SetId(name + junos.IDSeparator + routingInstance)
		}
	} else {
		protocolsRipGroup := "rip group"
		if ripNg {
			protocolsRipGroup = "ripng group"
		}
		if routingInstance != junos.DefaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				protocolsRipGroup+" %v not exists in routing-instance %v after commit "+
					"=> check your config", name, routingInstance))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(protocolsRipGroup+" %v not exists after commit "+
			"=> check your config", name))...)
	}

	return append(diagWarns, resourceRipGroupReadWJunSess(d, junSess)...)
}

func resourceRipGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceRipGroupReadWJunSess(d, junSess)
}

func resourceRipGroupReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	ripGroupOptions, err := readRipGroup(
		d.Get("name").(string),
		d.Get("ng").(bool),
		d.Get("routing_instance").(string),
		junSess,
	)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ripGroupOptions.name == "" {
		d.SetId("")
	} else {
		fillRipGroupData(d, ripGroupOptions)
	}

	return nil
}

func resourceRipGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delRipGroup(
			d.Get("name").(string),
			d.Get("ng").(bool),
			d.Get("routing_instance").(string),
			false,
			junSess,
		); err != nil {
			return diag.FromErr(err)
		}
		if err := setRipGroup(d, junSess); err != nil {
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
	if err := delRipGroup(
		d.Get("name").(string),
		d.Get("ng").(bool),
		d.Get("routing_instance").(string),
		false,
		junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setRipGroup(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_rip_group")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRipGroupReadWJunSess(d, junSess)...)
}

func resourceRipGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delRipGroup(
			d.Get("name").(string),
			d.Get("ng").(bool),
			d.Get("routing_instance").(string),
			true,
			junSess,
		); err != nil {
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
	if err := delRipGroup(
		d.Get("name").(string),
		d.Get("ng").(bool),
		d.Get("routing_instance").(string),
		true,
		junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_rip_group")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceRipGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), junos.IDSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", junos.IDSeparator)
	}
	if len(idSplit) == 2 {
		ripGroupExists, err := checkRipGroupExists(idSplit[0], false, idSplit[1], junSess)
		if err != nil {
			return nil, err
		}
		if !ripGroupExists {
			return nil, fmt.Errorf("don't find rip group id '%v' "+
				"(id must be <name>"+junos.IDSeparator+"<routing_instance> or "+
				"<name>"+junos.IDSeparator+"ng"+junos.IDSeparator+"<routing_instance>",
				d.Id(),
			)
		}
		ripGroupOptions, err := readRipGroup(idSplit[0], false, idSplit[1], junSess)
		if err != nil {
			return nil, err
		}
		fillRipGroupData(d, ripGroupOptions)

		result[0] = d

		return result, nil
	}
	if idSplit[1] != "ng" {
		return nil, fmt.Errorf("id must be <name>" + junos.IDSeparator + "<routing_instance> or " +
			"<name>" + junos.IDSeparator + "ng" + junos.IDSeparator + "<routing_instance>",
		)
	}
	ripGroupExists, err := checkRipGroupExists(idSplit[0], true, idSplit[2], junSess)
	if err != nil {
		return nil, err
	}
	if !ripGroupExists {
		return nil, fmt.Errorf("don't find ripng group with id '%v' "+
			"(id must be <name>"+junos.IDSeparator+"<routing_instance> or "+
			"<name>"+junos.IDSeparator+"ng"+junos.IDSeparator+"<routing_instance>",
			d.Id(),
		)
	}
	ripGroupOptions, err := readRipGroup(idSplit[0], true, idSplit[2], junSess)
	if err != nil {
		return nil, err
	}
	fillRipGroupData(d, ripGroupOptions)
	result[0] = d

	return result, nil
}

func checkRipGroupExists(name string, ripNg bool, routingInstance string, junSess *junos.Session,
) (_ bool, err error) {
	var showConfig string
	protoRipGroup := "protocols rip group \"" + name + "\""
	if ripNg {
		protoRipGroup = "protocols ripng group \"" + name + "\""
	}
	if routingInstance == junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			protoRipGroup + junos.PipeDisplaySet)
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			protoRipGroup + junos.PipeDisplaySet)
	}
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setRipGroup(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)

	setPrefix := junos.SetLS
	if rI := d.Get("routing_instance").(string); rI != junos.DefaultW {
		setPrefix = junos.SetRoutingInstances + rI + " "
	}
	if d.Get("ng").(bool) {
		setPrefix += "protocols ripng group "
	} else {
		setPrefix += "protocols rip group "
	}
	setPrefix += "\"" + d.Get("name").(string) + "\" "

	configSet = append(configSet, setPrefix)
	for _, mBFDLivDet := range d.Get("bfd_liveness_detection").([]interface{}) {
		if mBFDLivDet == nil {
			return fmt.Errorf("bfd_liveness_detection block is empty")
		}
		setPrefixBfd := setPrefix + "bfd-liveness-detection "
		bfdLiveDetect := mBFDLivDet.(map[string]interface{})
		if v := bfdLiveDetect["authentication_algorithm"].(string); v != "" {
			configSet = append(configSet, setPrefixBfd+"authentication algorithm "+v)
		}
		if v := bfdLiveDetect["authentication_key_chain"].(string); v != "" {
			configSet = append(configSet, setPrefixBfd+"authentication key-chain \""+v+"\"")
		}
		if bfdLiveDetect["authentication_loose_check"].(bool) {
			configSet = append(configSet, setPrefixBfd+"authentication loose-check")
		}
		if v := bfdLiveDetect["detection_time_threshold"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+"detection-time threshold "+strconv.Itoa(v))
		}
		if v := bfdLiveDetect["minimum_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+"minimum-interval "+strconv.Itoa(v))
		}
		if v := bfdLiveDetect["minimum_receive_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+"minimum-receive-interval "+strconv.Itoa(v))
		}
		if v := bfdLiveDetect["multiplier"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+"multiplier "+strconv.Itoa(v))
		}
		if bfdLiveDetect["no_adaptation"].(bool) {
			configSet = append(configSet, setPrefixBfd+"no-adaptation")
		}
		if v := bfdLiveDetect["transmit_interval_minimum_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+
				"transmit-interval minimum-interval "+strconv.Itoa(v))
		}
		if v := bfdLiveDetect["transmit_interval_threshold"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+
				"transmit-interval threshold "+strconv.Itoa(v))
		}
		if v := bfdLiveDetect["version"].(string); v != "" {
			configSet = append(configSet, setPrefixBfd+"version "+v)
		}
		if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefixBfd) {
			return fmt.Errorf("bfd_liveness_detection block is empty")
		}
	}
	if d.Get("demand_circuit").(bool) {
		configSet = append(configSet, setPrefix+"demand-circuit")
	}
	for _, exportPolicy := range d.Get("export").([]interface{}) {
		configSet = append(configSet, setPrefix+"export "+exportPolicy.(string))
	}
	for _, importPolicy := range d.Get("import").([]interface{}) {
		configSet = append(configSet, setPrefix+"import "+importPolicy.(string))
	}
	if v := d.Get("max_retrans_time").(int); v != 0 {
		configSet = append(configSet, setPrefix+"max-retrans-time "+strconv.Itoa(v))
	}
	if v := d.Get("metric_out").(int); v != 0 {
		configSet = append(configSet, setPrefix+"metric-out "+strconv.Itoa(v))
	}
	if v := d.Get("preference").(int); v != -1 {
		configSet = append(configSet, setPrefix+"preference "+strconv.Itoa(v))
	}
	if v := d.Get("route_timeout").(int); v != 0 {
		configSet = append(configSet, setPrefix+"route-timeout "+strconv.Itoa(v))
	}
	if v := d.Get("update_interval").(int); v != 0 {
		configSet = append(configSet, setPrefix+"update-interval "+strconv.Itoa(v))
	}

	return junSess.ConfigSet(configSet)
}

func readRipGroup(name string, ripNg bool, routingInstance string, junSess *junos.Session,
) (confRead ripGroupOptions, err error) {
	// default -1
	confRead.preference = -1
	var showConfig string
	protoRipGroup := "protocols rip group \"" + name + "\""
	if ripNg {
		protoRipGroup = "protocols ripng group \"" + name + "\""
	}
	if routingInstance == junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			protoRipGroup + junos.PipeDisplaySetRelative)
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			protoRipGroup + junos.PipeDisplaySetRelative)
	}
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		confRead.ng = ripNg
		confRead.routingInstance = routingInstance
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "bfd-liveness-detection "):
				if len(confRead.bfdLivenessDetection) == 0 {
					confRead.bfdLivenessDetection = append(confRead.bfdLivenessDetection, map[string]interface{}{
						"authentication_algorithm":           "",
						"authentication_key_chain":           "",
						"authentication_loose_check":         false,
						"detection_time_threshold":           0,
						"minimum_interval":                   0,
						"minimum_receive_interval":           0,
						"multiplier":                         0,
						"no_adaptation":                      false,
						"transmit_interval_minimum_interval": 0,
						"transmit_interval_threshold":        0,
						"version":                            "",
					})
				}
				if err := readRipGroupBfd(itemTrim, confRead.bfdLivenessDetection[0]); err != nil {
					return confRead, err
				}
			case itemTrim == "demand-circuit":
				confRead.demandCircuit = true
			case balt.CutPrefixInString(&itemTrim, "export "):
				confRead.exportPolicy = append(confRead.exportPolicy, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "import "):
				confRead.importPolicy = append(confRead.importPolicy, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "max-retrans-time "):
				confRead.maxRetransTime, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "metric-out "):
				confRead.metricOut, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "preference "):
				confRead.preference, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "route-timeout "):
				confRead.routeTimeout, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "update-interval "):
				confRead.updateInterval, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func readRipGroupBfd(itemTrim string, bfd map[string]interface{}) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "authentication algorithm "):
		bfd["authentication_algorithm"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "authentication key-chain "):
		bfd["authentication_key_chain"] = strings.Trim(itemTrim, "\"")
	case itemTrim == "authentication loose-check":
		bfd["authentication_loose_check"] = true
	case balt.CutPrefixInString(&itemTrim, "detection-time threshold "):
		bfd["detection_time_threshold"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "minimum-interval "):
		bfd["minimum_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "minimum-receive-interval "):
		bfd["minimum_receive_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "multiplier "):
		bfd["multiplier"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "no-adaptation":
		bfd["no_adaptation"] = true
	case balt.CutPrefixInString(&itemTrim, "transmit-interval minimum-interval "):
		bfd["transmit_interval_minimum_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "transmit-interval threshold "):
		bfd["transmit_interval_threshold"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "version "):
		bfd["version"] = itemTrim
	}

	return nil
}

func delRipGroup(
	name string,
	ripNg bool,
	routingInstance string,
	deleteAll bool,
	junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if routingInstance != junos.DefaultW {
		delPrefix = junos.DelRoutingInstances + routingInstance + " "
	}
	if ripNg {
		delPrefix += "protocols ripng group "
	} else {
		delPrefix += "protocols rip group "
	}
	delPrefix += "\"" + name + "\" "

	if deleteAll {
		return junSess.ConfigSet([]string{delPrefix})
	}
	listLinesToDelete := []string{
		"bfd-liveness-detection",
		"demand-circuit",
		"export",
		"import",
		"max-retrans-time",
		"metric-out",
		"preference",
		"route-timeout",
		"update-interval",
	}
	configSet := make([]string, len(listLinesToDelete))
	for k, line := range listLinesToDelete {
		configSet[k] = delPrefix + line
	}

	return junSess.ConfigSet(configSet)
}

func fillRipGroupData(d *schema.ResourceData, ripGroupOptions ripGroupOptions) {
	if tfErr := d.Set("name", ripGroupOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", ripGroupOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ng", ripGroupOptions.ng); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bfd_liveness_detection", ripGroupOptions.bfdLivenessDetection); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("demand_circuit", ripGroupOptions.demandCircuit); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("export", ripGroupOptions.exportPolicy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("import", ripGroupOptions.importPolicy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("max_retrans_time", ripGroupOptions.maxRetransTime); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out", ripGroupOptions.metricOut); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("preference", ripGroupOptions.preference); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("route_timeout", ripGroupOptions.routeTimeout); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("update_interval", ripGroupOptions.updateInterval); tfErr != nil {
		panic(tfErr)
	}
}
