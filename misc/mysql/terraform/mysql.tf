variable "env" {}
variable "root_password" {}
variable "password" {}

provider "mysql" {
  endpoint = "localhost:3306"
  username = "root"
  password = "${var.root_password}"
}

resource "mysql_user" "web" {
  user               = "web"
  host               = "${var.env == "dev" ? "172.17.0.1" : "localhost"}"
  plaintext_password = "${var.password}"
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