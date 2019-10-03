package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type policyOptions struct {
	fromZone string
	toZone   string
	policy   []map[string]interface{}
}

func resourceSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityPolicyCreate,
		Read:   resourceSecurityPolicyRead,
		Update: resourceSecurityPolicyUpdate,
		Delete: resourceSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"from_zone": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"to_zone": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"policy": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateNameObjectJunos(),
						},
						"match_source_address": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_destination_address": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_application": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"then": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "permit",
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"permit", "reject", "deny"}) {
									errors = append(errors, fmt.Errorf(
										"%q %q invalid action", value, k))
								}
								return
							},
						},
						"permit_tunnel_ipsec_vpn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"count": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"log_init": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"log_close": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceSecurityPolicyCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security policy not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	securityPolicyExists, err := checkSecurityPolicyExists(d.Get("from_zone").(string), d.Get("to_zone").(string),
		m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if securityPolicyExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("security policy from %v to %v already exists",
			d.Get("from_zone").(string), d.Get("to_zone").(string))
	}

	err = setSecurityPolicy(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	securityPolicyExists, err = checkSecurityPolicyExists(d.Get("from_zone").(string), d.Get("to_zone").(string),
		m, jnprSess)
	if err != nil {
		return err
	}
	if securityPolicyExists {
		d.SetId(d.Get("from_zone").(string) + idSeparator + d.Get("to_zone").(string))
	} else {
		return fmt.Errorf("security policy from %v to %v not exists after commit "+
			"=> check your config", d.Get("from_zone").(string), d.Get("to_zone").(string))
	}
	return resourceSecurityPolicyRead(d, m)
}
func resourceSecurityPolicyRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	policyOptions, err := readSecurityPolicy(d.Get("from_zone").(string)+idSeparator+d.Get("to_zone").(string),
		m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if len(policyOptions.policy) == 0 {
		d.SetId("")
	} else {
		fillSecurityPolicyData(d, policyOptions)
	}
	return nil
}
func resourceSecurityPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}

	err = delSecurityPolicy(d.Get("from_zone").(string), d.Get("to_zone").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}

	err = setSecurityPolicy(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourceSecurityPolicyRead(d, m)
}
func resourceSecurityPolicyDelete(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delSecurityPolicy(d.Get("from_zone").(string), d.Get("to_zone").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
}
func resourceSecurityPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idList := strings.Split(d.Id(), idSeparator)
	if len(idList) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	securityPolicyExists, err := checkSecurityPolicyExists(idList[0], idList[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityPolicyExists {
		return nil, fmt.Errorf("don't find policy with id '%v' (id must be <from_zone>"+idSeparator+"<to_zone>)", d.Id())
	}
	policyOptions, err := readSecurityPolicy(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityPolicyData(d, policyOptions)

	result[0] = d
	return result, nil

}

func checkSecurityPolicyExists(fromZone, toZone string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	zoneConfig, err := sess.command("show configuration"+
		" security policies from-zone "+fromZone+" to-zone "+toZone+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if zoneConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setSecurityPolicy(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security policies" +
		" from-zone " + d.Get("from_zone").(string) +
		" to-zone " + d.Get("to_zone").(string) +
		" policy "
	for _, v := range d.Get("policy").([]interface{}) {
		policy := v.(map[string]interface{})
		setPrefixPolicy := setPrefix + policy["name"].(string)
		if len(policy["match_source_address"].([]interface{})) != 0 {
			for _, address := range policy["match_source_address"].([]interface{}) {
				configSet = append(configSet, setPrefixPolicy+" match source-address "+address.(string)+"\n")
			}
		} else {
			configSet = append(configSet, setPrefixPolicy+" match source-address any\n")
		}
		if len(policy["match_destination_address"].([]interface{})) != 0 {
			for _, address := range policy["match_destination_address"].([]interface{}) {
				configSet = append(configSet, setPrefixPolicy+" match destination-address "+address.(string)+"\n")
			}
		} else {
			configSet = append(configSet, setPrefixPolicy+" match destination-address any\n")
		}
		if len(policy["match_application"].([]interface{})) != 0 {
			for _, app := range policy["match_application"].([]interface{}) {
				configSet = append(configSet, setPrefixPolicy+" match application "+app.(string)+"\n")
			}
		} else {
			configSet = append(configSet, setPrefixPolicy+" match application any\n")
		}
		configSet = append(configSet, setPrefixPolicy+" then "+policy["then"].(string)+"\n")
		if policy["permit_tunnel_ipsec_vpn"].(string) != "" {
			if policy["then"].(string) != "permit" {
				return fmt.Errorf("conflict policy then %v and policy permit_tunnel_ipsec_vpn",
					policy["then"].(string))
			}
			configSet = append(configSet, setPrefixPolicy+" then permit tunnel ipsec-vpn "+
				policy["permit_tunnel_ipsec_vpn"].(string)+"\n")
		}
		if policy["count"].(bool) {
			configSet = append(configSet, setPrefixPolicy+" then count\n")
		}
		if policy["log_init"].(bool) {
			configSet = append(configSet, setPrefixPolicy+" then log session-init\n")
		}
		if policy["log_close"].(bool) {
			configSet = append(configSet, setPrefixPolicy+" then log session-close\n")
		}
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readSecurityPolicy(idPolicy string, m interface{}, jnprSess *NetconfObject) (policyOptions, error) {
	zone := strings.Split(idPolicy, idSeparator)
	fromZone := zone[0]
	toZone := zone[1]

	sess := m.(*Session)
	var confRead policyOptions

	policyConfig, err := sess.command("show configuration"+
		" security policies from-zone "+fromZone+" to-zone "+toZone+" | display set relative ", jnprSess)
	if err != nil {
		return confRead, err
	}
	policyList := make([]map[string]interface{}, 0)
	if policyConfig != emptyWord {
		confRead.fromZone = fromZone
		confRead.toZone = toZone
		for _, item := range strings.Split(policyConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			if strings.Contains(itemTrim, " match ") || strings.Contains(itemTrim, " then ") {
				policyLineCut := strings.Split(itemTrim, " ")
				m := genMapPolicyWithName(policyLineCut[1])
				m, policyList = copyAndRemoveItemMapList("name", false, m, policyList)
				itemTrimPolicy := strings.TrimPrefix(itemTrim, "policy "+policyLineCut[1]+" ")
				switch {
				case strings.HasPrefix(itemTrimPolicy, "match source-address "):
					m["match_source_address"] = append(m["match_source_address"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match source-address "))
				case strings.HasPrefix(itemTrimPolicy, "match destination-address "):
					m["match_destination_address"] = append(m["match_destination_address"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match destination-address "))
				case strings.HasPrefix(itemTrimPolicy, "match application "):
					m["match_application"] = append(m["match_application"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match application "))
				case strings.HasPrefix(itemTrimPolicy, "then "):
					switch {
					case strings.HasSuffix(itemTrimPolicy, "permit"),
						strings.HasSuffix(itemTrimPolicy, "deny"),
						strings.HasSuffix(itemTrimPolicy, "reject"):
						m["then"] = strings.TrimPrefix(itemTrimPolicy, "then ")
					case strings.HasSuffix(itemTrimPolicy, "count"):
						m["count"] = true
					case strings.HasSuffix(itemTrimPolicy, "log session-init"):
						m["log_init"] = true
					case strings.HasSuffix(itemTrimPolicy, "log session-close"):
						m["log_close"] = true
					case strings.HasPrefix(itemTrimPolicy, "then permit tunnel ipsec-vpn "):
						m["permit_tunnel_ipsec_vpn"] = strings.TrimPrefix(itemTrimPolicy,
							"then permit tunnel ipsec-vpn ")
					}
				}
				policyList = append(policyList, m)
			}
		}
	}
	confRead.policy = policyList
	return confRead, nil
}
func delSecurityPolicy(fromZone string, toZone string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security policies from-zone "+fromZone+" to-zone "+toZone+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}

func fillSecurityPolicyData(d *schema.ResourceData, policyOptions policyOptions) {
	tfErr := d.Set("from_zone", policyOptions.fromZone)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("to_zone", policyOptions.toZone)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("policy", policyOptions.policy)
	if tfErr != nil {
		panic(tfErr)
	}
}

func genMapPolicyWithName(name string) map[string]interface{} {
	return map[string]interface{}{
		"name":                      name,
		"match_source_address":      make([]string, 0),
		"match_destination_address": make([]string, 0),
		"match_application":         make([]string, 0),
		"then":                      "",
		"count":                     false,
		"log_init":                  false,
		"log_close":                 false,
		"permit_tunnel_ipsec_vpn":   "",
	}
}
