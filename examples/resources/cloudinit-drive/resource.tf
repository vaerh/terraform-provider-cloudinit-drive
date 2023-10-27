resource "cloudinit_drive" "vm-104-cloudinit" {
  drive_name = "vm-101-cloud-init.iso"
  drive_path = "file://./"
  drive_type = "nocloud"
  iso_maker  = "mkisofs"

  custom_files {
    meta_data        = "configs/meta_data.json"
    network_data     = "configs/network.yaml"
    user_data        = "configs/user_data.yaml"
    vendor_data      = "configs/vendor_data"
    scripts_per_boot = ["scripts/make_static_routes.sh"]
  }

  network_v1 {}
}
