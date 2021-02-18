package junos

import (
	"context"
	"fmt"
	"strings"

        "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
        "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
        "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type evpnOptions struct {
	encapsulation	string
	multicastMode	string
	routingInstance		string
        routeDistinguisher string
        vrfExport []string
        vrfImport []string
        vrfTarget string
        vtepSourceInterface string

}

func resourceProtocolEvpn() *schema.Resource {
	return &schema.Resource{
                CreateContext: resourceProtocolEvpnCreate,
                ReadContext:   resourceProtocolEvpnRead,
                UpdateContext: resourceProtocolEvpnUpdate,
                DeleteContext: resourceProtocolEvpnDelete,
                Importer: &schema.ResourceImporter{
                        State: resourceProtocolEvpnImport,
                },
		Schema: map[string]*schema.Schema{
			"encapsulation": {
				Type:	schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"mpls","vxlan"}, false),
			},
			"multicast_mode": {
				Type:	schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"ingress-replication"}, false),
			},
			"routing_instance": {
				Type:	schema.TypeString,
				Optional: true,
				Default: defaultWord,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
                        "route_distinguisher": {
                                Type:             schema.TypeString,
                                Required:         true,
                        },
                        "vrf_import": {
                                Type:     schema.TypeList,
                                Optional: true,
                                MinItems: 1,
                                Elem:     &schema.Schema{Type: schema.TypeString},
                        },
                        "vrf_export": {
                                Type:     schema.TypeList,
                                Optional: true,
                                MinItems: 1,
                                Elem:     &schema.Schema{Type: schema.TypeString},
                        },
                        "vrf_target": {
                                Type:             schema.TypeString,
                                Optional:         true,
                        },
			"vtep_source_interface": {
                                Type:     schema.TypeString,
                                Optional: true,
                        },

		},
	}
}

func resourceProtocolEvpnCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
        sess := m.(*Session)
        jnprSess, err := sess.startNewSession()
        if err != nil {
                return diag.FromErr(err)
        }
        defer sess.closeSession(jnprSess)
        sess.configLock(jnprSess)
        if d.Get("routing_instance").(string) != defaultWord {
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
        protocolEvpnExists, err := checkProtocolEvpnExists(d.Get("routing_instance").(string), m, jnprSess)
        if err != nil {
                sess.configClear(jnprSess)

                return diag.FromErr(err)
        }
        if protocolEvpnExists {
                sess.configClear(jnprSess)

                return diag.FromErr(fmt.Errorf("protocol evpn already exists in routing-instance %v",
                        d.Get("routing_instance").(string)))
        }
        if err := setProtocolEvpn(d, m, jnprSess); err != nil {
                sess.configClear(jnprSess)

                return diag.FromErr(err)
        }
        var diagWarns diag.Diagnostics
        warns, err := sess.commitConf("create resource protocol_evpn", jnprSess)
        appendDiagWarns(&diagWarns, warns)
        if err != nil {
                sess.configClear(jnprSess)

                return append(diagWarns, diag.FromErr(err)...)
        }
        protocolEvpnExists, err = checkProtocolEvpnExists(d.Get("routing_instance").(string), m, jnprSess)
        if err != nil {
                return append(diagWarns, diag.FromErr(err)...)
        }
        if protocolEvpnExists {
                d.SetId("protocol_evpn" + idSeparator + d.Get("routing_instance").(string))
        } else {
                return append(diagWarns, diag.FromErr(fmt.Errorf("protocol evpn not exists in routing-instance %v after commit "+
                        "=> check your config", d.Get("routing_instance").(string)))...)
        }

        return append(diagWarns, resourceProtocolEvpnReadWJnprSess(d, m, jnprSess)...)
}

func resourceProtocolEvpnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
        sess := m.(*Session)
        jnprSess, err := sess.startNewSession()
        if err != nil {
                return diag.FromErr(err)
        }
        defer sess.closeSession(jnprSess)

        return resourceProtocolEvpnReadWJnprSess(d, m, jnprSess)
}
func resourceProtocolEvpnReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
        mutex.Lock()
        protocolEvpnOptions, err := readProtocolEvpn(d.Get("routing_instance").(string), m, jnprSess)
        mutex.Unlock()
        if err != nil {
                return diag.FromErr(err)
        }
        fillProtocolEvpnData(d, protocolEvpnOptions)

        return nil
}

func resourceProtocolEvpnUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
        d.Partial(true)
        sess := m.(*Session)
        jnprSess, err := sess.startNewSession()
        if err != nil {
                return diag.FromErr(err)
        }
        defer sess.closeSession(jnprSess)
        sess.configLock(jnprSess)
        //if err := delProtocolEvpnOpts(d, m, jnprSess); err != nil {
        if err := delProtocolEvpn(d, m, jnprSess); err != nil {
                sess.configClear(jnprSess)

                return diag.FromErr(err)
        }
        if err := setProtocolEvpn(d, m, jnprSess); err != nil {
                sess.configClear(jnprSess)

                return diag.FromErr(err)
        }
        var diagWarns diag.Diagnostics
        warns, err := sess.commitConf("update resource protocol_evpn", jnprSess)
        appendDiagWarns(&diagWarns, warns)
        if err != nil {
                sess.configClear(jnprSess)

                return append(diagWarns, diag.FromErr(err)...)
        }
        d.Partial(false)

        return append(diagWarns, resourceProtocolEvpnReadWJnprSess(d, m, jnprSess)...)
}

func resourceProtocolEvpnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
        sess := m.(*Session)
        jnprSess, err := sess.startNewSession()
        if err != nil {
                return diag.FromErr(err)
        }
        defer sess.closeSession(jnprSess)
        sess.configLock(jnprSess)
        if err := delProtocolEvpn(d, m, jnprSess); err != nil {
                sess.configClear(jnprSess)

                return diag.FromErr(err)
        }
        var diagWarns diag.Diagnostics
        warns, err := sess.commitConf("delete resource protocol_evpn", jnprSess)
        appendDiagWarns(&diagWarns, warns)
        if err != nil {
                sess.configClear(jnprSess)

                return append(diagWarns, diag.FromErr(err)...)
        }

        return diagWarns
}

func resourceProtocolEvpnImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
        sess := m.(*Session)
        jnprSess, err := sess.startNewSession()
        if err != nil {
                return nil, err
        }
        defer sess.closeSession(jnprSess)
        result := make([]*schema.ResourceData, 1)
        idSplit := strings.Split(d.Id(), idSeparator)
        if len(idSplit) < 2 {
                return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
        }
        protocolEvpnExists, err := checkProtocolEvpnExists(idSplit[1], m, jnprSess)
        if err != nil {
                return nil, err
        }
        if !protocolEvpnExists {
                return nil, fmt.Errorf("don't find protocol evpn with id '%v' "+
                        "(id must be protocol_evpn"+idSeparator+"<routing_instance>)", d.Id())
        }
        protocolEvpnOptions, err := readProtocolEvpn(idSplit[1], m, jnprSess)
        if err != nil {
                return nil, err
        }
        fillProtocolEvpnData(d, protocolEvpnOptions)
        result[0] = d

        return result, nil
}

func checkProtocolEvpnExists(instance string, m interface{}, jnprSess *NetconfObject) (bool, error) {
        sess := m.(*Session)
        var protocolEvpnConfig string
        var err error
        if instance == defaultWord {
                protocolEvpnConfig, err = sess.command("show configuration protocols evpn | display set", jnprSess)
                if err != nil {
                        return false, err
                }
        } else {
                protocolEvpnConfig, err = sess.command("show configuration routing-instances "+
                        instance+" protocols evpn | display set", jnprSess)
                if err != nil {
                        return false, err
                }
        }
        if protocolEvpnConfig == emptyWord {
                return false, nil
        }

        return true, nil
}

func readProtocolEvpn(instance string, m interface{}, jnprSess *NetconfObject) (evpnOptions, error) {
	sess := m.(*Session)
	var confRead evpnOptions
	var protocolEvpnConfig string
	var switchOptionsConfig string
	var err error

	// Read protocol evpn settings
        if instance == defaultWord {
                protocolEvpnConfig, err = sess.command("show configuration protocols evpn | display set relative", jnprSess)
                if err != nil {
                        return confRead, err
                }
        } else {
                protocolEvpnConfig, err = sess.command("show configuration routing-instances "+
                        instance+" protocols evpn | display set relative", jnprSess)
                if err != nil {
                        return confRead, err
                }
        }
	if protocolEvpnConfig != emptyWord {
		confRead.routingInstance = instance
		for _, item := range strings.Split(protocolEvpnConfig, "\n") {
                        if strings.Contains(item, "<configuration-output>") {
                                continue
                        }
                        if strings.Contains(item, "</configuration-output>") {
                                break
                        }
                        itemTrim := strings.TrimPrefix(item, setLineStart)
                        switch {
			case strings.HasPrefix(itemTrim, "encapsulation "):
				confRead.encapsulation = strings.TrimPrefix(itemTrim, "encapsulation ")
			case strings.HasPrefix(itemTrim, "multicast-mode "):
				confRead.multicastMode = strings.TrimPrefix(itemTrim, "multicast-mode ")
			default:
				continue
			}
		}
	}
	// Read switchoptions settings
	switchOptionsConfig, err = sess.command("show configuration switch-options | display set relative", jnprSess)
	if switchOptionsConfig != emptyWord {
		for _, item := range strings.Split(switchOptionsConfig, "\n") {
                        if strings.Contains(item, "<configuration-output>") {
                                continue
                        }
                        if strings.Contains(item, "</configuration-output>") {
                                break
                        }
			itemTrim := strings.TrimPrefix(item, setLineStart)
                        switch {
                        case strings.HasPrefix(itemTrim, "route-distinguisher "):
                                confRead.routeDistinguisher = strings.TrimPrefix(itemTrim, "route-distinguisher ")
                        case strings.HasPrefix(itemTrim, "vtep-source-interface "):
                                confRead.vtepSourceInterface = strings.TrimPrefix(itemTrim, "vtep-source-interface ")
                        case strings.HasPrefix(itemTrim, "vrf-target "):
                                confRead.vrfTarget = strings.TrimPrefix(itemTrim, "vrf-target ")
                        case strings.HasPrefix(itemTrim, "vrf-import "):
                                confRead.vrfImport = append(confRead.vrfImport, strings.TrimPrefix(itemTrim, "vrf-import "))
                        case strings.HasPrefix(itemTrim, "vrf-export "):
                                confRead.vrfExport = append(confRead.vrfExport, strings.TrimPrefix(itemTrim, "vrf-export "))
                        }

		}
	}

	return confRead, nil
}
func delProtocolEvpn(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
        sess := m.(*Session)
        configSet := make([]string, 0, 1)
	switchOptLinesToDelete := []string{
                "route-distinguisher",
                "vrf-import",
                "vrf-export",
                "vrf-target",
                "vtep-source-interface",
	}
        if d.Get("routing_instance").(string) == defaultWord {
                configSet = append(configSet, "delete protocols evpn")
        } else {
                configSet = append(configSet, "delete routing-instances "+d.Get("routing_instance").(string)+
                        " protocols evpn")
        }
	delPrefix := "delete switch-options "
	for _, line := range switchOptLinesToDelete {
		configSet = append(configSet, delPrefix + line)
	}
        if err := sess.configSet(configSet, jnprSess); err != nil {
                return err
        }

        return nil
}


