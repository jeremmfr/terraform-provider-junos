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
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type staticRouteOptions struct {
	active                   bool
	asPathAtomicAggregate    bool
	discard                  bool
	install                  bool
	noInstall                bool
	passive                  bool
	readvertise              bool
	noReadvertise            bool
	receive                  bool
	reject                   bool
	resolve                  bool
	noResolve                bool
	retain                   bool
	noRetain                 bool
	preference               int
	metric                   int
	asPathAggregatorAddress  string
	asPathAggregatorAsNumber string
	asPathOrigin             string
	asPathPath               string
	destination              string
	routingInstance          string
	nextTable                string
	community                []string
	nextHop                  []string
	qualifiedNextHop         []map[string]interface{}
}

func resourceStaticRoute() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceStaticRouteCreate,
		ReadWithoutTimeout:   resourceStaticRouteRead,
		UpdateWithoutTimeout: resourceStaticRouteUpdate,
		DeleteWithoutTimeout: resourceStaticRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceStaticRouteImport,
		},
		Schema: map[string]*schema.Schema{
			"destination": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 128),
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          junos.DefaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"active": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"passive"},
			},
			"as_path_aggregator_address": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"as_path_aggregator_as_number"},
				ValidateFunc: validation.IsIPAddress,
			},
			"as_path_aggregator_as_number": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"as_path_aggregator_address"},
			},
			"as_path_atomic_aggregate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"as_path_origin": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"egp", "igp", "incomplete"}, false),
			},
			"as_path_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"community": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"discard": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"receive", "reject", "next_hop", "next_table", "qualified_next_hop"},
			},
			"install": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"no_install"},
			},
			"no_install": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"install"},
			},
			"metric": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"next_hop": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"next_table", "discard", "receive", "reject"},
			},
			"next_table": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"next_hop", "qualified_next_hop", "discard", "receive", "reject"},
			},
			"passive": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"active"},
			},
			"preference": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"qualified_next_hop": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"next_table", "discard", "receive", "reject"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"next_hop": {
							Type:     schema.TypeString,
							Required: true,
						},
						"interface": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"metric": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"preference": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"readvertise": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"no_readvertise"},
			},
			"no_readvertise": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"readvertise"},
			},
			"receive": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"discard", "reject", "next_hop", "next_table", "qualified_next_hop"},
			},
			"reject": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"discard", "receive", "next_hop", "next_table", "qualified_next_hop"},
			},
			"resolve": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"no_resolve", "retain", "no_retain"},
			},
			"no_resolve": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"resolve"},
			},
			"retain": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"no_retain", "resolve"},
			},
			"no_retain": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"retain", "resolve"},
			},
		},
	}
}

func resourceStaticRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setStaticRoute(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("destination").(string) + junos.IDSeparator + d.Get("routing_instance").(string))

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
	staticRouteExists, err := checkStaticRouteExists(
		d.Get("destination").(string),
		d.Get("routing_instance").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if staticRouteExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("static route %v already exists on table %s",
			d.Get("destination").(string), d.Get("routing_instance").(string)))...)
	}
	if err := setStaticRoute(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_static_route")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	staticRouteExists, err = checkStaticRouteExists(
		d.Get("destination").(string),
		d.Get("routing_instance").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if staticRouteExists {
		d.SetId(d.Get("destination").(string) + junos.IDSeparator + d.Get("routing_instance").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("static route %v not exists in routing_instance %v after commit "+
			"=> check your config", d.Get("destination").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceStaticRouteReadWJunSess(d, junSess)...)
}

func resourceStaticRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceStaticRouteReadWJunSess(d, junSess)
}

func resourceStaticRouteReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	staticRouteOptions, err := readStaticRoute(
		d.Get("destination").(string),
		d.Get("routing_instance").(string),
		junSess,
	)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if staticRouteOptions.destination == "" {
		d.SetId("")
	} else {
		fillStaticRouteData(d, staticRouteOptions)
	}

	return nil
}

func resourceStaticRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delStaticRoute(d.Get("destination").(string), d.Get("routing_instance").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setStaticRoute(d, junSess); err != nil {
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
	if err := delStaticRoute(
		d.Get("destination").(string),
		d.Get("routing_instance").(string),
		junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setStaticRoute(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_static_route")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	d.Partial(false)

	return append(diagWarns, resourceStaticRouteReadWJunSess(d, junSess)...)
}

func resourceStaticRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delStaticRoute(d.Get("destination").(string), d.Get("routing_instance").(string), junSess); err != nil {
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
	if err := delStaticRoute(
		d.Get("destination").(string),
		d.Get("routing_instance").(string),
		junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_static_route")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceStaticRouteImport(ctx context.Context, d *schema.ResourceData, m interface{},
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
	staticRouteExists, err := checkStaticRouteExists(idSplit[0], idSplit[1], junSess)
	if err != nil {
		return nil, err
	}
	if !staticRouteExists {
		return nil, fmt.Errorf("don't find static route with id '%v' (id must be "+
			"<destination>"+junos.IDSeparator+"<routing_instance>)", d.Id())
	}
	staticRouteOptions, err := readStaticRoute(idSplit[0], idSplit[1], junSess)
	if err != nil {
		return nil, err
	}
	fillStaticRouteData(d, staticRouteOptions)

	result[0] = d

	return result, nil
}

func checkStaticRouteExists(destination, instance string, junSess *junos.Session,
) (_ bool, err error) {
	var showConfig string
	if instance == junos.DefaultW {
		if !strings.Contains(destination, ":") {
			showConfig, err = junSess.Command(junos.CmdShowConfig +
				"routing-options static route " + destination + junos.PipeDisplaySet)
			if err != nil {
				return false, err
			}
		} else {
			showConfig, err = junSess.Command(junos.CmdShowConfig +
				"routing-options rib inet6.0 static route " + destination + junos.PipeDisplaySet)
			if err != nil {
				return false, err
			}
		}
	} else {
		if !strings.Contains(destination, ":") {
			showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + instance + " " +
				"routing-options static route " + destination + junos.PipeDisplaySet)
			if err != nil {
				return false, err
			}
		} else {
			showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + instance + " " +
				"routing-options rib " + instance + ".inet6.0 static route " + destination + junos.PipeDisplaySet)
			if err != nil {
				return false, err
			}
		}
	}

	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setStaticRoute(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	var setPrefix string
	if d.Get("routing_instance").(string) == junos.DefaultW {
		if !strings.Contains(d.Get("destination").(string), ":") {
			setPrefix = "set routing-options static route " + d.Get("destination").(string)
		} else {
			setPrefix = "set routing-options rib inet6.0 static route " + d.Get("destination").(string)
		}
	} else {
		if !strings.Contains(d.Get("destination").(string), ":") {
			setPrefix = junos.SetRoutingInstances + d.Get("routing_instance").(string) +
				" routing-options static route " + d.Get("destination").(string)
		} else {
			setPrefix = junos.SetRoutingInstances + d.Get("routing_instance").(string) +
				" routing-options rib " + d.Get("routing_instance").(string) + ".inet6.0 " +
				"static route " + d.Get("destination").(string)
		}
	}
	if d.Get("active").(bool) {
		configSet = append(configSet, setPrefix+" active")
	}
	if d.Get("as_path_aggregator_address").(string) != "" &&
		d.Get("as_path_aggregator_as_number").(string) != "" {
		configSet = append(configSet, setPrefix+" as-path aggregator "+
			d.Get("as_path_aggregator_as_number").(string)+" "+
			d.Get("as_path_aggregator_address").(string))
	}
	if d.Get("as_path_atomic_aggregate").(bool) {
		configSet = append(configSet, setPrefix+" as-path atomic-aggregate")
	}
	if v := d.Get("as_path_origin").(string); v != "" {
		configSet = append(configSet, setPrefix+" as-path origin "+v)
	}
	if v := d.Get("as_path_path").(string); v != "" {
		configSet = append(configSet, setPrefix+" as-path path \""+v+"\"")
	}
	for _, v := range d.Get("community").([]interface{}) {
		configSet = append(configSet, setPrefix+" community "+v.(string))
	}
	if d.Get("discard").(bool) {
		configSet = append(configSet, setPrefix+" discard")
	}
	if d.Get("install").(bool) {
		configSet = append(configSet, setPrefix+" install")
	}
	if d.Get("no_install").(bool) {
		configSet = append(configSet, setPrefix+" no-install")
	}
	if d.Get("metric").(int) > 0 {
		configSet = append(configSet, setPrefix+" metric "+strconv.Itoa(d.Get("metric").(int)))
	}
	for _, nextHop := range d.Get("next_hop").([]interface{}) {
		configSet = append(configSet, setPrefix+" next-hop "+nextHop.(string))
	}
	if d.Get("next_table").(string) != "" {
		configSet = append(configSet, setPrefix+" next-table "+d.Get("next_table").(string))
	}
	if d.Get("passive").(bool) {
		configSet = append(configSet, setPrefix+" passive")
	}
	if d.Get("preference").(int) > 0 {
		configSet = append(configSet, setPrefix+" preference "+strconv.Itoa(d.Get("preference").(int)))
	}
	qualifiedNextHopList := make([]string, 0)
	for _, qualifiedNextHop := range d.Get("qualified_next_hop").([]interface{}) {
		qualifiedNextHopMap := qualifiedNextHop.(map[string]interface{})
		if bchk.InSlice(qualifiedNextHopMap["next_hop"].(string), qualifiedNextHopList) {
			return fmt.Errorf("multiple blocks qualified_next_hop with the same next_hop %s",
				qualifiedNextHopMap["next_hop"].(string))
		}
		qualifiedNextHopList = append(qualifiedNextHopList, qualifiedNextHopMap["next_hop"].(string))
		configSet = append(configSet, setPrefix+" qualified-next-hop "+qualifiedNextHopMap["next_hop"].(string))
		if qualifiedNextHopMap["interface"] != "" {
			configSet = append(configSet, setPrefix+
				" qualified-next-hop "+qualifiedNextHopMap["next_hop"].(string)+
				" interface "+qualifiedNextHopMap["interface"].(string))
		}
		if qualifiedNextHopMap["metric"].(int) > 0 {
			configSet = append(configSet, setPrefix+
				" qualified-next-hop "+qualifiedNextHopMap["next_hop"].(string)+
				" metric "+strconv.Itoa(qualifiedNextHopMap["metric"].(int)))
		}
		if qualifiedNextHopMap["preference"].(int) > 0 {
			configSet = append(configSet, setPrefix+
				" qualified-next-hop "+qualifiedNextHopMap["next_hop"].(string)+
				" preference "+strconv.Itoa(qualifiedNextHopMap["preference"].(int)))
		}
	}
	if d.Get("readvertise").(bool) {
		configSet = append(configSet, setPrefix+" readvertise")
	}
	if d.Get("no_readvertise").(bool) {
		configSet = append(configSet, setPrefix+" no-readvertise")
	}
	if d.Get("receive").(bool) {
		configSet = append(configSet, setPrefix+" receive")
	}
	if d.Get("reject").(bool) {
		configSet = append(configSet, setPrefix+" reject")
	}
	if d.Get("resolve").(bool) {
		configSet = append(configSet, setPrefix+" resolve")
	}
	if d.Get("no_resolve").(bool) {
		configSet = append(configSet, setPrefix+" no-resolve")
	}
	if d.Get("retain").(bool) {
		configSet = append(configSet, setPrefix+" retain")
	}
	if d.Get("no_retain").(bool) {
		configSet = append(configSet, setPrefix+" no-retain")
	}

	return junSess.ConfigSet(configSet)
}

func readStaticRoute(destination, instance string, junSess *junos.Session,
) (confRead staticRouteOptions, err error) {
	var showConfig string
	if instance == junos.DefaultW {
		if !strings.Contains(destination, ":") {
			showConfig, err = junSess.Command(junos.CmdShowConfig +
				"routing-options static route " + destination + junos.PipeDisplaySetRelative)
		} else {
			showConfig, err = junSess.Command(junos.CmdShowConfig +
				"routing-options rib inet6.0 static route " + destination + junos.PipeDisplaySetRelative)
		}
	} else {
		if !strings.Contains(destination, ":") {
			showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + instance + " " +
				"routing-options static route " + destination + junos.PipeDisplaySetRelative)
		} else {
			showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + instance + " " +
				"routing-options rib " + instance + ".inet6.0 static route " + destination + junos.PipeDisplaySetRelative)
		}
	}
	if err != nil {
		return confRead, err
	}

	if showConfig != junos.EmptyW {
		confRead.destination = destination
		confRead.routingInstance = instance
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == "active":
				confRead.active = true
			case balt.CutPrefixInString(&itemTrim, "as-path aggregator "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 2 { // <as_number> <address>
					return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "as-path aggregator", itemTrim)
				}
				confRead.asPathAggregatorAsNumber = itemTrimFields[0]
				confRead.asPathAggregatorAddress = itemTrimFields[1]
			case itemTrim == "as-path atomic-aggregate":
				confRead.asPathAtomicAggregate = true
			case balt.CutPrefixInString(&itemTrim, "as-path origin "):
				confRead.asPathOrigin = itemTrim
			case balt.CutPrefixInString(&itemTrim, "as-path path "):
				confRead.asPathPath = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "community "):
				confRead.community = append(confRead.community, itemTrim)
			case itemTrim == junos.DiscardW:
				confRead.discard = true
			case itemTrim == "install":
				confRead.install = true
			case itemTrim == "no-install":
				confRead.noInstall = true
			case balt.CutPrefixInString(&itemTrim, "metric "):
				confRead.metric, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "next-hop "):
				confRead.nextHop = append(confRead.nextHop, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "next-table "):
				confRead.nextTable = itemTrim
			case itemTrim == "passive":
				confRead.passive = true
			case balt.CutPrefixInString(&itemTrim, "preference "):
				confRead.preference, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "qualified-next-hop "):
				itemTrimFields := strings.Split(itemTrim, " ")
				qualifiedNextHopOptions := map[string]interface{}{
					"next_hop":   itemTrimFields[0],
					"interface":  "",
					"metric":     0,
					"preference": 0,
				}
				confRead.qualifiedNextHop = copyAndRemoveItemMapList("next_hop", qualifiedNextHopOptions, confRead.qualifiedNextHop)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "interface "):
					qualifiedNextHopOptions["interface"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, "metric "):
					qualifiedNextHopOptions["metric"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "preference "):
					qualifiedNextHopOptions["preference"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				}
				confRead.qualifiedNextHop = append(confRead.qualifiedNextHop, qualifiedNextHopOptions)
			case itemTrim == "readvertise":
				confRead.readvertise = true
			case itemTrim == "no-readvertise":
				confRead.noReadvertise = true
			case itemTrim == "receive":
				confRead.receive = true
			case itemTrim == "reject":
				confRead.reject = true
			case itemTrim == "resolve":
				confRead.resolve = true
			case itemTrim == "no-resolve":
				confRead.noResolve = true
			case itemTrim == "retain":
				confRead.retain = true
			case itemTrim == "no-retain":
				confRead.noRetain = true
			}
		}
	}

	return confRead, nil
}

