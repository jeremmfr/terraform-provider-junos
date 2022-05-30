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
)

type igmpSnoopingVlanOptions struct {
	immediateLeave          bool
	proxy                   bool
	queryInterval           int
	robustCount             int
	l2QuerierSrcAddress     string
	name                    string
	proxySrcAddress         string
	queryLastMemberInterval string
	queryResponseInterval   string
	routingInstance         string
	interFace               []map[string]interface{}
}

func resourceIgmpSnoopingVlan() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceIgmpSnoopingVlanCreate,
		ReadWithoutTimeout:   resourceIgmpSnoopingVlanRead,
		UpdateWithoutTimeout: resourceIgmpSnoopingVlanUpdate,
		DeleteWithoutTimeout: resourceIgmpSnoopingVlanDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIgmpSnoopingVlanImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          defaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"immediate_leave": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"interface": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if strings.Count(value, ".") != 1 {
									errors = append(errors, fmt.Errorf(
										"%q in %q need to have 1 dot", value, k))
								}

								return
							},
						},
						"group_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"host_only_interface": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"immediate_leave": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"multicast_router_interface": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"static_group": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
									"source": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "",
										ValidateFunc: validation.IsIPv4Address,
									},
								},
							},
						},
					},
				},
			},
			"l2_querier_source_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"proxy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"proxy_source_address": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"proxy"},
				ValidateFunc: validation.IsIPv4Address,
			},
			"query_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 1024),
			},
			"query_last_member_interval": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^\d+(\.\d+)?$`), "must be a number with optional decimal"),
			},
			"query_response_interval": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^\d+(\.\d+)?$`), "must be a number with optional decimal"),
			},
			"robust_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 10),
			},
		},
	}
}

func resourceIgmpSnoopingVlanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setIgmpSnoopingVlan(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string) + idSeparator + d.Get("routing_instance").(string))

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if d.Get("routing_instance").(string) != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	igmpSnoopingVlanExists, err := checkIgmpSnoopingVlanExists(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if igmpSnoopingVlanExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))
		if d.Get("routing_instance").(string) == defaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf("protocols igmp-snooping vlan %v already exists",
				d.Get("name").(string)))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"protocols igmp-snooping vlan %v already exists in routing-instance %v",
			d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}

	if err := setIgmpSnoopingVlan(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_igmp_snooping_vlan", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	igmpSnoopingVlanExists, err = checkIgmpSnoopingVlanExists(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if igmpSnoopingVlanExists {
		d.SetId(d.Get("name").(string) + idSeparator + d.Get("routing_instance").(string))
	} else {
		if d.Get("routing_instance").(string) == defaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf("protocols igmp-snooping vlan %v not exists after commit "+
				"=> check your config", d.Get("name").(string)))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"protocols igmp-snooping vlan %v not exists in routing-instance %v after commit "+
				"=> check your config", d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceIgmpSnoopingVlanReadWJunSess(d, clt, junSess)...)
}

func resourceIgmpSnoopingVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceIgmpSnoopingVlanReadWJunSess(d, clt, junSess)
}

func resourceIgmpSnoopingVlanReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	igmpSnoopingVlanOptions, err := readIgmpSnoopingVlan(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if igmpSnoopingVlanOptions.name == "" {
		d.SetId("")
	} else {
		fillIgmpSnoopingVlanData(d, igmpSnoopingVlanOptions)
	}

	return nil
}

func resourceIgmpSnoopingVlanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delIgmpSnoopingVlan(d.Get("name").(string), d.Get("routing_instance").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setIgmpSnoopingVlan(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delIgmpSnoopingVlan(d.Get("name").(string), d.Get("routing_instance").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setIgmpSnoopingVlan(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_igmp_snooping_vlan", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIgmpSnoopingVlanReadWJunSess(d, clt, junSess)...)
}

func resourceIgmpSnoopingVlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delIgmpSnoopingVlan(d.Get("name").(string), d.Get("routing_instance").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delIgmpSnoopingVlan(d.Get("name").(string), d.Get("routing_instance").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_igmp_snooping_vlan", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceIgmpSnoopingVlanImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	igmpSnoopingVlanExists, err := checkIgmpSnoopingVlanExists(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	if !igmpSnoopingVlanExists {
		return nil, fmt.Errorf("don't find protocols igmp-snooping vlan with id '%v' "+
			"(id must be <name>%s<routing_instance>)", d.Id(), idSeparator)
	}
	igmpSnoopingVlanOptions, err := readIgmpSnoopingVlan(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	fillIgmpSnoopingVlanData(d, igmpSnoopingVlanOptions)

	result[0] = d

	return result, nil
}

func checkIgmpSnoopingVlanExists(name, routingInstance string, clt *Client, junSess *junosSession) (bool, error) {
	var showConfig string
	var err error
	if routingInstance == defaultW {
		showConfig, err = clt.command(cmdShowConfig+
			"protocols igmp-snooping vlan "+name+pipeDisplaySet, junSess)
	} else {
		showConfig, err = clt.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols igmp-snooping vlan "+name+pipeDisplaySet, junSess)
	}
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setIgmpSnoopingVlan(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := setLS
	if rI := d.Get("routing_instance").(string); rI != defaultW {
		setPrefix = setRoutingInstances + rI + " "
	}
	setPrefix += "protocols igmp-snooping vlan " + d.Get("name").(string) + " "

	configSet = append(configSet, setPrefix)
	if d.Get("immediate_leave").(bool) {
		configSet = append(configSet, setPrefix+"immediate-leave")
	}
	interfaceList := make([]string, 0)
	for _, mIntface := range d.Get("interface").([]interface{}) {
		intFace := mIntface.(map[string]interface{})
		if bchk.StringInSlice(intFace["name"].(string), interfaceList) {
			return fmt.Errorf("multiple blocks interface with the same name '%s'", intFace["name"].(string))
		}
		interfaceList = append(interfaceList, intFace["name"].(string))
		setPrefixIntface := setPrefix + "interface " + intFace["name"].(string) + " "
		configSet = append(configSet, setPrefixIntface)
		if v := intFace["group_limit"].(int); v != -1 {
			configSet = append(configSet, setPrefixIntface+"group-limit "+strconv.Itoa(v))
		}
		if intFace["host_only_interface"].(bool) {
			configSet = append(configSet, setPrefixIntface+"host-only-interface")
		}
		if intFace["immediate_leave"].(bool) {
			configSet = append(configSet, setPrefixIntface+"immediate-leave")
		}
		if intFace["multicast_router_interface"].(bool) {
			configSet = append(configSet, setPrefixIntface+"multicast-router-interface")
		}
		staticGroupList := make([]string, 0)
		for _, mStaticGrp := range intFace["static_group"].(*schema.Set).List() {
			staticGroup := mStaticGrp.(map[string]interface{})
			if bchk.StringInSlice(staticGroup["address"].(string), staticGroupList) {
				return fmt.Errorf("multiple blocks static_group with the same address '%s'", staticGroup["address"].(string))
			}
			staticGroupList = append(staticGroupList, staticGroup["address"].(string))
			configSet = append(configSet, setPrefixIntface+"static group "+staticGroup["address"].(string))
			if v := staticGroup["source"].(string); v != "" {
				configSet = append(configSet, setPrefixIntface+"static group "+staticGroup["address"].(string)+" source "+v)
			}
		}
	}
	if v := d.Get("l2_querier_source_address").(string); v != "" {
		configSet = append(configSet, setPrefix+"l2-querier source-address "+v)
	}
	if d.Get("proxy").(bool) {
		configSet = append(configSet, setPrefix+"proxy")
		if v := d.Get("proxy_source_address").(string); v != "" {
			configSet = append(configSet, setPrefix+"proxy source-address "+v)
		}
	} else if d.Get("proxy_source_address").(string) != "" {
		return fmt.Errorf("proxy need to be true with proxy_source_address")
	}
	if v := d.Get("query_interval").(int); v != 0 {
		configSet = append(configSet, setPrefix+"query-interval "+strconv.Itoa(v))
	}
	if v := d.Get("query_last_member_interval").(string); v != "" {
		configSet = append(configSet, setPrefix+"query-last-member-interval "+v)
	}
	if v := d.Get("query_response_interval").(string); v != "" {
		configSet = append(configSet, setPrefix+"query-response-interval "+v)
	}
	if v := d.Get("robust_count").(int); v != 0 {
		configSet = append(configSet, setPrefix+"robust-count "+strconv.Itoa(v))
	}

	return clt.configSet(configSet, junSess)
}

func readIgmpSnoopingVlan(name, routingInstance string, clt *Client, junSess *junosSession,
) (igmpSnoopingVlanOptions, error) {
	var confRead igmpSnoopingVlanOptions
	var showConfig string
	var err error
	if routingInstance == defaultW {
		showConfig, err = clt.command(cmdShowConfig+
			"protocols igmp-snooping vlan "+name+pipeDisplaySetRelative, junSess)
	} else {
		showConfig, err = clt.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols igmp-snooping vlan "+name+pipeDisplaySetRelative, junSess)
	}
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
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
			case itemTrim == "immediate-leave":
				confRead.immediateLeave = true
			case strings.HasPrefix(itemTrim, "interface "):
				itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "interface "), " ")
				intFace := map[string]interface{}{
					"name":                       itemTrimSplit[0],
					"group_limit":                -1,
					"host_only_interface":        false,
					"immediate_leave":            false,
					"multicast_router_interface": false,
					"static_group":               make([]map[string]interface{}, 0),
				}
				confRead.interFace = copyAndRemoveItemMapList("name", intFace, confRead.interFace)
				itemTrimIntface := strings.TrimPrefix(itemTrim, "interface "+itemTrimSplit[0]+" ")
				switch {
				case strings.HasPrefix(itemTrimIntface, "group-limit "):
					var err error
					intFace["group_limit"], err = strconv.Atoi(strings.TrimPrefix(itemTrimIntface, "group-limit "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case itemTrimIntface == "host-only-interface":
					intFace["host_only_interface"] = true
				case itemTrimIntface == "immediate-leave":
					intFace["immediate_leave"] = true
				case itemTrimIntface == "multicast-router-interface":
					intFace["multicast_router_interface"] = true
				case strings.HasPrefix(itemTrimIntface, "static group "):
					itemTrimIntfaceSplit := strings.Split(strings.TrimPrefix(itemTrimIntface, "static group "), " ")
					staticGrp := map[string]interface{}{
						"address": itemTrimIntfaceSplit[0],
						"source":  "",
					}
					intFace["static_group"] = copyAndRemoveItemMapList("address", staticGrp,
						intFace["static_group"].([]map[string]interface{}))
					if strings.HasPrefix(itemTrimIntface, "static group "+itemTrimIntfaceSplit[0]+" source ") {
						staticGrp["source"] = strings.TrimPrefix(itemTrimIntface, "static group "+itemTrimIntfaceSplit[0]+" source ")
					}
					intFace["static_group"] = append(intFace["static_group"].([]map[string]interface{}), staticGrp)
				}
				confRead.interFace = append(confRead.interFace, intFace)
			case strings.HasPrefix(itemTrim, "l2-querier source-address "):
				confRead.l2QuerierSrcAddress = strings.TrimPrefix(itemTrim, "l2-querier source-address ")
			case strings.HasPrefix(itemTrim, "proxy"):
				confRead.proxy = true
				if strings.HasPrefix(itemTrim, "proxy source-address ") {
					confRead.proxySrcAddress = strings.TrimPrefix(itemTrim, "proxy source-address ")
				}
			case strings.HasPrefix(itemTrim, "query-interval "):
				var err error
				confRead.queryInterval, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "query-interval "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "query-last-member-interval "):
				confRead.queryLastMemberInterval = strings.TrimPrefix(itemTrim, "query-last-member-interval ")
			case strings.HasPrefix(itemTrim, "query-response-interval "):
				confRead.queryResponseInterval = strings.TrimPrefix(itemTrim, "query-response-interval ")
			case strings.HasPrefix(itemTrim, "robust-count "):
				var err error
				confRead.robustCount, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "robust-count "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delIgmpSnoopingVlan(name, routingInstance string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)

	if routingInstance == defaultW {
		configSet = append(configSet, "delete protocols igmp-snooping vlan "+name)
	} else {
		configSet = append(configSet, delRoutingInstances+routingInstance+" protocols igmp-snooping vlan "+name)
	}

	return clt.configSet(configSet, junSess)
}

func fillIgmpSnoopingVlanData(d *schema.ResourceData, igmpSnoopingVlanOptions igmpSnoopingVlanOptions) {
	if tfErr := d.Set("name", igmpSnoopingVlanOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", igmpSnoopingVlanOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("immediate_leave", igmpSnoopingVlanOptions.immediateLeave); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interface", igmpSnoopingVlanOptions.interFace); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("l2_querier_source_address", igmpSnoopingVlanOptions.l2QuerierSrcAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("proxy", igmpSnoopingVlanOptions.proxy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("proxy_source_address", igmpSnoopingVlanOptions.proxySrcAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("query_interval", igmpSnoopingVlanOptions.queryInterval); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("query_last_member_interval", igmpSnoopingVlanOptions.queryLastMemberInterval); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("query_response_interval", igmpSnoopingVlanOptions.queryResponseInterval); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("robust_count", igmpSnoopingVlanOptions.robustCount); tfErr != nil {
		panic(tfErr)
	}
}
