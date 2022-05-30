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

type accessAddressAssignPoolOptions struct {
	activeDrain     bool
	holdDown        bool
	name            string
	link            string
	routingInstance string
	family          []map[string]interface{}
}

func resourceAccessAddressAssignPool() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAccessAddressAssignPoolCreate,
		ReadWithoutTimeout:   resourceAccessAddressAssignPoolRead,
		UpdateWithoutTimeout: resourceAccessAddressAssignPoolUpdate,
		DeleteWithoutTimeout: resourceAccessAddressAssignPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAccessAddressAssignPoolImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          defaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"family": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{inetW, inet6W}, false),
						},
						"network": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
						},
						"dhcp_attributes": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"boot_file": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"boot_server": {
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefAndDots),
									},
									"dns_server": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validateIsIPv6Address,
										},
									},
									"domain_name": {
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefAndDots),
									},
									"exclude_prefix_len": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 128),
									},
									"grace_period": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"maximum_lease_time": {
										Type:     schema.TypeInt,
										Optional: true,
										ConflictsWith: []string{
											"family.0.dhcp_attributes.0.maximum_lease_time_infinite",
											"family.0.dhcp_attributes.0.preferred_lifetime",
											"family.0.dhcp_attributes.0.preferred_lifetime_infinite",
											"family.0.dhcp_attributes.0.valid_lifetime",
											"family.0.dhcp_attributes.0.valid_lifetime_infinite",
										},
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"maximum_lease_time_infinite": {
										Type:     schema.TypeBool,
										Optional: true,
										ConflictsWith: []string{
											"family.0.dhcp_attributes.0.maximum_lease_time",
											"family.0.dhcp_attributes.0.preferred_lifetime",
											"family.0.dhcp_attributes.0.preferred_lifetime_infinite",
											"family.0.dhcp_attributes.0.valid_lifetime",
											"family.0.dhcp_attributes.0.valid_lifetime_infinite",
										},
									},
									"name_server": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.IsIPv4Address,
										},
									},
									"netbios_node_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"b-node", "h-node", "m-node", "p-node"}, false),
									},
									"next_server": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
									"option": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringMatch(regexp.MustCompile(
												`^\d+ (array )?(byte|flag|hex-string|integer|ip-address|short|string|unsigned-integer|unsigned-short) .*$`),
												"need to match '^\\d+ (array )?"+
													"(byte|flag|hex-string|integer|ip-address|short|string|unsigned-integer|unsigned-short) .*$'"),
										},
									},
									"option_match_82_circuit_id": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"range": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"option_match_82_remote_id": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"range": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"preferred_lifetime": {
										Type:     schema.TypeInt,
										Optional: true,
										ConflictsWith: []string{
											"family.0.dhcp_attributes.0.preferred_lifetime_infinite",
											"family.0.dhcp_attributes.0.maximum_lease_time",
											"family.0.dhcp_attributes.0.maximum_lease_time_infinite",
										},
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"preferred_lifetime_infinite": {
										Type:     schema.TypeBool,
										Optional: true,
										ConflictsWith: []string{
											"family.0.dhcp_attributes.0.preferred_lifetime",
											"family.0.dhcp_attributes.0.maximum_lease_time",
											"family.0.dhcp_attributes.0.maximum_lease_time_infinite",
										},
									},
									"propagate_ppp_settings": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"propagate_settings": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"router": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.IsIPv4Address,
										},
									},
									"server_identifier": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
									"sip_server_inet_address": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.IsIPv4Address,
										},
									},
									"sip_server_inet_domain_name": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"sip_server_inet6_address": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validateIsIPv6Address,
										},
									},
									"sip_server_inet6_domain_name": {
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefAndDots),
									},
									"t1_percentage": {
										Type:     schema.TypeInt,
										Optional: true,
										ConflictsWith: []string{
											"family.0.dhcp_attributes.0.t1_renewal_time",
											"family.0.dhcp_attributes.0.t2_rebinding_time",
										},
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 100),
									},
									"t1_renewal_time": {
										Type:     schema.TypeInt,
										Optional: true,
										ConflictsWith: []string{
											"family.0.dhcp_attributes.0.t1_percentage",
											"family.0.dhcp_attributes.0.t2_percentage",
										},
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"t2_percentage": {
										Type:     schema.TypeInt,
										Optional: true,
										ConflictsWith: []string{
											"family.0.dhcp_attributes.0.t1_renewal_time",
											"family.0.dhcp_attributes.0.t2_rebinding_time",
										},
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 100),
									},
									"t2_rebinding_time": {
										Type:     schema.TypeInt,
										Optional: true,
										ConflictsWith: []string{
											"family.0.dhcp_attributes.0.t1_percentage",
											"family.0.dhcp_attributes.0.t2_percentage",
										},
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"tftp_server": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
									"valid_lifetime": {
										Type:     schema.TypeInt,
										Optional: true,
										ConflictsWith: []string{
											"family.0.dhcp_attributes.0.valid_lifetime_infinite",
											"family.0.dhcp_attributes.0.maximum_lease_time",
											"family.0.dhcp_attributes.0.maximum_lease_time_infinite",
										},
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"valid_lifetime_infinite": {
										Type:     schema.TypeBool,
										Optional: true,
										ConflictsWith: []string{
											"family.0.dhcp_attributes.0.valid_lifetime",
											"family.0.dhcp_attributes.0.maximum_lease_time",
											"family.0.dhcp_attributes.0.maximum_lease_time_infinite",
										},
									},
									"wins_server": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.IsIPv4Address,
										},
									},
								},
							},
						},
						"excluded_address": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"excluded_range": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validateNameObjectJunos([]string{}, 63, formatDefault),
									},
									"low": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsIPAddress,
									},
									"high": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsIPAddress,
									},
								},
							},
						},
						"host": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validateNameObjectJunos([]string{}, 63, formatDefault),
									},
									"hardware_address": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsMACAddress,
									},
									"ip_address": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
								},
							},
						},
						"inet_range": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"family.0.inet6_range"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validateNameObjectJunos([]string{}, 63, formatDefault),
									},
									"low": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
									"high": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
								},
							},
						},
						"inet6_range": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"family.0.inet_range"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validateNameObjectJunos([]string{}, 63, formatDefault),
									},
									"low": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsCIDR,
									},
									"high": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsCIDR,
									},
									"prefix_length": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 128),
									},
								},
							},
						},
						"xauth_attributes_primary_dns": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateIPMaskFunc(),
						},
						"xauth_attributes_primary_wins": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateIPMaskFunc(),
						},
						"xauth_attributes_secondary_dns": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateIPMaskFunc(),
						},
						"xauth_attributes_secondary_wins": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateIPMaskFunc(),
						},
					},
				},
			},
			"active_drain": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"hold_down": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"link": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
		},
	}
}

func resourceAccessAddressAssignPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setAccessAddressAssignPool(d, clt, nil); err != nil {
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
	accessAddressAssignPoolExists, err := checkAccessAddressAssignPoolExists(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if accessAddressAssignPoolExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(
			fmt.Errorf("access address-assignment pool %v already exists in routing-instance %s",
				d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}
	if err := setAccessAddressAssignPool(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_access_address_assignment_pool", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	accessAddressAssignPoolExists, err = checkAccessAddressAssignPoolExists(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if accessAddressAssignPoolExists {
		d.SetId(d.Get("name").(string) + idSeparator + d.Get("routing_instance").(string))
	} else {
		return append(diagWarns,
			diag.FromErr(fmt.Errorf("access address-assignment pool %v not exists in routing_instance %s after commit "+
				"=> check your config", d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceAccessAddressAssignPoolReadWJunSess(d, clt, junSess)...)
}

func resourceAccessAddressAssignPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceAccessAddressAssignPoolReadWJunSess(d, clt, junSess)
}

func resourceAccessAddressAssignPoolReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	accessAddressAssignPoolOptions, err := readAccessAddressAssignPool(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if accessAddressAssignPoolOptions.name == "" {
		d.SetId("")
	} else {
		fillAccessAddressAssignPoolData(d, accessAddressAssignPoolOptions)
	}

	return nil
}

func resourceAccessAddressAssignPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delAccessAddressAssignPool(
			d.Get("name").(string),
			d.Get("routing_instance").(string),
			clt, nil,
		); err != nil {
			return diag.FromErr(err)
		}
		if err := setAccessAddressAssignPool(d, clt, nil); err != nil {
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
	if err := delAccessAddressAssignPool(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		clt, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setAccessAddressAssignPool(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_access_address_assignment_pool", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceAccessAddressAssignPoolReadWJunSess(d, clt, junSess)...)
}

func resourceAccessAddressAssignPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delAccessAddressAssignPool(
			d.Get("name").(string),
			d.Get("routing_instance").(string),
			clt, nil,
		); err != nil {
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
	if err := delAccessAddressAssignPool(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		clt, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_access_address_assignment_pool", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceAccessAddressAssignPoolImport(ctx context.Context, d *schema.ResourceData, m interface{},
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
	accessAddressAssignPoolExists, err := checkAccessAddressAssignPoolExists(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	if !accessAddressAssignPoolExists {
		return nil, fmt.Errorf("don't find access address-assignment pool with id '%v' (id must be "+
			"<name>"+idSeparator+"<routing_instance>)", d.Id())
	}
	accessAddressAssignPoolOptions, err := readAccessAddressAssignPool(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	fillAccessAddressAssignPoolData(d, accessAddressAssignPoolOptions)

	result[0] = d

	return result, nil
}

func checkAccessAddressAssignPoolExists(name, instance string, clt *Client, junSess *junosSession) (bool, error) {
	var showConfig string
	var err error
	if instance == defaultW {
		showConfig, err = clt.command(cmdShowConfig+
			"access address-assignment pool "+name+pipeDisplaySet, junSess)
		if err != nil {
			return false, err
		}
	} else {
		showConfig, err = clt.command(cmdShowConfig+routingInstancesWS+instance+" "+
			"access address-assignment pool "+name+pipeDisplaySet, junSess)
		if err != nil {
			return false, err
		}
	}

	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setAccessAddressAssignPool(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := setLS
	if d.Get("routing_instance").(string) != defaultW {
		setPrefix = setRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	setPrefix += "access address-assignment pool " + d.Get("name").(string) + " "

	for _, fi := range d.Get("family").([]interface{}) {
		family := fi.(map[string]interface{})

		configSetFamily, err := setAccessAddressAssignPoolFamily(family, setPrefix)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetFamily...)
	}
	if d.Get("active_drain").(bool) {
		configSet = append(configSet, setPrefix+"active-drain")
	}
	if d.Get("hold_down").(bool) {
		configSet = append(configSet, setPrefix+"hold-down")
	}
	if v := d.Get("link").(string); v != "" {
		configSet = append(configSet, setPrefix+"link "+v)
	}

	return clt.configSet(configSet, junSess)
}

func setAccessAddressAssignPoolFamily(family map[string]interface{}, setPrefix string) ([]string, error) {
	configSet := make([]string, 0)

	setPrefixFamily := setPrefix + "family inet "

	switch family["type"].(string) {
	case inetW:
		configSet = append(configSet, setPrefixFamily+"network "+family["network"].(string))
	case inet6W:
		setPrefixFamily = setPrefix + "family inet6 "
		configSet = append(configSet, setPrefixFamily+"prefix "+family["network"].(string))
	}
	for _, da := range family["dhcp_attributes"].([]interface{}) {
		dhcpAttr := da.(map[string]interface{})
		setPrefixDhcpAttr := setPrefixFamily + "dhcp-attributes "

		configSetDhcpAttributes, err := setAccessAddressAssignPoolFamilyDhcpAttributes(
			dhcpAttr, family["type"].(string), setPrefixDhcpAttr)
		if err != nil {
			return configSet, err
		}
		configSet = append(configSet, configSetDhcpAttributes...)
	}
	for _, v := range sortSetOfString(family["excluded_address"].(*schema.Set).List()) {
		switch family["type"].(string) {
		case inetW:
			if _, errs := validation.IsIPv4Address(v, "family.0.excluded_address"); len(errs) > 0 {
				return configSet, errs[0]
			}
		case inet6W:
			if _, errs := validateIsIPv6Address(v, "family.0.excluded_address"); len(errs) > 0 {
				return configSet, errs[0]
			}
		}
		configSet = append(configSet, setPrefixFamily+"excluded-address "+v)
	}
	excludedRangeNameList := make([]string, 0)
	for _, v := range family["excluded_range"].([]interface{}) {
		excludedRange := v.(map[string]interface{})
		if bchk.StringInSlice(excludedRange["name"].(string), excludedRangeNameList) {
			return configSet, fmt.Errorf("multiple blocks excluded_range with the same name %s", excludedRange["name"].(string))
		}
		excludedRangeNameList = append(excludedRangeNameList, excludedRange["name"].(string))
		configSet = append(configSet,
			setPrefixFamily+"excluded-range "+excludedRange["name"].(string)+" low "+excludedRange["low"].(string))
		configSet = append(configSet,
			setPrefixFamily+"excluded-range "+excludedRange["name"].(string)+" high "+excludedRange["high"].(string))
	}
	hostNameList := make([]string, 0)
	for _, v := range family["host"].([]interface{}) {
		if family["type"].(string) == inet6W {
			return configSet, fmt.Errorf("host not compatible when type = inet6")
		}
		host := v.(map[string]interface{})
		if bchk.StringInSlice(host["name"].(string), hostNameList) {
			return configSet, fmt.Errorf("multiple blocks host with the same name %s", host["name"].(string))
		}
		hostNameList = append(hostNameList, host["name"].(string))
		configSet = append(configSet,
			setPrefixFamily+"host "+host["name"].(string)+" hardware-address "+host["hardware_address"].(string))
		configSet = append(configSet,
			setPrefixFamily+"host "+host["name"].(string)+" ip-address "+host["ip_address"].(string))
	}
	rangeNameList := make([]string, 0)
	switch family["type"].(string) {
	case inetW:
		if len(family["inet6_range"].([]interface{})) > 0 {
			return configSet, fmt.Errorf("inet6_range not compatible when type = inet")
		}
		for _, v := range family["inet_range"].([]interface{}) {
			rangeBlck := v.(map[string]interface{})
			if bchk.StringInSlice(rangeBlck["name"].(string), rangeNameList) {
				return configSet, fmt.Errorf("multiple blocks inet_range with the same name %s", rangeBlck["name"].(string))
			}
			rangeNameList = append(rangeNameList, rangeBlck["name"].(string))
			configSet = append(configSet,
				setPrefixFamily+"range "+rangeBlck["name"].(string)+" low "+rangeBlck["low"].(string))
			configSet = append(configSet,
				setPrefixFamily+"range "+rangeBlck["name"].(string)+" high "+rangeBlck["high"].(string))
		}
	case inet6W:
		if len(family["inet_range"].([]interface{})) > 0 {
			return configSet, fmt.Errorf("inet_range not compatible when type = inet6")
		}
		for _, v := range family["inet6_range"].([]interface{}) {
			rangeBlck := v.(map[string]interface{})
			if bchk.StringInSlice(rangeBlck["name"].(string), rangeNameList) {
				return configSet, fmt.Errorf("multiple blocks inet6_range with the same name %s", rangeBlck["name"].(string))
			}
			rangeNameList = append(rangeNameList, rangeBlck["name"].(string))
			switch {
			case rangeBlck["prefix_length"].(int) != 0 &&
				(rangeBlck["low"].(string) != "" || rangeBlck["high"].(string) != ""):
				return configSet,
					fmt.Errorf("conflict between prefix_length and low/high in inet6_range %s", rangeBlck["name"].(string))
			case rangeBlck["prefix_length"].(int) != 0:
				configSet = append(configSet, setPrefixFamily+"range "+rangeBlck["name"].(string)+
					" prefix-length "+strconv.Itoa(rangeBlck["prefix_length"].(int)))
			case rangeBlck["low"].(string) != "" && rangeBlck["high"].(string) != "":
				configSet = append(configSet, setPrefixFamily+"range "+rangeBlck["name"].(string)+
					" low "+rangeBlck["low"].(string))
				configSet = append(configSet, setPrefixFamily+"range "+rangeBlck["name"].(string)+
					" high "+rangeBlck["high"].(string))
			default:
				return configSet, fmt.Errorf("missing prefix_length or low & high for inet6_range %s", rangeBlck["name"].(string))
			}
		}
	}
	if v := family["xauth_attributes_primary_dns"].(string); v != "" {
		if family["type"].(string) == inet6W {
			return configSet, fmt.Errorf("xauth_attributes_primary_dns not compatible when type = inet6")
		}
		if _, errs := validation.IsIPv4Address(strings.Split(v, "/")[0], ""); len(errs) > 0 {
			return configSet, fmt.Errorf("%s is not a IPv4", v)
		}
		configSet = append(configSet, setPrefixFamily+"xauth-attributes primary-dns "+v)
	}
	if v := family["xauth_attributes_primary_wins"].(string); v != "" {
		if family["type"].(string) == inet6W {
			return configSet, fmt.Errorf("xauth_attributes_primary_wins not compatible when type = inet6")
		}
		if _, errs := validation.IsIPv4Address(strings.Split(v, "/")[0], ""); len(errs) > 0 {
			return configSet, fmt.Errorf("%s is not a IPv4", v)
		}
		configSet = append(configSet, setPrefixFamily+"xauth-attributes primary-wins "+v)
	}
	if v := family["xauth_attributes_secondary_dns"].(string); v != "" {
		if family["type"].(string) == inet6W {
			return configSet, fmt.Errorf("xauth_attributes_secondary_dns not compatible when type = inet6")
		}
		if _, errs := validation.IsIPv4Address(strings.Split(v, "/")[0], ""); len(errs) > 0 {
			return configSet, fmt.Errorf("%s is not a IPv4", v)
		}
		configSet = append(configSet, setPrefixFamily+"xauth-attributes secondary-dns "+v)
	}
	if v := family["xauth_attributes_secondary_wins"].(string); v != "" {
		if family["type"].(string) == inet6W {
			return configSet, fmt.Errorf("xauth_attributes_secondary_wins not compatible when type = inet6")
		}
		if _, errs := validation.IsIPv4Address(strings.Split(v, "/")[0], ""); len(errs) > 0 {
			return configSet, fmt.Errorf("%s is not a IPv4", v)
		}
		configSet = append(configSet, setPrefixFamily+"xauth-attributes secondary-wins "+v)
	}

	return configSet, nil
}

func setAccessAddressAssignPoolFamilyDhcpAttributes(dhcpAttr map[string]interface{}, familyType, setPrefix string,
) ([]string, error) {
	configSet := make([]string, 0)

	if v := dhcpAttr["boot_file"].(string); v != "" {
		configSet = append(configSet, setPrefix+"boot-file \""+v+"\"")
	}
	if v := dhcpAttr["boot_server"].(string); v != "" {
		configSet = append(configSet, setPrefix+"boot-server "+v)
	}
	for _, v := range dhcpAttr["dns_server"].([]interface{}) {
		if familyType == inetW {
			return configSet, fmt.Errorf("dhcp_attributes.0.dns_server not compatible when type = inet")
		}

		configSet = append(configSet, setPrefix+"dns-server "+v.(string))
	}
	if v := dhcpAttr["domain_name"].(string); v != "" {
		configSet = append(configSet, setPrefix+"domain-name "+v)
	}
	if v := dhcpAttr["exclude_prefix_len"].(int); v != 0 {
		if familyType == inetW {
			return configSet, fmt.Errorf("dhcp_attributes.0.exclude_prefix_len not compatible when type = inet")
		}
		configSet = append(configSet, setPrefix+"exclude-prefix-len "+strconv.Itoa(v))
	}
	if v := dhcpAttr["grace_period"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"grace-period "+strconv.Itoa(v))
	}
	if v := dhcpAttr["maximum_lease_time"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"maximum-lease-time "+strconv.Itoa(v))
	}
	if dhcpAttr["maximum_lease_time_infinite"].(bool) {
		configSet = append(configSet, setPrefix+"maximum-lease-time infinite")
	}
	for _, v := range dhcpAttr["name_server"].([]interface{}) {
		configSet = append(configSet, setPrefix+"name-server "+v.(string))
	}
	if v := dhcpAttr["netbios_node_type"].(string); v != "" {
		configSet = append(configSet, setPrefix+"netbios-node-type "+v)
	}
	if v := dhcpAttr["next_server"].(string); v != "" {
		configSet = append(configSet, setPrefix+"next-server "+v)
	}
	for _, v := range sortSetOfString(dhcpAttr["option"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"option "+v)
	}
	optionMatch82CircuitIDValueList := make([]string, 0)
	for _, v := range dhcpAttr["option_match_82_circuit_id"].([]interface{}) {
		opt := v.(map[string]interface{})
		if bchk.StringInSlice(opt["value"].(string), optionMatch82CircuitIDValueList) {
			return configSet,
				fmt.Errorf("multiple blocks option_match_82_circuit_id with the same value %s", opt["value"].(string))
		}
		optionMatch82CircuitIDValueList = append(optionMatch82CircuitIDValueList, opt["value"].(string))
		configSet = append(configSet,
			setPrefix+"option-match option-82 "+
				"circuit-id \""+opt["value"].(string)+"\" range \""+opt["range"].(string)+"\"")
	}
	optionMatch82RemoteIDValueList := make([]string, 0)
	for _, v := range dhcpAttr["option_match_82_remote_id"].([]interface{}) {
		opt := v.(map[string]interface{})
		if bchk.StringInSlice(opt["value"].(string), optionMatch82RemoteIDValueList) {
			return configSet,
				fmt.Errorf("multiple blocks option_match_82_remote_id with the same value %s", opt["value"].(string))
		}
		optionMatch82RemoteIDValueList = append(optionMatch82RemoteIDValueList, opt["value"].(string))
		configSet = append(configSet,
			setPrefix+"option-match option-82 "+
				"remote-id \""+opt["value"].(string)+"\" range \""+opt["range"].(string)+"\"")
	}
	if v := dhcpAttr["preferred_lifetime"].(int); v != -1 {
		if familyType == inetW {
			return configSet, fmt.Errorf("dhcp_attributes.0.preferred_lifetime not compatible when type = inet")
		}
		configSet = append(configSet, setPrefix+"preferred-lifetime "+strconv.Itoa(v))
	}
	if dhcpAttr["preferred_lifetime_infinite"].(bool) {
		if familyType == inetW {
			return configSet, fmt.Errorf("dhcp_attributes.0.preferred_lifetime_infinite not compatible when type = inet")
		}
		configSet = append(configSet, setPrefix+"preferred-lifetime infinite")
	}
	for _, v := range sortSetOfString(dhcpAttr["propagate_ppp_settings"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"propagate-ppp-settings "+v)
	}
	if v := dhcpAttr["propagate_settings"].(string); v != "" {
		configSet = append(configSet, setPrefix+"propagate-settings \""+v+"\"")
	}
	for _, v := range dhcpAttr["router"].([]interface{}) {
		configSet = append(configSet, setPrefix+"router "+v.(string))
	}
	if v := dhcpAttr["server_identifier"].(string); v != "" {
		configSet = append(configSet, setPrefix+"server-identifier "+v)
	}
	for _, v := range dhcpAttr["sip_server_inet_address"].([]interface{}) {
		configSet = append(configSet, setPrefix+"sip-server ip-address "+v.(string))
	}
	for _, v := range dhcpAttr["sip_server_inet6_address"].([]interface{}) {
		if familyType == inetW {
			return configSet, fmt.Errorf("dhcp_attributes.0.sip_server_inet6_address not compatible when type = inet")
		}
		configSet = append(configSet, setPrefix+"sip-server-address "+v.(string))
	}
	for _, v := range dhcpAttr["sip_server_inet_domain_name"].([]interface{}) {
		configSet = append(configSet, setPrefix+"sip-server name \""+v.(string)+"\"")
	}
	if v := dhcpAttr["sip_server_inet6_domain_name"].(string); v != "" {
		if familyType == inetW {
			return configSet, fmt.Errorf("dhcp_attributes.0.sip_server_inet6_domain_name not compatible when type = inet")
		}
		configSet = append(configSet, setPrefix+"sip-server-domain-name \""+v+"\"")
	}
	if v := dhcpAttr["t1_percentage"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"t1-percentage "+strconv.Itoa(v))
	}
	if v := dhcpAttr["t1_renewal_time"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"t1-renewal-time "+strconv.Itoa(v))
	}
	if v := dhcpAttr["t2_percentage"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"t2-percentage "+strconv.Itoa(v))
	}
	if v := dhcpAttr["t2_rebinding_time"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"t2-rebinding-time "+strconv.Itoa(v))
	}
	if v := dhcpAttr["tftp_server"].(string); v != "" {
		configSet = append(configSet, setPrefix+"tftp-server "+v)
	}
	if v := dhcpAttr["valid_lifetime"].(int); v != -1 {
		if familyType == inetW {
			return configSet, fmt.Errorf("dhcp_attributes.0.valid_lifetime not compatible when type = inet")
		}
		configSet = append(configSet, setPrefix+"valid-lifetime "+strconv.Itoa(v))
	}
	if dhcpAttr["valid_lifetime_infinite"].(bool) {
		if familyType == inetW {
			return configSet, fmt.Errorf("dhcp_attributes.0.valid_lifetime_infinite not compatible when type = inet")
		}
		configSet = append(configSet, setPrefix+"valid-lifetime infinite")
	}
	for _, v := range dhcpAttr["wins_server"].([]interface{}) {
		configSet = append(configSet, setPrefix+"wins-server "+v.(string))
	}
	if len(configSet) == 0 {
		return configSet, fmt.Errorf("family.0.dhcp_attributes block is empty")
	}

	return configSet, nil
}

func readAccessAddressAssignPool(name, instance string, clt *Client, junSess *junosSession,
) (accessAddressAssignPoolOptions, error) {
	var confRead accessAddressAssignPoolOptions
	var showConfig string
	var err error

	if instance == defaultW {
		showConfig, err = clt.command(cmdShowConfig+
			"access address-assignment pool "+name+pipeDisplaySetRelative, junSess)
	} else {
		showConfig, err = clt.command(cmdShowConfig+routingInstancesWS+instance+" "+
			"access address-assignment pool "+name+pipeDisplaySetRelative, junSess)
	}
	if err != nil {
		return confRead, err
	}

	if showConfig != emptyW {
		confRead.name = name
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
			case strings.HasPrefix(itemTrim, "family "):
				itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "family "), " ")
				if len(confRead.family) == 0 {
					confRead.family = append(confRead.family, map[string]interface{}{
						"type":                            itemTrimSplit[0],
						"network":                         "",
						"dhcp_attributes":                 make([]map[string]interface{}, 0),
						"excluded_address":                make([]string, 0),
						"excluded_range":                  make([]map[string]interface{}, 0),
						"host":                            make([]map[string]interface{}, 0),
						"inet_range":                      make([]map[string]interface{}, 0),
						"inet6_range":                     make([]map[string]interface{}, 0),
						"xauth_attributes_primary_dns":    "",
						"xauth_attributes_primary_wins":   "",
						"xauth_attributes_secondary_dns":  "",
						"xauth_attributes_secondary_wins": "",
					})
				}
				if err := readAccessAddressAssignPoolFamily(
					strings.TrimPrefix(itemTrim, "family "+itemTrimSplit[0]+" "),
					confRead.family[0],
				); err != nil {
					return confRead, err
				}
			case itemTrim == "active-drain":
				confRead.activeDrain = true
			case itemTrim == "hold-down":
				confRead.holdDown = true
			case strings.HasPrefix(itemTrim, "link "):
				confRead.link = strings.TrimPrefix(itemTrim, "link ")
			}
		}
	}

	return confRead, nil
}

