---
page_title: "kubiya_webhook Resource - Kubiya"
subcategory: ""
description: |-
  The kubiya_webhook resource manages webhook triggers for agents and workflows in the Kubiya platform.
---

# kubiya_webhook (Resource)

The `kubiya_webhook` resource allows you to create and manage webhooks in the Kubiya platform. Webhooks enable external systems to trigger Kubiya agents or workflows via HTTP requests, facilitating event-driven automation.

## Prerequisites

Before using this resource, ensure you have:
1. A Kubiya account with API access
2. An API key (generated from Kubiya dashboard under Admin → Kubiya API Keys)
3. At least one configured agent or workflow to trigger
4. Appropriate integrations configured for notification methods (Slack, Teams)

## Example Usage

### 1. Basic Webhook with Agent

Create a simple webhook that triggers an agent:

```hcl
terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
    }
  }
}

provider "kubiya" {
  # API key will be taken from KUBIYA_API_KEY environment variable
}

resource "kubiya_agent" "support_agent" {
  name         = "support-assistant"
  runner       = "kubiya-hosted"
  description  = "Customer support agent"
  instructions = "You are a support assistant that helps resolve customer issues."
}

resource "kubiya_webhook" "support_webhook" {
  name   = "customer-support-webhook"
  agent  = kubiya_agent.support_agent.name
  source = "zendesk"
  prompt = "New support ticket received. Please analyze and provide initial response."
}
```

**Expected Outcome**: Creates a webhook that triggers the support agent when called.

### 2. Webhook with Slack Notifications

Configure a webhook that sends notifications to Slack:

```hcl
resource "kubiya_webhook" "github_webhook" {
  name        = "github-pr-webhook"
  agent       = "code-review-agent"
  source      = "github"
  filter      = "pull_request"
  prompt      = "Review the new pull request and provide feedback"
  
  # Slack notification configuration
  method      = "Slack"
  destination = "#code-reviews"
}
```

**Expected Outcome**: Creates a webhook that triggers the agent and sends notifications to a Slack channel.

### 3. Webhook with Agent - Advanced Configuration

Create a webhook with comprehensive configuration for incident response:

```hcl
resource "kubiya_agent" "incident_agent" {
  name         = "incident-responder"
  runner       = "kubiya-hosted"
  description  = "Automated incident response agent"
  instructions = <<-EOT
    You are an incident response agent. When triggered:
    1. Analyze the incident severity
    2. Gather relevant logs and metrics
    3. Create a JIRA ticket
    4. Notify the on-call team
    5. Start initial remediation steps
  EOT
  
  integrations = ["pagerduty", "jira_cloud", "datadog", "slack_integration"]
}

resource "kubiya_webhook" "incident_webhook" {
  name        = "critical-incident-webhook"
  agent       = kubiya_agent.incident_agent.name
  source      = "datadog"
  filter      = "severity:critical"
  prompt      = "Critical incident detected. Initiate incident response protocol."
  
  method      = "Slack"
  destination = "#incidents-critical"
}
```

**Expected Outcome**: Creates a webhook that triggers comprehensive incident response automation.

### 4. Webhook with Microsoft Teams Notifications

Configure a webhook with Teams notifications:

```hcl
resource "kubiya_webhook" "deployment_webhook" {
  name        = "deployment-notification"
  agent       = "deployment-agent"
  source      = "github-actions"
  filter      = "workflow:deploy"
  prompt      = "Deployment completed. Verify deployment health and report status."
  
  # Microsoft Teams notification
  method      = "Teams"
  team_name   = "DevOps Team"
  destination = "Deployments"  # Channel name in Teams
}
```

**Expected Outcome**: Creates a webhook that sends notifications to a Microsoft Teams channel.

### 5. Webhook with Workflow - Simple

Create a webhook that directly triggers a workflow:

```hcl
resource "kubiya_webhook" "backup_webhook" {
  name   = "nightly-backup"
  source = "scheduler"
  prompt = "Execute nightly backup workflow"
  
  workflow = jsonencode({
    name        = "nightly-backup-workflow"
    description = "Automated nightly backup process"
    steps = [
      {
        name        = "backup-database"
        description = "Backup production database"
        executor = {
          type = "command"
          config = {
            command = "pg_dump production_db > /backups/db_$(date +%Y%m%d).sql"
          }
        }
      },
      {
        name        = "upload-to-s3"
        description = "Upload backup to S3"
        depends     = ["backup-database"]
        executor = {
          type = "command"
          config = {
            command = "aws s3 cp /backups/db_$(date +%Y%m%d).sql s3://company-backups/"
          }
        }
      },
      {
        name        = "verify-backup"
        description = "Verify backup integrity"
        depends     = ["upload-to-s3"]
        executor = {
          type = "command"
          config = {
            command = "aws s3 ls s3://company-backups/db_$(date +%Y%m%d).sql"
          }
        }
      }
    ]
  })
  
  runner = "kubiya-hosted"
}
```

**Expected Outcome**: Creates a webhook that executes a backup workflow when triggered.

### 6. Webhook with Complex Workflow

Create a webhook with a sophisticated multi-step workflow:

```hcl
resource "kubiya_webhook" "cicd_webhook" {
  name   = "cicd-pipeline"
  source = "github"
  filter = "push:main"
  prompt = "Execute CI/CD pipeline for main branch push"
  
  workflow = jsonencode({
    name        = "cicd-pipeline"
    description = "Complete CI/CD pipeline with testing and deployment"
    steps = [
      {
        name        = "run-tests"
        description = "Execute unit and integration tests"
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "test-runner"
              description = "Run test suite"
              type        = "docker"
              image       = "node:18"
              content     = "npm test && npm run test:integration"
            }
          }
        }
        output = "TEST_RESULTS"
      },
      {
        name        = "build-application"
        description = "Build Docker image"
        depends     = ["run-tests"]
        executor = {
          type = "command"
          config = {
            command = "docker build -t myapp:$${GITHUB_SHA} ."
          }
        }
        output = "BUILD_STATUS"
      },
      {
        name        = "security-scan"
        description = "Run security vulnerability scan"
        depends     = ["build-application"]
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "security-scanner"
              description = "Scan for vulnerabilities"
              type        = "docker"
              image       = "aquasec/trivy"
              content     = "trivy image myapp:$${GITHUB_SHA}"
            }
          }
        }
        output = "SCAN_RESULTS"
      },
      {
        name        = "deploy-staging"
        description = "Deploy to staging environment"
        depends     = ["security-scan"]
        executor = {
          type = "command"
          config = {
            command = "kubectl set image deployment/myapp myapp=myapp:$${GITHUB_SHA} -n staging"
          }
        }
      },
      {
        name        = "smoke-tests"
        description = "Run smoke tests on staging"
        depends     = ["deploy-staging"]
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "smoke-tester"
              description = "Execute smoke tests"
              type        = "docker"
              image       = "postman/newman"
              content     = "newman run staging-smoke-tests.json"
            }
          }
        }
        output = "SMOKE_TEST_RESULTS"
      },
      {
        name        = "notify-team"
        description = "Send deployment notification"
        depends     = ["smoke-tests"]
        executor = {
          type = "agent"
          config = {
            teammate_name = "notification-agent"
            message       = "Deployment to staging complete. Test results: $${SMOKE_TEST_RESULTS}"
          }
        }
      }
    ]
  })
  
  runner      = "kubiya-hosted"
  method      = "Slack"
  destination = "#deployments"
}
```

**Expected Outcome**: Creates a webhook that executes a complete CI/CD pipeline workflow.

### 7. Webhook with Data Processing Workflow

Create a webhook for data processing workflows:

```hcl
resource "kubiya_webhook" "data_processing_webhook" {
  name   = "etl-pipeline"
  source = "s3"
  filter = "bucket:raw-data"
  prompt = "New data file uploaded, process ETL pipeline"
  
  workflow = jsonencode({
    name        = "etl-workflow"
    description = "Extract, Transform, Load data pipeline"
    steps = [
      {
        name        = "extract-data"
        description = "Extract data from source"
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "data-extractor"
              description = "Extract and validate data"
              type        = "docker"
              image       = "python:3.11-slim"
              with_files = [
                {
                  destination = "/extract.py"
                  content     = <<-PYTHON
                    import pandas as pd
                    import sys
                    
                    # Read data from S3
                    df = pd.read_csv(sys.argv[1])
                    
                    # Validate data
                    assert not df.empty, "Data is empty"
                    assert df.columns.tolist() == ['id', 'name', 'value'], "Invalid schema"
                    
                    # Save validated data
                    df.to_csv('/tmp/validated_data.csv', index=False)
                    print(f"Extracted {len(df)} records")
                  PYTHON
                }
              ]
              content = "python /extract.py $${DATA_FILE}"
              args = [
                {
                  name        = "DATA_FILE"
                  type        = "string"
                  description = "Input data file path"
                  required    = true
                }
              ]
            }
          }
        }
        output = "EXTRACTED_COUNT"
      },
      {
        name        = "transform-data"
        description = "Apply business transformations"
        depends     = ["extract-data"]
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "data-transformer"
              description = "Transform data according to business rules"
              type        = "docker"
              image       = "python:3.11-slim"
              with_files = [
                {
                  destination = "/transform.py"
                  content     = <<-PYTHON
                    import pandas as pd
                    import numpy as np
                    
                    # Load validated data
                    df = pd.read_csv('/tmp/validated_data.csv')
                    
                    # Apply transformations
                    df['value_normalized'] = (df['value'] - df['value'].mean()) / df['value'].std()
                    df['category'] = pd.cut(df['value'], bins=3, labels=['low', 'medium', 'high'])
                    df['processed_date'] = pd.Timestamp.now()
                    
                    # Save transformed data
                    df.to_parquet('/tmp/transformed_data.parquet')
                    print(f"Transformed {len(df)} records")
                  PYTHON
                }
              ]
              content = "python /transform.py"
            }
          }
        }
        output = "TRANSFORMED_COUNT"
      },
      {
        name        = "load-to-warehouse"
        description = "Load data to data warehouse"
        depends     = ["transform-data"]
        executor = {
          type = "command"
          config = {
            command = "aws s3 cp /tmp/transformed_data.parquet s3://data-warehouse/processed/$(date +%Y%m%d)/"
          }
        }
      }
    ]
  })
  
  runner = "kubiya-hosted"
}
```

