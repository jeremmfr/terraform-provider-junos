package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type zoneOptions struct {
	name             string
	inboundServices  []string
	inboundProtocols []string
	addressBook      []map[string]interface{}
	addressBookSet   []map[string]interface{}
}

func resourceSecurityZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityZoneCreate,
		Read:   resourceSecurityZoneRead,
		Update: resourceSecurityZoneUpdate,
		Delete: resourceSecurityZoneDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityZoneImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"inbound_services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"inbound_protocols": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"address_book": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateNameObjectJunos(),
						},
						"network": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								err := validateNetwork(value)
								if err != nil {
									errors = append(errors, fmt.Errorf(
										"%q error for validate %q : %q", k, value, err))
								}
								return
							},
						},
					},
				},
			},
			"address_book_set": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateNameObjectJunos(),
						},
						"address": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceSecurityZoneCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security zone not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	securityZoneExists, err := checkSecurityZonesExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if securityZoneExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("security zone %v already exists", d.Get("name").(string))
	}

	err = setSecurityZone(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("create resource junos_security_zone", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	mutex.Lock()
	securityZoneExists, err = checkSecurityZonesExists(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if securityZoneExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("security zone %v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourceSecurityZoneRead(d, m)
}
func resourceSecurityZoneRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	zoneOptions, err := readSecurityZone(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if zoneOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityZoneData(d, zoneOptions)
	}
	return nil
}
func resourceSecurityZoneUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	if d.HasChange("inbound_services") {
		err = delSecurityZoneElement("host-inbound-traffic system-services", d.Get("name").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)
			return err
		}
	}
	if d.HasChange("inbound_protocols") {
		err = delSecurityZoneElement("host-inbound-traffic protocols", d.Get("name").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)
			return err
		}
	}
	if d.HasChange("address_book") || d.HasChange("address_book_set") {
		err = delSecurityZoneElement("address-book", d.Get("name").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)
			return err
		}
	}
	err = setSecurityZone(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("update resource junos_security_zone", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourceSecurityZoneRead(d, m)
}
func resourceSecurityZoneDelete(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delSecurityZone(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("delete resource junos_security_zone", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
}
func resourceSecurityZoneImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	securityZoneExists, err := checkSecurityZonesExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityZoneExists {
		return nil, fmt.Errorf("don't find zone with id '%v' (id must be <name>)", d.Id())
	}
	zoneOptions, err := readSecurityZone(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityZoneData(d, zoneOptions)

	result[0] = d
	return result, nil
}

func checkSecurityZonesExists(zone string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	zoneConfig, err := sess.command("show configuration security zones security-zone "+zone+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if zoneConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setSecurityZone(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security zones security-zone " + d.Get("name").(string)
	configSet = append(configSet, setPrefix+"\n")

	for _, v := range d.Get("inbound_services").([]interface{}) {
		configSet = append(configSet, setPrefix+" host-inbound-traffic system-services "+v.(string)+"\n")
	}
	for _, v := range d.Get("inbound_protocols").([]interface{}) {
		configSet = append(configSet, setPrefix+" host-inbound-traffic protocols "+v.(string)+"\n")
	}

	for _, v := range d.Get("address_book").([]interface{}) {
		addressBook := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+" address-book address "+
			addressBook["name"].(string)+" "+addressBook["network"].(string)+"\n")
	}

	for _, v := range d.Get("address_book_set").([]interface{}) {
		addressBookSet := v.(map[string]interface{})
		for _, addressBookSetAddress := range addressBookSet["address"].([]interface{}) {
			configSet = append(configSet, setPrefix+" address-book address-set "+addressBookSet["name"].(string)+
				" address "+addressBookSetAddress.(string)+"\n")
		}
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readSecurityZone(zone string, m interface{}, jnprSess *NetconfObject) (zoneOptions, error) {
	sess := m.(*Session)
	var confRead zoneOptions

	zoneConfig, err := sess.command("show configuration"+
		" security zones security-zone "+zone+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	inboundServices := make([]string, 0)
	inboundProtocols := make([]string, 0)
	addressBook := make([]map[string]interface{}, 0)
	addressBookSet := make([]map[string]interface{}, 0)
	if zoneConfig != emptyWord {
		confRead.name = zone
		for _, item := range strings.Split(zoneConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "host-inbound-traffic system-services "):
				inboundServices = append(inboundServices, strings.TrimPrefix(itemTrim,
					"host-inbound-traffic system-services "))
			case strings.HasPrefix(itemTrim, "host-inbound-traffic protocols "):
				inboundProtocols = append(inboundProtocols, strings.TrimPrefix(itemTrim,
					"host-inbound-traffic protocols "))
			case strings.HasPrefix(itemTrim, "address-book address "):
				address := strings.TrimPrefix(itemTrim, "address-book address ")
				addressWords := strings.Split(address, " ")
				// addressWords[0] = name of address
				// addressWords[1] = network
				m := make(map[string]interface{})
				m["name"] = addressWords[0]
				m["network"] = addressWords[1]
				addressBook = append(addressBook, m)
			case strings.HasPrefix(itemTrim, "address-book address-set "):
				address := strings.TrimPrefix(itemTrim, "address-book address-set ")
				addressWords := strings.Split(address, " ")
				// addressWords[0] = name of address-set
				// addressWords[1] = "address"
				// addressWords[2] = name of address
				m := map[string]interface{}{
					"name":    addressWords[0],
					"address": make([]string, 0),
				}
				// search if name of address-set already create
				m, addressBookSet = copyAndRemoveItemMapList("name", false, m, addressBookSet)
				// append new address find
				m["address"] = append(m["address"].([]string), addressWords[2])
				addressBookSet = append(addressBookSet, m)
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	confRead.inboundServices = inboundServices
	confRead.inboundProtocols = inboundProtocols
	confRead.addressBook = addressBook
	confRead.addressBookSet = addressBookSet

	return confRead, nil
}
func delSecurityZoneElement(element string, zone string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security zones security-zone "+zone+" "+element+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func delSecurityZone(zone string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security zones security-zone "+zone+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}

func fillSecurityZoneData(d *schema.ResourceData, zoneOptions zoneOptions) {
	tfErr := d.Set("name", zoneOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("inbound_services", zoneOptions.inboundServices)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("inbound_protocols", zoneOptions.inboundProtocols)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("address_book", zoneOptions.addressBook)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("address_book_set", zoneOptions.addressBookSet)
	if tfErr != nil {
		panic(tfErr)
	}
}
