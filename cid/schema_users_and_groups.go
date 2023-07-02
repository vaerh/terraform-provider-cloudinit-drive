package cid

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var Groups = schema.MapAttribute{
	Optional: true,
	MarkdownDescription: `Groups to add to the system can be specified as a string list.  
	Each item in the list should either contain a string of a single group to create, or a dictionary with the 
	group name as the key and string of a single user as a member of that group or a list of users who should be 
	members of the group.  [Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#users-and-groups)`,
	ElementType: types.ListType{
		ElemType: types.StringType,
	},
}

type UsersType struct {
	Name              types.String `tfsdk:"name"`
	Default           types.Bool   `tfsdk:"default"`
	Expiredate        types.String `tfsdk:"expiredate"`
	Gecos             types.String `tfsdk:"gecos"`
	Groups            types.String `tfsdk:"groups"`
	Homedir           types.String `tfsdk:"homedir"`
	Inactive          types.String `tfsdk:"inactive"`
	LockPasswd        types.Bool   `tfsdk:"lock_passwd"`
	NoCreateHome      types.Bool   `tfsdk:"no_create_home"`
	NoLogInit         types.Bool   `tfsdk:"no_log_init"`
	NoUserGroup       types.Bool   `tfsdk:"no_user_group"`
	Passwd            types.String `tfsdk:"passwd"`
	HashedPasswd      types.String `tfsdk:"hashed_passwd"`
	PlainTextPasswd   types.String `tfsdk:"plain_text_passwd"`
	CreateGroups      types.Bool   `tfsdk:"create_groups"`
	PrimaryGroup      types.String `tfsdk:"primary_group"`
	SelinuxUser       types.String `tfsdk:"selinux_user"`
	Shell             types.String `tfsdk:"shell"`
	Snapuser          types.String `tfsdk:"snapuser"`
	SshAuthorizedKeys types.List   `tfsdk:"ssh_authorized_keys"`
	SshImportID       types.List   `tfsdk:"ssh_import_id"`
	SshRedirectUser   types.Bool   `tfsdk:"ssh_redirect_user"`
	System            types.Bool   `tfsdk:"system"`
	Sudo              types.String `tfsdk:"sudo"`
	Uid               types.Int64  `tfsdk:"uid"`
}

var Users = schema.ListNestedBlock{
	MarkdownDescription: `Users to add to the system.  
	[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#users-and-groups)`,

	NestedObject: schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The user's login name. Required otherwise user creation will be skipped for " +
					"this user.",
			},
			"default": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Add the `default_user` from /etc/cloud/cloud.cfg.",
			},
			"expiredate": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Date on which the user's account will be disabled.",
			},
			"gecos": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Comment about the user, usually a comma-separated string of real name and " +
					"contact information.",
			},
			"groups": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Comma-separated string of groups to add the user to.",
			},
			"homedir": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Home dir for user. `Default: /home/<username>.`",
			},
			"inactive": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "String representing the number of days until the user is disabled.",
			},
			"lock_passwd": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "<span style=\"color:red\">Disable password login. **Default: true.**</span>",
			},
			"no_create_home": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Do not create home directory. **Default: false.**",
			},
			"no_log_init": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Do not initialize lastlog and faillog for user. **Default: false.**",
			},
			"no_user_group": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Do not create group named after user. **Default: false.**",
			},
			"passwd": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				MarkdownDescription: "Hash of user password applied when user does not exist. This will NOT be " +
					"applied if the user already exists. To generate this hash, run: `mkpasswd -method=SHA-512 " +
					"-rounds=4096`.\n" +
					" > *Note:*\nWhile hashed password is better than plain text, using passwd in user-data represents a " +
					"security risk as user-data could be accessible by third-parties depending on your cloud platform.",
			},
			"hashed_passwd": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				MarkdownDescription: "Hash of user password to be applied. This will be applied even if the user " +
					"is pre-existing. To generate this hash, run: `mkpasswd -method=SHA-512 -rounds=4096`.\n" +
					" > *Note:*\nWhile hashed_password is better than plain_text_passwd, using passwd in user-data " +
					"represents a security risk as user-data could be accessible by third-parties depending on your " +
					"cloud platform.",
			},
			"plain_text_passwd": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				MarkdownDescription: "Clear text of user password to be applied. This will be applied even if the " +
					"user is pre-existing. There are many more secure options than using plain text passwords, such " +
					"as ssh_import_id or hashed_passwd. Do not use this in production as user-data and your password " +
					"can be exposed.\n" +
					"    >*<span style=\"color:red\">Security Notice:</span>*\n" +
					"    The password will be stored *unencrypted* in your Terraform state file.\n" +
					"    **Use of this attribute for production deployments is *not* recommended!**",
			},
			"create_groups": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Boolean set false to disable creation of specified user groups. **Default: true.**",
			},
			"primary_group": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Primary group for user. `Default: <username>.`",
			},
			"selinux_user": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SELinux user for user's login. Default to default SELinux user.",
			},
			"shell": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Path to the user's login shell. The default is to set no shell, which results " +
					"in a system-specific default being used.",
			},
			"snapuser": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Specify an email address to create the user as a Snappy user through snap " +
					"create-user. If an Ubuntu SSO account is associated with the address, username and SSH keys " +
					"will be requested from there.",
			},
			"ssh_authorized_keys": schema.ListAttribute{
				Optional: true,
				MarkdownDescription: "List of SSH keys to add to user's authkeys file. *Can not be combined with " +
					"ssh_redirect_user.*",
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("users").AtAnyListIndex().AtName("ssh_redirect_user")),
				},
			},
			"ssh_import_id": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "List of SSH IDs to import for user. *Can not be combined with ssh_redirect_user.*",
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("users").AtAnyListIndex().AtName("ssh_redirect_user")),
				},
			},
			"ssh_redirect_user": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Boolean set to true to disable SSH logins for this user. When specified, all " +
					"cloud meta-data public SSH keys will be set up in a disabled state for this username. Any SSH " +
					"login as this username will timeout and prompt with a message to login instead as the " +
					"`default_username` for this instance. **Default: **false. This key can not be combined with " +
					"`ssh_import_id` or `ssh_authorized_keys`.",
			},
			"system": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Create user as system user with no home directory. **Default: false.**",
			},
			"sudo": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Sudo rule to use or false. Absence of a sudo value or `null` will result in no " +
					"sudo rules added for this user.",
			},
			"uid": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The user's ID. Default is next available value.",
			},
		},
	},
}
