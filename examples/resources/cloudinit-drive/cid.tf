resource "cloudinit-drive" "vm-test-cloudinit-drive" {
  drive_type = "configdrive2"
  drive_name = "vm-101-cloud-init.iso"
  drive_path = "file://./"
  iso_maker  = "genisoimage"

  hostname = "testhost.fqdn"

  device_aliases = {
    my_alias  = "/dev/sdb"
    swap_disk = "/dev/sdc"
  }

  fs_setup {
    label      = "fs1"
    filesystem = "ext4"
    device     = "my_alias.1"
    cmd        = "mkfs -t %(filesystem)s -L %(label)s %(device)s"
  }

  # fs_setup {
  #   label      = "fs2"
  #   filesystem = "ext4"
  #   device     = "my_alias.2"
  # }

  mounts = [
    [ "/dev/ephemeral0", "/mnt", "auto", "defaults,noexec" ],
    [ "sdc", "/opt/data" ],
    [ "xvdh", "/opt/data", "auto", "defaults,nofail", "0", "0" ]
  ]

  ca_certs {
    filename = [
      "/home/terraform/cid/root-ca.pem"
    ]
  }

  custom_files {
    # user_data = "delvrun.sh"
    # network_data = "/custom/file/path"
    # scripts_per_boot = ["cid/file1", "cid/nofile"]
    scripts_per_boot = ["cid/file1"]
  }

  write_files {
    encoding   = "gz+b64"
    local_file = "/home/terraform/hello.sh"
    path       = "/tmp/aaa/bbb"
  }

  network_v2 {
    ethernets {
      match {
        macaddress = "00:11:22:33:44:55"
      }
      alias     = "id0"
      wakeonlan = true
      dhcp4     = false
      addresses = ["192.168.14.2/24", "2001:1::1/64"]
      gateway4  = "192.168.14.1"
      gateway6  = "2001:1::2"
      nameservers {
        search    = ["foo.local", "bar.local"]
        addresses = ["8.8.8.8"]
      }
      routes {
        to     = "192.0.2.0/24"
        via    = "11.0.0.1"
        metric = 3
      }
    }
    
    ethernets {
      alias = "lom"
      match {
        driver = "ixgbe"
      }
      set_name = "lom1"
      dhcp6    = true
    }
    
    ethernets {
      alias = "switchports"
      match {
        name = "enp2*"
      }
      mtu = 1280
    }

    bonds {
      alias      = "bond0"
      interfaces = ["id0", "lom"]
    }

    bridges {
      alias      = "br0"
      interfaces = ["wlp1s0", "switchports"]
      dhcp4      = true
    }

    vlans {
      alias = "en-intra"
      id    = 1
      link  = "id0"
      dhcp4 = true
    }
  }
}

output "cid-drive-size" {
  value = cloudinit-drive.vm-test-cloudinit-drive.size
}
