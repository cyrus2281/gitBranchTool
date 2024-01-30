#!/bin/bash

# A bash tool to facilitate managing git branch with long cryptic names with aliases

# Author: Cyrus Mobini - https://github.com/cyrus2281
# Github Repository: https://github.com/cyrus2281/gitBranchTool
# License: MIT License

G_VERSION="2.0.2"

__g_help(){
  echo -e "\nThe following commands can be used with gitBranchTool.\n"
  echo -e "   g <command> [...<args>]\n"
  echo -e "\n*  create <id> <alias> [<note>] \t  Creates a branch with id, alias, and note, and checks into it"
  echo -e "   c      <id> <alias> [<note>]    \t\t Uses the git command \"git checkout -b <id>\""
  echo -e "\n*  check  <id|alias> \t\t\t  Checks into a branch base on an id or an alias"
  echo -e "   switch <id|alias> \t\t\t\t Uses the git command \"git checkout <id>\""
  echo -e "   s      <id|alias>"
  echo -e "\n*  del [...<id|alias>] \t\t\t  Deletes listed branches base on ID or alias (requires at least one ID/alias)"
  echo -e "   d   [...<id|alias>] \t\t\t\t Uses the git command \"git branch -D [...<id>] \""
  echo -e "\n*  list \t\t\t\t  Lists all branches with their id, alias, and notes"
  echo -e "   l"
  echo -e "\n*  resolve-alias <alias> \t\t  Resolves the branch name from an alias"
  echo -e "   r             <alias>"
  echo -e "\n*  add-alias <id> <alias> [<note>] \t  Adds alias and note to a branch that is not stored yet"
  echo -e "   a         <id> <alias> [<note>]"
  echo -e "\n*  update-branch-alias <id> <alias> \t  Updates the alias for the given branch id"
  echo -e "\n*  update-branch-note <id|alias> <note> \t  Adds/updates the notes for a branch base on id/alias"
  echo -e "\n*  current-branch \t\t\t  Returns the name of active branch with alias and note"
  echo -e "\n*  edit-repository-config \t\t\t  Opens active repository config file in vim for manual editing"
  echo -e "\n*  help \t\t\t\t  Shows this help menu"
  echo -e "   h"
  echo -e ""
  echo -e "You can set the following parameters in your terminal profile:"
  echo -e "  * G_DEFAULT_BRANCH \t\t\t  Default branch name, usually master or main"
  echo -e "  * G_CUSTOMIZED_PROMPT \t\t  To whether customize the prompt or not"
  echo -e "  * G_BRANCH_DELIMITER \t\t\t  Delimiter for branch info (default '|')"
  echo -e "                      \t\t\t    This character should not be in your branch or alias names"
  echo -e "\nGitBranchTool Version: $G_VERSION"
  echo -e "Created by Cyrus Mobini"
}

__DEFAULT_G_BRANCH_PATH=~/.gitBranchTool/.g

G_BRANCH_DELIMITER=${G_BRANCH_DELIMITER:-'|'}
G_BRANCH_PATH=${G_BRANCH_PATH:-"$__DEFAULT_G_BRANCH_PATH"}
G_CUSTOMIZED_PROMPT=${G_CUSTOMIZED_PROMPT:-true}
G_DEFAULT_BRANCH=${G_DEFAULT_BRANCH:-"main"}

if [ -n "$ZSH_VERSION" ]; then
  # Current shell is ZSH
  autoload -U +X bashcompinit && bashcompinit
  autoload -U +X compinit && compinit
fi

# Git Branch Snippets
mkdir -p ~/.gitBranchTool/

__g_get_branch_name() {
  echo $(git branch 2> /dev/null | grep \* | cut -d "*" -f2 | cut -d " " -f2)
}

