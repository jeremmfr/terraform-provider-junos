package providersdk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
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
				Default:          junos.DefaultW,
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setIgmpSnoopingVlan(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string) + junos.IDSeparator + d.Get("routing_instance").(string))

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
	igmpSnoopingVlanExists, err := checkIgmpSnoopingVlanExists(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if igmpSnoopingVlanExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())
		if d.Get("routing_instance").(string) == junos.DefaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf("protocols igmp-snooping vlan %v already exists",
				d.Get("name").(string)))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"protocols igmp-snooping vlan %v already exists in routing-instance %v",
			d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}

	if err := setIgmpSnoopingVlan(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_igmp_snooping_vlan")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	igmpSnoopingVlanExists, err = checkIgmpSnoopingVlanExists(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if igmpSnoopingVlanExists {
		d.SetId(d.Get("name").(string) + junos.IDSeparator + d.Get("routing_instance").(string))
	} else {
		if d.Get("routing_instance").(string) == junos.DefaultW {
			return append(diagWarns, diag.FromErr(fmt.Errorf("protocols igmp-snooping vlan %v not exists after commit "+
				"=> check your config", d.Get("name").(string)))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"protocols igmp-snooping vlan %v not exists in routing-instance %v after commit "+
				"=> check your config", d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceIgmpSnoopingVlanReadWJunSess(d, junSess)...)
}

func resourceIgmpSnoopingVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceIgmpSnoopingVlanReadWJunSess(d, junSess)
}

func resourceIgmpSnoopingVlanReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	igmpSnoopingVlanOptions, err := readIgmpSnoopingVlan(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		junSess,
	)
	junos.MutexUnlock()
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
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delIgmpSnoopingVlan(d.Get("name").(string), d.Get("routing_instance").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setIgmpSnoopingVlan(d, junSess); err != nil {
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
	if err := delIgmpSnoopingVlan(d.Get("name").(string), d.Get("routing_instance").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setIgmpSnoopingVlan(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_igmp_snooping_vlan")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIgmpSnoopingVlanReadWJunSess(d, junSess)...)
}

func resourceIgmpSnoopingVlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delIgmpSnoopingVlan(d.Get("name").(string), d.Get("routing_instance").(string), junSess); err != nil {
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
	if err := delIgmpSnoopingVlan(d.Get("name").(string), d.Get("routing_instance").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_igmp_snooping_vlan")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceIgmpSnoopingVlanImport(ctx context.Context, d *schema.ResourceData, m interface{},
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
	igmpSnoopingVlanExists, err := checkIgmpSnoopingVlanExists(idSplit[0], idSplit[1], junSess)
	if err != nil {
		return nil, err
	}
	if !igmpSnoopingVlanExists {
		return nil, fmt.Errorf("don't find protocols igmp-snooping vlan with id '%v' "+
			"(id must be <name>"+junos.IDSeparator+"<routing_instance>)", d.Id())
	}
	igmpSnoopingVlanOptions, err := readIgmpSnoopingVlan(idSplit[0], idSplit[1], junSess)
	if err != nil {
		return nil, err
	}
	fillIgmpSnoopingVlanData(d, igmpSnoopingVlanOptions)

	result[0] = d

	return result, nil
}

func checkIgmpSnoopingVlanExists(name, routingInstance string, junSess *junos.Session,
) (_ bool, err error) {
	var showConfig string
	if routingInstance == junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"protocols igmp-snooping vlan " + name + junos.PipeDisplaySet)
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			"protocols igmp-snooping vlan " + name + junos.PipeDisplaySet)
	}
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setIgmpSnoopingVlan(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := junos.SetLS
	if rI := d.Get("routing_instance").(string); rI != junos.DefaultW {
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
		if slices.Contains(interfaceList, intFace["name"].(string)) {
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
			if slices.Contains(staticGroupList, staticGroup["address"].(string)) {
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
		return errors.New("proxy need to be true with proxy_source_address")
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

	return junSess.ConfigSet(configSet)
}

func readIgmpSnoopingVlan(name, routingInstance string, junSess *junos.Session,
) (confRead igmpSnoopingVlanOptions, err error) {
	var showConfig string
	if routingInstance == junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"protocols igmp-snooping vlan " + name + junos.PipeDisplaySetRelative)
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			"protocols igmp-snooping vlan " + name + junos.PipeDisplaySetRelative)
	}
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
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
			case itemTrim == "immediate-leave":
				confRead.immediateLeave = true
			case balt.CutPrefixInString(&itemTrim, "interface "):
				itemTrimFields := strings.Split(itemTrim, " ")
				intFace := map[string]interface{}{
					"name":                       itemTrimFields[0],
					"group_limit":                -1,
					"host_only_interface":        false,
					"immediate_leave":            false,
					"multicast_router_interface": false,
					"static_group":               make([]map[string]interface{}, 0),
				}
				confRead.interFace = copyAndRemoveItemMapList("name", intFace, confRead.interFace)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "group-limit "):
					intFace["group_limit"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case itemTrim == "host-only-interface":
					intFace["host_only_interface"] = true
				case itemTrim == "immediate-leave":
					intFace["immediate_leave"] = true
				case itemTrim == "multicast-router-interface":
					intFace["multicast_router_interface"] = true
				case balt.CutPrefixInString(&itemTrim, "static group "):
					itemTrimStaticGrpFields := strings.Split(itemTrim, " ")
					staticGrp := map[string]interface{}{
						"address": itemTrimStaticGrpFields[0],
						"source":  "",
					}
					intFace["static_group"] = copyAndRemoveItemMapList("address", staticGrp,
						intFace["static_group"].([]map[string]interface{}))
					if balt.CutPrefixInString(&itemTrim, itemTrimStaticGrpFields[0]+" source ") {
						staticGrp["source"] = itemTrim
					}
					intFace["static_group"] = append(intFace["static_group"].([]map[string]interface{}), staticGrp)
				}
				confRead.interFace = append(confRead.interFace, intFace)
			case balt.CutPrefixInString(&itemTrim, "l2-querier source-address "):
				confRead.l2QuerierSrcAddress = itemTrim
			case balt.CutPrefixInString(&itemTrim, "proxy"):
				confRead.proxy = true
				if balt.CutPrefixInString(&itemTrim, " source-address ") {
					confRead.proxySrcAddress = itemTrim
				}
			case balt.CutPrefixInString(&itemTrim, "query-interval "):
				confRead.queryInterval, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "query-last-member-interval "):
				confRead.queryLastMemberInterval = itemTrim
			case balt.CutPrefixInString(&itemTrim, "query-response-interval "):
				confRead.queryResponseInterval = itemTrim
			case balt.CutPrefixInString(&itemTrim, "robust-count "):
				confRead.robustCount, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delIgmpSnoopingVlan(name, routingInstance string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)

	if routingInstance == junos.DefaultW {
		configSet = append(configSet, "delete protocols igmp-snooping vlan "+name)
	} else {
		configSet = append(configSet, delRoutingInstances+routingInstance+" protocols igmp-snooping vlan "+name)
	}

	return junSess.ConfigSet(configSet)
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
