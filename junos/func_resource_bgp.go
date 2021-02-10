package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	jdecode "github.com/jeremmfr/junosdecode"
)

type bgpOptions struct {
	acceptRemoteNexthop          bool
	advertiseExternal            bool
	advertiseExternalConditional bool
	advertiseInactive            bool
	advertisePeerAs              bool
	asOverride                   bool
	damping                      bool
	localAsPrivate               bool
	localAsAlias                 bool
	localAsNoPrependGlobalAs     bool
	logUpdown                    bool
	metricOutIgp                 bool
	metricOutIgpDelayMedUpdate   bool
	metricOutMinimumIgp          bool
	mtuDiscovery                 bool
	multihop                     bool
	multipath                    bool
	noAdvertisePeerAs            bool
	removePrivate                bool
	passive                      bool
	holdTime                     int
	localAsLoops                 int
	localPreference              int
	metricOut                    int
	metricOutIgpOffset           int
	metricOutMinimumIgpOffset    int
	outDelay                     int
	preference                   int
	authenticationAlgorithm      string
	authenticationKey            string
	authenticationKeyChain       string
	bgpType                      string // group only
	localAddress                 string
	localAs                      string
	localInterface               string
	name                         string // group parameter for neighbor
	ip                           string // for neighbor only
	peerAs                       string
	routingInstance              string
	exportPolicy                 []string
	importPolicy                 []string
	bfdLivenessDetection         []map[string]interface{}
	familyInet                   []map[string]interface{}
	familyInet6                  []map[string]interface{}
	gracefulRestart              []map[string]interface{}
}

