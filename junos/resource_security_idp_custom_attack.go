package junos

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type idpCustomAttackOptions struct {
	timeBindingCount    int
	name                string
	recommendedAction   string
	severity            string
	timeBindingScope    string
	attackTypeAnomaly   []map[string]interface{}
	attackTypeChain     []map[string]interface{}
	attackTypeSignature []map[string]interface{}
}

func resourceSecurityIdpCustomAttack() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityIdpCustomAttackCreate,
		ReadContext:   resourceSecurityIdpCustomAttackRead,
		UpdateContext: resourceSecurityIdpCustomAttackUpdate,
		DeleteContext: resourceSecurityIdpCustomAttackDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityIdpCustomAttackImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"recommended_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"close",
					"close-client",
					"close-server",
					"drop",
					"drop-packet",
					"ignore",
					"none",
				}, false),
			},
			"severity": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"critical", "info", "major", "minor", "warning"}, false),
			},
			"attack_type_anomaly": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"attack_type_anomaly", "attack_type_chain", "attack_type_signature"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: schemaSecurityIdpCustomAttackTypeAnomaly(false),
				},
			},
			"attack_type_chain": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"attack_type_anomaly", "attack_type_chain", "attack_type_signature"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 2,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringDoesNotContainAny(" "),
									},
									"attack_type_anomaly": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: schemaSecurityIdpCustomAttackTypeAnomaly(true),
										},
									},
									"attack_type_signature": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: schemaSecurityIdpCustomAttackTypeSignature(true),
										},
									},
								},
							},
						},
						"expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"order": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"protocol_binding": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^(application|icmp|ip|rpc|tcp|udp)`),
								"must have valid protocol (application|icmp|ip|rpc|tcp|udp) with optional option"),
						},
						"reset": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"scope": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"session", "transaction"}, false),
						},
					},
				},
			},
			"attack_type_signature": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"attack_type_anomaly", "attack_type_chain", "attack_type_signature"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: schemaSecurityIdpCustomAttackTypeSignature(false),
				},
			},
			"time_binding_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(0, 4294967295),
			},
			"time_binding_scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"destination", "peer", "source"}, false),
			},
		},
	}
}

func schemaSecurityIdpCustomAttackTypeAnomaly(chain bool) map[string]*schema.Schema {
	r := map[string]*schema.Schema{
		"direction": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"any", "client-to-server", "server-to-client"}, false),
		},
		"service": {
			Type:     schema.TypeString,
			Required: true,
		},
		"test": {
			Type:     schema.TypeString,
			Required: true,
		},
		"shellcode": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"all", "intel", "no-shellcode", "sparc"}, false),
		},
	}
	if chain {
		delete(r, "service")
	}

	return r
}

func schemaSecurityIdpCustomAttackTypeSignature(chain bool) map[string]*schema.Schema {
	r := map[string]*schema.Schema{
		"context": {
			Type:     schema.TypeString,
			Required: true,
		},
		"direction": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"any", "client-to-server", "server-to-client"}, false),
		},
		"negate": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"pattern": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"pattern_pcre": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"protocol_icmp": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"checksum_validate_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"checksum_validate_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"code_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"code_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
					"data_length_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"data_length_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"identification_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"identification_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"sequence_number_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"sequence_number_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"type_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"type_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
				},
			},
		},
		"protocol_icmpv6": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"checksum_validate_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"checksum_validate_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"code_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"code_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
					"data_length_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"data_length_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
					"identification_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"identification_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"sequence_number_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"sequence_number_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"type_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"type_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
				},
			},
		},
		"protocol_ipv4": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"checksum_validate_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"checksum_validate_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"destination_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"destination_value": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsIPv4Address,
					},
					"identification_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"identification_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"ihl_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"ihl_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 15),
					},
					"ip_flags": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"protocol_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"protocol_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
					"source_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"source_value": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsIPv4Address,
					},
					"tos_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"tos_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
					"total_length_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"total_length_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"ttl_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"ttl_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
				},
			},
		},
		"protocol_ipv6": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"destination_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"destination_value": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsIPv6Address,
					},
					"extension_header_destination_option_home_address_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"extension_header_destination_option_home_address_value": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsIPv6Address,
					},
					"extension_header_destination_option_type_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"extension_header_destination_option_type_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
					"extension_header_routing_header_type_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"extension_header_routing_header_type_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
					"flow_label_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"flow_label_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 1048575),
					},
					"hop_limit_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"hop_limit_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
					"next_header_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"next_header_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
					"payload_length_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"payload_length_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"source_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"source_value": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsIPv6Address,
					},
					"traffic_class_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"traffic_class_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
				},
			},
		},
		"protocol_tcp": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ack_number_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"ack_number_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 4294967295),
					},
					"checksum_validate_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"checksum_validate_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"data_length_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"data_length_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(2, 255),
					},
					"destination_port_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"destination_port_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"header_length_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"header_length_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 15),
					},
					"mss_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"mss_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"option_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"option_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
					"reserved_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"reserved_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 7),
					},
					"sequence_number_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"sequence_number_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 4294967295),
					},
					"source_port_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"source_port_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"tcp_flags": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"urgent_pointer_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"urgent_pointer_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"window_scale_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"window_scale_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 255),
					},
					"window_size_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"window_size_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
				},
			},
		},
		"protocol_udp": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"checksum_validate_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"checksum_validate_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"data_length_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"data_length_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"destination_port_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"destination_port_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"source_port_match": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"equal", "greater-than", "less-than", "not-equal"}, false),
					},
					"source_port_value": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
				},
			},
		},
		"protocol_binding": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(
				`^(application|icmp|ip|rpc|tcp|udp)`),
				"must have valid protocol (application|icmp|ip|rpc|tcp|udp) with optional option"),
		},
		"regexp": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"shellcode": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"all", "intel", "no-shellcode", "sparc"}, false),
		},
	}
	if chain {
		delete(r, "protocol_binding")
	}

	return r
}

func resourceSecurityIdpCustomAttackCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSecurityIdpCustomAttack(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security idp custom-attack not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	idpCustomAttackExists, err := checkSecurityIdpCustomAttackExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if idpCustomAttackExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security idp custom-attack %v already exists", d.Get("name").(string)))...)
	}
	if err := setSecurityIdpCustomAttack(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_idp_custom_attack", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	idpCustomAttackExists, err = checkSecurityIdpCustomAttackExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if idpCustomAttackExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security idp custom-attack %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityIdpCustomAttackReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityIdpCustomAttackRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityIdpCustomAttackReadWJnprSess(d, m, jnprSess)
}

func resourceSecurityIdpCustomAttackReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	idpCustomAttackOptions, err := readSecurityIdpCustomAttack(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if idpCustomAttackOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityIdpCustomAttackData(d, idpCustomAttackOptions)
	}

	return nil
}

