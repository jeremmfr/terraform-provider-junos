package junos

import (
	"context"

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
	sess := m.(*Session)
	j, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(j)

	// Catches case where hostname is not set
	if j.SystemInformation.HostName != "" {
		d.SetId(j.SystemInformation.HostName)
	} else {
		d.SetId("Null-Hostname")
	}

	if tfErr := d.Set("hardware_model", j.SystemInformation.HardwareModel); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("os_name", j.SystemInformation.OsName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("os_version", j.SystemInformation.OsVersion); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("serial_number", j.SystemInformation.SerialNumber); tfErr != nil {
		panic(tfErr)
	}
	// Fix recommended in https://stackoverflow.com/a/23725010 due ot the lack of being able to Unmarshal self-closing
	// xml tags. This issue is tracked here - https://github.com/golang/go/issues/21399
	// Pointer to bool in sysInfo struct will be nil if the tag does not exist
	if j.SystemInformation.ClusterNode != nil {
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
