package cid

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var Mounts = schema.ListAttribute{
	Optional: true,
	MarkdownDescription: "List of lists. Each inner list entry is a list of /etc/fstab mount declarations of the " +
		"format: `[ fs_spec, fs_file, fs_vfstype, fs_mntops, fs-freq, fs_passno ]`. A mount declaration with " +
		"less than 6 items will get remaining values from mount_default_fields. A mount declaration with only" +
		"fs_spec and no fs_file mountpoint will be skipped. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#mounts)",
	ElementType: types.ListType{
		ElemType: types.StringType,
	},
}

var MountDefaultFields = schema.ListAttribute{
	Optional: true,
	MarkdownDescription: "Default mount configuration for any mount entry with less than 6 options provided. When " +
		"specified, 6 items are required and represent /etc/fstab entries. " +
		"**Default: defaults,nofail,x-systemd.requires=cloud-init.service,_netdev**. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#mounts)",
	ElementType: types.StringType,
}

type SwapType struct {
	Filename types.String `tfsdk:"filename"`
	Size     types.String `tfsdk:"size"`
	MaxSize  types.String `tfsdk:"maxsize"`
}

var Swap = schema.SingleNestedBlock{
	MarkdownDescription: "Swap files can be configured by setting the path to the swap file to create with filename, " +
		"the size of the swap file with size maximum size of the swap file if using an size: auto with " +
		"maxsize. By default no swap file is created. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#mounts)",

	Attributes: map[string]schema.Attribute{
		"filename": schema.StringAttribute{
			Optional:    true,
			Description: "Path to the swap file to create.",
		},
		"size": schema.StringAttribute{
			Optional: true,
			Description: "The size in bytes of the swap file, 'auto' or a human-readable size " +
				"abbreviation of the format <float_size><units> where units are one of B, K, M, G or T.",
		},
		"maxsize": schema.StringAttribute{
			Optional:    true,
			Description: "The maxsize in bytes of the swap file.",
		},
	},
}
