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
				configSet = append(configSet, setPrefix+"esi all-active")
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
	identifier, _ := regexp.MatchString(`^(\d+:){9,9}\d+`, itemTrim)
	if identifier {
		grRead["identifier"] = itemTrim
	}
	if err != nil {
		return []map[string]interface{}{grRead}, fmt.Errorf("an error occurred: %s", itemTrim)
	}

	return []map[string]interface{}{grRead}, nil
}
