package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type lldpInterfaceOptions struct {
	disable                 bool
	enable                  bool
	trapNotificationDisable bool
	trapNotificationEnable  bool
	name                    string
	powerNegotiation        []map[string]interface{}
}

func resourceLldpInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLldpInterfaceCreate,
		ReadContext:   resourceLldpInterfaceRead,
		UpdateContext: resourceLldpInterfaceUpdate,
		DeleteContext: resourceLldpInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceLldpInterfaceImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if strings.Count(value, ".") > 0 {
						errors = append(errors, fmt.Errorf(
							"%q in %q cannot have a dot", value, k))
					}

					return
				},
			},
			"disable": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"enable"},
			},
			"enable": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"disable"},
			},
			"power_negotiation": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"power_negotiation.0.enable"},
						},
						"enable": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"power_negotiation.0.disable"},
						},
					},
				},
			},
			"trap_notification_disable": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"trap_notification_enable"},
			},
			"trap_notification_enable": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"trap_notification_disable"},
			},
		},
	}
}

func resourceLldpInterfaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setLldpInterface(d, m, nil); err != nil {
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
	lldpInterfaceExists, err := checkLldpInterfaceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if lldpInterfaceExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"protocols lldp interface %v already exists", d.Get("name").(string)))...)
	}

	if err := setLldpInterface(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_lldp_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	lldpInterfaceExists, err = checkLldpInterfaceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if lldpInterfaceExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("protocols lldp interface %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceLldpInterfaceReadWJnprSess(d, m, jnprSess)...)
}

func resourceLldpInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceLldpInterfaceReadWJnprSess(d, m, jnprSess)
}

func resourceLldpInterfaceReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	lldpInterfaceOptions, err := readLldpInterface(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if lldpInterfaceOptions.name == "" {
		d.SetId("")
	} else {
		fillLldpInterfaceData(d, lldpInterfaceOptions)
	}

	return nil
}

func resourceLldpInterfaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delLldpInterface(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setLldpInterface(d, m, nil); err != nil {
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
	if err := delLldpInterface(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setLldpInterface(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_lldp_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceLldpInterfaceReadWJnprSess(d, m, jnprSess)...)
}

func resourceLldpInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delLldpInterface(d.Get("name").(string), m, nil); err != nil {
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
	if err := delLldpInterface(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_lldp_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceLldpInterfaceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	lldpInterfaceExists, err := checkLldpInterfaceExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !lldpInterfaceExists {
		return nil, fmt.Errorf("don't find protocols lldp interface with id '%v' (id must be <name>)", d.Id())
	}
	lldpInterfaceOptions, err := readLldpInterface(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillLldpInterfaceData(d, lldpInterfaceOptions)

	result[0] = d

	return result, nil
}

func checkLldpInterfaceExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command(
		cmdShowConfig+"protocols lldp interface "+name+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setLldpInterface(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set protocols lldp interface " + d.Get("name").(string) + " "
	configSet := make([]string, 0)

	configSet = append(configSet, setPrefix)
	if d.Get("disable").(bool) {
		configSet = append(configSet, setPrefix+"disable")
	}
	if d.Get("enable").(bool) {
		configSet = append(configSet, setPrefix+"enable")
	}
	for _, mPwNego := range d.Get("power_negotiation").([]interface{}) {
		configSet = append(configSet, setPrefix+"power-negotiation")
		if mPwNego != nil {
			powerNegotiation := mPwNego.(map[string]interface{})
			if powerNegotiation["disable"].(bool) {
				configSet = append(configSet, setPrefix+"power-negotiation disable")
			}
			if powerNegotiation["enable"].(bool) {
				configSet = append(configSet, setPrefix+"power-negotiation enable")
			}
		}
	}
	if d.Get("trap_notification_disable").(bool) {
		configSet = append(configSet, setPrefix+"trap-notification disable")
	}
	if d.Get("trap_notification_enable").(bool) {
		configSet = append(configSet, setPrefix+"trap-notification enable")
	}

	return sess.configSet(configSet, jnprSess)
}

func readLldpInterface(name string, m interface{}, jnprSess *NetconfObject,
) (lldpInterfaceOptions, error) {
	sess := m.(*Session)
	var confRead lldpInterfaceOptions

	showConfig, err := sess.command(
		cmdShowConfig+"protocols lldp interface "+name+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case itemTrim == disableW:
				confRead.disable = true
			case itemTrim == "enable":
				confRead.enable = true
			case strings.HasPrefix(itemTrim, "power-negotiation"):
				if len(confRead.powerNegotiation) == 0 {
					confRead.powerNegotiation = append(confRead.powerNegotiation, map[string]interface{}{
						"disable": false,
						"enable":  false,
					})
				}
				switch {
				case itemTrim == "power-negotiation disable":
					confRead.powerNegotiation[0]["disable"] = true
				case itemTrim == "power-negotiation enable":
					confRead.powerNegotiation[0]["enable"] = true
				}
			case itemTrim == "trap-notification disable":
				confRead.trapNotificationDisable = true
			case itemTrim == "trap-notification enable":
				confRead.trapNotificationEnable = true
			}
		}
	}

	return confRead, nil
}

func delLldpInterface(name string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	configSet := []string{"delete protocols lldp interface " + name}

	return sess.configSet(configSet, jnprSess)
}

func fillLldpInterfaceData(d *schema.ResourceData, lldpInterfaceOptions lldpInterfaceOptions) {
	if tfErr := d.Set("name", lldpInterfaceOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("disable", lldpInterfaceOptions.disable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("enable", lldpInterfaceOptions.enable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("power_negotiation", lldpInterfaceOptions.powerNegotiation); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("trap_notification_disable", lldpInterfaceOptions.trapNotificationDisable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("trap_notification_enable", lldpInterfaceOptions.trapNotificationEnable); tfErr != nil {
		panic(tfErr)
	}
}
