package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type samplingInstanceOptions struct {
	disable           bool
	name              string
	familyInetInput   []map[string]interface{}
	familyInetOutput  []map[string]interface{}
	familyInet6Input  []map[string]interface{}
	familyInet6Output []map[string]interface{}
	familyMplsInput   []map[string]interface{}
	familyMplsOutput  []map[string]interface{}
	input             []map[string]interface{}
}

func resourceForwardingOptionsSamplingInstance() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceForwardingOptionsSamplingInstanceCreate,
		ReadWithoutTimeout:   resourceForwardingOptionsSamplingInstanceRead,
		UpdateWithoutTimeout: resourceForwardingOptionsSamplingInstanceUpdate,
		DeleteWithoutTimeout: resourceForwardingOptionsSamplingInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceForwardingOptionsSamplingInstanceImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"disable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"family_inet_input": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"family_inet_input", "family_inet6_input", "family_mpls_input", "input"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_packets_per_second": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"maximum_packet_length": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 9192),
						},
						"rate": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 16000000),
						},
						"run_length": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 20),
						},
					},
				},
			},
			"family_inet_output": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"family_inet_output", "family_inet6_output", "family_mpls_output"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aggregate_export_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(90, 1800),
						},
						"extension_service": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"flow_active_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(60, 1800),
						},
						"flow_inactive_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(15, 1800),
						},
						"flow_server": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"aggregation_autonomous_system": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_destination_prefix": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_protocol_port": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_source_destination_prefix": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_source_destination_prefix_caida_compliant": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_source_prefix": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"autonomous_system_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"origin", "peer"}, false),
									},
									"dscp": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 63),
									},
									"forwarding_class": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"local_dump": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"no_local_dump": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"routing_instance": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"source_address": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"version": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntInSlice([]int{5, 8}),
									},
									"version_ipfix_template": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"version9_template": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"inline_jflow_export_rate": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 3200),
							RequiredWith: []string{"family_inet_output.0.inline_jflow_source_address"},
						},
						"inline_jflow_source_address": {
							Type:         schema.TypeString,
							Optional:     true,
							RequiredWith: []string{"family_inet_output.0.flow_server"},
						},
						"interface": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"engine_id": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 255),
									},
									"engine_type": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 255),
									},
									"source_address": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"family_inet6_input": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"family_inet_input", "family_inet6_input", "family_mpls_input", "input"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_packets_per_second": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"maximum_packet_length": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 9192),
						},
						"rate": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 16000000),
						},
						"run_length": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 20),
						},
					},
				},
			},
			"family_inet6_output": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"family_inet_output", "family_inet6_output", "family_mpls_output"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aggregate_export_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(90, 1800),
						},
						"extension_service": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"flow_active_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(60, 1800),
						},
						"flow_inactive_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(15, 1800),
						},
						"flow_server": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"aggregation_autonomous_system": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_destination_prefix": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_protocol_port": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_source_destination_prefix": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_source_destination_prefix_caida_compliant": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_source_prefix": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"autonomous_system_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"origin", "peer"}, false),
									},
									"dscp": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 63),
									},
									"forwarding_class": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"local_dump": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"no_local_dump": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"routing_instance": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"source_address": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"version_ipfix_template": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"version9_template": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"inline_jflow_export_rate": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 3200),
							RequiredWith: []string{"family_inet6_output.0.inline_jflow_source_address"},
						},
						"inline_jflow_source_address": {
							Type:         schema.TypeString,
							Optional:     true,
							RequiredWith: []string{"family_inet6_output.0.flow_server"},
						},
						"interface": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"engine_id": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 255),
									},
									"engine_type": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 255),
									},
									"source_address": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"family_mpls_input": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"family_inet_input", "family_inet6_input", "family_mpls_input", "input"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_packets_per_second": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"maximum_packet_length": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 9192),
						},
						"rate": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 16000000),
						},
						"run_length": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 20),
						},
					},
				},
			},
			"family_mpls_output": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"family_inet_output", "family_inet6_output", "family_mpls_output"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aggregate_export_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(90, 1800),
						},
						"flow_active_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(60, 1800),
						},
						"flow_inactive_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(15, 1800),
						},
						"flow_server": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"aggregation_autonomous_system": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_destination_prefix": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_protocol_port": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_source_destination_prefix": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_source_destination_prefix_caida_compliant": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"aggregation_source_prefix": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"autonomous_system_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"origin", "peer"}, false),
									},
									"dscp": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 63),
									},
									"forwarding_class": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"local_dump": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"no_local_dump": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"routing_instance": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"source_address": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"version_ipfix_template": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"version9_template": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"inline_jflow_export_rate": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 3200),
							RequiredWith: []string{"family_mpls_output.0.inline_jflow_source_address"},
						},
						"inline_jflow_source_address": {
							Type:         schema.TypeString,
							Optional:     true,
							RequiredWith: []string{"family_mpls_output.0.flow_server"},
						},
						"interface": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"engine_id": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 255),
									},
									"engine_type": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 255),
									},
									"source_address": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"input": {
				Type:          schema.TypeList,
				Optional:      true,
				AtLeastOneOf:  []string{"family_inet_input", "family_inet6_input", "family_mpls_input", "input"},
				ConflictsWith: []string{"family_inet_input", "family_inet6_input", "family_mpls_input"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_packets_per_second": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"maximum_packet_length": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 9192),
						},
						"rate": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 16000000),
						},
						"run_length": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 20),
						},
					},
				},
			},
		},
	}
}

func resourceForwardingOptionsSamplingInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setForwardingOptionsSamplingInstance(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

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
	fwdoptsSamplingInstanceExists, err := checkForwardingOptionsSamplingInstanceExists(
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if fwdoptsSamplingInstanceExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("forwarding-options sampling instance %v already exists", d.Get("name").(string)))...)
	}

	if err := setForwardingOptionsSamplingInstance(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	warns, err := junSess.CommitConf("create resource junos_forwardingoptions_sampling_instance")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	fwdoptsSamplingInstanceExists, err = checkForwardingOptionsSamplingInstanceExists(
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if fwdoptsSamplingInstanceExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("forwarding-options sampling instance %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceForwardingOptionsSamplingInstanceReadWJunSess(d, junSess)...)
}

func resourceForwardingOptionsSamplingInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceForwardingOptionsSamplingInstanceReadWJunSess(d, junSess)
}

func resourceForwardingOptionsSamplingInstanceReadWJunSess(
	d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	samplingInstanceOptions, err := readForwardingOptionsSamplingInstance(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if samplingInstanceOptions.name == "" {
		d.SetId("")
	} else {
		fillForwardingOptionsSamplingInstanceData(d, samplingInstanceOptions)
	}

	return nil
}

func resourceForwardingOptionsSamplingInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delForwardingOptionsSamplingInstance(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setForwardingOptionsSamplingInstance(d, junSess); err != nil {
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
	if err := delForwardingOptionsSamplingInstance(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setForwardingOptionsSamplingInstance(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	warns, err := junSess.CommitConf("update resource junos_forwardingoptions_sampling_instance")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceForwardingOptionsSamplingInstanceReadWJunSess(d, junSess)...)
}

func resourceForwardingOptionsSamplingInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delForwardingOptionsSamplingInstance(d.Get("name").(string), junSess); err != nil {
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
	if err := delForwardingOptionsSamplingInstance(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_forwardingoptions_sampling_instance")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceForwardingOptionsSamplingInstanceImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	fwdoptsSamplingInstanceExists, err := checkForwardingOptionsSamplingInstanceExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !fwdoptsSamplingInstanceExists {
		return nil, fmt.Errorf("don't find forwarding-options sampling instance with id '%v' (id must be <name>)", d.Id())
	}
	samplingInstanceOptions, err := readForwardingOptionsSamplingInstance(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillForwardingOptionsSamplingInstanceData(d, samplingInstanceOptions)

	result[0] = d

	return result, nil
}

func checkForwardingOptionsSamplingInstanceExists(name string, junSess *junos.Session,
) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"forwarding-options sampling instance \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setForwardingOptionsSamplingInstance(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set forwarding-options sampling instance \"" + d.Get("name").(string) + "\" "
	if d.Get("disable").(bool) {
		configSet = append(configSet, setPrefix+"disable")
	}
	for _, v := range d.Get("family_inet_input").([]interface{}) {
		if err := setForwardingOptionsSamplingInstanceInput(setPrefix,
			v.(map[string]interface{}), junos.InetW, junSess); err != nil {
			return err
		}
	}
	for _, v := range d.Get("family_inet_output").([]interface{}) {
		if v == nil {
			return fmt.Errorf("family_inet_output block is empty")
		}
		if err := setForwardingOptionsSamplingInstanceOutput(setPrefix,
			v.(map[string]interface{}), junos.InetW, junSess); err != nil {
			return err
		}
	}
	for _, v := range d.Get("family_inet6_input").([]interface{}) {
		if err := setForwardingOptionsSamplingInstanceInput(setPrefix,
			v.(map[string]interface{}), junos.Inet6W, junSess); err != nil {
			return err
		}
	}
	for _, v := range d.Get("family_inet6_output").([]interface{}) {
		if v == nil {
			return fmt.Errorf("family_inet6_output block is empty")
		}
		if err := setForwardingOptionsSamplingInstanceOutput(setPrefix,
			v.(map[string]interface{}), junos.Inet6W, junSess); err != nil {
			return err
		}
	}
	for _, v := range d.Get("family_mpls_input").([]interface{}) {
		if err := setForwardingOptionsSamplingInstanceInput(setPrefix,
			v.(map[string]interface{}), junos.MplsW, junSess); err != nil {
			return err
		}
	}
	for _, v := range d.Get("family_mpls_output").([]interface{}) {
		if v == nil {
			return fmt.Errorf("family_mpls_output block is empty")
		}
		if err := setForwardingOptionsSamplingInstanceOutput(setPrefix,
			v.(map[string]interface{}), junos.MplsW, junSess); err != nil {
			return err
		}
	}
	for _, v := range d.Get("input").([]interface{}) {
		if err := setForwardingOptionsSamplingInstanceInput(setPrefix,
			v.(map[string]interface{}), "", junSess); err != nil {
			return err
		}
	}

	return junSess.ConfigSet(configSet)
}

func setForwardingOptionsSamplingInstanceInput(
	setPrefix string, input map[string]interface{}, family string, junSess *junos.Session,
) error {
	configSet := make([]string, 0)
	switch family {
	case junos.InetW:
		setPrefix += "family inet input "
	case junos.Inet6W:
		setPrefix += "family inet6 input "
	case junos.MplsW:
		setPrefix += "family mpls input "
	default:
		setPrefix += "input "
	}
	if v := input["max_packets_per_second"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"max-packets-per-second "+strconv.Itoa(v))
	}
	if v := input["maximum_packet_length"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"maximum-packet-length "+strconv.Itoa(v))
	}
	if v := input["rate"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"rate "+strconv.Itoa(v))
	}
	if v := input["run_length"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"run-length "+strconv.Itoa(v))
	}
	if len(configSet) == 0 {
		switch family {
		case junos.InetW:
			return fmt.Errorf("family_inet_input block is empty")
		case junos.Inet6W:
			return fmt.Errorf("family_inet6_input block is empty")
		case junos.MplsW:
			return fmt.Errorf("family_mpls_input block is empty")
		default:
			return fmt.Errorf("input block is empty")
		}
	}

	return junSess.ConfigSet(configSet)
}

