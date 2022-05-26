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
	jdecode "github.com/jeremmfr/junosdecode"
)

type ripNeighborOptions struct {
	anySender                  bool
	checkZero                  bool
	demandCircuit              bool
	dynamicPeers               bool
	interfaceTypeP2mp          bool
	ng                         bool
	noCheckZero                bool
	maxRetransTime             int
	messageSize                int
	metricIn                   int
	routeTimeout               int
	updateInterval             int
	authenticationKey          string
	authenticationType         string
	name                       string
	group                      string
	receive                    string
	routingInstance            string
	send                       string
	importPolicy               []string
	peer                       []string
	authenticationSelectiveMD5 []map[string]interface{}
	bfdLivenessDetection       []map[string]interface{}
}

func resourceRipNeighbor() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceRipNeighborCreate,
		ReadWithoutTimeout:   resourceRipNeighborRead,
		UpdateWithoutTimeout: resourceRipNeighborUpdate,
		DeleteWithoutTimeout: resourceRipNeighborDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRipNeighborImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if strings.Count(value, ".") != 1 {
						errors = append(errors, fmt.Errorf(
							"%q in %q need to have 1 dot", value, k))
					}

					return
				},
			},
			"group": {
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
			"any_sender": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"ng"},
			},
			"authentication_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"authentication_selective_md5", "ng"},
				Sensitive:     true,
			},
			"authentication_selective_md5": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"authentication_key", "authentication_type", "ng"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 255),
						},
						"key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^\d{4}\-\d\d?\-\d\d?\.\d{2}:\d{2}:\d{2}$`),
								"must be in the format 'YYYY-MM-DD.HH:MM:SS'"),
						},
					},
				},
			},
			"authentication_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"authentication_selective_md5", "ng"},
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
			"check_zero": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"no_check_zero", "ng"},
			},
			"demand_circuit": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"ng"},
			},
			"dynamic_peers": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"ng"},
				RequiredWith:  []string{"interface_type_p2mp"},
			},
			"import": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"interface_type_p2mp": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"ng"},
			},
			"max_retrans_time": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"ng"},
				ValidateFunc:  validation.IntBetween(5, 180),
			},
			"message_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"ng"},
				ValidateFunc:  validation.IntBetween(25, 255),
			},
			"metric_in": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 15),
			},
			"no_check_zero": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"check_zero", "ng"},
			},
			"peer": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"ng"},
				RequiredWith:  []string{"interface_type_p2mp"},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPAddress,
				},
			},
			"receive": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"both", "none", "version-1", "version-2"}, false),
			},
			"route_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(30, 360),
			},
			"send": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"broadcast", "multicast", "none", "version-1"}, false),
			},
			"update_interval": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"ng"},
				ValidateFunc:  validation.IntBetween(10, 60),
			},
		},
	}
}

func resourceRipNeighborCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	name := d.Get("name").(string)
	group := d.Get("group").(string)
	ripNg := d.Get("ng").(bool)
	routingInstance := d.Get("routing_instance").(string)
	if sess.junosFakeCreateSetFile != "" {
		if err := setRipNeighbor(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		if ripNg {
			d.SetId(name + idSeparator + group + idSeparator + "ng" + idSeparator + routingInstance)
		} else {
			d.SetId(name + idSeparator + group + idSeparator + routingInstance)
		}

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if routingInstance != defaultW {
		instanceExists, err := checkRoutingInstanceExists(routingInstance, sess, jnprSess)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", routingInstance))...)
		}
	}
	ripGroupExists, err := checkRipGroupExists(group, ripNg, routingInstance, sess, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !ripGroupExists {
		if ripNg {
			return append(diagWarns, diag.FromErr(fmt.Errorf("ripng group %v doesn't exist", group))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf("rip group %v doesn't exist", group))...)
	}
	ripNeighborExists, err := checkRipNeighborExists(name, group, ripNg, routingInstance, sess, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ripNeighborExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))
		protocolsRipNeighbor := "rip group " + group + " neighbor"
		if ripNg {
			protocolsRipNeighbor = "ripng group " + group + " neighbor"
		}
		if routingInstance != defaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				protocolsRipNeighbor+" %v already exists in routing-instance %v", name, routingInstance))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(protocolsRipNeighbor+" %v already exists", name))...)
	}
	if err := setRipNeighbor(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_rip_neighbor", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ripNeighborExists, err = checkRipNeighborExists(name, group, ripNg, routingInstance, sess, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ripNeighborExists {
		if ripNg {
			d.SetId(name + idSeparator + group + idSeparator + "ng" + idSeparator + routingInstance)
		} else {
			d.SetId(name + idSeparator + group + idSeparator + routingInstance)
		}
	} else {
		protocolsRipNeighbor := "rip group " + group + " neighbor"
		if ripNg {
			protocolsRipNeighbor = "ripng group " + group + " neighbor"
		}
		if routingInstance != defaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				protocolsRipNeighbor+" %v not exists in routing-instance %v after commit "+
					"=> check your config", name, routingInstance))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(protocolsRipNeighbor+" %v not exists after commit "+
			"=> check your config", name))...)
	}

	return append(diagWarns, resourceRipNeighborReadWJnprSess(d, sess, jnprSess)...)
}

func resourceRipNeighborRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceRipNeighborReadWJnprSess(d, sess, jnprSess)
}

