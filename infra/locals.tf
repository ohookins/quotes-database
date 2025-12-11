locals {
  subnet_cidrs = cidrsubnets(var.vpc_cidr, 2, 2)
  # subnet_cidrs[0] and [1]: 2 blocks, then split each into 2 for 4 subnets
  public_subnet_cidrs  = [cidrsubnets(local.subnet_cidrs[0], 1)[0], cidrsubnets(local.subnet_cidrs[0], 1)[1]]
  private_subnet_cidrs = [cidrsubnets(local.subnet_cidrs[1], 1)[0], cidrsubnets(local.subnet_cidrs[1], 1)[1]]
}
