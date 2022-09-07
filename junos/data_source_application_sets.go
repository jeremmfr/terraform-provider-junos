package junos

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApplicationSets() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceApplicationSetsRead,
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
			"match_applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"application_sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"applications": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceApplicationSetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	mutex.Lock()
	applicationSets, err := dataSourceApplicationSetsSearch(clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if err := dataSourceApplicationSetsFilter(d, applicationSets); err != nil {
		return diag.FromErr(err)
	}
	fillDataApplicationSetsData(d, applicationSets)
	d.SetId("match_name=" + d.Get("match_name").(string) +
		idSeparator +
		"match_applications_n=" + fmt.Sprintf("%d", len(d.Get("match_applications").(*schema.Set).List())),
	)

	return nil
}

func dataSourceApplicationSetsSearch(clt *Client, junSess *junosSession) (map[string]applicationSetOptions, error) {
	results := make(map[string]applicationSetOptions, 0)
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
			if !strings.HasPrefix(item, "set application-set ") {
				continue
			}
			itemTrim := strings.TrimPrefix(item, "set application-set ")
			itemTrimSplit := strings.Split(itemTrim, " ")
			if _, ok := results[itemTrimSplit[0]]; !ok {
				results[itemTrimSplit[0]] = applicationSetOptions{name: itemTrimSplit[0]}
			}
			appSetOpts := results[itemTrimSplit[0]]
			itemTrim = strings.TrimPrefix(itemTrim, itemTrimSplit[0]+" ")
			if strings.HasPrefix(itemTrim, "application ") {
				appSetOpts.applications = append(appSetOpts.applications, strings.TrimPrefix(itemTrim, "application "))
			}
			results[itemTrimSplit[0]] = appSetOpts
		}
	}

	return results, nil
}

func dataSourceApplicationSetsFilter(
	d *schema.ResourceData, results map[string]applicationSetOptions,
) error {
	if mName := d.Get("match_name").(string); mName != "" {
		for appSetKey, appSet := range results {
			matched, err := regexp.MatchString(mName, appSet.name)
			if err != nil {
				return fmt.Errorf("failed to regexp with '%s': %w", mName, err)
			}
			if !matched {
				delete(results, appSetKey)
			}
		}
	}
	if matchApps := d.Get("match_applications").(*schema.Set).List(); len(matchApps) > 0 {
		// for each app-set, check if all applications is matched
		for appSetKey, appSet := range results {
			if len(appSet.applications) != len(matchApps) {
				delete(results, appSetKey)

				continue
			}
			matchAppsOk := make(map[string]struct{})
		each_match:
			for _, matchApp := range matchApps {
				matchAppStr := matchApp.(string)
				for _, app := range appSet.applications {
					if matchAppStr == app {
						matchAppsOk[matchAppStr] = struct{}{}

						continue each_match
					}
				}
			}
			if len(matchAppsOk) != len(matchApps) {
				delete(results, appSetKey)
			}
		}
	}

	return nil
}

func fillDataApplicationSetsData(d *schema.ResourceData, results map[string]applicationSetOptions) {
	resultsSets := make([]map[string]interface{}, 0, len(results))
	for _, appSet := range results {
		resultsSets = append(resultsSets, map[string]interface{}{
			"name":         appSet.name,
			"applications": appSet.applications,
		})
	}
	if tfErr := d.Set("application_sets", resultsSets); tfErr != nil {
		panic(tfErr)
	}
}
