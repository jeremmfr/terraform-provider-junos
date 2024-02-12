package providersdk

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type lldpMedInterfaceOptions struct {
	disable  bool
	enable   bool
	name     string
	location []map[string]interface{}
}

func resourceLldpMedInterface() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceLldpMedInterfaceCreate,
		ReadWithoutTimeout:   resourceLldpMedInterfaceRead,
		UpdateWithoutTimeout: resourceLldpMedInterfaceUpdate,
		DeleteWithoutTimeout: resourceLldpMedInterfaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLldpMedInterfaceImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if strings.Count(value, ".") > 0 {
						errors = append(errors, fmt.Errorf(
							"%q in %q cannot have a dot", value, k))
					}

					return
				},
			},
			"disable": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"enable"},
			},
			"enable": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"disable"},
			},
			"location": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"civic_based_ca_type": {
							Type:         schema.TypeList,
							Optional:     true,
							RequiredWith: []string{"location.0.civic_based_country_code"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ca_type": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 255),
									},
									"ca_value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"civic_based_country_code": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"location.0.co_ordinate_latitude",
								"location.0.co_ordinate_longitude",
								"location.0.elin",
							},
							ValidateFunc: validation.StringLenBetween(2, 2),
						},
						"civic_based_what": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							RequiredWith: []string{"location.0.civic_based_country_code"},
							ValidateFunc: validation.IntBetween(0, 2),
						},
						"co_ordinate_latitude": {
							Type:     schema.TypeInt,
							Optional: true,
							ConflictsWith: []string{
								"location.0.civic_based_country_code",
								"location.0.elin",
							},
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 360),
						},
						"co_ordinate_longitude": {
							Type:     schema.TypeInt,
							Optional: true,
							ConflictsWith: []string{
								"location.0.civic_based_country_code",
								"location.0.elin",
							},
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 360),
						},
						"elin": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"location.0.civic_based_country_code",
								"location.0.co_ordinate_latitude",
								"location.0.co_ordinate_longitude",
							},
						},
					},
				},
			},
		},
	}
}

func resourceLldpMedInterfaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setLldpMedInterface(d, junSess); err != nil {
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
	lldpMedInterfaceExists, err := checkLldpMedInterfaceExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if lldpMedInterfaceExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"protocols lldp-med interface %v already exists", d.Get("name").(string)))...)
	}

	if err := setLldpMedInterface(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_lldpmed_interface")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	lldpMedInterfaceExists, err = checkLldpMedInterfaceExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if lldpMedInterfaceExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("protocols lldp-med interface %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceLldpMedInterfaceReadWJunSess(d, junSess)...)
}

func resourceLldpMedInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceLldpMedInterfaceReadWJunSess(d, junSess)
}

func resourceLldpMedInterfaceReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	lldpMedInterfaceOptions, err := readLldpMedInterface(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if lldpMedInterfaceOptions.name == "" {
		d.SetId("")
	} else {
		fillLldpMedInterfaceData(d, lldpMedInterfaceOptions)
	}

	return nil
}

func resourceLldpMedInterfaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delLldpMedInterface(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setLldpMedInterface(d, junSess); err != nil {
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
	if err := delLldpMedInterface(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setLldpMedInterface(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_lldpmed_interface")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceLldpMedInterfaceReadWJunSess(d, junSess)...)
}

func resourceLldpMedInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delLldpMedInterface(d.Get("name").(string), junSess); err != nil {
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
	if err := delLldpMedInterface(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_lldpmed_interface")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceLldpMedInterfaceImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	lldpMedInterfaceExists, err := checkLldpMedInterfaceExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !lldpMedInterfaceExists {
		return nil, fmt.Errorf("don't find protocols lldp-med interface with id '%v' (id must be <name>)", d.Id())
	}
	lldpMedInterfaceOptions, err := readLldpMedInterface(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillLldpMedInterfaceData(d, lldpMedInterfaceOptions)

	result[0] = d

	return result, nil
}

func checkLldpMedInterfaceExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "protocols lldp-med interface " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setLldpMedInterface(d *schema.ResourceData, junSess *junos.Session) error {
	setPrefix := "set protocols lldp-med interface " + d.Get("name").(string) + " "
	configSet := make([]string, 0)

	configSet = append(configSet, setPrefix)
	if d.Get("disable").(bool) {
		configSet = append(configSet, setPrefix+"disable")
	}
	if d.Get("enable").(bool) {
		configSet = append(configSet, setPrefix+"enable")
	}
	for _, mLocation := range d.Get("location").([]interface{}) {
		location := mLocation.(map[string]interface{})
		setPrefixLocation := setPrefix + "location "
		configSet = append(configSet, setPrefixLocation)
		if cCode := location["civic_based_country_code"].(string); cCode != "" {
			configSet = append(configSet, setPrefixLocation+"civic-based country-code "+cCode)
			civicBasedCaTypeList := make([]int, 0)
			for _, mCaT := range location["civic_based_ca_type"].([]interface{}) {
				caType := mCaT.(map[string]interface{})
				if slices.Contains(civicBasedCaTypeList, caType["ca_type"].(int)) {
					return fmt.Errorf("multiple blocks civic_based_ca_type with the same ca_type '%d'", caType["ca_type"].(int))
				}
				civicBasedCaTypeList = append(civicBasedCaTypeList, caType["ca_type"].(int))
				configSet = append(configSet, setPrefixLocation+"civic-based ca-type "+strconv.Itoa(caType["ca_type"].(int)))
				if v := caType["ca_value"].(string); v != "" {
					configSet = append(configSet, setPrefixLocation+
						"civic-based ca-type "+strconv.Itoa(caType["ca_type"].(int))+" ca-value \""+v+"\"")
				}
			}
			if v := location["civic_based_what"].(int); v != -1 {
				configSet = append(configSet, setPrefixLocation+"civic-based what "+strconv.Itoa(v))
			}
		} else if len(location["civic_based_ca_type"].([]interface{})) > 0 ||
			location["civic_based_what"].(int) != -1 {
			return errors.New("civic_based_country_code need to be set with civic_based_ca_type and civic_based_what")
		}
		if v := location["co_ordinate_latitude"].(int); v != -1 {
			configSet = append(configSet, setPrefixLocation+"co-ordinate lattitude "+strconv.Itoa(v)) //nolint:misspell
			configSet = append(configSet, setPrefixLocation+"co-ordinate latitude "+strconv.Itoa(v))
		}
		if v := location["co_ordinate_longitude"].(int); v != -1 {
			configSet = append(configSet, setPrefixLocation+"co-ordinate longitude "+strconv.Itoa(v))
		}
		if v := location["elin"].(string); v != "" {
			configSet = append(configSet, setPrefixLocation+"elin \""+v+"\"")
		}
	}

	return junSess.ConfigSet(configSet)
}

func readLldpMedInterface(name string, junSess *junos.Session,
) (confRead lldpMedInterfaceOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"protocols lldp-med interface " + name + junos.PipeDisplaySetRelative)
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
			case itemTrim == junos.DisableW:
				confRead.disable = true
			case itemTrim == "enable":
				confRead.enable = true
			case balt.CutPrefixInString(&itemTrim, "location"):
				if len(confRead.location) == 0 {
					confRead.location = append(confRead.location, map[string]interface{}{
						"civic_based_ca_type":      make([]map[string]interface{}, 0),
						"civic_based_country_code": "",
						"civic_based_what":         -1,
						"co_ordinate_latitude":     -1,
						"co_ordinate_longitude":    -1,
						"elin":                     "",
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " civic-based ca-type "):
					itemTrimFields := strings.Split(itemTrim, " ")
					caType, err := strconv.Atoi(itemTrimFields[0])
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
					switch len(itemTrimFields) {
					case 1: // <ca_type>
						confRead.location[0]["civic_based_ca_type"] = append(
							confRead.location[0]["civic_based_ca_type"].([]map[string]interface{}),
							map[string]interface{}{
								"ca_type":  caType,
								"ca_value": "",
							})
					case 3: // <ca_type> ca-value <ca_value>
						confRead.location[0]["civic_based_ca_type"] = append(
							confRead.location[0]["civic_based_ca_type"].([]map[string]interface{}),
							map[string]interface{}{
								"ca_type":  caType,
								"ca_value": itemTrimFields[2],
							})
					default:
						return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "civic-based ca-type", itemTrim)
					}
				case balt.CutPrefixInString(&itemTrim, " civic-based country-code "):
					confRead.location[0]["civic_based_country_code"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, " civic-based what "):
					confRead.location[0]["civic_based_what"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, " co-ordinate lattitude "): //nolint:misspell
					confRead.location[0]["co_ordinate_latitude"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, " co-ordinate latitude "):
					confRead.location[0]["co_ordinate_latitude"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, " co-ordinate longitude "):
					confRead.location[0]["co_ordinate_longitude"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, " elin "):
					confRead.location[0]["elin"] = strings.Trim(itemTrim, "\"")
				}
			}
		}
	}

	return confRead, nil
}

func delLldpMedInterface(name string, junSess *junos.Session) error {
	configSet := []string{"delete protocols lldp-med interface " + name}

	return junSess.ConfigSet(configSet)
}

func fillLldpMedInterfaceData(d *schema.ResourceData, lldpMedInterfaceOptions lldpMedInterfaceOptions) {
	if tfErr := d.Set("name", lldpMedInterfaceOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("disable", lldpMedInterfaceOptions.disable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("enable", lldpMedInterfaceOptions.enable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("location", lldpMedInterfaceOptions.location); tfErr != nil {
		panic(tfErr)
	}
}
