package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type ipsecPolicyOptions struct {
	name        string
	pfsKeys     string
	proposalSet string
	proposals   []string
}

func listProposalSet() []string {
	return []string{
		"basic",
		"compatible",
		"prime-128",
		"prime-256",
		"standard",
		"suiteb-gcm-128",
		"suiteb-gcm-256",
	}
}

func resourceIpsecPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceIpsecPolicyCreate,
		ReadWithoutTimeout:   resourceIpsecPolicyRead,
		UpdateWithoutTimeout: resourceIpsecPolicyUpdate,
		DeleteWithoutTimeout: resourceIpsecPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIpsecPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
			},
			"pfs_keys": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proposals": {
				Type:         schema.TypeList,
				Optional:     true,
				MinItems:     1,
				Elem:         &schema.Schema{Type: schema.TypeString},
				ExactlyOneOf: []string{"proposals", "proposal_set"},
			},
			"proposal_set": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listProposalSet(), false),
				ExactlyOneOf: []string{"proposals", "proposal_set"},
			},
		},
	}
}

func resourceIpsecPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setIpsecPolicy(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security ipsec policy not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ipsecPolicyExists, err := checkIpsecPolicyExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecPolicyExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security ipsec policy %v already exists", d.Get("name").(string)))...)
	}
	if err := setIpsecPolicy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_ipsec_policy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	ipsecPolicyExists, err = checkIpsecPolicyExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecPolicyExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ipsec policy %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIpsecPolicyReadWJunSess(d, junSess)...)
}

func resourceIpsecPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceIpsecPolicyReadWJunSess(d, junSess)
}

func resourceIpsecPolicyReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	ipsecPolicyOptions, err := readIpsecPolicy(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ipsecPolicyOptions.name == "" {
		d.SetId("")
	} else {
		fillIpsecPolicyData(d, ipsecPolicyOptions)
	}

	return nil
}

func resourceIpsecPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delIpsecPolicy(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setIpsecPolicy(d, junSess); err != nil {
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
	if err := delIpsecPolicy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setIpsecPolicy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_ipsec_policy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIpsecPolicyReadWJunSess(d, junSess)...)
}

func resourceIpsecPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delIpsecPolicy(d, junSess); err != nil {
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
	if err := delIpsecPolicy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_ipsec_policy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceIpsecPolicyImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	ipsecPolicyExists, err := checkIpsecPolicyExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !ipsecPolicyExists {
		return nil, fmt.Errorf("don't find security ipsec policy with id '%v' (id must be <name>)", d.Id())
	}
	ipsecPolicyOptions, err := readIpsecPolicy(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillIpsecPolicyData(d, ipsecPolicyOptions)
	result[0] = d

	return result, nil
}

func checkIpsecPolicyExists(ipsecPolicy string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ipsec policy " + ipsecPolicy + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setIpsecPolicy(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security ipsec policy " + d.Get("name").(string)
	if d.Get("pfs_keys").(string) != "" {
		configSet = append(configSet, setPrefix+" perfect-forward-secrecy keys "+d.Get("pfs_keys").(string))
	}
	for _, v := range d.Get("proposals").([]interface{}) {
		configSet = append(configSet, setPrefix+" proposals "+v.(string))
	}
	if d.Get("proposal_set").(string) != "" {
		configSet = append(configSet, setPrefix+" proposal-set "+d.Get("proposal_set").(string))
	}

	return junSess.ConfigSet(configSet)
}

func readIpsecPolicy(ipsecPolicy string, junSess *junos.Session,
) (confRead ipsecPolicyOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ipsec policy " + ipsecPolicy + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = ipsecPolicy
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "perfect-forward-secrecy keys "):
				confRead.pfsKeys = itemTrim
			case balt.CutPrefixInString(&itemTrim, "proposals "):
				confRead.proposals = append(confRead.proposals, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "proposal-set "):
				confRead.proposalSet = itemTrim
			}
		}
	}

	return confRead, nil
}

func delIpsecPolicy(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec policy "+d.Get("name").(string))

	return junSess.ConfigSet(configSet)
}

func fillIpsecPolicyData(d *schema.ResourceData, ipsecPolicyOptions ipsecPolicyOptions) {
	if tfErr := d.Set("name", ipsecPolicyOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("pfs_keys", ipsecPolicyOptions.pfsKeys); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("proposals", ipsecPolicyOptions.proposals); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("proposal_set", ipsecPolicyOptions.proposalSet); tfErr != nil {
		panic(tfErr)
	}
}
