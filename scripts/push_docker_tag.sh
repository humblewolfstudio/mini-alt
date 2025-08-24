#!/bin/bash

# Function to display usage
usage() {
    echo "Usage: $0 <tag>"
    echo "Example: $0 v1.0.0"
    exit 1
}

# Check for the correct number of arguments
if [ "$#" -ne 1 ]; then
    usage
fi

# Input arguments
TAG="$1"

# Validate tag format (e.g., v1.0.0)
if [[ ! "$TAG" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Tag format is invalid. Use the format v<major>.<minor>.<patch> (e.g., v1.0.0)."
    exit 1
fi

# Validate Git status
if ! git rev-parse --is-inside-work-tree &>/dev/null; then
    echo "Error: Not inside a Git repository."
    exit 1
fi

# Check if the tag already exists
if git rev-parse "$TAG" >/dev/null 2>&1; then
    echo "Error: Tag '$TAG' already exists in the repository."
    exit 1
fi

# Moving to parent folder
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.." || { echo "Error: Failed to move to the parent directory."; exit 1; }


git checkout main

# Pull content
git pull origin main

# Create the Git tag
echo "Creating Git tag '$TAG'..."
git tag "$TAG"
if [ $? -ne 0 ]; then
    echo "Error: Failed to create Git tag."
    exit 1
fi

# Push the Git tag to the remote repository
echo "Pushing Git tag '$TAG' to remote repository..."
git push origin "$TAG"
if [ $? -ne 0 ]; then
    echo "Error: Failed to push Git tag."
    exit 1
fi

DOCKER_TAG="${TAG#v}"
IMAGE_NAME="verdeloro/mini-alt:$DOCKER_TAG"

echo "Building Docker image '$IMAGE_NAME'..."
docker build -t "$IMAGE_NAME" .
if [ $? -ne 0 ]; then
    echo "Error: Failed to build Docker image."
    exit 1
fi

# Build the Docker image
LATEST_IMAGE_NAME="verdeloro/mini-alt:latest"
echo "Tagging Docker image '$IMAGE_NAME' as 'latest'..."
docker tag "$IMAGE_NAME" "$LATEST_IMAGE_NAME"

# Push the Docker image to Docker Hub
echo "Pushing Docker image '$IMAGE_NAME' to Docker Hub..."
docker login
docker push "$IMAGE_NAME"
docker push "$LATEST_IMAGE_NAME"
if [ $? -ne 0 ]; then
    echo "Error: Failed to push Docker image to Docker Hub."
    exit 1
fi

echo "Successfully created and pushed Docker image '$IMAGE_NAME' with tag '$TAG'."