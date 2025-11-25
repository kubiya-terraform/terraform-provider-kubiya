---
page_title: "kubiya_inline_source Resource - Kubiya"
subcategory: ""
description: |-
  The kubiya_inline_source resource manages inline tools and workflows in the Kubiya platform.
---

# kubiya_inline_source (Resource)

The `kubiya_inline_source` resource allows you to define tools and workflows directly in your Terraform configuration. This is ideal for custom tools, quick prototypes, and workflows that don't require a separate Git repository. For Git-based sources, use the `kubiya_source` resource instead.

## Container-First Architecture

Kubiya uses a container-first architecture where every tool is backed by a Docker image. This ensures:
- Secure, isolated execution environments
- Predictable and reproducible results
- Language-agnostic tool implementation
- Easy integration with existing containerized workflows

## Prerequisites

Before using this resource, ensure you have:
1. A Kubiya account with API access
2. An API key (generated from Kubiya dashboard under Admin → Kubiya API Keys)
3. Docker images accessible from your runner environment
4. Understanding of tool and workflow structure in Kubiya

## Example Usage

### 1. Basic Inline Tool

Create a simple inline tool:

```hcl
terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
    }
  }
}

provider "kubiya" {
  # API key is automatically read from KUBIYA_API_KEY environment variable
}

resource "kubiya_inline_source" "hello_tool" {
  name   = "hello-world-tool"
  runner = "kubiya-hosted"
  
  tools = jsonencode([
    {
      name        = "say-hello"
      description = "A simple greeting tool"
      type        = "docker"
      image       = "alpine:latest"
      content     = "echo 'Hello, World!'"
    }
  ])
}

resource "kubiya_agent" "basic_agent" {
  name         = "greeting-agent"
  runner       = "kubiya-hosted"
  description  = "Agent with hello world tool"
  instructions = "You are a friendly agent that can greet users."
  
  sources = [kubiya_inline_source.hello_tool.id]
}
```

**Expected Outcome**: Creates an inline source with a simple greeting tool.

### 2. System Monitoring Tools

Define system monitoring and diagnostic tools:

```hcl
resource "kubiya_inline_source" "monitoring_tools" {
  name   = "system-monitoring-tools"
  runner = "kubiya-hosted"
  
  tools = jsonencode([
    {
      name        = "system-health-check"
      description = "Check system health metrics"
      type        = "docker"
      image       = "alpine:latest"
      content     = <<-BASH
        echo "=== System Health Check ==="
        echo "CPU Usage:"
        top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1
        echo ""
        echo "Memory Usage:"
        free -m | awk 'NR==2{printf "%.2f%%\n", $3*100/$2}'
        echo ""
        echo "Disk Usage:"
        df -h | grep -E '^/dev/' | awk '{print $1 " - " $5}'
      BASH
    },
    {
      name        = "log-analyzer"
      description = "Analyze application logs for errors"
      type        = "docker"
      image       = "python:3.11-slim"
      with_files = [
        {
          destination = "/analyzer.py"
          content     = <<-PYTHON
            import re
            import sys
            from collections import Counter
            
            def analyze_logs(file_path):
                error_pattern = re.compile(r'(ERROR|CRITICAL|FATAL)', re.IGNORECASE)
                warn_pattern = re.compile(r'(WARNING|WARN)', re.IGNORECASE)
                
                errors = 0
                warnings = 0
                
                try:
                    with open(file_path, 'r') as f:
                        for line in f:
                            if error_pattern.search(line):
                                errors += 1
                            elif warn_pattern.search(line):
                                warnings += 1
                    
                    print(f"Log Analysis Results:")
                    print(f"Errors: {errors}")
                    print(f"Warnings: {warnings}")
                    print(f"Total Issues: {errors + warnings}")
                    
                except FileNotFoundError:
                    print(f"File not found: {file_path}")
                except Exception as e:
                    print(f"Error analyzing logs: {e}")
            
            if __name__ == "__main__":
                log_file = sys.argv[1] if len(sys.argv) > 1 else "/var/log/app.log"
                analyze_logs(log_file)
          PYTHON
        }
      ]
      content = "python /analyzer.py $${LOG_FILE}"
      args = [
        {
          name        = "LOG_FILE"
          type        = "string"
          description = "Path to log file to analyze"
          required    = false
          default     = "/var/log/app.log"
        }
      ]
    }
  ])
}

resource "kubiya_agent" "monitoring_agent" {
  name         = "monitoring-specialist"
  runner       = "kubiya-hosted"
  description  = "System monitoring and diagnostics agent"
  instructions = "You are a monitoring specialist with tools for system health checks and log analysis."
  
  sources = [kubiya_inline_source.monitoring_tools.id]
}
```

