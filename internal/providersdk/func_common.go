package providersdk

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func appendDiagWarns(diags *diag.Diagnostics, warns []error) {
	for _, w := range warns {
		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  w.Error(),
		})
	}
}

func validateNameObjectJunos(length int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		v := i.(string)
		if strings.Count(v, "") > length {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("%s invalid name (too long)", i),
				AttributePath: path,
			})
		}

		resultRune := strings.IndexFunc(v, func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' && r != '_'
		})
		if resultRune != -1 {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("%s invalid name (bad character)", i),
				AttributePath: path,
			})
		}

		return diags
	}
}

func copyAndRemoveItemMapList(
	identifier string, m map[string]interface{}, list []map[string]interface{},
) []map[string]interface{} {
	if m[identifier] == nil {
		panic(fmt.Errorf("internal error: can't find identifier %s in map", identifier))
	}
	for i, element := range list {
		if element[identifier] == m[identifier] {
			for key, value := range element {
				m[key] = value
			}
			list = append(list[:i], list[i+1:]...)

			break
		}
	}

	return list
}
