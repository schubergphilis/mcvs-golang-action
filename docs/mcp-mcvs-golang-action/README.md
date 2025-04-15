# MCP MCVS Golang Action

Model Context Protocol (MCP).

## Installation

```zsh
go install
```

## Configuration

Modify the settings.json of Visual Studio Code:

```json
"mcp": {
  "servers": {
    "golang-actions": {
      "command": "mcvs-golang-action",
      "args": []
    }
  }
},
"chat.agent.enabled": true,
```

## Usage

Open the GitHub Copilot chat and ensure that the `agent` mode is enabled.

Type: `lint project` and check the output of the Copilot. An answer could be:

```zsh
The linting process completed successfully, and no issues were found in the
project. However, there is a warning indicating that the mcvs-golang-action
version in the Taskfile.yml is outdated. It is currently set to v1.0.1, but it
should be updated to v1.1.2. You can resolve this by updating the
REMOTE_URL_REF variable in the Taskfile.yml file to v1.1.2. Let me know if
you'd like assistance with this update!
```

Then type: `fix` and the agent should fix the issue.
