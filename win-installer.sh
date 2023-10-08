#!/bin/bash

BASE_URL="https://raw.githubusercontent.com/orangekame3/winget-pkgs/main/manifests/g/orangekame3/ghfetch"
FILES=("orangekame3.ghfetch.yaml" "orangekame3.ghfetch.installer.yaml" "orangekame3.ghfetch.locale.en-US.yaml")

mkdir -p ./tmp

for file in "${FILES[@]}"; do
    curl -L "$BASE_URL/$file" -o "./tmp/$file"
done

winget install -m ./tmp/


for file in "${FILES[@]}"; do
    rm "./tmp/$file"
done
