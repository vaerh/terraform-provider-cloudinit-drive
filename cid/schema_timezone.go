package cid

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

var TimeZone = schema.StringAttribute{
	Optional: true,
	MarkdownDescription: "The timezone to use as represented in /usr/share/zoneinfo. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#timezone)",
}
