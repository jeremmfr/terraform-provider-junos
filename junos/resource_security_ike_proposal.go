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
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setIkeProposal(d, sess, nil); err != nil {
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
		return diag.FromErr(fmt.Errorf("security ike proposal not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ikeProposalExists, err := checkIkeProposalExists(d.Get("name").(string), sess, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ikeProposalExists {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security ike proposal %v already exists", d.Get("name").(string)))...)
	}
	if err := setIkeProposal(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_ike_proposal", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ikeProposalExists, err = checkIkeProposalExists(d.Get("name").(string), sess, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ikeProposalExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ike proposal %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIkeProposalReadWJunSess(d, sess, junSess)...)
}

func resourceIkeProposalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)

	return resourceIkeProposalReadWJunSess(d, sess, junSess)
}

func resourceIkeProposalReadWJunSess(d *schema.ResourceData, sess *Session, junSess *junosSession) diag.Diagnostics {
	mutex.Lock()
	ikeProposalOptions, err := readIkeProposal(d.Get("name").(string), sess, junSess)
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
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delIkeProposal(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setIkeProposal(d, sess, nil); err != nil {
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
	if err := delIkeProposal(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setIkeProposal(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_ike_proposal", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIkeProposalReadWJunSess(d, sess, junSess)...)
}

func resourceIkeProposalDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delIkeProposal(d, sess, nil); err != nil {
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
	if err := delIkeProposal(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_ike_proposal", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceIkeProposalImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	ikeProposalExists, err := checkIkeProposalExists(d.Id(), sess, junSess)
	if err != nil {
		return nil, err
	}
	if !ikeProposalExists {
		return nil, fmt.Errorf("don't find security ike proposal with id '%v' (id must be <name>)", d.Id())
	}
	ikeProposalOptions, err := readIkeProposal(d.Id(), sess, junSess)
	if err != nil {
		return nil, err
	}
	fillIkeProposalData(d, ikeProposalOptions)
	result[0] = d

	return result, nil
}

func checkIkeProposalExists(ikeProposal string, sess *Session, junSess *junosSession) (bool, error) {
	showConfig, err := sess.command(cmdShowConfig+
		"security ike proposal "+ikeProposal+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setIkeProposal(d *schema.ResourceData, sess *Session, junSess *junosSession) error {
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

	return sess.configSet(configSet, junSess)
}

func readIkeProposal(ikeProposal string, sess *Session, junSess *junosSession) (ikeProposalOptions, error) {
	var confRead ikeProposalOptions

	showConfig, err := sess.command(cmdShowConfig+
		"security ike proposal "+ikeProposal+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = ikeProposal
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
				confRead.authenticationAlgorithm = strings.TrimPrefix(itemTrim, "authentication-algorithm ")
			case strings.HasPrefix(itemTrim, "authentication-method "):
				confRead.authenticationMethod = strings.TrimPrefix(itemTrim, "authentication-method ")
			case strings.HasPrefix(itemTrim, "dh-group "):
				confRead.dhGroup = strings.TrimPrefix(itemTrim, "dh-group ")
			case strings.HasPrefix(itemTrim, "encryption-algorithm"):
				confRead.encryptionAlgorithm = strings.TrimPrefix(itemTrim, "encryption-algorithm ")
			case strings.HasPrefix(itemTrim, "lifetime-seconds"):
				confRead.lifetimeSeconds, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "lifetime-seconds "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delIkeProposal(d *schema.ResourceData, sess *Session, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ike proposal "+d.Get("name").(string))

	return sess.configSet(configSet, junSess)
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
