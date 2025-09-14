# Plan: Injecting Sentry DSN at Build Time

## Overview

Instead of hardcoding the DSN in the source code, we'll inject it during the build process using goreleaser's ldflags.
This provides better security and flexibility.

## Implementation Strategy

### 1. Update Constants Package

Modify `internal/sentry/constants.go` to accept DSN from build variables:

```go
package sentry

// DSN will be set at build time via ldflags
var DSN string = ""

// Fallback DSN for development (optional)
const DefaultDSN = ""
```

### 2. Update Sentry Initialization

Modify `internal/sentry/sentry.go` to handle build-time DSN:

```go
func getConfig(providerVersion string) *Config {
dsn := DSN
if dsn == "" {
dsn = DefaultDSN // Use fallback for local development
}

config := &Config{
DSN: dsn,
// ... rest of config
}
// ...
}
```

### 3. Update .goreleaser.yml

Add DSN injection to ldflags:

#### Option A: Direct DSN in ldflags

```yaml
builds:
  - env:
      - CGO_ENABLED=0
    mod_timestamp: '{{ .CommitTimestamp }}'
    flags:
      - -trimpath
    ldflags:
      - '-s -w -X main.version={{.Version}} -X main.commit={{.Commit}}'
      - '-X terraform-provider-kubiya/internal/sentry.DSN={{.Env.SENTRY_DSN}}'
    # ... rest of config
```

#### Option B: Multiple build configurations

```yaml
builds:
  - id: "production"
    env:
      - CGO_ENABLED=0
    mod_timestamp: '{{ .CommitTimestamp }}'
    flags:
      - -trimpath
    ldflags:
      - '-s -w -X main.version={{.Version}} -X main.commit={{.Commit}}'
      - '-X terraform-provider-kubiya/internal/sentry.DSN={{.Env.SENTRY_DSN_PROD}}'
      - '-X terraform-provider-kubiya/internal/sentry.Environment=production'
    # ... rest of config

  - id: "staging"
    env:
      - CGO_ENABLED=0
    mod_timestamp: '{{ .CommitTimestamp }}'
    flags:
      - -trimpath
    ldflags:
      - '-s -w -X main.version={{.Version}} -X main.commit={{.Commit}}'
      - '-X terraform-provider-kubiya/internal/sentry.DSN={{.Env.SENTRY_DSN_STAGING}}'
      - '-X terraform-provider-kubiya/internal/sentry.Environment=staging'
    # ... rest of config
```

## Build Process

### Local Development Build

```bash
# Without Sentry (uses empty DSN)
go build

# With Sentry DSN
go build -ldflags "-X terraform-provider-kubiya/internal/sentry.DSN=https://your-dsn@sentry.io/project"
```

### CI/CD Build with Goreleaser

```bash
# Set environment variable
export SENTRY_DSN="https://your-dsn@sentry.io/project"

# Run goreleaser
goreleaser release --clean
```

### GitHub Actions Example

```yaml
- name: Run GoReleaser
  uses: goreleaser/goreleaser-action@v4
  with:
    version: latest
    args: release --clean
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    SENTRY_DSN: ${{ secrets.SENTRY_DSN }}
    GPG_FINGERPRINT: ${{ secrets.GPG_FINGERPRINT }}
```

## Security Benefits

1. **No DSN in Source Code**: DSN never appears in the repository
2. **Environment-Specific DSNs**: Different DSNs for staging/production
3. **CI/CD Secret Management**: DSN stored as secret in CI/CD
4. **Build-Time Validation**: Can validate DSN during build

## Additional Considerations

### 1. Multiple DSN Support

We could support different DSNs for different environments:

```go
var (
DSNProduction string
DSNStaging    string
)

func getConfig(providerVersion string) *Config {
environment := getEnvironment()
dsn := ""

switch environment {
case EnvironmentStaging:
dsn = DSNStaging
case EnvironmentProduction:
dsn = DSNProduction
}

// ... rest of config
}
```

### 2. Build Tags

Alternative approach using build tags:

```go
// +build production

package sentry

const DSN = "production-dsn"
```

```go
// +build staging

package sentry

const DSN = "staging-dsn"
```

### 3. Obfuscation (Optional)

For additional security, we could obfuscate the DSN:

```go
var DSNEncoded string // Base64 encoded DSN

func getConfig(providerVersion string) *Config {
dsn := decodeDSN(DSNEncoded)
// ...
}
```

## Recommended Approach

**Use Option A (Single DSN with environment variable)** because:

1. Simpler configuration
2. One build artifact for all environments
3. Environment determined by `KUBIYA_ENV` at runtime
4. DSN injected at build time via CI/CD secrets

## Implementation Steps

1. **Update `internal/sentry/constants.go`**
    - Change DSN from const to var
    - Remove hardcoded placeholder

2. **Update `internal/sentry/sentry.go`**
    - Add validation for empty DSN
    - Log warning if DSN is missing (but don't fail)

3. **Update `.goreleaser.yml`**
    - Add DSN ldflags
    - Document environment variable requirement

4. **Update CI/CD**
    - Add SENTRY_DSN secret
    - Update build scripts

5. **Update Documentation**
    - Build instructions with DSN
    - CI/CD setup guide
    - Local development guide

## Testing Plan

1. **Build without DSN**: Verify graceful degradation
2. **Build with DSN**: Verify Sentry initialization
3. **CI/CD build**: Test with secrets
4. **Multiple environments**: Test environment switching

## Rollback Plan

If issues arise:

1. Keep current hardcoded approach as fallback
2. Use DefaultDSN constant for emergency builds
3. Document manual DSN update process