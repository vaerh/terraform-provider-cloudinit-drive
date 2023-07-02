package cid

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var ManageEtcHosts = schema.StringAttribute{
	Optional: true,
	Validators: []validator.String{
		stringvalidator.OneOf([]string{"true", "false", "localhost"}...),
	},
	Description: `Update the hosts file ("true", "false", "localhost").  
	If set to true, cloud-init will generate the hosts file using the template located
	in /etc/cloud/templates/hosts.tmpl. In the /etc/cloud/templates/hosts.tmpl template,
	the strings \$hostname and \$fqdn will be replaced with the hostname and fqdn respectively.
	  
	If manage_etc_hosts is set to localhost, then cloud-init will not rewrite the hosts file entirely,
	but rather will ensure that a entry for the fqdn with a distribution dependent ip is present
	(i.e. ping <hostname> will ping 127.0.0.1 or 127.0.1.1 or other ip).
	  
	> *Note:*
	If manage_etc_hosts is set true, the contents of the hosts file will be updated every boot.
	To make any changes to the hosts file persistent they must be made in /etc/cloud/templates/hosts.tmpl  
	[Info](https://cloudinit.readthedocs.io/en/latest/topics/modules.html#update-etc-hosts)
	`,
}
