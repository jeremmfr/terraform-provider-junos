package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jeremmfr/go-utils/basiccheck"
)

type snmpV3VacmSecurityToGroupOptions struct {
	name  string
	model string
	group string
}

func resourceSnmpV3VacmSecurityToGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnmpV3VacmSecurityToGroupCreate,
		ReadContext:   resourceSnmpV3VacmSecurityToGroupRead,
		UpdateContext: resourceSnmpV3VacmSecurityToGroupUpdate,
		DeleteContext: resourceSnmpV3VacmSecurityToGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSnmpV3VacmSecurityToGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"model": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"usm", "v1", "v2c"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSnmpV3VacmSecurityToGroupCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSnmpV3VacmSecurityToGroup(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("model").(string) + idSeparator + d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	snmpV3VacmSecurityToGroupExists, err := checkSnmpV3VacmSecurityToGroupExists(
		d.Get("model").(string), d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3VacmSecurityToGroupExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"snmp v3 vacm security-to-group security-model %v security-name %v already exists",
			d.Get("model").(string), d.Get("name").(string)))...)
	}

	if err := setSnmpV3VacmSecurityToGroup(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_snmp_v3_vacm_securitytogroup", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	snmpV3VacmSecurityToGroupExists, err = checkSnmpV3VacmSecurityToGroupExists(
		d.Get("model").(string), d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3VacmSecurityToGroupExists {
		d.SetId(d.Get("model").(string) + idSeparator + d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"snmp v3 vacm security-to-group security-model %v security-name %v not exists after commit "+
				"=> check your config", d.Get("model").(string), d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSnmpV3VacmSecurityToGroupReadWJnprSess(d, m, jnprSess)...)
}

func resourceSnmpV3VacmSecurityToGroupRead(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSnmpV3VacmSecurityToGroupReadWJnprSess(d, m, jnprSess)
}

func resourceSnmpV3VacmSecurityToGroupReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	snmpV3VacmSecurityToGroupOptions, err := readSnmpV3VacmSecurityToGroup(
		d.Get("model").(string), d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if snmpV3VacmSecurityToGroupOptions.name == "" {
		d.SetId("")
	} else {
		fillSnmpV3VacmSecurityToGroupData(d, snmpV3VacmSecurityToGroupOptions)
	}

	return nil
}

func resourceSnmpV3VacmSecurityToGroupUpdate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delSnmpV3VacmSecurityToGroup(d.Get("model").(string), d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSnmpV3VacmSecurityToGroup(d, m, nil); err != nil {
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
	if err := delSnmpV3VacmSecurityToGroup(d.Get("model").(string), d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmpV3VacmSecurityToGroup(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_snmp_v3_vacm_securitytogroup", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpV3VacmSecurityToGroupReadWJnprSess(d, m, jnprSess)...)
}

func resourceSnmpV3VacmSecurityToGroupDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delSnmpV3VacmSecurityToGroup(d.Get("model").(string), d.Get("name").(string), m, nil); err != nil {
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
	if err := delSnmpV3VacmSecurityToGroup(d.Get("model").(string), d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_snmp_v3_vacm_securitytogroup", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSnmpV3VacmSecurityToGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) != 2 {
		return nil, fmt.Errorf("can't find snmp v3 vacm security-to-group "+
			"with id '%v' (id must be <model>%s<name>)", d.Id(), idSeparator)
	}
	if !basiccheck.StringInSlice(idSplit[0], []string{"usm", "v1", "v2c"}) {
		return nil, fmt.Errorf("can't find snmp v3 vacm security-to-group "+
			"with id '%v' (id must be <model>%s<name>)", d.Id(), idSeparator)
	}
	snmpV3VacmSecurityToGroupExists, err := checkSnmpV3VacmSecurityToGroupExists(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !snmpV3VacmSecurityToGroupExists {
		return nil, fmt.Errorf("don't find snmp v3 vacm security-to-group "+
			"with id '%v' (id must be <model>%s<name>)", d.Id(), idSeparator)
	}
	snmpV3VacmSecurityToGroupOptions, err := readSnmpV3VacmSecurityToGroup(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSnmpV3VacmSecurityToGroupData(d, snmpV3VacmSecurityToGroupOptions)

	result[0] = d

	return result, nil
}

func checkSnmpV3VacmSecurityToGroupExists(model, name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration snmp v3 vacm security-to-group "+
		"security-model "+model+" security-name \""+name+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setSnmpV3VacmSecurityToGroup(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	configSet := make([]string, 1)
	if group := d.Get("group").(string); group != "" {
		configSet[0] = "set snmp v3 vacm security-to-group" +
			" security-model " + d.Get("model").(string) +
			" security-name \"" + d.Get("name").(string) + "\"" +
			" group \"" + d.Get("group").(string) + "\""
	} else {
		return fmt.Errorf("group need to be set")
	}

	return sess.configSet(configSet, jnprSess)
}

func readSnmpV3VacmSecurityToGroup(model, name string, m interface{}, jnprSess *NetconfObject,
) (snmpV3VacmSecurityToGroupOptions, error) {
	sess := m.(*Session)
	var confRead snmpV3VacmSecurityToGroupOptions

	showConfig, err := sess.command("show configuration snmp v3 vacm security-to-group "+
		"security-model "+model+" security-name \""+name+"\"  | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		confRead.model = model
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			if strings.HasPrefix(item, setLineStart+"group ") {
				confRead.group = strings.Trim(strings.TrimPrefix(item, setLineStart+"group "), "\"")
			}
		}
	}

	return confRead, nil
}

func delSnmpV3VacmSecurityToGroup(model, name string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete snmp v3 vacm security-to-group " +
		"security-model " + model + " security-name \"" + name + "\""}

	return sess.configSet(configSet, jnprSess)
}

func fillSnmpV3VacmSecurityToGroupData(
	d *schema.ResourceData, snmpV3VacmSecurityToGroupOptions snmpV3VacmSecurityToGroupOptions) {
	if tfErr := d.Set("model", snmpV3VacmSecurityToGroupOptions.model); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("name", snmpV3VacmSecurityToGroupOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("group", snmpV3VacmSecurityToGroupOptions.group); tfErr != nil {
		panic(tfErr)
	}
}
