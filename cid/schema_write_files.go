package cid

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WriteFilesType struct {
	LocalFile   types.String `tfsdk:"local_file" cid:"skip"`
	Path        types.String `tfsdk:"path"`
	Content     types.String `tfsdk:"content" json:"-" yaml:"-"`
	Owner       types.String `tfsdk:"owner"`
	Permissions types.String `tfsdk:"permissions"`
	Encoding    types.String `tfsdk:"encoding"`
	Append      types.Bool   `tfsdk:"append"`
	Defer       types.Bool   `tfsdk:"defer"`
	FileData    types.String `tfsdk:"-" json:"content,omitempty" yaml:"content,omitempty"`
}

var WriteFiles = schema.ListNestedBlock{
	MarkdownDescription: "Write arbitrary files. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#write-files)",

	NestedObject: schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"local_file": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Path of the local file data from which will be written to the attribute 'content'.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("content")),
					fileValidator{},
				},
			},
			"path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Path of the file to which content is decoded and written.",
			},
			"content": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional content to write to the provided path. When content is present and encoding is " +
					"not `text/plain`, decode the content prior to writing. Binaries will be gzip+base64 encoded",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("local_file")),
				},
			},
			"owner": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional owner:group to chown on the file. `Default: root:root`.",
			},
			"permissions": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional file permissions to set on path represented as an octal string `'0###'`. " +
					"`Default: 0644`",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^'0\d{3}'$`),
						"string must be in the format: \"'0###'\"",
					),
				},
			},
			"encoding": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional encoding type of the content. Default is base64 and no content decoding is " +
					"performed. Supported encoding types are: `gz`, `gzip`, `gz+base64`, `gzip+base64`, `gz+b64`, " +
					"`gzip+b64`, `b64`, `base64`, `text/plain`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						[]string{"gz", "gzip", "gz+base64", "gzip+base64", "gz+b64", "gzip+b64", "b64", "base64", "text/plain"}...,
					),
				},
			},
			"append": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to append content to existing file if path exists. `Default: false.`",
			},
			"defer": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Defer writing the file until 'final' stage, after users were created, and packages were " +
					"installed. `Default: false.`",
			},
		},
	},
}
