package cid

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var Locale = schema.StringAttribute{
	Optional: true,
	Description: "The locale to set as the system's locale. " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#locale)",
}

var LocaleConfig = schema.StringAttribute{
	Optional: true,
	Description: "The file in which to write the locale configuration (defaults to the distro's default location). " +
		"[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#locale)",
}
