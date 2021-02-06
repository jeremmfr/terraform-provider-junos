package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/ssh"
)

type systemRootAuthOptions struct {
	noPublicKeys      bool
	encryptedPassword string
	sshPublicKeys     []string
}

func resourceSystemRootAuthentication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemRootAuthenticationCreate,
		ReadContext:   resourceSystemRootAuthenticationRead,
		UpdateContext: resourceSystemRootAuthenticationUpdate,
		DeleteContext: resourceSystemRootAuthenticationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSystemRootAuthenticationImport,
		},
		Schema: map[string]*schema.Schema{
			"encrypted_password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"no_public_keys": {
				Type:     schema.TypeBool,
				Optional: true,
				ConflictsWith: []string{
					"ssh_public_keys",
				},
			},
			"ssh_public_keys": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"no_public_keys"},
			},
		},
	}
}

func resourceSystemRootAuthenticationCreate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)

	if err := setSystemRootAuthentication(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_system_root_authentication", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("system_root_authentication")

	return append(diagWarns, resourceSystemRootAuthenticationReadWJnprSess(d, m, jnprSess)...)
}
func resourceSystemRootAuthenticationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSystemRootAuthenticationReadWJnprSess(d, m, jnprSess)
}
func resourceSystemRootAuthenticationReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	systemRootAuthOptions, err := readSystemRootAuthentication(m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillSystemRootAuthenticationData(d, systemRootAuthOptions)

	return nil
}
func resourceSystemRootAuthenticationUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSystemRootAuthentication(m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setSystemRootAuthentication(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_system_root_authentication", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemRootAuthenticationReadWJnprSess(d, m, jnprSess)...)
}
func resourceSystemRootAuthenticationDelete(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
func resourceSystemRootAuthenticationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	systemRootAuthOptions, err := readSystemRootAuthentication(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSystemRootAuthenticationData(d, systemRootAuthOptions)
	d.SetId("system_root_authentication")
	result[0] = d

	return result, nil
}

func setSystemRootAuthentication(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := "set system root-authentication "

	configSet = append(configSet, setPrefix+"encrypted-password \""+d.Get("encrypted_password").(string)+"\"")
	if d.Get("no_public_keys").(bool) {
		configSet = append(configSet, setPrefix+"no-public-keys")
	}
	for _, v := range d.Get("ssh_public_keys").(*schema.Set).List() {
		switch {
		case strings.HasPrefix(v.(string), ssh.KeyAlgoDSA):
			configSet = append(configSet, setPrefix+"ssh-dsa \""+v.(string)+"\"")
		case strings.HasPrefix(v.(string), ssh.KeyAlgoRSA):
			configSet = append(configSet, setPrefix+"ssh-rsa \""+v.(string)+"\"")
		case strings.HasPrefix(v.(string), ssh.KeyAlgoECDSA256),
			strings.HasPrefix(v.(string), ssh.KeyAlgoECDSA384),
			strings.HasPrefix(v.(string), ssh.KeyAlgoECDSA521):
			configSet = append(configSet, setPrefix+"ssh-ecdsa \""+v.(string)+"\"")
		case strings.HasPrefix(v.(string), ssh.KeyAlgoED25519):
			configSet = append(configSet, setPrefix+"ssh-ed25519 \""+v.(string)+"\"")
		default:
			return fmt.Errorf("format in public key '%v' not supported", v.(string))
		}
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readSystemRootAuthentication(m interface{}, jnprSess *NetconfObject) (systemRootAuthOptions, error) {
	sess := m.(*Session)
	var confRead systemRootAuthOptions

	systemRootAuthConfig, err := sess.command("show configuration system root-authentication"+
		" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if systemRootAuthConfig != emptyWord {
		for _, item := range strings.Split(systemRootAuthConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "encrypted-password "):
				confRead.encryptedPassword = strings.Trim(strings.TrimPrefix(itemTrim, "encrypted-password "), "\"")
			case itemTrim == "no-public-keys":
				confRead.noPublicKeys = true
			case strings.HasPrefix(itemTrim, "ssh-dsa "):
				confRead.sshPublicKeys = append(confRead.sshPublicKeys,
					strings.Trim(strings.TrimPrefix(itemTrim, "ssh-dsa "), "\""))
			case strings.HasPrefix(itemTrim, "ssh-ecdsa "):
				confRead.sshPublicKeys = append(confRead.sshPublicKeys,
					strings.Trim(strings.TrimPrefix(itemTrim, "ssh-ecdsa "), "\""))
			case strings.HasPrefix(itemTrim, "ssh-ed25519 "):
				confRead.sshPublicKeys = append(confRead.sshPublicKeys,
					strings.Trim(strings.TrimPrefix(itemTrim, "ssh-ed25519 "), "\""))
			case strings.HasPrefix(itemTrim, "ssh-rsa "):
				confRead.sshPublicKeys = append(confRead.sshPublicKeys,
					strings.Trim(strings.TrimPrefix(itemTrim, "ssh-rsa "), "\""))
			}
		}
	}

	return confRead, nil
}

func delSystemRootAuthentication(m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system root-authentication")
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func fillSystemRootAuthenticationData(d *schema.ResourceData, systemRootAuthOptions systemRootAuthOptions) {
	if tfErr := d.Set("encrypted_password", systemRootAuthOptions.encryptedPassword); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_public_keys", systemRootAuthOptions.noPublicKeys); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ssh_public_keys", systemRootAuthOptions.sshPublicKeys); tfErr != nil {
		panic(tfErr)
	}
}
