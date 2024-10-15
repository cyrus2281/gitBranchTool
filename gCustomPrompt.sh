#!/bin/bash

# This script is used to customize the prompt for the gitBranchTool.

# Author: Cyrus Mobini - https://github.com/cyrus2281
# Github Repository: https://github.com/cyrus2281/gitBranchTool
# License: MIT License

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
  branch=$(git branch 2> /dev/null | grep \* | cut -d "*" -f2 | cut -d " " -f2)
  if [[ -n $branch ]]; then
    alias=$(g get-branch-alias $branch -N)
    if [[ -n $alias ]]; then
      branch="$branch ($alias)"
    fi
    topPath="$(git rev-parse --show-toplevel)"
    repo=$(basename $topPath)
    subpath=$(__g_get_subdirectory $repo)
    brn="$repo$subpath ⌥ $branch "
  else
    brn="$(pwd) "
  fi
  echo $brn
}

__update_prompt() {
  PS1="$(whoami) ➤ $(__g_get_name) ❖ "
}
PROMPT_COMMAND=__update_prompt
precmd() { eval "$PROMPT_COMMAND"; }
