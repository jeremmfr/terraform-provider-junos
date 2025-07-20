resource "junos_services_proxy_profile" "testacc_services" {
  lifecycle {
    create_before_destroy = true
  }

  name               = "testacc_services"
  protocol_http_host = "192.0.2.1"
  protocol_http_port = 3128
}
resource "junos_services_security_intelligence_profile" "testacc_services" {
  lifecycle {
    create_before_destroy = true
  }
  name     = "testacc_services"
  category = "IPFilter"
  rule {
    name = "rule_1"
    match {
      threat_level = [1]
    }
    then_action = "permit"
  }
}
resource "junos_services_ssl_initiation_profile" "testacc_services" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_services"
}
resource "junos_services" "testacc" {
  advanced_anti_malware {
    connection {
      auth_tls_profile = junos_services_ssl_initiation_profile.testacc_services.name
      source_interface = "fxp0.0"
      url              = "https://example.com/api/test.xml"
    }
  }
  application_identification {
    application_system_cache {
      security_services = true
    }
    application_system_cache_timeout = 120
    download {
      automatic_interval       = 120
      automatic_start_time     = "12-24.22:00"
      ignore_server_validation = true
      proxy_profile            = junos_services_proxy_profile.testacc_services.name
      url                      = "https://example.com/"
    }
    enable_performance_mode {
      max_packet_threshold = 50
    }
    imap_cache_size     = 120
    imap_cache_timeout  = 120
    max_transactions    = 10
    statistics_interval = 120
  }
  security_intelligence {
    authentication_token = "abcdefghijklmnopqrstuvwxyz123400"
    category_disable     = ["CC"]
    default_policy {
      category_name = "IPFilter"
      profile_name  = junos_services_security_intelligence_profile.testacc_services.name
    }
    proxy_profile = junos_services_proxy_profile.testacc_services.name
    url           = "https://example.com/api/manifest.xml"
    url_parameter = "test_param_update"
  }
  user_identification {
    device_info_auth_source = "network-access-controller"
    identity_management {
      connection {
        primary_address        = "192.0.2.254"
        primary_client_id      = "clientID"
        primary_client_secret  = "mySecret2"
        connect_method         = "https"
        port                   = 2000
        primary_ca_certificate = "ca@2"
        query_api              = "user_query/v2"
        token_api              = "oauth_token/oauth"
      }
    }
  }
}
