package cid

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var RunCmd = schema.ListAttribute{
	Optional: true,
	MarkdownDescription: "Run arbitrary commands at a rc.local like time-frame with output to the console. " +
		"Each item can be either a list or a string. The item type affects how it is executed:  " +
		"* If the item is a string, it will be interpreted by sh." +
		"* If the item is a list, the items will be executed as if passed to execve(3) (with the first arg as the command).  " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#runcmd)",
	ElementType: types.StringType,
}
