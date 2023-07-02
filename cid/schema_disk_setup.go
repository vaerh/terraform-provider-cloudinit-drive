package cid

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var DeviceAliases = schema.MapAttribute{
	Optional: true,
	MarkdownDescription: "Path to disk to be aliased by this name. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#disk-setup)",
	ElementType: types.StringType,
}

type DiskSetupType map[string]struct {
	TableType types.String `tfsdk:"table_type"`
	Layout    types.String `tfsdk:"layout"`
	Overwrite types.Bool   `tfsdk:"overwrite"`
}

var DiskSetup = schema.MapNestedAttribute{
	Optional: true,
	MarkdownDescription: "Disk partitioning. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#disk-setup)",

	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"table_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("mbr"),
				MarkdownDescription: "The partition table type, either mbr or gpt. *Default: mbr.*",
				Validators: []validator.String{
					stringvalidator.OneOf("mbr", "gpt"),
				},
			},
			"layout": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Partitions can be specified by providing a list to layout, where each entry in " +
					"the list is either a size or a list containing a size and the numerical value for a partition " +
					"type. The size for partitions is specified in percentage of disk space, not in bytes (e.g. a " +
					"size of 33 would take up 1/3 of the disk space). The partition type defaults to '83' (Linux " +
					"partition), for other types of partition, such as Linux swap, the type must be passed as part " +
					"of a list along with the size. **Boolean not supported**",
			},
			"overwrite": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Controls whether this module tries to be safe about writing partition tables or " +
					"not. If overwrite: false is set, the device will be checked for a partition table and for a " +
					"file system and if either is found, the operation will be skipped. If overwrite: true is set, " +
					"no checks will be performed. Using overwrite: true is dangerous and can lead to data loss, so " +
					"double check that the correct device has been specified if using this option. *Default: false.*",
			},
		},
	},
}

type FsSetupType struct {
	Label      types.String `tfsdk:"label"`
	Filesystem types.String `tfsdk:"filesystem"`
	Device     types.String `tfsdk:"device"`
	Partition  types.String `tfsdk:"partition"`
	Overwrite  types.Bool   `tfsdk:"overwrite"`
	ReplaceFs  types.String `tfsdk:"replace_fs"`
	ExtraOpts  types.String `tfsdk:"extra_opts"`
	Cmd        types.String `tfsdk:"cmd"`
}

var FsSetup = schema.ListNestedBlock{
	MarkdownDescription: "File system configuration . " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#disk-setup)",

	NestedObject: schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"label": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Label for the filesystem.",
			},
			"filesystem": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filesystem type to create. E.g., ext4 or btrfs.",
			},
			"device": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Specified either as a path or as an alias in the format <alias name>.<y> where " +
					"<y> denotes the partition number on the device. If specifying device using the <device name>." +
					"<partition number> format, the value of partition will be overwritten.",
			},
			"partition": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "*(string/auto/any/none)* The partition option can be set to auto, in which " +
					"this module will search for the existence of a filesystem matching the label, type and device " +
					"of the fs_setup entry and will skip creating the filesystem if one is found. The partition " +
					"option may also be set to any, in which case any file system that matches type and device will " +
					"cause this module to skip filesystem creation for the fs_setup entry, regardless of label " +
					"matching or not. To write a filesystem directly to a device, use partition: none. partition: " +
					"none will always write the filesystem, even when the label and filesystem are matched, and " +
					"overwrite is false. *Integer not supported*",
			},
			"overwrite": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If true, overwrite any existing filesystem. Using overwrite: true for " +
					"filesystems is dangerous and can lead to data loss, so double check the entry in fs_setup. " +
					"Default: false.",
			},
			"replace_fs": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Ignored unless partition is auto or any.",
			},
			"extra_opts": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional options to pass to the filesystem creation command. Ignored if you " +
					"using cmd directly.",
			},
			"cmd": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional command to run to create the filesystem. Can include string " +
					"substitutions  of the other fs_setup config keys. This is only necessary if you need to " +
					"override the default command.",
			},
		},
	},
}