__g_current_G_BRANCH_PATH(){
  if [[ -z $(__g_get_branch_name) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  currentPath="$G_BRANCH_PATH.$(basename $(git rev-parse --show-toplevel))"
  if [[ ! -e $currentPath ]]; then
    touch $currentPath
  fi
  echo $currentPath
}

__g_list(){
  if [[ -z $(__g_get_branch_name) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  cat -n $(__g_current_G_BRANCH_PATH) | tr "$G_BRANCH_DELIMITER" '\t'
}

__g_current_branch() {
  if [[ -z $(__g_get_branch_name) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  cat $(__g_current_G_BRANCH_PATH) | tr "$G_BRANCH_DELIMITER" '\t' | grep $(__g_get_branch_name)
  if [ $?  != 0 ]; then
      echo "master (or unregistered branch)"
  fi
}

__g_resolve_alias(){
  if [[ -z $(__g_get_branch_name) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 1 )); then
    echo "Wrong Usage"
    echo -e "\tg resolve-alias <alias>"
    return 1
  fi
  while IFS="$G_BRANCH_DELIMITER" read -r id als desc; do
    if [[ $1 == $als ]]; then
        echo $id
        return 0
    fi
  done < $(__g_current_G_BRANCH_PATH)
  echo "-- Alias not found --"
  return 1
}

__g_edit_repo_config() {
  vim $(__g_current_G_BRANCH_PATH)
}

__g_add_alias(){
  if [[ -z $(__g_get_branch_name) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 2 )); then
    echo "Wrong Usage"
    echo -e "\tg add-alias <id> <alias> [<note>]"
    return 1
  fi
  # checking for alias to be unique
  while IFS="$G_BRANCH_DELIMITER" read -r id als desc; do
    if [[ $2 == $als ]]; then
        echo '-- Alias should be unique --'
        echo '-- FAILED --'
        return 1
    fi
  done < $(__g_current_G_BRANCH_PATH)
  # Adding branch, alias and note to list
  id=$1
  alias=$2
  shift 2
  note=$@
  # Adding branch, alias and note to list
  echo "$id$G_BRANCH_DELIMITER$alias$G_BRANCH_DELIMITER$note" >> "$(__g_current_G_BRANCH_PATH)"
  echo "-- Added alias '$alias' for branch '$id' --"
}

__g_update_branch_note() {
  if [[ -z $(__g_get_branch_name) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 2 )); then
    echo "Wrong Usage"
    echo -e "\tg update-branch-note <id|alias> <note>"
    return 1
  fi
  branchPath=$(__g_current_G_BRANCH_PATH)
  searchParam=$1
  shift
  newNote="$@"
  newContent=()
  message=""
  found=false
  while IFS="$G_BRANCH_DELIMITER" read -r id als note; do
    # update note
    if [[ $searchParam == $id ||  $searchParam == $als ]]; then
      newContent+=("$id$G_BRANCH_DELIMITER$als$G_BRANCH_DELIMITER$newNote")
      message="-- update note for '$id  $als' to '$newNote' --"
      found=true
    else
      newContent+=("$id$G_BRANCH_DELIMITER$als$G_BRANCH_DELIMITER$note")
    fi
  done < $(__g_current_G_BRANCH_PATH)
  # Updating file content
  echo -n "" > "$branchPath"
  for line in "${newContent[@]}"; do
    echo "$line" >> "$branchPath"
  done

  if [[ $found == false ]]; then
      echo "-- branch not found: $searchParam --"
  else
    echo $message
  fi
}

__g_update_branch_alias() {
  if [[ -z $(__g_get_branch_name) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 2 )); then
    echo "Wrong Usage"
    echo -e "\tg update-branch-alias <id> <note>"
    return 1
  fi
  newContent=()
  message=""
  found=false
  while IFS="$G_BRANCH_DELIMITER" read -r id als note; do
    # update note
    if [[ $1 == $id ]]; then
      newContent+=("$id$G_BRANCH_DELIMITER$2$G_BRANCH_DELIMITER$note")
      message="-- update alias for '$id' from  '$als' to '$2' --"
      found=true
    else
      newContent+=("$id$G_BRANCH_DELIMITER$als$G_BRANCH_DELIMITER$note")
    fi
  done < $(__g_current_G_BRANCH_PATH)
  # Updating file content
  echo -n "" > "$branchPath"
  for line in "${newContent[@]}"; do
    echo "$line" >> "$branchPath"
  done

  if [[ $found == false ]]; then
      echo "-- branch not found: $1 --"
  else
    echo $message
  fi
}

__g_create(){
  if [[ -z $(__g_get_branch_name) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 2 )); then
    echo "Wrong Usage"
    echo -e "\tg create <id> <alias> [<note>]"
    return 1
  fi
  # checking for alias to be unique
  while IFS="$G_BRANCH_DELIMITER" read -r id als desc; do
    if [[ $2 == $als ]]; then
        echo '-- Alias should be unique --'
        echo '-- FAILED --'
        return 1
    fi
  done < $(__g_current_G_BRANCH_PATH)

  # creating and checking out to branch
  git checkout -b $1
  if [ $? -eq 0 ]; then
    id=$1
    alias=$2
    shift 2
    note=$@
    # Adding branch, alias and note to list only if operation was successful
    echo "$id$G_BRANCH_DELIMITER$alias$G_BRANCH_DELIMITER$note" >> "$(__g_current_G_BRANCH_PATH)"
    return 0
  else
    echo '-- FAILED --'
    return 1
  fi
}

__g_switch(){
  if [[ -z $(__g_get_branch_name) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 1 )); then
    echo "Wrong Usage"
    echo -e "\tg switch <id|alias>"
    return 1
  fi
  if [[ $1 == $G_DEFAULT_BRANCH ]]; then
    git checkout $G_DEFAULT_BRANCH
    return 0
  fi
  # Searching for alias/id in the list
  while IFS="$G_BRANCH_DELIMITER" read -r id als desc; do
    if [[ $1 == $als || $1 == $id ]]; then
        git checkout $id
        return 0
    fi
  done < $(__g_current_G_BRANCH_PATH)
  # Branch not in list, trying to checkout
  git checkout $1
  # if successful
  if [[ $?  == 0 ]]; then
      # If provided alias
      if [[ -n $2 ]]; then
        while IFS="$G_BRANCH_DELIMITER" read -r id als desc; do
          if [[ $2 == $als ]]; then
              echo '-- Alias should be unique --'
              echo '-- Branch Switch successfull. Failed to add branch alias. --'
              return 0 
          fi
        done < $(__g_current_G_BRANCH_PATH)
        # Adding branch, alias and note to list
        id=$1
        alias=$2
        shift 2
        note=$@
        # Adding branch, alias and note to list
        echo "$id$G_BRANCH_DELIMITER$alias$G_BRANCH_DELIMITER$note" >> "$(__g_current_G_BRANCH_PATH)"
        echo "-- Branch $id has been registered with alias $alias"
        return 0
      else
        echo "-- branch \"$1\" is not registered with alias --"
        return 0
      fi
  fi
  echo "-- branch not found: $1 --"
  return 1
}

__g_del(){
  if [[ -z $(__g_get_branch_name) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 1 )); then
    echo "Wrong Usage"
    echo -e "\tg del <alias|id> [...<alias|id>]"
    return 1
  fi
  branchPath=$(__g_current_G_BRANCH_PATH)
  for value in "$@"; do
    newContent=()
    found=false
    while IFS="$G_BRANCH_DELIMITER" read -r id als note; do
      # Delete base on id/alias
      if [[ $value == $id || $value == $als ]]; then
        echo "-- deleting branch: $id $als --"
        git branch -D $id
        if [ $? -ne 0 ]; then
          echo "-- FAILED to delete $id --"
          newContent+=("$id$G_BRANCH_DELIMITER$als$G_BRANCH_DELIMITER$note")
        fi
        found=true
      else 
        newContent+=("$id$G_BRANCH_DELIMITER$als$G_BRANCH_DELIMITER$note")
      fi
    done < $(__g_current_G_BRANCH_PATH)
    # Updating file content
    echo -n "" > "$branchPath"
    for line in "${newContent[@]}"; do
      echo "$line" >> "$branchPath"
    done
    # branch not found
    if [[ $found == false ]]; then
        echo "-- branch not found: $value --"
    fi
  done
}

__g_get_ids(){
  if [[ -z $(__g_get_branch_name) ]]; then
    return 0
  fi
  IDs=()
  # Read each line of the file
  while IFS="$G_BRANCH_DELIMITER" read -r id name desc; do
    # Add values to the array
    IDs+=("$id")
  done < $(__g_current_G_BRANCH_PATH)
  echo ${IDs[@]}
}

__g_get_aliases(){
  if [[ -z $(__g_get_branch_name) ]]; then
    return 0
  fi
  names=()
  # Read each line of the file
  while IFS="$G_BRANCH_DELIMITER" read -r id name desc; do
    # Add values to the array
    names+=("$name")
  done < $(__g_current_G_BRANCH_PATH)
  echo ${names[@]} ${G_DEFAULT_BRANCH}
}

g() {
  if (( $# < 1 )); then
    echo "Missing arguments!"
    echo -e "\tEnter 'g help' to get a list of all command."
    return 1
  fi
  command=$1
  shift
  case $command in
    "create"|"c")
      __g_create $@
      ;;
    "check"|"switch"|"s")
      __g_switch $@
      ;;
    "del"|"d")
      __g_del $@
      ;;
    "resolve-alias"|"r")
      __g_resolve_alias $@
      ;;
    "add-alias"|"a")
    __g_add_alias $@
      ;;
    "update-branch-alias")
    __g_update_branch_alias $@
      ;;
    "update-branch-note")
    __g_update_branch_note $@
      ;;
    "list"|"l")
      __g_list
      ;;
    "current-branch")
    __g_current_branch
      ;;
    "edit-repository-config")
    __g_edit_repo_config
      ;;
    "help"|"h")
      __g_help
      ;;
    *)
      echo "g command not found!"
      echo -e "\tEnter 'g help' to get a list of all command."
      return 1
      ;;
  esac
}

