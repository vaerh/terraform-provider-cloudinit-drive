package cid

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CustomFilesType struct {
	ScriptsPerBoot     types.Set    `tfsdk:"scripts_per_boot"`
	ScriptsPerInstance types.Set    `tfsdk:"scripts_per_instance"`
	ScriptsPerOnce     types.Set    `tfsdk:"scripts_per_once"`
	UserData           types.String `tfsdk:"user_data" cid:"skip"`
	MetaData           types.String `tfsdk:"meta_data" cid:"skip"`
	NetworkData        types.String `tfsdk:"network_data" cid:"skip"`
	VendorData         types.String `tfsdk:"vendor_data" cid:"skip"`
	OpennebulaContext  types.String `tfsdk:"opennebula_context" cid:"skip"`
	Files              []FilesType  `tfsdk:"files"`
}

type FilesType struct {
	Src types.String `tfsdk:"src"`
	Dst types.String `tfsdk:"dst"`
}

var CustomFiles = schema.SingleNestedBlock{
	MarkdownDescription: "Overriding settings by existing files.",

	Attributes: map[string]schema.Attribute{
		"scripts_per_boot": schema.SetAttribute{
			Optional: true,
			MarkdownDescription: "Any scripts in the *scripts/per-boot* directory on the datasource will be run " +
				"every time the system boots. Scripts will be run in alphabetical order. This module " +
				"does not accept any config keys. " +
				"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#scripts-per-boot)",
			ElementType: types.StringType,
		},
		"scripts_per_instance": schema.SetAttribute{
			Optional: true,
			MarkdownDescription: "Any scripts in the *scripts/per-instance* directory on the datasource will be " +
				"run when a new instance is first booted. Scripts will be run in alphabetical order. This " +
				"module does not accept any config keys.  " +
				"Some cloud platforms change instance-id if a significant change was made to the system. " +
				"As a result per-instance scripts will run again." +
				"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#scripts-per-instance)",
			ElementType: types.StringType,
		},
		"scripts_per_once": schema.SetAttribute{
			Optional: true,
			MarkdownDescription: "Any scripts in the *scripts/per-once directory* on the datasource will be run " +
				"only once. Changes to the instance will not force a re-run. The only way to re-run " +
				"these scripts is to run the clean subcommand and reboot. Scripts will be run in " +
				"alphabetical order. This module does not accept any config keys." +
				"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#scripts-per-once)",
			ElementType: types.StringType,
		},
		"user_data": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "Use custom user data file. " +
				"[User Data](https://cloudinit.readthedocs.io/en/latest/topics/format.html)",
		},
		"meta_data": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "Use custom metadata file. " +
				"[Instance Metadata](https://cloudinit.readthedocs.io/en/latest/topics/instancedata.html)",
		},
		"network_data": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "Use custom network configuration file. " +
				"[Network Configuration](https://cloudinit.readthedocs.io/en/latest/topics/network-config.html)",
		},
		"vendor_data": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "Use custom vendor data file. " +
				"[Vendor Data](https://cloudinit.readthedocs.io/en/latest/topics/vendordata.html)",
		},
		// TODO Open Nebula
		"opennebula_context": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "[Use custom 'context.sh' file. OpenNebula contextualization variables]" +
				"(https://cloudinit.readthedocs.io/en/latest/topics/datasources/opennebula.html)",
		},
	},

	Blocks: map[string]schema.Block{
		"files": schema.ListNestedBlock{
			MarkdownDescription: "Create a disk image from custom files. The generated configuration files will be " +
				"overwritten by the current files.",

			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"src": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Source file.",
					},
					"dst": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Destination path.",
					},
				},
			},
		},
	},
}
