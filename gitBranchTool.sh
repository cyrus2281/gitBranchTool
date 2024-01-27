# Custom prompt
autoload -Uz vcs_info
precmd() { vcs_info }
zstyle ':vcs_info:git:*' formats '%b'
setopt PROMPT_SUBST
getName() {
  brn=""
  if [[ -n ${vcs_info_msg_0_} ]]; then
    name=""
    cat "$(currentBranchPath)" | while read line || [ -n "$line" ]; do
      id=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 1)
      als=$(echo "$line" | cut -d "$BRANCH_DELIMITER" -f 2)
      if [[ ${vcs_info_msg_0_} == $id ]]; then
          name=" ($als)"
          break
      fi
    done
    brn="%1~ ⌥ ${vcs_info_msg_0_}$name "
  else
    brn="$(pwd) "
  fi
  echo $brn
}
PROMPT='%n ➤ $(getName)❖ '

# Git Branch Snippets

BRANCH_DELIMITER="|"
BRANCH_PATH=~/.branchAliases/.branchAliases

currentBranchPath(){
  if [[ -z ${vcs_info_msg_0_} ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  currentPath="$BRANCH_PATH.${PWD##*/}"
  if [[ ! -e $currentPath ]]; then
    touch $currentPath
  fi
  echo $currentPath
}

list(){
  cat -n $(currentBranchPath) | tr '|' '\t'
}

editCurrentPath() {
  vim $(currentBranchPath)
}

branchAlias(){
  if [[ -z ${vcs_info_msg_0_} ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 2 )); then
    echo "Wrong Usage"
    echo "\tcheck <id> <alias> [<note>]"
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
  if [[ -z ${vcs_info_msg_0_} ]]; then
    echo "-- Not a git repository --"
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
  if [[ -z ${vcs_info_msg_0_} ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if (( $# < 2 )); then
    echo "Wrong Usage"
    echo "\tcheck <id> <alias> [<note>]"
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
  if [[ -z ${vcs_info_msg_0_} ]]; then
    echo "-- Not a git repository --"
    return 1
  fi
  if [[ $1 == master ]]; then
    git checkout master
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
  echo '-- branch not found --'
  return 1
}

del(){
  if [[ -z ${vcs_info_msg_0_} ]]; then
    echo "-- Not a git repository --"
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
  echo "*  \$BRANCH_DELIMITER \t\t\t : Delimiter for branch info"
  echo "*  \$BRANCH_PATH \t\t\t : The path where branch info are stored"
  echo "*  currentBranchPath \t\t\t : returns the branch path for the current repository"
  echo "*  editCurrentPath \t\t\t : opens current branch path in vim for manual editing"
  echo "*  list \t\t\t\t : Lists all branches with their id, alias, and notes"
  echo "*  branchAlias <id> <alias> [<note>] \t : Adds alias and note to a branch that is not stored yet"
  echo "*  updateBranchNote <id|alias> <note> \t : adds/updates the notes for a branch base on id/alias"
  echo "*  branch <id> <alias> [<note>] \t : creates a branch with id, alias, and note, and checks in"
  echo "*  check <id|alias> \t\t\t : checks into a branch base on an id or an alias"
  echo "*  del [...<id|alias>] \t\t\t : delete listed branches base on id or alias"
  echo "*  branchHelp \t\t\t\t : shows this list"
}
