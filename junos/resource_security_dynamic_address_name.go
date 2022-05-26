package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
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
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSecurityDynamicAddressName(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if !checkCompatibilitySecurity(junSess) {
		return diag.FromErr(fmt.Errorf("security dynamic-address address-name "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityDynamicAddressNameExists, err := checkSecurityDynamicAddressNamesExists(d.Get("name").(string), sess, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityDynamicAddressNameExists {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security dynamic-address address-name %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityDynamicAddressName(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_dynamic_address_name", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityDynamicAddressNameExists, err = checkSecurityDynamicAddressNamesExists(d.Get("name").(string), sess, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityDynamicAddressNameExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security dynamic-address address-name %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityDynamicAddressNameReadWJunSess(d, sess, junSess)...)
}

func resourceSecurityDynamicAddressNameRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)

	return resourceSecurityDynamicAddressNameReadWJunSess(d, sess, junSess)
}

func resourceSecurityDynamicAddressNameReadWJunSess(d *schema.ResourceData, sess *Session, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	dynamicAddressNameOptions, err := readSecurityDynamicAddressName(d.Get("name").(string), sess, junSess)
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
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delSecurityDynamicAddressName(d.Get("name").(string), sess, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityDynamicAddressName(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityDynamicAddressName(d.Get("name").(string), sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityDynamicAddressName(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_dynamic_address_name", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityDynamicAddressNameReadWJunSess(d, sess, junSess)...)
}

func resourceSecurityDynamicAddressNameDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delSecurityDynamicAddressName(d.Get("name").(string), sess, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityDynamicAddressName(d.Get("name").(string), sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_dynamic_address_name", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityDynamicAddressNameImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	securityDynamicAddressNameExists, err := checkSecurityDynamicAddressNamesExists(d.Id(), sess, junSess)
	if err != nil {
		return nil, err
	}
	if !securityDynamicAddressNameExists {
		return nil, fmt.Errorf("security dynamic-address address-name with id '%v' (id must be <name>)", d.Id())
	}
	dynamicAddressNameOptions, err := readSecurityDynamicAddressName(d.Id(), sess, junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityDynamicAddressNameData(d, dynamicAddressNameOptions)

	result[0] = d

	return result, nil
}

func checkSecurityDynamicAddressNamesExists(name string, sess *Session, junSess *junosSession) (bool, error) {
	showConfig, err := sess.command(cmdShowConfig+
		"security dynamic-address address-name "+name+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityDynamicAddressName(d *schema.ResourceData, sess *Session, junSess *junosSession) error {
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
			if bchk.StringInSlice(property["name"].(string), propertyNameList) {
				return fmt.Errorf("multiple blocks property with the same name %s", property["name"].(string))
			}
			propertyNameList = append(propertyNameList, property["name"].(string))
			for _, str := range property["string"].([]interface{}) {
				configSet = append(configSet, setPrefixProfileCategory+"property "+
					"\""+property["name"].(string)+"\" string \""+str.(string)+"\"")
			}
		}
	}

	return sess.configSet(configSet, junSess)
}

func readSecurityDynamicAddressName(name string, sess *Session, junSess *junosSession,
) (dynamicAddressNameOptions, error) {
	var confRead dynamicAddressNameOptions

	showConfig, err := sess.command(cmdShowConfig+
		"security dynamic-address address-name "+name+pipeDisplaySetRelative, junSess)
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
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "profile feed-name "):
				confRead.profileFeedName = strings.TrimPrefix(itemTrim, "profile feed-name ")
			case strings.HasPrefix(itemTrim, "profile category "):
				itemTrimProfileCategorySplit := strings.Split(strings.TrimPrefix(itemTrim, "profile category "), " ")
				if len(confRead.profileCategory) == 0 {
					confRead.profileCategory = append(confRead.profileCategory, map[string]interface{}{
						"name":     itemTrimProfileCategorySplit[0],
						"feed":     "",
						"property": make([]map[string]interface{}, 0),
					})
				}
				itemTrimProfileCategory := strings.TrimPrefix(itemTrim, "profile category "+itemTrimProfileCategorySplit[0]+" ")
				switch {
				case strings.HasPrefix(itemTrimProfileCategory, "feed "):
					confRead.profileCategory[0]["feed"] = strings.TrimPrefix(itemTrimProfileCategory, "feed ")
				case strings.HasPrefix(itemTrimProfileCategory, "property "):
					itemTrimPropertySplit := strings.Split(strings.TrimPrefix(itemTrimProfileCategory, "property "), " ")
					property := map[string]interface{}{
						"name":   strings.Trim(itemTrimPropertySplit[0], "\""),
						"string": make([]string, 0),
					}
					confRead.profileCategory[0]["property"] = copyAndRemoveItemMapList(
						"name",
						property,
						confRead.profileCategory[0]["property"].([]map[string]interface{}),
					)
					property["string"] = append(
						property["string"].([]string),
						strings.Trim(strings.TrimPrefix(itemTrimProfileCategory, "property "+itemTrimPropertySplit[0]+" string "), "\""),
					)
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

func delSecurityDynamicAddressName(name string, sess *Session, junSess *junosSession) error {
	configSet := []string{"delete security dynamic-address address-name " + name}

	return sess.configSet(configSet, junSess)
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
