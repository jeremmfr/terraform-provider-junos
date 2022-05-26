package junos

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type instanceOptions struct {
	vrfTargetAuto      bool
	as                 string
	description        string
	instanceType       string
	name               string
	routeDistinguisher string
	vrfTarget          string
	vrfTargetExport    string
	vrfTargetImport    string
	vtepSourceIf       string
	instanceExport     []string
	instanceImport     []string
	interFace          []string // to data_source
	vrfExport          []string
	vrfImport          []string
}

func resourceRoutingInstance() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceRoutingInstanceCreate,
		ReadWithoutTimeout:   resourceRoutingInstanceRead,
		UpdateWithoutTimeout: resourceRoutingInstanceUpdate,
		DeleteWithoutTimeout: resourceRoutingInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRoutingInstanceImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64, formatDefault),
			},
			"configure_rd_vrfopts_singly": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"route_distinguisher",
					"vrf_export",
					"vrf_import",
					"vrf_target",
					"vrf_target_auto",
					"vrf_target_export",
					"vrf_target_import",
				},
			},
			"configure_type_singly": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "virtual-router",
			},
			"as": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_export": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"instance_import": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"route_distinguisher": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^(\d|\.)+L?:\d+$`), "must have valid route distinguisher. Use format 'x:y'"),
			},
			"vrf_export": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"vrf_import": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"vrf_target": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^target:(\d|\.)+L?:\d+$`), "must have valid target. Use format 'target:x:y'"),
			},
			"vrf_target_auto": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vrf_target_export": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^target:(\d|\.)+L?:\d+$`), "must have valid target. Use format 'target:x:y'"),
			},
			"vrf_target_import": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^target:(\d|\.)+L?:\d+$`), "must have valid target. Use format 'target:x:y'"),
			},
			"vtep_source_interface": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if strings.Count(value, ".") != 1 {
						errors = append(errors, fmt.Errorf(
							"%q in %q need to have 1 dot", value, k))
					}

					return
				},
			},
		},
	}
}

func resourceRoutingInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setRoutingInstance(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

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
	routingInstanceExists, err := checkRoutingInstanceExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if routingInstanceExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("routing-instance %v already exists", d.Get("name").(string)))...)
	}
	if err := setRoutingInstance(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_routing_instance", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	routingInstanceExists, err = checkRoutingInstanceExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if routingInstanceExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("routing-instance %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceRoutingInstanceReadWJunSess(d, clt, junSess)...)
}

func resourceRoutingInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceRoutingInstanceReadWJunSess(d, clt, junSess)
}

func resourceRoutingInstanceReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	instanceOptions, err := readRoutingInstance(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if instanceOptions.name == "" {
		d.SetId("")
	} else {
		fillRoutingInstanceData(d, instanceOptions)
	}

	return nil
}

func resourceRoutingInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delRoutingInstanceOpts(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setRoutingInstance(d, clt, nil); err != nil {
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
	if err := delRoutingInstanceOpts(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setRoutingInstance(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_routing_instance", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRoutingInstanceReadWJunSess(d, clt, junSess)...)
}

func resourceRoutingInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delRoutingInstance(d, clt, nil); err != nil {
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
	if err := delRoutingInstance(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_routing_instance", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceRoutingInstanceImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	routingInstanceExists, err := checkRoutingInstanceExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !routingInstanceExists {
		return nil, fmt.Errorf("don't find routing instance with id '%v' (id must be <name>)", d.Id())
	}
	instanceOptions, err := readRoutingInstance(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillRoutingInstanceData(d, instanceOptions)
	result[0] = d

	return result, nil
}

func checkRoutingInstanceExists(instance string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+routingInstancesWS+instance+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setRoutingInstance(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := setRoutingInstances + d.Get("name").(string) + " "
	if d.Get("configure_type_singly").(bool) {
		if v := d.Get("type").(string); v != "" {
			return fmt.Errorf("if `configure_type_singly` = true, `type` need to be set to empty value to avoid confusion")
		}
	} else {
		if v := d.Get("type").(string); v != "" {
			configSet = append(configSet, setPrefix+"instance-type "+v)
		}
	}
	if !d.Get("configure_rd_vrfopts_singly").(bool) {
		if v := d.Get("route_distinguisher").(string); v != "" {
			configSet = append(configSet, setPrefix+"route-distinguisher "+v)
		}
		for _, v := range d.Get("vrf_export").([]interface{}) {
			configSet = append(configSet, setPrefix+"vrf-export \""+v.(string)+"\"")
		}
		for _, v := range d.Get("vrf_import").([]interface{}) {
			configSet = append(configSet, setPrefix+"vrf-import \""+v.(string)+"\"")
		}
		if v := d.Get("vrf_target").(string); v != "" {
			configSet = append(configSet, setPrefix+"vrf-target "+v)
		}
		if d.Get("vrf_target_auto").(bool) {
			configSet = append(configSet, setPrefix+"vrf-target auto")
		}
		if v := d.Get("vrf_target_export").(string); v != "" {
			configSet = append(configSet, setPrefix+"vrf-target export "+v)
		}
		if v := d.Get("vrf_target_import").(string); v != "" {
			configSet = append(configSet, setPrefix+"vrf-target import "+v)
		}
	}
	if v := d.Get("as").(string); v != "" {
		configSet = append(configSet, setPrefix+"routing-options autonomous-system "+v)
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	for _, v := range d.Get("instance_export").([]interface{}) {
		configSet = append(configSet, setPrefix+"routing-options instance-export "+v.(string))
	}
	for _, v := range d.Get("instance_import").([]interface{}) {
		configSet = append(configSet, setPrefix+"routing-options instance-import "+v.(string))
	}
	if v := d.Get("vtep_source_interface").(string); v != "" {
		configSet = append(configSet, setPrefix+"vtep-source-interface "+v)
	}

	return clt.configSet(configSet, junSess)
}

func readRoutingInstance(instance string, clt *Client, junSess *junosSession) (instanceOptions, error) {
	var confRead instanceOptions

	showConfig, err := clt.command(cmdShowConfig+routingInstancesWS+instance+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = instance
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "instance-type "):
				confRead.instanceType = strings.TrimPrefix(itemTrim, "instance-type ")
			case strings.HasPrefix(itemTrim, "route-distinguisher "):
				confRead.routeDistinguisher = strings.TrimPrefix(itemTrim, "route-distinguisher ")
			case strings.HasPrefix(itemTrim, "routing-options autonomous-system "):
				confRead.as = strings.TrimPrefix(itemTrim, "routing-options autonomous-system ")
			case strings.HasPrefix(itemTrim, "routing-options instance-export "):
				confRead.instanceExport = append(confRead.instanceExport,
					strings.TrimPrefix(itemTrim, "routing-options instance-export "))
			case strings.HasPrefix(itemTrim, "routing-options instance-import "):
				confRead.instanceImport = append(confRead.instanceImport,
					strings.TrimPrefix(itemTrim, "routing-options instance-import "))
			case strings.HasPrefix(itemTrim, "vrf-export "):
				confRead.vrfExport = append(confRead.vrfExport, strings.Trim(strings.TrimPrefix(itemTrim, "vrf-export "), "\""))
			case strings.HasPrefix(itemTrim, "vrf-import "):
				confRead.vrfImport = append(confRead.vrfImport, strings.Trim(strings.TrimPrefix(itemTrim, "vrf-import "), "\""))
			case itemTrim == "vrf-target auto":
				confRead.vrfTargetAuto = true
			case strings.HasPrefix(itemTrim, "vrf-target export "):
				confRead.vrfTargetExport = strings.TrimPrefix(itemTrim, "vrf-target export ")
			case strings.HasPrefix(itemTrim, "vrf-target import "):
				confRead.vrfTargetImport = strings.TrimPrefix(itemTrim, "vrf-target import ")
			case strings.HasPrefix(itemTrim, "vrf-target "):
				confRead.vrfTarget = strings.TrimPrefix(itemTrim, "vrf-target ")
			case strings.HasPrefix(itemTrim, "vtep-source-interface "):
				confRead.vtepSourceIf = strings.TrimPrefix(itemTrim, "vtep-source-interface ")
			case strings.HasPrefix(itemTrim, "interface "):
				confRead.interFace = append(confRead.interFace,
					strings.Split(strings.TrimPrefix(itemTrim, "interface "), " ")[0])
			}
		}
	}

	return confRead, nil
}

func delRoutingInstanceOpts(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)
	setPrefix := delRoutingInstances + d.Get("name").(string) + " "
	configSet = append(configSet,
		setPrefix+"description",
		setPrefix+"routing-options autonomous-system",
		setPrefix+"routing-options instance-export",
		setPrefix+"routing-options instance-import",
		setPrefix+"vtep-source-interface",
	)
	if !d.Get("configure_type_singly").(bool) {
		configSet = append(configSet, setPrefix+"instance-type")
	}
	if !d.Get("configure_rd_vrfopts_singly").(bool) {
		configSet = append(configSet,
			setPrefix+"route-distinguisher",
			setPrefix+"vrf-export",
			setPrefix+"vrf-import",
			setPrefix+"vrf_target",
		)
	}

	return clt.configSet(configSet, junSess)
}

func delRoutingInstance(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, delRoutingInstances+d.Get("name").(string))

	return clt.configSet(configSet, junSess)
}

func fillRoutingInstanceData(d *schema.ResourceData, instanceOptions instanceOptions) {
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
	if tfErr := d.Set("vtep_source_interface", instanceOptions.vtepSourceIf); tfErr != nil {
		panic(tfErr)
	}
	if !d.Get("configure_type_singly").(bool) {
		if tfErr := d.Set("type", instanceOptions.instanceType); tfErr != nil {
			panic(tfErr)
		}
	}
	if !d.Get("configure_rd_vrfopts_singly").(bool) {
		if tfErr := d.Set("route_distinguisher", instanceOptions.routeDistinguisher); tfErr != nil {
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
	}
}
