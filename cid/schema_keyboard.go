package cid

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KeyboardType struct {
	Layout  types.String `tfsdk:"layout"`
	Model   types.String `tfsdk:"model"`
	Variant types.String `tfsdk:"variant"`
	Options types.String `tfsdk:"options"`
}

var Keyboard = schema.SingleNestedBlock{
	MarkdownDescription: "Set keyboard layout. [Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#keyboard)",

	Attributes: map[string]schema.Attribute{
		"layout": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Keyboard layout. Corresponds to XKBLAYOUT.",
		},
		"model": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("pc105"),
			MarkdownDescription: "Keyboard model. Corresponds to XKBMODEL. **Default: pc105**.",
		},
		"variant": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Keyboard variant. Corresponds to XKBVARIANT.",
		},
		"options": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Keyboard options. Corresponds to XKBOPTIONS.",
		},
	},
}
