package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
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
		CreateWithoutTimeout: resourceLldpInterfaceCreate,
		ReadWithoutTimeout:   resourceLldpInterfaceRead,
		UpdateWithoutTimeout: resourceLldpInterfaceUpdate,
		DeleteWithoutTimeout: resourceLldpInterfaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLldpInterfaceImport,
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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setLldpInterface(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	lldpInterfaceExists, err := checkLldpInterfaceExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if lldpInterfaceExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"protocols lldp interface %v already exists", d.Get("name").(string)))...)
	}

	if err := setLldpInterface(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_lldp_interface", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	lldpInterfaceExists, err = checkLldpInterfaceExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if lldpInterfaceExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("protocols lldp interface %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceLldpInterfaceReadWJunSess(d, clt, junSess)...)
}

func resourceLldpInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceLldpInterfaceReadWJunSess(d, clt, junSess)
}

func resourceLldpInterfaceReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	lldpInterfaceOptions, err := readLldpInterface(d.Get("name").(string), clt, junSess)
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
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delLldpInterface(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setLldpInterface(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delLldpInterface(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setLldpInterface(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_lldp_interface", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceLldpInterfaceReadWJunSess(d, clt, junSess)...)
}

func resourceLldpInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delLldpInterface(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delLldpInterface(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_lldp_interface", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceLldpInterfaceImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)

	lldpInterfaceExists, err := checkLldpInterfaceExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !lldpInterfaceExists {
		return nil, fmt.Errorf("don't find protocols lldp interface with id '%v' (id must be <name>)", d.Id())
	}
	lldpInterfaceOptions, err := readLldpInterface(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillLldpInterfaceData(d, lldpInterfaceOptions)

	result[0] = d

	return result, nil
}

func checkLldpInterfaceExists(name string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(
		cmdShowConfig+"protocols lldp interface "+name+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setLldpInterface(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
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

	return clt.configSet(configSet, junSess)
}

func readLldpInterface(name string, clt *Client, junSess *junosSession,
) (confRead lldpInterfaceOptions, err error) {
	showConfig, err := clt.command(
		cmdShowConfig+"protocols lldp interface "+name+pipeDisplaySetRelative, junSess)
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
			case balt.CutPrefixInString(&itemTrim, "power-negotiation"):
				if len(confRead.powerNegotiation) == 0 {
					confRead.powerNegotiation = append(confRead.powerNegotiation, map[string]interface{}{
						"disable": false,
						"enable":  false,
					})
				}
				switch {
				case itemTrim == " disable":
					confRead.powerNegotiation[0]["disable"] = true
				case itemTrim == " enable":
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

func delLldpInterface(name string, clt *Client, junSess *junosSession) error {
	configSet := []string{"delete protocols lldp interface " + name}

	return clt.configSet(configSet, junSess)
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
