---
layout: "junos"
page_title: "Junos: junos_null_commit_file"
sidebar_current: "docs-junos-resource-null-commmit-file"
description: |-
  Load a file with set/delete lines on device and commit
---

# junos_null_commit_file

Load a file with set/delete lines on device and commit

~> **NOTE:** Not provide a real resource, just load content of file with set/delete lines to candidate configuration on device, and commit  

## Example Usage

```hcl
# Load file and commit
variable "setfile" { default = "~/junos/setfile" }
resource junos_null_commit_file "setfile" {
  filename = var.setfile
  triggers = {
    md5 = filemd5(var.setfile)
  }
  # clear_file_after_commit = true
}
```

## Argument Reference

The following arguments are supported:

* `filename` - (Required, Forces new resource)(`String`) The path of the file to load
* `append_lines` - (Optional, Forces new resource)(`ListOfString`) List of lines append to lines in the loaded file.
* `clear_file_after_commit` - (Optional, Forces new resource)(`Bool`) Truncate file after successful commit.
* `triggers` - (Optional, Forces new resource)(`Map`) A map of arbitrary strings that, when changed, will force the resource to be replaced.
