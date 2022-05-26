package junos

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

func dataSourceInterfaceLogical() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceInterfaceLogicalRead,
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
			"disable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"family_inet": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"preferred": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"primary": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"vrrp_group": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"identifier": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"virtual_address": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"accept_data": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"advertise_interval": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"advertisements_threshold": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"authentication_key": {
													Type:      schema.TypeString,
													Computed:  true,
													Sensitive: true,
												},
												"authentication_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"no_accept_data": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"no_preempt": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"preempt": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"priority": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"track_interface": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"interface": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"priority_cost": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"track_route": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"route": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"routing_instance": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"priority_cost": {
																Type:     schema.TypeInt,
																Computed: true,
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
						"dhcp": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"srx_old_option_name": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"client_identifier_ascii": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_identifier_hexadecimal": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_identifier_prefix_hostname": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"client_identifier_prefix_routing_instance_name": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"client_identifier_use_interface_description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_identifier_userid_ascii": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_identifier_userid_hexadecimal": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"force_discover": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"lease_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"lease_time_infinite": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"metric": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"no_dns_install": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"options_no_hostname": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"retransmission_attempt": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"retransmission_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"server_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"update_server": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"vendor_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"filter_input": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_output": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mtu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rpf_check": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fail_filter": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mode_loose": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"sampling_input": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"sampling_output": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"family_inet6": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"preferred": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"primary": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"vrrp_group": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"identifier": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"virtual_address": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"virtual_link_local_address": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"accept_data": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"advertise_interval": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"no_accept_data": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"no_preempt": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"preempt": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"priority": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"track_interface": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"interface": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"priority_cost": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"track_route": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"route": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"routing_instance": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"priority_cost": {
																Type:     schema.TypeInt,
																Computed: true,
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
						"dad_disable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"dhcpv6_client": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client_identifier_duid_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_ia_type_na": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"client_ia_type_pd": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"no_dns_install": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"prefix_delegating_preferred_prefix_length": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"prefix_delegating_sub_prefix_length": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rapid_commit": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"req_option": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"retransmission_attempt": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"update_router_advertisement_interface": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"update_server": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"filter_input": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_output": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mtu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rpf_check": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fail_filter": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mode_loose": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"sampling_input": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"sampling_output": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"routing_instance": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_inbound_protocols": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_inbound_services": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allow_fragmentation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"do_not_fragment": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"flow_label": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"no_path_mtu_discovery": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"path_mtu_discovery": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"routing_instance_destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_class": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceInterfaceLogicalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("config_interface").(string) == "" && d.Get("match").(string) == "" {
		return diag.FromErr(fmt.Errorf("no arguments provided, 'config_interface' and 'match' empty"))
	}
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	mutex.Lock()
	nameFound, err := searchInterfaceLogicalID(d.Get("config_interface").(string), d.Get("match").(string), sess, junSess)
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	if nameFound == "" {
		mutex.Unlock()

		return diag.FromErr(fmt.Errorf("no logical interface found with arguments provided"))
	}
	interfaceOpt, err := readInterfaceLogical(nameFound, sess, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(nameFound)
	if tfErr := d.Set("name", nameFound); tfErr != nil {
		panic(tfErr)
	}
	fillInterfaceLogicalData(d, interfaceOpt)

	return nil
}

func searchInterfaceLogicalID(configInterface, match string, sess *Session, junSess *junosSession,
) (string, error) {
	intConfigList := make([]string, 0)
	showConfig, err := sess.command(cmdShowConfig+"interfaces "+configInterface+pipeDisplaySet, junSess)
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
			return "", fmt.Errorf("failed to regexp with '%s': %w", match, err)
		}
		if !matched {
			continue
		}
		itemTrimSplit := strings.Split(itemTrim, " ")
		switch len(itemTrimSplit) {
		case 0, 1, 2:
			continue
		default:
			if itemTrimSplit[1] == "unit" && !bchk.StringInSlice("ethernet-switching", itemTrimSplit) {
				intConfigList = append(intConfigList, itemTrimSplit[0]+"."+itemTrimSplit[2])
			}
		}
	}
	intConfigList = balt.UniqueStrings(intConfigList)
	if len(intConfigList) == 0 {
		return "", nil
	}
	if len(intConfigList) > 1 {
		return "", fmt.Errorf("too many different logical interfaces found")
	}

	return intConfigList[0], nil
}
