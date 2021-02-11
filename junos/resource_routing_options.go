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

type routingOptionsOptions struct {
	autonomousSystem []map[string]interface{}
	gracefulRestart  []map[string]interface{}
}

func resourceRoutingOptions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoutingOptionsCreate,
		ReadContext:   resourceRoutingOptionsRead,
		UpdateContext: resourceRoutingOptionsUpdate,
		DeleteContext: resourceRoutingOptionsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRoutingOptionsImport,
		},
		Schema: map[string]*schema.Schema{
			"autonomous_system": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number": {
							Type:     schema.TypeString,
							Required: true,
						},
						"asdot_notation": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"loops": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 10),
						},
					},
				},
			},
			"graceful_restart": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"restart_duration": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(120, 10000),
						},
					},
				},
			},
		},
	}
}

func resourceRoutingOptionsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)

	if err := setRoutingOptions(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_routing_options", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("routing_options")

	return append(diagWarns, resourceRoutingOptionsReadWJnprSess(d, m, jnprSess)...)
}
func resourceRoutingOptionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceRoutingOptionsReadWJnprSess(d, m, jnprSess)
}
func resourceRoutingOptionsReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	routingOptionsOptions, err := readRoutingOptions(m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillRoutingOptions(d, routingOptionsOptions)

	return nil
}
func resourceRoutingOptionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delRoutingOptions(m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setRoutingOptions(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_routing_options", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRoutingOptionsReadWJnprSess(d, m, jnprSess)...)
}
func resourceRoutingOptionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
func resourceRoutingOptionsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	routingOptionsOptions, err := readRoutingOptions(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillRoutingOptions(d, routingOptionsOptions)
	d.SetId("routing_options")
	result[0] = d

	return result, nil
}

func setRoutingOptions(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set routing-options "
	configSet := make([]string, 0)

	for _, as := range d.Get("autonomous_system").([]interface{}) {
		asM := as.(map[string]interface{})
		configSet = append(configSet, setPrefix+"autonomous-system "+asM["number"].(string))
		if asM["asdot_notation"].(bool) {
			configSet = append(configSet, setPrefix+"autonomous-system asdot-notation")
		}
		if asM["loops"].(int) > 0 {
			configSet = append(configSet, setPrefix+"autonomous-system loops "+strconv.Itoa(asM["loops"].(int)))
		}
	}
	for _, grR := range d.Get("graceful_restart").([]interface{}) {
		configSet = append(configSet, setPrefix+"graceful-restart")
		if grR != nil {
			grRM := grR.(map[string]interface{})
			if grRM["disable"].(bool) {
				configSet = append(configSet, setPrefix+"graceful-restart disable")
			}
			if grRM["restart_duration"].(int) > 0 {
				configSet = append(configSet, setPrefix+"graceful-restart restart-duration "+
					strconv.Itoa(grRM["restart_duration"].(int)))
			}
		}
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func delRoutingOptions(m interface{}, jnprSess *NetconfObject) error {
	listLinesToDelete := []string{
		"autonomous-system",
		"graceful-restart",
	}
	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := "delete routing-options "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readRoutingOptions(m interface{}, jnprSess *NetconfObject) (routingOptionsOptions, error) {
	sess := m.(*Session)
	var confRead routingOptionsOptions

	routingOptionsConfig, err := sess.command("show configuration routing-options"+
		" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if routingOptionsConfig != emptyWord {
		for _, item := range strings.Split(routingOptionsConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "autonomous-system "):
				if len(confRead.autonomousSystem) == 0 {
					confRead.autonomousSystem = append(confRead.autonomousSystem, map[string]interface{}{
						"number":         "",
						"asdot_notation": false,
						"loops":          0,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "autonomous-system loops "):
					var err error
					confRead.autonomousSystem[0]["loops"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrim, "autonomous-system loops "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case itemTrim == "autonomous-system asdot-notation":
					confRead.autonomousSystem[0]["asdot_notation"] = true
				default:
					confRead.autonomousSystem[0]["number"] = strings.TrimPrefix(itemTrim, "autonomous-system ")
				}
			case strings.HasPrefix(itemTrim, "graceful-restart"):
				if len(confRead.gracefulRestart) == 0 {
					confRead.gracefulRestart = append(confRead.gracefulRestart, map[string]interface{}{
						"disable":          false,
						"restart_duration": 0,
					})
				}
				switch {
				case itemTrim == "graceful-restart disable":
					confRead.gracefulRestart[0]["disable"] = true
				case strings.HasPrefix(itemTrim, "graceful-restart restart-duration "):
					var err error
					confRead.gracefulRestart[0]["restart_duration"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrim, "graceful-restart restart-duration "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				}
			}
		}
	}

	return confRead, nil
}

func fillRoutingOptions(d *schema.ResourceData, routingOptionsOptions routingOptionsOptions) {
	if tfErr := d.Set("autonomous_system", routingOptionsOptions.autonomousSystem); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("graceful_restart", routingOptionsOptions.gracefulRestart); tfErr != nil {
		panic(tfErr)
	}
}
