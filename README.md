# gitBranchTool

![Version](https://img.shields.io/badge/version-v2.1.5-blue)
[![License](https://img.shields.io/github/license/cyrus2281/gitBranchTool)](https://github.com/cyrus2281/gitBranchTool/blob/main/LICENSE)
[![buyMeACoffee](https://img.shields.io/badge/BuyMeACoffee-cyrus2281-yellow?logo=buymeacoffee)](https://www.buymeacoffee.com/cyrus2281)
[![GitHub issues](https://img.shields.io/github/issues/cyrus2281/gitBranchTool?color=red)](https://github.com/cyrus2281/gitBranchTool/issues)
[![GitHub stars](https://img.shields.io/github/stars/cyrus2281/gitBranchTool?style=social)](https://github.com/cyrus2281/gitBranchTool/stargazers)

**Tested and supported on: Bash, ZSH, and Git Bash**

> A bash tool to facilitate managing git branches with long cryptic names with aliases

The `gitBranchTool.sh` bash script adds `g` command to your terminal. This command provides additional functionalities around working with *git* branches. 

If you frequently work with long branch names that include developer names, project names, issue numbers, and etc this tool is for you. With `gitBranchTool` or `g`, you'll be able to assign **alias names** for **each branch**.

You can monitor, switch, and delete your branches **using the aliases** instead of the long and confusing branch names (`g` commands **support both branch names and aliases**, so you wouldn't have to use *git* for any of your branch switching needs).

Additionally, You can add notes to each branch to fully remember what they were about so you can keep the aliases shorter. You can list branches with their aliases and notes at anytime.

`g` also provides **auto-completion** for branch names and aliases, so you wouldn't even have to type the full alias name.

On top of all these, `g` provides a **custom prompt** that displays the name of the current repository, sub-directory, branch name, and its alias. (You can turn this off if you want to use your own custom prompt).

## Installation

Run the installation script using the following script:

- To download using `curl`:

```bash  
bash -c "$(curl -fsSL https://raw.githubusercontent.com/cyrus2281/gitBranchTool/main/install.sh)"
```

- To download using `wget`:

```bash
bash -c "$(wget -O- https://raw.githubusercontent.com/cyrus2281/gitBranchTool/main/install.sh)"
```

To activate the G customized prompt, type `yes` (or press enter) when prompted. You can change this later by re-running the installation script and providing the same values except for G customized prompt, or by manually changing the `G_CUSTOMIZED_PROMPT` variable in your terminal profile file.

The script will be installed in the `~/.gitBranchTool` directory. You can change this by setting the environment variable `GIT_BRANCH_TOOL_DIR` to **the absolute path** you want to install the script in before running the installation script. (This path will also be used to store the repository config files).

```bash
export GIT_BRANCH_TOOL_DIR=~/.gitBranchTool; curl ...
```

- Note: If you are using a custom directory, you should remember to use the same directory again or clean up manually if you're re-running the install.sh script.

By default, the script will be added to your bash terminal profile (`~/.bashrc`), and ZSH terminal profile (`~/.zshrc`) if one MacOS.

You will be prompt if you want to load the script in any other terminal profiles. You need to provide the absolute or relative path to the terminal profile file.
Press enter with no value to break the loop.


## Usage:
```md
The following commands can be used with gitBranchTool.

   g <command> [...<args>]


*  create <id> <alias> [<note>]           Creates a branch with id, alias, and note, and checks into it
   c      <id> <alias> [<note>]                  Uses the git command "git checkout -b <id>"

*  check  <id|alias>                      Checks into a branch base on an id or an alias
   switch <id|alias>                             Uses the git command "git checkout <id>"
   s      <id|alias>

*  del [...<id|alias>]                    Deletes listed branches base on ID or alias (requires at least one ID/alias)
   d   [...<id|alias>]                           Uses the git command "git branch -D [...<id>] "

*  list                                   Lists all branches with their id, alias, and notes
   l

*  resolve-alias <alias>                  Resolves the branch name from an alias
   r             <alias>

*  add-alias <id> <alias> [<note>]        Adds alias and note to a branch that is not stored yet
   a         <id> <alias> [<note>]

*  update-branch-alias <id> <alias>       Updates the alias for the given branch id

*  update-branch-note <id|alias> <note>   Adds/updates the notes for a branch base on id/alias

*  current-branch                         Returns the name of active branch with alias and note

*  edit-repository-config                 Opens active repository config file in vim for manual editing

*  update-check                           Checks for new version of gitBranchTool and prompts for update

*  help                                   Shows this help menu
   h

You can set the following parameters in your terminal profile:
  * G_DEFAULT_BRANCH                        Default branch name, usually master or main
  * G_DIRECTORY                             Where the gitBranchTool.sh script is and where the branch info should be stored
  * G_CUSTOMIZED_PROMPT                     To whether customize the prompt or not
  * G_BRANCH_DELIMITER                      Delimiter for branch info (default '|')
                                            This character should not be in your branch or alias names
```

## Contributing

This repository is open for contributions.
If you have any suggestions or issues, please open an issue or a pull request.

In your pull request, please include a description of the changes you made and why you made them, and update the [CHANGE_LOGS.md](./CHANGE_LOGS.md), [VERSIONS.md](./VERSIONS.md), and [README.md](./README.md) (to contributors section) files accordingly.

For versioning
- Patch version: Bug fixes, performance improvements, etc.
- Minor version: New features, new commands.
- Major version: Major breaking changes, change of interface, backward incompatible changes

Version need to be updated in following files (Use search and replace all):
- [CHANGE_LOGS.md line 3](./CHANGE_LOGS.md#L3) [You need to add your own, don't delete the existing one]
- [VERSION line 1](./VERSION#L1)
- [README.md line 3](./README.md#L3)
- [gitBranchTool.sh line 9](./gitBranchTool.sh#L9)

### Contributors
- [Cyrus Mobini (@cyrus2281)](https://github.com/cyrus2281)

## License

This project is licensed under the
[MIT License](./LICENSE)

Copyright 2024 - Cyrus Mobini (@cyrus2281)