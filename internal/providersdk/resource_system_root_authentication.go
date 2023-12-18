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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		// To be able detect a plain text password not accepted by system
		if d.Get("plain_text_password").(string) != "" {
			if err := delSystemRootAuthenticationPassword(junSess); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := setSystemRootAuthentication(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("system_root_authentication")

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
	// To be able detect a plain text password not accepted by system
	if d.Get("plain_text_password").(string) != "" {
		if err := delSystemRootAuthenticationPassword(junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	if err := setSystemRootAuthentication(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_system_root_authentication")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("system_root_authentication")

	return append(diagWarns, resourceSystemRootAuthenticationReadWJunSess(d, junSess)...)
}

func resourceSystemRootAuthenticationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSystemRootAuthenticationReadWJunSess(d, junSess)
}

func resourceSystemRootAuthenticationReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	systemRootAuthOptions, err := readSystemRootAuthentication(junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillSystemRootAuthenticationData(d, systemRootAuthOptions)

	return nil
}

func resourceSystemRootAuthenticationUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSystemRootAuthentication(junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystemRootAuthentication(d, junSess); err != nil {
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
	if err := delSystemRootAuthentication(junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemRootAuthentication(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_system_root_authentication")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemRootAuthenticationReadWJunSess(d, junSess)...)
}

func resourceSystemRootAuthenticationDelete(_ context.Context, _ *schema.ResourceData, _ interface{},
) diag.Diagnostics {
	return nil
}

func resourceSystemRootAuthenticationImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	systemRootAuthOptions, err := readSystemRootAuthentication(junSess)
	if err != nil {
		return nil, err
	}
	fillSystemRootAuthenticationData(d, systemRootAuthOptions)
	d.SetId("system_root_authentication")
	result[0] = d

	return result, nil
}

func setSystemRootAuthentication(d *schema.ResourceData, junSess *junos.Session) error {
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

	return junSess.ConfigSet(configSet)
}

func readSystemRootAuthentication(junSess *junos.Session,
) (confRead systemRootAuthOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "system root-authentication" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "encrypted-password "):
				confRead.encryptedPassword = strings.Trim(itemTrim, "\"")
			case itemTrim == "no-public-keys":
				confRead.noPublicKeys = true
			case balt.CutPrefixInString(&itemTrim, "ssh-dsa "),
				balt.CutPrefixInString(&itemTrim, "ssh-ecdsa "),
				balt.CutPrefixInString(&itemTrim, "ssh-ed25519 "),
				balt.CutPrefixInString(&itemTrim, "ssh-rsa "):
				confRead.sshPublicKeys = append(confRead.sshPublicKeys, strings.Trim(itemTrim, "\""))
			}
		}
	}

	return confRead, nil
}

func delSystemRootAuthentication(junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system root-authentication")

	return junSess.ConfigSet(configSet)
}

func delSystemRootAuthenticationPassword(junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system root-authentication encrypted-password")

	return junSess.ConfigSet(configSet)
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
