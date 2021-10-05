package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type snmpViewOptions struct {
	name       string
	oidInclude []string
	oidExclude []string
}

func resourceSnmpView() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnmpViewCreate,
		ReadContext:   resourceSnmpViewRead,
		UpdateContext: resourceSnmpViewUpdate,
		DeleteContext: resourceSnmpViewDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSnmpViewImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oid_include": {
				Type:         schema.TypeSet,
				Optional:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"oid_include", "oid_exclude"},
			},
			"oid_exclude": {
				Type:         schema.TypeSet,
				Optional:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"oid_include", "oid_exclude"},
			},
		},
	}
}

func resourceSnmpViewCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSnmpView(d, m, nil); err != nil {
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
	snmpViewExists, err := checkSnmpViewExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpViewExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp view %v already exists", d.Get("name").(string)))...)
	}

	if err := setSnmpView(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_snmp_view", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	snmpViewExists, err = checkSnmpViewExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpViewExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp view %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSnmpViewReadWJnprSess(d, m, jnprSess)...)
}

func resourceSnmpViewRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSnmpViewReadWJnprSess(d, m, jnprSess)
}

func resourceSnmpViewReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	snmpViewOptions, err := readSnmpView(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if snmpViewOptions.name == "" {
		d.SetId("")
	} else {
		fillSnmpViewData(d, snmpViewOptions)
	}

	return nil
}

func resourceSnmpViewUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSnmpView(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmpView(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_snmp_view", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpViewReadWJnprSess(d, m, jnprSess)...)
}

func resourceSnmpViewDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSnmpView(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_snmp_view", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSnmpViewImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	snmpViewExists, err := checkSnmpViewExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !snmpViewExists {
		return nil, fmt.Errorf("don't find snmp view with id '%v' (id must be <name>)", d.Id())
	}
	snmpViewOptions, err := readSnmpView(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSnmpViewData(d, snmpViewOptions)

	result[0] = d

	return result, nil
}

func checkSnmpViewExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration snmp view \""+name+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setSnmpView(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set snmp view \"" + d.Get("name").(string) + "\" "
	configSet := make([]string, 0)

	for _, v := range sortSetOfString(d.Get("oid_include").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"oid "+v+" include")
	}
	for _, v := range sortSetOfString(d.Get("oid_exclude").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"oid "+v+" exclude")
	}

	return sess.configSet(configSet, jnprSess)
}

func readSnmpView(name string, m interface{}, jnprSess *NetconfObject) (snmpViewOptions, error) {
	sess := m.(*Session)
	var confRead snmpViewOptions

	showConfig, err := sess.command("show configuration snmp view \""+name+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			itemTrimSplit := strings.Split(itemTrim, " ")
			switch {
			case strings.HasPrefix(itemTrim, "oid ") && strings.HasSuffix(itemTrim, " include"):
				confRead.oidInclude = append(confRead.oidInclude, itemTrimSplit[1])
			case strings.HasPrefix(itemTrim, "oid ") && strings.HasSuffix(itemTrim, " exclude"):
				confRead.oidExclude = append(confRead.oidExclude, itemTrimSplit[1])
			case strings.HasPrefix(itemTrim, "oid "):
				confRead.oidInclude = append(confRead.oidInclude, itemTrimSplit[1])
			}
		}
	}

	return confRead, nil
}

func delSnmpView(name string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete snmp view \"" + name + "\""}

	return sess.configSet(configSet, jnprSess)
}

func fillSnmpViewData(d *schema.ResourceData, snmpViewOptions snmpViewOptions) {
	if tfErr := d.Set("name", snmpViewOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("oid_include", snmpViewOptions.oidInclude); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("oid_exclude", snmpViewOptions.oidExclude); tfErr != nil {
		panic(tfErr)
	}
}
