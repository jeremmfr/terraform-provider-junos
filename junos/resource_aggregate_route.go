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
	active          bool
	brief           bool
	discard         bool
	full            bool
	passive         bool
	metric          int
	preference      int
	destination     string
	routingInstance string
	community       []string
	policy          []string
}

func resourceAggregateRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAggregateRouteCreate,
		ReadContext:   resourceAggregateRouteRead,
		UpdateContext: resourceAggregateRouteUpdate,
		DeleteContext: resourceAggregateRouteDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAggregateRouteImport,
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
				Default:          defaultWord,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"active": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"passive"},
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
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if d.Get("routing_instance").(string) != defaultWord {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)

			return diag.FromErr(err)
		}
		if !instanceExists {
			sess.configClear(jnprSess)

			return diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))
		}
	}
	aggregateRouteExists, err := checkAggregateRouteExists(
		d.Get("destination").(string), d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if aggregateRouteExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("aggregate route %v already exists on table %s",
			d.Get("destination").(string), d.Get("routing_instance").(string)))
	}
	if err := setAggregateRoute(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_aggregate_route", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	aggregateRouteExists, err = checkAggregateRouteExists(
		d.Get("destination").(string), d.Get("routing_instance").(string), m, jnprSess)
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

	return append(diagWarns, resourceAggregateRouteReadWJnprSess(d, m, jnprSess)...)
}
func resourceAggregateRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceAggregateRouteReadWJnprSess(d, m, jnprSess)
}
func resourceAggregateRouteReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	aggregateRouteOptions, err := readAggregateRoute(d.Get("destination").(string), d.Get("routing_instance").(string),
		m, jnprSess)
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
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delAggregateRouteOpts(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}

	if err := setAggregateRoute(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_aggregate_route", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceAggregateRouteReadWJnprSess(d, m, jnprSess)...)
}
func resourceAggregateRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delAggregateRoute(d.Get("destination").(string), d.Get("routing_instance").(string),
		m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_aggregate_route", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceAggregateRouteImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	aggregateRouteExists, err := checkAggregateRouteExists(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !aggregateRouteExists {
		return nil, fmt.Errorf("don't find aggregate route with id '%v' (id must be "+
			"<destination>"+idSeparator+"<routing_instance>)", d.Id())
	}
	aggregateRouteOptions, err := readAggregateRoute(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillAggregateRouteData(d, aggregateRouteOptions)

	result[0] = d

	return result, nil
}

func checkAggregateRouteExists(destination string, instance string, m interface{},
	jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	var aggregateRouteConfig string
	var err error
	if instance == defaultWord {
		aggregateRouteConfig, err = sess.command("show configuration"+
			" routing-options aggregate route "+destination+" | display set", jnprSess)
		if err != nil {
			return false, err
		}
	} else {
		aggregateRouteConfig, err = sess.command("show configuration routing-instances "+instance+
			" routing-options aggregate route "+destination+" | display set", jnprSess)
		if err != nil {
			return false, err
		}
	}

	if aggregateRouteConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setAggregateRoute(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	var setPrefix string
	if d.Get("routing_instance").(string) == defaultWord {
		setPrefix = "set routing-options aggregate route " + d.Get("destination").(string)
	} else {
		setPrefix = "set routing-instances " + d.Get("routing_instance").(string) +
			" routing-options aggregate route " + d.Get("destination").(string)
	}
	configSet = append(configSet, setPrefix)
	if d.Get("active").(bool) {
		configSet = append(configSet, setPrefix+" active")
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
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readAggregateRoute(destination string, instance string, m interface{},
	jnprSess *NetconfObject) (aggregateRouteOptions, error) {
	sess := m.(*Session)
	var confRead aggregateRouteOptions
	var destinationConfig string
	var err error

	if instance == defaultWord {
		destinationConfig, err = sess.command("show configuration"+
			" routing-options aggregate route "+destination+" | display set relative", jnprSess)
	} else {
		destinationConfig, err = sess.command("show configuration routing-instances "+instance+
			" routing-options aggregate route "+destination+" | display set relative", jnprSess)
	}
	if err != nil {
		return confRead, err
	}

	if destinationConfig != emptyWord {
		confRead.destination = destination
		confRead.routingInstance = instance
		for _, item := range strings.Split(destinationConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case itemTrim == "active":
				confRead.active = true
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
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case itemTrim == passiveW:
				confRead.passive = true
			case strings.HasPrefix(itemTrim, "policy "):
				confRead.policy = append(confRead.policy, strings.TrimPrefix(itemTrim, "policy "))
			case strings.HasPrefix(itemTrim, "preference "):
				confRead.preference, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "preference "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delAggregateRouteOpts(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := "delete "
	if d.Get("routing_instance").(string) == defaultWord {
		delPrefix += "routing-options aggregate route "
	} else {
		delPrefix += "routing-instances " + d.Get("routing_instance").(string) + " routing-options aggregate route "
	}
	delPrefix += d.Get("destination").(string) + " "
	configSet = append(configSet,
		delPrefix+"active",
		delPrefix+"brief",
		delPrefix+"community",
		delPrefix+"discard",
		delPrefix+"full",
		delPrefix+"metric",
		delPrefix+"passive",
		delPrefix+"policy",
		delPrefix+"preference",
	)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func delAggregateRoute(destination string, instance string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	if instance == defaultWord {
		configSet = append(configSet, "delete routing-options aggregate route "+destination)
	} else {
		configSet = append(configSet, "delete routing-instances "+instance+" routing-options aggregate route "+destination)
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
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
