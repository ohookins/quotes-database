resource "random_pet" "database_username" {
  separator = ""
}

resource "random_password" "database_password" {
  length  = 32
  special = false
}

resource "aws_rds_cluster" "main" {
  cluster_identifier   = "quotes"
  engine               = "aurora-postgresql"
  engine_mode          = "provisioned"
  engine_version       = "17.4"
  database_name        = "quotes"
  db_subnet_group_name = aws_db_subnet_group.main.name
  master_username      = random_pet.database_username.id
  master_password      = random_password.database_password.result


  serverlessv2_scaling_configuration {
    max_capacity             = 1.0
    min_capacity             = 0.0
    seconds_until_auto_pause = 3600
  }
}

resource "aws_rds_cluster_instance" "main" {
  cluster_identifier = aws_rds_cluster.main.id
  instance_class     = "db.serverless"
  engine             = aws_rds_cluster.main.engine
  engine_version     = aws_rds_cluster.main.engine_version
}

resource "aws_db_subnet_group" "main" {
  name       = "main"
  subnet_ids = aws_subnet.private.*.id

  tags = {
    Name = "Private DB subnet group"
  }
}
