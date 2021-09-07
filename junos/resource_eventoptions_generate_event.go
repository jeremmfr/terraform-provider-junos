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
)

type eventoptionsGenerateEventOptions struct {
	noDrift      bool
	timeInterval int
	name         string
	timeOfDay    string
}

func resourceEventoptionsGenerateEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventoptionsGenerateEventCreate,
		ReadContext:   resourceEventoptionsGenerateEventRead,
		UpdateContext: resourceEventoptionsGenerateEventUpdate,
		DeleteContext: resourceEventoptionsGenerateEventDelete,
		Importer: &schema.ResourceImporter{
			State: resourceEventoptionsGenerateEventImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"time_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(60, 2592000),
				ExactlyOneOf: []string{"time_interval", "time_of_day"},
			},
			"time_of_day": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`), "must be in the format 'HH:MM:SS'"),
				ExactlyOneOf: []string{"time_interval", "time_of_day"},
			},
			"no_drift": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceEventoptionsGenerateEventCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setEventoptionsGenerateEvent(d, m, nil); err != nil {
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
	eventoptionsGenerateEventExists, err := checkEventoptionsGenerateEventExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if eventoptionsGenerateEventExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("event-options generate-event %v already exists", d.Get("name").(string)))...)
	}

	if err := setEventoptionsGenerateEvent(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_eventoptions_generate_event", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	eventoptionsGenerateEventExists, err = checkEventoptionsGenerateEventExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if eventoptionsGenerateEventExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("event-options generate-event %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceEventoptionsGenerateEventReadWJnprSess(d, m, jnprSess)...)
}

func resourceEventoptionsGenerateEventRead(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceEventoptionsGenerateEventReadWJnprSess(d, m, jnprSess)
}

func resourceEventoptionsGenerateEventReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	eventoptionsGenerateEventOptions, err := readEventoptionsGenerateEvent(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if eventoptionsGenerateEventOptions.name == "" {
		d.SetId("")
	} else {
		fillEventoptionsGenerateEventData(d, eventoptionsGenerateEventOptions)
	}

	return nil
}

func resourceEventoptionsGenerateEventUpdate(ctx context.Context,
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
	if err := delEventoptionsGenerateEvent(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setEventoptionsGenerateEvent(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_eventoptions_generate_event", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceEventoptionsGenerateEventReadWJnprSess(d, m, jnprSess)...)
}

func resourceEventoptionsGenerateEventDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delEventoptionsGenerateEvent(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_eventoptions_generate_event", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceEventoptionsGenerateEventImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	eventoptionsGenerateEventExists, err := checkEventoptionsGenerateEventExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !eventoptionsGenerateEventExists {
		return nil, fmt.Errorf("don't find event-options generate-event with id '%v' (id must be <name>)", d.Id())
	}
	eventoptionsGenerateEventOptions, err := readEventoptionsGenerateEvent(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillEventoptionsGenerateEventData(d, eventoptionsGenerateEventOptions)

	result[0] = d

	return result, nil
}

func checkEventoptionsGenerateEventExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	eventoptionsGenerateEventConfig, err :=
		sess.command("show configuration event-options generate-event \""+name+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if eventoptionsGenerateEventConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setEventoptionsGenerateEvent(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := "set event-options generate-event \"" + d.Get("name").(string) + "\" "

	if v := d.Get("time_interval").(int); v != 0 {
		configSet = append(configSet, setPrefix+"time-interval "+strconv.Itoa(v))
	}
	if v := d.Get("time_of_day").(string); v != "" {
		configSet = append(configSet, setPrefix+"time-of-day "+v)
	}
	if d.Get("no_drift").(bool) {
		configSet = append(configSet, setPrefix+"no-drift")
	}

	return sess.configSet(configSet, jnprSess)
}

func readEventoptionsGenerateEvent(
	event string, m interface{}, jnprSess *NetconfObject) (eventoptionsGenerateEventOptions, error) {
	sess := m.(*Session)
	var confRead eventoptionsGenerateEventOptions

	eventoptionsGenerateEventConfig, err := sess.command("show configuration event-options generate-event \""+
		event+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if eventoptionsGenerateEventConfig != emptyWord {
		confRead.name = event
		for _, item := range strings.Split(eventoptionsGenerateEventConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "time-interval "):
				var err error
				confRead.timeInterval, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "time-interval "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "time-of-day "):
				confRead.timeOfDay = strings.Split(strings.Trim(strings.TrimPrefix(itemTrim, "time-of-day "), "\""), " ")[0]
			case itemTrim == "no-drift":
				confRead.noDrift = true
			}
		}
	}

	return confRead, nil
}

func delEventoptionsGenerateEvent(event string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete event-options generate-event \""+event+"\"")

	return sess.configSet(configSet, jnprSess)
}

func fillEventoptionsGenerateEventData(
	d *schema.ResourceData, eventoptionsGenerateEventOptions eventoptionsGenerateEventOptions) {
	if tfErr := d.Set("name", eventoptionsGenerateEventOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("time_interval", eventoptionsGenerateEventOptions.timeInterval); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("time_of_day", eventoptionsGenerateEventOptions.timeOfDay); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_drift", eventoptionsGenerateEventOptions.noDrift); tfErr != nil {
		panic(tfErr)
	}
}
