#!/bin/bash

# This script generates a code coverage badge and adds it to the README.md file.

# Find the coverage percentage
go test ./... -coverprofile coverage.out
COVERAGE=$(go tool cover -func=coverage.out | grep total: | grep -Eo '[0-9]+\.[0-9]+')
COLOR=orange
if (( $(echo "$COVERAGE <= 50" | bc -l) )) ; then
    COLOR="FF0000"
elif (( $(echo "$COVERAGE > 80" | bc -l) )); then
    COLOR="31C754"
fi
rm coverage.out

# Generate the SVG image URL
svg_image_url="https://img.shields.io/badge/coverage-$COVERAGE%25-$COLOR"

readme_file="README.md"

# Find the line number where the CI badge is located
badge_line_number=$(grep -n "^\[!\[ci\]" "$readme_file" | cut -d: -f1)
if [[ -z "$badge_line_number" ]]; then
  echo "CI badge not found in $readme_file"
  exit 1
fi

# Insert the SVG image URL below the CI badge
cp "$readme_file" "$readme_file.tmp"
awk -v line="$badge_line_number" -v url="$svg_image_url" 'NR==line{print; print "[![Code Coverage](" url ")](" url ")"; next}1' "$readme_file.tmp" > "$readme_file"
rm "$readme_file.tmp"

echo "SVG image added below the CI badge in $readme_file"
