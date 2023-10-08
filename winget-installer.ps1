$BASE_URL = "https://raw.githubusercontent.com/orangekame3/winget-pkgs/main/manifests/g/orangekame3/cobra-template"
$FILES = @("orangekame3.cobra-template.yaml", "orangekame3.cobra-template.installer.yaml", "orangekame3.cobra-template.locale.en-US.yaml")

if (-not (Test-Path ./tmp)) {
    New-Item -ItemType Directory -Path ./tmp
}

foreach ($file in $FILES) {
    Invoke-WebRequest -Uri "$BASE_URL/$file" -OutFile "./tmp/$file"
}

winget install -m ./tmp/orangekame3.cobra-template.yaml


foreach ($file in $FILES) {
    Remove-Item "./tmp/$file"
}
