variable "config_json" {
  description = "List of Kubiya integrations to enable. Supports multiple values. For AWS integration, the main account must be provided."
  type        = string
  default     = <<-EOT
    {
        "access_configs": {
            "DB Access to Staging": {
                "name": "Database Access to Staging",
                "description": "Grants access to all staging RDS databases",
                "account_id": "***",
                "permission_set": "ECRReadOnly",
                "session_duration": "PT5M"
            },
            "Power User to SandBox": {
                "name": "Power User Access to SandBox",
                "description": "Grants poweruser permissions on Sandbox",
                "account_id": "****",
                "permission_set": "PowerUserAccess",
                "session_duration": "PT5M"
            }
        },
        "s3_configs": {
            "Data Lake Read Access": {
                "name": "data_lake_read 4",
                "description": "Grants read-only access to data lake buckets",
                "buckets": [
                    "company-data-lake-prod",
                    "company-data-lake-staging"
                ],
                "policy_template": "S3ReadOnlyPolicy",
                "session_duration": "PT1H"
            }
        }
    }
  EOT
}