func readAccessAddressAssignPoolFamily(itemTrim string, family map[string]interface{}) error {
	var err error
	switch {
	case strings.HasPrefix(itemTrim, "network "):
		family["network"] = strings.TrimPrefix(itemTrim, "network ")
	case strings.HasPrefix(itemTrim, "prefix "):
		family["network"] = strings.TrimPrefix(itemTrim, "prefix ")
	case strings.HasPrefix(itemTrim, "dhcp-attributes "):
		if len(family["dhcp_attributes"].([]map[string]interface{})) == 0 {
			family["dhcp_attributes"] = append(family["dhcp_attributes"].([]map[string]interface{}), map[string]interface{}{
				"boot_file":                    "",
				"boot_server":                  "",
				"dns_server":                   make([]string, 0),
				"domain_name":                  "",
				"exclude_prefix_len":           0,
				"grace_period":                 -1,
				"maximum_lease_time":           -1,
				"maximum_lease_time_infinite":  false,
				"name_server":                  make([]string, 0),
				"netbios_node_type":            "",
				"next_server":                  "",
				"option":                       make([]string, 0),
				"option_match_82_circuit_id":   make([]map[string]interface{}, 0),
				"option_match_82_remote_id":    make([]map[string]interface{}, 0),
				"preferred_lifetime":           -1,
				"preferred_lifetime_infinite":  false,
				"propagate_ppp_settings":       make([]string, 0),
				"propagate_settings":           "",
				"router":                       make([]string, 0),
				"server_identifier":            "",
				"sip_server_inet_address":      make([]string, 0),
				"sip_server_inet_domain_name":  make([]string, 0),
				"sip_server_inet6_address":     make([]string, 0),
				"sip_server_inet6_domain_name": "",
				"t1_percentage":                -1,
				"t1_renewal_time":              -1,
				"t2_percentage":                -1,
				"t2_rebinding_time":            -1,
				"tftp_server":                  "",
				"valid_lifetime":               -1,
				"valid_lifetime_infinite":      false,
				"wins_server":                  make([]string, 0),
			})
		}
		itemTrimAttr := strings.TrimPrefix(itemTrim, "dhcp-attributes ")
		dhcpAttr := family["dhcp_attributes"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrimAttr, "boot-file "):
			dhcpAttr["boot_file"] = strings.Trim(strings.TrimPrefix(itemTrimAttr, "boot-file "), "\"")
		case strings.HasPrefix(itemTrimAttr, "boot-server "):
			dhcpAttr["boot_server"] = strings.TrimPrefix(itemTrimAttr, "boot-server ")
		case strings.HasPrefix(itemTrimAttr, "dns-server "):
			dhcpAttr["dns_server"] = append(dhcpAttr["dns_server"].([]string), strings.TrimPrefix(itemTrimAttr, "dns-server "))
		case strings.HasPrefix(itemTrimAttr, "domain-name "):
			dhcpAttr["domain_name"] = strings.TrimPrefix(itemTrimAttr, "domain-name ")
		case strings.HasPrefix(itemTrimAttr, "exclude-prefix-len "):
			dhcpAttr["exclude_prefix_len"], err = strconv.Atoi(strings.TrimPrefix(itemTrimAttr, "exclude-prefix-len "))
		case strings.HasPrefix(itemTrimAttr, "grace-period "):
			dhcpAttr["grace_period"], err = strconv.Atoi(strings.TrimPrefix(itemTrimAttr, "grace-period "))
		case itemTrimAttr == "maximum-lease-time infinite":
			dhcpAttr["maximum_lease_time_infinite"] = true
		case strings.HasPrefix(itemTrimAttr, "maximum-lease-time "):
			dhcpAttr["maximum_lease_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrimAttr, "maximum-lease-time "))
		case strings.HasPrefix(itemTrimAttr, "name-server "):
			dhcpAttr["name_server"] = append(dhcpAttr["name_server"].([]string),
				strings.TrimPrefix(itemTrimAttr, "name-server "))
		case strings.HasPrefix(itemTrimAttr, "netbios-node-type "):
			dhcpAttr["netbios_node_type"] = strings.TrimPrefix(itemTrimAttr, "netbios-node-type ")
		case strings.HasPrefix(itemTrimAttr, "next-server "):
			dhcpAttr["next_server"] = strings.TrimPrefix(itemTrimAttr, "next-server ")
		case strings.HasPrefix(itemTrimAttr, "option "):
			dhcpAttr["option"] = append(dhcpAttr["option"].([]string), strings.TrimPrefix(itemTrimAttr, "option "))
		case strings.HasPrefix(itemTrimAttr, "option-match option-82 circuit-id "):
			itemTrimAttrSplit := strings.Split(strings.TrimPrefix(itemTrimAttr, "option-match option-82 circuit-id "), " ")
			if len(itemTrimAttrSplit) < 3 {
				return fmt.Errorf("can't find range from '%s'", itemTrimAttr)
			}
			dhcpAttr["option_match_82_circuit_id"] = append(
				dhcpAttr["option_match_82_circuit_id"].([]map[string]interface{}),
				map[string]interface{}{
					"value": itemTrimAttrSplit[0],
					"range": itemTrimAttrSplit[2],
				},
			)
		case strings.HasPrefix(itemTrimAttr, "option-match option-82 remote-id "):
			itemTrimAttrSplit := strings.Split(strings.TrimPrefix(itemTrimAttr, "option-match option-82 remote-id "), " ")
			if len(itemTrimAttrSplit) < 3 {
				return fmt.Errorf("can't find range from '%s'", itemTrimAttr)
			}
			dhcpAttr["option_match_82_remote_id"] = append(
				dhcpAttr["option_match_82_remote_id"].([]map[string]interface{}),
				map[string]interface{}{
					"value": itemTrimAttrSplit[0],
					"range": itemTrimAttrSplit[2],
				},
			)
		case itemTrimAttr == "preferred-lifetime infinite":
			dhcpAttr["preferred_lifetime_infinite"] = true
		case strings.HasPrefix(itemTrimAttr, "preferred-lifetime "):
			dhcpAttr["preferred_lifetime"], err = strconv.Atoi(strings.TrimPrefix(itemTrimAttr, "preferred-lifetime "))
		case strings.HasPrefix(itemTrimAttr, "propagate-ppp-settings "):
			dhcpAttr["propagate_ppp_settings"] = append(dhcpAttr["propagate_ppp_settings"].([]string),
				strings.TrimPrefix(itemTrimAttr, "propagate-ppp-settings "))
		case strings.HasPrefix(itemTrimAttr, "propagate-settings "):
			dhcpAttr["propagate_settings"] = strings.TrimPrefix(itemTrimAttr, "propagate-settings ")
		case strings.HasPrefix(itemTrimAttr, "router "):
			dhcpAttr["router"] = append(dhcpAttr["router"].([]string), strings.TrimPrefix(itemTrimAttr, "router "))
		case strings.HasPrefix(itemTrimAttr, "server-identifier "):
			dhcpAttr["server_identifier"] = strings.TrimPrefix(itemTrimAttr, "server-identifier ")
		case strings.HasPrefix(itemTrimAttr, "sip-server ip-address "):
			dhcpAttr["sip_server_inet_address"] = append(dhcpAttr["sip_server_inet_address"].([]string),
				strings.TrimPrefix(itemTrimAttr, "sip-server ip-address "))
		case strings.HasPrefix(itemTrimAttr, "sip-server-address "):
			dhcpAttr["sip_server_inet6_address"] = append(dhcpAttr["sip_server_inet6_address"].([]string),
				strings.TrimPrefix(itemTrimAttr, "sip-server-address "))
		case strings.HasPrefix(itemTrimAttr, "sip-server name "):
			dhcpAttr["sip_server_inet_domain_name"] = append(dhcpAttr["sip_server_inet_domain_name"].([]string),
				strings.Trim(strings.TrimPrefix(itemTrimAttr, "sip-server name "), "\""))
		case strings.HasPrefix(itemTrimAttr, "sip-server-domain-name "):
			dhcpAttr["sip_server_inet6_domain_name"] = strings.Trim(strings.TrimPrefix(
				itemTrimAttr, "sip-server-domain-name "), "\"")
		case strings.HasPrefix(itemTrimAttr, "t1-percentage "):
			dhcpAttr["t1_percentage"], err = strconv.Atoi(strings.TrimPrefix(itemTrimAttr, "t1-percentage "))
		case strings.HasPrefix(itemTrimAttr, "t1-renewal-time "):
			dhcpAttr["t1_renewal_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrimAttr, "t1-renewal-time "))
		case strings.HasPrefix(itemTrimAttr, "t2-percentage "):
			dhcpAttr["t2_percentage"], err = strconv.Atoi(strings.TrimPrefix(itemTrimAttr, "t2-percentage "))
		case strings.HasPrefix(itemTrimAttr, "t2-rebinding-time "):
			dhcpAttr["t2_rebinding_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrimAttr, "t2-rebinding-time "))
		case strings.HasPrefix(itemTrimAttr, "tftp-server "):
			dhcpAttr["tftp_server"] = strings.TrimPrefix(itemTrimAttr, "tftp-server ")
		case itemTrimAttr == "valid-lifetime infinite":
			dhcpAttr["valid_lifetime_infinite"] = true
		case strings.HasPrefix(itemTrimAttr, "valid-lifetime "):
			dhcpAttr["valid_lifetime"], err = strconv.Atoi(strings.TrimPrefix(itemTrimAttr, "valid-lifetime "))
		case strings.HasPrefix(itemTrimAttr, "wins-server "):
			dhcpAttr["wins_server"] = append(dhcpAttr["wins_server"].([]string),
				strings.TrimPrefix(itemTrimAttr, "wins-server "))
		}
	case strings.HasPrefix(itemTrim, "excluded-address "):
		family["excluded_address"] = append(family["excluded_address"].([]string),
			strings.TrimPrefix(itemTrim, "excluded-address "))
	case strings.HasPrefix(itemTrim, "excluded-range "):
		itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "excluded-range "), " ")
		familyExcludedRange := map[string]interface{}{
			"name": itemTrimSplit[0],
			"low":  "",
			"high": "",
		}
		family["excluded_range"] = copyAndRemoveItemMapList(
			"name", familyExcludedRange, family["excluded_range"].([]map[string]interface{}))
		itemTrimExcludedRange := strings.TrimPrefix(itemTrim, "excluded-range "+itemTrimSplit[0]+" ")
		switch {
		case strings.HasPrefix(itemTrimExcludedRange, "low "):
			familyExcludedRange["low"] = strings.TrimPrefix(itemTrimExcludedRange, "low ")
		case strings.HasPrefix(itemTrimExcludedRange, "high "):
			familyExcludedRange["high"] = strings.TrimPrefix(itemTrimExcludedRange, "high ")
		}
		family["excluded_range"] = append(family["excluded_range"].([]map[string]interface{}), familyExcludedRange)
	case strings.HasPrefix(itemTrim, "host "):
		itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "host "), " ")
		familyHost := map[string]interface{}{
			"name":             itemTrimSplit[0],
			"hardware_address": "",
			"ip_address":       "",
		}
		family["host"] = copyAndRemoveItemMapList(
			"name", familyHost, family["host"].([]map[string]interface{}))
		itemTrimHost := strings.TrimPrefix(itemTrim, "host "+itemTrimSplit[0]+" ")
		switch {
		case strings.HasPrefix(itemTrimHost, "hardware-address "):
			familyHost["hardware_address"] = strings.TrimPrefix(itemTrimHost, "hardware-address ")
		case strings.HasPrefix(itemTrimHost, "ip-address "):
			familyHost["ip_address"] = strings.TrimPrefix(itemTrimHost, "ip-address ")
		}
		family["host"] = append(family["host"].([]map[string]interface{}), familyHost)
	case strings.HasPrefix(itemTrim, "range "):
		if family["type"] == inetW {
			itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "range "), " ")
			familyInetRange := map[string]interface{}{
				"name": itemTrimSplit[0],
				"low":  "",
				"high": "",
			}
			family["inet_range"] = copyAndRemoveItemMapList(
				"name", familyInetRange, family["inet_range"].([]map[string]interface{}))
			itemTrimRange := strings.TrimPrefix(itemTrim, "range "+itemTrimSplit[0]+" ")
			switch {
			case strings.HasPrefix(itemTrimRange, "low "):
				familyInetRange["low"] = strings.TrimPrefix(itemTrimRange, "low ")
			case strings.HasPrefix(itemTrimRange, "high "):
				familyInetRange["high"] = strings.TrimPrefix(itemTrimRange, "high ")
			}
			family["inet_range"] = append(family["inet_range"].([]map[string]interface{}), familyInetRange)
		} else if family["type"] == inet6W {
			itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "range "), " ")
			familyInet6Range := map[string]interface{}{
				"name":          itemTrimSplit[0],
				"low":           "",
				"high":          "",
				"prefix_length": 0,
			}
			family["inet6_range"] = copyAndRemoveItemMapList(
				"name", familyInet6Range, family["inet6_range"].([]map[string]interface{}))
			itemTrimRange := strings.TrimPrefix(itemTrim, "range "+itemTrimSplit[0]+" ")
			switch {
			case strings.HasPrefix(itemTrimRange, "low "):
				familyInet6Range["low"] = strings.TrimPrefix(itemTrimRange, "low ")
			case strings.HasPrefix(itemTrimRange, "high "):
				familyInet6Range["high"] = strings.TrimPrefix(itemTrimRange, "high ")
			case strings.HasPrefix(itemTrimRange, "prefix-length "):
				familyInet6Range["prefix_length"], err = strconv.Atoi(strings.TrimPrefix(itemTrimRange, "prefix-length "))
			}
			family["inet6_range"] = append(family["inet6_range"].([]map[string]interface{}), familyInet6Range)
		}
	case strings.HasPrefix(itemTrim, "xauth-attributes primary-dns "):
		family["xauth_attributes_primary_dns"] = strings.TrimPrefix(itemTrim, "xauth-attributes primary-dns ")
	case strings.HasPrefix(itemTrim, "xauth-attributes primary-wins "):
		family["xauth_attributes_primary_wins"] = strings.TrimPrefix(itemTrim, "xauth-attributes primary-wins ")
	case strings.HasPrefix(itemTrim, "xauth-attributes secondary-dns "):
		family["xauth_attributes_secondary_dns"] = strings.TrimPrefix(itemTrim, "xauth-attributes secondary-dns ")
	case strings.HasPrefix(itemTrim, "xauth-attributes secondary-wins "):
		family["xauth_attributes_secondary_wins"] = strings.TrimPrefix(itemTrim, "xauth-attributes secondary-wins ")
	}
	if err != nil {
		return fmt.Errorf(failedConvAtoiError, itemTrim, err)
	}

	return nil
}

func delAccessAddressAssignPool(name, instance string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	if instance == defaultW {
		configSet = append(configSet, "delete access address-assignment pool "+name)
	} else {
		configSet = append(configSet, delRoutingInstances+instance+" access address-assignment pool "+name)
	}

	return clt.configSet(configSet, junSess)
}

func fillAccessAddressAssignPoolData(
	d *schema.ResourceData, accessAddressAssignPoolOptions accessAddressAssignPoolOptions,
) {
	if tfErr := d.Set("name", accessAddressAssignPoolOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", accessAddressAssignPoolOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family", accessAddressAssignPoolOptions.family); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("active_drain", accessAddressAssignPoolOptions.activeDrain); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("hold_down", accessAddressAssignPoolOptions.holdDown); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("link", accessAddressAssignPoolOptions.link); tfErr != nil {
		panic(tfErr)
	}
}
