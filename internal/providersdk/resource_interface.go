package providersdk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	jdecode "github.com/jeremmfr/junosdecode"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
		CreateWithoutTimeout: resourceInterfaceCreate,
		ReadWithoutTimeout:   resourceInterfaceRead,
		UpdateWithoutTimeout: resourceInterfaceUpdate,
		DeleteWithoutTimeout: resourceInterfaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceInterfaceImport,
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
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
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
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"inet_filter_output": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"inet6_filter_input": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"inet6_filter_output": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
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
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
		},
	}
}

func resourceInterfaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if clt.GroupInterfaceDelete() != "" {
			if err := delInterfaceElement("apply-groups "+clt.GroupInterfaceDelete(), d, clt, nil); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := delInterfaceElement("disable", d, clt, nil); err != nil {
				return diag.FromErr(err)
			}
			if err := delInterfaceElement("description", d, clt, nil); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := setInterface(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	intExists, err := checkInterfaceExistsOld(d.Get("name").(string), clt, junSess)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if intExists {
		ncInt, emptyInt, err := checkInterfaceNC(d.Get("name").(string), clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !ncInt && !emptyInt {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(fmt.Errorf("interface %s already configured", d.Get("name").(string)))...)
		}
		if clt.GroupInterfaceDelete() != "" {
			err = delInterfaceElement("apply-groups "+clt.GroupInterfaceDelete(), d, clt, junSess)
			if err != nil {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
		} else {
			err = delInterfaceElement("disable", d, clt, junSess)
			if err != nil {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
			err = delInterfaceElement("description", d, clt, junSess)
			if err != nil {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
		}
	}
	if d.Get("security_zone").(string) != "" {
		if !junos.CheckCompatibilitySecurity(junSess) {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(fmt.Errorf("security zone not compatible with Junos device %s",
				junSess.SystemInformation.HardwareModel))...)
		}
		zonesExists, err := checkSecurityZonesExists(d.Get("security_zone").(string), clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !zonesExists {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("security zone %v doesn't exist", d.Get("security_zone").(string)))...)
		}
	}
	if d.Get("routing_instance").(string) != "" {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	if err := setInterface(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_interface", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	intExists, err = checkInterfaceExistsOld(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if intExists {
		ncInt, _, err := checkInterfaceNC(d.Get("name").(string), clt, junSess)
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

	return append(diagWarns, resourceInterfaceReadWJunSess(d, clt, junSess)...)
}

func resourceInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceInterfaceReadWJunSess(d, clt, junSess)
}

func resourceInterfaceReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) diag.Diagnostics {
	mutex.Lock()
	intExists, err := checkInterfaceExistsOld(d.Get("name").(string), clt, junSess)
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	if !intExists {
		d.SetId("")
		mutex.Unlock()

		return nil
	}
	ncInt, _, err := checkInterfaceNC(d.Get("name").(string), clt, junSess)
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	if ncInt {
		d.SetId("")
		mutex.Unlock()

		return nil
	}
	interfaceOpt, err := readInterface(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillInterfaceData(d, interfaceOpt)

	return nil
}

func resourceInterfaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delInterfaceOpts(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if d.HasChange("ether802_3ad") {
		oAE, nAE := d.GetChange("ether802_3ad")
		if oAE.(string) != "" {
			newAE := "ae-1"
			if nAE.(string) != "" {
				newAE = nAE.(string)
			}
			lastAEchild, err := aggregatedLastChild(oAE.(string), d.Get("name").(string), clt, junSess)
			if err != nil {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
			if lastAEchild {
				aggregatedCount, err := aggregatedCountSearchMax(newAE, oAE.(string), d.Get("name").(string), clt, junSess)
				if err != nil {
					appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

					return append(diagWarns, diag.FromErr(err)...)
				}
				if aggregatedCount == "0" {
					err = clt.ConfigSet([]string{"delete chassis aggregated-devices ethernet device-count"}, junSess)
					if err != nil {
						appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

						return append(diagWarns, diag.FromErr(err)...)
					}
					oAEintNC, oAEintEmpty, err := checkInterfaceNC(oAE.(string), clt, junSess)
					if err != nil {
						appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

						return append(diagWarns, diag.FromErr(err)...)
					}
					if oAEintNC || oAEintEmpty {
						err = clt.ConfigSet([]string{"delete interfaces " + oAE.(string)}, junSess)
						if err != nil {
							appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

							return append(diagWarns, diag.FromErr(err)...)
						}
					}
				} else {
					oldAEInt, err := strconv.Atoi(strings.TrimPrefix(oAE.(string), "ae"))
					if err != nil {
						appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

						return append(diagWarns, diag.FromErr(err)...)
					}
					aggregatedCountInt, err := strconv.Atoi(aggregatedCount)
					if err != nil {
						appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

						return append(diagWarns, diag.FromErr(err)...)
					}
					if aggregatedCountInt < oldAEInt+1 {
						oAEintNC, oAEintEmpty, err := checkInterfaceNC(oAE.(string), clt, junSess)
						if err != nil {
							appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

							return append(diagWarns, diag.FromErr(err)...)
						}
						if oAEintNC || oAEintEmpty {
							err = clt.ConfigSet([]string{"delete interfaces " + oAE.(string)}, junSess)
							if err != nil {
								appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

								return append(diagWarns, diag.FromErr(err)...)
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
			if !junos.CheckCompatibilitySecurity(junSess) {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(fmt.Errorf("security zone not compatible with Junos device %s",
					junSess.SystemInformation.HardwareModel))...)
			}
			zonesExists, err := checkSecurityZonesExists(nSecurityZone.(string), clt, junSess)
			if err != nil {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
			if !zonesExists {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(fmt.Errorf("security zone %v doesn't exist", nSecurityZone.(string)))...)
			}
		}
		if oSecurityZone.(string) != "" {
			err = delZoneInterface(oSecurityZone.(string), d, clt, junSess)
			if err != nil {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
		}
	}
	if d.HasChange("routing_instance") {
		oRoutingInstance, nRoutingInstance := d.GetChange("routing_instance")
		if nRoutingInstance.(string) != "" {
			instanceExists, err := checkRoutingInstanceExists(nRoutingInstance.(string), clt, junSess)
			if err != nil {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
			if !instanceExists {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns,
					diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", nRoutingInstance.(string)))...)
			}
		}
		if oRoutingInstance.(string) != "" {
			err = delRoutingInstanceInterface(oRoutingInstance.(string), d, clt, junSess)
			if err != nil {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
		}
	}
	if err := setInterface(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_interface", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceInterfaceReadWJunSess(d, clt, junSess)...)
}

func resourceInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delInterface(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_interface", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !d.Get("complete_destroy").(bool) {
		intExists, err := checkInterfaceExistsOld(d.Get("name").(string), clt, junSess)
		if err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		if intExists {
			err = addInterfaceNC(d.Get("name").(string), clt, junSess)
			if err != nil {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
			_, err = clt.CommitConf("disable(NC) resource junos_interface", junSess)
			if err != nil {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
		}
	}

	return diagWarns
}

func resourceInterfaceImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)
	intExists, err := checkInterfaceExistsOld(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !intExists {
		return nil, fmt.Errorf("don't find interface with id '%v' (id must be <name>)", d.Id())
	}
	interfaceOpt, err := readInterface(d.Id(), clt, junSess)
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

func checkInterfaceNC(interFace string, clt *junos.Client, junSess *junos.Session,
) (ncInt, emtyInt bool, errFunc error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+"interfaces "+interFace+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return false, false, err
	}
	showConfigLines := make([]string, 0)
	// remove unused lines
	for _, item := range strings.Split(showConfig, "\n") {
		// show parameters root on interface exclude unit parameters (except ethernet-switching)
		if !strings.Contains(interFace, ".") && strings.HasPrefix(item, "set unit") &&
			!strings.Contains(item, "ethernet-switching") {
			continue
		}
		if strings.Contains(item, junos.XMLStartTagConfigOut) {
			continue
		}
		if strings.Contains(item, junos.XMLEndTagConfigOut) {
			break
		}
		if item == "" {
			continue
		}
		showConfigLines = append(showConfigLines, item)
	}
	if len(showConfigLines) == 0 {
		return false, true, nil
	}
	showConfig = strings.Join(showConfigLines, "\n")
	if clt.GroupInterfaceDelete() != "" {
		if showConfig == "set apply-groups "+clt.GroupInterfaceDelete() {
			return true, false, nil
		}
	}
	if showConfig == "set description NC\nset disable" ||
		showConfig == "set disable\nset description NC" {
		return true, false, nil
	}
	if showConfig == junos.SetLS ||
		showConfig == junos.EmptyW {
		return false, true, nil
	}

	return false, false, nil
}

func addInterfaceNC(interFace string, clt *junos.Client, junSess *junos.Session) (err error) {
	intCut := make([]string, 0, 2)
	var setName string
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
	if intCut[0] == junos.St0Word || clt.GroupInterfaceDelete() == "" {
		err = clt.ConfigSet([]string{"set interfaces " + setName + " disable description NC"}, junSess)
	} else {
		err = clt.ConfigSet([]string{"set interfaces " + setName +
			" apply-groups " + clt.GroupInterfaceDelete()}, junSess)
	}
	if err != nil {
		return err
	}

	return nil
}

func checkInterfaceExistsOld(interFace string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	rpcIntName := "<get-interface-information><interface-name>" + interFace +
		"</interface-name></get-interface-information>"
	reply, err := clt.CommandXML(rpcIntName, junSess)
	if err != nil {
		if strings.Contains(err.Error(), " not found\n") ||
			strings.HasSuffix(err.Error(), " not found") {
			return false, nil
		}

		return false, err
	}
	if strings.Contains(reply, " not found\n") {
		return false, nil
	}

	return true, nil
}

func setInterface(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) (err error) {
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
	} else if len(intCut) == 2 && intCut[0] != junos.St0Word && intCut[1] != "0" {
		configSet = append(configSet, setPrefix+"vlan-id "+intCut[1])
	}
	if d.Get("inet").(bool) {
		configSet = append(configSet, setPrefix+"family inet")
	}
	if d.Get("inet6").(bool) {
		configSet = append(configSet, setPrefix+"family inet6")
	}
	for _, address := range d.Get("inet_address").([]interface{}) {
		configSet, err = setFamilyAddressOld(address, intCut, configSet, setName, junos.InetW)
		if err != nil {
			return err
		}
	}
	for _, address := range d.Get("inet6_address").([]interface{}) {
		configSet, err = setFamilyAddressOld(address, intCut, configSet, setName, junos.Inet6W)
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
		if junSess != nil {
			aggregatedCount, err := aggregatedCountSearchMax(
				d.Get("ether802_3ad").(string),
				oldAE,
				d.Get("name").(string),
				clt, junSess)
			if err != nil {
				return err
			}
			configSet = append(configSet, "set chassis aggregated-devices ethernet device-count "+aggregatedCount)
		}
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
	if d.Get("security_zone").(string) != "" {
		configSet = append(configSet, "set security zones security-zone "+
			d.Get("security_zone").(string)+" interfaces "+d.Get("name").(string))
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, junos.SetRoutingInstances+d.Get("routing_instance").(string)+
			" interface "+d.Get("name").(string))
	}

	return clt.ConfigSet(configSet, junSess)
}

func readInterface(interFace string, clt *junos.Client, junSess *junos.Session) (confRead interfaceOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+"interfaces "+interFace+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if !strings.Contains(interFace, ".") && strings.Contains(item, " unit ") &&
				!strings.Contains(item, "ethernet-switching") {
				continue
			}
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			case itemTrim == "vlan-tagging":
				confRead.vlanTagging = true
			case balt.CutPrefixInString(&itemTrim, "vlan-id "):
				confRead.vlanTaggingID, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "family inet6"):
				confRead.inet6 = true
				switch {
				case balt.CutPrefixInString(&itemTrim, " address "):
					confRead.inet6Address, err = readFamilyInetAddressOld(itemTrim, confRead.inet6Address, junos.Inet6W)
					if err != nil {
						return confRead, err
					}
				case balt.CutPrefixInString(&itemTrim, " mtu "):
					confRead.inet6Mtu, err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, " filter input "):
					confRead.inet6FilterInput = itemTrim
				case balt.CutPrefixInString(&itemTrim, " filter output "):
					confRead.inet6FilterOutput = itemTrim
				case balt.CutPrefixInString(&itemTrim, " rpf-check"):
					if len(confRead.inet6RpfCheck) == 0 {
						confRead.inet6RpfCheck = append(confRead.inet6RpfCheck, map[string]interface{}{
							"fail_filter": "",
							"mode_loose":  false,
						})
					}
					switch {
					case balt.CutPrefixInString(&itemTrim, " fail-filter "):
						confRead.inet6RpfCheck[0]["fail_filter"] = strings.Trim(itemTrim, "\"")
					case itemTrim == " mode loose":
						confRead.inet6RpfCheck[0]["mode_loose"] = true
					}
				}
			case balt.CutPrefixInString(&itemTrim, "family inet"):
				confRead.inet = true
				switch {
				case balt.CutPrefixInString(&itemTrim, " address "):
					confRead.inetAddress, err = readFamilyInetAddressOld(itemTrim, confRead.inetAddress, junos.InetW)
					if err != nil {
						return confRead, err
					}
				case balt.CutPrefixInString(&itemTrim, " mtu "):
					confRead.inetMtu, err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, " filter input "):
					confRead.inetFilterInput = itemTrim
				case balt.CutPrefixInString(&itemTrim, " filter output "):
					confRead.inetFilterOutput = itemTrim
				case balt.CutPrefixInString(&itemTrim, " rpf-check"):
					if len(confRead.inetRpfCheck) == 0 {
						confRead.inetRpfCheck = append(confRead.inetRpfCheck, map[string]interface{}{
							"fail_filter": "",
							"mode_loose":  false,
						})
					}
					switch {
					case balt.CutPrefixInString(&itemTrim, " fail-filter "):
						confRead.inetRpfCheck[0]["fail_filter"] = strings.Trim(itemTrim, "\"")
					case itemTrim == " mode loose":
						confRead.inetRpfCheck[0]["mode_loose"] = true
					}
				}
			case balt.CutPrefixInString(&itemTrim, "ether-options 802.3ad "):
				confRead.v8023ad = itemTrim
			case balt.CutPrefixInString(&itemTrim, "gigether-options 802.3ad "):
				confRead.v8023ad = itemTrim
			case itemTrim == "unit 0 family ethernet-switching interface-mode trunk":
				confRead.trunk = true
			case balt.CutPrefixInString(&itemTrim, "unit 0 family ethernet-switching vlan members "):
				confRead.vlanMembers = append(confRead.vlanMembers, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "native-vlan-id "):
				confRead.vlanNative, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "aggregated-ether-options lacp "):
				confRead.aeLacp = itemTrim
			case balt.CutPrefixInString(&itemTrim, "aggregated-ether-options link-speed "):
				confRead.aeLinkSpeed = itemTrim
			case balt.CutPrefixInString(&itemTrim, "aggregated-ether-options minimum-links "):
				confRead.aeMinLink, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			default:
				continue
			}
		}
	}
	if junos.CheckCompatibilitySecurity(junSess) {
		showConfigSecurityZones, err := clt.Command(
			junos.CmdShowConfig+"security zones"+junos.PipeDisplaySetRelative,
			junSess,
		)
		if err != nil {
			return confRead, err
		}
		regexpInts := regexp.MustCompile(`set security-zone \S+ interfaces ` + interFace + `$`)
		for _, item := range strings.Split(showConfigSecurityZones, "\n") {
			intMatch := regexpInts.MatchString(item)
			if intMatch {
				confRead.securityZones = strings.TrimPrefix(strings.TrimSuffix(item, " interfaces "+interFace),
					"set security-zone ")

				break
			}
		}
	}
	showConfigRoutingInstances, err := clt.Command(
		junos.CmdShowConfig+"routing-instances"+junos.PipeDisplaySetRelative,
		junSess,
	)
	if err != nil {
		return confRead, err
	}
	regexpInt := regexp.MustCompile(`set \S+ interface ` + interFace + `$`)
	for _, item := range strings.Split(showConfigRoutingInstances, "\n") {
		intMatch := regexpInt.MatchString(item)
		if intMatch {
			confRead.routingInstances = strings.TrimPrefix(strings.TrimSuffix(item, " interface "+interFace), junos.SetLS)

			break
		}
	}

	return confRead, nil
}

func delInterface(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
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
		err := checkInterfaceContainsUnit(setName, clt, junSess)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("the name %s contains too dots", d.Get("name").(string))
	}

	if err := clt.ConfigSet([]string{"delete interfaces " + setName}, junSess); err != nil {
		return err
	}
	if strings.Contains(d.Get("name").(string), "st0.") && !d.Get("complete_destroy").(bool) {
		// interface totally delete by junos_interface_st0_unit resource
		// else there is an interface st0.x empty
		err := clt.ConfigSet([]string{"set interfaces " + setName}, junSess)
		if err != nil {
			return err
		}
	}
	if d.Get("ether802_3ad").(string) != "" {
		lastAEchild, err := aggregatedLastChild(d.Get("ether802_3ad").(string), d.Get("name").(string), clt, junSess)
		if err != nil {
			return err
		}
		if lastAEchild {
			aggregatedCount, err := aggregatedCountSearchMax(
				"ae-1",
				d.Get("ether802_3ad").(string),
				d.Get("name").(string),
				clt, junSess)
			if err != nil {
				return err
			}
			if aggregatedCount == "0" {
				err = clt.ConfigSet([]string{"delete chassis aggregated-devices ethernet device-count"}, junSess)
				if err != nil {
					return err
				}
			} else {
				err = clt.ConfigSet([]string{"set chassis aggregated-devices ethernet device-count " +
					aggregatedCount}, junSess)
				if err != nil {
					return err
				}
			}
			aeInt, err := strconv.Atoi(strings.TrimPrefix(d.Get("ether802_3ad").(string), "ae"))
			if err != nil {
				return fmt.Errorf("failed to convert AE id of ether802_3ad argument '%s' in integer: %w",
					d.Get("ether802_3ad").(string), err)
			}
			aggregatedCountInt, err := strconv.Atoi(aggregatedCount)
			if err != nil {
				return fmt.Errorf("failed to convert internal variable aggregatedCountInt in integer: %w", err)
			}
			if aggregatedCountInt < aeInt+1 {
				oAEintNC, oAEintEmpty, err := checkInterfaceNC(d.Get("ether802_3ad").(string), clt, junSess)
				if err != nil {
					return err
				}
				if oAEintNC || oAEintEmpty {
					err = clt.ConfigSet([]string{"delete interfaces " + d.Get("ether802_3ad").(string)}, junSess)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	if junos.CheckCompatibilitySecurity(junSess) && d.Get("security_zone").(string) != "" {
		if err := delZoneInterface(d.Get("security_zone").(string), d, clt, junSess); err != nil {
			return err
		}
	}
	if d.Get("routing_instance").(string) != "" {
		if err := delRoutingInstanceInterface(d.Get("routing_instance").(string), d, clt, junSess); err != nil {
			return err
		}
	}

	return nil
}

func checkInterfaceContainsUnit(interFace string, clt *junos.Client, junSess *junos.Session) error {
	showConfig, err := clt.Command(junos.CmdShowConfig+"interfaces "+interFace+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return err
	}
	for _, item := range strings.Split(showConfig, "\n") {
		if strings.Contains(item, junos.XMLStartTagConfigOut) {
			continue
		}
		if strings.Contains(item, junos.XMLEndTagConfigOut) {
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

func delInterfaceElement(element string, d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
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

	return clt.ConfigSet(configSet, junSess)
}

func delInterfaceOpts(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
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
		delPrefix+"aggregated-ether-options",
	)

	return clt.ConfigSet(configSet, junSess)
}

func delZoneInterface(zone string, d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security zones security-zone "+zone+" interfaces "+d.Get("name").(string))

	return clt.ConfigSet(configSet, junSess)
}

func delRoutingInstanceInterface(instance string, d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, junos.DelRoutingInstances+instance+" interface "+d.Get("name").(string))

	return clt.ConfigSet(configSet, junSess)
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

func readFamilyInetAddressOld(itemTrim string, inetAddress []map[string]interface{}, family string,
) ([]map[string]interface{}, error) {
	itemTrimFields := strings.Split(itemTrim, " ")
	balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
	mAddr := genFamilyInetAddressOld(itemTrimFields[0])
	inetAddress = copyAndRemoveItemMapList("address", mAddr, inetAddress)

	if balt.CutPrefixInString(&itemTrim, "vrrp-group ") || balt.CutPrefixInString(&itemTrim, "vrrp-inet6-group ") {
		if len(itemTrimFields) < 3 { // <address> (vrrp-group|vrrp-inet6-group) <vrrpID>
			return inetAddress, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "vrrp-group|vrrp-inet6-group", itemTrim)
		}
		vrrpGroup := genVRRPGroupOld(family)
		vrrpID, err := strconv.Atoi(itemTrimFields[2])
		if err != nil {
			return inetAddress, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
		balt.CutPrefixInString(&itemTrim, itemTrimFields[2]+" ")
		vrrpGroup["identifier"] = vrrpID
		mAddr["vrrp_group"] = copyAndRemoveItemMapList("identifier", vrrpGroup,
			mAddr["vrrp_group"].([]map[string]interface{}))
		switch {
		case balt.CutPrefixInString(&itemTrim, "virtual-address "):
			vrrpGroup["virtual_address"] = append(vrrpGroup["virtual_address"].([]string), itemTrim)
		case balt.CutPrefixInString(&itemTrim, "virtual-inet6-address "):
			vrrpGroup["virtual_address"] = append(vrrpGroup["virtual_address"].([]string), itemTrim)
		case balt.CutPrefixInString(&itemTrim, "virtual-link-local-address "):
			vrrpGroup["virtual_link_local_address"] = itemTrim
		case itemTrim == "accept-data":
			vrrpGroup["accept_data"] = true
		case balt.CutPrefixInString(&itemTrim, "advertise-interval "):
			vrrpGroup["advertise_interval"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return inetAddress, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "inet6-advertise-interval "):
			vrrpGroup["advertise_interval"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return inetAddress, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "advertisements-threshold "):
			vrrpGroup["advertisements_threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return inetAddress, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "authentication-key "):
			vrrpGroup["authentication_key"], err = jdecode.Decode(strings.Trim(itemTrim, "\""))
			if err != nil {
				return inetAddress, fmt.Errorf("failed to decode authentication-key: %w", err)
			}
		case balt.CutPrefixInString(&itemTrim, "authentication-type "):
			vrrpGroup["authentication_type"] = itemTrim
		case itemTrim == "no-accept-data":
			vrrpGroup["no_accept_data"] = true
		case itemTrim == "no-preempt":
			vrrpGroup["no_preempt"] = true
		case itemTrim == "preempt":
			vrrpGroup["preempt"] = true
		case balt.CutPrefixInString(&itemTrim, "priority "):
			vrrpGroup["priority"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return inetAddress, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "track interface "):
			itemTrackFields := strings.Split(itemTrim, " ")
			if len(itemTrackFields) < 3 { // <interface> priority-cost <priority_cost>
				return inetAddress, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "track interface", itemTrim)
			}
			cost, err := strconv.Atoi(itemTrackFields[2])
			if err != nil {
				return inetAddress, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
			trackInt := map[string]interface{}{
				"interface":     itemTrackFields[0],
				"priority_cost": cost,
			}
			vrrpGroup["track_interface"] = append(vrrpGroup["track_interface"].([]map[string]interface{}), trackInt)
		case balt.CutPrefixInString(&itemTrim, "track route "):
			itemTrackFields := strings.Split(itemTrim, " ")
			if len(itemTrackFields) < 5 { // <route> routing-instance <routing_instance> priority-cost <priority_cost>
				return inetAddress, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "track route", itemTrim)
			}
			cost, err := strconv.Atoi(itemTrackFields[4])
			if err != nil {
				return inetAddress, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
			trackRoute := map[string]interface{}{
				"route":            itemTrackFields[0],
				"routing_instance": itemTrackFields[2],
				"priority_cost":    cost,
			}
			vrrpGroup["track_route"] = append(vrrpGroup["track_route"].([]map[string]interface{}), trackRoute)
		}
		mAddr["vrrp_group"] = append(mAddr["vrrp_group"].([]map[string]interface{}), vrrpGroup)
	}
	inetAddress = append(inetAddress, mAddr)

	return inetAddress, nil
}

func setFamilyAddressOld(inetAddress interface{}, intCut, configSet []string, setName, family string,
) ([]string, error) {
	if family != junos.InetW && family != junos.Inet6W {
		return configSet, fmt.Errorf("setFamilyAddressOld() unknown family %v", family)
	}
	inetAddressMap := inetAddress.(map[string]interface{})
	configSet = append(configSet, "set interfaces "+setName+" family "+family+
		" address "+inetAddressMap["address"].(string))
	for _, vrrpGroup := range inetAddressMap["vrrp_group"].([]interface{}) {
		if intCut[0] == junos.St0Word {
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
		case junos.InetW:
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
		case junos.Inet6W:
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

func aggregatedLastChild(ae, interFace string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+"interfaces"+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return false, err
	}
	lastAE := true
	for _, item := range strings.Split(showConfig, "\n") {
		if strings.HasSuffix(item, "ether-options 802.3ad "+ae) &&
			!strings.HasPrefix(item, junos.SetLS+interFace+" ") {
			lastAE = false
		}
	}

	return lastAE, nil
}

func aggregatedCountSearchMax(newAE, oldAE, interFace string, clt *junos.Client, junSess *junos.Session,
) (string, error) {
	newAENum := strings.TrimPrefix(newAE, "ae")
	newAENumInt, err := strconv.Atoi(newAENum)
	if err != nil {
		return "", fmt.Errorf("failed to convert internal variable newAENum to integer: %w", err)
	}
	intShowInt, err := clt.Command("show interfaces terse", junSess)
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
	lastOldAE, err := aggregatedLastChild(oldAE, interFace, clt, junSess)
	if err != nil {
		return "", err
	}
	if !lastOldAE {
		intShowAE = append(intShowAE, oldAE)
	}
	if len(intShowAE) > 0 {
		lastAeInt, err := strconv.Atoi(strings.TrimPrefix(intShowAE[len(intShowAE)-1], "ae"))
		if err != nil {
			return "", fmt.Errorf("failed to convert internal variable lastAeInt to integer: %w", err)
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
	if family == junos.InetW {
		vrrpGroup["authentication_key"] = ""
		vrrpGroup["authentication_type"] = ""
	}
	if family == junos.Inet6W {
		vrrpGroup["virtual_link_local_address"] = ""
	}

	return vrrpGroup
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
