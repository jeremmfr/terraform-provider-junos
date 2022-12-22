package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type snmpV3VacmSecurityToGroupOptions struct {
	name  string
	model string
	group string
}

func resourceSnmpV3VacmSecurityToGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSnmpV3VacmSecurityToGroupCreate,
		ReadWithoutTimeout:   resourceSnmpV3VacmSecurityToGroupRead,
		UpdateWithoutTimeout: resourceSnmpV3VacmSecurityToGroupUpdate,
		DeleteWithoutTimeout: resourceSnmpV3VacmSecurityToGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSnmpV3VacmSecurityToGroupImport,
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

func resourceSnmpV3VacmSecurityToGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSnmpV3VacmSecurityToGroup(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("model").(string) + idSeparator + d.Get("name").(string))

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
	snmpV3VacmSecurityToGroupExists, err := checkSnmpV3VacmSecurityToGroupExists(
		d.Get("model").(string),
		d.Get("name").(string),
		clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3VacmSecurityToGroupExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"snmp v3 vacm security-to-group security-model %v security-name %v already exists",
			d.Get("model").(string), d.Get("name").(string)))...)
	}

	if err := setSnmpV3VacmSecurityToGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_snmp_v3_vacm_securitytogroup", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	snmpV3VacmSecurityToGroupExists, err = checkSnmpV3VacmSecurityToGroupExists(
		d.Get("model").(string),
		d.Get("name").(string),
		clt, junSess)
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

	return append(diagWarns, resourceSnmpV3VacmSecurityToGroupReadWJunSess(d, clt, junSess)...)
}

func resourceSnmpV3VacmSecurityToGroupRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSnmpV3VacmSecurityToGroupReadWJunSess(d, clt, junSess)
}

func resourceSnmpV3VacmSecurityToGroupReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	snmpV3VacmSecurityToGroupOptions, err := readSnmpV3VacmSecurityToGroup(
		d.Get("model").(string),
		d.Get("name").(string),
		clt, junSess)
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

func resourceSnmpV3VacmSecurityToGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSnmpV3VacmSecurityToGroup(d.Get("model").(string), d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSnmpV3VacmSecurityToGroup(d, clt, nil); err != nil {
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
	if err := delSnmpV3VacmSecurityToGroup(d.Get("model").(string), d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmpV3VacmSecurityToGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_snmp_v3_vacm_securitytogroup", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpV3VacmSecurityToGroupReadWJunSess(d, clt, junSess)...)
}

func resourceSnmpV3VacmSecurityToGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delSnmpV3VacmSecurityToGroup(d.Get("model").(string), d.Get("name").(string), clt, nil); err != nil {
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
	if err := delSnmpV3VacmSecurityToGroup(d.Get("model").(string), d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_snmp_v3_vacm_securitytogroup", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSnmpV3VacmSecurityToGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) != 2 {
		return nil, fmt.Errorf("can't find snmp v3 vacm security-to-group "+
			"with id '%v' (id must be <model>%s<name>)", d.Id(), idSeparator)
	}
	if !bchk.InSlice(idSplit[0], []string{"usm", "v1", "v2c"}) {
		return nil, fmt.Errorf("can't find snmp v3 vacm security-to-group "+
			"with id '%v' (id must be <model>%s<name>)", d.Id(), idSeparator)
	}
	snmpV3VacmSecurityToGroupExists, err := checkSnmpV3VacmSecurityToGroupExists(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	if !snmpV3VacmSecurityToGroupExists {
		return nil, fmt.Errorf("don't find snmp v3 vacm security-to-group "+
			"with id '%v' (id must be <model>%s<name>)", d.Id(), idSeparator)
	}
	snmpV3VacmSecurityToGroupOptions, err := readSnmpV3VacmSecurityToGroup(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSnmpV3VacmSecurityToGroupData(d, snmpV3VacmSecurityToGroupOptions)

	result[0] = d

	return result, nil
}

func checkSnmpV3VacmSecurityToGroupExists(model, name string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+"snmp v3 vacm security-to-group "+
		"security-model "+model+" security-name \""+name+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSnmpV3VacmSecurityToGroup(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 1)
	if group := d.Get("group").(string); group != "" {
		configSet[0] = "set snmp v3 vacm security-to-group" +
			" security-model " + d.Get("model").(string) +
			" security-name \"" + d.Get("name").(string) + "\"" +
			" group \"" + d.Get("group").(string) + "\""
	} else {
		return fmt.Errorf("group need to be set")
	}

	return clt.configSet(configSet, junSess)
}

func readSnmpV3VacmSecurityToGroup(model, name string, clt *Client, junSess *junosSession,
) (confRead snmpV3VacmSecurityToGroupOptions, err error) {
	showConfig, err := clt.command(cmdShowConfig+"snmp v3 vacm security-to-group "+
		"security-model "+model+" security-name \""+name+"\""+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.model = model
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			if balt.CutPrefixInString(&itemTrim, "group ") {
				confRead.group = strings.Trim(itemTrim, "\"")
			}
		}
	}

	return confRead, nil
}

func delSnmpV3VacmSecurityToGroup(model, name string, clt *Client, junSess *junosSession) error {
	configSet := []string{"delete snmp v3 vacm security-to-group " +
		"security-model " + model + " security-name \"" + name + "\""}

	return clt.configSet(configSet, junSess)
}

func fillSnmpV3VacmSecurityToGroupData(
	d *schema.ResourceData, snmpV3VacmSecurityToGroupOptions snmpV3VacmSecurityToGroupOptions,
) {
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
