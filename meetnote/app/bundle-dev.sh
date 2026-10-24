#!/bin/zsh
# Builds and opens a debug example.invalid. Needed to test system audio capture (online meetings):
# macOS only grants "System Audio Recording" to an app bundle, so recording from dev.sh, which runs
# under the terminal, captures silence. The debug app finds the repo's .env and spike/models itself.
set -e
cd "${0:A:h}"
export PATH=/opt/homebrew/opt/rustup/bin:$PATH
npm run tauri build -- --debug --bundles app
open src-tauri/target/debug/bundle/macos/example.invalid
