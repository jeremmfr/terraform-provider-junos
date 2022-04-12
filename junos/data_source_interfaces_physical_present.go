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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	mutex.Lock()
	iPresent, err := searchInterfacesPhysicalPresent(d, m, jnprSess)
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

func searchInterfacesPhysicalPresent(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
) (interfacesPresentOpts, error) {
	sess := m.(*Session)
	var result interfacesPresentOpts
	replyData, err := sess.commandXML(rpcGetInterfaceInformationTerse, jnprSess)
	if err != nil {
		return result, err
	}
	var iface getInterfaceTerseReply
	err = xml.Unmarshal([]byte(replyData), &iface.InterfaceInfo)
	if err != nil {
		return result, fmt.Errorf("failed to xml unmarshal reply data '%s': %w", replyData, err)
	}
	for _, iFace := range iface.InterfaceInfo.PhysicalInterface {
		if mName := d.Get("match_name").(string); mName != "" {
			matched, err := regexp.MatchString(mName, strings.Trim(iFace.Name, " \n\t"))
			if err != nil {
				return result, fmt.Errorf("failed to regexp with '%s': %w", mName, err)
			}
			if !matched {
				continue
			}
		}
		if d.Get("match_admin_up").(bool) && strings.Trim(iFace.AdminStatus, " \n\t") != "up" {
			continue
		}
		if d.Get("match_oper_up").(bool) && strings.Trim(iFace.OperStatus, " \n\t") != "up" {
			continue
		}
		result.interfaceNames = append(result.interfaceNames, strings.Trim(iFace.Name, " \n\t"))
		result.interfaceStatuses = append(result.interfaceStatuses, map[string]interface{}{
			"name":         strings.Trim(iFace.Name, " \n\t"),
			"admin_status": strings.Trim(iFace.AdminStatus, " \n\t"),
			"oper_status":  strings.Trim(iFace.OperStatus, " \n\t"),
		})
	}

	return result, nil
}
