package junos

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInterfaceRead,
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
			"vlan_tagging": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"inet": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"inet6": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"inet_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
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
										Type:     schema.TypeString,
										Computed: true,
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
			"inet6_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
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
										MinItems: 1,
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
			"inet_mtu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"inet6_mtu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"inet_filter_input": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inet_filter_output": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inet6_filter_input": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inet6_filter_output": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ether802_3ad": {
				Type:     schema.TypeString,
				Computed: true,
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
			"ae_lacp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ae_link_speed": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ae_minimum_links": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"security_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routing_instance": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceInterfaceRead(d *schema.ResourceData, m interface{}) error {
	if d.Get("config_interface").(string) == "" && d.Get("match").(string) == "" {
		return fmt.Errorf("no arguments provided, 'config_interface' and 'match' empty")
	}
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	nameFound, err := searchInterfaceID(d.Get("config_interface").(string), d.Get("match").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if nameFound == "" {
		return fmt.Errorf("no interface found with arguments provided")
	}
	interfaceOpt, err := readInterface(nameFound, m, jnprSess)
	if err != nil {
		return err
	}
	d.SetId(nameFound)
	tfErr := d.Set("name", nameFound)
	if tfErr != nil {
		panic(tfErr)
	}
	fillInterfaceData(d, interfaceOpt)

	return nil
}

func searchInterfaceID(configInterface string, match string,
	m interface{}, jnprSess *NetconfObject) (string, error) {
	sess := m.(*Session)
	intConfigList := make([]string, 0)
	intConfig, err := sess.command("show configuration interfaces "+configInterface+" | display set", jnprSess)
	if err != nil {
		return "", err
	}
	for _, item := range strings.Split(intConfig, "\n") {
		if strings.Contains(item, "<configuration-output>") {
			continue
		}
		if strings.Contains(item, "</configuration-output>") {
			break
		}
		if item == "" {
			continue
		}
		itemTrim := strings.TrimPrefix(item, "set interfaces ")
		matched, err := regexp.MatchString(match, itemTrim)
		if err != nil {
			return "", err
		}
		if !matched {
			continue
		}
		itemTrimSplit := strings.Split(itemTrim, " ")
		switch len(itemTrimSplit) {
		case 0:
			continue
		case 1:
			intConfigList = append(intConfigList, itemTrimSplit[0])
		case 2:
			intConfigList = append(intConfigList, itemTrimSplit[0])
		default:
			if itemTrimSplit[1] == "unit" {
				intConfigList = append(intConfigList, itemTrimSplit[0]+"."+itemTrimSplit[2])
			} else {
				intConfigList = append(intConfigList, itemTrimSplit[0])
			}
		}
	}
	intConfigList = uniqueListString(intConfigList)
	if len(intConfigList) == 0 {
		return "", nil
	}
	if len(intConfigList) > 1 {
		return "", fmt.Errorf("too many different interfaces found")
	}

	return intConfigList[0], nil
}
