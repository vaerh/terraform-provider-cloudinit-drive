package cid

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var PreserveHostname = schema.BoolAttribute{
	Optional: true,
	MarkdownDescription: "If true, the hostname will not be changed. **Default: false.** " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#set-hostname)",
	Validators: []validator.Bool{
		boolvalidator.ConflictsWith(path.MatchRoot("hostname")),
		boolvalidator.ConflictsWith(path.MatchRoot("fqdn")),
		boolvalidator.ConflictsWith(path.MatchRoot("prefer_fqdn_over_hostname")),
	},
}

var Hostname = schema.StringAttribute{
	Optional: true,
	Computed: true,
	Default:  stringdefault.StaticString("localhost"),
	MarkdownDescription: "Instance hostname. If fqdn is set, the hostname extracted from fqdn overrides hostname." +
		"**Default: 'localhost'.** [Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#set-hostname)",
}

var Fqdn = schema.StringAttribute{
	Optional: true,
	MarkdownDescription: "Fully qualified domain name of the instance. Preferred over hostname if both are provided. " +
		"In absence of hostname and fqdn in cloud-config, the local-hostname value will be used " +
		"from datasource metadata. [Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#set-hostname)",
}

var PreferFqdnOverHostname = schema.BoolAttribute{
	Optional: true,
	MarkdownDescription: "By default, it is distro-dependent whether cloud-init uses the short hostname or fully " +
		"qualified domain name when both local-hostname` and ``fqdn are both present in instance metadata. " +
		"When set true, use fully qualified domain name if present as hostname instead of short hostname. " +
		"When set false, use hostname config value if present, otherwise fallback to fqdn. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#set-hostname)",
}
