package junos

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type applicationOptions struct {
	inactivityTimeoutNever bool
	inactivityTimeout      int
	name                   string
	applicationProtocol    string
	description            string
	destinationPort        string
	etherType              string
	protocol               string
	rpcProgramNumber       string
	sourcePort             string
	uuid                   string
	term                   []map[string]interface{}
}

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceApplicationCreate,
		ReadWithoutTimeout:   resourceApplicationRead,
		UpdateWithoutTimeout: resourceApplicationUpdate,
		DeleteWithoutTimeout: resourceApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"application_protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ether_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^0[xX][0-9a-fA-F]{4}$`), "must be in hex (example: 0x8906)"),
			},
			"inactivity_timeout": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"inactivity_timeout_never"},
				ValidateFunc:  validation.IntBetween(4, 86400),
			},
			"inactivity_timeout_never": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"inactivity_timeout"},
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rpc_program_number": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^\d+(-\d+)?$`), "must be an integer or a range of integers"),
			},
			"source_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"term": {
				Type:     schema.TypeList,
				Optional: true,
				ConflictsWith: []string{
					"application_protocol",
					"destination_port",
					"inactivity_timeout",
					"inactivity_timeout_never",
					"protocol",
					"rpc_program_number",
					"source_port",
					"uuid",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"alg": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp6_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp6_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"inactivity_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(4, 86400),
						},
						"inactivity_timeout_never": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rpc_program_number": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^\d+(-\d+)?$`), "must be an integer or a range of integers"),
						},
						"source_port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
								"must be of the form xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"),
						},
					},
				},
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
					"must be of the form xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"),
			},
		},
	}
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setApplication(d, clt, nil); err != nil {
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
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	appExists, err := checkApplicationExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if appExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("application %v already exists", d.Get("name").(string)))...)
	}
	if err := setApplication(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_application", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	appExists, err = checkApplicationExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if appExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("application %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceApplicationReadWJunSess(d, clt, junSess)...)
}

func resourceApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceApplicationReadWJunSess(d, clt, junSess)
}

func resourceApplicationReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession) diag.Diagnostics {
	mutex.Lock()
	applicationOptions, err := readApplication(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if applicationOptions.name == "" {
		d.SetId("")
	} else {
		fillApplicationData(d, applicationOptions)
	}

	return nil
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delApplication(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setApplication(d, clt, nil); err != nil {
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
	if err := delApplication(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setApplication(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_application", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceApplicationReadWJunSess(d, clt, junSess)...)
}

func resourceApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delApplication(d, clt, nil); err != nil {
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
	if err := delApplication(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_application", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceApplicationImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	appExists, err := checkApplicationExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !appExists {
		return nil, fmt.Errorf("don't find application with id '%v' (id must be <name>)", d.Id())
	}
	applicationOptions, err := readApplication(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillApplicationData(d, applicationOptions)
	result[0] = d

	return result, nil
}

func checkApplicationExists(application string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+
		"applications application "+application+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setApplication(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := "set applications application " + d.Get("name").(string) + " "
	if v := d.Get("application_protocol").(string); v != "" {
		configSet = append(configSet, setPrefix+"application-protocol "+v)
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := d.Get("destination_port").(string); v != "" {
		configSet = append(configSet, setPrefix+"destination-port "+v)
	}
	if v := d.Get("ether_type").(string); v != "" {
		configSet = append(configSet, setPrefix+"ether-type "+v)
	}
	if v := d.Get("inactivity_timeout").(int); v != 0 {
		configSet = append(configSet, setPrefix+"inactivity-timeout "+strconv.Itoa(v))
	} else if d.Get("inactivity_timeout_never").(bool) {
		configSet = append(configSet, setPrefix+"inactivity-timeout never")
	}
	if v := d.Get("protocol").(string); v != "" {
		configSet = append(configSet, setPrefix+"protocol "+v)
	}
	if v := d.Get("rpc_program_number").(string); v != "" {
		configSet = append(configSet, setPrefix+"rpc-program-number "+v)
	}
	if v := d.Get("source_port").(string); v != "" {
		configSet = append(configSet, setPrefix+"source-port "+v)
	}
	termName := make([]string, 0)
	for _, v := range d.Get("term").([]interface{}) {
		term := v.(map[string]interface{})
		if bchk.InSlice(term["name"].(string), termName) {
			return fmt.Errorf("multiple blocks term with the same name %s", term["name"].(string))
		}
		termName = append(termName, term["name"].(string))
		if err := setApplicationTerm(setPrefix, term, clt, junSess); err != nil {
			return err
		}
	}
	if v := d.Get("uuid").(string); v != "" {
		configSet = append(configSet, setPrefix+"uuid "+v)
	}

	return clt.configSet(configSet, junSess)
}

func setApplicationTerm(setApp string, term map[string]interface{}, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)
	setPrefix := setApp + "term " + term["name"].(string) + " "

	configSet = append(configSet, setPrefix)
	configSet = append(configSet, setPrefix+"protocol "+term["protocol"].(string))
	if v := term["alg"].(string); v != "" {
		configSet = append(configSet, setPrefix+"alg "+v)
	}
	if v := term["destination_port"].(string); v != "" {
		configSet = append(configSet, setPrefix+"destination-port "+v)
	}
	if v := term["icmp_code"].(string); v != "" {
		configSet = append(configSet, setPrefix+"icmp-code "+v)
	}
	if v := term["icmp_type"].(string); v != "" {
		configSet = append(configSet, setPrefix+"icmp-type "+v)
	}
	if v := term["icmp6_code"].(string); v != "" {
		configSet = append(configSet, setPrefix+"icmp6-code "+v)
	}
	if v := term["icmp6_type"].(string); v != "" {
		configSet = append(configSet, setPrefix+"icmp6-type "+v)
	}
	if v := term["inactivity_timeout"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"inactivity-timeout "+strconv.Itoa(v))
		if term["inactivity_timeout_never"].(bool) {
			return fmt.Errorf("conflict between 'inactivity_timeout' and 'inactivity_timeout_never' "+
				"in term %s", term["name"].(string))
		}
	} else if term["inactivity_timeout_never"].(bool) {
		configSet = append(configSet, setPrefix+"inactivity-timeout never")
	}
	if v := term["rpc_program_number"].(string); v != "" {
		configSet = append(configSet, setPrefix+"rpc-program-number "+v)
	}
	if v := term["source_port"].(string); v != "" {
		configSet = append(configSet, setPrefix+"source-port "+v)
	}
	if v := term["uuid"].(string); v != "" {
		configSet = append(configSet, setPrefix+"uuid "+v)
	}

	return clt.configSet(configSet, junSess)
}

func readApplication(application string, clt *Client, junSess *junosSession) (confRead applicationOptions, err error) {
	showConfig, err := clt.command(cmdShowConfig+
		"applications application "+application+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = application
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			if err := confRead.readLine(itemTrim); err != nil {
				return confRead, err
			}
		}
	}

	return confRead, nil
}

func (confRead *applicationOptions) readLine(itemTrim string) (err error) {
	balt.CutPrefixInString(&itemTrim, setLS)
	switch {
	case balt.CutPrefixInString(&itemTrim, "application-protocol "):
		confRead.applicationProtocol = itemTrim
	case balt.CutPrefixInString(&itemTrim, "description "):
		confRead.description = strings.Trim(itemTrim, "\"")
	case balt.CutPrefixInString(&itemTrim, "destination-port "):
		confRead.destinationPort = itemTrim
	case balt.CutPrefixInString(&itemTrim, "ether-type "):
		confRead.etherType = itemTrim
	case itemTrim == "inactivity-timeout never":
		confRead.inactivityTimeoutNever = true
	case balt.CutPrefixInString(&itemTrim, "inactivity-timeout "):
		confRead.inactivityTimeout, err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "protocol "):
		confRead.protocol = itemTrim
	case balt.CutPrefixInString(&itemTrim, "rpc-program-number "):
		confRead.rpcProgramNumber = itemTrim
	case balt.CutPrefixInString(&itemTrim, "source-port "):
		confRead.sourcePort = itemTrim
	case balt.CutPrefixInString(&itemTrim, "term "):
		itemTrimFields := strings.Split(itemTrim, " ")
		termOpts := map[string]interface{}{
			"name":                     itemTrimFields[0],
			"protocol":                 "",
			"alg":                      "",
			"destination_port":         "",
			"icmp_code":                "",
			"icmp_type":                "",
			"icmp6_code":               "",
			"icmp6_type":               "",
			"inactivity_timeout":       0,
			"inactivity_timeout_never": false,
			"rpc_program_number":       "",
			"source_port":              "",
			"uuid":                     "",
		}
		confRead.term = copyAndRemoveItemMapList("name", termOpts, confRead.term)
		if err := readApplicationTerm(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" "), termOpts); err != nil {
			return err
		}
		confRead.term = append(confRead.term, termOpts)
	case balt.CutPrefixInString(&itemTrim, "uuid "):
		confRead.uuid = itemTrim
	}

	return nil
}

func readApplicationTerm(itemTrim string, term map[string]interface{}) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "protocol "):
		term["protocol"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "alg "):
		term["alg"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "destination-port "):
		term["destination_port"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "icmp-code "):
		term["icmp_code"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "icmp-type "):
		term["icmp_type"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "icmp6-code "):
		term["icmp6_code"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "icmp6-type "):
		term["icmp6_type"] = itemTrim
	case itemTrim == "inactivity-timeout never":
		term["inactivity_timeout_never"] = true
	case balt.CutPrefixInString(&itemTrim, "inactivity-timeout "):
		term["inactivity_timeout"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "rpc-program-number "):
		term["rpc_program_number"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "source-port "):
		term["source_port"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "uuid "):
		term["uuid"] = itemTrim
	}

	return nil
}

func delApplication(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete applications application "+d.Get("name").(string))

	return clt.configSet(configSet, junSess)
}

func fillApplicationData(d *schema.ResourceData, applicationOptions applicationOptions) {
	if tfErr := d.Set("name", applicationOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("application_protocol", applicationOptions.applicationProtocol); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", applicationOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("destination_port", applicationOptions.destinationPort); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ether_type", applicationOptions.etherType); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inactivity_timeout", applicationOptions.inactivityTimeout); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inactivity_timeout_never", applicationOptions.inactivityTimeoutNever); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("protocol", applicationOptions.protocol); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("rpc_program_number", applicationOptions.rpcProgramNumber); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("source_port", applicationOptions.sourcePort); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("term", applicationOptions.term); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("uuid", applicationOptions.uuid); tfErr != nil {
		panic(tfErr)
	}
}
