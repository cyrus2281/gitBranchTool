# Version Change Logs

### Version 3.1.2
> Added remote branch deletion support for the delete command

### Version 3.1.1
> Fixed the bug where the tool config was not being automatically created

### Version 3.1.0
> Added `g get` command, (moved getHome to get) - Tool ready for public use milestone

### Version 3.0.8
> Added PowerShell Custom Prompt support + fixed completion repetition issue 

### Version 3.0.7
> Added set command with local and global branch prefix & moved setDefaultBranch to here

### Version 3.0.6
> Add prompt string command

### Version 3.0.5
> Added a check to ensure branch exists before adding an alias to it

### Version 3.0.4
> Fixed updated binary issue in PowerShell Windows

### Version 3.0.3
> Added removeEntry command + Added safe check around ensuring alias uniqueness

### Version 3.0.2
> Refactored code base to use the new go-logger package

### Version 3.0.1
> Added upgrade command to update the tool to latest version, `g uc -y`

### Version 3.0.0
> Rewrote codebase in GoLang instead of Bash

### Version 2.1.8
> Updated the help documentation

### Version 2.1.7
> Checking if the branch exists for 'g add-alias' + allowing deleting branches that are not registered

### Version 2.1.6
> Fixed updating env variables issue on Mac for install script

### Version 2.1.5
> Fixed issue of g tool failing on paths with spaces in them

### Version 2.1.4
> Updated install script so it can update the value of G_CUSTOMIZED_PROMPT

### Version 2.1.3
> Changed remote urls from master to release branch for updates

### Version 2.1.2
> Fixed issues of not being able to read the prompt outside bash

### Version 2.1.1
> Updating terminal session after an update install

### Version 2.1.0
> Added update check (and install) command

### Version 2.0.4
> Added auto complete on branch names for add-alias command

### Version 2.0.3
> Fixed issues of not using G_DIRECTORY path for config data

### Version 2.0.2
> Added install script + env variable rename

### Version 2.0.1
> Added Git bash support + simplification of the code

### Version 2.0.0
> Added G tool - Functions wrapper + autocompletion + new commands

### Version 1.1.2
> Added directory subpath to custom prompt

### Version 1.1.1
> Added bash support (in addition to zsh)

### Version 1.0.1
> Added gitBranchTool basic version

### Version 1.0.0
> Repository created
