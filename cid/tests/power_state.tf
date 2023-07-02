power_state {
  delay     = 30
  mode      = "reboot"
  timeout   = 2
  message   = "Rebooting machine"
  condition = "test -f /var/tmp/reboot_me"
}