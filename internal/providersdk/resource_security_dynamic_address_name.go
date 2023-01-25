package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
		if err := setSecurityDynamicAddressName(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if !junos.CheckCompatibilitySecurity(junSess) {
		return diag.FromErr(fmt.Errorf("security dynamic-address address-name "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityDynamicAddressNameExists, err := checkSecurityDynamicAddressNamesExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityDynamicAddressNameExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security dynamic-address address-name %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityDynamicAddressName(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_security_dynamic_address_name", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityDynamicAddressNameExists, err = checkSecurityDynamicAddressNamesExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityDynamicAddressNameExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security dynamic-address address-name %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityDynamicAddressNameReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityDynamicAddressNameRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceSecurityDynamicAddressNameReadWJunSess(d, clt, junSess)
}

func resourceSecurityDynamicAddressNameReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	dynamicAddressNameOptions, err := readSecurityDynamicAddressName(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
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
		if err := delSecurityDynamicAddressName(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityDynamicAddressName(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityDynamicAddressName(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityDynamicAddressName(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_security_dynamic_address_name", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityDynamicAddressNameReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityDynamicAddressNameDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delSecurityDynamicAddressName(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityDynamicAddressName(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_security_dynamic_address_name", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

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
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)
	securityDynamicAddressNameExists, err := checkSecurityDynamicAddressNamesExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !securityDynamicAddressNameExists {
		return nil, fmt.Errorf("security dynamic-address address-name with id '%v' (id must be <name>)", d.Id())
	}
	dynamicAddressNameOptions, err := readSecurityDynamicAddressName(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityDynamicAddressNameData(d, dynamicAddressNameOptions)

	result[0] = d

	return result, nil
}

func checkSecurityDynamicAddressNamesExists(name string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"security dynamic-address address-name "+name+junos.PipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityDynamicAddressName(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
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
			if bchk.InSlice(property["name"].(string), propertyNameList) {
				return fmt.Errorf("multiple blocks property with the same name %s", property["name"].(string))
			}
			propertyNameList = append(propertyNameList, property["name"].(string))
			for _, str := range property["string"].([]interface{}) {
				configSet = append(configSet, setPrefixProfileCategory+"property "+
					"\""+property["name"].(string)+"\" string \""+str.(string)+"\"")
			}
		}
	}

	return clt.ConfigSet(configSet, junSess)
}

func readSecurityDynamicAddressName(name string, clt *junos.Client, junSess *junos.Session,
) (confRead dynamicAddressNameOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"security dynamic-address address-name "+name+junos.PipeDisplaySetRelative, junSess)
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

func delSecurityDynamicAddressName(name string, clt *junos.Client, junSess *junos.Session) error {
	configSet := []string{"delete security dynamic-address address-name " + name}

	return clt.ConfigSet(configSet, junSess)
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