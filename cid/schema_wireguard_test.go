package cid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAttributeWireguarg(t *testing.T) {
	isoSuffix, _ := uuid.NewUUID()
	attribute := "wireguard"

	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: MakeResourceConfig(isoSuffix.String(), "tests/"+attribute+".tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// CopyIsoForDebug(isoSuffix.String()),
					ImageUserdataEqual(
						isoSuffix.String(),
						"readinessprobe",
						"tests/"+attribute+".yml",
					),
					ImageUserdataEqual(
						isoSuffix.String(),
						attribute,
						"tests/"+attribute+".yml",
					),
				),
			},
		},
	})
}
