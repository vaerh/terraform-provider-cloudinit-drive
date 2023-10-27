package cid

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Network configuration v1
type NetConfTypeV1 struct {
	Version types.Int64          `tfsdk:"version"`
	Config  []NetConfV1Interface `tfsdk:"interface" yaml:"config,omitempty"`
}

type NetConfV1Interface struct {
	Type             types.String       `tfsdk:"type" yaml:"type,omitempty"`
	Name             types.String       `tfsdk:"name" yaml:"name,omitempty"`
	MacAddress       types.String       `tfsdk:"mac_address" yaml:"mac_address,omitempty"`
	Mtu              types.Int64        `tfsdk:"mtu" yaml:"mtu,omitempty"`
	Params           types.Map          `tfsdk:"params" yaml:"params,omitempty"`                  // Bond, Bridge
	BondInterfaces   types.Set          `tfsdk:"bond_interfaces" yaml:"bond_interfaces,flow"`     // Bond
	BridgeInterfaces types.Set          `tfsdk:"bridge_interfaces" yaml:"bridge_interfaces,flow"` // Bridge
	VlanId           types.Int64        `tfsdk:"vlan_id" yaml:"vlan_id,flow"`                     // VLAN
	VlanLink         types.String       `tfsdk:"vlan_link" yaml:"vlan_link,omitempty"`            // VLAN
	DNSAddress       types.Set          `tfsdk:"dns_nameservers" yaml:"address,flow,omitempty"`   // Nameserver
	DNSSearch        types.Set          `tfsdk:"dns_search" yaml:"search,flow,omitempty"`         // Nameserver
	DNSInterface     types.Set          `tfsdk:"dns_interface" yaml:"interface,flow"`             // Nameserver
	RouteDestination types.String       `tfsdk:"route_destination" yaml:"destination,omitempty"`  // Route
	RouteGateway     types.String       `tfsdk:"route_gateway" yaml:"gateway,omitempty"`          // Route
	RouteMetric      types.Int64        `tfsdk:"route_metric" yaml:"metric,omitempty"`            // Route
	Subnets          []NetConfV1Subnets `tfsdk:"subnets" yaml:"subnets,omitempty"`
}

type NetConfV1Subnets struct {
	Type       types.String `tfsdk:"type" yaml:"type,omitempty"`
	Control    types.String `tfsdk:"control" yaml:"control,omitempty"`
	Address    types.String `tfsdk:"address" yaml:"address,flow,omitempty"`
	Netmask    types.String `tfsdk:"netmask" yaml:"netmask,omitempty"`
	Gateway    types.String `tfsdk:"gateway" yaml:"gateway,omitempty"`
	DNSAddress types.Set    `tfsdk:"dns_nameservers" yaml:"dns_nameservers,flow,omitempty"`
	DNSSearch  types.Set    `tfsdk:"dns_search" yaml:"dns_search,flow,omitempty"`
	Routes     []struct {
		Gateway     types.String `tfsdk:"gateway" yaml:"gateway,omitempty"`
		Netmask     types.String `tfsdk:"netmask" yaml:"netmask,omitempty"`
		Destination types.String `tfsdk:"destination" yaml:"destination,omitempty"`
	} `tfsdk:"routes" yaml:"routes,flow,omitempty"`
}

func NetConfV1() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		MarkdownDescription: "Networking Config Version 1. " +
			"[Info](https://canonical-cloud-init.readthedocs-hosted.com/en/latest/reference/network-config-format-v1.html)",

		Attributes: map[string]schema.Attribute{
			"version": schema.Int64Attribute{
				Computed: true,
				Default:  int64default.StaticInt64(1),
			},
		},

		Blocks: map[string]schema.Block{
			"interface": schema.ListNestedBlock{
				MarkdownDescription: "",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "",
						},
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "",
						},
						"mac_address": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "",
						},
						"mtu": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "",
						},
						"params": schema.MapAttribute{
							Optional:            true,
							MarkdownDescription: "",
							ElementType:         types.StringType,
						},
						"bond_interfaces": schema.SetAttribute{
							Optional:            true,
							MarkdownDescription: "",
							ElementType:         types.StringType,
						},
						"bridge_interfaces": schema.SetAttribute{
							Optional:            true,
							MarkdownDescription: "",
							ElementType:         types.StringType,
						},
						"vlan_id": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "",
						},
						"vlan_link": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "",
						},
						"dns_nameservers": schema.SetAttribute{
							Optional:            true,
							MarkdownDescription: "",
							ElementType:         types.StringType,
						},
						"dns_search": schema.SetAttribute{
							Optional:            true,
							MarkdownDescription: "",
							ElementType:         types.StringType,
						},
						"dns_interface": schema.SetAttribute{
							Optional:            true,
							MarkdownDescription: "",
							ElementType:         types.StringType,
						},
						"route_destination": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "",
						},
						"route_gateway": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "",
						},
						"route_metric": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "",
						},
					},

					Blocks: map[string]schema.Block{
						"subnets": schema.ListNestedBlock{
							MarkdownDescription: "",

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "",
									},
									"control": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "",
									},
									"address": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "",
									},
									"netmask": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "",
									},
									"gateway": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "",
									},
									"dns_nameservers": schema.SetAttribute{
										Optional:            true,
										MarkdownDescription: "",
										ElementType:         types.StringType,
									},
									"dns_search": schema.SetAttribute{
										Optional:            true,
										MarkdownDescription: "",
										ElementType:         types.StringType,
									},
								},
								Blocks: map[string]schema.Block{
									"routes": schema.ListNestedBlock{
										MarkdownDescription: "",

										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"gateway": schema.StringAttribute{
													Optional:            true,
													MarkdownDescription: "",
												},
												"netmask": schema.StringAttribute{
													Optional:            true,
													MarkdownDescription: "",
												},
												"destination": schema.StringAttribute{
													Optional:            true,
													MarkdownDescription: "",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Network configuration v2
type NetConfTypeV2 struct {
	Version   types.Int64           `tfsdk:"version"`
	Ethernets []NetConfV2Ethernet   `tfsdk:"ethernets" yaml:"ethernets,omitempty"`
	Bonds     []NetConfV2BondBridge `tfsdk:"bonds" yaml:"bonds,omitempty"`
	Bridges   []NetConfV2BondBridge `tfsdk:"bridges" yaml:"bridges,omitempty"`
	Vlans     []NetConfV2Vlan       `tfsdk:"vlans" yaml:"vlans,omitempty"`
}

type NetConfV2Ethernet struct {
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

type NetConfV2BondBridge struct {
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

type NetConfV2Vlan struct {
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

func NetConfV2() schema.SingleNestedBlock {
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
