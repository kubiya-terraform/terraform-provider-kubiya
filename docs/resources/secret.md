---
page_title: "kubiya_secret Resource - terraform-provider-kubiya"
description: |-
  Provides a Kubiya Secret resource to manage sensitive information.
---

# kubiya_secret (Resource)

Provides a Kubiya Secret resource. This allows secrets to be created, updated, and deleted on the Kubiya platform. Secrets store sensitive information that can be securely accessed by Kubiya agents.

## Example Usage

```hcl
resource "kubiya_secret" "aws_credentials" {
  name        = "aws-credentials"
  description = "AWS access credentials for the production account"
  data = {
    aws_access_key_id     = var.aws_access_key_id
    aws_secret_access_key = var.aws_secret_access_key
    aws_region            = "us-west-2"
  }
}

resource "kubiya_secret" "database_credentials" {
  name        = "database-credentials"
  description = "Production database credentials"
  data = {
    username = "admin"
    password = var.db_password
    host     = "db.example.com"
    port     = "5432"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the secret.
* `description` - (Required) A description of the secret.
* `data` - (Required) Map of key-value pairs containing the secret data. These values are stored securely and can be accessed by agents that have permissions to use this secret.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the secret.
* `created_at` - The timestamp when the secret was created.
* `updated_at` - The timestamp when the secret was last updated.

## Import

Secrets can be imported using the `id`:

```
$ terraform import kubiya_secret.example SECRET_ID
```

## Security Note

Secrets should never be stored in plain text in your Terraform files. Instead, use variables or environment variables to pass sensitive values to your Terraform configuration.

Example with environment variables:

```hcl
variable "aws_access_key_id" {
  description = "AWS Access Key ID"
  sensitive   = true
}

variable "aws_secret_access_key" {
  description = "AWS Secret Access Key"
  sensitive   = true
}

resource "kubiya_secret" "aws_credentials" {
  name        = "aws-credentials"
  description = "AWS access credentials"
  data = {
    aws_access_key_id     = var.aws_access_key_id
    aws_secret_access_key = var.aws_secret_access_key
  }
}
```

Then run Terraform with:

```bash
export TF_VAR_aws_access_key_id="AKIAIOSFODNN7EXAMPLE"
export TF_VAR_aws_secret_access_key="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
terraform apply
``` 