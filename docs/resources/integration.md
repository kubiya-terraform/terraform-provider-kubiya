---
page_title: "kubiya_integration Resource - Kubiya"
subcategory: ""
description: |-
  The kubiya_integration resource manages connections to external systems in the Kubiya platform.
---

# kubiya_integration (Resource)

The `kubiya_integration` resource allows you to create and manage integrations in the Kubiya platform. Integrations enable connections to external systems like AWS, GitHub, Kubernetes, Jira, and more, allowing agents to interact with these services.

## Prerequisites

Before using this resource, ensure you have:
1. A Kubiya account with API access
2. An API key (generated from Kubiya dashboard under Admin → Kubiya API Keys)
3. Appropriate credentials and permissions for the external systems you want to integrate

## Example Usage

### 1. Basic AWS Integration

Connect to a single AWS account:

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

resource "kubiya_integration" "aws_basic" {
  name            = "aws-production"
  description     = "Production AWS account"
  integration_type = "aws"
  
  configs = [
    {
      name       = "us-west-2"
      is_default = true
      vendor_specific = {
        arn    = "arn:aws:iam::123456789012:role/KubiyaRole"
        region = "us-west-2"
      }
    }
  ]
}
```

**Expected Outcome**: Creates an AWS integration with a single configuration for the us-west-2 region.

### 2. Multi-Region AWS Integration

Configure AWS integration with multiple regions:

```hcl
resource "kubiya_integration" "aws_multi_region" {
  name            = "aws-global"
  description     = "Global AWS integration with multiple regions"
  integration_type = "aws"
  auth_type       = "global"
  
  configs = [
    {
      name       = "us-west-2"
      is_default = true
      vendor_specific = {
        arn    = "arn:aws:iam::123456789012:role/KubiyaRole"
        region = "us-west-2"
      }
    },
    {
      name       = "eu-west-1"
      is_default = false
      vendor_specific = {
        arn    = "arn:aws:iam::123456789012:role/KubiyaRole"
        region = "eu-west-1"
      }
    },
    {
      name       = "ap-south-1"
      is_default = false
      vendor_specific = {
        arn    = "arn:aws:iam::123456789012:role/KubiyaRole"
        region = "ap-south-1"
      }
    }
  ]
}

resource "kubiya_agent" "aws_agent" {
  name         = "aws-multi-region-agent"
  runner       = "kubiya-hosted"
  description  = "AWS agent with multi-region access"
  instructions = "You are an AWS agent with access to multiple regions. Always specify the region when performing operations."
  
  integrations = [kubiya_integration.aws_multi_region.name]
}
```

**Expected Outcome**: Creates an AWS integration with multiple regional configurations and an agent to use it.

### 3. AWS Organization Integration

Set up integration for AWS Organizations:

```hcl
resource "kubiya_integration" "aws_org" {
  name             = "aws-organization"
  description      = "AWS Organization with multiple accounts"
  integration_type = "aws_organization"
  
  configs = [
    {
      name       = "master-account"
      is_default = true
      vendor_specific = {
        arn               = "arn:aws:iam::111111111111:role/OrganizationRole"
        region            = "us-east-1"
        organization_role = "OrganizationAccessRole"
      }
    },
    {
      name       = "dev-account"
      is_default = false
      vendor_specific = {
        arn     = "arn:aws:iam::222222222222:role/KubiyaRole"
        region  = "us-east-1"
        account = "222222222222"
      }
    },
    {
      name       = "prod-account"
      is_default = false
      vendor_specific = {
        arn     = "arn:aws:iam::333333333333:role/KubiyaRole"
        region  = "us-east-1"
        account = "333333333333"
      }
    }
  ]
}
```

**Expected Outcome**: Creates an AWS Organization integration with access to multiple AWS accounts.

### 4. GitHub Integration

Configure GitHub integration:

```hcl
resource "kubiya_integration" "github" {
  name             = "github-org"
  description      = "GitHub organization integration"
  integration_type = "github"
  
  configs = [
    {
      name       = "main-org"
      is_default = true
      vendor_specific = {
        org_name     = "my-organization"
        api_endpoint = "https://api.github.com"
      }
    }
  ]
}

