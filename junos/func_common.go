package junos

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func logFile(message string, file string) {
	// create your file with desired read/write permissions
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	// set output of logs to f
	log.SetOutput(f)
	log.SetPrefix(time.Now().Format("2006-01-02 15:04:05"))

	log.Printf("%s", message)
}

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
		return fmt.Errorf("%v is not a valide IP/mask, is a network", ip)
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

func validateNameObjectJunos(exclude []string, length int) schema.SchemaValidateDiagFunc {
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
		f := func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' && r != '_'
		}
		if strings.IndexFunc(v, f) != -1 {
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

func checkStringHasPrefixInList(s string, list []string) bool {
	for _, item := range list {
		if strings.HasPrefix(s, item) {
			return true
		}
	}

	return false
}