func resourceSecurityIdpCustomAttackUpdate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSecurityIdpCustomAttack(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityIdpCustomAttack(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_idp_custom_attack", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityIdpCustomAttackReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityIdpCustomAttackDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSecurityIdpCustomAttack(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_idp_custom_attack", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityIdpCustomAttackImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idpCustomAttackExists, err := checkSecurityIdpCustomAttackExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !idpCustomAttackExists {
		return nil, fmt.Errorf("don't find security idp custom-attack with id '%v' (id must be <name>)", d.Id())
	}
	idpCustomAttackOptions, err := readSecurityIdpCustomAttack(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityIdpCustomAttackData(d, idpCustomAttackOptions)

	result[0] = d

	return result, nil
}

func checkSecurityIdpCustomAttackExists(customAttack string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	customAttackConfig, err := sess.command("show configuration security idp custom-attack \""+
		customAttack+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if customAttackConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setSecurityIdpCustomAttack(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security idp custom-attack \"" + d.Get("name").(string) + "\" "
	configSet = append(configSet, setPrefix+"recommended-action "+d.Get("recommended_action").(string))
	configSet = append(configSet, setPrefix+"severity "+d.Get("severity").(string))
	for _, v := range d.Get("attack_type_anomaly").([]interface{}) {
		attackAnomaly := v.(map[string]interface{})
		configSet = append(configSet, setSecurityIdpCustomAttackTypeAnomaly(setPrefix, attackAnomaly)...)
	}
	for _, v := range d.Get("attack_type_chain").([]interface{}) {
		attackChain := v.(map[string]interface{})
		for _, v2 := range attackChain["member"].([]interface{}) {
			attackChainMember := v2.(map[string]interface{})
			if len(attackChainMember["attack_type_anomaly"].([]interface{})) != 0 &&
				len(attackChainMember["attack_type_signature"].([]interface{})) != 0 {
				return fmt.Errorf(
					"only one attack type is permitted in member %s for attack_type_chain", attackChainMember["name"].(string))
			} else if len(attackChainMember["attack_type_anomaly"].([]interface{})) == 0 &&
				len(attackChainMember["attack_type_signature"].([]interface{})) == 0 {
				return fmt.Errorf("missing one attack type in member %s for attack_type_chain", attackChainMember["name"].(string))
			}
			for _, v3 := range attackChainMember["attack_type_anomaly"].([]interface{}) {
				attackAnomaly := v3.(map[string]interface{})
				configSet = append(configSet, setSecurityIdpCustomAttackTypeAnomaly(
					setPrefix+"attack-type chain member \""+attackChainMember["name"].(string)+"\" ", attackAnomaly)...)
			}
			for _, v3 := range attackChainMember["attack_type_signature"].([]interface{}) {
				attackSignature := v3.(map[string]interface{})
				sets, err := setSecurityIdpCustomAttackTypeSignature(
					setPrefix+"attack-type chain member \""+attackChainMember["name"].(string)+"\" ", attackSignature)
				if err != nil {
					return err
				}
				configSet = append(configSet, sets...)
			}
		}
		if v2 := attackChain["expression"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"attack-type chain expression \""+v2+"\"")
		}
		if attackChain["order"].(bool) {
			configSet = append(configSet, setPrefix+"attack-type chain order")
		}
		if v2 := attackChain["protocol_binding"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"attack-type chain protocol-binding "+v2)
		}
		if attackChain["reset"].(bool) {
			configSet = append(configSet, setPrefix+"attack-type chain reset")
		}
		if v2 := attackChain["scope"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"attack-type chain scope "+v2)
		}
	}
	for _, v := range d.Get("attack_type_signature").([]interface{}) {
		attackSignature := v.(map[string]interface{})
		sets, err := setSecurityIdpCustomAttackTypeSignature(setPrefix, attackSignature)
		if err != nil {
			return err
		}
		configSet = append(configSet, sets...)
	}
	if v := d.Get("time_binding_count").(int); v != -1 {
		configSet = append(configSet, setPrefix+"time-binding count "+strconv.Itoa(v))
	}
	if v := d.Get("time_binding_scope").(string); v != "" {
		configSet = append(configSet, setPrefix+"time-binding scope "+v)
	}

	return sess.configSet(configSet, jnprSess)
}

func setSecurityIdpCustomAttackTypeAnomaly(
	setPrefixOrigin string, attackAnomaly map[string]interface{}) []string {
	configSet := make([]string, 0)

	setPrefix := setPrefixOrigin + "attack-type anomaly "
	configSet = append(configSet, setPrefix+"direction "+attackAnomaly["direction"].(string))
	if !strings.Contains(setPrefixOrigin, "attack-type chain member") {
		configSet = append(configSet, setPrefix+"service "+attackAnomaly["service"].(string))
	}
	configSet = append(configSet, setPrefix+"test "+attackAnomaly["test"].(string))
	if v := attackAnomaly["shellcode"].(string); v != "" {
		configSet = append(configSet, setPrefix+"shellcode "+v)
	}

	return configSet
}

func setSecurityIdpCustomAttackTypeSignature(
	setPrefixOrigin string, attackSignature map[string]interface{}) ([]string, error) {
	configSet := make([]string, 0)
	setPrefix := setPrefixOrigin + "attack-type signature "

	if v := attackSignature["context"].(string); v != "" {
		configSet = append(configSet, setPrefix+"context \""+v+"\"")
	}
	if v := attackSignature["direction"].(string); v != "" {
		configSet = append(configSet, setPrefix+"direction "+v)
	}
	if attackSignature["negate"].(bool) {
		configSet = append(configSet, setPrefix+"negate")
	}
	if v := attackSignature["pattern"].(string); v != "" {
		configSet = append(configSet, setPrefix+"pattern \""+v+"\"")
	}
	if v := attackSignature["pattern_pcre"].(string); v != "" {
		configSet = append(configSet, setPrefix+"pattern-pcre \""+v+"\"")
	}
	for _, v := range attackSignature["protocol_icmp"].([]interface{}) {
		if len(attackSignature["protocol_icmpv6"].([]interface{})) != 0 ||
			len(attackSignature["protocol_tcp"].([]interface{})) != 0 ||
			len(attackSignature["protocol_udp"].([]interface{})) != 0 {
			return configSet, fmt.Errorf("protocol_icmp cannot be specified with " +
				"protocol_icmpv6 or protocol_tcp or protocol_udp")
		}
		sets := setSecurityIdpCustomAttackTypeSignatureProtoICMP(false, setPrefix, v.(map[string]interface{}))
		if len(sets) == 0 {
			return configSet, fmt.Errorf("protocol_icmp block is empty")
		}
		configSet = append(configSet, sets...)
	}
	for _, v := range attackSignature["protocol_icmpv6"].([]interface{}) {
		if len(attackSignature["protocol_icmp"].([]interface{})) != 0 ||
			len(attackSignature["protocol_tcp"].([]interface{})) != 0 ||
			len(attackSignature["protocol_udp"].([]interface{})) != 0 {
			return configSet, fmt.Errorf("protocol_icmpv6 cannot be specified with " +
				"protocol_icmp or protocol_tcp or protocol_udp")
		}
		sets := setSecurityIdpCustomAttackTypeSignatureProtoICMP(true, setPrefix, v.(map[string]interface{}))
		if len(sets) == 0 {
			return configSet, fmt.Errorf("protocol_icmpv6 block is empty")
		}
		configSet = append(configSet, sets...)
	}
	for _, v := range attackSignature["protocol_ipv4"].([]interface{}) {
		sets := setSecurityIdpCustomAttackTypeSignatureProtoIPv4(setPrefix, v.(map[string]interface{}))
		if len(sets) == 0 {
			return configSet, fmt.Errorf("protocol_ipv4 block is empty")
		}
		configSet = append(configSet, sets...)
	}
	for _, v := range attackSignature["protocol_ipv6"].([]interface{}) {
		sets := setSecurityIdpCustomAttackTypeSignatureProtoIPv6(setPrefix, v.(map[string]interface{}))
		if len(sets) == 0 {
			return configSet, fmt.Errorf("protocol_ipv6 block is empty")
		}
		configSet = append(configSet, sets...)
	}
	for _, v := range attackSignature["protocol_tcp"].([]interface{}) {
		if len(attackSignature["protocol_icmp"].([]interface{})) != 0 ||
			len(attackSignature["protocol_icmpv6"].([]interface{})) != 0 ||
			len(attackSignature["protocol_udp"].([]interface{})) != 0 {
			return configSet, fmt.Errorf("protocol_tcp cannot be specified with " +
				"protocol_icmp or protocol_icmpv6 or protocol_udp")
		}
		sets := setSecurityIdpCustomAttackTypeSignatureProtoTCP(setPrefix, v.(map[string]interface{}))
		if len(sets) == 0 {
			return configSet, fmt.Errorf("protocol_tcp block is empty")
		}
		configSet = append(configSet, sets...)
	}
	for _, v := range attackSignature["protocol_udp"].([]interface{}) {
		if len(attackSignature["protocol_icmp"].([]interface{})) != 0 ||
			len(attackSignature["protocol_icmpv6"].([]interface{})) != 0 ||
			len(attackSignature["protocol_tcp"].([]interface{})) != 0 {
			return configSet, fmt.Errorf("protocol_udp cannot be specified with " +
				"protocol_icmp or protocol_icmpv6 or protocol_tcp")
		}
		sets := setSecurityIdpCustomAttackTypeSignatureProtoUDP(setPrefix, v.(map[string]interface{}))
		if len(sets) == 0 {
			return configSet, fmt.Errorf("protocol_udp block is empty")
		}
		configSet = append(configSet, sets...)
	}
	if !strings.Contains(setPrefixOrigin, "attack-type chain member") &&
		attackSignature["protocol_binding"].(string) != "" {
		configSet = append(configSet, setPrefix+"protocol-binding "+attackSignature["protocol_binding"].(string))
	}
	if v := attackSignature["regexp"].(string); v != "" {
		configSet = append(configSet, setPrefix+"regexp \""+v+"\"")
	}
	if v := attackSignature["shellcode"].(string); v != "" {
		configSet = append(configSet, setPrefix+"shellcode "+v)
	}

	return configSet, nil
}

func setSecurityIdpCustomAttackTypeSignatureProtoICMP(
	v6 bool, setPrefixOrigin string, protoICMP map[string]interface{}) []string {
	configSet := make([]string, 0)
	setPrefix := setPrefixOrigin + "protocol icmp "
	if v6 {
		setPrefix = setPrefixOrigin + "protocol icmpv6 "
	}
	if match := protoICMP["checksum_validate_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"checksum-validate match "+match)
	}
	if value := protoICMP["checksum_validate_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"checksum-validate value "+strconv.Itoa(value))
	}
	if match := protoICMP["code_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"code match "+match)
	}
	if value := protoICMP["code_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"code value "+strconv.Itoa(value))
	}
	if match := protoICMP["data_length_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"data-length match "+match)
	}
	if value := protoICMP["data_length_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"data-length value "+strconv.Itoa(value))
	}
	if match := protoICMP["identification_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"identification match "+match)
	}
	if value := protoICMP["identification_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"identification value "+strconv.Itoa(value))
	}
	if match := protoICMP["sequence_number_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"sequence-number match "+match)
	}
	if value := protoICMP["sequence_number_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"sequence-number value "+strconv.Itoa(value))
	}
	if match := protoICMP["type_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"type match "+match)
	}
	if value := protoICMP["type_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"type value "+strconv.Itoa(value))
	}

	return configSet
}

func setSecurityIdpCustomAttackTypeSignatureProtoIPv4(
	setPrefixOrigin string, protoIPv4 map[string]interface{}) []string {
	configSet := make([]string, 0)
	setPrefix := setPrefixOrigin + "protocol ipv4 "
	if match := protoIPv4["checksum_validate_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"checksum-validate match "+match)
	}
	if value := protoIPv4["checksum_validate_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"checksum-validate value "+strconv.Itoa(value))
	}
	if match := protoIPv4["destination_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"destination match "+match)
	}
	if value := protoIPv4["destination_value"].(string); value != "" {
		configSet = append(configSet, setPrefix+"destination value "+value)
	}
	if match := protoIPv4["identification_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"identification match "+match)
	}
	if value := protoIPv4["identification_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"identification value "+strconv.Itoa(value))
	}
	if match := protoIPv4["ihl_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"ihl match "+match)
	}
	if value := protoIPv4["ihl_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"ihl value "+strconv.Itoa(value))
	}
	for _, flags := range sortSetOfString(protoIPv4["ip_flags"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"ip-flags "+flags)
	}
	if match := protoIPv4["protocol_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"protocol match "+match)
	}
	if value := protoIPv4["protocol_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"protocol value "+strconv.Itoa(value))
	}
	if match := protoIPv4["source_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"source match "+match)
	}
	if value := protoIPv4["source_value"].(string); value != "" {
		configSet = append(configSet, setPrefix+"source value "+value)
	}
	if match := protoIPv4["tos_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"tos match "+match)
	}
	if value := protoIPv4["tos_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"tos value "+strconv.Itoa(value))
	}
	if match := protoIPv4["total_length_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"total-length match "+match)
	}
	if value := protoIPv4["total_length_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"total-length value "+strconv.Itoa(value))
	}
	if match := protoIPv4["ttl_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"ttl match "+match)
	}
	if value := protoIPv4["ttl_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"ttl value "+strconv.Itoa(value))
	}

	return configSet
}

func setSecurityIdpCustomAttackTypeSignatureProtoIPv6(
	setPrefixOrigin string, protoIPv6 map[string]interface{}) []string {
	configSet := make([]string, 0)
	setPrefix := setPrefixOrigin + "protocol ipv6 "
	if match := protoIPv6["destination_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"destination match "+match)
	}
	if value := protoIPv6["destination_value"].(string); value != "" {
		configSet = append(configSet, setPrefix+"destination value "+value)
	}
	if match := protoIPv6["extension_header_destination_option_home_address_match"].(string); match != "" {
		configSet = append(configSet,
			setPrefix+"extension-header destination-option home-address match "+match)
	}
	if value := protoIPv6["extension_header_destination_option_home_address_value"].(string); value != "" {
		configSet = append(configSet,
			setPrefix+"extension-header destination-option home-address value "+value)
	}
	if match := protoIPv6["extension_header_destination_option_type_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"extension-header destination-option option-type match "+match)
	}
	if value := protoIPv6["extension_header_destination_option_type_value"].(int); value != -1 {
		configSet = append(configSet,
			setPrefix+"extension-header destination-option option-type value "+strconv.Itoa(value))
	}
	if match := protoIPv6["extension_header_routing_header_type_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"extension-header routing-header header-type match "+match)
	}
	if value := protoIPv6["extension_header_routing_header_type_value"].(int); value != -1 {
		configSet = append(configSet,
			setPrefix+"extension-header routing-header header-type value "+strconv.Itoa(value))
	}
	if match := protoIPv6["flow_label_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"flow-label match "+match)
	}
	if value := protoIPv6["flow_label_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"flow-label value "+strconv.Itoa(value))
	}
	if match := protoIPv6["hop_limit_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"hop-limit match "+match)
	}
	if value := protoIPv6["hop_limit_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"hop-limit value "+strconv.Itoa(value))
	}
	if match := protoIPv6["next_header_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"next-header match "+match)
	}
	if value := protoIPv6["next_header_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"next-header value "+strconv.Itoa(value))
	}
	if match := protoIPv6["payload_length_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"payload-length match "+match)
	}
	if value := protoIPv6["payload_length_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"payload-length value "+strconv.Itoa(value))
	}
	if match := protoIPv6["source_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"source match "+match)
	}
	if value := protoIPv6["source_value"].(string); value != "" {
		configSet = append(configSet, setPrefix+"source value "+value)
	}
	if match := protoIPv6["traffic_class_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"traffic-class match "+match)
	}
	if value := protoIPv6["traffic_class_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"traffic-class value "+strconv.Itoa(value))
	}

	return configSet
}

func setSecurityIdpCustomAttackTypeSignatureProtoTCP(
	setPrefixOrigin string, protoTCP map[string]interface{}) []string {
	configSet := make([]string, 0)
	setPrefix := setPrefixOrigin + "protocol tcp "
	if match := protoTCP["ack_number_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"ack-number match "+match)
	}
	if value := protoTCP["ack_number_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"ack-number value "+strconv.Itoa(value))
	}
	if match := protoTCP["checksum_validate_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"checksum-validate match "+match)
	}
	if value := protoTCP["checksum_validate_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"checksum-validate value "+strconv.Itoa(value))
	}
	if match := protoTCP["data_length_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"data-length match "+match)
	}
	if value := protoTCP["data_length_value"].(int); value != 0 {
		configSet = append(configSet, setPrefix+"data-length value "+strconv.Itoa(value))
	}
	if match := protoTCP["destination_port_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"destination-port match "+match)
	}
	if value := protoTCP["destination_port_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"destination-port value "+strconv.Itoa(value))
	}
	if match := protoTCP["header_length_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"header-length match "+match)
	}
	if value := protoTCP["header_length_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"header-length value "+strconv.Itoa(value))
	}
	if match := protoTCP["mss_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"mss match "+match)
	}
	if value := protoTCP["mss_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"mss value "+strconv.Itoa(value))
	}
	if match := protoTCP["option_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"option match "+match)
	}
	if value := protoTCP["option_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"option value "+strconv.Itoa(value))
	}
	if match := protoTCP["reserved_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"reserved match "+match)
	}
	if value := protoTCP["reserved_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"reserved value "+strconv.Itoa(value))
	}
	if match := protoTCP["sequence_number_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"sequence-number match "+match)
	}
	if value := protoTCP["sequence_number_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"sequence-number value "+strconv.Itoa(value))
	}
	if match := protoTCP["source_port_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"source-port match "+match)
	}
	if value := protoTCP["source_port_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"source-port value "+strconv.Itoa(value))
	}
	for _, flags := range sortSetOfString(protoTCP["tcp_flags"].(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"tcp-flags "+flags)
	}
	if match := protoTCP["urgent_pointer_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"urgent-pointer match "+match)
	}
	if value := protoTCP["urgent_pointer_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"urgent-pointer value "+strconv.Itoa(value))
	}
	if match := protoTCP["window_scale_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"window-scale match "+match)
	}
	if value := protoTCP["window_scale_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"window-scale value "+strconv.Itoa(value))
	}
	if match := protoTCP["window_size_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"window-size match "+match)
	}
	if value := protoTCP["window_size_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"window-size value "+strconv.Itoa(value))
	}

	return configSet
}

func setSecurityIdpCustomAttackTypeSignatureProtoUDP(
	setPrefixOrigin string, protoUDP map[string]interface{}) []string {
	configSet := make([]string, 0)
	setPrefix := setPrefixOrigin + "protocol udp "
	if match := protoUDP["checksum_validate_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"checksum-validate match "+match)
	}
	if value := protoUDP["checksum_validate_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"checksum-validate value "+strconv.Itoa(value))
	}
	if match := protoUDP["data_length_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"data-length match "+match)
	}
	if value := protoUDP["data_length_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"data-length value "+strconv.Itoa(value))
	}
	if match := protoUDP["destination_port_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"destination-port match "+match)
	}
	if value := protoUDP["destination_port_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"destination-port value "+strconv.Itoa(value))
	}
	if match := protoUDP["source_port_match"].(string); match != "" {
		configSet = append(configSet, setPrefix+"source-port match "+match)
	}
	if value := protoUDP["source_port_value"].(int); value != -1 {
		configSet = append(configSet, setPrefix+"source-port value "+strconv.Itoa(value))
	}

	return configSet
}

func readSecurityIdpCustomAttack(customAttack string, m interface{}, jnprSess *NetconfObject) (
	idpCustomAttackOptions, error) {
	sess := m.(*Session)
	var confRead idpCustomAttackOptions
	confRead.timeBindingCount = -1 // default to -1

	customAttackConfig, err := sess.command("show configuration"+
		" security idp custom-attack \""+customAttack+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if customAttackConfig != emptyWord {
		confRead.name = customAttack
		for _, item := range strings.Split(customAttackConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "recommended-action "):
				confRead.recommendedAction = strings.TrimPrefix(itemTrim, "recommended-action ")
			case strings.HasPrefix(itemTrim, "severity "):
				confRead.severity = strings.TrimPrefix(itemTrim, "severity ")
			case strings.HasPrefix(itemTrim, "attack-type anomaly "):
				if len(confRead.attackTypeAnomaly) == 0 {
					confRead.attackTypeAnomaly = append(confRead.attackTypeAnomaly, genSecurityIdpCustomAttackTypeAnomaly(false))
				}
				readSecurityIdpCustomAttackTypeAnomaly(
					strings.TrimPrefix(itemTrim, "attack-type anomaly "), confRead.attackTypeAnomaly[0])
			case strings.HasPrefix(itemTrim, "attack-type chain "):
				if len(confRead.attackTypeChain) == 0 {
					confRead.attackTypeChain = append(confRead.attackTypeChain, map[string]interface{}{
						"member":           make([]map[string]interface{}, 0),
						"expression":       "",
						"order":            false,
						"protocol_binding": "",
						"reset":            false,
						"scope":            "",
					})
				}
				if err := readSecurityIdpCustomAttackTypeChain(
					strings.TrimPrefix(itemTrim, "attack-type chain "), confRead.attackTypeChain[0]); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "attack-type signature "):
				if len(confRead.attackTypeSignature) == 0 {
					confRead.attackTypeSignature = append(confRead.attackTypeSignature, genSecurityIdpCustomAttackTypeSignature(false))
				}
				if err := readSecurityIdpCustomAttackTypeSignature(
					strings.TrimPrefix(itemTrim, "attack-type signature "), confRead.attackTypeSignature[0]); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "time-binding count "):
				var err error
				confRead.timeBindingCount, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "time-binding count "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "time-binding scope "):
				confRead.timeBindingScope = strings.TrimPrefix(itemTrim, "time-binding scope ")
			}
		}
	}

	return confRead, nil
}

