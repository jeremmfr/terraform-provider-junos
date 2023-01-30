package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	jdecode "github.com/jeremmfr/junosdecode"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type ikePolicyOptions struct {
	name             string
	mode             string
	preSharedKeyText string
	preSharedKeyHexa string
	proposalSet      string
	proposals        []string
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

func resourceIkePolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceIkePolicyCreate,
		ReadWithoutTimeout:   resourceIkePolicyRead,
		UpdateWithoutTimeout: resourceIkePolicyUpdate,
		DeleteWithoutTimeout: resourceIkePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIkePolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
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
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "main",
				ValidateFunc: validation.StringInSlice([]string{"main", "aggressive"}, false),
			},
			"pre_shared_key_text": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"pre_shared_key_hexa"},
				Sensitive:     true,
			},
			"pre_shared_key_hexa": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"pre_shared_key_text"},
				Sensitive:     true,
			},
		},
	}
}

func resourceIkePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setIkePolicy(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security ike policy not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ikePolicyExists, err := checkIkePolicyExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ikePolicyExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("security ike policy %v already exists", d.Get("name").(string)))...)
	}
	if err := setIkePolicy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_ike_policy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	ikePolicyExists, err = checkIkePolicyExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ikePolicyExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ike policy %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIkePolicyReadWJunSess(d, junSess)...)
}

func resourceIkePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceIkePolicyReadWJunSess(d, junSess)
}

func resourceIkePolicyReadWJunSess(d *schema.ResourceData, junSess *junos.Session) diag.Diagnostics {
	junos.MutexLock()
	ikePolicyOptions, err := readIkePolicy(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ikePolicyOptions.name == "" {
		d.SetId("")
	} else {
		fillIkePolicyData(d, ikePolicyOptions)
	}

	return nil
}

func resourceIkePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delIkePolicy(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setIkePolicy(d, junSess); err != nil {
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
	if err := delIkePolicy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setIkePolicy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_ike_policy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	d.Partial(false)

	return append(diagWarns, resourceIkePolicyReadWJunSess(d, junSess)...)
}

func resourceIkePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delIkePolicy(d, junSess); err != nil {
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
	if err := delIkePolicy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_ike_policy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceIkePolicyImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	ikePolicyExists, err := checkIkePolicyExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !ikePolicyExists {
		return nil, fmt.Errorf("don't find security ike policy with id '%v' (id must be <name>)", d.Id())
	}
	ikePolicyOptions, err := readIkePolicy(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillIkePolicyData(d, ikePolicyOptions)
	result[0] = d

	return result, nil
}

func checkIkePolicyExists(ikePolicy string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "security ike policy " + ikePolicy + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setIkePolicy(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security ike policy " + d.Get("name").(string)
	if d.Get("mode").(string) != "" {
		if d.Get("mode").(string) != "main" && d.Get("mode").(string) != "aggressive" {
			return fmt.Errorf("unknown ike mode %v", d.Get("mode").(string))
		}
		configSet = append(configSet, setPrefix+" mode "+d.Get("mode").(string))
	}
	for _, v := range d.Get("proposals").([]interface{}) {
		configSet = append(configSet, setPrefix+" proposals "+v.(string))
	}
	if d.Get("proposal_set").(string) != "" {
		configSet = append(configSet, setPrefix+" proposal-set "+d.Get("proposal_set").(string))
	}
	if d.Get("pre_shared_key_text").(string) != "" {
		configSet = append(configSet, setPrefix+" pre-shared-key ascii-text "+d.Get("pre_shared_key_text").(string))
	}
	if d.Get("pre_shared_key_hexa").(string) != "" {
		configSet = append(configSet, setPrefix+" pre-shared-key hexadecimal "+d.Get("pre_shared_key_hexa").(string))
	}

	return junSess.ConfigSet(configSet)
}

func readIkePolicy(ikePolicy string, junSess *junos.Session,
) (confRead ikePolicyOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ike policy " + ikePolicy + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = ikePolicy
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "mode "):
				confRead.mode = itemTrim
			case balt.CutPrefixInString(&itemTrim, "proposals "):
				confRead.proposals = append(confRead.proposals, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "proposal-set "):
				confRead.proposalSet = itemTrim
			case balt.CutPrefixInString(&itemTrim, "pre-shared-key hexadecimal "):
				confRead.preSharedKeyHexa, err = jdecode.Decode(strings.Trim(itemTrim, "\""))
				if err != nil {
					return confRead, fmt.Errorf("decoding pre-shared-key hexadecimal: %w", err)
				}
			case balt.CutPrefixInString(&itemTrim, "pre-shared-key ascii-text "):
				confRead.preSharedKeyText, err = jdecode.Decode(strings.Trim(itemTrim, "\""))
				if err != nil {
					return confRead, fmt.Errorf("decoding pre-shared-key ascii-text: %w", err)
				}
			}
		}
	}

	return confRead, nil
}

func delIkePolicy(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ike policy "+d.Get("name").(string))

	return junSess.ConfigSet(configSet)
}

func fillIkePolicyData(d *schema.ResourceData, ikePolicyOptions ikePolicyOptions) {
	if tfErr := d.Set("name", ikePolicyOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("mode", ikePolicyOptions.mode); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("pre_shared_key_text", ikePolicyOptions.preSharedKeyText); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("pre_shared_key_hexa", ikePolicyOptions.preSharedKeyHexa); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("proposals", ikePolicyOptions.proposals); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("proposal_set", ikePolicyOptions.proposalSet); tfErr != nil {
		panic(tfErr)
	}
}
