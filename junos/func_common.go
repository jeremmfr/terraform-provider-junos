package junos

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
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
		if !bchk.StringInSlice(octet, []string{"255", "254", "252", "248", "240", "224", "192", "128", "0"}) {
			return fmt.Errorf("mask %v must be in subnet mask format, octet [%v] is not", mask, octet)
		}
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
		if bchk.StringInSlice(v, exclude) {
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

		if err != nil || fileMode > 0o777 || fileMode < 0 {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("bad mode for file - must be three octal digits: %s", v),
				AttributePath: path,
			})
		}

		return diags
	}
}

func sortSetOfString(list []interface{}) []string {
	s := make([]string, 0)
	for _, e := range list {
		s = append(s, e.(string))
	}
	sort.Strings(s)

	return s
}

func copyAndRemoveItemMapList(identifier string,
	m map[string]interface{}, list []map[string]interface{}) []map[string]interface{} {
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

func checkCompatibilityRouter(jnprSess *NetconfObject) bool {
	if strings.HasPrefix(strings.ToLower(jnprSess.SystemInformation.HardwareModel), "mx") {
		return true
	}
	if strings.HasPrefix(strings.ToLower(jnprSess.SystemInformation.HardwareModel), "vmx") {
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

func stringLenBetweenSensitive(min, max int) schema.SchemaValidateDiagFunc {
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

		if len(v) < min || len(v) > max {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected length to be in the range (%d - %d), got %d", min, max, len(v)),
				AttributePath: path,
			})
		}

		return diags
	}
}
