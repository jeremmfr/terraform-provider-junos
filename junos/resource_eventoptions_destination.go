package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	jdecode "github.com/jeremmfr/junosdecode"
)

type eventoptionsDestinationOptions struct {
	transferDelay int
	name          string
	archiveSite   []map[string]interface{}
}

func resourceEventoptionsDestination() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventoptionsDestinationCreate,
		ReadContext:   resourceEventoptionsDestinationRead,
		UpdateContext: resourceEventoptionsDestinationUpdate,
		DeleteContext: resourceEventoptionsDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceEventoptionsDestinationImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"archive_site": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringDoesNotContainAny(" "),
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"transfer_delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(0, 4294967295),
			},
		},
	}
}

func resourceEventoptionsDestinationCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setEventoptionsDestination(d, m, nil); err != nil {
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
	eventoptionsDestinationExists, err := checkEventoptionsDestinationExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if eventoptionsDestinationExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("event-options destinations %v already exists", d.Get("name").(string)))...)
	}

	if err := setEventoptionsDestination(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_eventoptions_destination", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	eventoptionsDestinationExists, err = checkEventoptionsDestinationExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if eventoptionsDestinationExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("event-options destinations %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceEventoptionsDestinationReadWJnprSess(d, m, jnprSess)...)
}

func resourceEventoptionsDestinationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceEventoptionsDestinationReadWJnprSess(d, m, jnprSess)
}

func resourceEventoptionsDestinationReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	eventoptionsDestinationOptions, err := readEventoptionsDestination(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if eventoptionsDestinationOptions.name == "" {
		d.SetId("")
	} else {
		fillEventoptionsDestinationData(d, eventoptionsDestinationOptions)
	}

	return nil
}

func resourceEventoptionsDestinationUpdate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delEventoptionsDestination(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setEventoptionsDestination(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_eventoptions_destination", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceEventoptionsDestinationReadWJnprSess(d, m, jnprSess)...)
}

func resourceEventoptionsDestinationDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delEventoptionsDestination(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_eventoptions_destination", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceEventoptionsDestinationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	eventoptionsDestinationExists, err := checkEventoptionsDestinationExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !eventoptionsDestinationExists {
		return nil, fmt.Errorf("don't find event-options destinations with id '%v' (id must be <name>)", d.Id())
	}
	eventoptionsDestinationOptions, err := readEventoptionsDestination(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillEventoptionsDestinationData(d, eventoptionsDestinationOptions)

	result[0] = d

	return result, nil
}

func checkEventoptionsDestinationExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	eventoptionsDestinationConfig, err :=
		sess.command("show configuration event-options destinations \""+name+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if eventoptionsDestinationConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setEventoptionsDestination(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := "set event-options destinations \"" + d.Get("name").(string) + "\" "

	archiveSiteURLList := make([]string, 0)
	for _, v := range d.Get("archive_site").([]interface{}) {
		archiveSite := v.(map[string]interface{})
		if stringInSlice(archiveSite["url"].(string), archiveSiteURLList) {
			return fmt.Errorf("multiple archive_site blocks with the same url")
		}
		archiveSiteURLList = append(archiveSiteURLList, archiveSite["url"].(string))
		configSet = append(configSet, setPrefix+"archive-sites \""+archiveSite["url"].(string)+"\"")
		if v2 := archiveSite["password"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"archive-sites \""+archiveSite["url"].(string)+"\" password \""+v2+"\"")
		}
	}
	if v := d.Get("transfer_delay").(int); v != -1 {
		configSet = append(configSet, setPrefix+"transfer-delay "+strconv.Itoa(v))
	}

	return sess.configSet(configSet, jnprSess)
}

func readEventoptionsDestination(
	destination string, m interface{}, jnprSess *NetconfObject) (eventoptionsDestinationOptions, error) {
	sess := m.(*Session)
	var confRead eventoptionsDestinationOptions
	confRead.transferDelay = -1 // default value

	eventoptionsDestinationConfig, err := sess.command("show configuration event-options destinations \""+
		destination+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if eventoptionsDestinationConfig != emptyWord {
		confRead.name = destination
		for _, item := range strings.Split(eventoptionsDestinationConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "archive-sites "):
				itemTrimSplit := strings.Split(itemTrim, " ")
				if len(itemTrimSplit) > 2 {
					password, err := jdecode.Decode(strings.Trim(itemTrimSplit[3], "\""))
					if err != nil {
						return confRead, fmt.Errorf("failed to decode secret : %w", err)
					}
					confRead.archiveSite = append(confRead.archiveSite, map[string]interface{}{
						"url":      strings.Trim(itemTrimSplit[1], "\""),
						"password": password,
					})
				} else {
					confRead.archiveSite = append(confRead.archiveSite, map[string]interface{}{
						"url":      strings.Trim(itemTrimSplit[1], "\""),
						"password": "",
					})
				}
			case strings.HasPrefix(itemTrim, "transfer-delay "):
				var err error
				confRead.transferDelay, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "transfer-delay "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delEventoptionsDestination(destination string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete event-options destinations \""+destination+"\"")

	return sess.configSet(configSet, jnprSess)
}

func fillEventoptionsDestinationData(
	d *schema.ResourceData, eventoptionsDestinationOptions eventoptionsDestinationOptions) {
	if tfErr := d.Set("name", eventoptionsDestinationOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("archive_site", eventoptionsDestinationOptions.archiveSite); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("transfer_delay", eventoptionsDestinationOptions.transferDelay); tfErr != nil {
		panic(tfErr)
	}
}
