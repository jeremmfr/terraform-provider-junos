package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type fwdOptsDhcpRelSrvGrpOptions struct {
	name            string
	routingInstance string
	version         string
	ipAddress       []string
}

func resourceForwardingOptionsDhcpRelayServerGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceForwardingOptionsDhcpRelayServerGroupCreate,
		ReadWithoutTimeout:   resourceForwardingOptionsDhcpRelayServerGroupRead,
		UpdateWithoutTimeout: resourceForwardingOptionsDhcpRelayServerGroupUpdate,
		DeleteWithoutTimeout: resourceForwardingOptionsDhcpRelayServerGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceForwardingOptionsDhcpRelayServerGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          defaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "v4",
				ValidateFunc: validation.StringInSlice([]string{"v4", "v6"}, false),
			},
			"ip_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPAddress,
				},
			},
		},
	}
}

func resourceForwardingOptionsDhcpRelayServerGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	nameArg := d.Get("name").(string)
	routingInstanceArg := d.Get("routing_instance").(string)
	versionArg := d.Get("version").(string)
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setForwardingOptionsDhcpRelayServerGroup(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(nameArg + idSeparator + routingInstanceArg + idSeparator + versionArg)

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if routingInstanceArg != defaultW {
		instanceExists, err := checkRoutingInstanceExists(routingInstanceArg, clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", routingInstanceArg))...)
		}
	}
	fwdOptsDhcpRelSrvGrpExists, err := checkForwardingOptionsDhcpRelayServerGroupExists(
		nameArg,
		routingInstanceArg,
		versionArg,
		clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if fwdOptsDhcpRelSrvGrpExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))
		if versionArg == "v6" {
			return append(diagWarns,
				diag.FromErr(fmt.Errorf("forwarding-options dhcp-relay dhcpv6 server-group %v"+
					" already exists in routing-instance %s", nameArg, routingInstanceArg))...)
		}

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("forwarding-options dhcp-relay server-group %v"+
				" already exists in routing-instance %s", nameArg, routingInstanceArg))...)
	}

	if err := setForwardingOptionsDhcpRelayServerGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	warns, err := clt.commitConf("create resource junos_forwardingoptions_dhcprelay_servergroup", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	fwdOptsDhcpRelSrvGrpExists, err = checkForwardingOptionsDhcpRelayServerGroupExists(
		nameArg,
		routingInstanceArg,
		versionArg,
		clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if fwdOptsDhcpRelSrvGrpExists {
		d.SetId(nameArg + idSeparator + routingInstanceArg + idSeparator + versionArg)
	} else {
		if versionArg == "v6" {
			return append(diagWarns, diag.FromErr(fmt.Errorf(
				"forwarding-options dhcp-relay dhcpv6 server-group %v not exists in routing_instance %s after commit "+
					"=> check your config", nameArg, routingInstanceArg))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"forwarding-options dhcp-relay server-group %v not exists in routing_instance %s after commit "+
				"=> check your config", nameArg, routingInstanceArg))...)
	}

	return append(diagWarns, resourceForwardingOptionsDhcpRelayServerGroupReadWJunSess(d, clt, junSess)...)
}

func resourceForwardingOptionsDhcpRelayServerGroupRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceForwardingOptionsDhcpRelayServerGroupReadWJunSess(d, clt, junSess)
}

