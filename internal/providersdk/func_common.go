package providersdk

import (
	"fmt"
	"net"
	"slices"
	"sort"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type formatName int

const (
	formatDefault formatName = iota
	formatAddressName
	formatDefAndDots
)

func appendDiagWarns(diags *diag.Diagnostics, warns []error) {
	for _, w := range warns {
		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  w.Error(),
		})
	}
}

func validateIPMaskFunc() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		v := i.(string)
		err := validateIPwithMask(v)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       err.Error(),
				AttributePath: path,
			})
		}

		return diags
	}
}

func validateIPwithMask(ip string) error {
	if !strings.Contains(ip, "/") {
		return fmt.Errorf("%v missing mask", ip)
	}
	_, ipnet, err := net.ParseCIDR(ip)
	if err != nil || ipnet == nil {
		return fmt.Errorf("%v is not a valid CIDR", ip)
	}

	return nil
}

func validateCIDRNetworkFunc() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		v := i.(string)
		err := validateCIDRNetwork(v)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       err.Error(),
				AttributePath: path,
			})
		}

		return diags
	}
}

func validateCIDRNetwork(network string) error {
	if !strings.Contains(network, "/") {
		return fmt.Errorf("%v missing mask", network)
	}
	_, ipnet, err := net.ParseCIDR(network)
	if err != nil || ipnet == nil {
		return fmt.Errorf("%v is not a valid CIDR", network)
	}
	if network != ipnet.String() {
		return fmt.Errorf("%v is not a valid network CIDR", network)
	}

	return nil
}

func validateNameObjectJunos(exclude []string, length int, format formatName) schema.SchemaValidateDiagFunc {
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
		f1 := func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' && r != '_'
		}
		f2 := func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') &&
				r != '-' && r != '_' && r != ':' && r != '.' && r != '/'
		}
		f3 := func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') &&
				r != '-' && r != '_' && r != '.'
		}
		resultRune := -1
		switch format {
		case formatDefault:
			resultRune = strings.IndexFunc(v, f1)
		case formatAddressName:
			resultRune = strings.IndexFunc(v, f2)
		case formatDefAndDots:
			resultRune = strings.IndexFunc(v, f3)
		default:
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "internal error: validateNameObjectJunos function called with a bad argument",
				AttributePath: path,
			})
		}
		if resultRune != -1 {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("%s invalid name (bad character)", i),
				AttributePath: path,
			})
		}
		if slices.Contains(exclude, v) {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected value to not be one of %q, got %v", exclude, i),
				AttributePath: path,
			})
		}

		return diags
	}
}

func sortSetOfString(list []interface{}) []string {
	s := make([]string, len(list))
	for k, e := range list {
		s[k] = e.(string)
	}
	sort.Strings(s)

	return s
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

func validateIsIPv6Address(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))

		return warnings, errors
	}

	ip := net.ParseIP(v)
	if four, six := ip.To4(), ip.To16(); four != nil || six == nil {
		errors = append(errors, fmt.Errorf("expected %s to contain a valid IPv6 address, got: %s", k, v))
	}

	return warnings, errors
}
