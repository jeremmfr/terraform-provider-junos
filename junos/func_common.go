package junos

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
func validateIPMaskFunc() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v := i.(string)
		err := validateIPwithMask(v)
		if err != nil {
			es = append(es, err)
		}

		return
	}
}
func validateIPFunc() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v := i.(string)
		err := validateIP(v)
		if err != nil {
			es = append(es, err)
		}

		return
	}
}
func validateNetworkFunc() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v := i.(string)
		err := validateNetwork(v)
		if err != nil {
			es = append(es, err)
		}

		return
	}
}

func validateIP(ip string) error {
	ipMask := ip + "/32"
	_, ipnet, err := net.ParseCIDR(ipMask)
	if err != nil || ipnet == nil {
		return fmt.Errorf("%v is not a valid IP", ip)
	}

	return nil
}

func validateIPwithMask(ip string) error {
	if !strings.Contains(ip, "/") {
		return fmt.Errorf("%v missing mask", ip)
	}
	_, ipnet, err := net.ParseCIDR(ip)
	if err != nil || ipnet == nil {
		return fmt.Errorf("%v is not a valid IP/mask", ip)
	}
	if (strings.Contains(ip, ":") && strings.Contains(ip, "/128")) ||
		(!strings.Contains(ip, ":") && strings.Contains(ip, "/32")) {
		return nil
	}
	if ip == ipnet.String() {
		return fmt.Errorf("%v is not a valide IP/mask, is a network", ip)
	}

	return nil
}

func validateNetwork(network string) error {
	if !strings.Contains(network, "/") {
		return fmt.Errorf("%v missing mask", network)
	}
	_, ipnet, err := net.ParseCIDR(network)
	if err != nil || ipnet == nil {
		return fmt.Errorf("%v is not a valid IP/mask", network)
	}
	if network != ipnet.String() {
		return fmt.Errorf("%v is not a network, is a IP/mask", network)
	}

	return nil
}

func validateNameObjectJunos() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v := i.(string)
		if strings.Count(v, "") > 32 {
			es = append(es, fmt.Errorf(
				"%q %q invalid name (too long)", k, i))
		}
		f := func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' && r != '_'
		}
		if strings.IndexFunc(v, f) != -1 {
			es = append(es, fmt.Errorf(
				"%q %q invalid name (bad character)", k, i))
		}

		return
	}
}
func validateAddress() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v := i.(string)

		f := func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' && r != '.'
		}
		if strings.IndexFunc(v, f) != -1 {
			es = append(es, fmt.Errorf(
				"%q %q invalid address (bad character)", k, i))
		}

		return
	}
}

func validateIntRange(start int, end int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < start || value > end {
			errors = append(errors, fmt.Errorf(
				"%d for %q is not valid (%d-%d)", value, k, start, end))
		}

		return
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
	if strings.HasPrefix(strings.ToLower(jnprSess.Platform[0].Model), "srx") {
		return true
	}
	if strings.HasPrefix(strings.ToLower(jnprSess.Platform[0].Model), "vsrx") {
		return true
	}
	if strings.HasPrefix(strings.ToLower(jnprSess.Platform[0].Model), "j") {
		return true
	}

	return false
}

func validateSyslogSeverity() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if !stringInSlice(value, []string{"alert", "any", "critical",
			"emergency", "error", "info", "none", "notice", "warning"}) {
			errors = append(errors, fmt.Errorf(
				"%q %q invalid severity", value, k))
		}

		return
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
