# gitBranchTool

[![Version](https://img.shields.io/github/v/release/cyrus2281/gitBranchTool?label=Version)](https://github.com/cyrus2281/gitBranchTool/releases/latest)
[![License](https://img.shields.io/github/license/cyrus2281/gitBranchTool)](https://github.com/cyrus2281/gitBranchTool/blob/main/LICENSE)
[![buyMeACoffee](https://img.shields.io/badge/BuyMeACoffee-cyrus2281-yellow?logo=buymeacoffee)](https://www.buymeacoffee.com/cyrus2281)
[![GitHub issues](https://img.shields.io/github/issues/cyrus2281/gitBranchTool?color=red)](https://github.com/cyrus2281/gitBranchTool/issues)
[![GitHub stars](https://img.shields.io/github/stars/cyrus2281/gitBranchTool?style=social)](https://github.com/cyrus2281/gitBranchTool/stargazers)

- [gitBranchTool](#gitbranchtool)
  - [Overview](#overview)
  - [Installation](#installation)
    - [Linux](#linux)
    - [MacOS](#macos)
    - [Windows (PowerShell)](#windows-powershell)
  - [Auto-Completion](#auto-completion)
    - [Linux Bash](#linux-bash)
    - [MacOS Bash](#macos-bash)
    - [Linux ZSH](#linux-zsh)
    - [MacOS ZSH](#macos-zsh)
    - [Windows PowerShell](#windows-powershell-1)
  - [Custom Prompt](#custom-prompt)
    - [Linux/MacOS](#linuxmacos)
    - [Windows PowerShell](#windows-powershell-2)
  - [Commands](#commands)
    - [Examples](#examples)
  - [Contributing](#contributing)
    - [Contributors](#contributors)
  - [License](#license)

## Overview

> **A bash tool to facilitate managing git branches with long cryptic names with aliases**

The `gitBranchTool`, `g`, command provides additional functionalities around working with _git_ branches.

If you frequently work with long branch names that include developer names, project names, issue numbers, and etc this tool is for you. With `gitBranchTool` or `g`, you'll be able to assign **alias names** for **each branch**.

You can monitor, switch, and delete your branches **using the aliases** instead of the long and confusing branch names (`g` commands **support both branch names and aliases**, so you wouldn't have to use _git_ for any of your branch switching needs).

Additionally, You can add notes to each branch to fully remember what they were about so you can keep the aliases shorter. You can list branches with their aliases and notes at anytime.

`g` also provides **auto-completion** for branch names and aliases, so you wouldn't even have to type the full alias name.

On top of all these, `g` provides a **custom prompt** that displays the name of the current repository, sub-directory, branch name, and its alias.

## Installation

### Linux

1. Download the artifact `g-linux-vX.X.X` from the latest release [here](https://github.com/cyrus2281/gitBranchTool/releases)

2. Add the binary to your PATH environment variable (or to a directory that is already in your PATH)
   > A directory that is already in your PATH is `/usr/local/bin/` or `/usr/bin/`

```bash
sudo mv g-linux-vX.X.X /usr/local/bin/g
```

3. Ensure the binary has the correct permissions

```bash
sudo chmod 755 /usr/local/bin/g
```

### MacOS

1. Download the artifact `g-macos-vX.X.X` from the latest release [here](https://github.com/cyrus2281/gitBranchTool/releases)

2. Add the binary to your PATH environment variable (or to a directory that is already in your PATH)
   > A directory that is already in your PATH is `/usr/local/bin/` or `/usr/bin/`

```bash
sudo mv g-macos-vX.X.X /usr/local/bin/g
```

3. Ensure the binary has the correct permissions

```bash
sudo chmod 755 /usr/local/bin/g
```

### Windows (PowerShell)

1. Download the artifact `g-win-vX.X.X` from the latest release [here](https://github.com/cyrus2281/gitBranchTool/releases)

2. Add the binary to your PATH environment variable (or to a directory that is already in your PATH)
   > A directory that is already in your PATH is `C:\Windows\System32\` (you may need to run PowerShell as an administrator)

```powershell
Move-Item -Force -Path .\g-win-vX.X.X.exe -Destination C:\Windows\System32\g.exe
```

## Auto-Completion

For auto-completion, you can run the following commands based on your shell:

### Linux Bash

```bash
sudo mkdir -p /etc/bash_completion.d && sudo touch /etc/bash_completion.d/g && USER=$(whoami); sudo chown $USER /etc/bash_completion.d/g && sudo chmod 755 /etc/bash_completion.d/g && sudo g completion bash > /etc/bash_completion.d/g"
```

### MacOS Bash

```bash
echo "\nautoload -U compinit; compinit" >> ~/.bashrc
sudo mkdir -p /etc/bash_completion.d && sudo touch /etc/bash_completion.d/g && USER=$(whoami); sudo chown $USER /etc/bash_completion.d/g && sudo chmod 755 /etc/bash_completion.d/_g && sudo g completion bash > /etc/bash_completion.d/g"
```

### Linux ZSH

```bash
sudo mkdir -p ${fpath[1]} && sudo touch ${fpath[1]}/_g && USER=$(whoami); sudo chown $USER ${fpath[1]}/_g && sudo chmod 755 ${fpath[1]}/_g && sudo g completion zsh > "${fpath[1]}/_g"
```

### MacOS ZSH

```bash
echo "\nautoload -U compinit; compinit" >> ~/.zshrc
sudo mkdir -p ${fpath[1]} && sudo touch ${fpath[1]}/_g && USER=$(whoami); sudo chown $USER ${fpath[1]}/_g && sudo chmod 755 ${fpath[1]}/_g && sudo g completion zsh > "${fpath[1]}/_g"
```

### Windows PowerShell

For default PowerShell profile file, run the following command:

```powershell
echo "g completion powershell | Out-String | Invoke-Expression" >> $PROFILE
```

## Custom Prompt

### Linux/MacOS

Run the following command to add the gitBranchTool custom prompt to your shell profile file:

- Change `.bashrc` with `.zshrc` if you use ZSH, or the profile file you use if it's different.

```bash
echo -e "\nPROMPT_COMMAND='export PS1=\"\$(g _ps)\"'\nprecmd() { eval \"\$PROMPT_COMMAND\"; }" >> ~/.bashrc
```

### Windows PowerShell

Run the following command to add the gitBranchTool custom prompt to your PowerShell profile file:

```powershell
echo "function prompt { g _ps }" >> $PROFILE
```

## Commands

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
  get              Get configuration options (Run `g get --help` for more information)
  getBranchAlias   Gets the branch alias
  help             Help about any command
  list             Lists all branches with their name, alias, and notes
  removeEntry      Removes a registered branch entry without deleting the branch
  rename           Updates the alias for the given branch name
  resolveAlias     Resolves the branch name from an alias
  set              Set configuration options (Run `g set --help` for more information)
  switch           Switches to the branch with the given name or alias
  updateBranchNote Adds/updates the notes for a branch base on name/alias
  updateCheck      Checks if a newer version is available

Flags:
  -h, --help      help for g
  -N, --no-log    no logs
  -V, --verbose   verbose output
  -v, --version   version for g

Use "g [command] --help" for more information about a command.
```

### Examples

- **Create Branch**

```bash
g c cyrus/jira-60083 banner "Adds banner to the home page"
```

- **Add alias to existing branch**

```bash
g a cyrus/jira-60083 banner "Adds banner to the home page"
```

- **List all branches**

```bash
g l
```

- **Switch to branch**

```bash
g s banner
```

- **Delete branch**

```bash
g d banner
```

- **Delete multiple branches**

```bash
g d banner cyrus/jira-50930
```

- **Set branch name prefix for current repository**

```bash
g set local-prefix dev/
```

- **Upgrade the tool to latest version**

```bash
g uc -y
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

- [internal/version.go line 3](./internal/version.go#L3)
- [CHANGE_LOGS.md line 3](./CHANGE_LOGS.md#L3) [You need to add your own, don't delete the existing one]

### Contributors

- [Cyrus Mobini (@cyrus2281)](https://github.com/cyrus2281)

## License

This project is licensed under the
[MIT License](./LICENSE)

Copyright 2024 - Cyrus Mobini (@cyrus2281)
