---
page_title: "kubiya_integration Resource - terraform-provider-kubiya"
description: |-
  Provides a Kubiya Integration resource to connect with external systems.
---

# kubiya_integration (Resource)

Provides a Kubiya Integration resource. This allows integrations to be created, updated, and deleted on the Kubiya platform, enabling connections to external systems like GitHub, AWS, Kubernetes, Jira, and more.

## Example Usage

```hcl
# AWS integration example
resource "kubiya_integration" "aws" {
  name        = "production-aws"
  type        = "aws"
  description = "Integration with our production AWS account"
  
  configuration = {
    region          = "us-west-2"
    assume_role_arn = "arn:aws:iam::123456789012:role/KubiyaAssumeRole"
  }
}

# GitHub integration example
resource "kubiya_integration" "github" {
  name        = "github-org"
  type        = "github"
  description = "Organization GitHub account"
  
  configuration = {
    org_name = "example-org"
  }
}

# Kubernetes integration example
resource "kubiya_integration" "kubernetes" {
  name        = "production-k8s"
  type        = "kubernetes"
  description = "Production Kubernetes cluster"
  
  configuration = {
    context = "production"
  }
}

# Jira integration example
resource "kubiya_integration" "jira" {
  name        = "jira-cloud"
  type        = "jira"
  description = "Jira Cloud instance"
  
  configuration = {
    url = "https://example.atlassian.net"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the integration.
* `type` - (Required) The type of integration. Values include "aws", "github", "kubernetes", "jira", "slack", etc.

### Optional Arguments

* `description` - (Optional) A description of the integration.
* `configuration` - (Optional) Map of configuration values specific to the integration type. The required configuration keys depend on the integration type:
  * For AWS: "region", "assume_role_arn"
  * For GitHub: "org_name"
  * For Kubernetes: "context"
  * For Jira: "url"
  * For Slack: Configuration generally handled through OAuth

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the integration.
* `status` - The status of the integration.
* `created_at` - The timestamp when the integration was created.

## Import

Integrations can be imported using the `id`:

```
$ terraform import kubiya_integration.example INTEGRATION_ID
``` 