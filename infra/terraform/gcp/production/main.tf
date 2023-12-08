resource "google_artifact_registry_repository" "my_repository" {
  provider = google
  location = var.default_region
  repository_id = var.project_id
  description = "Docker repository"
  format = "DOCKER"
}

resource "google_project_iam_member" "cloudbuild_cloudrun_admin" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${var.cloudbuild_service_account}"
}