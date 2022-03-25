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

type ntpServerOptions struct {
	prefer          bool
	key             int
	version         int
	address         string
	routingInstance string
}

func resourceSystemNtpServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemNtpServerCreate,
		ReadContext:   resourceSystemNtpServerRead,
		UpdateContext: resourceSystemNtpServerUpdate,
		DeleteContext: resourceSystemNtpServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSystemNtpServerImport,
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"key": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65534),
			},
			"prefer": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 4),
			},
		},
	}
}

func resourceSystemNtpServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSystemNtpServer(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("address").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	ntpServerExists, err := checkSystemNtpServerExists(d.Get("address").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ntpServerExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("system ntp server %v already exists", d.Get("address").(string)))...)
	}

	if err := setSystemNtpServer(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_system_ntp_server", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ntpServerExists, err = checkSystemNtpServerExists(d.Get("address").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ntpServerExists {
		d.SetId(d.Get("address").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system ntp server %v not exists after commit "+
			"=> check your config", d.Get("address").(string)))...)
	}

	return append(diagWarns, resourceSystemNtpServerReadWJnprSess(d, m, jnprSess)...)
}

func resourceSystemNtpServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSystemNtpServerReadWJnprSess(d, m, jnprSess)
}

func resourceSystemNtpServerReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
) diag.Diagnostics {
	mutex.Lock()
	ntpServerOptions, err := readSystemNtpServer(d.Get("address").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ntpServerOptions.address == "" {
		d.SetId("")
	} else {
		fillSystemNtpServerData(d, ntpServerOptions)
	}

	return nil
}

func resourceSystemNtpServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delSystemNtpServer(d.Get("address").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystemNtpServer(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSystemNtpServer(d.Get("address").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemNtpServer(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_system_ntp_server", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	d.Partial(false)

	return append(diagWarns, resourceSystemNtpServerReadWJnprSess(d, m, jnprSess)...)
}

func resourceSystemNtpServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delSystemNtpServer(d.Get("address").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSystemNtpServer(d.Get("address").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_system_ntp_server", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSystemNtpServerImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	ntpServerExists, err := checkSystemNtpServerExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ntpServerExists {
		return nil, fmt.Errorf("don't find system ntp server with id '%v' (id must be <address>)", d.Id())
	}
	ntpServerOptions, err := readSystemNtpServer(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSystemNtpServerData(d, ntpServerOptions)

	result[0] = d

	return result, nil
}

func checkSystemNtpServerExists(address string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command(cmdShowConfig+"system ntp server "+address+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSystemNtpServer(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set system ntp server " + d.Get("address").(string)
	configSet := []string{setPrefix}

	if d.Get("key").(int) != 0 {
		configSet = append(configSet, setPrefix+" key "+strconv.Itoa(d.Get("key").(int)))
	}
	if d.Get("prefer").(bool) {
		configSet = append(configSet, setPrefix+" prefer")
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, setPrefix+" routing-instance "+d.Get("routing_instance").(string))
	}
	if d.Get("version").(int) != 0 {
		configSet = append(configSet, setPrefix+" version "+strconv.Itoa(d.Get("version").(int)))
	}

	return sess.configSet(configSet, jnprSess)
}

func readSystemNtpServer(address string, m interface{}, jnprSess *NetconfObject) (ntpServerOptions, error) {
	sess := m.(*Session)
	var confRead ntpServerOptions

	showConfig, err := sess.command(cmdShowConfig+"system ntp server "+address+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.address = address
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "key "):
				var err error
				confRead.key, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "key "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "prefer":
				confRead.prefer = true
			case strings.HasPrefix(itemTrim, "routing-instance "):
				confRead.routingInstance = strings.TrimPrefix(itemTrim, "routing-instance ")
			case strings.HasPrefix(itemTrim, "version "):
				var err error
				confRead.version, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "version "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delSystemNtpServer(address string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system ntp server "+address)

	return sess.configSet(configSet, jnprSess)
}

func fillSystemNtpServerData(d *schema.ResourceData, ntpServerOptions ntpServerOptions) {
	if tfErr := d.Set("address", ntpServerOptions.address); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("key", ntpServerOptions.key); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("prefer", ntpServerOptions.prefer); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", ntpServerOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("version", ntpServerOptions.version); tfErr != nil {
		panic(tfErr)
	}
}
