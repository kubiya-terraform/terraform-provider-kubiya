# Quick Setup Guide for GitHub Secrets with Sentry

## 🚀 Quick Start (5 minutes)

### Step 1: Get Your Sentry DSNs

1. Login to [Sentry](https://sentry.io)
2. Go to your project → Settings → Client Keys (DSN)
3. Copy your DSN (format: `https://KEY@sentry.io/PROJECT_ID`)

### Step 2: Add GitHub Secrets

1. Go to your repo → **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**
3. Add these secrets:

| Secret Name             | Value                                      | Usage               |
|-------------------------|--------------------------------------------|---------------------|
| `SENTRY_DSN_PRODUCTION` | `https://prod-key@sentry.io/prod-id`       | Production releases |
| `SENTRY_DSN_STAGING`    | `https://staging-key@sentry.io/staging-id` | Development builds  |

### Step 3: Test the Setup

1. Go to **Actions** tab
2. Click **Build with Sentry** workflow
3. Click **Run workflow** → Choose "staging" → **Run workflow**
4. Wait for completion ✅

### Step 4: Create a Release

```bash
git tag v1.0.0
git push origin v1.0.0
```

Done! 🎉 Your provider will now build with Sentry integration.

## 🔍 Verification

### Check Build Artifacts

- Go to **Actions** → Recent workflow → **Artifacts**
- Download and verify binary was built

### Check Sentry Dashboard

- Login to Sentry
- Look for initialization events
- Verify environment is correct

## 🚨 Common Issues

| Issue                | Solution                                     |
|----------------------|----------------------------------------------|
| "Secret not found"   | Check secret name spelling                   |
| "Invalid DSN"        | Verify DSN format `https://KEY@sentry.io/ID` |
| Build without Sentry | DSN is empty - check secret value            |

## 📚 Full Documentation

- [Complete GitHub Secrets Setup](GITHUB_SECRETS_SETUP.md)
- [Build Instructions](BUILD_INSTRUCTIONS.md)
- [Sentry Integration Guide](../../docs/SENTRY_INTEGRATION.md)

## ⚡ Workflow Files

The setup automatically uses these workflow files:

- `.github/workflows/release.yml` - Production releases
- `.github/workflows/build-with-sentry.yml` - Development builds

No changes needed! They're already configured to use your secrets.

---
*Need help? Check the [troubleshooting section](BUILD_INSTRUCTIONS.md#troubleshooting) or contact the platform team.*