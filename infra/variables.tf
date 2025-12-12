variable "image_url" {
  description = "Full URL to the application image in ECR"
  type        = string
}

variable "vpc_cidr" {
  description = "The IPv4 CIDR block for the VPC"
  type        = string
}
