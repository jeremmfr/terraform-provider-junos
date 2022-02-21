package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	jdecode "github.com/jeremmfr/junosdecode"
)

type snmpV3CommunityOptions struct {
	communityIndex string
	communityName  string
	securityName   string
	context        string
	tag            string
}

func resourceSnmpV3Community() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnmpV3CommunityCreate,
		ReadContext:   resourceSnmpV3CommunityRead,
		UpdateContext: resourceSnmpV3CommunityUpdate,
		DeleteContext: resourceSnmpV3CommunityDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSnmpV3CommunityImport,
		},
		Schema: map[string]*schema.Schema{
			"community_index": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"community_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"context": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSnmpV3CommunityCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSnmpV3Community(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("community_index").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	snmpV3CommunityExists, err := checkSnmpV3CommunityExists(d.Get("community_index").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3CommunityExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"snmp v3 snmp-community %v already exists", d.Get("community_index").(string)))...)
	}

	if err := setSnmpV3Community(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_snmp_v3_community", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	snmpV3CommunityExists, err = checkSnmpV3CommunityExists(d.Get("community_index").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3CommunityExists {
		d.SetId(d.Get("community_index").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp v3 snmp-community %v not exists after commit "+
			"=> check your config", d.Get("community_index").(string)))...)
	}

	return append(diagWarns, resourceSnmpV3CommunityReadWJnprSess(d, m, jnprSess)...)
}

func resourceSnmpV3CommunityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSnmpV3CommunityReadWJnprSess(d, m, jnprSess)
}

func resourceSnmpV3CommunityReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	snmpV3CommunityOptions, err := readSnmpV3Community(d.Get("community_index").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if snmpV3CommunityOptions.communityIndex == "" {
		d.SetId("")
	} else {
		fillSnmpV3CommunityData(d, snmpV3CommunityOptions)
	}

	return nil
}

func resourceSnmpV3CommunityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delSnmpV3Community(d.Get("community_index").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSnmpV3Community(d, m, nil); err != nil {
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
	if err := delSnmpV3Community(d.Get("community_index").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmpV3Community(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_snmp_v3_community", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpV3CommunityReadWJnprSess(d, m, jnprSess)...)
}

func resourceSnmpV3CommunityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delSnmpV3Community(d.Get("community_index").(string), m, nil); err != nil {
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
	if err := delSnmpV3Community(d.Get("community_index").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_snmp_v3_community", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSnmpV3CommunityImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	snmpV3CommunityExists, err := checkSnmpV3CommunityExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !snmpV3CommunityExists {
		return nil, fmt.Errorf("don't find snmp v3 snmp-community with id '%v' (id must be <community_index>)", d.Id())
	}
	snmpV3CommunityOptions, err := readSnmpV3Community(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSnmpV3CommunityData(d, snmpV3CommunityOptions)

	result[0] = d

	return result, nil
}

func checkSnmpV3CommunityExists(communityIndex string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command(
		"show configuration snmp v3 snmp-community \""+communityIndex+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setSnmpV3Community(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set snmp v3 snmp-community \"" + d.Get("community_index").(string) + "\" "
	configSet := make([]string, 0)

	configSet = append(configSet, setPrefix+"security-name \""+d.Get("security_name").(string)+"\"")
	if v := d.Get("community_name").(string); v != "" {
		configSet = append(configSet, setPrefix+"community-name \""+v+"\"")
	}
	if v := d.Get("context").(string); v != "" {
		configSet = append(configSet, setPrefix+"context \""+v+"\"")
	}
	if v := d.Get("tag").(string); v != "" {
		configSet = append(configSet, setPrefix+"tag \""+v+"\"")
	}

	return sess.configSet(configSet, jnprSess)
}

func readSnmpV3Community(communityIndex string, m interface{}, jnprSess *NetconfObject,
) (snmpV3CommunityOptions, error) {
	sess := m.(*Session)
	var confRead snmpV3CommunityOptions

	showConfig, err := sess.command(
		"show configuration snmp v3 snmp-community \""+communityIndex+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		confRead.communityIndex = communityIndex
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "security-name "):
				confRead.securityName = strings.Trim(strings.TrimPrefix(itemTrim, "security-name "), "\"")
			case strings.HasPrefix(itemTrim, "community-name "):
				var err error
				confRead.communityName, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim, "community-name "), "\""))
				if err != nil {
					return confRead, fmt.Errorf("failed to decode community-name : %w", err)
				}
			case strings.HasPrefix(itemTrim, "context "):
				confRead.context = strings.Trim(strings.TrimPrefix(itemTrim, "context "), "\"")
			case strings.HasPrefix(itemTrim, "tag "):
				confRead.tag = strings.Trim(strings.TrimPrefix(itemTrim, "tag "), "\"")
			}
		}
	}

	return confRead, nil
}

func delSnmpV3Community(communityIndex string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	configSet := []string{"delete snmp v3 snmp-community \"" + communityIndex + "\""}

	return sess.configSet(configSet, jnprSess)
}

func fillSnmpV3CommunityData(d *schema.ResourceData, snmpV3CommunityOptions snmpV3CommunityOptions) {
	if tfErr := d.Set("community_index", snmpV3CommunityOptions.communityIndex); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("security_name", snmpV3CommunityOptions.securityName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("community_name", snmpV3CommunityOptions.communityName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("context", snmpV3CommunityOptions.context); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("tag", snmpV3CommunityOptions.tag); tfErr != nil {
		panic(tfErr)
	}
}
