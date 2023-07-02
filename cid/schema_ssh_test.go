package cid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAttributesSsh(t *testing.T) {
	isoSuffix, _ := uuid.NewUUID()
	attribute := "ssh"

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
						"ssh_keys",
						"tests/"+attribute+".yml",
					),
					ImageUserdataEqual(
						isoSuffix.String(),
						"ssh_authorized_keys",
						"tests/"+attribute+".yml",
					),
					ImageUserdataEqual(
						isoSuffix.String(),
						"ssh_deletekeys",
						"tests/"+attribute+".yml",
					),
					ImageUserdataEqual(
						isoSuffix.String(),
						"ssh_genkeytypes",
						"tests/"+attribute+".yml",
					),
					ImageUserdataEqual(
						isoSuffix.String(),
						"disable_root",
						"tests/"+attribute+".yml",
					),
					ImageUserdataEqual(
						isoSuffix.String(),
						"disable_root_opts",
						"tests/"+attribute+".yml",
					),
					ImageUserdataEqual(
						isoSuffix.String(),
						"allow_public_ssh_keys",
						"tests/"+attribute+".yml",
					),
					ImageUserdataEqual(
						isoSuffix.String(),
						"ssh_quiet_keygen",
						"tests/"+attribute+".yml",
					),
					ImageUserdataEqual(
						isoSuffix.String(),
						"ssh_publish_hostkeys",
						"tests/"+attribute+".yml",
					),
				),
			},
		},
	})
}
