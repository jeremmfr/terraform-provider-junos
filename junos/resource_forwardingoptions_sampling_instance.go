package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
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

func resourceForwardingoptionsSamplingInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceForwardingoptionsSamplingInstanceCreate,
		ReadContext:   resourceForwardingoptionsSamplingInstanceRead,
		UpdateContext: resourceForwardingoptionsSamplingInstanceUpdate,
		DeleteContext: resourceForwardingoptionsSamplingInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceForwardingoptionsSamplingInstanceImport,
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

func resourceForwardingoptionsSamplingInstanceCreate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setForwardingoptionsSamplingInstance(d, m, nil); err != nil {
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
	fwdoptsSamplingInstanceExists, err := checkForwardingoptionsSamplingInstanceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if fwdoptsSamplingInstanceExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("forwarding-options sampling instance %v already exists", d.Get("name").(string)))...)
	}

	if err := setForwardingoptionsSamplingInstance(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	warns, err := sess.commitConf("create resource junos_forwardingoptions_sampling_instance", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	fwdoptsSamplingInstanceExists, err = checkForwardingoptionsSamplingInstanceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if fwdoptsSamplingInstanceExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("forwarding-options sampling instance %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceForwardingoptionsSamplingInstanceReadWJnprSess(d, m, jnprSess)...)
}

func resourceForwardingoptionsSamplingInstanceRead(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceForwardingoptionsSamplingInstanceReadWJnprSess(d, m, jnprSess)
}

func resourceForwardingoptionsSamplingInstanceReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	samplingInstanceOptions, err := readForwardingoptionsSamplingInstance(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if samplingInstanceOptions.name == "" {
		d.SetId("")
	} else {
		fillForwardingoptionsSamplingInstanceData(d, samplingInstanceOptions)
	}

	return nil
}

func resourceForwardingoptionsSamplingInstanceUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delForwardingoptionsSamplingInstance(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setForwardingoptionsSamplingInstance(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	warns, err := sess.commitConf("update resource junos_forwardingoptions_sampling_instance", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceForwardingoptionsSamplingInstanceReadWJnprSess(d, m, jnprSess)...)
}

func resourceForwardingoptionsSamplingInstanceDelete(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delForwardingoptionsSamplingInstance(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_forwardingoptions_sampling_instance", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceForwardingoptionsSamplingInstanceImport(d *schema.ResourceData,
	m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	fwdoptsSamplingInstanceExists, err := checkForwardingoptionsSamplingInstanceExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !fwdoptsSamplingInstanceExists {
		return nil, fmt.Errorf("don't find forwarding-options sampling instance with id '%v' (id must be <name>)", d.Id())
	}
	samplingInstanceOptions, err := readForwardingoptionsSamplingInstance(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillForwardingoptionsSamplingInstanceData(d, samplingInstanceOptions)

	result[0] = d

	return result, nil
}

func checkForwardingoptionsSamplingInstanceExists(
	name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	samplingInstanceConfig, err := sess.command(
		"show configuration forwarding-options sampling instance \""+name+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if samplingInstanceConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setForwardingoptionsSamplingInstance(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set forwarding-options sampling instance \"" + d.Get("name").(string) + "\" "
	if d.Get("disable").(bool) {
		configSet = append(configSet, setPrefix+"disable")
	}
	for _, v := range d.Get("family_inet_input").([]interface{}) {
		if err := setForwardingoptionsSamplingInstanceInput(setPrefix,
			v.(map[string]interface{}), inetWord, sess, jnprSess); err != nil {
			return err
		}
	}
	for _, v := range d.Get("family_inet_output").([]interface{}) {
		if v == nil {
			return fmt.Errorf("family_inet_output block is empty")
		}
		if err := setForwardingoptionsSamplingInstanceOutput(setPrefix,
			v.(map[string]interface{}), inetWord, sess, jnprSess); err != nil {
			return err
		}
	}
	for _, v := range d.Get("family_inet6_input").([]interface{}) {
		if err := setForwardingoptionsSamplingInstanceInput(setPrefix,
			v.(map[string]interface{}), inet6Word, sess, jnprSess); err != nil {
			return err
		}
	}
	for _, v := range d.Get("family_inet6_output").([]interface{}) {
		if v == nil {
			return fmt.Errorf("family_inet6_output block is empty")
		}
		if err := setForwardingoptionsSamplingInstanceOutput(setPrefix,
			v.(map[string]interface{}), inet6Word, sess, jnprSess); err != nil {
			return err
		}
	}
	for _, v := range d.Get("family_mpls_input").([]interface{}) {
		if err := setForwardingoptionsSamplingInstanceInput(setPrefix,
			v.(map[string]interface{}), mplsWord, sess, jnprSess); err != nil {
			return err
		}
	}
	for _, v := range d.Get("family_mpls_output").([]interface{}) {
		if v == nil {
			return fmt.Errorf("family_mpls_output block is empty")
		}
		if err := setForwardingoptionsSamplingInstanceOutput(setPrefix,
			v.(map[string]interface{}), mplsWord, sess, jnprSess); err != nil {
			return err
		}
	}
	for _, v := range d.Get("input").([]interface{}) {
		if err := setForwardingoptionsSamplingInstanceInput(setPrefix,
			v.(map[string]interface{}), "", sess, jnprSess); err != nil {
			return err
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func setForwardingoptionsSamplingInstanceInput(
	setPrefix string, input map[string]interface{}, family string, sess *Session, jnprSess *NetconfObject) error {
	configSet := make([]string, 0)
	switch family {
	case inetWord:
		setPrefix += "family inet input "
	case inet6Word:
		setPrefix += "family inet6 input "
	case mplsWord:
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
		case inetWord:
			return fmt.Errorf("family_inet_input block is empty")
		case inet6Word:
			return fmt.Errorf("family_inet6_input block is empty")
		case mplsWord:
			return fmt.Errorf("family_mpls_input block is empty")
		default:
			return fmt.Errorf("input block is empty")
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func setForwardingoptionsSamplingInstanceOutput(
	setPrefix string, output map[string]interface{}, family string, sess *Session, jnprSess *NetconfObject) error {
	configSet := make([]string, 0)
	switch family {
	case inetWord:
		setPrefix += "family inet output "
	case inet6Word:
		setPrefix += "family inet6 output "
	case mplsWord:
		setPrefix += "family mpls output "
	default:
		return fmt.Errorf("internal error: setForwardingoptionsSamplingInstanceOutput call with bad family")
	}
	if v := output["aggregate_export_interval"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"aggregate-export-interval "+strconv.Itoa(v))
	}
	if family == inetWord || family == inet6Word {
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
		if bchk.StringInSlice(flowServer["hostname"].(string), flowServerHostnameList) {
			return fmt.Errorf("multiple flow_server blocks with the same hostname")
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
		if family == inetWord {
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
		if bchk.StringInSlice(interFace["name"].(string), interfaceNameList) {
			return fmt.Errorf("multiple interface blocks with the same name")
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

	return sess.configSet(configSet, jnprSess)
}

func readForwardingoptionsSamplingInstance(
	samplingInstance string, m interface{}, jnprSess *NetconfObject) (samplingInstanceOptions, error) {
	sess := m.(*Session)
	var confRead samplingInstanceOptions

	samplingInstanceConfig, err := sess.command("show configuration forwarding-options sampling instance \""+
		samplingInstance+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if samplingInstanceConfig != emptyWord {
		confRead.name = samplingInstance
		for _, item := range strings.Split(samplingInstanceConfig, "\n") {
			itemTrim := strings.TrimPrefix(item, setLineStart)
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			switch {
			case itemTrim == "disable":
				confRead.disable = true
			case strings.HasPrefix(itemTrim, "family inet input "):
				if len(confRead.familyInetInput) == 0 {
					confRead.familyInetInput = append(confRead.familyInetInput, map[string]interface{}{
						"max_packets_per_second": -1,
						"maximum_packet_length":  -1,
						"rate":                   0,
						"run_length":             -1,
					})
				}
				if err := readForwardingoptionsSamplingInstanceInput(confRead.familyInetInput[0],
					strings.TrimPrefix(itemTrim, "family inet input ")); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "family inet6 input "):
				if len(confRead.familyInet6Input) == 0 {
					confRead.familyInet6Input = append(confRead.familyInet6Input, map[string]interface{}{
						"max_packets_per_second": -1,
						"maximum_packet_length":  -1,
						"rate":                   0,
						"run_length":             -1,
					})
				}
				if err := readForwardingoptionsSamplingInstanceInput(confRead.familyInet6Input[0],
					strings.TrimPrefix(itemTrim, "family inet6 input ")); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "family mpls input "):
				if len(confRead.familyMplsInput) == 0 {
					confRead.familyMplsInput = append(confRead.familyMplsInput, map[string]interface{}{
						"max_packets_per_second": -1,
						"maximum_packet_length":  -1,
						"rate":                   0,
						"run_length":             -1,
					})
				}
				if err := readForwardingoptionsSamplingInstanceInput(confRead.familyMplsInput[0],
					strings.TrimPrefix(itemTrim, "family mpls input ")); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "input "):
				if len(confRead.input) == 0 {
					confRead.input = append(confRead.input, map[string]interface{}{
						"max_packets_per_second": -1,
						"maximum_packet_length":  -1,
						"rate":                   0,
						"run_length":             -1,
					})
				}
				if err := readForwardingoptionsSamplingInstanceInput(confRead.input[0],
					strings.TrimPrefix(itemTrim, "input ")); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "family inet output "):
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
				if err := readForwardingoptionsSamplingInstanceOutput(confRead.familyInetOutput[0],
					strings.TrimPrefix(itemTrim, "family inet output "), inetWord); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "family inet6 output "):
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
				if err := readForwardingoptionsSamplingInstanceOutput(confRead.familyInet6Output[0],
					strings.TrimPrefix(itemTrim, "family inet6 output "), inet6Word); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "family mpls output "):
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
				if err := readForwardingoptionsSamplingInstanceOutput(confRead.familyMplsOutput[0],
					strings.TrimPrefix(itemTrim, "family mpls output "), mplsWord); err != nil {
					return confRead, err
				}
			}
		}
	}

	return confRead, nil
}

func readForwardingoptionsSamplingInstanceInput(inputRead map[string]interface{}, itemTrim string) error {
	switch {
	case strings.HasPrefix(itemTrim, "max-packets-per-second "):
		var err error
		inputRead["max_packets_per_second"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-packets-per-second "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "maximum-packet-length "):
		var err error
		inputRead["maximum_packet_length"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "maximum-packet-length "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "rate "):
		var err error
		inputRead["rate"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "rate "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "run-length "):
		var err error
		inputRead["run_length"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "run-length "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}

	return nil
}

func readForwardingoptionsSamplingInstanceOutput(
	outputRead map[string]interface{}, itemTrim string, family string) error {
	switch {
	case strings.HasPrefix(itemTrim, "aggregate-export-interval "):
		var err error
		outputRead["aggregate_export_interval"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "aggregate-export-interval "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "extension-service "):
		outputRead["extension_service"] = append(outputRead["extension_service"].([]string),
			strings.Trim(strings.TrimPrefix(itemTrim, "extension-service "), "\""))
	case strings.HasPrefix(itemTrim, "flow-active-timeout "):
		var err error
		outputRead["flow_active_timeout"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "flow-active-timeout "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "flow-inactive-timeout "):
		var err error
		outputRead["flow_inactive_timeout"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "flow-inactive-timeout "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "flow-server "):
		flowServerLineCut := strings.Split(itemTrim, " ")
		flowServer := map[string]interface{}{
			"hostname":                              flowServerLineCut[1],
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
		if family == inetWord {
			flowServer["version"] = 0
		}
		outputRead["flow_server"] = copyAndRemoveItemMapList(
			"hostname", flowServer, outputRead["flow_server"].([]map[string]interface{}))
		itemTrimFlowServer := strings.TrimPrefix(itemTrim, "flow-server "+flowServerLineCut[1]+" ")
		switch {
		case strings.HasPrefix(itemTrimFlowServer, "port "):
			var err error
			flowServer["port"], err = strconv.Atoi(strings.TrimPrefix(itemTrimFlowServer, "port "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrimFlowServer, err)
			}
		case itemTrimFlowServer == "aggregation autonomous-system":
			flowServer["aggregation_autonomous_system"] = true
		case itemTrimFlowServer == "aggregation destination-prefix":
			flowServer["aggregation_destination_prefix"] = true
		case itemTrimFlowServer == "aggregation protocol-port":
			flowServer["aggregation_protocol_port"] = true
		case strings.HasPrefix(itemTrimFlowServer, "aggregation source-destination-prefix"):
			flowServer["aggregation_source_destination_prefix"] = true
			if itemTrimFlowServer == "aggregation source-destination-prefix caida-compliant" {
				flowServer["aggregation_source_destination_prefix_caida_compliant"] = true
			}
		case itemTrimFlowServer == "aggregation source-prefix":
			flowServer["aggregation_source_prefix"] = true
		case strings.HasPrefix(itemTrimFlowServer, "autonomous-system-type "):
			flowServer["autonomous_system_type"] = strings.TrimPrefix(itemTrimFlowServer, "autonomous-system-type ")
		case strings.HasPrefix(itemTrimFlowServer, "dscp "):
			var err error
			flowServer["dscp"], err = strconv.Atoi(strings.TrimPrefix(itemTrimFlowServer, "dscp "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrimFlowServer, err)
			}
		case strings.HasPrefix(itemTrimFlowServer, "forwarding-class "):
			flowServer["forwarding_class"] = strings.TrimPrefix(itemTrimFlowServer, "forwarding-class ")
		case itemTrimFlowServer == "local-dump":
			flowServer["local_dump"] = true
		case itemTrimFlowServer == "no-local-dump":
			flowServer["no_local_dump"] = true
		case strings.HasPrefix(itemTrimFlowServer, "routing-instance "):
			flowServer["routing_instance"] = strings.TrimPrefix(itemTrimFlowServer, "routing-instance ")
		case strings.HasPrefix(itemTrimFlowServer, "source-address "):
			flowServer["source_address"] = strings.TrimPrefix(itemTrimFlowServer, "source-address ")
		case strings.HasPrefix(itemTrimFlowServer, "version "):
			var err error
			flowServer["version"], err = strconv.Atoi(strings.TrimPrefix(itemTrimFlowServer, "version "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrimFlowServer, err)
			}
		case strings.HasPrefix(itemTrimFlowServer, "version-ipfix template "):
			flowServer["version_ipfix_template"] =
				strings.Trim(strings.TrimPrefix(itemTrimFlowServer, "version-ipfix template "), "\"")
		case strings.HasPrefix(itemTrimFlowServer, "version9 template "):
			flowServer["version9_template"] = strings.Trim(strings.TrimPrefix(itemTrimFlowServer, "version9 template "), "\"")
		}
		outputRead["flow_server"] = append(outputRead["flow_server"].([]map[string]interface{}), flowServer)
	case strings.HasPrefix(itemTrim, "inline-jflow flow-export-rate "):
		var err error
		outputRead["inline_jflow_export_rate"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "inline-jflow flow-export-rate "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "inline-jflow source-address "):
		outputRead["inline_jflow_source_address"] = strings.TrimPrefix(itemTrim, "inline-jflow source-address ")
	case strings.HasPrefix(itemTrim, "interface "):
		interfaceLineCut := strings.Split(itemTrim, " ")
		iFace := map[string]interface{}{
			"name":           interfaceLineCut[1],
			"engine_id":      -1,
			"engine_type":    -1,
			"source_address": "",
		}
		outputRead["interface"] = copyAndRemoveItemMapList(
			"name", iFace, outputRead["interface"].([]map[string]interface{}))
		itemTrimInterface := strings.TrimPrefix(itemTrim, "interface "+interfaceLineCut[1]+" ")
		switch {
		case strings.HasPrefix(itemTrimInterface, "engine-id "):
			var err error
			iFace["engine_id"], err = strconv.Atoi(strings.TrimPrefix(itemTrimInterface, "engine-id "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrimInterface, err)
			}
		case strings.HasPrefix(itemTrimInterface, "engine-type "):
			var err error
			iFace["engine_type"], err = strconv.Atoi(strings.TrimPrefix(itemTrimInterface, "engine-type "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrimInterface, err)
			}
		case strings.HasPrefix(itemTrimInterface, "source-address "):
			iFace["source_address"] = strings.TrimPrefix(itemTrimInterface, "source-address ")
		}
		outputRead["interface"] = append(outputRead["interface"].([]map[string]interface{}), iFace)
	}

	return nil
}

func delForwardingoptionsSamplingInstance(samplingInstance string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete forwarding-options sampling instance \"" + samplingInstance + "\""}

	return sess.configSet(configSet, jnprSess)
}

func fillForwardingoptionsSamplingInstanceData(
	d *schema.ResourceData, samplingInstanceOptions samplingInstanceOptions) {
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
