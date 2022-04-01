package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type snmpClientlistOptions struct {
	name   string
	prefix []string
}

func resourceSnmpClientlist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnmpClientlistCreate,
		ReadContext:   resourceSnmpClientlistRead,
		UpdateContext: resourceSnmpClientlistUpdate,
		DeleteContext: resourceSnmpClientlistDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSnmpClientlistImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"prefix": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSnmpClientlistCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSnmpClientlist(d, m, nil); err != nil {
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
	snmpClientlistExists, err := checkSnmpClientlistExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpClientlistExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp client-list %v already exists", d.Get("name").(string)))...)
	}

	if err := setSnmpClientlist(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_snmp_clientlist", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	snmpClientlistExists, err = checkSnmpClientlistExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpClientlistExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp client-list %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSnmpClientlistReadWJnprSess(d, m, jnprSess)...)
}

func resourceSnmpClientlistRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSnmpClientlistReadWJnprSess(d, m, jnprSess)
}

func resourceSnmpClientlistReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
) diag.Diagnostics {
	mutex.Lock()
	snmpClientlistOptions, err := readSnmpClientlist(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if snmpClientlistOptions.name == "" {
		d.SetId("")
	} else {
		fillSnmpClientlistData(d, snmpClientlistOptions)
	}

	return nil
}

func resourceSnmpClientlistUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delSnmpClientlist(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSnmpClientlist(d, m, nil); err != nil {
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
	if err := delSnmpClientlist(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmpClientlist(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_snmp_clientlist", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpClientlistReadWJnprSess(d, m, jnprSess)...)
}

func resourceSnmpClientlistDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delSnmpClientlist(d.Get("name").(string), m, nil); err != nil {
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
	if err := delSnmpClientlist(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_snmp_clientlist", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSnmpClientlistImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	snmpClientlistExists, err := checkSnmpClientlistExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !snmpClientlistExists {
		return nil, fmt.Errorf("don't find snmp client-list with id '%v' (id must be <name>)", d.Id())
	}
	snmpClientlistOptions, err := readSnmpClientlist(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSnmpClientlistData(d, snmpClientlistOptions)

	result[0] = d

	return result, nil
}

func checkSnmpClientlistExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command(cmdShowConfig+"snmp client-list \""+name+"\""+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSnmpClientlist(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set snmp client-list \"" + d.Get("name").(string) + "\" "
	configSet := make([]string, 0)

	configSet = append(configSet, setPrefix)
	for _, v := range sortSetOfString(d.Get("prefix").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+v)
	}

	return sess.configSet(configSet, jnprSess)
}

func readSnmpClientlist(name string, m interface{}, jnprSess *NetconfObject) (snmpClientlistOptions, error) {
	sess := m.(*Session)
	var confRead snmpClientlistOptions

	showConfig, err := sess.command(cmdShowConfig+"snmp client-list \""+name+"\""+pipeDisplaySetRelative, jnprSess)
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
			if itemTrim != "" {
				confRead.prefix = append(confRead.prefix, itemTrim)
			}
		}
	}

	return confRead, nil
}

func delSnmpClientlist(name string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete snmp client-list \"" + name + "\""}

	return sess.configSet(configSet, jnprSess)
}

func fillSnmpClientlistData(d *schema.ResourceData, snmpClientlistOptions snmpClientlistOptions) {
	if tfErr := d.Set("name", snmpClientlistOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("prefix", snmpClientlistOptions.prefix); tfErr != nil {
		panic(tfErr)
	}
}
