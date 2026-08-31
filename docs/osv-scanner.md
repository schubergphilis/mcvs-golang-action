# osv-scanner

## Overview

The `mcvs-golang-action` uses [osv-scanner](https://github.com/google/osv-scanner)
by Google to scan for vulnerabilities in Go modules. osv-scanner is actively
maintained and provides robust vulnerability scanning using the OSV
(Open Source Vulnerabilities) database.

## Ignoring Vulnerabilities

Add an `osv-scanner.toml` file to your project to ignore certain vulnerabilities
that cannot be fixed right away. This allows you to acknowledge known issues while
preventing the CI/CD pipeline from failing.

### Configuration Format

Create an `osv-scanner.toml` file in your project root:

```toml
# osv-scanner.toml
# Documentation: https://google.github.io/osv-scanner/configuration/

# Ignore specific vulnerabilities
[[IgnoredVulns]]
id = "GO-2025-4020"
ignoreUntil = 2025-12-20
reason = "Waiting for upstream fix: https://github.com/anchore/syft/issues/4338"

[[IgnoredVulns]]
id = "GO-2024-1234"
ignoreUntil = 2025-12-20
reason = "False positive - not applicable to our usage"
```

### Important Notes

- Each ignored vulnerability should have a clear `reason` explaining why it's ignored
- Review and update the ignore list regularly
- Ignored vulnerabilities should be temporary - aim to fix or update dependencies

## Call Analysis

By default osv-scanner runs govulncheck based Go call analysis, which type checks
the whole build graph. On larger modules this requires more than 8GB of memory and
the process is OOM killed by the runner (exit code 137, output ends with
`Killed`). As the reported vulnerabilities are identical with and without it,
call analysis is disabled in this action.

Set `OSV_SCANNER_CALL_ANALYSIS` to `true` to re-enable it, and make sure enough
memory is available:

```yaml
jobs:
  mcvs-golang-action:
    runs-on: ubuntu-24.04
    env:
      OSV_SCANNER_CALL_ANALYSIS: "true"
    steps:
      - uses: schubergphilis/mcvs-golang-action@vX.Y.Z
        with:
          testing-type: security-golang-modules
```

## Additional Resources

- [osv-scanner GitHub Repository](https://github.com/google/osv-scanner)
- [osv-scanner Documentation](https://google.github.io/osv-scanner/)
- [OSV Database](https://osv.dev/)