resource "kubiya_agent" "github_agent" {
  name         = "github-bot"
  runner       = "kubiya-hosted"
  description  = "GitHub automation agent"
  instructions = "You are a GitHub agent that can manage repositories, pull requests, and issues."
  
  integrations = [kubiya_integration.github.name]
}
```

**Expected Outcome**: Creates a GitHub integration for repository management.

### 5. Kubernetes Integration

Set up Kubernetes cluster integration:

```hcl
resource "kubiya_integration" "kubernetes" {
  name             = "k8s-clusters"
  description      = "Kubernetes clusters integration"
  integration_type = "kubernetes"
  
  configs = [
    {
      name       = "production-cluster"
      is_default = true
      vendor_specific = {
        context   = "production"
        namespace = "default"
        server    = "https://k8s-prod.example.com"
      }
    },
    {
      name       = "staging-cluster"
      is_default = false
      vendor_specific = {
        context   = "staging"
        namespace = "default"
        server    = "https://k8s-staging.example.com"
      }
    }
  ]
}
```

**Expected Outcome**: Creates Kubernetes integration with multiple cluster configurations.

### 6. Jira Integration

Configure Jira integration:

```hcl
resource "kubiya_integration" "jira" {
  name             = "jira-cloud"
  description      = "Jira Cloud integration"
  integration_type = "jira"
  
  configs = [
    {
      name       = "main-instance"
      is_default = true
      vendor_specific = {
        url      = "https://mycompany.atlassian.net"
        project  = "OPS"
        api_type = "cloud"
      }
    }
  ]
}
```

**Expected Outcome**: Creates a Jira integration for issue tracking.

### 7. Multiple Integration Types

Create multiple integrations for a comprehensive DevOps setup:

```hcl
# AWS Integration
resource "kubiya_integration" "devops_aws" {
  name             = "devops-aws"
  description      = "DevOps AWS account"
  integration_type = "aws"
  
  configs = [
    {
      name       = "primary"
      is_default = true
      vendor_specific = {
        arn    = "arn:aws:iam::123456789012:role/DevOpsRole"
        region = "us-west-2"
      }
    }
  ]
}

# GCP Integration
resource "kubiya_integration" "devops_gcp" {
  name             = "devops-gcp"
  description      = "DevOps GCP project"
  integration_type = "gcp"
  
  configs = [
    {
      name       = "main-project"
      is_default = true
      vendor_specific = {
        project_id        = "my-gcp-project"
        service_account   = "kubiya@my-gcp-project.iam.gserviceaccount.com"
        region           = "us-central1"
      }
    }
  ]
}

# Azure Integration
resource "kubiya_integration" "devops_azure" {
  name             = "devops-azure"
  description      = "DevOps Azure subscription"
  integration_type = "azure"
  
  configs = [
    {
      name       = "main-subscription"
      is_default = true
      vendor_specific = {
        subscription_id = "12345678-1234-1234-1234-123456789012"
        tenant_id      = "87654321-4321-4321-4321-210987654321"
        client_id      = "abcdef12-3456-7890-abcd-ef1234567890"
        region         = "eastus"
      }
    }
  ]
}

