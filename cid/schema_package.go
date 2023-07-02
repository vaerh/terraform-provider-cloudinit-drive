package cid

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var Packages = schema.ListAttribute{
	Optional: true,
	MarkdownDescription: "A list of packages to install. *Package version selection is not supported!* " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#package-update-upgrade-install)",
	ElementType: types.StringType,
}

var PackageUpdate = schema.BoolAttribute{
	Optional: true,
	Computed: true,
	Default:  booldefault.StaticBool(false),
	MarkdownDescription: "Set true to update packages. **Default: false.** " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#package-update-upgrade-install)",
}

var PackageUpgrade = schema.BoolAttribute{
	Optional: true,
	Computed: true,
	Default:  booldefault.StaticBool(false),
	MarkdownDescription: " Set true to upgrade packages. **Default: false.** " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#package-update-upgrade-install)",
}

var PackageRebootIfRequired = schema.BoolAttribute{
	Optional: true,
	Computed: true,
	Default:  booldefault.StaticBool(false),
	MarkdownDescription: "Set true to reboot the system if required by presence of /var/run/reboot-required. " +
		"*Default: false.* " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#package-update-upgrade-install)",
}
