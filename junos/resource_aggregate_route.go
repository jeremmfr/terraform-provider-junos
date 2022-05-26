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

type aggregateRouteOptions struct {
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
	routingInstance          string
	community                []string
	policy                   []string
}

func resourceAggregateRoute() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAggregateRouteCreate,
		ReadWithoutTimeout:   resourceAggregateRouteRead,
		UpdateWithoutTimeout: resourceAggregateRouteUpdate,
		DeleteWithoutTimeout: resourceAggregateRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAggregateRouteImport,
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
				Type:     schema.TypeBool,
				Optional: true,
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

func resourceAggregateRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setAggregateRoute(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("destination").(string) + idSeparator + d.Get("routing_instance").(string))

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
	if d.Get("routing_instance").(string) != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), sess, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, sess.configClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	aggregateRouteExists, err := checkAggregateRouteExists(
		d.Get("destination").(string),
		d.Get("routing_instance").(string),
		sess, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if aggregateRouteExists {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("aggregate route %v already exists on table %s",
			d.Get("destination").(string), d.Get("routing_instance").(string)))...)
	}
	if err := setAggregateRoute(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_aggregate_route", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	aggregateRouteExists, err = checkAggregateRouteExists(
		d.Get("destination").(string),
		d.Get("routing_instance").(string),
		sess, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if aggregateRouteExists {
		d.SetId(d.Get("destination").(string) + idSeparator + d.Get("routing_instance").(string))
	} else {
		return append(diagWarns,
			diag.FromErr(fmt.Errorf("aggregate route %v not exists in routing_instance %v after commit "+
				"=> check your config", d.Get("destination").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceAggregateRouteReadWJunSess(d, sess, junSess)...)
}

func resourceAggregateRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)

	return resourceAggregateRouteReadWJunSess(d, sess, junSess)
}

func resourceAggregateRouteReadWJunSess(d *schema.ResourceData, sess *Session, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	aggregateRouteOptions, err := readAggregateRoute(
		d.Get("destination").(string),
		d.Get("routing_instance").(string),
		sess, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if aggregateRouteOptions.destination == "" {
		d.SetId("")
	} else {
		fillAggregateRouteData(d, aggregateRouteOptions)
	}

	return nil
}

func resourceAggregateRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delAggregateRoute(
			d.Get("destination").(string),
			d.Get("routing_instance").(string),
			sess, nil,
		); err != nil {
			return diag.FromErr(err)
		}
		if err := setAggregateRoute(d, sess, nil); err != nil {
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
	if err := delAggregateRoute(
		d.Get("destination").(string),
		d.Get("routing_instance").(string),
		sess, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setAggregateRoute(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_aggregate_route", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceAggregateRouteReadWJunSess(d, sess, junSess)...)
}

func resourceAggregateRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delAggregateRoute(
			d.Get("destination").(string),
			d.Get("routing_instance").(string),
			sess, nil,
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
	if err := delAggregateRoute(
		d.Get("destination").(string),
		d.Get("routing_instance").(string),
		sess, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_aggregate_route", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceAggregateRouteImport(ctx context.Context, d *schema.ResourceData, m interface{},
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
	aggregateRouteExists, err := checkAggregateRouteExists(idSplit[0], idSplit[1], sess, junSess)
	if err != nil {
		return nil, err
	}
	if !aggregateRouteExists {
		return nil, fmt.Errorf("don't find aggregate route with id '%v' (id must be "+
			"<destination>"+idSeparator+"<routing_instance>)", d.Id())
	}
	aggregateRouteOptions, err := readAggregateRoute(idSplit[0], idSplit[1], sess, junSess)
	if err != nil {
		return nil, err
	}
	fillAggregateRouteData(d, aggregateRouteOptions)

	result[0] = d

	return result, nil
}

func checkAggregateRouteExists(destination, instance string, sess *Session, junSess *junosSession) (bool, error) {
	var showConfig string
	var err error
	if instance == defaultW {
		if !strings.Contains(destination, ":") {
			showConfig, err = sess.command(cmdShowConfig+
				"routing-options aggregate route "+destination+pipeDisplaySet, junSess)
			if err != nil {
				return false, err
			}
		} else {
			showConfig, err = sess.command(cmdShowConfig+
				"routing-options rib inet6.0 aggregate route "+destination+pipeDisplaySet, junSess)
			if err != nil {
				return false, err
			}
		}
	} else {
		if !strings.Contains(destination, ":") {
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+instance+" "+
				"routing-options aggregate route "+destination+pipeDisplaySet, junSess)
			if err != nil {
				return false, err
			}
		} else {
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+instance+" "+
				"routing-options rib "+instance+".inet6.0 aggregate route "+destination+pipeDisplaySet, junSess)
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

func setAggregateRoute(d *schema.ResourceData, sess *Session, junSess *junosSession) error {
	configSet := make([]string, 0)

	var setPrefix string
	if d.Get("routing_instance").(string) == defaultW {
		if !strings.Contains(d.Get("destination").(string), ":") {
			setPrefix = "set routing-options aggregate route " + d.Get("destination").(string)
		} else {
			setPrefix = "set routing-options rib inet6.0 aggregate route " + d.Get("destination").(string)
		}
	} else {
		if !strings.Contains(d.Get("destination").(string), ":") {
			setPrefix = setRoutingInstances + d.Get("routing_instance").(string) +
				" routing-options aggregate route " + d.Get("destination").(string)
		} else {
			setPrefix = setRoutingInstances + d.Get("routing_instance").(string) +
				" routing-options rib " + d.Get("routing_instance").(string) + ".inet6.0 " +
				"aggregate route " + d.Get("destination").(string)
		}
	}
	configSet = append(configSet, setPrefix)
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
	if d.Get("brief").(bool) {
		configSet = append(configSet, setPrefix+" brief")
	}
	for _, v := range d.Get("community").([]interface{}) {
		configSet = append(configSet, setPrefix+" community "+v.(string))
	}
	if d.Get("discard").(bool) {
		configSet = append(configSet, setPrefix+" discard")
	}
	if d.Get("full").(bool) {
		configSet = append(configSet, setPrefix+" full")
	}
	if d.Get("metric").(int) > 0 {
		configSet = append(configSet, setPrefix+" metric "+strconv.Itoa(d.Get("metric").(int)))
	}
	if d.Get("passive").(bool) {
		configSet = append(configSet, setPrefix+" passive")
	}
	for _, v := range d.Get("policy").([]interface{}) {
		configSet = append(configSet, setPrefix+" policy "+v.(string))
	}
	if d.Get("preference").(int) > 0 {
		configSet = append(configSet, setPrefix+" preference "+strconv.Itoa(d.Get("preference").(int)))
	}

	return sess.configSet(configSet, junSess)
}

func readAggregateRoute(destination, instance string, sess *Session, junSess *junosSession,
) (aggregateRouteOptions, error) {
	var confRead aggregateRouteOptions
	var showConfig string
	var err error

	if instance == defaultW {
		if !strings.Contains(destination, ":") {
			showConfig, err = sess.command(cmdShowConfig+
				"routing-options aggregate route "+destination+pipeDisplaySetRelative, junSess)
		} else {
			showConfig, err = sess.command(cmdShowConfig+
				"routing-options rib inet6.0 aggregate route "+destination+pipeDisplaySetRelative, junSess)
		}
	} else {
		if !strings.Contains(destination, ":") {
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+instance+" "+
				"routing-options aggregate route "+destination+pipeDisplaySetRelative, junSess)
		} else {
			showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+instance+" "+
				"routing-options rib "+instance+".inet6.0 aggregate route "+destination+pipeDisplaySetRelative, junSess)
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

func delAggregateRoute(destination, instance string, sess *Session, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	if instance == defaultW {
		if !strings.Contains(destination, ":") {
			configSet = append(configSet, "delete routing-options aggregate route "+destination)
		} else {
			configSet = append(configSet, "delete routing-options rib inet6.0 aggregate route "+destination)
		}
	} else {
		if !strings.Contains(destination, ":") {
			configSet = append(configSet, delRoutingInstances+instance+" "+
				"routing-options aggregate route "+destination)
		} else {
			configSet = append(configSet, delRoutingInstances+instance+" "+
				"routing-options rib "+instance+".inet6.0 aggregate route "+destination)
		}
	}

	return sess.configSet(configSet, junSess)
}

func fillAggregateRouteData(d *schema.ResourceData, aggregateRouteOptions aggregateRouteOptions) {
	if tfErr := d.Set("destination", aggregateRouteOptions.destination); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", aggregateRouteOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("active", aggregateRouteOptions.active); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_aggregator_address", aggregateRouteOptions.asPathAggregatorAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_aggregator_as_number", aggregateRouteOptions.asPathAggregatorAsNumber); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_atomic_aggregate", aggregateRouteOptions.asPathAtomicAggregate); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_origin", aggregateRouteOptions.asPathOrigin); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path_path", aggregateRouteOptions.asPathPath); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("brief", aggregateRouteOptions.brief); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("community", aggregateRouteOptions.community); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("discard", aggregateRouteOptions.discard); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("full", aggregateRouteOptions.full); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric", aggregateRouteOptions.metric); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("passive", aggregateRouteOptions.passive); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("policy", aggregateRouteOptions.policy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("preference", aggregateRouteOptions.preference); tfErr != nil {
		panic(tfErr)
	}
}
