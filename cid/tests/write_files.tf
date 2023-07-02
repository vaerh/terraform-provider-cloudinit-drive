write_files {
  encoding    = "b64"
  content     = "CiMgVGhpcyBmaWxlIGNvbnRyb2xzIHRoZSBzdGF0ZSBvZiBTRUxpbnV4..."
  owner       = "root:root"
  path        = "/etc/sysconfig/selinux"
  permissions = "'0644'"
}

write_files {
  content = <<EOT
    15 * * * * root ship_logs
EOT
  path    = "/etc/crontab"
  append  = true
}

write_files {
  encoding    = "gzip"
  content     = "bin__H4sIAIDb/U8C/1NW1E/KzNMvzuBKTc7IV8hIzcnJVyjPL8pJ4QIA6N+MVxsAAAA="
  path        = "/usr/bin/hello"
  permissions = "'0755'"
}

write_files {
  path = "/root/CLOUD_INIT_WAS_HERE"
}

write_files {
  path        = "/etc/nginx/conf.d/example.com.conf"
  content     = <<EOT
    server {
        server_name example.com;
        listen 80;
        root /var/www;
        location / {
            try_files $uri $uri/ $uri.html =404;
        }
    }
EOT
  owner       = "nginx:nginx"
  permissions = "'0640'"
  defer       = true
}

write_files {
  local_file = "tests/test_file.txt"
  path       = "/tmp/write_files.tf"
  encoding   = "gz+b64"
}