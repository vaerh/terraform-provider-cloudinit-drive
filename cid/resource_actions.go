package cid

import (
	"bytes"
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	ospath "path"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (r *cloudInitDriveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading CloudInit-Drive resource")

	var state cloudInitDriveResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if strings.HasPrefix(state.DrivePath.ValueString(), "ssh://") {
		resp.Diagnostics.Append(r.client.Connect())
		if resp.Diagnostics.HasError() {
			return
		}
		defer r.client.Close()
	}

	diags = r.ReadCloudInitDrive(ctx, &state)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		// This is a hack to force re-create the resource for the field being computed.
		state.DriveType = types.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *cloudInitDriveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating CloudInit-Drive resource")

	var resurcePlan cloudInitDriveResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &resurcePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.CreateCloudInitDrive(ctx, &resurcePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &resurcePlan)...)
}

func (r *cloudInitDriveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating CloudInit-Drive")

	var resurcePlan cloudInitDriveResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &resurcePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.CreateCloudInitDrive(ctx, &resurcePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &resurcePlan)...)
}

func (r *cloudInitDriveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Removing CloudInit-Drive resource from state")

	var state cloudInitDriveResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// This does not take into consideration cross-platform for where the disk is created and where it is copied to!
	isoFile, err := url.JoinPath(state.DrivePath.ValueString(), state.DriveName.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diag.NewAttributeWarningDiagnostic(
			path.Empty(), err.Error(), fmt.Sprintf("Error deleting the '%v' file.", isoFile)))
	}

	switch {
	case strings.HasPrefix(state.DrivePath.ValueString(), "ssh://"):
		resp.Diagnostics.Append(r.client.Connect())
		if resp.Diagnostics.HasError() {
			return
		}
		defer r.client.Close()

		err = r.client.scp.Remove(isoFile[6:])
	case strings.HasPrefix(state.DrivePath.ValueString(), "file://"):
		err = os.Remove(isoFile[7:])
	default:
		resp.Diagnostics.Append(diag.NewAttributeWarningDiagnostic(
			path.Empty(), "Unknown URI scheme: "+state.DrivePath.ValueString(), fmt.Sprintf("Error deleting the '%v' file.", isoFile)))
		return
	}

	if err != nil {
		resp.Diagnostics.Append(diag.NewAttributeWarningDiagnostic(
			path.Empty(), err.Error(), fmt.Sprintf("Error deleting the '%v' file.", isoFile)))
	}

}

//	func (r *cloudInitDriveResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
//		resource.ImportStatePassthroughID(ctx, path.Root("instance_id"), req, resp)
//	}
//
// FIXME
func (r *cloudInitDriveResource) ReadCloudInitDrive(ctx context.Context, stateData *cloudInitDriveResourceModel) (diags diag.Diagnostics) {
	isoFile, err := url.JoinPath(stateData.DrivePath.ValueString(), stateData.DriveName.ValueString())
	if err != nil {
		diags.AddAttributeWarning(path.Empty(), err.Error(),
			"Error reading 'iso' file. The 'iso' file will be re-created.")
		return
	}

	img, err := DriveOpen(r.client, isoFile, false)
	if err != nil {
		diags.AddAttributeWarning(path.Empty(), err.Error(),
			"Error reading 'iso' file. The 'iso' file will be re-created.")
		return
	}
	defer img.Close()

	// Checksum.

	hash := sha256.New()
	_, err = io.Copy(hash, img)

	if err != nil {
		diags.AddAttributeWarning(path.Empty(), err.Error(),
			"Checksum calculation error. The 'iso' file will be re-created.")
		return
	}

	if !stateData.Checksum.Equal(types.StringValue(hex.EncodeToString(hash.Sum(nil)))) {
		diags.AddAttributeWarning(path.Root("checksum"), "The checksum of the ISO image has changed.",
			"The checksum of the ISO file does not match the state value. The 'iso' file will be re-created.")
		// tflog.Info(ctx, "The checksum of the ISO file does not match the state value. The 'iso' file will be re-created.")
		return
	}

	return
}

