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
		CreateWithoutTimeout: resourceEventoptionsGenerateEventCreate,
		ReadWithoutTimeout:   resourceEventoptionsGenerateEventRead,
		UpdateWithoutTimeout: resourceEventoptionsGenerateEventUpdate,
		DeleteWithoutTimeout: resourceEventoptionsGenerateEventDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceEventoptionsGenerateEventImport,
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

func resourceEventoptionsGenerateEventCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setEventoptionsGenerateEvent(d, clt, nil); err != nil {
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
	eventoptionsGenerateEventExists, err := checkEventoptionsGenerateEventExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if eventoptionsGenerateEventExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("event-options generate-event %v already exists", d.Get("name").(string)))...)
	}

	if err := setEventoptionsGenerateEvent(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_eventoptions_generate_event", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	eventoptionsGenerateEventExists, err = checkEventoptionsGenerateEventExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if eventoptionsGenerateEventExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("event-options generate-event %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceEventoptionsGenerateEventReadWJunSess(d, clt, junSess)...)
}

func resourceEventoptionsGenerateEventRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceEventoptionsGenerateEventReadWJunSess(d, clt, junSess)
}

func resourceEventoptionsGenerateEventReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	eventoptionsGenerateEventOptions, err := readEventoptionsGenerateEvent(d.Get("name").(string), clt, junSess)
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

func resourceEventoptionsGenerateEventUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delEventoptionsGenerateEvent(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setEventoptionsGenerateEvent(d, clt, nil); err != nil {
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
	if err := delEventoptionsGenerateEvent(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setEventoptionsGenerateEvent(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_eventoptions_generate_event", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceEventoptionsGenerateEventReadWJunSess(d, clt, junSess)...)
}

func resourceEventoptionsGenerateEventDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delEventoptionsGenerateEvent(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delEventoptionsGenerateEvent(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_eventoptions_generate_event", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceEventoptionsGenerateEventImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)

	eventoptionsGenerateEventExists, err := checkEventoptionsGenerateEventExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !eventoptionsGenerateEventExists {
		return nil, fmt.Errorf("don't find event-options generate-event with id '%v' (id must be <name>)", d.Id())
	}
	eventoptionsGenerateEventOptions, err := readEventoptionsGenerateEvent(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillEventoptionsGenerateEventData(d, eventoptionsGenerateEventOptions)

	result[0] = d

	return result, nil
}

func checkEventoptionsGenerateEventExists(name string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+"event-options generate-event \""+name+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setEventoptionsGenerateEvent(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
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

	return clt.configSet(configSet, junSess)
}

func readEventoptionsGenerateEvent(name string, clt *Client, junSess *junosSession,
) (eventoptionsGenerateEventOptions, error) {
	var confRead eventoptionsGenerateEventOptions

	showConfig, err := clt.command(cmdShowConfig+
		"event-options generate-event \""+name+"\""+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "time-interval "):
				var err error
				confRead.timeInterval, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "time-interval "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
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

func delEventoptionsGenerateEvent(event string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete event-options generate-event \""+event+"\"")

	return clt.configSet(configSet, junSess)
}

func fillEventoptionsGenerateEventData(
	d *schema.ResourceData, eventoptionsGenerateEventOptions eventoptionsGenerateEventOptions,
) {
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
