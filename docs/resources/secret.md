---
page_title: "kubiya_secret Resource - Kubiya"
subcategory: ""
description: |-
  The kubiya_secret resource manages sensitive information securely in the Kubiya platform.
---

# kubiya_secret (Resource)

The `kubiya_secret` resource allows you to create and manage secrets in the Kubiya platform. Secrets store sensitive information like API keys, passwords, and tokens that can be securely accessed by Kubiya agents during execution.

## Prerequisites

Before using this resource, ensure you have:
1. A Kubiya account with API access
2. An API key (generated from Kubiya dashboard under Admin → Kubiya API Keys)
3. Proper security practices for handling sensitive data in Terraform

## Example Usage

### 1. Basic API Key

Store a simple API key:

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

resource "kubiya_secret" "api_key" {
  name        = "external-api-key"
  value       = var.external_api_key
  description = "API key for external service"
}

resource "kubiya_agent" "api_agent" {
  name         = "api-integration-agent"
  runner       = "kubiya-hosted"
  description  = "Agent that uses external API"
  instructions = "You are an agent that integrates with external APIs using secure credentials."
  
  secrets = [kubiya_secret.api_key.name]
}
```

**Expected Outcome**: Creates a secret that can be referenced by agents for API authentication.

### 2. Database Password

Store database connection password:

```hcl
variable "db_password" {
  type        = string
  sensitive   = true
  description = "Database password"
}

resource "kubiya_secret" "db_password" {
  name        = "production-db-password"
  value       = var.db_password
  description = "Production database password"
}

resource "kubiya_agent" "db_manager" {
  name         = "database-manager"
  runner       = "kubiya-hosted"
  description  = "Database management agent"
  instructions = "You are a database administrator. Use the production-db-password secret for database connections."
  
  secrets = [kubiya_secret.db_password.name]
  
  environment_variables = {
    DB_HOST     = "db.example.com"
    DB_PORT     = "5432"
    DB_NAME     = "production"
    DB_USER     = "admin"
    SECRET_NAME = kubiya_secret.db_password.name
  }
}
```

**Expected Outcome**: Creates a database password secret that an agent can use for database operations.

### 3. GitHub Token

Store GitHub personal access token:

```hcl
resource "kubiya_secret" "github_token" {
  name        = "github-pat"
  value       = var.github_token
  description = "GitHub Personal Access Token"
}

resource "kubiya_agent" "github_automation" {
  name         = "github-automation"
  runner       = "kubiya-hosted"
  description  = "GitHub automation agent"
  instructions = <<-EOT
    You are a GitHub automation agent. 
    Use the github-pat secret for API authentication.
    You can manage repositories, create pull requests, and handle issues.
  EOT
  
  secrets = [kubiya_secret.github_token.name]
  
  integrations = ["github"]
}
```

**Expected Outcome**: Creates a GitHub token for repository automation.

### 4. AWS Secret Access Key

Store AWS credentials:

```hcl
resource "kubiya_secret" "aws_secret_key" {
  name        = "aws-secret-access-key"
  value       = var.aws_secret_access_key
  description = "AWS Secret Access Key for automation"
}

resource "kubiya_agent" "aws_automation" {
  name         = "aws-automation"
  runner       = "kubiya-hosted"
  description  = "AWS automation agent"
  instructions = "You are an AWS automation agent. Use the AWS credentials securely."
  
  secrets = [kubiya_secret.aws_secret_key.name]
  
  environment_variables = {
    AWS_ACCESS_KEY_ID = var.aws_access_key_id  # Non-sensitive
    AWS_REGION        = "us-west-2"
  }
}
```

**Expected Outcome**: Creates AWS secret key for cloud automation tasks.

### 5. Multiple Secrets for an Agent

Configure an agent with multiple secrets:

```hcl
resource "kubiya_secret" "slack_token" {
  name        = "slack-bot-token"
  value       = var.slack_bot_token
  description = "Slack bot token for notifications"
}

resource "kubiya_secret" "jira_token" {
  name        = "jira-api-token"
  value       = var.jira_api_token
  description = "Jira API token"
}

resource "kubiya_secret" "datadog_api_key" {
  name        = "datadog-api-key"
  value       = var.datadog_api_key
  description = "Datadog API key for monitoring"
}

