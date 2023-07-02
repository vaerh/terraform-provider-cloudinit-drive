package cid

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetConfType struct {
	Version   types.Int64         `tfsdk:"version"`
	Ethernets []NetConfEthernet   `tfsdk:"ethernets" yaml:"ethernets,omitempty"`
	Bonds     []NetConfBondBridge `tfsdk:"bonds" yaml:"bonds,omitempty"`
	Bridges   []NetConfBondBridge `tfsdk:"bridges" yaml:"bridges,omitempty"`
	Vlans     []NetConfVlan       `tfsdk:"vlans" yaml:"vlans,omitempty"`
}

type NetConfEthernet struct {
	Alias types.String `tfsdk:"alias" cid:"skip"`
	Match *struct {
		Driver     types.String `tfsdk:"driver" yaml:"driver,omitempty"`
		Macaddress types.String `tfsdk:"macaddress" yaml:"macaddress,omitempty"`
		Name       types.String `tfsdk:"name" yaml:"name,omitempty"`
	} `tfsdk:"match" yaml:"match,omitempty"`
	SetName        types.String `tfsdk:"set_name" yaml:"set-name,omitempty"`
	Wakeonlan      types.Bool   `tfsdk:"wakeonlan" yaml:"wakeonlan,omitempty"`
	Addresses      types.Set    `tfsdk:"addresses" yaml:"addresses,omitempty"`
	Dhcp4          types.Bool   `tfsdk:"dhcp4" yaml:"dhcp4,omitempty"`
	Dhcp6          types.Bool   `tfsdk:"dhcp6" yaml:"dhcp6,omitempty"`
	Dhcp4Overrides types.Map    `tfsdk:"dhcp4_overrides" yaml:"dhcp4-overrides,omitempty"`
	Dhcp6Overrides types.Map    `tfsdk:"dhcp6_overrides" yaml:"dhcp6-overrides,omitempty"`
	Gateway4       types.String `tfsdk:"gateway4" yaml:"gateway4,omitempty"`
	Gateway6       types.String `tfsdk:"gateway6" yaml:"gateway6,omitempty"`
	Mtu            types.Int64  `tfsdk:"mtu" yaml:"mtu,omitempty"`
	Nameservers    *struct {
		Search    types.Set `tfsdk:"search" yaml:"search,flow,omitempty"`
		Addresses types.Set `tfsdk:"addresses" yaml:"addresses,flow,omitempty"`
	} `tfsdk:"nameservers" yaml:"nameservers,omitempty"`
	Renderer types.String `tfsdk:"renderer" yaml:"renderer,omitempty"`
	Routes   []struct {
		To     types.String `tfsdk:"to" yaml:"to,omitempty"`
		Via    types.String `tfsdk:"via" yaml:"via,omitempty"`
		Metric types.Int64  `tfsdk:"metric" yaml:"metric,omitempty"`
	} `tfsdk:"routes" yaml:"routes,omitempty"`
}

type NetConfBondBridge struct {
	Alias          types.String `tfsdk:"alias" cid:"skip"`
	Interfaces     types.Set    `tfsdk:"interfaces" yaml:"interfaces,flow"`
	Parameters     types.Map    `tfsdk:"parameters" yaml:"parameters,omitempty"`
	Addresses      types.Set    `tfsdk:"addresses" yaml:"addresses,omitempty"`
	Dhcp4          types.Bool   `tfsdk:"dhcp4" yaml:"dhcp4,omitempty"`
	Dhcp6          types.Bool   `tfsdk:"dhcp6" yaml:"dhcp6,omitempty"`
	Dhcp4Overrides types.Map    `tfsdk:"dhcp4_overrides" yaml:"dhcp4-overrides,omitempty"`
	Dhcp6Overrides types.Map    `tfsdk:"dhcp6_overrides" yaml:"dhcp6-overrides,omitempty"`
	Gateway4       types.String `tfsdk:"gateway4" yaml:"gateway4,omitempty"`
	Gateway6       types.String `tfsdk:"gateway6" yaml:"gateway6,omitempty"`
	Mtu            types.Int64  `tfsdk:"mtu" yaml:"mtu,omitempty"`
	Nameservers    *struct {
		Search    types.Set `tfsdk:"search" yaml:"search,flow,omitempty"`
		Addresses types.Set `tfsdk:"addresses" yaml:"addresses,flow,omitempty"`
	} `tfsdk:"nameservers" yaml:"nameservers,omitempty"`
	Renderer types.String `tfsdk:"renderer" yaml:"renderer,omitempty"`
	Routes   []struct {
		To     types.String `tfsdk:"to" yaml:"to,omitempty"`
		Via    types.String `tfsdk:"via" yaml:"via,omitempty"`
		Metric types.Int64  `tfsdk:"metric" yaml:"metric,omitempty"`
	} `tfsdk:"routes" yaml:"routes,omitempty"`
}

