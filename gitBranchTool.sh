__DEFAULT_BRANCH_PATH=~/.branchAliases/.branchAliases

BRANCH_DELIMITER=${BRANCH_DELIMITER:-'|'}
BRANCH_PATH=${BRANCH_PATH:-"$__DEFAULT_BRANCH_PATH"}
CUSTOMIZED_GIT_PROMPT=${CUSTOMIZED_GIT_PROMPT:-true}
DEFAULT_BRANCH=main

# Git Branch Snippets
mkdir -p ~/.branchAliases/

__getBranchName() {
  echo $(git branch 2> /dev/null | grep \* | cut -d "*" -f2 | cut -d " " -f2)
}

currentBranchPath(){
  if [[ -z $(__getBranchName) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  currentPath="$BRANCH_PATH.$(basename $(git rev-parse --show-toplevel))"
  if [[ ! -e $currentPath ]]; then
    touch $currentPath
  fi
  echo $currentPath
}

list(){
  if [[ -z $(__getBranchName) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  cat -n $(currentBranchPath) | tr '|' '\t'
}

currentBranch() {
  if [[ -z $(__getBranchName) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  cat $(currentBranchPath) | tr '|' '\t' | grep $(__getBranchName)
  if [ $?  != 0 ]; then
      echo "master (or unregistered branch)"
  fi
}

branchNameFromAlias(){
  if [[ -z $(__getBranchName) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 1 )); then
    echo "Wrong Usage"
    echo "\tbranchNameFromAlias <alias>"
    return 1
  fi
  cat "$(currentBranchPath)" | while read line || [ -n "$line" ]; do
    id=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 1)
    als=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 2)
    if [[ $1 == $als ]]; then
        echo $id
        return 0
    fi
  done
  echo "-- Alias not found --"
  return 1
}

editCurrentPath() {
  vim $(currentBranchPath)
}

addBranchAlias(){
  if [[ -z $(__getBranchName) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 2 )); then
    echo "Wrong Usage"
    echo "\taddBranchAlias <id> <alias> [<note>]"
    return 1
  fi
  # checking for alias to be unique
  cat "$(currentBranchPath)" | while read line || [ -n "$line" ]; do
    als=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 2)
    if [[ $2 == $als ]]; then
        echo '-- Alias should be unique --'
        echo '-- FAILED --'
        return 1
    fi
  done
  # Adding branch, alias and note to list only if operation was successful
  echo "$1$BRANCH_DELIMITER$2$BRANCH_DELIMITER$3" >> "$(currentBranchPath)"
  echo "-- Added alias '$2' for branch '$1' --"
}

updateBranchNote() {
  if [[ -z $(__getBranchName) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 2 )); then
    echo "Wrong Usage"
    echo "\tupdateBranchNote <id|alias> <note>"
    return 1
  fi
  found=false
  first=true
  branchPath=$(currentBranchPath)
  cat "$branchPath" | while read line || [ -n "$line" ]; do
    if [[ $first == true ]]; then
        echo -n "" > "$branchPath"
        first=false
    fi
    id=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 1)
    als=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 2)
    note=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 3)
    # update note
    if [[ $1 == $id ||  $1 == $als ]]; then
      echo "$id$BRANCH_DELIMITER$als$BRANCH_DELIMITER$2" >> "$branchPath"
      echo "-- update note for '$id  $als' to '$2' --"
      found=true
    else 
      echo "$id$BRANCH_DELIMITER$als$BRANCH_DELIMITER$note" >> "$branchPath"
    fi
  done
  if [[ $found == false ]]; then
      echo "-- branch not found: $1 --"
  fi
}

branch(){
  if [[ -z $(__getBranchName) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 2 )); then
    echo "Wrong Usage"
    echo "\tbranch <id> <alias> [<note>]"
    return 1
  fi
  # checking for alias to be unique
  cat "$(currentBranchPath)" | while read line || [ -n "$line" ]; do
    als=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 2)
    if [[ $2 == $als ]]; then
        echo '-- Alias should be unique --'
        echo '-- FAILED --'
        return 1
    fi
  done
  # creating and checking out to branch
  git checkout -b $1
  if [ $? -eq 0 ]; then
    # Adding branch, alias and note to list only if operation was successful
    echo "$1$BRANCH_DELIMITER$2$BRANCH_DELIMITER$3" >> "$(currentBranchPath)"
    return 0
  else
    echo '-- FAILED --'
    return 1
  fi
}