# Create a multi-cloud agent
resource "kubiya_agent" "multi_cloud_agent" {
  name         = "multi-cloud-orchestrator"
  runner       = "kubiya-hosted"
  description  = "Multi-cloud infrastructure management agent"
  instructions = <<-EOT
    You are a multi-cloud infrastructure agent with access to:
    - AWS for compute and storage
    - GCP for data analytics and ML
    - Azure for enterprise applications
    
    Coordinate resources across all three cloud providers.
  EOT
  
  integrations = [
    kubiya_integration.devops_aws.name,
    kubiya_integration.devops_gcp.name,
    kubiya_integration.devops_azure.name
  ]
}
```

**Expected Outcome**: Creates multiple cloud provider integrations with a multi-cloud agent.

### 8. Environment-Specific Configurations

Set up environment-specific integration configurations:

```hcl
locals {
  environments = {
    dev = {
      arn    = "arn:aws:iam::111111111111:role/DevRole"
      region = "us-west-2"
    }
    staging = {
      arn    = "arn:aws:iam::222222222222:role/StagingRole"
      region = "us-east-1"
    }
    prod = {
      arn    = "arn:aws:iam::333333333333:role/ProdRole"
      region = "us-east-1"
    }
  }
}

resource "kubiya_integration" "env_aws" {
  name             = "aws-all-environments"
  description      = "AWS integration for all environments"
  integration_type = "aws"
  
  configs = [
    for env, config in local.environments : {
      name       = env
      is_default = env == "dev"
      vendor_specific = {
        arn    = config.arn
        region = config.region
        env    = env
      }
    }
  ]
}
```

**Expected Outcome**: Creates an AWS integration with configurations for multiple environments using dynamic blocks.

## Argument Reference

### Required Arguments

* `name` - (Required, String) The name of the integration. Must be unique within your organization.
* `configs` - (Required, List of Objects) List of configuration objects. Each config must have:
  - `name` - (Required, String) Name of the configuration
  - `is_default` - (Required, Boolean) Whether this is the default configuration
  - `vendor_specific` - (Required, Map) Vendor-specific configuration parameters

### Optional Arguments

* `description` - (Optional, String) A description of the integration's purpose.
* `auth_type` - (Optional, String) Authentication type. Defaults to empty string. Options: "global", "per_user".
* `integration_type` - (Optional, String) Type of integration. Defaults to "aws". Options include:
  - `aws` - Amazon Web Services
  - `aws_organization` - AWS Organizations
  - `gcp` - Google Cloud Platform
  - `azure` - Microsoft Azure
  - `github` - GitHub
  - `kubernetes` - Kubernetes
  - `jira` - Atlassian Jira
  - `confluence` - Atlassian Confluence

### Vendor-Specific Configuration

The `vendor_specific` map within each config varies by integration type:

#### AWS
* `arn` - IAM role ARN
* `region` - AWS region

#### AWS Organization
* `arn` - Organization role ARN
* `region` - AWS region
* `organization_role` - Cross-account role name
* `account` - Specific account ID (optional)

#### GCP
* `project_id` - GCP project ID
* `service_account` - Service account email
* `region` - GCP region

#### Azure
* `subscription_id` - Azure subscription ID
* `tenant_id` - Azure AD tenant ID
* `client_id` - Service principal client ID
* `region` - Azure region

#### GitHub
* `org_name` - GitHub organization name
* `api_endpoint` - API endpoint URL

#### Kubernetes
* `context` - Kubernetes context name
* `namespace` - Default namespace
* `server` - Kubernetes API server URL

#### Jira/Confluence
* `url` - Instance URL
* `project` - Default project (Jira)
* `api_type` - API type ("cloud" or "server")

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the integration.

## Import

Integrations can be imported using their ID:

```shell
terraform import kubiya_integration.example <integration-id>
```

## Compatibility Notes

* Requires Kubiya Terraform Provider version >= 1.0.0
* Compatible with Terraform >= 1.0
* At least one config must have `is_default = true`
* Integration names must match exactly when referenced by agents
* Some integration types may require additional setup in Kubiya dashboard

## Best Practices

1. **Least Privilege**: Configure integrations with minimal required permissions
2. **Environment Separation**: Use separate configs for different environments
3. **Default Configuration**: Always set one configuration as default
4. **Naming Convention**: Use clear names that indicate the integration's purpose
5. **Documentation**: Document vendor-specific parameters and their purpose
6. **Security**: Never hardcode sensitive credentials; use secure credential management
7. **Testing**: Test integrations in non-production environments first