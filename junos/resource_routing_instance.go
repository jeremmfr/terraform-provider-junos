package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type instanceOptions struct {
	name         string
	instanceType string
	as           string
}

func resourceRoutingInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoutingInstanceCreate,
		ReadContext:   resourceRoutingInstanceRead,
		UpdateContext: resourceRoutingInstanceUpdate,
		DeleteContext: resourceRoutingInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRoutingInstanceImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{"default"}, 64),
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
		},
	}
}

func resourceRoutingInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	routingInstanceExists, err := checkRoutingInstanceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if routingInstanceExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("routing-instance %v already exists", d.Get("name").(string)))
	}
	if err := setRoutingInstance(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_routing_instance", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	routingInstanceExists, err = checkRoutingInstanceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if routingInstanceExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("routing-instance %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceRoutingInstanceReadWJnprSess(d, m, jnprSess)...)
}
func resourceRoutingInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceRoutingInstanceReadWJnprSess(d, m, jnprSess)
}
func resourceRoutingInstanceReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	instanceOptions, err := readRoutingInstance(d.Get("name").(string), m, jnprSess)
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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)

	if err := delRoutingInstanceOpts(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setRoutingInstance(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_routing_instance", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRoutingInstanceReadWJnprSess(d, m, jnprSess)...)
}
func resourceRoutingInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delRoutingInstance(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_routing_instance", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceRoutingInstanceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	routingInstanceExists, err := checkRoutingInstanceExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !routingInstanceExists {
		return nil, fmt.Errorf("don't find routing instance with id '%v' (id must be <name>)", d.Id())
	}
	instanceOptions, err := readRoutingInstance(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillRoutingInstanceData(d, instanceOptions)
	result[0] = d

	return result, nil
}

func checkRoutingInstanceExists(instance string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	routingInstanceConfig, err := sess.command("show configuration routing-instances "+instance+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if routingInstanceConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setRoutingInstance(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set routing-instances " + d.Get("name").(string) + " "
	if d.Get("type").(string) != "" {
		configSet = append(configSet, setPrefix+"instance-type "+d.Get("type").(string))
	} else {
		configSet = append(configSet, setPrefix)
	}
	if d.Get("as").(string) != "" {
		configSet = append(configSet, setPrefix+
			"routing-options autonomous-system "+d.Get("as").(string))
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readRoutingInstance(instance string, m interface{}, jnprSess *NetconfObject) (instanceOptions, error) {
	sess := m.(*Session)
	var confRead instanceOptions

	instanceConfig, err := sess.command("show configuration"+
		" routing-instances "+instance+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if instanceConfig != emptyWord {
		confRead.name = instance
		for _, item := range strings.Split(instanceConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "instance-type "):
				confRead.instanceType = strings.TrimPrefix(itemTrim, "instance-type ")
			case strings.HasPrefix(itemTrim, "routing-options autonomous-system "):
				confRead.as = strings.TrimPrefix(itemTrim, "routing-options autonomous-system ")
			}
		}
	}

	return confRead, nil
}
func delRoutingInstanceOpts(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := "delete routing-instances " + d.Get("name").(string) + " "
	configSet = append(configSet,
		setPrefix+"instance-type",
		setPrefix+"routing-options autonomous-system")
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func delRoutingInstance(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete routing-instances "+d.Get("name").(string))
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func fillRoutingInstanceData(d *schema.ResourceData, instanceOptions instanceOptions) {
	if tfErr := d.Set("name", instanceOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("type", instanceOptions.instanceType); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as", instanceOptions.as); tfErr != nil {
		panic(tfErr)
	}
}
