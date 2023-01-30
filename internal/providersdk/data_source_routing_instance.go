package providersdk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

func dataSourceRoutingInstance() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceRoutingInstanceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64, formatDefault),
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"as": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_export": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_import": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"interface": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"route_distinguisher": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_export": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vrf_import": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vrf_target": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_target_auto": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"vrf_target_export": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_target_import": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vtep_source_interface": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRoutingInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	mutex.Lock()
	instanceOptions, err := readRoutingInstance(d.Get("name").(string), junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if instanceOptions.name == "" {
		return diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("name").(string)))
	}
	d.SetId(instanceOptions.name)
	fillRoutingInstanceDataSource(d, instanceOptions)

	return nil
}

func fillRoutingInstanceDataSource(d *schema.ResourceData, instanceOptions instanceOptions) {
	if tfErr := d.Set("name", instanceOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as", instanceOptions.as); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", instanceOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("instance_export", instanceOptions.instanceExport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("instance_import", instanceOptions.instanceImport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interface", instanceOptions.interFace); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("route_distinguisher", instanceOptions.routeDistinguisher); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("router_id", instanceOptions.routerID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("type", instanceOptions.instanceType); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vrf_export", instanceOptions.vrfExport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vrf_import", instanceOptions.vrfImport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vrf_target", instanceOptions.vrfTarget); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vrf_target_auto", instanceOptions.vrfTargetAuto); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vrf_target_export", instanceOptions.vrfTargetExport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vrf_target_import", instanceOptions.vrfTargetImport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vtep_source_interface", instanceOptions.vtepSourceIf); tfErr != nil {
		panic(tfErr)
	}
}
