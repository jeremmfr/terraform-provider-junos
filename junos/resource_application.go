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
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceApplicationImport,
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
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setApplication(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	appExists, err := checkApplicationExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if appExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("application %v already exists", d.Get("name").(string)))...)
	}
	if err := setApplication(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_application", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	appExists, err = checkApplicationExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if appExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("application %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceApplicationReadWJnprSess(d, m, jnprSess)...)
}

func resourceApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceApplicationReadWJnprSess(d, m, jnprSess)
}

func resourceApplicationReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	applicationOptions, err := readApplication(d.Get("name").(string), m, jnprSess)
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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delApplication(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setApplication(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_application", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceApplicationReadWJnprSess(d, m, jnprSess)...)
}

func resourceApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delApplication(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_application", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceApplicationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	appExists, err := checkApplicationExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !appExists {
		return nil, fmt.Errorf("don't find application with id '%v' (id must be <name>)", d.Id())
	}
	applicationOptions, err := readApplication(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillApplicationData(d, applicationOptions)
	result[0] = d

	return result, nil
}

func checkApplicationExists(application string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration"+
		" applications application "+application+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setApplication(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
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
		if bchk.StringInSlice(term["name"].(string), termName) {
			return fmt.Errorf("multiple blocks term with the same name %s", term["name"].(string))
		}
		termName = append(termName, term["name"].(string))
		if err := setApplicationTerm(setPrefix, term, sess, jnprSess); err != nil {
			return err
		}
	}
	if v := d.Get("uuid").(string); v != "" {
		configSet = append(configSet, setPrefix+"uuid "+v)
	}

	return sess.configSet(configSet, jnprSess)
}

func setApplicationTerm(setApp string, term map[string]interface{}, sess *Session, jnprSess *NetconfObject) error {
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

	return sess.configSet(configSet, jnprSess)
}

func readApplication(application string, m interface{}, jnprSess *NetconfObject) (applicationOptions, error) {
	sess := m.(*Session)
	var confRead applicationOptions

	showConfig, err := sess.command("show configuration"+
		" applications application "+application+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		confRead.name = application
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "application-protocol "):
				confRead.applicationProtocol = strings.TrimPrefix(itemTrim, "application-protocol ")
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "destination-port "):
				confRead.destinationPort = strings.TrimPrefix(itemTrim, "destination-port ")
			case strings.HasPrefix(itemTrim, "ether-type "):
				confRead.etherType = strings.TrimPrefix(itemTrim, "ether-type ")
			case itemTrim == "inactivity-timeout never":
				confRead.inactivityTimeoutNever = true
			case strings.HasPrefix(itemTrim, "inactivity-timeout "):
				var err error
				confRead.inactivityTimeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "inactivity-timeout "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "protocol "):
				confRead.protocol = strings.TrimPrefix(itemTrim, "protocol ")
			case strings.HasPrefix(itemTrim, "rpc-program-number "):
				confRead.rpcProgramNumber = strings.TrimPrefix(itemTrim, "rpc-program-number ")
			case strings.HasPrefix(itemTrim, "source-port "):
				confRead.sourcePort = strings.TrimPrefix(itemTrim, "source-port ")
			case strings.HasPrefix(itemTrim, "term "):
				itemTermList := strings.Split(strings.TrimPrefix(itemTrim, "term "), " ")
				termOpts := map[string]interface{}{
					"name":                     itemTermList[0],
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
				if err := readApplicationTerm(strings.TrimPrefix(itemTrim, "term "+itemTermList[0]+" "), termOpts); err != nil {
					return confRead, err
				}
				confRead.term = append(confRead.term, termOpts)
			case strings.HasPrefix(itemTrim, "uuid "):
				confRead.uuid = strings.TrimPrefix(itemTrim, "uuid ")
			}
		}
	}

	return confRead, nil
}

func readApplicationTerm(itemTrim string, term map[string]interface{}) error {
	switch {
	case strings.HasPrefix(itemTrim, "protocol "):
		term["protocol"] = strings.TrimPrefix(itemTrim, "protocol ")
	case strings.HasPrefix(itemTrim, "alg "):
		term["alg"] = strings.TrimPrefix(itemTrim, "alg ")
	case strings.HasPrefix(itemTrim, "destination-port "):
		term["destination_port"] = strings.TrimPrefix(itemTrim, "destination-port ")
	case strings.HasPrefix(itemTrim, "icmp-code "):
		term["icmp_code"] = strings.TrimPrefix(itemTrim, "icmp-code ")
	case strings.HasPrefix(itemTrim, "icmp-type "):
		term["icmp_type"] = strings.TrimPrefix(itemTrim, "icmp-type ")
	case strings.HasPrefix(itemTrim, "icmp6-code "):
		term["icmp6_code"] = strings.TrimPrefix(itemTrim, "icmp6-code ")
	case strings.HasPrefix(itemTrim, "icmp6-type "):
		term["icmp6_type"] = strings.TrimPrefix(itemTrim, "icmp6-type ")
	case itemTrim == "inactivity-timeout never":
		term["inactivity_timeout_never"] = true
	case strings.HasPrefix(itemTrim, "inactivity-timeout "):
		var err error
		term["inactivity_timeout"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "inactivity-timeout "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "rpc-program-number "):
		term["rpc_program_number"] = strings.TrimPrefix(itemTrim, "rpc-program-number ")
	case strings.HasPrefix(itemTrim, "source-port "):
		term["source_port"] = strings.TrimPrefix(itemTrim, "source-port ")
	case strings.HasPrefix(itemTrim, "uuid "):
		term["uuid"] = strings.TrimPrefix(itemTrim, "uuid ")
	}

	return nil
}

func delApplication(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete applications application "+d.Get("name").(string))

	return sess.configSet(configSet, jnprSess)
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
