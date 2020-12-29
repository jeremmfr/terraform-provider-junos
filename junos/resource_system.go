package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type systemOptions struct {
	nameServer                           []string
	services                             []map[string]interface{}
	syslog                               []map[string]interface{}
	tracingDestinationOverrideSyslogHost string
}

func resourceSystem() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemCreate,
		ReadContext:   resourceSystemRead,
		UpdateContext: resourceSystemUpdate,
		DeleteContext: resourceSystemDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSystemImport,
		},
		Schema: map[string]*schema.Schema{
			"name_server": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"services": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ssh": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authentication_order": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"ciphers": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"client_alive_count_max": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 255),
									},
									"client_alive_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 65535),
									},
									"connection_limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 250),
									},
									"fingerprint_hash": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"md5", "sha2-256"}, false),
									},
									"hostkey_algorithm": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"key_exchange": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"log_key_changes": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"macs": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"max_pre_authentication_packets": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(20, 2147483647),
									},
									"max_sessions_per_connection": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"no_passwords": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"services.0.ssh.0.no_public_keys"},
									},
									"no_public_keys": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"services.0.ssh.0.no_passwords"},
									},
									"port": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"protocol_version": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"rate_limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 250),
									},
									"root_login": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"allow", "deny", "deny-password"}, false),
									},
									"no_tcp_forwarding": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"tcp_forwarding": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"syslog": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"archive": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"binary_data": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"no_binary_data": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"files": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 1000),
									},
									"size": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(65536, 1073741824),
									},
									"no_world_readable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"world_readable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"log_rotate_frequency": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 59),
						},
						"source_address": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateAddress(),
						},
					},
				},
			},
			"tracing_dest_override_syslog_host": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateAddress(),
			},
		},
	}
}

func resourceSystemCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)

	if err := setSystem(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := sess.commitConf("create resource junos_system", jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}

	d.SetId("system")

	return resourceSystemRead(ctx, d, m)
}
func resourceSystemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	systemOptions, err := readSystem(m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillSystem(d, systemOptions)

	return nil
}
func resourceSystemUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSystem(m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setSystem(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := sess.commitConf("update resource junos_system", jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	d.Partial(false)

	return resourceSystemRead(ctx, d, m)
}
func resourceSystemDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
func resourceSystemImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	systemOptions, err := readSystem(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSystem(d, systemOptions)
	d.SetId("system")
	result[0] = d

	return result, nil
}

func setSystem(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set system "
	configSet := make([]string, 0)

	for _, nameServer := range d.Get("name_server").([]interface{}) {
		configSet = append(configSet, setPrefix+"name-server "+nameServer.(string))
	}
	if err := setSystemServices(d, m, jnprSess); err != nil {
		return err
	}
	for _, syslog := range d.Get("syslog").([]interface{}) {
		if syslog != nil {
			syslogM := syslog.(map[string]interface{})
			for _, archive := range syslogM["archive"].([]interface{}) {
				configSet = append(configSet, setPrefix+"syslog archive")
				if archive != nil {
					archiveM := archive.(map[string]interface{})
					if archiveM["binary_data"].(bool) && archiveM["no_binary_data"].(bool) {
						return fmt.Errorf("conflict between 'binary_data' and 'no_binary_data' for syslog archive")
					}
					if archiveM["binary_data"].(bool) {
						configSet = append(configSet, setPrefix+"syslog archive binary-data")
					}
					if archiveM["no_binary_data"].(bool) {
						configSet = append(configSet, setPrefix+"syslog archive no-binary-data")
					}
					if archiveM["files"].(int) > 0 {
						configSet = append(configSet, setPrefix+"syslog archive files "+strconv.Itoa(archiveM["files"].(int)))
					}
					if archiveM["size"].(int) > 0 {
						configSet = append(configSet, setPrefix+"syslog archive size "+strconv.Itoa(archiveM["size"].(int)))
					}
					if archiveM["no_world_readable"].(bool) && archiveM["world_readable"].(bool) {
						return fmt.Errorf("conflict between 'world_readable' and 'no_world_readable' for syslog archive")
					}
					if archiveM["no_world_readable"].(bool) {
						configSet = append(configSet, setPrefix+"syslog archive no-world-readable")
					}
					if archiveM["world_readable"].(bool) {
						configSet = append(configSet, setPrefix+"syslog archive world-readable")
					}
				}
			}
			if syslogM["log_rotate_frequency"].(int) > 0 {
				configSet = append(configSet, setPrefix+"syslog log-rotate-frequency "+
					strconv.Itoa(syslogM["log_rotate_frequency"].(int)))
			}
			if syslogM["source_address"].(string) != "" {
				configSet = append(configSet, setPrefix+"syslog source-address "+syslogM["source_address"].(string))
			}
		}
	}
	if d.Get("tracing_dest_override_syslog_host").(string) != "" {
		configSet = append(configSet, setPrefix+"tracing destination-override syslog host "+
			d.Get("tracing_dest_override_syslog_host").(string))
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func setSystemServices(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	setPrefix := "set system services "
	configSet := make([]string, 0)

	for _, services := range d.Get("services").([]interface{}) {
		if services != nil {
			servicesM := services.(map[string]interface{})
			for _, servicesSSH := range servicesM["ssh"].([]interface{}) {
				if servicesSSH != nil {
					servicesSSHM := servicesSSH.(map[string]interface{})
					for _, auth := range servicesSSHM["authentication_order"].([]interface{}) {
						configSet = append(configSet, setPrefix+"ssh authentication-order "+auth.(string))
					}
					for _, ciphers := range servicesSSHM["ciphers"].([]interface{}) {
						configSet = append(configSet, setPrefix+"ssh ciphers "+ciphers.(string))
					}
					if servicesSSHM["client_alive_count_max"].(int) > -1 {
						configSet = append(configSet, setPrefix+"ssh client-alive-count-max "+
							strconv.Itoa(servicesSSHM["client_alive_count_max"].(int)))
					}
					if servicesSSHM["client_alive_interval"].(int) > -1 {
						configSet = append(configSet, setPrefix+"ssh client-alive-interval "+
							strconv.Itoa(servicesSSHM["client_alive_interval"].(int)))
					}
					if servicesSSHM["connection_limit"].(int) > 0 {
						configSet = append(configSet, setPrefix+"ssh connection-limit "+
							strconv.Itoa(servicesSSHM["connection_limit"].(int)))
					}
					if servicesSSHM["fingerprint_hash"].(string) != "" {
						configSet = append(configSet, setPrefix+"ssh fingerprint-hash "+
							servicesSSHM["fingerprint_hash"].(string))
					}
					for _, hostkeyAlgo := range servicesSSHM["hostkey_algorithm"].([]interface{}) {
						configSet = append(configSet, setPrefix+"ssh hostkey-algorithm "+hostkeyAlgo.(string))
					}
					for _, keyExchange := range servicesSSHM["key_exchange"].([]interface{}) {
						configSet = append(configSet, setPrefix+"ssh key-exchange "+keyExchange.(string))
					}
					if servicesSSHM["log_key_changes"].(bool) {
						configSet = append(configSet, setPrefix+"ssh log-key-changes")
					}
					for _, macs := range servicesSSHM["macs"].([]interface{}) {
						configSet = append(configSet, setPrefix+"ssh macs "+macs.(string))
					}
					if servicesSSHM["max_pre_authentication_packets"].(int) > 0 {
						configSet = append(configSet, setPrefix+"ssh max-pre-authentication-packets "+
							strconv.Itoa(servicesSSHM["max_pre_authentication_packets"].(int)))
					}
					if servicesSSHM["max_sessions_per_connection"].(int) > 0 {
						configSet = append(configSet, setPrefix+"ssh max-sessions-per-connection "+
							strconv.Itoa(servicesSSHM["max_sessions_per_connection"].(int)))
					}
					if servicesSSHM["no_passwords"].(bool) {
						configSet = append(configSet, setPrefix+"ssh no-passwords")
					}
					if servicesSSHM["no_public_keys"].(bool) {
						configSet = append(configSet, setPrefix+"ssh no-public-keys")
					}
					if servicesSSHM["port"].(int) > 0 {
						configSet = append(configSet, setPrefix+"ssh port "+
							strconv.Itoa(servicesSSHM["port"].(int)))
					}
					for _, version := range servicesSSHM["protocol_version"].([]interface{}) {
						configSet = append(configSet, setPrefix+"ssh protocol-version "+version.(string))
					}
					if servicesSSHM["rate_limit"].(int) > 0 {
						configSet = append(configSet, setPrefix+"ssh rate-limit "+
							strconv.Itoa(servicesSSHM["rate_limit"].(int)))
					}
					if servicesSSHM["root_login"].(string) != "" {
						configSet = append(configSet, setPrefix+"ssh root-login "+servicesSSHM["root_login"].(string))
					}
					if servicesSSHM["no_tcp_forwarding"].(bool) && servicesSSHM["tcp_forwarding"].(bool) {
						return fmt.Errorf("conflict between 'no_tcp_forwarding' and 'tcp_forwarding' for services ssh")
					}
					if servicesSSHM["no_tcp_forwarding"].(bool) {
						configSet = append(configSet, setPrefix+"ssh no-tcp-forwarding")
					}
					if servicesSSHM["tcp_forwarding"].(bool) {
						configSet = append(configSet, setPrefix+"ssh tcp-forwarding")
					}
				}
			}
		}
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func listLinesServices() []string {
	ls := make([]string, 0)
	ls = append(ls, listLinesServicesSSH()...)

	return ls
}
func listLinesServicesSSH() []string {
	return []string{
		"services ssh authentication-order",
		"services ssh ciphers",
		"services ssh client-alive-count-max",
		"services ssh client-alive-interval",
		"services ssh connection-limit",
		"services ssh fingerprint-hash",
		"services ssh hostkey-algorithm",
		"services ssh key-exchange",
		"services ssh log-key-changes",
		"services ssh macs",
		"services ssh max-pre-authentication-packets",
		"services ssh max-sessions-per-connection",
		"services ssh no-passwords",
		"services ssh no-public-keys",
		"services ssh port",
		"services ssh protocol-version",
		"services ssh rate-limit",
		"services ssh root-login",
		"services ssh no-tcp-forwarding",
		"services ssh tcp-forwarding",
	}
}
func listLinesSyslog() []string {
	return []string{
		"syslog archive",
		"syslog log-rotate-frequency",
		"syslog source-address",
	}
}
func delSystem(m interface{}, jnprSess *NetconfObject) error {
	listLinesToDelete := make([]string, 0)
	listLinesToDelete = append(listLinesToDelete, "name-server")
	listLinesToDelete = append(listLinesToDelete, listLinesServices()...)
	listLinesToDelete = append(listLinesToDelete, listLinesSyslog()...)
	listLinesToDelete = append(listLinesToDelete,
		"tracing destination-override syslog host",
	)
	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := "delete system "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readSystem(m interface{}, jnprSess *NetconfObject) (systemOptions, error) {
	sess := m.(*Session)
	var confRead systemOptions

	systemConfig, err := sess.command("show configuration system"+
		" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if systemConfig != emptyWord {
		for _, item := range strings.Split(systemConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "name-server "):
				confRead.nameServer = append(confRead.nameServer, strings.TrimPrefix(itemTrim, "name-server "))
			case checkStringHasPrefixInList(itemTrim, listLinesServices()):
				if len(confRead.services) == 0 {
					confRead.services = append(confRead.services, map[string]interface{}{
						"ssh": make([]map[string]interface{}, 0),
					})
				}
				if checkStringHasPrefixInList(itemTrim, listLinesServicesSSH()) {
					if len(confRead.services[0]["ssh"].([]map[string]interface{})) == 0 {
						confRead.services[0]["ssh"] = append(confRead.services[0]["ssh"].([]map[string]interface{}),
							map[string]interface{}{
								"authentication_order":           make([]string, 0),
								"ciphers":                        make([]string, 0),
								"client_alive_count_max":         -1,
								"client_alive_interval":          -1,
								"connection_limit":               0,
								"fingerprint_hash":               "",
								"hostkey_algorithm":              make([]string, 0),
								"key_exchange":                   make([]string, 0),
								"log_key_changes":                false,
								"macs":                           make([]string, 0),
								"max_pre_authentication_packets": 0,
								"max_sessions_per_connection":    0,
								"no_passwords":                   false,
								"no_public_keys":                 false,
								"port":                           0,
								"protocol_version":               make([]string, 0),
								"rate_limit":                     0,
								"root_login":                     "",
								"no_tcp_forwarding":              false,
								"tcp_forwarding":                 false,
							})
					}
					switch {
					case strings.HasPrefix(itemTrim, "services ssh authentication-order "):
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["authentication_order"] = append(
							confRead.services[0]["ssh"].([]map[string]interface{})[0]["authentication_order"].([]string),
							strings.TrimPrefix(itemTrim, "services ssh authentication-order "))
					case strings.HasPrefix(itemTrim, "services ssh ciphers "):
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["ciphers"] = append(
							confRead.services[0]["ssh"].([]map[string]interface{})[0]["ciphers"].([]string),
							strings.TrimPrefix(itemTrim, "services ssh ciphers "))
					case strings.HasPrefix(itemTrim, "services ssh client-alive-count-max "):
						var err error
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["client_alive_count_max"], err =
							strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh client-alive-count-max "))
						if err != nil {
							return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
						}
					case strings.HasPrefix(itemTrim, "services ssh client-alive-interval "):
						var err error
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["client_alive_interval"], err =
							strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh client-alive-interval "))
						if err != nil {
							return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
						}
					case strings.HasPrefix(itemTrim, "services ssh connection-limit "):
						var err error
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["connection_limit"], err =
							strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh connection-limit "))
						if err != nil {
							return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
						}
					case strings.HasPrefix(itemTrim, "services ssh fingerprint-hash "):
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["fingerprint_hash"] = strings.TrimPrefix(
							itemTrim, "services ssh fingerprint-hash ")
					case strings.HasPrefix(itemTrim, "services ssh hostkey-algorithm "):
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["hostkey_algorithm"] = append(
							confRead.services[0]["ssh"].([]map[string]interface{})[0]["hostkey_algorithm"].([]string),
							strings.TrimPrefix(itemTrim, "services ssh hostkey-algorithm "))
					case strings.HasPrefix(itemTrim, "services ssh key-exchange "):
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["key_exchange"] = append(
							confRead.services[0]["ssh"].([]map[string]interface{})[0]["key_exchange"].([]string),
							strings.TrimPrefix(itemTrim, "services ssh key-exchange "))
					case itemTrim == "services ssh log-key-changes":
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["log_key_changes"] = true
					case strings.HasPrefix(itemTrim, "services ssh macs "):
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["macs"] = append(
							confRead.services[0]["ssh"].([]map[string]interface{})[0]["macs"].([]string),
							strings.TrimPrefix(itemTrim, "services ssh macs "))
					case strings.HasPrefix(itemTrim, "services ssh max-pre-authentication-packets "):
						var err error
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["max_pre_authentication_packets"], err =
							strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh max-pre-authentication-packets "))
						if err != nil {
							return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
						}
					case strings.HasPrefix(itemTrim, "services ssh max-sessions-per-connection "):
						var err error
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["max_sessions_per_connection"], err =
							strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh max-sessions-per-connection "))
						if err != nil {
							return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
						}
					case itemTrim == "services ssh no-passwords":
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["no_passwords"] = true
					case itemTrim == "services ssh no-public-keys":
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["no_public_keys"] = true
					case strings.HasPrefix(itemTrim, "services ssh port "):
						var err error
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["port"], err =
							strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh port "))
						if err != nil {
							return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
						}
					case strings.HasPrefix(itemTrim, "services ssh protocol-version "):
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["protocol_version"] = append(
							confRead.services[0]["ssh"].([]map[string]interface{})[0]["protocol_version"].([]string),
							strings.TrimPrefix(itemTrim, "services ssh protocol-version "))
					case strings.HasPrefix(itemTrim, "services ssh rate-limit "):
						var err error
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["rate_limit"], err =
							strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh rate-limit "))
						if err != nil {
							return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
						}
					case strings.HasPrefix(itemTrim, "services ssh root-login "):
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["root_login"] =
							strings.TrimPrefix(itemTrim, "services ssh root-login ")
					case itemTrim == "services ssh no-tcp-forwarding":
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["no_tcp_forwarding"] = true
					case itemTrim == "services ssh tcp-forwarding":
						confRead.services[0]["ssh"].([]map[string]interface{})[0]["tcp_forwarding"] = true
					}
				}
			case checkStringHasPrefixInList(itemTrim, listLinesSyslog()):
				if len(confRead.syslog) == 0 {
					confRead.syslog = append(confRead.syslog, map[string]interface{}{
						"archive":              make([]map[string]interface{}, 0),
						"log_rotate_frequency": 0,
						"source_address":       "",
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "syslog archive"):
					if len(confRead.syslog[0]["archive"].([]map[string]interface{})) == 0 {
						confRead.syslog[0]["archive"] = append(confRead.syslog[0]["archive"].([]map[string]interface{}),
							map[string]interface{}{
								"binary_data":       false,
								"no_binary_data":    false,
								"files":             0,
								"size":              0,
								"no_world_readable": false,
								"world_readable":    false,
							})
					}
					switch {
					case itemTrim == "syslog archive binary-data":
						confRead.syslog[0]["archive"].([]map[string]interface{})[0]["binary_data"] = true
					case itemTrim == "syslog archive no-binary-data":
						confRead.syslog[0]["archive"].([]map[string]interface{})[0]["no_binary_data"] = true
					case strings.HasPrefix(itemTrim, "syslog archive files "):
						var err error
						confRead.syslog[0]["archive"].([]map[string]interface{})[0]["files"], err = strconv.Atoi(
							strings.TrimPrefix(itemTrim, "syslog archive files "))
						if err != nil {
							return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
						}
					case strings.HasPrefix(itemTrim, "syslog archive size "):
						var err error
						confRead.syslog[0]["archive"].([]map[string]interface{})[0]["size"], err = strconv.Atoi(
							strings.TrimPrefix(itemTrim, "syslog archive size "))
						if err != nil {
							return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
						}
					case itemTrim == "syslog archive no-world-readable":
						confRead.syslog[0]["archive"].([]map[string]interface{})[0]["no_world_readable"] = true
					case itemTrim == "syslog archive world-readable":
						confRead.syslog[0]["archive"].([]map[string]interface{})[0]["world_readable"] = true
					}
				case strings.HasPrefix(itemTrim, "syslog log-rotate-frequency "):
					var err error
					confRead.syslog[0]["log_rotate_frequency"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrim, "syslog log-rotate-frequency "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "syslog source-address "):
					confRead.syslog[0]["source_address"] = strings.TrimPrefix(
						itemTrim, "syslog source-address ")
				}
			case strings.HasPrefix(itemTrim, "tracing destination-override syslog host "):
				confRead.tracingDestinationOverrideSyslogHost = strings.TrimPrefix(itemTrim,
					"tracing destination-override syslog host ")
			}
		}
	}

	return confRead, nil
}

func fillSystem(d *schema.ResourceData, systemOptions systemOptions) {
	if tfErr := d.Set("name_server", systemOptions.nameServer); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("services", systemOptions.services); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("syslog", systemOptions.syslog); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("tracing_dest_override_syslog_host",
		systemOptions.tracingDestinationOverrideSyslogHost); tfErr != nil {
		panic(tfErr)
	}
}
