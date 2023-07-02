package cid

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAttributesCustomFiles(t *testing.T) {
	isoSuffix, _ := uuid.NewUUID()
	attribute := "custom_files"

	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: LoadConfig(isoSuffix.String(), "tests/"+attribute+".tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					CopyIsoForDebug(isoSuffix.String()),
				),
			},
		},
	})
}

func LoadConfig(suffix, attributeFileName string) string {
	f, err := os.Open(attributeFileName)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(string(b), suffix)
}
