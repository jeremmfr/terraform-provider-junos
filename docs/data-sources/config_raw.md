---
page_title: "Junos: junos_config_raw"
---

# junos_config_raw

Get raw configuration from the Junos device in the specified format.

~> **Warning**
The `config` attribute may contain secrets that are hashed using weak hashing algorithms (`$9$`).

## Example Usage

```hcl
# Get configuration in JSON format
data "junos_config_raw" "config_json" {
  format = "json"
}

# Get configuration in minified JSON format
data "junos_config_raw" "config_json_minified" {
  format = "json-minified"
}

# Get configuration in set format
data "junos_config_raw" "config_set" {
  format = "set"
}

# Get configuration in text format (default)
data "junos_config_raw" "config_text" {}

# Get configuration in XML format
data "junos_config_raw" "config_xml" {
  format = "xml"
}

# Get configuration in minified XML format
data "junos_config_raw" "config_xml_minified" {
  format = "xml-minified"
}
```

## Argument Reference

The following arguments are supported:

- **format** (Optional, String)  
  Configuration format.  
  Need to be `json`, `json-minified`, `set`, `text`, `xml` or `xml-minified`.  
  Defaults to `text`.

## Attribute Reference

The following attributes are exported:

- **id** (String)
  An identifier for the data source with format `<format>`.
- **config** (String)
  The raw configuration output in the requested format.
