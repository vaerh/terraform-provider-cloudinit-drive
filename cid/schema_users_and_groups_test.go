package cid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAttributesUserAndGroups(t *testing.T) {
	isoSuffix, _ := uuid.NewUUID()
	attribute := "users_and_groups"

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
						"groups",
						"tests/"+attribute+".yml",
					),
					ImageUserdataEqual(
						isoSuffix.String(),
						"users",
						"tests/"+attribute+".yml",
					),
				),
			},
		},
	})
}
