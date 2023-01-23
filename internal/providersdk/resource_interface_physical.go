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
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type interfacePhysicalOptions struct {
	disable         bool
	trunk           bool
	vlanTagging     bool
	aeMinLink       int
	mtu             int
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
		CreateWithoutTimeout: resourceInterfacePhysicalCreate,
		ReadWithoutTimeout:   resourceInterfacePhysicalRead,
		UpdateWithoutTimeout: resourceInterfacePhysicalUpdate,
		DeleteWithoutTimeout: resourceInterfacePhysicalDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceInterfacePhysicalImport,
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
			"disable": {
				Type:     schema.TypeBool,
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
			"mtu": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 9500),
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := delInterfaceNC(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setInterfacePhysical(d, clt, nil); err != nil {
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
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ncInt, emptyInt, err := checkInterfacePhysicalNCEmpty(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !ncInt && !emptyInt {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("interface %s already configured", d.Get("name").(string)))...)
	}
	if ncInt {
		if err := delInterfaceNC(d, clt, junSess); err != nil {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	if err := setInterfacePhysical(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_interface_physical", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ncInt, emptyInt, err = checkInterfacePhysicalNCEmpty(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ncInt {
		return append(diagWarns, diag.FromErr(fmt.Errorf("interface %v always disable (NC) after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}
	if emptyInt {
		intExists, err := checkInterfaceExists(d.Get("name").(string), clt, junSess)
		if err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		if !intExists {
			return append(diagWarns, diag.FromErr(fmt.Errorf("interface %v not exists and "+
				"config can't found after commit => check your config", d.Get("name").(string)))...)
		}
	}
	d.SetId(d.Get("name").(string))

	return append(diagWarns, resourceInterfacePhysicalReadWJunSess(d, clt, junSess)...)
}

func resourceInterfacePhysicalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceInterfacePhysicalReadWJunSess(d, clt, junSess)
}

func resourceInterfacePhysicalReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	ncInt, emptyInt, err := checkInterfacePhysicalNCEmpty(d.Get("name").(string), clt, junSess)
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
		intExists, err := checkInterfaceExists(d.Get("name").(string), clt, junSess)
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
	interfaceOpt, err := readInterfacePhysical(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillInterfacePhysicalData(d, interfaceOpt)

	return nil
}

func resourceInterfacePhysicalUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delInterfacePhysicalOpts(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setInterfacePhysical(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delInterfacePhysicalOpts(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := unsetInterfacePhysicalAE(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setInterfacePhysical(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_interface_physical", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceInterfacePhysicalReadWJunSess(d, clt, junSess)...)
}

func resourceInterfacePhysicalDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delInterfacePhysical(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delInterfacePhysical(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_interface_physical", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !d.Get("no_disable_on_destroy").(bool) {
		intExists, err := checkInterfaceExists(d.Get("name").(string), clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, []error{err})
		} else if intExists {
			err = addInterfacePhysicalNC(d.Get("name").(string), clt, junSess)
			if err != nil {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
			_, err = clt.CommitConf("disable(NC) resource junos_interface_physical", junSess)
			if err != nil {
				appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
		}
	}

	return diagWarns
}

func resourceInterfacePhysicalImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	if strings.Count(d.Id(), ".") != 0 {
		return nil, fmt.Errorf("name of interface %s need to doesn't have a dot", d.Id())
	}
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)
	ncInt, emptyInt, err := checkInterfacePhysicalNCEmpty(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if ncInt {
		return nil, fmt.Errorf("interface '%v' is disabled (NC), import is not possible", d.Id())
	}
	if emptyInt {
		intExists, err := checkInterfaceExists(d.Id(), clt, junSess)
		if err != nil {
			return nil, err
		}
		if !intExists {
			return nil, fmt.Errorf("don't find interface with id '%v' (id must be <name>)", d.Id())
		}
	}
	interfaceOpt, err := readInterfacePhysical(d.Id(), clt, junSess)
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

func checkInterfacePhysicalNCEmpty(interFace string, clt *junos.Client, junSess *junos.Session,
) (ncInt, emtyInt bool, errFunc error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+"interfaces "+interFace+junos.PipeDisplaySetRelative, junSess)
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
	if showConfig == junos.EmptyW {
		return false, true, nil
	}

	return false, false, nil
}

func addInterfacePhysicalNC(interFace string, clt *junos.Client, junSess *junos.Session) (err error) {
	if clt.GroupInterfaceDelete() == "" {
		err = clt.ConfigSet([]string{"set interfaces " + interFace + " disable description NC"}, junSess)
	} else {
		err = clt.ConfigSet([]string{"set interfaces " + interFace +
			" apply-groups " + clt.GroupInterfaceDelete()}, junSess)
	}
	if err != nil {
		return err
	}

	return nil
}

func checkInterfaceExists(interFace string, clt *junos.Client, junSess *junos.Session) (bool, error) {
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

func unsetInterfacePhysicalAE(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
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
		aggregatedCount, err := interfaceAggregatedCountSearchMax("ae-1", oldAE, d.Get("name").(string), clt, junSess)
		if err != nil {
			return err
		}
		if aggregatedCount == "0" {
			return clt.ConfigSet([]string{"delete chassis aggregated-devices ethernet device-count"}, junSess)
		}

		return clt.ConfigSet([]string{"set chassis aggregated-devices ethernet device-count " + aggregatedCount}, junSess)
	}

	return nil
}

func setInterfacePhysical(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
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
	if d.Get("disable").(bool) {
		if d.Get("description").(string) == "NC" {
			return fmt.Errorf("disable=true and description=NC is not allowed " +
				"because the provider might consider the resource deleted")
		}
		configSet = append(configSet, setPrefix+"disable")
	}
	if err := setInterfacePhysicalEsi(setPrefix, d.Get("esi").([]interface{}), clt, junSess); err != nil {
		return err
	}
	if v := d.Get("name").(string); strings.HasPrefix(v, "ae") && junSess != nil {
		aggregatedCount, err := interfaceAggregatedCountSearchMax(v, "ae-1", v, clt, junSess)
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
		if newAE != "" && junSess != nil {
			aggregatedCount, err := interfaceAggregatedCountSearchMax(
				newAE,
				oldAE,
				d.Get("name").(string),
				clt, junSess)
			if err != nil {
				return err
			}
			configSet = append(configSet, "set chassis aggregated-devices ethernet device-count "+aggregatedCount)
		}
	}
	if v := d.Get("mtu").(int); v != 0 {
		configSet = append(configSet, setPrefix+"mtu "+strconv.Itoa(v))
	}
	for _, v := range d.Get("parent_ether_opts").([]interface{}) {
		if v == nil {
			return fmt.Errorf("parent_ether_opts block is empty")
		}
		if err := setInterfacePhysicalParentEtherOpts(
			v.(map[string]interface{}),
			d.Get("name").(string),
			clt, junSess,
		); err != nil {
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

	return clt.ConfigSet(configSet, junSess)
}

func setInterfacePhysicalEsi(setPrefix string, esiParams []interface{}, clt *junos.Client, junSess *junos.Session,
) error {
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

	return clt.ConfigSet(configSet, junSess)
}

func setInterfacePhysicalParentEtherOpts(
	ethOpts map[string]interface{}, interfaceName string, clt *junos.Client, junSess *junos.Session,
) error {
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
		configSet = append(configSet, setPrefix+"flow-control")
	}
	if ethOpts["no_flow_control"].(bool) {
		configSet = append(configSet, setPrefix+"no-flow-control")
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
		configSet = append(configSet, setPrefix+"loopback")
	}
	if ethOpts["no_loopback"].(bool) {
		configSet = append(configSet, setPrefix+"no-loopback")
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

	return clt.ConfigSet(configSet, junSess)
}

func readInterfacePhysical(interFace string, clt *junos.Client, junSess *junos.Session,
) (confRead interfacePhysicalOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+"interfaces "+interFace+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, " unit ") && !strings.Contains(item, "ethernet-switching") {
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
			case balt.CutPrefixInString(&itemTrim, "aggregated-ether-options "):
				itemTrimToLegacy := itemTrim
				switch {
				case balt.CutPrefixInString(&itemTrimToLegacy, "lacp "):
					confRead.aeLacp = itemTrimToLegacy
				case balt.CutPrefixInString(&itemTrimToLegacy, "link-speed "):
					confRead.aeLinkSpeed = itemTrimToLegacy
				case balt.CutPrefixInString(&itemTrimToLegacy, "minimum-links "):
					confRead.aeMinLink, err = strconv.Atoi(itemTrimToLegacy)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrimToLegacy, err)
					}
				}
				if err := confRead.readInterfacePhysicalParentEtherOpts(itemTrim); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "redundant-ether-options "):
				if err := confRead.readInterfacePhysicalParentEtherOpts(itemTrim); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			case itemTrim == "disable":
				confRead.disable = true
			case balt.CutPrefixInString(&itemTrim, "esi "):
				if err := confRead.readInterfacePhysicalEsi(itemTrim); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "ether-options "):
				confRead.readInterfacePhysicalEtherOpts(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "gigether-options "):
				confRead.readInterfacePhysicalGigetherOpts(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "mtu "):
				confRead.mtu, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "native-vlan-id "):
				confRead.vlanNative, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "unit 0 family ethernet-switching interface-mode trunk":
				confRead.trunk = true
			case balt.CutPrefixInString(&itemTrim, "unit 0 family ethernet-switching vlan members "):
				confRead.vlanMembers = append(confRead.vlanMembers, itemTrim)
			case itemTrim == "vlan-tagging":
				confRead.vlanTagging = true
			default:
				continue
			}
		}
	}

	return confRead, nil
}

func (confRead *interfacePhysicalOptions) readInterfacePhysicalEsi(itemTrim string) error {
	if len(confRead.esi) == 0 {
		confRead.esi = append(confRead.esi, map[string]interface{}{
			"mode":             "",
			"auto_derive_lacp": false,
			"df_election_type": "",
			"identifier":       "",
			"source_bmac":      "",
		})
	}
	identifier, err := regexp.MatchString(`^([\d\w]{2}:){9}[\d\w]{2}`, itemTrim)
	if err != nil {
		return fmt.Errorf("esi_identifier regexp error: %w", err)
	}
	switch {
	case identifier:
		confRead.esi[0]["identifier"] = itemTrim
	case itemTrim == "all-active", itemTrim == "single-active":
		confRead.esi[0]["mode"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "df-election-type "):
		confRead.esi[0]["df_election_type"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "source-bmac "):
		confRead.esi[0]["source_bmac"] = itemTrim
	case itemTrim == "auto-derive lacp":
		confRead.esi[0]["auto_derive_lacp"] = true
	}

	return nil
}

func (confRead *interfacePhysicalOptions) readInterfacePhysicalEtherOpts(itemTrim string) {
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
	case balt.CutPrefixInString(&itemTrim, "802.3ad "):
		confRead.v8023ad = itemTrim
		confRead.etherOpts[0]["ae_8023ad"] = itemTrim
	case itemTrim == "auto-negotiation":
		confRead.etherOpts[0]["auto_negotiation"] = true
	case itemTrim == "no-auto-negotiation":
		confRead.etherOpts[0]["no_auto_negotiation"] = true
	case itemTrim == "flow-control":
		confRead.etherOpts[0]["flow_control"] = true
	case itemTrim == "no-flow-control":
		confRead.etherOpts[0]["no_flow_control"] = true
	case itemTrim == "loopback":
		confRead.etherOpts[0]["loopback"] = true
	case itemTrim == "no-loopback":
		confRead.etherOpts[0]["no_loopback"] = true
	case balt.CutPrefixInString(&itemTrim, "redundant-parent "):
		confRead.etherOpts[0]["redundant_parent"] = itemTrim
	}
}

func (confRead *interfacePhysicalOptions) readInterfacePhysicalGigetherOpts(itemTrim string) {
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
	case balt.CutPrefixInString(&itemTrim, "802.3ad "):
		confRead.v8023ad = itemTrim
		confRead.gigetherOpts[0]["ae_8023ad"] = itemTrim
	case itemTrim == "auto-negotiation":
		confRead.gigetherOpts[0]["auto_negotiation"] = true
	case itemTrim == "no-auto-negotiation":
		confRead.gigetherOpts[0]["no_auto_negotiation"] = true
	case itemTrim == "flow-control":
		confRead.gigetherOpts[0]["flow_control"] = true
	case itemTrim == "no-flow-control":
		confRead.gigetherOpts[0]["no_flow_control"] = true
	case itemTrim == "loopback":
		confRead.gigetherOpts[0]["loopback"] = true
	case itemTrim == "no-loopback":
		confRead.gigetherOpts[0]["no_loopback"] = true
	case balt.CutPrefixInString(&itemTrim, "redundant-parent "):
		confRead.gigetherOpts[0]["redundant_parent"] = itemTrim
	}
}

func (confRead *interfacePhysicalOptions) readInterfacePhysicalParentEtherOpts(itemTrim string) (err error) {
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
	case balt.CutPrefixInString(&itemTrim, "bfd-liveness-detection "):
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
		switch {
		case balt.CutPrefixInString(&itemTrim, "local-address "):
			parentEtherOptsBFDLiveDetect["local_address"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "authentication algorithm "):
			parentEtherOptsBFDLiveDetect["authentication_algorithm"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "authentication key-chain "):
			parentEtherOptsBFDLiveDetect["authentication_key_chain"] = itemTrim
		case itemTrim == "authentication loose-check":
			parentEtherOptsBFDLiveDetect["authentication_loose_check"] = true
		case balt.CutPrefixInString(&itemTrim, "detection-time threshold "):
			parentEtherOptsBFDLiveDetect["detection_time_threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "holddown-interval "):
			parentEtherOptsBFDLiveDetect["holddown_interval"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-interval "):
			parentEtherOptsBFDLiveDetect["minimum_interval"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-receive-interval "):
			parentEtherOptsBFDLiveDetect["minimum_receive_interval"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "multiplier "):
			parentEtherOptsBFDLiveDetect["multiplier"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "neighbor "):
			parentEtherOptsBFDLiveDetect["neighbor"] = itemTrim
		case itemTrim == "no-adaptation":
			parentEtherOptsBFDLiveDetect["no_adaptation"] = true
		case balt.CutPrefixInString(&itemTrim, "transmit-interval minimum-interval "):
			parentEtherOptsBFDLiveDetect["transmit_interval_minimum_interval"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "transmit-interval threshold "):
			parentEtherOptsBFDLiveDetect["transmit_interval_threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "version "):
			parentEtherOptsBFDLiveDetect["version"] = itemTrim
		}
	case itemTrim == "flow-control":
		confRead.parentEtherOpts[0]["flow_control"] = true
	case itemTrim == "no-flow-control":
		confRead.parentEtherOpts[0]["no_flow_control"] = true
	case balt.CutPrefixInString(&itemTrim, "lacp "):
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
		lacp := confRead.parentEtherOpts[0]["lacp"].([]map[string]interface{})[0]
		switch {
		case itemTrim == "active", itemTrim == "passive":
			lacp["mode"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "admin-key "):
			lacp["admin_key"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "periodic "):
			lacp["periodic"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "sync-reset "):
			lacp["sync_reset"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "system-id "):
			lacp["system_id"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "system-priority "):
			lacp["system_priority"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		}
	case itemTrim == "loopback":
		confRead.parentEtherOpts[0]["loopback"] = true
	case itemTrim == "no-loopback":
		confRead.parentEtherOpts[0]["no_loopback"] = true
	case balt.CutPrefixInString(&itemTrim, "link-speed "):
		confRead.parentEtherOpts[0]["link_speed"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "minimum-bandwidth bw-value "):
		confRead.parentEtherOpts[0]["minimum_bandwidth"] = itemTrim +
			confRead.parentEtherOpts[0]["minimum_bandwidth"].(string)
	case balt.CutPrefixInString(&itemTrim, "minimum-bandwidth bw-unit "):
		confRead.parentEtherOpts[0]["minimum_bandwidth"] = confRead.parentEtherOpts[0]["minimum_bandwidth"].(string) +
			" " + itemTrim
	case balt.CutPrefixInString(&itemTrim, "minimum-links "):
		confRead.parentEtherOpts[0]["minimum_links"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "redundancy-group "):
		confRead.parentEtherOpts[0]["redundancy_group"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "source-address-filter "):
		confRead.parentEtherOpts[0]["source_address_filter"] = append(
			confRead.parentEtherOpts[0]["source_address_filter"].([]string),
			itemTrim,
		)
	case itemTrim == "source-filtering":
		confRead.parentEtherOpts[0]["source_filtering"] = true
	}

	return nil
}

func delInterfacePhysical(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	if junSess != nil {
		if containsUnit, err := checkInterfacePhysicalContainsUnit(d.Get("name").(string), clt, junSess); err != nil {
			return err
		} else if containsUnit {
			return fmt.Errorf("interface %s is used for a logical unit interface", d.Get("name").(string))
		}
	}
	if err := clt.ConfigSet([]string{"delete interfaces " + d.Get("name").(string)}, junSess); err != nil {
		return err
	}
	if junSess == nil {
		return nil
	}
	if v := d.Get("name").(string); strings.HasPrefix(v, "ae") {
		aggregatedCount, err := interfaceAggregatedCountSearchMax("ae-1", v, v, clt, junSess)
		if err != nil {
			return err
		}
		if aggregatedCount == "0" {
			err = clt.ConfigSet([]string{"delete chassis aggregated-devices ethernet device-count"}, junSess)
			if err != nil {
				return err
			}
		} else {
			err = clt.ConfigSet([]string{"set chassis aggregated-devices ethernet device-count " + aggregatedCount}, junSess)
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
			lastAEchild, err := interfaceAggregatedLastChild(aeDel, d.Get("name").(string), clt, junSess)
			if err != nil {
				return err
			}
			if lastAEchild {
				aggregatedCount, err := interfaceAggregatedCountSearchMax(
					"ae-1",
					aeDel,
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
					err = clt.ConfigSet([]string{"set chassis aggregated-devices ethernet device-count " + aggregatedCount}, junSess)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func checkInterfacePhysicalContainsUnit(interFace string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+"interfaces "+interFace+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return false, err
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

			return true, nil
		}
	}

	return false, nil
}

func delInterfaceNC(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 3)
	delPrefix := "delete interfaces " + d.Get("name").(string) + " "
	if clt.GroupInterfaceDelete() != "" {
		configSet = append(configSet, delPrefix+"apply-groups "+clt.GroupInterfaceDelete())
	}
	configSet = append(configSet, delPrefix+"description")
	configSet = append(configSet, delPrefix+"disable")

	return clt.ConfigSet(configSet, junSess)
}

func delInterfacePhysicalOpts(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	delPrefix := "delete interfaces " + d.Get("name").(string) + " "
	configSet = append(configSet,
		delPrefix+"aggregated-ether-options",
		delPrefix+"description",
		delPrefix+"disable",
		delPrefix+"esi",
		delPrefix+"ether-options",
		delPrefix+"gigether-options",
		delPrefix+"native-vlan-id",
		delPrefix+"redundant-ether-options",
		delPrefix+"unit 0 family ethernet-switching interface-mode",
		delPrefix+"unit 0 family ethernet-switching vlan members",
		delPrefix+"vlan-tagging",
	)

	return clt.ConfigSet(configSet, junSess)
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
	if tfErr := d.Set("disable", interfaceOpt.disable); tfErr != nil {
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
	if tfErr := d.Set("mtu", interfaceOpt.mtu); tfErr != nil {
		panic(tfErr)
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

func interfaceAggregatedLastChild(ae, interFace string, clt *junos.Client, junSess *junos.Session) (bool, error) {
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

func interfaceAggregatedCountSearchMax(newAE, oldAE, interFace string, clt *junos.Client, junSess *junos.Session,
) (string, error) {
	newAENum := strings.TrimPrefix(newAE, "ae")
	newAENumInt, err := strconv.Atoi(newAENum)
	if err != nil {
		return "", fmt.Errorf("failed to convert ae interaface '%v' to integer: %w", newAE, err)
	}
	showConfig, err := clt.Command(junos.CmdShowConfig+"interfaces"+junos.PipeDisplaySetRelative, junSess)
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
	lastOldAE, err := interfaceAggregatedLastChild(oldAE, interFace, clt, junSess)
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
			return "", fmt.Errorf("failed to convert internal variable lastAeInt to integer: %w", err)
		}
		if lastAeInt > newAENumInt {
			return strconv.Itoa(lastAeInt + 1), nil
		}
	}

	return strconv.Itoa(newAENumInt + 1), nil
}
