package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	jdecode "github.com/jeremmfr/junosdecode"
)

type snmpV3UsmUserOptions struct {
	name                   string
	engineID               string
	engineType             string
	authenticationType     string
	authenticationKey      string
	authenticationPassword string
	privacyType            string
	privacyKey             string
	privacyPassword        string
}

func resourceSnmpV3UsmUser() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSnmpV3UsmUserCreate,
		ReadWithoutTimeout:   resourceSnmpV3UsmUserRead,
		UpdateWithoutTimeout: resourceSnmpV3UsmUserUpdate,
		DeleteWithoutTimeout: resourceSnmpV3UsmUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSnmpV3UsmUserImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "local",
				ValidateFunc: validation.StringInSlice([]string{"local", "remote"}, false),
			},
			"engine_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"authentication_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"authentication_password"},
			},
			"authentication_password": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ConflictsWith:    []string{"authentication_key"},
				ValidateDiagFunc: stringLenBetweenSensitive(8, 1024),
			},
			"authentication_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "authentication-none",
				ValidateFunc: validation.StringInSlice([]string{
					"authentication-md5",
					"authentication-sha",
					"authentication-none",
				}, false),
			},
			"privacy_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"privacy_password"},
			},
			"privacy_password": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ConflictsWith:    []string{"privacy_key"},
				ValidateDiagFunc: stringLenBetweenSensitive(8, 1024),
			},
			"privacy_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "privacy-none",
				ValidateFunc: validation.StringInSlice([]string{
					"privacy-3des",
					"privacy-aes128",
					"privacy-des",
					"privacy-none",
				}, false),
			},
		},
	}
}

