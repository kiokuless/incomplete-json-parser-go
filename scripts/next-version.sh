#!/bin/bash

set -euo pipefail

# Get the current version tag
current_version=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Remove the 'v' prefix if it exists
version_number=${current_version#v}

# Split version into major, minor, patch
IFS='.' read -r major minor patch <<< "$version_number"

# Default values if version parsing fails
major=${major:-0}
minor=${minor:-0}
patch=${patch:-0}

# Calculate next version based on argument
case "${1:-patch}" in
    "major")
        new_major=$((major + 1))
        new_minor=0
        new_patch=0
        ;;
    "minor")
        new_major=$major
        new_minor=$((minor + 1))
        new_patch=0
        ;;
    "patch")
        new_major=$major
        new_minor=$minor
        new_patch=$((patch + 1))
        ;;
    *)
        echo "Usage: $0 [major|minor|patch]" >&2
        exit 1
        ;;
esac

echo "v${new_major}.${new_minor}.${new_patch}"
