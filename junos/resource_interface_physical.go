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
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type interfacePhysicalOptions struct {
	trunk           bool
	vlanTagging     bool
	aeMinLink       int
	vlanNative      int
	aeLacp          string
	aeLinkSpeed     string
	description     string
	v8023ad         string
	vlanMembers     []string
	esi             []map[string]interface{}
	etherOpts       []map[string]interface{}
	gigetherOpts    []map[string]interface{}
	parentEtherOpts []map[string]interface{}
}

func resourceInterfacePhysical() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInterfacePhysicalCreate,
		ReadContext:   resourceInterfacePhysicalRead,
		UpdateContext: resourceInterfacePhysicalUpdate,
		DeleteContext: resourceInterfacePhysicalDelete,
		Importer: &schema.ResourceImporter{
			State: resourceInterfacePhysicalImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if strings.Count(value, ".") > 0 {
						errors = append(errors, fmt.Errorf(
							"%q in %q cannot have a dot", value, k))
					}

					return
				},
			},
			"no_disable_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ae_lacp": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringInSlice([]string{"active", "passive"}, false),
				ConflictsWith: []string{"parent_ether_opts"},
				Deprecated:    "use parent_ether_opts { lacp } instead",
			},
			"ae_link_speed": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringInSlice([]string{"100m", "1g", "8g", "10g", "40g", "50g", "80g", "100g"}, false),
				ConflictsWith: []string{"parent_ether_opts"},
				Deprecated:    "use parent_ether_opts { link_speed } instead",
			},
			"ae_minimum_links": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"parent_ether_opts"},
				Deprecated:    "use parent_ether_opts { minimum_links } instead",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"esi": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"all-active", "single-active"}, false),
						},
						"auto_derive_lacp": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"esi.0.identifier"},
						},
						"df_election_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"mod", "preference"}, false),
						},
						"identifier": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"esi.0.auto_derive_lacp"},
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^([\d\w]{2}:){9}[\d\w]{2}$`), "bad format or length"),
						},
						"source_bmac": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsMACAddress,
						},
					},
				},
			},
			"ether_opts": {
				Type:     schema.TypeList,
				Optional: true,
				ConflictsWith: []string{
					"ae_lacp", "ae_link_speed", "ae_minimum_links",
					"gigether_opts", "parent_ether_opts",
				},
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ae_8023ad": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"ether_opts.0.redundant_parent"},
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !strings.HasPrefix(value, "ae") {
									errors = append(errors, fmt.Errorf(
										"%q in %q isn't an ae interface", value, k))
								}

								return
							},
						},
						"auto_negotiation": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.no_auto_negotiation"},
						},
						"no_auto_negotiation": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.auto_negotiation"},
						},
						"flow_control": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.no_flow_control"},
						},
						"no_flow_control": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.flow_control"},
						},
						"loopback": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.no_loopback"},
						},
						"no_loopback": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.loopback"},
						},
						"redundant_parent": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"ether_opts.0.ae_8023ad"},
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !strings.HasPrefix(value, "reth") {
									errors = append(errors, fmt.Errorf(
										"%q in %q isn't an reth interface", value, k))
								}

								return
							},
						},
					},
				},
			},
			"ether802_3ad": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use ether_opts { ae_8023ad } or gigether_opts { ae_8023ad } instead",
				ConflictsWith: []string{
					"ae_lacp", "ae_link_speed", "ae_minimum_links",
					"ether_opts", "gigether_opts",
				},
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !strings.HasPrefix(value, "ae") {
						errors = append(errors, fmt.Errorf(
							"%q in %q isn't an ae interface", value, k))
					}

					return
				},
			},
			"gigether_opts": {
				Type:     schema.TypeList,
				Optional: true,
				ConflictsWith: []string{
					"ae_lacp", "ae_link_speed", "ae_minimum_links",
					"ether_opts", "ether802_3ad", "parent_ether_opts",
				},
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ae_8023ad": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.redundant_parent"},
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !strings.HasPrefix(value, "ae") {
									errors = append(errors, fmt.Errorf(
										"%q in %q isn't an ae interface", value, k))
								}

								return
							},
						},
						"auto_negotiation": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.no_auto_negotiation"},
						},
						"no_auto_negotiation": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.auto_negotiation"},
						},
						"flow_control": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.no_flow_control"},
						},
						"no_flow_control": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.flow_control"},
						},
						"loopback": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.no_loopback"},
						},
						"no_loopback": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.loopback"},
						},
						"redundant_parent": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"gigether_opts.0.ae_8023ad"},
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !strings.HasPrefix(value, "reth") {
									errors = append(errors, fmt.Errorf(
										"%q in %q isn't an reth interface", value, k))
								}

								return
							},
						},
					},
				},
			},
			"parent_ether_opts": {
				Type:     schema.TypeList,
				Optional: true,
				ConflictsWith: []string{
					"ae_lacp", "ae_link_speed", "ae_minimum_links",
					"ether_opts", "ether802_3ad", "gigether_opts",
				},
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bfd_liveness_detection": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"local_address": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsIPAddress,
									},
									"authentication_algorithm": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"keyed-md5", "keyed-sha-1", "meticulous-keyed-md5", "meticulous-keyed-sha-1", "simple-password",
										}, false),
									},
									"authentication_key_chain": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"authentication_loose_check": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"detection_time_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 4294967295),
									},
									"holddown_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 255000),
									},
									"minimum_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 255000),
									},
									"minimum_receive_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 255000),
									},
									"multiplier": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 255),
									},
									"neighbor": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPAddress,
									},
									"no_adaptation": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"transmit_interval_minimum_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 255000),
									},
									"transmit_interval_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 4294967295),
									},
									"version": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"0", "1", "automatic"}, false),
									},
								},
							},
						},
						"flow_control": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"parent_ether_opts.0.no_flow_control"},
						},
						"no_flow_control": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"parent_ether_opts.0.flow_control"},
						},
						"lacp": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"active", "passive"}, false),
									},
									"admin_key": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 65535),
									},
									"periodic": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"fast", "slow"}, false),
									},
									"sync_reset": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"disable", "enable"}, false),
									},
									"system_id": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsMACAddress,
									},
									"system_priority": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 65535),
									},
								},
							},
						},
						"loopback": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"parent_ether_opts.0.no_loopback"},
						},
						"no_loopback": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"parent_ether_opts.0.loopback"},
						},
						"link_speed": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"100m", "1g", "2.5g", "5g", "8g",
								"10g", "25g", "40g", "50g", "80g",
								"100g", "400g", "mixed",
							}, false),
						},
						"minimum_bandwidth": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"parent_ether_opts.0.minimum_links"},
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^[0-9]+ (k|g|m)?bps$`), "must be 'N (k|g|m)?bps' format"),
						},
						"minimum_links": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"parent_ether_opts.0.minimum_bandwidth"},
							ValidateFunc:  validation.IntBetween(1, 64),
						},
						"redundancy_group": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 128),
						},
						"source_address_filter": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"source_filtering": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
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
			"vlan_tagging": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceInterfacePhysicalCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := delInterfaceNC(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setInterfacePhysical(d, m, nil); err != nil {
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
	ncInt, emptyInt, err := checkInterfacePhysicalNCEmpty(d.Get("name").(string), m, jnprSess)
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
	if err := setInterfacePhysical(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_interface_physical", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ncInt, emptyInt, err = checkInterfacePhysicalNCEmpty(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ncInt {
		return append(diagWarns, diag.FromErr(fmt.Errorf("interface %v always disable after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}
	if emptyInt {
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

	return append(diagWarns, resourceInterfacePhysicalReadWJnprSess(d, m, jnprSess)...)
}

func resourceInterfacePhysicalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceInterfacePhysicalReadWJnprSess(d, m, jnprSess)
}

func resourceInterfacePhysicalReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	ncInt, emptyInt, err := checkInterfacePhysicalNCEmpty(d.Get("name").(string), m, jnprSess)
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	if ncInt {
		d.SetId("")
		mutex.Unlock()

		return nil
	}
	if emptyInt {
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
	interfaceOpt, err := readInterfacePhysical(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillInterfacePhysicalData(d, interfaceOpt)

	return nil
}

func resourceInterfacePhysicalUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delInterfacePhysicalOpts(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setInterfacePhysical(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delInterfacePhysicalOpts(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := unsetInterfacePhysicalAE(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setInterfacePhysical(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_interface_physical", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceInterfacePhysicalReadWJnprSess(d, m, jnprSess)...)
}

func resourceInterfacePhysicalDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delInterfacePhysical(d, m, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delInterfacePhysical(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_interface_physical", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !d.Get("no_disable_on_destroy").(bool) {
		intExists, err := checkInterfaceExists(d.Get("name").(string), m, jnprSess)
		if err != nil {
			appendDiagWarns(&diagWarns, []error{err})
		} else if intExists {
			err = addInterfacePhysicalNC(d.Get("name").(string), m, jnprSess)
			if err != nil {
				appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
			_, err = sess.commitConf("disable(NC) resource junos_interface_physical", jnprSess)
			if err != nil {
				appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
		}
	}

	return diagWarns
}

func resourceInterfacePhysicalImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if strings.Count(d.Id(), ".") != 0 {
		return nil, fmt.Errorf("name of interface %s need to doesn't have a dot", d.Id())
	}
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	ncInt, emptyInt, err := checkInterfacePhysicalNCEmpty(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if ncInt {
		return nil, fmt.Errorf("interface '%v' is disabled, import is not possible", d.Id())
	}
	if emptyInt {
		intExists, err := checkInterfaceExists(d.Id(), m, jnprSess)
		if err != nil {
			return nil, err
		}
		if !intExists {
			return nil, fmt.Errorf("don't find interface with id '%v' (id must be <name>)", d.Id())
		}
	}
	interfaceOpt, err := readInterfacePhysical(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if tfErr := d.Set("name", d.Id()); tfErr != nil {
		panic(tfErr)
	}
	fillInterfacePhysicalData(d, interfaceOpt)

	result[0] = d

	return result, nil
}

func checkInterfacePhysicalNCEmpty(interFace string, m interface{}, jnprSess *NetconfObject) (
	ncInt bool, emtyInt bool, errFunc error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration interfaces "+interFace+" | display set relative", jnprSess)
	if err != nil {
		return false, false, err
	}
	showConfigLines := make([]string, 0)
	// remove unused lines
	for _, item := range strings.Split(showConfig, "\n") {
		// show parameters root on interface exclude unit parameters (except ethernet-switching)
		if strings.HasPrefix(item, "set unit") && !strings.Contains(item, "ethernet-switching") {
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
		showConfigLines = append(showConfigLines, item)
	}
	if len(showConfigLines) == 0 {
		return false, true, nil
	}
	showConfig = strings.Join(showConfigLines, "\n")
	if sess.junosGroupIntDel != "" {
		if showConfig == "set apply-groups "+sess.junosGroupIntDel {
			return true, false, nil
		}
	}
	if showConfig == "set description NC\nset disable" ||
		showConfig == "set disable\nset description NC" {
		return true, false, nil
	}
	if showConfig == emptyWord {
		return false, true, nil
	}

	return false, false, nil
}

func addInterfacePhysicalNC(interFace string, m interface{}, jnprSess *NetconfObject) error {
	var err error
	if sess := m.(*Session); sess.junosGroupIntDel == "" {
		err = sess.configSet([]string{"set interfaces " + interFace + " disable description NC"}, jnprSess)
	} else {
		err = sess.configSet([]string{"set interfaces " + interFace +
			" apply-groups " + sess.junosGroupIntDel}, jnprSess)
	}
	if err != nil {
		return err
	}

	return nil
}

func checkInterfaceExists(interFace string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	rpcIntName := "<get-interface-information><interface-name>" + interFace +
		"</interface-name></get-interface-information>"
	reply, err := sess.commandXML(rpcIntName, jnprSess)
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

func unsetInterfacePhysicalAE(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	var oldAE string
	switch {
	case d.HasChange("ether802_3ad"):
		oldAEtf, _ := d.GetChange("ether802_3ad")
		if oldAEtf.(string) != "" {
			oldAE = oldAEtf.(string)
		}
	case d.HasChange("ether_opts"):
		oldEthOpts, _ := d.GetChange("ether_opts")
		if len(oldEthOpts.([]interface{})) != 0 {
			v := oldEthOpts.([]interface{})[0]
			if o := v.(map[string]interface{})["ae_8023ad"].(string); o != "" {
				oldAE = o
			}
		}
	case d.HasChange("gigether_opts"):
		oldGigethOpts, _ := d.GetChange("gigether_opts")
		if len(oldGigethOpts.([]interface{})) != 0 {
			v := oldGigethOpts.([]interface{})[0]
			if o := v.(map[string]interface{})["ae_8023ad"].(string); o != "" {
				oldAE = o
			}
		}
	}
	if oldAE != "" {
		aggregatedCount, err := interfaceAggregatedCountSearchMax("ae-1", oldAE, d.Get("name").(string), m, jnprSess)
		if err != nil {
			return err
		}
		if aggregatedCount == "0" {
			return sess.configSet([]string{"delete chassis aggregated-devices ethernet device-count"}, jnprSess)
		}

		return sess.configSet([]string{"set chassis aggregated-devices ethernet device-count " + aggregatedCount}, jnprSess)
	}

	return nil
}

func setInterfacePhysical(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := "set interfaces " + d.Get("name").(string) + " "
	configSet = append(configSet, setPrefix)
	if d.Get("ae_lacp").(string) != "" {
		if !strings.HasPrefix(d.Get("name").(string), "ae") {
			return fmt.Errorf("ae_lacp invalid for this interface")
		}
		configSet = append(configSet, setPrefix+
			"aggregated-ether-options lacp "+d.Get("ae_lacp").(string))
	}
	if d.Get("ae_link_speed").(string) != "" {
		if !strings.HasPrefix(d.Get("name").(string), "ae") {
			return fmt.Errorf("ae_link_speed invalid for this interface")
		}
		configSet = append(configSet, setPrefix+
			"aggregated-ether-options link-speed "+d.Get("ae_link_speed").(string))
	}
	if d.Get("ae_minimum_links").(int) > 0 {
		if !strings.HasPrefix(d.Get("name").(string), "ae") {
			return fmt.Errorf("ae_minimum_links invalid for this interface")
		}
		configSet = append(configSet, setPrefix+
			"aggregated-ether-options minimum-links "+strconv.Itoa(d.Get("ae_minimum_links").(int)))
	}
	if d.Get("description").(string) != "" {
		configSet = append(configSet, setPrefix+"description \""+d.Get("description").(string)+"\"")
	}
	if err := setInterfacePhysicalEsi(setPrefix, d.Get("esi").([]interface{}), m, jnprSess); err != nil {
		return err
	}
	if v := d.Get("name").(string); strings.HasPrefix(v, "ae") && jnprSess != nil {
		aggregatedCount, err := interfaceAggregatedCountSearchMax(v, "ae-1", v, m, jnprSess)
		if err != nil {
			return err
		}
		configSet = append(configSet, "set chassis aggregated-devices ethernet device-count "+aggregatedCount)
	} else if d.Get("ether802_3ad").(string) != "" ||
		len(d.Get("ether_opts").([]interface{})) != 0 ||
		len(d.Get("gigether_opts").([]interface{})) != 0 {
		oldAE := "ae-1"
		var newAE string
		switch {
		case d.Get("ether802_3ad").(string) != "":
			newAE = d.Get("ether802_3ad").(string)
			configSet = append(configSet, setPrefix+"ether-options 802.3ad "+
				d.Get("ether802_3ad").(string))
			configSet = append(configSet, setPrefix+"gigether-options 802.3ad "+
				d.Get("ether802_3ad").(string))
		case len(d.Get("ether_opts").([]interface{})) != 0:
			for _, v := range d.Get("ether_opts").([]interface{}) {
				if v == nil {
					return fmt.Errorf("ether_opts block is empty")
				}
				eOpts := v.(map[string]interface{})
				if eOpts["ae_8023ad"].(string) != "" {
					newAE = eOpts["ae_8023ad"].(string)
					configSet = append(configSet, setPrefix+"ether-options 802.3ad "+
						eOpts["ae_8023ad"].(string))
				}
				if eOpts["auto_negotiation"].(bool) {
					configSet = append(configSet, setPrefix+"ether-options auto-negotiation")
				}
				if eOpts["no_auto_negotiation"].(bool) {
					configSet = append(configSet, setPrefix+"ether-options no-auto-negotiation")
				}
				if eOpts["flow_control"].(bool) {
					configSet = append(configSet, setPrefix+"ether-options flow-control")
				}
				if eOpts["no_flow_control"].(bool) {
					configSet = append(configSet, setPrefix+"ether-options no-flow-control")
				}
				if eOpts["loopback"].(bool) {
					configSet = append(configSet, setPrefix+"ether-options loopback")
				}
				if eOpts["no_loopback"].(bool) {
					configSet = append(configSet, setPrefix+"ether-options no-loopback")
				}
				if eOpts["redundant_parent"].(string) != "" {
					configSet = append(configSet, setPrefix+"ether-options redundant-parent "+
						eOpts["redundant_parent"].(string))
				}
			}
		case len(d.Get("gigether_opts").([]interface{})) != 0:
			for _, v := range d.Get("gigether_opts").([]interface{}) {
				if v == nil {
					return fmt.Errorf("gigether_opts block is empty")
				}
				geOpts := v.(map[string]interface{})
				if geOpts["ae_8023ad"].(string) != "" {
					newAE = geOpts["ae_8023ad"].(string)
					configSet = append(configSet, setPrefix+"gigether-options 802.3ad "+
						geOpts["ae_8023ad"].(string))
				}
				if geOpts["auto_negotiation"].(bool) {
					configSet = append(configSet, setPrefix+"gigether-options auto-negotiation")
				}
				if geOpts["no_auto_negotiation"].(bool) {
					configSet = append(configSet, setPrefix+"gigether-options no-auto-negotiation")
				}
				if geOpts["flow_control"].(bool) {
					configSet = append(configSet, setPrefix+"gigether-options flow-control")
				}
				if geOpts["no_flow_control"].(bool) {
					configSet = append(configSet, setPrefix+"gigether-options no-flow-control")
				}
				if geOpts["loopback"].(bool) {
					configSet = append(configSet, setPrefix+"gigether-options loopback")
				}
				if geOpts["no_loopback"].(bool) {
					configSet = append(configSet, setPrefix+"gigether-options no-loopback")
				}
				if geOpts["redundant_parent"].(string) != "" {
					configSet = append(configSet, setPrefix+"gigether-options redundant-parent "+
						geOpts["redundant_parent"].(string))
				}
			}
		}
		switch {
		case d.HasChange("ether802_3ad"):
			oldAEtf, _ := d.GetChange("ether802_3ad")
			if oldAEtf.(string) != "" {
				oldAE = oldAEtf.(string)
			}
		case d.HasChange("ether_opts"):
			oldEthOpts, _ := d.GetChange("ether_opts")
			if len(oldEthOpts.([]interface{})) != 0 {
				v := oldEthOpts.([]interface{})[0]
				if o := v.(map[string]interface{})["ae_8023ad"].(string); o != "" {
					oldAE = o
				}
			}
		case d.HasChange("gigether_opts"):
			oldGigethOpts, _ := d.GetChange("gigether_opts")
			if len(oldGigethOpts.([]interface{})) != 0 {
				v := oldGigethOpts.([]interface{})[0]
				if o := v.(map[string]interface{})["ae_8023ad"].(string); o != "" {
					oldAE = o
				}
			}
		}
		if newAE != "" && jnprSess != nil {
			aggregatedCount, err := interfaceAggregatedCountSearchMax(newAE, oldAE,
				d.Get("name").(string), m, jnprSess)
			if err != nil {
				return err
			}
			configSet = append(configSet, "set chassis aggregated-devices ethernet device-count "+aggregatedCount)
		}
	}
	for _, v := range d.Get("parent_ether_opts").([]interface{}) {
		if v == nil {
			return fmt.Errorf("parent_ether_opts block is empty")
		}
		if err := setInterfacePhysicalParentEtherOpts(
			v.(map[string]interface{}), d.Get("name").(string), m, jnprSess); err != nil {
			return err
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
	if d.Get("vlan_tagging").(bool) {
		configSet = append(configSet, setPrefix+"vlan-tagging")
	}

	return sess.configSet(configSet, jnprSess)
}

func setInterfacePhysicalEsi(setPrefix string, esiParams []interface{},
	m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	for _, v := range esiParams {
		esi := v.(map[string]interface{})
		if esi["mode"].(string) != "" {
			configSet = append(configSet, setPrefix+"esi "+esi["mode"].(string))
		}
		if esi["auto_derive_lacp"].(bool) {
			configSet = append(configSet, setPrefix+"esi auto-derive lacp")
		}
		if esi["df_election_type"].(string) != "" {
			configSet = append(configSet, setPrefix+"esi df-election-type "+esi["df_election_type"].(string))
		}
		if esi["identifier"].(string) != "" {
			configSet = append(configSet, setPrefix+"esi "+esi["identifier"].(string))
		}
		if esi["source_bmac"].(string) != "" {
			configSet = append(configSet, setPrefix+"esi source-bmac "+esi["source_bmac"].(string))
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func setInterfacePhysicalParentEtherOpts(
	ethOpts map[string]interface{}, interfaceName string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := "set interfaces " + interfaceName + " "
	switch {
	case strings.HasPrefix(interfaceName, "ae"):
		setPrefix += "aggregated-ether-options "
	case strings.HasPrefix(interfaceName, "reth"):
		setPrefix += "redundant-ether-options "
	default:
		return fmt.Errorf("parent_ether_opts not compatible with this interface %s "+
			"(need to ae* or reth*)", interfaceName)
	}

	for _, v := range ethOpts["bfd_liveness_detection"].([]interface{}) {
		bfdLiveDetect := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+
			"bfd-liveness-detection local-address "+bfdLiveDetect["local_address"].(string))
		if v2 := bfdLiveDetect["authentication_algorithm"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"bfd-liveness-detection authentication algorithm "+v2)
		}
		if v2 := bfdLiveDetect["authentication_key_chain"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"bfd-liveness-detection authentication key-chain "+v2)
		}
		if bfdLiveDetect["authentication_loose_check"].(bool) {
			configSet = append(configSet, setPrefix+"bfd-liveness-detection authentication loose-check")
		}
		if v2 := bfdLiveDetect["detection_time_threshold"].(int); v2 != 0 {
			configSet = append(configSet, setPrefix+"bfd-liveness-detection detection-time threshold "+strconv.Itoa(v2))
		}
		if v2 := bfdLiveDetect["holddown_interval"].(int); v2 != 0 {
			configSet = append(configSet, setPrefix+"bfd-liveness-detection holddown-interval "+strconv.Itoa(v2))
		}
		if v2 := bfdLiveDetect["minimum_interval"].(int); v2 != 0 {
			configSet = append(configSet, setPrefix+"bfd-liveness-detection minimum-interval "+strconv.Itoa(v2))
		}
		if v2 := bfdLiveDetect["minimum_receive_interval"].(int); v2 != 0 {
			configSet = append(configSet, setPrefix+"bfd-liveness-detection minimum-receive-interval "+strconv.Itoa(v2))
		}
		if v2 := bfdLiveDetect["multiplier"].(int); v2 != 0 {
			configSet = append(configSet, setPrefix+"bfd-liveness-detection multiplier "+strconv.Itoa(v2))
		}
		if v2 := bfdLiveDetect["neighbor"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"bfd-liveness-detection neighbor "+v2)
		}
		if bfdLiveDetect["no_adaptation"].(bool) {
			configSet = append(configSet, setPrefix+"bfd-liveness-detection no-adaptation")
		}
		if v2 := bfdLiveDetect["transmit_interval_minimum_interval"].(int); v2 != 0 {
			configSet = append(configSet, setPrefix+
				"bfd-liveness-detection transmit-interval minimum-interval "+strconv.Itoa(v2))
		}
		if v2 := bfdLiveDetect["transmit_interval_threshold"].(int); v2 != 0 {
			configSet = append(configSet, setPrefix+"bfd-liveness-detection transmit-interval threshold "+strconv.Itoa(v2))
		}
		if v2 := bfdLiveDetect["version"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"bfd-liveness-detection version "+v2)
		}
	}
	if ethOpts["flow_control"].(bool) {
		configSet = append(configSet, setPrefix+flowControlWords)
	}
	if ethOpts["no_flow_control"].(bool) {
		configSet = append(configSet, setPrefix+noFlowControlWords)
	}
	for _, v := range ethOpts["lacp"].([]interface{}) {
		lacp := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+"lacp "+lacp["mode"].(string))
		if v2 := lacp["admin_key"].(int); v2 != -1 {
			configSet = append(configSet, setPrefix+"lacp admin-key "+strconv.Itoa(v2))
		}
		if v2 := lacp["periodic"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"lacp periodic "+v2)
		}
		if v2 := lacp["sync_reset"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"lacp sync-reset "+v2)
		}
		if v2 := lacp["system_id"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"lacp system-id "+v2)
		}
		if v2 := lacp["system_priority"].(int); v2 != -1 {
			configSet = append(configSet, setPrefix+"lacp system-priority "+strconv.Itoa(v2))
		}
	}
	if ethOpts["loopback"].(bool) {
		configSet = append(configSet, setPrefix+loopbackWord)
	}
	if ethOpts["no_loopback"].(bool) {
		configSet = append(configSet, setPrefix+noLoopbackWord)
	}
	if v := ethOpts["link_speed"].(string); v != "" {
		configSet = append(configSet, setPrefix+"link-speed "+v)
	}
	if v := ethOpts["minimum_bandwidth"].(string); v != "" {
		vS := strings.Split(v, " ")
		configSet = append(configSet, setPrefix+"minimum-bandwidth bw-value "+vS[0])
		if len(vS) > 1 {
			configSet = append(configSet, setPrefix+"minimum-bandwidth bw-unit "+vS[1])
		}
	}
	if v := ethOpts["minimum_links"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"minimum-links "+strconv.Itoa(v))
	}
	if v := ethOpts["redundancy_group"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"redundancy-group "+strconv.Itoa(v))
	}
	for _, v := range ethOpts["source_address_filter"].([]interface{}) {
		configSet = append(configSet, setPrefix+"source-address-filter "+v.(string))
	}
	if ethOpts["source_filtering"].(bool) {
		configSet = append(configSet, setPrefix+"source-filtering")
	}

	return sess.configSet(configSet, jnprSess)
}

func readInterfacePhysical(interFace string, m interface{}, jnprSess *NetconfObject) (interfacePhysicalOptions, error) {
	sess := m.(*Session)
	var confRead interfacePhysicalOptions

	showConfig, err := sess.command("show configuration interfaces "+interFace+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, " unit ") && !strings.Contains(item, "ethernet-switching") {
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
			case strings.HasPrefix(itemTrim, "aggregated-ether-options lacp "):
				confRead.aeLacp = strings.TrimPrefix(itemTrim, "aggregated-ether-options lacp ")
				if err := readInterfacePhysicalParentEtherOpts(&confRead,
					strings.TrimPrefix(itemTrim, "aggregated-ether-options ")); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "aggregated-ether-options link-speed "):
				confRead.aeLinkSpeed = strings.TrimPrefix(itemTrim, "aggregated-ether-options link-speed ")
				if err := readInterfacePhysicalParentEtherOpts(&confRead,
					strings.TrimPrefix(itemTrim, "aggregated-ether-options ")); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "aggregated-ether-options minimum-links "):
				confRead.aeMinLink, err = strconv.Atoi(strings.TrimPrefix(itemTrim,
					"aggregated-ether-options minimum-links "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
				if err := readInterfacePhysicalParentEtherOpts(&confRead,
					strings.TrimPrefix(itemTrim, "aggregated-ether-options ")); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "aggregated-ether-options "):
				if err := readInterfacePhysicalParentEtherOpts(&confRead,
					strings.TrimPrefix(itemTrim, "aggregated-ether-options ")); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "redundant-ether-options "):
				if err := readInterfacePhysicalParentEtherOpts(&confRead,
					strings.TrimPrefix(itemTrim, "redundant-ether-options ")); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "esi "):
				if err := readInterfacePhysicalEsi(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "ether-options "):
				readInterfacePhysicalEtherOpts(&confRead, strings.TrimPrefix(itemTrim, "ether-options "))
			case strings.HasPrefix(itemTrim, "gigether-options "):
				readInterfacePhysicalGigetherOpts(&confRead, strings.TrimPrefix(itemTrim, "gigether-options "))
			case strings.HasPrefix(itemTrim, "native-vlan-id"):
				confRead.vlanNative, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "native-vlan-id "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case itemTrim == "unit 0 family ethernet-switching interface-mode trunk":
				confRead.trunk = true
			case strings.HasPrefix(itemTrim, "unit 0 family ethernet-switching vlan members"):
				confRead.vlanMembers = append(confRead.vlanMembers, strings.TrimPrefix(itemTrim,
					"unit 0 family ethernet-switching vlan members "))
			case itemTrim == "vlan-tagging":
				confRead.vlanTagging = true
			default:
				continue
			}
		}
	}

	return confRead, nil
}

func readInterfacePhysicalEsi(confRead *interfacePhysicalOptions, item string) error {
	itemTrim := strings.TrimPrefix(item, "esi ")
	if len(confRead.esi) == 0 {
		confRead.esi = append(confRead.esi, map[string]interface{}{
			"mode":             "",
			"auto_derive_lacp": false,
			"df_election_type": "",
			"identifier":       "",
			"source_bmac":      "",
		})
	}
	var err error
	identifier, err := regexp.MatchString(`^([\d\w]{2}:){9}[\d\w]{2}`, itemTrim)
	if err != nil {
		return fmt.Errorf("esi_identifier regexp error : %w", err)
	}
	switch {
	case identifier:
		confRead.esi[0]["identifier"] = itemTrim
	case itemTrim == "all-active" || itemTrim == "single-active":
		confRead.esi[0]["mode"] = itemTrim
	case strings.HasPrefix(itemTrim, "df-election-type "):
		confRead.esi[0]["df_election_type"] = strings.TrimPrefix(itemTrim, "df-election-type ")
	case strings.HasPrefix(itemTrim, "source-bmac "):
		confRead.esi[0]["source_bmac"] = strings.TrimPrefix(itemTrim, "source-bmac ")
	case itemTrim == "auto-derive lacp":
		confRead.esi[0]["auto_derive_lacp"] = true
	}

	return nil
}

func readInterfacePhysicalEtherOpts(confRead *interfacePhysicalOptions, itemTrim string) {
	if len(confRead.etherOpts) == 0 {
		confRead.etherOpts = append(confRead.etherOpts, map[string]interface{}{
			"ae_8023ad":           "",
			"auto_negotiation":    false,
			"no_auto_negotiation": false,
			"flow_control":        false,
			"no_flow_control":     false,
			"loopback":            false,
			"no_loopback":         false,
			"redundant_parent":    "",
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "802.3ad "):
		confRead.v8023ad = strings.TrimPrefix(itemTrim, "802.3ad ")
		confRead.etherOpts[0]["ae_8023ad"] = strings.TrimPrefix(itemTrim, "802.3ad ")
	case itemTrim == "auto-negotiation":
		confRead.etherOpts[0]["auto_negotiation"] = true
	case itemTrim == "no-auto-negotiation":
		confRead.etherOpts[0]["no_auto_negotiation"] = true
	case itemTrim == flowControlWords:
		confRead.etherOpts[0]["flow_control"] = true
	case itemTrim == noFlowControlWords:
		confRead.etherOpts[0]["no_flow_control"] = true
	case itemTrim == loopbackWord:
		confRead.etherOpts[0]["loopback"] = true
	case itemTrim == noLoopbackWord:
		confRead.etherOpts[0]["no_loopback"] = true
	case strings.HasPrefix(itemTrim, "redundant-parent "):
		confRead.etherOpts[0]["redundant_parent"] = strings.TrimPrefix(itemTrim, "redundant-parent ")
	}
}

func readInterfacePhysicalGigetherOpts(confRead *interfacePhysicalOptions, itemTrim string) {
	if len(confRead.gigetherOpts) == 0 {
		confRead.gigetherOpts = append(confRead.gigetherOpts, map[string]interface{}{
			"ae_8023ad":           "",
			"auto_negotiation":    false,
			"no_auto_negotiation": false,
			"flow_control":        false,
			"no_flow_control":     false,
			"loopback":            false,
			"no_loopback":         false,
			"redundant_parent":    "",
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "802.3ad "):
		confRead.v8023ad = strings.TrimPrefix(itemTrim, "802.3ad ")
		confRead.gigetherOpts[0]["ae_8023ad"] = strings.TrimPrefix(itemTrim, "802.3ad ")
	case itemTrim == "auto-negotiation":
		confRead.gigetherOpts[0]["auto_negotiation"] = true
	case itemTrim == "no-auto-negotiation":
		confRead.gigetherOpts[0]["no_auto_negotiation"] = true
	case itemTrim == flowControlWords:
		confRead.gigetherOpts[0]["flow_control"] = true
	case itemTrim == noFlowControlWords:
		confRead.gigetherOpts[0]["no_flow_control"] = true
	case itemTrim == loopbackWord:
		confRead.gigetherOpts[0]["loopback"] = true
	case itemTrim == noLoopbackWord:
		confRead.gigetherOpts[0]["no_loopback"] = true
	case strings.HasPrefix(itemTrim, "redundant-parent "):
		confRead.gigetherOpts[0]["redundant_parent"] = strings.TrimPrefix(itemTrim, "redundant-parent ")
	}
}

func readInterfacePhysicalParentEtherOpts(confRead *interfacePhysicalOptions, itemTrim string) error {
	if len(confRead.parentEtherOpts) == 0 {
		confRead.parentEtherOpts = append(confRead.parentEtherOpts, map[string]interface{}{
			"bfd_liveness_detection": make([]map[string]interface{}, 0),
			"flow_control":           false,
			"no_flow_control":        false,
			"lacp":                   make([]map[string]interface{}, 0),
			"loopback":               false,
			"no_loopback":            false,
			"link_speed":             "",
			"minimum_bandwidth":      "",
			"minimum_links":          0,
			"redundancy_group":       0,
			"source_address_filter":  make([]string, 0),
			"source_filtering":       false,
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "bfd-liveness-detection "):
		if len(confRead.parentEtherOpts[0]["bfd_liveness_detection"].([]map[string]interface{})) == 0 {
			confRead.parentEtherOpts[0]["bfd_liveness_detection"] = append(
				confRead.parentEtherOpts[0]["bfd_liveness_detection"].([]map[string]interface{}),
				map[string]interface{}{
					"local_address":                      "",
					"authentication_algorithm":           "",
					"authentication_key_chain":           "",
					"authentication_loose_check":         false,
					"detection_time_threshold":           0,
					"holddown_interval":                  0,
					"minimum_interval":                   0,
					"minimum_receive_interval":           0,
					"multiplier":                         0,
					"neighbor":                           "",
					"no_adaptation":                      false,
					"transmit_interval_minimum_interval": 0,
					"transmit_interval_threshold":        0,
					"version":                            "",
				})
		}
		parentEtherOptsBFDLiveDetect := confRead.parentEtherOpts[0]["bfd_liveness_detection"].([]map[string]interface{})[0]
		itemTrimBfdLiveDet := strings.TrimPrefix(itemTrim, "bfd-liveness-detection ")
		switch {
		case strings.HasPrefix(itemTrimBfdLiveDet, "local-address "):
			parentEtherOptsBFDLiveDetect["local_address"] = strings.TrimPrefix(itemTrimBfdLiveDet, "local-address ")
		case strings.HasPrefix(itemTrimBfdLiveDet, "authentication algorithm "):
			parentEtherOptsBFDLiveDetect["authentication_algorithm"] = strings.TrimPrefix(
				itemTrimBfdLiveDet, "authentication algorithm ")
		case strings.HasPrefix(itemTrimBfdLiveDet, "authentication key-chain "):
			parentEtherOptsBFDLiveDetect["authentication_key_chain"] = strings.TrimPrefix(
				itemTrimBfdLiveDet, "authentication key-chain ")
		case itemTrimBfdLiveDet == authenticationLooseCheck:
			parentEtherOptsBFDLiveDetect["authentication_loose_check"] = true
		case strings.HasPrefix(itemTrimBfdLiveDet, "detection-time threshold "):
			var err error
			parentEtherOptsBFDLiveDetect["detection_time_threshold"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrimBfdLiveDet, "detection-time threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimBfdLiveDet, "holddown-interval "):
			var err error
			parentEtherOptsBFDLiveDetect["holddown_interval"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrimBfdLiveDet, "holddown-interval "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimBfdLiveDet, "minimum-interval "):
			var err error
			parentEtherOptsBFDLiveDetect["minimum_interval"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrimBfdLiveDet, "minimum-interval "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimBfdLiveDet, "minimum-receive-interval "):
			var err error
			parentEtherOptsBFDLiveDetect["minimum_receive_interval"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrimBfdLiveDet, "minimum-receive-interval "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimBfdLiveDet, "multiplier "):
			var err error
			parentEtherOptsBFDLiveDetect["multiplier"], err = strconv.Atoi(strings.TrimPrefix(itemTrimBfdLiveDet, "multiplier "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimBfdLiveDet, "neighbor "):
			parentEtherOptsBFDLiveDetect["neighbor"] = strings.TrimPrefix(itemTrimBfdLiveDet, "neighbor ")
		case itemTrimBfdLiveDet == noAdaptation:
			parentEtherOptsBFDLiveDetect["no_adaptation"] = true
		case strings.HasPrefix(itemTrimBfdLiveDet, "transmit-interval minimum-interval "):
			var err error
			parentEtherOptsBFDLiveDetect["transmit_interval_minimum_interval"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrimBfdLiveDet, "transmit-interval minimum-interval "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimBfdLiveDet, "transmit-interval threshold "):
			var err error
			parentEtherOptsBFDLiveDetect["transmit_interval_threshold"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrimBfdLiveDet, "transmit-interval threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimBfdLiveDet, "version "):
			parentEtherOptsBFDLiveDetect["version"] = strings.TrimPrefix(itemTrimBfdLiveDet, "version ")
		}
	case itemTrim == flowControlWords:
		confRead.parentEtherOpts[0]["flow_control"] = true
	case itemTrim == noFlowControlWords:
		confRead.parentEtherOpts[0]["no_flow_control"] = true
	case strings.HasPrefix(itemTrim, "lacp "):
		if len(confRead.parentEtherOpts[0]["lacp"].([]map[string]interface{})) == 0 {
			confRead.parentEtherOpts[0]["lacp"] = append(confRead.parentEtherOpts[0]["lacp"].([]map[string]interface{}),
				map[string]interface{}{
					"mode":            "",
					"admin_key":       -1,
					"periodic":        "",
					"sync_reset":      "",
					"system_id":       "",
					"system_priority": -1,
				})
		}
		itemTrimLacp := strings.TrimPrefix(itemTrim, "lacp ")
		lacp := confRead.parentEtherOpts[0]["lacp"].([]map[string]interface{})[0]
		switch {
		case itemTrimLacp == activeW || itemTrimLacp == "passive":
			lacp["mode"] = itemTrimLacp
		case strings.HasPrefix(itemTrimLacp, "admin-key "):
			var err error
			lacp["admin_key"], err = strconv.Atoi(strings.TrimPrefix(itemTrimLacp, "admin-key "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimLacp, "periodic "):
			lacp["periodic"] = strings.TrimPrefix(itemTrimLacp, "periodic ")
		case strings.HasPrefix(itemTrimLacp, "sync-reset "):
			lacp["sync_reset"] = strings.TrimPrefix(itemTrimLacp, "sync-reset ")
		case strings.HasPrefix(itemTrimLacp, "system-id "):
			lacp["system_id"] = strings.TrimPrefix(itemTrimLacp, "system-id ")
		case strings.HasPrefix(itemTrimLacp, "system-priority "):
			var err error
			lacp["system_priority"], err = strconv.Atoi(strings.TrimPrefix(itemTrimLacp, "system-priority "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case itemTrim == loopbackWord:
		confRead.parentEtherOpts[0]["loopback"] = true
	case itemTrim == noLoopbackWord:
		confRead.parentEtherOpts[0]["no_loopback"] = true
	case strings.HasPrefix(itemTrim, "link-speed "):
		confRead.parentEtherOpts[0]["link_speed"] = strings.TrimPrefix(itemTrim, "link-speed ")
	case strings.HasPrefix(itemTrim, "minimum-bandwidth bw-value "):
		confRead.parentEtherOpts[0]["minimum_bandwidth"] = strings.TrimPrefix(itemTrim, "minimum-bandwidth bw-value ") +
			confRead.parentEtherOpts[0]["minimum_bandwidth"].(string)
	case strings.HasPrefix(itemTrim, "minimum-bandwidth bw-unit "):
		confRead.parentEtherOpts[0]["minimum_bandwidth"] = confRead.parentEtherOpts[0]["minimum_bandwidth"].(string) +
			" " + strings.TrimPrefix(itemTrim, "minimum-bandwidth bw-unit ")
	case strings.HasPrefix(itemTrim, "minimum-links "):
		var err error
		confRead.parentEtherOpts[0]["minimum_links"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "minimum-links "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "redundancy-group "):
		var err error
		confRead.parentEtherOpts[0]["redundancy_group"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "redundancy-group "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "source-address-filter "):
		confRead.parentEtherOpts[0]["source_address_filter"] = append(
			confRead.parentEtherOpts[0]["source_address_filter"].([]string),
			strings.TrimPrefix(itemTrim, "source-address-filter "))
	case itemTrim == "source-filtering":
		confRead.parentEtherOpts[0]["source_filtering"] = true
	}

	return nil
}

func delInterfacePhysical(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	if jnprSess != nil {
		if containsUnit, err := checkInterfacePhysicalContainsUnit(d.Get("name").(string), m, jnprSess); err != nil {
			return err
		} else if containsUnit {
			return fmt.Errorf("interface %s is used for a logical unit interface", d.Get("name").(string))
		}
	}
	if err := sess.configSet([]string{"delete interfaces " + d.Get("name").(string)}, jnprSess); err != nil {
		return err
	}
	if jnprSess == nil {
		return nil
	}
	if v := d.Get("name").(string); strings.HasPrefix(v, "ae") {
		aggregatedCount, err := interfaceAggregatedCountSearchMax("ae-1", v, v, m, jnprSess)
		if err != nil {
			return err
		}
		if aggregatedCount == "0" {
			err = sess.configSet([]string{"delete chassis aggregated-devices ethernet device-count"}, jnprSess)
			if err != nil {
				return err
			}
		} else {
			err = sess.configSet([]string{"set chassis aggregated-devices ethernet device-count " + aggregatedCount}, jnprSess)
			if err != nil {
				return err
			}
		}
	} else if d.Get("ether802_3ad").(string) != "" ||
		len(d.Get("ether_opts").([]interface{})) != 0 ||
		len(d.Get("gigether_opts").([]interface{})) != 0 {
		var aeDel string
		switch {
		case d.Get("ether802_3ad").(string) != "":
			aeDel = d.Get("ether802_3ad").(string)
		case len(d.Get("ether_opts").([]interface{})) != 0 && d.Get("ether_opts").([]interface{})[0] != nil:
			v := d.Get("ether_opts").([]interface{})[0].(map[string]interface{})
			aeDel = v["ae_8023ad"].(string)
		case len(d.Get("gigether_opts").([]interface{})) != 0 && d.Get("gigether_opts").([]interface{})[0] != nil:
			v := d.Get("gigether_opts").([]interface{})[0].(map[string]interface{})
			aeDel = v["ae_8023ad"].(string)
		}
		if aeDel != "" {
			lastAEchild, err := interfaceAggregatedLastChild(aeDel, d.Get("name").(string), m, jnprSess)
			if err != nil {
				return err
			}
			if lastAEchild {
				aggregatedCount, err := interfaceAggregatedCountSearchMax("ae-1", aeDel,
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
					err = sess.configSet([]string{"set chassis aggregated-devices ethernet device-count " + aggregatedCount}, jnprSess)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func checkInterfacePhysicalContainsUnit(interFace string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration interfaces "+interFace+" | display set relative", jnprSess)
	if err != nil {
		return false, err
	}
	for _, item := range strings.Split(showConfig, "\n") {
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

			return true, nil
		}
	}

	return false, nil
}

func delInterfaceNC(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	delPrefix := "delete interfaces " + d.Get("name").(string) + " "
	if sess.junosGroupIntDel != "" {
		configSet = append(configSet, delPrefix+"apply-groups "+sess.junosGroupIntDel)
	}
	configSet = append(configSet, delPrefix+"description")
	configSet = append(configSet, delPrefix+"disable")

	return sess.configSet(configSet, jnprSess)
}

func delInterfacePhysicalOpts(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	delPrefix := "delete interfaces " + d.Get("name").(string) + " "
	configSet = append(configSet,
		delPrefix+"aggregated-ether-options",
		delPrefix+"description",
		delPrefix+"esi",
		delPrefix+"ether-options",
		delPrefix+"gigether-options",
		delPrefix+"native-vlan-id",
		delPrefix+"redundant-ether-options",
		delPrefix+"unit 0 family ethernet-switching interface-mode",
		delPrefix+"unit 0 family ethernet-switching vlan members",
		delPrefix+"vlan-tagging",
	)

	return sess.configSet(configSet, jnprSess)
}

func fillInterfacePhysicalData(d *schema.ResourceData, interfaceOpt interfacePhysicalOptions) {
	_, okAeLacp := d.GetOk("ae_lacp")
	if okAeLacp {
		if tfErr := d.Set("ae_lacp", interfaceOpt.aeLacp); tfErr != nil {
			panic(tfErr)
		}
	}
	_, okAeLinkSpeed := d.GetOk("ae_link_speed")
	if okAeLinkSpeed {
		if tfErr := d.Set("ae_link_speed", interfaceOpt.aeLinkSpeed); tfErr != nil {
			panic(tfErr)
		}
	}
	_, okAeMinLinks := d.GetOk("ae_minimum_links")
	if okAeMinLinks {
		if tfErr := d.Set("ae_minimum_links", interfaceOpt.aeMinLink); tfErr != nil {
			panic(tfErr)
		}
	}
	if tfErr := d.Set("esi", interfaceOpt.esi); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", interfaceOpt.description); tfErr != nil {
		panic(tfErr)
	}
	if _, ok := d.GetOk("ether802_3ad"); ok {
		if tfErr := d.Set("ether802_3ad", interfaceOpt.v8023ad); tfErr != nil {
			panic(tfErr)
		}
	} else {
		if tfErr := d.Set("ether_opts", interfaceOpt.etherOpts); tfErr != nil {
			panic(tfErr)
		}
		if tfErr := d.Set("gigether_opts", interfaceOpt.gigetherOpts); tfErr != nil {
			panic(tfErr)
		}
	}
	if !okAeLacp && !okAeLinkSpeed && !okAeMinLinks {
		if tfErr := d.Set("parent_ether_opts", interfaceOpt.parentEtherOpts); tfErr != nil {
			panic(tfErr)
		}
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
	if tfErr := d.Set("vlan_tagging", interfaceOpt.vlanTagging); tfErr != nil {
		panic(tfErr)
	}
}

func interfaceAggregatedLastChild(ae, interFace string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration interfaces | display set relative", jnprSess)
	if err != nil {
		return false, err
	}
	lastAE := true
	for _, item := range strings.Split(showConfig, "\n") {
		if strings.HasSuffix(item, "ether-options 802.3ad "+ae) &&
			!strings.HasPrefix(item, "set "+interFace+" ") {
			lastAE = false
		}
	}

	return lastAE, nil
}

func interfaceAggregatedCountSearchMax(
	newAE, oldAE, interFace string, m interface{}, jnprSess *NetconfObject) (string, error) {
	sess := m.(*Session)
	newAENum := strings.TrimPrefix(newAE, "ae")
	newAENumInt, err := strconv.Atoi(newAENum)
	if err != nil {
		return "", fmt.Errorf("failed to convert ae interaface '%v' to integer : %w", newAE, err)
	}
	showConfig, err := sess.command("show configuration interfaces | display set relative", jnprSess)
	if err != nil {
		return "", err
	}
	listAEFound := make([]string, 0)
	regexpAEchild := regexp.MustCompile(`ether-options 802\.3ad ae\d+$`)
	regexpAEparent := regexp.MustCompile(`^set ae\d+ `)
	for _, line := range strings.Split(showConfig, "\n") {
		aeMatchChild := regexpAEchild.MatchString(line)
		aeMatchParent := regexpAEparent.MatchString(line)
		switch {
		case aeMatchChild:
			wordsLine := strings.Fields(line)
			if interFace == oldAE {
				// interfaceAggregatedCountSearchMax called for delete parent interface
				listAEFound = append(listAEFound, wordsLine[len(wordsLine)-1])
			} else if wordsLine[len(wordsLine)-1] != oldAE {
				listAEFound = append(listAEFound, wordsLine[len(wordsLine)-1])
			}
		case aeMatchParent:
			wordsLine := strings.Fields(line)
			if interFace != oldAE {
				// interfaceAggregatedCountSearchMax called for child interface or new parent
				listAEFound = append(listAEFound, wordsLine[1])
			} else if wordsLine[1] != oldAE {
				listAEFound = append(listAEFound, wordsLine[1])
			}
		}
	}
	lastOldAE, err := interfaceAggregatedLastChild(oldAE, interFace, m, jnprSess)
	if err != nil {
		return "", err
	}
	if !lastOldAE {
		listAEFound = append(listAEFound, oldAE)
	}
	if len(listAEFound) > 0 {
		balt.SortStringsByLengthInc(listAEFound)
		lastAeInt, err := strconv.Atoi(strings.TrimPrefix(listAEFound[len(listAEFound)-1], "ae"))
		if err != nil {
			return "", fmt.Errorf("failed to convert internal variable lastAeInt to integer : %w", err)
		}
		if lastAeInt > newAENumInt {
			return strconv.Itoa(lastAeInt + 1), nil
		}
	}

	return strconv.Itoa(newAENumInt + 1), nil
}
