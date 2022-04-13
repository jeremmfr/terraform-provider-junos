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

type generateRouteOptions struct {
	active                   bool
	asPathAtomicAggregate    bool
	brief                    bool
	discard                  bool
	full                     bool
	passive                  bool
	metric                   int
	preference               int
	asPathAggregatorAddress  string
	asPathAggregatorAsNumber string
	asPathOrigin             string
	asPathPath               string
	destination              string
	nextTable                string
	routingInstance          string
	community                []string
	policy                   []string
}

func resourceGenerateRoute() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceGenerateRouteCreate,
		ReadWithoutTimeout:   resourceGenerateRouteRead,
		UpdateWithoutTimeout: resourceGenerateRouteUpdate,
		DeleteWithoutTimeout: resourceGenerateRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGenerateRouteImport,
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
				Default:          defaultW,
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
			"brief": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"full"},
			},
			"community": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"discard": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"next_table"},
			},
			"full": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"brief"},
			},
			"metric": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"next_table": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"discard"},
			},
			"passive": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"active"},
			},
			"policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"preference": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceGenerateRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setGenerateRoute(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("destination").(string) + idSeparator + d.Get("routing_instance").(string))

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
	if d.Get("routing_instance").(string) != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), m, jnprSess)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	generateRouteExists, err := checkGenerateRouteExists(d.Get("destination").(string), d.Get("routing_instance").(string),
		m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if generateRouteExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("generate route %v already exists on table %s",
			d.Get("destination").(string), d.Get("routing_instance").(string)))...)
	}
	if err := setGenerateRoute(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_generate_route", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	generateRouteExists, err = checkGenerateRouteExists(d.Get("destination").(string), d.Get("routing_instance").(string),
		m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if generateRouteExists {
		d.SetId(d.Get("destination").(string) + idSeparator + d.Get("routing_instance").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("generate route %v not exists in routing_instance %v after commit "+
			"=> check your config", d.Get("destination").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceGenerateRouteReadWJnprSess(d, m, jnprSess)...)
}

func resourceGenerateRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceGenerateRouteReadWJnprSess(d, m, jnprSess)
}

func resourceGenerateRouteReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
) diag.Diagnostics {
	mutex.Lock()
	generateRouteOptions, err := readGenerateRoute(d.Get("destination").(string), d.Get("routing_instance").(string),
		m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if generateRouteOptions.destination == "" {
		d.SetId("")
	} else {
		fillGenerateRouteData(d, generateRouteOptions)
	}

	return nil
}

func resourceGenerateRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delGenerateRoute(
			d.Get("destination").(string), d.Get("routing_instance").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setGenerateRoute(d, m, nil); err != nil {
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
	if err := delGenerateRoute(
		d.Get("destination").(string), d.Get("routing_instance").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setGenerateRoute(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_generate_route", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	d.Partial(false)

	return append(diagWarns, resourceGenerateRouteReadWJnprSess(d, m, jnprSess)...)
}

func resourceGenerateRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delGenerateRoute(
			d.Get("destination").(string), d.Get("routing_instance").(string), m, nil); err != nil {
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
	if err := delGenerateRoute(
		d.Get("destination").(string), d.Get("routing_instance").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_generate_route", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceGenerateRouteImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	generateRouteExists, err := checkGenerateRouteExists(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !generateRouteExists {
		return nil, fmt.Errorf("don't find generate route with id '%v' (id must be "+
			"<destination>"+idSeparator+"<routing_instance>)", d.Id())
	}
	generateRouteOptions, err := readGenerateRoute(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillGenerateRouteData(d, generateRouteOptions)

	result[0] = d

	return result, nil
}

func checkGenerateRouteExists(destination, instance string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	var showConfig string
	var err error
	if instance == defaultW {
		if !strings.Contains(destination, ":") {
			showConfig, err = sess.command(cmdShowConfig+
				"routing-options generate route "+destination+pipeDisplaySet, jnprSess)
			if err != nil {
				return false, err
			}
		} else {
			showConfig, err = sess.command(cmdShowConfig+
				"routing-options rib inet6.0 "+"generate route "+destination+pipeDisplaySet, jnprSess)
			if err != nil {
				return false, err
			}
		}
	} else {
		if !strings.Contains(destination, ":") {
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+instance+" "+
				"routing-options generate route "+destination+pipeDisplaySet, jnprSess)
			if err != nil {
				return false, err
			}
		} else {
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+instance+" "+
				"routing-options rib "+instance+".inet6.0 generate route "+destination+pipeDisplaySet, jnprSess)
			if err != nil {
				return false, err
			}
		}
	}

	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setGenerateRoute(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	var setPrefix string
	if d.Get("routing_instance").(string) == defaultW {
		if !strings.Contains(d.Get("destination").(string), ":") {
			setPrefix = "set routing-options generate route " + d.Get("destination").(string) + " "
		} else {
			setPrefix = "set routing-options rib inet6.0 generate route " + d.Get("destination").(string) + " "
		}
	} else {
		if !strings.Contains(d.Get("destination").(string), ":") {
			setPrefix = setRoutingInstances + d.Get("routing_instance").(string) +
				" routing-options generate route " + d.Get("destination").(string) + " "
		} else {
			setPrefix = setRoutingInstances + d.Get("routing_instance").(string) +
				" routing-options rib " + d.Get("routing_instance").(string) + ".inet6.0 " +
				"generate route " + d.Get("destination").(string) + " "
		}
	}
	if d.Get("active").(bool) {
		configSet = append(configSet, setPrefix+"active")
	}
	if d.Get("as_path_aggregator_address").(string) != "" &&
		d.Get("as_path_aggregator_as_number").(string) != "" {
		configSet = append(configSet, setPrefix+"as-path aggregator "+
			d.Get("as_path_aggregator_as_number").(string)+" "+
			d.Get("as_path_aggregator_address").(string))
	}
	if d.Get("as_path_atomic_aggregate").(bool) {
		configSet = append(configSet, setPrefix+"as-path atomic-aggregate")
	}
	if v := d.Get("as_path_origin").(string); v != "" {
		configSet = append(configSet, setPrefix+"as-path origin "+v)
	}
	if v := d.Get("as_path_path").(string); v != "" {
		configSet = append(configSet, setPrefix+"as-path path \""+v+"\"")
	}
	if d.Get("brief").(bool) {
		configSet = append(configSet, setPrefix+"brief")
	}
	for _, v := range d.Get("community").([]interface{}) {
		configSet = append(configSet, setPrefix+"community "+v.(string))
	}
	if d.Get("discard").(bool) {
		configSet = append(configSet, setPrefix+"discard")
	}
	if d.Get("full").(bool) {
		configSet = append(configSet, setPrefix+"full")
	}
	if d.Get("metric").(int) > 0 {
		configSet = append(configSet, setPrefix+"metric "+strconv.Itoa(d.Get("metric").(int)))
	}
	if d.Get("next_table").(string) != "" {
		configSet = append(configSet, setPrefix+"next-table "+d.Get("next_table").(string))
	}
	if d.Get("passive").(bool) {
		configSet = append(configSet, setPrefix+"passive")
	}
	for _, v := range d.Get("policy").([]interface{}) {
		configSet = append(configSet, setPrefix+"policy "+v.(string))
	}
	if d.Get("preference").(int) > 0 {
		configSet = append(configSet, setPrefix+"preference "+strconv.Itoa(d.Get("preference").(int)))
	}

	return sess.configSet(configSet, jnprSess)
}

func readGenerateRoute(destination, instance string, m interface{}, jnprSess *NetconfObject,
) (generateRouteOptions, error) {
	sess := m.(*Session)
	var confRead generateRouteOptions
	var showConfig string
	var err error

	if instance == defaultW {
		if !strings.Contains(destination, ":") {
			showConfig, err = sess.command(cmdShowConfig+
				"routing-options generate route "+destination+pipeDisplaySetRelative, jnprSess)
		} else {
			showConfig, err = sess.command(cmdShowConfig+
				"routing-options rib inet6.0 generate route "+destination+pipeDisplaySetRelative, jnprSess)
		}
	} else {
		if !strings.Contains(destination, ":") {
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+instance+" "+
				"routing-options generate route "+destination+pipeDisplaySetRelative, jnprSess)
		} else {
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+instance+" "+
				"routing-options rib "+instance+".inet6.0 generate route "+destination+pipeDisplaySetRelative, jnprSess)
		}
	}
	if err != nil {
		return confRead, err
	}

	if showConfig != emptyW {
		confRead.destination = destination
		confRead.routingInstance = instance
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case itemTrim == "active":
				confRead.active = true
			case strings.HasPrefix(itemTrim, "as-path aggregator "):
				itemTrimSplit := strings.Split(itemTrim, " ")
				confRead.asPathAggregatorAsNumber = itemTrimSplit[2]
				confRead.asPathAggregatorAddress = itemTrimSplit[3]
			case itemTrim == "as-path atomic-aggregate":
				confRead.asPathAtomicAggregate = true
			case strings.HasPrefix(itemTrim, "as-path origin "):
				confRead.asPathOrigin = strings.TrimPrefix(itemTrim, "as-path origin ")
			case strings.HasPrefix(itemTrim, "as-path path "):
				confRead.asPathPath = strings.Trim(strings.TrimPrefix(itemTrim, "as-path path "), "\"")
			case itemTrim == "brief":
				confRead.brief = true
			case strings.HasPrefix(itemTrim, "community "):
				confRead.community = append(confRead.community, strings.TrimPrefix(itemTrim, "community "))
			case itemTrim == discardW:
				confRead.discard = true
			case itemTrim == "full":
				confRead.full = true
			case strings.HasPrefix(itemTrim, "metric "):
				confRead.metric, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "metric "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "next-table "):
				confRead.nextTable = strings.TrimPrefix(itemTrim, "next-table ")
			case itemTrim == "passive":
				confRead.passive = true
			case strings.HasPrefix(itemTrim, "policy "):
				confRead.policy = append(confRead.policy, strings.TrimPrefix(itemTrim, "policy "))
			case strings.HasPrefix(itemTrim, "preference "):
				confRead.preference, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "preference "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delGenerateRoute(destination, instance string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	if instance == defaultW {
		if !strings.Contains(destination, ":") {
			configSet = append(configSet, "delete routing-options generate route "+destination)
		} else {
			configSet = append(configSet, "delete routing-options rib inet6.0 generate route "+destination)
		}
	} else {
		if !strings.Contains(destination, ":") {
			configSet = append(configSet, delRoutingInstances+instance+" "+
				"routing-options generate route "+destination)
		} else {
			configSet = append(configSet, delRoutingInstances+instance+" "+
				"routing-options rib "+instance+".inet6.0 generate route "+destination)
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func fillGenerateRouteData(d *schema.ResourceData, generateRouteOptions generateRouteOptions) {
	if tfErr := d.Set("destination", generateRouteOptions.destination); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", generateRouteOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("active", generateRouteOptions.active); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_aggregator_address", generateRouteOptions.asPathAggregatorAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_aggregator_as_number", generateRouteOptions.asPathAggregatorAsNumber); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_atomic_aggregate", generateRouteOptions.asPathAtomicAggregate); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_origin", generateRouteOptions.asPathOrigin); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_path", generateRouteOptions.asPathPath); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("brief", generateRouteOptions.brief); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("community", generateRouteOptions.community); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("discard", generateRouteOptions.discard); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("discard", generateRouteOptions.discard); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("full", generateRouteOptions.full); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric", generateRouteOptions.metric); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("next_table", generateRouteOptions.nextTable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("passive", generateRouteOptions.passive); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("policy", generateRouteOptions.policy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("preference", generateRouteOptions.preference); tfErr != nil {
		panic(tfErr)
	}
}
