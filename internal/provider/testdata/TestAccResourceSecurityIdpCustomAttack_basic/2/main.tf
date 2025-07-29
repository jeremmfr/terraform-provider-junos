resource "junos_security_idp_custom_attack" "testacc_idpCustomAttack" {
  name               = "testacc/#1_"
  severity           = "minor"
  time_binding_count = 120
  time_binding_scope = "peer"
  attack_type_chain {
    protocol_binding = "application HTTP"
    member {
      name = "testacc/#1_chain_member1"
      attack_type_anomaly {
        direction = "any"
        test      = "MISSING_HOST"
        shellcode = "all"
      }
    }
    member {
      name = "testacc/#1_chain_member2"
      attack_type_anomaly {
        direction = "any"
        test      = "ACCEPT_LANGUAGE_OVERFLOW"
        shellcode = "all"
      }
    }
    member {
      name = "testacc/#1_chain_member3"
      attack_type_signature {
        context   = "http-url"
        direction = "any"
        protocol_tcp {
          tcp_flags = ["syn"]
        }
      }
    }
  }
}
