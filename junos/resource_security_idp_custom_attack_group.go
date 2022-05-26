package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSecurityIdpCustomAttackGroup(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if !checkCompatibilitySecurity(junSess) {
		return diag.FromErr(fmt.Errorf("security idp custom-attack-group not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	idpCustomAttackGroupExists, err := checkSecurityIdpCustomAttackGroupExists(d.Get("name").(string), sess, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if idpCustomAttackGroupExists {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security idp custom-attack-group %v already exists", d.Get("name").(string)))...)
	}
	if err := setSecurityIdpCustomAttackGroup(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_idp_custom_attack_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	idpCustomAttackGroupExists, err = checkSecurityIdpCustomAttackGroupExists(d.Get("name").(string), sess, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if idpCustomAttackGroupExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security idp custom-attack-group %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityIdpCustomAttackGroupReadWJunSess(d, sess, junSess)...)
}

func resourceSecurityIdpCustomAttackGroupRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)

	return resourceSecurityIdpCustomAttackGroupReadWJunSess(d, sess, junSess)
}

func resourceSecurityIdpCustomAttackGroupReadWJunSess(d *schema.ResourceData, sess *Session, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	idpCustomAttackGroupOptions, err := readSecurityIdpCustomAttackGroup(d.Get("name").(string), sess, junSess)
	mutex.Unlock()
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
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delSecurityIdpCustomAttackGroup(d.Get("name").(string), sess, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityIdpCustomAttackGroup(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}

	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityIdpCustomAttackGroup(d.Get("name").(string), sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityIdpCustomAttackGroup(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_idp_custom_attack_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityIdpCustomAttackGroupReadWJunSess(d, sess, junSess)...)
}

func resourceSecurityIdpCustomAttackGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delSecurityIdpCustomAttackGroup(d.Get("name").(string), sess, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityIdpCustomAttackGroup(d.Get("name").(string), sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_idp_custom_attack_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityIdpCustomAttackGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	idpCustomAttackGroupExists, err := checkSecurityIdpCustomAttackGroupExists(d.Id(), sess, junSess)
	if err != nil {
		return nil, err
	}
	if !idpCustomAttackGroupExists {
		return nil, fmt.Errorf("don't find security idp custom-attack-group with id '%v' (id must be <name>)", d.Id())
	}
	idpCustomAttackGroupOptions, err := readSecurityIdpCustomAttackGroup(d.Id(), sess, junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityIdpCustomAttackGroupData(d, idpCustomAttackGroupOptions)

	result[0] = d

	return result, nil
}

func checkSecurityIdpCustomAttackGroupExists(customAttackGroup string, sess *Session, junSess *junosSession,
) (bool, error) {
	showConfig, err := sess.command(cmdShowConfig+
		"security idp custom-attack-group \""+customAttackGroup+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityIdpCustomAttackGroup(d *schema.ResourceData, sess *Session, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := "set security idp custom-attack-group \"" + d.Get("name").(string) + "\" "
	configSet = append(configSet, setPrefix)
	for _, v := range sortSetOfString(d.Get("member").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"group-members \""+v+"\"")
	}

	return sess.configSet(configSet, junSess)
}

func readSecurityIdpCustomAttackGroup(customAttackGroup string, sess *Session, junSess *junosSession,
) (idpCustomAttackGroupOptions, error) {
	var confRead idpCustomAttackGroupOptions

	showConfig, err := sess.command(cmdShowConfig+
		"security idp custom-attack-group \""+customAttackGroup+"\""+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = customAttackGroup
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			if strings.HasPrefix(itemTrim, "group-members ") {
				confRead.member = append(confRead.member, strings.Trim(strings.TrimPrefix(itemTrim, "group-members "), "\""))
			}
		}
	}

	return confRead, nil
}

func delSecurityIdpCustomAttackGroup(customAttack string, sess *Session, junSess *junosSession) error {
	configSet := []string{"delete security idp custom-attack-group \"" + customAttack + "\""}

	return sess.configSet(configSet, junSess)
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
