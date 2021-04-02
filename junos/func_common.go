package junos

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FormatName int

const (
	FormatDefault FormatName = iota
	FormatAddressName
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
		return fmt.Errorf("%v is not a valid IP/mask", ip)
	}
	if (strings.Contains(ip, ":") && (strings.Contains(ip, "/128") || strings.Contains(ip, "/127"))) ||
		(!strings.Contains(ip, ":") && (strings.Contains(ip, "/32") || strings.Contains(ip, "/31"))) {
		return nil
	}
	if ip == ipnet.String() {
		return fmt.Errorf("%v is not a valid IP/mask, is a network", ip)
	}

	return nil
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
func validateCIDR(cidr string) error {
	if !strings.Contains(cidr, "/") {
		return fmt.Errorf("%v missing mask", cidr)
	}
	if _, _, err := net.ParseCIDR(cidr); err != nil {
		return fmt.Errorf("%v is not a valid CIDR", cidr)
	}

	return nil
}

func validateWildcardFunc() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		v := i.(string)
		err := validateWildcardWithMask(v)
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

func validateWildcardWithMask(wildcard string) error {
	if !strings.Contains(wildcard, "/") {
		return fmt.Errorf("%v missing mask", wildcard)
	}
	if strings.Contains(wildcard, ":") {
		return fmt.Errorf("wildcards do not support IPv6 addresses, %v is IPv6", wildcard)
	}
	wildcardSplit := strings.Split(wildcard, "/")
	ip := net.ParseIP(wildcardSplit[0]).To4()
	if ip == nil {
		return fmt.Errorf("ip %v not a valid ip address", wildcardSplit[0])
	}
	mask := net.ParseIP(wildcardSplit[1]).To4()
	if mask == nil {
		return fmt.Errorf("mask %v is improperly formatted, must be in x.x.x.x notation", wildcardSplit[1])
	}
	for _, octet := range strings.Split(mask.String(), ".") {
		if !stringInSlice(octet, []string{"255", "254", "252", "248", "240", "224", "192", "128", "0"}) {
			return fmt.Errorf("mask %v must be in subnet mask format, octet [%v] is not", mask, octet)
		}
	}

	return nil
}

func validateNameObjectJunos(exclude []string, length int, format FormatName) schema.SchemaValidateDiagFunc {
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
		resultRune := -1
		switch format {
		case FormatDefault:
			resultRune = strings.IndexFunc(v, f1)
		case FormatAddressName:
			resultRune = strings.IndexFunc(v, f2)
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
		if stringInSlice(v, exclude) {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected value to not be one of %q, got %v", exclude, i),
				AttributePath: path,
			})
		}

		return diags
	}
}
func validateAddress() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		v := i.(string)

		f := func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' && r != '.'
		}
		if strings.IndexFunc(v, f) != -1 {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("%s invalid address (bad character)", v),
				AttributePath: path,
			})
		}

		return diags
	}
}
func validateFilePermission() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		v, ok := i.(string)

		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "expected type to be string",
				AttributePath: path,
			})

			return diags
		}

		if len(v) > 4 || len(v) < 3 {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("bad mode for file - string length should be 3 or 4 digits: %s", v),
				AttributePath: path,
			})
		}

		fileMode, err := strconv.ParseInt(v, 8, 64)

		if err != nil || fileMode > 0777 || fileMode < 0 {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("bad mode for file - must be three octal digits: %s", v),
				AttributePath: path,
			})
		}

		return diags
	}
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}

	return false
}
func copyAndRemoveItemMapList(identifier string, integer bool,
	m map[string]interface{}, list []map[string]interface{}) (map[string]interface{}, []map[string]interface{}) {
	for i, element := range list {
		if integer {
			if element[identifier].(int) == m[identifier].(int) {
				for key, value := range element {
					m[key] = value
				}
				list = append(list[:i], list[i+1:]...)

				break
			}
		} else {
			if element[identifier].(string) == m[identifier].(string) {
				for key, value := range element {
					m[key] = value
				}
				list = append(list[:i], list[i+1:]...)

				break
			}
		}
	}

	return m, list
}

func checkCompatibilitySecurity(jnprSess *NetconfObject) bool {
	if strings.HasPrefix(strings.ToLower(jnprSess.SystemInformation.HardwareModel), "srx") {
		return true
	}
	if strings.HasPrefix(strings.ToLower(jnprSess.SystemInformation.HardwareModel), "vsrx") {
		return true
	}
	if strings.HasPrefix(strings.ToLower(jnprSess.SystemInformation.HardwareModel), "j") {
		return true
	}

	return false
}

func listOfSyslogSeverity() []string {
	return []string{
		"alert", "any", "critical",
		"emergency", "error", "info", "none", "notice", "warning",
	}
}
func listOfSyslogFacility() []string {
	return []string{
		"authorization", "daemon", "ftp", "kernel", "user",
		"local0", "local1", "local2", "local3", "local4", "local5", "local6", "local7",
	}
}

func uniqueListString(s []string) []string {
	k := make(map[string]bool)
	r := []string{}
	for _, v := range s {
		if _, value := k[v]; !value {
			k[v] = true
			r = append(r, v)
		}
	}

	return r
}

type sortStringsLength []string

func (s sortStringsLength) Len() int {
	return len(s)
}
func (s sortStringsLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sortStringsLength) Less(i, j int) bool {
	if len(s[i]) < len(s[j]) {
		return true
	}
	if len(s[j]) < len(s[i]) {
		return false
	}

	return s[i] < s[j]
}

func checkStringHasPrefixInList(s string, list []string) bool {
	for _, item := range list {
		if strings.HasPrefix(s, item) {
			return true
		}
	}

	return false
}

func replaceTildeToHomeDir(path *string) error {
	if strings.HasPrefix(*path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to read user home directory : %w", err)
		}
		*path = homeDir + strings.TrimPrefix(*path, "~")
	}

	return nil
}
