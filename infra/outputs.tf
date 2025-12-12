output "ecr_url" {
  value = aws_ecr_repository.main.repository_url
}

output "service_url" {
  description = "Publicly accessible URL of the AppRunner service"
  value       = aws_apprunner_service.main.service_url
}
