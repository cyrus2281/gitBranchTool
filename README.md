# gitBranchTool

**Tested and supported on: Bash, ZSH, and Git Bash**

> A bash tool to facilitate managing git branch with long cryptic names with aliases

The `gitBranchTool.sh` bash script adds `g` command to your terminal. This command provides additional functionalities around working with *git* branches. 

If you frequently work with long branch names that include developer names, project names, issue numbers, and etc this tool is for you. With `gitBranchTool` or `g`, you'll be able to assign **alias names** for **each branch**.

You can monitor, switch, and delete your branches **using the aliases** instead of the long and confusing branch names (`g` commands **support both branch names and aliases**, so you wouldn't have to use *git* for any of your branch switching needs).

Additionally, You can add notes to each branch to fully remember what they were about so you can keep the aliases shorter. You can list branches with their aliases and notes at anytime.

`g` also provides **auto-completion** for branch names and aliases, so you wouldn't even have to type the full alias name.

On top of all these, `g` provides a **custom prompt** that displays the name of the current repository, sub-directory, branch name, and its alias. (You can turn this off if you want to use your own custom prompt).

## Installation

Download the script and run:

```bash
echo "source PATH_TO_SCRIPT/gitBranchTool.sh" >> ~/.bashrc
```
- Replace `PATH_TO_SCRIPT` with the path to the script file.
- Replace `~/.bashrc` with the path to your bash profile file.
    - `~/.bashrc` is the default bash profile file for most Linux distributions.
    - `~/.zshrc` is the default zsh profile file for most MacOS distributions.

### Install without custom prompt
Alternatively, you can install the script without the custom prompt by setting the environment variable `CUSTOMIZED_GIT_PROMPT` to `false` in your bash profile file before sourcing the script.

```bash
echo "export CUSTOMIZED_GIT_PROMPT=false" >> ~/.bashrc
echo "source PATH_TO_SCRIPT/gitBranchTool.sh" >> ~/.bashrc
```
- Replace `PATH_TO_SCRIPT` with the path to the script file.
- Replace `~/.bashrc` with the path to your bash profile file.



## Usage:
```bash
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

*  update-branch-note <id|alias> <note>           Adds/updates the notes for a branch base on id/alias

*  current-branch                         Returns the name of active branch with alias and note

*  edit-repository-config                         Opens active repository config file in vim for manual editing

*  help                                   Shows this help menu
   h

You can set the following parameters in your terminal profile:
  * DEFAULT_BRANCH                        Default branch name, usually master or main
  * CUSTOMIZED_GIT_PROMPT                 To whether customize the prompt or not
  * BRANCH_DELIMITER                      Delimiter for branch info (default '|')
                                            This character should not be in your branch or alias names
```

Copyright 2024 Cyrus Mobini (@cyrus2281)