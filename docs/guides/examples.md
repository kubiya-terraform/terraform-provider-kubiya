---
page_title: "Comprehensive Examples Guide - Kubiya"
subcategory: ""
description: |-
  Complete examples demonstrating all Kubiya Terraform Provider resource configurations.
---

# Comprehensive Examples Guide

This guide provides detailed examples for all Kubiya Terraform Provider resources, organized by use case.

## Table of Contents

1. [Provider Setup](#provider-setup)
2. [Agent Examples](#agent-examples)
3. [Webhook Examples](#webhook-examples)
4. [Trigger Examples](#trigger-examples)
5. [Complete Solutions](#complete-solutions)

## Provider Setup

### Prerequisites

1. **Kubiya Account**: Sign up at [kubiya.ai](https://kubiya.ai)
2. **API Key**: Generate from Kubiya dashboard (Admin → Kubiya API Keys)
3. **Terraform**: Version 1.0 or higher

### Basic Configuration

```hcl
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    kubiya = {
      source  = "kubiya-terraform/kubiya"
      version = ">= 1.0.0"
    }
  }
}

provider "kubiya" {
  # API key is automatically read from KUBIYA_API_KEY environment variable
}
```

### Environment Setup

```bash
# Set API key
export KUBIYA_API_KEY="your-api-key-here"

# Initialize Terraform
terraform init

# Plan changes
terraform plan

# Apply configuration
terraform apply
```

## Agent Examples

### 1. Basic Agent

Simple agent with minimal configuration:

```hcl
resource "kubiya_agent" "basic" {
  name         = "basic-assistant"
  runner       = "kubiya-hosted"
  description  = "A helpful AI assistant"
  instructions = "You are a helpful assistant. Provide clear and concise responses."
}
```

### 2. Agent with Inline Tools

Agent with custom tools defined inline:

```hcl
resource "kubiya_source" "custom_tools" {
  name = "devops-tools"
  
  tools = jsonencode([
    {
      name        = "system-info"
      description = "Get system information"
      type        = "docker"
      image       = "alpine:latest"
      content     = <<-BASH
        echo "System Information:"
        echo "==================="
        uname -a
        echo ""
        echo "CPU Info:"
        nproc
        echo ""
        echo "Memory Info:"
        free -h
        echo ""
        echo "Disk Usage:"
        df -h
      BASH
    },
    {
      name        = "log-analyzer"
      description = "Analyze application logs"
      type        = "docker"
      image       = "python:3.11-slim"
      with_files = [
        {
          destination = "/analyzer.py"
          content     = <<-PYTHON
            import re
            import sys
            from collections import Counter
            
            def analyze_logs(log_file):
                error_patterns = {
                    'critical': r'CRITICAL|FATAL',
                    'error': r'ERROR',
                    'warning': r'WARNING|WARN',
                    'info': r'INFO'
                }
                
                counts = Counter()
                
                with open(log_file, 'r') as f:
                    for line in f:
                        for level, pattern in error_patterns.items():
                            if re.search(pattern, line, re.IGNORECASE):
                                counts[level] += 1
                
                print("Log Analysis Summary")
                print("=" * 40)
                for level, count in counts.items():
                    print(f"{level.upper()}: {count}")
                
                return counts
            
            if __name__ == "__main__":
                if len(sys.argv) > 1:
                    analyze_logs(sys.argv[1])
                else:
                    print("Usage: python analyzer.py <log_file>")
          PYTHON
        }
      ]
      content = "python /analyzer.py /var/log/app.log"
      args = [
        {
          name        = "log_file"
          type        = "string"
          description = "Path to log file"
          required    = false
          default     = "/var/log/app.log"
        }
      ]
    }
  ])
  
  runner = "kubiya-hosted"
}

resource "kubiya_agent" "tooled_agent" {
  name         = "devops-assistant"
  runner       = "kubiya-hosted"
  description  = "DevOps agent with custom tools"
  instructions = "You are a DevOps assistant with access to system monitoring and log analysis tools."
  
  tool_sources = [kubiya_source.custom_tools.id]
  
  integrations = ["slack_integration"]
}
```

### 3. Webhook with Agent

Webhook that triggers an agent for incident response:

```hcl
resource "kubiya_agent" "incident_handler" {
  name         = "incident-response-agent"
  runner       = "kubiya-hosted"
  description  = "Automated incident response"
  instructions = <<-EOT
    You are an incident response agent. When triggered:
    1. Assess the severity of the incident
    2. Gather relevant logs and metrics
    3. Create tickets in JIRA
    4. Notify the appropriate teams
    5. Suggest remediation steps
  EOT
  
  integrations = [
    "pagerduty",
    "jira_cloud",
    "datadog",
    "slack_integration"
  ]
  
  environment_variables = {
    INCIDENT_CHANNEL = "#incidents"
    JIRA_PROJECT     = "OPS"
  }
}

resource "kubiya_webhook" "incident_webhook" {
  name        = "datadog-incident-webhook"
  agent       = kubiya_agent.incident_handler.name
  source      = "datadog"
  filter      = "alert_type:error"
  prompt      = "Handle Datadog alert: analyze the issue and initiate incident response"
  
  method      = "Slack"
  destination = "#incidents"
}

output "incident_webhook_url" {
  value       = kubiya_webhook.incident_webhook.url
  sensitive   = true
  description = "Webhook URL for Datadog integration"
}
```

### 4. Agent with Workflow

Agent that executes predefined workflows:

```hcl
resource "kubiya_source" "deployment_workflows" {
  name = "deployment-automation"
  
  workflows = jsonencode([
    {
      name        = "blue-green-deployment"
      description = "Perform blue-green deployment"
      steps = [
        {
          name        = "validate-prerequisites"
          description = "Check deployment prerequisites"
          executor = {
            type = "command"
            config = {
              command = <<-BASH
                echo "Checking prerequisites..."
                kubectl get nodes
                kubectl get deployments -n production
              BASH
            }
          }
        },
        {
          name        = "deploy-green"
          description = "Deploy green environment"
          depends     = ["validate-prerequisites"]
          executor = {
            type = "command"
            config = {
              command = "kubectl apply -f green-deployment.yaml -n production"
            }
          }
        },
        {
          name        = "health-check"
          description = "Verify green deployment health"
          depends     = ["deploy-green"]
          executor = {
            type = "command"
            config = {
              command = "kubectl wait --for=condition=ready pod -l version=green -n production --timeout=300s"
            }
          }
        },
        {
          name        = "switch-traffic"
          description = "Route traffic to green"
          depends     = ["health-check"]
          executor = {
            type = "command"
            config = {
              command = "kubectl patch service app-service -n production -p '{\"spec\":{\"selector\":{\"version\":\"green\"}}}'"
            }
          }
        },
        {
          name        = "cleanup-blue"
          description = "Remove old blue deployment"
          depends     = ["switch-traffic"]
          executor = {
            type = "command"
            config = {
              command = "kubectl delete deployment app-blue -n production"
            }
          }
        }
      ]
    },
    {
      name        = "rollback-deployment"
      description = "Rollback to previous version"
      steps = [
        {
          name = "rollback"
          executor = {
            type = "command"
            config = {
              command = "kubectl rollout undo deployment/app -n production"
            }
          }
        },
        {
          name = "verify-rollback"
          depends = ["rollback"]
          executor = {
            type = "command"
            config = {
              command = "kubectl rollout status deployment/app -n production"
            }
          }
        }
      ]
    }
  ])
  
  runner = "kubiya-hosted"
}

resource "kubiya_agent" "deployment_agent" {
  name         = "deployment-orchestrator"
  runner       = "kubiya-hosted"
  description  = "Manages application deployments"
  instructions = <<-EOT
    You are a deployment orchestrator. You can:
    1. Execute blue-green deployments
    2. Perform rollbacks when issues are detected
    3. Monitor deployment health
    4. Send status updates
    
    Available workflows:
    - blue-green-deployment: For safe production deployments
    - rollback-deployment: For emergency rollbacks
  EOT
  
  sources = [kubiya_source.deployment_workflows.id]
  
  integrations = ["kubernetes", "slack_integration"]
  
  starters = [
    {
      name    = "Deploy Latest"
      command = "Execute blue-green deployment for the latest version"
    },
    {
      name    = "Rollback"
      command = "Rollback to the previous stable version"
    },
    {
      name    = "Check Status"
      command = "Show current deployment status"
    }
  ]
}
```

## Webhook Examples

### 5. Webhook with Workflow

Webhook that directly executes a workflow:

```hcl
resource "kubiya_webhook" "etl_webhook" {
  name   = "data-processing-webhook"
  source = "s3"
  filter = "event:ObjectCreated"
  prompt = "Process new data file"
  
  workflow = jsonencode({
    name        = "data-etl-pipeline"
    description = "Process incoming data files"
    steps = [
      {
        name        = "validate-file"
        description = "Validate incoming file format"
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "file-validator"
              description = "Validate data file"
              type        = "docker"
              image       = "python:3.11-slim"
              with_files = [
                {
                  destination = "/validate.py"
                  content     = <<-PYTHON
                    import json
                    import sys
                    
                    def validate_json(file_path):
                        try:
                            with open(file_path, 'r') as f:
                                data = json.load(f)
                            
                            required_fields = ['id', 'timestamp', 'data']
                            for field in required_fields:
                                if field not in data:
                                    raise ValueError(f"Missing required field: {field}")
                            
                            print(f"✓ File validated successfully")
                            print(f"  Records: {len(data.get('data', []))}")
                            return True
                        except Exception as e:
                            print(f"✗ Validation failed: {e}")
                            return False
                    
                    if __name__ == "__main__":
                        validate_json(sys.argv[1] if len(sys.argv) > 1 else "/tmp/data.json")
                  PYTHON
                }
              ]
              content = "python /validate.py $${FILE_PATH}"
              args = [
                {
                  name        = "FILE_PATH"
                  type        = "string"
                  description = "Path to file to validate"
                  required    = true
                }
              ]
            }
            args = {
              FILE_PATH = "$${S3_FILE_PATH}"
            }
          }
        }
        output = "VALIDATION_STATUS"
      },
      {
        name        = "transform-data"
        description = "Transform and enrich data"
        depends     = ["validate-file"]
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "data-transformer"
              description = "Transform data"
              type        = "docker"
              image       = "python:3.11-slim"
              with_files = [
                {
                  destination = "/transform.py"
                  content     = <<-PYTHON
                    import json
                    import pandas as pd
                    from datetime import datetime
                    
                    # Load data
                    with open('/tmp/data.json', 'r') as f:
                        raw_data = json.load(f)
                    
                    # Transform
                    df = pd.DataFrame(raw_data['data'])
                    df['processed_at'] = datetime.now().isoformat()
                    df['quality_score'] = df['value'].apply(lambda x: 'high' if x > 80 else 'low')
                    
                    # Save
                    df.to_parquet('/tmp/transformed.parquet')
                    print(f"Transformed {len(df)} records")
                  PYTHON
                }
              ]
              content = "python /transform.py"
            }
          }
        }
        output = "TRANSFORM_STATUS"
      },
      {
        name        = "load-to-warehouse"
        description = "Load to data warehouse"
        depends     = ["transform-data"]
        executor = {
          type = "command"
          config = {
            command = "aws s3 cp /tmp/transformed.parquet s3://data-warehouse/processed/$(date +%Y%m%d)/"
          }
        }
      },
      {
        name        = "notify-completion"
        description = "Send completion notification"
        depends     = ["load-to-warehouse"]
        executor = {
          type = "agent"
          config = {
            teammate_name = "notification-agent"
            message       = "ETL pipeline completed. Status: $${TRANSFORM_STATUS}"
          }
        }
      }
    ]
  })
  
  runner      = "kubiya-hosted"
  method      = "Slack"
  destination = "#data-team"
}
```

## Trigger Examples

### 6. Trigger with Workflow

HTTP trigger for workflow execution:

```hcl
resource "kubiya_trigger" "cicd_trigger" {
  name   = "github-cicd-pipeline"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "CI/CD Pipeline"
    version = 1
    steps = [
      {
        name        = "checkout"
        description = "Checkout code"
        executor = {
          type = "command"
          config = {
            command = "git clone $${REPO_URL} /tmp/repo && cd /tmp/repo && git checkout $${BRANCH}"
          }
        }
      },
      {
        name        = "test"
        description = "Run tests"
        depends     = ["checkout"]
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "test-runner"
              description = "Execute test suite"
              type        = "docker"
              image       = "node:18"
              content     = "cd /tmp/repo && npm install && npm test"
            }
          }
        }
        output = "TEST_RESULTS"
      },
      {
        name        = "build"
        description = "Build application"
        depends     = ["test"]
        executor = {
          type = "command"
          config = {
            command = "cd /tmp/repo && docker build -t app:$${VERSION} ."
          }
        }
      },
      {
        name        = "scan"
        description = "Security scan"
        depends     = ["build"]
        executor = {
          type = "tool"
          config = {
            tool_def = {
              name        = "security-scanner"
              description = "Scan for vulnerabilities"
              type        = "docker"
              image       = "aquasec/trivy"
              content     = "trivy image app:$${VERSION}"
            }
          }
        }
        output = "SCAN_RESULTS"
      },
      {
        name        = "deploy"
        description = "Deploy to environment"
        depends     = ["scan"]
        executor = {
          type = "command"
          config = {
            command = <<-BASH
              kubectl set image deployment/app app=app:$${VERSION} -n $${ENVIRONMENT}
              kubectl rollout status deployment/app -n $${ENVIRONMENT}
            BASH
          }
        }
      }
    ]
  })
}

output "cicd_trigger_url" {
  value       = kubiya_trigger.cicd_trigger.url
  sensitive   = true
  description = "CI/CD pipeline trigger URL"
}
```

## Complete Solutions

### 7. Full DevOps Platform

Complete DevOps automation platform:

```hcl
# Shared tools and workflows
resource "kubiya_source" "devops_tools" {
  name = "devops-toolkit"
  
  tools = jsonencode([
    {
      name        = "k8s-diagnostics"
      description = "Kubernetes cluster diagnostics"
      type        = "docker"
      image       = "bitnami/kubectl:latest"
      content     = <<-BASH
        echo "=== Cluster Health Check ==="
        kubectl get nodes
        echo ""
        echo "=== Pod Status ==="
        kubectl get pods --all-namespaces | grep -v Running
        echo ""
        echo "=== Resource Usage ==="
        kubectl top nodes
        kubectl top pods --all-namespaces | head -20
      BASH
    },
    {
      name        = "database-health"
      description = "Check database health"
      type        = "docker"
      image       = "postgres:14"
      content     = <<-SQL
        psql -h $${DB_HOST} -U $${DB_USER} -d $${DB_NAME} -c "
          SELECT 
            pg_database.datname,
            pg_size_pretty(pg_database_size(pg_database.datname)) as size,
            count(pg_stat_activity.pid) as connections
          FROM pg_database
          LEFT JOIN pg_stat_activity ON pg_database.datname = pg_stat_activity.datname
          GROUP BY pg_database.datname
          ORDER BY pg_database_size(pg_database.datname) DESC;
        "
      SQL
      args = [
        {
          name        = "DB_HOST"
          type        = "string"
          description = "Database host"
          required    = true
        },
        {
          name        = "DB_USER"
          type        = "string"
          description = "Database user"
          required    = true
        },
        {
          name        = "DB_NAME"
          type        = "string"
          description = "Database name"
          required    = true
        }
      ]
    }
  ])
  
  workflows = jsonencode([
    {
      name        = "incident-response"
      description = "Automated incident response"
      steps = [
        {
          name = "diagnose"
          executor = {
            type = "tool"
            config = {
              tool_name = "k8s-diagnostics"
            }
          }
          output = "DIAGNOSIS"
        },
        {
          name    = "create-ticket"
          depends = ["diagnose"]
          executor = {
            type = "agent"
            config = {
              teammate_name = "jira-agent"
              message       = "Create incident ticket with diagnosis: $${DIAGNOSIS}"
            }
          }
          output = "TICKET_ID"
        },
        {
          name    = "notify"
          depends = ["create-ticket"]
          executor = {
            type = "agent"
            config = {
              teammate_name = "slack-agent"
              message       = "Incident $${TICKET_ID} created. Diagnosis: $${DIAGNOSIS}"
            }
          }
        }
      ]
    }
  ])
  
  runner = "kubiya-hosted"
}

# Central DevOps Agent
resource "kubiya_agent" "devops_central" {
  name         = "devops-central"
  runner       = "kubiya-hosted"
  description  = "Central DevOps automation agent"
  instructions = <<-EOT
    You are the central DevOps automation agent. You can:
    1. Monitor infrastructure health
    2. Respond to incidents
    3. Execute deployments
    4. Manage databases
    5. Coordinate with other agents
    
    Use available tools and workflows to automate tasks.
  EOT
  
  tool_sources = [kubiya_source.devops_tools.id]
  sources      = [kubiya_source.devops_tools.id]
  
  integrations = [
    "kubernetes",
    "github_app",
    "slack_integration",
    "jira_cloud",
    "datadog",
    "pagerduty"
  ]
  
  users  = ["devops@company.com", "sre@company.com"]
  groups = ["DevOps", "SRE", "Platform"]
  
  environment_variables = {
    DEFAULT_NAMESPACE = "production"
    SLACK_CHANNEL     = "#devops"
    JIRA_PROJECT      = "OPS"
  }
  
  starters = [
    {
      name    = "Health Check"
      command = "Run full infrastructure health check"
    },
    {
      name    = "Deploy App"
      command = "Deploy application to production"
    },
    {
      name    = "Incident Response"
      command = "Start incident response workflow"
    }
  ]
}

# Monitoring Webhook
resource "kubiya_webhook" "monitoring" {
  name        = "monitoring-alerts"
  agent       = kubiya_agent.devops_central.name
  source      = "datadog"
  filter      = "severity:high"
  prompt      = "Handle monitoring alert"
  
  method      = "Slack"
  destination = "#alerts"
}

# Deployment Trigger
resource "kubiya_trigger" "deployment" {
  name   = "deploy-trigger"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "Deployment Pipeline"
    version = 1
    steps = [
      {
        name = "deploy"
        executor = {
          type = "agent"
          config = {
            teammate_name = kubiya_agent.devops_central.name
            message       = "Deploy version $${VERSION} to $${ENVIRONMENT}"
          }
        }
      }
    ]
  })
}

# Outputs
output "devops_webhook_url" {
  value       = kubiya_webhook.monitoring.url
  sensitive   = true
  description = "Monitoring webhook URL"
}

output "deployment_trigger_url" {
  value       = kubiya_trigger.deployment.url
  sensitive   = true
  description = "Deployment trigger URL"
}
```

### 8. Data Processing Platform

Complete data processing and analytics platform:

```hcl
# Data processing workflows
resource "kubiya_source" "data_workflows" {
  name = "data-processing-workflows"
  
  workflows = jsonencode([
    {
      name        = "daily-etl"
      description = "Daily ETL pipeline"
      steps = [
        {
          name = "extract"
          executor = {
            type = "command"
            config = {
              command = "aws s3 sync s3://raw-data/$(date +%Y/%m/%d)/ /tmp/raw/"
            }
          }
        },
        {
          name    = "transform"
          depends = ["extract"]
          executor = {
            type = "tool"
            config = {
              tool_def = {
                name        = "data-transformer"
                description = "Transform raw data"
                type        = "docker"
                image       = "apache/spark:3.4.0"
                content     = <<-SPARK
                  spark-submit --master local[*] /transform.py \
                    --input /tmp/raw \
                    --output /tmp/processed \
                    --date $(date +%Y-%m-%d)
                SPARK
              }
            }
          }
        },
        {
          name    = "load"
          depends = ["transform"]
          executor = {
            type = "command"
            config = {
              command = "aws s3 sync /tmp/processed/ s3://processed-data/$(date +%Y/%m/%d)/"
            }
          }
        },
        {
          name    = "quality-check"
          depends = ["load"]
          executor = {
            type = "tool"
            config = {
              tool_def = {
                name        = "quality-checker"
                description = "Validate data quality"
                type        = "docker"
                image       = "python:3.11"
                content     = <<-PYTHON
                  python -c "
                  import pandas as pd
                  import sys
                  
                  df = pd.read_parquet('/tmp/processed/')
                  
                  # Quality checks
                  assert not df.empty, 'Data is empty'
                  assert df.isnull().sum().sum() == 0, 'Found null values'
                  assert len(df) > 1000, 'Insufficient records'
                  
                  print(f'✓ Quality check passed: {len(df)} records processed')
                  "
                PYTHON
              }
            }
          }
        }
      ]
    }
  ])
  
  runner = "kubiya-hosted"
}

# Data Analytics Agent
resource "kubiya_agent" "data_analyst" {
  name         = "data-analyst"
  runner       = "kubiya-hosted"
  description  = "Data analytics and reporting agent"
  instructions = <<-EOT
    You are a data analytics agent. You can:
    1. Execute ETL pipelines
    2. Generate reports
    3. Perform data quality checks
    4. Create visualizations
    5. Answer data-related questions
  EOT
  
  sources = [kubiya_source.data_workflows.id]
  
  integrations = ["slack_integration"]
  
  environment_variables = {
    DATA_WAREHOUSE = "s3://processed-data"
    REPORT_BUCKET  = "s3://reports"
  }
}

# Scheduled ETL Trigger
resource "kubiya_trigger" "scheduled_etl" {
  name   = "scheduled-etl"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "Scheduled ETL"
    version = 1
    steps = [
      {
        name = "run-etl"
        executor = {
          type = "agent"
          config = {
            teammate_name = kubiya_agent.data_analyst.name
            message       = "Execute daily ETL pipeline"
          }
        }
      }
    ]
  })
}

output "etl_trigger_url" {
  value       = kubiya_trigger.scheduled_etl.url
  sensitive   = true
  description = "ETL trigger URL for cron job"
}
```

## Best Practices

### Security

1. **API Keys**: Always use environment variables for API keys
2. **Secrets**: Use Kubiya secrets for sensitive data
3. **Access Control**: Implement proper user and group restrictions
4. **Audit**: Enable logging for all critical operations

### Organization

1. **Naming Convention**: Use consistent naming (e.g., `environment-service-resource`)
2. **Modules**: Create reusable Terraform modules for common patterns
3. **Version Control**: Store configurations in Git
4. **Documentation**: Document custom tools and workflows

### Testing

1. **Development Environment**: Test in non-production first
2. **Validation**: Validate workflows before deployment
3. **Monitoring**: Set up alerts for failed executions
4. **Rollback Plan**: Always have a rollback strategy

### Performance

1. **Runner Selection**: Choose appropriate runners for workload
2. **Parallel Execution**: Use step dependencies for parallelism
3. **Resource Limits**: Set appropriate resource limits
4. **Caching**: Implement caching where possible

## Troubleshooting

### Common Issues

1. **Authentication Errors**
   ```bash
   # Verify API key
   echo $KUBIYA_API_KEY
   
   # Test connection
   curl -H "Authorization: UserKey $KUBIYA_API_KEY" \
        https://api.kubiya.ai/v1/agents
   ```

2. **Runner Issues**
   ```hcl
   # Use hosted runner for testing
   runner = "kubiya-hosted"
   ```

3. **Workflow Errors**
   ```hcl
   # Enable debug mode
   resource "kubiya_agent" "debug" {
     # ...
     is_debug_mode = true
   }
   ```

### Getting Help

- Documentation: [docs.kubiya.ai](https://docs.kubiya.ai)
- Support: support@kubiya.ai
- Community: [community.kubiya.ai](https://community.kubiya.ai)

## Conclusion

This guide demonstrates the full capabilities of the Kubiya Terraform Provider. Use these examples as templates for your own automation needs, adapting them to your specific requirements and infrastructure.