---
page_title: "kubiya_agent Resource - Kubiya"
subcategory: ""
description: |-
  The kubiya_agent resource manages AI agents in the Kubiya platform.
---

# kubiya_agent (Resource)

The `kubiya_agent` resource allows you to create and manage AI agents in the Kubiya platform. Agents are intelligent assistants that can perform various tasks, integrate with external systems, and execute workflows.

## Prerequisites

Before using this resource, ensure you have:
1. A Kubiya account with API access
2. An API key (generated from Kubiya dashboard under Admin → Kubiya API Keys)
3. At least one runner configured (or use "kubiya-hosted" for cloud execution)

## Example Usage

### 1. Basic Agent Configuration

This example creates a simple agent with minimal configuration:

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

resource "kubiya_agent" "basic_agent" {
  name         = "my-basic-agent"
  runner       = "kubiya-hosted"
  description  = "A basic AI assistant for general tasks"
  instructions = "You are a helpful AI assistant. Provide clear and concise responses to user queries."
}
```

**Expected Outcome**: Creates a basic AI agent accessible to all users in your organization.

### 2. Agent with User and Group Access Control

Configure an agent with specific user and group permissions:

```hcl
resource "kubiya_agent" "team_agent" {
  name         = "team-assistant"
  runner       = "kubiya-hosted"
  description  = "Team collaboration assistant"
  instructions = "You are a team assistant helping with project management and collaboration tasks."
  
  # Access control
  users  = ["alice@company.com", "bob@company.com"]
  groups = ["Engineering", "DevOps"]
  
  # Optional: Use a specific AI model
  model = "gpt-4o"
}
```

**Expected Outcome**: Creates an agent accessible only to specified users and groups.

### 3. DevOps Agent with Integrations

Create an agent with access to external integrations:

```hcl
resource "kubiya_agent" "devops_agent" {
  name         = "devops-assistant"
  runner       = "kubiya-hosted"
  description  = "DevOps automation assistant"
  instructions = <<-EOT
    You are a DevOps assistant specialized in:
    - Managing Kubernetes clusters
    - GitHub repository operations
    - Slack notifications
    - JIRA ticket management
    
    Always confirm before making changes to production systems.
  EOT
  
  integrations = [
    "github_app",
    "slack_integration",
    "jira_cloud",
    "kubernetes"
  ]
  
  # Environment variables for the agent
  environment_variables = {
    DEFAULT_NAMESPACE = "production"
    ALERT_CHANNEL     = "#devops-alerts"
  }
}
```

**Expected Outcome**: Creates a DevOps agent with access to multiple integrations and custom environment variables.

### 4. Agent with Predefined Tasks

Define an agent with specific tasks that users can execute:

```hcl
resource "kubiya_agent" "ops_agent" {
  name         = "operations-agent"
  runner       = "kubiya-hosted"
  description  = "Operations management agent"
  instructions = "You are an operations agent that helps with infrastructure management tasks."
  
  tasks = [
    {
      name        = "check-cluster-health"
      prompt      = "Check the health status of all Kubernetes clusters and report any issues"
      description = "Performs comprehensive health checks on K8s clusters"
    },
    {
      name        = "scale-deployment"
      prompt      = "Scale the specified deployment to the requested number of replicas"
      description = "Scales Kubernetes deployments"
    },
    {
      name        = "backup-database"
      prompt      = "Create a backup of the specified database and store it in S3"
      description = "Performs database backup operations"
    }
  ]
  
  integrations = ["kubernetes", "aws"]
}
```

**Expected Outcome**: Creates an agent with predefined tasks that users can easily trigger.

### 5. Agent with Inline Tools

Create an agent with custom inline tools defined directly in the configuration:

```hcl
resource "kubiya_source" "inline_tools" {
  name = "custom-tools"
  
  tools = jsonencode([
    {
      name        = "disk-usage-checker"
      description = "Check disk usage on servers"
      type        = "docker"
      image       = "alpine:latest"
      content     = "df -h"
    },
    {
      name        = "log-analyzer"
      description = "Analyze application logs"
      type        = "docker"
      image       = "python:3.11-slim"
      with_files = [
        {
          destination = "/analyze.py"
          content     = <<-PYTHON
            import sys
            import json
            
            def analyze_logs(log_file):
                errors = []
                warnings = []
                with open(log_file, 'r') as f:
                    for line in f:
                        if 'ERROR' in line:
                            errors.append(line.strip())
                        elif 'WARNING' in line:
                            warnings.append(line.strip())
                
                return {
                    'total_errors': len(errors),
                    'total_warnings': len(warnings),
                    'errors': errors[:10],
                    'warnings': warnings[:10]
                }
            
            if __name__ == "__main__":
                result = analyze_logs('/var/log/app.log')
                print(json.dumps(result, indent=2))
          PYTHON
        }
      ]
      content = "python /analyze.py"
    }
  ])
  
  runner = "kubiya-hosted"
}

resource "kubiya_agent" "agent_with_tools" {
  name         = "monitoring-agent"
  runner       = "kubiya-hosted"
  description  = "System monitoring agent with custom tools"
  instructions = "You are a monitoring agent with access to custom tools for system analysis."
  
  tool_sources = [kubiya_source.inline_tools.id]
}
```

**Expected Outcome**: Creates an agent with access to custom inline tools for specialized operations.

### 6. Agent with Workflow Execution

Configure an agent that can execute predefined workflows:

```hcl
resource "kubiya_source" "deployment_workflow" {
  name = "deployment-workflows"
  
  workflows = jsonencode([
    {
      name        = "blue-green-deployment"
      description = "Perform blue-green deployment"
      steps = [
        {
          name        = "prepare-green-env"
          description = "Prepare green environment"
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
          depends     = ["prepare-green-env"]
          executor = {
            type = "command"
            config = {
              command = "kubectl wait --for=condition=ready pod -l version=green --timeout=300s"
            }
          }
        },
        {
          name        = "switch-traffic"
          description = "Switch traffic to green environment"
          depends     = ["health-check"]
          executor = {
            type = "command"
            config = {
              command = "kubectl patch service myapp -p '{\"spec\":{\"selector\":{\"version\":\"green\"}}}'"
            }
          }
        }
      ]
    }
  ])
  
  runner = "kubiya-hosted"
}

