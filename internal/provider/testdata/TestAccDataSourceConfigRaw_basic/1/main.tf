data "junos_config_raw" "test_json" {
  format = "json"
}

data "junos_config_raw" "test_json_minified" {
  format = "json-minified"
}

data "junos_config_raw" "test_set" {
  format = "set"
}

data "junos_config_raw" "test_text" {
  // format = "text" // defaults to text 
}

data "junos_config_raw" "test_xml" {
  format = "xml"
}

data "junos_config_raw" "test_xml_minified" {
  format = "xml-minified"
}

