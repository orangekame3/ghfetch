#!/bin/bash

BASE_URL="https://raw.githubusercontent.com/orangekame3/winget-pkgs/main/manifests/g/orangekame3/cobra-template"
FILES=("orangekame3.cobra-template.yaml" "orangekame3.cobra-template.installer.yaml" "orangekame3.cobra-template.locale.en-US.yaml")

mkdir -p ./tmp

for file in "${FILES[@]}"; do
    curl -L "$BASE_URL/$file" -o "./tmp/$file"
done

winget install -m ./tmp/


for file in "${FILES[@]}"; do
    rm "./tmp/$file"
done