func resourceRipNeighborReadWJnprSess(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject,
) diag.Diagnostics {
	mutex.Lock()
	ripNeighborOptions, err := readRipNeighbor(
		d.Get("name").(string),
		d.Get("group").(string),
		d.Get("ng").(bool),
		d.Get("routing_instance").(string),
		sess, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ripNeighborOptions.name == "" {
		d.SetId("")
	} else {
		fillRipNeighborData(d, ripNeighborOptions)
	}

	return nil
}

func resourceRipNeighborUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delRipNeighbor(
			d.Get("name").(string),
			d.Get("group").(string),
			d.Get("ng").(bool),
			d.Get("routing_instance").(string),
			sess, nil,
		); err != nil {
			return diag.FromErr(err)
		}
		if err := setRipNeighbor(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delRipNeighbor(
		d.Get("name").(string),
		d.Get("group").(string),
		d.Get("ng").(bool),
		d.Get("routing_instance").(string),
		sess, jnprSess,
	); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setRipNeighbor(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_rip_neighbor", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRipNeighborReadWJnprSess(d, sess, jnprSess)...)
}

func resourceRipNeighborDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delRipNeighbor(
			d.Get("name").(string),
			d.Get("group").(string),
			d.Get("ng").(bool),
			d.Get("routing_instance").(string),
			sess, nil,
		); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delRipNeighbor(
		d.Get("name").(string),
		d.Get("group").(string),
		d.Get("ng").(bool),
		d.Get("routing_instance").(string),
		sess, jnprSess,
	); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_rip_neighbor", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceRipNeighborImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 3 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	if len(idSplit) == 3 {
		ripNeighborExists, err := checkRipNeighborExists(idSplit[0], idSplit[1], false, idSplit[2], sess, jnprSess)
		if err != nil {
			return nil, err
		}
		if !ripNeighborExists {
			return nil, fmt.Errorf("don't find rip neighbor id '%v' "+
				"(id must be <name>%s<group>%s<routing_instance> or <name>%s<group>%sng%s<routing_instance>",
				d.Id(), idSeparator, idSeparator, idSeparator, idSeparator, idSeparator,
			)
		}
		ripNeighborOptions, err := readRipNeighbor(idSplit[0], idSplit[1], false, idSplit[2], sess, jnprSess)
		if err != nil {
			return nil, err
		}
		fillRipNeighborData(d, ripNeighborOptions)

		result[0] = d

		return result, nil
	}
	if idSplit[2] != "ng" {
		return nil, fmt.Errorf("id must be <name>%s<group>%s<routing_instance> or <name>%s<group>%sng%s<routing_instance>",
			idSeparator, idSeparator, idSeparator, idSeparator, idSeparator,
		)
	}
	ripNeighborExists, err := checkRipNeighborExists(idSplit[0], idSplit[1], true, idSplit[3], sess, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ripNeighborExists {
		return nil, fmt.Errorf("don't find ripng neighbor with id '%v' "+
			"(id must be <name>%s<group>%s<routing_instance> or <name>%s<group>%sng%s<routing_instance>",
			d.Id(), idSeparator, idSeparator, idSeparator, idSeparator, idSeparator,
		)
	}
	ripNeighborOptions, err := readRipNeighbor(idSplit[0], idSplit[1], true, idSplit[3], sess, jnprSess)
	if err != nil {
		return nil, err
	}
	fillRipNeighborData(d, ripNeighborOptions)
	result[0] = d

	return result, nil
}

func checkRipNeighborExists(
	name, group string,
	ripNg bool,
	routingInstance string,
	sess *Session,
	jnprSess *NetconfObject,
) (bool, error) {
	var showConfig string
	var err error
	protoRipNeighbor := "protocols rip group \"" + group + "\" neighbor " + name
	if ripNg {
		protoRipNeighbor = "protocols ripng group \"" + group + "\" neighbor " + name
	}
	if routingInstance == defaultW {
		showConfig, err = sess.command(cmdShowConfig+
			protoRipNeighbor+pipeDisplaySet, jnprSess)
	} else {
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			protoRipNeighbor+pipeDisplaySet, jnprSess)
	}
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setRipNeighbor(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject) error {
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
	setPrefix += "\"" + d.Get("group").(string) + "\" neighbor " + d.Get("name").(string) + " "

	configSet = append(configSet, setPrefix)
	if d.Get("any_sender").(bool) {
		configSet = append(configSet, setPrefix+"any-sender")
	}
	if v := d.Get("authentication_key").(string); v != "" {
		configSet = append(configSet, setPrefix+"authentication-key \""+v+"\"")
	}
	authSimpleMD5List := make([]int, 0)
	for _, authSimpMd5Block := range d.Get("authentication_selective_md5").([]interface{}) {
		authSimpMd5 := authSimpMd5Block.(map[string]interface{})
		if bchk.IntInSlice(authSimpMd5["key_id"].(int), authSimpleMD5List) {
			return fmt.Errorf("multiple blocks authentication_selective_md5 "+
				"with the same key_id %d", authSimpMd5["key_id"].(int))
		}
		authSimpleMD5List = append(authSimpleMD5List, authSimpMd5["key_id"].(int))
		configSet = append(configSet, setPrefix+"authentication-selective-md5 "+
			strconv.Itoa(authSimpMd5["key_id"].(int))+" key \""+authSimpMd5["key"].(string)+"\"")
		if v := authSimpMd5["start_time"].(string); v != "" {
			configSet = append(configSet, setPrefix+"authentication-selective-md5 "+
				strconv.Itoa(authSimpMd5["key_id"].(int))+" start-time "+v)
		}
	}
	if v := d.Get("authentication_type").(string); v != "" {
		configSet = append(configSet, setPrefix+"authentication-type "+v)
	}
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
	if d.Get("check_zero").(bool) {
		configSet = append(configSet, setPrefix+"check-zero")
	}
	if d.Get("demand_circuit").(bool) {
		configSet = append(configSet, setPrefix+"demand-circuit")
	}
	if d.Get("dynamic_peers").(bool) {
		configSet = append(configSet, setPrefix+"dynamic-peers")
	}
	for _, importPolicy := range d.Get("import").([]interface{}) {
		configSet = append(configSet, setPrefix+"import "+importPolicy.(string))
	}
	if d.Get("interface_type_p2mp").(bool) {
		configSet = append(configSet, setPrefix+"interface-type p2mp")
	}
	if v := d.Get("max_retrans_time").(int); v != 0 {
		configSet = append(configSet, setPrefix+"max-retrans-time "+strconv.Itoa(v))
	}
	if v := d.Get("message_size").(int); v != 0 {
		configSet = append(configSet, setPrefix+"message-size "+strconv.Itoa(v))
	}
	if v := d.Get("metric_in").(int); v != 0 {
		configSet = append(configSet, setPrefix+"metric-in "+strconv.Itoa(v))
	}
	if d.Get("no_check_zero").(bool) {
		configSet = append(configSet, setPrefix+"no-check-zero")
	}
	for _, peer := range sortSetOfString(d.Get("peer").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"peer "+peer)
	}
	if v := d.Get("receive").(string); v != "" {
		configSet = append(configSet, setPrefix+"receive "+v)
	}
	if v := d.Get("route_timeout").(int); v != 0 {
		configSet = append(configSet, setPrefix+"route-timeout "+strconv.Itoa(v))
	}
	if v := d.Get("send").(string); v != "" {
		configSet = append(configSet, setPrefix+"send "+v)
	}
	if v := d.Get("update_interval").(int); v != 0 {
		configSet = append(configSet, setPrefix+"update-interval "+strconv.Itoa(v))
	}

	return sess.configSet(configSet, jnprSess)
}

