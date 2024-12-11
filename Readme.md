## Kubiya Terraform provider

Manage your kubiya cloud with our most extensive kubiya provider.


## For local Development

### Pre-requisits:
```
# Configure GOBIN (replace with your own value)
export GOBIN=<go bin value in your local env> (e.g /Users/michael.bauer/go/bin)

# Configure Kubiya BE instance:
export KUBIYA_ENV=staging

# Set KUBIYA_API_KEY
export KUBIYA_API_KEY=<generate at https://app.kubiya.ai/api-keys>
```

### Run
1. Create a simple main.tf file like:
```
terraform {
  required_providers {
    kubiya = {
#       source = "kubiya-terraform/kubiya"
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_source" "item" {
  url = "https://github.com/finebee/terraform-golden-usecases"
  dynamic_config = {
    michael = "hello"
  }
}

output "output" {
  value = kubiya_source.item
}
```

2. Create a ~/.terraformrc filem pay attention to fill dev_overrides value:
```
provider_installation {

  dev_overrides {
    "hashicorp.com/edu/Kubiya" = <paste $GOBIN value here>
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

3. Compile local code
```
go install .
```

4. Run terraform
```
terraform plan
terraform apply
terraform destroy
```
