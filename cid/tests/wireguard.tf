wireguard {
  interfaces {
    name        = "wg0"
    config_path = "/etc/wireguard/wg0.conf"
    content     = <<EOT
        [Interface]
        PrivateKey = <private_key>
        Address = <address>
        [Peer]
        PublicKey = <public_key>
        Endpoint = <endpoint_ip>:<endpoint_ip_port>
        AllowedIPs = <allowedip1>, <allowedip2>, ...
EOT
  }
  interfaces {
    name        = "wg1"
    config_path = "/etc/wireguard/wg1.conf"
    content     = <<EOT
        [Interface]
        PrivateKey = <private_key>
        Address = <address>
        [Peer]
        PublicKey = <public_key>
        Endpoint = <endpoint_ip>:<endpoint_ip_port>
        AllowedIPs = <allowedip1>
EOT
  }
  readinessprobe = [
    "systemctl restart service",
    "curl https://webhook.endpoint/example",
    "nc -zv some-service-fqdn 443"
  ]
}