**Expected Outcome**: Creates monitoring tools for system health checks and log analysis.

### 3. Deployment Workflows

Create deployment automation workflows:

```hcl
resource "kubiya_inline_source" "deployment_workflows" {
  name   = "deployment-automation"
  runner = "kubiya-hosted"
  
  workflows = jsonencode([
    {
      name        = "blue-green-deployment"
      description = "Blue-green deployment strategy"
      steps = [
        {
          name        = "validate-config"
          description = "Validate deployment configuration"
          executor = {
            type = "command"
            config = {
              command = "kubectl apply --dry-run=client -f deployment.yaml"
            }
          }
        },
        {
          name        = "deploy-green"
          description = "Deploy green environment"
          depends     = ["validate-config"]
          executor = {
            type = "command"
            config = {
              command = "kubectl apply -f green-deployment.yaml"
            }
          }
        },
        {
          name        = "health-check"
          description = "Check green environment health"
          depends     = ["deploy-green"]
          executor = {
            type = "command"
            config = {
              command = "kubectl wait --for=condition=ready pod -l version=green --timeout=300s"
            }
          }
          output = "HEALTH_STATUS"
        },
        {
          name        = "switch-traffic"
          description = "Switch traffic to green"
          depends     = ["health-check"]
          executor = {
            type = "command"
            config = {
              command = "kubectl patch service app -p '{\"spec\":{\"selector\":{\"version\":\"green\"}}}'"
            }
          }
        }
      ]
    }
  ])
}

resource "kubiya_agent" "deployment_agent" {
  name         = "deployment-orchestrator"
  runner       = "kubiya-hosted"
  description  = "Deployment automation agent"
  instructions = "You orchestrate deployments using blue-green strategy."
  
  sources = [kubiya_inline_source.deployment_workflows.id]
}
```

**Expected Outcome**: Creates workflows for automated blue-green deployments.

### 4. Data Processing Pipeline

Create a data processing workflow with multiple steps:

```hcl
resource "kubiya_inline_source" "data_pipeline" {
  name   = "data-processing-pipeline"
  runner = "kubiya-hosted"
  
  workflows = jsonencode([
    {
      name        = "etl-pipeline"
      description = "Extract, Transform, and Load data pipeline"
      steps = [
        {
          name        = "extract-data"
          description = "Extract data from source"
          executor = {
            type = "tool"
            config = {
              tool_def = {
                name        = "data-extractor"
                description = "Extract data from API"
                type        = "docker"
                image       = "python:3.11-slim"
                with_files = [
                  {
                    destination = "/extract.py"
                    content     = <<-PYTHON
                      import json
                      import random
                      
                      # Simulate data extraction
                      data = {
                          "records": [
                              {"id": i, "value": random.randint(100, 1000)}
                              for i in range(10)
                          ]
                      }
                      
                      print(json.dumps(data))
                    PYTHON
                  }
                ]
                content = "python /extract.py"
              }
            }
          }
          output = "RAW_DATA"
        },
        {
          name        = "transform-data"
          description = "Transform extracted data"
          depends     = ["extract-data"]
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
                      import os
                      import json
                      
                      raw_data = os.environ.get('data', '{}')
                      data = json.loads(raw_data)
                      
                      # Transform data
                      if 'records' in data:
                          for record in data['records']:
                              record['transformed'] = True
                              record['value_doubled'] = record.get('value', 0) * 2
                      
                      print(json.dumps(data))
                    PYTHON
                  }
                ]
                content = "python /transform.py"
                args = [
                  {
                    name        = "data"
                    type        = "string"
                    description = "Raw data to transform"
                    required    = true
                  }
                ]
              }
              args = {
                data = "$${RAW_DATA}"
              }
            }
          }
          output = "TRANSFORMED_DATA"
        },
        {
          name        = "load-data"
          description = "Load data to destination"
          depends     = ["transform-data"]
          executor = {
            type = "agent"
            config = {
              teammate_name = "data-loader"
              message       = "Load the following data: $${TRANSFORMED_DATA}"
            }
          }
        }
      ]
    }
  ])
}
```

**Expected Outcome**: Creates a complete ETL pipeline workflow.

### 5. Mixed Tools and Workflows

Combine tools and workflows in one source:

