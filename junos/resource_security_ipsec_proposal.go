package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setIpsecProposal(d, clt, nil); err != nil {
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
	if !checkCompatibilitySecurity(junSess) {
		return diag.FromErr(fmt.Errorf("security ipsec proposal not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ipsecProposalExists, err := checkIpsecProposalExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecProposalExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security ipsec proposal %v already exists", d.Get("name").(string)))...)
	}
	if err := setIpsecProposal(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_security_ipsec_proposal", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ipsecProposalExists, err = checkIpsecProposalExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecProposalExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ipsec proposal %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIpsecProposalReadWJunSess(d, clt, junSess)...)
}

func resourceIpsecProposalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceIpsecProposalReadWJunSess(d, clt, junSess)
}

func resourceIpsecProposalReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	ipsecProposalOptions, err := readIpsecProposal(d.Get("name").(string), clt, junSess)
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
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delIpsecProposal(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setIpsecProposal(d, clt, nil); err != nil {
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
	if err := delIpsecProposal(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setIpsecProposal(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_security_ipsec_proposal", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIpsecProposalReadWJunSess(d, clt, junSess)...)
}

func resourceIpsecProposalDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delIpsecProposal(d, clt, nil); err != nil {
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
	if err := delIpsecProposal(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_security_ipsec_proposal", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceIpsecProposalImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	ipsecProposalExists, err := checkIpsecProposalExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !ipsecProposalExists {
		return nil, fmt.Errorf("don't find security ipsec proposal with id '%v' (id must be <name>)", d.Id())
	}
	ipsecProposalOptions, err := readIpsecProposal(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillIpsecProposalData(d, ipsecProposalOptions)
	result[0] = d

	return result, nil
}

func checkIpsecProposalExists(ipsecProposal string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+
		"security ipsec proposal "+ipsecProposal+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setIpsecProposal(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
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

	return clt.configSet(configSet, junSess)
}

func readIpsecProposal(ipsecProposal string, clt *Client, junSess *junosSession) (ipsecProposalOptions, error) {
	var confRead ipsecProposalOptions

	showConfig, err := clt.command(cmdShowConfig+
		"security ipsec proposal "+ipsecProposal+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = ipsecProposal
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "authentication-algorithm "):
				confRead.authenticatioAlgorithm = strings.TrimPrefix(itemTrim, "authentication-algorithm ")
			case strings.HasPrefix(itemTrim, "encryption-algorithm "):
				confRead.encryptionAlgorithm = strings.TrimPrefix(itemTrim, "encryption-algorithm ")
			case strings.HasPrefix(itemTrim, "lifetime-kilobytes "):
				confRead.lifetimeKilobytes, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "lifetime-kilobytes "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "lifetime-seconds "):
				confRead.lifetimeSeconds, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "lifetime-seconds "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "protocol "):
				confRead.protocol = strings.TrimPrefix(itemTrim, "protocol ")
			}
		}
	}

	return confRead, nil
}

func delIpsecProposal(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec proposal "+d.Get("name").(string))

	return clt.configSet(configSet, junSess)
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
