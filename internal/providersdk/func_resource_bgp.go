package providersdk

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	jdecode "github.com/jeremmfr/junosdecode"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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

func delBgpOpts(d *schema.ResourceData, typebgp string, junSess *junos.Session,
) error {
	configSet := make([]string, 0)
	delPrefix := junos.DeleteLS
	if d.Get("routing_instance").(string) != junos.DefaultW {
		delPrefix = junos.DelRoutingInstances + d.Get("routing_instance").(string) + " "
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

	return junSess.ConfigSet(configSet)
}

func setBgpOptsSimple(setPrefix string, d *schema.ResourceData, junSess *junos.Session) error {
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

	return junSess.ConfigSet(configSet)
}

func (confRead *bgpOptions) readBgpOptsSimple(itemTrim string) (err error) {
	switch {
	case itemTrim == "accept-remote-nexthop":
		confRead.acceptRemoteNexthop = true
	case itemTrim == "advertise-external":
		confRead.advertiseExternal = true
	case itemTrim == "advertise-external conditional":
		confRead.advertiseExternalConditional = true
	case itemTrim == "advertise-inactive":
		confRead.advertiseInactive = true
	case itemTrim == "advertise-peer-as":
		confRead.advertisePeerAs = true
	case itemTrim == "as-override":
		confRead.asOverride = true
	case balt.CutPrefixInString(&itemTrim, "authentication-algorithm "):
		confRead.authenticationAlgorithm = itemTrim
	case balt.CutPrefixInString(&itemTrim, "authentication-key "):
		confRead.authenticationKey, err = jdecode.Decode(strings.Trim(itemTrim, "\""))
		if err != nil {
			return fmt.Errorf("failed to decode authentication-key: %w", err)
		}
	case balt.CutPrefixInString(&itemTrim, "authentication-key-chain "):
		confRead.authenticationKeyChain = itemTrim
	case balt.CutPrefixInString(&itemTrim, "cluster "):
		confRead.cluster = itemTrim
	case itemTrim == "damping":
		confRead.damping = true
	case balt.CutPrefixInString(&itemTrim, "export "):
		confRead.exportPolicy = append(confRead.exportPolicy, itemTrim)
	case balt.CutPrefixInString(&itemTrim, "hold-time "):
		confRead.holdTime, err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "import "):
		confRead.importPolicy = append(confRead.importPolicy, itemTrim)
	case itemTrim == "keep all":
		confRead.keepAll = true
	case itemTrim == "keep none":
		confRead.keepNone = true
	case balt.CutPrefixInString(&itemTrim, "local-address "):
		confRead.localAddress = itemTrim
	case balt.CutPrefixInString(&itemTrim, "local-as "):
		switch {
		case itemTrim == "private":
			confRead.localAsPrivate = true
		case itemTrim == "alias":
			confRead.localAsAlias = true
		case itemTrim == "no-prepend-global-as":
			confRead.localAsNoPrependGlobalAs = true
		case balt.CutPrefixInString(&itemTrim, "loops "):
			confRead.localAsLoops, err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		default:
			confRead.localAs = itemTrim
		}
	case balt.CutPrefixInString(&itemTrim, "local-interface "):
		confRead.localInterface = itemTrim
	case balt.CutPrefixInString(&itemTrim, "local-preference "):
		confRead.localPreference, err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "log-updown":
		confRead.logUpdown = true
	case balt.CutPrefixInString(&itemTrim, "metric-out "):
		switch {
		case balt.CutPrefixInString(&itemTrim, "igp"):
			confRead.metricOutIgp = true
			switch {
			case itemTrim == " delay-med-update":
				confRead.metricOutIgpDelayMedUpdate = true
			case balt.CutPrefixInString(&itemTrim, " "):
				confRead.metricOutIgpOffset, err = strconv.Atoi(itemTrim)
				if err != nil {
					return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-igp"):
			confRead.metricOutMinimumIgp = true
			if balt.CutPrefixInString(&itemTrim, " ") {
				confRead.metricOutMinimumIgpOffset, err = strconv.Atoi(itemTrim)
				if err != nil {
					return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			}
		default:
			confRead.metricOut, err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		}
	case itemTrim == "mtu-discovery":
		confRead.mtuDiscovery = true
	case itemTrim == "multihop":
		confRead.multihop = true
	case balt.CutPrefixInString(&itemTrim, "multipath"):
		if len(confRead.bgpMultipath) == 0 {
			confRead.bgpMultipath = append(confRead.bgpMultipath, map[string]interface{}{
				"allow_protection": false,
				"disable":          false,
				"multiple_as":      false,
			})
		}
		switch {
		case itemTrim == " allow-protection":
			confRead.bgpMultipath[0]["allow_protection"] = true
		case itemTrim == " disable":
			confRead.bgpMultipath[0]["disable"] = true
		case itemTrim == " multiple-as":
			confRead.bgpMultipath[0]["multiple_as"] = true
		}
	case itemTrim == "no-advertise-peer-as":
		confRead.noAdvertisePeerAs = true
	case balt.CutPrefixInString(&itemTrim, "out-delay "):
		confRead.outDelay, err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "passive":
		confRead.passive = true
	case balt.CutPrefixInString(&itemTrim, "peer-as "):
		confRead.peerAs = itemTrim
	case balt.CutPrefixInString(&itemTrim, "preference "):
		confRead.preference, err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "remove-private":
		confRead.removePrivate = true
	case balt.CutPrefixInString(&itemTrim, "type "):
		confRead.bgpType = itemTrim
	}

	return nil
}

func setBgpOptsBfd(setPrefix string, bfdLivenessDetection []interface{}, junSess *junos.Session,
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
		err := junSess.ConfigSet(configSet)
		if err != nil {
			return err
		}
	}

	return nil
}

func readBgpOptsBfd(itemTrim string, bfdRead map[string]interface{}) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "authentication algorithm "):
		bfdRead["authentication_algorithm"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "authentication key-chain "):
		bfdRead["authentication_key_chain"] = itemTrim
	case itemTrim == "authentication loose-check":
		bfdRead["authentication_loose_check"] = true
	case balt.CutPrefixInString(&itemTrim, "detection-time threshold "):
		bfdRead["detection_time_threshold"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "holddown-interval "):
		bfdRead["holddown_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "minimum-interval "):
		bfdRead["minimum_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "minimum-receive-interval "):
		bfdRead["minimum_receive_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "multiplier "):
		bfdRead["multiplier"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "session-mode "):
		bfdRead["session_mode"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "transmit-interval threshold "):
		bfdRead["transmit_interval_threshold"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "transmit-interval minimum-interval "):
		bfdRead["transmit_interval_minimum_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "version "):
		bfdRead["version"] = itemTrim
	}

	return nil
}

func setBgpOptsFamily(
	setPrefix, familyType string, familyOptsList []interface{}, junSess *junos.Session,
) error {
	configSet := make([]string, 0)
	setPrefixFamily := setPrefix + "family "
	switch familyType {
	case junos.EvpnW:
		setPrefixFamily += "evpn "
	case junos.InetW:
		setPrefixFamily += "inet "
	case junos.Inet6W:
		setPrefixFamily += "inet6 "
	}
	familyNlriTypeList := make([]string, 0)
	for _, familyOpts := range familyOptsList {
		familyOptsM := familyOpts.(map[string]interface{})
		if bchk.InSlice(familyOptsM["nlri_type"].(string), familyNlriTypeList) {
			switch familyType {
			case junos.EvpnW:
				return fmt.Errorf("multiple blocks family_evpn with the same nlri_type %s", familyOptsM["nlri_type"].(string))
			case junos.InetW:
				return fmt.Errorf("multiple blocks family_inet with the same nlri_type %s", familyOptsM["nlri_type"].(string))
			case junos.Inet6W:
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
		err := junSess.ConfigSet(configSet)
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

func readBgpOptsFamily(itemTrim string, opts []map[string]interface{}) (_ []map[string]interface{}, err error) {
	itemTrimFields := strings.Split(itemTrim, " ")
	readOpts := map[string]interface{}{
		"nlri_type":             itemTrimFields[0],
		"accepted_prefix_limit": make([]map[string]interface{}, 0),
		"prefix_limit":          make([]map[string]interface{}, 0),
	}
	opts = copyAndRemoveItemMapList("nlri_type", readOpts, opts)
	balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
	switch {
	case balt.CutPrefixInString(&itemTrim, "accepted-prefix-limit "):
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
		case balt.CutPrefixInString(&itemTrim, "maximum "):
			readOptsPL["maximum"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return append(opts, readOpts), fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}

		case balt.CutPrefixInString(&itemTrim, "teardown idle-timeout "):
			if itemTrim == "forever" {
				readOptsPL["teardown_idle_timeout_forever"] = true
			} else {
				readOptsPL["teardown_idle_timeout"], err = strconv.Atoi(itemTrim)
				if err != nil {
					return append(opts, readOpts), fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			}
		case balt.CutPrefixInString(&itemTrim, "teardown "):
			readOptsPL["teardown"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return append(opts, readOpts), fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		}
	case balt.CutPrefixInString(&itemTrim, "prefix-limit "):
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
		case balt.CutPrefixInString(&itemTrim, "maximum "):
			readOptsPL["maximum"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return append(opts, readOpts), fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "teardown idle-timeout "):
			if itemTrim == "forever" {
				readOptsPL["teardown_idle_timeout_forever"] = true
			} else {
				readOptsPL["teardown_idle_timeout"], err = strconv.Atoi(itemTrim)
				if err != nil {
					return append(opts, readOpts), fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			}
		case balt.CutPrefixInString(&itemTrim, "teardown "):
			readOptsPL["teardown"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return append(opts, readOpts), fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
			}
		}
	}

	return append(opts, readOpts), nil
}

func setBgpOptsGrafefulRestart(
	setPrefix string, gracefulRestarts []interface{}, junSess *junos.Session,
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
		err := junSess.ConfigSet(configSet)
		if err != nil {
			return err
		}
	}

	return nil
}

func readBgpOptsGracefulRestart(itemTrim string, grRead map[string]interface{}) (err error) {
	switch {
	case itemTrim == junos.DisableW:
		grRead["disable"] = true
	case balt.CutPrefixInString(&itemTrim, "restart-time "):
		grRead["restart_time"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "stale-routes-time "):
		grRead["stale_route_time"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	}

	return nil
}
