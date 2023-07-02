package cid

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestMakeMetadata(t *testing.T) {
	type args struct {
		ctx  context.Context
		plan *cloudInitDriveResourceModel
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 diag.Diagnostics
	}{
		{
			"NoCloud Hostname",
			args{
				context.Background(),
				&cloudInitDriveResourceModel{
					InstanceId: types.StringValue("4728a390-08c8-4c8b-91c5-30ffd92b623b"),
					Fqdn:       types.StringNull(),
					Hostname:   types.StringValue("hostname"),
					DriveType:  types.StringValue("nocloud"),
				},
			},
			"instance-id: 4728a390-08c8-4c8b-91c5-30ffd92b623b\nlocal-hostname: hostname\n",
			nil,
		},
		{
			"NoCloud FQDN",
			args{
				context.Background(),
				&cloudInitDriveResourceModel{
					InstanceId: types.StringValue("4728a390-08c8-4c8b-91c5-30ffd92b623b"),
					Fqdn:       types.StringValue("fqdn"),
					Hostname:   types.StringNull(),
					DriveType:  types.StringValue("nocloud"),
				},
			},
			"instance-id: 4728a390-08c8-4c8b-91c5-30ffd92b623b\nlocal-hostname: fqdn\n",
			nil,
		},
		{
			"ConfigDrive2 Net Default",
			args{
				context.Background(),
				&cloudInitDriveResourceModel{
					InstanceId: types.StringValue("4728a390-08c8-4c8b-91c5-30ffd92b623b"),
					Fqdn:       types.StringNull(),
					Hostname:   types.StringValue("hostname"),
					DriveType:  types.StringValue("configdrive2"),
					// DriveType: types.StringValue("nocloud"),
					// DriveType:  types.StringValue("opennebula"),
				},
			},
			`{"uuid":"4728a390-08c8-4c8b-91c5-30ffd92b623b","hostname":"hostname","network_config":{"content_path":"/content/0000"}}`,
			nil,
		},
		{
			"ConfigDrive2 Net Custom",
			args{
				context.Background(),
				&cloudInitDriveResourceModel{
					InstanceId: types.StringValue("4728a390-08c8-4c8b-91c5-30ffd92b623b"),
					Fqdn:       basetypes.NewStringUnknown(),
					Hostname:   types.StringValue("hostname"),
					DriveType:  types.StringValue("configdrive2"),
					// DriveType:  types.StringValue("opennebula"),
					CustomFiles: &CustomFilesType{
						NetworkData: types.StringValue("/content/custom/path"),
					},
				},
			},
			`{"uuid":"4728a390-08c8-4c8b-91c5-30ffd92b623b","hostname":"hostname","network_config":{"content_path":"/content/custom/path"}}`,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := MakeMetadata(tt.args.ctx, tt.args.plan)
			if string(got) != tt.want {
				t.Errorf("MakeMetadata() got = %s, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("MakeMetadata() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
