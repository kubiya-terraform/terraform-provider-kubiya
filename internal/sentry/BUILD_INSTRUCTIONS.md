# Build Instructions with Sentry DSN

## Overview

The Kubiya Terraform Provider now supports injecting the Sentry DSN at build time, ensuring the DSN is never exposed in
source code.

## Prerequisites

- Go 1.22 or later
- GoReleaser (for release builds)
- Sentry DSN from your Sentry project

## Build Methods

### 1. Local Development Build (No Sentry)

For local development without Sentry integration:

```bash
go build
```

This creates a binary without Sentry integration (DSN will be empty).

### 2. Local Development Build (With Sentry)

To build with Sentry for local testing:

```bash
go build -ldflags "-X terraform-provider-kubiya/internal/sentry.DSN=https://YOUR_DSN@sentry.io/PROJECT_ID"
```

Replace `YOUR_DSN` and `PROJECT_ID` with your actual Sentry values.

### 3. Production Build with GoReleaser

#### Set Environment Variable

```bash
export SENTRY_DSN="https://YOUR_DSN@sentry.io/PROJECT_ID"
```

#### Run GoReleaser

```bash
# For a snapshot build (no tag required)
goreleaser build --snapshot --clean

# For a release build (requires git tag)
goreleaser release --clean
```

### 4. GitHub Actions Build

#### Setup GitHub Secrets

1. Go to your repository Settings → Secrets and variables → Actions
2. Add these secrets:
    - `SENTRY_DSN_PRODUCTION` - Your production Sentry DSN
    - `SENTRY_DSN_STAGING` - Your staging Sentry DSN

For detailed setup instructions, see [GITHUB_SECRETS_SETUP.md](GITHUB_SECRETS_SETUP.md).

#### Automatic Builds

The repository includes two workflows:

**Release Workflow** (`.github/workflows/release.yml`):

- Triggers on git tags (v*)
- Uses production Sentry DSN
- Creates signed releases

**Build with Sentry Workflow** (`.github/workflows/build-with-sentry.yml`):

- Triggers on pushes, PRs, and manual dispatch
- Supports both staging and production DSNs
- Creates test artifacts

#### Manual Build Trigger

1. Go to **Actions** tab in GitHub
2. Select **Build with Sentry** workflow
3. Click **Run workflow**
4. Choose environment (staging/production)
5. Click **Run workflow**

#### Creating a Release

```bash
# Tag and push for release
git tag v1.2.3
git push origin v1.2.3
```

This automatically triggers the release workflow with production Sentry DSN.

## Build Verification

### Check if Sentry DSN is Embedded

After building, you can verify if the DSN was properly embedded:

```bash
# Check if DSN string exists in binary (won't show actual DSN)
strings terraform-provider-kubiya | grep -i sentry

# Run with debug logging to see if Sentry initializes
TF_LOG=DEBUG terraform apply
```

### Test Sentry Integration

1. Build with your DSN
2. Set environment: `export KUBIYA_ENV=staging`
3. Run the provider with a test configuration
4. Check your Sentry dashboard for events

## Environment-Specific Builds

### Staging Build

```bash
export SENTRY_DSN="https://STAGING_DSN@sentry.io/PROJECT_ID"
export KUBIYA_ENV=staging
goreleaser build --snapshot --clean
```

### Production Build

```bash
export SENTRY_DSN="https://PRODUCTION_DSN@sentry.io/PROJECT_ID"
export KUBIYA_ENV=production
goreleaser release --clean
```

## Security Best Practices

### 1. Never Commit DSN

- Never hardcode DSN in source files
- Don't commit `.env` files with DSN
- Use `.gitignore` to exclude sensitive files

### 2. Use Secret Management

- **GitHub Actions**: Use GitHub Secrets
- **GitLab CI**: Use GitLab CI/CD Variables
- **Jenkins**: Use Credentials Plugin
- **Local**: Use environment variables or secret managers

### 3. Rotate DSNs Regularly

- Periodically rotate your Sentry DSNs
- Update CI/CD secrets when rotating
- Keep separate DSNs for staging/production

## Troubleshooting

### GitHub Actions Issues

#### Secret Not Found Error

