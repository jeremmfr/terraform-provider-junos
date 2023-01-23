package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := setServicesUserIdentDeviceIdentityProfile(d, clt, nil); err != nil {
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
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	svcUserIdentDevIdentProfileExists, err := checkServicesUserIdentDeviceIdentityProfileExists(
		d.Get("name").(string),
		clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if svcUserIdentDevIdentProfileExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf(
				"services user-identification device-information end-user-profile %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesUserIdentDeviceIdentityProfile(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_services_user_identification_device_identity_profile", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

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
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceServicesUserIdentDeviceIdentityProfileReadWJunSess(d, clt, junSess)
}

func resourceServicesUserIdentDeviceIdentityProfileReadWJunSess(
	d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
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
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delServicesUserIdentDeviceIdentityProfile(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setServicesUserIdentDeviceIdentityProfile(d, clt, nil); err != nil {
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
	if err := delServicesUserIdentDeviceIdentityProfile(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesUserIdentDeviceIdentityProfile(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_services_user_identification_device_identity_profile", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesUserIdentDeviceIdentityProfileReadWJunSess(d, clt, junSess)...)
}

func resourceServicesUserIdentDeviceIdentityProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delServicesUserIdentDeviceIdentityProfile(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delServicesUserIdentDeviceIdentityProfile(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_services_user_identification_device_identity_profile", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesUserIdentDeviceIdentityProfileImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
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

func checkServicesUserIdentDeviceIdentityProfileExists(profile string, clt *junos.Client, junSess *junos.Session,
) (bool, error) {
	showConfig, err := clt.Command(
		junos.CmdShowConfig+
			"services user-identification device-information end-user-profile profile-name "+profile+junos.PipeDisplaySet,
		junSess,
	)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setServicesUserIdentDeviceIdentityProfile(
	d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) error {
	configSet := make([]string, 0)

	setPrefix := "set services user-identification device-information end-user-profile profile-name " +
		d.Get("name").(string) + " "
	configSet = append(configSet, setPrefix+"domain-name "+d.Get("domain").(string))
	attributeNameList := make([]string, 0)
	for _, v := range d.Get("attribute").([]interface{}) {
		attribute := v.(map[string]interface{})
		if bchk.InSlice(attribute["name"].(string), attributeNameList) {
			return fmt.Errorf("multiple blocks attribute with the same name %s", attribute["name"].(string))
		}
		attributeNameList = append(attributeNameList, attribute["name"].(string))
		for _, v2 := range sortSetOfString(attribute["value"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefix+"attribute "+attribute["name"].(string)+
				" string \""+v2+"\"")
		}
	}

	return clt.ConfigSet(configSet, junSess)
}

func readServicesUserIdentDeviceIdentityProfile(profile string, clt *junos.Client, junSess *junos.Session,
) (confRead svcUserIdentDevIdentProfileOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"services user-identification device-information end-user-profile"+
		" profile-name "+profile+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = profile
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "domain-name "):
				confRead.domain = itemTrim
			case balt.CutPrefixInString(&itemTrim, "attribute "):
				itemTrimFields := strings.Split(itemTrim, " ")
				attribute := map[string]interface{}{
					"name":  itemTrimFields[0],
					"value": make([]string, 0),
				}
				confRead.attribute = copyAndRemoveItemMapList("name", attribute, confRead.attribute)
				if balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" string ") {
					attribute["value"] = append(attribute["value"].([]string), strings.Trim(itemTrim, "\""))
				}
				confRead.attribute = append(confRead.attribute, attribute)
			}
		}
	}

	return confRead, nil
}

func delServicesUserIdentDeviceIdentityProfile(profile string, clt *junos.Client, junSess *junos.Session) error {
	configSet := []string{
		"delete services user-identification device-information end-user-profile profile-name " + profile,
	}

	return clt.ConfigSet(configSet, junSess)
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
