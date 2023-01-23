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

type natDestinationPoolOptions struct {
	addressPort     int
	name            string
	address         string
	addressTo       string
	description     string
	routingInstance string
}

func resourceSecurityNatDestinationPool() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityNatDestinationPoolCreate,
		ReadWithoutTimeout:   resourceSecurityNatDestinationPoolRead,
		UpdateWithoutTimeout: resourceSecurityNatDestinationPoolUpdate,
		DeleteWithoutTimeout: resourceSecurityNatDestinationPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityNatDestinationPoolImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
			},
			"address": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateIPMaskFunc(),
			},
			"address_port": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validation.IntBetween(1, 65535),
				ConflictsWith: []string{"address_to"},
			},
			"address_to": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateIPMaskFunc(),
				ConflictsWith:    []string{"address_port"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
		},
	}
}

func resourceSecurityNatDestinationPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := setSecurityNatDestinationPool(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if !junos.CheckCompatibilitySecurity(junSess) {
		return diag.FromErr(fmt.Errorf("security nat destination pool not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityNatDestinationPoolExists, err := checkSecurityNatDestinationPoolExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatDestinationPoolExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security nat destination pool %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityNatDestinationPool(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_security_nat_destination_pool", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityNatDestinationPoolExists, err = checkSecurityNatDestinationPoolExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatDestinationPoolExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security nat destination pool %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityNatDestinationPoolReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityNatDestinationPoolRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceSecurityNatDestinationPoolReadWJunSess(d, clt, junSess)
}

func resourceSecurityNatDestinationPoolReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	natDestinationPoolOptions, err := readSecurityNatDestinationPool(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if natDestinationPoolOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityNatDestinationPoolData(d, natDestinationPoolOptions)
	}

	return nil
}

func resourceSecurityNatDestinationPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delSecurityNatDestinationPool(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityNatDestinationPool(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityNatDestinationPool(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityNatDestinationPool(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_security_nat_destination_pool", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityNatDestinationPoolReadWJunSess(d, clt, junSess)...)
}

func resourceSecurityNatDestinationPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delSecurityNatDestinationPool(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityNatDestinationPool(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_security_nat_destination_pool", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityNatDestinationPoolImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)

	securityNatDestinationPoolExists, err := checkSecurityNatDestinationPoolExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !securityNatDestinationPoolExists {
		return nil, fmt.Errorf("don't find nat destination pool with id '%v' (id must be <name>)", d.Id())
	}
	natDestinationPoolOptions, err := readSecurityNatDestinationPool(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityNatDestinationPoolData(d, natDestinationPoolOptions)

	result[0] = d

	return result, nil
}

func checkSecurityNatDestinationPoolExists(name string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"security nat destination pool "+name+junos.PipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityNatDestinationPool(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security nat destination pool " + d.Get("name").(string)
	configSet = append(configSet, setPrefix+" address "+d.Get("address").(string))
	if d.Get("address_port").(int) != 0 {
		configSet = append(configSet, setPrefix+" address port "+strconv.Itoa(d.Get("address_port").(int)))
	}
	if d.Get("address_to").(string) != "" {
		configSet = append(configSet, setPrefix+" address to "+d.Get("address_to").(string))
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+" description \""+v+"\"")
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, setPrefix+" routing-instance "+d.Get("routing_instance").(string))
	}

	return clt.ConfigSet(configSet, junSess)
}

func readSecurityNatDestinationPool(name string, clt *junos.Client, junSess *junos.Session,
) (confRead natDestinationPoolOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"security nat destination pool "+name+junos.PipeDisplaySetRelative, junSess)
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
			case balt.CutPrefixInString(&itemTrim, "address port "):
				confRead.addressPort, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "address to "):
				confRead.addressTo = itemTrim
			case balt.CutPrefixInString(&itemTrim, "address "):
				confRead.address = itemTrim
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "routing-instance "):
				confRead.routingInstance = itemTrim
			}
		}
	}

	return confRead, nil
}

func delSecurityNatDestinationPool(natDestinationPool string, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat destination pool "+natDestinationPool)

	return clt.ConfigSet(configSet, junSess)
}

func fillSecurityNatDestinationPoolData(d *schema.ResourceData, natDestinationPoolOptions natDestinationPoolOptions) {
	if tfErr := d.Set("name", natDestinationPoolOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address", natDestinationPoolOptions.address); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_port", natDestinationPoolOptions.addressPort); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_to", natDestinationPoolOptions.addressTo); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", natDestinationPoolOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", natDestinationPoolOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
}
