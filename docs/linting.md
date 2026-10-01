# Automatically Fixing Linting Issues

When golangci-lint reports issues that can be automatically fixed, you can use the `fix-linting-issues` task:

```zsh
task remote:fix-linting-issues --yes
```

This task automatically fixes common linting issues using two tools:

- **golines**: Reformats Go code to meet line length requirements (wraps long lines)
- **wsl**: Fixes whitespace linting issues (adds/removes blank lines according to style rules)

After running this task, review the changes and commit them. Note that some linting issues may still require manual fixes.
