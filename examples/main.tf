terraform {
  required_providers {
    kubiya = {
      source = "local/provider/example"
    }
  }
}

provider "kubiya" {
  //Your Kubiya API Key will be taken from the
  //environment variable KUBIYA_API_KEY
  //To set the key, please use export KUBIYA_API_KEY="YOUR_API_KEY"
}

resource "kubiya_scheduled_task" "example" {
  repeat = ""
  // Optional. Allowed values: hourly, daily, weekly, monthly. Leaving this value empty or omitting it will cause the task to be executed only once.
  channel_id     = "C082X4R0FL0"
  agent          = "Which Agent should perform the task?"
  scheduled_time = "2024-12-01T05:00:00"
  description    = "Describe the task"
}


resource "kubiya_source" "item" {
  url = "https://github.com/finebee/terraform-golden-usecases"
  runner = "avi-stg-test"
}

resource "kubiya_source" "item_config" {
  url = "https://github.com/finebee/terraform-golden-usecases"
  dynamic_config = var.s3_configs_json
  runner = "avi-stg-test"
}



variable "s3_configs_json" {
  description = "List of Kubiya integrations to enable. Supports multiple values. \n For AWS integration, the main account must be provided."
  type        = string
  default     = <<-EOT
    {
        "access_configs": {
            "DB Access to Staging": {
                "name": "Database Access to Staging 4",
                "description": "Grants access to all staging RDS databases",
                "account_id": "876809951775",
                "permission_set": "ECRReadOnly",
                "session_duration": "PT5M"
            },
            "Power User to SandBox": {
                "name": "Database Access to SandBox",
                "description": "Grants poweruser permissions on Sandbox",
                "account_id": "110327817829",
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

output "output" {
  value = kubiya_source.item
}