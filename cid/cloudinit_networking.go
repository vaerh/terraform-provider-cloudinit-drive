package cid

import (
	"context"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// Version 1.
// https://canonical-cloud-init.readthedocs-hosted.com/en/latest/reference/network-config-format-v1.html

func MakeNetConfigV1(ctx context.Context, planData *cloudInitDriveResourceModel) ([]byte, diag.Diagnostics) {
	var diags diag.Diagnostics

	buf, err := yaml.MarshalWithOptions(&struct {
		Data any `yaml:"network"`
	}{TerraformToGo(planData.NetConfV1)},
		yaml.UseLiteralStyleIfMultiline(true),
		yaml.UseSingleQuote(true),
	)

	if err != nil {
		diags.AddAttributeError(path.Empty(), err.Error(), "")
	}

	return buf, diags
}

// Version 2.
// https://canonical-cloud-init.readthedocs-hosted.com/en/latest/reference/network-config-format-v2.html

type NetConfV2Transform struct {
	Version   int
	Ethernets map[string]any `yaml:"ethernets,omitempty"`
	Bonds     map[string]any `yaml:"bonds,omitempty"`
	Bridges   map[string]any `yaml:"bridges,omitempty"`
	Vlans     map[string]any `yaml:"vlans,omitempty"`
}

func MakeNetConfigV2(ctx context.Context, planData *cloudInitDriveResourceModel) ([]byte, diag.Diagnostics) {
	var diags diag.Diagnostics

	var Network = NetConfV2Transform{
		Version: int(planData.NetConfV2.Version.ValueInt64()),
	}

	if l := len(planData.NetConfV2.Ethernets); l > 0 {
		Network.Ethernets = make(map[string]any, l)
		for _, iface := range planData.NetConfV2.Ethernets {
			Network.Ethernets[iface.Alias.ValueString()] = TerraformToGo(iface)
		}
	}

	if l := len(planData.NetConfV2.Bonds); l > 0 {
		Network.Bonds = make(map[string]any, l)
		for _, iface := range planData.NetConfV2.Bonds {
			Network.Bonds[iface.Alias.ValueString()] = TerraformToGo(iface)
		}
	}

	if l := len(planData.NetConfV2.Bridges); l > 0 {
		Network.Bridges = make(map[string]any, l)
		for _, iface := range planData.NetConfV2.Bridges {
			Network.Bridges[iface.Alias.ValueString()] = TerraformToGo(iface)
		}
	}

	if l := len(planData.NetConfV2.Vlans); l > 0 {
		Network.Vlans = make(map[string]any, l)
		for _, iface := range planData.NetConfV2.Vlans {
			Network.Vlans[iface.Alias.ValueString()] = TerraformToGo(iface)
		}
	}

	buf, err := yaml.MarshalWithOptions(&struct {
		*NetConfV2Transform `yaml:"network"`
	}{&Network},
		yaml.UseLiteralStyleIfMultiline(true),
		yaml.UseSingleQuote(true),
	)

	if err != nil {
		diags.AddAttributeError(path.Empty(), err.Error(), "")
	}

	return buf, diags
}
