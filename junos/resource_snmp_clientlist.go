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
		CreateWithoutTimeout: resourceSnmpClientlistCreate,
		ReadWithoutTimeout:   resourceSnmpClientlistRead,
		UpdateWithoutTimeout: resourceSnmpClientlistUpdate,
		DeleteWithoutTimeout: resourceSnmpClientlistDelete,
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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSnmpClientlist(d, clt, nil); err != nil {
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
	snmpClientlistExists, err := checkSnmpClientlistExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpClientlistExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp client-list %v already exists", d.Get("name").(string)))...)
	}

	if err := setSnmpClientlist(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_snmp_clientlist", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	snmpClientlistExists, err = checkSnmpClientlistExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpClientlistExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp client-list %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSnmpClientlistReadWJunSess(d, clt, junSess)...)
}

func resourceSnmpClientlistRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSnmpClientlistReadWJunSess(d, clt, junSess)
}

func resourceSnmpClientlistReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	snmpClientlistOptions, err := readSnmpClientlist(d.Get("name").(string), clt, junSess)
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
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSnmpClientlist(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSnmpClientlist(d, clt, nil); err != nil {
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
	if err := delSnmpClientlist(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmpClientlist(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_snmp_clientlist", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpClientlistReadWJunSess(d, clt, junSess)...)
}

func resourceSnmpClientlistDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delSnmpClientlist(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delSnmpClientlist(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_snmp_clientlist", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSnmpClientlistImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)

	snmpClientlistExists, err := checkSnmpClientlistExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !snmpClientlistExists {
		return nil, fmt.Errorf("don't find snmp client-list with id '%v' (id must be <name>)", d.Id())
	}
	snmpClientlistOptions, err := readSnmpClientlist(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSnmpClientlistData(d, snmpClientlistOptions)

	result[0] = d

	return result, nil
}

func checkSnmpClientlistExists(name string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+"snmp client-list \""+name+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSnmpClientlist(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	setPrefix := "set snmp client-list \"" + d.Get("name").(string) + "\" "
	configSet := make([]string, 0)

	configSet = append(configSet, setPrefix)
	for _, v := range sortSetOfString(d.Get("prefix").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+v)
	}

	return clt.configSet(configSet, junSess)
}

func readSnmpClientlist(name string, clt *Client, junSess *junosSession) (snmpClientlistOptions, error) {
	var confRead snmpClientlistOptions

	showConfig, err := clt.command(cmdShowConfig+"snmp client-list \""+name+"\""+pipeDisplaySetRelative, junSess)
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

func delSnmpClientlist(name string, clt *Client, junSess *junosSession) error {
	configSet := []string{"delete snmp client-list \"" + name + "\""}

	return clt.configSet(configSet, junSess)
}

func fillSnmpClientlistData(d *schema.ResourceData, snmpClientlistOptions snmpClientlistOptions) {
	if tfErr := d.Set("name", snmpClientlistOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("prefix", snmpClientlistOptions.prefix); tfErr != nil {
		panic(tfErr)
	}
}