func resourceSnmpV3UsmUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSnmpV3UsmUser(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if d.Get("engine_type").(string) == "local" {
			d.SetId("local" + junos.IDSeparator + d.Get("name").(string))
		} else {
			d.SetId("remote" + junos.IDSeparator + d.Get("engine_id").(string) + junos.IDSeparator + d.Get("name").(string))
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
	snmpV3UsmUserExists, err := checkSnmpV3UsmUserExists(
		d.Get("name").(string),
		d.Get("engine_type").(string),
		d.Get("engine_id").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3UsmUserExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		if d.Get("engine_type").(string) != "local" {
			return append(diagWarns, diag.FromErr(fmt.Errorf("snmp v3 usm user %v in remote-engine %s already exists",
				d.Get("name").(string), d.Get("engine_id").(string)))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp v3 usm user %v in local-engine already exists",
			d.Get("name").(string)))...)
	}

	if err := setSnmpV3UsmUser(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_snmp_v3_usm_user")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	snmpV3UsmUserExists, err = checkSnmpV3UsmUserExists(
		d.Get("name").(string),
		d.Get("engine_type").(string),
		d.Get("engine_id").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3UsmUserExists {
		if d.Get("engine_type").(string) == "local" {
			d.SetId("local" + junos.IDSeparator + d.Get("name").(string))
		} else {
			d.SetId("remote" + junos.IDSeparator + d.Get("engine_id").(string) + junos.IDSeparator + d.Get("name").(string))
		}
	} else {
		if d.Get("engine_type").(string) != "local" {
			return append(diagWarns, diag.FromErr(fmt.Errorf("snmp v3 usm user %v in remote-engine %s not exists after commit "+
				"=> check your config", d.Get("name").(string), d.Get("engine_id").(string)))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp v3 usm user %v in local-engine not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSnmpV3UsmUserReadWJunSess(d, junSess)...)
}

func resourceSnmpV3UsmUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSnmpV3UsmUserReadWJunSess(d, junSess)
}

func resourceSnmpV3UsmUserReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	configSrc := snmpV3UsmUserOptions{
		name:                   d.Get("name").(string),
		engineType:             d.Get("engine_type").(string),
		engineID:               d.Get("engine_id").(string),
		authenticationPassword: d.Get("authentication_password").(string),
		authenticationType:     d.Get("authentication_type").(string),
		privacyPassword:        d.Get("privacy_password").(string),
		privacyType:            d.Get("privacy_type").(string),
	}
	junos.MutexLock()
	snmpV3UsmUserOptions, err := readSnmpV3UsmUser(configSrc, junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if snmpV3UsmUserOptions.name == "" {
		d.SetId("")
	} else {
		fillSnmpV3UsmUserData(d, snmpV3UsmUserOptions)
	}

	return nil
}

func resourceSnmpV3UsmUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSnmpV3UsmUser(
			d.Get("name").(string),
			d.Get("engine_type").(string),
			d.Get("engine_id").(string),
			junSess,
		); err != nil {
			return diag.FromErr(err)
		}
		if err := setSnmpV3UsmUser(d, junSess); err != nil {
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
	if err := delSnmpV3UsmUser(
		d.Get("name").(string),
		d.Get("engine_type").(string),
		d.Get("engine_id").(string),
		junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmpV3UsmUser(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_snmp_v3_usm_user")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpV3UsmUserReadWJunSess(d, junSess)...)
}

func resourceSnmpV3UsmUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSnmpV3UsmUser(
			d.Get("name").(string),
			d.Get("engine_type").(string),
			d.Get("engine_id").(string),
			junSess,
		); err != nil {
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
	if err := delSnmpV3UsmUser(
		d.Get("name").(string),
		d.Get("engine_type").(string),
		d.Get("engine_id").(string),
		junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_snmp_v3_usm_user")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSnmpV3UsmUserImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), junos.IDSeparator)
	var configImport snmpV3UsmUserOptions
	switch {
	case len(idSplit) == 2 && idSplit[0] == "local":
		snmpV3UsmUserExists, err := checkSnmpV3UsmUserExists(idSplit[1], idSplit[0], "", junSess)
		if err != nil {
			return nil, err
		}
		if !snmpV3UsmUserExists {
			return nil, fmt.Errorf(
				"don't find snmp v3 usm user with id '%v' (id must be local"+junos.IDSeparator+"<name>)",
				d.Id(),
			)
		}
		configImport.name = idSplit[1]
		configImport.engineType = idSplit[0]
	case len(idSplit) == 3 && idSplit[0] == "remote":
		snmpV3UsmUserExists, err := checkSnmpV3UsmUserExists(idSplit[2], idSplit[0], idSplit[1], junSess)
		if err != nil {
			return nil, err
		}
		if !snmpV3UsmUserExists {
			return nil, fmt.Errorf("don't find snmp v3 usm user with id "+
				"'%v' (id must be remote"+junos.IDSeparator+"<engine_id>"+junos.IDSeparator+"<name>)",
				d.Id(),
			)
		}
		configImport.name = idSplit[2]
		configImport.engineType = idSplit[0]
		configImport.engineID = idSplit[1]
	default:
		return nil, fmt.Errorf("can't find snmp v3 usm user with id "+
			"'%v' (id must be local"+junos.IDSeparator+"<name> or "+
			"remote"+junos.IDSeparator+"<engine_id>"+junos.IDSeparator+"<name>)",
			d.Id(),
		)
	}
	snmpV3UsmUserOptions, err := readSnmpV3UsmUser(configImport, junSess)
	if err != nil {
		return nil, err
	}
	fillSnmpV3UsmUserData(d, snmpV3UsmUserOptions)

	result[0] = d

	return result, nil
}

func checkSnmpV3UsmUserExists(name, engineType, engineID string, junSess *junos.Session,
) (bool, error) {
	if engineType == "local" {
		showConfig, err := junSess.Command(junos.CmdShowConfig +
			"snmp v3 usm local-engine user \"" + name + "\"" + junos.PipeDisplaySet)
		if err != nil {
			return false, err
		}
		if showConfig == junos.EmptyW {
			return false, nil
		}
	} else {
		showConfig, err := junSess.Command(junos.CmdShowConfig +
			"snmp v3 usm remote-engine \"" + engineID + "\" user \"" + name + "\"" + junos.PipeDisplaySet)
		if err != nil {
			return false, err
		}
		if showConfig == junos.EmptyW {
			return false, nil
		}
	}

	return true, nil
}

func setSnmpV3UsmUser(d *schema.ResourceData, junSess *junos.Session) error {
	setPrefix := "set snmp v3 usm local-engine user \"" + d.Get("name").(string) + "\" "
	if d.Get("engine_type").(string) != "local" {
		engineID := d.Get("engine_id").(string)
		if engineID == "" {
			return fmt.Errorf("engine_id need to set when engine_type != local")
		}
		setPrefix = "set snmp v3 usm remote-engine \"" + engineID + "\" user \"" + d.Get("name").(string) + "\" "
	} else if d.Get("engine_id").(string) != "" {
		return fmt.Errorf("engine_id not compatible if engine_type = local")
	}
	configSet := make([]string, 0)

	if authType := d.Get("authentication_type").(string); authType != "authentication-none" {
		if d.Get("authentication_key").(string) == "" && d.Get("authentication_password").(string) == "" {
			return fmt.Errorf("authentication_key or authentication_password need to set " +
				"when authentication_type != authentication-none")
		}
		setPrefixAuth := setPrefix + authType + " "
		if authKey := d.Get("authentication_key").(string); authKey != "" {
			configSet = append(configSet, setPrefixAuth+"authentication-key \""+authKey+"\"")
		}
		if authPass := d.Get("authentication_password").(string); authPass != "" {
			configSet = append(configSet, setPrefixAuth+"authentication-password \""+authPass+"\"")
		}
	} else {
		if d.Get("privacy_type").(string) != "privacy-none" {
			return fmt.Errorf("authentication should be configured before configuring the privacy")
		}
		if d.Get("authentication_key").(string) != "" {
			return fmt.Errorf("authentication_key not compatible when authentication_type = authentication-none")
		}
		if d.Get("authentication_password").(string) != "" {
			return fmt.Errorf("authentication_password not compatible when authentication_type = authentication-none")
		}
		configSet = append(configSet, setPrefix+"authentication-none")
	}
	if privType := d.Get("privacy_type").(string); privType != "privacy-none" {
		if d.Get("privacy_key").(string) == "" && d.Get("privacy_password").(string) == "" {
			return fmt.Errorf("privacy_key or privacy_password need to set when privacy_type != privacy-none")
		}
		setPrefixPriv := setPrefix + privType + " "
		if privKey := d.Get("privacy_key").(string); privKey != "" {
			configSet = append(configSet, setPrefixPriv+"privacy-key \""+privKey+"\"")
		}
		if privPass := d.Get("privacy_password").(string); privPass != "" {
			configSet = append(configSet, setPrefixPriv+"privacy-password \""+privPass+"\"")
		}
	} else {
		if d.Get("privacy_key").(string) != "" {
			return fmt.Errorf("privacy_key not compatible when privacy_type = privacy-none")
		}
		if d.Get("privacy_password").(string) != "" {
			return fmt.Errorf("privacy_password not compatible when privacy_type = privacy-none")
		}
		configSet = append(configSet, setPrefix+"privacy-none")
	}

	return junSess.ConfigSet(configSet)
}

func readSnmpV3UsmUser(confSrc snmpV3UsmUserOptions, junSess *junos.Session,
) (confRead snmpV3UsmUserOptions, err error) {
	showCommand := junos.CmdShowConfig +
		"snmp v3 usm local-engine user \"" + confSrc.name + "\"" + junos.PipeDisplaySetRelative
	if confSrc.engineType != "local" {
		showCommand = junos.CmdShowConfig + "snmp v3 usm remote-engine \"" + confSrc.engineID +
			"\" user \"" + confSrc.name + "\"" + junos.PipeDisplaySetRelative
	}
	showConfig, err := junSess.Command(showCommand)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = confSrc.name
		confRead.engineType = confSrc.engineType
		confRead.engineID = confSrc.engineID
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "authentication-md5 authentication-key "):
				confRead.authenticationType = "authentication-md5"
				if confSrc.authenticationPassword != "" && confSrc.authenticationType == confRead.authenticationType {
					confRead.authenticationPassword = confSrc.authenticationPassword
				} else {
					confRead.authenticationKey, err = jdecode.Decode(strings.Trim(itemTrim, "\""))
					if err != nil {
						return confRead, fmt.Errorf("decoding authentication-key: %w", err)
					}
				}
			case itemTrim == "authentication-none":
				confRead.authenticationType = itemTrim
			case balt.CutPrefixInString(&itemTrim, "authentication-sha authentication-key "):
				confRead.authenticationType = "authentication-sha"
				if confSrc.authenticationPassword != "" && confSrc.authenticationType == confRead.authenticationType {
					confRead.authenticationPassword = confSrc.authenticationPassword
				} else {
					confRead.authenticationKey, err = jdecode.Decode(strings.Trim(itemTrim, "\""))
					if err != nil {
						return confRead, fmt.Errorf("decoding authentication-key: %w", err)
					}
				}
			case balt.CutPrefixInString(&itemTrim, "privacy-3des privacy-key "):
				confRead.privacyType = "privacy-3des"
				if confSrc.privacyPassword != "" && confSrc.privacyType == confRead.privacyType {
					confRead.privacyPassword = confSrc.privacyPassword
				} else {
					confRead.privacyKey, err = jdecode.Decode(strings.Trim(itemTrim, "\""))
					if err != nil {
						return confRead, fmt.Errorf("decoding privacy-key: %w", err)
					}
				}
			case balt.CutPrefixInString(&itemTrim, "privacy-aes128 privacy-key "):
				confRead.privacyType = "privacy-aes128"
				if confSrc.privacyPassword != "" && confSrc.privacyType == confRead.privacyType {
					confRead.privacyPassword = confSrc.privacyPassword
				} else {
					confRead.privacyKey, err = jdecode.Decode(strings.Trim(itemTrim, "\""))
					if err != nil {
						return confRead, fmt.Errorf("decoding privacy-key: %w", err)
					}
				}
			case balt.CutPrefixInString(&itemTrim, "privacy-des privacy-key "):
				confRead.privacyType = "privacy-des"
				if confSrc.privacyPassword != "" && confSrc.privacyType == confRead.privacyType {
					confRead.privacyPassword = confSrc.privacyPassword
				} else {
					confRead.privacyKey, err = jdecode.Decode(strings.Trim(itemTrim, "\""))
					if err != nil {
						return confRead, fmt.Errorf("decoding privacy-key: %w", err)
					}
				}
			case itemTrim == "privacy-none":
				confRead.privacyType = itemTrim
			}
		}
	}

	return confRead, nil
}

func delSnmpV3UsmUser(name, engineType, engineID string, junSess *junos.Session) error {
	var configSet []string
	if engineType == "local" {
		configSet = append(configSet, "delete snmp v3 usm local-engine user \""+name+"\"")
	} else {
		configSet = append(configSet, "delete snmp v3 usm remote-engine \""+engineID+"\" user \""+name+"\"")
	}

	return junSess.ConfigSet(configSet)
}

func fillSnmpV3UsmUserData(d *schema.ResourceData, snmpV3UsmUserOptions snmpV3UsmUserOptions) {
	if tfErr := d.Set("name", snmpV3UsmUserOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("engine_type", snmpV3UsmUserOptions.engineType); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("engine_id", snmpV3UsmUserOptions.engineID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_key", snmpV3UsmUserOptions.authenticationKey); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_password", snmpV3UsmUserOptions.authenticationPassword); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_type", snmpV3UsmUserOptions.authenticationType); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("privacy_key", snmpV3UsmUserOptions.privacyKey); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("privacy_password", snmpV3UsmUserOptions.privacyPassword); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("privacy_type", snmpV3UsmUserOptions.privacyType); tfErr != nil {
		panic(tfErr)
	}
}