func genSecurityIdpCustomAttackTypeAnomaly(chain bool) map[string]interface{} {
	r := map[string]interface{}{
		"direction": "",
		"service":   "",
		"test":      "",
		"shellcode": "",
	}
	if chain {
		delete(r, "service")
	}

	return r
}

func genSecurityIdpCustomAttackTypeSignature(chain bool) map[string]interface{} {
	r := map[string]interface{}{
		"context":          "",
		"direction":        "",
		"negate":           false,
		"pattern":          "",
		"pattern_pcre":     "",
		"protocol_icmp":    make([]map[string]interface{}, 0),
		"protocol_icmpv6":  make([]map[string]interface{}, 0),
		"protocol_ipv4":    make([]map[string]interface{}, 0),
		"protocol_ipv6":    make([]map[string]interface{}, 0),
		"protocol_tcp":     make([]map[string]interface{}, 0),
		"protocol_udp":     make([]map[string]interface{}, 0),
		"protocol_binding": "",
		"regexp":           "",
		"shellcode":        "",
	}
	if chain {
		delete(r, "protocol_binding")
	}

	return r
}

func readSecurityIdpCustomAttackTypeAnomaly(itemTrim string, attackTypeAnomaly map[string]interface{}) {
	switch {
	case strings.HasPrefix(itemTrim, "direction "):
		attackTypeAnomaly["direction"] = strings.TrimPrefix(itemTrim, "direction ")
	case strings.HasPrefix(itemTrim, "service "):
		attackTypeAnomaly["service"] = strings.TrimPrefix(itemTrim, "service ")
	case strings.HasPrefix(itemTrim, "test "):
		attackTypeAnomaly["test"] = strings.TrimPrefix(itemTrim, "test ")
	case strings.HasPrefix(itemTrim, "shellcode "):
		attackTypeAnomaly["shellcode"] = strings.TrimPrefix(itemTrim, "shellcode ")
	}
}

