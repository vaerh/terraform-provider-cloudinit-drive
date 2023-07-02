package cid

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

var FinalMessage = schema.StringAttribute{
	Optional: true,
	MarkdownDescription: "This module configures the final message that cloud-init writes. The message is specified " +
		"as a jinja template with the following variables set " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#final-message):\n" +
		"  * `version`: cloud-init version\n" +
		"  * `timestamp`: time at cloud-init finish\n" +
		"  * `datasource`: cloud-init data source\n" +
		"  * `uptime`: system uptime",
}