```hcl
resource "kubiya_inline_source" "devops_toolkit" {
  name   = "complete-devops-toolkit"
  runner = "kubiya-hosted"
  
  tools = jsonencode([
    {
      name        = "k8s-diagnostics"
      description = "Kubernetes cluster diagnostics"
      type        = "docker"
      image       = "bitnami/kubectl:latest"
      content     = <<-BASH
        echo "=== Cluster Status ==="
        kubectl cluster-info
        echo ""
        echo "=== Node Status ==="
        kubectl get nodes
        echo ""
        echo "=== Pod Issues ==="
        kubectl get pods --all-namespaces | grep -v Running
      BASH
    },
    {
      name        = "docker-cleanup"
      description = "Clean up Docker resources"
      type        = "docker"
      image       = "docker:latest"
      content     = <<-BASH
        echo "Cleaning up Docker resources..."
        docker system prune -f
        docker image prune -a -f
        echo "Cleanup completed"
      BASH
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
          name    = "cleanup"
          depends = ["diagnose"]
          executor = {
            type = "tool"
            config = {
              tool_name = "docker-cleanup"
            }
          }
        },
        {
          name    = "notify"
          depends = ["cleanup"]
          executor = {
            type = "agent"
            config = {
              teammate_name = "slack-notifier"
              message       = "Incident resolved. Diagnosis: $${DIAGNOSIS}"
            }
          }
        }
      ]
    }
  ])
}

resource "kubiya_agent" "devops_agent" {
  name         = "devops-specialist"
  runner       = "kubiya-hosted"
  description  = "DevOps agent with tools and workflows"
  instructions = "You are a DevOps specialist with diagnostic tools and incident response workflows."
  
  sources = [kubiya_inline_source.devops_toolkit.id]
}
```

**Expected Outcome**: Creates a comprehensive toolkit with both tools and workflows.

### 6. API Testing Tools

Create API testing and validation tools:

```hcl
resource "kubiya_inline_source" "api_testing" {
  name   = "api-testing-tools"
  runner = "kubiya-hosted"
  
  tools = jsonencode([
    {
      name        = "api-health-check"
      description = "Check API endpoint health"
      type        = "docker"
      image       = "curlimages/curl:latest"
      content     = <<-BASH
        curl -s -o /dev/null -w "Status: %{http_code}\nResponse Time: %{time_total}s\n" $${API_URL}
      BASH
      args = [
        {
          name        = "API_URL"
          type        = "string"
          description = "API endpoint URL"
          required    = true
        }
      ]
    },
    {
      name        = "load-test"
      description = "Simple load testing"
      type        = "docker"
      image       = "alpine:latest"
      content     = <<-BASH
        apk add --no-cache curl
        
        echo "Starting load test..."
        for i in $(seq 1 $${REQUESTS}); do
          curl -s -o /dev/null -w "%{http_code} " $${API_URL}
        done
        echo ""
        echo "Load test completed: $${REQUESTS} requests sent"
      BASH
      args = [
        {
          name        = "API_URL"
          type        = "string"
          description = "API endpoint URL"
          required    = true
        },
        {
          name        = "REQUESTS"
          type        = "string"
          description = "Number of requests"
          required    = false
          default     = "10"
        }
      ]
    }
  ])
}
```

**Expected Outcome**: Creates API testing and validation tools.

### 7. Database Management Tools

Create database management and migration tools:

```hcl
resource "kubiya_inline_source" "db_tools" {
  name   = "database-management-tools"
  runner = "kubiya-hosted"
  
  tools = jsonencode([
    {
      name        = "db-backup"
      description = "Backup PostgreSQL database"
      type        = "docker"
      image       = "postgres:15-alpine"
      content     = <<-BASH
        PGPASSWORD=$${DB_PASSWORD} pg_dump \
          -h $${DB_HOST} \
          -U $${DB_USER} \
          -d $${DB_NAME} \
          > /backup/backup_$(date +%Y%m%d_%H%M%S).sql
        
        echo "Backup completed successfully"
      BASH
      args = [
        {
          name        = "DB_HOST"
          type        = "string"
          description = "Database host"
          required    = true
        },
        {
          name        = "DB_NAME"
          type        = "string"
          description = "Database name"
          required    = true
        },
        {
          name        = "DB_USER"
          type        = "string"
          description = "Database user"
          required    = true
        },
        {
          name        = "DB_PASSWORD"
          type        = "string"
          description = "Database password"
          required    = true
        }
      ]
    }
  ])
  
  workflows = jsonencode([
    {
      name        = "database-migration"
      description = "Database migration workflow"
      steps = [
        {
          name = "backup-current"
          executor = {
            type = "tool"
            config = {
              tool_name = "db-backup"
            }
          }
        },
        {
          name    = "run-migration"
          depends = ["backup-current"]
          executor = {
            type = "command"
            config = {
              command = "flyway migrate"
            }
          }
        },
        {
          name    = "verify-migration"
          depends = ["run-migration"]
          executor = {
            type = "command"
            config = {
              command = "flyway info"
            }
          }
        }
      ]
    }
  ])
}
```

