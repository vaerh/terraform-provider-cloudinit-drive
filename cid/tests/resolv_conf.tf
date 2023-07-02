manage_resolv_conf = true
resolv_conf {
  nameservers   = ["8.8.8.8", "8.8.4.4"]
  searchdomains = ["foo.example.com", "bar.example.com"]
  domain        = "example.com"
  sortlist      = ["10.0.0.1/255", "10.0.0.2"]
  options = {
    rotate  = true
    timeout = 1
  }
}