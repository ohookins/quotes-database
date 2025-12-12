resource "aws_secretsmanager_secret" "quotes_database_dsn" {
  name = "quotes-database-dsn"
}

resource "aws_secretsmanager_secret_version" "quotes_database_dsn" {
  secret_id     = aws_secretsmanager_secret.quotes_database_dsn.id
  secret_string = local.dsn
}
