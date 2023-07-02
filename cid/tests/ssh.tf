ssh_keys = {
  rsa_private     = <<EOT
    -----BEGIN RSA PRIVATE KEY-----
    MIIBxwIBAAJhAKD0YSHy73nUgysO13XsJmd4fHiFyQ+00R7VVu2iV9Qco
    ...
    -----END RSA PRIVATE KEY-----
EOT
  rsa_public      = "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAGEAoPRhIfLvedSDKw7Xd ..."
  rsa_certificate = <<EOT
    ssh-rsa-cert-v01@openssh.com AAAAIHNzaC1lZDI1NTE5LWNlcnQt ...
EOT
  dsa_private     = <<EOT
    -----BEGIN DSA PRIVATE KEY-----
    MIIBxwIBAAJhAKD0YSHy73nUgysO13XsJmd4fHiFyQ+00R7VVu2iV9Qco
    ...
    -----END DSA PRIVATE KEY-----
EOT
  dsa_public      = "ssh-dsa AAAAB3NzaC1yc2EAAAABIwAAAGEAoPRhIfLvedSDKw7Xd ..."
  dsa_certificate = <<EOT
    ssh-dsa-cert-v01@openssh.com AAAAIHNzaC1lZDI1NTE5LWNlcnQt ...
EOT
}
ssh_authorized_keys = [
  "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAGEA3FSyQwBI6Z+nCSjUU ...",
  "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA3I7VUf2l5gSn5uavROsc5HRDpZ ..."
]
ssh_deletekeys        = true
ssh_genkeytypes       = ["rsa", "dsa", "ecdsa", "ed25519"]
disable_root          = true
disable_root_opts     = "no-port-forwarding,no-agent-forwarding,no-X11-forwarding"
allow_public_ssh_keys = true
ssh_quiet_keygen      = true
ssh_publish_hostkeys {
  enabled   = true
  blacklist = ["dsa"]
}