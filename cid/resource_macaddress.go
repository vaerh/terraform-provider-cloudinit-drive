package cid

import (
	"context"
	crypt "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var macPrefixRe = regexp.MustCompile(`[0-9a-zA-Z]{2}`)

func init() {
	var b [8]byte
	if _, err := crypt.Read(b[:]); err != nil {
		panic("cannot read crypto/rand")
	}

	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

var _ resource.Resource = &macAddressResource{}

type macAddressResource struct{}

type macAddressResourceModel struct {
	Mac    types.List   `tfsdk:"mac"`
	Prefix types.String `tfsdk:"prefix"`
	Suffix types.String `tfsdk:"suffix"`
	LikeIP types.Bool   `tfsdk:"like_ip"`
	Number types.Int64  `tfsdk:"number"`
	Seq    types.Bool   `tfsdk:"sequential"`
}

func NewMacAddressResource() resource.Resource {
	return &macAddressResource{}
}

func (r *macAddressResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mac"
}

func (r *macAddressResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"mac": schema.ListAttribute{
				Computed:            true,
				MarkdownDescription: "Generated MAC address.",
				ElementType:         types.StringType,
			},
			"prefix": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The first 3..5 octets of a MAC address to which the remaining " +
					"part will be generated.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^([0-9a-fA-F]{2}([:-])?){2,5}$`),
						"the attribute must be in the format: 'xx:yy:zz...' "+
							"or 'xx-yy-zz...' or 'xxyyzz...'",
					)},
			},
			"suffix": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "IPv4 address to generate the last 3 octets of the MAC address, which" +
					"avoids collisions when generating MAC addresses in the same broadcast domain.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("number")),
					stringvalidator.ConflictsWith(path.MatchRoot("sequential")),
				},
			},
			"like_ip": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Do not convert octets of the MAC address, but copy from the IP.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("number")),
					boolvalidator.ConflictsWith(path.MatchRoot("sequential")),
				},
			},
			"number": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
				MarkdownDescription: "Number of addresses to be generated.",
				Validators: []validator.Int64{
					int64validator.Between(1, 32),
				},
			},
			"sequential": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "The addresses are generated sequentially.",
			},
			// TODO Maybe add a list of addresses to exclude so that we can keep MAC unique
			// between different resources?
		},
	}
}

func getMAC() ([6]byte, error) {
	var buf [6]byte

	_, err := rand.Read(buf[:])

	// Unicasts
	buf[0] &= 0xfe

	if buf[5] == 0 {
		buf[5] = 1
	}

	return buf, err
}

func (m *macAddressResourceModel) MakeAddresses() (diags diag.Diagnostics) {
	var prefix, suffix []byte
	var count int = 1
	var seq bool
	var res []attr.Value

	if !m.Prefix.IsNull() {
		for _, oct := range macPrefixRe.FindAllString(m.Prefix.ValueString(), -1) {
			b, err := hex.DecodeString(oct)
			if err != nil {
				diags.AddError(err.Error(), "hex.DecodeString() error")
				return
			}
			prefix = append(prefix, b[0])
		}
	}

	if !m.Suffix.IsNull() {
		ip := net.ParseIP(m.Suffix.ValueString()).To4()
		if ip == nil {
			diags.AddAttributeError(path.Root("suffix"), "Attribute isn't an IPv4 address.", "")
			return
		}
		suffix = ip[1:]
	} else {
		if !m.Number.IsNull() {
			count = int(m.Number.ValueInt64())
			if count < 1 {
				count = 1
			}
		}
		if !m.Seq.IsNull() {
			seq = m.Seq.ValueBool()
		}
	}

	if seq {
		mac, err := getMAC()
		if err != nil {
			diags.AddError(err.Error(), "getMAC() error")
			return
		}
		// len(vendor) 2 .. 5
		if len(prefix) != 0 {
			copy(mac[:], prefix)
		}

		oct5 := uint8(mac[5])
		if oct5+uint8(count) > 0xFE {
			oct5 = 0xFE - uint8(count)
		}

		for i := 0; i < count; i++ {
			res = append(res, types.StringValue(
				fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], oct5+uint8(i)),
			))
		}

		m.Mac, diags = types.ListValue(types.StringType, res)
		return
	}

	var uniq = map[string]struct{}{}
	for len(uniq) < count {
		mac, err := getMAC()
		if err != nil {
			diags.AddError(err.Error(), "getMAC() error")
			return
		}
		// len(vendor) 2 .. 5
		if len(prefix) != 0 {
			copy(mac[:], prefix)
		}

		// https://github.com/c-seeger/mac-gen-go/blob/465118e656da9ce72ed2e17ae59e2227a494215f/mac.go#L64
		if len(suffix) > 0 {
			if m.LikeIP.ValueBool() {
				copy(mac[3:], suffix[:])
			} else {
				var n = int(suffix[0]+1)*int(suffix[1]+1)*int(suffix[2]+1) - 1

				mac[3] = uint8(n >> 16)
				mac[4] = uint8(n >> 8)
				mac[5] = uint8(n)
			}
		}

		uniq[fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
			mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])] = struct{}{}
	}

	for k := range uniq {
		res = append(res, types.StringValue(k))
	}

	m.Mac, _ = types.ListValue(types.StringType, res)
	return
}

func (r *macAddressResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state macAddressResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Make a MAC address only if it has not been previously created.
	if state.Mac.IsNull() {
		resp.Diagnostics.Append(state.MakeAddresses()...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *macAddressResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resurcePlan macAddressResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &resurcePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resurcePlan.MakeAddresses()...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &resurcePlan)...)
}

func (r *macAddressResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resurcePlan macAddressResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &resurcePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resurcePlan.MakeAddresses()...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &resurcePlan)...)
}

func (r *macAddressResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