func delProtocolEvpnOpts(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
        sess := m.(*Session)
        configSet := make([]string, 0)
        delPrefix := deleteWord + " "
        if d.Get("routing_instance").(string) == defaultWord {
                delPrefix += "protocols evpn "
        } else {
                delPrefix += "routing-instances " + d.Get("routing_instance").(string) +
                        " protocols evpn "
        }
        configSet = append(configSet,
                delPrefix+"encapsulation",
                delPrefix+"multicast-mode",
        )

        if err := sess.configSet(configSet, jnprSess); err != nil {
                return err
        }

        return nil
}

func setProtocolEvpn(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
        setPrefix := setLineStart
	// Set protocol evpn settings
        if d.Get("routing_instance").(string) == defaultWord {
                setPrefix += "protocols evpn "
        } else {
                setPrefix += "routing-instances " + d.Get("routing_instance").(string) +
                        " protocols evpn "
        }
	if d.Get("encapsulation").(string) != "" {
		configSet = append(configSet, setPrefix + "encapsulation " + d.Get("encapsulation").(string))
	}
	if d.Get("multicast_mode").(string) != "" {
		configSet = append(configSet,  setPrefix + "multicast-mode " + d.Get("multicast_mode").(string))
	}
	// Set protocol switch-options settings
	setPrefix = "set switch-options "
        if d.Get("route_distinguisher").(string) != "" {
                configSet = append(configSet, setPrefix+"route-distinguisher "+d.Get("route_distinguisher").(string))
        }
        for _, v := range d.Get("vrf_import").([]interface{}) {
                configSet = append(configSet, setPrefix+" vrf-import "+v.(string))
        }
        for _, v := range d.Get("vrf_export").([]interface{}) {
                        configSet = append(configSet, setPrefix+" vrf-export "+v.(string))
        }
        if d.Get("vrf_target").(string) != "" {
                configSet = append(configSet, setPrefix+"vrf-target "+d.Get("vrf_target").(string))
        }
        if d.Get("vtep_source_interface").(string) != "" {
                configSet = append(configSet, setPrefix+"vtep-source-interface "+d.Get("vtep_source_interface").(string))
        }


	if err := sess.configSet(configSet, jnprSess); err != nil {
                return err
        }


	return nil
}
func fillProtocolEvpnData(d *schema.ResourceData, protocolEvpnOptions evpnOptions) {
        if tfErr := d.Set("encapsulation", protocolEvpnOptions.encapsulation); tfErr != nil {
                panic(tfErr)
        }
        if tfErr := d.Set("multicast_mode", protocolEvpnOptions.multicastMode); tfErr != nil {
                panic(tfErr)
        }
	// Fill switch-options
        if tfErr := d.Set("route_distinguisher", protocolEvpnOptions.routeDistinguisher); tfErr != nil {
                panic(tfErr)
        }
        if tfErr := d.Set("vrf_import", protocolEvpnOptions.vrfImport); tfErr != nil {
                panic(tfErr)
        }
        if tfErr := d.Set("vrf_export", protocolEvpnOptions.vrfExport); tfErr != nil {
                panic(tfErr)
        }
        if tfErr := d.Set("vrf_target", protocolEvpnOptions.vrfTarget); tfErr != nil {
                panic(tfErr)
        }
        if tfErr := d.Set("vtep_source_interface", protocolEvpnOptions.vtepSourceInterface); tfErr != nil {
                panic(tfErr)
        }




}
