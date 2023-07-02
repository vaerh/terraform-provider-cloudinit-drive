package cid

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PowerStateType struct {
	Delay     types.Int64  `tfsdk:"delay"`
	Mode      types.String `tfsdk:"mode"`
	Message   types.String `tfsdk:"message"`
	Timeout   types.Int64  `tfsdk:"timeout"`
	Condition types.String `tfsdk:"condition"`
}

var PowerState = schema.SingleNestedBlock{
	MarkdownDescription: "Change power state. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#power-state-change)",

	Attributes: map[string]schema.Attribute{
		"delay": schema.Int64Attribute{
			Optional: true,
			MarkdownDescription: "Time in minutes to delay after cloud-init has finished. " +
				"If no delay time is specified, the action will take place immediately.",
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"mode": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Must be one of poweroff, halt, or reboot.",
			Validators: []validator.String{
				stringvalidator.OneOf([]string{"poweroff", "halt", "reboot"}...),
			},
		},
		"message": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Optional message to display to the user when the system is powering off or rebooting.",
		},
		"timeout": schema.Int64Attribute{
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(30),
			MarkdownDescription: "Time in seconds to wait for the cloud-init process to finish before executing " +
				"shutdown. **Default: 30.**",
		},
		"condition": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "Apply state change only if condition is met. May be true (always met), " +
				"false (never met), or a command string(s) (in list representation) to be executed. For command " +
				"formatting, see the " +
				"[documentation](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#runcmd). " +
				"If exit code is 0, condition is met, otherwise not.",
		},
	},
}
