# External Knowledge Resource Testing Guide

This directory contains comprehensive tests for the `kubiya_external_knowledge` Terraform resource.

## Test Files

1. **main.tf** - Basic example configuration
2. **test_comprehensive.tf** - Comprehensive test cases covering various scenarios
3. **test_operations.sh** - Shell script to test all CRUD operations
4. **test_all_operations.tf** - Terraform configuration for testing all operations

## Running the Tests

### Prerequisites

1. Ensure the provider is built:
   ```bash
   cd ../..
   go build -o ~/go/bin/terraform-provider-kubiya
   ```

2. Set your API key (or use a test key):
   ```bash
   export KUBIYA_API_KEY="your-api-key"
   # Or for testing without a real API:
   export KUBIYA_API_KEY="test-api-key"
   ```

### Test 1: Basic Configuration Test

Run the basic example:
```bash
terraform init
terraform plan
terraform apply
```

### Test 2: Comprehensive Test Cases

Test various configuration scenarios:
```bash
# Remove main.tf to avoid conflicts
mv main.tf main.tf.bak

# Run comprehensive tests
terraform init
terraform plan -var-file=test_comprehensive.tf
terraform apply -var-file=test_comprehensive.tf
```

### Test 3: CRUD Operations Test

Run the shell script to test all CRUD operations:
```bash
chmod +x test_operations.sh
./test_operations.sh
```

### Test 4: Manual Testing Steps

#### CREATE Operation
```bash
terraform apply -auto-approve
```

#### READ Operation (Refresh)
```bash
terraform refresh
```

#### UPDATE Operation
1. Modify the `channel_ids` in your configuration
2. Run:
   ```bash
   terraform plan
   terraform apply
   ```

#### DELETE Operation
```bash
terraform destroy -auto-approve
```

#### IMPORT Operation (if you have an existing resource)
```bash
terraform import kubiya_external_knowledge.example <resource-id>
```

## Test Scenarios Covered

### In test_comprehensive.tf:

1. **Single Channel** - Basic configuration with one channel
2. **Multiple Channels** - Configuration with multiple channels
3. **Dynamic Variables** - Using Terraform variables for channel lists
4. **Local Values** - Using Terraform locals and functions
5. **Conditional Creation** - Using count for conditional resources
6. **For Each** - Creating multiple resources with for_each
7. **Type Validation** - Outputs that validate data types

### Validation Points

The tests validate:
- Resource creation with valid configuration
- Proper handling of single and multiple channels
- Computed fields are populated correctly
- Type consistency (strings, lists, objects)
- Dynamic configuration support
- Terraform language features compatibility

## Expected Behavior

### With Valid API Key
- All operations should complete successfully
- Resources should be created in the Kubiya platform
- Computed fields should be populated with actual values

### With Test API Key
- Operations will fail at the API level with 401 errors
- This is expected and demonstrates the resource is working correctly
- The Terraform resource logic and type handling can still be validated

## Troubleshooting

1. **Provider not found**: Ensure the provider is built and in the correct location
2. **Type errors**: Check that you're using the latest build of the provider
3. **API errors**: Verify your API key is valid and has the necessary permissions

## Adding New Tests

To add new test cases:

1. Add new resource blocks to `test_comprehensive.tf`
2. Add corresponding outputs to validate the behavior
3. Document the test case in this README

## Clean Up

After testing:
```bash
terraform destroy -auto-approve
rm -f terraform.tfstate* .terraform.lock.hcl
rm -rf .terraform/
``` 