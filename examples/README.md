# Testing Terraform Configuration Locally

This guide explains how to test the `main.tf` Terraform configuration located in the `examples` directory using your locally built Terraform provider.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19
- Locally built Terraform provider

Additional Configurations
```
export GOBIN=<go bin value in your local env> (e.g /Users/michael.bauer/go/bin)
export KUBIYA_ENV=staging
```

## Setup

### 1. Build the Provider

Ensure your provider is built and available locally:

```bash
cd /path/to/terraform-provider-example
go build -o terraform-provider-example
```


### 2. Configure Terraform to Use Local Provider

Edit or create the `~/.terraformrc` file to point to your local provider:
```
provider_installation {
  dev_overrides {
    "local/provider/example" = "/Users/avi.rosenberg/projects/terraform-provider-kubiya"
  }
  direct {}
}
```

replace the path with the relvant path to your provider.


### 3. Set Environment Variables

Set any necessary environment variables required by your provider:

```bash
export KUBIYA_API_KEY="your-api-key"
```

```bash
export TF_LOG=DEBUG # Optional: for detailed logging
export TF_LOG_PATH=./terraform.log # Optional: for log file output
```


## Testing the Configuration

### 1. Navigate to the Examples Directory

```bash
cd examples
```

### 2. Run Terraform Init

```bash
terraform apply
```

### 5. Verify the Results

After applying, verify that the resources have been created or modified as expected. You can check the Terraform output or directly inspect the resources in your environment.

### 6. Clean Up

To remove the resources created by this configuration, run:

```bash
terraform destroy
```
