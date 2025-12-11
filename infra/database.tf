resource "random_pet" "database_username" {
  separator = ""
}

resource "random_password" "database_password" {
  length  = 32
  special = false
}

resource "aws_rds_cluster" "main" {
  cluster_identifier     = "quotes"
  engine                 = "aurora-postgresql"
  engine_mode            = "provisioned"
  engine_version         = "17.4"
  database_name          = "quotes"
  db_subnet_group_name   = aws_db_subnet_group.main.name
  master_username        = random_pet.database_username.id
  master_password        = random_password.database_password.result
  vpc_security_group_ids = [aws_security_group.database.id]

  serverlessv2_scaling_configuration {
    max_capacity             = 1.0
    min_capacity             = 0.0
    seconds_until_auto_pause = 3600
  }
}

resource "aws_rds_cluster_instance" "main" {
  count = 2

  cluster_identifier = aws_rds_cluster.main.id
  identifier_prefix  = "quotes-instance-"
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

resource "aws_security_group" "database" {
  name        = "quotes-database-sg"
  description = "Security group for Aurora PostgreSQL cluster"
  vpc_id      = aws_vpc.main.id
  tags = {
    Name = "quotes-database-sg"
  }
}

resource "aws_vpc_security_group_ingress_rule" "database_ingress" {
  security_group_id = aws_security_group.database.id
  description       = "PostgreSQL from VPC"
  from_port         = 5432
  to_port           = 5432
  ip_protocol       = "tcp"
  cidr_ipv4         = local.private_block
}