```
Error: Required secret 'SENTRY_DSN_PRODUCTION' not found
```

**Solution**:

1. Verify secret exists in repository settings
2. Check secret name matches workflow exactly
3. Ensure you have repository admin access

#### Workflow Permission Issues

```
Error: Resource not accessible by integration
```

**Solution**:

1. Check repository permissions in Settings → Actions → General
2. Ensure "Read and write permissions" is selected
3. Verify workflow has necessary permissions in YAML

#### Build Fails with DSN

```
Error: Invalid DSN format
```

**Solution**:

1. Verify DSN format: `https://KEY@sentry.io/PROJECT_ID`
2. Check for extra spaces or quotes in secret
3. Test DSN manually with curl

### Local Build Issues

#### DSN Not Being Injected

1. Verify environment variable is set:
   ```bash
   echo $SENTRY_DSN
   ```

2. Check GoReleaser is using the correct ldflags:
   ```bash
   goreleaser build --snapshot --clean --debug
   ```

3. Verify the build command includes ldflags:
   ```bash
   go build -v -ldflags "-X terraform-provider-kubiya/internal/sentry.DSN=$SENTRY_DSN"
   ```

### Sentry Not Initializing

1. Check provider logs:
   ```bash
   TF_LOG=DEBUG terraform apply 2>&1 | grep -i sentry
   ```

2. Verify DSN format is correct:
    - Should be: `https://KEY@sentry.io/PROJECT_ID`
    - Or: `https://KEY@CUSTOM_DOMAIN/PROJECT_ID`

3. Check network connectivity to Sentry

### Build Fails with ldflags

1. Ensure package path is correct:
   ```bash
   go list -f '{{.ImportPath}}' ./internal/sentry
   ```

2. Verify variable name matches:
    - Check `internal/sentry/constants.go` has `var DSN string`

## Manual Build Script

Create a `build.sh` script for consistent builds:

```bash
#!/bin/bash

# build.sh - Build script with Sentry DSN injection

# Check if SENTRY_DSN is set
if [ -z "$SENTRY_DSN" ]; then
    echo "Warning: SENTRY_DSN not set, building without Sentry integration"
fi

# Get version from git
VERSION=$(git describe --tags --always --dirty)
COMMIT=$(git rev-parse HEAD)

# Build with ldflags
go build -v \
    -ldflags "-s -w \
        -X main.version=${VERSION} \
        -X main.commit=${COMMIT} \
        -X terraform-provider-kubiya/internal/sentry.DSN=${SENTRY_DSN}" \
    -o terraform-provider-kubiya

echo "Build complete: terraform-provider-kubiya"
echo "Version: ${VERSION}"
echo "Commit: ${COMMIT}"
echo "Sentry: $([ -z "$SENTRY_DSN" ] && echo "Disabled" || echo "Enabled")"
```

Make it executable:

```bash
chmod +x build.sh
```

Use it:

```bash
export SENTRY_DSN="your-dsn-here"
./build.sh
```

## Testing the Build

### Unit Test with DSN

```bash
export SENTRY_DSN="test-dsn"
go test -ldflags "-X terraform-provider-kubiya/internal/sentry.DSN=${SENTRY_DSN}" ./...
```

### Integration Test

```bash
# Build with DSN
export SENTRY_DSN="your-dsn"
go build -ldflags "-X terraform-provider-kubiya/internal/sentry.DSN=${SENTRY_DSN}"

# Test with Terraform
export KUBIYA_API_KEY="your-api-key"
export KUBIYA_ENV="staging"
terraform init
terraform plan
```

## Release Checklist

- [ ] Set SENTRY_DSN environment variable
- [ ] Verify DSN is correct for target environment
- [ ] Run tests with DSN injected
- [ ] Create git tag for release
- [ ] Run GoReleaser
- [ ] Verify release artifacts
- [ ] Check Sentry dashboard for initialization events
- [ ] Document DSN used in release notes (just mention environment, not actual DSN)

## Support

For issues with the build process:

1. Check this documentation
2. Review `.goreleaser.yml` configuration
3. Check GoReleaser documentation
4. Contact the platform team for Sentry DSN access