func readRipNeighbor(name, group string, ripNg bool, routingInstance string, sess *Session, jnprSess *NetconfObject,
) (ripNeighborOptions, error) {
	var confRead ripNeighborOptions
	var showConfig string
	var err error
	protoRipNeighbor := "protocols rip group \"" + group + "\" neighbor " + name
	if ripNg {
		protoRipNeighbor = "protocols ripng group \"" + group + "\" neighbor " + name
	}
	if routingInstance == defaultW {
		showConfig, err = sess.command(cmdShowConfig+
			protoRipNeighbor+pipeDisplaySetRelative, jnprSess)
	} else {
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			protoRipNeighbor+pipeDisplaySetRelative, jnprSess)
	}
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		confRead.group = group
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
			case itemTrim == "any-sender":
				confRead.anySender = true
			case strings.HasPrefix(itemTrim, "authentication-key "):
				confRead.authenticationKey, err = jdecode.Decode(
					strings.Trim(strings.TrimPrefix(itemTrim, "authentication-key "), "\""))
				if err != nil {
					return confRead, fmt.Errorf("failed to decode authentication-key: %w", err)
				}
			case strings.HasPrefix(itemTrim, "authentication-selective-md5 "):
				itemTrimplit := strings.Split(strings.TrimPrefix(itemTrim, "authentication-selective-md5 "), " ")
				keyID, err := strconv.Atoi(itemTrimplit[0])
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrimplit[0], err)
				}
				authSelectMD5 := map[string]interface{}{
					"key_id":     keyID,
					"key":        "",
					"start_time": "",
				}
				confRead.authenticationSelectiveMD5 = copyAndRemoveItemMapList(
					"key_id", authSelectMD5, confRead.authenticationSelectiveMD5)
				itemTrimAuthSelectMD5 := strings.TrimPrefix(itemTrim, "authentication-selective-md5 "+itemTrimplit[0]+" ")
				switch {
				case strings.HasPrefix(itemTrimAuthSelectMD5, "key "):
					var err error
					authSelectMD5["key"], err = jdecode.Decode(strings.Trim(strings.TrimPrefix(
						itemTrimAuthSelectMD5, "key "), "\""))
					if err != nil {
						return confRead, fmt.Errorf("failed to decode authentication-selective-md5 key: %w", err)
					}
				case strings.HasPrefix(itemTrimAuthSelectMD5, "start-time "):
					authSelectMD5["start_time"] = strings.Split(strings.Trim(strings.TrimPrefix(
						itemTrimAuthSelectMD5, "start-time "), "\""), " ")[0]
				}
				confRead.authenticationSelectiveMD5 = append(confRead.authenticationSelectiveMD5, authSelectMD5)
			case strings.HasPrefix(itemTrim, "authentication-type "):
				confRead.authenticationType = strings.TrimPrefix(itemTrim, "authentication-type ")
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
				if err := readRipNeighborBfd(
					strings.TrimPrefix(itemTrim, "bfd-liveness-detection "),
					confRead.bfdLivenessDetection[0],
				); err != nil {
					return confRead, err
				}
			case itemTrim == "check-zero":
				confRead.checkZero = true
			case itemTrim == "demand-circuit":
				confRead.demandCircuit = true
			case itemTrim == "dynamic-peers":
				confRead.dynamicPeers = true
			case strings.HasPrefix(itemTrim, "import "):
				confRead.importPolicy = append(confRead.importPolicy, strings.TrimPrefix(itemTrim, "import "))
			case itemTrim == "interface-type p2mp":
				confRead.interfaceTypeP2mp = true
			case strings.HasPrefix(itemTrim, "max-retrans-time "):
				confRead.maxRetransTime, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-retrans-time "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "message-size "):
				confRead.messageSize, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "message-size "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "metric-in "):
				confRead.metricIn, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "metric-in "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "no-check-zero":
				confRead.noCheckZero = true
			case strings.HasPrefix(itemTrim, "peer "):
				confRead.peer = append(confRead.peer, strings.TrimPrefix(itemTrim, "peer "))
			case strings.HasPrefix(itemTrim, "receive "):
				confRead.receive = strings.TrimPrefix(itemTrim, "receive ")
			case strings.HasPrefix(itemTrim, "route-timeout "):
				confRead.routeTimeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "route-timeout "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "send "):
				confRead.send = strings.TrimPrefix(itemTrim, "send ")
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