resource "kubiya_agent" "deployment_agent" {
  name         = "deployment-automation"
  runner       = "kubiya-hosted"
  description  = "Automated deployment agent"
  instructions = <<-EOT
    You are a deployment automation agent. You can:
    1. Execute blue-green deployments
    2. Perform rollbacks if needed
    3. Monitor deployment status
    
    Always verify the target environment before deploying.
  EOT
  
  sources = [kubiya_source.deployment_workflow.id]
  
  integrations = ["kubernetes", "slack_integration"]
}
```

**Expected Outcome**: Creates an agent capable of executing complex deployment workflows.

### 7. Agent with Conversation Starters

Configure an agent with predefined conversation starters for better user experience:

```hcl
resource "kubiya_agent" "support_agent" {
  name         = "support-assistant"
  runner       = "kubiya-hosted"
  description  = "Customer support and troubleshooting assistant"
  instructions = "You are a support assistant helping users troubleshoot issues and answer questions."
  
  starters = [
    {
      name    = "Check System Status"
      command = "Show me the current system status and any ongoing incidents"
    },
    {
      name    = "Recent Errors"
      command = "Display the most recent error logs from the application"
    },
    {
      name    = "Performance Metrics"
      command = "Show current performance metrics and resource utilization"
    }
  ]
  
  integrations = ["datadog", "pagerduty"]
}
```

**Expected Outcome**: Creates an agent with quick-start commands for common support tasks.

### 8. Agent with Custom Docker Image

Use a custom Docker image for specialized agent functionality:

```hcl
resource "kubiya_agent" "custom_image_agent" {
  name         = "python-data-agent"
  runner       = "kubiya-hosted"
  description  = "Data analysis agent with Python libraries"
  instructions = "You are a data analysis agent with access to pandas, numpy, and matplotlib."
  
  # Custom Docker image with pre-installed libraries
  image = "ghcr.io/company/kubiya-python-agent:latest"
  
  model = "gpt-4o"
  
  environment_variables = {
    PYTHON_VERSION = "3.11"
    DATA_PATH      = "/data"
  }
  
  is_debug_mode = true  # Enable for troubleshooting
}
```

**Expected Outcome**: Creates an agent running in a custom Docker container with specialized tools.

## Argument Reference

### Required Arguments

* `name` - (Required, String) The name of the agent. Must be unique within your organization.
* `runner` - (Required, String) The runner to use for agent execution. Use "kubiya-hosted" for cloud execution or specify your own runner name.
* `description` - (Required, String) A detailed description of the agent's purpose and capabilities.
* `instructions` - (Required, String) System instructions that define the agent's behavior and capabilities.

### Optional Arguments

* `model` - (Optional, String) The LLM model to use. Defaults to "gpt-4o". Options include:
  - `gpt-4o` - GPT-4 Optimized (default)
  - `gpt-4` - GPT-4
  - `gpt-3.5-turbo` - GPT-3.5 Turbo
  - `azure/gpt-4` - Azure OpenAI GPT-4

* `image` - (Optional, String) Docker image for the agent. Defaults to "ghcr.io/kubiyabot/kubiya-agent:stable".

* `is_debug_mode` - (Optional, Boolean) Enable debug mode for detailed logging. Defaults to `false`.

* `integrations` - (Optional, List of Strings) List of integration names the agent can access.

* `users` - (Optional, List of Strings) List of user emails who can access the agent.

* `groups` - (Optional, List of Strings) List of group names that can access the agent.

* `sources` - (Optional, List of Strings) List of source IDs for knowledge bases and workflows.

* `tool_sources` - (Optional, List of Strings) List of tool source URLs or IDs.

* `secrets` - (Optional, List of Strings) List of secret names the agent can access.

* `environment_variables` - (Optional, Map of Strings) Map of environment variables available to the agent.

* `tasks` - (Optional, List of Objects) List of predefined tasks. Each task requires:
  - `name` - (String) Task identifier
  - `prompt` - (String) The prompt to execute
  - `description` - (String) Task description

* `starters` - (Optional, List of Objects) List of conversation starters. Each starter requires:
  - `name` - (String) Display name
  - `command` - (String) Command to execute

* `links` - (Optional, List of Strings) List of reference links for the agent.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the agent.
* `owner` - The email of the user who created the agent.
* `created_at` - The timestamp when the agent was created.

## Import

Agents can be imported using their ID:

```shell
terraform import kubiya_agent.example <agent-id>
```

## Compatibility Notes

* Requires Kubiya Terraform Provider version >= 1.0.0
* Compatible with Terraform >= 1.0
* Some features may require specific Kubiya platform tier (Enterprise features)
* Custom Docker images must be accessible from the runner environment
* Integration names must match exactly with configured integrations in your Kubiya account

## Best Practices

1. **Security**: Store sensitive information in secrets, not in environment variables
2. **Access Control**: Use groups for team-based access rather than individual users
3. **Instructions**: Be specific and clear in agent instructions to ensure consistent behavior
4. **Testing**: Test agents in a non-production runner before deploying to production
5. **Version Control**: Store your Terraform configurations in version control
6. **Naming**: Use descriptive names that indicate the agent's purpose
7. **Documentation**: Include comments in your Terraform code explaining complex configurations