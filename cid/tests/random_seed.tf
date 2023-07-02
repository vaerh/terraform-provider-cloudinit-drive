random_seed {
  file             = "/dev/urandom"
  data             = <<EOT
bin__H4sIAAAAAAAAAwGAAH//lgRNKp3SkVVNP2jbhM1PVjvWE2DCjV7kzNg80ga5640qlc8KT64tsyyw
1GCCQf5e9ScfMpU/pLQtha9mfaqcrdj6uIAm7qbzD9TzzbNCwuOFl5yeexlpEHFFmH2JmjqlEBMB
YbUEKw49XDkxMQQuWUjh5EcmrakXgQrvc4wPex2vUqGkgAAAAA==
EOT
  encoding         = "gz"
  command          = ["sh", "-c", "dd if=/dev/urandom of=$RANDOM_SEED_FILE"]
  command_required = true
}