**Expected Outcome**: Creates database management tools and migration workflows.

### 8. Security Scanning Tools

Create security scanning and compliance tools:

```hcl
resource "kubiya_inline_source" "security_tools" {
  name   = "security-scanning-tools"
  runner = "kubiya-hosted"
  
  tools = jsonencode([
    {
      name        = "dependency-check"
      description = "Check for vulnerable dependencies"
      type        = "docker"
      image       = "owasp/dependency-check:latest"
      content     = <<-BASH
        dependency-check.sh \
          --scan /src \
          --format JSON \
          --out /tmp/report.json \
          --suppression /tmp/suppressions.xml
        
        cat /tmp/report.json | jq '.dependencies[] | select(.vulnerabilities != null)'
      BASH
    },
    {
      name        = "secrets-scan"
      description = "Scan for exposed secrets"
      type        = "docker"
      image       = "trufflesecurity/trufflehog:latest"
      content     = <<-BASH
        trufflehog filesystem /src \
          --json \
          --only-verified
      BASH
    }
  ])
  
  workflows = jsonencode([
    {
      name        = "security-audit"
      description = "Complete security audit"
      steps = [
        {
          name = "scan-dependencies"
          executor = {
            type = "tool"
            config = {
              tool_name = "dependency-check"
            }
          }
          output = "DEPENDENCY_REPORT"
        },
        {
          name = "scan-secrets"
          executor = {
            type = "tool"
            config = {
              tool_name = "secrets-scan"
            }
          }
          output = "SECRETS_REPORT"
        },
        {
          name    = "generate-report"
          depends = ["scan-dependencies", "scan-secrets"]
          executor = {
            type = "agent"
            config = {
              teammate_name = "security-reporter"
              message       = "Generate security report from: Dependencies: $${DEPENDENCY_REPORT}, Secrets: $${SECRETS_REPORT}"
            }
          }
        }
      ]
    }
  ])
}
```

**Expected Outcome**: Creates security scanning tools and audit workflows.

## Argument Reference

### Required Arguments

* `name` - (Required, String) The name of the inline source. Must be unique within your organization.

### Optional Arguments

* `runner` - (Optional, String) The runner to use for executing tools and workflows. Defaults to "kubiya-hosted".
* `tools` - (Optional, String) JSON-encoded array of inline tool definitions.
* `workflows` - (Optional, String) JSON-encoded array of workflow definitions.
* `dynamic_config` - (Optional, String) JSON-encoded configuration for dynamic parameters.

### Tool Definition Structure

Each tool in the `tools` array should have:
* `name` - Tool identifier
* `description` - Tool description
* `type` - Execution type (usually "docker")
* `image` - Docker image to use
* `content` - Command or script to execute
* `with_files` - (Optional) Files to create in the container:
  - `destination` - File path in container
  - `content` - File content
* `args` - (Optional) Tool arguments:
  - `name` - Argument name
  - `type` - Argument type
  - `description` - Argument description
  - `required` - Whether required
  - `default` - Default value

### Workflow Definition Structure

Each workflow in the `workflows` array should have:
* `name` - Workflow name
* `description` - Workflow description
* `steps` - Array of workflow steps:
  - `name` - Step name
  - `description` - Step description (optional)
  - `executor` - Executor configuration:
    - `type` - Executor type ("tool", "command", "agent")
    - `config` - Executor-specific configuration
  - `depends` - (Optional) Array of step names this step depends on
  - `output` - (Optional) Variable name to store output
  - `args` - (Optional) Arguments to pass to the executor

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the inline source.
* `type` - The computed type of the source (always "inline").

## Import

Inline sources can be imported using their ID:

```shell
terraform import kubiya_inline_source.example <inline-source-id>
```

## Compatibility Notes

* Requires Kubiya Terraform Provider version >= 1.0.0
* Compatible with Terraform >= 1.0
* Docker images must be accessible from the runner environment
* Tools and workflows are defined inline, not from Git repositories
* Inline sources must be created before agents can reference them

## Best Practices

1. **Image Management**: Use specific image tags rather than "latest" for reproducibility
2. **Security**: Never hardcode credentials in tool definitions; use Kubiya secrets
3. **Testing**: Test tools and workflows in development before production use
4. **Documentation**: Include clear descriptions for tools and their parameters
5. **Error Handling**: Include error handling in tool scripts
6. **Resource Limits**: Consider resource requirements when choosing Docker images
7. **Modularity**: Create focused, single-purpose tools that can be composed in workflows
8. **Version Control**: Consider moving complex tools to Git repositories as they mature
9. **Debugging**: Use echo statements and output variables for debugging workflows
10. **Idempotency**: Design tools to be idempotent where possible