package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
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
	keepAll                      bool
	keepNone                     bool
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
	cluster                      string
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
	bgpMultipath                 []map[string]interface{}
	familyEvpn                   []map[string]interface{}
	familyInet                   []map[string]interface{}
	familyInet6                  []map[string]interface{}
	gracefulRestart              []map[string]interface{}
}

func delBgpOpts(d *schema.ResourceData, typebgp string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)
	delPrefix := deleteLS
	if d.Get("routing_instance").(string) != defaultW {
		delPrefix = delRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	switch typebgp {
	case "group":
		delPrefix += "protocols bgp group " + d.Get("name").(string) + " "
	case "neighbor":
		delPrefix += "protocols bgp group " + d.Get("group").(string) +
			" neighbor " + d.Get("ip").(string) + " "
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
		delPrefix+"cluster",
		delPrefix+"damping",
		delPrefix+"export",
		delPrefix+"family evpn",
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

	return clt.configSet(configSet, junSess)
}

func setBgpOptsSimple(setPrefix string, d *schema.ResourceData, clt *Client, junSess *junosSession) error {
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
	for _, v := range d.Get("bgp_multipath").([]interface{}) {
		configSet = append(configSet, setPrefix+"multipath")
		if v != nil {
			bgpMultipah := v.(map[string]interface{})
			if bgpMultipah["allow_protection"].(bool) {
				configSet = append(configSet, setPrefix+"multipath allow-protection")
			}
			if bgpMultipah["disable"].(bool) {
				configSet = append(configSet, setPrefix+"multipath disable")
			}
			if bgpMultipah["multiple_as"].(bool) {
				configSet = append(configSet, setPrefix+"multipath multiple-as")
			}
		}
	}
	if v := d.Get("cluster").(string); v != "" {
		configSet = append(configSet, setPrefix+"cluster "+v)
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
	if d.Get("keep_all").(bool) {
		configSet = append(configSet, setPrefix+"keep all")
	}
	if d.Get("keep_none").(bool) {
		configSet = append(configSet, setPrefix+"keep none")
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

	return clt.configSet(configSet, junSess)
}

func readBgpOptsSimple(item string, confRead *bgpOptions) error {
	switch {
	case item == "accept-remote-nexthop":
		confRead.acceptRemoteNexthop = true
	case item == "advertise-external":
		confRead.advertiseExternal = true
	case item == "advertise-external conditional":
		confRead.advertiseExternalConditional = true
	case item == "advertise-inactive":
		confRead.advertiseInactive = true
	case item == "advertise-peer-as":
		confRead.advertisePeerAs = true
	case item == "as-override":
		confRead.asOverride = true
	case strings.HasPrefix(item, "authentication-algorithm "):
		confRead.authenticationAlgorithm = strings.TrimPrefix(item, "authentication-algorithm ")
	case strings.HasPrefix(item, "authentication-key "):
		var err error
		confRead.authenticationKey, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(item, "authentication-key "), "\""))
		if err != nil {
			return fmt.Errorf("failed to decode authentication-key: %w", err)
		}
	case strings.HasPrefix(item, "authentication-key-chain "):
		confRead.authenticationKeyChain = strings.TrimPrefix(item, "authentication-key-chain ")
	case strings.HasPrefix(item, "cluster "):
		confRead.cluster = strings.TrimPrefix(item, "cluster ")
	case item == "damping":
		confRead.damping = true
	case strings.HasPrefix(item, "export "):
		confRead.exportPolicy = append(confRead.exportPolicy, strings.TrimPrefix(item, "export "))
	case strings.HasPrefix(item, "hold-time "):
		var err error
		confRead.holdTime, err = strconv.Atoi(strings.TrimPrefix(item, "hold-time "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, item, err)
		}
	case strings.HasPrefix(item, "import "):
		confRead.importPolicy = append(confRead.importPolicy, strings.TrimPrefix(item, "import "))
	case item == "keep all":
		confRead.keepAll = true
	case item == "keep none":
		confRead.keepNone = true
	case strings.HasPrefix(item, "local-address "):
		confRead.localAddress = strings.TrimPrefix(item, "local-address ")
	case strings.HasPrefix(item, "local-as "):
		switch {
		case strings.HasSuffix(item, " private"):
			confRead.localAsPrivate = true
		case strings.HasSuffix(item, " alias"):
			confRead.localAsAlias = true
		case strings.HasSuffix(item, " no-prepend-global-as"):
			confRead.localAsNoPrependGlobalAs = true
		case strings.HasPrefix(item, "local-as loops "):
			var err error
			confRead.localAsLoops, err = strconv.Atoi(strings.TrimPrefix(item, "local-as loops "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, item, err)
			}
		default:
			confRead.localAs = strings.TrimPrefix(item, "local-as ")
		}
	case strings.HasPrefix(item, "local-interface "):
		confRead.localInterface = strings.TrimPrefix(item, "local-interface ")
	case strings.HasPrefix(item, "local-preference "):
		var err error
		confRead.localPreference, err = strconv.Atoi(strings.TrimPrefix(item, "local-preference "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, item, err)
		}
	case item == "log-updown":
		confRead.logUpdown = true
	case strings.HasPrefix(item, "metric-out "):
		if !strings.Contains(item, "igp") {
			var err error
			confRead.metricOut, err = strconv.Atoi(strings.TrimPrefix(item, "metric-out "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, item, err)
			}
		} else {
			if strings.HasPrefix(item, "metric-out igp") {
				confRead.metricOutIgp = true
				if item == "metric-out igp delay-med-update" {
					confRead.metricOutIgpDelayMedUpdate = true
				} else if strings.HasPrefix(item, "metric-out igp ") {
					var err error
					confRead.metricOutIgpOffset, err = strconv.Atoi(strings.TrimPrefix(item, "metric-out igp "))
					if err != nil {
						return fmt.Errorf(failedConvAtoiError, item, err)
					}
				}
			} else {
				confRead.metricOutMinimumIgp = true
				if strings.HasPrefix(item, "metric-out minimum-igp ") {
					var err error
					confRead.metricOutMinimumIgpOffset, err = strconv.Atoi(strings.TrimPrefix(item, "metric-out minimum-igp "))
					if err != nil {
						return fmt.Errorf(failedConvAtoiError, item, err)
					}
				}
			}
		}
	case item == "mtu-discovery":
		confRead.mtuDiscovery = true
	case item == "multihop":
		confRead.multihop = true
	case strings.HasPrefix(item, "multipath"):
		confRead.multipath = true
		if len(confRead.bgpMultipath) == 0 {
			confRead.bgpMultipath = append(confRead.bgpMultipath, map[string]interface{}{
				"allow_protection": false,
				"disable":          false,
				"multiple_as":      false,
			})
		}
		switch {
		case item == "multipath allow-protection":
			confRead.bgpMultipath[0]["allow_protection"] = true
		case item == "multipath disable":
			confRead.bgpMultipath[0]["disable"] = true
		case item == "multipath multiple-as":
			confRead.bgpMultipath[0]["multiple_as"] = true
		}
	case item == "no-advertise-peer-as":
		confRead.noAdvertisePeerAs = true
	case strings.HasPrefix(item, "out-delay "):
		var err error
		confRead.outDelay, err = strconv.Atoi(strings.TrimPrefix(item, "out-delay "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, item, err)
		}
	case item == "passive":
		confRead.passive = true
	case strings.HasPrefix(item, "peer-as "):
		confRead.peerAs = strings.TrimPrefix(item, "peer-as ")
	case strings.HasPrefix(item, "preference "):
		var err error
		confRead.preference, err = strconv.Atoi(strings.TrimPrefix(item, "preference "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, item, err)
		}
	case item == "remove-private":
		confRead.removePrivate = true
	case strings.HasPrefix(item, "type "):
		confRead.bgpType = strings.TrimPrefix(item, "type ")
	}

	return nil
}

func setBgpOptsBfd(setPrefix string, bfdLivenessDetection []interface{}, clt *Client, junSess *junosSession,
) error {
	configSet := make([]string, 0)

	setPrefixBfd := setPrefix + "bfd-liveness-detection "
	for _, v := range bfdLivenessDetection {
		if v != nil {
			bfdLD := v.(map[string]interface{})
			if bfdLD["authentication_algorithm"].(string) != "" {
				configSet = append(configSet, setPrefixBfd+"authentication algorithm "+bfdLD["authentication_algorithm"].(string))
			}
			if bfdLD["authentication_key_chain"].(string) != "" {
				configSet = append(configSet, setPrefixBfd+"authentication key-chain "+bfdLD["authentication_key_chain"].(string))
			}
			if bfdLD["authentication_loose_check"].(bool) {
				configSet = append(configSet, setPrefixBfd+"authentication loose-check")
			}
			if bfdLD["detection_time_threshold"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"detection-time threshold "+
					strconv.Itoa(bfdLD["detection_time_threshold"].(int)))
			}
			if bfdLD["holddown_interval"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"holddown-interval "+
					strconv.Itoa(bfdLD["holddown_interval"].(int)))
			}
			if bfdLD["minimum_interval"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"minimum-interval "+
					strconv.Itoa(bfdLD["minimum_interval"].(int)))
			}
			if bfdLD["minimum_receive_interval"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"minimum-receive-interval "+
					strconv.Itoa(bfdLD["minimum_receive_interval"].(int)))
			}
			if bfdLD["multiplier"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"multiplier "+
					strconv.Itoa(bfdLD["multiplier"].(int)))
			}
			if bfdLD["session_mode"].(string) != "" {
				configSet = append(configSet, setPrefixBfd+"session-mode "+bfdLD["session_mode"].(string))
			}
			if bfdLD["transmit_interval_minimum_interval"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"transmit-interval minimum-interval "+
					strconv.Itoa(bfdLD["transmit_interval_minimum_interval"].(int)))
			}
			if bfdLD["transmit_interval_threshold"].(int) != 0 {
				configSet = append(configSet, setPrefixBfd+"transmit-interval threshold "+
					strconv.Itoa(bfdLD["transmit_interval_threshold"].(int)))
			}
			if bfdLD["version"].(string) != "" {
				configSet = append(configSet, setPrefixBfd+"version "+bfdLD["version"].(string))
			}
		}
	}
	if len(configSet) > 0 {
		err := clt.configSet(configSet, junSess)
		if err != nil {
			return err
		}
	}

	return nil
}

