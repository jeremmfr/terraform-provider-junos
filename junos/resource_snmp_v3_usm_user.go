package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSnmpV3UsmUser(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		if d.Get("engine_type").(string) == "local" {
			d.SetId("local" + idSeparator + d.Get("name").(string))
		} else {
			d.SetId("remote" + idSeparator + d.Get("engine_id").(string) + idSeparator + d.Get("name").(string))
		}

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	snmpV3UsmUserExists, err := checkSnmpV3UsmUserExists(
		d.Get("name").(string), d.Get("engine_type").(string), d.Get("engine_id").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3UsmUserExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		if d.Get("engine_type").(string) != "local" {
			return append(diagWarns, diag.FromErr(fmt.Errorf("snmp v3 usm user %v in remote-engine %s already exists",
				d.Get("name").(string), d.Get("engine_id").(string)))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp v3 usm user %v in local-engine already exists",
			d.Get("name").(string)))...)
	}

	if err := setSnmpV3UsmUser(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_snmp_v3_usm_user", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	snmpV3UsmUserExists, err = checkSnmpV3UsmUserExists(
		d.Get("name").(string), d.Get("engine_type").(string), d.Get("engine_id").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3UsmUserExists {
		if d.Get("engine_type").(string) == "local" {
			d.SetId("local" + idSeparator + d.Get("name").(string))
		} else {
			d.SetId("remote" + idSeparator + d.Get("engine_id").(string) + idSeparator + d.Get("name").(string))
		}
	} else {
		if d.Get("engine_type").(string) != "local" {
			return append(diagWarns, diag.FromErr(fmt.Errorf("snmp v3 usm user %v in remote-engine %s not exists after commit "+
				"=> check your config", d.Get("name").(string), d.Get("engine_id").(string)))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp v3 usm user %v in local-engine not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSnmpV3UsmUserReadWJnprSess(d, m, jnprSess)...)
}

func resourceSnmpV3UsmUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSnmpV3UsmUserReadWJnprSess(d, m, jnprSess)
}

func resourceSnmpV3UsmUserReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
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
	mutex.Lock()
	snmpV3UsmUserOptions, err := readSnmpV3UsmUser(configSrc, m, jnprSess)
	mutex.Unlock()
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
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delSnmpV3UsmUser(
			d.Get("name").(string), d.Get("engine_type").(string), d.Get("engine_id").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSnmpV3UsmUser(d, m, nil); err != nil {
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
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSnmpV3UsmUser(
		d.Get("name").(string), d.Get("engine_type").(string), d.Get("engine_id").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmpV3UsmUser(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_snmp_v3_usm_user", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpV3UsmUserReadWJnprSess(d, m, jnprSess)...)
}

func resourceSnmpV3UsmUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delSnmpV3UsmUser(
			d.Get("name").(string), d.Get("engine_type").(string), d.Get("engine_id").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSnmpV3UsmUser(
		d.Get("name").(string), d.Get("engine_type").(string), d.Get("engine_id").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_snmp_v3_usm_user", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSnmpV3UsmUserImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	var configImport snmpV3UsmUserOptions
	switch {
	case len(idSplit) == 2 && idSplit[0] == "local":
		snmpV3UsmUserExists, err := checkSnmpV3UsmUserExists(idSplit[1], idSplit[0], "", m, jnprSess)
		if err != nil {
			return nil, err
		}
		if !snmpV3UsmUserExists {
			return nil, fmt.Errorf("don't find snmp v3 usm user with id '%v' (id must be local%s<name>)", d.Id(), idSeparator)
		}
		configImport.name = idSplit[1]
		configImport.engineType = idSplit[0]
	case len(idSplit) == 3 && idSplit[0] == "remote":
		snmpV3UsmUserExists, err := checkSnmpV3UsmUserExists(idSplit[2], idSplit[0], idSplit[1], m, jnprSess)
		if err != nil {
			return nil, err
		}
		if !snmpV3UsmUserExists {
			return nil, fmt.Errorf("don't find snmp v3 usm user with id "+
				"'%v' (id must be remote%s<engine_id>%s<name>)", d.Id(), idSeparator, idSeparator)
		}
		configImport.name = idSplit[2]
		configImport.engineType = idSplit[0]
		configImport.engineID = idSplit[1]
	default:
		return nil, fmt.Errorf("can't find snmp v3 usm user with id "+
			"'%v' (id must be local%s<name> or remote%s<engine_id>%s<name>)", d.Id(), idSeparator, idSeparator, idSeparator)
	}
	snmpV3UsmUserOptions, err := readSnmpV3UsmUser(configImport, m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSnmpV3UsmUserData(d, snmpV3UsmUserOptions)

	result[0] = d

	return result, nil
}

func checkSnmpV3UsmUserExists(name, engineType, engineID string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	if engineType == "local" {
		showConfig, err := sess.command(cmdShowConfig+
			"snmp v3 usm local-engine user \""+name+"\""+pipeDisplaySet, jnprSess)
		if err != nil {
			return false, err
		}
		if showConfig == emptyW {
			return false, nil
		}
	} else {
		showConfig, err := sess.command(cmdShowConfig+
			"snmp v3 usm remote-engine \""+engineID+"\" user \""+name+"\""+pipeDisplaySet, jnprSess)
		if err != nil {
			return false, err
		}
		if showConfig == emptyW {
			return false, nil
		}
	}

	return true, nil
}

func setSnmpV3UsmUser(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

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

	return sess.configSet(configSet, jnprSess)
}

func readSnmpV3UsmUser(confSrc snmpV3UsmUserOptions, m interface{}, jnprSess *NetconfObject,
) (snmpV3UsmUserOptions, error) {
	sess := m.(*Session)
	var confRead snmpV3UsmUserOptions

	showCommand := cmdShowConfig + "snmp v3 usm local-engine user \"" + confSrc.name + "\"" + pipeDisplaySetRelative
	if confSrc.engineType != "local" {
		showCommand = cmdShowConfig + "snmp v3 usm remote-engine \"" + confSrc.engineID +
			"\" user \"" + confSrc.name + "\"" + pipeDisplaySetRelative
	}
	showConfig, err := sess.command(showCommand, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = confSrc.name
		confRead.engineType = confSrc.engineType
		confRead.engineID = confSrc.engineID
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "authentication-md5 authentication-key "):
				confRead.authenticationType = "authentication-md5"
				if confSrc.authenticationPassword != "" && confSrc.authenticationType == confRead.authenticationType {
					confRead.authenticationPassword = confSrc.authenticationPassword
				} else {
					var err error
					confRead.authenticationKey, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(
						itemTrim, "authentication-md5 authentication-key "), "\""))
					if err != nil {
						return confRead, fmt.Errorf("failed to decode authentication-key : %w", err)
					}
				}
			case itemTrim == "authentication-none":
				confRead.authenticationType = itemTrim
			case strings.HasPrefix(itemTrim, "authentication-sha authentication-key "):
				confRead.authenticationType = "authentication-sha"
				if confSrc.authenticationPassword != "" && confSrc.authenticationType == confRead.authenticationType {
					confRead.authenticationPassword = confSrc.authenticationPassword
				} else {
					var err error
					confRead.authenticationKey, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(
						itemTrim, "authentication-sha authentication-key "), "\""))
					if err != nil {
						return confRead, fmt.Errorf("failed to decode authentication-key : %w", err)
					}
				}
			case strings.HasPrefix(itemTrim, "privacy-3des privacy-key "):
				confRead.privacyType = "privacy-3des"
				if confSrc.privacyPassword != "" && confSrc.privacyType == confRead.privacyType {
					confRead.privacyPassword = confSrc.privacyPassword
				} else {
					var err error
					confRead.privacyKey, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(
						itemTrim, "privacy-3des privacy-key "), "\""))
					if err != nil {
						return confRead, fmt.Errorf("failed to decode privacy-key : %w", err)
					}
				}
			case strings.HasPrefix(itemTrim, "privacy-aes128 privacy-key "):
				confRead.privacyType = "privacy-aes128"
				if confSrc.privacyPassword != "" && confSrc.privacyType == confRead.privacyType {
					confRead.privacyPassword = confSrc.privacyPassword
				} else {
					var err error
					confRead.privacyKey, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(
						itemTrim, "privacy-aes128 privacy-key "), "\""))
					if err != nil {
						return confRead, fmt.Errorf("failed to decode privacy-key : %w", err)
					}
				}
			case strings.HasPrefix(itemTrim, "privacy-des privacy-key "):
				confRead.privacyType = "privacy-des"
				if confSrc.privacyPassword != "" && confSrc.privacyType == confRead.privacyType {
					confRead.privacyPassword = confSrc.privacyPassword
				} else {
					var err error
					confRead.privacyKey, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(
						itemTrim, "privacy-des privacy-key "), "\""))
					if err != nil {
						return confRead, fmt.Errorf("failed to decode privacy-key : %w", err)
					}
				}
			case itemTrim == "privacy-none":
				confRead.privacyType = itemTrim
			}
		}
	}

	return confRead, nil
}

func delSnmpV3UsmUser(name, engineType, engineID string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	var configSet []string
	if engineType == "local" {
		configSet = append(configSet, "delete snmp v3 usm local-engine user \""+name+"\"")
	} else {
		configSet = append(configSet, "delete snmp v3 usm remote-engine \""+engineID+"\" user \""+name+"\"")
	}

	return sess.configSet(configSet, jnprSess)
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
