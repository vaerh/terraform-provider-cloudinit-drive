ssh_pwauth = false
# password = "password1"

chpasswd {
  expire = false
  users {
    name     = "user1"
    password = "password1"
    type     = "text"
  }
  users {
    name     = "user2"
    password = "$6$rounds=4096$5DJ8a9WMTEzIo5J4$Yms6imfeBvf3Yfu84mQBerh18l7OR1Wm1BJXZqFSpJ6BVas0AYJqIjP7czkOaAZHZi1kxQ5Y1IhgWN8K9NgxR1"
  }
  users {
    name = "user3"
    type = "RANDOM"
  }
}
  