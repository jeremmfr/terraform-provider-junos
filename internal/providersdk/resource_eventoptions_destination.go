package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	jdecode "github.com/jeremmfr/junosdecode"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setEventoptionsDestination(d, junSess); err != nil {
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
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	eventoptionsDestinationExists, err := checkEventoptionsDestinationExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if eventoptionsDestinationExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("event-options destinations %v already exists", d.Get("name").(string)))...)
	}

	if err := setEventoptionsDestination(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_eventoptions_destination")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	eventoptionsDestinationExists, err = checkEventoptionsDestinationExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if eventoptionsDestinationExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("event-options destinations %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceEventoptionsDestinationReadWJunSess(d, junSess)...)
}

func resourceEventoptionsDestinationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceEventoptionsDestinationReadWJunSess(d, junSess)
}

func resourceEventoptionsDestinationReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	eventoptionsDestinationOptions, err := readEventoptionsDestination(d.Get("name").(string), junSess)
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
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delEventoptionsDestination(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setEventoptionsDestination(d, junSess); err != nil {
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
	if err := delEventoptionsDestination(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setEventoptionsDestination(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_eventoptions_destination")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceEventoptionsDestinationReadWJunSess(d, junSess)...)
}

func resourceEventoptionsDestinationDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delEventoptionsDestination(d.Get("name").(string), junSess); err != nil {
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
	if err := delEventoptionsDestination(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_eventoptions_destination")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceEventoptionsDestinationImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	eventoptionsDestinationExists, err := checkEventoptionsDestinationExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !eventoptionsDestinationExists {
		return nil, fmt.Errorf("don't find event-options destinations with id '%v' (id must be <name>)", d.Id())
	}
	eventoptionsDestinationOptions, err := readEventoptionsDestination(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillEventoptionsDestinationData(d, eventoptionsDestinationOptions)

	result[0] = d

	return result, nil
}

func checkEventoptionsDestinationExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"event-options destinations \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setEventoptionsDestination(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)
	setPrefix := "set event-options destinations \"" + d.Get("name").(string) + "\" "

	archiveSiteURLList := make([]string, 0)
	for _, v := range d.Get("archive_site").([]interface{}) {
		archiveSite := v.(map[string]interface{})
		if bchk.InSlice(archiveSite["url"].(string), archiveSiteURLList) {
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

	return junSess.ConfigSet(configSet)
}

func readEventoptionsDestination(name string, junSess *junos.Session,
) (confRead eventoptionsDestinationOptions, err error) {
	// default -1
	confRead.transferDelay = -1
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"event-options destinations \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "archive-sites "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) > 2 { // <url> password <password>
					password, err := jdecode.Decode(strings.Trim(itemTrimFields[2], "\""))
					if err != nil {
						return confRead, fmt.Errorf("failed to decode secret: %w", err)
					}
					confRead.archiveSite = append(confRead.archiveSite, map[string]interface{}{
						"url":      strings.Trim(itemTrimFields[0], "\""),
						"password": password,
					})
				} else { // <url>
					confRead.archiveSite = append(confRead.archiveSite, map[string]interface{}{
						"url":      strings.Trim(itemTrimFields[0], "\""),
						"password": "",
					})
				}
			case balt.CutPrefixInString(&itemTrim, "transfer-delay "):
				confRead.transferDelay, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func delEventoptionsDestination(destination string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete event-options destinations \""+destination+"\"")

	return junSess.ConfigSet(configSet)
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
