package cid

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ManageResolvConf = schema.BoolAttribute{
	Optional: true,
	Computed: true,
	Default:  booldefault.StaticBool(false),
	MarkdownDescription: "Whether to manage the resolv.conf file. resolv_conf block will be ignored unless this is " +
		"set to true. **Default: false.** " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#resolv-conf)",
}

type ResolvConfType struct {
	NameServers    types.List   `tfsdk:"nameservers"`
	SearchDomains types.List   `tfsdk:"searchdomains"`
	Domain        types.String `tfsdk:"domain"`
	SortList      types.List   `tfsdk:"sortlist"`
	Options       types.Map    `tfsdk:"options"`
}

var ResolvConf = schema.SingleNestedBlock{
	MarkdownDescription: "resolv.conf file" +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#resolv-conf)",

	Attributes: map[string]schema.Attribute{
		"nameservers": schema.ListAttribute{
			Optional:            true,
			MarkdownDescription: "A list of nameservers to use to be added as nameserver lines.",
			ElementType:         types.StringType,
		},
		"searchdomains": schema.ListAttribute{
			Optional:            true,
			MarkdownDescription: "A list of domains to be added search line.",
			ElementType:         types.StringType,
		},
		"domain": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The domain to be added as domain line.",
		},
		"sortlist": schema.ListAttribute{
			Optional:            true,
			MarkdownDescription: "A list of IP addresses to be added to sortlist line.",
			ElementType:         types.StringType,
		},
		"options": schema.MapAttribute{
			Optional:            true,
			MarkdownDescription: "Key/value pairs of options to go under options heading.",
			ElementType:         types.StringType,
		},
	},
}
