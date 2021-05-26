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
	jdecode "github.com/jeremmfr/junosdecode"
)

type interfaceLogicalOptions struct {
	vlanID                   int
	description              string
	routingInstance          string
	securityZone             string
	securityInboundProtocols []string
	securityInboundServices  []string
	familyInet               []map[string]interface{}
	familyInet6              []map[string]interface{}
}

func resourceInterfaceLogical() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInterfaceLogicalCreate,
		ReadContext:   resourceInterfaceLogicalRead,
		UpdateContext: resourceInterfaceLogicalUpdate,
		DeleteContext: resourceInterfaceLogicalDelete,
		Importer: &schema.ResourceImporter{
			State: resourceInterfaceLogicalImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
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
			"st0_also_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"family_inet": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_ip": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validateIPMaskFunc(),
									},
									"preferred": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"primary": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"vrrp_group": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"identifier": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(1, 255),
												},
												"virtual_address": {
													Type:     schema.TypeList,
													Required: true,
													MinItems: 1,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"accept_data": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"advertise_interval": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(1, 255),
												},
												"advertisements_threshold": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(1, 15),
												},
												"authentication_key": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"authentication_type": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"md5", "simple"}, false),
												},
												"no_accept_data": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"no_preempt": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"preempt": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"priority": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(1, 255),
												},
												"track_interface": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"interface": {
																Type:     schema.TypeString,
																Required: true,
															},
															"priority_cost": {
																Type:         schema.TypeInt,
																Required:     true,
																ValidateFunc: validation.IntBetween(1, 254),
															},
														},
													},
												},
												"track_route": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"route": {
																Type:     schema.TypeString,
																Required: true,
															},
															"routing_instance": {
																Type:     schema.TypeString,
																Required: true,
															},
															"priority_cost": {
																Type:         schema.TypeInt,
																Required:     true,
																ValidateFunc: validation.IntBetween(1, 254),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"filter_input": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"filter_output": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"mtu": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(500, 9192),
						},
						"rpf_check": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fail_filter": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"mode_loose": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"sampling_input": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sampling_output": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"family_inet6": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_ip": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validateIPMaskFunc(),
									},
									"preferred": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"primary": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"vrrp_group": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"identifier": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(1, 255),
												},
												"virtual_address": {
													Type:     schema.TypeList,
													Required: true,
													MinItems: 1,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"virtual_link_local_address": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.IsIPAddress,
												},
												"accept_data": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"advertise_interval": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(100, 40000),
												},
												"advertisements_threshold": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(1, 15),
												},
												"no_accept_data": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"no_preempt": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"preempt": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"priority": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(1, 255),
												},
												"track_interface": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"interface": {
																Type:     schema.TypeString,
																Required: true,
															},
															"priority_cost": {
																Type:         schema.TypeInt,
																Required:     true,
																ValidateFunc: validation.IntBetween(1, 254),
															},
														},
													},
												},
												"track_route": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"route": {
																Type:     schema.TypeString,
																Required: true,
															},
															"routing_instance": {
																Type:     schema.TypeString,
																Required: true,
															},
															"priority_cost": {
																Type:         schema.TypeInt,
																Required:     true,
																ValidateFunc: validation.IntBetween(1, 254),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"filter_input": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"filter_output": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"mtu": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(500, 9192),
						},
						"rpf_check": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fail_filter": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"mode_loose": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"sampling_input": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sampling_output": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"security_inbound_protocols": {
				Type:         schema.TypeList,
				Optional:     true,
				RequiredWith: []string{"security_zone"},
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"security_inbound_services": {
				Type:         schema.TypeList,
				Optional:     true,
				RequiredWith: []string{"security_zone"},
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"security_zone": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"vlan_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 4094),
			},
		},
	}
}

func resourceInterfaceLogicalCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := delInterfaceNC(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setInterfaceLogical(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	ncInt, emptyInt, _, err := checkInterfaceLogicalNCEmpty(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !ncInt && !emptyInt {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("interface %s already configured", d.Get("name").(string)))...)
	}
	if ncInt {
		if err := delInterfaceNC(d, m, jnprSess); err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	if d.Get("security_zone").(string) != "" {
		if !checkCompatibilitySecurity(jnprSess) {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(fmt.Errorf("security zone not compatible with Junos device %s",
				jnprSess.SystemInformation.HardwareModel))...)
		}
		zonesExists, err := checkSecurityZonesExists(d.Get("security_zone").(string), m, jnprSess)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !zonesExists {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("security zone %v doesn't exist", d.Get("security_zone").(string)))...)
		}
	}
	if d.Get("routing_instance").(string) != "" {
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
	if err := setInterfaceLogical(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_interface_logical", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ncInt, emptyInt, setInt, err := checkInterfaceLogicalNCEmpty(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ncInt {
		return append(diagWarns, diag.FromErr(fmt.Errorf("interface %v always disable after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}
	if emptyInt && !setInt {
		intExists, err := checkInterfaceExists(d.Get("name").(string), m, jnprSess)
		if err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		if !intExists {
			return append(diagWarns, diag.FromErr(fmt.Errorf("interface %v not exists and "+
				"config can't found after commit => check your config", d.Get("name").(string)))...)
		}
	}
	d.SetId(d.Get("name").(string))

	return append(diagWarns, resourceInterfaceLogicalReadWJnprSess(d, m, jnprSess)...)
}

func resourceInterfaceLogicalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceInterfaceLogicalReadWJnprSess(d, m, jnprSess)
}

func resourceInterfaceLogicalReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	ncInt, emptyInt, setInt, err := checkInterfaceLogicalNCEmpty(d.Get("name").(string), m, jnprSess)
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	if ncInt {
		d.SetId("")
		mutex.Unlock()

		return nil
	}
	if emptyInt && !setInt {
		intExists, err := checkInterfaceExists(d.Get("name").(string), m, jnprSess)
		if err != nil {
			mutex.Unlock()

			return diag.FromErr(err)
		}
		if !intExists {
			d.SetId("")
			mutex.Unlock()

			return nil
		}
	}
	interfaceLogicalOpt, err := readInterfaceLogical(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillInterfaceLogicalData(d, interfaceLogicalOpt)

	return nil
}

func resourceInterfaceLogicalUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delInterfaceLogicalOpts(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if d.HasChange("security_zone") {
		oSecurityZone, nSecurityZone := d.GetChange("security_zone")
		if nSecurityZone.(string) != "" {
			if !checkCompatibilitySecurity(jnprSess) {
				appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

				return append(diagWarns, diag.FromErr(fmt.Errorf("security zone not compatible with Junos device %s",
					jnprSess.SystemInformation.HardwareModel))...)
			}
			zonesExists, err := checkSecurityZonesExists(nSecurityZone.(string), m, jnprSess)
			if err != nil {
				appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
			if !zonesExists {
				appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

				return append(diagWarns, diag.FromErr(fmt.Errorf("security zone %v doesn't exist", nSecurityZone.(string)))...)
			}
		}
		if oSecurityZone.(string) != "" {
			err = delZoneInterfaceLogical(oSecurityZone.(string), d, m, jnprSess)
			if err != nil {
				appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
		}
	} else if v := d.Get("security_zone").(string); v != "" {
		if err := delZoneInterfaceLogical(v, d, m, jnprSess); err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	if d.HasChange("routing_instance") {
		oRoutingInstance, nRoutingInstance := d.GetChange("routing_instance")
		if nRoutingInstance.(string) != "" {
			instanceExists, err := checkRoutingInstanceExists(nRoutingInstance.(string), m, jnprSess)
			if err != nil {
				appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
			if !instanceExists {
				appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

				return append(diagWarns,
					diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", nRoutingInstance.(string)))...)
			}
		}
		if oRoutingInstance.(string) != "" {
			err = delRoutingInstanceInterfaceLogical(oRoutingInstance.(string), d, m, jnprSess)
			if err != nil {
				appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
		}
	}
	if err := setInterfaceLogical(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_interface_logical", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceInterfaceLogicalReadWJnprSess(d, m, jnprSess)...)
}

func resourceInterfaceLogicalDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delInterfaceLogical(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_interface_logical", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceInterfaceLogicalImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if strings.Count(d.Id(), ".") != 1 {
		return nil, fmt.Errorf("name of interface %s need to have 1 dot", d.Id())
	}
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	ncInt, emptyInt, setInt, err := checkInterfaceLogicalNCEmpty(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if ncInt {
		return nil, fmt.Errorf("interface '%v' is disabled, import is not possible", d.Id())
	}
	if emptyInt && !setInt {
		intExists, err := checkInterfaceExists(d.Id(), m, jnprSess)
		if err != nil {
			return nil, err
		}
		if !intExists {
			return nil, fmt.Errorf("don't find interface with id '%v' (id must be <name>)", d.Id())
		}
	}
	interfaceLogicalOpt, err := readInterfaceLogical(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if tfErr := d.Set("name", d.Id()); tfErr != nil {
		panic(tfErr)
	}
	fillInterfaceLogicalData(d, interfaceLogicalOpt)

	result[0] = d

	return result, nil
}

func checkInterfaceLogicalNCEmpty(interFace string, m interface{}, jnprSess *NetconfObject) (
	ncInt bool, emtyInt bool, justSet bool, _err error) {
	sess := m.(*Session)
	intConfig, err := sess.command("show configuration interfaces "+interFace+" | display set relative", jnprSess)
	if err != nil {
		return false, false, false, err
	}
	intConfigLines := make([]string, 0)
	// remove unused lines
	for _, item := range strings.Split(intConfig, "\n") {
		// exclude ethernet-switching (parameters in junos_interface_physical)
		if strings.Contains(item, "ethernet-switching") {
			continue
		}
		if strings.Contains(item, "<configuration-output>") {
			continue
		}
		if strings.Contains(item, "</configuration-output>") {
			break
		}
		if item == "" {
			continue
		}
		intConfigLines = append(intConfigLines, item)
	}
	if len(intConfigLines) == 0 {
		return false, true, true, nil
	}
	intConfig = strings.Join(intConfigLines, "\n")
	if sess.junosGroupIntDel != "" {
		if intConfig == "set apply-groups "+sess.junosGroupIntDel {
			return true, false, false, nil
		}
	}
	if intConfig == "set description NC\nset disable" ||
		intConfig == "set disable\nset description NC" {
		return true, false, false, nil
	}
	switch {
	case intConfig == setLineStart:
		return false, true, true, nil
	case intConfig == emptyWord:
		return false, true, false, nil
	default:
		return false, false, false, nil
	}
}

func setInterfaceLogical(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	intCut := strings.Split(d.Get("name").(string), ".")
	if len(intCut) != 2 {
		return fmt.Errorf("the name %s doesn't contain one dot", d.Get("name").(string))
	}
	configSet := make([]string, 0)
	setPrefix := "set interfaces " + d.Get("name").(string) + " "
	configSet = append(configSet, setPrefix)
	if d.Get("description").(string) != "" {
		configSet = append(configSet, setPrefix+"description \""+d.Get("description").(string)+"\"")
	}
	for _, v := range d.Get("family_inet").([]interface{}) {
		configSet = append(configSet, setPrefix+"family inet")
		if v != nil {
			familyInet := v.(map[string]interface{})
			for _, address := range familyInet["address"].([]interface{}) {
				var err error
				configSet, err = setFamilyAddress(address, configSet, setPrefix, inetWord)
				if err != nil {
					return err
				}
			}
			if familyInet["filter_input"].(string) != "" {
				configSet = append(configSet, setPrefix+"family inet filter input "+
					familyInet["filter_input"].(string))
			}
			if familyInet["filter_output"].(string) != "" {
				configSet = append(configSet, setPrefix+"family inet filter output "+
					familyInet["filter_output"].(string))
			}
			if familyInet["mtu"].(int) > 0 {
				configSet = append(configSet, setPrefix+"family inet mtu "+
					strconv.Itoa(familyInet["mtu"].(int)))
			}
			for _, v2 := range familyInet["rpf_check"].([]interface{}) {
				configSet = append(configSet, setPrefix+"family inet rpf-check")
				if v2 != nil {
					rpfCheck := v2.(map[string]interface{})
					if rpfCheck["fail_filter"].(string) != "" {
						configSet = append(configSet, setPrefix+"family inet rpf-check fail-filter "+
							"\""+rpfCheck["fail_filter"].(string)+"\"")
					}
					if rpfCheck["mode_loose"].(bool) {
						configSet = append(configSet, setPrefix+"family inet rpf-check mode loose ")
					}
				}
			}
			if familyInet["sampling_input"].(bool) {
				configSet = append(configSet, setPrefix+"family inet sampling input")
			}
			if familyInet["sampling_output"].(bool) {
				configSet = append(configSet, setPrefix+"family inet sampling output")
			}
		}
	}
	for _, v := range d.Get("family_inet6").([]interface{}) {
		configSet = append(configSet, setPrefix+"family inet6")
		if v != nil {
			familyInet6 := v.(map[string]interface{})
			for _, address := range familyInet6["address"].([]interface{}) {
				var err error
				configSet, err = setFamilyAddress(address, configSet, setPrefix, inet6Word)
				if err != nil {
					return err
				}
			}
			if familyInet6["filter_input"].(string) != "" {
				configSet = append(configSet, setPrefix+"family inet6 filter input "+
					familyInet6["filter_input"].(string))
			}
			if familyInet6["filter_output"].(string) != "" {
				configSet = append(configSet, setPrefix+"family inet6 filter output "+
					familyInet6["filter_output"].(string))
			}
			if familyInet6["mtu"].(int) > 0 {
				configSet = append(configSet, setPrefix+"family inet6 mtu "+
					strconv.Itoa(familyInet6["mtu"].(int)))
			}
			for _, v2 := range familyInet6["rpf_check"].([]interface{}) {
				configSet = append(configSet, setPrefix+"family inet6 rpf-check")
				if v2 != nil {
					rpfCheck := v2.(map[string]interface{})
					if rpfCheck["fail_filter"].(string) != "" {
						configSet = append(configSet, setPrefix+"family inet6 rpf-check fail-filter "+
							"\""+rpfCheck["fail_filter"].(string)+"\"")
					}
					if rpfCheck["mode_loose"].(bool) {
						configSet = append(configSet, setPrefix+"family inet6 rpf-check mode loose ")
					}
				}
			}
			if familyInet6["sampling_input"].(bool) {
				configSet = append(configSet, setPrefix+"family inet6 sampling input")
			}
			if familyInet6["sampling_output"].(bool) {
				configSet = append(configSet, setPrefix+"family inet6 sampling output")
			}
		}
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, "set routing-instances "+d.Get("routing_instance").(string)+
			" interface "+d.Get("name").(string))
	}
	if d.Get("security_zone").(string) != "" {
		configSet = append(configSet, "set security zones security-zone "+
			d.Get("security_zone").(string)+" interfaces "+d.Get("name").(string))
		for _, v := range d.Get("security_inbound_protocols").([]interface{}) {
			configSet = append(configSet, "set security zones security-zone "+
				d.Get("security_zone").(string)+" interfaces "+d.Get("name").(string)+
				" host-inbound-traffic protocols "+v.(string))
		}
		for _, v := range d.Get("security_inbound_services").([]interface{}) {
			configSet = append(configSet, "set security zones security-zone "+
				d.Get("security_zone").(string)+" interfaces "+d.Get("name").(string)+
				" host-inbound-traffic system-services "+v.(string))
		}
	}
	if d.Get("vlan_id").(int) != 0 {
		configSet = append(configSet, setPrefix+"vlan-id "+strconv.Itoa(d.Get("vlan_id").(int)))
	} else if intCut[0] != st0Word && intCut[1] != "0" {
		configSet = append(configSet, setPrefix+"vlan-id "+intCut[1])
	}

	return sess.configSet(configSet, jnprSess)
}

func readInterfaceLogical(interFace string, m interface{}, jnprSess *NetconfObject) (interfaceLogicalOptions, error) {
	sess := m.(*Session)
	var confRead interfaceLogicalOptions

	intConfig, err := sess.command("show configuration interfaces "+interFace+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}

	if intConfig != emptyWord {
		for _, item := range strings.Split(intConfig, "\n") {
			// exclude ethernet-switching (parameters in junos_interface_physical)
			if strings.Contains(item, "ethernet-switching") {
				continue
			}
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "family inet6"):
				if len(confRead.familyInet6) == 0 {
					confRead.familyInet6 = append(confRead.familyInet6, map[string]interface{}{
						"address":         make([]map[string]interface{}, 0),
						"filter_input":    "",
						"filter_output":   "",
						"mtu":             0,
						"rpf_check":       make([]map[string]interface{}, 0),
						"sampling_input":  false,
						"sampling_output": false,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "family inet6 address "):
					var err error
					confRead.familyInet6[0]["address"], err = readFamilyInetAddress(
						itemTrim, confRead.familyInet6[0]["address"].([]map[string]interface{}), inet6Word)
					if err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrim, "family inet6 filter input "):
					confRead.familyInet6[0]["filter_input"] = strings.TrimPrefix(itemTrim, "family inet6 filter input ")
				case strings.HasPrefix(itemTrim, "family inet6 filter output "):
					confRead.familyInet6[0]["filter_output"] = strings.TrimPrefix(itemTrim, "family inet6 filter output ")
				case strings.HasPrefix(itemTrim, "family inet6 mtu"):
					var err error
					confRead.familyInet6[0]["mtu"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "family inet6 mtu "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "family inet6 rpf-check"):
					if len(confRead.familyInet6[0]["rpf_check"].([]map[string]interface{})) == 0 {
						confRead.familyInet6[0]["rpf_check"] = append(
							confRead.familyInet6[0]["rpf_check"].([]map[string]interface{}), map[string]interface{}{
								"fail_filter": "",
								"mode_loose":  false,
							})
					}
					switch {
					case strings.HasPrefix(itemTrim, "family inet6 rpf-check fail-filter "):
						confRead.familyInet6[0]["rpf_check"].([]map[string]interface{})[0]["fail_filter"] = strings.Trim(
							strings.TrimPrefix(itemTrim, "family inet6 rpf-check fail-filter "), "\"")
					case itemTrim == "family inet6 rpf-check mode loose":
						confRead.familyInet6[0]["rpf_check"].([]map[string]interface{})[0]["mode_loose"] = true
					}
				case itemTrim == "family inet6 sampling input":
					confRead.familyInet6[0]["sampling_input"] = true
				case itemTrim == "family inet6 sampling output":
					confRead.familyInet6[0]["sampling_output"] = true
				}
			case strings.HasPrefix(itemTrim, "family inet"):
				if len(confRead.familyInet) == 0 {
					confRead.familyInet = append(confRead.familyInet, map[string]interface{}{
						"address":         make([]map[string]interface{}, 0),
						"mtu":             0,
						"filter_input":    "",
						"filter_output":   "",
						"rpf_check":       make([]map[string]interface{}, 0),
						"sampling_input":  false,
						"sampling_output": false,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "family inet address "):
					var err error
					confRead.familyInet[0]["address"], err = readFamilyInetAddress(
						itemTrim, confRead.familyInet[0]["address"].([]map[string]interface{}), inetWord)
					if err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrim, "family inet filter input "):
					confRead.familyInet[0]["filter_input"] = strings.TrimPrefix(itemTrim, "family inet filter input ")
				case strings.HasPrefix(itemTrim, "family inet filter output "):
					confRead.familyInet[0]["filter_output"] = strings.TrimPrefix(itemTrim, "family inet filter output ")
				case strings.HasPrefix(itemTrim, "family inet mtu"):
					var err error
					confRead.familyInet[0]["mtu"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "family inet mtu "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "family inet rpf-check"):
					if len(confRead.familyInet[0]["rpf_check"].([]map[string]interface{})) == 0 {
						confRead.familyInet[0]["rpf_check"] = append(
							confRead.familyInet[0]["rpf_check"].([]map[string]interface{}), map[string]interface{}{
								"fail_filter": "",
								"mode_loose":  false,
							})
					}
					switch {
					case strings.HasPrefix(itemTrim, "family inet rpf-check fail-filter "):
						confRead.familyInet[0]["rpf_check"].([]map[string]interface{})[0]["fail_filter"] = strings.Trim(
							strings.TrimPrefix(itemTrim, "family inet rpf-check fail-filter "), "\"")
					case itemTrim == "family inet rpf-check mode loose":
						confRead.familyInet[0]["rpf_check"].([]map[string]interface{})[0]["mode_loose"] = true
					}
				case itemTrim == "family inet sampling input":
					confRead.familyInet[0]["sampling_input"] = true
				case itemTrim == "family inet sampling output":
					confRead.familyInet[0]["sampling_output"] = true
				}
			case strings.HasPrefix(itemTrim, "vlan-id "):
				var err error
				confRead.vlanID, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "vlan-id "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			default:
				continue
			}
		}
	}
	routingConfig, err := sess.command("show configuration routing-instances | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	regexpInt := regexp.MustCompile(`set \S+ interface ` + interFace + `$`)
	for _, item := range strings.Split(routingConfig, "\n") {
		intMatch := regexpInt.MatchString(item)
		if intMatch {
			confRead.routingInstance = strings.TrimPrefix(strings.TrimSuffix(item, " interface "+interFace),
				"set ")

			break
		}
	}
	if checkCompatibilitySecurity(jnprSess) {
		zonesConfig, err := sess.command("show configuration security zones | display set relative", jnprSess)
		if err != nil {
			return confRead, err
		}
		regexpInts := regexp.MustCompile(`set security-zone \S+ interfaces ` + interFace + `( host-inbound-traffic .*)?$`)
		for _, item := range strings.Split(zonesConfig, "\n") {
			intMatch := regexpInts.MatchString(item)
			if intMatch {
				itemTrimSplit := strings.Split(strings.TrimPrefix(item, "set security-zone "), " ")
				confRead.securityZone = itemTrimSplit[0]
				if err := readInterfaceLogicalSecurityInboundTraffic(interFace, &confRead, m, jnprSess); err != nil {
					return confRead, err
				}

				break
			}
		}
	}

	return confRead, nil
}

func readInterfaceLogicalSecurityInboundTraffic(interFace string, confRead *interfaceLogicalOptions,
	m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	intConfig, err := sess.command("show configuration security zones security-zone "+confRead.securityZone+
		" interfaces "+interFace+" | display set relative", jnprSess)
	if err != nil {
		return err
	}

	if intConfig != emptyWord {
		for _, item := range strings.Split(intConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "host-inbound-traffic protocols "):
				confRead.securityInboundProtocols = append(confRead.securityInboundProtocols,
					strings.TrimPrefix(itemTrim, "host-inbound-traffic protocols "))
			case strings.HasPrefix(itemTrim, "host-inbound-traffic system-services "):
				confRead.securityInboundServices = append(confRead.securityInboundServices,
					strings.TrimPrefix(itemTrim, "host-inbound-traffic system-services "))
			}
		}
	}

	return nil
}

func delInterfaceLogical(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	if err := sess.configSet([]string{"delete interfaces " + d.Get("name").(string)}, jnprSess); err != nil {
		return err
	}
	if strings.HasPrefix(d.Get("name").(string), "st0.") && !d.Get("st0_also_on_destroy").(bool) {
		// interface totally delete by
		// - junos_interface_st0_unit resource
		// else there is an interface st0.x empty
		err := sess.configSet([]string{"set interfaces " + d.Get("name").(string)}, jnprSess)
		if err != nil {
			return err
		}
	}
	if d.Get("routing_instance").(string) != "" {
		if err := delRoutingInstanceInterfaceLogical(d.Get("routing_instance").(string), d, m, jnprSess); err != nil {
			return err
		}
	}
	if checkCompatibilitySecurity(jnprSess) && d.Get("security_zone").(string) != "" {
		if err := delZoneInterfaceLogical(d.Get("security_zone").(string), d, m, jnprSess); err != nil {
			return err
		}
	}

	return nil
}

func delInterfaceLogicalOpts(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	delPrefix := "delete interfaces " + d.Get("name").(string) + " "
	configSet = append(configSet,
		delPrefix+"description",
		delPrefix+"family inet",
		delPrefix+"family inet6")

	return sess.configSet(configSet, jnprSess)
}

func delZoneInterfaceLogical(zone string, d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security zones security-zone "+zone+" interfaces "+d.Get("name").(string))

	return sess.configSet(configSet, jnprSess)
}

func delRoutingInstanceInterfaceLogical(instance string, d *schema.ResourceData,
	m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete routing-instances "+instance+" interface "+d.Get("name").(string))

	return sess.configSet(configSet, jnprSess)
}

func fillInterfaceLogicalData(d *schema.ResourceData, interfaceLogicalOpt interfaceLogicalOptions) {
	if tfErr := d.Set("description", interfaceLogicalOpt.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_inet", interfaceLogicalOpt.familyInet); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_inet6", interfaceLogicalOpt.familyInet6); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", interfaceLogicalOpt.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("security_inbound_protocols", interfaceLogicalOpt.securityInboundProtocols); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("security_inbound_services", interfaceLogicalOpt.securityInboundServices); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("security_zone", interfaceLogicalOpt.securityZone); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vlan_id", interfaceLogicalOpt.vlanID); tfErr != nil {
		panic(tfErr)
	}
}

func readFamilyInetAddress(item string, inetAddress []map[string]interface{},
	family string) ([]map[string]interface{}, error) {
	var addressConfig []string
	var itemTrim string
	switch family {
	case inetWord:
		addressConfig = strings.Split(strings.TrimPrefix(item, "family inet address "), " ")
		itemTrim = strings.TrimPrefix(item, "family inet address "+addressConfig[0]+" ")
	case inet6Word:
		addressConfig = strings.Split(strings.TrimPrefix(item, "family inet6 address "), " ")
		itemTrim = strings.TrimPrefix(item, "family inet6 address "+addressConfig[0]+" ")
	}

	mAddr := genFamilyInetAddress(addressConfig[0])
	inetAddress = copyAndRemoveItemMapList("cidr_ip", mAddr, inetAddress)

	switch {
	case itemTrim == "primary":
		mAddr["primary"] = true
	case itemTrim == "preferred":
		mAddr["preferred"] = true
	case strings.HasPrefix(itemTrim, "vrrp-group ") || strings.HasPrefix(itemTrim, "vrrp-inet6-group "):
		vrrpGroup := genVRRPGroup(family)
		vrrpID, err := strconv.Atoi(addressConfig[2])
		if err != nil {
			return inetAddress, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
		itemTrimVrrp := strings.TrimPrefix(itemTrim, "vrrp-group "+strconv.Itoa(vrrpID)+" ")
		if strings.HasPrefix(itemTrim, "vrrp-inet6-group ") {
			itemTrimVrrp = strings.TrimPrefix(itemTrim, "vrrp-inet6-group "+strconv.Itoa(vrrpID)+" ")
		}
		vrrpGroup["identifier"] = vrrpID
		mAddr["vrrp_group"] = copyAndRemoveItemMapList("identifier", vrrpGroup,
			mAddr["vrrp_group"].([]map[string]interface{}))
		switch {
		case strings.HasPrefix(itemTrimVrrp, "virtual-address "):
			vrrpGroup["virtual_address"] = append(vrrpGroup["virtual_address"].([]string),
				strings.TrimPrefix(itemTrimVrrp, "virtual-address "))
		case strings.HasPrefix(itemTrimVrrp, "virtual-inet6-address "):
			vrrpGroup["virtual_address"] = append(vrrpGroup["virtual_address"].([]string),
				strings.TrimPrefix(itemTrimVrrp, "virtual-inet6-address "))
		case strings.HasPrefix(itemTrimVrrp, "virtual-link-local-address "):
			vrrpGroup["virtual_link_local_address"] = strings.TrimPrefix(itemTrimVrrp,
				"virtual-link-local-address ")
		case itemTrimVrrp == "accept-data":
			vrrpGroup["accept_data"] = true
		case strings.HasPrefix(itemTrimVrrp, "advertise-interval "):
			vrrpGroup["advertise_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrimVrrp,
				"advertise-interval "))
			if err != nil {
				return inetAddress, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrimVrrp, err)
			}
		case strings.HasPrefix(itemTrimVrrp, "inet6-advertise-interval "):
			vrrpGroup["advertise_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrimVrrp,
				"inet6-advertise-interval "))
			if err != nil {
				return inetAddress, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrimVrrp, err)
			}
		case strings.HasPrefix(itemTrimVrrp, "advertisements-threshold "):
			vrrpGroup["advertisements_threshold"], err = strconv.Atoi(strings.TrimPrefix(itemTrimVrrp,
				"advertisements-threshold "))
			if err != nil {
				return inetAddress, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrimVrrp, err)
			}
		case strings.HasPrefix(itemTrimVrrp, "authentication-key "):
			vrrpGroup["authentication_key"], err = jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrimVrrp,
				"authentication-key "), "\""))
			if err != nil {
				return inetAddress, fmt.Errorf("failed to decode authentication-key : %w", err)
			}
		case strings.HasPrefix(itemTrimVrrp, "authentication-type "):
			vrrpGroup["authentication_type"] = strings.TrimPrefix(itemTrimVrrp, "authentication-type ")
		case itemTrimVrrp == "no-accept-data":
			vrrpGroup["no_accept_data"] = true
		case itemTrimVrrp == "no-preempt":
			vrrpGroup["no_preempt"] = true
		case itemTrimVrrp == "preempt":
			vrrpGroup["preempt"] = true
		case strings.HasPrefix(itemTrimVrrp, "priority"):
			vrrpGroup["priority"], err = strconv.Atoi(strings.TrimPrefix(itemTrimVrrp, "priority "))
			if err != nil {
				return inetAddress, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrimVrrp, err)
			}
		case strings.HasPrefix(itemTrimVrrp, "track interface "):
			vrrpSlit := strings.Split(itemTrimVrrp, " ")
			cost, err := strconv.Atoi(vrrpSlit[len(vrrpSlit)-1])
			if err != nil {
				return inetAddress, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrimVrrp, err)
			}
			trackInt := map[string]interface{}{
				"interface":     vrrpSlit[2],
				"priority_cost": cost,
			}
			vrrpGroup["track_interface"] = append(vrrpGroup["track_interface"].([]map[string]interface{}), trackInt)
		case strings.HasPrefix(itemTrimVrrp, "track route "):
			vrrpSlit := strings.Split(itemTrimVrrp, " ")
			cost, err := strconv.Atoi(vrrpSlit[len(vrrpSlit)-1])
			if err != nil {
				return inetAddress, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrimVrrp, err)
			}
			trackRoute := map[string]interface{}{
				"route":            vrrpSlit[2],
				"routing_instance": vrrpSlit[4],
				"priority_cost":    cost,
			}
			vrrpGroup["track_route"] = append(vrrpGroup["track_route"].([]map[string]interface{}), trackRoute)
		}
		mAddr["vrrp_group"] = append(mAddr["vrrp_group"].([]map[string]interface{}), vrrpGroup)
	}
	inetAddress = append(inetAddress, mAddr)

	return inetAddress, nil
}