func delStaticRoute(destination, instance string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	switch {
	case instance == junos.DefaultW && !strings.Contains(destination, ":"):
		configSet = append(configSet, "delete routing-options static route "+destination)
	case instance == junos.DefaultW && strings.Contains(destination, ":"):
		configSet = append(configSet, "delete routing-options rib inet6.0 static route "+destination)
	case instance != junos.DefaultW && !strings.Contains(destination, ":"):
		configSet = append(configSet, junos.DelRoutingInstances+instance+" "+
			"routing-options static route "+destination)
	case instance != junos.DefaultW && strings.Contains(destination, ":"):
		configSet = append(configSet, junos.DelRoutingInstances+instance+" "+
			"routing-options rib "+instance+".inet6.0 static route "+destination)
	}

	return junSess.ConfigSet(configSet)
}

func fillStaticRouteData(d *schema.ResourceData, staticRouteOptions staticRouteOptions) {
	if tfErr := d.Set("destination", staticRouteOptions.destination); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", staticRouteOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("active", staticRouteOptions.active); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_aggregator_address", staticRouteOptions.asPathAggregatorAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_aggregator_as_number", staticRouteOptions.asPathAggregatorAsNumber); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_atomic_aggregate", staticRouteOptions.asPathAtomicAggregate); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_origin", staticRouteOptions.asPathOrigin); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_path", staticRouteOptions.asPathPath); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("community", staticRouteOptions.community); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("discard", staticRouteOptions.discard); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("install", staticRouteOptions.install); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_install", staticRouteOptions.noInstall); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric", staticRouteOptions.metric); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("next_hop", staticRouteOptions.nextHop); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("next_table", staticRouteOptions.nextTable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("passive", staticRouteOptions.passive); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("preference", staticRouteOptions.preference); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("qualified_next_hop", staticRouteOptions.qualifiedNextHop); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("readvertise", staticRouteOptions.readvertise); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_readvertise", staticRouteOptions.noReadvertise); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("receive", staticRouteOptions.receive); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("reject", staticRouteOptions.reject); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("resolve", staticRouteOptions.resolve); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_resolve", staticRouteOptions.noResolve); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("retain", staticRouteOptions.retain); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_retain", staticRouteOptions.noRetain); tfErr != nil {
		panic(tfErr)
	}
}
