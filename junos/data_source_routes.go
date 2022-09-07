package junos

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type routesTableOpts struct {
	table []map[string]interface{}
}

func dataSourceRoutes() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceRoutesRead,
		Schema: map[string]*schema.Schema{
			"table_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"entry": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"as_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"current_active": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"local_preference": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"metric": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"next_hop": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"local_interface": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"selected_next_hop": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"to": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"via": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"next_hop_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"preference": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"protocol": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	mutex.Lock()
	routesTable, err := searchRoutes(d, clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if tfErr := d.Set("table", routesTable.table); tfErr != nil {
		panic(tfErr)
	}
	if v := d.Get("table_name").(string); v != "" {
		d.SetId(v)
	} else {
		d.SetId("all")
	}

	return nil
}

func searchRoutes(d *schema.ResourceData, clt *Client, junSess *junosSession,
) (routesTableOpts, error) {
	var result routesTableOpts
	rpcReq := rpcGetRouteAllInformation
	if v := d.Get("table_name").(string); v != "" {
		rpcReq = fmt.Sprintf(rpcGetRouteAllTableInformation, v)
	}
	replyData, err := clt.commandXML(rpcReq, junSess)
	if err != nil {
		return result, err
	}
	var routeTable getRouteInformationReply
	err = xml.Unmarshal([]byte(replyData), &routeTable.RouteInfo)
	if err != nil {
		return result, fmt.Errorf("failed to xml unmarshal reply data '%s': %w", replyData, err)
	}
	for _, tableInfo := range routeTable.RouteInfo.RouteTable {
		table := map[string]interface{}{
			"name":  tableInfo.TableName,
			"route": make([]map[string]interface{}, 0, len(tableInfo.Route)),
		}
		for _, routeInfo := range tableInfo.Route {
			route := map[string]interface{}{
				"destination": routeInfo.Destination,
				"entry":       make([]map[string]interface{}, 0, len(routeInfo.Entry)),
			}
			for _, entryInfo := range routeInfo.Entry {
				entry := map[string]interface{}{
					"as_path":          strings.Trim(entryInfo.ASPath, "\n"),
					"current_active":   entryInfo.CurrentActive != nil,
					"local_preference": entryInfo.LocalPreference,
					"metric":           entryInfo.Metric,
					"next_hop":         make([]map[string]interface{}, 0, len(entryInfo.NextHop)),
					"next_hop_type":    entryInfo.NextHopType,
					"preference":       entryInfo.Preference,
					"protocol":         entryInfo.Protocol,
				}
				for _, nextHopInfo := range entryInfo.NextHop {
					entry["next_hop"] = append(
						entry["next_hop"].([]map[string]interface{}),
						map[string]interface{}{
							"local_interface":   nextHopInfo.LocalInterface,
							"selected_next_hop": nextHopInfo.SelectedNextHop != nil,
							"to":                nextHopInfo.To,
							"via":               nextHopInfo.Via,
						})
				}
				route["entry"] = append(route["entry"].([]map[string]interface{}), entry)
			}
			table["route"] = append(table["route"].([]map[string]interface{}), route)
		}
		result.table = append(result.table, table)
	}

	return result, nil
}