func readSecurityIdpCustomAttackTypeChain(itemTrim string, attackTypeChain map[string]interface{}) error {
	switch {
	case strings.HasPrefix(itemTrim, "member "):
		itemTrimSplit := strings.Split(itemTrim, " ")
		attack := map[string]interface{}{
			"name":                  strings.Trim(itemTrimSplit[1], "\""),
			"attack_type_anomaly":   make([]map[string]interface{}, 0),
			"attack_type_signature": make([]map[string]interface{}, 0),
		}
		itemTrimMember := strings.TrimPrefix(itemTrim, "member "+itemTrimSplit[1]+" ")
		attackTypeChain["member"] = copyAndRemoveItemMapList(
			"name", attack, attackTypeChain["member"].([]map[string]interface{}))
		switch {
		case strings.HasPrefix(itemTrimMember, "attack-type anomaly "):
			if len(attack["attack_type_anomaly"].([]map[string]interface{})) == 0 {
				attack["attack_type_anomaly"] = append(
					attack["attack_type_anomaly"].([]map[string]interface{}),
					genSecurityIdpCustomAttackTypeAnomaly(true))
			}
			readSecurityIdpCustomAttackTypeAnomaly(
				strings.TrimPrefix(itemTrimMember, "attack-type anomaly "),
				attack["attack_type_anomaly"].([]map[string]interface{})[0])
		case strings.HasPrefix(itemTrimMember, "attack-type signature "):
			if len(attack["attack_type_signature"].([]map[string]interface{})) == 0 {
				attack["attack_type_signature"] = append(
					attack["attack_type_signature"].([]map[string]interface{}),
					genSecurityIdpCustomAttackTypeSignature(true))
			}
			if err := readSecurityIdpCustomAttackTypeSignature(
				strings.TrimPrefix(itemTrimMember, "attack-type signature "),
				attack["attack_type_signature"].([]map[string]interface{})[0]); err != nil {
				return err
			}
		}
		attackTypeChain["member"] = append(attackTypeChain["member"].([]map[string]interface{}), attack)
	case strings.HasPrefix(itemTrim, "expression "):
		attackTypeChain["expression"] = strings.Trim(strings.TrimPrefix(itemTrim, "expression "), "\"")
	case itemTrim == "order":
		attackTypeChain["order"] = true
	case strings.HasPrefix(itemTrim, "protocol-binding "):
		attackTypeChain["protocol_binding"] = strings.TrimPrefix(itemTrim, "protocol-binding ")
	case itemTrim == "reset":
		attackTypeChain["reset"] = true
	case strings.HasPrefix(itemTrim, "scope "):
		attackTypeChain["scope"] = strings.TrimPrefix(itemTrim, "scope ")
	}

	return nil
}