func setForwardingOptionsSamplingInstanceOutput(
	setPrefix string, output map[string]interface{}, family string, junSess *junos.Session,
) error {
	configSet := make([]string, 0)
	switch family {
	case junos.InetW:
		setPrefix += "family inet output "
	case junos.Inet6W:
		setPrefix += "family inet6 output "
	case junos.MplsW:
		setPrefix += "family mpls output "
	default:
		return fmt.Errorf("internal error: setForwardingOptionsSamplingInstanceOutput call with bad family")
	}
	if v := output["aggregate_export_interval"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"aggregate-export-interval "+strconv.Itoa(v))
	}
	if family == junos.InetW || family == junos.Inet6W {
		for _, v := range output["extension_service"].([]interface{}) {
			configSet = append(configSet, setPrefix+"extension-service \""+v.(string)+"\"")
		}
	}
	if v := output["flow_active_timeout"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"flow-active-timeout "+strconv.Itoa(v))
	}
	if v := output["flow_inactive_timeout"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"flow-inactive-timeout "+strconv.Itoa(v))
	}
	flowServerHostnameList := make([]string, 0)
	for _, vFS := range output["flow_server"].([]interface{}) {
		flowServer := vFS.(map[string]interface{})
		if bchk.InSlice(flowServer["hostname"].(string), flowServerHostnameList) {
			return fmt.Errorf("multiple blocks flow_server with the same hostname %s", flowServer["hostname"].(string))
		}
		flowServerHostnameList = append(flowServerHostnameList, flowServer["hostname"].(string))
		setPrefixFlowServer := setPrefix + "flow-server " + flowServer["hostname"].(string) + " "
		configSet = append(configSet, setPrefixFlowServer+"port "+strconv.Itoa(flowServer["port"].(int)))
		if flowServer["aggregation_autonomous_system"].(bool) {
			configSet = append(configSet, setPrefixFlowServer+"aggregation autonomous-system")
		}
		if flowServer["aggregation_destination_prefix"].(bool) {
			configSet = append(configSet, setPrefixFlowServer+"aggregation destination-prefix")
		}
		if flowServer["aggregation_protocol_port"].(bool) {
			configSet = append(configSet, setPrefixFlowServer+"aggregation protocol-port")
		}
		if flowServer["aggregation_source_destination_prefix"].(bool) {
			configSet = append(configSet, setPrefixFlowServer+"aggregation source-destination-prefix")
			if flowServer["aggregation_source_destination_prefix_caida_compliant"].(bool) {
				configSet = append(configSet, setPrefixFlowServer+"aggregation source-destination-prefix caida-compliant")
			}
		} else if flowServer["aggregation_source_destination_prefix_caida_compliant"].(bool) {
			return fmt.Errorf("aggregation_source_destination_prefix_caida_compliant = true " +
				"without aggregation_source_destination_prefix on flow-server " + flowServer["hostname"].(string))
		}
		if flowServer["aggregation_source_prefix"].(bool) {
			configSet = append(configSet, setPrefixFlowServer+"aggregation source-prefix")
		}
		if v := flowServer["autonomous_system_type"].(string); v != "" {
			configSet = append(configSet, setPrefixFlowServer+"autonomous-system-type "+v)
		}
		if v := flowServer["dscp"].(int); v != -1 {
			configSet = append(configSet, setPrefixFlowServer+"dscp "+strconv.Itoa(v))
		}
		if v := flowServer["forwarding_class"].(string); v != "" {
			configSet = append(configSet, setPrefixFlowServer+"forwarding-class "+v)
		}
		if flowServer["local_dump"].(bool) {
			configSet = append(configSet, setPrefixFlowServer+"local-dump")
		}
		if flowServer["no_local_dump"].(bool) {
			configSet = append(configSet, setPrefixFlowServer+"no-local-dump")
		}
		if v := flowServer["routing_instance"].(string); v != "" {
			configSet = append(configSet, setPrefixFlowServer+"routing-instance "+v)
		}
		if v := flowServer["source_address"].(string); v != "" {
			configSet = append(configSet, setPrefixFlowServer+"source-address "+v)
		}
		if family == junos.InetW {
			if v := flowServer["version"].(int); v != 0 {
				configSet = append(configSet, setPrefixFlowServer+"version "+strconv.Itoa(v))
			}
		}
		if v := flowServer["version_ipfix_template"].(string); v != "" {
			configSet = append(configSet, setPrefixFlowServer+"version-ipfix template \""+v+"\"")
		}
		if v := flowServer["version9_template"].(string); v != "" {
			configSet = append(configSet, setPrefixFlowServer+"version9 template \""+v+"\"")
		}
	}
	if v := output["inline_jflow_export_rate"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"inline-jflow flow-export-rate "+strconv.Itoa(v))
	}
	if v := output["inline_jflow_source_address"].(string); v != "" {
		configSet = append(configSet, setPrefix+"inline-jflow source-address "+v)
	}
	interfaceNameList := make([]string, 0)
	for _, vIF := range output["interface"].([]interface{}) {
		interFace := vIF.(map[string]interface{})
		if bchk.InSlice(interFace["name"].(string), interfaceNameList) {
			return fmt.Errorf("multiple blocks interface with the same name %s", interFace["name"].(string))
		}
		interfaceNameList = append(interfaceNameList, interFace["name"].(string))
		setPrefixInterface := setPrefix + "interface " + interFace["name"].(string) + " "
		configSet = append(configSet, setPrefixInterface)
		if v := interFace["engine_id"].(int); v != -1 {
			configSet = append(configSet, setPrefixInterface+"engine-id "+strconv.Itoa(v))
		}
		if v := interFace["engine_type"].(int); v != -1 {
			configSet = append(configSet, setPrefixInterface+"engine-type "+strconv.Itoa(v))
		}
		if v := interFace["source_address"].(string); v != "" {
			configSet = append(configSet, setPrefixInterface+"source-address "+v)
		}
	}

	return junSess.ConfigSet(configSet)
}

