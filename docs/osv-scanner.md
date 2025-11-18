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
reason = "Waiting for upstream fix: https://github.com/anchore/syft/issues/4338"

[[IgnoredVulns]]
id = "GO-2024-1234"
reason = "False positive - not applicable to our usage"
```

### Important Notes

- Each ignored vulnerability should have a clear `reason` explaining why it's ignored
- Review and update the ignore list regularly
- Ignored vulnerabilities should be temporary - aim to fix or update dependencies

## Additional Resources

- [osv-scanner GitHub Repository](https://github.com/google/osv-scanner)
- [osv-scanner Documentation](https://google.github.io/osv-scanner/)
- [OSV Database](https://osv.dev/)
