package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
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
		CreateWithoutTimeout: resourceSnmpV3CommunityCreate,
		ReadWithoutTimeout:   resourceSnmpV3CommunityRead,
		UpdateWithoutTimeout: resourceSnmpV3CommunityUpdate,
		DeleteWithoutTimeout: resourceSnmpV3CommunityDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSnmpV3CommunityImport,
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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSnmpV3Community(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("community_index").(string))

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
	snmpV3CommunityExists, err := checkSnmpV3CommunityExists(d.Get("community_index").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3CommunityExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"snmp v3 snmp-community %v already exists", d.Get("community_index").(string)))...)
	}

	if err := setSnmpV3Community(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_snmp_v3_community", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	snmpV3CommunityExists, err = checkSnmpV3CommunityExists(d.Get("community_index").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3CommunityExists {
		d.SetId(d.Get("community_index").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp v3 snmp-community %v not exists after commit "+
			"=> check your config", d.Get("community_index").(string)))...)
	}

	return append(diagWarns, resourceSnmpV3CommunityReadWJunSess(d, clt, junSess)...)
}

func resourceSnmpV3CommunityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSnmpV3CommunityReadWJunSess(d, clt, junSess)
}

func resourceSnmpV3CommunityReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	snmpV3CommunityOptions, err := readSnmpV3Community(d.Get("community_index").(string), clt, junSess)
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
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSnmpV3Community(d.Get("community_index").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSnmpV3Community(d, clt, nil); err != nil {
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
	if err := delSnmpV3Community(d.Get("community_index").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmpV3Community(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_snmp_v3_community", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpV3CommunityReadWJunSess(d, clt, junSess)...)
}

func resourceSnmpV3CommunityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delSnmpV3Community(d.Get("community_index").(string), clt, nil); err != nil {
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
	if err := delSnmpV3Community(d.Get("community_index").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_snmp_v3_community", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSnmpV3CommunityImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)

	snmpV3CommunityExists, err := checkSnmpV3CommunityExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !snmpV3CommunityExists {
		return nil, fmt.Errorf("don't find snmp v3 snmp-community with id '%v' (id must be <community_index>)", d.Id())
	}
	snmpV3CommunityOptions, err := readSnmpV3Community(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSnmpV3CommunityData(d, snmpV3CommunityOptions)

	result[0] = d

	return result, nil
}

func checkSnmpV3CommunityExists(communityIndex string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+
		"snmp v3 snmp-community \""+communityIndex+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSnmpV3Community(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
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

	return clt.configSet(configSet, junSess)
}

func readSnmpV3Community(communityIndex string, clt *Client, junSess *junosSession,
) (confRead snmpV3CommunityOptions, err error) {
	showConfig, err := clt.command(cmdShowConfig+
		"snmp v3 snmp-community \""+communityIndex+"\""+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.communityIndex = communityIndex
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "security-name "):
				confRead.securityName = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "community-name "):
				confRead.communityName, err = jdecode.Decode(strings.Trim(itemTrim, "\""))
				if err != nil {
					return confRead, fmt.Errorf("failed to decode community-name: %w", err)
				}
			case balt.CutPrefixInString(&itemTrim, "context "):
				confRead.context = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "tag "):
				confRead.tag = strings.Trim(itemTrim, "\"")
			}
		}
	}

	return confRead, nil
}

func delSnmpV3Community(communityIndex string, clt *Client, junSess *junosSession) error {
	configSet := []string{"delete snmp v3 snmp-community \"" + communityIndex + "\""}

	return clt.configSet(configSet, junSess)
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
