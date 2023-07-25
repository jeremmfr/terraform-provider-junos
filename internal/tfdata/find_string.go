package tfdata

import "strings"

func FirstElementOfJunosLine(itemTrim string) string {
	if !strings.HasPrefix(itemTrim, `"`) {
		return strings.Split(itemTrim, ` `)[0]
	}

	itemTrimFields := strings.Split(itemTrim, ` `)
	var first string
	for i, v := range itemTrimFields {
		first += v
		if i == 0 && v == `"` {
			first += ` `

			continue
		}
		if strings.HasSuffix(v, `"`) {
			return first
		}
		first += ` `
	}

	return first
}
