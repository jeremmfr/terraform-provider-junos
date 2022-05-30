package junos

import (
	"context"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type interfacesPresentOpts struct {
	interfaceNames    []string
	interfaceStatuses []map[string]interface{}
}

func dataSourceInterfacesPhysicalPresent() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceInterfacesPhysicalPresentRead,
		Schema: map[string]*schema.Schema{
			"match_name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if _, err := regexp.Compile(value); err != nil {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not valid regexp", value, k))
					}

					return
				},
			},
			"match_admin_up": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"match_oper_up": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"interface_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"interface_statuses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"admin_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oper_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceInterfacesPhysicalPresentRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	mutex.Lock()
	iPresent, err := searchInterfacesPhysicalPresent(d, clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if tfErr := d.Set("interface_names", iPresent.interfaceNames); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interface_statuses", iPresent.interfaceStatuses); tfErr != nil {
		panic(tfErr)
	}
	idString := "match=" + d.Get("match_name").(string)
	if d.Get("match_admin_up").(bool) {
		idString += idSeparator + "admin_up=true"
	}
	if d.Get("match_oper_up").(bool) {
		idString += idSeparator + "oper_up=true"
	}
	d.SetId(idString)

	return nil
}

func searchInterfacesPhysicalPresent(d *schema.ResourceData, clt *Client, junSess *junosSession,
) (interfacesPresentOpts, error) {
	var result interfacesPresentOpts
	replyData, err := clt.commandXML(rpcGetInterfacesInformationTerse, junSess)
	if err != nil {
		return result, err
	}
	var iface getPhysicalInterfaceTerseReply
	err = xml.Unmarshal([]byte(replyData), &iface.InterfaceInfo)
	if err != nil {
		return result, fmt.Errorf("failed to xml unmarshal reply data '%s': %w", replyData, err)
	}
	for _, iFace := range iface.InterfaceInfo.PhysicalInterface {
		if mName := d.Get("match_name").(string); mName != "" {
			matched, err := regexp.MatchString(mName, strings.TrimSpace(iFace.Name))
			if err != nil {
				return result, fmt.Errorf("failed to regexp with '%s': %w", mName, err)
			}
			if !matched {
				continue
			}
		}
		if d.Get("match_admin_up").(bool) && strings.TrimSpace(iFace.AdminStatus) != "up" {
			continue
		}
		if d.Get("match_oper_up").(bool) && strings.TrimSpace(iFace.OperStatus) != "up" {
			continue
		}
		result.interfaceNames = append(result.interfaceNames, strings.TrimSpace(iFace.Name))
		result.interfaceStatuses = append(result.interfaceStatuses, map[string]interface{}{
			"name":         strings.TrimSpace(iFace.Name),
			"admin_status": strings.TrimSpace(iFace.AdminStatus),
			"oper_status":  strings.TrimSpace(iFace.OperStatus),
		})
	}

	return result, nil
}
