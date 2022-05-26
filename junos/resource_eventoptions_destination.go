package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	jdecode "github.com/jeremmfr/junosdecode"
)

type eventoptionsDestinationOptions struct {
	transferDelay int
	name          string
	archiveSite   []map[string]interface{}
}

func resourceEventoptionsDestination() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceEventoptionsDestinationCreate,
		ReadWithoutTimeout:   resourceEventoptionsDestinationRead,
		UpdateWithoutTimeout: resourceEventoptionsDestinationUpdate,
		DeleteWithoutTimeout: resourceEventoptionsDestinationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceEventoptionsDestinationImport,
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

func resourceEventoptionsDestinationCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setEventoptionsDestination(d, clt, nil); err != nil {
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
	eventoptionsDestinationExists, err := checkEventoptionsDestinationExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if eventoptionsDestinationExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("event-options destinations %v already exists", d.Get("name").(string)))...)
	}

	if err := setEventoptionsDestination(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_eventoptions_destination", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	eventoptionsDestinationExists, err = checkEventoptionsDestinationExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if eventoptionsDestinationExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("event-options destinations %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceEventoptionsDestinationReadWJunSess(d, clt, junSess)...)
}

func resourceEventoptionsDestinationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceEventoptionsDestinationReadWJunSess(d, clt, junSess)
}

func resourceEventoptionsDestinationReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	eventoptionsDestinationOptions, err := readEventoptionsDestination(d.Get("name").(string), clt, junSess)
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

func resourceEventoptionsDestinationUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delEventoptionsDestination(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setEventoptionsDestination(d, clt, nil); err != nil {
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
	if err := delEventoptionsDestination(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setEventoptionsDestination(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_eventoptions_destination", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceEventoptionsDestinationReadWJunSess(d, clt, junSess)...)
}

func resourceEventoptionsDestinationDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delEventoptionsDestination(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delEventoptionsDestination(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_eventoptions_destination", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceEventoptionsDestinationImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)

	eventoptionsDestinationExists, err := checkEventoptionsDestinationExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !eventoptionsDestinationExists {
		return nil, fmt.Errorf("don't find event-options destinations with id '%v' (id must be <name>)", d.Id())
	}
	eventoptionsDestinationOptions, err := readEventoptionsDestination(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillEventoptionsDestinationData(d, eventoptionsDestinationOptions)

	result[0] = d

	return result, nil
}

func checkEventoptionsDestinationExists(name string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+"event-options destinations \""+name+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setEventoptionsDestination(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)
	setPrefix := "set event-options destinations \"" + d.Get("name").(string) + "\" "

	archiveSiteURLList := make([]string, 0)
	for _, v := range d.Get("archive_site").([]interface{}) {
		archiveSite := v.(map[string]interface{})
		if bchk.StringInSlice(archiveSite["url"].(string), archiveSiteURLList) {
			return fmt.Errorf("multiple blocks archive_site with the same url %s", archiveSite["url"].(string))
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

	return clt.configSet(configSet, junSess)
}

func readEventoptionsDestination(name string, clt *Client, junSess *junosSession,
) (eventoptionsDestinationOptions, error) {
	var confRead eventoptionsDestinationOptions
	confRead.transferDelay = -1 // default value

	showConfig, err := clt.command(cmdShowConfig+
		"event-options destinations \""+name+"\""+pipeDisplaySetRelative, junSess)
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
			case strings.HasPrefix(itemTrim, "archive-sites "):
				itemTrimSplit := strings.Split(itemTrim, " ")
				if len(itemTrimSplit) > 2 {
					password, err := jdecode.Decode(strings.Trim(itemTrimSplit[3], "\""))
					if err != nil {
						return confRead, fmt.Errorf("failed to decode secret: %w", err)
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
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delEventoptionsDestination(destination string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete event-options destinations \""+destination+"\"")

	return clt.configSet(configSet, junSess)
}

func fillEventoptionsDestinationData(
	d *schema.ResourceData, eventoptionsDestinationOptions eventoptionsDestinationOptions,
) {
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
