package providersdk

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type layer2ControlOptions struct {
	nonstopBridging     bool
	bpduBlock           []map[string]interface{}
	macRewriteInterface []map[string]interface{}
}

func resourceLayer2Control() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceLayer2ControlCreate,
		ReadWithoutTimeout:   resourceLayer2ControlRead,
		UpdateWithoutTimeout: resourceLayer2ControlUpdate,
		DeleteWithoutTimeout: resourceLayer2ControlDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLayer2ControlImport,
		},
		Schema: map[string]*schema.Schema{
			"bpdu_block": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(10, 3600),
						},
						"interface": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											value := v.(string)
											if strings.Count(value, ".") > 0 {
												errors = append(errors, fmt.Errorf(
													"%q in %q cannot have a dot", value, k))
											}

											return
										},
									},
									"disable": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"drop": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
					},
				},
			},
			"mac_rewrite_interface": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if strings.Count(value, ".") > 0 {
									errors = append(errors, fmt.Errorf(
										"%q in %q cannot have a dot", value, k))
								}

								return
							},
						},
						"enable_all_ifl": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"protocol": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"nonstop_bridging": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceLayer2ControlCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setLayer2Control(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("layer2_control")

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
	if err := setLayer2Control(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_layer2_control")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("layer2_control")

	return append(diagWarns, resourceLayer2ControlReadWJunSess(d, junSess)...)
}

func resourceLayer2ControlRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceLayer2ControlReadWJunSess(d, junSess)
}

func resourceLayer2ControlReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	layer2ControlOptions, err := readLayer2Control(junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillLayer2ControlData(d, layer2ControlOptions)

	return nil
}

func resourceLayer2ControlUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delLayer2Control(junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setLayer2Control(d, junSess); err != nil {
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
	if err := delLayer2Control(junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setLayer2Control(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_layer2_control")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceLayer2ControlReadWJunSess(d, junSess)...)
}

func resourceLayer2ControlDelete(ctx context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delLayer2Control(junSess); err != nil {
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
	if err := delLayer2Control(junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_layer2_control")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceLayer2ControlImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	layer2ControlOptions, err := readLayer2Control(junSess)
	if err != nil {
		return nil, err
	}
	fillLayer2ControlData(d, layer2ControlOptions)
	d.SetId("layer2_control")
	result[0] = d

	return result, nil
}

func setLayer2Control(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)
	setPrefix := "set protocols layer2-control "

	for _, mBpduBlock := range d.Get("bpdu_block").([]interface{}) {
		configSet = append(configSet, setPrefix+"bpdu-block")
		if mBpduBlock != nil {
			bpduBlock := mBpduBlock.(map[string]interface{})
			if v := bpduBlock["disable_timeout"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"bpdu-block disable-timeout "+strconv.Itoa(v))
			}
			interfaceList := make([]string, 0)
			for _, mIntFace := range bpduBlock["interface"].(*schema.Set).List() {
				intFace := mIntFace.(map[string]interface{})
				if slices.Contains(interfaceList, intFace["name"].(string)) {
					return fmt.Errorf("multiple blocks interface with the same name '%s'", intFace["name"].(string))
				}
				interfaceList = append(interfaceList, intFace["name"].(string))
				setPrefixIntFace := setPrefix + "bpdu-block interface " + intFace["name"].(string) + " "
				configSet = append(configSet, setPrefixIntFace)
				if intFace["disable"].(bool) {
					configSet = append(configSet, setPrefixIntFace+"disable")
				}
				if intFace["drop"].(bool) {
					configSet = append(configSet, setPrefixIntFace+"drop")
				}
			}
		}
	}
	macRewriteInterfaceList := make([]string, 0)
	for _, mMacRewIntFace := range d.Get("mac_rewrite_interface").(*schema.Set).List() {
		macRewIntFace := mMacRewIntFace.(map[string]interface{})
		if slices.Contains(macRewriteInterfaceList, macRewIntFace["name"].(string)) {
			return fmt.Errorf("multiple blocks mac_rewrite_interface with the same name '%s'", macRewIntFace["name"].(string))
		}
		macRewriteInterfaceList = append(macRewriteInterfaceList, macRewIntFace["name"].(string))
		setPrefixIntFace := setPrefix + "mac-rewrite interface " + macRewIntFace["name"].(string) + " "
		configSet = append(configSet, setPrefixIntFace)
		if macRewIntFace["enable_all_ifl"].(bool) {
			configSet = append(configSet, setPrefixIntFace+"enable-all-ifl")
		}
		for _, sProto := range macRewIntFace["protocol"].(*schema.Set).List() {
			configSet = append(configSet, setPrefixIntFace+"protocol "+sProto.(string))
		}
	}
	if d.Get("nonstop_bridging").(bool) {
		configSet = append(configSet, setPrefix+"nonstop-bridging")
	}

	return junSess.ConfigSet(configSet)
}

func readLayer2Control(junSess *junos.Session,
) (confRead layer2ControlOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"protocols layer2-control" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}

	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "bpdu-block"):
				if len(confRead.bpduBlock) == 0 {
					confRead.bpduBlock = append(confRead.bpduBlock, map[string]interface{}{
						"disable_timeout": 0,
						"interface":       make([]map[string]interface{}, 0),
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " disable-timeout "):
					confRead.bpduBlock[0]["disable_timeout"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, " interface "):
					itemTrimFields := strings.Split(itemTrim, " ")
					bpduBlockIntFace := map[string]interface{}{
						"name":    itemTrimFields[0],
						"disable": false,
						"drop":    false,
					}
					confRead.bpduBlock[0]["interface"] = copyAndRemoveItemMapList(
						"name",
						bpduBlockIntFace,
						confRead.bpduBlock[0]["interface"].([]map[string]interface{}),
					)
					balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
					switch {
					case itemTrim == "disable":
						bpduBlockIntFace["disable"] = true
					case itemTrim == "drop":
						bpduBlockIntFace["drop"] = true
					}
					confRead.bpduBlock[0]["interface"] = append(confRead.bpduBlock[0]["interface"].([]map[string]interface{}),
						bpduBlockIntFace)
				}
			case balt.CutPrefixInString(&itemTrim, "mac-rewrite interface "):
				itemTrimFields := strings.Split(itemTrim, " ")
				macRewIntFace := map[string]interface{}{
					"name":           itemTrimFields[0],
					"enable_all_ifl": false,
					"protocol":       make([]string, 0),
				}
				confRead.macRewriteInterface = copyAndRemoveItemMapList("name", macRewIntFace, confRead.macRewriteInterface)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case itemTrim == "enable-all-ifl":
					macRewIntFace["enable_all_ifl"] = true
				case balt.CutPrefixInString(&itemTrim, "protocol "):
					macRewIntFace["protocol"] = append(macRewIntFace["protocol"].([]string), itemTrim)
				}
				confRead.macRewriteInterface = append(confRead.macRewriteInterface, macRewIntFace)
			case itemTrim == "nonstop-bridging":
				confRead.nonstopBridging = true
			}
		}
	}

	return confRead, nil
}

func delLayer2Control(junSess *junos.Session) error {
	configSet := []string{"delete protocols layer2-control"}

	return junSess.ConfigSet(configSet)
}

func fillLayer2ControlData(d *schema.ResourceData, layer2ControlOptions layer2ControlOptions) {
	if tfErr := d.Set("bpdu_block", layer2ControlOptions.bpduBlock); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("mac_rewrite_interface", layer2ControlOptions.macRewriteInterface); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("nonstop_bridging", layer2ControlOptions.nonstopBridging); tfErr != nil {
		panic(tfErr)
	}
}