**Expected Outcome**: Creates a webhook that triggers a complete ETL pipeline for data processing.

### 8. Webhook with Conditional Workflow

Create a webhook with conditional logic in the workflow:

```hcl
resource "kubiya_webhook" "conditional_webhook" {
  name   = "smart-deployment"
  source = "github"
  filter = "release"
  prompt = "New release created, execute smart deployment"
  
  workflow = jsonencode({
    name        = "conditional-deployment"
    description = "Deployment with environment-based conditions"
    steps = [
      {
        name        = "check-environment"
        description = "Determine target environment based on tag"
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "env-checker"
              description = "Check deployment environment"
              type        = "docker"
              image       = "alpine:latest"
              content     = <<-BASH
                if [[ "$${RELEASE_TAG}" == *"-prod"* ]]; then
                  echo "production"
                elif [[ "$${RELEASE_TAG}" == *"-staging"* ]]; then
                  echo "staging"
                else
                  echo "development"
                fi
              BASH
              args = [
                {
                  name        = "RELEASE_TAG"
                  type        = "string"
                  description = "Release tag name"
                  required    = true
                }
              ]
            }
            args = {
              RELEASE_TAG = "$${GITHUB_RELEASE_TAG}"
            }
          }
        }
        output = "TARGET_ENV"
      },
      {
        name        = "deploy-to-environment"
        description = "Deploy to determined environment"
        depends     = ["check-environment"]
        executor = {
          type = "agent"
          config = {
            teammate_name = "deployment-agent"
            message       = "Deploy release $${GITHUB_RELEASE_TAG} to $${TARGET_ENV} environment"
          }
        }
      },
      {
        name        = "run-environment-tests"
        description = "Run environment-specific tests"
        depends     = ["deploy-to-environment"]
        executor = {
          type = "command"
          config = {
            command = "npm run test:$${TARGET_ENV}"
          }
        }
      }
    ]
  })
  
  runner      = "kubiya-hosted"
  method      = "Slack"
  destination = "#releases"
}
```

**Expected Outcome**: Creates a webhook with a workflow that adapts based on conditions.

## Argument Reference

### Required Arguments

* `name` - (Required, String) The name of the webhook. Must be unique within your organization.

### Conditional Required Arguments

Either `agent` or `workflow` must be specified:

* `agent` - (Required if no workflow, String) The name of the agent to trigger with this webhook.
* `workflow` - (Required if no agent, String) JSON-encoded workflow definition to execute when webhook is triggered.

### Optional Arguments

* `source` - (Optional, String) Source identification for the webhook (e.g., "github", "datadog", "custom").

* `filter` - (Optional, String) Filter expression for events that will trigger the webhook.

* `prompt` - (Optional, String) Prompt to send to the agent or workflow when triggered.

* `runner` - (Optional, String) The runner to use for workflow execution. Required when using `workflow`.

* `method` - (Optional, String) Notification method. Values: "Slack", "Teams", "Http". Defaults to "Slack".

* `destination` - (Optional, String) Destination for notifications:
  - For Slack: Channel name with "#" prefix (e.g., "#alerts") or username with "@" prefix
  - For Teams: Channel name within the team specified by `team_name`
  - For Http: Not required

* `team_name` - (Optional, String) Team name for Microsoft Teams notifications. Required when `method` is "Teams".

### Workflow Structure

When using the `workflow` parameter, the JSON structure should include:

* `name` - Workflow name
* `description` - Workflow description
* `steps` - Array of workflow steps, each containing:
  - `name` - Step name
  - `description` - Step description
  - `executor` - Executor configuration with `type` and `config`
  - `depends` - (Optional) Array of step names this step depends on
  - `output` - (Optional) Variable name to store step output

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the webhook.
* `url` - The URL to trigger the webhook (sensitive).
* `created_at` - The timestamp when the webhook was created.
* `created_by` - The user who created the webhook.
* `status` - The current status of the webhook.
* `workflow_id` - The ID of the associated workflow (when using workflow parameter).

## Import

Webhooks can be imported using their ID:

```shell
terraform import kubiya_webhook.example <webhook-id>
```

## Compatibility Notes

* Requires Kubiya Terraform Provider version >= 1.0.0
* Compatible with Terraform >= 1.0
* Workflow support requires appropriate runner configuration
* Notification methods require corresponding integrations to be configured
* Some features may require specific Kubiya platform tier

## Best Practices

1. **Security**: Use filters to restrict webhook triggers to specific events
2. **Naming**: Use descriptive names that indicate the webhook's purpose and source
3. **Error Handling**: Include error handling steps in workflows
4. **Testing**: Test webhooks with sample payloads before production use
5. **Documentation**: Document expected payload formats and trigger conditions
6. **Monitoring**: Set up appropriate notifications for webhook execution status
7. **Rate Limiting**: Be aware of rate limits when designing high-frequency webhooks