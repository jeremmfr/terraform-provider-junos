package junos

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInterfacePhysical() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInterfacePhysicalRead,
		Schema: map[string]*schema.Schema{
			"config_interface": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"match": {
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ae_lacp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ae_link_speed": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ae_minimum_links": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ether802_3ad": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trunk": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"vlan_members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vlan_native": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vlan_tagging": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceInterfacePhysicalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("config_interface").(string) == "" && d.Get("match").(string) == "" {
		return diag.FromErr(fmt.Errorf("no arguments provided, 'config_interface' and 'match' empty"))
	}
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	mutex.Lock()
	nameFound, err := searchInterfacePhysicalID(d.Get("config_interface").(string), d.Get("match").(string), m, jnprSess)
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	if nameFound == "" {
		mutex.Unlock()

		return diag.FromErr(fmt.Errorf("no physical interface found with arguments provided"))
	}
	interfaceOpt, err := readInterfacePhysical(nameFound, m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(nameFound)
	if tfErr := d.Set("name", nameFound); tfErr != nil {
		panic(tfErr)
	}
	fillInterfacePhysicalData(d, interfaceOpt)

	return nil
}

func searchInterfacePhysicalID(configInterface string, match string,
	m interface{}, jnprSess *NetconfObject) (string, error) {
	sess := m.(*Session)
	intConfigList := make([]string, 0)
	intConfig, err := sess.command("show configuration interfaces "+configInterface+" | display set", jnprSess)
	if err != nil {
		return "", err
	}
	for _, item := range strings.Split(intConfig, "\n") {
		if strings.Contains(item, "<configuration-output>") {
			continue
		}
		if strings.Contains(item, "</configuration-output>") {
			break
		}
		if item == "" {
			continue
		}
		itemTrim := strings.TrimPrefix(item, "set interfaces ")
		matched, err := regexp.MatchString(match, itemTrim)
		if err != nil {
			return "", fmt.Errorf("failed to regexp with %s : %w", match, err)
		}
		if !matched {
			continue
		}
		itemTrimSplit := strings.Split(itemTrim, " ")
		if len(itemTrimSplit) == 0 {
			continue
		}
		intConfigList = append(intConfigList, itemTrimSplit[0])
	}
	intConfigList = uniqueListString(intConfigList)
	if len(intConfigList) == 0 {
		return "", nil
	}
	if len(intConfigList) > 1 {
		return "", fmt.Errorf("too many different physical interfaces found")
	}

	return intConfigList[0], nil
}
