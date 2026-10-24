#!/bin/zsh
# Builds the spike CLI with Command Line Tools only (no Xcode needed).
set -e
cd "${0:A:h}"
swiftc -O -swift-version 5 -parse-as-library -o transcribe Sources/main.swift
echo "built ./transcribe"
