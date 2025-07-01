#!/bin/bash

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're on main branch
current_branch=$(git branch --show-current)
if [ "$current_branch" != "main" ]; then
    print_error "You must be on the main branch to release. Current branch: $current_branch"
    exit 1
fi

# Check if working directory is clean
if ! git diff-index --quiet HEAD --; then
    print_error "Working directory is not clean. Please commit or stash your changes."
    exit 1
fi

# Get release type
release_type="${1:-patch}"

# Validate release type
case "$release_type" in
    "major"|"minor"|"patch")
        ;;
    *)
        print_error "Invalid release type: $release_type"
        echo "Usage: $0 [major|minor|patch]"
        exit 1
        ;;
esac

# Get current and next version
current_version=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
next_version=$(./scripts/next-version.sh "$release_type")

print_status "Current version: $current_version"
print_status "Next version: $next_version"

# Confirm release
echo
read -p "Do you want to release $next_version? (y/N): " -r
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_warning "Release cancelled."
    exit 0
fi

# Run tests before release
print_status "Running tests..."
if ! go test -v ./...; then
    print_error "Tests failed. Release cancelled."
    exit 1
fi

# Run linter if available
if command -v golangci-lint &> /dev/null; then
    print_status "Running linter..."
    if ! golangci-lint run; then
        print_error "Linter failed. Release cancelled."
        exit 1
    fi
else
    print_warning "golangci-lint not found. Skipping linter check."
fi

# Format and tidy
print_status "Formatting and tidying..."
go fmt ./...
go mod tidy

# Commit any changes from formatting/tidying
if ! git diff-index --quiet HEAD --; then
    print_status "Committing formatting changes..."
    git add .
    git commit -m "chore: format and tidy before release $next_version"
fi

# Create and push tag
print_status "Creating tag $next_version..."
git tag -a "$next_version" -m "Release $next_version"

print_status "Pushing tag to origin..."
git push origin "$next_version"

print_status "✨ Release $next_version completed successfully!"
print_status "GitHub Actions will now create the release automatically."
print_status "Check https://github.com/kiokuless/incomplete-json-parser-go/actions"
