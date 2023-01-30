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

type ipsecProposalOptions struct {
	lifetimeSeconds        int
	lifetimeKilobytes      int
	name                   string
	authenticatioAlgorithm string
	encryptionAlgorithm    string
	protocol               string
}

func resourceIpsecProposal() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceIpsecProposalCreate,
		ReadWithoutTimeout:   resourceIpsecProposalRead,
		UpdateWithoutTimeout: resourceIpsecProposalUpdate,
		DeleteWithoutTimeout: resourceIpsecProposalDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIpsecProposalImport,
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
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(180, 86400),
			},
			"lifetime_kilobytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(64, 4294967294),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"esp", "ah"}, false),
			},
		},
	}
}

func resourceIpsecProposalCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setIpsecProposal(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security ipsec proposal not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ipsecProposalExists, err := checkIpsecProposalExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecProposalExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security ipsec proposal %v already exists", d.Get("name").(string)))...)
	}
	if err := setIpsecProposal(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_ipsec_proposal")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	ipsecProposalExists, err = checkIpsecProposalExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecProposalExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ipsec proposal %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIpsecProposalReadWJunSess(d, junSess)...)
}

func resourceIpsecProposalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceIpsecProposalReadWJunSess(d, junSess)
}

func resourceIpsecProposalReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	ipsecProposalOptions, err := readIpsecProposal(d.Get("name").(string), junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ipsecProposalOptions.name == "" {
		d.SetId("")
	} else {
		fillIpsecProposalData(d, ipsecProposalOptions)
	}

	return nil
}

func resourceIpsecProposalUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delIpsecProposal(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setIpsecProposal(d, junSess); err != nil {
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
	if err := delIpsecProposal(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setIpsecProposal(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_ipsec_proposal")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIpsecProposalReadWJunSess(d, junSess)...)
}

func resourceIpsecProposalDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delIpsecProposal(d, junSess); err != nil {
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
	if err := delIpsecProposal(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_ipsec_proposal")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceIpsecProposalImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	ipsecProposalExists, err := checkIpsecProposalExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !ipsecProposalExists {
		return nil, fmt.Errorf("don't find security ipsec proposal with id '%v' (id must be <name>)", d.Id())
	}
	ipsecProposalOptions, err := readIpsecProposal(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillIpsecProposalData(d, ipsecProposalOptions)
	result[0] = d

	return result, nil
}

func checkIpsecProposalExists(ipsecProposal string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ipsec proposal " + ipsecProposal + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setIpsecProposal(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security ipsec proposal " + d.Get("name").(string)
	if d.Get("authentication_algorithm").(string) != "" {
		configSet = append(configSet, setPrefix+" authentication-algorithm "+d.Get("authentication_algorithm").(string))
	}
	if d.Get("encryption_algorithm").(string) != "" {
		configSet = append(configSet, setPrefix+" encryption-algorithm "+d.Get("encryption_algorithm").(string))
	}
	if d.Get("lifetime_seconds").(int) != 0 {
		configSet = append(configSet, setPrefix+" lifetime-seconds "+strconv.Itoa(d.Get("lifetime_seconds").(int)))
	}
	if d.Get("lifetime_kilobytes").(int) != 0 {
		configSet = append(configSet, setPrefix+" lifetime-kilobytes "+strconv.Itoa(d.Get("lifetime_kilobytes").(int)))
	}
	if d.Get("protocol").(string) != "" {
		configSet = append(configSet, setPrefix+" protocol "+d.Get("protocol").(string))
	}

	return junSess.ConfigSet(configSet)
}

func readIpsecProposal(ipsecProposal string, junSess *junos.Session,
) (confRead ipsecProposalOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ipsec proposal " + ipsecProposal + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = ipsecProposal
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
				confRead.authenticatioAlgorithm = itemTrim
			case balt.CutPrefixInString(&itemTrim, "encryption-algorithm "):
				confRead.encryptionAlgorithm = itemTrim
			case balt.CutPrefixInString(&itemTrim, "lifetime-kilobytes "):
				confRead.lifetimeKilobytes, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "lifetime-seconds "):
				confRead.lifetimeSeconds, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "protocol "):
				confRead.protocol = itemTrim
			}
		}
	}

	return confRead, nil
}

func delIpsecProposal(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec proposal "+d.Get("name").(string))

	return junSess.ConfigSet(configSet)
}

func fillIpsecProposalData(d *schema.ResourceData, ipsecProposalOptions ipsecProposalOptions) {
	if tfErr := d.Set("name", ipsecProposalOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_algorithm", ipsecProposalOptions.authenticatioAlgorithm); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("encryption_algorithm", ipsecProposalOptions.encryptionAlgorithm); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("lifetime_kilobytes", ipsecProposalOptions.lifetimeKilobytes); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("lifetime_seconds", ipsecProposalOptions.lifetimeSeconds); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("protocol", ipsecProposalOptions.protocol); tfErr != nil {
		panic(tfErr)
	}
}
