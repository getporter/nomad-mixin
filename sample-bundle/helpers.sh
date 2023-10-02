#!/usr/bin/env bash
set -euo pipefail

echo() {
  command echo Eval id is $1
}

# Call the requested function and pass the arguments as-is
"$@"