func readRipNeighborBfd(itemTrim string, bfd map[string]interface{}) error {
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

func delRipNeighbor(
	name, group string,
	ripNg bool,
	routingInstance string,
	sess *Session,
	jnprSess *NetconfObject,
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
	delPrefix += "\"" + group + "\" neighbor " + name + " "

	return sess.configSet([]string{delPrefix}, jnprSess)
}

func fillRipNeighborData(d *schema.ResourceData, ripNeighborOptions ripNeighborOptions) {
	if tfErr := d.Set("name", ripNeighborOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("group", ripNeighborOptions.group); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", ripNeighborOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ng", ripNeighborOptions.ng); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("any_sender", ripNeighborOptions.anySender); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_key", ripNeighborOptions.authenticationKey); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_selective_md5", ripNeighborOptions.authenticationSelectiveMD5); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_type", ripNeighborOptions.authenticationType); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bfd_liveness_detection", ripNeighborOptions.bfdLivenessDetection); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("check_zero", ripNeighborOptions.checkZero); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("demand_circuit", ripNeighborOptions.demandCircuit); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dynamic_peers", ripNeighborOptions.dynamicPeers); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("import", ripNeighborOptions.importPolicy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interface_type_p2mp", ripNeighborOptions.interfaceTypeP2mp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("max_retrans_time", ripNeighborOptions.maxRetransTime); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("message_size", ripNeighborOptions.messageSize); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_in", ripNeighborOptions.metricIn); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_check_zero", ripNeighborOptions.noCheckZero); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("peer", ripNeighborOptions.peer); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("receive", ripNeighborOptions.receive); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("route_timeout", ripNeighborOptions.routeTimeout); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("send", ripNeighborOptions.send); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("update_interval", ripNeighborOptions.updateInterval); tfErr != nil {
		panic(tfErr)
	}
}
