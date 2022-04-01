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
)

type lldpMedInterfaceOptions struct {
	disable  bool
	enable   bool
	name     string
	location []map[string]interface{}
}

func resourceLldpMedInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLldpMedInterfaceCreate,
		ReadContext:   resourceLldpMedInterfaceRead,
		UpdateContext: resourceLldpMedInterfaceUpdate,
		DeleteContext: resourceLldpMedInterfaceDelete,
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
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setLldpMedInterface(d, m, nil); err != nil {
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
	lldpMedInterfaceExists, err := checkLldpMedInterfaceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if lldpMedInterfaceExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"protocols lldp-med interface %v already exists", d.Get("name").(string)))...)
	}

	if err := setLldpMedInterface(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_lldpmed_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	lldpMedInterfaceExists, err = checkLldpMedInterfaceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if lldpMedInterfaceExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("protocols lldp-med interface %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceLldpMedInterfaceReadWJnprSess(d, m, jnprSess)...)
}

func resourceLldpMedInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceLldpMedInterfaceReadWJnprSess(d, m, jnprSess)
}

func resourceLldpMedInterfaceReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
) diag.Diagnostics {
	mutex.Lock()
	lldpMedInterfaceOptions, err := readLldpMedInterface(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
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
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delLldpMedInterface(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setLldpMedInterface(d, m, nil); err != nil {
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
	if err := delLldpMedInterface(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setLldpMedInterface(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_lldpmed_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceLldpMedInterfaceReadWJnprSess(d, m, jnprSess)...)
}

func resourceLldpMedInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delLldpMedInterface(d.Get("name").(string), m, nil); err != nil {
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
	if err := delLldpMedInterface(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_lldpmed_interface", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceLldpMedInterfaceImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	lldpMedInterfaceExists, err := checkLldpMedInterfaceExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !lldpMedInterfaceExists {
		return nil, fmt.Errorf("don't find protocols lldp-med interface with id '%v' (id must be <name>)", d.Id())
	}
	lldpMedInterfaceOptions, err := readLldpMedInterface(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillLldpMedInterfaceData(d, lldpMedInterfaceOptions)

	result[0] = d

	return result, nil
}

func checkLldpMedInterfaceExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command(
		cmdShowConfig+"protocols lldp-med interface "+name+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setLldpMedInterface(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

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
				if bchk.IntInSlice(caType["ca_type"].(int), civicBasedCaTypeList) {
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
			return fmt.Errorf("civic_based_country_code need to be set with civic_based_ca_type and civic_based_what")
		}
		if v := location["co_ordinate_latitude"].(int); v != -1 {
			configSet = append(configSet, setPrefixLocation+"co-ordinate lattitude "+strconv.Itoa(v)) // nolint: misspell
			configSet = append(configSet, setPrefixLocation+"co-ordinate latitude "+strconv.Itoa(v))
		}
		if v := location["co_ordinate_longitude"].(int); v != -1 {
			configSet = append(configSet, setPrefixLocation+"co-ordinate longitude "+strconv.Itoa(v))
		}
		if v := location["elin"].(string); v != "" {
			configSet = append(configSet, setPrefixLocation+"elin \""+v+"\"")
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func readLldpMedInterface(name string, m interface{}, jnprSess *NetconfObject,
) (lldpMedInterfaceOptions, error) {
	sess := m.(*Session)
	var confRead lldpMedInterfaceOptions

	showConfig, err := sess.command(
		cmdShowConfig+"protocols lldp-med interface "+name+pipeDisplaySetRelative, jnprSess)
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
			case itemTrim == disableW:
				confRead.disable = true
			case itemTrim == "enable":
				confRead.enable = true
			case strings.HasPrefix(itemTrim, "location"):
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
				itemTrimLocation := strings.TrimPrefix(itemTrim, "location ")
				switch {
				case strings.HasPrefix(itemTrimLocation, "civic-based ca-type "):
					itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrimLocation, "civic-based ca-type "), " ")
					if len(itemTrimSplit) > 0 {
						caType, err := strconv.Atoi(itemTrimSplit[0])
						if err != nil {
							return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
						}
						switch {
						case len(itemTrimSplit) == 1:
							confRead.location[0]["civic_based_ca_type"] = append(
								confRead.location[0]["civic_based_ca_type"].([]map[string]interface{}),
								map[string]interface{}{
									"ca_type":  caType,
									"ca_value": "",
								})
						case len(itemTrimSplit) == 3:
							confRead.location[0]["civic_based_ca_type"] = append(
								confRead.location[0]["civic_based_ca_type"].([]map[string]interface{}),
								map[string]interface{}{
									"ca_type":  caType,
									"ca_value": itemTrimSplit[2],
								})
						default:
							return confRead, fmt.Errorf("can't find ca-type and ca-value in %s", itemTrim)
						}
					} else {
						return confRead, fmt.Errorf("can't find ca-type and ca-value in %s", itemTrim)
					}
				case strings.HasPrefix(itemTrimLocation, "civic-based country-code "):
					confRead.location[0]["civic_based_country_code"] = strings.TrimPrefix(
						itemTrimLocation, "civic-based country-code ")
				case strings.HasPrefix(itemTrimLocation, "civic-based what "):
					var err error
					confRead.location[0]["civic_based_what"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrimLocation, "civic-based what "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrimLocation, "co-ordinate lattitude "): // nolint: misspell
					var err error
					confRead.location[0]["co_ordinate_latitude"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrimLocation, "co-ordinate lattitude ")) // nolint: misspell
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrimLocation, "co-ordinate latitude "):
					var err error
					confRead.location[0]["co_ordinate_latitude"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrimLocation, "co-ordinate latitude "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrimLocation, "co-ordinate longitude "):
					var err error
					confRead.location[0]["co_ordinate_longitude"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrimLocation, "co-ordinate longitude "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrimLocation, "elin "):
					confRead.location[0]["elin"] = strings.Trim(strings.TrimPrefix(itemTrimLocation, "elin "), "\"")
				}
			}
		}
	}

	return confRead, nil
}

func delLldpMedInterface(name string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	configSet := []string{"delete protocols lldp-med interface " + name}

	return sess.configSet(configSet, jnprSess)
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