func resourceForwardingOptionsDhcpRelayServerGroupReadWJunSess(
	d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	fwdOptsDhcpRelSrvGrpOptions, err := readForwardingOptionsDhcpRelayServerGroup(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		d.Get("version").(string),
		clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if fwdOptsDhcpRelSrvGrpOptions.name == "" {
		d.SetId("")
	} else {
		fillForwardingOptionsDhcpRelayServerGroupData(d, fwdOptsDhcpRelSrvGrpOptions)
	}

	return nil
}

func resourceForwardingOptionsDhcpRelayServerGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	nameArg := d.Get("name").(string)
	routingInstanceArg := d.Get("routing_instance").(string)
	versionArg := d.Get("version").(string)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delForwardingOptionsDhcpRelayServerGroup(nameArg, routingInstanceArg, versionArg, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setForwardingOptionsDhcpRelayServerGroup(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delForwardingOptionsDhcpRelayServerGroup(nameArg, routingInstanceArg, versionArg, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setForwardingOptionsDhcpRelayServerGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	warns, err := clt.commitConf("update resource junos_forwardingoptions_dhcprelay_servergroup", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceForwardingOptionsDhcpRelayServerGroupReadWJunSess(d, clt, junSess)...)
}

func resourceForwardingOptionsDhcpRelayServerGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	nameArg := d.Get("name").(string)
	routingInstanceArg := d.Get("routing_instance").(string)
	versionArg := d.Get("version").(string)
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delForwardingOptionsDhcpRelayServerGroup(nameArg, routingInstanceArg, versionArg, clt, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delForwardingOptionsDhcpRelayServerGroup(nameArg, routingInstanceArg, versionArg, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_forwardingoptions_dhcprelay_servergroup", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceForwardingOptionsDhcpRelayServerGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 3 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	if idSplit[2] != "v4" && idSplit[2] != "v6" {
		return nil, fmt.Errorf("bad version '%s' in id, need to be 'v4' or 'v6' (id must be "+
			"<name>"+idSeparator+"<routing_instance>"+idSeparator+"<version>)", idSplit[2])
	}
	fwdOptsDhcpRelSrvGrpExists, err := checkForwardingOptionsDhcpRelayServerGroupExists(
		idSplit[0],
		idSplit[1],
		idSplit[2],
		clt, junSess)
	if err != nil {
		return nil, err
	}
	if !fwdOptsDhcpRelSrvGrpExists {
		if idSplit[2] == "v6" {
			return nil, fmt.Errorf("don't find forwarding-options dhcp-relay dhcpv6 server-group with id '%v' (id must be "+
				"<name>"+idSeparator+"<routing_instance>"+idSeparator+"<version>)", d.Id())
		}

		return nil, fmt.Errorf("don't find forwarding-options dhcp-relay server-group with id '%v' (id must be "+
			"<name>"+idSeparator+"<routing_instance>"+idSeparator+"<version>)", d.Id())
	}
	fwdOptsDhcpRelSrvGrpOptions, err := readForwardingOptionsDhcpRelayServerGroup(
		idSplit[0],
		idSplit[1],
		idSplit[2],
		clt, junSess)
	if err != nil {
		return nil, err
	}
	fillForwardingOptionsDhcpRelayServerGroupData(d, fwdOptsDhcpRelSrvGrpOptions)

	result[0] = d

	return result, nil
}

func checkForwardingOptionsDhcpRelayServerGroupExists(
	name, instance, version string, clt *Client, junSess *junosSession,
) (bool, error) {
	showCmd := cmdShowConfig
	if instance != defaultW {
		showCmd += routingInstancesWS + instance + " "
	}
	showCmd += "forwarding-options dhcp-relay "
	if version == "v6" {
		showCmd += "dhcpv6 server-group " + name
	} else {
		showCmd += "server-group " + name
	}
	showConfig, err := clt.command(showCmd+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setForwardingOptionsDhcpRelayServerGroup(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := setLS
	if d.Get("routing_instance").(string) != defaultW {
		setPrefix = setRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	if d.Get("version").(string) == "v6" {
		setPrefix += "forwarding-options dhcp-relay dhcpv6 server-group " + d.Get("name").(string) + " "
	} else {
		setPrefix += "forwarding-options dhcp-relay server-group " + d.Get("name").(string) + " "
	}

	configSet = append(configSet, setPrefix)
	for _, v := range d.Get("ip_address").([]interface{}) {
		configSet = append(configSet, setPrefix+v.(string))
	}

	return clt.configSet(configSet, junSess)
}

func readForwardingOptionsDhcpRelayServerGroup(name, instance, version string, clt *Client, junSess *junosSession,
) (fwdOptsDhcpRelSrvGrpOptions, error) {
	var confRead fwdOptsDhcpRelSrvGrpOptions

	showCmd := cmdShowConfig
	if instance != defaultW {
		showCmd += routingInstancesWS + instance + " "
	}
	showCmd += "forwarding-options dhcp-relay "
	if version == "v6" {
		showCmd += "dhcpv6 server-group " + name
	} else {
		showCmd += "server-group " + name
	}
	showConfig, err := clt.command(showCmd+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		confRead.routingInstance = instance
		confRead.version = version
		for _, item := range strings.Split(showConfig, "\n") {
			itemTrim := strings.TrimPrefix(item, setLS)
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			if itemTrim != "" {
				confRead.ipAddress = append(confRead.ipAddress, strings.TrimPrefix(itemTrim, " "))
			}
		}
	}

	return confRead, nil
}

func delForwardingOptionsDhcpRelayServerGroup(name, instance, version string, clt *Client, junSess *junosSession,
) error {
	configSet := make([]string, 0, 1)
	switch {
	case instance == defaultW && version == "v6":
		configSet = append(configSet, deleteLS+"forwarding-options dhcp-relay dhcpv6 server-group "+name)
	case instance == defaultW && version == "v4":
		configSet = append(configSet, deleteLS+"forwarding-options dhcp-relay server-group "+name)
	case instance != defaultW && version == "v6":
		configSet = append(configSet, delRoutingInstances+instance+" "+
			"forwarding-options dhcp-relay dhcpv6 server-group "+name)
	case instance != defaultW && version == "v4":
		configSet = append(configSet, delRoutingInstances+instance+" "+
			"forwarding-options dhcp-relay server-group "+name)
	}

	return clt.configSet(configSet, junSess)
}

func fillForwardingOptionsDhcpRelayServerGroupData(
	d *schema.ResourceData, fwdOptsDhcpRelSrvGrpOptions fwdOptsDhcpRelSrvGrpOptions,
) {
	if tfErr := d.Set("name", fwdOptsDhcpRelSrvGrpOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", fwdOptsDhcpRelSrvGrpOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("version", fwdOptsDhcpRelSrvGrpOptions.version); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ip_address", fwdOptsDhcpRelSrvGrpOptions.ipAddress); tfErr != nil {
		panic(tfErr)
	}
}
