groups = {
  admingroup  = ["root", "sys"],
  cloud-users = []
}

users {
  name        = "newsuper"
  gecos       = "Big Stuff"
  groups      = "users, admin"
  sudo        = "ALL=(ALL) NOPASSWD:ALL"
  shell       = "/bin/bash"
  lock_passwd = true
  ssh_import_id = [
    "lp:falcojr",
    "gh:TheRealFalcon"
  ]
}

users {
  name         = "youruser"
  selinux_user = "staff_u"
}