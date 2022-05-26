package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type svcUserIdentDevIdentProfileOptions struct {
	name      string
	domain    string
	attribute []map[string]interface{}
}

func resourceServicesUserIdentDeviceIdentityProfile() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceServicesUserIdentDeviceIdentityProfileCreate,
		ReadWithoutTimeout:   resourceServicesUserIdentDeviceIdentityProfileRead,
		UpdateWithoutTimeout: resourceServicesUserIdentDeviceIdentityProfileUpdate,
		DeleteWithoutTimeout: resourceServicesUserIdentDeviceIdentityProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServicesUserIdentDeviceIdentityProfileImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{"any"}, 64, formatDefAndDots),
			},
			"domain": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefAndDots),
			},
			"attribute": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefAndDots),
						},
						"value": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 20,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceServicesUserIdentDeviceIdentityProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setServicesUserIdentDeviceIdentityProfile(d, clt, nil); err != nil {
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
	svcUserIdentDevIdentProfileExists, err := checkServicesUserIdentDeviceIdentityProfileExists(
		d.Get("name").(string),
		clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if svcUserIdentDevIdentProfileExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf(
				"services user-identification device-information end-user-profile %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesUserIdentDeviceIdentityProfile(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_services_user_identification_device_identity_profile", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	svcUserIdentDevIdentProfileExists, err = checkServicesUserIdentDeviceIdentityProfileExists(
		d.Get("name").(string),
		clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if svcUserIdentDevIdentProfileExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"services user-identification device-information end-user-profile %v "+
				"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesUserIdentDeviceIdentityProfileReadWJunSess(d, clt, junSess)...)
}

func resourceServicesUserIdentDeviceIdentityProfileRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceServicesUserIdentDeviceIdentityProfileReadWJunSess(d, clt, junSess)
}

func resourceServicesUserIdentDeviceIdentityProfileReadWJunSess(
	d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	svcUserIdentDevIdentProfileOptions, err := readServicesUserIdentDeviceIdentityProfile(
		d.Get("name").(string),
		clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if svcUserIdentDevIdentProfileOptions.name == "" {
		d.SetId("")
	} else {
		fillServicesUserIdentDeviceIdentityProfileData(d, svcUserIdentDevIdentProfileOptions)
	}

	return nil
}

func resourceServicesUserIdentDeviceIdentityProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delServicesUserIdentDeviceIdentityProfile(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setServicesUserIdentDeviceIdentityProfile(d, clt, nil); err != nil {
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
	if err := delServicesUserIdentDeviceIdentityProfile(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesUserIdentDeviceIdentityProfile(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_services_user_identification_device_identity_profile", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesUserIdentDeviceIdentityProfileReadWJunSess(d, clt, junSess)...)
}

func resourceServicesUserIdentDeviceIdentityProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delServicesUserIdentDeviceIdentityProfile(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delServicesUserIdentDeviceIdentityProfile(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_services_user_identification_device_identity_profile", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesUserIdentDeviceIdentityProfileImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	svcUserIdentDevIdentProfileExists, err := checkServicesUserIdentDeviceIdentityProfileExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !svcUserIdentDevIdentProfileExists {
		return nil, fmt.Errorf("don't find services user-identification "+
			"device-information end-user-profile with id '%v' (id must be <name>)", d.Id())
	}
	svcUserIdentDevIdentProfileOptions, err := readServicesUserIdentDeviceIdentityProfile(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillServicesUserIdentDeviceIdentityProfileData(d, svcUserIdentDevIdentProfileOptions)

	result[0] = d

	return result, nil
}

func checkServicesUserIdentDeviceIdentityProfileExists(profile string, clt *Client, junSess *junosSession,
) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+
		"services user-identification device-information end-user-profile profile-name "+profile+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setServicesUserIdentDeviceIdentityProfile(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := "set services user-identification device-information end-user-profile profile-name " +
		d.Get("name").(string) + " "
	configSet = append(configSet, setPrefix+"domain-name "+d.Get("domain").(string))
	attributeNameList := make([]string, 0)
	for _, v := range d.Get("attribute").([]interface{}) {
		attribute := v.(map[string]interface{})
		if bchk.StringInSlice(attribute["name"].(string), attributeNameList) {
			return fmt.Errorf("multiple blocks attribute with the same name %s", attribute["name"].(string))
		}
		attributeNameList = append(attributeNameList, attribute["name"].(string))
		for _, v2 := range sortSetOfString(attribute["value"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefix+"attribute "+attribute["name"].(string)+
				" string \""+v2+"\"")
		}
	}

	return clt.configSet(configSet, junSess)
}

func readServicesUserIdentDeviceIdentityProfile(profile string, clt *Client, junSess *junosSession,
) (svcUserIdentDevIdentProfileOptions, error) {
	var confRead svcUserIdentDevIdentProfileOptions

	showConfig, err := clt.command(cmdShowConfig+
		"services user-identification device-information end-user-profile"+
		" profile-name "+profile+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = profile
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "domain-name "):
				confRead.domain = strings.TrimPrefix(itemTrim, "domain-name ")
			case strings.HasPrefix(itemTrim, "attribute "):
				itemTrimCut := strings.Split(itemTrim, " ")
				attribute := map[string]interface{}{
					"name":  itemTrimCut[1],
					"value": make([]string, 0),
				}
				confRead.attribute = copyAndRemoveItemMapList("name", attribute, confRead.attribute)
				attribute["value"] = append(attribute["value"].([]string), strings.Trim(strings.TrimPrefix(
					itemTrim, "attribute "+itemTrimCut[1]+" string "), "\""))
				confRead.attribute = append(confRead.attribute, attribute)
			}
		}
	}

	return confRead, nil
}

func delServicesUserIdentDeviceIdentityProfile(profile string, clt *Client, junSess *junosSession) error {
	configSet := []string{
		"delete services user-identification device-information end-user-profile profile-name " + profile,
	}

	return clt.configSet(configSet, junSess)
}

func fillServicesUserIdentDeviceIdentityProfileData(
	d *schema.ResourceData, svcUserIdentDevIdentProfileOptions svcUserIdentDevIdentProfileOptions,
) {
	if tfErr := d.Set("name", svcUserIdentDevIdentProfileOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("domain", svcUserIdentDevIdentProfileOptions.domain); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("attribute", svcUserIdentDevIdentProfileOptions.attribute); tfErr != nil {
		panic(tfErr)
	}
}
