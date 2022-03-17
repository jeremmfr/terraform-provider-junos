package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type groupDualSystemOptions struct {
	applyGroups    bool
	name           string
	interfaceFxp0  []map[string]interface{}
	routingOptions []map[string]interface{}
	security       []map[string]interface{}
	system         []map[string]interface{}
}

func resourceGroupDualSystem() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupDualSystemCreate,
		ReadContext:   resourceGroupDualSystemRead,
		UpdateContext: resourceGroupDualSystemUpdate,
		DeleteContext: resourceGroupDualSystemDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGroupDualSystemImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"node0", "node1", "re0", "re1"}, false),
			},
			"apply_groups": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"interface_fxp0": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"family_inet_address": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_ip": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validateIPMaskFunc(),
									},
									"master_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"preferred": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"primary": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"family_inet6_address": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_ip": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validateIPMaskFunc(),
									},
									"master_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"preferred": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"primary": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"routing_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"static_route": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsCIDRNetwork(0, 128),
									},
									"next_hop": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"security": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_source_address": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"system": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"backup_router_address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPv4Address,
						},
						"backup_router_destination": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"inet6_backup_router_address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateIsIPv6Address,
						},
						"inet6_backup_router_destination": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceGroupDualSystemCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setGroupDualSystem(d, m, nil); err != nil {
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
	groupDualSystemExists, err := checkGroupDualSystemExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if groupDualSystemExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("group %v already exists", d.Get("name").(string)))...)
	}

	if err := setGroupDualSystem(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_group_dual_system", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	groupDualSystemExists, err = checkGroupDualSystemExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if groupDualSystemExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("group %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceGroupDualSystemReadWJnprSess(d, m, jnprSess)...)
}

func resourceGroupDualSystemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceGroupDualSystemReadWJnprSess(d, m, jnprSess)
}

func resourceGroupDualSystemReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	groupDualSystemOpts, err := readGroupDualSystem(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if groupDualSystemOpts.name == "" {
		d.SetId("")
	} else {
		fillGroupDualSystemData(d, groupDualSystemOpts)
	}

	return nil
}

func resourceGroupDualSystemUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delGroupDualSystem(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if strings.HasPrefix(d.Get("name").(string), "node") {
			if err := sess.configSet([]string{"delete apply-groups \"${node}\""}, nil); err != nil {
				return diag.FromErr(err)
			}
		} else if err := sess.configSet([]string{"delete apply-groups " + d.Get("name").(string)}, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setGroupDualSystem(d, m, nil); err != nil {
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
	if err := delGroupDualSystem(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if strings.HasPrefix(d.Get("name").(string), "node") {
		if err := sess.configSet([]string{"delete apply-groups \"${node}\""}, jnprSess); err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	} else if err := sess.configSet([]string{"delete apply-groups " + d.Get("name").(string)}, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setGroupDualSystem(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_group_dual_system", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceGroupDualSystemReadWJnprSess(d, m, jnprSess)...)
}

func resourceGroupDualSystemDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delGroupDualSystem(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if strings.HasPrefix(d.Get("name").(string), "node") {
			if err := sess.configSet([]string{"delete apply-groups \"${node}\""}, nil); err != nil {
				return diag.FromErr(err)
			}
		} else if err := sess.configSet([]string{"delete apply-groups " + d.Get("name").(string)}, nil); err != nil {
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
	if err := delGroupDualSystem(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if strings.HasPrefix(d.Get("name").(string), "node") {
		if err := sess.configSet([]string{"delete apply-groups \"${node}\""}, jnprSess); err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	} else if err := sess.configSet([]string{"delete apply-groups " + d.Get("name").(string)}, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_group_dual_system", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceGroupDualSystemImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	if !bchk.StringInSlice(d.Id(), []string{"node0", "node1", "re0", "re1"}) {
		return nil, fmt.Errorf("invalid group id '%v' (id must be <name>)", d.Id())
	}
	groupDualSystemExists, err := checkGroupDualSystemExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !groupDualSystemExists {
		return nil, fmt.Errorf("don't find group with id '%v' (id must be <name>)", d.Id())
	}
	groupDualSystemOptions, err := readGroupDualSystem(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillGroupDualSystemData(d, groupDualSystemOptions)

	result[0] = d

	return result, nil
}

func checkGroupDualSystemExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command(cmdShowConfig+"groups "+name+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setGroupDualSystem(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	if d.Get("apply_groups").(bool) {
		if strings.HasPrefix(d.Get("name").(string), "node") {
			configSet = append(configSet, "set apply-groups \"${node}\"")
		} else {
			configSet = append(configSet, "set apply-groups "+d.Get("name").(string))
		}
	}
	setPrefix := "set groups " + d.Get("name").(string) + " "
	for _, v := range d.Get("interface_fxp0").([]interface{}) {
		if v == nil {
			return fmt.Errorf("interface_fxp0 block is empty")
		}
		interfaceFxp0 := v.(map[string]interface{})
		if v2 := interfaceFxp0["description"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"interfaces fxp0 description \""+v2+"\"")
		}
		familyInetAddressCIDRIPList := make([]string, 0)
		for _, v2 := range interfaceFxp0["family_inet_address"].([]interface{}) {
			familyInetAddress := v2.(map[string]interface{})
			if bchk.StringInSlice(familyInetAddress["cidr_ip"].(string), familyInetAddressCIDRIPList) {
				return fmt.Errorf("multiple blocks family_inet_address with the same cidr_ip %s",
					familyInetAddress["cidr_ip"].(string))
			}
			familyInetAddressCIDRIPList = append(familyInetAddressCIDRIPList, familyInetAddress["cidr_ip"].(string))
			configSet = append(configSet, setPrefix+"interfaces fxp0 unit 0 family inet address "+
				familyInetAddress["cidr_ip"].(string))
			if familyInetAddress["master_only"].(bool) {
				configSet = append(configSet, setPrefix+"interfaces fxp0 unit 0 family inet address "+
					familyInetAddress["cidr_ip"].(string)+" master-only")
			}
			if familyInetAddress["preferred"].(bool) {
				configSet = append(configSet, setPrefix+"interfaces fxp0 unit 0 family inet address "+
					familyInetAddress["cidr_ip"].(string)+" preferred")
			}
			if familyInetAddress["primary"].(bool) {
				configSet = append(configSet, setPrefix+"interfaces fxp0 unit 0 family inet address "+
					familyInetAddress["cidr_ip"].(string)+" primary")
			}
		}
		familyInet6AddressCIDRIPList := make([]string, 0)
		for _, v2 := range interfaceFxp0["family_inet6_address"].([]interface{}) {
			familyInet6Address := v2.(map[string]interface{})
			if bchk.StringInSlice(familyInet6Address["cidr_ip"].(string), familyInet6AddressCIDRIPList) {
				return fmt.Errorf("multiple blocks family_inet6_address with the same cidr_ip %s",
					familyInet6Address["cidr_ip"].(string))
			}
			familyInet6AddressCIDRIPList = append(familyInet6AddressCIDRIPList, familyInet6Address["cidr_ip"].(string))
			configSet = append(configSet, setPrefix+"interfaces fxp0 unit 0 family inet6 address "+
				familyInet6Address["cidr_ip"].(string))
			if familyInet6Address["master_only"].(bool) {
				configSet = append(configSet, setPrefix+"interfaces fxp0 unit 0 family inet6 address "+
					familyInet6Address["cidr_ip"].(string)+" master-only")
			}
			if familyInet6Address["preferred"].(bool) {
				configSet = append(configSet, setPrefix+"interfaces fxp0 unit 0 family inet6 address "+
					familyInet6Address["cidr_ip"].(string)+" preferred")
			}
			if familyInet6Address["primary"].(bool) {
				configSet = append(configSet, setPrefix+"interfaces fxp0 unit 0 family inet6 address "+
					familyInet6Address["cidr_ip"].(string)+" primary")
			}
		}
	}
	for _, v := range d.Get("routing_options").([]interface{}) {
		routingOptions := v.(map[string]interface{})
		staticRouteDestList := make([]string, 0)
		for _, v2 := range routingOptions["static_route"].([]interface{}) {
			staticRoute := v2.(map[string]interface{})
			if bchk.StringInSlice(staticRoute["destination"].(string), staticRouteDestList) {
				return fmt.Errorf("multiple blocks static_route with the same destination %s", staticRoute["destination"].(string))
			}
			staticRouteDestList = append(staticRouteDestList, staticRoute["destination"].(string))
			for _, v3 := range staticRoute["next_hop"].([]interface{}) {
				configSet = append(configSet, setPrefix+"routing-options static route "+
					staticRoute["destination"].(string)+" next-hop "+v3.(string))
			}
		}
	}
	for _, v := range d.Get("security").([]interface{}) {
		security := v.(map[string]interface{})
		if v2 := security["log_source_address"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"security log source-address "+v2)
		}
	}
	for _, v := range d.Get("system").([]interface{}) {
		if v == nil {
			return fmt.Errorf("system block is empty")
		}
		system := v.(map[string]interface{})
		if v2 := system["host_name"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+" system host-name \""+v2+"\"")
		}
		if v2 := system["backup_router_address"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+" system backup-router "+v2)
		}
		for _, v2 := range sortSetOfString(system["backup_router_destination"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefix+" system backup-router destination "+v2)
		}
		if v2 := system["inet6_backup_router_address"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+" system inet6-backup-router "+v2)
		}
		for _, v2 := range sortSetOfString(system["inet6_backup_router_destination"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefix+" system inet6-backup-router destination "+v2)
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func readGroupDualSystem(group string, m interface{}, jnprSess *NetconfObject) (groupDualSystemOptions, error) {
	sess := m.(*Session)
	var confRead groupDualSystemOptions

	showConfig, err := sess.command(cmdShowConfig+"groups "+group+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = group
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "interfaces fxp0 "):
				if len(confRead.interfaceFxp0) == 0 {
					confRead.interfaceFxp0 = append(confRead.interfaceFxp0, map[string]interface{}{
						"description":          "",
						"family_inet_address":  make([]map[string]interface{}, 0),
						"family_inet6_address": make([]map[string]interface{}, 0),
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "interfaces fxp0 description "):
					confRead.interfaceFxp0[0]["description"] = strings.TrimPrefix(itemTrim, "interfaces fxp0 description ")
				case strings.HasPrefix(itemTrim, "interfaces fxp0 unit 0 family inet address "):
					itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "interfaces fxp0 unit 0 family inet address "), " ")
					familyInetAddress := map[string]interface{}{
						"cidr_ip":     itemTrimSplit[0],
						"master_only": false,
						"preferred":   false,
						"primary":     false,
					}
					confRead.interfaceFxp0[0]["family_inet_address"] = copyAndRemoveItemMapList(
						"cidr_ip", familyInetAddress, confRead.interfaceFxp0[0]["family_inet_address"].([]map[string]interface{}))
					switch {
					case strings.HasSuffix(itemTrim, "master-only"):
						familyInetAddress["master_only"] = true
					case strings.HasSuffix(itemTrim, "preferred"):
						familyInetAddress["preferred"] = true
					case strings.HasSuffix(itemTrim, "primary"):
						familyInetAddress["primary"] = true
					}
					confRead.interfaceFxp0[0]["family_inet_address"] = append(
						confRead.interfaceFxp0[0]["family_inet_address"].([]map[string]interface{}), familyInetAddress)
				case strings.HasPrefix(itemTrim, "interfaces fxp0 unit 0 family inet6 address "):
					itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "interfaces fxp0 unit 0 family inet6 address "), " ")
					familyInet6Address := map[string]interface{}{
						"cidr_ip":     itemTrimSplit[0],
						"master_only": false,
						"preferred":   false,
						"primary":     false,
					}
					confRead.interfaceFxp0[0]["family_inet6_address"] = copyAndRemoveItemMapList(
						"cidr_ip", familyInet6Address, confRead.interfaceFxp0[0]["family_inet6_address"].([]map[string]interface{}))
					switch {
					case strings.HasSuffix(itemTrim, "master-only"):
						familyInet6Address["master_only"] = true
					case strings.HasSuffix(itemTrim, "preferred"):
						familyInet6Address["preferred"] = true
					case strings.HasSuffix(itemTrim, "primary"):
						familyInet6Address["primary"] = true
					}
					confRead.interfaceFxp0[0]["family_inet6_address"] = append(
						confRead.interfaceFxp0[0]["family_inet6_address"].([]map[string]interface{}), familyInet6Address)
				}
			case strings.HasPrefix(itemTrim, "routing-options static route "):
				if len(confRead.routingOptions) == 0 {
					confRead.routingOptions = append(confRead.routingOptions, map[string]interface{}{
						"static_route": make([]map[string]interface{}, 0),
					})
				}
				routeTrim := strings.TrimPrefix(itemTrim, "routing-options static route ")
				routeTrimSplit := strings.Split(routeTrim, " ")
				destOptions := map[string]interface{}{
					"destination": routeTrimSplit[0],
					"next_hop":    make([]string, 0),
				}
				confRead.routingOptions[0]["static_route"] = copyAndRemoveItemMapList(
					"destination", destOptions, confRead.routingOptions[0]["static_route"].([]map[string]interface{}))
				if strings.HasPrefix(routeTrim, routeTrimSplit[0]+" next-hop ") {
					destOptions["next_hop"] = append(destOptions["next_hop"].([]string),
						strings.TrimPrefix(routeTrim, routeTrimSplit[0]+" next-hop "))
				}
				confRead.routingOptions[0]["static_route"] = append(
					confRead.routingOptions[0]["static_route"].([]map[string]interface{}), destOptions)
			case strings.HasPrefix(itemTrim, "security"):
				if len(confRead.security) == 0 {
					confRead.security = append(confRead.security, map[string]interface{}{
						"log_source_address": "",
					})
				}
				if strings.HasPrefix(itemTrim, "security log source-address ") {
					confRead.security[0]["log_source_address"] = strings.TrimPrefix(
						itemTrim, "security log source-address ")
				}
			case strings.HasPrefix(itemTrim, "system"):
				if len(confRead.system) == 0 {
					confRead.system = append(confRead.system, map[string]interface{}{
						"host_name":                       "",
						"backup_router_address":           "",
						"backup_router_destination":       make([]string, 0),
						"inet6_backup_router_address":     "",
						"inet6_backup_router_destination": make([]string, 0),
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "system host-name "):
					confRead.system[0]["host_name"] = strings.Trim(strings.TrimPrefix(itemTrim, "system host-name "), "\"")
				case strings.HasPrefix(itemTrim, "system backup-router destination "):
					confRead.system[0]["backup_router_destination"] = append(
						confRead.system[0]["backup_router_destination"].([]string),
						strings.TrimPrefix(itemTrim, "system backup-router destination "))
				case strings.HasPrefix(itemTrim, "system backup-router "):
					confRead.system[0]["backup_router_address"] = strings.TrimPrefix(itemTrim, "system backup-router ")
				case strings.HasPrefix(itemTrim, "system inet6-backup-router destination "):
					confRead.system[0]["inet6_backup_router_destination"] = append(
						confRead.system[0]["inet6_backup_router_destination"].([]string),
						strings.TrimPrefix(itemTrim, "system inet6-backup-router destination "))
				case strings.HasPrefix(itemTrim, "system inet6-backup-router "):
					confRead.system[0]["inet6_backup_router_address"] = strings.TrimPrefix(itemTrim, "system inet6-backup-router ")
				}
			}
		}
	}
	showConfigApplyGroups, err := sess.command(cmdShowConfig+"apply-groups"+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfigApplyGroups != emptyW {
		confRead.name = group
		for _, item := range strings.Split(showConfigApplyGroups, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			switch {
			case item == "set "+confRead.name+" ":
				confRead.applyGroups = true
			case item == "set \"${node}\" " && strings.HasPrefix(confRead.name, "node"):
				confRead.applyGroups = true
			}
		}
	}

	return confRead, nil
}

func delGroupDualSystem(group string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete groups "+group)

	return sess.configSet(configSet, jnprSess)
}

func fillGroupDualSystemData(d *schema.ResourceData, groupDualSystemOptions groupDualSystemOptions) {
	if tfErr := d.Set("name", groupDualSystemOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("apply_groups", groupDualSystemOptions.applyGroups); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interface_fxp0", groupDualSystemOptions.interfaceFxp0); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_options", groupDualSystemOptions.routingOptions); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("security", groupDualSystemOptions.security); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("system", groupDualSystemOptions.system); tfErr != nil {
		panic(tfErr)
	}
}
