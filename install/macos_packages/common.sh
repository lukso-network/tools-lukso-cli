#!/bin/bash

error() {
  echo "[ERROR]: $1"
  exit 1
}

lukso_version() {
  lukso_program_file="$1"
  # Check if this is a relative path: if so, prepend './' to it
  if [ "$1" = "${1#/}" ]; then
    lukso_program_file="./$lukso_program_file"
  fi
  >&2 echo "prog: $lukso_program_file"
  >&2 ls -l $lukso_program_file
  if "$lukso_program_file" version 2>&1 >/dev/null; then
    "$lukso_program_file" version | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+'
  else
    # try to get this information from Cargo.toml
    pushd "$(git rev-parse --show-toplevel)" || true
    grep '^version = \"\(.*\)\"' Cargo.toml | cut -f 2 -d '"'
    popd
  fi
}
