package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"golang.org/x/crypto/ssh"
)

type systemRootAuthOptions struct {
	noPublicKeys      bool
	encryptedPassword string
	sshPublicKeys     []string
}

func resourceSystemRootAuthentication() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSystemRootAuthenticationCreate,
		ReadWithoutTimeout:   resourceSystemRootAuthenticationRead,
		UpdateWithoutTimeout: resourceSystemRootAuthenticationUpdate,
		DeleteWithoutTimeout: resourceSystemRootAuthenticationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSystemRootAuthenticationImport,
		},
		Schema: map[string]*schema.Schema{
			"encrypted_password": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"encrypted_password", "plain_text_password"},
			},
			"plain_text_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"encrypted_password", "plain_text_password"},
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

func resourceSystemRootAuthenticationCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		// To be able detect a plain text password not accepted by system
		if d.Get("plain_text_password").(string) != "" {
			if err := delSystemRootAuthenticationPassword(clt, nil); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := setSystemRootAuthentication(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("system_root_authentication")

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
	// To be able detect a plain text password not accepted by system
	if d.Get("plain_text_password").(string) != "" {
		if err := delSystemRootAuthenticationPassword(clt, junSess); err != nil {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	if err := setSystemRootAuthentication(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_system_root_authentication", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("system_root_authentication")

	return append(diagWarns, resourceSystemRootAuthenticationReadWJunSess(d, clt, junSess)...)
}

func resourceSystemRootAuthenticationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSystemRootAuthenticationReadWJunSess(d, clt, junSess)
}

func resourceSystemRootAuthenticationReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	systemRootAuthOptions, err := readSystemRootAuthentication(clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillSystemRootAuthenticationData(d, systemRootAuthOptions)

	return nil
}

func resourceSystemRootAuthenticationUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSystemRootAuthentication(clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystemRootAuthentication(d, clt, nil); err != nil {
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
	if err := delSystemRootAuthentication(clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemRootAuthentication(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_system_root_authentication", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemRootAuthenticationReadWJunSess(d, clt, junSess)...)
}

func resourceSystemRootAuthenticationDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	return nil
}

func resourceSystemRootAuthenticationImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)

	systemRootAuthOptions, err := readSystemRootAuthentication(clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSystemRootAuthenticationData(d, systemRootAuthOptions)
	d.SetId("system_root_authentication")
	result[0] = d

	return result, nil
}

func setSystemRootAuthentication(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)
	setPrefix := "set system root-authentication "

	if v := d.Get("plain_text_password").(string); v != "" {
		configSet = append(configSet, setPrefix+"plain-text-password-value \""+v+"\"")
	} else {
		configSet = append(configSet, setPrefix+"encrypted-password \""+d.Get("encrypted_password").(string)+"\"")
	}
	if d.Get("no_public_keys").(bool) {
		configSet = append(configSet, setPrefix+"no-public-keys")
	}
	for _, v := range sortSetOfString(d.Get("ssh_public_keys").(*schema.Set).List()) {
		switch {
		case strings.HasPrefix(v, ssh.KeyAlgoDSA):
			configSet = append(configSet, setPrefix+"ssh-dsa \""+v+"\"")
		case strings.HasPrefix(v, ssh.KeyAlgoRSA):
			configSet = append(configSet, setPrefix+"ssh-rsa \""+v+"\"")
		case strings.HasPrefix(v, ssh.KeyAlgoECDSA256),
			strings.HasPrefix(v, ssh.KeyAlgoECDSA384),
			strings.HasPrefix(v, ssh.KeyAlgoECDSA521):
			configSet = append(configSet, setPrefix+"ssh-ecdsa \""+v+"\"")
		case strings.HasPrefix(v, ssh.KeyAlgoED25519):
			configSet = append(configSet, setPrefix+"ssh-ed25519 \""+v+"\"")
		default:
			return fmt.Errorf("format in public key '%v' not supported", v)
		}
	}

	return clt.configSet(configSet, junSess)
}

func readSystemRootAuthentication(clt *Client, junSess *junosSession) (systemRootAuthOptions, error) {
	var confRead systemRootAuthOptions

	showConfig, err := clt.command(cmdShowConfig+"system root-authentication"+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
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

func delSystemRootAuthentication(clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system root-authentication")

	return clt.configSet(configSet, junSess)
}

func delSystemRootAuthenticationPassword(clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system root-authentication encrypted-password")

	return clt.configSet(configSet, junSess)
}

func fillSystemRootAuthenticationData(d *schema.ResourceData, systemRootAuthOptions systemRootAuthOptions) {
	if d.Get("plain_text_password").(string) == "" {
		if tfErr := d.Set("encrypted_password", systemRootAuthOptions.encryptedPassword); tfErr != nil {
			panic(tfErr)
		}
	}
	if tfErr := d.Set("no_public_keys", systemRootAuthOptions.noPublicKeys); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ssh_public_keys", systemRootAuthOptions.sshPublicKeys); tfErr != nil {
		panic(tfErr)
	}
}
