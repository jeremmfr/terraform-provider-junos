package providersdk

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNullCommitFile() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceNullCommitFileCreate,
		ReadWithoutTimeout:   resourceNullCommitFileRead,
		DeleteWithoutTimeout: resourceNullCommitFileDelete,
		Schema: map[string]*schema.Schema{
			"filename": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"append_lines": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"clear_file_after_commit": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"triggers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     nil,
			},
		},
	}
}

func resourceNullCommitFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	fileName := d.Get("filename").(string)
	configSet, err := readNullCommitFile(fileName)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	for _, v := range d.Get("append_lines").([]interface{}) {
		configSet = append(configSet, v.(string))
	}
	if err := junSess.ConfigSet(configSet); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "commit a file with resource junos_null_commit_file")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId(fileName)
	if d.Get("clear_file_after_commit").(bool) {
		if err := cleanNullCommitFile(fileName, clt); err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
	}

	return diagWarns
}

func resourceNullCommitFileRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNullCommitFileDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}

func readNullCommitFile(filename string) ([]string, error) {
	if err := utils.ReplaceTildeToHomeDir(&filename); err != nil {
		return []string{}, err
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []string{}, fmt.Errorf("file '%s' doesn't exist", filename)
	}
	fileReadByte, err := os.ReadFile(filename)
	if err != nil {
		return []string{}, fmt.Errorf("could not read file '%s': %w", filename, err)
	}

	return strings.Split(string(fileReadByte), "\n"), nil
}

func cleanNullCommitFile(filename string, clt *junos.Client) error {
	if err := utils.ReplaceTildeToHomeDir(&filename); err != nil {
		return err
	}
	f, err := os.OpenFile(filename, os.O_TRUNC, os.FileMode(clt.FilePermission()))
	if err != nil {
		return fmt.Errorf("could not open file '%s' to truncate after commit: %w", filename, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close file handler for '%s' after truncation: %w", filename, err)
	}

	return nil
}
