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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSecurityLogStream(d, clt, nil); err != nil {
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
		return diag.FromErr(fmt.Errorf("security log stream "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityLogStreamExists, err := checkSecurityLogStreamExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityLogStreamExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("security log stream %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityLogStream(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_security_log_stream", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityLogStreamExists, err = checkSecurityLogStreamExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityLogStreamExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security log stream %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityLogStreamReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityLogStreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSecurityLogStreamReadWJunSess(d, clt, junSess)
}

func resourceSecurityLogStreamReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	securityLogStreamOptions, err := readSecurityLogStream(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
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
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delLogStream(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityLogStream(d, clt, nil); err != nil {
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
	if err := delLogStream(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityLogStream(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_security_log_stream", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityLogStreamReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityLogStreamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delLogStream(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delLogStream(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_security_log_stream", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityLogStreamImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	securityLogStreamExists, err := checkSecurityLogStreamExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !securityLogStreamExists {
		return nil, fmt.Errorf("don't find security log stream with id '%v' (id must be <name>)", d.Id())
	}
	securityLogStreamOptions, err := readSecurityLogStream(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityLogStreamData(d, securityLogStreamOptions)

	result[0] = d

	return result, nil
}

func checkSecurityLogStreamExists(securityLogStream string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+
		"security log stream \""+securityLogStream+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityLogStream(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
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

	return clt.configSet(configSet, junSess)
}

func readSecurityLogStream(securityLogStream string, clt *Client, junSess *junosSession,
) (securityLogStreamOptions, error) {
	var confRead securityLogStreamOptions

	showConfig, err := clt.command(cmdShowConfig+
		"security log stream \""+securityLogStream+"\""+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = securityLogStream
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "category "):
				confRead.category = append(confRead.category, strings.TrimPrefix(itemTrim, "category "))
			case strings.HasPrefix(itemTrim, "file "):
				if len(confRead.file) == 0 {
					confRead.file = append(confRead.file, map[string]interface{}{
						"name":             "",
						"allow_duplicates": false,
						"rotation":         0,
						"size":             0,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "file name "):
					confRead.file[0]["name"] = strings.TrimPrefix(itemTrim, "file name ")
				case itemTrim == "file allow-duplicates":
					confRead.file[0]["allow_duplicates"] = true
				case strings.HasPrefix(itemTrim, "file rotation "):
					var err error
					confRead.file[0]["rotation"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "file rotation "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "file size "):
					var err error
					confRead.file[0]["size"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "file size "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			case itemTrim == "filter threat-attack":
				confRead.filterThreatAttack = true
			case strings.HasPrefix(itemTrim, "format "):
				confRead.format = strings.TrimPrefix(itemTrim, "format ")
			case strings.HasPrefix(itemTrim, "host "):
				if len(confRead.host) == 0 {
					confRead.host = append(confRead.host, map[string]interface{}{
						"ip_address":       "",
						"port":             0,
						"routing_instance": "",
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "host port "):
					var err error
					confRead.host[0]["port"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "host port "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "host routing-instance "):
					confRead.host[0]["routing_instance"] = strings.TrimPrefix(itemTrim, "host routing-instance ")
				default:
					confRead.host[0]["ip_address"] = strings.TrimPrefix(itemTrim, "host ")
				}
			case strings.HasPrefix(itemTrim, "rate-limit "):
				var err error
				confRead.rateLimit, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "rate-limit "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "severity "):
				confRead.severity = strings.TrimPrefix(itemTrim, "severity ")
			}
		}
	}

	return confRead, nil
}

func delLogStream(securityLogStream string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security log stream \""+securityLogStream+"\"")

	return clt.configSet(configSet, junSess)
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
