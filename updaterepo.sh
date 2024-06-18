#!/bin/bash

# Define variables
GITHUB_USERNAME="gsantopaolo" #"YOUR-USERNAME"
GITHUB_GENMIND="gen-mind"
OLD_REPO_NAME="cognix-services"
NEW_REPO_NAME="cognix"


# to generate te access token go to https://github.com/settings/tokens
# and paste the token here (token classic)
# Repository access all
GITHUB_TOKEN="GEENRATE YOUR TOKEN"

# Navigate to your local repository
#cd $LOCAL_REPO_PATH || { echo "Local repository path not found"; exit 1; }

# Update the remote URL to use the token for authentication
git remote set-url origin https://$GITHUB_USERNAME:$GITHUB_TOKEN@github.com/$GITHUB_GENMIND/$NEW_REPO_NAME.git

# Verify the new remote URL
git remote -v

# Fetch the latest changes from the new remote repository
git fetch origin

# Push any local changes to the remote repository (adjust the branch name as necessary)
git push origin main  # Change 'main' to your branch name if it's different

echo "Repository URL updated and changes pushed successfully."


