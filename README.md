# DocDrift 📄⚡

**PR-level documentation drift detection. Catch stale docs before they ship.**

DocDrift scans your `git diff`, extracts changed function signatures, CLI flags, and API symbols, then checks if your documentation files still reference outdated versions. Runs as a single static binary — no runtime, no dependencies.

## 🚀 Quick Start

```bash
# Install
go install github.com/docdrift/docdrift@latest

# Run against last commit
docdrift --base HEAD~1

# Run against main branch in CI
docdrift --base origin/main --threshold 80 --format json

# Scan specific docs directory
docdrift --base HEAD~1 --docs ./docs
```

### GitHub Action

```yaml
- uses: actions/checkout@v4
  with:
    fetch-depth: 0
- run: go install github.com/docdrift/docdrift@latest
- run: docdrift --base origin/${{ github.base_ref }} --threshold 80
```

## 📊 How It Works

1. Parses `git diff` to find changed code files
2. Extracts function names, CLI flags (`--options`), class names from diff hunks
3. Scans all `.md`, `.rst`, `.adoc` files for references to those symbols
4. If a doc references a changed symbol but wasn't itself updated → **drift detected**
5. Calculates freshness score. Fails CI if below threshold.

## 💰 Pricing

| Feature | Free (OSS) | Pro ($49/mo) | Enterprise ($399/mo) |
|---|---|---|---|
| Symbol extraction (Python/Go/JS/TS/Java) | ✅ | ✅ | ✅ |
| Doc drift detection & scoring | ✅ | ✅ | ✅ |
| CI gate (threshold blocking) | ✅ | ✅ | ✅ |
| JSON/text output | ✅ | ✅ | ✅ |
| **Semantic similarity matching** | ❌ | ✅ | ✅ |
| **Auto-fix suggestions (AI)** | ❌ | ✅ | ✅ |
| **SARIF output for GitHub Security tab** | ❌ | ✅ | ✅ |
| **Slack/Teams notifications** | ❌ | ✅ | ✅ |
| **Multi-repo dashboard & trends** | ❌ | ❌ | ✅ |
| **SSO / SAML** | ❌ | ❌ | ✅ |
| **Audit trail & compliance reports** | ❌ | ❌ | ✅ |
| **Self-hosted option** | ❌ | ❌ | ✅ |

## 🤔 Why Pay?

- A single P1 support ticket from a customer hitting a stale API doc costs **$500+ in eng time**
- New hire following outdated CONTRIBUTING.md wastes **2+ days of onboarding**
- DocDrift Pro pays for itself after preventing **one** stale-doc incident per month
- ROI: $49/mo saves $6,000+/year in documentation-related incidents

## Supported Languages

Python (`def`), Go (`func`), JavaScript/TypeScript (`function`, `export`), Java/C# (`public/private/protected`), CLI flags (`--flag-name`)

## License

MIT — free core, paid features via license key.
