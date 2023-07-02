package cid

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SshPublishHostkeysType struct {
	Enabled   types.Bool `tfsdk:"enabled"`
	BlackList types.List `tfsdk:"blacklist"`
}

var SshKeys = schema.MapAttribute{
	Optional: true,
	MarkdownDescription: "A dictionary entries for the public and private host keys of each desired key type. " +
		"Entries in the ssh_keys config dict should have keys in the format <key type>_private, <key type>_public, " +
		"and, optionally, <key type>_certificate, e.g. rsa_private: <key>, rsa_public: <key>, and rsa_certificate: " +
		"<key>. Not all key types have to be specified, ones left unspecified will not be used. If this config " +
		"option is used, then separate keys will not be automatically generated. In order to specify multiline " +
		"private host keys and certificates, use multiline syntax.",
	ElementType: types.StringType,
	Validators: []validator.Map{
		mapvalidator.KeysAre(stringvalidator.RegexMatches(
			regexp.MustCompile(`_(private|public|certificate)$`),
			"the keys must contain one of thise suffixes: _private, _public or _certificate",
		)),
	},
}

var SshAuthorizedKeys = schema.SetAttribute{
	Optional: true,
	MarkdownDescription: "The SSH public keys to add .ssh/authorized_keys in the default user's home directory. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#ssh)",
	ElementType: types.StringType,
}

var SshDeleteKeys = schema.BoolAttribute{
	Optional: true,
	Computed: true,
	Default:  booldefault.StaticBool(true),
	MarkdownDescription: "Remove host SSH keys. This prevents re-use of a private host key from an image with " +
		"default host SSH keys. **Default: true.** " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#ssh)",
}

var SshGenKeyTypes = schema.ListAttribute{
	Optional: true,
	Computed: true,
	Default: listdefault.StaticValue(types.ListValueMust(
		types.StringType,
		[]attr.Value{
			types.StringValue("rsd"),
			types.StringValue("dsa"),
			types.StringValue("ecdsa"),
			types.StringValue("ed25519"),
		},
	),
	),
	MarkdownDescription: `The SSH key types to generate. **Default: ["rsa", "dsa", "ecdsa", "ed25519"].** ` +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#ssh)",
	ElementType: types.StringType,
}

var DisableRoot = schema.BoolAttribute{
	Optional: true,
	Computed: true,
	Default:  booldefault.StaticBool(true),
	MarkdownDescription: "Disable root login. **Default: true.** " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#ssh)",
}

var DisableRootOpts = schema.StringAttribute{
	Optional: true,
	MarkdownDescription: `Disable root login options. If disable_root_opts is specified and contains the string $USER, it ` +
		`will be replaced with the username of the default user. *Default: no-port-forwarding,no-agent-forwarding,` +
		`no-X11-forwarding,command="echo 'Please login as the user \"\$USER\" rather than the user \"\$DISABLE_USER\".` +
		`';echo;sleep 10;exit 142"* [Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#ssh).`,
}

var AllowPublicSshKeys = schema.BoolAttribute{
	Optional: true,
	Computed: true,
	Default:  booldefault.StaticBool(true),
	MarkdownDescription: "If true, will import the public SSH keys from the datasource's metadata to the user's " +
		".ssh/authorized_keys file. **Default: true.** " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#ssh).",
}

var SshQuietKeygen = schema.BoolAttribute{
	Optional:            true,
	MarkdownDescription: "If true, will suppress the output of key generation to the console.",
}

var SshPublishHostkeys = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			Optional: true,
			MarkdownDescription: "If true, will read host keys from /etc/ssh/*.pub and publish them to the " +
				"datasource (if supported). Default: true",
		},
		"blacklist": schema.ListAttribute{
			Optional:            true,
			MarkdownDescription: " The SSH key types to ignore when publishing. Default: [dsa]",
			ElementType:         types.StringType,
		},
	},
}
