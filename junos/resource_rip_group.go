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
				Default:          defaultW,
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
	sess := m.(*Session)
	routingInstance := d.Get("routing_instance").(string)
	name := d.Get("name").(string)
	ripNg := d.Get("ng").(bool)
	if sess.junosFakeCreateSetFile != "" {
		if err := setRipGroup(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		if ripNg {
			d.SetId(name + idSeparator + "ng" + idSeparator + routingInstance)
		} else {
			d.SetId(name + idSeparator + routingInstance)
		}

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if routingInstance != defaultW {
		instanceExists, err := checkRoutingInstanceExists(routingInstance, sess, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, sess.configClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", routingInstance))...)
		}
	}
	ripGroupExists, err := checkRipGroupExists(name, ripNg, routingInstance, sess, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ripGroupExists {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))
		protocolsRipGroup := "rip group"
		if ripNg {
			protocolsRipGroup = "ripng group"
		}
		if routingInstance != defaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				protocolsRipGroup+" %v already exists in routing-instance %v", name, routingInstance))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(protocolsRipGroup+" %v already exists", name))...)
	}
	if err := setRipGroup(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_rip_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ripGroupExists, err = checkRipGroupExists(name, ripNg, routingInstance, sess, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ripGroupExists {
		if ripNg {
			d.SetId(name + idSeparator + "ng" + idSeparator + routingInstance)
		} else {
			d.SetId(name + idSeparator + routingInstance)
		}
	} else {
		protocolsRipGroup := "rip group"
		if ripNg {
			protocolsRipGroup = "ripng group"
		}
		if routingInstance != defaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				protocolsRipGroup+" %v not exists in routing-instance %v after commit "+
					"=> check your config", name, routingInstance))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(protocolsRipGroup+" %v not exists after commit "+
			"=> check your config", name))...)
	}

	return append(diagWarns, resourceRipGroupReadWJunSess(d, sess, junSess)...)
}

func resourceRipGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)

	return resourceRipGroupReadWJunSess(d, sess, junSess)
}

