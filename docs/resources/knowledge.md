---
page_title: "kubiya_knowledge Resource - terraform-provider-kubiya"
description: |-
  Provides a Kubiya Knowledge resource to manage specific information for agents.
---

# kubiya_knowledge (Resource)

Provides a Kubiya Knowledge resource. This allows knowledge resources to be created, updated, and deleted on the Kubiya platform. Knowledge resources add specific information that can be used by Kubiya agents.

## Example Usage

```hcl
# Simple knowledge resource
resource "kubiya_knowledge" "deployment_procedures" {
  name        = "deployment-procedures"
  description = "Standard procedures for application deployment"
  content     = file("${path.module}/deployment-procedures.md")
  format      = "markdown"
}

# Knowledge resource with metadata
resource "kubiya_knowledge" "architecture_diagram" {
  name        = "architecture-diagram"
  description = "System architecture diagram and explanation"
  content     = file("${path.module}/architecture.md")
  format      = "markdown"
  
  metadata = {
    department = "Engineering"
    version    = "2.1"
    author     = "DevOps Team"
  }
}

# JSON format knowledge resource
resource "kubiya_knowledge" "api_schemas" {
  name        = "api-schemas"
  description = "API schemas for our services"
  content     = file("${path.module}/api-schemas.json")
  format      = "json"
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the knowledge resource.
* `description` - (Required) A description of the knowledge resource.
* `content` - (Required) The content of the knowledge resource. This can be provided directly or loaded from a file using the `file()` function.

### Optional Arguments

* `format` - (Optional) The format of the content. Values: "markdown", "text", "json", etc. Defaults to "text".
* `metadata` - (Optional) Map of metadata key-value pairs. This can be used to store additional information about the knowledge resource.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the knowledge resource.
* `created_at` - The timestamp when the knowledge resource was created.
* `updated_at` - The timestamp when the knowledge resource was last updated.

## Import

Knowledge resources can be imported using the `id`:

```
$ terraform import kubiya_knowledge.example KNOWLEDGE_ID
``` 