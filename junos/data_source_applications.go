package junos

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApplications() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceApplicationsRead,
		Schema: map[string]*schema.Schema{
			"match_name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if _, err := regexp.Compile(value); err != nil {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not valid regexp", value, k))
					}

					return
				},
			},
			"match_options": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alg": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"application_protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ether_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp6_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp6_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"inactivity_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"inactivity_timeout_never": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rpc_program_number": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source_port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ether_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inactivity_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"inactivity_timeout_never": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rpc_program_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"term": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"alg": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"icmp_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"icmp_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"icmp6_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"icmp6_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"inactivity_timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"inactivity_timeout_never": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rpc_program_number": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApplicationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	mutex.Lock()
	applications, err := dataSourceApplicationsSearch(clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if err := dataSourceApplicationsFilter(d, applications); err != nil {
		return diag.FromErr(err)
	}
	fillDataApplicationsData(d, applications)
	d.SetId("match_name=" + d.Get("match_name").(string) +
		idSeparator +
		"match_options_n=" + fmt.Sprintf("%d", len(d.Get("match_options").(*schema.Set).List())),
	)

	return nil
}

func dataSourceApplicationsSearch(clt *Client, junSess *junosSession) (map[string]applicationOptions, error) {
	results := make(map[string]applicationOptions, 0)
	for _, config := range []string{
		"groups junos-defaults applications",
		"applications",
	} {
		showConfig, err := clt.command(cmdShowConfig+config+pipeDisplaySetRelative, junSess)
		if err != nil {
			return results, err
		}
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			if item == "" {
				continue
			}
			if !strings.HasPrefix(item, "set application ") {
				continue
			}
			itemTrim := strings.TrimPrefix(item, "set application ")
			itemTrimFields := strings.Split(itemTrim, " ")
			if _, ok := results[itemTrimFields[0]]; !ok {
				results[itemTrimFields[0]] = applicationOptions{name: itemTrimFields[0]}
			}
			appOpts := results[itemTrimFields[0]]
			if err := appOpts.readLine(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" ")); err != nil {
				return results, err
			}
			results[itemTrimFields[0]] = appOpts
		}
	}

	return results, nil
}

func dataSourceApplicationsFilter( //nolint: gocognit
	d *schema.ResourceData, results map[string]applicationOptions,
) error {
	if mName := d.Get("match_name").(string); mName != "" {
		for appKey, app := range results {
			matched, err := regexp.MatchString(mName, app.name)
			if err != nil {
				return fmt.Errorf("failed to regexp with '%s': %w", mName, err)
			}
			if !matched {
				delete(results, appKey)
			}
		}
	}
	if matchOpts := d.Get("match_options").(*schema.Set).List(); len(matchOpts) > 0 {
		// for each app, check if all options is matched
		for appKey, app := range results {
			matchOk := true
			// application defined with term or not but not with both
			if len(app.term) > 0 {
				listStringOpts := []string{
					"alg",
					"destination_port",
					"icmp_code",
					"icmp_type",
					"icmp6_code",
					"icmp6_type",
					"protocol",
					"rpc_program_number",
					"source_port",
					"uuid",
				}
				listIntOpts := []string{
					"inactivity_timeout",
				}
				listBoolOpts := []string{
					"inactivity_timeout_never",
				}
				// check if a term match a options block
				matchOptsNum := 0
			each_opts:
				for _, optsT := range matchOpts {
					optsT := optsT.(map[string]interface{})
				each_term:
					for _, term := range app.term {
						for _, optStr := range listStringOpts {
							if v := optsT[optStr].(string); v != "" {
								if term[optStr].(string) != v {
									continue each_term
								}
							}
						}
						for _, optInt := range listIntOpts {
							if v := optsT[optInt].(int); v != 0 {
								if term[optInt].(int) != v {
									continue each_term
								}
							}
						}
						for _, optBool := range listBoolOpts {
							if optsT[optBool].(bool) {
								if !term[optBool].(bool) {
									continue each_term
								}
							}
						}
						// current term match current options block
						matchOptsNum++

						continue each_opts
					}
				}
				// all options block has been validated
				if matchOptsNum != len(matchOpts) {
					matchOk = false
				}
			} else {
				for _, optsT := range matchOpts {
					optsT := optsT.(map[string]interface{})
					if v := optsT["application_protocol"].(string); v != "" {
						if app.applicationProtocol != v {
							matchOk = false

							break
						}
					}
					if v := optsT["destination_port"].(string); v != "" {
						if app.destinationPort != v {
							matchOk = false

							break
						}
					}
					if v := optsT["ether_type"].(string); v != "" {
						if app.etherType != v {
							matchOk = false

							break
						}
					}
					if v := optsT["inactivity_timeout"].(int); v != 0 {
						if app.inactivityTimeout != v {
							matchOk = false

							break
						}
					}
					if optsT["inactivity_timeout_never"].(bool) {
						if !app.inactivityTimeoutNever {
							matchOk = false

							break
						}
					}
					if v := optsT["protocol"].(string); v != "" {
						if app.protocol != v {
							matchOk = false

							break
						}
					}
					if v := optsT["rpc_program_number"].(string); v != "" {
						if app.rpcProgramNumber != v {
							matchOk = false

							break
						}
					}
					if v := optsT["source_port"].(string); v != "" {
						if app.sourcePort != v {
							matchOk = false

							break
						}
					}
					if v := optsT["uuid"].(string); v != "" {
						if app.uuid != v {
							matchOk = false

							break
						}
					}
				}
			}
			if !matchOk {
				delete(results, appKey)
			}
		}
	}

	return nil
}

func fillDataApplicationsData(d *schema.ResourceData, results map[string]applicationOptions) {
	resultsSets := make([]map[string]interface{}, 0, len(results))
	for _, app := range results {
		resultsSets = append(resultsSets, map[string]interface{}{
			"name":                     app.name,
			"application_protocol":     app.applicationProtocol,
			"description":              app.description,
			"destination_port":         app.destinationPort,
			"ether_type":               app.etherType,
			"inactivity_timeout":       app.inactivityTimeout,
			"inactivity_timeout_never": app.inactivityTimeoutNever,
			"protocol":                 app.protocol,
			"rpc_program_number":       app.rpcProgramNumber,
			"source_port":              app.sourcePort,
			"term":                     app.term,
			"uuid":                     app.uuid,
		})
	}
	if tfErr := d.Set("applications", resultsSets); tfErr != nil {
		panic(tfErr)
	}
}
