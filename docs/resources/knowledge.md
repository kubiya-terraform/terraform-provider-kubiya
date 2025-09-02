---
page_title: "kubiya_knowledge Resource - Kubiya"
subcategory: ""
description: |-
  The kubiya_knowledge resource manages knowledge bases for agents in the Kubiya platform.
---

# kubiya_knowledge (Resource)

The `kubiya_knowledge` resource allows you to create and manage knowledge resources in the Kubiya platform. Knowledge resources provide specific information, documentation, and context that agents can reference when performing tasks.

## Prerequisites

Before using this resource, ensure you have:
1. A Kubiya account with API access
2. An API key (generated from Kubiya dashboard under Admin → Kubiya API Keys)
3. At least one group configured in your Kubiya organization
4. Content files prepared (markdown, text, or other formats)

## Example Usage

### 1. Basic Knowledge Resource

Create a simple knowledge resource:

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

resource "kubiya_knowledge" "basic_docs" {
  name        = "basic-documentation"
  description = "Basic system documentation"
  content     = "This is the basic documentation for our system."
  format      = "markdown"
  
  # Groups are required
  groups = ["Engineering"]
}
```

**Expected Outcome**: Creates a basic knowledge resource accessible to the Engineering group.

### 2. Deployment Procedures

Store deployment procedures documentation:

```hcl
resource "kubiya_knowledge" "deployment_guide" {
  name        = "deployment-procedures"
  description = "Standard procedures for application deployment"
  content     = file("${path.module}/docs/deployment-guide.md")
  format      = "markdown"
  
  groups = ["DevOps", "Engineering"]
  
  labels = ["deployment", "procedures", "documentation"]
}

resource "kubiya_agent" "deployment_assistant" {
  name         = "deployment-assistant"
  runner       = "kubiya-hosted"
  description  = "Deployment assistance agent"
  instructions = <<-EOT
    You are a deployment assistant with access to deployment procedures.
    Reference the deployment-procedures knowledge base for standard procedures.
    Always follow the documented deployment steps.
  EOT
  
  groups = ["DevOps"]
}
```

**Expected Outcome**: Creates deployment documentation that agents can reference during deployments.

### 3. Architecture Documentation

Store system architecture documentation:

```hcl
resource "kubiya_knowledge" "architecture" {
  name        = "system-architecture"
  description = "System architecture diagrams and explanations"
  content     = <<-MARKDOWN
    # System Architecture
    
    ## Overview
    Our system follows a microservices architecture with the following components:
    
    ### Frontend Layer
    - React application
    - Nginx reverse proxy
    - CloudFront CDN
    
    ### API Layer
    - REST API Gateway
    - GraphQL endpoint
    - WebSocket server
    
    ### Backend Services
    - User Service
    - Order Service
    - Payment Service
    - Notification Service
    
    ### Data Layer
    - PostgreSQL (primary database)
    - Redis (caching)
    - Elasticsearch (search)
    - S3 (object storage)
    
    ## Communication Patterns
    - Synchronous: REST/GraphQL
    - Asynchronous: RabbitMQ/Kafka
    - Real-time: WebSockets
    
    ## Security
    - JWT authentication
    - API rate limiting
    - WAF protection
  MARKDOWN
  format = "markdown"
  
  groups = ["Engineering", "DevOps", "Product", "Security"]
  
  labels = ["architecture", "documentation", "system-design", "technical"]
}
```

**Expected Outcome**: Creates comprehensive architecture documentation accessible to multiple teams.

### 4. Troubleshooting Runbooks

Create troubleshooting runbooks:

```hcl
resource "kubiya_knowledge" "troubleshooting" {
  name        = "troubleshooting-runbooks"
  description = "Runbooks for common issues and their resolution"
  content     = <<-MARKDOWN
    # Troubleshooting Runbooks
    
    ## High Memory Usage
    
    ### Symptoms
    - Memory usage above 90%
    - Application slowness
    - OOM kills
    
    ### Diagnosis Steps
    1. Check memory usage: `free -m`
    2. Identify top consumers: `ps aux --sort=-%mem | head`
    3. Check for memory leaks: `jmap -heap <pid>`
    
    ### Resolution
    1. Restart affected services
    2. Scale horizontally if needed
    3. Investigate code for memory leaks
    
    ## Database Connection Issues
    
    ### Symptoms
    - Connection timeout errors
    - "Too many connections" errors
    
    ### Diagnosis Steps
    1. Check connection count: `SELECT count(*) FROM pg_stat_activity;`
    2. Check for long-running queries
    3. Verify network connectivity
    
    ### Resolution
    1. Kill idle connections
    2. Increase connection pool size
    3. Optimize queries
  MARKDOWN
  format = "markdown"
  
  groups = ["DevOps", "SRE", "Engineering"]
  
  labels = ["troubleshooting", "runbooks", "operations", "incident-response"]
}

