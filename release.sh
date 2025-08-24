#!/bin/bash
: <<'COMMENT'
# Relase script for micro-agent-go
Usage:
1. Update the version in release.env
2. Run this script: ./release.sh
3. Create a GitHub release with the new tag and description
COMMENT

set -o allexport; source release.env; set +o allexport

echo "Generating release: ${TAG} ${ABOUT}"

find . -name '.DS_Store' -type f -delete

echo "📝 Replacing ${PREVIOUS_TAG} by ${TAG} in files..."

# Update all go.mod files in examples subdirectories
for dir in examples/*/; do
  if [ -f "${dir}go.mod" ]; then
    echo "Updating ${dir}go.mod"
    go run release.go -old="github.com/micro-agent/micro-agent-go ${PREVIOUS_TAG}" -new="github.com/micro-agent/micro-agent-go ${TAG}" -file="${dir}go.mod"
  fi
done

# Update all go.mod files in cmd subdirectories
for dir in cmd/*/; do
  if [ -f "${dir}go.mod" ]; then
    echo "Updating ${dir}go.mod"
    go run release.go -old="github.com/micro-agent/micro-agent-go ${PREVIOUS_TAG}" -new="github.com/micro-agent/micro-agent-go ${TAG}" -file="${dir}go.mod"
  fi
done

#go run release.go -old="${PREVIOUS_TAG}" -new="${TAG}" -file="./README.md"

git add .
git commit -m "📦 ${ABOUT}"
git push origin main

git tag -a ${TAG} -m "${ABOUT}"
git push origin ${TAG}

