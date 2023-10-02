#!/usr/bin/env bash
set -euo pipefail

echo() {
  command echo "$1" "$2"
}

# Call the requested function and pass the arguments as-is
"$@"
