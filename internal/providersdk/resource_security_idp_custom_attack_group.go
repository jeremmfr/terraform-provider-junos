package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type idpCustomAttackGroupOptions struct {
	name   string
	member []string
}

func resourceSecurityIdpCustomAttackGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityIdpCustomAttackGroupCreate,
		ReadWithoutTimeout:   resourceSecurityIdpCustomAttackGroupRead,
		UpdateWithoutTimeout: resourceSecurityIdpCustomAttackGroupUpdate,
		DeleteWithoutTimeout: resourceSecurityIdpCustomAttackGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityIdpCustomAttackGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"member": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSecurityIdpCustomAttackGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityIdpCustomAttackGroup(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilitySecurity() {
		return diag.FromErr(fmt.Errorf("security idp custom-attack-group not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	idpCustomAttackGroupExists, err := checkSecurityIdpCustomAttackGroupExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if idpCustomAttackGroupExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security idp custom-attack-group %v already exists", d.Get("name").(string)))...)
	}
	if err := setSecurityIdpCustomAttackGroup(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_security_idp_custom_attack_group")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	idpCustomAttackGroupExists, err = checkSecurityIdpCustomAttackGroupExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if idpCustomAttackGroupExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security idp custom-attack-group %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityIdpCustomAttackGroupReadWJunSess(d, junSess)...)
}

func resourceSecurityIdpCustomAttackGroupRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityIdpCustomAttackGroupReadWJunSess(d, junSess)
}

func resourceSecurityIdpCustomAttackGroupReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	idpCustomAttackGroupOptions, err := readSecurityIdpCustomAttackGroup(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if idpCustomAttackGroupOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityIdpCustomAttackGroupData(d, idpCustomAttackGroupOptions)
	}

	return nil
}

func resourceSecurityIdpCustomAttackGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityIdpCustomAttackGroup(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityIdpCustomAttackGroup(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}

	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityIdpCustomAttackGroup(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityIdpCustomAttackGroup(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_security_idp_custom_attack_group")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityIdpCustomAttackGroupReadWJunSess(d, junSess)...)
}

func resourceSecurityIdpCustomAttackGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityIdpCustomAttackGroup(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityIdpCustomAttackGroup(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_security_idp_custom_attack_group")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityIdpCustomAttackGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	idpCustomAttackGroupExists, err := checkSecurityIdpCustomAttackGroupExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !idpCustomAttackGroupExists {
		return nil, fmt.Errorf("don't find security idp custom-attack-group with id '%v' (id must be <name>)", d.Id())
	}
	idpCustomAttackGroupOptions, err := readSecurityIdpCustomAttackGroup(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityIdpCustomAttackGroupData(d, idpCustomAttackGroupOptions)

	result[0] = d

	return result, nil
}

func checkSecurityIdpCustomAttackGroupExists(customAttackGroup string, junSess *junos.Session,
) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security idp custom-attack-group \"" + customAttackGroup + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityIdpCustomAttackGroup(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security idp custom-attack-group \"" + d.Get("name").(string) + "\" "
	configSet = append(configSet, setPrefix)
	for _, v := range sortSetOfString(d.Get("member").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"group-members \""+v+"\"")
	}

	return junSess.ConfigSet(configSet)
}

func readSecurityIdpCustomAttackGroup(customAttackGroup string, junSess *junos.Session,
) (confRead idpCustomAttackGroupOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security idp custom-attack-group \"" + customAttackGroup + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = customAttackGroup
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if balt.CutPrefixInString(&itemTrim, "group-members ") {
				confRead.member = append(confRead.member, strings.Trim(itemTrim, "\""))
			}
		}
	}

	return confRead, nil
}

func delSecurityIdpCustomAttackGroup(customAttack string, junSess *junos.Session) error {
	configSet := []string{"delete security idp custom-attack-group \"" + customAttack + "\""}

	return junSess.ConfigSet(configSet)
}

func fillSecurityIdpCustomAttackGroupData(
	d *schema.ResourceData, idpCustomAttackGroupOptions idpCustomAttackGroupOptions,
) {
	if tfErr := d.Set("name", idpCustomAttackGroupOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("member", idpCustomAttackGroupOptions.member); tfErr != nil {
		panic(tfErr)
	}
}
