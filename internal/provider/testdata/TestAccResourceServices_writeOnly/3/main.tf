resource "junos_services" "testacc_wo" {
  security_intelligence {
    authentication_token     = "abcdefghijklmnopqrstuvwxyz123456"
    url                      = "https://example.com/api/manifest.xml"
    url_parameter_wo         = "test_param2"
    url_parameter_wo_version = 2
  }
  user_identification {
    identity_management {
      connection {
        primary_address                    = "192.0.2.254"
        primary_client_id                  = "clientID"
        primary_client_secret_wo           = "mySecretBis"
        primary_client_secret_wo_version   = 2
        connect_method                     = "https"
        port                               = 2000
        query_api                          = "user_query/v2"
        secondary_address                  = "192.0.2.253"
        secondary_client_id                = "clientID2"
        secondary_client_secret_wo         = "mySecret2"
        secondary_client_secret_wo_version = 1
        token_api                          = "oauth_token/oauth"
      }
    }
  }
}
