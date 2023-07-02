data "vault_secret" "pve_cid_key" {
  path = "secret/pve"
}

provider "cloudinit-drive" {
  ssh {
    host        = "proxmox-node001.local"                        # env SSH_HOST
    port        = 22                                             # env SSH_PORT
    user        = "cid"                                          # env SSH_USER
    private_key = data.vault_secret.pve_cid_key.data["id_ecdsa"] # env SSH_PRIVATE_KEY
    auth_socket = "$TMPDIR/ssh-XXXXXXXXXX/agent.<ppid>"          # env SSH_AUTH_SOCK
  }
}