func readForwardingOptionsSamplingInstance(name string, junSess *junos.Session,
) (confRead samplingInstanceOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"forwarding-options sampling instance \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			switch {
			case itemTrim == junos.DisableW:
				confRead.disable = true
			case balt.CutPrefixInString(&itemTrim, "family inet input "):
				if len(confRead.familyInetInput) == 0 {
					confRead.familyInetInput = append(confRead.familyInetInput, map[string]interface{}{
						"max_packets_per_second": -1,
						"maximum_packet_length":  -1,
						"rate":                   0,
						"run_length":             -1,
					})
				}
				if err := readForwardingOptionsSamplingInstanceInput(itemTrim, confRead.familyInetInput[0]); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "family inet6 input "):
				if len(confRead.familyInet6Input) == 0 {
					confRead.familyInet6Input = append(confRead.familyInet6Input, map[string]interface{}{
						"max_packets_per_second": -1,
						"maximum_packet_length":  -1,
						"rate":                   0,
						"run_length":             -1,
					})
				}
				if err := readForwardingOptionsSamplingInstanceInput(itemTrim, confRead.familyInet6Input[0]); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "family mpls input "):
				if len(confRead.familyMplsInput) == 0 {
					confRead.familyMplsInput = append(confRead.familyMplsInput, map[string]interface{}{
						"max_packets_per_second": -1,
						"maximum_packet_length":  -1,
						"rate":                   0,
						"run_length":             -1,
					})
				}
				if err := readForwardingOptionsSamplingInstanceInput(itemTrim, confRead.familyMplsInput[0]); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "input "):
				if len(confRead.input) == 0 {
					confRead.input = append(confRead.input, map[string]interface{}{
						"max_packets_per_second": -1,
						"maximum_packet_length":  -1,
						"rate":                   0,
						"run_length":             -1,
					})
				}
				if err := readForwardingOptionsSamplingInstanceInput(itemTrim, confRead.input[0]); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "family inet output "):
				if len(confRead.familyInetOutput) == 0 {
					confRead.familyInetOutput = append(confRead.familyInetOutput, map[string]interface{}{
						"aggregate_export_interval":   0,
						"extension_service":           make([]string, 0),
						"flow_active_timeout":         0,
						"flow_inactive_timeout":       0,
						"flow_server":                 make([]map[string]interface{}, 0),
						"inline_jflow_export_rate":    0,
						"inline_jflow_source_address": "",
						"interface":                   make([]map[string]interface{}, 0),
					})
				}
				if err := readForwardingOptionsSamplingInstanceOutput(
					itemTrim,
					confRead.familyInetOutput[0],
					junos.InetW,
				); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "family inet6 output "):
				if len(confRead.familyInet6Output) == 0 {
					confRead.familyInet6Output = append(confRead.familyInet6Output, map[string]interface{}{
						"aggregate_export_interval":   0,
						"extension_service":           make([]string, 0),
						"flow_active_timeout":         0,
						"flow_inactive_timeout":       0,
						"flow_server":                 make([]map[string]interface{}, 0),
						"inline_jflow_export_rate":    0,
						"inline_jflow_source_address": "",
						"interface":                   make([]map[string]interface{}, 0),
					})
				}
				if err := readForwardingOptionsSamplingInstanceOutput(
					itemTrim,
					confRead.familyInet6Output[0],
					junos.Inet6W,
				); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "family mpls output "):
				if len(confRead.familyMplsOutput) == 0 {
					confRead.familyMplsOutput = append(confRead.familyMplsOutput, map[string]interface{}{
						"aggregate_export_interval":   0,
						"flow_active_timeout":         0,
						"flow_inactive_timeout":       0,
						"flow_server":                 make([]map[string]interface{}, 0),
						"inline_jflow_export_rate":    0,
						"inline_jflow_source_address": "",
						"interface":                   make([]map[string]interface{}, 0),
					})
				}
				if err := readForwardingOptionsSamplingInstanceOutput(
					itemTrim,
					confRead.familyMplsOutput[0],
					junos.MplsW,
				); err != nil {
					return confRead, err
				}
			}
		}
	}

	return confRead, nil
}

