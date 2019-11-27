package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type staticRouteOptions struct {
	preference       int
	metric           int
	destination      string
	routingInstance  string
	nextHop          []string
	qualifiedNextHop []map[string]interface{}
}

func resourceStaticRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceStaticRouteCreate,
		Read:   resourceStaticRouteRead,
		Update: resourceStaticRouteUpdate,
		Delete: resourceStaticRouteDelete,
		Importer: &schema.ResourceImporter{
			State: resourceStaticRouteImport,
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
			"preference": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"metric": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"next_hop": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"qualified_next_hop": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"next_hop": {
							Type:     schema.TypeString,
							Required: true,
						},
						"preference": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"metric": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceStaticRouteCreate(d *schema.ResourceData, m interface{}) error {
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
	staticRouteExists, err := checkStaticRouteExists(d.Get("destination").(string), d.Get("routing_instance").(string),
		m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if staticRouteExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("static route %v already exists on table %s",
			d.Get("destination").(string), d.Get("routing_instance").(string))
	}
	err = setStaticRoute(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("create resource junos_static_route", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	staticRouteExists, err = checkStaticRouteExists(d.Get("destination").(string), d.Get("routing_instance").(string),
		m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if staticRouteExists {
		d.SetId(d.Get("destination").(string) + idSeparator + d.Get("routing_instance").(string))
	} else {
		return fmt.Errorf("static route %v not exists in routing_instance %v after commit "+
			"=> check your config", d.Get("destination").(string), d.Get("routing_instance").(string))
	}
	return resourceStaticRouteRead(d, m)
}
func resourceStaticRouteRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	staticRouteOptions, err := readStaticRoute(d.Get("destination").(string), d.Get("routing_instance").(string),
		m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if staticRouteOptions.destination == "" {
		d.SetId("")
	} else {
		fillStaticRouteData(d, staticRouteOptions)
	}
	return nil
}
func resourceStaticRouteUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delStaticRouteOpts(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}

	err = setStaticRoute(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("update resource junos_static_route", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourceStaticRouteRead(d, m)
}
func resourceStaticRouteDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delStaticRoute(d.Get("destination").(string), d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("delete resource junos_static_route", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
}
func resourceStaticRouteImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
	staticRouteExists, err := checkStaticRouteExists(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !staticRouteExists {
		return nil, fmt.Errorf("don't find static route with id '%v' (id must be "+
			"<destination>"+idSeparator+"<routing_instance>)", d.Id())
	}
	staticRouteOptions, err := readStaticRoute(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillStaticRouteData(d, staticRouteOptions)

	result[0] = d
	return result, nil
}

func checkStaticRouteExists(destination string, instance string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	var staticRouteConfig string
	var err error
	if instance == defaultWord {
		staticRouteConfig, err = sess.command("show configuration"+
			" routing-options static route "+destination+" | display set", jnprSess)
		if err != nil {
			return false, err
		}
	} else {
		staticRouteConfig, err = sess.command("show configuration routing-instances "+instance+
			" routing-options static route "+destination+" | display set", jnprSess)
		if err != nil {
			return false, err
		}
	}

	if staticRouteConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setStaticRoute(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	var setPrefix string
	if d.Get("routing_instance").(string) == defaultWord {
		setPrefix = "set routing-options static route " + d.Get("destination").(string)
	} else {
		setPrefix = "set routing-instances " + d.Get("routing_instance").(string) +
			" routing-options static route " + d.Get("destination").(string)
	}
	if d.Get("preference").(int) > 0 {
		configSet = append(configSet, setPrefix+" preference "+strconv.Itoa(d.Get("preference").(int))+"\n")
	}
	if d.Get("metric").(int) > 0 {
		configSet = append(configSet, setPrefix+" metric "+strconv.Itoa(d.Get("metric").(int))+"\n")
	}
	for _, nextHop := range d.Get("next_hop").([]interface{}) {
		configSet = append(configSet, setPrefix+" next-hop "+nextHop.(string)+"\n")
	}
	for _, qualifiedNextHop := range d.Get("qualified_next_hop").([]interface{}) {
		qualifiedNextHopMap := qualifiedNextHop.(map[string]interface{})
		configSet = append(configSet, setPrefix+" qualified-next-hop "+qualifiedNextHopMap["next_hop"].(string)+"\n")
		if qualifiedNextHopMap["preference"].(int) > 0 {
			configSet = append(configSet, setPrefix+
				" qualified-next-hop "+qualifiedNextHopMap["next_hop"].(string)+
				" preference "+strconv.Itoa(qualifiedNextHopMap["preference"].(int))+"\n")
		}
		if qualifiedNextHopMap["metric"].(int) > 0 {
			configSet = append(configSet, setPrefix+
				" qualified-next-hop "+qualifiedNextHopMap["next_hop"].(string)+
				" metric "+strconv.Itoa(qualifiedNextHopMap["metric"].(int))+"\n")
		}
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readStaticRoute(destination string, instance string, m interface{},
	jnprSess *NetconfObject) (staticRouteOptions, error) {
	sess := m.(*Session)
	var confRead staticRouteOptions
	var destinationConfig string
	var err error

	if instance == defaultWord {
		destinationConfig, err = sess.command("show configuration"+
			" routing-options static route "+destination+" | display set relative", jnprSess)
	} else {
		destinationConfig, err = sess.command("show configuration routing-instances "+instance+
			" routing-options static route "+destination+" | display set relative", jnprSess)
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
			case strings.HasPrefix(itemTrim, "next-hop "):
				confRead.nextHop = append(confRead.nextHop, strings.TrimPrefix(itemTrim, "next-hop "))
			case strings.HasPrefix(itemTrim, "qualified-next-hop "):
				nextHop := strings.TrimPrefix(itemTrim, "qualified-next-hop ")
				nextHopWords := strings.Split(nextHop, " ")
				qualifiedNextHopOptions := map[string]interface{}{
					"next_hop":   nextHopWords[0],
					"metric":     0,
					"preference": 0,
				}
				qualifiedNextHopOptions, confRead.qualifiedNextHop = copyAndRemoveItemMapList("next_hop",
					false, qualifiedNextHopOptions, confRead.qualifiedNextHop)
				itemTrimQnh := strings.TrimPrefix(itemTrim, "qualified-next-hop "+nextHopWords[0]+" ")
				switch {
				case strings.HasPrefix(itemTrimQnh, "metric "):
					qualifiedNextHopOptions["metric"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrimQnh, "metric "))
					if err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrimQnh, "preference "):
					qualifiedNextHopOptions["preference"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrimQnh, "preference "))
					if err != nil {
						return confRead, err
					}
				}
				confRead.qualifiedNextHop = append(confRead.qualifiedNextHop, qualifiedNextHopOptions)
			}
		}
	} else {
		confRead.destination = ""
		return confRead, nil
	}
	return confRead, nil
}

func delStaticRouteOpts(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := "delete "
	if d.Get("routing_instance").(string) == defaultWord {
		delPrefix += "routing-options static route "
	} else {
		delPrefix += "routing-instances " + d.Get("routing_instance").(string) + " routing-options static route "
	}
	delPrefix += d.Get("destination").(string) + " "
	configSet = append(configSet,
		delPrefix+"preference\n",
		delPrefix+"metric\n",
		delPrefix+"next-hop\n")
	if d.HasChange("qualified_next_hop") {
		oQualifiedNextHop, _ := d.GetChange("qualified_next_hop")
		for _, v := range oQualifiedNextHop.([]interface{}) {
			qualifiedNextHop := v.(map[string]interface{})
			configSet = append(configSet, delPrefix+"qualified-next-hop "+qualifiedNextHop["next_hop"].(string)+"\n")
		}
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func delStaticRoute(destination string, instance string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	if instance == defaultWord {
		configSet = append(configSet, "del routing-options static route "+destination)
	} else {
		configSet = append(configSet, "del routing-instances "+instance+" routing-options static route "+destination)
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}

func fillStaticRouteData(d *schema.ResourceData, staticRouteOptions staticRouteOptions) {
	tfErr := d.Set("destination", staticRouteOptions.destination)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("routing_instance", staticRouteOptions.routingInstance)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("preference", staticRouteOptions.preference)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric", staticRouteOptions.metric)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("next_hop", staticRouteOptions.nextHop)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("qualified_next_hop", staticRouteOptions.qualifiedNextHop)
	if tfErr != nil {
		panic(tfErr)
	}
}
