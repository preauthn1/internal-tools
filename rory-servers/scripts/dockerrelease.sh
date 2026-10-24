#!/bin/bash
release_folder=$(ls -d ./release/*/|head -n 1)
build_folder="$release_folder/linux"
cp -r "$build_folder" approot