provider "mysql" {
  endpoint = "localhost:3306"
  username = "root"
  password = "root"
}

resource "mysql_user" "web" {
  user               = "web"
  host               = "localhost"
  plaintext_password = "web"
}

resource "mysql_database" "web" {
  name                  = "web"
  default_character_set = "utf8mb4"
  default_collation     = "utf8mb4_bin"
}

resource "mysql_grant" "web" {
  user       = "${mysql_user.web.user}"
  host       = "${mysql_user.web.host}"
  database   = "${mysql_database.web.name}"
  privileges = ["SELECT", "UPDATE", "INSERT", "DELETE"]
}