func readSecurityIdpCustomAttackTypeSignature(itemTrim string, attackTypeSignature map[string]interface{}) error {
	switch {
	case strings.HasPrefix(itemTrim, "context "):
		attackTypeSignature["context"] = strings.Trim(strings.TrimPrefix(itemTrim, "context "), "\"")
	case strings.HasPrefix(itemTrim, "direction "):
		attackTypeSignature["direction"] = strings.TrimPrefix(itemTrim, "direction ")
	case itemTrim == "negate":
		attackTypeSignature["negate"] = true
	case strings.HasPrefix(itemTrim, "pattern "):
		attackTypeSignature["pattern"] = strings.Trim(strings.TrimPrefix(itemTrim, "pattern "), "\"")
	case strings.HasPrefix(itemTrim, "pattern-pcre "):
		attackTypeSignature["pattern_pcre"] = strings.Trim(strings.TrimPrefix(itemTrim, "pattern-pcre "), "\"")
	case strings.HasPrefix(itemTrim, "protocol icmp "):
		if len(attackTypeSignature["protocol_icmp"].([]map[string]interface{})) == 0 {
			attackTypeSignature["protocol_icmp"] = append(
				attackTypeSignature["protocol_icmp"].([]map[string]interface{}),
				genSecurityIdpCustomAttackTypeSignatureProtoICMP())
		}
		if err := readSecurityIdpCustomAttackTypeSignatureProtoICMP(
			strings.TrimPrefix(itemTrim, "protocol icmp "),
			attackTypeSignature["protocol_icmp"].([]map[string]interface{})[0]); err != nil {
			return err
		}
	case strings.HasPrefix(itemTrim, "protocol icmpv6 "):
		if len(attackTypeSignature["protocol_icmpv6"].([]map[string]interface{})) == 0 {
			attackTypeSignature["protocol_icmpv6"] = append(
				attackTypeSignature["protocol_icmpv6"].([]map[string]interface{}),
				genSecurityIdpCustomAttackTypeSignatureProtoICMP())
		}
		if err := readSecurityIdpCustomAttackTypeSignatureProtoICMP(
			strings.TrimPrefix(itemTrim, "protocol icmpv6 "),
			attackTypeSignature["protocol_icmpv6"].([]map[string]interface{})[0]); err != nil {
			return err
		}
	case strings.HasPrefix(itemTrim, "protocol ipv4 "):
		if len(attackTypeSignature["protocol_ipv4"].([]map[string]interface{})) == 0 {
			attackTypeSignature["protocol_ipv4"] = append(
				attackTypeSignature["protocol_ipv4"].([]map[string]interface{}), genSecurityIdpCustomAttackTypeSignatureProtoIPv4())
		}
		if err := readSecurityIdpCustomAttackTypeSignatureProtoIPv4(
			strings.TrimPrefix(itemTrim, "protocol ipv4 "),
			attackTypeSignature["protocol_ipv4"].([]map[string]interface{})[0]); err != nil {
			return err
		}
	case strings.HasPrefix(itemTrim, "protocol ipv6 "):
		if len(attackTypeSignature["protocol_ipv6"].([]map[string]interface{})) == 0 {
			attackTypeSignature["protocol_ipv6"] = append(
				attackTypeSignature["protocol_ipv6"].([]map[string]interface{}), genSecurityIdpCustomAttackTypeSignatureProtoIPv6())
		}
		if err := readSecurityIdpCustomAttackTypeSignatureProtoIPv6(
			strings.TrimPrefix(itemTrim, "protocol ipv6 "),
			attackTypeSignature["protocol_ipv6"].([]map[string]interface{})[0]); err != nil {
			return err
		}
	case strings.HasPrefix(itemTrim, "protocol tcp "):
		if len(attackTypeSignature["protocol_tcp"].([]map[string]interface{})) == 0 {
			attackTypeSignature["protocol_tcp"] = append(
				attackTypeSignature["protocol_tcp"].([]map[string]interface{}), genSecurityIdpCustomAttackTypeSignatureProtoTCP())
		}
		if err := readSecurityIdpCustomAttackTypeSignatureProtoTCP(
			strings.TrimPrefix(itemTrim, "protocol tcp "),
			attackTypeSignature["protocol_tcp"].([]map[string]interface{})[0]); err != nil {
			return err
		}
	case strings.HasPrefix(itemTrim, "protocol udp "):
		if len(attackTypeSignature["protocol_udp"].([]map[string]interface{})) == 0 {
			attackTypeSignature["protocol_udp"] = append(
				attackTypeSignature["protocol_udp"].([]map[string]interface{}), genSecurityIdpCustomAttackTypeSignatureProtoUDP())
		}
		if err := readSecurityIdpCustomAttackTypeSignatureProtoUDP(
			strings.TrimPrefix(itemTrim, "protocol udp "),
			attackTypeSignature["protocol_udp"].([]map[string]interface{})[0]); err != nil {
			return err
		}
	case strings.HasPrefix(itemTrim, "protocol-binding "):
		attackTypeSignature["protocol_binding"] = strings.TrimPrefix(itemTrim, "protocol-binding ")
	case strings.HasPrefix(itemTrim, "regexp "):
		attackTypeSignature["regexp"] = strings.Trim(strings.TrimPrefix(itemTrim, "regexp "), "\"")
	case strings.HasPrefix(itemTrim, "shellcode "):
		attackTypeSignature["shellcode"] = strings.TrimPrefix(itemTrim, "shellcode ")
	}

	return nil
}

