package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

func resourceInterfaceSt0Unit() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceInterfaceSt0UnitCreate,
		ReadWithoutTimeout:   resourceInterfaceSt0UnitRead,
		DeleteWithoutTimeout: resourceInterfaceSt0UnitDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceInterfaceSt0UnitImport,
		},
	}
}

func resourceInterfaceSt0UnitCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	newSt0, err := searchInterfaceSt0UnitToCreate(clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("error to find new st0 unit interface: %w", err))...)
	}
	if err := clt.ConfigSet([]string{"set interfaces " + newSt0}, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_interface_st0_unit", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ncInt, emptyInt, setInt, err := checkInterfaceLogicalNCEmpty(newSt0, clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ncInt {
		return append(diagWarns, diag.FromErr(fmt.Errorf("create new %v always disable after commit "+
			"=> check your config", newSt0))...)
	}
	if emptyInt && !setInt {
		return append(diagWarns, diag.FromErr(fmt.Errorf("create new st0 unit interface doesn't works, "+
			"can't find the new interface %s after commit", newSt0))...)
	}
	d.SetId(newSt0)

	return diagWarns
}

func resourceInterfaceSt0UnitRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	mutex.Lock()
	ncInt, emptyInt, setInt, err := checkInterfaceLogicalNCEmpty(d.Id(), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ncInt || (emptyInt && !setInt) {
		d.SetId("")
	}

	return nil
}

func resourceInterfaceSt0UnitDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := clt.ConfigSet([]string{"delete interfaces " + d.Id()}, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ncInt, emptyInt, _, err := checkInterfaceLogicalNCEmpty(d.Id(), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !ncInt && !emptyInt {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("interface %s not empty or disable", d.Id()))...)
	}
	if err := clt.ConfigSet([]string{"delete interfaces " + d.Id()}, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_interface_st0_unit", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceInterfaceSt0UnitImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	if !strings.HasPrefix(d.Id(), "st0.") {
		return nil, fmt.Errorf("id must be start with 'st0.'")
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)
	ncInt, emptyInt, setInt, err := checkInterfaceLogicalNCEmpty(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if ncInt {
		return nil, fmt.Errorf("interface '%v' is disabled, import is not possible", d.Id())
	}
	if emptyInt && !setInt {
		return nil, fmt.Errorf("don't find interface with id '%v'"+
			" (id must be the name of st0 unit interface <st0.?>)", d.Id())
	}
	result[0] = d

	return result, nil
}

func searchInterfaceSt0UnitToCreate(clt *junos.Client, junSess *junos.Session) (string, error) {
	st0, err := clt.Command("show interfaces st0 terse", junSess)
	if err != nil {
		return "", err
	}
	st0Line := strings.Split(st0, "\n")
	st0int := make([]string, 0)
	for _, line := range st0Line {
		if strings.HasPrefix(line, "st0.") {
			lineFields := strings.Split(line, " ")
			st0int = append(st0int, lineFields[0])
		}
	}
	for i := 0; i <= 1073741823; i++ {
		if !bchk.InSlice("st0."+strconv.Itoa(i), st0int) {
			return "st0." + strconv.Itoa(i), nil
		}
	}

	return "", fmt.Errorf("error for find st0 unit to create")
}