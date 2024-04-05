package cid

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vaerh/iso9660"
)

func init() {
	// os.Setenv("TF_ACC", "1")
	os.Setenv("CID_LOG_COLOR", "1")
}

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"cloudinit-drive": providerserver.NewProtocol6WithError(NewCloudInitDrive()),
	}
)

func TestProvider(t *testing.T) {
	isoSuffix, _ := uuid.NewUUID()
	attribute := "provider"

	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: MakeResourceConfig(isoSuffix.String(), "tests/"+attribute+".tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudinit-drive.vm-test-cloudinit-drive", "hostname", "testhost.fqdn"),

					// ImageFileEqual(
					// 	"vm-101-cloud-init."+isoSuffix.String(),
					// 	"/openstack/latest/meta_data.json",
					// 	`{"uuid":"98208689-b17d-47c7-889f-b0f7ebb06a21","hostname":"testhost.fqdn",`+
					// 		`"network_config":{"content_path":"/openstack/content/0000"}}`,
					// ),
				),
			},
		},
	})
}

func CopyIsoForDebug(isoSuffix string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		data, err := os.ReadFile("vm-101-cloud-init." + isoSuffix)
		if err != nil {
			return err
		}
		return os.WriteFile("vm-101-cloud-init."+isoSuffix+".iso", data, 0644)
	}
}

func ImageUserdataEqual(isoSuffix, attrName, yamlName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		f, err := os.Open("vm-101-cloud-init." + isoSuffix)
		if err != nil {
			return err
		}
		iso, err := iso9660.OpenImage(f)
		if err != nil {
			return err
		}
		root, err := iso.RootDir()
		if err != nil {
			return err
		}
		b, err := getIsoFile(root, "/openstack/latest/user-data")
		if err != nil {
			return err
		}
		f.Close()

		got, want := make(map[string]any), make(map[string]any)
		err = yaml.Unmarshal(b, &got)
		if err != nil {
			return err
		}

		f, err = os.Open(yamlName)
		if err != nil {
			return err
		}

		b, err = io.ReadAll(f)
		if err != nil {
			return err
		}
		f.Close()

		err = yaml.Unmarshal(b, &want)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(got[attrName], want[attrName]) {
			return fmt.Errorf("ImageFileEqual = \n%v,\nwant\n%v", spew.Sdump(got[attrName]), spew.Sdump(want[attrName]))
		}
		return nil
	}
}

func MakeResourceConfig(suffix, attributeFileName string) string {
	f, err := os.Open(attributeFileName)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(`
resource "cloudinit-drive" "vm-test-cloudinit-drive" {
  drive_name = "vm-101-cloud-init.%s"
  drive_path = "file://./"
  drive_type = "configdrive2"
  iso_maker  = "genisoimage"

  %s

   network_v2 {}
}
`, suffix, b)
}

func getIsoFile(f *iso9660.File, targetFile string) ([]byte, error) {
	b, _, e := isoSearch(f, targetFile, "/")
	return b, e
}

func isoSearch(f *iso9660.File, targetFile, currentPath string) ([]byte, bool, error) {
	var b []byte
	var found bool

	if f.IsDir() {
		children, err := f.GetChildren()

		if err != nil {
			return nil, found, err
		}

		for _, c := range children {
			if b, found, err = isoSearch(c, targetFile, path.Join(currentPath, c.Name())); err != nil {
				return nil, found, err
			} else {

				if found {
					return b, found, nil
				}

			}
		}
	} else if targetFile == currentPath { // it's a file
		var buf bytes.Buffer

		if _, err := buf.ReadFrom(f.Reader()); err != nil {
			return nil, found, err
		}

		return buf.Bytes(), true, nil
	}

	return b, found, nil
}