func genSecurityIdpCustomAttackTypeSignatureProtoICMP() map[string]interface{} {
	return map[string]interface{}{
		"checksum_validate_match": "",
		"checksum_validate_value": -1,
		"code_match":              "",
		"code_value":              -1,
		"data_length_match":       "",
		"data_length_value":       -1,
		"identification_match":    "",
		"identification_value":    -1,
		"sequence_number_match":   "",
		"sequence_number_value":   -1,
		"type_match":              "",
		"type_value":              -1,
	}
}

func genSecurityIdpCustomAttackTypeSignatureProtoIPv4() map[string]interface{} {
	return map[string]interface{}{
		"checksum_validate_match": "",
		"checksum_validate_value": -1,
		"destination_match":       "",
		"destination_value":       "",
		"identification_match":    "",
		"identification_value":    -1,
		"ihl_match":               "",
		"ihl_value":               -1,
		"ip_flags":                make([]string, 0),
		"protocol_match":          "",
		"protocol_value":          -1,
		"source_match":            "",
		"source_value":            "",
		"tos_match":               "",
		"tos_value":               -1,
		"total_length_match":      "",
		"total_length_value":      -1,
		"ttl_match":               "",
		"ttl_value":               -1,
	}
}

func genSecurityIdpCustomAttackTypeSignatureProtoIPv6() map[string]interface{} {
	return map[string]interface{}{
		"destination_match": "",
		"destination_value": "",
		"extension_header_destination_option_home_address_match": "",
		"extension_header_destination_option_home_address_value": "",
		"extension_header_destination_option_type_match":         "",
		"extension_header_destination_option_type_value":         -1,
		"extension_header_routing_header_type_match":             "",
		"extension_header_routing_header_type_value":             -1,
		"flow_label_match":     "",
		"flow_label_value":     -1,
		"hop_limit_match":      "",
		"hop_limit_value":      -1,
		"next_header_match":    "",
		"next_header_value":    -1,
		"payload_length_match": "",
		"payload_length_value": -1,
		"source_match":         "",
		"source_value":         "",
		"traffic_class_match":  "",
		"traffic_class_value":  -1,
	}
}

func genSecurityIdpCustomAttackTypeSignatureProtoTCP() map[string]interface{} {
	return map[string]interface{}{
		"ack_number_match":        "",
		"ack_number_value":        -1,
		"checksum_validate_match": "",
		"checksum_validate_value": -1,
		"data_length_match":       "",
		"data_length_value":       0,
		"destination_port_match":  "",
		"destination_port_value":  -1,
		"header_length_match":     "",
		"header_length_value":     -1,
		"mss_match":               "",
		"mss_value":               -1,
		"option_match":            "",
		"option_value":            -1,
		"reserved_match":          "",
		"reserved_value":          -1,
		"sequence_number_match":   "",
		"sequence_number_value":   -1,
		"source_port_match":       "",
		"source_port_value":       -1,
		"tcp_flags":               make([]string, 0),
		"urgent_pointer_match":    "",
		"urgent_pointer_value":    -1,
		"window_scale_match":      "",
		"window_scale_value":      -1,
		"window_size_match":       "",
		"window_size_value":       -1,
	}
}

