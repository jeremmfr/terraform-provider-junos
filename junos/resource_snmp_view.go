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
		CreateWithoutTimeout: resourceSnmpViewCreate,
		ReadWithoutTimeout:   resourceSnmpViewRead,
		UpdateWithoutTimeout: resourceSnmpViewUpdate,
		DeleteWithoutTimeout: resourceSnmpViewDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSnmpViewImport,
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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSnmpView(d, clt, nil); err != nil {
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
	snmpViewExists, err := checkSnmpViewExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpViewExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp view %v already exists", d.Get("name").(string)))...)
	}

	if err := setSnmpView(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_snmp_view", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	snmpViewExists, err = checkSnmpViewExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpViewExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp view %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSnmpViewReadWJunSess(d, clt, junSess)...)
}

func resourceSnmpViewRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSnmpViewReadWJunSess(d, clt, junSess)
}

func resourceSnmpViewReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	snmpViewOptions, err := readSnmpView(d.Get("name").(string), clt, junSess)
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
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSnmpView(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSnmpView(d, clt, nil); err != nil {
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
	if err := delSnmpView(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmpView(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_snmp_view", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpViewReadWJunSess(d, clt, junSess)...)
}

func resourceSnmpViewDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delSnmpView(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delSnmpView(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_snmp_view", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSnmpViewImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)

	snmpViewExists, err := checkSnmpViewExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !snmpViewExists {
		return nil, fmt.Errorf("don't find snmp view with id '%v' (id must be <name>)", d.Id())
	}
	snmpViewOptions, err := readSnmpView(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSnmpViewData(d, snmpViewOptions)

	result[0] = d

	return result, nil
}

func checkSnmpViewExists(name string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+"snmp view \""+name+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSnmpView(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	setPrefix := "set snmp view \"" + d.Get("name").(string) + "\" "
	configSet := make([]string, 0)

	for _, v := range sortSetOfString(d.Get("oid_include").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"oid "+v+" include")
	}
	for _, v := range sortSetOfString(d.Get("oid_exclude").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"oid "+v+" exclude")
	}

	return clt.configSet(configSet, junSess)
}

func readSnmpView(name string, clt *Client, junSess *junosSession) (snmpViewOptions, error) {
	var confRead snmpViewOptions

	showConfig, err := clt.command(cmdShowConfig+"snmp view \""+name+"\""+pipeDisplaySetRelative, junSess)
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

func delSnmpView(name string, clt *Client, junSess *junosSession) error {
	configSet := []string{"delete snmp view \"" + name + "\""}

	return clt.configSet(configSet, junSess)
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
