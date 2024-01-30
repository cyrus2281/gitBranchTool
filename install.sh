#!/bin/bash

# This script will install gitBranchTool on your system.

# Author: Cyrus Mobini - https://github.com/cyrus2281
# Github Repository: https://github.com/cyrus2281/gitBranchTool
# License: MIT License


# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Function to check if a line exists in a file
line_exists_in_file() {
  touch $2
  grep -qF "$1" "$2"
}

# Adding gitBranchTool to profiles
# $1 gitBranchTool directory path
# $2 G_CUSTOMIZED_PROMPT value
# $3 profile path
add_gitBranchTool_to_profile() {
  file_path=$(eval echo "$3")
  # Check if the line already exists in the file
  if line_exists_in_file "source $1/gitBranchTool.sh" "$file_path"; then
      echo -e "\tThe source command is already present in $file_path. Skipping."
  else
      # Add export command for G_CUSTOMIZED_PROMPT
      echo -e "\n# Setting G command - gitBranchTool" >> "$file_path"
      echo "export G_CUSTOMIZED_PROMPT=$2" >> "$file_path"
      echo "export G_DIRECTORY=$1" >> "$file_path"
      echo -e "source $1/gitBranchTool.sh\n" >> "$file_path"
      echo -e "\tAdded gitBranchTool to $file_path."
  fi
}

gitBranchToolURL=https://raw.githubusercontent.com/cyrus2281/gitBranchTool/main/gitBranchTool.sh
gitBranchToolDir=$(eval echo "${GIT_BRANCH_TOOL_DIR:-"~/.gitBranchTool"}")
gitBranchToolScriptPath=$gitBranchToolDir/gitBranchTool.sh

# Create gitBranchTool directory if it doesn't exist
mkdir -p $gitBranchToolDir

# Check if curl exists
if command_exists "curl"; then
  echo "Downloading the gitBranchTool script..."
  curl -o $gitBranchToolScriptPath -fsSL $gitBranchToolURL
# Check if wget exists
elif command_exists "wget"; then
  echo "Downloading the gitBranchTool script..."
  wget -O $gitBranchToolScriptPath  $gitBranchToolURL
# Throw an error if neither curl nor wget is found
else
  echo "Error: Neither curl nor wget is installed. Please install either curl or wget and try again."
  exit 1
fi

# Add execute permission for the current user
chmod +x $gitBranchToolScriptPath

# Prompt user for 'G custom prompt' preference
read -p "Do you want the 'G custom prompt'? (yes/no, default: yes): " custom_prompt_preference
custom_prompt_preference=$(echo "$custom_prompt_preference" | tr '[:upper:]' '[:lower:]')  # Convert to lowercase
custom_prompt_preference=${custom_prompt_preference:-yes}

# Set G_CUSTOMIZED_PROMPT variable based on user preference
if [[ "$custom_prompt_preference" =~ ^(yes|y)$ ]]; then
  export G_CUSTOMIZED_PROMPT=true
else
  export G_CUSTOMIZED_PROMPT=false
fi

# Add script to bashrc 
add_gitBranchTool_to_profile $gitBranchToolDir $G_CUSTOMIZED_PROMPT ~/.bashrc

# Check if on macOS and add to zshrc
if [ "$(uname)" == "Darwin" ]; then
  add_gitBranchTool_to_profile $gitBranchToolDir $G_CUSTOMIZED_PROMPT ~/.zshrc
fi

# Ask user for additional profiles
while true; do
  read -p "Enter the path to an additional profile file (or press Enter to exit): " profile_path
  if [ -z "$profile_path" ]; then
    break
  else
    add_gitBranchTool_to_profile $gitBranchToolDir $G_CUSTOMIZED_PROMPT $profile_path
  fi
done


# Setup completed
echo -e "\nSetup completed successfully."
echo -e "\nPlease restart your terminal or enter the following command to start using gitBranchTool."
echo -e "\n\t\tsource $gitBranchToolScriptPath"
echo -e "\nOnce activated, you can run 'g help' to get a list of all g commands.\n"