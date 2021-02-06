package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInterfaceSt0Unit() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInterfaceSt0UnitCreate,
		ReadContext:   resourceInterfaceSt0UnitRead,
		DeleteContext: resourceInterfaceSt0UnitDelete,
		Importer: &schema.ResourceImporter{
			State: resourceInterfaceSt0UnitImport,
		},
	}
}

func resourceInterfaceSt0UnitCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	newSt0, err := searchInterfaceSt0UnitToCreate(m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("error for find new st0 unit interface: %w", err))
	}
	if err := sess.configSet([]string{"set interfaces " + newSt0}, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_interface_st0_unit", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	intExists, err := checkInterfaceExists(newSt0, m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if intExists {
		d.SetId(newSt0)
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("create new st0 unit interface doesn't works, "+
			"can't find the new interface %s after commit", newSt0))...)
	}

	return diagWarns
}
func resourceInterfaceSt0UnitRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	mutex.Lock()
	intExists, err := checkInterfaceExists(d.Id(), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if !intExists {
		d.SetId("")
	}

	return nil
}

func resourceInterfaceSt0UnitDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	ncInt, emptyInt, err := checkInterfaceLogicalNC(d.Id(), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if !ncInt && !emptyInt {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("interface %s not empty or disable", d.Id()))
	}
	if err := sess.configSet([]string{"delete interfaces " + d.Id()}, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_interface_st0_unit", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceInterfaceSt0UnitImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	if !strings.HasPrefix(d.Id(), "st0.") {
		return nil, fmt.Errorf("id must be start with 'st0.'")
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	intExists, err := checkInterfaceExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !intExists {
		return nil, fmt.Errorf("don't find interface with id '%v'"+
			" (id must be the name of st0 unit interface <st0.?>)", d.Id())
	}
	result[0] = d

	return result, nil
}

func searchInterfaceSt0UnitToCreate(m interface{}, jnprSess *NetconfObject) (string, error) {
	sess := m.(*Session)
	st0, err := sess.command("show interfaces st0 terse", jnprSess)
	if err != nil {
		return "", err
	}
	st0Line := strings.Split(st0, "\n")
	st0int := make([]string, 0)
	for _, line := range st0Line {
		if strings.HasPrefix(line, "st0.") {
			lineSplit := strings.Split(line, " ")
			st0int = append(st0int, lineSplit[0])
		}
	}
	for i := 0; i <= 1073741823; i++ {
		if !stringInSlice("st0."+strconv.Itoa(i), st0int) {
			return "st0." + strconv.Itoa(i), nil
		}
	}

	return "", fmt.Errorf("error for find st0 unit to create")
}
