provider "cloudinit-drive" {
}

resource "cloudinit-drive_mac" "mac-1" {
  prefix = "aabbcedf"
  number = 10
}

resource "cloudinit-drive_mac" "mac-2" {
  prefix = "aabbce"
  number = 10
  sequential = true
}

resource "cloudinit-drive_mac" "mac-3" {
  prefix = "aabbce"
  suffix = "192.168.123.32"
}

resource "cloudinit-drive_mac" "mac-4" {
  prefix = "aabbce"
  suffix = "192.168.123.32"
  like_ip = true
}

resource "cloudinit-drive_mac" "mac-5" {
  number = 10
}

resource "cloudinit-drive_mac" "mac-6" {
  suffix = "192.168.123.32"
}

resource "cloudinit-drive_mac" "mac-7" {
}

output "mac-1" {
  value = [
  resource.cloudinit-drive_mac.mac-1.mac,
  resource.cloudinit-drive_mac.mac-2.mac,
  resource.cloudinit-drive_mac.mac-3.mac,
  resource.cloudinit-drive_mac.mac-4.mac,
  resource.cloudinit-drive_mac.mac-5.mac,
  resource.cloudinit-drive_mac.mac-6.mac,
  resource.cloudinit-drive_mac.mac-7.mac,
  ]
}

