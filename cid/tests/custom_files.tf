resource "cloudinit-drive" "vm-test-cloudinit-drive" {
  drive_name = "vm-101-cloud-init.%s"
  drive_path = "file://./"
  drive_type = "configdrive2"
  iso_maker  = "mkisofs"

  custom_files {
    scripts_per_boot     = ["tests/ca_certs.tf", "tests/ca_certs.yml"]
    scripts_per_instance = ["tests/disk_setup.tf", "tests/disk_setup.yml"]
    scripts_per_once     = ["tests/final_message.tf", "tests/final_message.yml"]
    user_data            = "tests/growpart.tf"
    vendor_data          = "tests/hostname.tf"
    network_data         = "tests/keyboard.tf"
    # opennebula_context = "tests/locale.tf"
    files {
      src = "tests/mounts.tf"
      dst = "/files/1/file-1.tf"
    }
    files {
      src = "tests/mounts.yml"
      dst = "/files/2/file-2.yml"
    }
  }

  network_v2 {}
}