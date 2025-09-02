---
page_title: "kubiya_trigger Resource - Kubiya"
subcategory: ""
description: |-
  The kubiya_trigger resource manages workflow triggers with webhook capabilities in the Kubiya platform.
---

# kubiya_trigger (Resource)

The `kubiya_trigger` resource allows you to create and manage workflow triggers in the Kubiya platform. This resource creates a workflow and automatically publishes it with a webhook trigger, providing a URL that can be used to execute the workflow via HTTP requests.

## Prerequisites

Before using this resource, ensure you have:
1. A Kubiya account with API access
2. An API key (generated from Kubiya dashboard under Admin → Kubiya API Keys)
3. A configured runner (or use "kubiya-hosted" for cloud execution)

## Example Usage

### 1. Simple Echo Workflow Trigger

Create a basic trigger that executes a simple command:

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

resource "kubiya_trigger" "simple_trigger" {
  name   = "hello-world-trigger"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "Hello World Workflow"
    version = 1
    steps = [
      {
        name = "greeting"
        executor = {
          type = "command"
          config = {
            command = "echo 'Hello from Kubiya Workflow!'"
          }
        }
      }
    ]
  })
}

output "webhook_url" {
  value       = kubiya_trigger.simple_trigger.url
  sensitive   = true
  description = "POST to this URL to trigger the workflow"
}
```

**Expected Outcome**: Creates a simple workflow trigger that outputs a greeting message when called.

### 2. Multi-Step Sequential Workflow

Create a trigger with multiple sequential steps:

```hcl
resource "kubiya_trigger" "deployment_trigger" {
  name   = "deployment-pipeline"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "Deployment Pipeline"
    version = 1
    steps = [
      {
        name = "validate-config"
        executor = {
          type = "command"
          config = {
            command = "echo 'Validating configuration...'"
          }
        }
      },
      {
        name = "run-tests"
        executor = {
          type = "command"
          config = {
            command = "echo 'Running tests...' && sleep 2 && echo 'Tests passed!'"
          }
        }
      },
      {
        name = "deploy"
        executor = {
          type = "command"
          config = {
            command = "echo 'Deploying application...'"
          }
        }
      },
      {
        name = "notify"
        executor = {
          type = "command"
          config = {
            command = "echo 'Deployment completed successfully!'"
          }
        }
      }
    ]
  })
}
```

**Expected Outcome**: Creates a trigger that executes a deployment pipeline with validation, testing, and notification steps.

### 3. Data Processing Workflow Trigger

Create a trigger for data processing operations:

```hcl
resource "kubiya_trigger" "data_processing" {
  name   = "data-etl-pipeline"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "ETL Data Pipeline"
    version = 1
    steps = [
      {
        name = "fetch-data"
        executor = {
          type = "command"
          config = {
            command = "curl -s https://api.example.com/data | jq '.'"
          }
        }
      },
      {
        name = "transform-data"
        executor = {
          type = "command"
          config = {
            command = "echo 'Transforming data...' && jq '.items | map({id: .id, value: .value * 2})'"
          }
        }
      },
      {
        name = "load-data"
        executor = {
          type = "command"
          config = {
            command = "echo 'Loading data to warehouse...' && echo 'Data loaded successfully'"
          }
        }
      }
    ]
  })
}
```

**Expected Outcome**: Creates a trigger for an ETL pipeline that fetches, transforms, and loads data.

### 4. Workflow with Tool Executors

Create a trigger with custom tool executors:

```hcl
resource "kubiya_trigger" "monitoring_trigger" {
  name   = "system-health-check"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "System Health Monitoring"
    version = 1
    steps = [
      {
        name        = "check-disk-usage"
        description = "Check disk usage on all systems"
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "disk-checker"
              description = "Check disk usage"
              type        = "docker"
              image       = "alpine:latest"
              content     = "df -h | awk '$5+0 > 80 {print \"Warning: \" $1 \" is \" $5 \" full\"}'"
            }
          }
        }
        output = "DISK_STATUS"
      },
      {
        name        = "check-memory"
        description = "Check memory usage"
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "memory-checker"
              description = "Check memory usage"
              type        = "docker"
              image       = "alpine:latest"
              content     = "free -m | awk 'NR==2{printf \"Memory Usage: %.2f%%\\n\", $3*100/$2}'"
            }
          }
        }
        output = "MEMORY_STATUS"
      },
      {
        name        = "generate-report"
        description = "Generate health report"
        executor = {
          type = "command"
          config = {
            command = "echo 'System Health Report Generated'"
          }
        }
      }
    ]
  })
}
```

**Expected Outcome**: Creates a trigger that performs system health checks using Docker-based tools.

### 5. Workflow with Dependencies

Create a trigger with step dependencies:

```hcl
resource "kubiya_trigger" "parallel_workflow" {
  name   = "parallel-processing"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "Parallel Processing Workflow"
    version = 1
    steps = [
      {
        name = "prepare-environment"
        executor = {
          type = "command"
          config = {
            command = "echo 'Preparing environment...'"
          }
        }
      },
      {
        name    = "process-batch-1"
        depends = ["prepare-environment"]
        executor = {
          type = "command"
          config = {
            command = "echo 'Processing batch 1...'"
          }
        }
      },
      {
        name    = "process-batch-2"
        depends = ["prepare-environment"]
        executor = {
          type = "command"
          config = {
            command = "echo 'Processing batch 2...'"
          }
        }
      },
      {
        name    = "process-batch-3"
        depends = ["prepare-environment"]
        executor = {
          type = "command"
          config = {
            command = "echo 'Processing batch 3...'"
          }
        }
      },
      {
        name    = "aggregate-results"
        depends = ["process-batch-1", "process-batch-2", "process-batch-3"]
        executor = {
          type = "command"
          config = {
            command = "echo 'Aggregating results from all batches...'"
          }
        }
      }
    ]
  })
}
```

**Expected Outcome**: Creates a trigger with parallel processing capabilities where multiple steps can run simultaneously.

### 6. Advanced CI/CD Workflow Trigger

Create a comprehensive CI/CD pipeline trigger:

```hcl
resource "kubiya_trigger" "cicd_trigger" {
  name   = "cicd-pipeline"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "CI/CD Pipeline"
    version = 2
    steps = [
      {
        name        = "checkout-code"
        description = "Checkout latest code"
        executor = {
          type = "command"
          config = {
            command = "git clone https://github.com/example/repo.git /tmp/repo && cd /tmp/repo && git log -1"
          }
        }
        output = "COMMIT_SHA"
      },
      {
        name        = "run-unit-tests"
        description = "Execute unit tests"
        depends     = ["checkout-code"]
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "test-runner"
              description = "Run unit tests"
              type        = "docker"
              image       = "node:18"
              content     = "cd /tmp/repo && npm install && npm test"
            }
          }
        }
        output = "TEST_RESULTS"
      },
      {
        name        = "build-docker-image"
        description = "Build Docker image"
        depends     = ["run-unit-tests"]
        executor = {
          type = "command"
          config = {
            command = "docker build -t myapp:latest /tmp/repo"
          }
        }
        output = "IMAGE_ID"
      },
      {
        name        = "security-scan"
        description = "Scan for vulnerabilities"
        depends     = ["build-docker-image"]
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "security-scanner"
              description = "Scan Docker image for vulnerabilities"
              type        = "docker"
              image       = "aquasec/trivy"
              content     = "trivy image --severity HIGH,CRITICAL myapp:latest"
            }
          }
        }
        output = "SCAN_REPORT"
      },
      {
        name        = "deploy-to-kubernetes"
        description = "Deploy to Kubernetes"
        depends     = ["security-scan"]
        executor = {
          type = "command"
          config = {
            command = "kubectl apply -f /tmp/repo/k8s/deployment.yaml && kubectl rollout status deployment/myapp"
          }
        }
      },
      {
        name        = "run-smoke-tests"
        description = "Execute smoke tests"
        depends     = ["deploy-to-kubernetes"]
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "smoke-tester"
              description = "Run smoke tests against deployed application"
              type        = "docker"
              image       = "postman/newman"
              content     = "newman run /tmp/repo/tests/smoke-tests.json"
            }
          }
        }
        output = "SMOKE_TEST_RESULTS"
      }
    ]
  })
}