func (r *cloudInitDriveResource) CreateCloudInitDrive(ctx context.Context, resourcePlan *cloudInitDriveResourceModel) (diags diag.Diagnostics) {
	var files = make(map[string]string)
	var cdLabel = "unknown"

	if strings.HasPrefix(resourcePlan.DrivePath.ValueString(), "ssh://") {
		diags.Append(r.client.Connect())
		if diags.HasError() {
			return
		}
		defer r.client.Close()
	}

	// The attribute is mandatory and always defined.
	switch ParseCiDriveType(resourcePlan.DriveType.ValueString()) {
	case ConfigDrive2:
		cdLabel = "config-2"
		files["meta-data"] = "/openstack/latest/meta_data.json"
		files["user-data"] = "/openstack/latest/user_data"
		files["vendor-data"] = "/openstack/latest/vendor_data.json"
		files["network-data"] = "/content/0000"
	case NoCloud:
		cdLabel = "CIDATA"
		files["meta-data"] = "/meta-data"
		files["user-data"] = "/user-data"
		files["vendor-data"] = "/vendor-data"
		files["network-data"] = "/network-config"
	case OpenNebula:
		cdLabel = "CONTEXT"
		//TODO OpenNebula https://github.com/proxmox/qemu-server/blob/1b5706cd168fedc5e494e24300069ee4ff25761f/PVE/QemuServer/Cloudinit.pm#L255
		// https://cloudinit.readthedocs.io/en/latest/reference/datasources/opennebula.html
		// files: opennebula_context
	}

	// New ISO writer.
	iso, err := NewISOWriter(resourcePlan.ISOMaker.ValueString())

	if err != nil {
		diags.AddAttributeError(path.Empty(), err.Error(), "Error creating 'iso' file.")
		return
	}

	// ----- Metadata -----
	var md attr.Value

	if resourcePlan.CustomFiles != nil {
		md = resourcePlan.CustomFiles.MetaData
	}

	if fileName, ok := GetAttribute(md); ok {
		file, err := os.Open(fileName)

		if err != nil {
			diags.AddAttributeError(path.Root("custom_files").AtListIndex(0).AtName("meta_data"), err.Error(),
				"Error when reading a custom Metadata file.")
			return
		}
		diags.AddAttributeWarning(path.Root("custom_files").AtListIndex(0).AtName("meta_data"),
			"Overwrite settings warning",
			"All settings in child sections of the resource will be overwritten by the contents of this file")

		iso.AddFile(file, files["meta-data"])
	} else {

		metadata, d := MakeMetadata(ctx, resourcePlan)
		diags.Append(d...)
		if diags.HasError() {
			return
		}

		iso.AddFile(bytes.NewBuffer(metadata), files["meta-data"])
	}

	// ----- User data -----
	var ud attr.Value

	if resourcePlan.CustomFiles != nil {
		ud = resourcePlan.CustomFiles.UserData
	}

	if fileName, ok := GetAttribute(ud); ok {
		file, err := os.Open(fileName)

		if err != nil {
			diags.AddAttributeError(path.Root("custom_files").AtListIndex(0).AtName("user_data"), err.Error(),
				"Error when reading a custom User data file.")
			return
		}
		diags.AddAttributeWarning(path.Root("custom_files").AtListIndex(0).AtName("user_data"),
			"Overwrite settings warning",
			"All settings in child sections of the resource will be overwritten by the contents of this file")

		iso.AddFile(file, files["user-data"])
	} else {

		userData, d := MakeUserdata(ctx, resourcePlan)
		diags.Append(d...)
		if diags.HasError() {
			return
		}

		iso.AddFile(bytes.NewReader(userData), files["user-data"])
	}

	// ----- Network data -----
	var nd attr.Value

	if resourcePlan.CustomFiles != nil {
		nd = resourcePlan.CustomFiles.NetworkData
	}

	if fileName, ok := GetAttribute(nd); ok {
		file, err := os.Open(fileName)

		if err != nil {
			diags.AddAttributeError(path.Root("custom_files").AtListIndex(0).AtName("network_data"), err.Error(),
				"Error when reading a custom Network configuration file.")
			return
		}
		diags.AddAttributeWarning(path.Root("custom_files").AtListIndex(0).AtName("network_data"),
			"Overwrite settings warning",
			"All settings in child sections of the resource will be overwritten by the contents of this file")

		iso.AddFile(file, files["network-data"])
	} else {
		// Dublicates ExactlyOneOf check.
		// if resourcePlan.NetConfV1 == nil && resourcePlan.NetConfV2 == nil {
		// 	diags.Append(diag.NewErrorDiagnostic("The network settings block is not found in the configuration.",
		// 		"You must specify a file with network settings or an empty 'network_v1 {} or network_v2 {}' block."))
		// 	return
		// }

		var netConfig []byte
		var d diag.Diagnostics

		if resourcePlan.NetConfV1 != nil {
			netConfig, d = MakeNetConfigV1(ctx, resourcePlan)
			if diags.Append(d...); diags.HasError() {
				return
			}
		} else {
			netConfig, d = MakeNetConfigV2(ctx, resourcePlan)
			if diags.Append(d...); diags.HasError() {
				return
			}
		}

		iso.AddFile(bytes.NewReader(netConfig), files["network-data"])
	}

	// ----- Vendor data -----
	var vd attr.Value

	if resourcePlan.CustomFiles != nil {
		vd = resourcePlan.CustomFiles.VendorData
	}

	if fileName, ok := GetAttribute(vd); ok {
		file, err := os.Open(fileName)

		if err != nil {
			diags.AddAttributeError(path.Root("custom_files").AtListIndex(0).AtName("vendor_data"), err.Error(),
				"Error when reading a Vendor data file.")
			return
		}

		iso.AddFile(file, files["vendor-data"])
	}

	// ----- Custom files -----
	diags.Append(AddCustomFiles(resourcePlan, iso)...)
	if diags.HasError() {
		return
	}

	// ----- Make ISO file -----

	isoFile, err := url.JoinPath(resourcePlan.DrivePath.ValueString(), resourcePlan.DriveName.ValueString())
	if err != nil {
		diags.AddAttributeError(path.Empty(), err.Error(), "Error parsing the file path url")
		return
	}

	img, err := DriveOpen(r.client, isoFile, true)
	if err != nil {
		diags.AddAttributeError(path.Empty(), err.Error(), "Error creating 'iso' file.")
		return
	}
	defer img.Close()

	// ----- Hash calculation & writing -----

	hash := sha256.New()
	hashSha1 := sha1.New()
	buf := new(bytes.Buffer)
	w := io.MultiWriter(img, hash, hashSha1, buf)

	iso.SetLabel(cdLabel)
	_, err = iso.WriteTo(w)
	if err != nil {
		diags.AddAttributeError(path.Empty(), err.Error(), "Error writing 'iso' file.")
		return
	}

	resourcePlan.Checksum = types.StringValue(hex.EncodeToString(hash.Sum(nil)))
	resourcePlan.ID = types.StringValue(hex.EncodeToString(hashSha1.Sum(nil)))

	// Disk image size
	resourcePlan.Size = types.StringValue(fmt.Sprintf("%.0fK", float64(buf.Len())/1024.0))
	if float64(buf.Len())/(1<<20) > 4 {
		diags.AddWarning("Cloud-init drive size greater than 4 megabytes.",
			"The standard drive generated by Proxmox VE should not exceed 4 megabytes.")
	}

	// Cloud-init drive name
	if strings.Contains(resourcePlan.DriveName.ValueString(), "cloudinit") {
		diags.AddAttributeWarning(path.Root("drive_name"), "incorrect drive name",
			"If you are using Proxmox VE, the drive name must not contain the word 'cloudinit' or the VM startup error will occur.")
	}

	return
}