func readBgpOptsBfd(item string, bfdRead map[string]interface{}) error {
	itemTrim := strings.TrimPrefix(item, "bfd-liveness-detection ")
	switch {
	case strings.HasPrefix(itemTrim, "authentication algorithm "):
		bfdRead["authentication_algorithm"] = strings.TrimPrefix(itemTrim, "authentication algorithm ")
	case strings.HasPrefix(itemTrim, "authentication key-chain "):
		bfdRead["authentication_key_chain"] = strings.TrimPrefix(itemTrim, "authentication key-chain ")
	case itemTrim == "authentication loose-check":
		bfdRead["authentication_loose_check"] = true
	case strings.HasPrefix(itemTrim, "detection-time threshold "):
		var err error
		bfdRead["detection_time_threshold"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "detection-time threshold "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "holddown-interval "):
		var err error
		bfdRead["holddown_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "holddown-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "minimum-interval "):
		var err error
		bfdRead["minimum_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "minimum-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "minimum-receive-interval "):
		var err error
		bfdRead["minimum_receive_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "minimum-receive-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "multiplier "):
		var err error
		bfdRead["multiplier"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "multiplier "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "session-mode "):
		bfdRead["session_mode"] = strings.TrimPrefix(itemTrim, "session-mode ")
	case strings.HasPrefix(itemTrim, "transmit-interval threshold "):
		var err error
		bfdRead["transmit_interval_threshold"], err = strconv.Atoi(
			strings.TrimPrefix(itemTrim, "transmit-interval threshold "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "transmit-interval minimum-interval "):
		var err error
		bfdRead["transmit_interval_minimum_interval"], err = strconv.Atoi(
			strings.TrimPrefix(itemTrim, "transmit-interval minimum-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "version "):
		bfdRead["version"] = strings.TrimPrefix(itemTrim, "version ")
	}

	return nil
}

func setBgpOptsFamily(
	setPrefix, familyType string, familyOptsList []interface{}, clt *Client, junSess *junosSession,
) error {
	configSet := make([]string, 0)
	setPrefixFamily := setPrefix + "family "
	switch familyType {
	case evpnW:
		setPrefixFamily += "evpn "
	case inetW:
		setPrefixFamily += "inet "
	case inet6W:
		setPrefixFamily += "inet6 "
	}
	familyNlriTypeList := make([]string, 0)
	for _, familyOpts := range familyOptsList {
		familyOptsM := familyOpts.(map[string]interface{})
		if bchk.StringInSlice(familyOptsM["nlri_type"].(string), familyNlriTypeList) {
			switch familyType {
			case evpnW:
				return fmt.Errorf("multiple blocks family_evpn with the same nlri_type %s", familyOptsM["nlri_type"].(string))
			case inetW:
				return fmt.Errorf("multiple blocks family_inet with the same nlri_type %s", familyOptsM["nlri_type"].(string))
			case inet6W:
				return fmt.Errorf("multiple blocks family_inet6 with the same nlri_type %s", familyOptsM["nlri_type"].(string))
			}
		}
		familyNlriTypeList = append(familyNlriTypeList, familyOptsM["nlri_type"].(string))
		configSet = append(configSet, setPrefixFamily+familyOptsM["nlri_type"].(string))
		for _, v := range familyOptsM["accepted_prefix_limit"].([]interface{}) {
			mAccPrefixLimit := v.(map[string]interface{})
			setP := setPrefixFamily + familyOptsM["nlri_type"].(string) + " accepted-prefix-limit "
			configSetBgpOptsFamilyPrefixLimit, err := setBgpOptsFamilyPrefixLimit(setP, mAccPrefixLimit)
			if err != nil {
				return err
			}
			configSet = append(configSet, configSetBgpOptsFamilyPrefixLimit...)
		}
		for _, v := range familyOptsM["prefix_limit"].([]interface{}) {
			mPrefixLimit := v.(map[string]interface{})
			setP := setPrefixFamily + familyOptsM["nlri_type"].(string) + " prefix-limit "
			configSetBgpOptsFamilyPrefixLimit, err := setBgpOptsFamilyPrefixLimit(setP, mPrefixLimit)
			if err != nil {
				return err
			}
			configSet = append(configSet, configSetBgpOptsFamilyPrefixLimit...)
		}
	}
	if len(configSet) > 0 {
		err := clt.configSet(configSet, junSess)
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
		"accepted_prefix_limit": make([]map[string]interface{}, 0),
		"prefix_limit":          make([]map[string]interface{}, 0),
	}
	setPrefix := "family "
	switch familyType {
	case evpnW:
		setPrefix += "evpn "
	case inetW:
		setPrefix += "inet "
	case inet6W:
		setPrefix += "inet6 "
	}
	trimSplit := strings.Split(strings.TrimPrefix(item, setPrefix), " ")
	readOpts["nlri_type"] = trimSplit[0]
	opts = copyAndRemoveItemMapList("nlri_type", readOpts, opts)
	itemTrim := strings.TrimPrefix(item, setPrefix+readOpts["nlri_type"].(string)+" ")
	switch {
	case strings.HasPrefix(itemTrim, "accepted-prefix-limit "):
		if len(readOpts["accepted_prefix_limit"].([]map[string]interface{})) == 0 {
			readOpts["accepted_prefix_limit"] = append(readOpts["accepted_prefix_limit"].([]map[string]interface{}),
				map[string]interface{}{
					"maximum":                       0,
					"teardown":                      0,
					"teardown_idle_timeout":         0,
					"teardown_idle_timeout_forever": false,
				})
		}
		readOptsPL := readOpts["accepted_prefix_limit"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "accepted-prefix-limit maximum"):
			var err error
			readOptsPL["maximum"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "accepted-prefix-limit maximum "))
			if err != nil {
				return append(opts, readOpts), fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}

		case strings.HasPrefix(itemTrim, "accepted-prefix-limit teardown idle-timeout "):
			var err error
			if !strings.HasSuffix(itemTrim, " forever") {
				readOptsPL["teardown_idle_timeout"], err = strconv.Atoi(
					strings.TrimPrefix(itemTrim, "accepted-prefix-limit teardown idle-timeout "))
				if err != nil {
					return append(opts, readOpts), fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			} else {
				readOptsPL["teardown_idle_timeout_forever"] = true
			}
		case strings.HasPrefix(itemTrim, "accepted-prefix-limit teardown "):
			var err error
			readOptsPL["teardown"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "accepted-prefix-limit teardown "))
			if err != nil {
				return append(opts, readOpts), fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "prefix-limit "):
		if len(readOpts["prefix_limit"].([]map[string]interface{})) == 0 {
			readOpts["prefix_limit"] = append(readOpts["prefix_limit"].([]map[string]interface{}),
				map[string]interface{}{
					"maximum":                       0,
					"teardown":                      0,
					"teardown_idle_timeout":         0,
					"teardown_idle_timeout_forever": false,
				})
		}
		readOptsPL := readOpts["prefix_limit"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "prefix-limit maximum "):
			var err error
			readOptsPL["maximum"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "prefix-limit maximum "))
			if err != nil {
				return append(opts, readOpts), fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "prefix-limit teardown idle-timeout "):
			var err error
			if !strings.HasSuffix(itemTrim, " forever") {
				readOptsPL["teardown_idle_timeout"], err = strconv.Atoi(
					strings.TrimPrefix(itemTrim, "prefix-limit teardown idle-timeout "))
				if err != nil {
					return append(opts, readOpts), fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			} else {
				readOptsPL["teardown_idle_timeout_forever"] = true
			}
		case strings.HasPrefix(itemTrim, "prefix-limit teardown "):
			var err error
			readOptsPL["teardown"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "prefix-limit teardown "))
			if err != nil {
				return append(opts, readOpts), fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	}

	return append(opts, readOpts), nil
}

func setBgpOptsGrafefulRestart(setPrefix string, gracefulRestarts []interface{}, clt *Client, junSess *junosSession,
) error {
	configSet := make([]string, 0)

	for _, v := range gracefulRestarts {
		if v != nil {
			gRestart := v.(map[string]interface{})
			if gRestart["disable"].(bool) {
				configSet = append(configSet, setPrefix+"graceful-restart disable")
			}
			if gRestart["restart_time"].(int) != 0 {
				configSet = append(configSet, setPrefix+"graceful-restart restart-time "+
					strconv.Itoa(gRestart["restart_time"].(int)))
			}
			if gRestart["stale_route_time"].(int) != 0 {
				configSet = append(configSet, setPrefix+"graceful-restart stale-routes-time "+
					strconv.Itoa(gRestart["stale_route_time"].(int)))
			}
		}
	}
	if len(configSet) > 0 {
		err := clt.configSet(configSet, junSess)
		if err != nil {
			return err
		}
	}

	return nil
}

func readBgpOptsGracefulRestart(item string, grRead map[string]interface{}) error {
	itemTrim := strings.TrimPrefix(item, "graceful-restart ")
	switch {
	case itemTrim == disableW:
		grRead["disable"] = true
	case strings.HasPrefix(itemTrim, "restart-time "):
		var err error
		grRead["restart_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "restart-time "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "stale-routes-time "):
		var err error
		grRead["stale_route_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "stale-routes-time "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	}

	return nil
}
