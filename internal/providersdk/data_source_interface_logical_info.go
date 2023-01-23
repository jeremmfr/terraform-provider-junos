package providersdk

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type interfaceLogicalInfo struct {
	adminStatus string
	name        string
	operStatus  string
	familyInet  []map[string]interface{}
	familyInet6 []map[string]interface{}
}

func dataSourceInterfaceLogicalInfo() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceInterfaceLogicalInfoRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if strings.Count(value, ".") != 1 {
						errors = append(errors, fmt.Errorf(
							"%q in %q need to have 1 dot", value, k))
					}

					return
				},
			},
			"admin_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"oper_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"family_inet": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_cidr": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"family_inet6": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_cidr": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceInterfaceLogicalInfoRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	mutex.Lock()
	ifaceInfo, err := readInterfaceLogicalInfo(d, clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillInterfaceLogicalInfo(d, ifaceInfo)
	d.SetId(ifaceInfo.name)

	return nil
}

func readInterfaceLogicalInfo(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) (interfaceLogicalInfo, error) {
	var result interfaceLogicalInfo
	replyData, err := clt.CommandXML(fmt.Sprintf(junos.RPCGetInterfaceInformationTerse, d.Get("name").(string)), junSess)
	if err != nil {
		return result, err
	}
	var iface junos.GetLogicalInterfaceTerseReply
	err = xml.Unmarshal([]byte(replyData), &iface.InterfaceInfo)
	if err != nil {
		return result, fmt.Errorf("failed to xml unmarshal reply data '%s': %w", replyData, err)
	}
	if len(iface.InterfaceInfo.LogicalInterface) == 0 {
		return result, fmt.Errorf("logical-interface not found in xml: %v", replyData)
	}
	ifaceInfo := iface.InterfaceInfo.LogicalInterface[0]
	result.name = strings.TrimSpace(ifaceInfo.Name)
	result.adminStatus = strings.TrimSpace(ifaceInfo.AdminStatus)
	result.operStatus = strings.TrimSpace(ifaceInfo.OperStatus)
	for _, family := range ifaceInfo.AddressFamily {
		switch strings.TrimSpace(family.Name) {
		case junos.InetW:
			if len(result.familyInet) == 0 {
				result.familyInet = append(result.familyInet, map[string]interface{}{
					"address_cidr": make([]string, 0, len(family.Address)),
				})
			}
			for _, address := range family.Address {
				result.familyInet[0]["address_cidr"] = append(
					result.familyInet[0]["address_cidr"].([]string),
					strings.TrimSpace(address.Local),
				)
			}
		case junos.Inet6W:
			if len(result.familyInet6) == 0 {
				result.familyInet6 = append(result.familyInet6, map[string]interface{}{
					"address_cidr": make([]string, 0, len(family.Address)),
				})
			}
			for _, address := range family.Address {
				result.familyInet6[0]["address_cidr"] = append(
					result.familyInet6[0]["address_cidr"].([]string),
					strings.TrimSpace(address.Local),
				)
			}
		}
	}

	return result, nil
}

func fillInterfaceLogicalInfo(d *schema.ResourceData, ifaceInfo interfaceLogicalInfo) {
	if tfErr := d.Set("admin_status", ifaceInfo.adminStatus); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("oper_status", ifaceInfo.operStatus); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_inet", ifaceInfo.familyInet); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_inet6", ifaceInfo.familyInet6); tfErr != nil {
		panic(tfErr)
	}
}