__g_complete() {
  local cur_word
  local commands
  COMPREPLY=()
  cur_word="${COMP_WORDS[COMP_CWORD]}"
  prev_word="${COMP_WORDS[COMP_CWORD-1]}"
  command_word=${COMP_WORDS[1]}

  if [ $COMP_CWORD -eq 1 ]; then
      commands=("check" "list" "del" "help" "create" "switch" \
       "resolve-alias" "edit-repository-config" "add-alias" \
       "update-branch-note" "current-branch" "update-branch-alias" )
  else
    case "$command_word" in
      # Aliases & ID on all args
      "del"|"d")
        commands=($(__g_get_ids) $(__g_get_aliases))
        ;;
      # Aliases & ID first arg
      "check"|"switch"|"s"|"update-branch-note")
        if [ $COMP_CWORD -eq 2 ]; then
            commands=($(__g_get_ids) $(__g_get_aliases))
        fi
        ;;
      # Aliases only first arg
      "resolve-alias"|"r")
        if [ $COMP_CWORD -eq 2 ]; then
            commands=($(__g_get_aliases))
        fi
        ;;
      # IDs only first arg
      "update-branch-alias")
        if [ $COMP_CWORD -eq 2 ]; then
            commands=($(__g_get_ids))
        fi
        ;;
    esac
  fi

  # If the current word starts with the characters in one of the commands, suggest it
  COMPREPLY=($(compgen -W "${commands[*]}" -- "$cur_word"))
  return 0
}

