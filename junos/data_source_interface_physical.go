package junos

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

func dataSourceInterfacePhysical() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInterfacePhysicalRead,
		Schema: map[string]*schema.Schema{
			"config_interface": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"match": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if _, err := regexp.Compile(value); err != nil {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not valid regexp", value, k))
					}

					return
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ae_lacp": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "use parent_ether_opts { lacp } instead",
			},
			"ae_link_speed": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "use parent_ether_opts { link_speed } instead",
			},
			"ae_minimum_links": {
				Type:       schema.TypeInt,
				Computed:   true,
				Deprecated: "use parent_ether_opts { minimum_links } instead",
			},
			"esi": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_derive_lacp": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"df_election_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_bmac": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ether_opts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ae_8023ad": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_negotiation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"no_auto_negotiation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"flow_control": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"no_flow_control": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"loopback": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"no_loopback": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"redundant_parent": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ether802_3ad": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "use ether_opts { ae_8023ad } or gigether_opts { ae_8023ad } instead",
			},
			"gigether_opts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ae_8023ad": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_negotiation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"no_auto_negotiation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"flow_control": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"no_flow_control": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"loopback": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"no_loopback": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"redundant_parent": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"parent_ether_opts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bfd_liveness_detection": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"local_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"authentication_algorithm": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"authentication_key_chain": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"authentication_loose_check": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"detection_time_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"holddown_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"minimum_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"minimum_receive_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"multiplier": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"neighbor": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"no_adaptation": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"transmit_interval_minimum_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"transmit_interval_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"flow_control": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"no_flow_control": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"lacp": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"admin_key": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"periodic": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sync_reset": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"system_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"system_priority": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"loopback": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"no_loopback": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"link_speed": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"minimum_bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"minimum_links": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"redundancy_group": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source_address_filter": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"source_filtering": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"trunk": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"vlan_members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vlan_native": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vlan_tagging": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceInterfacePhysicalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("config_interface").(string) == "" && d.Get("match").(string) == "" {
		return diag.FromErr(fmt.Errorf("no arguments provided, 'config_interface' and 'match' empty"))
	}
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	mutex.Lock()
	nameFound, err := searchInterfacePhysicalID(d.Get("config_interface").(string), d.Get("match").(string), m, jnprSess)
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	if nameFound == "" {
		mutex.Unlock()

		return diag.FromErr(fmt.Errorf("no physical interface found with arguments provided"))
	}
	interfaceOpt, err := readInterfacePhysical(nameFound, m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(nameFound)
	if tfErr := d.Set("name", nameFound); tfErr != nil {
		panic(tfErr)
	}
	fillInterfacePhysicalData(d, interfaceOpt)

	return nil
}

func searchInterfacePhysicalID(configInterface string, match string,
	m interface{}, jnprSess *NetconfObject) (string, error) {
	sess := m.(*Session)
	intConfigList := make([]string, 0)
	showConfig, err := sess.command(cmdShowConfig+"interfaces "+configInterface+pipeDisplaySet, jnprSess)
	if err != nil {
		return "", err
	}
	for _, item := range strings.Split(showConfig, "\n") {
		if strings.Contains(item, xmlStartTagConfigOut) {
			continue
		}
		if strings.Contains(item, xmlEndTagConfigOut) {
			break
		}
		if item == "" {
			continue
		}
		itemTrim := strings.TrimPrefix(item, "set interfaces ")
		matched, err := regexp.MatchString(match, itemTrim)
		if err != nil {
			return "", fmt.Errorf("failed to regexp with %s : %w", match, err)
		}
		if !matched {
			continue
		}
		itemTrimSplit := strings.Split(itemTrim, " ")
		if len(itemTrimSplit) == 0 {
			continue
		}
		intConfigList = append(intConfigList, itemTrimSplit[0])
	}
	intConfigList = balt.UniqueStrings(intConfigList)
	if len(intConfigList) == 0 {
		return "", nil
	}
	if len(intConfigList) > 1 {
		return "", fmt.Errorf("too many different physical interfaces found")
	}

	return intConfigList[0], nil
}
