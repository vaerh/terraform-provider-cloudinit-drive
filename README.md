## terraform-provider-cloudinit-drive
This is a provider for making Cloud-Init drives. Developed for use with Proxmox VE.
Publication in the Hashicorp registry is not planned!

You can find a list of attributes on the [documentation page](docs/resources/cloudinit-drive.md).

Examples of HCL configuration files are in the directory with the test [files](cid/tests).

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

  hostname = "vm${var.vm_id}-test"

  #  network_v1 {
  #    interface {
  #      type = "physical"
  #      name = "eth0"
  #      macaddress = "00:11:22:33:44:55"
  #      subnets {
  #        type = "dhcp"
  #      }
  #      subnets {
  #        type    = "static"
  #        address = "10.184.225.122"
  #        netmask = "255.255.255.252"
  #        routes {
  #          gateway     = "10.184.225.121"
  #          netmask     = "255.240.0.0"
  #          destination = "10.176.0.0"
  #        }
  #        routes {
  #          gateway     = "10.184.225.121"
  #          netmask     = "255.240.0.0"
  #          destination = "10.208.0.0"
  #        }
  #      }
  #    }
  #  }

  network_v2 {
    ethernets {
      alias = "em0"
      dhcp4 = true
    }
  }
}
```