func genSecurityIdpCustomAttackTypeSignatureProtoUDP() map[string]interface{} {
	return map[string]interface{}{
		"checksum_validate_match": "",
		"checksum_validate_value": -1,
		"data_length_match":       "",
		"data_length_value":       -1,
		"destination_port_match":  "",
		"destination_port_value":  -1,
		"source_port_match":       "",
		"source_port_value":       -1,
	}
}

func readSecurityIdpCustomAttackTypeSignatureProtoICMP(itemTrim string, protoICMP map[string]interface{}) error {
	switch {
	case strings.HasPrefix(itemTrim, "checksum-validate match "):
		protoICMP["checksum_validate_match"] = strings.TrimPrefix(itemTrim, "checksum-validate match ")
	case strings.HasPrefix(itemTrim, "checksum-validate value "):
		var err error
		protoICMP["checksum_validate_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "checksum-validate value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "code match "):
		protoICMP["code_match"] = strings.TrimPrefix(itemTrim, "code match ")
	case strings.HasPrefix(itemTrim, "code value "):
		var err error
		protoICMP["code_value"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "code value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "data-length match "):
		protoICMP["data_length_match"] = strings.TrimPrefix(itemTrim, "data-length match ")
	case strings.HasPrefix(itemTrim, "data-length value "):
		var err error
		protoICMP["data_length_value"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "data-length value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "identification match "):
		protoICMP["identification_match"] = strings.TrimPrefix(itemTrim, "identification match ")
	case strings.HasPrefix(itemTrim, "identification value "):
		var err error
		protoICMP["identification_value"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "identification value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "sequence-number match "):
		protoICMP["sequence_number_match"] = strings.TrimPrefix(itemTrim, "sequence-number match ")
	case strings.HasPrefix(itemTrim, "sequence-number value "):
		var err error
		protoICMP["sequence_number_value"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "sequence-number value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "type match "):
		protoICMP["type_match"] = strings.TrimPrefix(itemTrim, "type match ")
	case strings.HasPrefix(itemTrim, "type value "):
		var err error
		protoICMP["type_value"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "type value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}

	return nil
}

func readSecurityIdpCustomAttackTypeSignatureProtoIPv4(itemTrim string, protoIPv4 map[string]interface{}) error {
	switch {
	case strings.HasPrefix(itemTrim, "checksum-validate match "):
		protoIPv4["checksum_validate_match"] = strings.TrimPrefix(itemTrim, "checksum-validate match ")
	case strings.HasPrefix(itemTrim, "checksum-validate value "):
		var err error
		protoIPv4["checksum_validate_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "checksum-validate value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "destination match "):
		protoIPv4["destination_match"] = strings.TrimPrefix(itemTrim, "destination match ")
	case strings.HasPrefix(itemTrim, "destination value "):
		protoIPv4["destination_value"] = strings.TrimPrefix(itemTrim, "destination value ")
	case strings.HasPrefix(itemTrim, "identification match "):
		protoIPv4["identification_match"] = strings.TrimPrefix(itemTrim, "identification match ")
	case strings.HasPrefix(itemTrim, "identification value "):
		var err error
		protoIPv4["identification_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "identification value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "ihl match "):
		protoIPv4["ihl_match"] = strings.TrimPrefix(itemTrim, "ihl match ")
	case strings.HasPrefix(itemTrim, "ihl value "):
		var err error
		protoIPv4["ihl_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "ihl value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "ip-flags "):
		protoIPv4["ip_flags"] = append(protoIPv4["ip_flags"].([]string), strings.TrimPrefix(itemTrim, "ip-flags "))
	case strings.HasPrefix(itemTrim, "protocol match "):
		protoIPv4["protocol_match"] = strings.TrimPrefix(itemTrim, "protocol match ")
	case strings.HasPrefix(itemTrim, "protocol value "):
		var err error
		protoIPv4["protocol_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "protocol value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "source match "):
		protoIPv4["source_match"] = strings.TrimPrefix(itemTrim, "source match ")
	case strings.HasPrefix(itemTrim, "source value "):
		protoIPv4["source_value"] = strings.TrimPrefix(itemTrim, "source value ")
	case strings.HasPrefix(itemTrim, "tos match "):
		protoIPv4["tos_match"] = strings.TrimPrefix(itemTrim, "tos match ")
	case strings.HasPrefix(itemTrim, "tos value "):
		var err error
		protoIPv4["tos_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "tos value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "total-length match "):
		protoIPv4["total_length_match"] = strings.TrimPrefix(itemTrim, "total-length match ")
	case strings.HasPrefix(itemTrim, "total-length value "):
		var err error
		protoIPv4["total_length_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "total-length value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "ttl match "):
		protoIPv4["ttl_match"] = strings.TrimPrefix(itemTrim, "ttl match ")
	case strings.HasPrefix(itemTrim, "ttl value "):
		var err error
		protoIPv4["ttl_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "ttl value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}

	return nil
}

func readSecurityIdpCustomAttackTypeSignatureProtoIPv6(itemTrim string, protoIPv6 map[string]interface{}) error {
	switch {
	case strings.HasPrefix(itemTrim, "destination match "):
		protoIPv6["destination_match"] = strings.TrimPrefix(itemTrim, "destination match ")
	case strings.HasPrefix(itemTrim, "destination value "):
		protoIPv6["destination_value"] = strings.TrimPrefix(itemTrim, "destination value ")
	case strings.HasPrefix(itemTrim, "extension-header destination-option home-address match "):
		protoIPv6["extension_header_destination_option_home_address_match"] =
			strings.TrimPrefix(itemTrim, "extension-header destination-option home-address match ")
	case strings.HasPrefix(itemTrim, "extension-header destination-option home-address value "):
		protoIPv6["extension_header_destination_option_home_address_value"] =
			strings.TrimPrefix(itemTrim, "extension-header destination-option home-address value ")
	case strings.HasPrefix(itemTrim, "extension-header destination-option option-type match "):
		protoIPv6["extension_header_destination_option_type_match"] =
			strings.TrimPrefix(itemTrim, "extension-header destination-option option-type match ")
	case strings.HasPrefix(itemTrim, "extension-header destination-option option-type value "):
		var err error
		protoIPv6["extension_header_destination_option_type_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "extension-header destination-option option-type value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "extension-header routing-header header-type match "):
		protoIPv6["extension_header_routing_header_type_match"] =
			strings.TrimPrefix(itemTrim, "extension-header routing-header header-type match ")
	case strings.HasPrefix(itemTrim, "extension-header routing-header header-type value "):
		var err error
		protoIPv6["extension_header_routing_header_type_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "extension-header routing-header header-type value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "flow-label match "):
		protoIPv6["flow_label_match"] = strings.TrimPrefix(itemTrim, "flow-label match ")
	case strings.HasPrefix(itemTrim, "flow-label value "):
		var err error
		protoIPv6["flow_label_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "flow-label value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "hop-limit match "):
		protoIPv6["hop_limit_match"] = strings.TrimPrefix(itemTrim, "hop-limit match ")
	case strings.HasPrefix(itemTrim, "hop-limit value "):
		var err error
		protoIPv6["hop_limit_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "hop-limit value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "next-header match "):
		protoIPv6["next_header_match"] = strings.TrimPrefix(itemTrim, "next-header match ")
	case strings.HasPrefix(itemTrim, "next-header value "):
		var err error
		protoIPv6["next_header_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "next-header value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "payload-length match "):
		protoIPv6["payload_length_match"] = strings.TrimPrefix(itemTrim, "payload-length match ")
	case strings.HasPrefix(itemTrim, "payload-length value "):
		var err error
		protoIPv6["payload_length_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "payload-length value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "source match "):
		protoIPv6["source_match"] = strings.TrimPrefix(itemTrim, "source match ")
	case strings.HasPrefix(itemTrim, "source value "):
		protoIPv6["source_value"] = strings.TrimPrefix(itemTrim, "source value ")
	case strings.HasPrefix(itemTrim, "traffic-class match "):
		protoIPv6["traffic_class_match"] = strings.TrimPrefix(itemTrim, "traffic-class match ")
	case strings.HasPrefix(itemTrim, "traffic-class value "):
		var err error
		protoIPv6["traffic_class_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "traffic-class value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}

	return nil
}

func readSecurityIdpCustomAttackTypeSignatureProtoTCP(itemTrim string, protoTCP map[string]interface{}) error {
	switch {
	case strings.HasPrefix(itemTrim, "ack-number match "):
		protoTCP["ack_number_match"] = strings.TrimPrefix(itemTrim, "ack-number match ")
	case strings.HasPrefix(itemTrim, "ack-number value "):
		var err error
		protoTCP["ack_number_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "ack-number value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "checksum-validate match "):
		protoTCP["checksum_validate_match"] = strings.TrimPrefix(itemTrim, "checksum-validate match ")
	case strings.HasPrefix(itemTrim, "checksum-validate value "):
		var err error
		protoTCP["checksum_validate_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "checksum-validate value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "data-length match "):
		protoTCP["data_length_match"] = strings.TrimPrefix(itemTrim, "data-length match ")
	case strings.HasPrefix(itemTrim, "data-length value "):
		var err error
		protoTCP["data_length_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "data-length value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "destination-port match "):
		protoTCP["destination_port_match"] = strings.TrimPrefix(itemTrim, "destination-port match ")
	case strings.HasPrefix(itemTrim, "destination-port value "):
		var err error
		protoTCP["destination_port_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "destination-port value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "header-length match "):
		protoTCP["header_length_match"] = strings.TrimPrefix(itemTrim, "header-length match ")
	case strings.HasPrefix(itemTrim, "header-length value "):
		var err error
		protoTCP["header_length_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "header-length value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "mss match "):
		protoTCP["mss_match"] = strings.TrimPrefix(itemTrim, "mss match ")
	case strings.HasPrefix(itemTrim, "mss value "):
		var err error
		protoTCP["mss_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "mss value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "option match "):
		protoTCP["option_match"] = strings.TrimPrefix(itemTrim, "option match ")
	case strings.HasPrefix(itemTrim, "option value "):
		var err error
		protoTCP["option_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "option value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "reserved match "):
		protoTCP["reserved_match"] = strings.TrimPrefix(itemTrim, "reserved match ")
	case strings.HasPrefix(itemTrim, "reserved value "):
		var err error
		protoTCP["reserved_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "reserved value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "sequence-number match "):
		protoTCP["sequence_number_match"] = strings.TrimPrefix(itemTrim, "sequence-number match ")
	case strings.HasPrefix(itemTrim, "sequence-number value "):
		var err error
		protoTCP["sequence_number_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "sequence-number value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "source-port match "):
		protoTCP["source_port_match"] = strings.TrimPrefix(itemTrim, "source-port match ")
	case strings.HasPrefix(itemTrim, "source-port value "):
		var err error
		protoTCP["source_port_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "source-port value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "tcp-flags "):
		protoTCP["tcp_flags"] = append(protoTCP["tcp_flags"].([]string), strings.TrimPrefix(itemTrim, "tcp-flags "))
	case strings.HasPrefix(itemTrim, "urgent-pointer match "):
		protoTCP["urgent_pointer_match"] = strings.TrimPrefix(itemTrim, "urgent-pointer match ")
	case strings.HasPrefix(itemTrim, "urgent-pointer value "):
		var err error
		protoTCP["urgent_pointer_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "urgent-pointer value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "window-scale match "):
		protoTCP["window_scale_match"] = strings.TrimPrefix(itemTrim, "window-scale match ")
	case strings.HasPrefix(itemTrim, "window-scale value "):
		var err error
		protoTCP["window_scale_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "window-scale value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "window-size match "):
		protoTCP["window_size_match"] = strings.TrimPrefix(itemTrim, "window-size match ")
	case strings.HasPrefix(itemTrim, "window-size value "):
		var err error
		protoTCP["window_size_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "window-size value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}

	return nil
}

func readSecurityIdpCustomAttackTypeSignatureProtoUDP(itemTrim string, protoUDP map[string]interface{}) error {
	switch {
	case strings.HasPrefix(itemTrim, "checksum-validate match "):
		protoUDP["checksum_validate_match"] = strings.TrimPrefix(itemTrim, "checksum-validate match ")
	case strings.HasPrefix(itemTrim, "checksum-validate value "):
		var err error
		protoUDP["checksum_validate_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "checksum-validate value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "data-length match "):
		protoUDP["data_length_match"] = strings.TrimPrefix(itemTrim, "data-length match ")
	case strings.HasPrefix(itemTrim, "data-length value "):
		var err error
		protoUDP["data_length_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "data-length value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "destination-port match "):
		protoUDP["destination_port_match"] = strings.TrimPrefix(itemTrim, "destination-port match ")
	case strings.HasPrefix(itemTrim, "destination-port value "):
		var err error
		protoUDP["destination_port_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "destination-port value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "source-port match "):
		protoUDP["source_port_match"] = strings.TrimPrefix(itemTrim, "source-port match ")
	case strings.HasPrefix(itemTrim, "source-port value "):
		var err error
		protoUDP["source_port_value"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "source-port value "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}

	return nil
}

func delSecurityIdpCustomAttack(customAttack string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete security idp custom-attack \"" + customAttack + "\""}

	return sess.configSet(configSet, jnprSess)
}

func fillSecurityIdpCustomAttackData(d *schema.ResourceData, idpCustomAttackOptions idpCustomAttackOptions) {
	if tfErr := d.Set("name", idpCustomAttackOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("recommended_action", idpCustomAttackOptions.recommendedAction); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("severity", idpCustomAttackOptions.severity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("attack_type_anomaly", idpCustomAttackOptions.attackTypeAnomaly); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("attack_type_chain", idpCustomAttackOptions.attackTypeChain); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("attack_type_signature", idpCustomAttackOptions.attackTypeSignature); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("time_binding_count", idpCustomAttackOptions.timeBindingCount); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("time_binding_scope", idpCustomAttackOptions.timeBindingScope); tfErr != nil {
		panic(tfErr)
	}
}
