package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type snmpCommunityOptions struct {
	authReadOnly    bool
	authReadWrite   bool
	name            string
	clientListName  string
	view            string
	clients         []string
	routingInstance []map[string]interface{}
}

func resourceSnmpCommunity() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSnmpCommunityCreate,
		ReadWithoutTimeout:   resourceSnmpCommunityRead,
		UpdateWithoutTimeout: resourceSnmpCommunityUpdate,
		DeleteWithoutTimeout: resourceSnmpCommunityDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSnmpCommunityImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"authorization_read_only": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"authorization_read_write"},
			},
			"authorization_read_write": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"authorization_read_only"},
			},
			"client_list_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"clients"},
			},
			"clients": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"client_list_name"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"routing_instance": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64, formatDefault),
						},
						"client_list_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"clients": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"view": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSnmpCommunityCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSnmpCommunity(d, clt, nil); err != nil {
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
	snmpCommunityExists, err := checkSnmpCommunityExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpCommunityExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp community %v already exists", d.Get("name").(string)))...)
	}

	if err := setSnmpCommunity(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_snmp_community", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	snmpCommunityExists, err = checkSnmpCommunityExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpCommunityExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp community %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSnmpCommunityReadWJunSess(d, clt, junSess)...)
}

func resourceSnmpCommunityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSnmpCommunityReadWJunSess(d, clt, junSess)
}

func resourceSnmpCommunityReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	snmpCommunityOptions, err := readSnmpCommunity(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if snmpCommunityOptions.name == "" {
		d.SetId("")
	} else {
		fillSnmpCommunityData(d, snmpCommunityOptions)
	}

	return nil
}

func resourceSnmpCommunityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSnmpCommunity(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSnmpCommunity(d, clt, nil); err != nil {
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
	if err := delSnmpCommunity(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmpCommunity(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_snmp_community", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpCommunityReadWJunSess(d, clt, junSess)...)
}

func resourceSnmpCommunityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delSnmpCommunity(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delSnmpCommunity(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_snmp_community", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSnmpCommunityImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)

	snmpCommunityExists, err := checkSnmpCommunityExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !snmpCommunityExists {
		return nil, fmt.Errorf("don't find snmp community with id '%v' (id must be <name>)", d.Id())
	}
	snmpCommunityOptions, err := readSnmpCommunity(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSnmpCommunityData(d, snmpCommunityOptions)

	result[0] = d

	return result, nil
}

func checkSnmpCommunityExists(name string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+"snmp community \""+name+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSnmpCommunity(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	setPrefix := "set snmp community \"" + d.Get("name").(string) + "\" "
	configSet := make([]string, 0)

	if d.Get("authorization_read_only").(bool) {
		configSet = append(configSet, setPrefix+"authorization read-only")
	}
	if d.Get("authorization_read_write").(bool) {
		configSet = append(configSet, setPrefix+"authorization read-write")
	}
	if v := d.Get("client_list_name").(string); v != "" {
		configSet = append(configSet, setPrefix+"client-list-name \""+v+"\"")
	}
	for _, v := range sortSetOfString(d.Get("clients").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"clients "+v)
	}
	routingInstanceNameList := make([]string, 0)
	for _, v := range d.Get("routing_instance").([]interface{}) {
		routingInstance := v.(map[string]interface{})
		if len(routingInstance["clients"].(*schema.Set).List()) > 0 && routingInstance["client_list_name"].(string) != "" {
			return fmt.Errorf("conflict between clients and client_list_name in routing-instance %s",
				routingInstance["name"].(string))
		}
		if bchk.StringInSlice(routingInstance["name"].(string), routingInstanceNameList) {
			return fmt.Errorf("multiple blocks routing_instance with the same name %s", routingInstance["name"].(string))
		}
		routingInstanceNameList = append(routingInstanceNameList, routingInstance["name"].(string))
		configSet = append(configSet, setPrefix+"routing-instance "+routingInstance["name"].(string))
		if cLNname := routingInstance["client_list_name"].(string); cLNname != "" {
			configSet = append(configSet,
				setPrefix+"routing-instance "+routingInstance["name"].(string)+" client-list-name \""+cLNname+"\"")
		}
		for _, clt := range sortSetOfString(routingInstance["clients"].(*schema.Set).List()) {
			configSet = append(configSet,
				setPrefix+"routing-instance "+routingInstance["name"].(string)+" clients "+clt)
		}
	}
	if v := d.Get("view").(string); v != "" {
		configSet = append(configSet, setPrefix+"view \""+v+"\"")
	}

	return clt.configSet(configSet, junSess)
}

func readSnmpCommunity(name string, clt *Client, junSess *junosSession) (snmpCommunityOptions, error) {
	var confRead snmpCommunityOptions

	showConfig, err := clt.command(cmdShowConfig+"snmp community \""+name+"\""+pipeDisplaySetRelative, junSess)
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
			case itemTrim == "authorization read-only":
				confRead.authReadOnly = true
			case itemTrim == "authorization read-write":
				confRead.authReadWrite = true
			case strings.HasPrefix(itemTrim, "client-list-name "):
				confRead.clientListName = strings.Trim(strings.TrimPrefix(itemTrim, "client-list-name "), "\"")
			case strings.HasPrefix(itemTrim, "clients "):
				confRead.clients = append(confRead.clients, strings.TrimPrefix(itemTrim, "clients "))
			case strings.HasPrefix(itemTrim, "routing-instance "):
				routingInstanceLineCut := strings.Split(itemTrim, " ")
				mRoutingInstance := map[string]interface{}{
					"name":             routingInstanceLineCut[1],
					"client_list_name": "",
					"clients":          make([]string, 0),
				}
				confRead.routingInstance = copyAndRemoveItemMapList("name", mRoutingInstance, confRead.routingInstance)
				itemTrimRoutingInstance := strings.TrimPrefix(itemTrim, "routing-instance "+routingInstanceLineCut[1]+" ")
				switch {
				case strings.HasPrefix(itemTrimRoutingInstance, "client-list-name "):
					mRoutingInstance["client_list_name"] = strings.Trim(strings.TrimPrefix(
						itemTrimRoutingInstance, "client-list-name "), "\"")
				case strings.HasPrefix(itemTrimRoutingInstance, "clients "):
					mRoutingInstance["clients"] = append(mRoutingInstance["clients"].([]string),
						strings.TrimPrefix(itemTrimRoutingInstance, "clients "))
				}
				confRead.routingInstance = append(confRead.routingInstance, mRoutingInstance)
			case strings.HasPrefix(itemTrim, "view "):
				confRead.view = strings.Trim(strings.TrimPrefix(itemTrim, "view "), "\"")
			}
		}
	}

	return confRead, nil
}

func delSnmpCommunity(name string, clt *Client, junSess *junosSession) error {
	configSet := []string{"delete snmp community \"" + name + "\""}

	return clt.configSet(configSet, junSess)
}

func fillSnmpCommunityData(d *schema.ResourceData, snmpCommunityOptions snmpCommunityOptions) {
	if tfErr := d.Set("name", snmpCommunityOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authorization_read_only", snmpCommunityOptions.authReadOnly); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authorization_read_write", snmpCommunityOptions.authReadWrite); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("client_list_name", snmpCommunityOptions.clientListName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("clients", snmpCommunityOptions.clients); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", snmpCommunityOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("view", snmpCommunityOptions.view); tfErr != nil {
		panic(tfErr)
	}
}
