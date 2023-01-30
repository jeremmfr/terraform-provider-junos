package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type ikeProposalOptions struct {
	lifetimeSeconds         int
	name                    string
	authenticationAlgorithm string
	authenticationMethod    string
	dhGroup                 string
	encryptionAlgorithm     string
}

func resourceIkeProposal() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceIkeProposalCreate,
		ReadWithoutTimeout:   resourceIkeProposalRead,
		UpdateWithoutTimeout: resourceIkeProposalUpdate,
		DeleteWithoutTimeout: resourceIkeProposalDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIkeProposalImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
			},
			"authentication_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authentication_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "pre-shared-keys",
			},
			"dh_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(180, 86400),
			},
		},
	}
}

func resourceIkeProposalCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setIkeProposal(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security ike proposal not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ikeProposalExists, err := checkIkeProposalExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ikeProposalExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security ike proposal %v already exists", d.Get("name").(string)))...)
	}
	if err := setIkeProposal(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_ike_proposal")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	ikeProposalExists, err = checkIkeProposalExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ikeProposalExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ike proposal %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIkeProposalReadWJunSess(d, junSess)...)
}

func resourceIkeProposalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceIkeProposalReadWJunSess(d, junSess)
}

func resourceIkeProposalReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	ikeProposalOptions, err := readIkeProposal(d.Get("name").(string), junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ikeProposalOptions.name == "" {
		d.SetId("")
	} else {
		fillIkeProposalData(d, ikeProposalOptions)
	}

	return nil
}

func resourceIkeProposalUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delIkeProposal(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setIkeProposal(d, junSess); err != nil {
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
	if err := delIkeProposal(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setIkeProposal(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_ike_proposal")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIkeProposalReadWJunSess(d, junSess)...)
}

func resourceIkeProposalDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delIkeProposal(d, junSess); err != nil {
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
	if err := delIkeProposal(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_ike_proposal")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceIkeProposalImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	ikeProposalExists, err := checkIkeProposalExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !ikeProposalExists {
		return nil, fmt.Errorf("don't find security ike proposal with id '%v' (id must be <name>)", d.Id())
	}
	ikeProposalOptions, err := readIkeProposal(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillIkeProposalData(d, ikeProposalOptions)
	result[0] = d

	return result, nil
}

func checkIkeProposalExists(ikeProposal string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ike proposal " + ikeProposal + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setIkeProposal(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security ike proposal " + d.Get("name").(string)
	if d.Get("authentication_method").(string) != "" {
		configSet = append(configSet, setPrefix+" authentication-method "+d.Get("authentication_method").(string))
	}
	if d.Get("authentication_algorithm").(string) != "" {
		configSet = append(configSet, setPrefix+" authentication-algorithm "+d.Get("authentication_algorithm").(string))
	}
	if d.Get("dh_group").(string) != "" {
		configSet = append(configSet, setPrefix+" dh-group "+d.Get("dh_group").(string))
	}
	if d.Get("encryption_algorithm").(string) != "" {
		configSet = append(configSet, setPrefix+" encryption-algorithm "+d.Get("encryption_algorithm").(string))
	}
	if d.Get("lifetime_seconds").(int) != 0 {
		configSet = append(configSet, setPrefix+" lifetime-seconds "+strconv.Itoa(d.Get("lifetime_seconds").(int)))
	}

	return junSess.ConfigSet(configSet)
}

func readIkeProposal(ikeProposal string, junSess *junos.Session,
) (confRead ikeProposalOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ike proposal " + ikeProposal + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = ikeProposal
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "authentication-algorithm "):
				confRead.authenticationAlgorithm = itemTrim
			case balt.CutPrefixInString(&itemTrim, "authentication-method "):
				confRead.authenticationMethod = itemTrim
			case balt.CutPrefixInString(&itemTrim, "dh-group "):
				confRead.dhGroup = itemTrim
			case balt.CutPrefixInString(&itemTrim, "encryption-algorithm "):
				confRead.encryptionAlgorithm = itemTrim
			case balt.CutPrefixInString(&itemTrim, "lifetime-seconds "):
				confRead.lifetimeSeconds, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delIkeProposal(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ike proposal "+d.Get("name").(string))

	return junSess.ConfigSet(configSet)
}

func fillIkeProposalData(d *schema.ResourceData, ikeProposalOptions ikeProposalOptions) {
	if tfErr := d.Set("name", ikeProposalOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_algorithm", ikeProposalOptions.authenticationAlgorithm); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_method", ikeProposalOptions.authenticationMethod); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dh_group", ikeProposalOptions.dhGroup); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("encryption_algorithm", ikeProposalOptions.encryptionAlgorithm); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("lifetime_seconds", ikeProposalOptions.lifetimeSeconds); tfErr != nil {
		panic(tfErr)
	}
}
