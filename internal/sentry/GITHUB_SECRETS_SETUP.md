# GitHub Secrets Setup for Sentry DSN

## Overview

This guide explains how to set up GitHub Secrets to securely inject Sentry DSNs into your Terraform Provider builds.

## Required Secrets

### 1. Sentry DSN Secrets

You need to configure the following GitHub Secrets in your repository:

#### For Single Environment Setup:

- `SENTRY_DSN` - Your Sentry DSN for production releases

#### For Multi-Environment Setup (Recommended):

- `SENTRY_DSN_PRODUCTION` - Production Sentry DSN
- `SENTRY_DSN_STAGING` - Staging Sentry DSN

### 2. Existing Secrets (Already Required)

- `GPG_PRIVATE_KEY` - For signing releases
- `PASSPHRASE` - GPG key passphrase

## How to Add GitHub Secrets

### Step 1: Navigate to Repository Settings

1. Go to your GitHub repository
2. Click on **Settings** tab
3. In the left sidebar, click **Secrets and variables** → **Actions**

### Step 2: Add New Repository Secret

1. Click **New repository secret**
2. Add each secret with the name and value:

#### Secret: `SENTRY_DSN_PRODUCTION`

- **Name**: `SENTRY_DSN_PRODUCTION`
- **Secret**: `https://YOUR_PRODUCTION_KEY@sentry.io/YOUR_PRODUCTION_PROJECT_ID`

#### Secret: `SENTRY_DSN_STAGING`

- **Name**: `SENTRY_DSN_STAGING`
- **Secret**: `https://YOUR_STAGING_KEY@sentry.io/YOUR_STAGING_PROJECT_ID`

### Step 3: Verify Secrets

After adding, you should see the secrets listed in the Actions secrets section.

## Sentry DSN Format

Your Sentry DSN should look like:

```
https://PUBLIC_KEY@sentry.io/PROJECT_ID
```

Or for custom Sentry installations:

```
https://PUBLIC_KEY@your-sentry-domain.com/PROJECT_ID
```

### Finding Your DSN in Sentry:

1. Login to your Sentry account
2. Go to **Settings** → **Projects**
3. Select your project
4. Go to **Client Keys (DSN)**
5. Copy the DSN

## GitHub Actions Workflows

### Current Workflows Using Secrets:

#### 1. Release Workflow (`.github/workflows/release.yml`)

- Triggers on: Git tags (v*)
- Uses: `SENTRY_DSN` secret
- Purpose: Production releases

#### 2. Build with Sentry Workflow (`.github/workflows/build-with-sentry.yml`)

- Triggers on: Push to main/develop, PRs, manual dispatch
- Uses: `SENTRY_DSN_STAGING` or `SENTRY_DSN_PRODUCTION` based on input
- Purpose: Testing and development builds

## Workflow Usage Examples

### Automatic Release Build

```bash
# Create and push a tag
git tag v1.2.3
git push origin v1.2.3
```

This will trigger the release workflow with the production Sentry DSN.

### Manual Development Build

1. Go to **Actions** tab in GitHub
2. Select **Build with Sentry** workflow
3. Click **Run workflow**
4. Choose environment (staging/production)
5. Click **Run workflow**

### Pull Request Builds

Pull requests automatically build with staging Sentry DSN for testing.

## Security Best Practices

### 1. Separate Environments

- Use different Sentry projects for staging and production
- Never use production DSN in development builds
- Regularly rotate your Sentry DSNs

### 2. Access Control

- Limit repository access to trusted team members
- Use GitHub's branch protection rules
- Require reviews for changes to workflows

### 3. Secret Management

- Never hardcode DSNs in source code
- Don't log or echo secrets in workflows
- Use GitHub's secret scanning alerts

### 4. Audit and Monitoring

- Monitor secret usage in Actions logs
- Set up Sentry alerts for unexpected usage
- Regularly review who has access to secrets

## Troubleshooting

### Secret Not Found Error

```
Error: Required secret 'SENTRY_DSN' not found
```

**Solution**: Verify the secret name matches exactly in both the workflow and GitHub settings.

### Invalid DSN Format

```
Error: Invalid Sentry DSN format
```

**Solution**: Verify your DSN follows the correct format:

- Must start with `https://`
- Must contain `@sentry.io` or your custom domain
- Must end with a numeric project ID

### Build Without Sentry

If secrets are not available (e.g., in forks), the build will work without Sentry integration:

```
Warning: SENTRY_DSN not set, building without Sentry integration
```

### Debugging Workflow Issues

1. Check the Actions tab for workflow runs
2. Look for error messages in build logs
3. Verify secret names match in workflows
4. Test DSN format manually

## Environment Variables in Workflows

### Current Usage:

```yaml
env:
  SENTRY_DSN: ${{ secrets.SENTRY_DSN_PRODUCTION }}
```

### Conditional Usage:

```yaml
env:
  SENTRY_DSN: ${{ github.ref == 'refs/heads/main' && secrets.SENTRY_DSN_PRODUCTION || secrets.SENTRY_DSN_STAGING }}
```

## Testing Your Setup

### 1. Manual Test Build

```bash
# Clone repo
git clone https://github.com/your-org/terraform-provider-kubiya.git
cd terraform-provider-kubiya

# Trigger manual workflow
gh workflow run build-with-sentry.yml -f environment=staging
```

### 2. Check Build Artifacts

1. Go to Actions tab
2. Click on a completed workflow
3. Download the artifacts
4. Verify the binary was built correctly

### 3. Verify Sentry Integration

After a successful build:

1. Check your Sentry dashboard
2. Look for initialization events
3. Verify the correct environment is reporting

## Multiple Repository Setup

If you have multiple repositories that need Sentry integration:

### Option 1: Organization Secrets

1. Go to Organization Settings → Secrets and variables → Actions
2. Add organization-level secrets
3. Grant repository access

### Option 2: Repository Templates

Create a template repository with the secrets setup and workflows configured.

## Migration from Hardcoded DSN

If you're migrating from hardcoded DSNs:

1. **Add secrets** to GitHub repository
2. **Update workflows** to use secrets
3. **Remove hardcoded DSNs** from source code
4. **Test builds** with new setup
5. **Deploy new version** with secret-based DSN

## Support

### Common Issues:

1. **Secret not found**: Check secret name spelling
2. **Build fails**: Verify DSN format
3. **Sentry not working**: Check network connectivity and DSN validity

### Getting Help:

1. Check GitHub Actions documentation
2. Review Sentry documentation
3. Contact your platform team for Sentry access
4. Open an issue in this repository

## Checklist for Setup

- [ ] Sentry project created (staging and/or production)
- [ ] DSN obtained from Sentry dashboard
- [ ] GitHub secrets added to repository
- [ ] Secret names match workflow files
- [ ] Workflows tested with manual trigger
- [ ] Build artifacts verified
- [ ] Sentry dashboard shows initialization events
- [ ] Team has access to secrets management
- [ ] Documentation updated for your team

## Example Secret Values

```
# Production
SENTRY_DSN_PRODUCTION=https://abc123def456@o123456.ingest.sentry.io/789012

# Staging  
SENTRY_DSN_STAGING=https://xyz789uvw012@o123456.ingest.sentry.io/345678
```

**Note**: These are example values. Use your actual Sentry DSNs.