type NetConfVlan struct {
	Alias          types.String `tfsdk:"alias" cid:"skip"`
	Id             types.Int64  `tfsdk:"id" yaml:"id,flow"`
	Link           types.String `tfsdk:"link" yaml:"link,omitempty"`
	Addresses      types.Set    `tfsdk:"addresses" yaml:"addresses,omitempty"`
	Dhcp4          types.Bool   `tfsdk:"dhcp4" yaml:"dhcp4,omitempty"`
	Dhcp6          types.Bool   `tfsdk:"dhcp6" yaml:"dhcp6,omitempty"`
	Dhcp4Overrides types.Map    `tfsdk:"dhcp4_overrides" yaml:"dhcp4-overrides,omitempty"`
	Dhcp6Overrides types.Map    `tfsdk:"dhcp6_overrides" yaml:"dhcp6-overrides,omitempty"`
	Gateway4       types.String `tfsdk:"gateway4" yaml:"gateway4,omitempty"`
	Gateway6       types.String `tfsdk:"gateway6" yaml:"gateway6,omitempty"`
	Mtu            types.Int64  `tfsdk:"mtu" yaml:"mtu,omitempty"`
	Nameservers    *struct {
		Search    types.Set `tfsdk:"search" yaml:"search,flow,omitempty"`
		Addresses types.Set `tfsdk:"addresses" yaml:"addresses,flow,omitempty"`
	} `tfsdk:"nameservers" yaml:"nameservers,omitempty"`
	Renderer types.String `tfsdk:"renderer" yaml:"renderer,omitempty"`
	Routes   []struct {
		To     types.String `tfsdk:"to" yaml:"to,omitempty"`
		Via    types.String `tfsdk:"via" yaml:"via,omitempty"`
		Metric types.Int64  `tfsdk:"metric" yaml:"metric,omitempty"`
	} `tfsdk:"routes" yaml:"routes,omitempty"`
}

func NetConf() schema.SingleNestedBlock {
	var netSchema = schema.SingleNestedBlock{
		MarkdownDescription: "Networking Config Version 2. " +
			"[Info](https://canonical-cloud-init.readthedocs-hosted.com/en/latest/reference/network-config-format-v2.html)",

		Attributes: map[string]schema.Attribute{
			"version": schema.Int64Attribute{
				Computed: true,
				Default:  int64default.StaticInt64(2),
			},
		},

		Blocks: map[string]schema.Block{
			"ethernets": schema.ListNestedBlock{
				MarkdownDescription: "",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"alias": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "Local interface identifier within the configuration file. " +
								"This identifier can be referenced in the following stanzas.",
						},
						"set_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "",
						},
						"wakeonlan": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "",
						},
					},

					Blocks: map[string]schema.Block{
						"match": schema.SingleNestedBlock{
							MarkdownDescription: "",

							Attributes: map[string]schema.Attribute{
								"driver": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "",
								},
								"macaddress": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "",
								},
								"name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "",
								},
							},
						},
					},
				},
			},
			"bonds": schema.ListNestedBlock{
				MarkdownDescription: "",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"alias": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "Local interface identifier within the configuration file. " +
								"This identifier can be referenced in the following stanzas.",
						},
						"interfaces": schema.SetAttribute{
							Optional:            true,
							MarkdownDescription: "",
							ElementType:         types.StringType,
						},
						"parameters": schema.MapAttribute{
							Optional:            true,
							MarkdownDescription: "",
							ElementType:         types.StringType,
						},
					},

					Blocks: map[string]schema.Block{},
				},
			},
			"bridges": schema.ListNestedBlock{
				MarkdownDescription: "",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"alias": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "Local interface identifier within the configuration file. " +
								"This identifier can be referenced in the following stanzas.",
						},
						"interfaces": schema.SetAttribute{
							Optional:            true,
							MarkdownDescription: "",
							ElementType:         types.StringType,
						},
						"parameters": schema.MapAttribute{
							Optional:            true,
							MarkdownDescription: "",
							ElementType:         types.StringType,
						},
					},

					Blocks: map[string]schema.Block{},
				},
			},
			"vlans": schema.ListNestedBlock{
				MarkdownDescription: "",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"alias": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "Local interface identifier within the configuration file. " +
								"This identifier can be referenced in the following stanzas.",
						},
						"id": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "",
						},
						"link": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "",
						},
					},

					Blocks: map[string]schema.Block{},
				},
			},
		},
	}
	for _, iface := range []string{"ethernets", "bonds", "bridges", "vlans"} {
		for k, v := range netAllDevices.Attributes {
			netSchema.Blocks[iface].(schema.ListNestedBlock).NestedObject.Attributes[k] = v
		}
		for k, v := range netAllDevices.Blocks {
			netSchema.Blocks[iface].(schema.ListNestedBlock).NestedObject.Blocks[k] = v
		}
	}

	return netSchema
}

var netAllDevices = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
		"addresses": schema.SetAttribute{
			Optional:            true,
			MarkdownDescription: "",
			ElementType:         types.StringType,
		},
		"dhcp4": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "",
		},
		"dhcp6": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "",
		},
		"dhcp4_overrides": schema.MapAttribute{
			Optional:            true,
			MarkdownDescription: "",
			ElementType:         types.StringType,
		},
		"dhcp6_overrides": schema.MapAttribute{
			Optional:            true,
			MarkdownDescription: "",
			ElementType:         types.StringType,
		},
		"gateway4": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "",
		},
		"gateway6": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "",
		},
		"mtu": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "",
		},
		"renderer": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "",
		},
	},

	Blocks: map[string]schema.Block{
		"nameservers": schema.SingleNestedBlock{
			MarkdownDescription: "",

			Attributes: map[string]schema.Attribute{
				"addresses": schema.SetAttribute{
					Optional:            true,
					MarkdownDescription: "",
					ElementType:         types.StringType,
				},
				"search": schema.SetAttribute{
					Optional:            true,
					MarkdownDescription: "",
					ElementType:         types.StringType,
				},
			},
		},
		"routes": schema.ListNestedBlock{
			MarkdownDescription: "",

			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"to": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "",
					},
					"via": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "",
					},
					"metric": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "",
					},
				},
			},
		},
	},
}
