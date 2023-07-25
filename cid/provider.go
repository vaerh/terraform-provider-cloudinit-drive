package cid

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ provider.Provider = &CloudInitDriveProvider{}

type CloudInitDriveProvider struct{}

type CloudInitDriveProviderModel struct {
	SSH *SSH `tfsdk:"ssh"`
	// Become *Become `tfsdk:"become"`
}

type SSH struct {
	Host       types.String `tfsdk:"host"`
	Port       types.Int64  `tfsdk:"port"`
	User       types.String `tfsdk:"user"`
	PrivateKey types.String `tfsdk:"private_key"`
	AuthSocket types.String `tfsdk:"auth_socket"`
}

// type Become struct {
// 	Method   types.String `tfsdk:"method"`
// 	User     types.String `tfsdk:"user"`
// 	Password types.String `tfsdk:"password"`
// }

func NewCloudInitDrive() provider.Provider {
	return &CloudInitDriveProvider{}
}

func (p *CloudInitDriveProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudinit-drive"
}

func (p *CloudInitDriveProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{

		Blocks: map[string]schema.Block{
			"ssh": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Address or hostname of the remote destination host of the Cloud-init drive.",
					},
					"port": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "SSH port.",
						Validators:          []validator.Int64{int64validator.Between(0, 65535)},
					},
					"user": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Remote host ssh login.",
					},
					"private_key": schema.StringAttribute{
						Optional:            true,
						Sensitive:           true,
						MarkdownDescription: "Path to the identity file (private key) for public key authentication.",
					},
					"auth_socket": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The OS socket path used to communicate with the SSH agent.",
					},
				},
			},
			// "become": schema.ListNestedBlock{
			// 	MarkdownDescription: "Privilege escalation",
			// 	NestedObject: schema.NestedBlockObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"method": schema.StringAttribute{
			// 				Optional:            true,
			// 				MarkdownDescription: "Privilege escalation method (sudo, su, doas).",
			// 				Validators: []validator.String{
			// 					stringvalidator.OneOf([]string{"sudo", "su", "doas"}...),
			// 				},
			// 			},
			// 			"user": schema.StringAttribute{
			// 				Optional:            true,
			// 				MarkdownDescription: "User with desired privileges.",
			// 			},
			// 			"password": schema.StringAttribute{
			// 				Optional:            true,
			// 				Sensitive:           true,
			// 				MarkdownDescription: "Password for privilege escalation.",
			// 			},
			// 		},
			// 	},
			// },
		},
	}
}

func (p *CloudInitDriveProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring CloudInit-Drive client")

	var cidProvider CloudInitDriveProviderModel
	var diags diag.Diagnostics

	resp.Diagnostics.Append(req.Config.Get(ctx, &cidProvider)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if cidProvider.SSH == nil {
		cidProvider.SSH = &SSH{}
	}

	// if cidProvider.Become == nil {
	// 	cidProvider.Become = &Become{}
	// }

	cidProvider.SSH.Host = LoadEnvVariable(cidProvider.SSH.Host, "SSH_HOST").(types.String)
	cidProvider.SSH.User = LoadEnvVariable(cidProvider.SSH.User, "SSH_USER").(types.String)
	cidProvider.SSH.PrivateKey = LoadEnvVariable(cidProvider.SSH.PrivateKey, "SSH_PRIVATE_KEY").(types.String)

	cidProvider.SSH.Port = LoadEnvVariable(cidProvider.SSH.Port, "SSH_PORT").(types.Int64)
	if cidProvider.SSH.Port.IsNull() {
		cidProvider.SSH.Port = types.Int64Value(22)
	}

	if !cidProvider.SSH.AuthSocket.IsNull() {
		os.Setenv("SSH_AUTH_SOCK", cidProvider.SSH.AuthSocket.ValueString())
	}

	// LoadEnvVariable(cidProvider.Become.Method, "CID_BECOME_METHOD")
	// LoadEnvVariable(cidProvider.Become.User, "CID_BECOME_USER")
	// LoadEnvVariable(cidProvider.Become.Password, "CID_BECOME_PASSWORD")

	/*
		Jump host
		https://github.com/appleboy/easyssh-proxy/blob/master/example/proxy/proxy.go
		https://github.com/loafoe/terraform-provider-ssh/blob/main/ssh/resource_resource.go
	*/

	resp.ResourceData, diags = NewClient(ctx, &cidProvider)
	resp.Diagnostics.Append(diags...)
}

func (p *CloudInitDriveProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCloudInitDriveResource,
		NewMacAddressResource,
	}
}

func (p *CloudInitDriveProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