output "cicd_webhook_url" {
  value       = kubiya_trigger.cicd_trigger.url
  sensitive   = true
  description = "CI/CD pipeline trigger URL"
}

output "workflow_status" {
  value = kubiya_trigger.cicd_trigger.status
}
```

**Expected Outcome**: Creates a complete CI/CD pipeline trigger with testing, building, scanning, and deployment steps.

### 7. Database Backup Workflow Trigger

Create a trigger for automated database backups:

```hcl
resource "kubiya_trigger" "backup_trigger" {
  name   = "database-backup"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "Database Backup Workflow"
    version = 1
    steps = [
      {
        name        = "create-backup"
        description = "Create database backup"
        executor = {
          type = "command"
          config = {
            command = "pg_dump -h localhost -U postgres -d production > /tmp/backup_$(date +%Y%m%d_%H%M%S).sql"
          }
        }
        output = "BACKUP_FILE"
      },
      {
        name        = "compress-backup"
        description = "Compress backup file"
        depends     = ["create-backup"]
        executor = {
          type = "command"
          config = {
            command = "gzip /tmp/backup_*.sql"
          }
        }
      },
      {
        name        = "upload-to-s3"
        description = "Upload backup to S3"
        depends     = ["compress-backup"]
        executor = {
          type = "command"
          config = {
            command = "aws s3 cp /tmp/backup_*.sql.gz s3://company-backups/postgres/$(date +%Y/%m/%d)/"
          }
        }
      },
      {
        name        = "cleanup-local"
        description = "Clean up local backup files"
        depends     = ["upload-to-s3"]
        executor = {
          type = "command"
          config = {
            command = "rm -f /tmp/backup_*.sql.gz"
          }
        }
      },
      {
        name        = "verify-backup"
        description = "Verify backup in S3"
        depends     = ["upload-to-s3"]
        executor = {
          type = "command"
          config = {
            command = "aws s3 ls s3://company-backups/postgres/$(date +%Y/%m/%d)/ --recursive"
          }
        }
        output = "BACKUP_VERIFICATION"
      }
    ]
  })
}
```

**Expected Outcome**: Creates a trigger for automated database backup with compression and S3 upload.

### 8. Incident Response Workflow Trigger

Create a trigger for incident response automation:

```hcl
resource "kubiya_trigger" "incident_response" {
  name   = "incident-response"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "Incident Response Automation"
    version = 1
    steps = [
      {
        name        = "gather-metrics"
        description = "Collect system metrics"
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "metrics-collector"
              description = "Gather system metrics"
              type        = "docker"
              image       = "python:3.11-slim"
              with_files = [
                {
                  destination = "/collect_metrics.py"
                  content     = <<-PYTHON
                    import json
                    import datetime
                    
                    metrics = {
                        "timestamp": datetime.datetime.now().isoformat(),
                        "cpu_usage": "75%",
                        "memory_usage": "82%",
                        "disk_usage": "65%",
                        "active_connections": 1250,
                        "error_rate": "2.5%"
                    }
                    
                    print(json.dumps(metrics, indent=2))
                  PYTHON
                }
              ]
              content = "python /collect_metrics.py"
            }
          }
        }
        output = "SYSTEM_METRICS"
      },
      {
        name        = "analyze-logs"
        description = "Analyze error logs"
        depends     = ["gather-metrics"]
        executor = {
          type = "command"
          config = {
            command = "tail -n 100 /var/log/application.log | grep -i error | head -10"
          }
        }
        output = "ERROR_LOGS"
      },
      {
        name        = "create-incident-ticket"
        description = "Create JIRA incident ticket"
        depends     = ["analyze-logs"]
        executor = {
          type = "agent"
          config = {
            teammate_name = "incident-agent"
            message       = "Create a JIRA ticket for incident with metrics: $${SYSTEM_METRICS} and errors: $${ERROR_LOGS}"
          }
        }
        output = "TICKET_ID"
      },
      {
        name        = "notify-oncall"
        description = "Notify on-call team"
        depends     = ["create-incident-ticket"]
        executor = {
          type = "agent"
          config = {
            teammate_name = "notification-agent"
            message       = "Send Slack alert to #incidents about ticket $${TICKET_ID}"
          }
        }
      },
      {
        name        = "initiate-remediation"
        description = "Start auto-remediation"
        depends     = ["notify-oncall"]
        executor = {
          type = "command"
          config = {
            command = "kubectl scale deployment/app --replicas=5 && kubectl rollout restart deployment/app"
          }
        }
      }
    ]
  })
}
```

**Expected Outcome**: Creates a trigger for automated incident response with metrics collection, log analysis, ticket creation, and remediation.

## Triggering the Workflow

Once the trigger resource is created, you can execute the workflow by making HTTP requests to the webhook URL.

### Basic Trigger

```bash
# Get the webhook URL from Terraform output
WEBHOOK_URL=$(terraform output -raw webhook_url)

