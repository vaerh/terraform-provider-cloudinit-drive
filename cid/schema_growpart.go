package cid

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GrowPartType struct {
	Mode                   types.String `tfsdk:"mode"`
	Devices                types.Set    `tfsdk:"devices"`
	IgnoreGrowrootDisabled types.Bool   `tfsdk:"ignore_growroot_disabled"`
}

var GrowPart = schema.SingleNestedBlock{
	MarkdownDescription: "Grow partitions. [Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#growpart)",

	Attributes: map[string]schema.Attribute{
		"mode": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("auto"),
			MarkdownDescription: "The utility to use for resizing. **Default: auto**.",
			Validators: []validator.String{
				stringvalidator.OneOf([]string{"auto", "growpart", "gpart", "off"}...),
			},
		},
		"devices": schema.SetAttribute{
			Optional: true,
			MarkdownDescription: "The devices to resize. Each entry can either be the path to the device's " +
				"mountpoint in the filesystem or a path to the block device in '/dev'. **Default: [/]**.",
			ElementType: types.StringType,
		},
		"ignore_growroot_disabled": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(false),
			Description: "If true, ignore the presence of /etc/growroot-disabled. If false and the file " +
				"exists, then don't resize. **Default: false**.",
		},
	},
}
