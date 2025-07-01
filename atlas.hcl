env "local" {
  url = "postgres://admin:030603@localhost:5000/core_bank?sslmode=disable&search_path=public"
  dev = "docker://postgres/latest/core_bank"
  migration {
    dir    = "file://db/migrations"
    format = atlas
  }
}