package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
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
		CreateWithoutTimeout: resourceSystemLoginUserCreate,
		ReadWithoutTimeout:   resourceSystemLoginUserRead,
		UpdateWithoutTimeout: resourceSystemLoginUserUpdate,
		DeleteWithoutTimeout: resourceSystemLoginUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSystemLoginUserImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefAndDots),
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
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"authentication.0.plain_text_password"},
						},
						"no_public_keys": {
							Type:     schema.TypeBool,
							Optional: true,
							ConflictsWith: []string{
								"authentication.0.ssh_public_keys",
							},
						},
						"plain_text_password": {
							Type:          schema.TypeString,
							Optional:      true,
							Sensitive:     true,
							ConflictsWith: []string{"authentication.0.encrypted_password"},
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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSystemLoginUser(d, clt, nil); err != nil {
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
	systemLoginUserExists, err := checkSystemLoginUserExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if systemLoginUserExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("system login user %v already exists", d.Get("name").(string)))...)
	}

	if err := setSystemLoginUser(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_system_login_user", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	systemLoginUserExists, err = checkSystemLoginUserExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if systemLoginUserExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system login user %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSystemLoginUserReadWJunSess(d, clt, junSess)...)
}

func resourceSystemLoginUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSystemLoginUserReadWJunSess(d, clt, junSess)
}

func resourceSystemLoginUserReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	plainTextPassword := readSystemLoginUserReadDataPlainTextPassword(d)
	mutex.Lock()
	systemLoginUserOptions, err := readSystemLoginUser(d.Get("name").(string), plainTextPassword, clt, junSess)
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
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSystemLoginUser(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystemLoginUser(d, clt, nil); err != nil {
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
	if err := delSystemLoginUser(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemLoginUser(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_system_login_user", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	d.Partial(false)

	return append(diagWarns, resourceSystemLoginUserReadWJunSess(d, clt, junSess)...)
}

func resourceSystemLoginUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delSystemLoginUser(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delSystemLoginUser(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_system_login_user", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSystemLoginUserImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)

	systemLoginUserExists, err := checkSystemLoginUserExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !systemLoginUserExists {
		return nil, fmt.Errorf("don't find system login user with id '%v' (id must be <name>)", d.Id())
	}
	systemLoginUserOptions, err := readSystemLoginUser(d.Id(), "", clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSystemLoginUserData(d, systemLoginUserOptions)

	result[0] = d

	return result, nil
}

func checkSystemLoginUserExists(name string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+"system login user "+name+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSystemLoginUser(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)
	setPrefix := "set system login user " + d.Get("name").(string) + " "

	configSet = append(configSet, setPrefix+"class "+d.Get("class").(string))

	if v := d.Get("uid").(int); v != 0 {
		configSet = append(configSet, setPrefix+"uid "+strconv.Itoa(v))
	}
	for _, block := range d.Get("authentication").([]interface{}) {
		if block == nil {
			return fmt.Errorf("authentication block is empty")
		}
		authentication := block.(map[string]interface{})
		if pass := authentication["plain_text_password"].(string); pass != "" {
			configSet = append(configSet, setPrefix+"authentication plain-text-password-value \""+pass+"\"")
		}
		if pass := authentication["encrypted_password"].(string); pass != "" {
			configSet = append(configSet, setPrefix+"authentication encrypted-password \""+pass+"\"")
		}
		if authentication["no_public_keys"].(bool) {
			configSet = append(configSet, setPrefix+"authentication no-public-keys")
		}
		for _, key := range sortSetOfString(authentication["ssh_public_keys"].(*schema.Set).List()) {
			switch {
			case strings.HasPrefix(key, ssh.KeyAlgoDSA):
				configSet = append(configSet, setPrefix+"authentication ssh-dsa \""+key+"\"")
			case strings.HasPrefix(key, ssh.KeyAlgoRSA):
				configSet = append(configSet, setPrefix+"authentication ssh-rsa \""+key+"\"")
			case strings.HasPrefix(key, ssh.KeyAlgoECDSA256),
				strings.HasPrefix(key, ssh.KeyAlgoECDSA384),
				strings.HasPrefix(key, ssh.KeyAlgoECDSA521):
				configSet = append(configSet, setPrefix+"authentication ssh-ecdsa \""+key+"\"")
			case strings.HasPrefix(key, ssh.KeyAlgoED25519):
				configSet = append(configSet, setPrefix+"authentication ssh-ed25519 \""+key+"\"")
			default:
				return fmt.Errorf("format in public key '%v' not supported", key)
			}
		}
	}
	if v := d.Get("cli_prompt").(string); v != "" {
		configSet = append(configSet, setPrefix+"cli prompt \""+v+"\"")
	}
	if v := d.Get("full_name").(string); v != "" {
		configSet = append(configSet, setPrefix+"full-name \""+v+"\"")
	}

	return clt.configSet(configSet, junSess)
}

func readSystemLoginUserReadDataPlainTextPassword(d *schema.ResourceData) string {
	if blocks := d.Get("authentication").([]interface{}); len(blocks) > 0 {
		if blocks[0] == nil {
			return ""
		}
		auth := blocks[0].(map[string]interface{})

		return auth["plain_text_password"].(string)
	}

	return ""
}

func readSystemLoginUser(name, plainTextPassword string, clt *Client, junSess *junosSession,
) (confRead systemLoginUserOptions, err error) {
	showConfig, err := clt.command(cmdShowConfig+"system login user "+name+pipeDisplaySetRelative, junSess)
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
			case balt.CutPrefixInString(&itemTrim, "class "):
				confRead.class = itemTrim
			case balt.CutPrefixInString(&itemTrim, "uid "):
				confRead.uid, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "authentication "):
				if len(confRead.authentication) == 0 {
					confRead.authentication = append(confRead.authentication, map[string]interface{}{
						"encrypted_password":  "",
						"no_public_keys":      false,
						"plain_text_password": "",
						"ssh_public_keys":     make([]string, 0),
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "encrypted-password "):
					if plainTextPassword != "" {
						confRead.authentication[0]["plain_text_password"] = plainTextPassword
					} else {
						confRead.authentication[0]["encrypted_password"] = strings.Trim(itemTrim, "\"")
					}
				case itemTrim == "no-public-keys":
					confRead.authentication[0]["no_public_keys"] = true
				case balt.CutPrefixInString(&itemTrim, "ssh-dsa "),
					balt.CutPrefixInString(&itemTrim, "ssh-ecdsa "),
					balt.CutPrefixInString(&itemTrim, "ssh-ed25519 "),
					balt.CutPrefixInString(&itemTrim, "ssh-rsa "):
					confRead.authentication[0]["ssh_public_keys"] = append(
						confRead.authentication[0]["ssh_public_keys"].([]string),
						strings.Trim(itemTrim, "\""),
					)
				}
			case balt.CutPrefixInString(&itemTrim, "cli prompt "):
				confRead.cliPrompt = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "full-name "):
				confRead.fullName = strings.Trim(itemTrim, "\"")
			}
		}
	}

	return confRead, nil
}

func delSystemLoginUser(systemLoginUser string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system login user "+systemLoginUser)

	return clt.configSet(configSet, junSess)
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
