package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type securityLogStreamOptions struct {
	filterThreatAttack bool
	rateLimit          int
	name               string
	format             string
	severity           string
	category           []string
	file               []map[string]interface{}
	host               []map[string]interface{}
}

func resourceSecurityLogStream() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityLogStreamCreate,
		ReadWithoutTimeout:   resourceSecurityLogStreamRead,
		UpdateWithoutTimeout: resourceSecurityLogStreamUpdate,
		DeleteWithoutTimeout: resourceSecurityLogStreamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityLogStreamImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"category": {
				Type:          schema.TypeList,
				Optional:      true,
				MinItems:      1,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"filter_threat_attack"},
			},
			"file": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"host"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"allow_duplicates": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rotation": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(2, 19),
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 3),
						},
					},
				},
			},
			"filter_threat_attack": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"category"},
			},
			"format": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"binary", "sd-syslog", "syslog", "welf",
				}, false),
			},
			"host": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"file"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"routing_instance": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64, formatDefault),
						},
					},
				},
			},
			"rate_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
		},
	}
}

func resourceSecurityLogStreamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityLogStream(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security log stream "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityLogStreamExists, err := checkSecurityLogStreamExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityLogStreamExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("security log stream %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityLogStream(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_security_log_stream")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityLogStreamExists, err = checkSecurityLogStreamExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityLogStreamExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security log stream %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityLogStreamReadWJunSess(d, junSess)...)
}

func resourceSecurityLogStreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityLogStreamReadWJunSess(d, junSess)
}

func resourceSecurityLogStreamReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	securityLogStreamOptions, err := readSecurityLogStream(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if securityLogStreamOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityLogStreamData(d, securityLogStreamOptions)
	}

	return nil
}

func resourceSecurityLogStreamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delLogStream(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityLogStream(d, junSess); err != nil {
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
	if err := delLogStream(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityLogStream(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_security_log_stream")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityLogStreamReadWJunSess(d, junSess)...)
}

func resourceSecurityLogStreamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delLogStream(d.Get("name").(string), junSess); err != nil {
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
	if err := delLogStream(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_security_log_stream")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityLogStreamImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	securityLogStreamExists, err := checkSecurityLogStreamExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !securityLogStreamExists {
		return nil, fmt.Errorf("don't find security log stream with id '%v' (id must be <name>)", d.Id())
	}
	securityLogStreamOptions, err := readSecurityLogStream(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityLogStreamData(d, securityLogStreamOptions)

	result[0] = d

	return result, nil
}

func checkSecurityLogStreamExists(securityLogStream string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security log stream \"" + securityLogStream + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityLogStream(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security log stream \"" + d.Get("name").(string) + "\" "
	for _, v := range d.Get("category").([]interface{}) {
		configSet = append(configSet, setPrefix+"category "+v.(string))
	}
	for _, v := range d.Get("file").([]interface{}) {
		file := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+"file name "+file["name"].(string))
		if file["allow_duplicates"].(bool) {
			configSet = append(configSet, setPrefix+"file allow-duplicates")
		}
		if file["rotation"].(int) != 0 {
			configSet = append(configSet, setPrefix+"file rotation "+strconv.Itoa(file["rotation"].(int)))
		}
		if file["size"].(int) != 0 {
			configSet = append(configSet, setPrefix+"file size "+strconv.Itoa(file["size"].(int)))
		}
	}
	if d.Get("filter_threat_attack").(bool) {
		configSet = append(configSet, setPrefix+"filter threat-attack")
	}
	if d.Get("format").(string) != "" {
		configSet = append(configSet, setPrefix+"format "+d.Get("format").(string))
	}
	for _, v := range d.Get("host").([]interface{}) {
		host := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+"host "+host["ip_address"].(string))
		if host["port"].(int) != 0 {
			configSet = append(configSet, setPrefix+"host port "+strconv.Itoa(
				host["port"].(int)))
		}
		if host["routing_instance"].(string) != "" {
			configSet = append(configSet, setPrefix+"host routing-instance "+
				host["routing_instance"].(string))
		}
	}
	if d.Get("rate_limit").(int) != 0 {
		configSet = append(configSet, setPrefix+"rate-limit "+
			strconv.Itoa(d.Get("rate_limit").(int)))
	}
	if d.Get("severity").(string) != "" {
		configSet = append(configSet, setPrefix+"severity "+
			d.Get("severity").(string))
	}

	return junSess.ConfigSet(configSet)
}

func readSecurityLogStream(securityLogStream string, junSess *junos.Session,
) (confRead securityLogStreamOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security log stream \"" + securityLogStream + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = securityLogStream
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "category "):
				confRead.category = append(confRead.category, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "file "):
				if len(confRead.file) == 0 {
					confRead.file = append(confRead.file, map[string]interface{}{
						"name":             "",
						"allow_duplicates": false,
						"rotation":         0,
						"size":             0,
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "name "):
					confRead.file[0]["name"] = itemTrim
				case itemTrim == "allow-duplicates":
					confRead.file[0]["allow_duplicates"] = true
				case balt.CutPrefixInString(&itemTrim, "rotation "):
					confRead.file[0]["rotation"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "size "):
					confRead.file[0]["size"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			case itemTrim == "filter threat-attack":
				confRead.filterThreatAttack = true
			case balt.CutPrefixInString(&itemTrim, "format "):
				confRead.format = itemTrim
			case balt.CutPrefixInString(&itemTrim, "host "):
				if len(confRead.host) == 0 {
					confRead.host = append(confRead.host, map[string]interface{}{
						"ip_address":       "",
						"port":             0,
						"routing_instance": "",
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "port "):
					confRead.host[0]["port"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "routing-instance "):
					confRead.host[0]["routing_instance"] = itemTrim
				default:
					confRead.host[0]["ip_address"] = itemTrim
				}
			case balt.CutPrefixInString(&itemTrim, "rate-limit "):
				confRead.rateLimit, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "severity "):
				confRead.severity = itemTrim
			}
		}
	}

	return confRead, nil
}

func delLogStream(securityLogStream string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security log stream \""+securityLogStream+"\"")

	return junSess.ConfigSet(configSet)
}

func fillSecurityLogStreamData(d *schema.ResourceData, securityLogStreamOptions securityLogStreamOptions) {
	if tfErr := d.Set("name", securityLogStreamOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("category", securityLogStreamOptions.category); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("file", securityLogStreamOptions.file); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("filter_threat_attack", securityLogStreamOptions.filterThreatAttack); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("format", securityLogStreamOptions.format); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("host", securityLogStreamOptions.host); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("rate_limit", securityLogStreamOptions.rateLimit); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("severity", securityLogStreamOptions.severity); tfErr != nil {
		panic(tfErr)
	}
}
