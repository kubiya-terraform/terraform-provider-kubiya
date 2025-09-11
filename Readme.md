## Kubiya Terraform Provider

Manage your Kubiya cloud with our most enhanced Kubiya provider.

### Features

- Full support for Kubiya resources (agents, runners, webhooks, triggers, etc.)
- Built-in error tracking and monitoring with Sentry
- Distributed tracing for performance monitoring
- Automatic error recovery and reporting

### Error Monitoring

This provider includes automatic error reporting to help improve reliability. Errors and performance metrics are
collected to identify and fix issues proactively. All sensitive data is automatically redacted before transmission.

Environment configuration:

- Set `KUBIYA_ENV=staging` for staging environment
- Default environment is `production`

#### Quick Setup for Contributors

If you're building releases, set up Sentry DSN in GitHub Secrets:

- See [Quick Setup Guide](internal/sentry/GITHUB_QUICK_SETUP.md) for 5-minute setup
- Full instructions in [GitHub Secrets Setup](internal/sentry/GITHUB_SECRETS_SETUP.md)

For detailed information about the monitoring integration, see [SENTRY_INTEGRATION.md](docs/SENTRY_INTEGRATION.md).

## For Local Development

Check out https://github.com/kubiya-terraform/terraform-provider-kubiya/tree/main/examples