resource "kubiya_agent" "incident_responder" {
  name         = "incident-responder"
  runner       = "kubiya-hosted"
  description  = "Incident response agent"
  instructions = <<-EOT
    You are an incident response agent. 
    Use the troubleshooting-runbooks knowledge base to diagnose and resolve issues.
    Follow the documented procedures exactly.
  EOT
  
  groups = ["DevOps", "SRE"]
}
```

**Expected Outcome**: Creates runbooks that agents can follow during incident response.

### 5. API Documentation

Store API documentation:

```hcl
resource "kubiya_knowledge" "api_docs" {
  name        = "api-documentation"
  description = "REST API documentation and examples"
  content     = <<-MARKDOWN
    # API Documentation
    
    ## Authentication
    All API requests require a Bearer token:
    ```
    Authorization: Bearer <token>
    ```
    
    ## Endpoints
    
    ### GET /api/v1/users
    Retrieve user list
    
    **Parameters:**
    - page (int): Page number
    - limit (int): Results per page
    
    **Response:**
    ```json
    {
      "users": [...],
      "total": 100,
      "page": 1
    }
    ```
    
    ### POST /api/v1/users
    Create new user
    
    **Body:**
    ```json
    {
      "email": "user@example.com",
      "name": "John Doe",
      "role": "user"
    }
    ```
    
    ### Error Codes
    - 400: Bad Request
    - 401: Unauthorized
    - 403: Forbidden
    - 404: Not Found
    - 500: Internal Server Error
  MARKDOWN
  format = "markdown"
  
  groups = ["Engineering", "QA", "Product"]
  
  labels = ["api", "documentation", "rest", "reference"]
}
```

**Expected Outcome**: Creates API documentation for reference by development teams.

### 6. Compliance and Security Policies

Store compliance documentation:

```hcl
resource "kubiya_knowledge" "compliance_policies" {
  name        = "compliance-security-policies"
  description = "Security and compliance policies and procedures"
  content     = file("${path.module}/docs/compliance-policies.md")
  format      = "markdown"
  
  groups = ["Security", "Compliance", "Engineering", "DevOps"]
  
  labels = ["compliance", "security", "policies", "governance"]
}

resource "kubiya_knowledge" "gdpr_procedures" {
  name        = "gdpr-procedures"
  description = "GDPR compliance procedures"
  content     = file("${path.module}/docs/gdpr-procedures.md")
  format      = "markdown"
  
  groups = ["Compliance", "Legal", "Engineering"]
  
  labels = ["gdpr", "compliance", "data-privacy", "legal"]
}
```

**Expected Outcome**: Creates compliance documentation accessible to relevant teams.

### 7. Environment-Specific Knowledge

Create environment-specific documentation:

```hcl
locals {
  environments = {
    dev = {
      name        = "dev-environment-guide"
      description = "Development environment setup and guidelines"
      groups      = ["Engineering", "QA"]
    }
    staging = {
      name        = "staging-environment-guide"
      description = "Staging environment procedures"
      groups      = ["Engineering", "QA", "Product"]
    }
    prod = {
      name        = "production-environment-guide"
      description = "Production environment operations"
      groups      = ["DevOps", "SRE", "Security"]
    }
  }
}

resource "kubiya_knowledge" "env_guides" {
  for_each = local.environments
  
  name        = each.value.name
  description = each.value.description
  content     = file("${path.module}/docs/${each.key}-guide.md")
  format      = "markdown"
  
  groups = each.value.groups
  
  labels = [each.key, "environment", "documentation", "guide"]
}

resource "kubiya_agent" "env_assistants" {
  for_each = local.environments
  
  name         = "${each.key}-assistant"
  runner       = "kubiya-hosted"
  description  = "Assistant for ${each.key} environment"
  instructions = "You are an assistant for the ${each.key} environment. Reference the ${each.value.name} knowledge base."
  
  groups = each.value.groups
}
```

**Expected Outcome**: Creates environment-specific documentation with corresponding agents.

## Argument Reference

### Required Arguments

* `name` - (Required, String) The name of the knowledge resource. Must be unique within your organization.
* `description` - (Required, String) A description of the knowledge resource's content and purpose.
* `content` - (Required, String) The actual content of the knowledge resource.
* `format` - (Required, String) The format of the content (e.g., "markdown", "text", "json").
* `groups` - (Required, List of Strings) List of groups that can access this knowledge resource. At least one group is required.

### Optional Arguments

* `labels` - (Optional, List of Strings) Labels for categorizing and searching knowledge resources.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the knowledge resource.
* `created_at` - The timestamp when the knowledge resource was created.
* `updated_at` - The timestamp when the knowledge resource was last updated.
* `created_by` - The user who created the knowledge resource.

## Import

Knowledge resources can be imported using their ID:

```shell
terraform import kubiya_knowledge.example <knowledge-id>
```

## Compatibility Notes

* Requires Kubiya Terraform Provider version >= 1.0.0
* Compatible with Terraform >= 1.0
* Groups must exist in your Kubiya organization before being referenced
* Maximum content size may be limited by platform tier
* Supported formats include markdown, text, json, yaml, and html

## Best Practices

1. **Version Control**: Store knowledge content in version-controlled files
2. **Format Consistency**: Use markdown for documentation for better readability
3. **Access Control**: Assign appropriate groups based on content sensitivity
4. **Labeling**: Use consistent labels for easy discovery and categorization
5. **Updates**: Keep documentation current with regular updates
6. **Structure**: Use clear headings and sections in documentation
7. **Examples**: Include practical examples in technical documentation
8. **Cross-References**: Link related knowledge resources using labels
9. **Review Process**: Implement review processes for critical documentation updates