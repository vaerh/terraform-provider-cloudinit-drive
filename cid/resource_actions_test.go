package cid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMarshal(t *testing.T) {
	v := cloudInitDriveResourceModel{
		DriveType:          types.StringValue("configdrive2"),
		AllowPublicSshKeys: types.BoolValue(true),
		DeviceAliases: types.MapValueMust(types.StringType, map[string]attr.Value{
			"my_alias":  types.StringValue("/dev/sdb"),
			"swap_disk": types.StringValue("/dev/sdc"),
		}),
		Hostname: types.StringValue("testhost.fqdn"),
		// Mounts: types.ListValueMust(types.StringType, []attr.Value{
		// 	types.StringValue("[ /dev/ephemeral0, /mnt, auto, \"defaults,noexec\" ]"),
		// 	types.StringValue("[ sdc, /opt/data ]"),
		// 	types.StringValue("[ xvdh, /opt/data, auto, \"defaults,nofail\", \"0\", \"0\" ]"),
		// }),
		// CaCerts: CaCertType{
		// 	RemoveDefaults: types.BoolValue(false),
		// 	Trusted: types.SetValueMust(types.StringType, []attr.Value{
		// 		types.StringValue("cid/ya.ru.pem"),
		// 	}),
		// },
		FsSetup: []FsSetupType{
			{
				Label:      types.StringValue("fs1"),
				Filesystem: types.StringValue("ext4"),
				Device:     types.StringValue("my_alias.1"),
				Cmd:        types.StringValue("mkfs -t %(filesystem)s -L %(label)s %(device)s"),
			},
		},
		// CustomFiles: CustomFilesType{
		// 	NetworkData: types.StringValue("/custom/file/path"),
		// },
		WriteFiles: []WriteFilesType{
			{
				Encoding:  types.StringValue("gz+b64"),
				LocalFile: types.StringValue("hello.sh"),
				Path:      types.StringValue("/tmp/aaa"),
			},
		},
	}

	TerraformToGo(v)
}
