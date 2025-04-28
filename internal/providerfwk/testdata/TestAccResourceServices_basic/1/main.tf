resource "junos_services_proxy_profile" "testacc_services" {
  name               = "testacc_services"
  protocol_http_host = "192.0.2.1"
  protocol_http_port = 3128
}
resource "junos_security_address_book" "testacc_services" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_services"
  network_address {
    name  = "testacc_services_add"
    value = "192.0.2.0/25"
  }
  address_set {
    name    = "testacc_services_set"
    address = ["testacc_services_add"]
  }
}
resource "junos_services_ssl_initiation_profile" "testacc_services" {
  name = "testacc_services"
}
resource "junos_services" "testacc" {
  depends_on = [
    junos_security_address_book.testacc_services
  ]
  advanced_anti_malware {
    connection {
      auth_tls_profile = junos_services_ssl_initiation_profile.testacc_services.name
      proxy_profile    = junos_services_proxy_profile.testacc_services.name
      source_address   = "192.0.2.1"
      url              = "https://example.com/api/test.xml"
    }
    default_policy {
      blacklist_notification_log        = true
      default_notification_log          = true
      fallback_options_action           = "permit"
      fallback_options_notification_log = true
      http_action                       = "block"
      http_inspection_profile           = "testacc_services"
      http_notification_log             = true
      imap_inspection_profile           = "testacc_services"
      imap_notification_log             = true
      smtp_inspection_profile           = "testacc_services"
      smtp_notification_log             = true
      verdict_threshold                 = 5
      whitelist_notification_log        = true
    }
  }
  application_identification {
    application_system_cache {}
    download {
      automatic_start_time = "12-24.22:00"
    }
    enable_performance_mode {}
    max_transactions = 10
  }
  security_intelligence {
    authentication_token = "abcdefghijklmnopqrstuvwxyz123456"
    category_disable     = ["all"]
    proxy_profile        = junos_services_proxy_profile.testacc_services.name
    url                  = "https://example.com/api/manifest.xml"
    url_parameter        = "test_param"
  }
  user_identification {
    device_info_auth_source = "network-access-controller"
    identity_management {
      connection {
        primary_address          = "192.0.2.254"
        primary_client_id        = "clientID"
        primary_client_secret    = "mySecret"
        connect_method           = "https"
        port                     = 2000
        primary_ca_certificate   = "ca"
        query_api                = "user_query/v2"
        secondary_address        = "192.0.2.253"
        secondary_ca_certificate = "ca2"
        secondary_client_id      = "clientID2"
        secondary_client_secret  = "mySecret2"
        token_api                = "oauth_token/oauth"
      }
      authentication_entry_timeout         = 60
      batch_query_items_per_batch          = 100
      batch_query_interval                 = 30
      filter_domain                        = ["test3", "test2"]
      filter_exclude_ip_address_book       = junos_security_address_book.testacc_services.name
      filter_exclude_ip_address_set        = "testacc_services_set"
      filter_include_ip_address_book       = junos_security_address_book.testacc_services.name
      filter_include_ip_address_set        = "testacc_services_set"
      invalid_authentication_entry_timeout = 60
      ip_query_disable                     = true
      ip_query_delay_time                  = 30
    }
  }
}
