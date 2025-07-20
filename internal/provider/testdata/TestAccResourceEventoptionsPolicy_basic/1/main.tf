resource "junos_eventoptions_destination" "testacc_evtopts_policy" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_evtopts_policy#1"
  archive_site {
    url = "https://example.com"
  }
}
resource "junos_system_login_user" "testacc_evtopts_policy" {
  lifecycle {
    create_before_destroy = true
  }
  name  = "testacc_evtopts_policy"
  class = "unauthorized"
}
resource "junos_eventoptions_policy" "testacc_evtopts_policy" {
  name   = "testacc_evtopts_policy#1"
  events = ["aaa_infra_fail", "acct_fork_err"]
  then {
    change_configuration {
      commands                         = ["test2 test1"]
      commit_options_check             = true
      commit_options_check_synchronize = true
    }
    event_script {
      filename = "filename_1"
      destination {
        name           = junos_eventoptions_destination.testacc_evtopts_policy.name
        retry_count    = 1
        retry_interval = 2
        transfer_delay = 3
      }
      output_filename = "filename_2"
      output_format   = "xml"
      user_name       = junos_system_login_user.testacc_evtopts_policy.name
    }
    execute_commands {
      commands = ["test4 test3"]
      destination {
        name           = junos_eventoptions_destination.testacc_evtopts_policy.name
        retry_count    = 1
        retry_interval = 2
        transfer_delay = 3
      }
      output_filename = "filename_3"
      output_format   = "xml"
      user_name       = junos_system_login_user.testacc_evtopts_policy.name
    }
    priority_override_facility = "external"
    priority_override_severity = "info"
    raise_trap                 = true
    upload {
      filename    = "filename_5"
      destination = junos_eventoptions_destination.testacc_evtopts_policy.name
    }
  }
  attributes_match {
    from    = "acct_fork_err.error-message"
    compare = "starts-with"
    to      = "aaa_infra_fail.test3"
  }
  within {
    time_interval = 7
    events        = ["aaa_infra_fail", "acct_fork_err"]
    not_events    = ["aaa_usage_err"]
    trigger_count = 8
    trigger_when  = "after"
  }
}
