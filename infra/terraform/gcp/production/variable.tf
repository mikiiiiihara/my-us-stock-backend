variable "credentials_key_path" {
    description = "credentials key path"
    type        = string
}

variable "project_id" {
  description = "project id"
  type        = string
}

variable "default_region" {
  description = "The default region for resources"
    type        = string
}

variable "cloudbuild_service_account" {
  description = "The service account email of Cloud Build."
  type        = string
}