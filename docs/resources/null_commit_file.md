---
page_title: "Junos: junos_null_commit_file"
---

# junos_null_commit_file

Load a file with set/delete lines on device and commit

~> **NOTE:** Not provide a real resource, just load content of file with set/delete lines to
candidate configuration on device, and commit  

## Example Usage

```hcl
# Load file and commit
variable "setfile" { default = "~/junos/setfile" }
resource "junos_null_commit_file" "setfile" {
  filename = var.setfile
  triggers = {
    md5 = filemd5(var.setfile)
  }
  # clear_file_after_commit = true
}
```

## Argument Reference

The following arguments are supported:

- **filename** (Required, String, Forces new resource)  
  The path of the file to load
- **append_lines** (Optional, List of String, Forces new resource)  
  List of lines append to lines in the loaded file.
- **clear_file_after_commit** (Optional, Boolean, Forces new resource)  
  Truncate file after successful commit.
- **triggers** (Optional, Map, Forces new resource)  
  A map of arbitrary strings that, when changed, will force the resource to be replaced.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<filename>`.
