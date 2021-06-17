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

type applicationOptions struct {
	inactivityTimeout int
	name              string
	destinationPort   string
	protocol          string
	sourcePort        string
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
			"destination_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inactivity_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(4, 86400),
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_port": {
				Type:     schema.TypeString,
				Optional: true,
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
	applicationConfig, err := sess.command("show configuration applications application "+
		application+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if applicationConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setApplication(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set applications application " + d.Get("name").(string) + " "
	if v := d.Get("destination_port").(string); v != "" {
		configSet = append(configSet, setPrefix+"destination-port "+v)
	}
	if v := d.Get("inactivity_timeout").(int); v != 0 {
		configSet = append(configSet, setPrefix+"inactivity-timeout "+strconv.Itoa(v))
	}
	if v := d.Get("protocol").(string); v != "" {
		configSet = append(configSet, setPrefix+"protocol "+v)
	}
	if v := d.Get("source_port").(string); v != "" {
		configSet = append(configSet, setPrefix+"source-port "+v)
	}

	return sess.configSet(configSet, jnprSess)
}

func readApplication(application string, m interface{}, jnprSess *NetconfObject) (applicationOptions, error) {
	sess := m.(*Session)
	var confRead applicationOptions

	applicationConfig, err := sess.command("show configuration applications application "+
		application+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if applicationConfig != emptyWord {
		confRead.name = application
		for _, item := range strings.Split(applicationConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "destination-port "):
				confRead.destinationPort = strings.TrimPrefix(itemTrim, "destination-port ")
			case strings.HasPrefix(itemTrim, "inactivity-timeout "):
				var err error
				confRead.inactivityTimeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "inactivity-timeout "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "protocol "):
				confRead.protocol = strings.TrimPrefix(itemTrim, "protocol ")
			case strings.HasPrefix(itemTrim, "source-port "):
				confRead.sourcePort = strings.TrimPrefix(itemTrim, "source-port ")
			}
		}
	}

	return confRead, nil
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
	if tfErr := d.Set("destination_port", applicationOptions.destinationPort); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inactivity_timeout", applicationOptions.inactivityTimeout); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("protocol", applicationOptions.protocol); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("source_port", applicationOptions.sourcePort); tfErr != nil {
		panic(tfErr)
	}
}