func setFamilyAddress(inetAddress interface{}, configSet []string, setPrefix string,
	family string) ([]string, error) {
	if family != inetWord && family != inet6Word {
		return configSet, fmt.Errorf("setFamilyAddress() unknown family %v", family)
	}
	inetAddressMap := inetAddress.(map[string]interface{})
	setPrefixAddress := setPrefix + "family " + family + " address " + inetAddressMap["cidr_ip"].(string)
	configSet = append(configSet, setPrefixAddress)
	if inetAddressMap["preferred"].(bool) {
		configSet = append(configSet, setPrefixAddress+" preferred")
	}
	if inetAddressMap["primary"].(bool) {
		configSet = append(configSet, setPrefixAddress+" primary")
	}
	for _, vrrpGroup := range inetAddressMap["vrrp_group"].([]interface{}) {
		if strings.Contains(setPrefix, "set interfaces st0 unit") {
			return configSet, fmt.Errorf("vrrp not available on st0")
		}
		vrrpGroupMap := vrrpGroup.(map[string]interface{})
		if vrrpGroupMap["no_preempt"].(bool) && vrrpGroupMap["preempt"].(bool) {
			return configSet, fmt.Errorf("ConflictsWith no_preempt and preempt")
		}
		if vrrpGroupMap["no_accept_data"].(bool) && vrrpGroupMap["accept_data"].(bool) {
			return configSet, fmt.Errorf("ConflictsWith no_accept_data and accept_data")
		}
		var setNameAddVrrp string
		switch family {
		case inetWord:
			setNameAddVrrp = setPrefixAddress + " vrrp-group " + strconv.Itoa(vrrpGroupMap["identifier"].(int))
			for _, ip := range vrrpGroupMap["virtual_address"].([]interface{}) {
				_, errs := validation.IsIPAddress(ip, "virtual_address")
				if len(errs) > 0 {
					return configSet, errs[0]
				}
				configSet = append(configSet, setNameAddVrrp+" virtual-address "+ip.(string))
			}
			if vrrpGroupMap["advertise_interval"].(int) != 0 {
				configSet = append(configSet, setNameAddVrrp+" advertise-interval "+
					strconv.Itoa(vrrpGroupMap["advertise_interval"].(int)))
			}
			if vrrpGroupMap["authentication_key"].(string) != "" {
				configSet = append(configSet, setNameAddVrrp+" authentication-key \""+
					vrrpGroupMap["authentication_key"].(string)+"\"")
			}
			if vrrpGroupMap["authentication_type"].(string) != "" {
				configSet = append(configSet, setNameAddVrrp+" authentication-type "+
					vrrpGroupMap["authentication_type"].(string))
			}
		case inet6Word:
			setNameAddVrrp = setPrefixAddress + " vrrp-inet6-group " + strconv.Itoa(vrrpGroupMap["identifier"].(int))
			for _, ip := range vrrpGroupMap["virtual_address"].([]interface{}) {
				_, errs := validation.IsIPAddress(ip, "virtual_address")
				if len(errs) > 0 {
					return configSet, errs[0]
				}
				configSet = append(configSet, setNameAddVrrp+" virtual-inet6-address "+ip.(string))
			}
			configSet = append(configSet, setNameAddVrrp+" virtual-link-local-address "+
				vrrpGroupMap["virtual_link_local_address"].(string))
			if vrrpGroupMap["advertise_interval"].(int) != 0 {
				configSet = append(configSet, setNameAddVrrp+" inet6-advertise-interval "+
					strconv.Itoa(vrrpGroupMap["advertise_interval"].(int)))
			}
		}
		if vrrpGroupMap["accept_data"].(bool) {
			configSet = append(configSet, setNameAddVrrp+" accept-data")
		}
		if vrrpGroupMap["advertisements_threshold"].(int) != 0 {
			configSet = append(configSet, setNameAddVrrp+" advertisements-threshold "+
				strconv.Itoa(vrrpGroupMap["advertisements_threshold"].(int)))
		}
		if vrrpGroupMap["no_accept_data"].(bool) {
			configSet = append(configSet, setNameAddVrrp+" no-accept-data")
		}
		if vrrpGroupMap["no_preempt"].(bool) {
			configSet = append(configSet, setNameAddVrrp+" no-preempt")
		}
		if vrrpGroupMap["preempt"].(bool) {
			configSet = append(configSet, setNameAddVrrp+" preempt")
		}
		if vrrpGroupMap["priority"].(int) != 0 {
			configSet = append(configSet, setNameAddVrrp+" priority "+strconv.Itoa(vrrpGroupMap["priority"].(int)))
		}
		for _, trackInterface := range vrrpGroupMap["track_interface"].([]interface{}) {
			trackInterfaceMap := trackInterface.(map[string]interface{})
			configSet = append(configSet, setNameAddVrrp+" track interface "+trackInterfaceMap["interface"].(string)+
				" priority-cost "+strconv.Itoa(trackInterfaceMap["priority_cost"].(int)))
		}
		for _, trackRoute := range vrrpGroupMap["track_route"].([]interface{}) {
			trackRouteMap := trackRoute.(map[string]interface{})
			configSet = append(configSet, setNameAddVrrp+" track route "+trackRouteMap["route"].(string)+
				" routing-instance "+trackRouteMap["routing_instance"].(string)+
				" priority-cost "+strconv.Itoa(trackRouteMap["priority_cost"].(int)))
		}
	}

	return configSet, nil
}

func genFamilyInetAddress(address string) map[string]interface{} {
	return map[string]interface{}{
		"cidr_ip":    address,
		"primary":    false,
		"preferred":  false,
		"vrrp_group": make([]map[string]interface{}, 0),
	}
}

func genVRRPGroup(family string) map[string]interface{} {
	vrrpGroup := map[string]interface{}{
		"identifier":               0,
		"virtual_address":          make([]string, 0),
		"accept_data":              false,
		"advertise_interval":       0,
		"advertisements_threshold": 0,
		"no_accept_data":           false,
		"no_preempt":               false,
		"preempt":                  false,
		"priority":                 0,
		"track_interface":          make([]map[string]interface{}, 0),
		"track_route":              make([]map[string]interface{}, 0),
	}
	if family == inetWord {
		vrrpGroup["authentication_key"] = ""
		vrrpGroup["authentication_type"] = ""
	}
	if family == inet6Word {
		vrrpGroup["virtual_link_local_address"] = ""
	}

	return vrrpGroup
}
