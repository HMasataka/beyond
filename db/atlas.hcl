env "docker" {
  src = "file://schema.hcl"
  url = "mysql://user:password@mysql:3306/beyond-db"
}

env "local" {
  src = "file://schema.hcl"
  url = "mysql://user:password@127.0.0.1:3306/beyond-db"
}
