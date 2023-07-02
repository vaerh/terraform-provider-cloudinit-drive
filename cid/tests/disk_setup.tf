device_aliases = {
  my_alias  = "/dev/sdb"
  swap_disk = "/dev/sdc"
}

disk_setup = {
  my_alias = {
    table_type = "gpt"
    layout     = "[50, 50]"
    overwrite  = true
  }
  swap_disk = {
    table_type = "gpt"
    layout     = "[[100, 82]]"
    overwrite  = true
  }
  "/dev/sdd" = {
    table_type = "mbr"
    layout     = "true"
    overwrite  = true
  }
}

fs_setup {
  label      = "fs1"
  filesystem = "ext4"
  device     = "my_alias.1"
  cmd        = "mkfs -t %(filesystem)s -L %(label)s %(device)s"
}

fs_setup {
  label      = "fs2"
  device     = "my_alias.2"
  filesystem = "ext4"
}

fs_setup {
  label      = "swap"
  device     = "swap_disk.1"
  filesystem = "swap"
}

fs_setup {
  label      = "fs3"
  device     = "/dev/sdd1"
  filesystem = "ext4"
}

mounts = [
  ["my_alias.1", "/mnt1"],
  ["my_alias.2", "/mnt2"],
  ["swap_disk.1", "none", "swap", "sw", "'0'", "'0'"],
  ["/dev/sdd1", "/mnt3"],
]