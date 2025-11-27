#!/usr/bin/env bash
set -euo pipefail

function printUsage() {
  echo "usage: $(basename "$0") [options] tag"
  echo
  echo "Pushes ObjectBox Generator binaries to a GitHub release (draft)"
  echo "  tag: the release version (e.g. 1.2.3); no 'v' prefix is required"
  echo
  echo "Options:"
  echo "  --force-release-upload: upload to an existing release (non-draft) if it exists"
  echo "  --linux-x64: upload only the Linux artifact (useful for testing)"
  echo "  --dry-run: do not upload any artifacts"

  exit 1
}

# Default values
force_release_upload=false
linux_x64_only=false
dry_run=false

# Parse command line arguments; e.g. override defaults
while test $# -gt 0; do
case "$1" in
-h|help|--help|usage)
    printUsage
    ;;
--force-release-upload)
    force_release_upload=true
    shift
    ;;
--linux-x64) # partial artifact upload: this is useful to test the script locally
    linux_x64_only=true
    shift
    ;;
--dryrun|--dry-run)
    dry_run=true
    shift
    ;;
"") # ignore any empty parameters ("artifacts" that can happen depending on how the caller constructs the parameters)
    shift
    ;;
*)
    break
    ;;
esac
done

if [[ "$#" -ne "1" ]]; then
  echo "$# arguments provided ('$*'), expected 1 "
  printUsage
fi

echo "Options: linux_x64_only=$linux_x64_only, dry_run=$dry_run"

# start in the repo root directory
cd "$(dirname "$0")/.."

github_tag=$1
github_repo="https://github.com/objectbox/objectbox-generator"

if ! [[ $github_tag == v* ]]; then
  github_tag="v$github_tag"
fi

echo "Checking GitHub authentication status..."
gh auth status

echo "Checking if release ${github_tag} exists..."
# https://cli.github.com/manual/gh_release_view
if gh release view --repo "${github_repo}" "${github_tag}"; then
    echo "" # gh may have not ended with a new line
    is_draft_json=$(gh release view --repo "${github_repo}" "${github_tag}" --json isDraft)
    echo "Existing release \"${github_tag}\" found; draft status: $is_draft_json"
    if [[ $is_draft_json == *"\"isDraft\":true"* ]]; then
      echo "Existing release draft identified: will try to upload artifacts."
      echo "Note: existing artifact will not be overwritten; delete them manually beforehand if required."
    elif [[ $is_draft_json == *"\"isDraft\":false"* ]]; then
      echo "WARN: Existing release is not a draft"
        if [ "$force_release_upload" = true ]; then
          echo "WARN: Upload to existing release is forced; continuing..."
        else
          echo "WARN: Aborting upload to avoid messing up existing releases."
          echo "WARN: Situation must be resolved manually, or force the upload via argument."
          exit 1
        fi
    else
      echo "ERROR: Must this script be updated? Found unknown draft status: $is_draft_json"
      exit 1
    fi
else
  echo "No release found for '${github_tag}'; creating a new release draft..."
  commit_id=$(git rev-parse HEAD)
  echo "Using commit ID (this commit must also exist in GitHub): $commit_id"
  gh release create --draft --repo "${github_repo}" "${github_tag}" \
    --target "${commit_id}" \
    --title "ObjectBox Generator ${github_tag}" \
    --notes "See CHANGELOG.md for details."
  echo "Uploading artifacts to new release draft..."
fi

if [ "$dry_run" = true ]; then
  echo "Dry run: skipping artifact preparation & upload"
  exit 0
fi

# Explicitly specify all files instead of using a wildcard glob to avoid unintended uploads (and catch missing files).
# These artifacts are built by the build jobs in GitLab CI. See .gitlab-ci.yml for details on how to build.
candidate_files=(
                 objectbox-generator-Linux.zip
                 objectbox-generator-macOS.zip
                 objectbox-generator-Windows.zip
                )
files_upload=()
for f in "${candidate_files[@]}"; do
  if [ -f "${f}" ]; then
    files_upload+=("${f}")
  elif [ -f "artifacts/${f}" ]; then
    files_upload+=("artifacts/${f}")
  else
    echo "ERROR: File not found: ${f} (also not in artifacts/${f})"
    exit 1
  fi
  if [ "$linux_x64_only" = true ]; then break; fi # Linux is the first array item
done

echo "${#files_upload[@]} files to upload:"
for file_upload in "${files_upload[@]}"; do
  ls -lh "$file_upload"
done

# https://cli.github.com/manual/gh_release_upload
gh release upload --repo "${github_repo}" "${github_tag}" "${files_upload[@]}"

echo "Done."