package cid

import (
	"context"
	"encoding/json"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// https://docs.openstack.org/nova/latest/user/metadata.html#config-drives
type ConfigDriveMeta struct {
	Id            string `json:"uuid"`
	Hostname      string `json:"hostname"`
	NetworkConfig struct {
		ContentPath string `json:"content_path"`
	} `json:"network_config"`
}

// https://cloudinit.readthedocs.io/en/latest/topics/datasources/nocloud.html
type NoCloudMeta struct {
	Id       string `yaml:"instance-id"`
	Hostname string `yaml:"local-hostname"`
	SeedFrom string `yaml:"seedfrom,omitempty"`
}

func MakeMetadata(ctx context.Context, planData *cloudInitDriveResourceModel) ([]byte, diag.Diagnostics) {
	var b []byte
	var err error
	var diags diag.Diagnostics

	/*
		https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html-single/configuring_and_managing_cloud-init_for_rhel_8/index
		When you create or restore an instance from a backup image, the instance ID changes.
		The change in instance ID can cause cloud-init to update configuration files.

		https://github.com/lxc/lxd/issues/9814
		Ok, so we introduce a new volatile.cloud-init.machine-id which gets auto-reset if:
			* Name is changed (copy or rename)
			* User/Vendor data has changed (cloud-init.vendor-data, cloud-init.user-data, user.vendor-data,
				user.user-data)
			* Network config have changed (cloud-init.network-config, user.network-config)
			* The set of nic type devices exposed to the container (ExpandedDevices) has changed in a way that would
				affect the list of interface names in the container
	*/
	var id, host string
	var ok bool

	if id, ok = GetAttribute(planData.InstanceId); !ok {
		id = uuid.New().String()
		planData.InstanceId = types.StringValue(id)
	}

	if host, ok = GetAttribute(planData.Fqdn); !ok {
		if host, ok = GetAttribute(planData.Hostname); !ok {
			host = "localhost"
		}
	}

	// The attribute is mandatory and always defined.
	switch ParseCiDriveType(planData.DriveType.ValueString()) {
	case ConfigDrive2:
		meta := &ConfigDriveMeta{
			Id:       id,
			Hostname: host,
		}

		meta.NetworkConfig.ContentPath = "/openstack/content/0000"
		if planData.CustomFiles != nil {
			if value, ok := GetAttribute(planData.CustomFiles.NetworkData); ok {
				meta.NetworkConfig.ContentPath = value
			}
		}

		b, err = json.Marshal(meta)
	case NoCloud:
		meta := &NoCloudMeta{
			Id:       id,
			Hostname: host,
		}
		b, err = yaml.MarshalWithOptions(meta,
			yaml.UseLiteralStyleIfMultiline(true),
			yaml.UseSingleQuote(true),
		)
	}

	if err != nil {
		diags.AddAttributeError(path.Empty(), err.Error(), "")
	}

	return b, diags
}
