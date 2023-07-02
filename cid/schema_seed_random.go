package cid

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SeedRandomType struct {
	File            types.String `tfsdk:"file"`
	Data            types.String `tfsdk:"data"`
	Encoding        types.String `tfsdk:"encoding"`
	CommandRequired types.Bool   `tfsdk:"command_required"`
	Command         types.List   `tfsdk:"command"`
}

var SeedRandom = schema.SingleNestedBlock{
	MarkdownDescription: "Provide random seed data. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#seed-random)",

	Attributes: map[string]schema.Attribute{
		"file": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "File to write random data to. Default: /dev/urandom.",
		},
		"data": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "This data will be written to file before data from the datasource.",
			Validators: []validator.String{
				base64Validator{},
			},
		},
		"encoding": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "Used to decode data provided. Allowed values are raw, base64, b64, gzip, gz. " +
				"If you use the gzip format, you need to convert the contents of the file into base64 encoding and " +
				"specify it with the binary data flag: \"**bin'**QmluYXJ5IGRhdGEK**'**\" " +
				"*Default: raw.*",
			Validators: []validator.String{
				stringvalidator.OneOf([]string{"raw", "base64", "b64", "gzip", "gz"}...),
			},
		},
		"command_required": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(false),
			MarkdownDescription: "If true, and command is not available to be run then an exception is raised and " +
				"cloud-init will record failure. Otherwise, only debug error is mentioned. **Default: false.**",
		},
		"command": schema.ListAttribute{
			Optional: true,
			MarkdownDescription: " Execute this command to seed random. The command will have RANDOM_SEED_FILE " +
				"in its environment set to the value of file above.",
			ElementType: types.StringType,
		},
	},
}

type base64Validator struct{}

func (b64 base64Validator) Description(ctx context.Context) string {
	return "string must contain the correct data in base64 encoding"
}

func (b64 base64Validator) MarkdownDescription(ctx context.Context) string {
	return "string must contain the correct data in base64 encoding"
}

func (b64 base64Validator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	var enc types.String
	diags := req.Config.GetAttribute(ctx, path.Root("random_seed").AtName("encoding"), &enc)
	if diags.HasError() {
		return
	}

	if enc.ValueString() == "raw" {
		return
	}

	var s = req.ConfigValue.ValueString()

	if enc.ValueString() == "gz" || enc.ValueString() == "gzip" {
		if len(s) < 5 || s[:5] != "bin__" {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("data"),
				"Invalid Base64 String",
				fmt.Sprintf("String is too small or does not contain a prefix. Expected \"bin__base64 string\", got: %v", s),
			)
			return
		}

		s = s[5:]
	}

	if _, err := base64.StdEncoding.DecodeString(s); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path.AtName("data"),
			"Invalid Base64 String",
			fmt.Sprintf("String must contain the correct data in base64 encoding, got: %v.", err),
		)
	}
}
