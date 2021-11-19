package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInterfacePhysicalDisable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInterfacePhysicalDisableCreate,
		ReadContext:   resourceInterfacePhysicalDisableRead,
		DeleteContext: resourceInterfacePhysicalDisableDelete,
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

func resourceInterfacePhysicalDisableCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := addInterfacePhysicalNC(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	ncInt, emptyInt, err := checkInterfacePhysicalNCEmpty(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !ncInt && !emptyInt {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("interface %s is configured", d.Get("name").(string)))...)
	}
	if ncInt {
		d.SetId(d.Get("name").(string))
		if errs := sess.configClear(jnprSess); len(errs) > 0 {
			return diagWarns
		}

		return nil
	}
	if emptyInt {
		if containsUnit, err := checkInterfacePhysicalContainsUnit(d.Get("name").(string), m, jnprSess); err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		} else if containsUnit {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(
				fmt.Errorf("interface %s is used for a logical unit interface", d.Get("name").(string)))...)
		}
	}
	if err := addInterfacePhysicalNC(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_interface_physical_disable", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ncInt, _, err = checkInterfacePhysicalNCEmpty(d.Get("name").(string), m, jnprSess)
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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	mutex.Lock()
	ncInt, _, err := checkInterfacePhysicalNCEmpty(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if !ncInt {
		d.SetId("")
	}

	return nil
}

func resourceInterfacePhysicalDisableDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