# Trigger the workflow
curl -X POST "$WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -H "Authorization: UserKey YOUR_API_KEY" \
  -d '{}'
```

### Trigger with Payload

```bash
# Send data to the workflow
curl -X POST "$WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -H "Authorization: UserKey YOUR_API_KEY" \
  -d '{
    "environment": "production",
    "version": "1.2.3",
    "user": "deploy-bot"
  }'
```

### Streaming Execution Output

For real-time execution output, append `?stream=true` to the webhook URL:

```bash
curl -X POST "$WEBHOOK_URL?stream=true" \
  -H "Content-Type: application/json" \
  -H "Authorization: UserKey YOUR_API_KEY" \
  -d '{}'
```

## Argument Reference

### Required Arguments

* `name` - (Required, String) Name of the trigger. This will be used as the workflow name.

* `runner` - (Required, String) Runner to use for executing the workflow. Common values:
  - `kubiya-hosted` - Use Kubiya's cloud-hosted runners
  - Custom runner names from your organization

* `workflow` - (Required, String) JSON-encoded workflow definition. Use `jsonencode()` for better readability. Structure:
  - `name` - (Required, String) Name of the workflow
  - `version` - (Required, Number) Version number of the workflow
  - `steps` - (Required, List) Array of workflow steps, each containing:
    - `name` - (Required, String) Unique name for the step
    - `description` - (Optional, String) Description of what the step does
    - `executor` - (Required, Object) Executor configuration:
      - `type` - (Required, String) Type of executor ("command", "tool", "agent")
      - `config` - (Required, Object) Configuration for the executor
    - `depends` - (Optional, List) Array of step names this step depends on
    - `output` - (Optional, String) Variable name to store step output

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the trigger
* `url` - The webhook URL for triggering the workflow (sensitive)
* `status` - Current status of the workflow (e.g., "draft", "published")
* `workflow_id` - The ID of the created workflow in Kubiya

## Import

Trigger resources can be imported using their ID:

```shell
terraform import kubiya_trigger.example <trigger-id>
```

## Compatibility Notes

* Requires Kubiya Terraform Provider version >= 1.0.0
* Compatible with Terraform >= 1.0
* The workflow is automatically published when the trigger resource is created
* The webhook URL remains stable across updates unless the resource is recreated
* Updating the workflow definition will update the published workflow
* Deleting the trigger resource will delete both the workflow and webhook

## Best Practices

1. **Security**: Store the webhook URL as a sensitive output to prevent accidental exposure
2. **Versioning**: Increment the workflow version when making significant changes
3. **Dependencies**: Use step dependencies to control execution order and parallelism
4. **Error Handling**: Include error checking and recovery steps in your workflows
5. **Monitoring**: Set up logging and monitoring for workflow executions
6. **Testing**: Test workflows in a non-production environment first
7. **Documentation**: Use descriptive names and descriptions for steps to improve maintainability
8. **Idempotency**: Design workflows to be idempotent where possible