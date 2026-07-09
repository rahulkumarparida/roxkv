#!/bin/sh
set -eu

HOME_DIR="${HOME:-/root}"
ROXKV_HOME="${HOME_DIR}/.roxkv"

mkdir -p \
  "${ROXKV_HOME}/roxdb" \
  "${ROXKV_HOME}/roxlogs" \
  "${ROXKV_HOME}/roxsnaps" \
  "${ROXKV_HOME}/roxmetrics"

exec "$@"
