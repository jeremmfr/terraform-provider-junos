package junos

import (
	"fmt"
	"regexp"
	"strings"
)

func setIntEsi(setPrefix string, esiParams []interface{},
	m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	for _, v := range esiParams {
		if v != nil {
			m := v.(map[string]interface{})
			if m["identifier"].(string) != "" {
				configSet = append(configSet, setPrefix+"esi "+m["identifier"].(string))
			}
			if m["mode"].(string) != "" {
				configSet = append(configSet, setPrefix+"esi "+m["mode"].(string))
			}
			if m["auto_derive_lacp"].(bool) {
				configSet = append(configSet, setPrefix+"esi auto-derive lacp")
			}
			if m["df_election_type"].(string) != "" {
				configSet = append(configSet, setPrefix+"esi df-election-type "+m["df_election_type"].(string))
			}
			if m["source_bmac"].(string) != "" {
				configSet = append(configSet, setPrefix+"esi source-bmac "+m["source_bmac"].(string))
			}
		}
	}
	if len(configSet) > 0 {
		err := sess.configSet(configSet, jnprSess)
		if err != nil {
			return err
		}
	}

	return nil
}

func readIntEsi(item string, grOpts []map[string]interface{}) ([]map[string]interface{}, error) {
	itemTrim := strings.TrimPrefix(item, "esi ")
	grRead := map[string]interface{}{
		"identifier": "",
	}
	if len(grOpts) > 0 {
		for k, v := range grOpts[0] {
			grRead[k] = v
		}
	}
	var err error
	identifier, _ := regexp.MatchString(`^([\d\w]{2}:){9}[\d\w]{2}`, itemTrim)
	esiMode, _ := regexp.MatchString(`^(all-active|single-active)`, itemTrim)
	if identifier {
		grRead["identifier"] = itemTrim
	}
	if esiMode {
		grRead["mode"] = itemTrim
	}
	if strings.HasPrefix(itemTrim, "df-election-type ") {
                grRead["df_election_type"] = strings.TrimPrefix(itemTrim, "df-election-type ")
        }
	if strings.HasPrefix(itemTrim, "source-bmac ") {
                grRead["source_bmac"] = strings.TrimPrefix(itemTrim, "source-bmac ")
        }
        if itemTrim == "auto-derive lacp" {
                grRead["auto_derive_lacp"] = true
        }
	if err != nil {
		return []map[string]interface{}{grRead}, fmt.Errorf("an error occurred: %s", itemTrim)
	}

	return []map[string]interface{}{grRead}, nil
}