func readForwardingOptionsSamplingInstanceInput(itemTrim string, inputRead map[string]interface{}) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "max-packets-per-second "):
		inputRead["max_packets_per_second"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "maximum-packet-length "):
		inputRead["maximum_packet_length"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "rate "):
		inputRead["rate"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "run-length "):
		inputRead["run_length"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	}

	return nil
}

func readForwardingOptionsSamplingInstanceOutput(itemTrim string, outputRead map[string]interface{}, family string,
) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "aggregate-export-interval "):
		outputRead["aggregate_export_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "extension-service "):
		outputRead["extension_service"] = append(outputRead["extension_service"].([]string), strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "flow-active-timeout "):
		outputRead["flow_active_timeout"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "flow-inactive-timeout "):
		outputRead["flow_inactive_timeout"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "flow-server "):
		itemTrimFields := strings.Split(itemTrim, " ")
		flowServer := map[string]interface{}{
			"hostname":                              itemTrimFields[0],
			"port":                                  0,
			"aggregation_autonomous_system":         false,
			"aggregation_destination_prefix":        false,
			"aggregation_protocol_port":             false,
			"aggregation_source_destination_prefix": false,
			"aggregation_source_destination_prefix_caida_compliant": false,
			"aggregation_source_prefix":                             false,
			"autonomous_system_type":                                "",
			"dscp":                                                  -1,
			"forwarding_class":                                      "",
			"local_dump":                                            false,
			"no_local_dump":                                         false,
			"routing_instance":                                      "",
			"source_address":                                        "",
			"version_ipfix_template":                                "",
			"version9_template":                                     "",
		}
		if family == junos.InetW {
			flowServer["version"] = 0
		}
		outputRead["flow_server"] = copyAndRemoveItemMapList(
			"hostname", flowServer, outputRead["flow_server"].([]map[string]interface{}))
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
		switch {
		case balt.CutPrefixInString(&itemTrim, "port "):
			flowServer["port"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case itemTrim == "aggregation autonomous-system":
			flowServer["aggregation_autonomous_system"] = true
		case itemTrim == "aggregation destination-prefix":
			flowServer["aggregation_destination_prefix"] = true
		case itemTrim == "aggregation protocol-port":
			flowServer["aggregation_protocol_port"] = true
		case balt.CutPrefixInString(&itemTrim, "aggregation source-destination-prefix"):
			flowServer["aggregation_source_destination_prefix"] = true
			if itemTrim == " caida-compliant" {
				flowServer["aggregation_source_destination_prefix_caida_compliant"] = true
			}
		case itemTrim == "aggregation source-prefix":
			flowServer["aggregation_source_prefix"] = true
		case balt.CutPrefixInString(&itemTrim, "autonomous-system-type "):
			flowServer["autonomous_system_type"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "dscp "):
			flowServer["dscp"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "forwarding-class "):
			flowServer["forwarding_class"] = itemTrim
		case itemTrim == "local-dump":
			flowServer["local_dump"] = true
		case itemTrim == "no-local-dump":
			flowServer["no_local_dump"] = true
		case balt.CutPrefixInString(&itemTrim, "routing-instance "):
			flowServer["routing_instance"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "source-address "):
			flowServer["source_address"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "version "):
			flowServer["version"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "version-ipfix template "):
			flowServer["version_ipfix_template"] = strings.Trim(itemTrim, "\"")
		case balt.CutPrefixInString(&itemTrim, "version9 template "):
			flowServer["version9_template"] = strings.Trim(itemTrim, "\"")
		}
		outputRead["flow_server"] = append(outputRead["flow_server"].([]map[string]interface{}), flowServer)
	case balt.CutPrefixInString(&itemTrim, "inline-jflow flow-export-rate "):
		outputRead["inline_jflow_export_rate"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "inline-jflow source-address "):
		outputRead["inline_jflow_source_address"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "interface "):
		itemTrimFields := strings.Split(itemTrim, " ")
		iFace := map[string]interface{}{
			"name":           itemTrimFields[0],
			"engine_id":      -1,
			"engine_type":    -1,
			"source_address": "",
		}
		outputRead["interface"] = copyAndRemoveItemMapList(
			"name", iFace, outputRead["interface"].([]map[string]interface{}))
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
		switch {
		case balt.CutPrefixInString(&itemTrim, "engine-id "):
			iFace["engine_id"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "engine-type "):
			iFace["engine_type"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "source-address "):
			iFace["source_address"] = itemTrim
		}
		outputRead["interface"] = append(outputRead["interface"].([]map[string]interface{}), iFace)
	}

	return nil
}

func delForwardingOptionsSamplingInstance(samplingInstance string, junSess *junos.Session) error {
	configSet := []string{"delete forwarding-options sampling instance \"" + samplingInstance + "\""}

	return junSess.ConfigSet(configSet)
}

func fillForwardingOptionsSamplingInstanceData(
	d *schema.ResourceData, samplingInstanceOptions samplingInstanceOptions,
) {
	if tfErr := d.Set("name", samplingInstanceOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("disable", samplingInstanceOptions.disable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_inet_input", samplingInstanceOptions.familyInetInput); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_inet_output", samplingInstanceOptions.familyInetOutput); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_inet6_input", samplingInstanceOptions.familyInet6Input); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_inet6_output", samplingInstanceOptions.familyInet6Output); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_mpls_input", samplingInstanceOptions.familyMplsInput); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_mpls_output", samplingInstanceOptions.familyMplsOutput); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("input", samplingInstanceOptions.input); tfErr != nil {
		panic(tfErr)
	}
}
