package cid

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/kdomanski/iso9660"
)

func LoadEnvVariable(v any, env string) any {

	switch v := v.(type) {
	case basetypes.StringValue:
		if !v.IsNull() {
			return v
		}

		return types.StringValue(os.Getenv(env))

	case basetypes.BoolValue:
		if !v.IsNull() {
			return v
		}

		if _, ok := os.LookupEnv(env); ok {
			return types.BoolValue(true)
		}

		return types.BoolValue(false)

	case basetypes.Int64Value:
		if !v.IsNull() {
			return v
		}

		i, err := strconv.Atoi(os.Getenv(env))
		if err != nil {
			return types.Int64Null()
		}

		return types.Int64Value(int64(i))
	}

	return nil
}

type Attribute interface {
	IsNull() bool
	IsUnknown() bool
	Type(_ context.Context) attr.Type
}

func GetAttribute(attr Attribute) (s string, known bool) {
	if attr == nil || attr.IsNull() || attr.IsUnknown() {
		return "", false
	}

	switch attr.Type(context.TODO()) {
	case basetypes.StringType{}:
		return attr.(basetypes.StringValue).ValueString(), true

	case basetypes.BoolType{}:
		return strconv.FormatBool(attr.(basetypes.BoolValue).ValueBool()), true

	default:
		panic("Not implemented yet.")
	}
}

type UniversalFile interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
	io.ReaderAt
}

func DriveOpen(c *SSHClient, uri string, rw bool) (UniversalFile, error) {
	switch {
	case strings.HasPrefix(uri, "ssh://"):
		if rw {

			path := path.Dir(uri[6:])

			if path != "." {
				if err := c.scp.MkdirAll(path); err != nil {
					return nil, err
				}
			}

			return c.RemoteWrite(uri[6:])
		}

		return c.RemoteRead(uri[6:])

	case strings.HasPrefix(uri, "file://"):
		if rw {

			path := path.Dir(uri[7:])

			if path != "." {
				if err := os.MkdirAll(path, os.ModePerm); err != nil {
					return nil, err
				}
			}

			return os.OpenFile(uri[7:], os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		}

		if _, err := os.Stat(uri[7:]); err != nil {
			return nil, err
		}

		return os.OpenFile(uri[7:], os.O_RDONLY, 0)
	}

	if i := strings.Index(uri, "://"); i > -1 {
		uri = uri[:i]
	}

	return nil, errors.New("unknown URI scheme: " + uri)
}

func GetIsoFile(f *iso9660.File, targetFile string) ([]byte, error) {
	b, _, e := isoSearch(f, targetFile, "/")
	return b, e
}

func isoSearch(f *iso9660.File, targetFile, currentPath string) ([]byte, bool, error) {
	var b []byte
	var found bool

	if f.IsDir() {
		children, err := f.GetChildren()

		if err != nil {
			return nil, found, err
		}

		for _, c := range children {
			if b, found, err = isoSearch(c, targetFile, path.Join(currentPath, c.Name())); err != nil {
				return nil, found, err
			} else {

				if found {
					return b, found, nil
				}

			}
		}
	} else if targetFile == currentPath { // it's a file
		var buf bytes.Buffer

		if _, err := buf.ReadFrom(f.Reader()); err != nil {
			return nil, found, err
		}

		return buf.Bytes(), true, nil
	}

	return b, found, nil
}

type fileValidator struct{}

func (f fileValidator) Description(ctx context.Context) string {
	return "the file from which the information will be read must be present in the system"
}

func (f fileValidator) MarkdownDescription(ctx context.Context) string {
	return f.Description(ctx)
}

func (f fileValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	_, err := os.Stat(req.ConfigValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path.AtName(req.Path.String()),
			"File stat error",
			fmt.Sprintf("Error %v", err),
		)
	}
}

func SplitStringForYaml(s string) (res string) {
	const lineLen = 70

	if len(s) <= lineLen {
		return s
	}

	var i int
	for i = lineLen; i < len(s); i += lineLen {
		res += s[i-lineLen:i] + "\n"
	}
	res += s[i-lineLen:]
	return
}
