package cid

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ resource.Resource              = &cloudInitDriveResource{}
	_ resource.ResourceWithConfigure = &cloudInitDriveResource{}
	// _ resource.ResourceWithImportState = &cloudInitDriveResource{}
)

type cloudInitDriveResource struct {
	client *SSHClient
}

func NewCloudInitDriveResource() resource.Resource {
	return &cloudInitDriveResource{}
}

func (r *cloudInitDriveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName
}

func (r *cloudInitDriveResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*SSHClient)
}

func (r *cloudInitDriveResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "Unique identifier for this resource: " +
					"hexadecimal representation of the SHA1 checksum of the resource.",
			},
			"drive_type": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "[ configdrive2 | nocloud | opennebula ] Specifies the cloud-init configuration format " +
					"(Proxmox VE use the nocloud format for Linux, and configdrive2 for Windows).",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{ConfigDrive2.String(), NoCloud.String(), OpenNebula.String()}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"drive_path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The path where the cloud-init drive will be saved.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(ssh|file)://`),
						"must start with the scheme 'ssh://' or 'file://'",
					)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"drive_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the cloud-init drive.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`cloudinit`),
						"the disk name must not contain the word 'cloudinit' or the VM startup error will occur",
					)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"checksum": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SHA256 checksum of the generated ISO image.",
			},
			// --- Metadata ---
			"instance_id": schema.StringAttribute{ // Metadata: uuid, instance-id
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Instance ID. If the field is empty, a UID will be generated.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// --- User data ---
			"allow_public_ssh_keys":      AllowPublicSshKeys,
			"device_aliases":             DeviceAliases,
			"disable_root":               DisableRoot,
			"disable_root_opts":          DisableRootOpts,
			"disk_setup":                 DiskSetup,
			"preserve_hostname":          PreserveHostname,
			"final_message":              FinalMessage,
			"fqdn":                       Fqdn,
			"groups":                     Groups,
			"hostname":                   Hostname,
			"prefer_fqdn_over_hostname":  PreferFqdnOverHostname,
			"locale":                     Locale,
			"locale_configfile":          LocaleConfig,
			"manage_etc_hosts":           ManageEtcHosts,
			"manage_resolv_conf":         ManageResolvConf,
			"mount_default_fields":       MountDefaultFields,
			"mounts":                     Mounts,
			"package_reboot_if_required": PackageRebootIfRequired,
			"package_update":             PackageUpdate,
			"package_upgrade":            PackageUpgrade,
			"packages":                   Packages,
			"password":                   Password,
			"runcmd":                     RunCmd,
			"ssh_authorized_keys":        SshAuthorizedKeys,
			"ssh_deletekeys":             SshDeleteKeys,
			"ssh_genkeytypes":            SshGenKeyTypes,
			"ssh_keys":                   SshKeys,
			"ssh_pwauth":                 SshPwAuth,
			"ssh_quiet_keygen":           SshQuietKeygen,
			"timezone":                   TimeZone,
		},

		Blocks: map[string]schema.Block{
			"ca_certs":             CaCerts,
			"chpasswd":             ChPasswd,
			"custom_files":         CustomFiles,
			"fs_setup":             FsSetup,
			"growpart":             GrowPart,
			"keyboard":             Keyboard,
			"power_state":          PowerState,
			"random_seed":          SeedRandom,
			"resolv_conf":          ResolvConf,
			"ssh_publish_hostkeys": SshPublishHostkeys,
			"swap":                 Swap,
			"users":                Users,
			"wireguard":            Wireguard,
			"write_files":          WriteFiles,

			// Network Config.
			"network": NetConf(),
		},
	}
}
