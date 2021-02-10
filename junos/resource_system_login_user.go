package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"golang.org/x/crypto/ssh"
)

type systemLoginUserOptions struct {
	uid            int
	name           string
	class          string
	cliPrompt      string
	fullName       string
	authentication []map[string]interface{}
}

func resourceSystemLoginUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemLoginUserCreate,
		ReadContext:   resourceSystemLoginUserRead,
		UpdateContext: resourceSystemLoginUserUpdate,
		DeleteContext: resourceSystemLoginUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSystemLoginUserImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uid": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(100, 64000),
			},
			"authentication": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encrypted_password": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"no_public_keys": {
							Type:     schema.TypeBool,
							Optional: true,
							ConflictsWith: []string{
								"authentication.0.ssh_public_keys",
							},
						},
						"ssh_public_keys": {
							Type:          schema.TypeSet,
							Optional:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							ConflictsWith: []string{"authentication.0.no_public_keys"},
						},
					},
				},
			},
			"cli_prompt": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSystemLoginUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	systemLoginUserExists, err := checkSystemLoginUserExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if systemLoginUserExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("system login user %v already exists", d.Get("name").(string)))
	}

	if err := setSystemLoginUser(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_system_login_user", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	systemLoginUserExists, err = checkSystemLoginUserExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if systemLoginUserExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system login user %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSystemLoginUserReadWJnprSess(d, m, jnprSess)...)
}
func resourceSystemLoginUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSystemLoginUserReadWJnprSess(d, m, jnprSess)
}
func resourceSystemLoginUserReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	systemLoginUserOptions, err := readSystemLoginUser(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if systemLoginUserOptions.name == "" {
		d.SetId("")
	} else {
		fillSystemLoginUserData(d, systemLoginUserOptions)
	}

	return nil
}
func resourceSystemLoginUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSystemLoginUser(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setSystemLoginUser(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_system_login_user", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	d.Partial(false)

	return append(diagWarns, resourceSystemLoginUserReadWJnprSess(d, m, jnprSess)...)
}
func resourceSystemLoginUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSystemLoginUser(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_system_login_user", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceSystemLoginUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	systemLoginUserExists, err := checkSystemLoginUserExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !systemLoginUserExists {
		return nil, fmt.Errorf("don't find system login user with id '%v' (id must be <name>)", d.Id())
	}
	systemLoginUserOptions, err := readSystemLoginUser(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSystemLoginUserData(d, systemLoginUserOptions)

	result[0] = d

	return result, nil
}

func checkSystemLoginUserExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	systemLoginUserConfig, err := sess.command("show configuration system login user "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if systemLoginUserConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setSystemLoginUser(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := "set system login user " + d.Get("name").(string) + " "

	configSet = append(configSet, setPrefix+"class "+d.Get("class").(string))

	if d.Get("uid").(int) != 0 {
		configSet = append(configSet, setPrefix+"uid "+strconv.Itoa(d.Get("uid").(int)))
	}
	for _, v := range d.Get("authentication").([]interface{}) {
		if v == nil {
			return fmt.Errorf("authentication block is empty")
		}
		authentication := v.(map[string]interface{})
		if authentication["encrypted_password"].(string) != "" {
			configSet = append(configSet, setPrefix+"authentication encrypted-password \""+
				authentication["encrypted_password"].(string)+"\"")
		}
		if authentication["no_public_keys"].(bool) {
			configSet = append(configSet, setPrefix+"authentication no-public-keys")
		}
		for _, v2 := range authentication["ssh_public_keys"].(*schema.Set).List() {
			switch {
			case strings.HasPrefix(v2.(string), ssh.KeyAlgoDSA):
				configSet = append(configSet, setPrefix+"authentication ssh-dsa \""+v2.(string)+"\"")
			case strings.HasPrefix(v2.(string), ssh.KeyAlgoRSA):
				configSet = append(configSet, setPrefix+"authentication ssh-rsa \""+v2.(string)+"\"")
			case strings.HasPrefix(v2.(string), ssh.KeyAlgoECDSA256),
				strings.HasPrefix(v2.(string), ssh.KeyAlgoECDSA384),
				strings.HasPrefix(v2.(string), ssh.KeyAlgoECDSA521):
				configSet = append(configSet, setPrefix+"authentication ssh-ecdsa \""+v2.(string)+"\"")
			case strings.HasPrefix(v2.(string), ssh.KeyAlgoED25519):
				configSet = append(configSet, setPrefix+"authentication ssh-ed25519 \""+v2.(string)+"\"")
			default:
				return fmt.Errorf("format in public key '%v' not supported", v2.(string))
			}
		}
	}
	if d.Get("cli_prompt").(string) != "" {
		configSet = append(configSet, setPrefix+"cli prompt \""+d.Get("cli_prompt").(string)+"\"")
	}
	if d.Get("full_name").(string) != "" {
		configSet = append(configSet, setPrefix+"full-name \""+d.Get("full_name").(string)+"\"")
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readSystemLoginUser(user string, m interface{}, jnprSess *NetconfObject) (systemLoginUserOptions, error) {
	sess := m.(*Session)
	var confRead systemLoginUserOptions

	systemLoginUserConfig, err := sess.command("show configuration system login user "+
		user+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if systemLoginUserConfig != emptyWord {
		confRead.name = user
		for _, item := range strings.Split(systemLoginUserConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "class "):
				confRead.class = strings.TrimPrefix(itemTrim, "class ")
			case strings.HasPrefix(itemTrim, "uid "):
				var err error
				confRead.uid, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "uid "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "authentication "):
				if len(confRead.authentication) == 0 {
					confRead.authentication = append(confRead.authentication, map[string]interface{}{
						"encrypted_password": "",
						"no_public_keys":     false,
						"ssh_public_keys":    make([]string, 0),
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "authentication encrypted-password "):
					confRead.authentication[0]["encrypted_password"] = strings.Trim(strings.TrimPrefix(
						itemTrim, "authentication encrypted-password "), "\"")
				case itemTrim == "authentication no-public-keys":
					confRead.authentication[0]["no_public_keys"] = true
				case strings.HasPrefix(itemTrim, "authentication ssh-dsa "):
					confRead.authentication[0]["ssh_public_keys"] = append(confRead.authentication[0]["ssh_public_keys"].([]string),
						strings.Trim(strings.TrimPrefix(itemTrim, "authentication ssh-dsa "), "\""))
				case strings.HasPrefix(itemTrim, "authentication ssh-ecdsa "):
					confRead.authentication[0]["ssh_public_keys"] = append(confRead.authentication[0]["ssh_public_keys"].([]string),
						strings.Trim(strings.TrimPrefix(itemTrim, "authentication ssh-ecdsa "), "\""))
				case strings.HasPrefix(itemTrim, "authentication ssh-ed25519 "):
					confRead.authentication[0]["ssh_public_keys"] = append(confRead.authentication[0]["ssh_public_keys"].([]string),
						strings.Trim(strings.TrimPrefix(itemTrim, "authentication ssh-ed25519 "), "\""))
				case strings.HasPrefix(itemTrim, "authentication ssh-rsa "):
					confRead.authentication[0]["ssh_public_keys"] = append(confRead.authentication[0]["ssh_public_keys"].([]string),
						strings.Trim(strings.TrimPrefix(itemTrim, "authentication ssh-rsa "), "\""))
				}
			case strings.HasPrefix(itemTrim, "cli prompt "):
				confRead.cliPrompt = strings.Trim(strings.TrimPrefix(itemTrim, "cli prompt "), "\"")
			case strings.HasPrefix(itemTrim, "full-name "):
				confRead.fullName = strings.Trim(strings.TrimPrefix(itemTrim, "full-name "), "\"")
			}
		}
	}

	return confRead, nil
}

func delSystemLoginUser(systemLoginUser string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system login user "+systemLoginUser)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func fillSystemLoginUserData(d *schema.ResourceData, systemLoginUserOptions systemLoginUserOptions) {
	if tfErr := d.Set("name", systemLoginUserOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("class", systemLoginUserOptions.class); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("uid", systemLoginUserOptions.uid); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication", systemLoginUserOptions.authentication); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("cli_prompt", systemLoginUserOptions.cliPrompt); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("full_name", systemLoginUserOptions.fullName); tfErr != nil {
		panic(tfErr)
	}
}