check(){
  if [[ -z $(__getBranchName) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 1 )); then
    echo "Wrong Usage"
    echo "\tcheck <id|alias>"
    return 1
  fi
  if [[ $1 == $DEFAULT_BRANCH ]]; then
    git checkout $DEFAULT_BRANCH
    return 0
  fi
  cat "$(currentBranchPath)" | while read line || [ -n "$line" ]; do
    id=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 1)
    als=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 2)
    if [[ $1 == $als || $1 == $id ]]; then
        git checkout $id
        return 0
    fi
  done
  git checkout $1
  if [[ $?  == 0 ]]; then
      echo "-- branch \"$1\" is not registered with alias --"
      return 0
  fi
  echo "-- branch not found: $1 --"
  return 1
}

del(){
  if [[ -z $(__getBranchName) ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 1 )); then
    echo "Wrong Usage"
    echo "\tdel <alias|id> [...<alias|id>]"
    return 1
  fi
  branchPath=$(currentBranchPath)
  for value in "$@"; do
    found=false
    first=true
    cat "$branchPath" | while read line || [ -n "$line" ]; do
      if [[ $first == true ]]; then
          echo -n "" > "$branchPath"
          first=false
      fi
      id=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 1)
      als=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 2)
      note=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 3)
      # Delete base on id/alias
      if [[ $value == $id || $value == $als ]]; then
        echo "-- deleting branch: $id $als --"
        git branch -D  $id
        if [ $? -ne 0 ]; then
          echo '-- FAILED --'
          echo "$id$BRANCH_DELIMITER$als$BRANCH_DELIMITER$note" >> "$branchPath"
        fi
        found=true
      else 
        echo "$id$BRANCH_DELIMITER$als$BRANCH_DELIMITER$note" >> "$branchPath"
      fi
    done
    if [[ $found == false ]]; then
        echo "-- branch not found: $value --"
    fi
  done
}

branchHelp(){
  echo "*  \$DEFAULT_BRANCH \t\t\t : Default branch name, usually master or main"
  echo "*  \$BRANCH_DELIMITER \t\t\t : Delimiter for branch info"
  echo "*  \$BRANCH_PATH \t\t\t : The path where branch info are stored"
  echo "*  \$CUSTOMIZED_GIT_PROMPT \t\t : To whether customize the prompt or not"
  echo "*  currentBranchPath \t\t\t : returns the branch path for the current repository"
  echo "*  editCurrentPath \t\t\t : opens current branch path in vim for manual editing"
  echo "*  list \t\t\t\t : Lists all branches with their id, alias, and notes"
  echo "*  addBranchAlias <id> <alias> [<note>]  : Adds alias and note to a branch that is not stored yet"
  echo "*  updateBranchNote <id|alias> <note> \t : adds/updates the notes for a branch base on id/alias"
  echo "*  branch <id> <alias> [<note>] \t : creates a branch with id, alias, and note, and checks in"
  echo "*  check <id|alias> \t\t\t : checks into a branch base on an id or an alias"
  echo "*  del [...<id|alias>] \t\t\t : delete listed branches base on id or alias"
  echo "*  currentBranch \t\t\t : Returns the name of current branch with alias and note"
  echo "*  branchNameFromAlias <alias> \t\t : Returns the branch name from an alias"
  echo "*  branchHelp \t\t\t\t : shows this list"
}

# Custom prompt
__getName() {
  brn=""
  if [[ -n $(__getBranchName) ]]; then
    name=""
    cat "$(currentBranchPath)" | while read line || [ -n "$line" ]; do
      id=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 1)
      als=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 2)
      if [[ $(__getBranchName) == $id ]]; then
          name=" ($als)"
          break
      fi
    done
    topPath="$(git rev-parse --show-toplevel)"
    repo=$(basename $topPath)
    subpath=${${PWD:$(($(echo "$(PWD)" | awk -v repo="$repo" '{print index($0, repo)}') + ${#repo}))}#/}
    if [ -n "$subpath" ]; then
      subpath=" [$subpath]"
    fi
    brn="$repo$subpath ⌥ $(__getBranchName)$name "
  else
    brn="$(pwd) "
  fi
  echo $brn
}

if [[ $CUSTOMIZED_GIT_PROMPT == true ]]; then
  __update_prompt() {
    PS1="%n ➤ $(__getName)❖ "
  }
  PROMPT_COMMAND=__update_prompt
  precmd() { eval "$PROMPT_COMMAND"; }
fi
