package providersdk

import (
	"context"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSystemInformation() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceSystemInformationRead,

		Schema: map[string]*schema.Schema{
			"hardware_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_node": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSystemInformationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	// Catches case where hostname is not set
	if junSess.SystemInformation.HostName != "" {
		d.SetId(junSess.SystemInformation.HostName)
	} else {
		d.SetId("Null-Hostname")
	}

	if tfErr := d.Set("hardware_model", junSess.SystemInformation.HardwareModel); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("os_name", junSess.SystemInformation.OsName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("os_version", junSess.SystemInformation.OsVersion); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("serial_number", junSess.SystemInformation.SerialNumber); tfErr != nil {
		panic(tfErr)
	}
	// Fix recommended in https://stackoverflow.com/a/23725010 due ot the lack of being able to Unmarshal self-closing
	// xml tags. This issue is tracked here - https://github.com/golang/go/issues/21399
	// Pointer to bool in sysInfo struct will be nil if the tag does not exist
	if junSess.SystemInformation.ClusterNode != nil {
		if tfErr := d.Set("cluster_node", true); tfErr != nil {
			panic(tfErr)
		}
	} else {
		if tfErr := d.Set("cluster_node", false); tfErr != nil {
			panic(tfErr)
		}
	}

	return nil
}
