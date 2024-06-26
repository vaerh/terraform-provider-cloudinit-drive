## terraform-provider-cloudinit-drive
This is a provider for making Cloud-Init drives. Developed for use with Proxmox VE.
Publication in the Hashicorp registry is not planned!

You can find a list of attributes on the [documentation page](docs/resources/cloudinit-drive.md) + [MAC](docs/resources/mac.md).

[Examples](docs/resources/cloudinit-drive.md#examples) of HCL configuration files are in the directory with the test [files](cid/tests).

Preparation of Porxmox VE
``` bash
# Add SSH user
useradd -m cid
install -o cid -g cid -m 700 -d /home/cid/.ssh
# Add SSH public key
cat << EOF >> /home/cid/.ssh/authorized_keys
ecdsa-sha2-nistp521 AAAA...
EOF
chown cid:cid /home/cid/.ssh/authorized_keys
chmod 600 /home/cid/.ssh/authorized_keys
apt install acl
# Locate the directory containing the VM disk images and add extended +rwx permissions.
setfacl -m u:cid:rwx /mnt/ocfs/images
```

Due to the unstable operation of the module for creating ISO9660 it is necessary to install one external program for creating ISO images:

* genisoimage 
* mkisofs 
* hdiutil 
* oscdimg 
* xorriso

Or specify `iso_maker = "none"` to get the output directory `drive_path + '/cid-raw/'` containing a set of Cloud-Init files.


Terraform configuration (FreeBSD VM):
``` hcl
resource "proxmox_vm_qemu" "test-vm" {
  name                   = "test-vm"
  target_node            = "pve"
  clone                  = "test-vm-template"
  full_clone             = true
  vmid                   = 1000
  define_connection_info = false
  agent                  = 1
  automatic_reboot       = true
  balloon                = 1
  qemu_os                = "other"
  cores                  = 1
  hotplug                = "network,disk,usb"
  numa                   = false

  network {
    bridge    = "internal"
    firewall  = false
    link_down = false
    model     = "e1000"
    tag       = 1000
  }

  disk {
    type    = "ide"
    media   = "cdrom"
    size    = cloudinit-drive.vm-test-cloudinit.size
    storage = "Storage"
    file    = "${var.vm_id}/${cloudinit-drive.vm-test-cloudinit.drive_name}"
    #    slot = 2
    volume  = "Storage:${var.vm_id}/${cloudinit-drive.vm-test-cloudinit.drive_name}"
  }

  depends_on = [
    cloudinit-drive.vm-test-cloudinit
  ]
}

resource "cloudinit-drive" "vm-test-cloudinit" {
  drive_name = "vm-${var.vm_id}-cloud-init.raw"
  #  drive_path = "file://./"
  drive_path = "ssh:///mnt/<storage path>/images/${var.vm_id}"
  drive_type = "nocloud"
  iso_maker  = "genisoimage"

  hostname = "vm${var.vm_id}-test"

  network_v2 {
    ethernets {
      alias = "em0"
      dhcp4 = true
    }
  }
}
```
