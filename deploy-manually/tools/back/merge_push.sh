#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
sh "$SCRIPT_DIR/merge_tpower.sh" && \
sh "$SCRIPT_DIR/merge_vworld.sh" && \
sh "$SCRIPT_DIR/merge_wintoto.sh" && \
sh "$SCRIPT_DIR/push_branch.sh"