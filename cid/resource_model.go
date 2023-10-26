package cid

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CIDriveType int

const (
	ConfigDrive2 CIDriveType = iota
	NoCloud
	OpenNebula
)

func (t CIDriveType) String() string {
	switch t {
	case ConfigDrive2:
		return "configdrive2"
	case NoCloud:
		return "nocloud"
	case OpenNebula:
		return "opennebula"
	}
	panic("Unexpected ciDriveType (t): " + strconv.Itoa(int(t)))
}

func ParseCiDriveType(s string) CIDriveType {
	switch s {
	case "configdrive2":
		return ConfigDrive2
	case "nocloud":
		return NoCloud
	case "opennebula":
		return OpenNebula
	}
	panic("Unexpected ciDriveType (s): " + s)
}

type cloudInitDriveResourceModel struct {
	// Attributes
	ID                      types.String `tfsdk:"id" cid:"skip"`
	DriveType               types.String `tfsdk:"drive_type" cid:"skip"`
	DrivePath               types.String `tfsdk:"drive_path" cid:"skip"`
	DriveName               types.String `tfsdk:"drive_name" cid:"skip"`
	ISOMaker                types.String `tfsdk:"iso_maker" cid:"skip"`
	Checksum                types.String `tfsdk:"checksum" cid:"skip"`
	Size                    types.String `tfsdk:"size" cid:"skip"`
	InstanceId              types.String `tfsdk:"instance_id" cid:"skip"`
	AllowPublicSshKeys      types.Bool   `tfsdk:"allow_public_ssh_keys"`
	DeviceAliases           types.Map    `tfsdk:"device_aliases"`
	DisableRoot             types.Bool   `tfsdk:"disable_root"`
	DisableRootOpts         types.String `tfsdk:"disable_root_opts"`
	FinalMessage            types.String `tfsdk:"final_message"`
	Fqdn                    types.String `tfsdk:"fqdn"`
	Groups                  types.Map    `tfsdk:"groups"`
	Hostname                types.String `tfsdk:"hostname"`
	PreferFqdnOverHostname  types.Bool   `tfsdk:"prefer_fqdn_over_hostname"`
	PreserveHostname        types.Bool   `tfsdk:"preserve_hostname"`
	Locale                  types.String `tfsdk:"locale"`
	LocaleConfig            types.String `tfsdk:"locale_configfile"`
	ManageEtcHosts          types.String `tfsdk:"manage_etc_hosts"`
	ManageResolvConf        types.Bool   `tfsdk:"manage_resolv_conf"`
	MountDefaultFields      types.List   `tfsdk:"mount_default_fields"`
	Mounts                  types.List   `tfsdk:"mounts"`
	PackageRebootIfRequired types.Bool   `tfsdk:"package_reboot_if_required"`
	PackageUpdate           types.Bool   `tfsdk:"package_update"`
	PackageUpgrade          types.Bool   `tfsdk:"package_upgrade"`
	Packages                types.List   `tfsdk:"packages"`
	Password                types.String `tfsdk:"password"`
	RunCmd                  types.List   `tfsdk:"runcmd"`
	SshAuthorizedKeys       types.Set    `tfsdk:"ssh_authorized_keys"`
	SshDeleteKeys           types.Bool   `tfsdk:"ssh_deletekeys"`
	SshGenKeyTypes          types.List   `tfsdk:"ssh_genkeytypes"`
	SshKeys                 types.Map    `tfsdk:"ssh_keys"`
	SshPwAuth               types.Bool   `tfsdk:"ssh_pwauth"`
	SshQuietKeygen          types.Bool   `tfsdk:"ssh_quiet_keygen"`
	TimeZone                types.String `tfsdk:"timezone"`

	// Blocks
	CaCerts                *CaCertType             `tfsdk:"ca_certs"`
	ChPasswd               *ChPasswdType           `tfsdk:"chpasswd"`
	CustomFiles            *CustomFilesType        `tfsdk:"custom_files" json:"custom_files,omitempty" yaml:"custom_files,omitempty"`
	DiskSetup              DiskSetupType           `tfsdk:"disk_setup"`
	FsSetup                []FsSetupType           `tfsdk:"fs_setup"`
	GrowPart               *GrowPartType           `tfsdk:"growpart"`
	Keyboard               *KeyboardType           `tfsdk:"keyboard"`
	PowerState             *PowerStateType         `tfsdk:"power_state"`
	ResolvConf             *ResolvConfType         `tfsdk:"resolv_conf"`
	SeedRandom             *SeedRandomType         `tfsdk:"random_seed"`
	SshPublishHostkeysType *SshPublishHostkeysType `tfsdk:"ssh_publish_hostkeys"`
	Swap                   *SwapType               `tfsdk:"swap"`
	Users                  []UsersType             `tfsdk:"users"`
	Wireguard              *WireguardType          `tfsdk:"wireguard"`
	WriteFiles             []WriteFilesType        `tfsdk:"write_files"`
	NetConf                *NetConfType            `tfsdk:"network" cid:"skip"`
}
