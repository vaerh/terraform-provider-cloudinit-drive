package cid

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WireguardType struct {
	Readinessprobe types.List `tfsdk:"readinessprobe"`
	Interfaces     []struct {
		Name       types.String `tfsdk:"name"`
		ConfigPath types.String `tfsdk:"config_path"`
		Content    types.String `tfsdk:"content"`
	} `tfsdk:"interfaces"`
}

var Wireguard = schema.SingleNestedBlock{
	MarkdownDescription: "Wireguard tunnel. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#wireguard)",

	Attributes: map[string]schema.Attribute{
		"readinessprobe": schema.ListAttribute{
			Optional:    true,
			Description: "List of shell commands to be executed as probes.",
			ElementType: types.StringType,
		},
	},

	Blocks: map[string]schema.Block{
		"interfaces": schema.ListNestedBlock{
			MarkdownDescription: "Interface list.",

			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Optional:    true,
						Description: "Name of the interface. Typically wgx (example: wg0).",
					},
					"config_path": schema.StringAttribute{
						Optional:    true,
						Description: "Path to configuration file of Wireguard interface.",
					},
					"content": schema.StringAttribute{
						Optional:    true,
						Description: "Wireguard interface configuration. Contains key, peer, etc.",
					},
				},
			},
		},
	},
}
