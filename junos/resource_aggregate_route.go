package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type aggregateRouteOptions struct {
	active          bool
	passive         bool
	brief           bool
	full            bool
	discard         bool
	preference      int
	metric          int
	destination     string
	routingInstance string
	policy          []string
	community       []string
}

func resourceAggregateRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceAggregateRouteCreate,
		Read:   resourceAggregateRouteRead,
		Update: resourceAggregateRouteUpdate,
		Delete: resourceAggregateRouteDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAggregateRouteImport,
		},
		Schema: map[string]*schema.Schema{
			"destination": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					err := validateNetwork(value)
					if err != nil {
						errors = append(errors, fmt.Errorf(
							"%q error for validate %q : %q", k, value, err))
					}

					return
				},
			},
			"routing_instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      defaultWord,
				ValidateFunc: validateNameObjectJunos(),
			},
			"active": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"passive"},
			},
			"passive": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"active"},
			},
			"brief": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"full"},
			},
			"full": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"brief"},
			},
			"discard": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"preference": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"metric": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"community": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAggregateRouteCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	if d.Get("routing_instance").(string) != defaultWord {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)

			return err
		}
		if !instanceExists {
			sess.configClear(jnprSess)

			return fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string))
		}
	}
	aggregateRouteExists, err := checkAggregateRouteExists(
		d.Get("destination").(string), d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	if aggregateRouteExists {
		sess.configClear(jnprSess)

		return fmt.Errorf("aggregate route %v already exists on table %s",
			d.Get("destination").(string), d.Get("routing_instance").(string))
	}
	err = setAggregateRoute(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("create resource junos_aggregate_route", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	aggregateRouteExists, err = checkAggregateRouteExists(
		d.Get("destination").(string), d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	if aggregateRouteExists {
		d.SetId(d.Get("destination").(string) + idSeparator + d.Get("routing_instance").(string))
	} else {
		return fmt.Errorf("aggregate route %v not exists in routing_instance %v after commit "+
			"=> check your config", d.Get("destination").(string), d.Get("routing_instance").(string))
	}

	return resourceAggregateRouteRead(d, m)
}
func resourceAggregateRouteRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return err
	}
	defer sess.closeSession(jnprSess)
	aggregateRouteOptions, err := readAggregateRoute(d.Get("destination").(string), d.Get("routing_instance").(string),
		m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if aggregateRouteOptions.destination == "" {
		d.SetId("")
	} else {
		fillAggregateRouteData(d, aggregateRouteOptions)
	}

	return nil
}
func resourceAggregateRouteUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delAggregateRouteOpts(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}

	err = setAggregateRoute(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("update resource junos_aggregate_route", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	d.Partial(false)

	return resourceAggregateRouteRead(d, m)
}
func resourceAggregateRouteDelete(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delAggregateRoute(d.Get("destination").(string), d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("delete resource junos_aggregate_route", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}

	return nil
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
	if d.Get("passive").(bool) {
		configSet = append(configSet, setPrefix+" passive")
	}
	if d.Get("brief").(bool) {
		configSet = append(configSet, setPrefix+" brief")
	}
	if d.Get("full").(bool) {
		configSet = append(configSet, setPrefix+" full")
	}
	if d.Get("discard").(bool) {
		configSet = append(configSet, setPrefix+" discard")
	}
	if d.Get("preference").(int) > 0 {
		configSet = append(configSet, setPrefix+" preference "+strconv.Itoa(d.Get("preference").(int)))
	}
	if d.Get("metric").(int) > 0 {
		configSet = append(configSet, setPrefix+" metric "+strconv.Itoa(d.Get("metric").(int)))
	}
	if len(d.Get("community").([]interface{})) > 0 {
		for _, v := range d.Get("community").([]interface{}) {
			configSet = append(configSet, setPrefix+" community "+v.(string))
		}
	}
	if len(d.Get("policy").([]interface{})) > 0 {
		for _, v := range d.Get("policy").([]interface{}) {
			configSet = append(configSet, setPrefix+" policy "+v.(string))
		}
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
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
			case strings.HasSuffix(itemTrim, "active"):
				confRead.active = true
			case strings.HasSuffix(itemTrim, "passive"):
				confRead.passive = true
			case strings.HasSuffix(itemTrim, "brief"):
				confRead.brief = true
			case strings.HasSuffix(itemTrim, "full"):
				confRead.full = true
			case strings.HasSuffix(itemTrim, "discard"):
				confRead.discard = true
			case strings.HasPrefix(itemTrim, "preference "):
				confRead.preference, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "preference "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "metric "):
				confRead.metric, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "metric "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "community "):
				confRead.community = append(confRead.community, strings.TrimPrefix(itemTrim, "community "))
			case strings.HasPrefix(itemTrim, "policy "):
				confRead.policy = append(confRead.policy, strings.TrimPrefix(itemTrim, "policy "))
			}
		}
	} else {
		confRead.destination = ""

		return confRead, nil
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
		delPrefix+"passive",
		delPrefix+"brief",
		delPrefix+"full",
		delPrefix+"discard",
		delPrefix+"preference",
		delPrefix+"metric",
		delPrefix+"community",
		delPrefix+"policy",
	)
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
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
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}

func fillAggregateRouteData(d *schema.ResourceData, aggregateRouteOptions aggregateRouteOptions) {
	tfErr := d.Set("destination", aggregateRouteOptions.destination)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("routing_instance", aggregateRouteOptions.routingInstance)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("active", aggregateRouteOptions.active)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("passive", aggregateRouteOptions.passive)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("brief", aggregateRouteOptions.brief)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("full", aggregateRouteOptions.full)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("discard", aggregateRouteOptions.discard)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("preference", aggregateRouteOptions.preference)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric", aggregateRouteOptions.metric)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("community", aggregateRouteOptions.community)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("policy", aggregateRouteOptions.policy)
	if tfErr != nil {
		panic(tfErr)
	}
}
