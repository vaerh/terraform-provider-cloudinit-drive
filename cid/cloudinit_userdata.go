package cid

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"regexp"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var reTrimSpaces = regexp.MustCompile(`(?m)^\s+`)

// https://cloudinit.readthedocs.io/en/latest/topics/examples.html
func MakeUserdata(ctx context.Context, planData *cloudInitDriveResourceModel) ([]byte, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Custom Files
	if planData.CustomFiles != nil {
		if file, ok := GetAttribute(planData.CustomFiles.UserData); ok {
			b, err := os.ReadFile(file)

			if err != nil {
				diags.AddAttributeError(path.Root("custom_files").AtListIndex(0).AtName("user_data"), err.Error(),
					"Error when reading a custom User Data file.")
				return nil, diags
			}

			diags.AddAttributeWarning(path.Root("custom_files").AtListIndex(0).AtName("user_data"),
				"Overwrite user data file.",
				"When using the 'user_data' option, all resource settings will be replaced by the contents of this file.")
			return b, diags
		}
	}

	// CA Certs
	var caSet []attr.Value

	if planData.CaCerts != nil && !planData.CaCerts.FileName.IsNull() && !planData.CaCerts.FileName.IsUnknown() {

		for i, elem := range planData.CaCerts.FileName.Elements() {
			var cert = elem.(basetypes.StringValue).ValueString() // Path

			// Load from file.
			b, err := os.ReadFile(cert)
			if err != nil {
				diags.AddAttributeError(path.Root("ca_certs").AtListIndex(0).AtName("filename"), err.Error(),
					"Error when reading the x509 certificate file.")
				return nil, diags
			}
			ColorizedMessage(ctx, "INFO", fmt.Sprintf("Load trusted CA certificate #%v from %v", i, cert))

			if pemBlock, _ := pem.Decode(b); pemBlock != nil {
				caSet = append(caSet, types.StringValue(string(b)))
			} else {
				diags.AddAttributeError(
					path.Root("ca_certs").AtListIndex(0).AtName("filename"),
					fmt.Sprintf("Failed to parse trusted CA certificate: %v", cert), "")
				return nil, diags
			}
		}
	}

	if planData.CaCerts != nil && !planData.CaCerts.Content.IsNull() && !planData.CaCerts.Content.IsUnknown() {
		for _, elem := range planData.CaCerts.Content.Elements() {
			var cert = elem.(basetypes.StringValue).ValueString() // PEM

			b := []byte(reTrimSpaces.ReplaceAllString(cert, ""))

			if pemBlock, _ := pem.Decode(b); pemBlock != nil {
				caSet = append(caSet, types.StringValue(string(b)))
			} else {
				diags.AddAttributeError(
					path.Root("ca_certs").AtListIndex(0).AtName("content"),
					fmt.Sprintf("Failed to parse trusted CA certificate:\n%v", cert), "")
				return nil, diags
			}
		}
	}

	if len(caSet) > 0 {
		planData.CaCerts.Trusted, diags = types.SetValue(types.StringType, caSet)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Write files
	for i, block := range planData.WriteFiles {
		if !block.LocalFile.IsNull() && !block.LocalFile.IsUnknown() {

			b, err := os.ReadFile(block.LocalFile.ValueString())

			if err != nil {
				diags.AddAttributeError(path.Root("write_files").AtListIndex(i).AtName("local_file"), err.Error(),
					"Error when reading an arbitrary file.")
				return nil, diags
			}

			switch block.Encoding.ValueString() {
			case "gz", "gzip":
				var buf bytes.Buffer
				gz, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)

				_, err := gz.Write(b)
				if err != nil {
					diags.AddAttributeError(path.Root("write_files").AtListIndex(i).AtName("encoding"), err.Error(),
						"File encoding error.")
					return nil, diags
				}

				gz.Close()
				planData.WriteFiles[i].Content = types.StringValue(buf.String())

			case "gz+base64", "gzip+base64", "gz+b64", "gzip+b64":
				var buf bytes.Buffer
				b64 := base64.NewEncoder(base64.StdEncoding, &buf)
				gz, _ := gzip.NewWriterLevel(b64, gzip.BestCompression)

				_, err := gz.Write(b)
				if err != nil {
					diags.AddAttributeError(path.Root("write_files").AtListIndex(i).AtName("encoding"), err.Error(),
						"File encoding error.")
					return nil, diags
				}

				gz.Close()
				planData.WriteFiles[i].FileData = types.StringValue(SplitStringForYaml(buf.String()))

			case "b64", "base64":
				planData.WriteFiles[i].FileData = types.StringValue(SplitStringForYaml(base64.StdEncoding.EncodeToString(b)))

			default: // "text/plain"
				planData.WriteFiles[i].FileData = types.StringValue(string(b))
			}
		} else {
			planData.WriteFiles[i].FileData = planData.WriteFiles[i].Content
		}
	}

	buf, err := yaml.MarshalWithOptions(TerraformToGo(planData),
		yaml.UseLiteralStyleIfMultiline(true),
		yaml.UseSingleQuote(true),
	)
	if err != nil {
		diags.AddAttributeError(path.Empty(), err.Error(), "yaml.Marshal error")
		return nil, diags
	}

	return append([]byte("#cloud-config\n"), buf...), diags
}