resource "kubiya_agent" "integration_manager" {
  name         = "integration-manager"
  runner       = "kubiya-hosted"
  description  = "Multi-service integration manager"
  instructions = <<-EOT
    You manage integrations with multiple services:
    - Use slack-bot-token for Slack notifications
    - Use jira-api-token for Jira operations
    - Use datadog-api-key for monitoring tasks
  EOT
  
  secrets = [
    kubiya_secret.slack_token.name,
    kubiya_secret.jira_token.name,
    kubiya_secret.datadog_api_key.name
  ]
  
  integrations = ["slack", "jira", "datadog"]
}
```

**Expected Outcome**: Creates multiple secrets for various service integrations.

### 6. Environment-Specific Secrets

Create secrets for different environments:

```hcl
locals {
  environments = ["dev", "staging", "prod"]
}

resource "kubiya_secret" "env_db_passwords" {
  for_each = toset(local.environments)
  
  name        = "${each.value}-db-password"
  value       = var.db_passwords[each.value]
  description = "Database password for ${each.value} environment"
}

resource "kubiya_agent" "env_agents" {
  for_each = toset(local.environments)
  
  name         = "${each.value}-agent"
  runner       = "kubiya-hosted"
  description  = "Agent for ${each.value} environment"
  instructions = "You are an agent managing the ${each.value} environment."
  
  secrets = [kubiya_secret.env_db_passwords[each.value].name]
  
  environment_variables = {
    ENVIRONMENT = each.value
    DB_HOST     = "${each.value}-db.example.com"
  }
}
```

**Expected Outcome**: Creates environment-specific secrets with corresponding agents.

### 7. JWT Secret

Store JWT signing secret:

```hcl
resource "kubiya_secret" "jwt_secret" {
  name        = "jwt-signing-secret"
  value       = var.jwt_secret
  description = "Secret key for JWT token signing"
}

resource "kubiya_agent" "auth_agent" {
  name         = "authentication-agent"
  runner       = "kubiya-hosted"
  description  = "Authentication and authorization agent"
  instructions = <<-EOT
    You manage authentication and JWT tokens.
    Use the jwt-signing-secret for token operations.
    Never expose the secret value in logs or outputs.
  EOT
  
  secrets = [kubiya_secret.jwt_secret.name]
}
```

**Expected Outcome**: Creates a JWT secret for authentication operations.

### 8. SSH Private Key

Store SSH private key (base64 encoded):

```hcl
resource "kubiya_secret" "ssh_key" {
  name        = "server-ssh-key"
  value       = base64encode(file("${path.module}/keys/id_rsa"))
  description = "SSH private key for server access (base64 encoded)"
}

resource "kubiya_agent" "server_admin" {
  name         = "server-administrator"
  runner       = "kubiya-hosted"
  description  = "Server administration agent"
  instructions = <<-EOT
    You are a server administrator with SSH access.
    The server-ssh-key secret contains the base64-encoded SSH private key.
    Decode it before use and handle it securely.
  EOT
  
  secrets = [kubiya_secret.ssh_key.name]
  
  environment_variables = {
    SSH_USER = "ubuntu"
    SSH_HOST = "server.example.com"
    SSH_PORT = "22"
  }
}
```

**Expected Outcome**: Creates an SSH key secret for secure server access.

## Argument Reference

### Required Arguments

* `name` - (Required, String) The name of the secret. Must be unique within your organization.
* `value` - (Required, String, Sensitive) The secret value. This is write-only and cannot be read back after creation.

### Optional Arguments

* `description` - (Optional, String) A description of the secret's purpose.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - The timestamp when the secret was created.
* `created_by` - The user who created the secret.

Note: The `value` attribute is write-only and will not be included in the state file or outputs for security reasons.

## Import

Secrets can be imported using their name:

```shell
terraform import kubiya_secret.example <secret-name>
```

Note: The actual secret value cannot be read back after import for security reasons. You will need to update the value in your configuration to match the existing secret.

## Compatibility Notes

* Requires Kubiya Terraform Provider version >= 1.0.0
* Compatible with Terraform >= 1.0
* Secret values are write-only and cannot be read back after creation
* Secrets must exist before agents can reference them
* Secret names must be unique within your organization

## Best Practices

1. **Never Hardcode**: Never hardcode sensitive values directly in Terraform files
2. **Use Variables**: Always use Terraform variables marked as `sensitive = true`
3. **State Security**: Ensure Terraform state is encrypted and stored securely
4. **Rotation**: Implement regular secret rotation policies
5. **Least Privilege**: Only grant secret access to agents that require it
6. **Naming Convention**: Use clear names that indicate the secret's purpose and environment
7. **Audit Trail**: Monitor secret access through Kubiya audit logs
8. **Environment Separation**: Keep production and non-production secrets separate
9. **Backup**: Maintain secure backups of critical secrets outside of Terraform
10. **Base64 Encoding**: Use base64 encoding for binary data like certificates or keys