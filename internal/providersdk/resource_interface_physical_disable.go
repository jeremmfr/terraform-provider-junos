package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

func resourceInterfacePhysicalDisable() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceInterfacePhysicalDisableCreate,
		ReadWithoutTimeout:   resourceInterfacePhysicalDisableRead,
		DeleteWithoutTimeout: resourceInterfacePhysicalDisableDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if strings.Count(value, ".") > 0 {
						errors = append(errors, fmt.Errorf(
							"%q in %q cannot have a dot", value, k))
					}

					return
				},
			},
		},
	}
}

func resourceInterfacePhysicalDisableCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := addInterfacePhysicalNC(d.Get("name").(string), clt.GroupInterfaceDelete(), junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

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
	ncInt, emptyInt, err := checkInterfacePhysicalNCEmpty(d.Get("name").(string), clt.GroupInterfaceDelete(), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !ncInt && !emptyInt {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("interface %s is configured", d.Get("name").(string)))...)
	}
	if ncInt {
		d.SetId(d.Get("name").(string))
		if errs := junSess.ConfigClear(); len(errs) > 0 {
			return diagWarns
		}

		return nil
	}
	if emptyInt {
		if containsUnit, err := checkInterfacePhysicalContainsUnit(d.Get("name").(string), junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		} else if containsUnit {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(
				fmt.Errorf("interface %s is used for a logical unit interface", d.Get("name").(string)))...)
		}
	}
	if err := addInterfacePhysicalNC(d.Get("name").(string), clt.GroupInterfaceDelete(), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_interface_physical_disable")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	ncInt, _, err = checkInterfacePhysicalNCEmpty(d.Get("name").(string), clt.GroupInterfaceDelete(), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if !ncInt {
		return append(diagWarns, diag.FromErr(fmt.Errorf("interface %v always not disable after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}
	d.SetId(d.Get("name").(string))

	return nil
}

func resourceInterfacePhysicalDisableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	mutex.Lock()
	ncInt, _, err := checkInterfacePhysicalNCEmpty(d.Get("name").(string), clt.GroupInterfaceDelete(), junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if !ncInt {
		d.SetId("")
	}

	return nil
}

func resourceInterfacePhysicalDisableDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	return nil
}
