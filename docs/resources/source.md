---
page_title: "kubiya_source Resource - Kubiya"
subcategory: ""
description: |-
  The kubiya_source resource manages Git-based tool and workflow repositories in the Kubiya platform.
---

# kubiya_source (Resource)

The `kubiya_source` resource allows you to connect Git repositories containing tools and workflows to the Kubiya platform. These repositories provide agents with the functionality needed to perform tasks. For inline tool and workflow definitions, use the `kubiya_inline_source` resource instead.

## Prerequisites

Before using this resource, ensure you have:
1. A Kubiya account with API access
2. An API key (generated from Kubiya dashboard under Admin → Kubiya API Keys)
3. A Git repository containing tools and workflows (public or with appropriate credentials)
4. The repository URL accessible from your Kubiya environment

## Example Usage

### 1. Basic GitHub Source

Connect to a GitHub repository containing tools:

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

resource "kubiya_source" "github_tools" {
  name        = "community-tools"
  description = "Kubiya community tools repository"
  url         = "https://github.com/kubiyabot/community-tools"
  runner      = "kubiya-hosted"
}

resource "kubiya_agent" "basic_agent" {
  name         = "basic-agent"
  runner       = "kubiya-hosted"
  description  = "Agent with community tools"
  instructions = "You are an agent with access to community tools."
  
  sources = [kubiya_source.github_tools.id]
}
```

**Expected Outcome**: Creates a source linked to a GitHub repository with tools available to agents.

### 2. Private GitHub Repository

Connect to a private GitHub repository:

```hcl
resource "kubiya_source" "private_tools" {
  url    = "https://github.com/myorg/private-tools"
  runner = "kubiya-hosted"
  
  # Authentication handled via Kubiya GitHub integration
  dynamic_config = jsonencode({
    branch = "main"
    auth   = "github-integration"
  })
}

resource "kubiya_agent" "private_agent" {
  name         = "private-tools-agent"
  runner       = "kubiya-hosted"
  description  = "Agent with access to private tools"
  instructions = "You have access to proprietary tools from our private repository."
  
  sources = [kubiya_source.private_tools.id]
}
```

**Expected Outcome**: Creates a source linked to a private GitHub repository.

### 3. GitLab Repository

Connect to a GitLab repository:

```hcl
resource "kubiya_source" "gitlab_tools" {
  url    = "https://gitlab.com/myorg/devops-tools"
  runner = "kubiya-hosted"
  
  dynamic_config = jsonencode({
    branch = "develop"
    path   = "tools/"
  })
}

resource "kubiya_agent" "gitlab_agent" {
  name         = "gitlab-tools-agent"
  runner       = "kubiya-hosted"
  description  = "Agent with GitLab-hosted tools"
  instructions = "You have access to DevOps tools from our GitLab repository."
  
  sources = [kubiya_source.gitlab_tools.id]
}
```

**Expected Outcome**: Creates a source linked to a GitLab repository.

### 4. Bitbucket Repository

Connect to a Bitbucket repository:

```hcl
resource "kubiya_source" "bitbucket_tools" {
  url    = "https://bitbucket.org/myworkspace/automation-tools"
  runner = "kubiya-hosted"
  
  dynamic_config = jsonencode({
    branch = "master"
    tag    = "v1.2.0"
  })
}

resource "kubiya_agent" "bitbucket_agent" {
  name         = "bitbucket-tools-agent"
  runner       = "kubiya-hosted"
  description  = "Agent with Bitbucket-hosted tools"
  instructions = "You have access to automation tools from our Bitbucket repository."
  
  sources = [kubiya_source.bitbucket_tools.id]
}
```

**Expected Outcome**: Creates a source linked to a Bitbucket repository.

### 5. Multi-Environment Sources

Create environment-specific sources:

```hcl
locals {
  environments = {
    dev = {
      url    = "https://github.com/org/dev-tools"
      branch = "develop"
    }
    staging = {
      url    = "https://github.com/org/staging-tools"
      branch = "staging"
    }
    prod = {
      url    = "https://github.com/org/prod-tools"
      branch = "main"
    }
  }
}

resource "kubiya_source" "env_sources" {
  for_each = local.environments
  
  url    = each.value.url
  runner = "kubiya-hosted"
  
  dynamic_config = jsonencode({
    environment = each.key
    branch      = each.value.branch
    restricted  = each.key == "prod" ? true : false
  })
}

resource "kubiya_agent" "env_agents" {
  for_each = local.environments
  
  name         = "${each.key}-agent"
  runner       = "kubiya-hosted"
  description  = "Agent for ${each.key} environment"
  instructions = "You manage the ${each.key} environment with appropriate tools."
  
  sources = [kubiya_source.env_sources[each.key].id]
  
  environment_variables = {
    ENVIRONMENT = each.key
  }
}
```

**Expected Outcome**: Creates environment-specific sources with corresponding agents.

### 6. Monorepo with Multiple Tool Paths

Connect to different paths within a monorepo:

```hcl
resource "kubiya_source" "monorepo_backend" {
  url    = "https://github.com/org/monorepo"
  runner = "kubiya-hosted"
  
  dynamic_config = jsonencode({
    path   = "backend/tools/"
    branch = "main"
  })
}

resource "kubiya_source" "monorepo_frontend" {
  url    = "https://github.com/org/monorepo"
  runner = "kubiya-hosted"
  
  dynamic_config = jsonencode({
    path   = "frontend/tools/"
    branch = "main"
  })
}

resource "kubiya_agent" "fullstack_agent" {
  name         = "fullstack-agent"
  runner       = "kubiya-hosted"
  description  = "Full-stack development agent"
  instructions = "You have access to both backend and frontend tools from our monorepo."
  
  sources = [
    kubiya_source.monorepo_backend.id,
    kubiya_source.monorepo_frontend.id
  ]
}
```

**Expected Outcome**: Creates multiple sources from different paths within a monorepo.

## Argument Reference

### Required Arguments

* `url` - (Required, String) URL to a Git repository containing tools and workflows.

### Optional Arguments

* `runner` - (Optional, String) The runner to use for executing tools from this source. Defaults to "kubiya-hosted".
* `dynamic_config` - (Optional, String) JSON-encoded configuration for the source. Common options:
  - `branch` - Git branch to use
  - `tag` - Git tag to use
  - `path` - Subdirectory within the repository
  - `auth` - Authentication method
  - `environment` - Environment identifier

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the source.
* `name` - The computed name of the source (derived from repository).

## Import

Sources can be imported using their ID:

```shell
terraform import kubiya_source.example <source-id>
```

## Compatibility Notes

* Requires Kubiya Terraform Provider version >= 1.0.0
* Compatible with Terraform >= 1.0
* Git repositories must be accessible (public or with appropriate credentials)
* Sources must be created before agents can reference them
* Supports GitHub, GitLab, Bitbucket, and other Git providers

## Best Practices

1. **Version Control**: Use specific branches or tags for production deployments
2. **Repository Structure**: Organize tools logically within your repository
3. **Security**: Use private repositories for proprietary tools
4. **Authentication**: Configure appropriate Git credentials via Kubiya integrations
5. **Documentation**: Include README files in your tool repositories
6. **Testing**: Test repository connections in development first
7. **Branch Strategy**: Use different branches for different environments
8. **Monorepo Support**: Use path configuration for monorepo structures