package cid

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CaCertType struct {
	RemoveDefaults types.Bool `tfsdk:"remove_defaults"`
	FileName       types.Set  `tfsdk:"filename" cid:"skip"`
	Content        types.Set  `tfsdk:"content" cid:"skip"`
	Trusted        types.Set  `tfsdk:"-" json:"trusted,omitempty" yaml:"trusted,omitempty"`
}

var CaCerts = schema.SingleNestedBlock{
	MarkdownDescription: "Add CA certificates to /etc/ca-certificates.conf and updates the ssl cert cache using " +
		"update-ca-certificates. [Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#ca-certificates)",

	Attributes: map[string]schema.Attribute{
		"remove_defaults": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Remove default CA certificates if true. Default: false.",
		},
		"filename": schema.SetAttribute{
			Optional:            true,
			MarkdownDescription: "List of paths to files containing SSL certificates.",
			ElementType:         types.StringType,
		},
		"content": schema.SetAttribute{
			Optional:            true,
			MarkdownDescription: "List of trusted CA certificates to add (encoded as base64).",
			ElementType:         types.StringType,
		},
	},
}