func delBgpOpts(d *schema.ResourceData, typebgp string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := deleteWord + " "
	switch typebgp {
	case "group":
		if d.Get("routing_instance").(string) == defaultWord {
			delPrefix += "protocols bgp group " + d.Get("name").(string) + " "
		} else {
			delPrefix += "routing-instances " + d.Get("routing_instance").(string) +
				" protocols bgp group " + d.Get("name").(string) + " "
		}
	case "neighbor":
		if d.Get("routing_instance").(string) == defaultWord {
			delPrefix += "protocols bgp group " + d.Get("group").(string) +
				" neighbor " + d.Get("ip").(string) + " "
		} else {
			delPrefix += "routing-instances " + d.Get("routing_instance").(string) +
				" protocols bgp group " + d.Get("group").(string) +
				" neighbor " + d.Get("ip").(string) + " "
		}
	}
	configSet = append(configSet,
		delPrefix+"accept-remote-nexthop",
		delPrefix+"advertise-external",
		delPrefix+"advertise-inactive",
		delPrefix+"advertise-peer-as",
		delPrefix+"no-advertise-peer-as",
		delPrefix+"as-override",
		delPrefix+"authentication-algorithm",
		delPrefix+"authentication-key",
		delPrefix+"authentication-key-chain",
		delPrefix+"bfd-liveness-detection",
		delPrefix+"damping",
		delPrefix+"export",
		delPrefix+"family inet",
		delPrefix+"family inet6",
		delPrefix+"graceful-restart",
		delPrefix+"hold-time",
		delPrefix+"import",
		delPrefix+"local-address",
		delPrefix+"local-as",
		delPrefix+"local-interface",
		delPrefix+"local-preference",
		delPrefix+"log-updown",
		delPrefix+"metric-out",
		delPrefix+"mtu-discovery",
		delPrefix+"multihop",
		delPrefix+"multipath",
		delPrefix+"out-delay",
		delPrefix+"passive",
		delPrefix+"peer-as",
		delPrefix+"preference",
		delPrefix+"remove-private",
		delPrefix+"type",
	)

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func setBgpOptsSimple(setPrefix string, d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{setPrefix}
	if d.Get("accept_remote_nexthop").(bool) {
		configSet = append(configSet, setPrefix+"accept-remote-nexthop")
	}
	if d.Get("advertise_external").(bool) {
		configSet = append(configSet, setPrefix+"advertise-external")
	}
	if d.Get("advertise_external_conditional").(bool) {
		configSet = append(configSet, setPrefix+"advertise-external conditional")
	}
	if d.Get("advertise_inactive").(bool) {
		configSet = append(configSet, setPrefix+"advertise-inactive")
	}
	if d.Get("advertise_peer_as").(bool) {
		configSet = append(configSet, setPrefix+"advertise-peer-as")
	}
	if d.Get("as_override").(bool) {
		configSet = append(configSet, setPrefix+"as-override")
	}
	if d.Get("authentication_algorithm").(string) != "" {
		configSet = append(configSet, setPrefix+"authentication-algorithm "+d.Get("authentication_algorithm").(string))
	}
	if d.Get("authentication_key").(string) != "" {
		configSet = append(configSet, setPrefix+"authentication-key "+d.Get("authentication_key").(string))
	}
	if d.Get("authentication_key_chain").(string) != "" {
		configSet = append(configSet, setPrefix+"authentication-key-chain "+d.Get("authentication_key_chain").(string))
	}
	if d.Get("damping").(bool) {
		configSet = append(configSet, setPrefix+"damping")
	}
	for _, v := range d.Get("export").([]interface{}) {
		configSet = append(configSet, setPrefix+"export "+v.(string))
	}
	if d.Get("hold_time").(int) != 0 {
		configSet = append(configSet, setPrefix+"hold-time "+strconv.Itoa(d.Get("hold_time").(int)))
	}
	for _, v := range d.Get("import").([]interface{}) {
		configSet = append(configSet, setPrefix+"import "+v.(string))
	}
	if d.Get("local_address").(string) != "" {
		configSet = append(configSet, setPrefix+"local-address "+d.Get("local_address").(string))
	}
	if d.Get("local_as").(string) != "" {
		configSet = append(configSet, setPrefix+"local-as "+d.Get("local_as").(string))
	}
	if d.Get("local_as_alias").(bool) {
		configSet = append(configSet, setPrefix+"local-as alias")
	}
	if d.Get("local_as_loops").(int) != 0 {
		configSet = append(configSet, setPrefix+"local-as loops "+strconv.Itoa(d.Get("local_as_loops").(int)))
	}
	if d.Get("local_as_no_prepend_global_as").(bool) {
		configSet = append(configSet, setPrefix+"local-as no-prepend-global-as")
	}
	if d.Get("local_as_private").(bool) {
		configSet = append(configSet, setPrefix+"local-as private")
	}
	if d.Get("local_interface").(string) != "" {
		configSet = append(configSet, setPrefix+"local-interface "+d.Get("local_interface").(string))
	}
	if d.Get("local_preference").(int) != -1 {
		configSet = append(configSet, setPrefix+"local-preference "+strconv.Itoa(d.Get("local_preference").(int)))
	}
	if d.Get("log_updown").(bool) {
		configSet = append(configSet, setPrefix+"log-updown")
	}
	if d.Get("metric_out").(int) != -1 {
		configSet = append(configSet, setPrefix+"metric-out "+strconv.Itoa(d.Get("metric_out").(int)))
	}
	if d.Get("metric_out_igp").(bool) {
		configSet = append(configSet, setPrefix+"metric-out igp")
	}
	if d.Get("metric_out_igp_delay_med_update").(bool) {
		tfErr := d.Set("metric_out_igp", true)
		if tfErr != nil {
			panic(tfErr)
		}
		configSet = append(configSet, setPrefix+"metric-out igp delay-med-update")
	}
	if d.Get("metric_out_igp_offset").(int) != 0 {
		tfErr := d.Set("metric_out_igp", true)
		if tfErr != nil {
			panic(tfErr)
		}
		configSet = append(configSet, setPrefix+"metric-out igp "+strconv.Itoa(d.Get("metric_out_igp_offset").(int)))
	}
	if d.Get("metric_out_minimum_igp").(bool) {
		configSet = append(configSet, setPrefix+"metric-out minimum-igp")
	}
	if d.Get("metric_out_minimum_igp_offset").(int) != 0 {
		tfErr := d.Set("metric_out_minimum_igp", true)
		if tfErr != nil {
			panic(tfErr)
		}
		configSet = append(configSet, setPrefix+"metric-out minimum-igp "+
			strconv.Itoa(d.Get("metric_out_minimum_igp_offset").(int)))
	}
	if d.Get("mtu_discovery").(bool) {
		configSet = append(configSet, setPrefix+"mtu-discovery")
	}
	if d.Get("multihop").(bool) {
		configSet = append(configSet, setPrefix+"multihop")
	}
	if d.Get("multipath").(bool) {
		configSet = append(configSet, setPrefix+"multipath")
	}
	if d.Get("no_advertise_peer_as").(bool) {
		configSet = append(configSet, setPrefix+"no-advertise-peer-as")
	}
	if d.Get("out_delay").(int) != 0 {
		configSet = append(configSet, setPrefix+"out-delay "+strconv.Itoa(d.Get("out_delay").(int)))
	}
	if d.Get("passive").(bool) {
		configSet = append(configSet, setPrefix+"passive")
	}
	if d.Get("peer_as").(string) != "" {
		configSet = append(configSet, setPrefix+"peer-as "+d.Get("peer_as").(string))
	}
	if d.Get("preference").(int) != -1 {
		configSet = append(configSet, setPrefix+"preference "+strconv.Itoa(d.Get("preference").(int)))
	}
	if d.Get("remove_private").(bool) {
		configSet = append(configSet, setPrefix+"remove-private")
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readBgpOptsSimple(item string, confRead *bgpOptions) error {
	var err error
	if item == "accept-remote-nexthop" {
		confRead.acceptRemoteNexthop = true
	}
	if item == "advertise-external" {
		confRead.advertiseExternal = true
	}
	if item == "advertise-external conditional" {
		confRead.advertiseExternalConditional = true
	}
	if item == "advertise-inactive" {
		confRead.advertiseInactive = true
	}
	if item == "advertise-peer-as" {
		confRead.advertisePeerAs = true
	}
	if item == "as-override" {
		confRead.asOverride = true
	}
	if strings.HasPrefix(item, "authentication-algorithm ") {
		confRead.authenticationAlgorithm = strings.TrimPrefix(item, "authentication-algorithm ")
	}
	if strings.HasPrefix(item, "authentication-key ") {
		confRead.authenticationKey, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(item, "authentication-key "), "\""))
		if err != nil {
			return fmt.Errorf("failed to decode authentication-key : %w", err)
		}
	}
	if strings.HasPrefix(item, "authentication-key-chain ") {
		confRead.authenticationKeyChain = strings.TrimPrefix(item, "authentication-key-chain ")
	}
	if item == "damping" {
		confRead.damping = true
	}
	if strings.HasPrefix(item, "export ") {
		confRead.exportPolicy = append(confRead.exportPolicy, strings.TrimPrefix(item, "export "))
	}
	if strings.HasPrefix(item, "hold-time ") {
		confRead.holdTime, err = strconv.Atoi(strings.TrimPrefix(item, "hold-time "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", item, err)
		}
	}
	if strings.HasPrefix(item, "import ") {
		confRead.importPolicy = append(confRead.importPolicy, strings.TrimPrefix(item, "import "))
	}
	if strings.HasPrefix(item, "local-address ") {
		confRead.localAddress = strings.TrimPrefix(item, "local-address ")
	}
	if strings.HasPrefix(item, "local-as ") {
		switch {
		case strings.HasSuffix(item, " private"):
			confRead.localAsPrivate = true
		case strings.HasSuffix(item, " alias"):
			confRead.localAsAlias = true
		case strings.HasSuffix(item, " no-prepend-global-as"):
			confRead.localAsNoPrependGlobalAs = true
		case strings.HasPrefix(item, "local-as loops "):
			confRead.localAsLoops, err = strconv.Atoi(strings.TrimPrefix(item, "local-as loops "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", item, err)
			}
		default:
			confRead.localAs = strings.TrimPrefix(item, "local-as ")
		}
	}
	if strings.HasPrefix(item, "local-interface ") {
		confRead.localInterface = strings.TrimPrefix(item, "local-interface ")
	}
	if strings.HasPrefix(item, "local-preference ") {
		confRead.localPreference, err = strconv.Atoi(strings.TrimPrefix(item, "local-preference "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", item, err)
		}
	}
	if item == "log-updown" {
		confRead.logUpdown = true
	}
	if strings.HasPrefix(item, "metric-out ") {
		if !strings.Contains(item, "igp") {
			confRead.metricOut, err = strconv.Atoi(strings.TrimPrefix(item, "metric-out "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", item, err)
			}
		} else {
			if strings.HasPrefix(item, "metric-out igp") {
				confRead.metricOutIgp = true
				if item == "metric-out igp delay-med-update" {
					confRead.metricOutIgpDelayMedUpdate = true
				} else if strings.HasPrefix(item, "metric-out igp ") {
					confRead.metricOutIgpOffset, err = strconv.Atoi(strings.TrimPrefix(item, "metric-out igp "))
					if err != nil {
						return fmt.Errorf("failed to convert value from '%s' to integer : %w", item, err)
					}
				}
			} else {
				confRead.metricOutMinimumIgp = true
				if strings.HasPrefix(item, "metric-out minimum-igp ") {
					confRead.metricOutMinimumIgpOffset, err = strconv.Atoi(strings.TrimPrefix(item, "metric-out minimum-igp "))
					if err != nil {
						return fmt.Errorf("failed to convert value from '%s' to integer : %w", item, err)
					}
				}
			}
		}
	}
	if item == "mtu-discovery" {
		confRead.mtuDiscovery = true
	}
	if item == "multihop" {
		confRead.multihop = true
	}
	if item == "multipath" {
		confRead.multipath = true
	}
	if item == "no-advertise-peer-as" {
		confRead.noAdvertisePeerAs = true
	}
	if strings.HasPrefix(item, "out-delay ") {
		confRead.outDelay, err = strconv.Atoi(strings.TrimPrefix(item, "out-delay "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", item, err)
		}
	}
	if item == passiveW {
		confRead.passive = true
	}
	if strings.HasPrefix(item, "peer-as ") {
		confRead.peerAs = strings.TrimPrefix(item, "peer-as ")
	}
	if strings.HasPrefix(item, "preference ") {
		confRead.preference, err = strconv.Atoi(strings.TrimPrefix(item, "preference "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", item, err)
		}
	}
	if item == "remove-private" {
		confRead.removePrivate = true
	}
	if strings.HasPrefix(item, "type ") {
		confRead.bgpType = strings.TrimPrefix(item, "type ")
	}

	return nil
}

func setBgpOptsBfd(setPrefix string, bfdLivenessDetection []interface{},
	m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefixBfd := setPrefix + "bfd-liveness-detection "
	for _, v := range bfdLivenessDetection {
		if v != nil {
			m := v.(map[string]interface{})
			if m["authentication_algorithm"].(string) != "" {
				configSet = append(configSet, setPrefixBfd+"authentication algorithm "+m["authentication_algorithm"].(string))
			}
			if m["authentication_key_chain"].(string) != "" {
				configSet = append(configSet, setPrefixBfd+"authentication key-chain "+m["authentication_key_chain"].(string))
			}
			if m["authentication_loose_check"].(bool) {
				configSet = append(configSet, setPrefixBfd+"authentication loose-check")
			}
			if m["detection_time_threshold"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"detection-time threshold "+
					strconv.Itoa(m["detection_time_threshold"].(int)))
			}
			if m["holddown_interval"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"holddown-interval "+
					strconv.Itoa(m["holddown_interval"].(int)))
			}
			if m["minimum_interval"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"minimum-interval "+
					strconv.Itoa(m["minimum_interval"].(int)))
			}
			if m["minimum_receive_interval"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"minimum-receive-interval "+
					strconv.Itoa(m["minimum_receive_interval"].(int)))
			}
			if m["multiplier"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"multiplier "+
					strconv.Itoa(m["multiplier"].(int)))
			}
			if m["session_mode"].(string) != "" {
				configSet = append(configSet, setPrefixBfd+"session-mode "+m["session_mode"].(string))
			}
			if m["transmit_interval_minimum_interval"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"transmit-interval minimum-interval "+
					strconv.Itoa(m["transmit_interval_minimum_interval"].(int)))
			}
			if m["transmit_interval_threshold"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"transmit-interval threshold "+
					strconv.Itoa(m["transmit_interval_threshold"].(int)))
			}
			if m["version"].(string) != "" {
				configSet = append(configSet, setPrefixBfd+"version "+m["version"].(string))
			}
		}
	}
	if len(configSet) > 0 {
		err := sess.configSet(configSet, jnprSess)
		if err != nil {
			return err
		}
	}

	return nil
}
func readBgpOptsBfd(item string, bfdOpts []map[string]interface{}) ([]map[string]interface{}, error) {
	itemTrim := strings.TrimPrefix(item, "bfd-liveness-detection ")
	bfdRead := map[string]interface{}{
		"authentication_algorithm":           "",
		"authentication_key_chain":           "",
		"authentication_loose_check":         false,
		"detection_time_threshold":           0,
		"holddown_interval":                  0,
		"minimum_interval":                   0,
		"minimum_receive_interval":           0,
		"multiplier":                         0,
		"session_mode":                       "",
		"transmit_interval_minimum_interval": 0,
		"transmit_interval_threshold":        0,
		"version":                            "",
	}
	if len(bfdOpts) > 0 {
		for k, v := range bfdOpts[0] {
			bfdRead[k] = v
		}
	}
	var err error
	if strings.HasPrefix(itemTrim, "authentication algorithm ") {
		bfdRead["authentication_algorithm"] = strings.TrimPrefix(itemTrim, "authentication algorithm ")
	}
	if strings.HasPrefix(itemTrim, "authentication key-chain ") {
		bfdRead["authentication_key_chain"] = strings.TrimPrefix(itemTrim, "authentication key-chain ")
	}
	if itemTrim == "authentication loose-check" {
		bfdRead["authentication_loose_check"] = true
	}
	if strings.HasPrefix(itemTrim, "detection-time threshold ") {
		bfdRead["detection_time_threshold"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "detection-time threshold "))
		if err != nil {
			return []map[string]interface{}{bfdRead},
				fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}
	if strings.HasPrefix(itemTrim, "holddown-interval ") {
		bfdRead["holddown_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "holddown-interval "))
		if err != nil {
			return []map[string]interface{}{bfdRead},
				fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}
	if strings.HasPrefix(itemTrim, "minimum-interval ") {
		bfdRead["minimum_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "minimum-interval "))
		if err != nil {
			return []map[string]interface{}{bfdRead},
				fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}
	if strings.HasPrefix(itemTrim, "minimum-receive-interval ") {
		bfdRead["minimum_receive_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "minimum-receive-interval "))
		if err != nil {
			return []map[string]interface{}{bfdRead},
				fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}
	if strings.HasPrefix(itemTrim, "multiplier ") {
		bfdRead["multiplier"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "multiplier "))
		if err != nil {
			return []map[string]interface{}{bfdRead},
				fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}
	if strings.HasPrefix(itemTrim, "session-mode ") {
		bfdRead["session_mode"] = strings.TrimPrefix(itemTrim, "session-mode ")
	}
	if strings.HasPrefix(itemTrim, "transmit-interval threshold ") {
		bfdRead["transmit_interval_threshold"], err = strconv.Atoi(
			strings.TrimPrefix(itemTrim, "transmit-interval threshold "))
		if err != nil {
			return []map[string]interface{}{bfdRead},
				fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}
	if strings.HasPrefix(itemTrim, "transmit-interval minimum-interval ") {
		bfdRead["transmit_interval_minimum_interval"], err = strconv.Atoi(
			strings.TrimPrefix(itemTrim, "transmit-interval minimum-interval "))
		if err != nil {
			return []map[string]interface{}{bfdRead},
				fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}
	if strings.HasPrefix(itemTrim, "version ") {
		bfdRead["version"] = strings.TrimPrefix(itemTrim, "version ")
	}

	// override (maxItem = 1)
	return []map[string]interface{}{bfdRead}, nil
}

func setBgpOptsFamily(setPrefix, familyType string, familyOptsList []interface{},
	m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefixFamily := setPrefix + "family "
	if familyType == evpnWord {
		setPrefixFamily += "evpn "
	} else if familyType == inetWord {
		setPrefixFamily += "inet "
	} else if familyType == inet6Word {
		setPrefixFamily += "inet6 "
	}
	for _, familyOpts := range familyOptsList {
		familyOptsM := familyOpts.(map[string]interface{})
		configSet = append(configSet, setPrefixFamily+familyOptsM["nlri_type"].(string))
		for _, v := range familyOptsM["accepted_prefix_limit"].([]interface{}) {
			m := v.(map[string]interface{})
			setP := setPrefixFamily + familyOptsM["nlri_type"].(string) + " accepted-prefix-limit "
			configSetBgpOptsFamilyPrefixLimit, err := setBgpOptsFamilyPrefixLimit(setP, m)
			if err != nil {
				return err
			}
			configSet = append(configSet, configSetBgpOptsFamilyPrefixLimit...)
		}
		for _, v := range familyOptsM["prefix_limit"].([]interface{}) {
			m := v.(map[string]interface{})
			setP := setPrefixFamily + familyOptsM["nlri_type"].(string) + " prefix-limit "
			configSetBgpOptsFamilyPrefixLimit, err := setBgpOptsFamilyPrefixLimit(setP, m)
			if err != nil {
				return err
			}
			configSet = append(configSet, configSetBgpOptsFamilyPrefixLimit...)
		}
	}
	if len(configSet) > 0 {
		err := sess.configSet(configSet, jnprSess)
		if err != nil {
			return err
		}
	}

	return nil
}
func setBgpOptsFamilyPrefixLimit(setPrefix string, prefixLimit map[string]interface{}) ([]string, error) {
	configSet := make([]string, 0)
	if prefixLimit["maximum"].(int) != 0 {
		configSet = append(configSet, setPrefix+"maximum "+strconv.Itoa(prefixLimit["maximum"].(int)))
	}
	if prefixLimit["teardown"].(int) != 0 {
		configSet = append(configSet, setPrefix+"teardown "+strconv.Itoa(prefixLimit["teardown"].(int)))
	}
	if prefixLimit["teardown_idle_timeout"].(int) != 0 {
		if prefixLimit["teardown_idle_timeout_forever"].(bool) {
			return configSet, fmt.Errorf("conflict between teardown_idle_timeout and teardown_idle_timeout_forever")
		}
		configSet = append(configSet, setPrefix+"teardown idle-timeout "+
			strconv.Itoa(prefixLimit["teardown_idle_timeout"].(int)))
	}
	if prefixLimit["teardown_idle_timeout_forever"].(bool) {
		configSet = append(configSet, setPrefix+"teardown idle-timeout forever")
	}

	return configSet, nil
}
func readBgpOptsFamily(item, familyType string, opts []map[string]interface{}) ([]map[string]interface{}, error) {
	readOpts := map[string]interface{}{
		"nlri_type":             "",
		"accepted_prefix_limit": make([]map[string]interface{}, 0, 1),
		"prefix_limit":          make([]map[string]interface{}, 0, 1),
	}
	setPrefix := "family "
	if familyType == inetWord {
		setPrefix += "inet "
	}
	if familyType == inet6Word {
		setPrefix += "inet6 "
	}
	trimSplit := strings.Split(strings.TrimPrefix(item, setPrefix), " ")
	readOpts["nlri_type"] = trimSplit[0]
	readOpts, opts = copyAndRemoveItemMapList("nlri_type", false, readOpts, opts)

	var err error
	itemTrim := strings.TrimPrefix(item, setPrefix+readOpts["nlri_type"].(string)+" ")
	if strings.HasPrefix(itemTrim, "accepted-prefix-limit ") {
		readOptsPL := map[string]interface{}{
			"maximum":                       0,
			"teardown":                      0,
			"teardown_idle_timeout":         0,
			"teardown_idle_timeout_forever": false,
		}
		if len(readOpts["accepted_prefix_limit"].([]map[string]interface{})) > 0 {
			for k, v := range readOpts["accepted_prefix_limit"].([]map[string]interface{})[0] {
				readOptsPL[k] = v
			}
		}
		if strings.HasPrefix(itemTrim, "accepted-prefix-limit maximum") {
			readOptsPL["maximum"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "accepted-prefix-limit maximum "))
			if err != nil {
				return append(opts, readOpts), fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
		if strings.HasPrefix(itemTrim, "accepted-prefix-limit teardown ") && !strings.Contains(itemTrim, " idle-timeout ") {
			readOptsPL["teardown"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "accepted-prefix-limit teardown "))
			if err != nil {
				return append(opts, readOpts), fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
		if strings.HasPrefix(itemTrim, "accepted-prefix-limit teardown idle-timeout ") {
			if !strings.HasSuffix(itemTrim, " forever") {
				readOptsPL["teardown_idle_timeout"], err = strconv.Atoi(
					strings.TrimPrefix(itemTrim, "accepted-prefix-limit teardown idle-timeout "))
				if err != nil {
					return append(opts, readOpts), fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			} else {
				readOptsPL["teardown_idle_timeout_forever"] = true
			}
		}
		// override (maxItem = 1)
		readOpts["accepted_prefix_limit"] = []map[string]interface{}{readOptsPL}
	}
	if strings.HasPrefix(itemTrim, "prefix-limit ") {
		readOptsPL := map[string]interface{}{
			"maximum":                       0,
			"teardown":                      0,
			"teardown_idle_timeout":         0,
			"teardown_idle_timeout_forever": false,
		}
		if len(readOpts["prefix_limit"].([]map[string]interface{})) > 0 {
			for k, v := range readOpts["prefix_limit"].([]map[string]interface{})[0] {
				readOptsPL[k] = v
			}
		}

		if strings.HasPrefix(itemTrim, "prefix-limit maximum ") {
			readOptsPL["maximum"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "prefix-limit maximum "))
			if err != nil {
				return append(opts, readOpts), fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
		if strings.HasPrefix(itemTrim, "prefix-limit teardown ") && !strings.Contains(itemTrim, " idle-timeout ") {
			readOptsPL["teardown"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "prefix-limit teardown "))
			if err != nil {
				return append(opts, readOpts), fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
		if strings.HasPrefix(itemTrim, "prefix-limit teardown idle-timeout ") {
			if !strings.HasSuffix(itemTrim, " forever") {
				readOptsPL["teardown_idle_timeout"], err = strconv.Atoi(
					strings.TrimPrefix(itemTrim, "prefix-limit teardown idle-timeout "))
				if err != nil {
					return append(opts, readOpts), fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			} else {
				readOptsPL["teardown_idle_timeout_forever"] = true
			}
		}
		// override (maxItem = 1)
		readOpts["prefix_limit"] = []map[string]interface{}{readOptsPL}
	}

	return append(opts, readOpts), nil
}
func setBgpOptsGrafefulRestart(setPrefix string, gracefulRestarts []interface{},
	m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	for _, v := range gracefulRestarts {
		if v != nil {
			m := v.(map[string]interface{})
			if m["disable"].(bool) {
				configSet = append(configSet, setPrefix+"graceful-restart disable")
			}
			if m["restart_time"].(int) != 0 {
				configSet = append(configSet, setPrefix+"graceful-restart restart-time "+
					strconv.Itoa(m["restart_time"].(int)))
			}
			if m["stale_route_time"].(int) != 0 {
				configSet = append(configSet, setPrefix+"graceful-restart stale-routes-time "+
					strconv.Itoa(m["stale_route_time"].(int)))
			}
		}
	}
	if len(configSet) > 0 {
		err := sess.configSet(configSet, jnprSess)
		if err != nil {
			return err
		}
	}

	return nil
}
func readBgpOptsGracefulRestart(item string, grOpts []map[string]interface{}) ([]map[string]interface{}, error) {
	itemTrim := strings.TrimPrefix(item, "graceful-restart ")
	grRead := map[string]interface{}{
		"disable":          false,
		"restart_time":     0,
		"stale_route_time": 0,
	}
	if len(grOpts) > 0 {
		for k, v := range grOpts[0] {
			grRead[k] = v
		}
	}
	var err error
	if itemTrim == disableW {
		grRead["disable"] = true
	}
	if strings.HasPrefix(itemTrim, "restart-time ") {
		grRead["restart_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "restart-time "))
		if err != nil {
			return []map[string]interface{}{grRead},
				fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}
	if strings.HasPrefix(itemTrim, "stale-routes-time ") {
		grRead["stale_route_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "stale-routes-time "))
		if err != nil {
			return []map[string]interface{}{grRead},
				fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}
	// override (maxItem = 1)
	return []map[string]interface{}{grRead}, nil
}
