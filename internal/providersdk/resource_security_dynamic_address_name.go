package providersdk

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type dynamicAddressNameOptions struct {
	description     string
	name            string
	profileFeedName string
	profileCategory []map[string]interface{}
}

func resourceSecurityDynamicAddressName() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityDynamicAddressNameCreate,
		ReadWithoutTimeout:   resourceSecurityDynamicAddressNameRead,
		UpdateWithoutTimeout: resourceSecurityDynamicAddressNameUpdate,
		DeleteWithoutTimeout: resourceSecurityDynamicAddressNameDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityDynamicAddressNameImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"profile_feed_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ExactlyOneOf:     []string{"profile_feed_name", "profile_category"},
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"profile_category": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"profile_feed_name", "profile_category"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
						},
						"feed": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"property": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 3,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringDoesNotContainAny(" "),
									},
									"string": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSecurityDynamicAddressNameCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityDynamicAddressName(d, junSess); err != nil {
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
	if !junSess.CheckCompatibilitySecurity() {
		return diag.FromErr(fmt.Errorf("security dynamic-address address-name "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityDynamicAddressNameExists, err := checkSecurityDynamicAddressNamesExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityDynamicAddressNameExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security dynamic-address address-name %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityDynamicAddressName(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_security_dynamic_address_name")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityDynamicAddressNameExists, err = checkSecurityDynamicAddressNamesExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityDynamicAddressNameExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security dynamic-address address-name %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityDynamicAddressNameReadWJunSess(d, junSess)...)
}

func resourceSecurityDynamicAddressNameRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityDynamicAddressNameReadWJunSess(d, junSess)
}

func resourceSecurityDynamicAddressNameReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	dynamicAddressNameOptions, err := readSecurityDynamicAddressName(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if dynamicAddressNameOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityDynamicAddressNameData(d, dynamicAddressNameOptions)
	}

	return nil
}

func resourceSecurityDynamicAddressNameUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityDynamicAddressName(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityDynamicAddressName(d, junSess); err != nil {
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
	if err := delSecurityDynamicAddressName(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityDynamicAddressName(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_security_dynamic_address_name")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityDynamicAddressNameReadWJunSess(d, junSess)...)
}

func resourceSecurityDynamicAddressNameDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityDynamicAddressName(d.Get("name").(string), junSess); err != nil {
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
	if err := delSecurityDynamicAddressName(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_security_dynamic_address_name")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityDynamicAddressNameImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	securityDynamicAddressNameExists, err := checkSecurityDynamicAddressNamesExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !securityDynamicAddressNameExists {
		return nil, fmt.Errorf("security dynamic-address address-name with id '%v' (id must be <name>)", d.Id())
	}
	dynamicAddressNameOptions, err := readSecurityDynamicAddressName(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityDynamicAddressNameData(d, dynamicAddressNameOptions)

	result[0] = d

	return result, nil
}

func checkSecurityDynamicAddressNamesExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security dynamic-address address-name " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityDynamicAddressName(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security dynamic-address address-name " + d.Get("name").(string) + " "

	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := d.Get("profile_feed_name").(string); v != "" {
		configSet = append(configSet, setPrefix+"profile feed-name "+v)
	}
	for _, pc := range d.Get("profile_category").([]interface{}) {
		profileCategory := pc.(map[string]interface{})
		setPrefixProfileCategory := setPrefix + "profile category " + profileCategory["name"].(string) + " "
		configSet = append(configSet, setPrefixProfileCategory)
		if v := profileCategory["feed"].(string); v != "" {
			configSet = append(configSet, setPrefixProfileCategory+"feed "+v)
		}
		propertyNameList := make([]string, 0)
		for _, pro := range profileCategory["property"].([]interface{}) {
			property := pro.(map[string]interface{})
			if slices.Contains(propertyNameList, property["name"].(string)) {
				return fmt.Errorf("multiple blocks property with the same name %s", property["name"].(string))
			}
			propertyNameList = append(propertyNameList, property["name"].(string))
			for _, str := range property["string"].([]interface{}) {
				configSet = append(configSet, setPrefixProfileCategory+"property "+
					"\""+property["name"].(string)+"\" string \""+str.(string)+"\"")
			}
		}
	}

	return junSess.ConfigSet(configSet)
}

func readSecurityDynamicAddressName(name string, junSess *junos.Session,
) (confRead dynamicAddressNameOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security dynamic-address address-name " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "profile feed-name "):
				confRead.profileFeedName = itemTrim
			case balt.CutPrefixInString(&itemTrim, "profile category "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(confRead.profileCategory) == 0 {
					confRead.profileCategory = append(confRead.profileCategory, map[string]interface{}{
						"name":     itemTrimFields[0],
						"feed":     "",
						"property": make([]map[string]interface{}, 0),
					})
				}
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "feed "):
					confRead.profileCategory[0]["feed"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, "property "):
					itemTrimPropertyFields := strings.Split(itemTrim, " ")
					property := map[string]interface{}{
						"name":   strings.Trim(itemTrimPropertyFields[0], "\""),
						"string": make([]string, 0),
					}
					confRead.profileCategory[0]["property"] = copyAndRemoveItemMapList(
						"name",
						property,
						confRead.profileCategory[0]["property"].([]map[string]interface{}),
					)
					if balt.CutPrefixInString(&itemTrim, itemTrimPropertyFields[0]+" string ") {
						property["string"] = append(property["string"].([]string), strings.Trim(itemTrim, "\""))
					}
					confRead.profileCategory[0]["property"] = append(
						confRead.profileCategory[0]["property"].([]map[string]interface{}),
						property,
					)
				}
			}
		}
	}

	return confRead, nil
}

func delSecurityDynamicAddressName(name string, junSess *junos.Session) error {
	configSet := []string{"delete security dynamic-address address-name " + name}

	return junSess.ConfigSet(configSet)
}

func fillSecurityDynamicAddressNameData(d *schema.ResourceData, dynamicAddressNameOptions dynamicAddressNameOptions,
) {
	if tfErr := d.Set("name", dynamicAddressNameOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", dynamicAddressNameOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("profile_feed_name", dynamicAddressNameOptions.profileFeedName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("profile_category", dynamicAddressNameOptions.profileCategory); tfErr != nil {
		panic(tfErr)
	}
}
