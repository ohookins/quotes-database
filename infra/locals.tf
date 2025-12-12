locals {
  subnet_cidrs  = cidrsubnets(var.vpc_cidr, 1, 1)
  public_block  = local.subnet_cidrs[0]
  private_block = local.subnet_cidrs[1]

  # subnet_cidrs[0] and [1]: 2 blocks, then split each into 2 for 4 subnets
  public_subnet_cidrs  = cidrsubnets(local.public_block, 1, 1)
  private_subnet_cidrs = cidrsubnets(local.private_block, 1, 1)

  dsn = "host=${aws_rds_cluster.main.endpoint} user=${random_pet.database_username.id} password=${random_password.database_password.result} dbname=${aws_rds_cluster.main.database_name} port=${aws_rds_cluster.main.port}"
}