func AddCustomFiles(resourcePlan *cloudInitDriveResourceModel, iso *Iso) diag.Diagnostics {
	if resourcePlan.CustomFiles == nil {
		return nil
	}

	var diags diag.Diagnostics

	for _, scripts := range []struct {
		files    types.Set
		attrName string
		errDescr string
		dstPath  string
	}{
		{resourcePlan.CustomFiles.ScriptsPerBoot, "scripts_per_boot", "script", "/scripts/per-boot/"},
		{resourcePlan.CustomFiles.ScriptsPerInstance, "scripts_per_instance", "script", "/scripts/per-instance/"},
		{resourcePlan.CustomFiles.ScriptsPerOnce, "scripts_per_once", "script", "/scripts/per-once/"},
	} {
		for i, fileName := range scripts.files.Elements() {
			if fileName.IsNull() || fileName.IsUnknown() {
				continue
			}

			file, err := os.Open(fileName.(basetypes.StringValue).ValueString())
			if err != nil {
				diags.AddAttributeError(path.Root("custom_files").AtListIndex(0).AtName(scripts.attrName).AtListIndex(i),
					err.Error(), "Error when reading the file: "+scripts.errDescr)
				continue
			}

			iso.AddFile(file, scripts.dstPath+ospath.Base(fileName.(basetypes.StringValue).ValueString()))
		}
	}

	for i, file := range resourcePlan.CustomFiles.Files {
		if file.Src.IsNull() || file.Dst.IsNull() || file.Src.IsUnknown() || file.Dst.IsUnknown() {
			continue
		}

		f, err := os.Open(file.Src.ValueString())
		if err != nil {
			diags.AddAttributeError(path.Root("custom_files").AtListIndex(0).AtName("files").AtListIndex(i),
				err.Error(), "Error when reading the file: "+file.Src.ValueString())
			continue
		}

		iso.AddFile(f, file.Dst.ValueString())
	}

	return diags
}