# Register the autocompletion function for sayWhat command
complete -F __g_complete g

__g_get_subdirectory() {
  repo=$1
  # Find the index of the repository name in the current working directory
  index=$(echo "$PWD" | awk -v repo="$repo" '{print index($0, repo)}')
  # Use expr to calculate the start position for the subpath
  start=$((index + ${#repo}))
  # Extract the subpath after the repository name
  subpath="${PWD:$start}"
  # Remove leading slash if present
  subpath=${subpath#/}
  # Check if subpath is not empty and set it to [subpath]
  if [ -n "$subpath" ]; then
      subpath=" [$subpath]"
  fi
  echo $subpath
}

# Custom prompt
__g_get_name() {
  brn=""
  if [[ -n $(__g_get_branch_name) ]]; then
    alias=""
    while IFS="$G_BRANCH_DELIMITER" read -r id als desc; do
      if [[ $(__g_get_branch_name) == $id ]]; then
          alias=" ($als)"
          break
      fi
    done < $(__g_current_G_BRANCH_PATH)
    topPath="$(git rev-parse --show-toplevel)"
    repo=$(basename $topPath)
    subpath=$(__g_get_subdirectory $repo)
    brn="$repo$subpath ⌥ $(__g_get_branch_name)$alias "
  else
    brn="$(pwd) "
  fi
  echo $brn
}

if [[ $G_CUSTOMIZED_PROMPT == true ]]; then
  __update_prompt() {
    PS1="$(whoami) ➤ $(__g_get_name) ❖ "
  }
  PROMPT_COMMAND=__update_prompt
  precmd() { eval "$PROMPT_COMMAND"; }
fi
