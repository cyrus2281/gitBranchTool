# gitBranchTool

[![GitHub Release](https://img.shields.io/github/v/release/cyrus2281/gitBranchTool?label=Version)](https://github.com/cyrus2281/gitBranchTool/releases/latest)
[![License](https://img.shields.io/github/license/cyrus2281/gitBranchTool)](https://github.com/cyrus2281/gitBranchTool/blob/main/LICENSE)
[![buyMeACoffee](https://img.shields.io/badge/BuyMeACoffee-cyrus2281-yellow?logo=buymeacoffee)](https://www.buymeacoffee.com/cyrus2281)
[![GitHub issues](https://img.shields.io/github/issues/cyrus2281/gitBranchTool?color=red)](https://github.com/cyrus2281/gitBranchTool/issues)
[![GitHub stars](https://img.shields.io/github/stars/cyrus2281/gitBranchTool?style=social)](https://github.com/cyrus2281/gitBranchTool/stargazers)


> **A bash tool to facilitate managing git branches with long cryptic names with aliases**

The `gitBranchTool`, `g`, command provides additional functionalities around working with *git* branches. 

If you frequently work with long branch names that include developer names, project names, issue numbers, and etc this tool is for you. With `gitBranchTool` or `g`, you'll be able to assign **alias names** for **each branch**.

You can monitor, switch, and delete your branches **using the aliases** instead of the long and confusing branch names (`g` commands **support both branch names and aliases**, so you wouldn't have to use *git* for any of your branch switching needs).

Additionally, You can add notes to each branch to fully remember what they were about so you can keep the aliases shorter. You can list branches with their aliases and notes at anytime.

`g` also provides **auto-completion** for branch names and aliases, so you wouldn't even have to type the full alias name.

On top of all these, `g` provides a **custom prompt** that displays the name of the current repository, sub-directory, branch name, and its alias. (You can install this by downloading and loading the [`gCustomPrompt.sh`](./gCustomPrompt.sh) file in your `.bashrc` or `.bash_profile`)

## Installation

### Linux/Unix

1. Download the latest non-`.exe` binary from the latest release [here](https://github.com/cyrus2281/gitBranchTool/releases)

2. Add the binary to your PATH environment variable (or to a directory that is already in your PATH)
- A directory that is already in your PATH is `/usr/local/bin/` or `/usr/bin/`

#### Auto-Completion
**For Bash on Linux:**
```bash
g completion bash > /etc/bash_completion.d/g
```

**For Bash on MacOs:**
```bash
g completion bash > $(brew --prefix)/etc/bash_completion.d/g
```

**For Zsh on Linux:**
```bash
g completion zsh > "${fpath[1]}/_g"
```

**For Zsh on MacOs:**
```bash
g completion zsh > $(brew --prefix)/share/zsh/site-functions/_g
```

### Windows (PowerShell)

1. Download the latest `.exe` binary from the latest release [here](https://github.com/cyrus2281/gitBranchTool/releases)

2. Add the binary to your PATH environment variable (or to a directory that is already in your PATH)

#### Auto-Completion
To add auto-completion to your PowerShell, you can add the following to your PowerShell profile file (`$PROFILE`):

```powershell
g completion powershell | Out-String | Invoke-Expression
```


## Usage:
```md
A bash tool to facilitate managing git branches with long cryptic names with aliases

Usage:
  g [command]

Available Commands:
  addAlias         Adds alias and note to a branch that is not stored yet
  completion       Generate the autocompletion script for the specified shell
  create           Creates a branch with name, alias, and note, and checks into it
  currentBranch    Returns the name of active branch with alias and note
  delete           Deletes listed branches base on name or alias
  getHome          Get the gitBranchTool's home directory path
  help             Help about any command
  list             Lists all branches with their name, alias, and notes
  rename           Updates the alias for the given branch name
  resolveAlias     Resolves the branch name from an alias
  setDefaultBranch Change the default branch, default is main
  switch           Switches to the branch with the given name or alias
  updateBranchNote Adds/updates the notes for a branch base on name/alias
  updateCheck      Checks if a newer version is available

Flags:
  -h, --help      help for g
  -v, --version   version for g

Use "g [command] --help" for more information about a command.
```

## Contributing

This repository is open for contributions.
If you have any suggestions or issues, please open an issue or a pull request.

In your pull request, please include a description of the changes you made and why you made them, and update the [CHANGE_LOGS.md](./CHANGE_LOGS.md), [VERSION](./VERSION), and [README.md](./README.md) (to contributors section) files accordingly.

For versioning
- Patch version: Bug fixes, performance improvements, etc.
- Minor version: New features, new commands.
- Major version: Major breaking changes, change of interface, backward incompatible changes

Version need to be updated in following files (Use search and replace all):
- [cmd/root.go line 22](./cmd/root.go#L22)
- [CHANGE_LOGS.md line 3](./CHANGE_LOGS.md#L3) [You need to add your own, don't delete the existing one]

### Contributors
- [Cyrus Mobini (@cyrus2281)](https://github.com/cyrus2281)

## License

This project is licensed under the
[MIT License](./LICENSE)

Copyright 2024 - Cyrus Mobini (@cyrus2281)