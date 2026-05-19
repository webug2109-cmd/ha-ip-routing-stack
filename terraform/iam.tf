resource "google_project_iam_binding" "k8s_admin" {
  project = var.project_id
  role    = "roles/container.admin"
  members = [
    "user:admin@example.com",
  ]
}