func resourceRipGroupReadWJunSess(d *schema.ResourceData, sess *Session, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	ripGroupOptions, err := readRipGroup(
		d.Get("name").(string),
		d.Get("ng").(bool),
		d.Get("routing_instance").(string),
		sess, junSess)
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
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delRipGroup(
			d.Get("name").(string),
			d.Get("ng").(bool),
			d.Get("routing_instance").(string),
			false, sess, nil,
		); err != nil {
			return diag.FromErr(err)
		}
		if err := setRipGroup(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delRipGroup(
		d.Get("name").(string),
		d.Get("ng").(bool),
		d.Get("routing_instance").(string),
		false, sess, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setRipGroup(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_rip_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRipGroupReadWJunSess(d, sess, junSess)...)
}

func resourceRipGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delRipGroup(
			d.Get("name").(string),
			d.Get("ng").(bool),
			d.Get("routing_instance").(string),
			true, sess, nil,
		); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delRipGroup(
		d.Get("name").(string),
		d.Get("ng").(bool),
		d.Get("routing_instance").(string),
		true, sess, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_rip_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceRipGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	if len(idSplit) == 2 {
		ripGroupExists, err := checkRipGroupExists(idSplit[0], false, idSplit[1], sess, junSess)
		if err != nil {
			return nil, err
		}
		if !ripGroupExists {
			return nil, fmt.Errorf("don't find rip group id '%v' "+
				"(id must be <name>%s<routing_instance> or <name>%sng%s<routing_instance>",
				d.Id(), idSeparator, idSeparator, idSeparator,
			)
		}
		ripGroupOptions, err := readRipGroup(idSplit[0], false, idSplit[1], sess, junSess)
		if err != nil {
			return nil, err
		}
		fillRipGroupData(d, ripGroupOptions)

		result[0] = d

		return result, nil
	}
	if idSplit[1] != "ng" {
		return nil, fmt.Errorf("id must be <name>%s<routing_instance> or <name>%sng%s<routing_instance>",
			idSeparator, idSeparator, idSeparator,
		)
	}
	ripGroupExists, err := checkRipGroupExists(idSplit[0], true, idSplit[2], sess, junSess)
	if err != nil {
		return nil, err
	}
	if !ripGroupExists {
		return nil, fmt.Errorf("don't find ripng group with id '%v' "+
			"(id must be <name>%s<routing_instance> or <name>%sng%s<routing_instance>",
			d.Id(), idSeparator, idSeparator, idSeparator,
		)
	}
	ripGroupOptions, err := readRipGroup(idSplit[0], true, idSplit[2], sess, junSess)
	if err != nil {
		return nil, err
	}
	fillRipGroupData(d, ripGroupOptions)
	result[0] = d

	return result, nil
}

func checkRipGroupExists(name string, ripNg bool, routingInstance string, sess *Session, junSess *junosSession,
) (bool, error) {
	var showConfig string
	var err error
	protoRipGroup := "protocols rip group \"" + name + "\""
	if ripNg {
		protoRipGroup = "protocols ripng group \"" + name + "\""
	}
	if routingInstance == defaultW {
		showConfig, err = sess.command(cmdShowConfig+
			protoRipGroup+pipeDisplaySet, junSess)
	} else {
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			protoRipGroup+pipeDisplaySet, junSess)
	}
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setRipGroup(d *schema.ResourceData, sess *Session, junSess *junosSession) error {
	configSet := make([]string, 0, 1)

	setPrefix := setLS
	if rI := d.Get("routing_instance").(string); rI != defaultW {
		setPrefix = setRoutingInstances + rI + " "
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

	return sess.configSet(configSet, junSess)
}

func readRipGroup(name string, ripNg bool, routingInstance string, sess *Session, junSess *junosSession,
) (ripGroupOptions, error) {
	var confRead ripGroupOptions
	var showConfig string
	var err error
	confRead.preference = -1
	protoRipGroup := "protocols rip group \"" + name + "\""
	if ripNg {
		protoRipGroup = "protocols ripng group \"" + name + "\""
	}
	if routingInstance == defaultW {
		showConfig, err = sess.command(cmdShowConfig+
			protoRipGroup+pipeDisplaySetRelative, junSess)
	} else {
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			protoRipGroup+pipeDisplaySetRelative, junSess)
	}
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		confRead.ng = ripNg
		confRead.routingInstance = routingInstance
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "bfd-liveness-detection "):
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
				if err := readRipGroupBfd(
					strings.TrimPrefix(itemTrim, "bfd-liveness-detection "),
					confRead.bfdLivenessDetection[0],
				); err != nil {
					return confRead, err
				}
			case itemTrim == "demand-circuit":
				confRead.demandCircuit = true
			case strings.HasPrefix(itemTrim, "export "):
				confRead.exportPolicy = append(confRead.exportPolicy, strings.TrimPrefix(itemTrim, "export "))
			case strings.HasPrefix(itemTrim, "import "):
				confRead.importPolicy = append(confRead.importPolicy, strings.TrimPrefix(itemTrim, "import "))
			case strings.HasPrefix(itemTrim, "max-retrans-time "):
				confRead.maxRetransTime, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-retrans-time "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "metric-out "):
				confRead.metricOut, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "metric-out "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "preference "):
				confRead.preference, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "preference "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "route-timeout "):
				confRead.routeTimeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "route-timeout "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "update-interval "):
				confRead.updateInterval, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "update-interval "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func readRipGroupBfd(itemTrim string, bfd map[string]interface{}) error {
	switch {
	case strings.HasPrefix(itemTrim, "authentication algorithm "):
		bfd["authentication_algorithm"] = strings.TrimPrefix(itemTrim, "authentication algorithm ")
	case strings.HasPrefix(itemTrim, "authentication key-chain "):
		bfd["authentication_key_chain"] = strings.Trim(strings.TrimPrefix(itemTrim, "authentication key-chain "), "\"")
	case itemTrim == "authentication loose-check":
		bfd["authentication_loose_check"] = true
	case strings.HasPrefix(itemTrim, "detection-time threshold "):
		var err error
		bfd["detection_time_threshold"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "detection-time threshold "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "minimum-interval "):
		var err error
		bfd["minimum_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "minimum-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "minimum-receive-interval "):
		var err error
		bfd["minimum_receive_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "minimum-receive-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "multiplier "):
		var err error
		bfd["multiplier"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "multiplier "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "no-adaptation":
		bfd["no_adaptation"] = true
	case strings.HasPrefix(itemTrim, "transmit-interval minimum-interval "):
		var err error
		bfd["transmit_interval_minimum_interval"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "transmit-interval minimum-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "transmit-interval threshold "):
		var err error
		bfd["transmit_interval_threshold"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "transmit-interval threshold "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "version "):
		bfd["version"] = strings.TrimPrefix(itemTrim, "version ")
	}

	return nil
}

func delRipGroup(
	name string,
	ripNg bool,
	routingInstance string,
	deleteAll bool,
	sess *Session,
	junSess *junosSession,
) error {
	delPrefix := deleteLS
	if routingInstance != defaultW {
		delPrefix = delRoutingInstances + routingInstance + " "
	}
	if ripNg {
		delPrefix += "protocols ripng group "
	} else {
		delPrefix += "protocols rip group "
	}
	delPrefix += "\"" + name + "\" "

	if deleteAll {
		return sess.configSet([]string{delPrefix}, junSess)
	}
	configSet := make([]string, 0, 10)
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
	for _, line := range listLinesToDelete {
		configSet = append(configSet, delPrefix+line)
	}

	return sess.configSet(configSet, junSess)
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
