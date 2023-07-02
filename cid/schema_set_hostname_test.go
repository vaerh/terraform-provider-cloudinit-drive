package cid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAttributesHostname(t *testing.T) {
	isoSuffix, _ := uuid.NewUUID()
	attribute := "hostname"

	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: MakeResourceConfig(isoSuffix.String(), "tests/"+attribute+".tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// CopyIsoForDebug(isoSuffix.String()),
					// ImageUserdataEqual(
					// 	isoSuffix.String(),
					// 	"preserve_hostname",
					// 	"tests/"+attribute+".yml",
					// ),
					ImageUserdataEqual(
						isoSuffix.String(),
						attribute,
						"tests/"+attribute+".yml",
					),
					ImageUserdataEqual(
						isoSuffix.String(),
						"fqdn",
						"tests/"+attribute+".yml",
					),
					ImageUserdataEqual(
						isoSuffix.String(),
						"prefer_fqdn_over_hostname",
						"tests/"+attribute+".yml",
					),
				),
			},
		},
	})
}
