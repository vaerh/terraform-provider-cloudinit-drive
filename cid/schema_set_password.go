package cid

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ChPasswdType struct {
	Expire types.Bool `tfsdk:"expire"`
	Users  []struct {
		Name     types.String `tfsdk:"name"`
		Password types.String `tfsdk:"password"`
		Type     types.String `tfsdk:"type"`
	} `tfsdk:"users"`
}

var ChPasswd = schema.SingleNestedBlock{
	MarkdownDescription: "Set user passwords. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#set-passwords)",

	Attributes: map[string]schema.Attribute{
		"expire": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(true),
			MarkdownDescription: "Whether to expire all user passwords such that a password will need to be " +
				"reset on the user's next login. **Default: true.",
		},
	},

	Blocks: map[string]schema.Block{
		"users": schema.ListNestedBlock{
			MarkdownDescription: "This key represents a list of existing users to set passwords for. Each item " +
				"under users contains the following required keys: **name** and **password** or in the " +
				"case of a randomly generated password, **name** and **type**. The **type** key has a " +
				"default value of **hash**, and may alternatively be set to **text** or **RANDOM**.",

			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Username.`,
					},
					"password": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
						MarkdownDescription: "Password.\n" +
							"    >*<span style=\"color:red\">Security Notice:</span>*\n" +
							"    The password will be stored *unencrypted* in your Terraform state file.\n" +
							"    **Use of this attribute for production deployments is *not* recommended!**",
					},
					"type": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("hash"),
						MarkdownDescription: `Password type (hash, text or RANDOM). *Default: hash.*`,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"hash", "text", "RANDOM"}...),
						},
					},
				},
			},
		},
	},
}

var Password = schema.StringAttribute{
	Optional:  true,
	Sensitive: true,
	MarkdownDescription: "Set the default user's password. **Ignored if chpasswd list is used.**\n" +
		"    >*<span style=\"color:red\">Security Notice:</span>*\n" +
		"    The password will be stored *unencrypted* in your Terraform state file.\n" +
		"    **Use of this attribute for production deployments is *not* recommended!**\n" +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#set-passwords)",
	Validators: []validator.String{
		stringvalidator.ConflictsWith(path.MatchRoot("chpasswd")),
	},
}

var SshPwAuth = schema.BoolAttribute{
	Optional: true,
	MarkdownDescription: " Sets whether or not to accept password authentication. true will enable password auth. " +
		"false will disable. Default is to leave the value unchanged. In order for this config to be applied, SSH " +
		"may need to be restarted. On systemd systems, this restart will only happen if the SSH service has already " +
		"been started. On non-systemd systems, a restart will be attempted regardless of the service state.",
}
