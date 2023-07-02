package cid

import (
	"context"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// Version 2.
// https://canonical-cloud-init.readthedocs-hosted.com/en/latest/reference/network-config-format-v2.html

type NetConfTransform struct {
	Version   int
	Ethernets map[string]any `yaml:"ethernets,omitempty"`
	Bonds     map[string]any `yaml:"bonds,omitempty"`
	Bridges   map[string]any `yaml:"bridges,omitempty"`
	Vlans     map[string]any `yaml:"vlans,omitempty"`
}

func MakeNetConfig(ctx context.Context, planData *cloudInitDriveResourceModel) ([]byte, diag.Diagnostics) {
	var diags diag.Diagnostics

	if planData.NetConf == nil {
		diags.AddAttributeError(path.Root("network"),
			"The network settings block is not found in the configuration.",
			"You must specify a file with network settings or an empty 'network {}' block.")
		return nil, diags
	}

	var Network = NetConfTransform{
		Version: int(planData.NetConf.Version.ValueInt64()),
	}

	if l := len(planData.NetConf.Ethernets); l > 0 {
		Network.Ethernets = make(map[string]any, l)
		for _, iface := range planData.NetConf.Ethernets {
			Network.Ethernets[iface.Alias.ValueString()] = TerraformToGo(iface)
		}
	}

	if l := len(planData.NetConf.Bonds); l > 0 {
		Network.Bonds = make(map[string]any, l)
		for _, iface := range planData.NetConf.Bonds {
			Network.Bonds[iface.Alias.ValueString()] = TerraformToGo(iface)
		}
	}

	if l := len(planData.NetConf.Bridges); l > 0 {
		Network.Bridges = make(map[string]any, l)
		for _, iface := range planData.NetConf.Bridges {
			Network.Bridges[iface.Alias.ValueString()] = TerraformToGo(iface)
		}
	}

	if l := len(planData.NetConf.Vlans); l > 0 {
		Network.Vlans = make(map[string]any, l)
		for _, iface := range planData.NetConf.Vlans {
			Network.Vlans[iface.Alias.ValueString()] = TerraformToGo(iface)
		}
	}

	buf, err := yaml.MarshalWithOptions(&struct {
		*NetConfTransform `yaml:"network"`
	}{&Network},
		yaml.UseLiteralStyleIfMultiline(true),
		yaml.UseSingleQuote(true),
	)

	if err != nil {
		diags.AddAttributeError(path.Empty(), err.Error(), "")
	}

	return buf, diags
}
