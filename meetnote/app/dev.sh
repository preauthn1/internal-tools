#!/bin/zsh
# Starts meetnote in development mode:
# - models come from spike/models (no download needed)
# - the Claude API key comes from ../.env, so dev builds never touch the keychain
#   (macOS asks for keychain access again after every rebuild, because the signature changes)
set -e
cd "${0:A:h}"
export PATH=/opt/homebrew/opt/rustup/bin:$PATH
export MEETNOTE_MODEL_DIR="${0:A:h}/../spike/models"
if [[ -f ../.env ]]; then
  key=$(grep '^ANTHROPIC_API_KEY=' ../.env | cut -d= -f2- | tr -d "\"'")
  [[ -n $key ]] && export ANTHROPIC_API_KEY=[REDACTED]
fi
exec npm run tauri dev
