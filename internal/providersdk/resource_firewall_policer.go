package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type policerOptions struct {
	filterSpecific bool
	name           string
	ifExceeding    []map[string]interface{}
	then           []map[string]interface{}
}

func resourceFirewallPolicer() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceFirewallPolicerCreate,
		ReadWithoutTimeout:   resourceFirewallPolicerRead,
		UpdateWithoutTimeout: resourceFirewallPolicerUpdate,
		DeleteWithoutTimeout: resourceFirewallPolicerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceFirewallPolicerImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"filter_specific": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"if_exceeding": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"burst_size_limit": {
							Type:     schema.TypeString,
							Required: true,
						},
						"bandwidth_percent": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(1, 100),
							ConflictsWith: []string{"if_exceeding.0.bandwidth_limit"},
						},
						"bandwidth_limit": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"if_exceeding.0.bandwidth_percent"},
						},
					},
				},
			},
			"then": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"discard": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"then.0.out_of_profile", "then.0.forwarding_class"},
						},
						"forwarding_class": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"then.0.discard"},
						},
						"loss_priority": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"then.0.discard"},
						},
						"out_of_profile": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"then.0.discard"},
						},
					},
				},
			},
		},
	}
}

func resourceFirewallPolicerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setFirewallPolicer(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	firewallPolicerExists, err := checkFirewallPolicerExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if firewallPolicerExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("firewall policer %v already exists", d.Get("name").(string)))...)
	}

	if err := setFirewallPolicer(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_firewall_policer")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	firewallPolicerExists, err = checkFirewallPolicerExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if firewallPolicerExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("firewall policer %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceFirewallPolicerReadWJunSess(d, junSess)...)
}

func resourceFirewallPolicerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceFirewallPolicerReadWJunSess(d, junSess)
}

func resourceFirewallPolicerReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	policerOptions, err := readFirewallPolicer(d.Get("name").(string), junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if policerOptions.name == "" {
		d.SetId("")
	} else {
		fillFirewallPolicerData(d, policerOptions)
	}

	return nil
}

func resourceFirewallPolicerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delFirewallPolicer(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setFirewallPolicer(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delFirewallPolicer(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setFirewallPolicer(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_firewall_policer")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceFirewallPolicerReadWJunSess(d, junSess)...)
}

func resourceFirewallPolicerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delFirewallPolicer(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delFirewallPolicer(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_firewall_policer")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceFirewallPolicerImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	firewallPolicerExists, err := checkFirewallPolicerExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !firewallPolicerExists {
		return nil, fmt.Errorf("don't find firewall policer with id '%v' (id must be <name>)", d.Id())
	}
	policerOptions, err := readFirewallPolicer(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillFirewallPolicerData(d, policerOptions)

	result[0] = d

	return result, nil
}

func checkFirewallPolicerExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "firewall policer " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setFirewallPolicer(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set firewall policer " + d.Get("name").(string)
	if d.Get("filter_specific").(bool) {
		configSet = append(configSet, setPrefix+" filter-specific")
	}
	for _, ifExceeding := range d.Get("if_exceeding").([]interface{}) {
		ifExceedingMap := ifExceeding.(map[string]interface{})
		configSet = append(configSet, setPrefix+
			" if-exceeding burst-size-limit "+ifExceedingMap["burst_size_limit"].(string))
		if ifExceedingMap["bandwidth_percent"].(int) != 0 {
			configSet = append(configSet, setPrefix+
				" if-exceeding bandwidth-percent "+strconv.Itoa(ifExceedingMap["bandwidth_percent"].(int)))
		}
		if ifExceedingMap["bandwidth_limit"].(string) != "" {
			configSet = append(configSet, setPrefix+
				" if-exceeding bandwidth-limit "+ifExceedingMap["bandwidth_limit"].(string))
		}
	}
	for _, then := range d.Get("then").([]interface{}) {
		if then != nil {
			thenMap := then.(map[string]interface{})
			if thenMap["discard"].(bool) {
				configSet = append(configSet, setPrefix+
					" then discard")
			}
			if thenMap["forwarding_class"].(string) != "" {
				configSet = append(configSet, setPrefix+
					" then forwarding-class "+thenMap["forwarding_class"].(string))
			}
			if thenMap["loss_priority"].(string) != "" {
				configSet = append(configSet, setPrefix+
					" then loss-priority "+thenMap["loss_priority"].(string))
			}
			if thenMap["out_of_profile"].(bool) {
				configSet = append(configSet, setPrefix+
					" then out-of-profile")
			}
		}
	}

	return junSess.ConfigSet(configSet)
}

func readFirewallPolicer(name string, junSess *junos.Session,
) (confRead policerOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "firewall policer " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == "filter-specific":
				confRead.filterSpecific = true
			case balt.CutPrefixInString(&itemTrim, "if-exceeding "):
				if len(confRead.ifExceeding) == 0 {
					confRead.ifExceeding = append(confRead.ifExceeding, map[string]interface{}{
						"burst_size_limit":  "",
						"bandwidth_percent": 0,
						"bandwidth_limit":   "",
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "burst-size-limit "):
					confRead.ifExceeding[0]["burst_size_limit"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, "bandwidth-percent "):
					confRead.ifExceeding[0]["bandwidth_percent"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "bandwidth-limit "):
					confRead.ifExceeding[0]["bandwidth_limit"] = itemTrim
				}
			case balt.CutPrefixInString(&itemTrim, "then "):
				if len(confRead.then) == 0 {
					confRead.then = append(confRead.then, map[string]interface{}{
						"discard":          false,
						"forwarding_class": "",
						"loss_priority":    "",
						"out_of_profile":   false,
					})
				}
				switch {
				case itemTrim == "discard":
					confRead.then[0]["discard"] = true
				case balt.CutPrefixInString(&itemTrim, "forwarding-class "):
					confRead.then[0]["forwarding_class"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, "loss-priority "):
					confRead.then[0]["loss_priority"] = itemTrim
				case itemTrim == "out-of-profile":
					confRead.then[0]["out_of_profile"] = true
				}
			}
		}
	}

	return confRead, nil
}

func delFirewallPolicer(policer string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete firewall policer "+policer)

	return junSess.ConfigSet(configSet)
}

func fillFirewallPolicerData(d *schema.ResourceData, policerOptions policerOptions) {
	if tfErr := d.Set("name", policerOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("filter_specific", policerOptions.filterSpecific); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("if_exceeding", policerOptions.ifExceeding); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("then", policerOptions.then); tfErr != nil {
		panic(tfErr)
	}
}
