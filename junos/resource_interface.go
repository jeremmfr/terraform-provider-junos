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

type interfaceOptions struct {
	vlanTagging       bool
	inet              bool
	inet6             bool
	trunk             bool
	vlanNative        int
	aeMinLink         int
	inetMtu           int
	inet6Mtu          int
	vlanTaggingID     int
	inetFilterInput   string
	inetFilterOutput  string
	inet6FilterInput  string
	inet6FilterOutput string
	description       string
	v8023ad           string
	aeLacp            string
	aeLinkSpeed       string
	securityZones     string
	routingInstances  string
	vlanMembers       []string
	inetRpfCheck      []map[string]interface{}
	inet6RpfCheck     []map[string]interface{}
	inetAddress       []map[string]interface{}
	inet6Address      []map[string]interface{}
}

func resourceInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInterfaceCreate,
		ReadContext:   resourceInterfaceRead,
		UpdateContext: resourceInterfaceUpdate,
		DeleteContext: resourceInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceInterfaceImport,
		},
		DeprecationMessage: "use junos_interface_physical or junos_interface_logical resource instead",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if strings.Count(value, ".") > 1 {
						errors = append(errors, fmt.Errorf(
							"%q in %q cannot have more of 1 dot", value, k))
					}

					return
				},
			},
			"complete_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan_tagging": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vlan_tagging_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 4094),
			},
			"inet": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"inet6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"inet_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateIPMaskFunc(),
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
			"inet6_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateIPMaskFunc(),
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
			"inet_mtu": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(500, 9192),
			},
			"inet6_mtu": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(500, 9192),
			},
			"inet_filter_input": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"inet_filter_output": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"inet6_filter_input": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"inet6_filter_output": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"inet_rpf_check": {
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
			"inet6_rpf_check": {
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
			"ether802_3ad": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !strings.HasPrefix(value, "ae") {
						errors = append(errors, fmt.Errorf(
							"%q in %q isn't an ae interface", value, k))
					}

					return
				},
			},
			"trunk": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vlan_members": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vlan_native": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 4094),
			},
			"ae_lacp": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.StringInSlice([]string{"active", "passive"}, false),
			},
			"ae_link_speed": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"100m", "1g", "8g", "10g", "40g", "50g", "80g", "100g"}, false),
			},
			"ae_minimum_links": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"security_zone": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
		},
	}
}

func resourceInterfaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	intExists, err := checkInterfaceExistsOld(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return diag.FromErr(err)
	}
	sess.configLock(jnprSess)
	if intExists {
		ncInt, emptyInt, err := checkInterfaceNC(d.Get("name").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)

			return diag.FromErr(err)
		}
		if !ncInt && !emptyInt {
			return diag.FromErr(fmt.Errorf("interface %s already configured", d.Get("name").(string)))
		}
		if sess.junosGroupIntDel != "" {
			err = delInterfaceElement("apply-groups "+sess.junosGroupIntDel, d, m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)

				return diag.FromErr(err)
			}
		} else {
			err = delInterfaceElement("disable", d, m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)

				return diag.FromErr(err)
			}
			err = delInterfaceElement("description", d, m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)

				return diag.FromErr(err)
			}
		}
	}
	if d.Get("security_zone").(string) != "" {
		if !checkCompatibilitySecurity(jnprSess) {
			sess.configClear(jnprSess)

			return diag.FromErr(fmt.Errorf("security zone not compatible with Junos device %s",
				jnprSess.SystemInformation.HardwareModel))
		}
		zonesExists, err := checkSecurityZonesExists(d.Get("security_zone").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)

			return diag.FromErr(err)
		}
		if !zonesExists {
			sess.configClear(jnprSess)

			return diag.FromErr(fmt.Errorf("security zones %v doesn't exist", d.Get("security_zone").(string)))
		}
	}
	if d.Get("routing_instance").(string) != "" {
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
	if err := setInterface(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	intExists, err = checkInterfaceExistsOld(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if intExists {
		ncInt, _, err := checkInterfaceNC(d.Get("name").(string), m, jnprSess)
		if err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		if ncInt {
			return append(diagWarns,
				diag.FromErr(fmt.Errorf("interface %v exists (because is a physical or internal default interface)"+
					" but always disable after commit => check your config", d.Get("name").(string)))...)
		}
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("interface %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceInterfaceReadWJnprSess(d, m, jnprSess)...)
}
func resourceInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceInterfaceReadWJnprSess(d, m, jnprSess)
}
func resourceInterfaceReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	intExists, err := checkInterfaceExistsOld(d.Get("name").(string), m, jnprSess)
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	if !intExists {
		d.SetId("")
		mutex.Unlock()

		return nil
	}
	ncInt, _, err := checkInterfaceNC(d.Get("name").(string), m, jnprSess)
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	if ncInt {
		d.SetId("")
		mutex.Unlock()

		return nil
	}
	interfaceOpt, err := readInterface(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillInterfaceData(d, interfaceOpt)

	return nil
}
func resourceInterfaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delInterfaceOpts(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if d.HasChange("ether802_3ad") {
		oAE, nAE := d.GetChange("ether802_3ad")
		if oAE.(string) != "" {
			newAE := "ae-1" // nolint: goconst
			if nAE.(string) != "" {
				newAE = nAE.(string)
			}
			lastAEchild, err := aggregatedLastChild(oAE.(string), d.Get("name").(string), m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)

				return diag.FromErr(err)
			}
			if lastAEchild {
				aggregatedCount, err := aggregatedCountSearchMax(newAE, oAE.(string), d.Get("name").(string), m, jnprSess)
				if err != nil {
					sess.configClear(jnprSess)

					return diag.FromErr(err)
				}
				if aggregatedCount == "0" {
					err = sess.configSet([]string{"delete chassis aggregated-devices ethernet device-count"}, jnprSess)
					if err != nil {
						sess.configClear(jnprSess)

						return diag.FromErr(err)
					}
					oAEintNC, oAEintEmpty, err := checkInterfaceNC(oAE.(string), m, jnprSess)
					if err != nil {
						sess.configClear(jnprSess)

						return diag.FromErr(err)
					}
					if oAEintNC || oAEintEmpty {
						err = sess.configSet([]string{"delete interfaces " + oAE.(string)}, jnprSess)
						if err != nil {
							sess.configClear(jnprSess)

							return diag.FromErr(err)
						}
					}
				} else {
					oldAEInt, err := strconv.Atoi(strings.TrimPrefix(oAE.(string), "ae"))
					if err != nil {
						sess.configClear(jnprSess)

						return diag.FromErr(err)
					}
					aggregatedCountInt, err := strconv.Atoi(aggregatedCount)
					if err != nil {
						sess.configClear(jnprSess)

						return diag.FromErr(err)
					}
					if aggregatedCountInt < oldAEInt+1 {
						oAEintNC, oAEintEmpty, err := checkInterfaceNC(oAE.(string), m, jnprSess)
						if err != nil {
							sess.configClear(jnprSess)

							return diag.FromErr(err)
						}
						if oAEintNC || oAEintEmpty {
							err = sess.configSet([]string{"delete interfaces " + oAE.(string)}, jnprSess)
							if err != nil {
								sess.configClear(jnprSess)

								return diag.FromErr(err)
							}
						}
					}
				}
			}
		}
	}
	if d.HasChange("security_zone") {
		oSecurityZone, nSecurityZone := d.GetChange("security_zone")
		if nSecurityZone.(string) != "" {
			if !checkCompatibilitySecurity(jnprSess) {
				sess.configClear(jnprSess)

				return diag.FromErr(fmt.Errorf("security zone not compatible with Junos device %s",
					jnprSess.SystemInformation.HardwareModel))
			}
			zonesExists, err := checkSecurityZonesExists(nSecurityZone.(string), m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)

				return diag.FromErr(err)
			}
			if !zonesExists {
				sess.configClear(jnprSess)

				return diag.FromErr(fmt.Errorf("security zones %v doesn't exist", nSecurityZone.(string)))
			}
		}
		if oSecurityZone.(string) != "" {
			err = delZoneInterface(oSecurityZone.(string), d, m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)

				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("routing_instance") {
		oRoutingInstance, nRoutingInstance := d.GetChange("routing_instance")
		if nRoutingInstance.(string) != "" {
			instanceExists, err := checkRoutingInstanceExists(nRoutingInstance.(string), m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)

				return diag.FromErr(err)
			}
			if !instanceExists {
				sess.configClear(jnprSess)

				return diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", nRoutingInstance.(string)))
			}
		}
		if oRoutingInstance.(string) != "" {
			err = delRoutingInstanceInterface(oRoutingInstance.(string), d, m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)

				return diag.FromErr(err)
			}
		}
	}
	if err := setInterface(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceInterfaceReadWJnprSess(d, m, jnprSess)...)
}
func resourceInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delInterface(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !d.Get("complete_destroy").(bool) {
		intExists, err := checkInterfaceExistsOld(d.Get("name").(string), m, jnprSess)
		if err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		if intExists {
			err = addInterfaceNC(d.Get("name").(string), m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)

				return append(diagWarns, diag.FromErr(err)...)
			}
			_, err = sess.commitConf("disable(NC) resource junos_interface", jnprSess)
			if err != nil {
				sess.configClear(jnprSess)

				return append(diagWarns, diag.FromErr(err)...)
			}
		}
	}

	return diagWarns
}
func resourceInterfaceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	intExists, err := checkInterfaceExistsOld(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !intExists {
		return nil, fmt.Errorf("don't find interface with id '%v' (id must be <name>)", d.Id())
	}
	interfaceOpt, err := readInterface(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if tfErr := d.Set("name", d.Id()); tfErr != nil {
		panic(tfErr)
	}
	fillInterfaceData(d, interfaceOpt)

	result[0] = d

	return result, nil
}

func checkInterfaceNC(interFace string, m interface{}, jnprSess *NetconfObject) (
	ncInt bool, emtyInt bool, errFunc error) {
	sess := m.(*Session)
	intConfig, err := sess.command("show configuration interfaces "+interFace+" | display set relative", jnprSess)
	if err != nil {
		return false, false, err
	}
	intConfigLines := make([]string, 0)
	// remove unused lines
	for _, item := range strings.Split(intConfig, "\n") {
		// show parameters root on interface exclude unit parameters (except ethernet-switching)
		if !strings.Contains(interFace, ".") && strings.HasPrefix(item, "set unit") &&
			!strings.Contains(item, "ethernet-switching") {
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
		return false, true, nil
	}
	intConfig = strings.Join(intConfigLines, "\n")
	if sess.junosGroupIntDel != "" {
		if intConfig == "set apply-groups "+sess.junosGroupIntDel {
			return true, false, nil
		}
	}
	if intConfig == "set description NC\nset disable" || // nolint: goconst
		intConfig == "set disable\nset description NC" { // nolint: goconst
		return true, false, nil
	}
	if intConfig == setLineStart ||
		intConfig == emptyWord {
		return false, true, nil
	}

	return false, false, nil
}

func addInterfaceNC(interFace string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	intCut := make([]string, 0, 2)
	var setName string
	var err error
	if strings.Contains(interFace, ".") {
		intCut = strings.Split(interFace, ".")
	} else {
		intCut = append(intCut, interFace)
	}
	switch len(intCut) {
	case 2:
		setName = intCut[0] + " unit " + intCut[1]
	case 1:
		setName = intCut[0]
	default:
		return fmt.Errorf("the name %s contains too dots", interFace)
	}
	if intCut[0] == st0Word || sess.junosGroupIntDel == "" {
		err = sess.configSet([]string{"set interfaces " + setName + " disable description NC"}, jnprSess)
	} else {
		err = sess.configSet([]string{"set interfaces " + setName +
			" apply-groups " + sess.junosGroupIntDel}, jnprSess)
	}
	if err != nil {
		return err
	}

	return nil
}

func checkInterfaceExistsOld(interFace string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	rpcIntName := "<get-interface-information><interface-name>" + interFace +
		"</interface-name></get-interface-information>"
	reply, err := sess.commandXML(rpcIntName, jnprSess)
	if err != nil {
		return false, err
	}
	if strings.Contains(reply, " not found\n") {
		return false, nil
	}

	return true, nil
}
func setInterface(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	var setName string
	intCut := make([]string, 0, 2)
	configSet := make([]string, 0)
	if strings.Contains(d.Get("name").(string), ".") {
		intCut = strings.Split(d.Get("name").(string), ".")
	} else {
		intCut = append(intCut, d.Get("name").(string))
	}
	switch len(intCut) {
	case 2:
		setName = intCut[0] + " unit " + intCut[1]
	case 1:
		setName = intCut[0]
	default:
		return fmt.Errorf("the name %s contains too dots", d.Get("name").(string))
	}
	if err := checkResourceInterfaceConfigAndName(len(intCut), d); err != nil {
		return err
	}
	setPrefix := "set interfaces " + setName + " "
	configSet = append(configSet, setPrefix)
	if d.Get("description").(string) != "" {
		configSet = append(configSet, setPrefix+"description \""+d.Get("description").(string)+"\"")
	}
	if d.Get("vlan_tagging").(bool) {
		configSet = append(configSet, setPrefix+"vlan-tagging")
	}
	if d.Get("vlan_tagging_id").(int) != 0 {
		configSet = append(configSet, setPrefix+"vlan-id "+strconv.Itoa(d.Get("vlan_tagging_id").(int)))
	} else if len(intCut) == 2 && intCut[0] != st0Word && intCut[1] != "0" {
		configSet = append(configSet, setPrefix+"vlan-id "+intCut[1])
	}
	if d.Get("inet").(bool) {
		configSet = append(configSet, setPrefix+"family inet")
	}
	if d.Get("inet6").(bool) {
		configSet = append(configSet, setPrefix+"family inet6")
	}
	for _, address := range d.Get("inet_address").([]interface{}) {
		var err error
		configSet, err = setFamilyAddressOld(address, intCut, configSet, setName, inetWord)
		if err != nil {
			return err
		}
	}
	for _, address := range d.Get("inet6_address").([]interface{}) {
		var err error
		configSet, err = setFamilyAddressOld(address, intCut, configSet, setName, inet6Word)
		if err != nil {
			return err
		}
	}
	if d.Get("inet_mtu").(int) > 0 {
		configSet = append(configSet, setPrefix+"family inet mtu "+
			strconv.Itoa(d.Get("inet_mtu").(int)))
	}
	if d.Get("inet6_mtu").(int) > 0 {
		configSet = append(configSet, setPrefix+"family inet6 mtu "+
			strconv.Itoa(d.Get("inet6_mtu").(int)))
	}
	for _, v := range d.Get("inet_rpf_check").([]interface{}) {
		configSet = append(configSet, setPrefix+"family inet rpf-check")
		if v != nil {
			rpfCheck := v.(map[string]interface{})
			if rpfCheck["fail_filter"].(string) != "" {
				configSet = append(configSet, setPrefix+"family inet rpf-check fail-filter "+
					"\""+rpfCheck["fail_filter"].(string)+"\"")
			}
			if rpfCheck["mode_loose"].(bool) {
				configSet = append(configSet, setPrefix+"family inet rpf-check mode loose ")
			}
		}
	}
	for _, v := range d.Get("inet6_rpf_check").([]interface{}) {
		configSet = append(configSet, setPrefix+"family inet6 rpf-check")
		if v != nil {
			rpfCheck := v.(map[string]interface{})
			if rpfCheck["fail_filter"].(string) != "" {
				configSet = append(configSet, setPrefix+"family inet6 rpf-check fail-filter "+
					"\""+rpfCheck["fail_filter"].(string)+"\"")
			}
			if rpfCheck["mode_loose"].(bool) {
				configSet = append(configSet, setPrefix+"family inet6 rpf-check mode loose ")
			}
		}
	}
	if d.Get("inet_filter_input").(string) != "" {
		configSet = append(configSet, setPrefix+"family inet filter input "+
			d.Get("inet_filter_input").(string))
	}
	if d.Get("inet_filter_output").(string) != "" {
		configSet = append(configSet, setPrefix+"family inet filter output "+
			d.Get("inet_filter_output").(string))
	}
	if d.Get("inet6_filter_input").(string) != "" {
		configSet = append(configSet, setPrefix+"family inet6 filter input "+
			d.Get("inet6_filter_input").(string))
	}
	if d.Get("inet6_filter_output").(string) != "" {
		configSet = append(configSet, setPrefix+"family inet6 filter output "+
			d.Get("inet6_filter_output").(string))
	}
	if d.Get("ether802_3ad").(string) != "" {
		configSet = append(configSet, setPrefix+"ether-options 802.3ad "+
			d.Get("ether802_3ad").(string))
		configSet = append(configSet, setPrefix+"gigether-options 802.3ad "+
			d.Get("ether802_3ad").(string))
		oldAE := "ae-1"
		if d.HasChange("ether802_3ad") {
			oldAEtf, _ := d.GetChange("ether802_3ad")
			if oldAEtf.(string) != "" {
				oldAE = oldAEtf.(string)
			}
		}
		aggregatedCount, err := aggregatedCountSearchMax(d.Get("ether802_3ad").(string), oldAE,
			d.Get("name").(string), m, jnprSess)
		if err != nil {
			return err
		}
		configSet = append(configSet, "set chassis aggregated-devices ethernet device-count "+aggregatedCount)
	}
	if d.Get("trunk").(bool) {
		configSet = append(configSet, setPrefix+"unit 0 family ethernet-switching interface-mode trunk")
	}
	for _, v := range d.Get("vlan_members").([]interface{}) {
		configSet = append(configSet, setPrefix+
			"unit 0 family ethernet-switching vlan members "+v.(string))
	}
	if d.Get("vlan_native").(int) != 0 {
		configSet = append(configSet, setPrefix+"native-vlan-id "+strconv.Itoa(d.Get("vlan_native").(int)))
	}
	if d.Get("ae_lacp").(string) != "" {
		if !strings.Contains(intCut[0], "ae") {
			return fmt.Errorf("ae_lacp invalid for this interface")
		}
		configSet = append(configSet, setPrefix+
			"aggregated-ether-options lacp "+d.Get("ae_lacp").(string))
	}
	if d.Get("ae_link_speed").(string) != "" {
		if !strings.Contains(intCut[0], "ae") {
			return fmt.Errorf("ae_link_speed invalid for this interface")
		}
		configSet = append(configSet, setPrefix+
			"aggregated-ether-options link-speed "+d.Get("ae_link_speed").(string))
	}
	if d.Get("ae_minimum_links").(int) > 0 {
		if !strings.Contains(intCut[0], "ae") {
			return fmt.Errorf("ae_minimum_links invalid for this interface")
		}
		configSet = append(configSet, setPrefix+
			"aggregated-ether-options minimum-links "+strconv.Itoa(d.Get("ae_minimum_links").(int)))
	}
	if checkCompatibilitySecurity(jnprSess) && d.Get("security_zone").(string) != "" {
		configSet = append(configSet, "set security zones security-zone "+
			d.Get("security_zone").(string)+" interfaces "+d.Get("name").(string))
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, "set routing-instances "+d.Get("routing_instance").(string)+
			" interface "+d.Get("name").(string))
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readInterface(interFace string, m interface{}, jnprSess *NetconfObject) (interfaceOptions, error) {
	sess := m.(*Session)
	var confRead interfaceOptions

	intConfig, err := sess.command("show configuration interfaces "+interFace+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	inetAddress := make([]map[string]interface{}, 0)
	inet6Address := make([]map[string]interface{}, 0)

	if intConfig != emptyWord {
		for _, item := range strings.Split(intConfig, "\n") {
			if !strings.Contains(interFace, ".") && strings.Contains(item, " unit ") &&
				!strings.Contains(item, "ethernet-switching") {
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

			case itemTrim == "vlan-tagging":
				confRead.vlanTagging = true
			case strings.HasPrefix(itemTrim, "vlan-id "):
				var err error
				confRead.vlanTaggingID, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "vlan-id "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "family inet6"):
				confRead.inet6 = true
				switch {
				case strings.HasPrefix(itemTrim, "family inet6 address "):
					inet6Address, err = fillFamilyInetAddressOld(itemTrim, inet6Address, inet6Word)
					if err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrim, "family inet6 mtu"):
					confRead.inet6Mtu, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "family inet6 mtu "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "family inet6 filter input "):
					confRead.inet6FilterInput = strings.TrimPrefix(itemTrim, "family inet6 filter input ")
				case strings.HasPrefix(itemTrim, "family inet6 filter output "):
					confRead.inet6FilterOutput = strings.TrimPrefix(itemTrim, "family inet6 filter output ")
				case strings.HasPrefix(itemTrim, "family inet6 rpf-check"):
					if len(confRead.inet6RpfCheck) == 0 {
						confRead.inet6RpfCheck = append(confRead.inet6RpfCheck, map[string]interface{}{
							"fail_filter": "",
							"mode_loose":  false,
						})
					}
					switch {
					case strings.HasPrefix(itemTrim, "family inet6 rpf-check fail-filter "):
						confRead.inet6RpfCheck[0]["fail_filter"] = strings.Trim(
							strings.TrimPrefix(itemTrim, "family inet6 rpf-check fail-filter "), "\"")
					case itemTrim == "family inet6 rpf-check mode loose":
						confRead.inet6RpfCheck[0]["mode_loose"] = true
					}
				}
			case strings.HasPrefix(itemTrim, "family inet"):
				confRead.inet = true
				switch {
				case strings.HasPrefix(itemTrim, "family inet address "):
					inetAddress, err = fillFamilyInetAddressOld(itemTrim, inetAddress, inetWord)
					if err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrim, "family inet mtu "):
					confRead.inetMtu, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "family inet mtu "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "family inet filter input "):
					confRead.inetFilterInput = strings.TrimPrefix(itemTrim, "family inet filter input ")
				case strings.HasPrefix(itemTrim, "family inet filter output "):
					confRead.inetFilterOutput = strings.TrimPrefix(itemTrim, "family inet filter output ")
				case strings.HasPrefix(itemTrim, "family inet rpf-check"):
					if len(confRead.inetRpfCheck) == 0 {
						confRead.inetRpfCheck = append(confRead.inetRpfCheck, map[string]interface{}{
							"fail_filter": "",
							"mode_loose":  false,
						})
					}
					switch {
					case strings.HasPrefix(itemTrim, "family inet rpf-check fail-filter "):
						confRead.inetRpfCheck[0]["fail_filter"] = strings.Trim(
							strings.TrimPrefix(itemTrim, "family inet rpf-check fail-filter "), "\"")
					case itemTrim == "family inet rpf-check mode loose":
						confRead.inetRpfCheck[0]["mode_loose"] = true
					}
				}
			case strings.HasPrefix(itemTrim, "ether-options 802.3ad "):
				confRead.v8023ad = strings.TrimPrefix(itemTrim, "ether-options 802.3ad ")
			case strings.HasPrefix(itemTrim, "gigether-options 802.3ad "):
				confRead.v8023ad = strings.TrimPrefix(itemTrim, "gigether-options 802.3ad ")
			case itemTrim == "unit 0 family ethernet-switching interface-mode trunk":
				confRead.trunk = true
			case strings.HasPrefix(itemTrim, "unit 0 family ethernet-switching vlan members"):
				confRead.vlanMembers = append(confRead.vlanMembers, strings.TrimPrefix(itemTrim,
					"unit 0 family ethernet-switching vlan members "))
			case strings.HasPrefix(itemTrim, "native-vlan-id"):
				confRead.vlanNative, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "native-vlan-id "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "aggregated-ether-options lacp "):
				confRead.aeLacp = strings.TrimPrefix(itemTrim, "aggregated-ether-options lacp ")
			case strings.HasPrefix(itemTrim, "aggregated-ether-options link-speed "):
				confRead.aeLinkSpeed = strings.TrimPrefix(itemTrim, "aggregated-ether-options link-speed ")
			case strings.HasPrefix(itemTrim, "aggregated-ether-options minimum-links "):
				confRead.aeMinLink, err = strconv.Atoi(strings.TrimPrefix(itemTrim,
					"aggregated-ether-options minimum-links "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			default:
				continue
			}
		}
		confRead.inetAddress = inetAddress
		confRead.inet6Address = inet6Address
	}
	if checkCompatibilitySecurity(jnprSess) {
		zonesConfig, err := sess.command("show configuration security zones | display set relative", jnprSess)
		if err != nil {
			return confRead, err
		}
		regexpInts := regexp.MustCompile(`set security-zone \S+ interfaces ` + interFace + `$`)
		for _, item := range strings.Split(zonesConfig, "\n") {
			intMatch := regexpInts.MatchString(item)
			if intMatch {
				confRead.securityZones = strings.TrimPrefix(strings.TrimSuffix(item, " interfaces "+interFace),
					"set security-zone ")

				break
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
			confRead.routingInstances = strings.TrimPrefix(strings.TrimSuffix(item, " interface "+interFace),
				"set ")

			break
		}
	}

	return confRead, nil
}
func delInterface(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	intCut := make([]string, 0, 2)
	var setName string
	if strings.Contains(d.Get("name").(string), ".") {
		intCut = strings.Split(d.Get("name").(string), ".")
	} else {
		intCut = append(intCut, d.Get("name").(string))
	}
	switch len(intCut) {
	case 2:
		setName = intCut[0] + " unit " + intCut[1]
	case 1:
		setName = intCut[0]
		err := checkInterfaceContainsUnit(setName, m, jnprSess)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("the name %s contains too dots", d.Get("name").(string))
	}

	if err := sess.configSet([]string{"delete interfaces " + setName}, jnprSess); err != nil {
		return err
	}
	if strings.Contains(d.Get("name").(string), "st0.") && !d.Get("complete_destroy").(bool) {
		// interface totally delete by
		// - junos_security_ipsec_vpn resource with the bind_interface_auto argument (deprecated)
		// or by
		// - junos_interface_st0_unit resource
		// else there is an interface st0.x empty
		err := sess.configSet([]string{"set interfaces " + setName}, jnprSess)
		if err != nil {
			return err
		}
	}
	if d.Get("ether802_3ad").(string) != "" {
		lastAEchild, err := aggregatedLastChild(d.Get("ether802_3ad").(string), d.Get("name").(string), m, jnprSess)
		if err != nil {
			return err
		}
		if lastAEchild {
			aggregatedCount, err := aggregatedCountSearchMax("ae-1", d.Get("ether802_3ad").(string),
				d.Get("name").(string), m, jnprSess)
			if err != nil {
				return err
			}
			if aggregatedCount == "0" {
				err = sess.configSet([]string{"delete chassis aggregated-devices ethernet device-count"}, jnprSess)
				if err != nil {
					return err
				}
			} else {
				err = sess.configSet([]string{"set chassis aggregated-devices ethernet device-count " +
					aggregatedCount}, jnprSess)
				if err != nil {
					return err
				}
			}
			aeInt, err := strconv.Atoi(strings.TrimPrefix(d.Get("ether802_3ad").(string), "ae"))
			if err != nil {
				return fmt.Errorf("failed to convert AE id of ether802_3ad argument '%s' in integer : %w",
					d.Get("ether802_3ad").(string), err)
			}
			aggregatedCountInt, err := strconv.Atoi(aggregatedCount)
			if err != nil {
				return fmt.Errorf("failed to convert internal variable aggregatedCountInt in integer : %w", err)
			}
			if aggregatedCountInt < aeInt+1 {
				oAEintNC, oAEintEmpty, err := checkInterfaceNC(d.Get("ether802_3ad").(string), m, jnprSess)
				if err != nil {
					return err
				}
				if oAEintNC || oAEintEmpty {
					err = sess.configSet([]string{"delete interfaces " + d.Get("ether802_3ad").(string)}, jnprSess)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	if checkCompatibilitySecurity(jnprSess) && d.Get("security_zone").(string) != "" {
		if err := delZoneInterface(d.Get("security_zone").(string), d, m, jnprSess); err != nil {
			return err
		}
	}
	if d.Get("routing_instance").(string) != "" {
		if err := delRoutingInstanceInterface(d.Get("routing_instance").(string), d, m, jnprSess); err != nil {
			return err
		}
	}

	return nil
}

func checkInterfaceContainsUnit(interFace string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	intConfig, err := sess.command("show configuration interfaces "+interFace+" | display set relative", jnprSess)
	if err != nil {
		return err
	}
	for _, item := range strings.Split(intConfig, "\n") {
		if strings.Contains(item, "<configuration-output>") {
			continue
		}
		if strings.Contains(item, "</configuration-output>") {
			break
		}
		if strings.HasPrefix(item, "set unit") {
			if strings.Contains(item, "ethernet-switching") {
				continue
			}

			return fmt.Errorf("interface %s is used for other son unit interface", interFace)
		}
	}

	return nil
}
func delInterfaceElement(element string, d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	intCut := make([]string, 0, 2)
	var setName string
	configSet := make([]string, 0, 1)
	if strings.Contains(d.Get("name").(string), ".") {
		intCut = strings.Split(d.Get("name").(string), ".")
	} else {
		intCut = append(intCut, d.Get("name").(string))
	}
	switch len(intCut) {
	case 2:
		setName = intCut[0] + " unit " + intCut[1]
	case 1:
		setName = intCut[0]
	default:
		return fmt.Errorf("the name %s contains too dots", d.Get("name").(string))
	}
	configSet = append(configSet, "delete interfaces "+setName+" "+element)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func delInterfaceOpts(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	intCut := make([]string, 0, 2)
	var setName string
	configSet := make([]string, 0, 1)
	if strings.Contains(d.Get("name").(string), ".") {
		intCut = strings.Split(d.Get("name").(string), ".")
	} else {
		intCut = append(intCut, d.Get("name").(string))
	}
	switch len(intCut) {
	case 2:
		setName = intCut[0] + " unit " + intCut[1]
	case 1:
		setName = intCut[0]
	default:
		return fmt.Errorf("the name %s contains too dots", d.Get("name").(string))
	}
	delPrefix := "delete interfaces " + setName + " "
	configSet = append(configSet,
		delPrefix+"vlan-tagging",
		delPrefix+"family inet",
		delPrefix+"family inet6",
		delPrefix+"ether-options 802.3ad",
		delPrefix+"gigether-options 802.3ad",
		delPrefix+"unit 0 family ethernet-switching interface-mode",
		delPrefix+"unit 0 family ethernet-switching vlan members",
		delPrefix+"native-vlan-id",
		delPrefix+"aggregated-ether-options")
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func delZoneInterface(zone string, d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security zones security-zone "+zone+" interfaces "+d.Get("name").(string))
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func delRoutingInstanceInterface(instance string, d *schema.ResourceData,
	m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete routing-instances "+instance+" interface "+d.Get("name").(string))
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func fillInterfaceData(d *schema.ResourceData, interfaceOpt interfaceOptions) {
	if tfErr := d.Set("description", interfaceOpt.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vlan_tagging", interfaceOpt.vlanTagging); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vlan_tagging_id", interfaceOpt.vlanTaggingID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet", interfaceOpt.inet); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet6", interfaceOpt.inet6); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet_address", interfaceOpt.inetAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet6_address", interfaceOpt.inet6Address); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet_mtu", interfaceOpt.inetMtu); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet6_mtu", interfaceOpt.inet6Mtu); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet_filter_input", interfaceOpt.inetFilterInput); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet_filter_output", interfaceOpt.inetFilterOutput); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet6_filter_input", interfaceOpt.inet6FilterInput); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet6_filter_output", interfaceOpt.inet6FilterOutput); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet_rpf_check", interfaceOpt.inetRpfCheck); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet6_rpf_check", interfaceOpt.inet6RpfCheck); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ether802_3ad", interfaceOpt.v8023ad); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("trunk", interfaceOpt.trunk); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vlan_members", interfaceOpt.vlanMembers); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vlan_native", interfaceOpt.vlanNative); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ae_lacp", interfaceOpt.aeLacp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ae_link_speed", interfaceOpt.aeLinkSpeed); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ae_minimum_links", interfaceOpt.aeMinLink); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("security_zone", interfaceOpt.securityZones); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", interfaceOpt.routingInstances); tfErr != nil {
		panic(tfErr)
	}
}
func fillFamilyInetAddressOld(item string, inetAddress []map[string]interface{},
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

	m := genFamilyInetAddressOld(addressConfig[0])
	m, inetAddress = copyAndRemoveItemMapList("address", false, m, inetAddress)

	if strings.HasPrefix(itemTrim, "vrrp-group ") || strings.HasPrefix(itemTrim, "vrrp-inet6-group ") {
		vrrpGroup := genVRRPGroupOld(family)
		vrrpID, err := strconv.Atoi(addressConfig[2])
		if err != nil {
			return inetAddress, nil
		}
		itemTrimVrrp := strings.TrimPrefix(itemTrim, "vrrp-group "+strconv.Itoa(vrrpID)+" ")
		if strings.HasPrefix(itemTrim, "vrrp-inet6-group ") {
			itemTrimVrrp = strings.TrimPrefix(itemTrim, "vrrp-inet6-group "+strconv.Itoa(vrrpID)+" ")
		}
		vrrpGroup["identifier"] = vrrpID
		vrrpGroup, m["vrrp_group"] = copyAndRemoveItemMapList("identifier", true, vrrpGroup,
			m["vrrp_group"].([]map[string]interface{}))
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
		m["vrrp_group"] = append(m["vrrp_group"].([]map[string]interface{}), vrrpGroup)
	}
	inetAddress = append(inetAddress, m)

	return inetAddress, nil
}
func setFamilyAddressOld(inetAddress interface{}, intCut []string, configSet []string, setName string,
	family string) ([]string, error) {
	if family != inetWord && family != inet6Word {
		return configSet, fmt.Errorf("setFamilyAddressOld() unknown family %v", family)
	}
	inetAddressMap := inetAddress.(map[string]interface{})
	configSet = append(configSet, "set interfaces "+setName+" family "+family+
		" address "+inetAddressMap["address"].(string))
	for _, vrrpGroup := range inetAddressMap["vrrp_group"].([]interface{}) {
		if intCut[0] == st0Word {
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
			setNameAddVrrp = "set interfaces " + setName + " family inet address " + inetAddressMap["address"].(string) +
				" vrrp-group " + strconv.Itoa(vrrpGroupMap["identifier"].(int))
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
			setNameAddVrrp = "set interfaces " + setName + " family inet6 address " + inetAddressMap["address"].(string) +
				" vrrp-inet6-group " + strconv.Itoa(vrrpGroupMap["identifier"].(int))
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

func aggregatedLastChild(ae, interFace string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConf, err := sess.command("show configuration interfaces | display set relative", jnprSess)
	if err != nil {
		return false, err
	}
	lastAE := true
	for _, item := range strings.Split(showConf, "\n") {
		if strings.HasSuffix(item, "ether-options 802.3ad "+ae) &&
			!strings.HasPrefix(item, "set "+interFace+" ") {
			lastAE = false
		}
	}

	return lastAE, nil
}
func aggregatedCountSearchMax(newAE, oldAE, interFace string, m interface{}, jnprSess *NetconfObject) (string, error) {
	sess := m.(*Session)
	newAENum := strings.TrimPrefix(newAE, "ae")
	newAENumInt, err := strconv.Atoi(newAENum)
	if err != nil {
		return "", fmt.Errorf("failed to convert internal variable newAENum to integer : %w", err)
	}
	intShowInt, err := sess.command("show interfaces terse", jnprSess)
	if err != nil {
		return "", err
	}

	intShowIntLines := strings.Split(intShowInt, "\n")
	intShowAE := make([]string, 0)
	regexpAE := regexp.MustCompile(`ae\d*\s`)
	for _, line := range intShowIntLines {
		aematch := regexpAE.MatchString(line)
		if aematch {
			wordsLine := strings.Fields(line)
			if wordsLine[0] != oldAE {
				if (len(intShowAE) > 0 && intShowAE[len(intShowAE)-1] != wordsLine[0]) || len(intShowAE) == 0 {
					intShowAE = append(intShowAE, wordsLine[0])
				}
			}
		}
	}
	lastOldAE, err := aggregatedLastChild(oldAE, interFace, m, jnprSess)
	if err != nil {
		return "", err
	}
	if !lastOldAE {
		intShowAE = append(intShowAE, oldAE)
	}
	if len(intShowAE) > 0 {
		lastAeInt, err := strconv.Atoi(strings.TrimPrefix(intShowAE[len(intShowAE)-1], "ae"))
		if err != nil {
			return "", fmt.Errorf("failed to convert internal variable lastAeInt to integer : %w", err)
		}
		if lastAeInt > newAENumInt {
			return strconv.Itoa(lastAeInt + 1), nil
		}
	}

	return strconv.Itoa(newAENumInt + 1), nil
}
func genFamilyInetAddressOld(address string) map[string]interface{} {
	return map[string]interface{}{
		"address":    address,
		"vrrp_group": make([]map[string]interface{}, 0),
	}
}
func genVRRPGroupOld(family string) map[string]interface{} {
	m := map[string]interface{}{
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
		m["authentication_key"] = ""
		m["authentication_type"] = ""
	}
	if family == inet6Word {
		m["virtual_link_local_address"] = ""
	}

	return m
}
func checkResourceInterfaceConfigAndName(length int, d *schema.ResourceData) error {
	if length == 1 {
		if d.Get("vlan_tagging_id").(int) != 0 {
			return fmt.Errorf("vlan_tagging_id invalid for this interface")
		}
		if d.Get("inet").(bool) {
			return fmt.Errorf("inet invalid for this interface")
		}
		if d.Get("inet6").(bool) {
			return fmt.Errorf("inet6 invalid for this interface")
		}
		if len(d.Get("inet_address").([]interface{})) > 0 {
			return fmt.Errorf("inet address invalid for this interface")
		}
		if len(d.Get("inet6_address").([]interface{})) > 0 {
			return fmt.Errorf("inet6 address invalid for this interface")
		}
		if d.Get("inet_mtu").(int) > 0 {
			return fmt.Errorf("inet_mtu invalid for this interface")
		}
		if d.Get("inet6_mtu").(int) > 0 {
			return fmt.Errorf("inet6_mtu invalid for this interface")
		}
		if d.Get("inet_filter_input").(string) != "" {
			return fmt.Errorf("inet_filter_input invalid for this interface")
		}
		if d.Get("inet_filter_output").(string) != "" {
			return fmt.Errorf("inet_filter_output invalid for this interface")
		}
		if d.Get("inet6_filter_input").(string) != "" {
			return fmt.Errorf("inet6_filter_input invalid for this interface")
		}
		if d.Get("inet6_filter_output").(string) != "" {
			return fmt.Errorf("inet6_filter_output invalid for this interface")
		}
		if len(d.Get("inet_rpf_check").([]interface{})) > 0 {
			return fmt.Errorf("inet_rpf_check invalid for this interface")
		}
		if len(d.Get("inet6_rpf_check").([]interface{})) > 0 {
			return fmt.Errorf("inet6_rpf_check invalid for this interface")
		}
		if d.Get("security_zone").(string) != "" {
			return fmt.Errorf("security_zone invalid for this interface")
		}
		if d.Get("routing_instance").(string) != "" {
			return fmt.Errorf("routing_instance invalid for this interface")
		}
	}
	if length == 2 {
		if d.Get("vlan_tagging").(bool) {
			return fmt.Errorf("vlan tagging invalid for this interface")
		}
		if d.Get("ether802_3ad").(string) != "" {
			return fmt.Errorf("ether802_3ad invalid for this interface")
		}
		if d.Get("trunk").(bool) {
			return fmt.Errorf("trunk invalid for this interface (remove .0)")
		}
		if len(d.Get("vlan_members").([]interface{})) > 0 {
			return fmt.Errorf("vlan_members invalid for this interface (remove .0)")
		}
		if d.Get("vlan_native").(int) != 0 {
			return fmt.Errorf("vlan_members invalid for this interface (remove .0)")
		}
		if d.Get("ae_lacp").(string) != "" {
			return fmt.Errorf("ae_lacp invalid for this interface")
		}
		if d.Get("ae_link_speed").(string) != "" {
			return fmt.Errorf("ae_link_speed invalid for this interface")
		}
		if d.Get("ae_minimum_links").(int) > 0 {
			return fmt.Errorf("ae_minimum_links invalid for this interface")
		}
	}

	return nil
}
