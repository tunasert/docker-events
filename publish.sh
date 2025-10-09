#!/usr/bin/env bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_OWNER="filippofinke"
REPO_NAME="docker-events"
DOCKER_IMAGE="filippofinke/docker-events"
VERSION_FILE="VERSION"

# Function to print colored messages
print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Function to show usage
usage() {
    cat << EOF
Usage: $0 [major|minor|patch]

Publish script for docker-events - Increments version, creates GitHub release, and publishes to Docker Hub.

Arguments:
  major    Increment major version (X.0.0)
  minor    Increment minor version (x.X.0)
  patch    Increment patch version (x.x.X)

Requirements:
  - git
  - GitHub CLI (gh)
  - docker
  - Clean git working directory
  - Logged in to Docker Hub (docker login)
  - Authenticated with GitHub CLI (gh auth login)

Examples:
  $0 patch    # 1.0.0 -> 1.0.1
  $0 minor    # 1.0.0 -> 1.1.0
  $0 major    # 1.0.0 -> 2.0.0
EOF
    exit 1
}

# Function to check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."
    
    # Check if git is installed
    if ! command -v git &> /dev/null; then
        print_error "git is not installed"
        exit 1
    fi
    
    # Check if GitHub CLI is installed
    if ! command -v gh &> /dev/null; then
        print_error "GitHub CLI (gh) is not installed"
        print_info "Install it with: brew install gh"
        exit 1
    fi
    
    # Check if docker is installed
    if ! command -v docker &> /dev/null; then
        print_error "docker is not installed"
        exit 1
    fi
    
    # Check if GitHub CLI is authenticated
    if ! gh auth status &> /dev/null; then
        print_error "GitHub CLI is not authenticated"
        print_info "Run: gh auth login"
        exit 1
    fi
    
    # Check if working directory is clean
    if [[ -n $(git status -s) ]]; then
        print_error "Working directory is not clean. Please commit or stash changes."
        git status -s
        exit 1
    fi
    
    # Check if we're on the main branch
    current_branch=$(git rev-parse --abbrev-ref HEAD)
    if [[ "$current_branch" != "main" ]]; then
        print_warning "You are not on the main branch (current: $current_branch)"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    print_success "All prerequisites met"
}

# Function to get current version
get_current_version() {
    if [[ -f "$VERSION_FILE" ]]; then
        cat "$VERSION_FILE"
    else
        # Try to get the latest git tag
        latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
        if [[ -n "$latest_tag" ]]; then
            echo "${latest_tag#v}"
        else
            echo "0.0.0"
        fi
    fi
}

# Function to increment version
increment_version() {
    local version=$1
    local bump_type=$2
    
    # Parse version
    IFS='.' read -r major minor patch <<< "$version"
    
    case $bump_type in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            print_error "Invalid bump type: $bump_type"
            exit 1
            ;;
    esac
    
    echo "${major}.${minor}.${patch}"
}

# Function to create or update VERSION file
update_version_file() {
    local version=$1
    echo "$version" > "$VERSION_FILE"
    git add "$VERSION_FILE"
}

# Function to create git tag
create_git_tag() {
    local version=$1
    local tag="v${version}"
    
    print_info "Creating git tag: $tag"
    git tag -a "$tag" -m "Release $tag"
    print_success "Tag created: $tag"
}

# Function to push to GitHub
push_to_github() {
    local version=$1
    local tag="v${version}"
    
    print_info "Pushing to GitHub..."
    git push origin main
    git push origin "$tag"
    print_success "Pushed to GitHub"
}

# Function to create GitHub release
create_github_release() {
    local version=$1
    local tag="v${version}"
    
    print_info "Creating GitHub release..."
    
    # Generate release notes
    gh release create "$tag" \
        --title "Release $tag" \
        --generate-notes \
        --repo "${REPO_OWNER}/${REPO_NAME}"
    
    print_success "GitHub release created: $tag"
}

# Function to build and push Docker image
build_and_push_docker() {
    local version=$1
    local tag_version="${DOCKER_IMAGE}:${version}"
    local tag_latest="${DOCKER_IMAGE}:latest"
    
    print_info "Building Docker image..."
    docker build -t "$tag_version" -t "$tag_latest" .
    print_success "Docker image built"
    
    print_info "Pushing Docker image with version tag: $version"
    docker push "$tag_version"
    print_success "Docker image pushed: $tag_version"
    
    print_info "Pushing Docker image with latest tag"
    docker push "$tag_latest"
    print_success "Docker image pushed: $tag_latest"
}

# Main script
main() {
    # Check if argument is provided
    if [[ $# -ne 1 ]]; then
        usage
    fi
    
    local bump_type=$1
    
    # Validate bump type
    if [[ ! "$bump_type" =~ ^(major|minor|patch)$ ]]; then
        print_error "Invalid version bump type: $bump_type"
        usage
    fi
    
    # Check prerequisites
    check_prerequisites
    
    # Get current version
    current_version=$(get_current_version)
    print_info "Current version: $current_version"
    
    # Calculate new version
    new_version=$(increment_version "$current_version" "$bump_type")
    print_info "New version: $new_version"
    
    # Confirm with user
    echo
    print_warning "This will:"
    echo "  1. Update VERSION file to $new_version"
    echo "  2. Create git tag v${new_version}"
    echo "  3. Push to GitHub"
    echo "  4. Create GitHub release v${new_version}"
    echo "  5. Build and push Docker image ${DOCKER_IMAGE}:${new_version}"
    echo "  6. Tag and push Docker image ${DOCKER_IMAGE}:latest"
    echo
    read -p "Continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_warning "Aborted"
        exit 0
    fi
    
    # Update VERSION file
    print_info "Updating VERSION file..."
    update_version_file "$new_version"
    
    # Commit VERSION file
    print_info "Committing VERSION file..."
    git commit -m "Bump version to $new_version"
    print_success "VERSION file committed"
    
    # Create git tag
    create_git_tag "$new_version"
    
    # Push to GitHub
    push_to_github "$new_version"
    
    # Create GitHub release
    create_github_release "$new_version"
    
    # Build and push Docker image
    build_and_push_docker "$new_version"
    
    echo
    print_success "🎉 Release $new_version published successfully!"
    echo
    print_info "GitHub Release: https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/tag/v${new_version}"
    print_info "Docker Image: https://hub.docker.com/r/${DOCKER_IMAGE}"
}

# Run main function
main "$@"
